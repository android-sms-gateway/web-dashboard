package handlers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
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
	DevicesActive   int `json:"devicesActive"`
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
//	@Failure		502	{object}	fiberfx.ErrorResponse
//	@Router			/stats [get]
func (h *DashboardHandler) stats(c *fiber.Ctx) error {
	client := client.Get(c)
	if client == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	ctx := c.Context()

	total, pending, failed, err := h.countMessages(ctx, client)
	if err != nil {
		h.logger.Warn("failed to list messages", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to list messages")
	}

	devices, err := client.ListDevices(ctx)
	if err != nil {
		h.logger.Warn("failed to list devices", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to list devices")
	}

	activeCount := 0
	onlineCount := 0
	for _, d := range devices {
		if d.DeletedAt != nil {
			continue
		}
		activeCount++
		if time.Since(d.LastSeen) < 15*time.Minute {
			onlineCount++
		}
	}

	return c.JSON(statsResponse{
		DevicesActive:   activeCount,
		DevicesOnline:   onlineCount,
		DevicesTotal:    len(devices),
		MessagesSent:    total - pending - failed,
		MessagesPending: pending,
		MessagesFailed:  failed,
	})
}

func (h *DashboardHandler) countMessages(
	ctx context.Context,
	client *smsgateway.Client,
) (int, int, int, error) {
	var wg sync.WaitGroup
	var t, p, f int
	var collectErr error
	var errMu sync.Mutex

	limit := 1
	pendingState := smsgateway.ProcessingStatePending
	failedState := smsgateway.ProcessingStateFailed

	recordErr := func(e error) {
		errMu.Lock()
		collectErr = errors.Join(collectErr, e)
		errMu.Unlock()
	}

	wg.Go(func() {
		_, n, listErr := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
			State:          nil,
			Limit:          &limit,
			From:           nil,
			To:             nil,
			DeviceID:       nil,
			Offset:         nil,
			IncludeContent: nil,
		})
		if listErr != nil {
			recordErr(listErr)
			return
		}
		t = n
	})

	wg.Go(func() {
		_, n, listErr := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
			State:          (*string)(&pendingState),
			Limit:          &limit,
			From:           nil,
			To:             nil,
			DeviceID:       nil,
			Offset:         nil,
			IncludeContent: nil,
		})
		if listErr != nil {
			recordErr(listErr)
			return
		}
		p = n
	})

	wg.Go(func() {
		_, n, listErr := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
			State:          (*string)(&failedState),
			Limit:          &limit,
			From:           nil,
			To:             nil,
			DeviceID:       nil,
			Offset:         nil,
			IncludeContent: nil,
		})
		if listErr != nil {
			recordErr(listErr)
			return
		}
		f = n
	})

	wg.Wait()

	return t, p, f, collectErr
}
