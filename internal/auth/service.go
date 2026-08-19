package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"golaninvite/internal/middleware"
	"golaninvite/internal/users"
)

const (
	// SessionCookieName adalah nama cookie session yang digunakan di HTTP.
	SessionCookieName = "golaninvite_session"

	// bcryptCost adalah cost factor bcrypt. Nilai 12 memberikan keseimbangan keamanan/performa.
	bcryptCost = 12
)

// Service mendefinisikan business logic autentikasi.
type Service interface {
	Login(ctx context.Context, email, password string) (*SessionUser, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionUser, error)
	InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error
	HashPassword(password string) (string, error)
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (*middleware.AuthenticatedUser, error)
}

type service struct {
	authRepo Repository
	userRepo users.Repository
}

func NewService(authRepo Repository, userRepo users.Repository) Service {
	return &service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*SessionUser, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, users.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, users.ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, users.ErrInvalidCredentials
	}

	now := time.Now().UTC()
	session := &Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpiresAt: now.Add(SessionDuration),
		CreatedAt: now,
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("auth.service.Login: gagal membuat session: %w", err)
	}

	return &SessionUser{
		SessionID: session.ID,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.authRepo.DeleteSession(ctx, sessionID); err != nil {
		return fmt.Errorf("auth.service.Logout: %w", err)
	}
	return nil
}

func (s *service) GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionUser, error) {
	session, err := s.authRepo.FindSession(ctx, sessionID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	if !user.IsActive {
		_ = s.authRepo.DeleteSession(ctx, sessionID)
		return nil, users.ErrAccountInactive
	}

	return &SessionUser{
		SessionID: session.ID,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *service) ValidateSession(ctx context.Context, sessionID uuid.UUID) (*middleware.AuthenticatedUser, error) {
	su, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &middleware.AuthenticatedUser{
		SessionID: su.SessionID,
		UserID:    su.UserID,
		Name:      su.Name,
		Email:     su.Email,
		Role:      string(su.Role),
		ExpiresAt: su.ExpiresAt,
	}, nil
}

func (s *service) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	if err := s.authRepo.DeleteAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("auth.service.InvalidateUserSessions: %w", err)
	}
	return nil
}

func (s *service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("auth.service.HashPassword: %w", err)
	}
	return string(hash), nil
}
