package handlers

import (
	"encoding/json"

	"github.com/android-sms-gateway/web-dashboard/internal/server/events"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	sseType, ok := eventMapping()[body.Event]
	if !ok {
		h.logger.Warn("unknown webhook event",
			zap.String("event", body.Event),
			zap.String("user", userID),
		)
		return c.JSON(fiber.Map{"status": "ok"})
	}

	event := events.Event{
		Type:    sseType,
		Payload: body.Payload,
	}

	if err := h.hub.SendToUser(userID, event); err != nil {
		h.logger.Error("failed to send event",
			zap.String("event", body.Event),
			zap.String("user", userID),
			zap.Error(err),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to send event")
	}

	h.logger.Debug("webhook event delivered",
		zap.String("user", userID),
		zap.String("source_event", body.Event),
		zap.String("sse_type", string(sseType)),
	)

	return c.JSON(fiber.Map{"status": "ok"})
}
