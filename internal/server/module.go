package server

import (
	"net/http"
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/server/docs"
	"github.com/android-sms-gateway/web-dashboard/internal/server/events"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/android-sms-gateway/web-dashboard/internal/server/webhooks"
	"github.com/android-sms-gateway/web-dashboard/internal/web"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/health"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		logger.WithNamedLogger("server"),

		fx.Provide(func(log *zap.Logger) fiberfx.Options {
			opts := fiberfx.Options{}
			opts.WithErrorHandler(fiberfx.NewJSONErrorHandler(log))
			opts.WithMetrics()
			return opts
		}),
		fx.Supply(docs.SwaggerInfo),

		fx.Provide(
			health.NewHandler,
			openapi.NewHandler,
			fx.Private,
		),

		handlers.Module(),
		events.Module(),
		webhooks.Module(),

		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, callbackHandler *handlers.CallbackHandler, healthHandler *health.Handler, openapiHandler *openapi.Handler, app *fiber.App) {
					// Health endpoint
					healthHandler.Register(app)
					callbackHandler.Register(app)

					v1 := app.Group("/api/v1")
					openapiHandler.Register(v1.Group("/docs"))

					v1.Use(
						validation.Middleware,
						session.New(nil),
					)

					for _, h := range handlers {
						h.Register(v1)
					}

					app.Use(filesystem.New(
						filesystem.Config{
							Root:               http.FS(web.StaticFS()),
							PathPrefix:         "",
							Browse:             false,
							Index:              "index.html",
							MaxAge:             int(time.Hour / time.Second),
							NotFoundFile:       "",
							ContentTypeCharset: "",
						},
					))
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
