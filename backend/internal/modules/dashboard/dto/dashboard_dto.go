package dto

import "time"

type DashboardMetrics struct {
	TotalProperties      int64   `json:"total_properties"`
	TotalRooms           int64   `json:"total_rooms"`
	OccupiedRooms        int64   `json:"occupied_rooms"`
	VacantRooms          int64   `json:"vacant_rooms"`
	OccupancyRate        float64 `json:"occupancy_rate"`
	ActiveTenants        int64   `json:"active_tenants"`
	MonthlyRevenue       float64 `json:"monthly_revenue"`
	UnpaidInvoicesAmount float64 `json:"unpaid_invoices_amount"`
	OpenWorkOrders       int64   `json:"open_work_orders"`
}

type OverdueInvoiceAlert struct {
	Id           string    `json:"id"`
	TenantId     string    `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TenantEmail  string    `json:"tenant_email"`
	TenantPhone  string    `json:"tenant_phone"`
	Amount       float64   `json:"amount"`
	DueDate      time.Time `json:"due_date"`
	DaysOverdue  int       `json:"days_overdue"`
	PropertyName string    `json:"property_name"`
	RoomName     string    `json:"room_name"`
}

type DueSoonInvoiceAlert struct {
	Id           string    `json:"id"`
	TenantId     string    `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TenantEmail  string    `json:"tenant_email"`
	TenantPhone  string    `json:"tenant_phone"`
	Amount       float64   `json:"amount"`
	DueDate      time.Time `json:"due_date"`
	DaysUntilDue int       `json:"days_until_due"`
	PropertyName string    `json:"property_name"`
	RoomName     string    `json:"room_name"`
}

type ExpiringContractAlert struct {
	Id                  string    `json:"id"`
	TenantId            string    `json:"tenant_id"`
	TenantName          string    `json:"tenant_name"`
	TenantEmail         string    `json:"tenant_email"`
	TenantPhone         string    `json:"tenant_phone"`
	PropertyName        string    `json:"property_name"`
	RoomName            string    `json:"room_name"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	DaysUntilExpiration int       `json:"days_until_expiration"`
	MonthlyRent         float64   `json:"monthly_rent"`
}

type ExpiredContractAlert struct {
	Id           string    `json:"id"`
	TenantId     string    `json:"tenant_id"`
	TenantName   string    `json:"tenant_name"`
	TenantEmail  string    `json:"tenant_email"`
	TenantPhone  string    `json:"tenant_phone"`
	PropertyName string    `json:"property_name"`
	RoomName     string    `json:"room_name"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	DaysExpired  int       `json:"days_expired"`
	MonthlyRent  float64   `json:"monthly_rent"`
}

type DashboardAlerts struct {
	OverdueInvoices   []OverdueInvoiceAlert   `json:"overdue_invoices"`
	DueSoonInvoices   []DueSoonInvoiceAlert   `json:"due_soon_invoices"`
	ExpiringContracts []ExpiringContractAlert `json:"expiring_contracts"`
	ExpiredContracts  []ExpiredContractAlert  `json:"expired_contracts"`
}

type RecentActivity struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type PropertyOccupancySummary struct {
	PropertyId   string  `json:"property_id"`
	PropertyName string  `json:"property_name"`
	TotalRooms   int64   `json:"total_rooms"`
	Occupied     int64   `json:"occupied"`
	Rate         float64 `json:"rate"`
}

type DashboardSummaryResponse struct {
	Metrics           DashboardMetrics           `json:"metrics"`
	Alerts            DashboardAlerts            `json:"alerts"`
	PropertyOccupancy []PropertyOccupancySummary `json:"property_occupancy"`
	RecentActivities  []RecentActivity           `json:"recent_activities"`
}
