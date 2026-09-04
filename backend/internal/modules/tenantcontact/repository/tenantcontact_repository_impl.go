package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/tenantcontact/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantContactRepositoryImpl implements TenantContactRepository using PostgreSQL.
type TenantContactRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewTenantContactRepositoryImpl creates a new TenantContactRepositoryImpl.
func NewTenantContactRepositoryImpl(db *pgxpool.Pool) *TenantContactRepositoryImpl {
	return &TenantContactRepositoryImpl{db: db}
}

// Ensure TenantContactRepositoryImpl implements domain repository interface.
var _ TenantContactRepository = (*TenantContactRepositoryImpl)(nil)

func (r *TenantContactRepositoryImpl) Save(ctx context.Context, e *entity.TenantContact) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO tenant_contacts (organization_id, tenant_id, contact_type, contact_value, is_primary)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.TenantId, e.ContactType, e.ContactValue, e.IsPrimary,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE tenant_contacts
		SET    tenant_id=$1, contact_type=$2, contact_value=$3, is_primary=$4
		WHERE  id=$5 AND organization_id=$6 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.TenantId, e.ContactType, e.ContactValue, e.IsPrimary, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *TenantContactRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.TenantContact, error) {
	e := &entity.TenantContact{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, tenant_id, contact_type, contact_value, is_primary, deleted_at, created_at, updated_at
		FROM   tenant_contacts
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.ContactType, &e.ContactValue, &e.IsPrimary, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("tenantcontact repository: find by id: %w", err)
	}
	return e, nil
}

func (r *TenantContactRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantContact, error) {
	query := `
		SELECT organization_id, id, tenant_id, contact_type, contact_value, is_primary, deleted_at, created_at, updated_at
		FROM   tenant_contacts
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (contact_type ILIKE $%d OR contact_value ILIKE $%d)", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tenantcontact repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.TenantContact
	for rows.Next() {
		e := &entity.TenantContact{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.ContactType, &e.ContactValue, &e.IsPrimary, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *TenantContactRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM tenant_contacts WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (contact_type ILIKE $%d OR contact_value ILIKE $%d)", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("tenantcontact repository: count: %w", err)
	}
	return count, nil
}

func (r *TenantContactRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tenant_contacts SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
