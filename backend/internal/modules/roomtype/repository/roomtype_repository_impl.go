package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/roomtype/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomTypeRepositoryImpl implements RoomTypeRepository using PostgreSQL.
type RoomTypeRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewRoomTypeRepositoryImpl creates a new RoomTypeRepositoryImpl.
func NewRoomTypeRepositoryImpl(db *pgxpool.Pool) *RoomTypeRepositoryImpl {
	return &RoomTypeRepositoryImpl{db: db}
}

// Ensure RoomTypeRepositoryImpl implements domain repository interface.
var _ RoomTypeRepository = (*RoomTypeRepositoryImpl)(nil)

func (r *RoomTypeRepositoryImpl) Save(ctx context.Context, e *entity.RoomType) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO room_types (organization_id, name, description, base_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.Name, e.Description, e.BasePrice,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE room_types
		SET    name=$1, description=$2, base_price=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.Name, e.Description, e.BasePrice, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *RoomTypeRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.RoomType, error) {
	e := &entity.RoomType{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, name, description, base_price, deleted_at, created_at, updated_at
		FROM   room_types
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.Name, &e.Description, &e.BasePrice, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("roomtype repository: find by id: %w", err)
	}
	return e, nil
}

func (r *RoomTypeRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.RoomType, error) {
	query := `
		SELECT organization_id, id, name, description, base_price, deleted_at, created_at, updated_at
		FROM   room_types
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("roomtype repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.RoomType
	for rows.Next() {
		e := &entity.RoomType{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.Name, &e.Description, &e.BasePrice, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *RoomTypeRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM room_types WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("roomtype repository: count: %w", err)
	}
	return count, nil
}

func (r *RoomTypeRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE room_types SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
