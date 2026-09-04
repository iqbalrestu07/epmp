package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/adjustment/entity"
)

// AdjustmentRepository is the contract for persisting Adjustment aggregates.
type AdjustmentRepository interface {
	// Save persists a Adjustment entity.
	Save(ctx context.Context, e *entity.Adjustment) error

	// FindByID retrieves a Adjustment by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Adjustment, error)

	// FindAll retrieves a paginated list of Adjustment entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Adjustment, error)

	// Count returns the total number of Adjustment entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Adjustment by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
