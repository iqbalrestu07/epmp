package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/workorder/entity"
)

// WorkOrderRepository is the contract for persisting WorkOrder aggregates.
type WorkOrderRepository interface {
	// Save persists a WorkOrder entity.
	Save(ctx context.Context, e *entity.WorkOrder) error

	// FindByID retrieves a WorkOrder by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.WorkOrder, error)

	// FindAll retrieves a paginated list of WorkOrder entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.WorkOrder, error)

	// Count returns the total number of WorkOrder entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a WorkOrder by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
