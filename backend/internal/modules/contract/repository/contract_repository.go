package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/contract/entity"
)

// ContractRepository is the contract for persisting Contract aggregates.
type ContractRepository interface {
	// Save persists a Contract entity.
	Save(ctx context.Context, e *entity.Contract) error

	// FindByID retrieves a Contract by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Contract, error)

	// FindAll retrieves a paginated list of Contract entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Contract, error)

	// Count returns the total number of Contract entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Contract by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
