package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// runServer orchestrates the server startup and graceful shutdown.
func runServer() error {
	ctx := context.Background()
	app, cleanup, err := bootstrap(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		app.Log.Info().Str("port", port).Msg("server starting")
		if err := app.Engine.Start(":" + port); err != nil && err != http.ErrServerClosed {
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
