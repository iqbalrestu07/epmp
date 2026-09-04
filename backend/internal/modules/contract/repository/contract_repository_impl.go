package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/contract/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ContractRepositoryImpl implements ContractRepository using PostgreSQL.
type ContractRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewContractRepositoryImpl creates a new ContractRepositoryImpl.
func NewContractRepositoryImpl(db *pgxpool.Pool) *ContractRepositoryImpl {
	return &ContractRepositoryImpl{db: db}
}

// Ensure ContractRepositoryImpl implements domain repository interface.
var _ ContractRepository = (*ContractRepositoryImpl)(nil)

func (r *ContractRepositoryImpl) Save(ctx context.Context, e *entity.Contract) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO contracts (organization_id, reservation_id, tenant_id, property_id, room_id, status, start_date, end_date, monthly_rent, deposit_amount, terms)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.ReservationId, e.TenantId, e.PropertyId, e.RoomId, e.Status, e.StartDate, e.EndDate, e.MonthlyRent, e.DepositAmount, e.Terms,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE contracts
		SET    reservation_id=$1, tenant_id=$2, property_id=$3, room_id=$4, status=$5, start_date=$6, end_date=$7, monthly_rent=$8, deposit_amount=$9, terms=$10
		WHERE  id=$11 AND organization_id=$12 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.ReservationId, e.TenantId, e.PropertyId, e.RoomId, e.Status, e.StartDate, e.EndDate, e.MonthlyRent, e.DepositAmount, e.Terms, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *ContractRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Contract, error) {
	e := &entity.Contract{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, reservation_id, tenant_id, property_id, room_id, status, start_date, end_date, monthly_rent, deposit_amount, terms, deleted_at, created_at, updated_at
		FROM   contracts
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.ReservationId, &e.TenantId, &e.PropertyId, &e.RoomId, &e.Status, &e.StartDate, &e.EndDate, &e.MonthlyRent, &e.DepositAmount, &e.Terms, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("contract repository: find by id: %w", err)
	}
	return e, nil
}

func (r *ContractRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Contract, error) {
	query := `
		SELECT organization_id, id, reservation_id, tenant_id, property_id, room_id, status, start_date, end_date, monthly_rent, deposit_amount, terms, deleted_at, created_at, updated_at
		FROM   contracts
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("contract repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Contract
	for rows.Next() {
		e := &entity.Contract{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.ReservationId, &e.TenantId, &e.PropertyId, &e.RoomId, &e.Status, &e.StartDate, &e.EndDate, &e.MonthlyRent, &e.DepositAmount, &e.Terms, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *ContractRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM contracts WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("contract repository: count: %w", err)
	}
	return count, nil
}

func (r *ContractRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE contracts SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
