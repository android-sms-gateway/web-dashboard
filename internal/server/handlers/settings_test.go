package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/validatorfx"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func newSettingsApp(t *testing.T, upstreamHandler http.HandlerFunc) *fiber.App {
	t.Helper()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/devices" {
			_, _ = w.Write([]byte("[]"))
			return
		}
		upstreamHandler(w, r)
	}))
	t.Cleanup(upstream.Close)

	app := fiber.New()
	v1 := app.Group("/api/v1", validation.Middleware, session.New(nil))

	validator := validatorfx.New()
	logger := zap.NewNop()
	factory := gateway.NewFactory(upstream.URL)

	authHandler := handlers.NewAuthHandler(factory, validator, logger)
	settingsHandler := handlers.NewSettingsHandler(factory, validator, logger)

	authHandler.Register(v1)
	settingsHandler.Register(v1)

	return app
}

func patchSettings(t *testing.T, app *fiber.App, cookie string, body string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "session_id="+cookie)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("settings PATCH failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	return resp
}

// The url constraint on gateway.cloud_url must keep rejecting malformed URLs.
func TestUpdateSettingsRejectsInvalidCloudURL(t *testing.T) {
	payload := `{"gateway": {"cloud_url": "not-a-url"}}`

	app := newSettingsApp(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("{}"))
	})

	cookie := login(t, app)

	resp := patchSettings(t, app, cookie, payload)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}
