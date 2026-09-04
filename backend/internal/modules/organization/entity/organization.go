package entity

import "time"

// Organization is the domain entity for organization.
type Organization struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewOrganization creates a new Organization instance.
func NewOrganization() *Organization {
	return &Organization{}
}
