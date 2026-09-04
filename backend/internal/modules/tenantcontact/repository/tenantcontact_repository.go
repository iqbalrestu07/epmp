package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/tenantcontact/entity"
)

// TenantContactRepository is the contract for persisting TenantContact aggregates.
type TenantContactRepository interface {
	// Save persists a TenantContact entity.
	Save(ctx context.Context, e *entity.TenantContact) error

	// FindByID retrieves a TenantContact by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.TenantContact, error)

	// FindAll retrieves a paginated list of TenantContact entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantContact, error)

	// Count returns the total number of TenantContact entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a TenantContact by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
