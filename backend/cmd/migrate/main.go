// cmd/migrate/main.go
//
// EPMP Migration Runner
//
// Modes:
//   manual  — run a specific migrate command (up, down, version, force, create)
//   watch   — watch the migrations/ directory and apply new .up.sql files automatically
//
// Usage:
//
//	# Apply all pending migrations
//	go run ./cmd/migrate up
//
//	# Rollback the last migration
//	go run ./cmd/migrate down 1
//
//	# Show current schema version
//	go run ./cmd/migrate version
//
//	# Force-set a version (fixes dirty state)
//	go run ./cmd/migrate force <version>
//
//	# Create a new migration pair
//	go run ./cmd/migrate create <name>
//
//	# Watch migrations/ and auto-apply when new files are detected
//	go run ./cmd/migrate watch

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsDir = "migrations"

func main() {
	// ── Logger (always text; migration runner is a dev/ops tool) ──────────────
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	dbURL := databaseURL()
	migrationsPath := resolveMigrationsDir()

	cmd := os.Args[1]

	switch cmd {
	case "up":
		runUp(dbURL, migrationsPath)
	case "down":
		steps := 1
		if len(os.Args) >= 3 {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil {
				slog.Error("down requires a numeric steps argument", "arg", os.Args[2])
				os.Exit(1)
			}
			steps = n
		}
		runDown(dbURL, migrationsPath, steps)
	case "version":
		runVersion(dbURL, migrationsPath)
	case "force":
		if len(os.Args) < 3 {
			slog.Error("force requires a version number")
			os.Exit(1)
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			slog.Error("force version must be a number", "arg", os.Args[2])
			os.Exit(1)
		}
		runForce(dbURL, migrationsPath, v)
	case "create":
		if len(os.Args) < 3 {
			slog.Error("create requires a migration name")
			os.Exit(1)
		}
		name := strings.Join(os.Args[2:], "_")
		runCreate(migrationsPath, name)
	case "watch":
		runWatch(dbURL, migrationsPath)
	default:
		slog.Error("unknown command", "cmd", cmd)
		printUsage()
		os.Exit(1)
	}
}

// ── Commands ──────────────────────────────────────────────────────────────────

func runUp(dbURL, migrationsPath string) {
	m := mustNewMigrate(dbURL, migrationsPath)
	defer m.Close()

	slog.Info("applying all pending migrations…")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("no pending migrations — schema is up to date")
			return
		}
		slog.Error("migration up failed", "err", err)
		os.Exit(1)
	}

	v, _, _ := m.Version()
	slog.Info("migrations applied", "current_version", v)
}

func runDown(dbURL, migrationsPath string, steps int) {
	m := mustNewMigrate(dbURL, migrationsPath)
	defer m.Close()

	slog.Info("rolling back migrations", "steps", steps)
	if err := m.Steps(-steps); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("nothing to roll back")
			return
		}
		slog.Error("migration down failed", "err", err)
		os.Exit(1)
	}

	v, _, _ := m.Version()
	slog.Info("rollback complete", "current_version", v)
}

func runVersion(dbURL, migrationsPath string) {
	m := mustNewMigrate(dbURL, migrationsPath)
	defer m.Close()

	v, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			slog.Info("no migrations applied yet (version 0)")
			return
		}
		slog.Error("failed to get version", "err", err)
		os.Exit(1)
	}
	slog.Info("current schema version", "version", v, "dirty", dirty)
}

func runForce(dbURL, migrationsPath string, version int) {
	m := mustNewMigrate(dbURL, migrationsPath)
	defer m.Close()

	slog.Warn("forcing schema version — use only to fix dirty state", "version", version)
	if err := m.Force(version); err != nil {
		slog.Error("force failed", "err", err)
		os.Exit(1)
	}
	slog.Info("version forced", "version", version)
}

// runCreate scaffolds a sequentially numbered migration pair (up + down).
func runCreate(migrationsPath, name string) {
	// Find next sequence number by scanning existing files.
	seq := nextSequence(migrationsPath)
	base := fmt.Sprintf("%06d_%s", seq, name)

	upFile := filepath.Join(migrationsPath, base+".up.sql")
	downFile := filepath.Join(migrationsPath, base+".down.sql")

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created:   %s\n\n-- TODO: write your UP migration here\n",
		base, time.Now().Format("2006-01-02"))
	downContent := fmt.Sprintf("-- Rollback: %s\n\n-- TODO: write your DOWN migration here\n", base)

	writeFile(upFile, upContent)
	writeFile(downFile, downContent)

	slog.Info("migration files created", "up", upFile, "down", downFile)
}

// runWatch monitors the migrations directory for new .up.sql files and applies
// them automatically. It compares against the current schema version so it
// only applies genuinely new files.
func runWatch(dbURL, migrationsPath string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("watching migrations directory for new files", "dir", migrationsPath)
	slog.Info("press Ctrl+C to stop")

	// Apply any already-pending migrations on startup.
	applyPending(dbURL, migrationsPath)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	var lastKnownCount int
	lastKnownCount = countUpFiles(migrationsPath)

	for {
		select {
		case <-ctx.Done():
			slog.Info("watcher stopped")
			return
		case <-ticker.C:
			current := countUpFiles(migrationsPath)
			if current > lastKnownCount {
				slog.Info("new migration file(s) detected", "previous", lastKnownCount, "current", current)
				applyPending(dbURL, migrationsPath)
				lastKnownCount = current
			}
		}
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func applyPending(dbURL, migrationsPath string) {
	m := mustNewMigrate(dbURL, migrationsPath)
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Debug("no new migrations to apply")
			return
		}
		slog.Error("auto-apply failed", "err", err)
		return
	}

	v, _, _ := m.Version()
	slog.Info("auto-applied pending migrations", "current_version", v)
}

func mustNewMigrate(dbURL, migrationsPath string) *migrate.Migrate {
	sourceURL := "file://" + migrationsPath
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		slog.Error("failed to initialise migrate", "err", err)
		os.Exit(1)
	}
	m.Log = &migrateLogger{}
	return m
}

func databaseURL() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	return "postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable"
}

// resolveMigrationsDir returns an absolute path to migrations/ relative to
// the binary's working directory (i.e. the backend/ root).
func resolveMigrationsDir() string {
	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		slog.Error("cannot resolve migrations dir", "err", err)
		os.Exit(1)
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		slog.Error("migrations directory does not exist", "path", abs)
		os.Exit(1)
	}
	return abs
}

// nextSequence returns (max existing sequence number + 1).
func nextSequence(dir string) int {
	entries, _ := os.ReadDir(dir)
	max := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		parts := strings.SplitN(e.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max + 1
}

// countUpFiles returns the number of .up.sql files in dir, sorted so we can
// detect additions.
func countUpFiles(dir string) int {
	entries, _ := os.ReadDir(dir)
	count := 0
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	count = len(names)
	return count
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		slog.Error("failed to write file", "path", path, "err", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
EPMP Migration Runner

Usage:
  go run ./cmd/migrate <command> [args]

Commands:
  up               Apply all pending migrations
  down [N]         Roll back N steps (default: 1)
  version          Print current schema version
  force <version>  Force-set version to fix dirty state
  create <name>    Scaffold a new numbered migration pair
  watch            Auto-apply new migration files as they appear

Environment:
  DATABASE_URL     PostgreSQL connection string
                   Default: postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable

Examples:
  go run ./cmd/migrate up
  go run ./cmd/migrate down 2
  go run ./cmd/migrate create add_column_to_tenants
  go run ./cmd/migrate watch
`)
}

// migrateLogger bridges golang-migrate's logger to slog.
type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(format, v...))
}
func (l *migrateLogger) Verbose() bool { return true }
