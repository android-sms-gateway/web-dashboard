package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/events"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	errUnknownEvent = errors.New("unknown webhook event")
	errInvalidJSON  = errors.New("invalid JSON payload")
)

func eventMapping() map[string]events.EventType {
	return map[string]events.EventType{
		"sms:received":      events.EventMessageReceived,
		"sms:sent":          events.EventMessageStateChanged,
		"sms:delivered":     events.EventMessageStateChanged,
		"sms:failed":        events.EventMessageStateChanged,
		"sms:data-received": events.EventMessageReceived,
		"mms:received":      events.EventMessageReceived,
		"mms:downloaded":    events.EventMessageReceived,
		"system:ping":       events.EventDeviceStatusChanged,
	}
}

func eventState(sourceEvent string) (string, bool) {
	state, ok := map[string]string{
		"sms:sent":      string(smsgateway.ProcessingStateSent),
		"sms:delivered": string(smsgateway.ProcessingStateDelivered),
		"sms:failed":    string(smsgateway.ProcessingStateFailed),
	}[sourceEvent]

	return state, ok
}

func injectState(sourceEvent string, payload json.RawMessage) (json.RawMessage, error) {
	data := make(map[string]any)
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &data); err != nil {
			return nil, errInvalidJSON
		}
	}

	state, ok := eventState(sourceEvent)
	if !ok {
		return nil, errUnknownEvent
	}
	data["state"] = state

	out, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return out, nil
}

func normalizeEvent(sourceEvent string, payload json.RawMessage) (events.Event, error) {
	sseType, ok := eventMapping()[sourceEvent]
	if !ok {
		return events.Event{}, errUnknownEvent
	}

	if sseType != events.EventMessageStateChanged {
		return events.Event{Type: sseType, Payload: payload}, nil
	}

	payload, err := injectState(sourceEvent, payload)
	if err != nil {
		return events.Event{}, err
	}

	return events.Event{Type: sseType, Payload: payload}, nil
}

type CallbackHandler struct {
	hub    *events.Hub
	logger *zap.Logger
}

func NewCallbackHandler(hub *events.Hub, logger *zap.Logger) *CallbackHandler {
	return &CallbackHandler{
		hub:    hub,
		logger: logger,
	}
}

func (h *CallbackHandler) Register(r fiber.Router) {
	r.Post("/api/webhooks/callback/:userId", h.handle)
}

type callbackBody struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func (h *CallbackHandler) handle(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "userId is required")
	}

	body := new(callbackBody)
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	event, err := normalizeEvent(body.Event, body.Payload)
	if err != nil {
		if errors.Is(err, errUnknownEvent) {
			h.logger.Warn("unknown webhook event",
				zap.String("event", body.Event),
				zap.String("user", userID),
			)
			return c.JSON(fiber.Map{"status": "ok"})
		}

		if errors.Is(err, errInvalidJSON) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid event payload")
		}

		h.logger.Error("failed to normalize webhook event",
			zap.String("event", body.Event),
			zap.String("user", userID),
			zap.Error(err),
		)

		return fiber.NewError(fiber.StatusInternalServerError, "failed to process event")
	}

	if sendErr := h.hub.SendToUser(userID, event); sendErr != nil {
		h.logger.Error("failed to send event",
			zap.String("event", body.Event),
			zap.String("user", userID),
			zap.Error(sendErr),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to send event")
	}

	if event.Type == events.EventMessageReceived || event.Type == events.EventMessageStateChanged {
		statsEvent := events.Event{
			Type:    events.EventStatsUpdated,
			Payload: json.RawMessage("{}"),
		}
		if sendErr := h.hub.SendToUser(userID, statsEvent); sendErr != nil {
			h.logger.Error("failed to send stats event",
				zap.String("user", userID),
				zap.Error(sendErr),
			)
		}
	}

	h.logger.Debug("webhook event delivered",
		zap.String("user", userID),
		zap.String("source_event", body.Event),
		zap.String("sse_type", string(event.Type)),
	)

	return c.JSON(fiber.Map{"status": "ok"})
}
