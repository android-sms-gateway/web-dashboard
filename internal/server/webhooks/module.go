package webhooks

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"webhook-lifecycle",
		logger.WithNamedLogger("webhook-lifecycle"),

		fx.Provide(NewService),
	)
}
