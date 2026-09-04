package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/occupancy/entity"
)

// OccupancyRepository is the contract for persisting Occupancy aggregates.
type OccupancyRepository interface {
	// Save persists a Occupancy entity.
	Save(ctx context.Context, e *entity.Occupancy) error

	// FindByID retrieves a Occupancy by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Occupancy, error)

	// FindAll retrieves a paginated list of Occupancy entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Occupancy, error)

	// Count returns the total number of Occupancy entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Occupancy by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
