package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/assetinspection/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AssetInspectionRepositoryImpl implements AssetInspectionRepository using PostgreSQL.
type AssetInspectionRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewAssetInspectionRepositoryImpl creates a new AssetInspectionRepositoryImpl.
func NewAssetInspectionRepositoryImpl(db *pgxpool.Pool) *AssetInspectionRepositoryImpl {
	return &AssetInspectionRepositoryImpl{db: db}
}

// Ensure AssetInspectionRepositoryImpl implements domain repository interface.
var _ AssetInspectionRepository = (*AssetInspectionRepositoryImpl)(nil)

func (r *AssetInspectionRepositoryImpl) Save(ctx context.Context, e *entity.AssetInspection) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO asset_inspections (organization_id, asset_id, inspection_date, condition, notes)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.AssetId, e.InspectionDate, e.Condition, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE asset_inspections
		SET    asset_id=$1, inspection_date=$2, condition=$3, notes=$4
		WHERE  id=$5 AND organization_id=$6 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.AssetId, e.InspectionDate, e.Condition, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *AssetInspectionRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.AssetInspection, error) {
	e := &entity.AssetInspection{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, asset_id, inspection_date, condition, notes, deleted_at, created_at, updated_at
		FROM   asset_inspections
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.AssetId, &e.InspectionDate, &e.Condition, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("assetinspection repository: find by id: %w", err)
	}
	return e, nil
}

func (r *AssetInspectionRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.AssetInspection, error) {
	query := `
		SELECT organization_id, id, asset_id, inspection_date, condition, notes, deleted_at, created_at, updated_at
		FROM   asset_inspections
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND condition ILIKE $%d", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("assetinspection repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.AssetInspection
	for rows.Next() {
		e := &entity.AssetInspection{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.AssetId, &e.InspectionDate, &e.Condition, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *AssetInspectionRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM asset_inspections WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND condition ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("assetinspection repository: count: %w", err)
	}
	return count, nil
}

func (r *AssetInspectionRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE asset_inspections SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
