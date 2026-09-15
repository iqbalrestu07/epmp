package dto

import (
	"encoding/json"
	"time"
)

// AuditLogResponse is the DTO for one audit log entry.
type AuditLogResponse struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	UserID         string          `json:"user_id"`
	UserEmail      string          `json:"user_email"`
	Action         string          `json:"action"`
	Module         string          `json:"module"`
	EntityID       string          `json:"entity_id"`
	Method         string          `json:"method"`
	Path           string          `json:"path"`
	StatusCode     int             `json:"status_code"`
	IPAddress      string          `json:"ip_address"`
	UserAgent      string          `json:"user_agent"`
	RequestBody    json.RawMessage `json:"request_body,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// AuditLogListResponse is the DTO for a paginated list of audit logs.
type AuditLogListResponse struct {
	Data       []AuditLogResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}
