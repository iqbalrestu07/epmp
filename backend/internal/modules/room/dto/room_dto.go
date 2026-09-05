package dto

import "time"

// CreateRoomRequest is the DTO for creating a Room.
type CreateRoomRequest struct {
	OrganizationId string  `json:"organization_id"`
	PropertyId     string  `json:"property_id"`
	FloorId        string  `json:"floor_id"`
	RoomTypeId     string  `json:"room_type_id,omitempty"`
	Name           string  `json:"name"`
	Capacity       int     `json:"capacity"`
	Price          float64 `json:"price"`
	IsAvailable    bool    `json:"is_available"`
	Status         string  `json:"status,omitempty"`
	Currency       string  `json:"currency,omitempty"`
}

// UpdateRoomRequest is the DTO for updating a Room.
type UpdateRoomRequest struct {
	OrganizationId string  `json:"organization_id"`
	PropertyId     string  `json:"property_id"`
	FloorId        string  `json:"floor_id"`
	RoomTypeId     string  `json:"room_type_id,omitempty"`
	Name           string  `json:"name"`
	Capacity       int     `json:"capacity"`
	Price          float64 `json:"price"`
	IsAvailable    bool    `json:"is_available"`
	Status         string  `json:"status,omitempty"`
	Currency       string  `json:"currency,omitempty"`
}

// RoomResponse is the DTO for returning a Room.
type RoomResponse struct {
	OrganizationId string    `json:"organization_id"`
	Id             string    `json:"id"`
	PropertyId     string    `json:"property_id"`
	FloorId        string    `json:"floor_id,omitempty"`
	RoomTypeId     string    `json:"room_type_id,omitempty"`
	Name           string    `json:"name"`
	Capacity       int       `json:"capacity"`
	Price          float64   `json:"price"`
	IsAvailable    bool      `json:"is_available"`
	Status         string    `json:"status"`
	Currency       string    `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RoomListResponse is the DTO for a paginated list of Room.
type RoomListResponse struct {
	Data       []RoomResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}
