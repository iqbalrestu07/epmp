package entity

import "time"

// Building is the domain entity for building.
type Building struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	PropertyId     string     `json:"property_id"`
	Name           string     `json:"name"`
	TotalFloors    int        `json:"total_floors"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// NewBuilding creates a new Building instance.
func NewBuilding() *Building {
	return &Building{}
}
