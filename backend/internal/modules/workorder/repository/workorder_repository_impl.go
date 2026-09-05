package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/workorder/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkOrderRepositoryImpl implements WorkOrderRepository using PostgreSQL.
type WorkOrderRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewWorkOrderRepositoryImpl creates a new WorkOrderRepositoryImpl.
func NewWorkOrderRepositoryImpl(db *pgxpool.Pool) *WorkOrderRepositoryImpl {
	return &WorkOrderRepositoryImpl{db: db}
}

// Ensure WorkOrderRepositoryImpl implements domain repository interface.
var _ WorkOrderRepository = (*WorkOrderRepositoryImpl)(nil)

func (r *WorkOrderRepositoryImpl) Save(ctx context.Context, e *entity.WorkOrder) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO work_orders (organization_id, property_id, room_id, description, status, priority)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.PropertyId, e.RoomId, e.Description, e.Status, e.Priority,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE work_orders
		SET    property_id=$1, room_id=$2, description=$3, status=$4, priority=$5
		WHERE  id=$6 AND organization_id=$7 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.PropertyId, e.RoomId, e.Description, e.Status, e.Priority, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *WorkOrderRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.WorkOrder, error) {
	e := &entity.WorkOrder{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, property_id, room_id, description, status, priority, deleted_at, created_at, updated_at
		FROM   work_orders
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &e.RoomId, &e.Description, &e.Status, &e.Priority, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("workorder repository: find by id: %w", err)
	}
	return e, nil
}

func (r *WorkOrderRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.WorkOrder, error) {
	query := `
		SELECT organization_id, id, property_id, room_id, description, status, priority, deleted_at, created_at, updated_at
		FROM   work_orders
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (status ILIKE $%d OR priority ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("workorder repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.WorkOrder
	for rows.Next() {
		e := &entity.WorkOrder{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.PropertyId, &e.RoomId, &e.Description, &e.Status, &e.Priority, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *WorkOrderRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM work_orders WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (status ILIKE $%d OR priority ILIKE $%d)", 2, 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("workorder repository: count: %w", err)
	}
	return count, nil
}

func (r *WorkOrderRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE work_orders SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
