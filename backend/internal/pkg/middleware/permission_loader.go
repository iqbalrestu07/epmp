package middleware

import (
	"github.com/epmp/backend/internal/modules/iam/repository"
	"github.com/labstack/echo/v4"
)

// PermissionLoader returns a middleware that loads the authenticated user's
// permissions from the database and stores them in the echo context.
//
// This must be placed AFTER AuthRequired in the middleware chain so that
// GetUserID(c) returns a valid user ID.
func PermissionLoader(userRoleRepo repository.UserRoleRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := GetUserID(c)
			if userID == "" {
				return next(c) // AuthRequired will handle the 401
			}

			perms, err := userRoleRepo.FindPermissionsByUserID(c.Request().Context(), userID)
			if err == nil {
				keys := make([]string, 0, len(perms))
				for _, p := range perms {
					keys = append(keys, p.Key())
				}
				SetPermissions(c, keys)
			}

			return next(c)
		}
	}
}
