package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New membuat connection pool PostgreSQL menggunakan pgx/v5.
// databaseURL dibaca dari config yang sumbernya adalah DATABASE_URL di .env.
// Sesuai SSOT §9 (PostgreSQL) dan §57 (Database Rules).
func New(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: gagal parse DATABASE_URL: %w", err)
	}

	// Konfigurasi connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("database: gagal membuat connection pool: %w", err)
	}

	// Verifikasi koneksi berhasil
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: gagal ping PostgreSQL: %w", err)
	}

	return pool, nil
}

// Close menutup semua koneksi di pool.
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
