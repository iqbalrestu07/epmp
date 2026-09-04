package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/reservation/entity"
)

// ReservationRepository is the contract for persisting Reservation aggregates.
type ReservationRepository interface {
	// Save persists a Reservation entity.
	Save(ctx context.Context, e *entity.Reservation) error

	// FindByID retrieves a Reservation by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Reservation, error)

	// FindAll retrieves a paginated list of Reservation entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Reservation, error)

	// Count returns the total number of Reservation entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Reservation by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
