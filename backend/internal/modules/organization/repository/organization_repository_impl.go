package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/organization/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OrganizationRepositoryImpl implements OrganizationRepository using PostgreSQL.
type OrganizationRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewOrganizationRepositoryImpl creates a new OrganizationRepositoryImpl.
func NewOrganizationRepositoryImpl(db *pgxpool.Pool) *OrganizationRepositoryImpl {
	return &OrganizationRepositoryImpl{db: db}
}

// Ensure OrganizationRepositoryImpl implements domain repository interface.
var _ OrganizationRepository = (*OrganizationRepositoryImpl)(nil)

func (r *OrganizationRepositoryImpl) Save(ctx context.Context, e *entity.Organization) error {
	if e.Id == "" {
		// INSERT
		err := r.db.QueryRow(ctx, `
			INSERT INTO organizations (name, domain, is_active)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, updated_at`,
			e.Name, e.Domain, e.IsActive,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE organizations
		SET    name=$1, domain=$2, is_active=$3, updated_at=now()
		WHERE  id=$4 AND deleted_at IS NULL`,
		e.Name, e.Domain, e.IsActive, e.Id,
	)
	return err
}

func (r *OrganizationRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Organization, error) {
	e := &entity.Organization{}
	err := r.db.QueryRow(ctx, `
		SELECT id, name, domain, is_active, created_at, updated_at, deleted_at
		FROM   organizations
		WHERE  id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&e.Id, &e.Name, &e.Domain, &e.IsActive, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)

	if err != nil {
		return nil, fmt.Errorf("organization repository: find by id: %w", err)
	}
	return e, nil
}

func (r *OrganizationRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.Organization, error) {
	query := `
		SELECT id, name, domain, is_active, created_at, updated_at, deleted_at
		FROM   organizations
		WHERE  deleted_at IS NULL`
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(` AND (name ILIKE $%d OR domain ILIKE $%d)`, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("organization repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.Organization
	for rows.Next() {
		e := &entity.Organization{}
		if err := rows.Scan(&e.Id, &e.Name, &e.Domain, &e.IsActive, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *OrganizationRepositoryImpl) Count(ctx context.Context, search string) (int64, error) {
	query := `SELECT COUNT(*) FROM organizations WHERE deleted_at IS NULL`
	args := []interface{}{}

	if search != "" {
		query += ` AND (name ILIKE $1 OR domain ILIKE $1)`
		args = append(args, "%"+search+"%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("organization repository: count: %w", err)
	}
	return count, nil
}

func (r *OrganizationRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE organizations SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}
