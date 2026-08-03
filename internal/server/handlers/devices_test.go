package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/validatorfx"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func newDevicesApp(t *testing.T, devicesJSON string) *fiber.App {
	t.Helper()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/devices" {
			_, _ = w.Write([]byte(devicesJSON))
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(upstream.Close)

	app := fiber.New()

	v1 := app.Group("/api/v1", validation.Middleware, session.New(nil))

	validator := validatorfx.New()
	logger := zap.NewNop()
	factory := gateway.NewFactory(upstream.URL)

	authHandler := handlers.NewAuthHandler(factory, validator, logger)
	devicesHandler := handlers.NewDevicesHandler(factory, validator, logger)

	authHandler.Register(v1)
	devicesHandler.Register(v1)

	return app
}

func TestListDevicesOnline(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	stale := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	deleted := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)

	devices := fmt.Sprintf(`[
		{"id":"online","name":"Online Device","lastSeen":%q,"createdAt":%q},
		{"id":"stale","name":"Stale Device","lastSeen":%q,"createdAt":%q},
		{"id":"deleted","name":"Deleted Device","lastSeen":%q,"deletedAt":%q,"createdAt":%q}
	]`, now, now, stale, stale, deleted, deleted, deleted)

	app := newDevicesApp(t, devices)
	cookie := login(t, app, "user", "pass")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)
	req.Header.Set("Cookie", "session_id="+cookie)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("devices request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var items []struct {
		ID       string `json:"id"`
		IsOnline bool   `json:"isOnline"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&items); decodeErr != nil {
		t.Fatalf("failed to decode response: %v", decodeErr)
	}

	if len(items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(items))
	}

	wantOnline := map[string]bool{
		"online":  true,
		"stale":   false,
		"deleted": false,
	}
	for _, item := range items {
		want, ok := wantOnline[item.ID]
		if !ok {
			continue
		}
		if item.IsOnline != want {
			t.Errorf("isOnline[%s] = %v, want %v", item.ID, item.IsOnline, want)
		}
	}
}

func TestListDevicesUnauthorized(t *testing.T) {
	app := newDevicesApp(t, `[]`)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("devices request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}
