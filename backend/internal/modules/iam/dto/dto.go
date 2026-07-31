package dto

// ─── Auth ────────────────────────────────────────────────────────────────────

// RegisterRequest is the payload for POST /auth/register.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest is the payload for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest is the payload for POST /auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse is the response after a successful login or refresh.
type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"` // seconds
	User         UserResponse `json:"user"`
}

// ─── User ────────────────────────────────────────────────────────────────────

// CreateUserRequest is the payload for POST /users.
type CreateUserRequest struct {
	Name     string   `json:"name" validate:"required,min=2"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	RoleIDs  []string `json:"role_ids"`
}

// UpdateUserRequest is the payload for PUT /users/:id.
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	IsActive *bool  `json:"is_active"`
}

// AssignRolesRequest is the payload for PUT /users/:id/roles.
type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required"`
}

// UserResponse is the standard user representation returned by the API.
type UserResponse struct {
	ID          string         `json:"id"`
	Email       string         `json:"email"`
	Name        string         `json:"name"`
	IsActive    bool           `json:"is_active"`
	Roles       []RoleResponse `json:"roles"`
	Permissions []string       `json:"permissions"` // flattened "resource:action" strings
}

// UserListResponse wraps a paginated list of users.
type UserListResponse struct {
	Data    []UserResponse `json:"data"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

// ─── Role ────────────────────────────────────────────────────────────────────

// CreateRoleRequest is the payload for POST /roles.
type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required,min=2"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

// UpdateRoleRequest is the payload for PUT /roles/:id.
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
}

// SetPermissionsRequest is the payload for PUT /roles/:id/permissions.
type SetPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" validate:"required"`
}

// RoleResponse is the standard role representation.
type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	IsSystem    bool                 `json:"is_system"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	UserCount   int                  `json:"user_count,omitempty"`
}

// ─── Permission ──────────────────────────────────────────────────────────────

// PermissionResponse is the standard permission representation.
type PermissionResponse struct {
	ID          string `json:"id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Key         string `json:"key"` // "resource:action"
}
