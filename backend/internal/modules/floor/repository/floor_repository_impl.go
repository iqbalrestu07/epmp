package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/floor/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FloorRepositoryImpl implements FloorRepository using PostgreSQL.
type FloorRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewFloorRepositoryImpl creates a new FloorRepositoryImpl.
func NewFloorRepositoryImpl(db *pgxpool.Pool) *FloorRepositoryImpl {
	return &FloorRepositoryImpl{db: db}
}

// Ensure FloorRepositoryImpl implements domain repository interface.
var _ FloorRepository = (*FloorRepositoryImpl)(nil)

func (r *FloorRepositoryImpl) Save(ctx context.Context, e *entity.Floor) error {
	if e.Id == "" {
		// INSERT
		err := r.db.QueryRow(ctx, `
			INSERT INTO floors (organization_id, building_id, name, floor_number, is_active)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.BuildingId, e.Name, e.FloorNumber, e.IsActive,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE floors
		SET    organization_id=$1, building_id=$2, name=$3, floor_number=$4, is_active=$5, updated_at=now()
		WHERE  id=$6 AND deleted_at IS NULL`,
		e.OrganizationId, e.BuildingId, e.Name, e.FloorNumber, e.IsActive, e.Id,
	)
	return err
}

func (r *FloorRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Floor, error) {
	e := &entity.Floor{}
	query := `
		SELECT id, organization_id, building_id, name, floor_number, is_active, created_at, updated_at, deleted_at
		FROM   floors
		WHERE  id = $1 AND deleted_at IS NULL`
	args := []interface{}{id}
	if orgID != "" {
		query += ` AND organization_id = $2`
		args = append(args, orgID)
	}
	err := r.db.QueryRow(ctx, query, args...).Scan(&e.Id, &e.OrganizationId, &e.BuildingId, &e.Name, &e.FloorNumber, &e.IsActive, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)

	if err != nil {
		return nil, fmt.Errorf("floor repository: find by id: %w", err)
	}
	return e, nil
}

func (r *FloorRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string, buildingId, orgID string) ([]*entity.Floor, error) {
	query := `
		SELECT id, organization_id, building_id, name, floor_number, is_active, created_at, updated_at, deleted_at
		FROM   floors
		WHERE  deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if buildingId != "" {
		query += fmt.Sprintf(` AND building_id = $%d`, argIdx)
		args = append(args, buildingId)
		argIdx++
	}
	if orgID != "" {
		query += fmt.Sprintf(` AND organization_id = $%d`, argIdx)
		args = append(args, orgID)
		argIdx++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("floor repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Floor
	for rows.Next() {
		e := &entity.Floor{}
		if err := rows.Scan(&e.Id, &e.OrganizationId, &e.BuildingId, &e.Name, &e.FloorNumber, &e.IsActive, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *FloorRepositoryImpl) Count(ctx context.Context, search string, buildingId, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM floors WHERE deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if buildingId != "" {
		query += fmt.Sprintf(` AND building_id = $%d`, argIdx)
		args = append(args, buildingId)
		argIdx++
	}
	if orgID != "" {
		query += fmt.Sprintf(` AND organization_id = $%d`, argIdx)
		args = append(args, orgID)
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("floor repository: count: %w", err)
	}
	return count, nil
}

func (r *FloorRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE floors SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}
