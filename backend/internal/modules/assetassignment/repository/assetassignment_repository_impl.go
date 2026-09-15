package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/assetassignment/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AssetAssignmentRepositoryImpl implements AssetAssignmentRepository using PostgreSQL.
type AssetAssignmentRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewAssetAssignmentRepositoryImpl creates a new AssetAssignmentRepositoryImpl.
func NewAssetAssignmentRepositoryImpl(db *pgxpool.Pool) *AssetAssignmentRepositoryImpl {
	return &AssetAssignmentRepositoryImpl{db: db}
}

// Ensure AssetAssignmentRepositoryImpl implements domain repository interface.
var _ AssetAssignmentRepository = (*AssetAssignmentRepositoryImpl)(nil)

func (r *AssetAssignmentRepositoryImpl) Save(ctx context.Context, e *entity.AssetAssignment) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO asset_assignments (organization_id, asset_id, room_id, assigned_date)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.AssetId, e.RoomId, e.AssignedDate,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE asset_assignments
		SET    asset_id=$1, room_id=$2, assigned_date=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.AssetId, e.RoomId, e.AssignedDate, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *AssetAssignmentRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.AssetAssignment, error) {
	e := &entity.AssetAssignment{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, asset_id, room_id, assigned_date, deleted_at, created_at, updated_at
		FROM   asset_assignments
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.AssetId, &e.RoomId, &e.AssignedDate, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("assetassignment repository: find by id: %w", err)
	}
	return e, nil
}

func (r *AssetAssignmentRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.AssetAssignment, error) {
	query := `
		SELECT organization_id, id, asset_id, room_id, assigned_date, deleted_at, created_at, updated_at
		FROM   asset_assignments
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND condition ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("assetassignment repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.AssetAssignment
	for rows.Next() {
		e := &entity.AssetAssignment{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.AssetId, &e.RoomId, &e.AssignedDate, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *AssetAssignmentRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM asset_assignments WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND condition ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("assetassignment repository: count: %w", err)
	}
	return count, nil
}

func (r *AssetAssignmentRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE asset_assignments SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
