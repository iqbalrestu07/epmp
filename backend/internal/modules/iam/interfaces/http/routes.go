package http

import (
	mw "github.com/epmp/backend/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

// RegisterIAMRoutes registers all IAM routes under the given API group.
// Auth routes are public; user/role routes require authentication.
func RegisterIAMRoutes(v1 *echo.Group, authH *AuthHandler, userH *UserHandler, roleH *RoleHandler, jwtSecret string) {
	// ── Public auth routes ──────────────────────────────────────────────────
	auth := v1.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)
	auth.POST("/refresh", authH.Refresh)
	auth.POST("/logout", authH.Logout)

	// ── Authenticated routes ────────────────────────────────────────────────
	protected := v1.Group("", mw.AuthRequired(jwtSecret))

	// Current user profile
	protected.GET("/auth/me", authH.Me)

	// Users (requires role management permission)
	users := protected.Group("/users", mw.RequirePermission("user:read"))
	users.GET("", userH.List)
	users.GET("/:id", userH.GetByID)
	users.POST("", userH.Create, mw.RequirePermission("user:write"))
	users.PUT("/:id", userH.Update, mw.RequirePermission("user:write"))
	users.DELETE("/:id", userH.Delete, mw.RequirePermission("user:delete"))
	users.PUT("/:id/roles", userH.AssignRoles, mw.RequirePermission("role:write"))

	// Roles
	roles := protected.Group("/roles", mw.RequirePermission("role:read"))
	roles.GET("", roleH.List)
	roles.GET("/:id", roleH.GetByID)
	roles.POST("", roleH.Create, mw.RequirePermission("role:write"))
	roles.PUT("/:id", roleH.Update, mw.RequirePermission("role:write"))
	roles.DELETE("/:id", roleH.Delete, mw.RequirePermission("role:delete"))
	roles.PUT("/:id/permissions", roleH.SetPermissions, mw.RequirePermission("role:write"))

	// Permissions (read-only catalogue)
	protected.GET("/permissions", roleH.ListPermissions)
}
