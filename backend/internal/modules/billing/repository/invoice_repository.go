package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/billing/entity"
)

// InvoiceRepository is the contract for persisting Invoice aggregates.
type InvoiceRepository interface {
	// Save persists a Invoice entity.
	Save(ctx context.Context, e *entity.Invoice) error

	// FindByID retrieves a Invoice by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Invoice, error)

	// FindAll retrieves a paginated list of Invoice entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Invoice, error)

	// Count returns the total number of Invoice entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Invoice by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
