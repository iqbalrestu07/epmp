package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/occupancy/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OccupancyRepositoryImpl implements OccupancyRepository using PostgreSQL.
type OccupancyRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewOccupancyRepositoryImpl creates a new OccupancyRepositoryImpl.
func NewOccupancyRepositoryImpl(db *pgxpool.Pool) *OccupancyRepositoryImpl {
	return &OccupancyRepositoryImpl{db: db}
}

// Ensure OccupancyRepositoryImpl implements domain repository interface.
var _ OccupancyRepository = (*OccupancyRepositoryImpl)(nil)

func (r *OccupancyRepositoryImpl) Save(ctx context.Context, e *entity.Occupancy) error {
	var checkOutTime interface{} = e.CheckOutTime
	if e.CheckOutTime.IsZero() {
		checkOutTime = nil
	}

	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO occupancies (organization_id, contract_id, room_id, tenant_id, status, check_in_time, check_out_time, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.ContractId, e.RoomId, e.TenantId, e.Status, e.CheckInTime, checkOutTime, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE occupancies
		SET    contract_id=$1, room_id=$2, tenant_id=$3, status=$4, check_in_time=$5, check_out_time=$6, notes=$7
		WHERE  id=$8 AND organization_id=$9 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.ContractId, e.RoomId, e.TenantId, e.Status, e.CheckInTime, checkOutTime, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *OccupancyRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Occupancy, error) {
	e := &entity.Occupancy{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, contract_id, room_id, tenant_id, status, check_in_time, COALESCE(check_out_time, '0001-01-01 00:00:00+00'::timestamptz), notes, deleted_at, created_at, updated_at
		FROM   occupancies
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.RoomId, &e.TenantId, &e.Status, &e.CheckInTime, &e.CheckOutTime, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("occupancy repository: find by id: %w", err)
	}
	return e, nil
}

func (r *OccupancyRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Occupancy, error) {
	query := `
		SELECT organization_id, id, contract_id, room_id, tenant_id, status, check_in_time, COALESCE(check_out_time, '0001-01-01 00:00:00+00'::timestamptz), notes, deleted_at, created_at, updated_at
		FROM   occupancies
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("occupancy repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Occupancy
	for rows.Next() {
		e := &entity.Occupancy{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.RoomId, &e.TenantId, &e.Status, &e.CheckInTime, &e.CheckOutTime, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *OccupancyRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM occupancies WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("occupancy repository: count: %w", err)
	}
	return count, nil
}

func (r *OccupancyRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE occupancies SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
