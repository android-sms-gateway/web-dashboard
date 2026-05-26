package handlers

import (
	"errors"

	"github.com/android-sms-gateway/client-go/rest"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewAuthHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &AuthHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *AuthHandler) Register(r fiber.Router) {
	g := r.Group("/auth")

	g.Post("/login", validation.DecorateWithBodyEx(h.Validator, h.login))
	g.Post("/logout", h.logout)
	g.Get("/me", session.AuthRequired(), h.me)
}

type loginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Login string `json:"login"`
}

// Login authenticates with the SMSGate and creates a session.
//
//	@Summary		Login
//	@Description	Authenticate with the SMSGate using basic credentials. Returns a session cookie.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginRequest	true	"Login credentials"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) login(c *fiber.Ctx, req *loginRequest) error {
	client := h.gatewaySvc.NewClient(req.Login, req.Password)
	if _, err := client.ListDevices(c.Context()); err != nil {
		if errors.Is(err, gateway.ErrUnauthorized) || rest.IsClientError(err) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		h.logger.Warn("login failed", zap.String("login", req.Login), zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "gateway unavailable")
	}

	if _, err := session.SetCredentials(c, req.Login, req.Password); err != nil {
		h.logger.Warn("failed to set credentials", zap.String("login", req.Login), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to set credentials")
	}

	return c.JSON(loginResponse{
		Login: req.Login,
	})
}

// Logout invalidates the current session.
//
//	@Summary		Logout
//	@Description	Invalidate the current session and clear the session cookie.
//	@Tags			auth
//	@Produce		json
//	@Success		204
//	@Router			/auth/logout [post]
func (h *AuthHandler) logout(c *fiber.Ctx) error {
	sess := session.Get(c)
	if sess == nil {
		h.logger.Warn("failed to get session")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get session")
	}

	if err := sess.Destroy(); err != nil {
		h.logger.Warn("failed to destroy session", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to destroy session")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type meResponse struct {
	Login string `json:"login"`
}

// Me returns the current authenticated user's login.
//
//	@Summary		Current user
//	@Description	Returns the login of the currently authenticated user.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	meResponse
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Router			/auth/me [get]
func (h *AuthHandler) me(c *fiber.Ctx) error {
	login, _, err := session.GetCredentials(c)
	if err != nil {
		h.logger.Warn("failed to get credentials", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get credentials")
	}

	return c.JSON(meResponse{Login: login})
}
