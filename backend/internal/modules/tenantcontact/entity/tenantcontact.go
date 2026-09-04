package entity

import "time"

// TenantContact is the domain entity for tenantcontact.
type TenantContact struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	TenantId             string       `json:"tenant_id"`
	ContactType          string       `json:"contact_type"`
	ContactValue         string       `json:"contact_value"`
	IsPrimary            bool         `json:"is_primary"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewTenantContact creates a new TenantContact instance.
func NewTenantContact() *TenantContact {
	return &TenantContact{}
}
