package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/payment/entity"
)

// PaymentRepository is the contract for persisting Payment aggregates.
type PaymentRepository interface {
	// Save persists a Payment entity.
	Save(ctx context.Context, e *entity.Payment) error

	// FindByID retrieves a Payment by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Payment, error)

	// FindAll retrieves a paginated list of Payment entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Payment, error)

	// Count returns the total number of Payment entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Payment by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
