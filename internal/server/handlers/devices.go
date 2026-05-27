package handlers

import (
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/client"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DevicesHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewDevicesHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &DevicesHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *DevicesHandler) Register(r fiber.Router) {
	g := r.Group("/devices", session.AuthRequired(), client.New(h.gatewaySvc))

	g.Get("", h.list)
	g.Delete(":id", h.delete)
}

type deviceListItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastSeen  time.Time `json:"lastSeen"`
	IsOnline  bool      `json:"isOnline"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListDevices returns registered devices.
//
//	@Summary		List devices
//	@Description	Returns all registered devices.
//	@Tags			devices
//	@Produce		json
//	@Success		200	{array}		deviceListItem
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Router			/devices [get]
func (h *DevicesHandler) list(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	devices, err := sdkClient.ListDevices(c.Context())
	if err != nil {
		h.logger.Warn("failed to list devices", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to list devices")
	}

	items := make([]deviceListItem, len(devices))
	for i, d := range devices {
		items[i] = deviceListItem{
			ID:        d.ID,
			Name:      d.Name,
			LastSeen:  d.LastSeen,
			IsOnline:  d.DeletedAt == nil,
			CreatedAt: d.CreatedAt,
		}
	}

	return c.JSON(items)
}

// DeleteDevice removes a device by ID.
//
//	@Summary		Delete device
//	@Description	Removes a device by its ID.
//	@Tags			devices
//	@Produce		json
//	@Param			id	path	string	true	"Device ID"
//	@Success		204
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Failure		404	{object}	fiberfx.ErrorResponse
//	@Router			/devices/{id} [delete]
func (h *DevicesHandler) delete(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "device ID is required")
	}

	if err := sdkClient.DeleteDevice(c.Context(), id); err != nil {
		h.logger.Warn("failed to delete device", zap.String("id", id), zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to delete device")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
