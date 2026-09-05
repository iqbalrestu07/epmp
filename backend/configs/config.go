package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all the centralized application configurations.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load reads the .env file (if available) and builds the Config struct from environment variables.
func Load() (*Config, error) {
	// Attempt to load .env file; ignore if it doesn't exist (e.g. in production).
	_ = godotenv.Load()

	cfg := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	// Validate required variables
	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT environment variable is not set")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Token TTL — defaults: access 2h, refresh 30d
	accessMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))
	if err != nil || accessMinutes <= 0 {
		accessMinutes = 120 // 2 hours
	}
	cfg.AccessTokenTTL = time.Duration(accessMinutes) * time.Minute

	refreshDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL_DAYS"))
	if err != nil || refreshDays <= 0 {
		refreshDays = 30 // 30 days
	}
	cfg.RefreshTokenTTL = time.Duration(refreshDays) * 24 * time.Hour

	return cfg, nil
}
