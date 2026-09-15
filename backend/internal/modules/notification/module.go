package notification

import (
	"github.com/epmp/backend/internal/modules/notification/delivery/http"
	"github.com/epmp/backend/internal/modules/notification/service"
	"github.com/epmp/backend/internal/pkg/email"
	"github.com/epmp/backend/internal/pkg/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Module wires together all Notification dependencies.
type Module struct {
	Handler   *http.NotificationHandler
	WSHandler *http.WSHandler
	Service   *service.NotificationService
	Hub       *websocket.Hub
}

// NewModule creates and wires all Notification dependencies.
// emailer may be nil — the email channel is then disabled.
func NewModule(db *pgxpool.Pool, log zerolog.Logger, jwtSecret string, emailer email.Sender) *Module {
	hub := websocket.NewHub(log)
	svc := service.NewNotificationService(db, hub, emailer, log)
	handler := http.NewNotificationHandler(svc)
	wsHandler := http.NewWSHandler(hub, jwtSecret)

	log.Info().Str("module", "notification").Msg("module initialized")

	return &Module{
		Handler:   handler,
		WSHandler: wsHandler,
		Service:   svc,
		Hub:       hub,
	}
}

// RegisterRoutes registers the REST notification routes (JWT-protected).
func (m *Module) RegisterRoutes(e *echo.Group) {
	g := e.Group("/notifications")
	http.RegisterNotificationRoutes(g, m.Handler)
}

// RegisterWS registers the websocket endpoint on an unauthenticated group;
// the handler performs its own JWT check via query param.
func (m *Module) RegisterWS(e *echo.Group) {
	e.GET("/ws", m.WSHandler.Connect)
}
