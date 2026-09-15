package entity

import "time"

// Facility is the domain entity for facility.
type Facility struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	PropertyId     string     `json:"property_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewFacility creates a new Facility instance.
func NewFacility() *Facility {
	return &Facility{}
}
