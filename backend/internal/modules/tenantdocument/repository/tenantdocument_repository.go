package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/tenantdocument/entity"
)

// TenantDocumentRepository is the contract for persisting TenantDocument aggregates.
type TenantDocumentRepository interface {
	// Save persists a TenantDocument entity.
	Save(ctx context.Context, e *entity.TenantDocument) error

	// FindByID retrieves a TenantDocument by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.TenantDocument, error)

	// FindAll retrieves a paginated list of TenantDocument entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantDocument, error)

	// Count returns the total number of TenantDocument entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a TenantDocument by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
