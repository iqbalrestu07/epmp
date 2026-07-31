package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func RequestLogger(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			latency := time.Since(start)
			log.Info().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Int("status", c.Response().Status).
				Dur("latency", latency).
				Str("remote_ip", c.RealIP()).
				Msg("request")

			return err
		}
	}
}

func Recover(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Error().
						Interface("panic", r).
						Str("path", c.Request().URL.Path).
						Msg("panic recovered")
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, "internal server error"))
				}
			}()
			return next(c)
		}
	}
}
