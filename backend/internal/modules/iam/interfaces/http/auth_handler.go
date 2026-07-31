package http

import (
	"net/http"

	"github.com/epmp/backend/internal/modules/iam/application/dto"
	"github.com/epmp/backend/internal/modules/iam/application/service"
	"github.com/epmp/backend/internal/shared"
	"github.com/epmp/backend/internal/shared/response"
	mw "github.com/epmp/backend/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register godoc
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.authSvc.Register(c.Request().Context(), &req)
	if err != nil {
		if shared.IsDomainError(err, "ALREADY_EXISTS") {
			return response.Conflict(c, "email already registered")
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"success": true, "data": res})
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.authSvc.Login(c.Request().Context(), &req)
	if err != nil {
		if shared.IsDomainError(err, "INVALID_CREDENTIALS") || shared.IsDomainError(err, "NOT_FOUND") {
			return response.Unauthorized(c, "invalid email or password")
		}
		if shared.IsDomainError(err, "ACCOUNT_INACTIVE") {
			return response.Forbidden(c, "account is inactive")
		}
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

// Refresh godoc
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.authSvc.Refresh(c.Request().Context(), &req)
	if err != nil {
		return response.Unauthorized(c, "invalid or expired refresh token")
	}
	return response.OK(c, res)
}

// Logout godoc
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	_ = h.authSvc.Logout(c.Request().Context(), &req)
	return response.NoContent(c)
}

// Me godoc
// GET /api/v1/auth/me  (requires auth)
func (h *AuthHandler) Me(c echo.Context) error {
	userID := mw.GetUserID(c)
	if userID == "" {
		return response.Unauthorized(c, "unauthorized")
	}

	res, err := h.authSvc.Me(c.Request().Context(), userID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, res)
}
