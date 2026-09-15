package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/tenant/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantRepositoryImpl implements TenantRepository using PostgreSQL.
type TenantRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewTenantRepositoryImpl creates a new TenantRepositoryImpl.
func NewTenantRepositoryImpl(db *pgxpool.Pool) *TenantRepositoryImpl {
	return &TenantRepositoryImpl{db: db}
}

// Ensure TenantRepositoryImpl implements domain repository interface.
var _ TenantRepository = (*TenantRepositoryImpl)(nil)

func (r *TenantRepositoryImpl) Save(ctx context.Context, e *entity.Tenant) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO tenants (organization_id, full_name, email, phone, identity_number, is_active)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.FullName, e.Email, e.Phone, e.IdentityNumber, e.IsActive,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE tenants
		SET    full_name=$1, email=$2, phone=$3, identity_number=$4, is_active=$5
		WHERE  id=$6 AND organization_id=$7 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.FullName, e.Email, e.Phone, e.IdentityNumber, e.IsActive, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *TenantRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Tenant, error) {
	e := &entity.Tenant{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, full_name, email, phone, identity_number, is_active, deleted_at, created_at, updated_at
		FROM   tenants
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.FullName, &e.Email, &e.Phone, &e.IdentityNumber, &e.IsActive, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("tenant repository: find by id: %w", err)
	}
	return e, nil
}

func (r *TenantRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Tenant, error) {
	query := `
		SELECT organization_id, id, full_name, email, phone, identity_number, is_active, deleted_at, created_at, updated_at
		FROM   tenants
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (full_name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tenant repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Tenant
	for rows.Next() {
		e := &entity.Tenant{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.FullName, &e.Email, &e.Phone, &e.IdentityNumber, &e.IsActive, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *TenantRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += " AND (full_name ILIKE $2 OR email ILIKE $2 OR phone ILIKE $2)"
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("tenant repository: count: %w", err)
	}
	return count, nil
}

func (r *TenantRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE tenants SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
