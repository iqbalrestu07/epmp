package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/iam/entity"
)

// UserRepository is the contract for user persistence.
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
	CountAll(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id string) error
	UpdateLastLogin(ctx context.Context, id string) error
}

// RoleRepository is the contract for role persistence.
type RoleRepository interface {
	Save(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id string) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindAll(ctx context.Context) ([]*entity.Role, error)
	Delete(ctx context.Context, id string) error

	// Permission management within roles
	SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	FindPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)
}

// PermissionRepository is the contract for permission persistence.
type PermissionRepository interface {
	FindAll(ctx context.Context) ([]*entity.Permission, error)
	FindByID(ctx context.Context, id string) (*entity.Permission, error)
}

// UserRoleRepository manages user ↔ role assignments.
type UserRoleRepository interface {
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	FindRolesByUserID(ctx context.Context, userID string) ([]*entity.Role, error)
	FindPermissionsByUserID(ctx context.Context, userID string) ([]*entity.Permission, error)
}

// RefreshTokenRepository manages refresh token persistence.
type RefreshTokenRepository interface {
	Save(ctx context.Context, rt *entity.RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}
