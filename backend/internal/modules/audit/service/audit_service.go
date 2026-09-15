package service

import (
	"context"
	"fmt"
	"math"

	"github.com/epmp/backend/internal/modules/audit/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditService provides read-only access to the audit trail.
type AuditService struct {
	db *pgxpool.Pool
}

// NewAuditService creates a new AuditService.
func NewAuditService(db *pgxpool.Pool) *AuditService {
	return &AuditService{db: db}
}

// List returns org-scoped, paginated audit logs.
// search matches module, user_email, path, or entity_id.
func (s *AuditService) List(ctx context.Context, page, perPage int, search, module, action, orgID string) (*dto.AuditLogListResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("audit service: organization ID required")
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	where := "organization_id = $1"
	args := []interface{}{orgID}
	idx := 2

	if search != "" {
		where += fmt.Sprintf(" AND (module ILIKE $%d OR user_email ILIKE $%d OR path ILIKE $%d OR entity_id ILIKE $%d)", idx, idx, idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}
	if module != "" {
		where += fmt.Sprintf(" AND module = $%d", idx)
		args = append(args, module)
		idx++
	}
	if action != "" {
		where += fmt.Sprintf(" AND action = $%d", idx)
		args = append(args, action)
		idx++
	}

	var total int64
	if err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM audit_logs WHERE "+where, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("audit service: count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, organization_id, COALESCE(user_id::text,''), COALESCE(user_email,''),
		       action, module, COALESCE(entity_id,''), method, path, status_code,
		       COALESCE(ip_address,''), COALESCE(user_agent,''), request_body, created_at
		FROM   audit_logs
		WHERE  %s
		ORDER  BY created_at DESC
		LIMIT  $%d OFFSET $%d`, where, idx, idx+1)
	args = append(args, perPage, offset)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("audit service: list: %w", err)
	}
	defer rows.Close()

	data := make([]dto.AuditLogResponse, 0)
	for rows.Next() {
		var r dto.AuditLogResponse
		if err := rows.Scan(&r.ID, &r.OrganizationID, &r.UserID, &r.UserEmail,
			&r.Action, &r.Module, &r.EntityID, &r.Method, &r.Path, &r.StatusCode,
			&r.IPAddress, &r.UserAgent, &r.RequestBody, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("audit service: scan: %w", err)
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("audit service: rows: %w", err)
	}

	return &dto.AuditLogListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	}, nil
}
