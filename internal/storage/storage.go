package storage

import (
	"context"
	"io"
)

// Storage mendefinisikan interface penyimpanan file (media/upload).
// Sesuai SSOT §87 — Upload Security.
type Storage interface {
	// Save menyimpan file dengan filename acak/terisolasi.
	Save(ctx context.Context, filename string, reader io.Reader, size int64, contentType string) (string, error)

	// Delete menghapus file dari storage.
	Delete(ctx context.Context, path string) error

	// GetURL mengembalikan URL akses publik/internal file.
	GetURL(path string) string
}
