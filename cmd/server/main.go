package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"golaninvite/internal/app"
	"golaninvite/internal/config"
	"golaninvite/internal/database"
)

func main() {
	// 1. Setup Structured Logger (slog) sesuai SSOT §90
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("memulai GolanInvite server...")

	// 2. Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("file .env tidak ditemukan, menggunakan variabel environment sistem")
	}

	// 3. Load Konfigurasi
	cfg, err := config.Load()
	if err != nil {
		logger.Error("gagal memuat konfigurasi", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 4. Inisialisasi Database Connection Pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("gagal menghubungkan ke database PostgreSQL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close(pool)

	logger.Info("berhasil terhubung ke database PostgreSQL")

	// 5. Inisialisasi Aplikasi Container
	application := app.New(cfg, pool, logger)

	// 6. Graceful Shutdown Setup
	serverCtx, serverStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer serverStop()

	go func() {
		if err := application.Run(serverCtx); err != nil {
			logger.Error("aplikasi berhenti karena error", slog.String("error", err.Error()))
			serverStop()
		}
	}()

	// Tunggu sinyal terminate / interrupt
	<-serverCtx.Done()
	logger.Info("menerima sinyal shutdown, membersihkan resource...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Error("error saat shutdown server", slog.String("error", err.Error()))
	}

	logger.Info("server GolanInvite berhasil dimatikan secara graceful")
}
