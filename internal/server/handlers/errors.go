package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/android-sms-gateway/client-go/rest"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/go-core-fx/fiberfx"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

const upstreamMessageLimit = 300

func ErrorFormatter(err error, code int) any {
	switch {
	case errors.Is(err, context.Canceled):
		return fiberfx.NewErrorResponse("client closed request", code, nil)
	case errors.Is(err, context.DeadlineExceeded):
		return fiberfx.NewErrorResponse(statusMessage(code), code, nil)
	}

	if code >= fiber.StatusInternalServerError {
		if fiberErr, ok := errors.AsType[*fiber.Error](err); ok && fiberErr.Message != "" {
			return fiberfx.NewErrorResponse(fiberErr.Message, code, nil)
		}

		return fiberfx.NewErrorResponse(statusMessage(code), code, nil)
	}

	return fiberfx.NewErrorResponse(err.Error(), code, nil)
}

func ErrorMapping(c *fiber.Ctx) error {
	err := c.Next()
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if apiErr, ok := errors.AsType[*rest.APIError](err); ok {
		return mapAPIError(apiErr)
	}

	if _, ok := errors.AsType[*url.Error](err); ok {
		return fiber.NewError(fiber.StatusBadGateway)
	}

	if errors.Is(err, gateway.ErrUpstream) {
		return fiber.NewError(fiber.StatusBadGateway)
	}

	return err
}

func mapAPIError(apiErr *rest.APIError) error {
	status := apiErr.StatusCode
	if !propagateStatus(status) {
		status = fiber.StatusBadGateway
	}

	return fiber.NewError(status, messageOr(upstreamMessage(apiErr), status))
}

func propagateStatus(statusCode int) bool {
	return statusCode < fiber.StatusInternalServerError || statusCode == http.StatusServiceUnavailable
}

func upstreamMessage(apiErr *rest.APIError) string {
	var body smsgateway.ErrorResponse
	if json.Unmarshal(apiErr.Body, &body) != nil {
		return ""
	}

	return lo.Substring(strings.TrimSpace(body.Message), 0, upstreamMessageLimit)
}

func messageOr(msg string, statusCode int) string {
	if msg != "" {
		return msg
	}

	if text := http.StatusText(statusCode); text != "" {
		return text
	}

	return "upstream error"
}

func statusMessage(code int) string {
	if text := http.StatusText(code); text != "" {
		return text
	}

	return "internal server error"
}
