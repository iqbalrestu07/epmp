package entity

import "time"

// Refund is the domain entity for refund.
type Refund struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	PaymentId      string     `json:"payment_id"`
	TenantId       string     `json:"tenant_id"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"`
	RefundDate     time.Time  `json:"refund_date"`
	Reason         string     `json:"reason"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewRefund creates a new Refund instance.
func NewRefund() *Refund {
	return &Refund{}
}
