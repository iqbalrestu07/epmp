package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/epmp/backend/internal/modules/report/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReportService provides read-only aggregation queries for reporting.
type ReportService struct {
	db *pgxpool.Pool
}

// NewReportService creates a new ReportService.
func NewReportService(db *pgxpool.Pool) *ReportService {
	return &ReportService{db: db}
}

func occupancyRate(occupied, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return math.Round((float64(occupied)/float64(total))*1000) / 10
}

// GetOccupancyReport returns room status breakdown per property (and per building).
func (s *ReportService) GetOccupancyReport(ctx context.Context, orgID, propertyID, buildingID string) (*dto.OccupancyReportResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("report service: occupancy: organization ID required")
	}

	res := &dto.OccupancyReportResponse{Properties: make([]dto.OccupancyPropertyRow, 0)}

	// Per-property aggregation via rooms.property_id.
	// When building_id is given, only rooms on floors of that building are counted.
	roomJoin := "r.property_id = p.id AND r.deleted_at IS NULL"
	args := []interface{}{orgID}
	argIdx := 2
	if buildingID != "" {
		roomJoin += fmt.Sprintf(" AND r.floor_id IN (SELECT id FROM floors WHERE building_id = $%d AND deleted_at IS NULL)", argIdx)
		args = append(args, buildingID)
		argIdx++
	}

	propQuery := fmt.Sprintf(`
		SELECT p.id, p.name,
			COUNT(r.id)                                            AS total,
			COUNT(r.id) FILTER (WHERE r.status = 'Occupied')      AS occupied,
			COUNT(r.id) FILTER (WHERE r.status = 'Reserved')      AS reserved,
			COUNT(r.id) FILTER (WHERE r.status = 'Maintenance')   AS maintenance,
			COUNT(r.id) FILTER (WHERE r.status = 'Available')     AS available
		FROM   properties p
		LEFT JOIN rooms r ON %s
		WHERE  p.organization_id = $1 AND p.deleted_at IS NULL`, roomJoin)
	if propertyID != "" {
		propQuery += fmt.Sprintf(" AND p.id = $%d", argIdx)
		args = append(args, propertyID)
		argIdx++
	}
	propQuery += " GROUP BY p.id, p.name ORDER BY p.name"

	propRows, err := s.db.Query(ctx, propQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("report service: occupancy: properties: %w", err)
	}
	propIndex := map[string]int{}
	for propRows.Next() {
		var row dto.OccupancyPropertyRow
		if err := propRows.Scan(&row.PropertyID, &row.PropertyName, &row.TotalRooms, &row.Occupied, &row.Reserved, &row.Maintenance, &row.Available); err != nil {
			propRows.Close()
			return nil, fmt.Errorf("report service: occupancy: scan property: %w", err)
		}
		row.OccupancyRate = occupancyRate(row.Occupied, row.TotalRooms)
		row.Buildings = make([]dto.OccupancyBuildingRow, 0)
		propIndex[row.PropertyID] = len(res.Properties)
		res.Properties = append(res.Properties, row)

		res.Totals.TotalRooms += row.TotalRooms
		res.Totals.Occupied += row.Occupied
		res.Totals.Reserved += row.Reserved
		res.Totals.Maintenance += row.Maintenance
		res.Totals.Available += row.Available
	}
	propRows.Close()
	if err := propRows.Err(); err != nil {
		return nil, fmt.Errorf("report service: occupancy: properties rows: %w", err)
	}
	res.Totals.OccupancyRate = occupancyRate(res.Totals.Occupied, res.Totals.TotalRooms)

	// Per-building aggregation via floors -> rooms.
	bldgQuery := `
		SELECT p.id, b.id, b.name,
			COUNT(r.id)                                            AS total,
			COUNT(r.id) FILTER (WHERE r.status = 'Occupied')      AS occupied,
			COUNT(r.id) FILTER (WHERE r.status = 'Reserved')      AS reserved,
			COUNT(r.id) FILTER (WHERE r.status = 'Maintenance')   AS maintenance,
			COUNT(r.id) FILTER (WHERE r.status = 'Available')     AS available
		FROM   buildings b
		JOIN   properties p ON p.id = b.property_id AND p.deleted_at IS NULL
		LEFT JOIN floors f ON f.building_id = b.id AND f.deleted_at IS NULL
		LEFT JOIN rooms   r ON r.floor_id   = f.id AND r.deleted_at IS NULL
		WHERE  b.organization_id = $1 AND b.deleted_at IS NULL`
	bArgs := []interface{}{orgID}
	bIdx := 2
	if propertyID != "" {
		bldgQuery += fmt.Sprintf(" AND b.property_id = $%d", bIdx)
		bArgs = append(bArgs, propertyID)
		bIdx++
	}
	if buildingID != "" {
		bldgQuery += fmt.Sprintf(" AND b.id = $%d", bIdx)
		bArgs = append(bArgs, buildingID)
		bIdx++
	}
	bldgQuery += " GROUP BY p.id, b.id, b.name ORDER BY b.name"

	bldgRows, err := s.db.Query(ctx, bldgQuery, bArgs...)
	if err != nil {
		return nil, fmt.Errorf("report service: occupancy: buildings: %w", err)
	}
	defer bldgRows.Close()
	for bldgRows.Next() {
		var propID string
		var row dto.OccupancyBuildingRow
		if err := bldgRows.Scan(&propID, &row.BuildingID, &row.BuildingName, &row.TotalRooms, &row.Occupied, &row.Reserved, &row.Maintenance, &row.Available); err != nil {
			return nil, fmt.Errorf("report service: occupancy: scan building: %w", err)
		}
		row.OccupancyRate = occupancyRate(row.Occupied, row.TotalRooms)
		if idx, ok := propIndex[propID]; ok {
			res.Properties[idx].Buildings = append(res.Properties[idx].Buildings, row)
		}
	}

	return res, nil
}

// GetRevenueReport returns invoiced vs received amounts per month and currency.
func (s *ReportService) GetRevenueReport(ctx context.Context, orgID string, from, to time.Time) (*dto.RevenueReportResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("report service: revenue: organization ID required")
	}

	type key struct {
		Period   string
		Currency string
	}
	rows := map[key]*dto.RevenueMonthlyRow{}

	// Invoiced amounts grouped by due_date month.
	invRows, err := s.db.Query(ctx, `
		SELECT to_char(date_trunc('month', due_date), 'YYYY-MM') AS period,
		       currency,
		       COUNT(*)              AS invoice_count,
		       COALESCE(SUM(amount), 0) AS invoiced,
		       COUNT(*) FILTER (WHERE status = 'Paid') AS paid_invoice_count
		FROM   invoices
		WHERE  organization_id = $1 AND deleted_at IS NULL
		  AND  due_date >= $2 AND due_date < $3
		GROUP  BY 1, 2`, orgID, from, to)
	if err != nil {
		return nil, fmt.Errorf("report service: revenue: invoices: %w", err)
	}
	for invRows.Next() {
		var period, currency string
		var count, paidCount int64
		var invoiced float64
		if err := invRows.Scan(&period, &currency, &count, &invoiced, &paidCount); err != nil {
			invRows.Close()
			return nil, fmt.Errorf("report service: revenue: scan invoice: %w", err)
		}
		rows[key{period, currency}] = &dto.RevenueMonthlyRow{
			Period: period, Currency: currency,
			InvoiceCount: count, InvoicedAmount: invoiced, PaidInvoiceCount: paidCount,
		}
	}
	invRows.Close()
	if err := invRows.Err(); err != nil {
		return nil, fmt.Errorf("report service: revenue: invoice rows: %w", err)
	}

	// Payments received grouped by payment_date month.
	payRows, err := s.db.Query(ctx, `
		SELECT to_char(date_trunc('month', payment_date), 'YYYY-MM') AS period,
		       currency,
		       COUNT(*)                 AS payment_count,
		       COALESCE(SUM(amount), 0) AS received
		FROM   payments
		WHERE  organization_id = $1 AND deleted_at IS NULL
		  AND  status IN ('Success', 'Completed')
		  AND  payment_date >= $2 AND payment_date < $3
		GROUP  BY 1, 2`, orgID, from, to)
	if err != nil {
		return nil, fmt.Errorf("report service: revenue: payments: %w", err)
	}
	for payRows.Next() {
		var period, currency string
		var count int64
		var received float64
		if err := payRows.Scan(&period, &currency, &count, &received); err != nil {
			payRows.Close()
			return nil, fmt.Errorf("report service: revenue: scan payment: %w", err)
		}
		k := key{period, currency}
		if r, ok := rows[k]; ok {
			r.PaymentCount = count
			r.ReceivedAmount = received
		} else {
			rows[k] = &dto.RevenueMonthlyRow{
				Period: period, Currency: currency,
				PaymentCount: count, ReceivedAmount: received,
			}
		}
	}
	payRows.Close()
	if err := payRows.Err(); err != nil {
		return nil, fmt.Errorf("report service: revenue: payment rows: %w", err)
	}

	out := make([]dto.RevenueMonthlyRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Period != out[j].Period {
			return out[i].Period < out[j].Period
		}
		return out[i].Currency < out[j].Currency
	})

	return &dto.RevenueReportResponse{
		From: from.Format("2006-01-02"),
		To:   to.Format("2006-01-02"),
		Rows: out,
	}, nil
}

// agingBucket maps days overdue to a bucket label.
func agingBucket(daysOverdue int) string {
	switch {
	case daysOverdue <= 0:
		return "Current"
	case daysOverdue <= 30:
		return "1-30"
	case daysOverdue <= 60:
		return "31-60"
	case daysOverdue <= 90:
		return "61-90"
	default:
		return "90+"
	}
}

// GetArAgingReport returns outstanding invoices with aging buckets.
func (s *ReportService) GetArAgingReport(ctx context.Context, orgID string) (*dto.ArAgingReportResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("report service: ar-aging: organization ID required")
	}

	rows, err := s.db.Query(ctx, `
		SELECT i.id, i.tenant_id, COALESCE(t.full_name, '-'),
		       i.amount, i.currency, i.status, i.due_date,
		       COALESCE(paid.total, 0) AS paid_amount
		FROM   invoices i
		LEFT JOIN tenants t ON t.id = i.tenant_id
		LEFT JOIN (
			SELECT invoice_id, SUM(amount) AS total
			FROM   payments
			WHERE  status IN ('Success', 'Completed') AND deleted_at IS NULL
			GROUP  BY invoice_id
		) paid ON paid.invoice_id = i.id
		WHERE  i.organization_id = $1 AND i.deleted_at IS NULL
		  AND  i.status NOT IN ('Paid', 'Cancelled')
		ORDER  BY i.due_date ASC`, orgID)
	if err != nil {
		return nil, fmt.Errorf("report service: ar-aging: %w", err)
	}
	defer rows.Close()

	now := time.Now()
	res := &dto.ArAgingReportResponse{
		Buckets:  make([]dto.ArAgingBucketRow, 0),
		Invoices: make([]dto.ArAgingInvoice, 0),
	}
	bucketOrder := []string{"Current", "1-30", "31-60", "61-90", "90+"}
	bucketIdx := map[string]int{}
	for i, b := range bucketOrder {
		bucketIdx[b] = i
	}
	type bKey struct {
		Bucket   string
		Currency string
	}
	buckets := map[bKey]*dto.ArAgingBucketRow{}

	for rows.Next() {
		var inv dto.ArAgingInvoice
		if err := rows.Scan(&inv.InvoiceID, &inv.TenantID, &inv.TenantName, &inv.Amount, &inv.Currency, &inv.Status, &inv.DueDate, &inv.PaidAmount); err != nil {
			return nil, fmt.Errorf("report service: ar-aging: scan: %w", err)
		}
		inv.Outstanding = inv.Amount - inv.PaidAmount
		if inv.Outstanding <= 0 {
			continue
		}
		inv.DaysOverdue = int(math.Max(0, math.Floor(now.Sub(inv.DueDate).Hours()/24)))
		inv.Bucket = agingBucket(inv.DaysOverdue)
		res.Invoices = append(res.Invoices, inv)

		k := bKey{inv.Bucket, inv.Currency}
		if b, ok := buckets[k]; ok {
			b.InvoiceCount++
			b.Outstanding += inv.Outstanding
		} else {
			buckets[k] = &dto.ArAgingBucketRow{Bucket: inv.Bucket, Currency: inv.Currency, InvoiceCount: 1, Outstanding: inv.Outstanding}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("report service: ar-aging: rows: %w", err)
	}

	for _, b := range buckets {
		res.Buckets = append(res.Buckets, *b)
	}
	sort.Slice(res.Buckets, func(i, j int) bool {
		bi, bj := bucketIdx[res.Buckets[i].Bucket], bucketIdx[res.Buckets[j].Bucket]
		if bi != bj {
			return bi < bj
		}
		return res.Buckets[i].Currency < res.Buckets[j].Currency
	})

	return res, nil
}
