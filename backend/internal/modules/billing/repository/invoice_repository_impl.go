package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/billing/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InvoiceRepositoryImpl implements InvoiceRepository using PostgreSQL.
type InvoiceRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewInvoiceRepositoryImpl creates a new InvoiceRepositoryImpl.
func NewInvoiceRepositoryImpl(db *pgxpool.Pool) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{db: db}
}

// Ensure InvoiceRepositoryImpl implements domain repository interface.
var _ InvoiceRepository = (*InvoiceRepositoryImpl)(nil)

func (r *InvoiceRepositoryImpl) Save(ctx context.Context, e *entity.Invoice) error {
	var paidDate interface{} = e.PaidDate
	if e.PaidDate.IsZero() {
		paidDate = nil
	}
	if e.Currency == "" {
		e.Currency = "IDR"
	}

	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO invoices (organization_id, contract_id, tenant_id, amount, currency, status, due_date, paid_date, payment_method, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.ContractId, e.TenantId, e.Amount, e.Currency, e.Status, e.DueDate, paidDate, e.PaymentMethod, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE invoices
		SET    contract_id=$1, tenant_id=$2, amount=$3, currency=$4, status=$5, due_date=$6, paid_date=$7, payment_method=$8, notes=$9
		WHERE  id=$10 AND organization_id=$11 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.ContractId, e.TenantId, e.Amount, e.Currency, e.Status, e.DueDate, paidDate, e.PaymentMethod, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *InvoiceRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Invoice, error) {
	e := &entity.Invoice{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, contract_id, tenant_id, amount, currency, status, due_date, COALESCE(paid_date, '0001-01-01 00:00:00+00'::timestamptz), payment_method, notes, deleted_at, created_at, updated_at
		FROM   invoices
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.TenantId, &e.Amount, &e.Currency, &e.Status, &e.DueDate, &e.PaidDate, &e.PaymentMethod, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("billing repository: find by id: %w", err)
	}
	return e, nil
}

func (r *InvoiceRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Invoice, error) {
	query := `
		SELECT organization_id, id, contract_id, tenant_id, amount, currency, status, due_date, COALESCE(paid_date, '0001-01-01 00:00:00+00'::timestamptz), payment_method, notes, deleted_at, created_at, updated_at
		FROM   invoices
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND (status ILIKE $%d OR payment_method ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("billing repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Invoice
	for rows.Next() {
		e := &entity.Invoice{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.TenantId, &e.Amount, &e.Currency, &e.Status, &e.DueDate, &e.PaidDate, &e.PaymentMethod, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *InvoiceRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM invoices WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (status ILIKE $%d OR payment_method ILIKE $%d)", 2, 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("billing repository: count: %w", err)
	}
	return count, nil
}

func (r *InvoiceRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE invoices SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
