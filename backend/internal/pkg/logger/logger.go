package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// New creates a zerolog.Logger configured for the current environment.
// In production (APP_ENV=production) it emits structured JSON.
// Otherwise it emits human-readable console output with colour.
func New() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	env := os.Getenv("APP_ENV")
	if env == "production" {
		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	}

	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Logger()
}
