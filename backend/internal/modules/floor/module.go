package floor

import (
	"github.com/epmp/backend/internal/modules/building/repository"
	"github.com/epmp/backend/internal/modules/floor/delivery/http"
	floorrepo "github.com/epmp/backend/internal/modules/floor/repository"
	"github.com/epmp/backend/internal/modules/floor/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all Floor dependencies.
type Module struct {
	Handler *http.FloorHandler
}

// NewModule creates and wires all Floor dependencies.
func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	var repo floorrepo.FloorRepository = floorrepo.NewFloorRepositoryImpl(db)
	buildingRepo := repository.NewBuildingRepositoryImpl(db)
	svc := service.NewFloorService(repo, buildingRepo)
	handler := http.NewFloorHandler(svc)

	log.Info().Str("module", "floor").Msg("module initialized")

	return &Module{
		Handler: handler,
	}
}

// RegisterRoutes registers all Floor routes on the given Echo instance.
func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/floors")
	http.RegisterFloorRoutes(g, m.Handler)
}
