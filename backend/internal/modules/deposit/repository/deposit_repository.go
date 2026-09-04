package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/deposit/entity"
)

// DepositRepository is the contract for persisting Deposit aggregates.
type DepositRepository interface {
	// Save persists a Deposit entity.
	Save(ctx context.Context, e *entity.Deposit) error

	// FindByID retrieves a Deposit by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Deposit, error)

	// FindAll retrieves a paginated list of Deposit entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Deposit, error)

	// Count returns the total number of Deposit entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Deposit by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
