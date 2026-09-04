package entity

import "time"

// Zone is the domain entity for zone.
type Zone struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	BuildingId           string       `json:"building_id"`
	Floor                int          `json:"floor"`
	Name                 string       `json:"name"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewZone creates a new Zone instance.
func NewZone() *Zone {
	return &Zone{}
}
