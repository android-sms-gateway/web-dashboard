package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/handlers"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/validatorfx"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const deviceList = `[{"id":"dev-1","name":"Device","isOnline":true}]`

const sendRequestBody = `{"phoneNumbers":["+15555550100"],"text":"hello"}`

type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func newMessagesApp(t *testing.T, upstream http.HandlerFunc) *fiber.App {
	t.Helper()

	srv := httptest.NewServer(upstream)
	t.Cleanup(srv.Close)

	app := fiber.New(
		fiber.Config{ErrorHandler: fiberfx.NewCustomJSONErrorHandler(zap.NewNop(), handlers.ErrorFormatter)},
	)

	v1 := app.Group("/api/v1", handlers.ErrorMapping, validation.Middleware, session.New(nil))

	validator := validatorfx.New()
	logger := zap.NewNop()
	factory := gateway.NewFactory(srv.URL)

	authHandler := handlers.NewAuthHandler(factory, validator, logger)
	messagesHandler := handlers.NewMessagesHandler(factory, validator, logger)

	authHandler.Register(v1)
	messagesHandler.Register(v1)

	return app
}

func upstreamWith(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/devices" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(deviceList))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

func postMessage(t *testing.T, app *fiber.App) *http.Response {
	t.Helper()

	cookie := login(t, app)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", strings.NewReader(sendRequestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "session_id="+cookie)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("send message request failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	return resp
}

func assertErrorBody(t *testing.T, resp *http.Response, wantStatus int, want errorBody) {
	t.Helper()

	if resp.StatusCode != wantStatus {
		t.Errorf("status = %d, want %d", resp.StatusCode, wantStatus)
	}

	var got errorBody
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode error body: %v", err)
	}

	if got.Message != want.Message {
		t.Errorf("message = %q, want %q", got.Message, want.Message)
	}
	if got.Code != want.Code {
		t.Errorf("code = %d, want %d", got.Code, want.Code)
	}
}

func TestSendMessagePropagatesUpstreamValidation(t *testing.T) {
	t.Parallel()

	app := newMessagesApp(t, upstreamWith(http.StatusBadRequest, `{"message":"invalid phone number"}`))
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusBadRequest, errorBody{Message: "invalid phone number", Code: 400})
}

func TestSendMessagePropagatesUpstreamQueueLimit(t *testing.T) {
	t.Parallel()

	app := newMessagesApp(t, upstreamWith(http.StatusServiceUnavailable, `{"message":"queue limits exceeded"}`))
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusServiceUnavailable, errorBody{Message: "queue limits exceeded", Code: 503})
}

func TestSendMessageMapsUpstreamServerFailureToBadGateway(t *testing.T) {
	t.Parallel()

	app := newMessagesApp(t, upstreamWith(http.StatusInternalServerError, `{"message":"failed to handle request"}`))
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusBadGateway, errorBody{Message: "failed to handle request", Code: 502})
}

func TestSendMessageMapsTransportFailureToBadGateway(t *testing.T) {
	t.Parallel()

	upstream := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/devices" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(deviceList))
			return
		}
		panic(http.ErrAbortHandler)
	}
	app := newMessagesApp(t, upstream)
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusBadGateway, errorBody{Message: "Bad Gateway", Code: 502})
}

func TestSendMessageMapsUnparseableUpstreamBodyToBadGateway(t *testing.T) {
	t.Parallel()

	app := newMessagesApp(t, upstreamWith(http.StatusInternalServerError, `<html>proxy error</html>`))
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusBadGateway, errorBody{Message: "Bad Gateway", Code: 502})
}

func TestSendMessageMapsInvalidSuccessBodyToBadGateway(t *testing.T) {
	t.Parallel()

	app := newMessagesApp(t, upstreamWith(http.StatusOK, `<html>proxy error</html>`))
	resp := postMessage(t, app)

	assertErrorBody(t, resp, http.StatusBadGateway, errorBody{Message: "Bad Gateway", Code: 502})
}

func TestErrorMappingPreservesDeadlineExceeded(t *testing.T) {
	t.Parallel()

	app := fiber.New(
		fiber.Config{ErrorHandler: fiberfx.NewCustomJSONErrorHandler(zap.NewNop(), handlers.ErrorFormatter)},
	)
	app.Use(handlers.ErrorMapping)
	app.Get("/", func(_ *fiber.Ctx) error {
		return fmt.Errorf("wrapped: %w", &url.Error{
			Op:  "Post",
			URL: "http://upstream.local",
			Err: context.DeadlineExceeded,
		})
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })

	assertErrorBody(t, resp, http.StatusRequestTimeout, errorBody{Message: "Request Timeout", Code: 408})
}

func TestSendMessageSucceedsAgainstUpstream(t *testing.T) {
	t.Parallel()

	upstream := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/devices" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(deviceList))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprint(w, `{"id":"msg-1","state":"Pending"}`)
	}
	app := newMessagesApp(t, upstream)
	resp := postMessage(t, app)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var state struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if state.ID != "msg-1" {
		t.Errorf("id = %q, want %q", state.ID, "msg-1")
	}
}
