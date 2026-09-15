package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/payment/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepositoryImpl implements PaymentRepository using PostgreSQL.
type PaymentRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewPaymentRepositoryImpl creates a new PaymentRepositoryImpl.
func NewPaymentRepositoryImpl(db *pgxpool.Pool) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db: db}
}

// Ensure PaymentRepositoryImpl implements domain repository interface.
var _ PaymentRepository = (*PaymentRepositoryImpl)(nil)

func (r *PaymentRepositoryImpl) Save(ctx context.Context, e *entity.Payment) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO payments (organization_id, invoice_id, tenant_id, amount, payment_date, payment_method, status, reference_number)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.InvoiceId, e.TenantId, e.Amount, e.PaymentDate, e.PaymentMethod, e.Status, e.ReferenceNumber,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE payments
		SET    invoice_id=$1, tenant_id=$2, amount=$3, payment_date=$4, payment_method=$5, status=$6, reference_number=$7
		WHERE  id=$8 AND organization_id=$9 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.InvoiceId, e.TenantId, e.Amount, e.PaymentDate, e.PaymentMethod, e.Status, e.ReferenceNumber, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *PaymentRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Payment, error) {
	e := &entity.Payment{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, invoice_id, tenant_id, amount, payment_date, payment_method, status, reference_number, deleted_at, created_at, updated_at
		FROM   payments
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.TenantId, &e.Amount, &e.PaymentDate, &e.PaymentMethod, &e.Status, &e.ReferenceNumber, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("payment repository: find by id: %w", err)
	}
	return e, nil
}

func (r *PaymentRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Payment, error) {
	query := `
		SELECT organization_id, id, invoice_id, tenant_id, amount, payment_date, payment_method, status, reference_number, deleted_at, created_at, updated_at
		FROM   payments
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
		return nil, fmt.Errorf("payment repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Payment
	for rows.Next() {
		e := &entity.Payment{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.TenantId, &e.Amount, &e.PaymentDate, &e.PaymentMethod, &e.Status, &e.ReferenceNumber, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *PaymentRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND (status ILIKE $%d OR payment_method ILIKE $%d)", 2, 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("payment repository: count: %w", err)
	}
	return count, nil
}

func (r *PaymentRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE payments SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
