package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/iam/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRoleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new PostgreSQL-backed RoleRepository.
func NewRoleRepository(db *pgxpool.Pool) *pgRoleRepository {
	return &pgRoleRepository{db: db}
}

func (r *pgRoleRepository) Save(ctx context.Context, role *entity.Role) error {
	if role.ID == "" {
		row := r.db.QueryRow(ctx, `
			INSERT INTO roles (name, description, is_system)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, updated_at`,
			role.Name, role.Description, role.IsSystem,
		)
		return row.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	}
	_, err := r.db.Exec(ctx, `
		UPDATE roles SET name=$1, description=$2, updated_at=now() WHERE id=$3`,
		role.Name, role.Description, role.ID,
	)
	return err
}

func (r *pgRoleRepository) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	role := &entity.Role{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, is_system, created_at, updated_at
		FROM   roles WHERE id = $1`, id,
	).Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	return role, err
}

func (r *pgRoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	role := &entity.Role{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, is_system, created_at, updated_at
		FROM   roles WHERE name = $1`, name,
	).Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	return role, err
}

func (r *pgRoleRepository) FindAll(ctx context.Context) ([]*entity.Role, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, is_system, created_at, updated_at
		FROM   roles ORDER BY is_system DESC, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role := &entity.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *pgRoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	return err
}

func (r *pgRoleRepository) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Remove all existing
	if _, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}

	// Re-insert
	for _, pid := range permissionIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`, roleID, pid,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgRoleRepository) FindPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.resource, p.action, p.description, p.created_at
		FROM   permissions p
		JOIN   role_permissions rp ON rp.permission_id = p.id
		WHERE  rp.role_id = $1
		ORDER  BY p.resource, p.action`, roleID,
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
