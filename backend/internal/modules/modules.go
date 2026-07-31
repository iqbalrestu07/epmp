package modules

import (
	"os"

	"github.com/epmp/backend/internal/modules/iam"
	"github.com/epmp/backend/internal/modules/property"
	"github.com/epmp/backend/internal/modules/room"
	"github.com/epmp/backend/internal/modules/tenant"
	mw "github.com/epmp/backend/internal/pkg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Register registers all modules to the Echo router.
func Register(e *echo.Echo, db *pgxpool.Pool, log zerolog.Logger) {
	jwtSecret := jwtSecretFromEnv()

	v1 := e.Group("/api/v1")

	// IAM module registers its own public and protected routes.
	iamMod := iam.NewModule(db, log, jwtSecret)
	iamMod.RegisterRoutes(v1)

	// Resource modules — all protected by JWT auth middleware.
	protected := v1.Group("", mw.AuthRequired(jwtSecret))
	property.NewModule(db, log).RegisterRoutes(protected)
	tenant.NewModule(db, log).RegisterRoutes(protected)
	room.NewModule(db, log).RegisterRoutes(protected)
}

func jwtSecretFromEnv() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "epmp-dev-secret-change-in-production"
}
