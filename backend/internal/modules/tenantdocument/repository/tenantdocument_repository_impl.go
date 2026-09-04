package repository

import (
	"context"
	"fmt"

	"github.com/epmp/backend/internal/modules/tenantdocument/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantDocumentRepositoryImpl implements TenantDocumentRepository using PostgreSQL.
type TenantDocumentRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewTenantDocumentRepositoryImpl creates a new TenantDocumentRepositoryImpl.
func NewTenantDocumentRepositoryImpl(db *pgxpool.Pool) *TenantDocumentRepositoryImpl {
	return &TenantDocumentRepositoryImpl{db: db}
}

// Ensure TenantDocumentRepositoryImpl implements domain repository interface.
var _ TenantDocumentRepository = (*TenantDocumentRepositoryImpl)(nil)

func (r *TenantDocumentRepositoryImpl) Save(ctx context.Context, e *entity.TenantDocument) error {
	if e.Id == "" {
		err := r.db.QueryRow(ctx, `
			INSERT INTO tenant_documents (organization_id, tenant_id, document_type, file_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			e.OrganizationId, e.TenantId, e.DocumentType, e.FileUrl,
		).Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt)
		return err
	}
	err := r.db.QueryRow(ctx, `
		UPDATE tenant_documents
		SET    tenant_id=$1, document_type=$2, file_url=$3
		WHERE  id=$4 AND organization_id=$5 AND deleted_at IS NULL
		RETURNING updated_at`,
		e.TenantId, e.DocumentType, e.FileUrl, e.Id, e.OrganizationId,
	).Scan(&e.UpdatedAt)
	return err
}

func (r *TenantDocumentRepositoryImpl) FindByID(ctx context.Context, id, orgID string) (*entity.TenantDocument, error) {
	e := &entity.TenantDocument{}
	err := r.db.QueryRow(ctx, `
		SELECT organization_id, id, tenant_id, document_type, file_url, deleted_at, created_at, updated_at
		FROM   tenant_documents
		WHERE  id = $1 AND organization_id = $2 AND deleted_at IS NULL`,
		id, orgID,
	).Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.DocumentType, &e.FileUrl, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("tenantdocument repository: find by id: %w", err)
	}
	return e, nil
}

func (r *TenantDocumentRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search, orgID string) ([]*entity.TenantDocument, error) {
	query := `
		SELECT organization_id, id, tenant_id, document_type, file_url, deleted_at, created_at, updated_at
		FROM   tenant_documents
		WHERE  deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if search != "" {
		query += fmt.Sprintf(" AND document_type ILIKE $%d", argIdx)
		args = append(args, "%" + search + "%")
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("tenantdocument repository: find all: %w", err)
	}
	defer rows.Close()

	var list []*entity.TenantDocument
	for rows.Next() {
		e := &entity.TenantDocument{}
		if err := rows.Scan(&e.OrganizationId, &e.Id, &e.TenantId, &e.DocumentType, &e.FileUrl, &e.DeletedAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *TenantDocumentRepositoryImpl) Count(ctx context.Context, search, orgID string) (int64, error) {
	query := `SELECT COUNT(*) FROM tenant_documents WHERE deleted_at IS NULL AND organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		query += fmt.Sprintf(" AND document_type ILIKE $%d", 2)
		args = append(args, "%" + search + "%")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("tenantdocument repository: count: %w", err)
	}
	return count, nil
}

func (r *TenantDocumentRepositoryImpl) Delete(ctx context.Context, id, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tenant_documents SET deleted_at = now() WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`, id, orgID)
	return err
}
