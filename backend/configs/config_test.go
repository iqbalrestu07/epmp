package configs

import (
	"os"
	"testing"
	"time"
)

func setRequiredEnvVars(t *testing.T) {
	t.Helper()
	t.Setenv("PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
}

func TestLoad_DefaultTTLs(t *testing.T) {
	setRequiredEnvVars(t)
	// Ensure TTL env vars are NOT set — should use defaults.
	os.Unsetenv("ACCESS_TOKEN_TTL_MINUTES")
	os.Unsetenv("REFRESH_TOKEN_TTL_DAYS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	expectedAccess := 120 * time.Minute
	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected default AccessTokenTTL %v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}

	expectedRefresh := 30 * 24 * time.Hour
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected default RefreshTokenTTL %v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}

func TestLoad_CustomTTLs(t *testing.T) {
	setRequiredEnvVars(t)
	t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "30")
	t.Setenv("REFRESH_TOKEN_TTL_DAYS", "7")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	expectedAccess := 30 * time.Minute
	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected AccessTokenTTL %v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}

	expectedRefresh := 7 * 24 * time.Hour
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected RefreshTokenTTL %v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}

func TestLoad_InvalidTTLs_FallbackToDefaults(t *testing.T) {
	setRequiredEnvVars(t)
	t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "not-a-number")
	t.Setenv("REFRESH_TOKEN_TTL_DAYS", "abc")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should fall back to defaults when env values are invalid.
	expectedAccess := 120 * time.Minute
	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected default AccessTokenTTL %v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}

	expectedRefresh := 30 * 24 * time.Hour
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected default RefreshTokenTTL %v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}

func TestLoad_NegativeTTLs_FallbackToDefaults(t *testing.T) {
	setRequiredEnvVars(t)
	t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "-10")
	t.Setenv("REFRESH_TOKEN_TTL_DAYS", "-5")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should fall back to defaults when env values are <= 0.
	expectedAccess := 120 * time.Minute
	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected default AccessTokenTTL %v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}

	expectedRefresh := 30 * 24 * time.Hour
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected default RefreshTokenTTL %v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}

func TestLoad_ZeroTTLs_FallbackToDefaults(t *testing.T) {
	setRequiredEnvVars(t)
	t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "0")
	t.Setenv("REFRESH_TOKEN_TTL_DAYS", "0")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Zero values should also fall back to defaults.
	expectedAccess := 120 * time.Minute
	if cfg.AccessTokenTTL != expectedAccess {
		t.Errorf("expected default AccessTokenTTL %v, got %v", expectedAccess, cfg.AccessTokenTTL)
	}

	expectedRefresh := 30 * 24 * time.Hour
	if cfg.RefreshTokenTTL != expectedRefresh {
		t.Errorf("expected default RefreshTokenTTL %v, got %v", expectedRefresh, cfg.RefreshTokenTTL)
	}
}
