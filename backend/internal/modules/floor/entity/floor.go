package entity

import "time"

// Floor is the domain entity for floor.
type Floor struct {
	Id             string     `json:"id"`
	OrganizationId string     `json:"organization_id"`
	BuildingId     string     `json:"building_id"`
	Name           string     `json:"name"`
	FloorNumber    int        `json:"floor_number"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// NewFloor creates a new Floor instance.
func NewFloor() *Floor {
	return &Floor{}
}
