package iam

import (
	iamhttp "github.com/epmp/backend/internal/modules/iam/delivery/http"
	"github.com/epmp/backend/internal/modules/iam/repository"
	iamrepo "github.com/epmp/backend/internal/modules/iam/repository"
	"github.com/epmp/backend/internal/modules/iam/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all IAM dependencies.
type Module struct {
	AuthHandler  *iamhttp.AuthHandler
	UserHandler  *iamhttp.UserHandler
	RoleHandler  *iamhttp.RoleHandler
	UserRoleRepo repository.UserRoleRepository // used by PermissionLoader
	JWTSecret    string
}

// NewModule creates and wires all IAM dependencies.
func NewModule(db *pgxpool.Pool, log zerolog.Logger, jwtSecret string) *Module {
	// Repositories
	userRepo := iamrepo.NewUserRepository(db)
	roleRepo := iamrepo.NewRoleRepository(db)
	permissionRepo := iamrepo.NewPermissionRepository(db)
	userRoleRepo := iamrepo.NewUserRoleRepository(db)
	refreshTokenRepo := iamrepo.NewRefreshTokenRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, userRoleRepo, refreshTokenRepo, roleRepo, jwtSecret)
	userSvc := service.NewUserService(userRepo, userRoleRepo)
	roleSvc := service.NewRoleService(roleRepo, permissionRepo)

	// Handlers
	authH := iamhttp.NewAuthHandler(authSvc)
	userH := iamhttp.NewUserHandler(userSvc)
	roleH := iamhttp.NewRoleHandler(roleSvc)

	log.Info().Str("module", "iam").Msg("module initialized")

	return &Module{
		AuthHandler:  authH,
		UserHandler:  userH,
		RoleHandler:  roleH,
		UserRoleRepo: userRoleRepo,
		JWTSecret:    jwtSecret,
	}
}

// RegisterRoutes registers all IAM routes under the given API group.
// The PermissionLoader middleware is applied to all authenticated routes
// so that per-route RBAC checks have the user's permissions available.
func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Attach permission loader to authenticated group
	v1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := mw.GetUserID(c)
			if userID != "" {
				// Only load permissions if the user is authenticated
				loader := mw.PermissionLoader(m.UserRoleRepo)
				return loader(next)(c)
			}
			return next(c)
		}
	})

	iamhttp.RegisterIAMRoutes(v1, m.AuthHandler, m.UserHandler, m.RoleHandler, m.JWTSecret)
}
