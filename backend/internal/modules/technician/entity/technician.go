package entity

import "time"

// Technician is the domain entity for technician.
type Technician struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	Phone          string     `json:"phone"`
	Specialty      string     `json:"specialty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewTechnician creates a new Technician instance.
func NewTechnician() *Technician {
	return &Technician{}
}
