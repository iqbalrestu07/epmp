package repository

import (
	"context"
	"fmt"
	"github.com/epmp/backend/internal/pkg/errs"
	"time"

	"github.com/epmp/backend/internal/modules/organization/dto"
	"github.com/epmp/backend/internal/modules/organization/entity"
	"github.com/epmp/backend/internal/pkg/uid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
		return fmt.Errorf("organization repository: save: id must be pre-set by caller (use uid.New())")
	}

	var exists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM organizations WHERE id=$1)`, e.Id).Scan(&exists); err != nil {
		return fmt.Errorf("organization repository: save: check exists: %w", err)
	}

	var createdBy interface{}
	if e.CreatedBy != "" {
		createdBy = e.CreatedBy
	}

	if !exists {
		// INSERT with caller-provided ULID
		err := r.db.QueryRow(ctx, `
			INSERT INTO organizations (id, name, domain, is_active, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING created_at, updated_at`,
			e.Id, e.Name, e.Domain, e.IsActive, createdBy,
		).Scan(&e.CreatedAt, &e.UpdatedAt)
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
	var createdBy *string
	err := r.db.QueryRow(ctx, `
		SELECT id, name, domain, is_active, created_by, created_at, updated_at, deleted_at
		FROM   organizations
		WHERE  id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&e.Id, &e.Name, &e.Domain, &e.IsActive, &createdBy, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if createdBy != nil {
		e.CreatedBy = *createdBy
	}

	if err != nil {
		return nil, fmt.Errorf("organization repository: find by id: %w", err)
	}
	return e, nil
}

func (r *OrganizationRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.Organization, error) {
	query := `
		SELECT id, name, domain, is_active, created_by, created_at, updated_at, deleted_at
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
		var createdBy *string
		if err := rows.Scan(&e.Id, &e.Name, &e.Domain, &e.IsActive, &createdBy, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		if createdBy != nil {
			e.CreatedBy = *createdBy
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// FindByUserID returns all organizations where the user is a member.
func (r *OrganizationRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*entity.Organization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT o.id, o.name, o.domain, o.is_active, o.created_by, o.created_at, o.updated_at, o.deleted_at
		FROM   organizations o
		INNER JOIN organization_members om ON om.organization_id = o.id
		WHERE  om.user_id = $1
		  AND  om.is_active = true
		  AND  om.deleted_at IS NULL
		  AND  o.deleted_at IS NULL
		ORDER  BY o.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("organization repository: find by user id: %w", err)
	}
	defer rows.Close()

	var list []*entity.Organization
	for rows.Next() {
		e := &entity.Organization{}
		var createdBy *string
		if err := rows.Scan(&e.Id, &e.Name, &e.Domain, &e.IsActive, &createdBy, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		if createdBy != nil {
			e.CreatedBy = *createdBy
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
	tag, err := r.db.Exec(ctx, `
		UPDATE organizations SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// GetMemberRole returns the caller's role in an organization, or ErrNotFound.
func (r *OrganizationRepositoryImpl) GetMemberRole(ctx context.Context, orgID, userID string) (string, error) {
	var role string
	err := r.db.QueryRow(ctx, `
		SELECT role FROM organization_members
		WHERE  organization_id = $1 AND user_id = $2 AND is_active = true AND deleted_at IS NULL`,
		orgID, userID).Scan(&role)
	if err != nil {
		return "", errs.ErrNotFound
	}
	return role, nil
}

// SaveMember inserts or updates an organization_members record.
func (r *OrganizationRepositoryImpl) SaveMember(ctx context.Context, m *entity.OrganizationMember) error {
	if m.Id == "" {
		return fmt.Errorf("organization repository: save member: id must be pre-set by caller")
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO organization_members (id, organization_id, user_id, role, invited_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (organization_id, user_id)
		DO UPDATE SET role=$4, is_active=$6, updated_at=now()`,
		m.Id, m.OrganizationId, m.UserId, string(m.Role), m.InvitedBy, m.IsActive,
	)
	return err
}

// FindMembersByOrgID returns all active members of an organization.
func (r *OrganizationRepositoryImpl) FindMembersByOrgID(ctx context.Context, orgID string) ([]*entity.OrganizationMember, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, organization_id, user_id, role, invited_by, joined_at, is_active, created_at, updated_at, deleted_at
		FROM   organization_members
		WHERE  organization_id = $1 AND is_active = true AND deleted_at IS NULL`,
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("organization repository: find members: %w", err)
	}
	defer rows.Close()

	var list []*entity.OrganizationMember
	for rows.Next() {
		m := &entity.OrganizationMember{}
		if err := rows.Scan(&m.Id, &m.OrganizationId, &m.UserId, &m.Role, &m.InvitedBy, &m.JoinedAt, &m.IsActive, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// FindMembersWithUserByOrgID returns active members with user details.
func (r *OrganizationRepositoryImpl) FindMembersWithUserByOrgID(ctx context.Context, orgID string) ([]*dto.OrganizationMemberResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT om.id, om.organization_id, om.user_id, COALESCE(u.name, ''), COALESCE(u.email, ''), om.role, om.joined_at, om.is_active
		FROM   organization_members om
		JOIN   users u ON u.id = om.user_id
		WHERE  om.organization_id = $1 AND om.is_active = true AND om.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER  BY om.joined_at ASC`,
		orgID,
	)
	if err != nil {
		return nil, fmt.Errorf("organization repository: find members with user: %w", err)
	}
	defer rows.Close()

	var list []*dto.OrganizationMemberResponse
	for rows.Next() {
		m := &dto.OrganizationMemberResponse{}
		if err := rows.Scan(&m.Id, &m.OrganizationId, &m.UserId, &m.UserName, &m.UserEmail, &m.Role, &m.JoinedAt, &m.IsActive); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// DeleteMember soft-deletes an organization member.
func (r *OrganizationRepositoryImpl) DeleteMember(ctx context.Context, orgID, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE organization_members
		SET    is_active = false, deleted_at = now()
		WHERE  organization_id = $1 AND user_id = $2`,
		orgID, userID,
	)
	return err
}

// AddMemberByEmail adds an existing user or creates a new user and joins them to the organization.
func (r *OrganizationRepositoryImpl) AddMemberByEmail(ctx context.Context, orgID string, email, name, password, role, invitedBy string) (*dto.OrganizationMemberResponse, error) {
	if role == "" {
		role = "member"
	}
	var userID, userName string
	// Check if user exists
	err := r.db.QueryRow(ctx, `SELECT id, name FROM users WHERE email = $1 AND deleted_at IS NULL`, email).Scan(&userID, &userName)
	if err != nil {
		// User doesn't exist, create user
		if password == "" {
			password = "Password123!"
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		if name == "" {
			name = email
		}
		err = r.db.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, name, is_active, organization_id)
			VALUES ($1, $2, $3, true, $4)
			RETURNING id, name`,
			email, string(hash), name, orgID,
		).Scan(&userID, &userName)
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	memberID := uid.New()
	var joinedAt time.Time
	var isActive bool
	err = r.db.QueryRow(ctx, `
		INSERT INTO organization_members (id, organization_id, user_id, role, invited_by, is_active)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, true)
		ON CONFLICT (organization_id, user_id)
		DO UPDATE SET role = $4, is_active = true, deleted_at = NULL, updated_at = now()
		RETURNING id, joined_at, is_active`,
		memberID, orgID, userID, role, invitedBy,
	).Scan(&memberID, &joinedAt, &isActive)
	if err != nil {
		return nil, fmt.Errorf("insert org member: %w", err)
	}

	return &dto.OrganizationMemberResponse{
		Id:             memberID,
		OrganizationId: orgID,
		UserId:         userID,
		UserName:       userName,
		UserEmail:      email,
		Role:           role,
		JoinedAt:       joinedAt,
		IsActive:       isActive,
	}, nil
}
