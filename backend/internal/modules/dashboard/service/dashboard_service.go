package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/dashboard/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardService struct {
	db *pgxpool.Pool
}

func NewDashboardService(db *pgxpool.Pool) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetSummary(ctx context.Context, orgID string) (*dto.DashboardSummaryResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("dashboard service: orgID required")
	}

	res := &dto.DashboardSummaryResponse{
		Metrics:           dto.DashboardMetrics{},
		Alerts:            dto.DashboardAlerts{},
		PropertyOccupancy: make([]dto.PropertyOccupancySummary, 0),
		RecentActivities:  make([]dto.RecentActivity, 0),
	}

	// 1. Metrics
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM properties WHERE organization_id = $1 AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.TotalProperties)
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM rooms WHERE organization_id = $1 AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.TotalRooms)
	_ = s.db.QueryRow(ctx, `SELECT COUNT(DISTINCT room_id) FROM occupancies WHERE organization_id = $1 AND status = 'CheckedIn' AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.OccupiedRooms)

	if res.Metrics.TotalRooms > 0 {
		res.Metrics.VacantRooms = res.Metrics.TotalRooms - res.Metrics.OccupiedRooms
		if res.Metrics.VacantRooms < 0 {
			res.Metrics.VacantRooms = 0
		}
		res.Metrics.OccupancyRate = math.Round((float64(res.Metrics.OccupiedRooms)/float64(res.Metrics.TotalRooms))*1000) / 10
	}

	_ = s.db.QueryRow(ctx, `SELECT COUNT(DISTINCT tenant_id) FROM contracts WHERE organization_id = $1 AND status = 'Active' AND end_date >= now() AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.ActiveTenants)
	_ = s.db.QueryRow(ctx, `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE organization_id = $1 AND status = 'Success' AND created_at >= date_trunc('month', now()) AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.MonthlyRevenue)
	_ = s.db.QueryRow(ctx, `SELECT COALESCE(SUM(amount), 0) FROM invoices WHERE organization_id = $1 AND status != 'Paid' AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.UnpaidInvoicesAmount)
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM work_orders WHERE organization_id = $1 AND status IN ('Open', 'Pending', 'In Progress', 'Assigned') AND deleted_at IS NULL`, orgID).Scan(&res.Metrics.OpenWorkOrders)

	// 2. Alerts - Overdue Invoices
	overdueRows, err := s.db.Query(ctx, `
		SELECT i.id, i.tenant_id, COALESCE(t.full_name, 'Penyewa'), COALESCE(t.email, ''), COALESCE(t.phone, ''),
		       i.amount, i.due_date, GREATEST(0, EXTRACT(DAY FROM (now() - i.due_date))::int) as days_overdue,
		       COALESCE(p.name, 'Properti'), COALESCE(r.name, '-')
		FROM   invoices i
		LEFT JOIN tenants t ON t.id = i.tenant_id
		LEFT JOIN contracts c ON c.id = i.contract_id
		LEFT JOIN properties p ON p.id = c.property_id
		LEFT JOIN rooms r ON r.id = c.room_id
		WHERE  i.organization_id = $1 AND i.status != 'Paid' AND i.due_date < now() AND i.deleted_at IS NULL
		ORDER  BY i.due_date ASC
		LIMIT  20`, orgID)
	if err == nil {
		defer overdueRows.Close()
		for overdueRows.Next() {
			var a dto.OverdueInvoiceAlert
			if err := overdueRows.Scan(&a.Id, &a.TenantId, &a.TenantName, &a.TenantEmail, &a.TenantPhone, &a.Amount, &a.DueDate, &a.DaysOverdue, &a.PropertyName, &a.RoomName); err == nil {
				res.Alerts.OverdueInvoices = append(res.Alerts.OverdueInvoices, a)
			}
		}
	}

	// Alerts - Due Soon Invoices (next 7 days)
	dueSoonRows, err := s.db.Query(ctx, `
		SELECT i.id, i.tenant_id, COALESCE(t.full_name, 'Penyewa'), COALESCE(t.email, ''), COALESCE(t.phone, ''),
		       i.amount, i.due_date, GREATEST(0, EXTRACT(DAY FROM (i.due_date - now()))::int) as days_until_due,
		       COALESCE(p.name, 'Properti'), COALESCE(r.name, '-')
		FROM   invoices i
		LEFT JOIN tenants t ON t.id = i.tenant_id
		LEFT JOIN contracts c ON c.id = i.contract_id
		LEFT JOIN properties p ON p.id = c.property_id
		LEFT JOIN rooms r ON r.id = c.room_id
		WHERE  i.organization_id = $1 AND i.status != 'Paid' AND i.due_date >= now() AND i.due_date <= now() + interval '7 days' AND i.deleted_at IS NULL
		ORDER  BY i.due_date ASC
		LIMIT  20`, orgID)
	if err == nil {
		defer dueSoonRows.Close()
		for dueSoonRows.Next() {
			var a dto.DueSoonInvoiceAlert
			if err := dueSoonRows.Scan(&a.Id, &a.TenantId, &a.TenantName, &a.TenantEmail, &a.TenantPhone, &a.Amount, &a.DueDate, &a.DaysUntilDue, &a.PropertyName, &a.RoomName); err == nil {
				res.Alerts.DueSoonInvoices = append(res.Alerts.DueSoonInvoices, a)
			}
		}
	}

	// Alerts - Expiring Contracts (in next 30 days)
	expiringRows, err := s.db.Query(ctx, `
		SELECT c.id, c.tenant_id, COALESCE(t.full_name, 'Penyewa'), COALESCE(t.email, ''), COALESCE(t.phone, ''),
		       COALESCE(p.name, 'Properti'), COALESCE(r.name, '-'), c.start_date, c.end_date,
		       GREATEST(0, EXTRACT(DAY FROM (c.end_date - now()))::int) as days_until_expiration,
		       c.monthly_rent
		FROM   contracts c
		LEFT JOIN tenants t ON t.id = c.tenant_id
		LEFT JOIN properties p ON p.id = c.property_id
		LEFT JOIN rooms r ON r.id = c.room_id
		WHERE  c.organization_id = $1 AND c.status = 'Active' AND c.end_date >= now() AND c.end_date <= now() + interval '30 days' AND c.deleted_at IS NULL
		ORDER  BY c.end_date ASC
		LIMIT  20`, orgID)
	if err == nil {
		defer expiringRows.Close()
		for expiringRows.Next() {
			var a dto.ExpiringContractAlert
			if err := expiringRows.Scan(&a.Id, &a.TenantId, &a.TenantName, &a.TenantEmail, &a.TenantPhone, &a.PropertyName, &a.RoomName, &a.StartDate, &a.EndDate, &a.DaysUntilExpiration, &a.MonthlyRent); err == nil {
				res.Alerts.ExpiringContracts = append(res.Alerts.ExpiringContracts, a)
			}
		}
	}

	// Alerts - Expired Contracts (end_date < now)
	expiredRows, err := s.db.Query(ctx, `
		SELECT c.id, c.tenant_id, COALESCE(t.full_name, 'Penyewa'), COALESCE(t.email, ''), COALESCE(t.phone, ''),
		       COALESCE(p.name, 'Properti'), COALESCE(r.name, '-'), c.start_date, c.end_date,
		       GREATEST(0, EXTRACT(DAY FROM (now() - c.end_date))::int) as days_expired,
		       c.monthly_rent
		FROM   contracts c
		LEFT JOIN tenants t ON t.id = c.tenant_id
		LEFT JOIN properties p ON p.id = c.property_id
		LEFT JOIN rooms r ON r.id = c.room_id
		WHERE  c.organization_id = $1 AND c.end_date < now() AND c.deleted_at IS NULL
		ORDER  BY c.end_date DESC
		LIMIT  20`, orgID)
	if err == nil {
		defer expiredRows.Close()
		for expiredRows.Next() {
			var a dto.ExpiredContractAlert
			if err := expiredRows.Scan(&a.Id, &a.TenantId, &a.TenantName, &a.TenantEmail, &a.TenantPhone, &a.PropertyName, &a.RoomName, &a.StartDate, &a.EndDate, &a.DaysExpired, &a.MonthlyRent); err == nil {
				res.Alerts.ExpiredContracts = append(res.Alerts.ExpiredContracts, a)
			}
		}
	}

	// 3. Property Occupancy Breakdown
	propRows, err := s.db.Query(ctx, `
		SELECT p.id, p.name,
		       COUNT(DISTINCT r.id) as total_rooms,
		       COUNT(DISTINCT CASE WHEN o.status = 'CheckedIn' THEN o.room_id END) as occupied_rooms
		FROM   properties p
		LEFT JOIN rooms r ON r.property_id = p.id AND r.deleted_at IS NULL
		LEFT JOIN occupancies o ON o.room_id = r.id AND o.status = 'CheckedIn' AND o.deleted_at IS NULL
		WHERE  p.organization_id = $1 AND p.deleted_at IS NULL
		GROUP  BY p.id, p.name
		ORDER  BY p.name ASC
		LIMIT  10`, orgID)
	if err == nil {
		defer propRows.Close()
		for propRows.Next() {
			var item dto.PropertyOccupancySummary
			if err := propRows.Scan(&item.PropertyId, &item.PropertyName, &item.TotalRooms, &item.Occupied); err == nil {
				if item.TotalRooms > 0 {
					item.Rate = math.Round((float64(item.Occupied)/float64(item.TotalRooms))*1000) / 10
				}
				res.PropertyOccupancy = append(res.PropertyOccupancy, item)
			}
		}
	}

	// 4. Recent Activities (from payments, check-ins, work orders)
	recentRows, err := s.db.Query(ctx, `
		(SELECT id, 'payment' as type, 'Pembayaran Masuk' as title, 'Rp ' || amount::text || ' via ' || payment_method as description, created_at as timestamp
		 FROM payments WHERE organization_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 5)
		UNION ALL
		(SELECT id, 'occupancy' as type, 'Check-in Penghuni' as title, 'Status: ' || status as description, created_at as timestamp
		 FROM occupancies WHERE organization_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 5)
		UNION ALL
		(SELECT id, 'work_order' as type, 'Tiket Perbaikan: ' || title as title, 'Prioritas: ' || priority as description, created_at as timestamp
		 FROM work_orders WHERE organization_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 5)
		ORDER BY timestamp DESC
		LIMIT 10`, orgID)
	if err == nil {
		defer recentRows.Close()
		for recentRows.Next() {
			var act dto.RecentActivity
			if err := recentRows.Scan(&act.Id, &act.Type, &act.Title, &act.Description, &act.Timestamp); err == nil {
				res.RecentActivities = append(res.RecentActivities, act)
			}
		}
	}

	return res, nil
}
