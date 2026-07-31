package repository

import (
	"context"

	"github.com/epmp/backend/internal/modules/iam/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgPermissionRepository struct {
	db *pgxpool.Pool
}

// NewPermissionRepository creates a new PostgreSQL-backed PermissionRepository.
func NewPermissionRepository(db *pgxpool.Pool) *pgPermissionRepository {
	return &pgPermissionRepository{db: db}
}

func (r *pgPermissionRepository) FindAll(ctx context.Context) ([]*entity.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, resource, action, description, created_at
		FROM   permissions
		ORDER  BY resource, action`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		p := &entity.Permission{}
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

func (r *pgPermissionRepository) FindByID(ctx context.Context, id string) (*entity.Permission, error) {
	p := &entity.Permission{}
	err := r.db.QueryRow(ctx, `
		SELECT id, resource, action, description, created_at
		FROM   permissions WHERE id = $1`, id,
	).Scan(&p.ID, &p.Resource, &p.Action, &p.Description, &p.CreatedAt)
	return p, err
}

// ─────────────────────────────────────────────────────────────────────────────

type pgUserRoleRepository struct {
	db *pgxpool.Pool
}

// NewUserRoleRepository creates a new PostgreSQL-backed UserRoleRepository.
func NewUserRoleRepository(db *pgxpool.Pool) *pgUserRoleRepository {
	return &pgUserRoleRepository{db: db}
}

func (r *pgUserRoleRepository) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Clear existing assignments
	if _, err := tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID); err != nil {
		return err
	}

	for _, rid := range roleIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO user_roles (user_id, role_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, rid,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgUserRoleRepository) FindRolesByUserID(ctx context.Context, userID string) ([]*entity.Role, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ro.id, ro.name, ro.description, ro.is_system, ro.created_at, ro.updated_at
		FROM   roles ro
		JOIN   user_roles ur ON ur.role_id = ro.id
		WHERE  ur.user_id = $1
		ORDER  BY ro.name`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		ro := &entity.Role{}
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.Description, &ro.IsSystem, &ro.CreatedAt, &ro.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, ro)
	}
	return roles, rows.Err()
}

func (r *pgUserRoleRepository) FindPermissionsByUserID(ctx context.Context, userID string) ([]*entity.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT p.id, p.resource, p.action, p.description, p.created_at
		FROM   permissions p
		JOIN   role_permissions rp ON rp.permission_id = p.id
		JOIN   user_roles ur       ON ur.role_id = rp.role_id
		WHERE  ur.user_id = $1
		ORDER  BY p.resource, p.action`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		p := &entity.Permission{}
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────

type pgRefreshTokenRepository struct {
	db *pgxpool.Pool
}

// NewRefreshTokenRepository creates a new PostgreSQL-backed RefreshTokenRepository.
func NewRefreshTokenRepository(db *pgxpool.Pool) *pgRefreshTokenRepository {
	return &pgRefreshTokenRepository{db: db}
}

func (r *pgRefreshTokenRepository) Save(ctx context.Context, rt *entity.RefreshToken) error {
	row := r.db.QueryRow(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		rt.UserID, rt.TokenHash, rt.ExpiresAt,
	)
	return row.Scan(&rt.ID, &rt.CreatedAt)
}

func (r *pgRefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	rt := &entity.RefreshToken{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM   refresh_tokens
		WHERE  token_hash = $1`, hash,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.RevokedAt, &rt.CreatedAt)
	return rt, err
}

func (r *pgRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = now() WHERE id = $1`, id)
	return err
}

func (r *pgRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE  user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}
