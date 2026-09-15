package entity

import "time"

// Occupancy is the domain entity for occupancy.
type Occupancy struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	ContractId     string     `json:"contract_id"`
	RoomId         string     `json:"room_id"`
	TenantId       string     `json:"tenant_id"`
	Status         string     `json:"status"`
	CheckInTime    time.Time  `json:"check_in_time"`
	CheckOutTime   time.Time  `json:"check_out_time"`
	Notes          string     `json:"notes"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewOccupancy creates a new Occupancy instance.
func NewOccupancy() *Occupancy {
	return &Occupancy{}
}
