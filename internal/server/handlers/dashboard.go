package handlers

import (
	"github.com/android-sms-gateway/web-dashboard/internal/dashboard"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const defaultTrendsDays = 7

type DashboardHandler struct {
	handler.Base

	dashboardSvc *dashboard.Service
	logger       *zap.Logger
}

func NewDashboardHandler(
	dashboardSvc *dashboard.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &DashboardHandler{
		Base: handler.Base{
			Validator: validator,
		},
		dashboardSvc: dashboardSvc,
		logger:       logger,
	}
}

func (h *DashboardHandler) Register(r fiber.Router) {
	g := r.Group("/stats", session.AuthRequired())

	g.Get("", h.stats)
	g.Get("/trends", h.trends)
}

// Stats returns dashboard statistics.
//
//	@Summary		Dashboard stats
//	@Description	Returns aggregated statistics for the dashboard (devices, messages).
//	@Tags			dashboard
//	@Produce		json
//	@Success		200	{object}	dashboard.Stats
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Failure		502	{object}	fiberfx.ErrorResponse
//	@Router			/stats [get]
func (h *DashboardHandler) stats(c *fiber.Ctx) error {
	login, password, err := session.GetCredentials(c)
	if err != nil {
		h.logger.Warn("failed to get credentials", zap.Error(err))
		return fiber.NewError(fiber.StatusUnauthorized, "failed to get credentials")
	}

	stats, err := h.dashboardSvc.Stats(c.Context(), login, password)
	if err != nil {
		h.logger.Warn("failed to get stats", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to get stats")
	}

	return c.JSON(stats)
}

type trendsQuery struct {
	Days *int `query:"days" validate:"omitempty,oneof=7 14 30"`
}

// Trends returns per-day dashboard trends.
//
//	@Summary		Dashboard trends
//	@Description	Returns per-day message volume and device activity for the last 7, 14, or 30 days.
//	@Tags			dashboard
//	@Produce		json
//	@Param			days	query		int	false	"Number of days (7, 14, or 30)"	default(7)
//	@Success		200		{object}	dashboard.Trends
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Failure		502		{object}	fiberfx.ErrorResponse
//	@Router			/stats/trends [get]
func (h *DashboardHandler) trends(c *fiber.Ctx) error {
	query := new(trendsQuery)
	if err := h.QueryParserValidator(c, query); err != nil {
		h.logger.Warn("failed to parse query", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "failed to parse query")
	}

	login, password, err := session.GetCredentials(c)
	if err != nil {
		h.logger.Warn("failed to get credentials", zap.Error(err))
		return fiber.NewError(fiber.StatusUnauthorized, "failed to get credentials")
	}

	days := defaultTrendsDays
	if query.Days != nil {
		days = *query.Days
	}

	trends, err := h.dashboardSvc.Trends(c.Context(), login, password, days)
	if err != nil {
		h.logger.Warn("failed to get trends", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to get trends")
	}

	return c.JSON(trends)
}
