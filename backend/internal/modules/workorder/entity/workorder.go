package entity

import "time"

// WorkOrder is the domain entity for workorder.
type WorkOrder struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	PropertyId           string       `json:"property_id"`
	RoomId               string       `json:"room_id"`
	Description          string       `json:"description"`
	Status               string       `json:"status"`
	Priority             string       `json:"priority"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewWorkOrder creates a new WorkOrder instance.
func NewWorkOrder() *WorkOrder {
	return &WorkOrder{}
}
