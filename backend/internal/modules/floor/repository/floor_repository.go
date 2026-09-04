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
	FindByID(ctx context.Context, id string) (*entity.Floor, error)

	// FindAll retrieves a paginated list of Floor entities.
	FindAll(ctx context.Context, limit, offset int, search string, buildingId string) ([]*entity.Floor, error)

	// Count returns the total number of non-deleted floors.
	Count(ctx context.Context, search string, buildingId string) (int64, error)

	// Delete removes a Floor by its primary key.
	Delete(ctx context.Context, id string) error
}
