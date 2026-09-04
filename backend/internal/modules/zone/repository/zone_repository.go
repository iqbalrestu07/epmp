package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/zone/entity"
)

// ZoneRepository is the contract for persisting Zone aggregates.
type ZoneRepository interface {
	// Save persists a Zone entity.
	Save(ctx context.Context, e *entity.Zone) error

	// FindByID retrieves a Zone by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Zone, error)

	// FindAll retrieves a paginated list of Zone entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Zone, error)

	// Count returns the total number of Zone entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Zone by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
