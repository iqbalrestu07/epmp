package main

import (
	"context"
	"net/http"

	"github.com/epmp/backend/internal/database/postgres"
	"github.com/epmp/backend/internal/modules"
	"github.com/epmp/backend/internal/pkg/logger"
	mw "github.com/epmp/backend/internal/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// App holds the initialized core dependencies for the application.
type App struct {
	Engine *echo.Echo
	Log    zerolog.Logger
}

// bootstrap initializes configuration, database, logger, middleware, and registers all modules.
func bootstrap(ctx context.Context) (*App, func(), error) {
	log := logger.New()

	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return nil, nil, err
	}
	cleanup := func() {
		db.Close()
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = errorHandler(log)

	// Global middleware
	e.Use(mw.RequestLogger(log))
	e.Use(mw.Recover(log))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Init modules
	modules.Register(e, db, log)

	return &App{
		Engine: e,
		Log:    log,
	}, cleanup, nil
}

// errorHandler is the centralised Echo error handler.
// It maps *echo.HTTPError and domain errors to a consistent JSON envelope.
func errorHandler(log zerolog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "internal server error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if s, ok := he.Message.(string); ok {
				msg = s
			}
		}

		log.Error().
			Err(err).
			Int("status", code).
			Str("method", c.Request().Method).
			Str("path", c.Request().URL.Path).
			Msg("request error")

		// Do not overwrite a response that has already started.
		if !c.Response().Committed {
			_ = c.JSON(code, map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    http.StatusText(code),
					"message": msg,
				},
			})
		}
	}
}
