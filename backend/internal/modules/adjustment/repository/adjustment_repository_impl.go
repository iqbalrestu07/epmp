package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"

	"github.com/epmp/backend/internal/modules/adjustment/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdjustmentRepositoryImpl implements AdjustmentRepository using PostgreSQL.
type AdjustmentRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewAdjustmentRepositoryImpl creates a new AdjustmentRepositoryImpl.
func NewAdjustmentRepositoryImpl(db *pgxpool.Pool) *AdjustmentRepositoryImpl {
	return &AdjustmentRepositoryImpl{db: db}
}

// Ensure AdjustmentRepositoryImpl implements domain repository interface.
var _ AdjustmentRepository = (*AdjustmentRepositoryImpl)(nil)

func (r *AdjustmentRepositoryImpl) Save(ctx context.Context, e *entity.Adjustment) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO adjustments (organization_id, invoice_id, adjustment_type, amount, adjustment_date, reason)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.InvoiceId, e.AdjustmentType, e.Amount, e.AdjustmentDate, e.Reason,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE adjustments
		SET    invoice_id=$1, adjustment_type=$2, amount=$3, adjustment_date=$4, reason=$5
		WHERE  id=$6 AND organization_id=$7 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.InvoiceId, e.AdjustmentType, e.Amount, e.AdjustmentDate, e.Reason, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *AdjustmentRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Adjustment, error) {
	e := &entity.Adjustment{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, invoice_id, adjustment_type, amount, adjustment_date, reason, deleted_at, created_at, updated_at
		FROM   adjustments
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.AdjustmentType, &e.Amount, &e.AdjustmentDate, &e.Reason, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("adjustment repository: find by id: %w", err)
	}
	return e, nil
}

func (r *AdjustmentRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Adjustment, error) {
	query := `
		SELECT organization_id, id, invoice_id, adjustment_type, amount, adjustment_date, reason, deleted_at, created_at, updated_at
		FROM   adjustments
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND adjustment_type ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("adjustment repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Adjustment
	for rows.Next() {
		e := &entity.Adjustment{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.AdjustmentType, &e.Amount, &e.AdjustmentDate, &e.Reason, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *AdjustmentRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM adjustments WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND adjustment_type ILIKE $%d", 2)
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("adjustment repository: count: %w", err)
	}
	return count, nil
}

func (r *AdjustmentRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE adjustments SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
