package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/bed/entity"
)

// BedRepository is the contract for persisting Bed aggregates.
type BedRepository interface {
	// Save persists a Bed entity.
	Save(ctx context.Context, e *entity.Bed) error

	// FindByID retrieves a Bed by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Bed, error)

	// FindAll retrieves a paginated list of Bed entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Bed, error)

	// Count returns the total number of Bed entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Bed by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
