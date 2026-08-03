package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/android-sms-gateway/web-dashboard/internal/dashboard"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/cachefx/cache"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/validatorfx"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func newUpstream(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/devices":
			_, _ = w.Write([]byte("[]"))
		case r.Method == http.MethodGet && r.URL.Path == "/messages":
			w.Header().Set("X-Total-Count", "0")
			_, _ = w.Write([]byte("[]"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	return server
}

func newTestApp(upstreamURL string) *fiber.App {
	app := fiber.New()

	v1 := app.Group("/api/v1", validation.Middleware, session.New(nil))

	validator := validatorfx.New()
	logger := zap.NewNop()
	factory := gateway.NewFactory(upstreamURL)

	authHandler := handlers.NewAuthHandler(factory, validator, logger)
	dashboardHandler := handlers.NewDashboardHandler(
		dashboard.NewService(factory, cache.NewMemory(0), logger),
		validator,
		logger,
	)

	authHandler.Register(v1)
	dashboardHandler.Register(v1)

	return app
}

func login(t *testing.T, app *fiber.App, login, password string) string {
	t.Helper()

	body := fmt.Sprintf(`{"login":%q,"password":%q}`, login, password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want 200", resp.StatusCode)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" {
			return cookie.Value
		}
	}
	t.Fatal("no session cookie set")

	return ""
}

func TestTrendsInvalidDays(t *testing.T) {
	app := newTestApp(newUpstream(t).URL)
	cookie := login(t, app, "user", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/trends?days=5", nil)
	req.Header.Set("Cookie", "session_id="+cookie)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("trends request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestTrendsUnauthorized(t *testing.T) {
	app := newTestApp(newUpstream(t).URL)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/trends", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("trends request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestTrendsOK(t *testing.T) {
	app := newTestApp(newUpstream(t).URL)
	cookie := login(t, app, "user", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/trends", nil)
	req.Header.Set("Cookie", "session_id="+cookie)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("trends request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var trends dashboard.Trends
	if decodeErr := json.NewDecoder(resp.Body).Decode(&trends); decodeErr != nil {
		t.Fatalf("failed to decode response: %v", decodeErr)
	}

	if trends.Days != 7 {
		t.Errorf("Days = %d, want 7", trends.Days)
	}
	if len(trends.MessageVolume) != 7 {
		t.Errorf("len(MessageVolume) = %d, want 7", len(trends.MessageVolume))
	}
}
