package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/property/dto"
	"github.com/epmp/backend/internal/modules/property/entity"
)

// PropertyRepository is the contract for persisting Property aggregates.
type PropertyRepository interface {
	// Save persists a Property entity.
	Save(ctx context.Context, e *entity.Property) error

	// FindByID retrieves a Property by its primary key.
	// If orgID is non-empty, the result must belong to that organization.
	FindByID(ctx context.Context, id, orgID string) (*entity.Property, error)

	// FindAll retrieves a paginated list of Property entities.
	// If orgID is non-empty, results are filtered to that organization.
	FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Property, error)

	// Count returns the total number of non-deleted Properties.
	// If orgID is non-empty, count is filtered to that organization.
	Count(ctx context.Context, search, orgID string) (int64, error)

	// Delete removes a Property by its primary key.
	Delete(ctx context.Context, id, orgID string) error

	// FindStaffByPropertyID returns all assigned staff for a property.
	FindStaffByPropertyID(ctx context.Context, propertyID, orgID string) ([]*dto.PropertyStaffResponse, error)

	// AssignStaff assigns a user with a role to a property.
	AssignStaff(ctx context.Context, orgID, propertyID, userID, roleID string) (*dto.PropertyStaffResponse, error)

	// RemoveStaff removes a staff assignment.
	RemoveStaff(ctx context.Context, propertyID, userID, roleID, orgID string) error
}
