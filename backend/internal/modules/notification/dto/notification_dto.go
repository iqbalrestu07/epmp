package dto

import "time"

// NotificationResponse is the DTO for one in-app notification.
type NotificationResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id,omitempty"`
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	Link           string    `json:"link,omitempty"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

// NotificationListResponse is the DTO for a paginated list of notifications.
type NotificationListResponse struct {
	Data        []NotificationResponse `json:"data"`
	Total       int64                  `json:"total"`
	UnreadCount int64                  `json:"unread_count"`
	Page        int                    `json:"page"`
	PerPage     int                    `json:"per_page"`
	TotalPages  int                    `json:"total_pages"`
}
