package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/bed/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BedRepositoryImpl implements BedRepository using PostgreSQL.
type BedRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewBedRepositoryImpl creates a new BedRepositoryImpl.
func NewBedRepositoryImpl(db *pgxpool.Pool) *BedRepositoryImpl {
	return &BedRepositoryImpl{db: db}
}

// Ensure BedRepositoryImpl implements domain repository interface.
var _ BedRepository = (*BedRepositoryImpl)(nil)

func (r *BedRepositoryImpl) Save(ctx context.Context, e *entity.Bed) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO beds (organization_id, room_id, name, status)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.RoomId, e.Name, e.Status,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE beds
		SET    room_id=$1, name=$2, status=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.RoomId, e.Name, e.Status, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *BedRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Bed, error) {
	e := &entity.Bed{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, room_id, name, status, deleted_at, created_at, updated_at
		FROM   beds
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.RoomId, &e.Name, &e.Status, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("bed repository: find by id: %w", err)
	}
	return e, nil
}

func (r *BedRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Bed, error) {
	query := `
		SELECT organization_id, id, room_id, name, status, deleted_at, created_at, updated_at
		FROM   beds
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR status ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("bed repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Bed
	for rows.Next() {
		e := &entity.Bed{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.RoomId, &e.Name, &e.Status, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *BedRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM beds WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR status ILIKE $%d)", 2, 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("bed repository: count: %w", err)
	}
	return count, nil
}

func (r *BedRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE beds SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
