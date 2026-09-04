package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/supplier/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SupplierRepositoryImpl implements SupplierRepository using PostgreSQL.
type SupplierRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewSupplierRepositoryImpl creates a new SupplierRepositoryImpl.
func NewSupplierRepositoryImpl(db *pgxpool.Pool) *SupplierRepositoryImpl {
	return &SupplierRepositoryImpl{db: db}
}

// Ensure SupplierRepositoryImpl implements domain repository interface.
var _ SupplierRepository = (*SupplierRepositoryImpl)(nil)

func (r *SupplierRepositoryImpl) Save(ctx context.Context, e *entity.Supplier) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO vendors (organization_id, name, contact_person, phone, service_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.Name, e.ContactPerson, e.Phone, e.ServiceType,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE vendors
		SET    name=$1, contact_person=$2, phone=$3, service_type=$4
		WHERE  id=$5 AND organization_id=$6 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.Name, e.ContactPerson, e.Phone, e.ServiceType, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *SupplierRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Supplier, error) {
	e := &entity.Supplier{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, name, contact_person, phone, service_type, deleted_at, created_at, updated_at
		FROM   vendors
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.Name, &e.ContactPerson, &e.Phone, &e.ServiceType, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("supplier repository: find by id: %w", err)
	}
	return e, nil
}

func (r *SupplierRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Supplier, error) {
	query := `
		SELECT organization_id, id, name, contact_person, phone, service_type, deleted_at, created_at, updated_at
		FROM   vendors
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR service_type ILIKE $%d)", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("supplier repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Supplier
	for rows.Next() {
		e := &entity.Supplier{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.Name, &e.ContactPerson, &e.Phone, &e.ServiceType, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *SupplierRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM vendors WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR service_type ILIKE $%d)", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("supplier repository: count: %w", err)
	}
	return count, nil
}

func (r *SupplierRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE vendors SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
