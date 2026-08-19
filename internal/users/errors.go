package users

import "errors"

// Error sentinel untuk domain users.
var (
	// ErrNotFound dikembalikan ketika user tidak ditemukan di database.
	ErrNotFound = errors.New("users: user tidak ditemukan")

	// ErrEmailAlreadyExists dikembalikan ketika email sudah terdaftar.
	ErrEmailAlreadyExists = errors.New("users: email sudah terdaftar")

	// ErrInvalidRole dikembalikan ketika role yang diberikan tidak valid.
	ErrInvalidRole = errors.New("users: role tidak valid")

	// ErrInvalidCredentials dikembalikan ketika email/password tidak cocok.
	ErrInvalidCredentials = errors.New("users: email atau password salah")

	// ErrAccountInactive dikembalikan ketika akun dinonaktifkan.
	ErrAccountInactive = errors.New("users: akun tidak aktif")
)
