package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/building/entity"
)

// BuildingRepository is the contract for persisting Building aggregates.
type BuildingRepository interface {
	// Save persists a Building entity.
	Save(ctx context.Context, e *entity.Building) error

	// FindByID retrieves a Building by its primary key.
	// If orgID is non-empty, the result must belong to that organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Building, error)

	// FindAll retrieves a paginated list of Building entities.
	// If orgID is non-empty, results are filtered to that organization.
	FindAll(ctx context.Context, limit, offset int, search, propertyId, orgID string) ([]*entity.Building, error)

	// Count returns the total number of non-deleted buildings.
	// If orgID is non-empty, count is filtered to that organization.
	Count(ctx context.Context, search, propertyId, orgID string) (int64, error)

	// Delete removes a Building by its primary key.
	Delete(ctx context.Context, id string) error
}
