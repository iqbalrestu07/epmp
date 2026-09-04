package communication

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Module struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewModule(db *pgxpool.Pool, log zerolog.Logger) *Module {
	return &Module{
		db:  db,
		log: log,
	}
}

func (m *Module) RegisterRoutes(router *echo.Group) {
	group := router.Group("/communication")
	
	// Devices
	group.GET("/devices", func(c echo.Context) error { return c.JSON(200, []string{}) })
	group.POST("/devices/qr", func(c echo.Context) error { return c.JSON(200, map[string]string{"qr": "base64..."}) })
	
	// Messages
	group.POST("/blast", func(c echo.Context) error { return c.JSON(200, map[string]string{"status": "queued"}) })
}
