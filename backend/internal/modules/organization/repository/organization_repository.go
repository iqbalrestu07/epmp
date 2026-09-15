package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/organization/dto"
	"github.com/epmp/backend/internal/modules/organization/entity"
)

// OrganizationRepository is the contract for persisting Organization aggregates.
type OrganizationRepository interface {
	// Save persists an Organization entity. Id must be pre-set by caller (uid.New()).
	Save(ctx context.Context, e *entity.Organization) error

	// FindByID retrieves an Organization by its primary key.
	FindByID(ctx context.Context, id string) (*entity.Organization, error)

	// FindAll retrieves a paginated list of Organization entities.
	FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.Organization, error)

	// FindByUserID retrieves all organizations where the user is an active member.
	FindByUserID(ctx context.Context, userID string) ([]*entity.Organization, error)

	// Count returns the total number of non-deleted organizations.
	Count(ctx context.Context, search string) (int64, error)

	// Delete soft-deletes an Organization by its primary key.
	Delete(ctx context.Context, id string) error

	// GetMemberRole returns the caller's role in an organization.
	GetMemberRole(ctx context.Context, orgID, userID string) (string, error)

	// SaveMember inserts or upserts an organization_members record.
	SaveMember(ctx context.Context, m *entity.OrganizationMember) error

	// FindMembersByOrgID returns all active members of an organization.
	FindMembersByOrgID(ctx context.Context, orgID string) ([]*entity.OrganizationMember, error)

	// FindMembersWithUserByOrgID returns active members with user details.
	FindMembersWithUserByOrgID(ctx context.Context, orgID string) ([]*dto.OrganizationMemberResponse, error)

	// AddMemberByEmail adds an existing user or creates a new user and joins them to the organization.
	AddMemberByEmail(ctx context.Context, orgID string, email, name, password, role, invitedBy string) (*dto.OrganizationMemberResponse, error)

	// DeleteMember removes a member from an organization.
	DeleteMember(ctx context.Context, orgID, userID string) error
}
