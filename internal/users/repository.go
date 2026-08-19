package users

import (
	"context"

	"github.com/google/uuid"
)

// Repository mendefinisikan kontrak akses data untuk modul users.
// Sesuai SSOT §15 — Struktur Modul Backend.
type Repository interface {
	// FindByID mencari user berdasarkan UUID-nya.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail mencari user berdasarkan email-nya.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Create menyimpan user baru ke database.
	Create(ctx context.Context, user *User) error

	// Update memperbarui data user yang sudah ada.
	Update(ctx context.Context, user *User) error

	// List mengembalikan daftar user dengan pagination sederhana.
	List(ctx context.Context, limit, offset int) ([]*User, int, error)

	// ExistsWithEmail memeriksa apakah email sudah digunakan, opsional mengabaikan ID tertentu.
	ExistsWithEmail(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
}
