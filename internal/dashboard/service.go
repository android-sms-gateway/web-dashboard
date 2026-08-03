package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/go-core-fx/cachefx/cache"
	"go.uber.org/zap"
)

const (
	statsTimeout  = 10 * time.Second
	trendsTimeout = 20 * time.Second
)

// Service provides dashboard statistics and trends.
type Service struct {
	gatewaySvc *gateway.Factory
	logger     *zap.Logger
	cache      *trendsCache
}

// NewService creates a new dashboard service.
func NewService(gatewaySvc *gateway.Factory, c cache.Cache, logger *zap.Logger) *Service {
	return &Service{
		gatewaySvc: gatewaySvc,
		logger:     logger,
		cache:      newTrendsCache(c, logger),
	}
}

// Stats returns aggregated statistics for the given user.
func (s *Service) Stats(ctx context.Context, login, password string) (Stats, error) {
	ctx, cancel := context.WithTimeout(ctx, statsTimeout)
	defer cancel()

	client := s.gatewaySvc.NewClient(login, password)

	total, pending, failed, err := countMessagesRange(ctx, client, nil, nil)
	if err != nil {
		return Stats{}, err
	}

	devices, err := client.ListDevices(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("failed to list devices: %w", err)
	}

	activeCount := 0
	onlineCount := 0
	for _, d := range devices {
		if d.DeletedAt != nil {
			continue
		}
		activeCount++
		if time.Since(d.LastSeen) < DeviceOnlineIn {
			onlineCount++
		}
	}

	return Stats{
		DevicesActive:   activeCount,
		DevicesOnline:   onlineCount,
		DevicesTotal:    len(devices),
		MessagesSent:    max(total-pending-failed, 0),
		MessagesPending: pending,
		MessagesFailed:  failed,
	}, nil
}

// Trends returns per-day message volume and device activity for the given user.
func (s *Service) Trends(ctx context.Context, login, password string, days int) (Trends, error) {
	if cached, ok := s.cache.get(ctx, login, days); ok {
		return cached, nil
	}

	ctx, cancel := context.WithTimeout(ctx, trendsTimeout)
	defer cancel()

	client := s.gatewaySvc.NewClient(login, password)

	buckets := trendsBuckets(time.Now().UTC(), days)

	volumes, activity, err := sweepDays(ctx, client, buckets)
	if err != nil {
		return Trends{}, err
	}

	trends := Trends{
		Days:           days,
		MessageVolume:  volumes,
		DeviceActivity: activity,
	}

	s.cache.set(ctx, login, days, trends)

	return trends, nil
}
