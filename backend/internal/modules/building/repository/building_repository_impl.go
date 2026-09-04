package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/building/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BuildingRepositoryImpl implements BuildingRepository using PostgreSQL.
type BuildingRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewBuildingRepositoryImpl creates a new BuildingRepositoryImpl.
func NewBuildingRepositoryImpl(db *pgxpool.Pool) *BuildingRepositoryImpl {
	return &BuildingRepositoryImpl{db: db}
}

// Ensure BuildingRepositoryImpl implements domain repository interface.
var _ BuildingRepository = (*BuildingRepositoryImpl)(nil)

func (r *BuildingRepositoryImpl) Save(ctx context.Context, e *entity.Building) error {
	if e.Id == "" {
		// INSERT
		err := r.db.QueryRow(ctx, `
			INSERT INTO buildings (property_id, name, total_floors)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, updated_at`,
			e.PropertyId, e.Name, e.TotalFloors,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE buildings
		SET    property_id=$1, name=$2, total_floors=$3, updated_at=now()
		WHERE  id=$4 AND deleted_at IS NULL`,
		e.PropertyId, e.Name, e.TotalFloors, e.Id,
	)
	return err
}

func (r *BuildingRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Building, error) {
	e := &entity.Building{}
	err := r.db.QueryRow(ctx, `
		SELECT id, property_id, name, total_floors, created_at, updated_at, deleted_at
		FROM   buildings
		WHERE  id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&e.Id, &e.PropertyId, &e.Name, &e.TotalFloors, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)

	if err != nil {
		return nil, fmt.Errorf("building repository: find by id: %w", err)
	}
	return e, nil
}

func (r *BuildingRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string, propertyId string) ([]*entity.Building, error) {
	query := `
		SELECT id, property_id, name, total_floors, created_at, updated_at, deleted_at
		FROM   buildings
		WHERE  deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if propertyId != "" {
		query += fmt.Sprintf(` AND property_id = $%d`, argIdx)
		args = append(args, propertyId)
		argIdx++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("building repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Building
	for rows.Next() {
		e := &entity.Building{}
		if err := rows.Scan(&e.Id, &e.PropertyId, &e.Name, &e.TotalFloors, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *BuildingRepositoryImpl) Count(ctx context.Context, search string, propertyId string) (int64, error) {
	query := `SELECT COUNT(*) FROM buildings WHERE deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if propertyId != "" {
		query += fmt.Sprintf(` AND property_id = $%d`, argIdx)
		args = append(args, propertyId)
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("building repository: count: %w", err)
	}
	return count, nil
}

func (r *BuildingRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE buildings SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}
