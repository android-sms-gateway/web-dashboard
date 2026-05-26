package handlers

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"handlers",
		logger.WithNamedLogger("handlers"),

		fx.Provide(
			fx.Annotate(NewAuthHandler, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(NewDashboardHandler, fx.ResultTags(`group:"handlers"`)),
		),
	)
}
