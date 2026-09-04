package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/tenant/entity"
)

// TenantRepository is the contract for persisting Tenant aggregates.
type TenantRepository interface {
	// Save persists a Tenant entity.
	Save(ctx context.Context, e *entity.Tenant) error

	// FindByID retrieves a Tenant by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Tenant, error)

	// FindAll retrieves a paginated list of Tenant entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Tenant, error)

	// Count returns the total number of Tenant entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Tenant by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
