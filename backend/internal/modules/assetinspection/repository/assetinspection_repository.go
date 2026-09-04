package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/assetinspection/entity"
)

// AssetInspectionRepository is the contract for persisting AssetInspection aggregates.
type AssetInspectionRepository interface {
	// Save persists a AssetInspection entity.
	Save(ctx context.Context, e *entity.AssetInspection) error

	// FindByID retrieves a AssetInspection by its primary key within an organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.AssetInspection, error)

	// FindAll retrieves a paginated list of AssetInspection entities within an organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.AssetInspection, error)

	// Count returns the total number of AssetInspection entities matching the filter.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a AssetInspection by its primary key.
	Delete(ctx context.Context, id, orgID string) error
}
