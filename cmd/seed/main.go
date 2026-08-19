package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"golaninvite/internal/config"
	"golaninvite/internal/database"
	"golaninvite/internal/users"
)

func main() {
	log.Println("[MIGRATOR & SEEDER] Memulai proses inisialisasi schema dan data...")

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[SEEDER] Error: Gagal membaca konfigurasi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[SEEDER] Error: Gagal terhubung ke database: %v", err)
	}
	defer database.Close(pool)

	log.Println("[SEEDER] Berhasil terhubung ke database PostgreSQL.")

	// Eksekusi seluruh file .up.sql di folder migrations
	if err := runMigrations(ctx, pool); err != nil {
		log.Fatalf("[SEEDER] Error saat eksekusi migrasi: %v", err)
	}

	// Buat Akun Admin jika belum ada
	adminEmail := cfg.AdminEmail
	if adminEmail == "" {
		adminEmail = "admin@golaninvite.com"
	}
	adminName := cfg.AdminName
	if adminName == "" {
		adminName = "Administrator"
	}
	adminPassword := cfg.AdminPassword
	if adminPassword == "" {
		adminPassword = "AdminPassword123!"
	}

	userRepo := users.NewPostgresRepository(pool)
	existingUser, err := userRepo.FindByEmail(ctx, adminEmail)
	if err == nil && existingUser != nil {
		log.Printf("[SEEDER] Akun admin '%s' sudah aktif. Seeding user dilewati.\n", existingUser.Email)
	} else {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), 12)
		now := time.Now().UTC()
		adminUser := &users.User{
			ID:           uuid.New(),
			Name:         adminName,
			Email:        adminEmail,
			PasswordHash: string(hashed),
			Role:         users.RoleAdmin,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := userRepo.Create(ctx, adminUser); err != nil {
			log.Printf("[SEEDER] Gagal membuat admin: %v\n", err)
		} else {
			log.Println("[SUCCESS] Akun admin pertama berhasil dibuat!")
		}
	}

	log.Println("[SEEDER] Inisialisasi Database Selesai dengan Sukses.")
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrationFiles := []string{
		"migrations/000001_create_users_and_sessions.up.sql",
		"migrations/000002_create_landing_and_packages_tables.up.sql",
		"migrations/000003_create_full_admin_tables.up.sql",
		"migrations/000004_enhance_packages_and_reviews.up.sql",
	}

	for _, file := range migrationFiles {
		cleanPath := filepath.Clean(file)
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			log.Printf("[MIGRATION] Peringatan: File %s tidak terbaca: %v\n", file, err)
			continue
		}

		log.Printf("[MIGRATION] Mengeksekusi: %s\n", file)
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("migrasi %s gagal: %w", file, err)
		}
	}
	return nil
}
