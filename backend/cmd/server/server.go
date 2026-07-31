package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epmp/backend/configs"
)

// runServer orchestrates the server startup and graceful shutdown.
func runServer() error {
	ctx := context.Background()

	// Load centralized configuration
	cfg, err := configs.Load()
	if err != nil {
		return err
	}

	app, cleanup, err := bootstrap(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	go func() {
		app.Log.Info().Str("port", cfg.Port).Msg("server starting")
		if err := app.Engine.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			app.Log.Fatal().Err(err).Msg("server stopped unexpectedly")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Log.Info().Msg("shutting down server…")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Engine.Shutdown(shutdownCtx); err != nil {
		app.Log.Error().Err(err).Msg("graceful shutdown failed")
		return err
	}

	app.Log.Info().Msg("server stopped")
	return nil
}
