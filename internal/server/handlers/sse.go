package handlers

import (
	"fmt"

	"github.com/android-sms-gateway/web-dashboard/internal/server/events"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/android-sms-gateway/web-dashboard/internal/server/webhooks"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SSEHandler struct {
	handler.Base

	hub         *events.Hub
	webhooksSvc *webhooks.Service
	logger      *zap.Logger
}

func NewSSEHandler(
	hub *events.Hub,
	webhooksSvc *webhooks.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &SSEHandler{
		Base: handler.Base{
			Validator: validator,
		},
		hub:         hub,
		webhooksSvc: webhooksSvc,
		logger:      logger,
	}
}

func (h *SSEHandler) Register(r fiber.Router) {
	g := r.Group("/events", session.AuthRequired())

	g.Get("", h.stream)
}

func (h *SSEHandler) stream(c *fiber.Ctx) error {
	login, password, err := session.GetCredentials(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "failed to get credentials")
	}

	id := uuid.NewString()

	conn := events.NewConnection(id, login)
	h.hub.Add(conn)
	h.webhooksSvc.OnConnect(login, password)

	h.logger.Info("sse client connected",
		zap.String("id", id),
		zap.String("user", login),
		zap.Int("active", h.hub.Len()),
	)

	defer func() {
		h.hub.Remove(id, login)
		h.webhooksSvc.OnDisconnect(login, password)
		h.logger.Info("sse client disconnected",
			zap.String("id", id),
			zap.String("user", login),
			zap.Int("active", h.hub.Len()),
		)
	}()

	if streamErr := conn.Stream(c); streamErr != nil {
		return fmt.Errorf("failed to stream sse: %w", streamErr)
	}

	return nil
}
