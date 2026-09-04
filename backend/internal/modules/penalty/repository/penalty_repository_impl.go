package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/penalty/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PenaltyRepositoryImpl implements PenaltyRepository using PostgreSQL.
type PenaltyRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewPenaltyRepositoryImpl creates a new PenaltyRepositoryImpl.
func NewPenaltyRepositoryImpl(db *pgxpool.Pool) *PenaltyRepositoryImpl {
	return &PenaltyRepositoryImpl{db: db}
}

// Ensure PenaltyRepositoryImpl implements domain repository interface.
var _ PenaltyRepository = (*PenaltyRepositoryImpl)(nil)

func (r *PenaltyRepositoryImpl) Save(ctx context.Context, e *entity.Penalty) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO penalties (organization_id, invoice_id, amount, status, penalty_date, description)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.InvoiceId, e.Amount, e.Status, e.PenaltyDate, e.Description,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE penalties
		SET    invoice_id=$1, amount=$2, status=$3, penalty_date=$4, description=$5
		WHERE  id=$6 AND organization_id=$7 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.InvoiceId, e.Amount, e.Status, e.PenaltyDate, e.Description, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *PenaltyRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.Penalty, error) {
	e := &entity.Penalty{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, invoice_id, amount, status, penalty_date, description, deleted_at, created_at, updated_at
		FROM   penalties
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.Amount, &e.Status, &e.PenaltyDate, &e.Description, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("penalty repository: find by id: %w", err)
	}
	return e, nil
}

func (r *PenaltyRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.Penalty, error) {
	query := `
		SELECT organization_id, id, invoice_id, amount, status, penalty_date, description, deleted_at, created_at, updated_at
		FROM   penalties
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
		return nil, fmt.Errorf("penalty repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Penalty
	for rows.Next() {
		e := &entity.Penalty{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.InvoiceId, &e.Amount, &e.Status, &e.PenaltyDate, &e.Description, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *PenaltyRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM penalties WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND status ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("penalty repository: count: %w", err)
	}
	return count, nil
}

func (r *PenaltyRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE penalties SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
