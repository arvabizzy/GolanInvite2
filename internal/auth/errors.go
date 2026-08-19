package auth

import "errors"

// Error sentinel untuk domain autentikasi.
var (
	// ErrSessionNotFound dikembalikan ketika session ID tidak ditemukan atau sudah expired.
	ErrSessionNotFound = errors.New("auth: session tidak ditemukan atau sudah expired")

	// ErrSessionExpired dikembalikan ketika session sudah melewati batas expires_at.
	ErrSessionExpired = errors.New("auth: session sudah expired")

	// ErrUnauthorized dikembalikan ketika request tidak memiliki session valid.
	ErrUnauthorized = errors.New("auth: tidak terautentikasi")

	// ErrForbidden dikembalikan ketika user tidak memiliki izin untuk resource/aksi tertentu.
	ErrForbidden = errors.New("auth: akses ditolak")
)
