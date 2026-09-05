package dashboard

import (
	"github.com/epmp/backend/internal/modules/dashboard/delivery/http"
	"github.com/epmp/backend/internal/modules/dashboard/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Module struct {
	Handler *http.DashboardHandler
}

func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	svc := service.NewDashboardService(db)
	handler := http.NewDashboardHandler(svc)

	log.Info().Str("module", "dashboard").Msg("module initialized")

	return &Module{
		Handler: handler,
	}
}

func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/dashboard")
	http.RegisterDashboardRoutes(g, m.Handler)
}
