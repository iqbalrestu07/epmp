package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/epmp/backend/internal/modules/notification/dto"
	"github.com/epmp/backend/internal/pkg/email"
	"github.com/epmp/backend/internal/pkg/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// NotificationService stores notifications and pushes them over websocket.
type NotificationService struct {
	db      *pgxpool.Pool
	hub     *websocket.Hub
	emailer email.Sender
	log     zerolog.Logger
}

// NewNotificationService creates a new NotificationService.
// emailer may be nil — the email channel is then disabled.
func NewNotificationService(db *pgxpool.Pool, hub *websocket.Hub, emailer email.Sender, log zerolog.Logger) *NotificationService {
	return &NotificationService{db: db, hub: hub, emailer: emailer, log: log}
}

// NotifyOrg persists an org-wide notification and broadcasts it to connected clients.
func (s *NotificationService) NotifyOrg(ctx context.Context, orgID, notifType, title, message, link string) error {
	if orgID == "" {
		return fmt.Errorf("notification service: notify org: organization ID required")
	}

	var id string
	var createdAt time.Time
	err := s.db.QueryRow(ctx, `
		INSERT INTO notifications (organization_id, user_id, type, title, message, link)
		VALUES ($1, NULL, $2, $3, $4, NULLIF($5,''))
		RETURNING id, created_at`,
		orgID, notifType, title, message, link,
	).Scan(&id, &createdAt)
	if err != nil {
		return fmt.Errorf("notification service: notify org: %w", err)
	}

	if s.hub != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"event": "notification",
			"data": dto.NotificationResponse{
				ID: id, OrganizationID: orgID, Type: notifType,
				Title: title, Message: message, Link: link, CreatedAt: createdAt,
			},
		})
		s.hub.BroadcastToOrg(orgID, payload)
	}

	if s.emailer != nil {
		go s.emailOrgMembers(orgID, title, message)
	}
	return nil
}

// emailOrgMembers sends the notification to every active member of the organization.
// Errors are logged; email delivery is best-effort and never fails the request.
func (s *NotificationService) emailOrgMembers(orgID, title, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := s.db.Query(ctx, `
		SELECT u.email
		FROM   organization_members om
		JOIN   users u ON u.id = om.user_id
		WHERE  om.organization_id = $1
		  AND  om.is_active = true
		  AND  om.deleted_at IS NULL
		  AND  u.deleted_at IS NULL`, orgID)
	if err != nil {
		s.log.Error().Err(err).Str("org_id", orgID).Msg("notification email: member lookup failed")
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var e string
		if err := rows.Scan(&e); err == nil && e != "" {
			emails = append(emails, e)
		}
	}

	if len(emails) == 0 {
		return
	}
	if err := s.emailer.Send(ctx, emails, "[EPMP] "+title, message); err != nil {
		s.log.Error().Err(err).Str("org_id", orgID).Msg("notification email: send failed")
	}
}

// List returns org-scoped notifications visible to a user (user-targeted or org-wide).
func (s *NotificationService) List(ctx context.Context, userID, orgID string, page, perPage int) (*dto.NotificationListResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("notification service: list: organization ID required")
	}
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var total, unread int64
	if err := s.db.QueryRow(ctx, `
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE is_read = false)
		FROM   notifications
		WHERE  organization_id = $1 AND (user_id IS NULL OR user_id = $2)`,
		orgID, userID).Scan(&total, &unread); err != nil {
		return nil, fmt.Errorf("notification service: list: count: %w", err)
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, organization_id, COALESCE(user_id::text,''), type, title, message,
		       COALESCE(link,''), is_read, created_at
		FROM   notifications
		WHERE  organization_id = $1 AND (user_id IS NULL OR user_id = $2)
		ORDER  BY created_at DESC
		LIMIT  $3 OFFSET $4`, orgID, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("notification service: list: %w", err)
	}
	defer rows.Close()

	data := make([]dto.NotificationResponse, 0)
	for rows.Next() {
		var n dto.NotificationResponse
		if err := rows.Scan(&n.ID, &n.OrganizationID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Link, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("notification service: scan: %w", err)
		}
		data = append(data, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification service: rows: %w", err)
	}

	return &dto.NotificationListResponse{
		Data:        data,
		Total:       total,
		UnreadCount: unread,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  int(math.Ceil(float64(total) / float64(perPage))),
	}, nil
}

// MarkRead marks one notification as read within an organization.
func (s *NotificationService) MarkRead(ctx context.Context, id, orgID string) error {
	tag, err := s.db.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE id = $1 AND organization_id = $2`, id, orgID)
	if err != nil {
		return fmt.Errorf("notification service: mark read: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification service: mark read: not found")
	}
	return nil
}

// MarkAllRead marks all visible notifications as read for a user in an organization.
func (s *NotificationService) MarkAllRead(ctx context.Context, userID, orgID string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE notifications SET is_read = true
		WHERE  organization_id = $1 AND (user_id IS NULL OR user_id = $2) AND is_read = false`, orgID, userID)
	if err != nil {
		return fmt.Errorf("notification service: mark all read: %w", err)
	}
	return nil
}
