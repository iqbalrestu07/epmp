package report

import (
	"github.com/epmp/backend/internal/modules/report/delivery/http"
	"github.com/epmp/backend/internal/modules/report/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all Report dependencies.
type Module struct {
	Handler *http.ReportHandler
}

// NewModule creates and wires all Report dependencies.
func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	svc := service.NewReportService(db)
	handler := http.NewReportHandler(svc)

	log.Info().Str("module", "report").Msg("module initialized")

	return &Module{
		Handler: handler,
	}
}

// RegisterRoutes registers all Report routes on the given Echo instance.
func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/reports")
	http.RegisterReportRoutes(g, m.Handler)
}
