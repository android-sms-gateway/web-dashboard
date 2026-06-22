package events

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"events",
		logger.WithNamedLogger("events"),

		fx.Provide(NewHub),
	)
}
