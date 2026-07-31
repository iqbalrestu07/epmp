package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/epmp/backend/internal/modules/iam/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgUserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL-backed UserRepository.
func NewUserRepository(db *pgxpool.Pool) *pgUserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Save(ctx context.Context, u *entity.User) error {
	if u.ID == "" {
		// INSERT
		row := r.db.QueryRow(ctx, `
			INSERT INTO users (email, password_hash, name, is_active)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			u.Email, u.PasswordHash, u.Name, u.IsActive,
		)
		return row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	}
	// UPDATE
	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET    name=$1, is_active=$2, updated_at=now()
		WHERE  id=$3 AND deleted_at IS NULL`,
		u.Name, u.IsActive, u.ID,
	)
	return err
}

func (r *pgUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	u := &entity.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, name, is_active, last_login_at,
		       created_at, updated_at, deleted_at
		FROM   users
		WHERE  id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.IsActive,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	u := &entity.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, name, is_active, last_login_at,
		       created_at, updated_at, deleted_at
		FROM   users
		WHERE  email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.IsActive,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *pgUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, email, password_hash, name, is_active, last_login_at,
		       created_at, updated_at, deleted_at
		FROM   users
		WHERE  deleted_at IS NULL
		ORDER  BY created_at DESC
		LIMIT  $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.IsActive,
			&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *pgUserRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&count)
	return count, err
}

func (r *pgUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *pgUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.db.Exec(ctx, `
		UPDATE users SET last_login_at = $1 WHERE id = $2`, now, id)
	return err
}
