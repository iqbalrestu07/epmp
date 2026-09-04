package entity

import "time"

// Payment is the domain entity for payment.
type Payment struct {
	OrganizationId string     `json:"organization_id"`
	Id             string     `json:"id"`
	InvoiceId            string       `json:"invoice_id"`
	TenantId             string       `json:"tenant_id"`
	Amount               float64      `json:"amount"`
	PaymentDate          time.Time    `json:"payment_date"`
	PaymentMethod        string       `json:"payment_method"`
	Status               string       `json:"status"`
	ReferenceNumber      string       `json:"reference_number"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// NewPayment creates a new Payment instance.
func NewPayment() *Payment {
	return &Payment{}
}
