package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	property "github.com/epmp/backend/internal/modules/property"
	room "github.com/epmp/backend/internal/modules/room"
	tenant "github.com/epmp/backend/internal/modules/tenant"
	"github.com/epmp/backend/internal/shared/logger"
	mw "github.com/epmp/backend/internal/shared/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.New()

	db, err := pgxpool.New(context.Background(), dbURL())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = errorHandler(log)

	// --- Global middleware ---
	e.Use(mw.RequestLogger(log))
	e.Use(mw.Recover(log))

	// --- Health check (unauthenticated) ---
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// --- API v1 group ---
	v1 := e.Group("/api/v1")

	property.NewModule(db, log).RegisterRoutes(v1)
	tenant.NewModule(db, log).RegisterRoutes(v1)
	room.NewModule(db, log).RegisterRoutes(v1)

	// --- Graceful shutdown ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Info().Str("port", port).Msg("server starting")
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server stopped unexpectedly")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("server stopped")
}

// dbURL returns the PostgreSQL connection string from the environment.
func dbURL() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	return "postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable"
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
