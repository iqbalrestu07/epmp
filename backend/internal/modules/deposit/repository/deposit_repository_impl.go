package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/deposit/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DepositRepositoryImpl implements DepositRepository using PostgreSQL.
type DepositRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewDepositRepositoryImpl creates a new DepositRepositoryImpl.
func NewDepositRepositoryImpl(db *pgxpool.Pool) *DepositRepositoryImpl {
	return &DepositRepositoryImpl{db: db}
}

// Ensure DepositRepositoryImpl implements domain repository interface.
var _ DepositRepository = (*DepositRepositoryImpl)(nil)

func (r *DepositRepositoryImpl) Save(ctx context.Context, e *entity.Deposit) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO deposits (organization_id, contract_id, tenant_id, amount, status, collection_date, refund_date, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.ContractId, e.TenantId, e.Amount, e.Status, e.CollectionDate, e.RefundDate, e.Notes,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE deposits
		SET    contract_id=$1, tenant_id=$2, amount=$3, status=$4, collection_date=$5, refund_date=$6, notes=$7
		WHERE  id=$8 AND organization_id=$9 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.ContractId, e.TenantId, e.Amount, e.Status, e.CollectionDate, e.RefundDate, e.Notes, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *DepositRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Deposit, error) {
	e := &entity.Deposit{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, contract_id, tenant_id, amount, status, collection_date, refund_date, notes, deleted_at, created_at, updated_at
		FROM   deposits
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.TenantId, &e.Amount, &e.Status, &e.CollectionDate, &e.RefundDate, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("deposit repository: find by id: %w", err)
	}
	return e, nil
}

func (r *DepositRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Deposit, error) {
	query := `
		SELECT organization_id, id, contract_id, tenant_id, amount, status, collection_date, refund_date, notes, deleted_at, created_at, updated_at
		FROM   deposits
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
		return nil, fmt.Errorf("deposit repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Deposit
	for rows.Next() {
		e := &entity.Deposit{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.ContractId, &e.TenantId, &e.Amount, &e.Status, &e.CollectionDate, &e.RefundDate, &e.Notes, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *DepositRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM deposits WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("deposit repository: count: %w", err)
	}
	return count, nil
}

func (r *DepositRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE deposits SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
