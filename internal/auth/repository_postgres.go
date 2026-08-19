package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository mendefinisikan kontrak akses data untuk sesi autentikasi.
type Repository interface {
	// CreateSession menyimpan session baru ke database.
	CreateSession(ctx context.Context, session *Session) error

	// FindSession mencari session berdasarkan ID-nya dan memastikan belum expired.
	FindSession(ctx context.Context, sessionID uuid.UUID) (*Session, error)

	// DeleteSession menghapus session (logout / revoke).
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error

	// DeleteAllUserSessions menghapus semua session milik user (invalidasi akun).
	DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error

	// DeleteExpiredSessions membersihkan session yang sudah expired dari database.
	DeleteExpiredSessions(ctx context.Context) (int64, error)
}

// postgresRepository adalah implementasi Repository menggunakan PostgreSQL.
type postgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository membuat instance baru auth postgresRepository.
func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

// CreateSession menyimpan session baru.
func (r *postgresRepository) CreateSession(ctx context.Context, session *Session) error {
	const q = `
		INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, q,
		session.ID, session.UserID, session.ExpiresAt, session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("auth.repo.CreateSession: %w", err)
	}

	return nil
}

// FindSession mencari session yang valid (belum expired) berdasarkan ID.
func (r *postgresRepository) FindSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	const q = `
		SELECT id, user_id, expires_at, created_at
		FROM sessions
		WHERE id = $1 AND expires_at > NOW()
	`

	session := &Session{}
	err := r.db.QueryRow(ctx, q, sessionID).Scan(
		&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("auth.repo.FindSession: %w", err)
	}

	return session, nil
}

// DeleteSession menghapus session berdasarkan ID (untuk logout / revoke).
func (r *postgresRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	const q = `DELETE FROM sessions WHERE id = $1`

	_, err := r.db.Exec(ctx, q, sessionID)
	if err != nil {
		return fmt.Errorf("auth.repo.DeleteSession: %w", err)
	}

	return nil
}

// DeleteAllUserSessions menghapus semua session milik user tertentu.
// Dipanggil saat akun dinonaktifkan sesuai SSOT §17.
func (r *postgresRepository) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	const q = `DELETE FROM sessions WHERE user_id = $1`

	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("auth.repo.DeleteAllUserSessions: %w", err)
	}

	return nil
}

// DeleteExpiredSessions membersihkan session expired. Dipanggil oleh background job/cron.
func (r *postgresRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	const q = `DELETE FROM sessions WHERE expires_at <= NOW()`

	tag, err := r.db.Exec(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("auth.repo.DeleteExpiredSessions: %w", err)
	}

	return tag.RowsAffected(), nil
}

// SessionDuration adalah masa aktif session default.
const SessionDuration = 24 * time.Hour * 7 // 7 hari
