package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/refund/entity"
)

// RefundRepository is the contract for persisting Refund aggregates.
type RefundRepository interface {
	// Save persists a Refund entity.
	Save(ctx context.Context, e *entity.Refund) error

	// FindByID retrieves a Refund by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Refund, error)

	// FindAll retrieves a paginated list of Refund entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Refund, error)

	// Count returns the total number of Refund entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Refund by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
