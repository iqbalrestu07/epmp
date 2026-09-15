package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/facility/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FacilityRepositoryImpl implements FacilityRepository using PostgreSQL.
type FacilityRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewFacilityRepositoryImpl creates a new FacilityRepositoryImpl.
func NewFacilityRepositoryImpl(db *pgxpool.Pool) *FacilityRepositoryImpl {
	return &FacilityRepositoryImpl{db: db}
}

// Ensure FacilityRepositoryImpl implements domain repository interface.
var _ FacilityRepository = (*FacilityRepositoryImpl)(nil)

func (r *FacilityRepositoryImpl) Save(ctx context.Context, e *entity.Facility) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO facilities (organization_id, property_id, name, description)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.PropertyId, e.Name, e.Description,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE facilities
		SET    property_id=$1, name=$2, description=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.PropertyId, e.Name, e.Description, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *FacilityRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Facility, error) {
	e := &entity.Facility{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, property_id, name, description, deleted_at, created_at, updated_at
		FROM   facilities
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &e.Name, &e.Description, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("facility repository: find by id: %w", err)
	}
	return e, nil
}

func (r *FacilityRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Facility, error) {
	query := `
		SELECT organization_id, id, property_id, name, description, deleted_at, created_at, updated_at
		FROM   facilities
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("facility repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Facility
	for rows.Next() {
		e := &entity.Facility{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &e.Name, &e.Description, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *FacilityRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM facilities WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("facility repository: count: %w", err)
	}
	return count, nil
}

func (r *FacilityRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE facilities SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
