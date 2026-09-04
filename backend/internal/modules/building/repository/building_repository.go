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
	FindByID(ctx context.Context, id string) (*entity.Building, error)

	// FindAll retrieves a paginated list of Building entities.
	FindAll(ctx context.Context, limit, offset int, search string, propertyId string) ([]*entity.Building, error)

	// Count returns the total number of non-deleted buildings.
	Count(ctx context.Context, search string, propertyId string) (int64, error)

	// Delete removes a Building by its primary key.
	Delete(ctx context.Context, id string) error
}
