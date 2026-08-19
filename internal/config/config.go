package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config menyimpan seluruh konfigurasi aplikasi yang dibaca dari environment variables.
// Sesuai SSOT §79 — Environment Variables.
type Config struct {
	// App
	AppEnv  string
	AppName string
	AppHost string
	AppPort int

	// Domain
	AppDomain        string
	AppBaseURL       string
	PublicBaseDomain string

	// Database
	DatabaseURL string

	// Security
	SessionSecret string
	CSRFSecret    string
	EncryptionKey string

	// Admin Seeder
	AdminEmail    string
	AdminPassword string
	AdminName     string

	// Upload
	MaxUploadSizeMB int

	// Logging
	LogLevel string
}

// Load membaca konfigurasi dari environment variables.
// godotenv.Load() harus dipanggil sebelum memanggil fungsi ini
// sehingga .env sudah termuat ke os.Environ.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		AppName:          getEnv("APP_NAME", "GolanInvite"),
		AppHost:          getEnv("APP_HOST", "0.0.0.0"),
		AppDomain:        getEnv("APP_DOMAIN", "localhost"),
		AppBaseURL:       getEnv("APP_BASE_URL", "http://localhost:8080"),
		PublicBaseDomain: getEnv("PUBLIC_BASE_DOMAIN", ""),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		SessionSecret:    getEnv("SESSION_SECRET", ""),
		CSRFSecret:       getEnv("CSRF_SECRET", ""),
		EncryptionKey:    getEnv("ENCRYPTION_KEY", ""),
		AdminEmail:       getEnv("ADMIN_EMAIL", "admin@golaninvite.com"),
		AdminPassword:    getEnv("ADMIN_PASSWORD", ""),
		AdminName:        getEnv("ADMIN_NAME", "Administrator"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Parse integer values
	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("config: APP_PORT bukan angka yang valid: %w", err)
	}
	cfg.AppPort = port

	maxUpload, err := strconv.Atoi(getEnv("MAX_UPLOAD_SIZE_MB", "10"))
	if err != nil {
		return nil, fmt.Errorf("config: MAX_UPLOAD_SIZE_MB bukan angka yang valid: %w", err)
	}
	cfg.MaxUploadSizeMB = maxUpload

	// Validasi wajib
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL wajib diisi")
	}

	return cfg, nil
}

// getEnv membaca env var atau mengembalikan nilai default.
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
