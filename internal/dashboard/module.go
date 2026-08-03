package dashboard

import (
	"github.com/android-sms-gateway/web-dashboard/internal/cache"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

// Module provides the dashboard service.
func Module() fx.Option {
	return fx.Module(
		"dashboard",
		logger.WithNamedLogger("dashboard"),

		fx.Provide(
			func(factory cache.Factory) (cache.Cache, error) {
				return factory.New("dashboard")
			},
			fx.Private,
		),
		fx.Provide(NewService),
	)
}
