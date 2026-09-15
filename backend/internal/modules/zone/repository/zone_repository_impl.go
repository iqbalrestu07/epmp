package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/zone/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ZoneRepositoryImpl implements ZoneRepository using PostgreSQL.
type ZoneRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewZoneRepositoryImpl creates a new ZoneRepositoryImpl.
func NewZoneRepositoryImpl(db *pgxpool.Pool) *ZoneRepositoryImpl {
	return &ZoneRepositoryImpl{db: db}
}

// Ensure ZoneRepositoryImpl implements domain repository interface.
var _ ZoneRepository = (*ZoneRepositoryImpl)(nil)

func (r *ZoneRepositoryImpl) Save(ctx context.Context, e *entity.Zone) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO zones (organization_id, building_id, floor, name)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.BuildingId, e.Floor, e.Name,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE zones
		SET    building_id=$1, floor=$2, name=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.BuildingId, e.Floor, e.Name, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *ZoneRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Zone, error) {
	e := &entity.Zone{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, building_id, floor, name, deleted_at, created_at, updated_at
		FROM   zones
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.BuildingId, &e.Floor, &e.Name, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("zone repository: find by id: %w", err)
	}
	return e, nil
}

func (r *ZoneRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Zone, error) {
	query := `
		SELECT organization_id, id, building_id, floor, name, deleted_at, created_at, updated_at
		FROM   zones
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
		return nil, fmt.Errorf("zone repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Zone
	for rows.Next() {
		e := &entity.Zone{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.BuildingId, &e.Floor, &e.Name, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *ZoneRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM zones WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("zone repository: count: %w", err)
	}
	return count, nil
}

func (r *ZoneRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE zones SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
