package audit

import (
	"github.com/epmp/backend/internal/modules/audit/delivery/http"
	"github.com/epmp/backend/internal/modules/audit/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all Audit dependencies.
type Module struct {
	Handler *http.AuditHandler
}

// NewModule creates and wires all Audit dependencies.
func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	svc := service.NewAuditService(db)
	handler := http.NewAuditHandler(svc)

	log.Info().Str("module", "audit").Msg("module initialized")

	return &Module{
		Handler: handler,
	}
}

// RegisterRoutes registers all Audit routes on the given Echo instance.
func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/audit-logs")
	http.RegisterAuditRoutes(g, m.Handler)
}
