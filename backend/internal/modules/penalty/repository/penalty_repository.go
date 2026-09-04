package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/penalty/entity"
)

// PenaltyRepository is the contract for persisting Penalty aggregates.
type PenaltyRepository interface {
	// Save persists a Penalty entity.
	Save(ctx context.Context, e *entity.Penalty) error

	// FindByID retrieves a Penalty by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Penalty, error)

	// FindAll retrieves a paginated list of Penalty entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Penalty, error)

	// Count returns the total number of Penalty entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Penalty by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
