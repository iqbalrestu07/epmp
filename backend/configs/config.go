package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the centralized application configurations.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
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

	return cfg, nil
}
