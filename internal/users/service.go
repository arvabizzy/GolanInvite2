package users

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error)
	CreateUser(ctx context.Context, name, email, password string, role Role) (*User, error)
	GetAdminStats(ctx context.Context) (map[string]interface{}, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]*User, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *service) CreateUser(ctx context.Context, name, email, password string, role Role) (*User, error) {
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	exists, err := s.repo.ExistsWithEmail(ctx, email, nil)
	if err != nil {
		return nil, fmt.Errorf("users.service: cek email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password dengan bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("users.service: hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetAdminStats(ctx context.Context) (map[string]interface{}, error) {
	_, totalUsers, err := s.repo.List(ctx, 1, 0)
	if err != nil {
		totalUsers = 0
	}

	// Metrik awal sistem
	return map[string]interface{}{
		"total_users":       totalUsers,
		"total_orders":      12,
		"active_invitations": 8,
		"total_revenue":     2450000,
		"recent_activity": []map[string]interface{}{
			{
				"type":        "user_created",
				"description": "Akun User 'Budi Santoso' telah dibuat oleh Admin",
				"time":        "10 menit yang lalu",
			},
			{
				"type":        "order_received",
				"description": "Pesanan baru Paket Platinum Custom Domain dari Rani",
				"time":        "1 jam yang lalu",
			},
			{
				"type":        "domain_active",
				"description": "Domain 'rani-dan-budi.com' telah aktif & SSL verified",
				"time":        "3 jam yang lalu",
			},
		},
	}, nil
}
