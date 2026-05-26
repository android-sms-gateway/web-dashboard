package internal

import (
	"context"

	"github.com/android-sms-gateway/web-dashboard/internal/config"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/validatorfx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// badgerfx.Module(),
		// bunfx.Module(),
		// cachefx.Module(),
		fiberfx.Module(),
		// gocqlfx.Module(),
		// gocqlxfx.Module(),
		// sqlfx.Module(),
		// goosefx.Module(),
		// gormfx.Module(),
		healthfx.Module(),
		// openrouterfx.Module(),
		// redisfx.Module(),
		// sqlxfx.Module(),
		// telegofx.Module(true),
		validatorfx.Module(),
		// watermillfx.Module(),
		//
		// APP MODULES
		config.Module(),
		server.Module(),
		//
		// BUSINESS MODULES
		gateway.Module(),

		fx.Supply(version),

		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger, cfg gateway.Config) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("web-dashboard started", zap.String("gateway_url", cfg.URL))
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("web-dashboard stopped")
					return nil
				},
			})
		}),
	).Run()
}
