package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/floor/entity"
)

// FloorRepository is the contract for persisting Floor aggregates.
type FloorRepository interface {
	// Save persists a Floor entity.
	Save(ctx context.Context, e *entity.Floor) error

	// FindByID retrieves a Floor by its primary key.
	// If orgID is non-empty, the result must belong to that organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Floor, error)

	// FindAll retrieves a paginated list of Floor entities.
	// If orgID is non-empty, results are filtered to that organization.
	FindAll(ctx context.Context, limit, offset int, search, buildingId, orgID string) ([]*entity.Floor, error)

	// Count returns the total number of non-deleted floors.
	// If orgID is non-empty, count is filtered to that organization.
	Count(ctx context.Context, search, buildingId, orgID string) (int64, error)

	// Delete removes a Floor by its primary key.
	Delete(ctx context.Context, id string) error
}
