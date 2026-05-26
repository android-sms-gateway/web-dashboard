package gateway

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"gateway",
		logger.WithNamedLogger("gateway"),
		fx.Provide(
			func(cfg Config) *Factory {
				return NewFactory(cfg.URL)
			},
		),
	)
}
