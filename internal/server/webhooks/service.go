package webhooks

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"go.uber.org/zap"
)

const operationTimeout = 10 * time.Second
const webhookPrefix = "web-dashboard-"

//nolint:gochecknoglobals // single source of truth for auto-registered events
var registeredEvents = []smsgateway.WebhookEvent{
	smsgateway.WebhookEventSmsReceived,
	smsgateway.WebhookEventSmsSent,
	smsgateway.WebhookEventSmsDelivered,
	smsgateway.WebhookEventSmsFailed,
}

type Service struct {
	config     Config
	gatewaySvc *gateway.Factory
	logger     *zap.Logger

	mu       sync.Mutex
	refCount map[string]int
}

func NewService(config Config, gatewaySvc *gateway.Factory, logger *zap.Logger) *Service {
	return &Service{
		config:     config,
		gatewaySvc: gatewaySvc,
		logger:     logger,

		mu:       sync.Mutex{},
		refCount: make(map[string]int),
	}
}

func (s *Service) OnConnect(login, password string) {
	if !s.incrementRef(login) {
		return
	}

	go s.register(login, password)
}

func (s *Service) OnDisconnect(login, password string) {
	if s.decrementRef(login) > 0 {
		return
	}

	go s.unregister(login, password)
}

func (s *Service) incrementRef(login string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.refCount[login] > 0 {
		s.refCount[login]++

		return false
	}
	s.refCount[login] = 1

	return true
}

func (s *Service) decrementRef(login string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := s.refCount[login]
	if count == 0 {
		return 0
	}
	count--
	if count == 0 {
		delete(s.refCount, login)

		return 0
	}
	s.refCount[login] = count

	return count
}

func (s *Service) register(login, password string) {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	client := s.gatewaySvc.NewClient(login, password)

	callbackURL := strings.TrimRight(s.config.URL, "/") + "/" + url.PathEscape(login)

	for _, event := range registeredEvents {
		if _, err := client.RegisterWebhook(ctx, smsgateway.Webhook{
			ID:       webhookPrefix + event,
			URL:      callbackURL,
			Event:    event,
			DeviceID: nil,
		}); err != nil {
			s.logger.Warn("register webhook failed",
				zap.String("user", login),
				zap.String("event", event),
				zap.Error(err),
			)
		}
	}
}

func (s *Service) unregister(login, password string) {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	client := s.gatewaySvc.NewClient(login, password)
	for _, event := range registeredEvents {
		if err := client.DeleteWebhook(ctx, webhookPrefix+event); err != nil {
			s.logger.Warn("delete webhook failed",
				zap.String("user", login),
				zap.String("event", event),
				zap.Error(err),
			)
		}
	}
}
