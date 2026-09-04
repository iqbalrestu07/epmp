package entity

import "time"

// Reservation is the domain entity for reservation.
type Reservation struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	TenantId             string       `json:"tenant_id"`
	PropertyId           string       `json:"property_id"`
	RoomId               string       `json:"room_id"`
	Status               string       `json:"status"`
	CheckInDate          time.Time    `json:"check_in_date"`
	CheckOutDate         time.Time    `json:"check_out_date"`
	BookingFee           float64      `json:"booking_fee"`
	Notes                string       `json:"notes"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewReservation creates a new Reservation instance.
func NewReservation() *Reservation {
	return &Reservation{}
}
