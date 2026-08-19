package users

import (
	"time"

	"github.com/google/uuid"
)

// Role merepresentasikan peran pengguna dalam sistem.
// Sesuai SSOT §16 — Authentication.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// IsValid memeriksa apakah role valid.
func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser
}

// User merepresentasikan entitas pengguna dalam sistem.
// Sesuai SSOT §56 (Database), §48 (User Management).
type User struct {
	ID           uuid.UUID  `db:"id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Role         Role       `db:"role"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
