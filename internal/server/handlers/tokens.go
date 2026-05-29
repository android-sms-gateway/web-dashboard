package handlers

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/web-dashboard/internal/gateway"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/client"
	"github.com/android-sms-gateway/web-dashboard/internal/server/middlewares/session"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TokensHandler struct {
	handler.Base

	gatewaySvc *gateway.Factory
	logger     *zap.Logger
}

func NewTokensHandler(
	gatewaySvc *gateway.Factory,
	validator *validator.Validate,
	logger *zap.Logger,
) handler.Handler {
	return &TokensHandler{
		Base: handler.Base{
			Validator: validator,
		},
		gatewaySvc: gatewaySvc,
		logger:     logger,
	}
}

func (h *TokensHandler) Register(r fiber.Router) {
	g := r.Group("/tokens", session.AuthRequired(), client.New(h.gatewaySvc))

	g.Post("", validation.DecorateWithBodyEx(h.Validator, h.create))
	g.Delete(":jti", h.revoke)
}

type createTokenRequest struct {
	TTL    *uint64  `json:"ttl,omitempty"`
	Scopes []string `json:"scopes"        validate:"required,min=1,dive,required"`
}

// CreateToken generates a new API token.
//
//	@Summary		Create token
//	@Description	Generates a new API token with specified scopes and optional TTL.
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createTokenRequest	true	"Token request"
//	@Success		200		{object}	smsgateway.TokenResponse
//	@Failure		400		{object}	fiberfx.ErrorResponse
//	@Failure		401		{object}	fiberfx.ErrorResponse
//	@Router			/tokens [post]
func (h *TokensHandler) create(c *fiber.Ctx, req *createTokenRequest) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	scopes := make([]smsgateway.JWTScope, len(req.Scopes))
	copy(scopes, req.Scopes)

	tokenReq := smsgateway.TokenRequest{
		Scopes: scopes,
		TTL:    0,
	}
	if req.TTL != nil {
		tokenReq.TTL = *req.TTL
	}

	result, err := sdkClient.GenerateToken(c.Context(), tokenReq)
	if err != nil {
		h.logger.Warn("failed to create token", zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to create token")
	}

	return c.JSON(result)
}

// RevokeToken invalidates an API token by JTI.
//
//	@Summary		Revoke token
//	@Description	Invalidates an API token by its JTI (token ID).
//	@Tags			tokens
//	@Produce		json
//	@Param			jti	path	string	true	"Token JTI"
//	@Success		204
//	@Failure		401	{object}	fiberfx.ErrorResponse
//	@Failure		404	{object}	fiberfx.ErrorResponse
//	@Router			/tokens/{jti} [delete]
func (h *TokensHandler) revoke(c *fiber.Ctx) error {
	sdkClient := client.Get(c)
	if sdkClient == nil {
		h.logger.Warn("failed to get client")
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get client")
	}

	jti := c.Params("jti")
	if jti == "" {
		return fiber.NewError(fiber.StatusBadRequest, "token JTI is required")
	}

	if err := sdkClient.RevokeToken(c.Context(), jti); err != nil {
		h.logger.Warn("failed to revoke token", zap.String("jti", jti), zap.Error(err))
		return fiber.NewError(fiber.StatusBadGateway, "failed to revoke token")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
