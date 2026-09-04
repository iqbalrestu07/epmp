package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/supplier/entity"
)

// SupplierRepository is the contract for persisting Supplier aggregates.
type SupplierRepository interface {
	// Save persists a Supplier entity.
	Save(ctx context.Context, e *entity.Supplier) error

	// FindByID retrieves a Supplier by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Supplier, error)

	// FindAll retrieves a paginated list of Supplier entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Supplier, error)

	// Count returns the total number of Supplier entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Supplier by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
