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

type SettingsHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewSettingsHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &SettingsHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *SettingsHandler) Register(r fiber.Router) {
	g := r.Group("/settings", session.AuthRequired(), client.New(h.gatewaySvc))

	g.Get("", h.get)
	g.Patch("", validation.DecorateWithBodyEx(h.Validator, h.update))
}

// GetSettings returns current device settings.
//
//	@Summary		Get settings
//	@Description	Returns current device settings.
//	@Tags			settings
//	@Produce		json
//	@Success		200	{object}	smsgateway.DeviceSettings
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Router			/settings [get]
func (h *SettingsHandler) get(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	settings, err := sdkClient.GetSettings(c.Context())
	if err != nil {
		h.logger.Warn("failed to get settings", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to get settings")
	}

	return c.JSON(settings)
}

// UpdateSettings partially updates device settings.
//
//	@Summary		Update settings
//	@Description	Partially updates device settings. Only included fields will be changed.
//	@Tags			settings
//	@Accept			json
//	@Produce		json
//	@Param			body	body		smsgateway.DeviceSettings	true	"Settings to update"
//	@Success		200		{object}	smsgateway.DeviceSettings
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Router			/settings [patch]
func (h *SettingsHandler) update(c *fiber.Ctx, req *smsgateway.DeviceSettings) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	settings, err := sdkClient.UpdateSettings(c.Context(), *req)
	if err != nil {
		h.logger.Warn("failed to update settings", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to update settings")
	}

	return c.JSON(settings)
}
