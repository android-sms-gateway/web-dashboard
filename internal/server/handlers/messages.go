package handlers

import (
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/client"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const previewLength = 100

type MessagesHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewMessagesHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &MessagesHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *MessagesHandler) Register(r fiber.Router) {
	g := r.Group("/messages", session.AuthRequired(), client.New(h.gatewaySvc))

	g.Get("", h.list)
	g.Post("", validation.DecorateWithBodyEx(h.Validator, h.send))
	g.Get(":id", h.get)
}

type messageListQuery struct {
	Limit    *int       `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset   *int       `query:"offset"   validate:"omitempty,min=0"`
	State    *string    `query:"state"`
	DeviceID *string    `query:"deviceId"`
	From     *time.Time `query:"from"`
	To       *time.Time `query:"to"`
}

func (q messageListQuery) toOptions() smsgateway.ListMessagesOptions {
	return smsgateway.ListMessagesOptions{
		Limit:          q.Limit,
		Offset:         q.Offset,
		State:          q.State,
		DeviceID:       q.DeviceID,
		From:           q.From,
		To:             q.To,
		IncludeContent: lo.ToPtr(true),
	}
}

type messageListItem struct {
	ID          string                     `json:"id"`
	DeviceID    string                     `json:"deviceId"`
	State       smsgateway.ProcessingState `json:"state"`
	Recipients  []recipientItem            `json:"recipients"`
	TextPreview string                     `json:"textPreview"`
	CreatedAt   time.Time                  `json:"createdAt"`
}

type recipientItem struct {
	PhoneNumber string                     `json:"phoneNumber"`
	State       smsgateway.ProcessingState `json:"state"`
	Error       *string                    `json:"error,omitempty"`
}

type listMessagesResponse struct {
	Items []messageListItem `json:"items"`
	Total int               `json:"total"`
}

// ListMessages returns paginated messages.
//
//	@Summary		List messages
//	@Description	Returns paginated list of messages with optional filters.
//	@Tags			messages
//	@Produce		json
//	@Param			limit		query		int		false	"Max results"	default(50)
//	@Param			offset		query		int		false	"Offset"		default(0)
//	@Param			state		query		string	false	"Filter by state"
//	@Param			deviceId	query		string	false	"Filter by device ID"
//	@Param			from		query		string	false	"Start date (RFC3339)"
//	@Param			to			query		string	false	"End date (RFC3339)"
//	@Success		200			{object}	listMessagesResponse
//	@Failure		401			{object}	fiberfx.ErrorResponse
//	@Router			/messages [get]
func (h *MessagesHandler) list(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	query := new(messageListQuery)
	if err := h.QueryParserValidator(c, query); err != nil {
		h.logger.Warn("failed to parse query", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "failed to parse query")
	}

	messages, total, err := sdkClient.ListMessages(c.Context(), query.toOptions())
	if err != nil {
		h.logger.Warn("failed to list messages", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to list messages")
	}

	items := make([]messageListItem, len(messages))
	for i, m := range messages {
		preview := ""
		if m.TextMessage != nil {
			preview = lo.Substring(m.TextMessage.Text, 0, previewLength)
		} else if m.HashedMessage != nil {
			preview = "hash: " + lo.Substring(m.HashedMessage.Hash, 0, previewLength)
		}

		var recipients []recipientItem
		if len(m.Recipients) > 0 {
			recipients = make([]recipientItem, len(m.Recipients))
			for j, r := range m.Recipients {
				recipients[j] = recipientItem{
					PhoneNumber: r.PhoneNumber,
					State:       r.State,
					Error:       r.Error,
				}
			}
		}

		createdAt := time.Time{}
		if t, ok := m.States[string(smsgateway.ProcessingStatePending)]; ok {
			createdAt = t
		}

		items[i] = messageListItem{
			ID:          m.ID,
			DeviceID:    m.DeviceID,
			State:       m.State,
			Recipients:  recipients,
			TextPreview: preview,
			CreatedAt:   createdAt,
		}
	}

	return c.JSON(listMessagesResponse{
		Items: items,
		Total: total,
	})
}

type sendMessageRequest struct {
	PhoneNumbers []string `json:"phoneNumbers"        validate:"required,min=1,max=100,dive,required"`
	Text         string   `json:"text"                validate:"required,min=1,max=65535"`
	SimNumber    *uint8   `json:"simNumber,omitempty" validate:"omitempty,min=1,max=3"`
	Priority     *int8    `json:"priority,omitempty"  validate:"omitempty,min=-128,max=127"`
}

// SendMessage sends a new text message.
//
//	@Summary		Send message
//	@Description	Sends a new text message to one or more recipients.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			body	body		sendMessageRequest	true	"Message to send"
//	@Success		200		{object}	smsgateway.MessageState
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Router			/messages [post]
func (h *MessagesHandler) send(c *fiber.Ctx, req *sendMessageRequest) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	msg := smsgateway.Message{
		PhoneNumbers: req.PhoneNumbers,
		TextMessage: &smsgateway.TextMessage{
			Text: req.Text,
		},
		SimNumber: req.SimNumber,

		ID:                 "",
		DeviceID:           "",
		Message:            "",
		DataMessage:        nil,
		IsEncrypted:        false,
		WithDeliveryReport: nil,
		Priority:           smsgateway.PriorityDefault,
		TTL:                nil,
		ValidUntil:         nil,
		ScheduleAt:         nil,
	}

	if req.Priority != nil {
		msg.Priority = smsgateway.MessagePriority(*req.Priority)
	}

	state, err := sdkClient.Send(c.Context(), msg)
	if err != nil {
		h.logger.Warn("failed to send message", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to send message")
	}

	return c.JSON(state)
}

// GetMessage returns message details by ID.
//
//	@Summary		Get message
//	@Description	Returns full message state by message ID.
//	@Tags			messages
//	@Produce		json
//	@Param			id	path		string	true	"Message ID"
//	@Success		200	{object}	smsgateway.MessageState
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Failure		404	{object}	fiberfx.ErrorResponse
//	@Router			/messages/{id} [get]
func (h *MessagesHandler) get(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "message ID is required")
	}

	state, err := sdkClient.GetState(c.Context(), id)
	if err != nil {
		h.logger.Warn("failed to get message", zap.String("id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to get message")
	}

	return c.JSON(state)
}
