package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/technician/entity"
)

// TechnicianRepository is the contract for persisting Technician aggregates.
type TechnicianRepository interface {
	// Save persists a Technician entity.
	Save(ctx context.Context, e *entity.Technician) error

	// FindByID retrieves a Technician by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Technician, error)

	// FindAll retrieves a paginated list of Technician entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Technician, error)

	// Count returns the total number of Technician entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Technician by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
