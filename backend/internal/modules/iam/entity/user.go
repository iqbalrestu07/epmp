package entity

import "time"

// User is the domain entity for an application user.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name           string     `json:"name"`
	OrganizationID *string    `json:"organization_id,omitempty"`
	IsActive       bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// Loaded relations (not stored directly on users table)
	Roles []Role `json:"roles,omitempty"`
}
