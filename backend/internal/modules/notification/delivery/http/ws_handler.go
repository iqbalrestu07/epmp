package http

import (
	"net/http"

	iamservice "github.com/epmp/backend/internal/modules/iam/service"
	ws "github.com/epmp/backend/internal/pkg/websocket"

	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

// WSHandler upgrades authenticated connections to websockets for live notifications.
// Browsers cannot set Authorization headers on WebSocket handshakes, so the JWT
// is passed as a query parameter: /api/v1/ws?token=<jwt>&org_id=<org>.
type WSHandler struct {
	hub       *ws.Hub
	jwtSecret string
}

// NewWSHandler creates a new WSHandler.
func NewWSHandler(hub *ws.Hub, jwtSecret string) *WSHandler {
	return &WSHandler{hub: hub, jwtSecret: jwtSecret}
}

func (h *WSHandler) Connect(c echo.Context) error {
	tokenStr := c.QueryParam("token")
	orgID := c.QueryParam("org_id")
	if tokenStr == "" || orgID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "token and org_id are required")
	}

	claims, err := iamservice.ParseAccessToken(tokenStr, h.jwtSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
	}

	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Origin check is handled by CORS/reverse proxy.
	})
	if err != nil {
		return err
	}

	client := h.hub.Register(claims.UserID, orgID, conn)
	defer func() {
		h.hub.Unregister(client)
		_ = conn.Close(websocket.StatusNormalClosure, "bye")
	}()

	// Read pump: client messages are unused; reading detects disconnects
	// and lets coder/websocket answer control frames (ping/pong).
	ctx := c.Request().Context()
	for {
		if _, _, err := conn.Read(ctx); err != nil {
			return nil
		}
	}
}
