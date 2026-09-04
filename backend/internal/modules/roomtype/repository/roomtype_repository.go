package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/roomtype/entity"
)

// RoomTypeRepository is the contract for persisting RoomType aggregates.
type RoomTypeRepository interface {
	// Save persists a RoomType entity.
	Save(ctx context.Context, e *entity.RoomType) error

	// FindByID retrieves a RoomType by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.RoomType, error)

	// FindAll retrieves a paginated list of RoomType entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.RoomType, error)

	// Count returns the total number of RoomType entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a RoomType by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
