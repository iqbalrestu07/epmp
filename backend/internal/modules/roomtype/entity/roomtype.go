package entity

import "time"

// RoomType is the domain entity for roomtype.
type RoomType struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	BasePrice            float64      `json:"base_price"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewRoomType creates a new RoomType instance.
func NewRoomType() *RoomType {
	return &RoomType{}
}
