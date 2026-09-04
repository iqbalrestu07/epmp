package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/assetassignment/entity"
)

// AssetAssignmentRepository is the contract for persisting AssetAssignment aggregates.
type AssetAssignmentRepository interface {
	// Save persists a AssetAssignment entity.
	Save(ctx context.Context, e *entity.AssetAssignment) error

	// FindByID retrieves a AssetAssignment by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.AssetAssignment, error)

	// FindAll retrieves a paginated list of AssetAssignment entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.AssetAssignment, error)

	// Count returns the total number of AssetAssignment entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a AssetAssignment by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
