package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/tenantidentity/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantIdentityRepositoryImpl implements TenantIdentityRepository using PostgreSQL.
type TenantIdentityRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewTenantIdentityRepositoryImpl creates a new TenantIdentityRepositoryImpl.
func NewTenantIdentityRepositoryImpl(db *pgxpool.Pool) *TenantIdentityRepositoryImpl {
	return &TenantIdentityRepositoryImpl{db: db}
}

// Ensure TenantIdentityRepositoryImpl implements domain repository interface.
var _ TenantIdentityRepository = (*TenantIdentityRepositoryImpl)(nil)

func (r *TenantIdentityRepositoryImpl) Save(ctx context.Context, e *entity.TenantIdentity) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO tenant_identities (organization_id, tenant_id, identity_type, identity_number, file_url)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.TenantId, e.IdentityType, e.IdentityNumber, e.FileUrl,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE tenant_identities
		SET    tenant_id=$1, identity_type=$2, identity_number=$3, file_url=$4
		WHERE  id=$5 AND organization_id=$6 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.TenantId, e.IdentityType, e.IdentityNumber, e.FileUrl, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *TenantIdentityRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.TenantIdentity, error) {
	e := &entity.TenantIdentity{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, tenant_id, identity_type, identity_number, file_url, deleted_at, created_at, updated_at
		FROM   tenant_identities
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.IdentityType, &e.IdentityNumber, &e.FileUrl, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("tenantidentity repository: find by id: %w", err)
	}
	return e, nil
}

func (r *TenantIdentityRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantIdentity, error) {
	query := `
		SELECT organization_id, id, tenant_id, identity_type, identity_number, file_url, deleted_at, created_at, updated_at
		FROM   tenant_identities
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (identity_type ILIKE $%d OR identity_number ILIKE $%d)", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tenantidentity repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.TenantIdentity
	for rows.Next() {
		e := &entity.TenantIdentity{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.IdentityType, &e.IdentityNumber, &e.FileUrl, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *TenantIdentityRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM tenant_identities WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (identity_type ILIKE $%d OR identity_number ILIKE $%d)", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("tenantidentity repository: count: %w", err)
	}
	return count, nil
}

func (r *TenantIdentityRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tenant_identities SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
