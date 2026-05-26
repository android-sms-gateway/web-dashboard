package session

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type localsKey string

const localsKeySession = localsKey("session")

func New(storage fiber.Storage) fiber.Handler {
	store := session.New(session.Config{
		Storage:           storage,
		KeyLookup:         "cookie:session_id",
		CookiePath:        "/",
		CookieHTTPOnly:    true,
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
	})

	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}

		c.Locals(localsKeySession, sess)

		if nextErr := c.Next(); nextErr != nil {
			return nextErr //nolint:wrapcheck // we don't want to wrap the error
		}

		if saveErr := sess.Save(); saveErr != nil {
			return fmt.Errorf("failed to save session: %w", saveErr)
		}

		return nil
	}
}

func Get(c *fiber.Ctx) *session.Session {
	sess, ok := c.Locals(localsKeySession).(*session.Session)
	if !ok {
		return nil
	}

	return sess
}
