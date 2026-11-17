package client

import (
	"fmt"

	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/gofiber/fiber/v2"
)

type localsKey string

const localsKeyClient = localsKey("client")

func New(factory *gateway.Factory) fiber.Handler {
	return func(c *fiber.Ctx) error {
		login, password, err := session.GetCredentials(c)
		if err != nil {
			return fmt.Errorf("failed to get credentials: %w", err)
		}

		client := factory.NewClient(login, password)
		c.Locals(localsKeyClient, client)

		return c.Next()
	}
}

func Get(c *fiber.Ctx) *gateway.Client {
	client, ok := c.Locals(localsKeyClient).(*gateway.Client)
	if !ok {
		return nil
	}

	return client
}
