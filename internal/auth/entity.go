package auth

import (
	"time"

	"github.com/google/uuid"
	"golaninvite/internal/users"
)

// Session merepresentasikan sesi autentikasi pengguna yang tersimpan di database.
// Sesuai SSOT §17 — Session Cookie dan §56 — Database (tabel sessions).
type Session struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
}

// SessionUser adalah gabungan session + data user yang dibutuhkan middleware.
type SessionUser struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	Name      string
	Email     string
	Role      users.Role
	ExpiresAt time.Time
}
