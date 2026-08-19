package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresRepository adalah implementasi Repository menggunakan PostgreSQL via pgx/v5.
type postgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository membuat instance baru postgresRepository.
func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

// FindByID mencari user berdasarkan UUID.
func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("users.repo.FindByID: %w", err)
	}

	return user, nil
}

// FindByEmail mencari user berdasarkan email (case-insensitive via LOWER).
func (r *postgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`

	user := &User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("users.repo.FindByEmail: %w", err)
	}

	return user, nil
}

// Create menyimpan user baru. Menggunakan parameterized query sesuai SSOT §85.
func (r *postgresRepository) Create(ctx context.Context, user *User) error {
	const q = `
		INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, q,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		// Deteksi unique constraint violation pada email
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("users.repo.Create: %w", err)
	}

	return nil
}

// Update memperbarui data user yang sudah ada.
func (r *postgresRepository) Update(ctx context.Context, user *User) error {
	const q = `
		UPDATE users
		SET name = $2, email = $3, password_hash = $4, role = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	tag, err := r.db.Exec(ctx, q,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.IsActive, user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("users.repo.Update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// List mengembalikan daftar user dengan total count untuk pagination.
func (r *postgresRepository) List(ctx context.Context, limit, offset int) ([]*User, int, error) {
	const countQ = `SELECT COUNT(*) FROM users`
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var total int
	if err := r.db.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("users.repo.List count: %w", err)
	}

	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("users.repo.List query: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.PasswordHash,
			&u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("users.repo.List scan: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("users.repo.List rows: %w", err)
	}

	return users, total, nil
}

// ExistsWithEmail memeriksa apakah email sudah dipakai, opsional mengabaikan satu ID.
func (r *postgresRepository) ExistsWithEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	var q string
	var args []interface{}

	if excludeID != nil {
		q = `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER($1) AND id != $2)`
		args = []interface{}{email, *excludeID}
	} else {
		q = `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER($1))`
		args = []interface{}{email}
	}

	var exists bool
	if err := r.db.QueryRow(ctx, q, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("users.repo.ExistsWithEmail: %w", err)
	}

	return exists, nil
}

// isUniqueViolation memeriksa apakah error adalah PostgreSQL unique constraint violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	// pgx v5 menggunakan PgError
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
