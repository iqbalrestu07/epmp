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
	var orgID interface{}
	if e.OrganizationId != "" {
		orgID = e.OrganizationId
	}

	if e.Id == "" {
		// INSERT
		err := r.db.QueryRow(ctx, `
			INSERT INTO buildings (organization_id, property_id, name, total_floors)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			orgID, e.PropertyId, e.Name, e.TotalFloors,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE buildings
		SET    organization_id=$1, property_id=$2, name=$3, total_floors=$4, updated_at=now()
		WHERE  id=$5 AND deleted_at IS NULL`,
		orgID, e.PropertyId, e.Name, e.TotalFloors, e.Id,
	)
	return err
}

func (r *BuildingRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Building, error) {
	e := &entity.Building{}
	var orgIDPtr *string
	query := `
		SELECT organization_id, id, property_id, name, total_floors, created_at, updated_at, deleted_at
		FROM   buildings
		WHERE  id = $1 AND deleted_at IS NULL`
	args := []interface{}{id}
	if orgID != "" {
		query += ` AND organization_id = $2`
		args = append(args, orgID)
	}
	err := r.db.QueryRow(ctx, query, args...).Scan(&orgIDPtr, &e.Id, &e.PropertyId, &e.Name, &e.TotalFloors, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if orgIDPtr != nil {
		e.OrganizationId = *orgIDPtr
	}

	if err != nil {
		return nil, fmt.Errorf("building repository: find by id: %w", err)
	}
	return e, nil
}

func (r *BuildingRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string, propertyId, orgID string) ([]*entity.Building, error) {
	query := `
		SELECT organization_id, id, property_id, name, total_floors, created_at, updated_at, deleted_at
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
	if orgID != "" {
		query += fmt.Sprintf(` AND organization_id = $%d`, argIdx)
		args = append(args, orgID)
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
		var orgIDPtr *string
		if err := rows.Scan(&orgIDPtr, &e.Id, &e.PropertyId, &e.Name, &e.TotalFloors, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		if orgIDPtr != nil {
			e.OrganizationId = *orgIDPtr
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *BuildingRepositoryImpl) Count(ctx context.Context, search string, propertyId, orgID string) (int64, error) {
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
		argIdx++
	}
	if orgID != "" {
		query += fmt.Sprintf(` AND organization_id = $%d`, argIdx)
		args = append(args, orgID)
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
