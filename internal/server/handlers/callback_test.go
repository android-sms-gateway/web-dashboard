package handlers_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/android-sms-gateway/web-dashboard/internal/server/events"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const streamTimeout = 5 * time.Second

func newCallbackApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	hub := events.NewHub()

	handlers.NewCallbackHandler(hub, zap.NewNop()).Register(app)

	app.Get("/test-stream/:user", func(c *fiber.Ctx) error {
		conn := events.NewConnection(uuid.NewString(), c.Params("user"))
		hub.Add(conn)
		if err := conn.SendEvent(events.Event{
			Type:    events.EventDeviceStatusChanged,
			Payload: json.RawMessage("{}"),
		}); err != nil {
			return err
		}

		return conn.Stream(c, func() {
			hub.Remove(conn.ID, conn.User)
		})
	})

	return app
}

func newCallbackServer(t *testing.T) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	t.Cleanup(func() {
		_ = ln.Close()
	})

	go func() {
		_ = newCallbackApp().Listener(ln)
	}()

	return "http://" + ln.Addr().String()
}

func openStream(t *testing.T, appURL string) chan string {
	t.Helper()

	resp, err := http.Get(appURL + "/test-stream/user")
	if err != nil {
		t.Fatalf("failed to open stream: %v", err)
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	frames := make(chan string, 32)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			frames <- scanner.Text()
		}
		close(frames)
	}()

	for {
		select {
		case line, ok := <-frames:
			if !ok {
				t.Fatal("stream closed before bootstrap event")
			}
			if strings.HasPrefix(line, "data: ") {
				return frames
			}
		case <-time.After(streamTimeout):
			t.Fatal("timed out waiting for bootstrap event")
		}
	}
}

func nextData(t *testing.T, frames chan string) string {
	t.Helper()

	for {
		select {
		case line, ok := <-frames:
			if !ok {
				t.Fatal("stream closed before event")
			}
			if data, found := strings.CutPrefix(line, "data: "); found {
				return data
			}
		case <-time.After(streamTimeout):
			t.Fatal("timed out waiting for SSE event")
		}
	}
}

func nextEvent(t *testing.T, frames chan string) events.Event {
	t.Helper()

	var got events.Event
	if err := json.Unmarshal([]byte(nextData(t, frames)), &got); err != nil {
		t.Fatalf("failed to decode SSE event: %v", err)
	}

	return got
}

func postCallback(t *testing.T, appURL string, body string) *http.Response {
	t.Helper()

	resp, err := http.Post(appURL+"/api/webhooks/callback/user", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	return resp
}

func TestCallbackMessageStateChanged(t *testing.T) {
	cases := []struct {
		event     string
		wantState string
	}{
		{event: "sms:sent", wantState: "Sent"},
		{event: "sms:delivered", wantState: "Delivered"},
		{event: "sms:failed", wantState: "Failed"},
	}

	for _, tc := range cases {
		t.Run(tc.event, func(t *testing.T) {
			appURL := newCallbackServer(t)
			frames := openStream(t, appURL)

			resp := postCallback(t, appURL, `{"event":"`+tc.event+`","payload":{"messageId":"m1"}}`)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200", resp.StatusCode)
			}

			got := nextEvent(t, frames)
			if got.Type != events.EventMessageStateChanged {
				t.Errorf("type = %q, want %q", got.Type, events.EventMessageStateChanged)
			}

			var payload struct {
				MessageID string `json:"messageId"`
				State     string `json:"state"`
			}
			if err := json.Unmarshal(got.Payload, &payload); err != nil {
				t.Fatalf("failed to decode payload: %v", err)
			}
			if payload.MessageID != "m1" {
				t.Errorf("messageId = %q, want %q", payload.MessageID, "m1")
			}
			if payload.State != tc.wantState {
				t.Errorf("state = %q, want %q", payload.State, tc.wantState)
			}

			stats := nextEvent(t, frames)
			if stats.Type != events.EventStatsUpdated {
				t.Errorf("type = %q, want %q", stats.Type, events.EventStatsUpdated)
			}
		})
	}
}

func TestCallbackMessageReceived(t *testing.T) {
	appURL := newCallbackServer(t)
	frames := openStream(t, appURL)

	resp := postCallback(t, appURL, `{"event":"sms:received","payload":{"sender":"+1234567890","message":"hello"}}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	got := nextEvent(t, frames)
	if got.Type != events.EventMessageReceived {
		t.Errorf("type = %q, want %q", got.Type, events.EventMessageReceived)
	}

	var payload struct {
		Sender  string `json:"sender"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(got.Payload, &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if payload.Sender != "+1234567890" {
		t.Errorf("sender = %q, want %q", payload.Sender, "+1234567890")
	}
	if payload.Message != "hello" {
		t.Errorf("message = %q, want %q", payload.Message, "hello")
	}

	stats := nextEvent(t, frames)
	if stats.Type != events.EventStatsUpdated {
		t.Errorf("type = %q, want %q", stats.Type, events.EventStatsUpdated)
	}
}

func TestCallbackPing(t *testing.T) {
	appURL := newCallbackServer(t)
	frames := openStream(t, appURL)

	resp := postCallback(t, appURL, `{"event":"system:ping","payload":{}}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	got := nextEvent(t, frames)
	if got.Type != events.EventDeviceStatusChanged {
		t.Errorf("type = %q, want %q", got.Type, events.EventDeviceStatusChanged)
	}
}

func TestCallbackUnknownEvent(t *testing.T) {
	appURL := newCallbackServer(t)
	frames := openStream(t, appURL)

	resp := postCallback(t, appURL, `{"event":"foo:bar","payload":{}}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	resp = postCallback(t, appURL, `{"event":"sms:received","payload":{"sender":"+1","message":"m"}}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	got := nextEvent(t, frames)
	if got.Type != events.EventMessageReceived {
		t.Errorf("type = %q, want %q (unknown event must not be delivered)", got.Type, events.EventMessageReceived)
	}
}

func TestCallbackInvalidPayload(t *testing.T) {
	appURL := newCallbackServer(t)

	resp := postCallback(t, appURL, `{"event":"sms:sent","payload":"not-an-object"}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCallbackInvalidBody(t *testing.T) {
	appURL := newCallbackServer(t)

	resp := postCallback(t, appURL, `not-json`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCallbackMissingUser(t *testing.T) {
	appURL := newCallbackServer(t)

	resp, err := http.Post(appURL+"/api/webhooks/callback/", "application/json", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatalf("callback request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404 (empty userId must not match the route)", resp.StatusCode)
	}
}
