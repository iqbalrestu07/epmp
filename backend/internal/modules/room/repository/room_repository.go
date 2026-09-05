package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/room/entity"
)

// RoomRepository is the contract for persisting Room aggregates.
type RoomRepository interface {
	// Save persists a Room entity.
	Save(ctx context.Context, e *entity.Room) error

	// FindByID retrieves a Room by its primary key.
	// If orgID is non-empty, the result must belong to that organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Room, error)

	// FindAll retrieves a paginated list of Room entities.
	// If orgID is non-empty, results are filtered to that organization.
	// propertyId and buildingId are optional cascading filters.
	FindAll(ctx context.Context, limit, offset int, search, floorId, propertyId, buildingId, orgID string) ([]*entity.Room, error)

	// Count returns the total number of non-deleted rooms.
	// If orgID is non-empty, count is filtered to that organization.
	// propertyId and buildingId are optional cascading filters.
	Count(ctx context.Context, search, floorId, propertyId, buildingId, orgID string) (int64, error)

	// Delete removes a Room by its primary key.
	Delete(ctx context.Context, id string) error
}
