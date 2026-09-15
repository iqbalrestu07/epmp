package payment

import (
	"github.com/epmp/backend/internal/modules/payment/delivery/http"
	"github.com/epmp/backend/internal/modules/payment/repository"
	"github.com/epmp/backend/internal/modules/payment/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all Payment dependencies.
type Module struct {
	Handler *http.PaymentHandler
}

// NewModule creates and wires all Payment dependencies.
func NewModule(db *pgxpool.Pool, log zerolog.Logger, notifier service.OrgNotifier) *Module {
	var repo repository.PaymentRepository = repository.NewPaymentRepositoryImpl(db)
	svc := service.NewPaymentService(repo, notifier)
	handler := http.NewPaymentHandler(svc)

	log.Info().Str("module", "payment").Msg("module initialized")

	return &Module{
		Handler: handler,
	}
}

// RegisterRoutes registers all Payment routes on the given Echo instance.
func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/payments")
	http.RegisterPaymentRoutes(g, m.Handler)
}
