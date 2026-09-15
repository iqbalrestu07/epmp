package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/charge/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ChargeRepositoryImpl implements ChargeRepository using PostgreSQL.
type ChargeRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewChargeRepositoryImpl creates a new ChargeRepositoryImpl.
func NewChargeRepositoryImpl(db *pgxpool.Pool) *ChargeRepositoryImpl {
	return &ChargeRepositoryImpl{db: db}
}

// Ensure ChargeRepositoryImpl implements domain repository interface.
var _ ChargeRepository = (*ChargeRepositoryImpl)(nil)

func (r *ChargeRepositoryImpl) Save(ctx context.Context, e *entity.Charge) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO charges (organization_id, contract_id, invoice_id, charge_type, amount, status, charge_date, notes)
			VALUES ($1, $2, NULLIF($3, '')::uuid, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.ContractId, e.InvoiceId, e.ChargeType, e.Amount, e.Status, e.ChargeDate, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE charges
		SET    contract_id=$1, invoice_id=NULLIF($2, '')::uuid, charge_type=$3, amount=$4, status=$5, charge_date=$6, notes=$7
		WHERE  id=$8 AND organization_id=$9 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.ContractId, e.InvoiceId, e.ChargeType, e.Amount, e.Status, e.ChargeDate, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *ChargeRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Charge, error) {
	e := &entity.Charge{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, contract_id, COALESCE(invoice_id::text, ''), charge_type, amount, status, charge_date, notes, deleted_at, created_at, updated_at
		FROM   charges
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.InvoiceId, &e.ChargeType, &e.Amount, &e.Status, &e.ChargeDate, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("charge repository: find by id: %w", err)
	}
	return e, nil
}

func (r *ChargeRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Charge, error) {
	query := `
		SELECT organization_id, id, contract_id, COALESCE(invoice_id::text, ''), charge_type, amount, status, charge_date, notes, deleted_at, created_at, updated_at
		FROM   charges
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (charge_type ILIKE $%d OR status ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("charge repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Charge
	for rows.Next() {
		e := &entity.Charge{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.InvoiceId, &e.ChargeType, &e.Amount, &e.Status, &e.ChargeDate, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *ChargeRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM charges WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (charge_type ILIKE $%d OR status ILIKE $%d)", 2, 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("charge repository: count: %w", err)
	}
	return count, nil
}

func (r *ChargeRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE charges SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
