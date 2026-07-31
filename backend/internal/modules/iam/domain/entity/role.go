package entity

import "time"

// Role is the domain entity for an RBAC role.
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	IsSystem    bool         `json:"is_system"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// Permission is the domain entity for a granular permission (resource:action).
type Permission struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Key returns the canonical permission key used in authorization checks.
// Format: "resource:action", e.g. "property:read".
func (p Permission) Key() string {
	return p.Resource + ":" + p.Action
}
