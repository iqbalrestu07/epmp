package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/facility/entity"
)

// FacilityRepository is the contract for persisting Facility aggregates.
type FacilityRepository interface {
	// Save persists a Facility entity.
	Save(ctx context.Context, e *entity.Facility) error

	// FindByID retrieves a Facility by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Facility, error)

	// FindAll retrieves a paginated list of Facility entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Facility, error)

	// Count returns the total number of Facility entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Facility by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
