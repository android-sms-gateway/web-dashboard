package handlers

import (
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/client"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewDashboardHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &DashboardHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *DashboardHandler) Register(r fiber.Router) {
	r.Get("/stats", session.AuthRequired(), client.New(h.gatewaySvc), h.stats)
}

type statsResponse struct {
	MessagesSent    int `json:"messagesSent"`
	MessagesPending int `json:"messagesPending"`
	MessagesFailed  int `json:"messagesFailed"`
	DevicesOnline   int `json:"devicesOnline"`
	DevicesTotal    int `json:"devicesTotal"`
}

// Stats returns dashboard statistics.
//
//	@Summary		Dashboard stats
//	@Description	Returns aggregated statistics for the dashboard (devices, messages).
//	@Tags			dashboard
//	@Produce		json
//	@Success		200	{object}	statsResponse
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Router			/stats [get]
func (h *DashboardHandler) stats(c *fiber.Ctx) error {
	client := client.Get(c)
	if client == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	devices, err := client.ListDevices(c.Context())
	if err != nil {
		h.logger.Warn("failed to list devices", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to fetch dashboard stats")
	}

	onlineCount := 0
	for _, d := range devices {
		if d.DeletedAt == nil {
			onlineCount++
		}
	}

	return c.JSON(statsResponse{
		DevicesOnline:   onlineCount,
		DevicesTotal:    len(devices),
		MessagesSent:    0,
		MessagesPending: 0,
		MessagesFailed:  0,
	})
}
