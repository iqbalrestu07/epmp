package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/technician/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TechnicianRepositoryImpl implements TechnicianRepository using PostgreSQL.
type TechnicianRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewTechnicianRepositoryImpl creates a new TechnicianRepositoryImpl.
func NewTechnicianRepositoryImpl(db *pgxpool.Pool) *TechnicianRepositoryImpl {
	return &TechnicianRepositoryImpl{db: db}
}

// Ensure TechnicianRepositoryImpl implements domain repository interface.
var _ TechnicianRepository = (*TechnicianRepositoryImpl)(nil)

func (r *TechnicianRepositoryImpl) Save(ctx context.Context, e *entity.Technician) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO technicians (organization_id, name, phone, specialty)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.Name, e.Phone, e.Specialty,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE technicians
		SET    name=$1, phone=$2, specialty=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.Name, e.Phone, e.Specialty, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *TechnicianRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Technician, error) {
	e := &entity.Technician{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, name, phone, specialty, deleted_at, created_at, updated_at
		FROM   technicians
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.Name, &e.Phone, &e.Specialty, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("technician repository: find by id: %w", err)
	}
	return e, nil
}

func (r *TechnicianRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Technician, error) {
	query := `
		SELECT organization_id, id, name, phone, specialty, deleted_at, created_at, updated_at
		FROM   technicians
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR specialty ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("technician repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Technician
	for rows.Next() {
		e := &entity.Technician{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.Name, &e.Phone, &e.Specialty, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *TechnicianRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM technicians WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR specialty ILIKE $%d)", 2, 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("technician repository: count: %w", err)
	}
	return count, nil
}

func (r *TechnicianRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE technicians SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
