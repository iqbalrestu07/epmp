package dto

import "time"

// OccupancyBuildingRow is the per-building occupancy breakdown.
type OccupancyBuildingRow struct {
	BuildingID    string  `json:"building_id"`
	BuildingName  string  `json:"building_name"`
	TotalRooms    int64   `json:"total_rooms"`
	Occupied      int64   `json:"occupied"`
	Reserved      int64   `json:"reserved"`
	Maintenance   int64   `json:"maintenance"`
	Available     int64   `json:"available"`
	OccupancyRate float64 `json:"occupancy_rate"`
}

// OccupancyPropertyRow is the per-property occupancy breakdown.
type OccupancyPropertyRow struct {
	PropertyID    string                 `json:"property_id"`
	PropertyName  string                 `json:"property_name"`
	TotalRooms    int64                  `json:"total_rooms"`
	Occupied      int64                  `json:"occupied"`
	Reserved      int64                  `json:"reserved"`
	Maintenance   int64                  `json:"maintenance"`
	Available     int64                  `json:"available"`
	OccupancyRate float64                `json:"occupancy_rate"`
	Buildings     []OccupancyBuildingRow `json:"buildings"`
}

// OccupancyTotals aggregates occupancy across all properties.
type OccupancyTotals struct {
	TotalRooms    int64   `json:"total_rooms"`
	Occupied      int64   `json:"occupied"`
	Reserved      int64   `json:"reserved"`
	Maintenance   int64   `json:"maintenance"`
	Available     int64   `json:"available"`
	OccupancyRate float64 `json:"occupancy_rate"`
}

// OccupancyReportResponse is the payload for GET /reports/occupancy.
type OccupancyReportResponse struct {
	Totals     OccupancyTotals        `json:"totals"`
	Properties []OccupancyPropertyRow `json:"properties"`
}

// RevenueMonthlyRow aggregates invoiced vs received amounts per month and currency.
type RevenueMonthlyRow struct {
	Period           string  `json:"period"` // YYYY-MM
	Currency         string  `json:"currency"`
	InvoiceCount     int64   `json:"invoice_count"`
	InvoicedAmount   float64 `json:"invoiced_amount"`
	PaidInvoiceCount int64   `json:"paid_invoice_count"`
	PaymentCount     int64   `json:"payment_count"`
	ReceivedAmount   float64 `json:"received_amount"`
}

// RevenueReportResponse is the payload for GET /reports/revenue.
type RevenueReportResponse struct {
	From string              `json:"from"`
	To   string              `json:"to"`
	Rows []RevenueMonthlyRow `json:"rows"`
}

// ArAgingInvoice is one outstanding invoice in the aging report.
type ArAgingInvoice struct {
	InvoiceID   string    `json:"invoice_id"`
	TenantID    string    `json:"tenant_id"`
	TenantName  string    `json:"tenant_name"`
	Amount      float64   `json:"amount"`
	PaidAmount  float64   `json:"paid_amount"`
	Outstanding float64   `json:"outstanding"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	DaysOverdue int       `json:"days_overdue"`
	Bucket      string    `json:"bucket"`
}

// ArAgingBucketRow summarizes outstanding amounts per aging bucket and currency.
type ArAgingBucketRow struct {
	Bucket       string  `json:"bucket"`
	Currency     string  `json:"currency"`
	InvoiceCount int64   `json:"invoice_count"`
	Outstanding  float64 `json:"outstanding"`
}

// ArAgingReportResponse is the payload for GET /reports/ar-aging.
type ArAgingReportResponse struct {
	Buckets  []ArAgingBucketRow `json:"buckets"`
	Invoices []ArAgingInvoice   `json:"invoices"`
}
