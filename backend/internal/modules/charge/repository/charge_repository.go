package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/charge/entity"
)

// ChargeRepository is the contract for persisting Charge aggregates.
type ChargeRepository interface {
	// Save persists a Charge entity.
	Save(ctx context.Context, e *entity.Charge) error

	// FindByID retrieves a Charge by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Charge, error)

	// FindAll retrieves a paginated list of Charge entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Charge, error)

	// Count returns the total number of Charge entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Charge by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
