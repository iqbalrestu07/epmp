package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect connects to the PostgreSQL database.
func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dbURL)
}
