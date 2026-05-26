package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	sessionKeyLogin    = "login"
	sessionKeyPassword = "password"
)

func SetCredentials(c *fiber.Ctx, login, password string) (*session.Session, error) {
	sess := Get(c)
	if sess == nil {
		return nil, ErrSessionNotFound
	}

	sess.Set(sessionKeyLogin, login)
	sess.Set(sessionKeyPassword, password)

	return sess, nil
}

func GetCredentials(c *fiber.Ctx) (string, string, error) {
	sess := Get(c)
	if sess == nil {
		return "", "", ErrSessionNotFound
	}

	login, ok := sess.Get(sessionKeyLogin).(string)
	if !ok {
		return "", "", ErrSessionNotFound
	}

	password, ok := sess.Get(sessionKeyPassword).(string)
	if !ok {
		return "", "", ErrSessionNotFound
	}

	return login, password, nil
}

func IsAuthenticated(c *fiber.Ctx) bool {
	_, _, err := GetCredentials(c)
	return err == nil
}

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !IsAuthenticated(c) {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}
}
