package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect connects to the PostgreSQL database.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, getURL())
}

func getURL() string {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		return u
	}
	return "postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable"
}
