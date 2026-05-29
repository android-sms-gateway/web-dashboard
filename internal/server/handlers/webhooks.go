package handlers

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/client"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WebhooksHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewWebhooksHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &WebhooksHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *WebhooksHandler) Register(r fiber.Router) {
	g := r.Group("/webhooks", session.AuthRequired(), client.New(h.gatewaySvc))

	g.Get("", h.list)
	g.Post("", validation.DecorateWithBodyEx(h.Validator, h.create))
	g.Delete(":id", h.delete)
}

type webhookListItem struct {
	ID       string  `json:"id"`
	URL      string  `json:"url"`
	Event    string  `json:"event"`
	DeviceID *string `json:"deviceId,omitempty"`
}

// ListWebhooks returns registered webhooks.
//
//	@Summary		List webhooks
//	@Description	Returns all registered webhooks.
//	@Tags			webhooks
//	@Produce		json
//	@Success		200	{array}		webhookListItem
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Router			/webhooks [get]
func (h *WebhooksHandler) list(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	webhooks, err := sdkClient.ListWebhooks(c.Context())
	if err != nil {
		h.logger.Warn("failed to list webhooks", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to list webhooks")
	}

	items := make([]webhookListItem, len(webhooks))
	for i, w := range webhooks {
		items[i] = webhookListItem{
			ID:       w.ID,
			URL:      w.URL,
			Event:    w.Event,
			DeviceID: w.DeviceID,
		}
	}

	return c.JSON(items)
}

type createWebhookRequest struct {
	URL      string  `json:"url"                validate:"required,http_url"`
	Event    string  `json:"event"              validate:"required"`
	DeviceID *string `json:"deviceId,omitempty"`
}

// CreateWebhook registers a new webhook.
//
//	@Summary		Create webhook
//	@Description	Registers a new webhook for the specified event.
//	@Tags			webhooks
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createWebhookRequest	true	"Webhook to create"
//	@Success		200		{object}	smsgateway.Webhook
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Router			/webhooks [post]
func (h *WebhooksHandler) create(c *fiber.Ctx, req *createWebhookRequest) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	result, err := sdkClient.RegisterWebhook(c.Context(), smsgateway.Webhook{
		ID:       "",
		URL:      req.URL,
		Event:    req.Event,
		DeviceID: req.DeviceID,
	})
	if err != nil {
		h.logger.Warn("failed to create webhook", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to create webhook")
	}

	return c.JSON(result)
}

// DeleteWebhook removes a webhook by ID.
//
//	@Summary		Delete webhook
//	@Description	Removes a webhook by its ID.
//	@Tags			webhooks
//	@Produce		json
//	@Param			id	path	string	true	"Webhook ID"
//	@Success		204
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Failure		404	{object}	fiberfx.ErrorResponse
//	@Router			/webhooks/{id} [delete]
func (h *WebhooksHandler) delete(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "webhook ID is required")
	}

	if err := sdkClient.DeleteWebhook(c.Context(), id); err != nil {
		h.logger.Warn("failed to delete webhook", zap.String("id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to delete webhook")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
