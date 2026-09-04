package entity

import "time"

// OrgMemberRole represents the role of a user within an organization.
type OrgMemberRole string

const (
	OrgRoleOwner  OrgMemberRole = "owner"
	OrgRoleAdmin  OrgMemberRole = "admin"
	OrgRoleMember OrgMemberRole = "member"
)

// OrganizationMember represents a user's membership in an organization.
type OrganizationMember struct {
	Id             string        `json:"id"`
	OrganizationId string        `json:"organization_id"`
	UserId         string        `json:"user_id"`
	Role           OrgMemberRole `json:"role"`
	InvitedBy      *string       `json:"invited_by,omitempty"`
	JoinedAt       time.Time     `json:"joined_at"`
	IsActive       bool          `json:"is_active"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      *time.Time    `json:"deleted_at,omitempty"`
}

// NewOrganizationMember creates a new OrganizationMember with role owner by default.
func NewOrganizationMember(orgID, userID string, role OrgMemberRole) *OrganizationMember {
	return &OrganizationMember{
		OrganizationId: orgID,
		UserId:         userID,
		Role:           role,
		IsActive:       true,
	}
}
