package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/room/entity"
	"github.com/epmp/backend/internal/pkg/uid"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomRepositoryImpl implements RoomRepository using PostgreSQL.
type RoomRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewRoomRepositoryImpl creates a new RoomRepositoryImpl.
func NewRoomRepositoryImpl(db *pgxpool.Pool) *RoomRepositoryImpl {
	return &RoomRepositoryImpl{db: db}
}

// Ensure RoomRepositoryImpl implements domain repository interface.
var _ RoomRepository = (*RoomRepositoryImpl)(nil)

func (r *RoomRepositoryImpl) Save(ctx context.Context, e *entity.Room) error {
	if e.Id == "" {
		return fmt.Errorf("room repository: save: id must be pre-set by caller (use uid.New())")
	}

	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM rooms WHERE id=$1)`, e.Id).Scan(&exists); err != nil {
		return fmt.Errorf("room repository: save: check exists: %w", err)
	}

	var floorID interface{}
	if e.FloorId != "" {
		floorID = e.FloorId
	}

	if !exists {
		// INSERT with caller-provided ULID
		err := r.db.QueryRow(ctx, `
			INSERT INTO rooms (id, organization_id, property_id, floor_id, name, capacity, price, is_available)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING created_at, updated_at`,
			e.Id, e.OrganizationId, e.PropertyId, floorID, e.Name, e.Capacity, e.Price, e.IsAvailable,
		).Scan(&e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE rooms
		SET    organization_id=$1, property_id=$2, floor_id=$3, name=$4, capacity=$5, price=$6, is_available=$7
		WHERE  id=$8 AND deleted_at IS NULL`,
		e.OrganizationId, e.PropertyId, floorID, e.Name, e.Capacity, e.Price, e.IsAvailable, e.Id,
	)
	return err
}

func (r *RoomRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Room, error) {
	e := &entity.Room{}
	var floorID *string
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, property_id, floor_id, name, capacity, price, is_available,
		       created_at, updated_at, deleted_at
		FROM   rooms
		WHERE  id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &floorID, &e.Name, &e.Capacity, &e.Price, &e.IsAvailable,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if floorID != nil {
		e.FloorId = *floorID
	}

	if err != nil {
		return nil, fmt.Errorf("room repository: find by id: %w", err)
	}
	return e, nil
}

func (r *RoomRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]*entity.Room, error) {
	rows, err := r.db.Query(ctx, `
		SELECT organization_id, id, property_id, floor_id, name, capacity, price, is_available,
		       created_at, updated_at, deleted_at
		FROM   rooms
		WHERE  deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("room repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Room
	for rows.Next() {
		e := &entity.Room{}
		var floorID *string
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &floorID, &e.Name, &e.Capacity, &e.Price, &e.IsAvailable,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		if floorID != nil {
			e.FloorId = *floorID
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *RoomRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE rooms SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

// newRoomID is a helper — kept for reference; call uid.New() directly in service.
var _ = uid.New
