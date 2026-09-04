package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/tenantidentity/entity"
)

// TenantIdentityRepository is the contract for persisting TenantIdentity aggregates.
type TenantIdentityRepository interface {
	// Save persists a TenantIdentity entity.
	Save(ctx context.Context, e *entity.TenantIdentity) error

	// FindByID retrieves a TenantIdentity by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.TenantIdentity, error)

	// FindAll retrieves a paginated list of TenantIdentity entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantIdentity, error)

	// Count returns the total number of TenantIdentity entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a TenantIdentity by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
