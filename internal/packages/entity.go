package packages

import (
	"time"

	"github.com/google/uuid"
)

// Package merepresentasikan entitas paket undangan digital.
// Sesuai SSOT §44 — Packages.
type Package struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Price        int64     `json:"price" db:"price"`
	Benefits     []string  `json:"benefits" db:"benefits"`
	Badge        string    `json:"badge" db:"badge"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
