package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/refund/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RefundRepositoryImpl implements RefundRepository using PostgreSQL.
type RefundRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewRefundRepositoryImpl creates a new RefundRepositoryImpl.
func NewRefundRepositoryImpl(db *pgxpool.Pool) *RefundRepositoryImpl {
	return &RefundRepositoryImpl{db: db}
}

// Ensure RefundRepositoryImpl implements domain repository interface.
var _ RefundRepository = (*RefundRepositoryImpl)(nil)

func (r *RefundRepositoryImpl) Save(ctx context.Context, e *entity.Refund) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO refunds (organization_id, payment_id, tenant_id, amount, status, refund_date, reason)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.PaymentId, e.TenantId, e.Amount, e.Status, e.RefundDate, e.Reason,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE refunds
		SET    payment_id=$1, tenant_id=$2, amount=$3, status=$4, refund_date=$5, reason=$6
		WHERE  id=$7 AND organization_id=$8 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.PaymentId, e.TenantId, e.Amount, e.Status, e.RefundDate, e.Reason, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *RefundRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Refund, error) {
	e := &entity.Refund{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, payment_id, tenant_id, amount, status, refund_date, reason, deleted_at, created_at, updated_at
		FROM   refunds
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.PaymentId, &e.TenantId, &e.Amount, &e.Status, &e.RefundDate, &e.Reason, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("refund repository: find by id: %w", err)
	}
	return e, nil
}

func (r *RefundRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Refund, error) {
	query := `
		SELECT organization_id, id, payment_id, tenant_id, amount, status, refund_date, reason, deleted_at, created_at, updated_at
		FROM   refunds
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
		return nil, fmt.Errorf("refund repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Refund
	for rows.Next() {
		e := &entity.Refund{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.PaymentId, &e.TenantId, &e.Amount, &e.Status, &e.RefundDate, &e.Reason, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *RefundRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM refunds WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("refund repository: count: %w", err)
	}
	return count, nil
}

func (r *RefundRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refunds SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
