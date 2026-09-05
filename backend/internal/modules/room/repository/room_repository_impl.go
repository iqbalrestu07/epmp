package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/room/entity"

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
	var roomTypeID interface{}
	if e.RoomTypeId != "" {
		roomTypeID = e.RoomTypeId
	}

	if !exists {
		// INSERT with caller-provided ULID
		err := r.db.QueryRow(ctx, `
			INSERT INTO rooms (id, organization_id, property_id, floor_id, room_type_id, name, capacity, price, is_available)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING created_at, updated_at`,
			e.Id, e.OrganizationId, e.PropertyId, floorID, roomTypeID, e.Name, e.Capacity, e.Price, e.IsAvailable,
		).Scan(&e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE rooms
		SET    organization_id=$1, property_id=$2, floor_id=$3, room_type_id=$4, name=$5, capacity=$6, price=$7, is_available=$8
		WHERE  id=$9 AND deleted_at IS NULL`,
		e.OrganizationId, e.PropertyId, floorID, roomTypeID, e.Name, e.Capacity, e.Price, e.IsAvailable, e.Id,
	)
	return err
}

func (r *RoomRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Room, error) {
	e := &entity.Room{}
	var floorID *string
	var roomTypeID *string
	query := `
		SELECT organization_id, id, property_id, floor_id, room_type_id, name, capacity, price, is_available,
		       created_at, updated_at, deleted_at
		FROM   rooms
		WHERE  id = $1 AND deleted_at IS NULL`
	args := []interface{}{id}
	if orgID != "" {
		query += ` AND organization_id = $2`
		args = append(args, orgID)
	}
	err := r.db.QueryRow(ctx, query, args...).Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &floorID, &roomTypeID, &e.Name, &e.Capacity, &e.Price, &e.IsAvailable,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if floorID != nil {
		e.FloorId = *floorID
	}
	if roomTypeID != nil {
		e.RoomTypeId = *roomTypeID
	}

	if err != nil {
		return nil, fmt.Errorf("room repository: find by id: %w", err)
	}
	return e, nil
}

func (r *RoomRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, floorId, propertyId, buildingId, orgID string) ([]*entity.Room, error) {
	query := `
		SELECT organization_id, id, property_id, floor_id, room_type_id, name, capacity, price, is_available,
		       created_at, updated_at, deleted_at
		FROM   rooms
		WHERE  deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if floorId != "" {
		query += fmt.Sprintf(` AND floor_id = $%d`, argIdx)
		args = append(args, floorId)
		argIdx++
	}
	if propertyId != "" {
		query += fmt.Sprintf(` AND property_id = $%d`, argIdx)
		args = append(args, propertyId)
		argIdx++
	}
	if buildingId != "" {
		query += fmt.Sprintf(` AND floor_id IN (SELECT id FROM floors WHERE building_id = $%d)`, argIdx)
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
		return nil, fmt.Errorf("room repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Room
	for rows.Next() {
		e := &entity.Room{}
		var floorID *string
		var roomTypeID *string
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &floorID, &roomTypeID, &e.Name, &e.Capacity, &e.Price, &e.IsAvailable,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		if floorID != nil {
			e.FloorId = *floorID
		}
		if roomTypeID != nil {
			e.RoomTypeId = *roomTypeID
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

// Count returns the total number of non-deleted rooms.
func (r *RoomRepositoryImpl) Count(ctx context.Context, search, floorId, propertyId, buildingId, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM rooms WHERE deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if floorId != "" {
		query += fmt.Sprintf(` AND floor_id = $%d`, argIdx)
		args = append(args, floorId)
		argIdx++
	}
	if propertyId != "" {
		query += fmt.Sprintf(` AND property_id = $%d`, argIdx)
		args = append(args, propertyId)
		argIdx++
	}
	if buildingId != "" {
		query += fmt.Sprintf(` AND floor_id IN (SELECT id FROM floors WHERE building_id = $%d)`, argIdx)
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
		return 0, fmt.Errorf("room repository: count: %w", err)
	}
	return count, nil
}
