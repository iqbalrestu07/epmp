package middleware

import (
	"net/http"
	"strings"

	iamservice "github.com/epmp/backend/internal/modules/iam/service"
	"github.com/labstack/echo/v4"
)

const (
	// ContextKeyUserID is the echo context key for the authenticated user's ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUserEmail is the echo context key for the authenticated user's email.
	ContextKeyUserEmail = "user_email"
	// ContextKeyPermissions is the echo context key for the user's permission set.
	ContextKeyPermissions = "permissions"
)

// AuthRequired is a JWT authentication middleware.
// It validates the Bearer token and injects user claims into the echo context.
// Routes marked with this middleware are inaccessible without a valid token.
func AuthRequired(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authorization header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := iamservice.ParseAccessToken(tokenStr, jwtSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyUserEmail, claims.Email)

			return next(c)
		}
	}
}

// RequirePermission returns a middleware that checks if the authenticated user
// has the given permission key (e.g. "property:read") stored in context.
//
// IMPORTANT: This middleware must run AFTER AuthRequired AND after the user's
// permissions have been loaded into the context. For the current JWT-only setup,
// permissions are checked lazily via the database (see PermissionLoader).
// In future, permissions can be embedded in the JWT claims for performance.
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			perms, ok := c.Get(ContextKeyPermissions).([]string)
			if !ok {
				// Permissions not loaded — deny access
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}
			for _, p := range perms {
				if p == permission {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

// GetUserID extracts the authenticated user's ID from the echo context.
func GetUserID(c echo.Context) string {
	id, _ := c.Get(ContextKeyUserID).(string)
	return id
}

// GetUserEmail extracts the authenticated user's email from the echo context.
func GetUserEmail(c echo.Context) string {
	email, _ := c.Get(ContextKeyUserEmail).(string)
	return email
}

// SetPermissions injects a user's permissions into the echo context.
// Call this from a permission-loading middleware before route handlers.
func SetPermissions(c echo.Context, permissions []string) {
	c.Set(ContextKeyPermissions, permissions)
}
