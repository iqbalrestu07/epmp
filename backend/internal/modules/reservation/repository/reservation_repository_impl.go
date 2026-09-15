package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/reservation/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReservationRepositoryImpl implements ReservationRepository using PostgreSQL.
type ReservationRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewReservationRepositoryImpl creates a new ReservationRepositoryImpl.
func NewReservationRepositoryImpl(db *pgxpool.Pool) *ReservationRepositoryImpl {
	return &ReservationRepositoryImpl{db: db}
}

// Ensure ReservationRepositoryImpl implements domain repository interface.
var _ ReservationRepository = (*ReservationRepositoryImpl)(nil)

func (r *ReservationRepositoryImpl) Save(ctx context.Context, e *entity.Reservation) error {
	var checkOutDate interface{} = e.CheckOutDate
	if e.CheckOutDate.IsZero() {
		checkOutDate = nil
	}

	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO reservations (organization_id, tenant_id, property_id, room_id, status, check_in_date, check_out_date, booking_fee, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.TenantId, e.PropertyId, e.RoomId, e.Status, e.CheckInDate, checkOutDate, e.BookingFee, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE reservations
		SET    tenant_id=$1, property_id=$2, room_id=$3, status=$4, check_in_date=$5, check_out_date=$6, booking_fee=$7, notes=$8
		WHERE  id=$9 AND organization_id=$10 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.TenantId, e.PropertyId, e.RoomId, e.Status, e.CheckInDate, checkOutDate, e.BookingFee, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *ReservationRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Reservation, error) {
	e := &entity.Reservation{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, tenant_id, property_id, room_id, status, check_in_date, COALESCE(check_out_date, '0001-01-01 00:00:00+00'::timestamptz), booking_fee, notes, deleted_at, created_at, updated_at
		FROM   reservations
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.PropertyId, &e.RoomId, &e.Status, &e.CheckInDate, &e.CheckOutDate, &e.BookingFee, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("reservation repository: find by id: %w", err)
	}
	return e, nil
}

func (r *ReservationRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Reservation, error) {
	query := `
		SELECT organization_id, id, tenant_id, property_id, room_id, status, check_in_date, COALESCE(check_out_date, '0001-01-01 00:00:00+00'::timestamptz), booking_fee, notes, deleted_at, created_at, updated_at
		FROM   reservations
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
		return nil, fmt.Errorf("reservation repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Reservation
	for rows.Next() {
		e := &entity.Reservation{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.PropertyId, &e.RoomId, &e.Status, &e.CheckInDate, &e.CheckOutDate, &e.BookingFee, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *ReservationRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM reservations WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("reservation repository: count: %w", err)
	}
	return count, nil
}

func (r *ReservationRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE reservations SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
