package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey contextKey = "session_user"
	RequestIDKey   contextKey = "request_id"
	SessionCookie  string     = "golaninvite_session"
)

// AuthenticatedUser merepresentasikan data user dalam konteks request.
type AuthenticatedUser struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionValidator mendefinisikan kontrak validasi session untuk middleware.
type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (*AuthenticatedUser, error)
}

// SetRequestID menyisipkan Request ID unik untuk tracing structured logging.
func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth memvalidasi session cookie dan memasukkan data authenticated user ke request context.
func RequireAuth(validator SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookie)
			if err != nil || cookie.Value == "" {
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Sesi tidak valid atau telah berakhir"}}`, http.StatusUnauthorized)
				return
			}

			sessionID, err := uuid.Parse(cookie.Value)
			if err != nil {
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Sesi tidak valid"}}`, http.StatusUnauthorized)
				return
			}

			user, err := validator.ValidateSession(r.Context(), sessionID)
			if err != nil {
				// Hapus cookie kadaluwarsa
				http.SetCookie(w, &http.Cookie{
					Name:     SessionCookie,
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					MaxAge:   -1,
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
				})
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Sesi telah kedaluwarsa"}}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole memvalidasi peran user dalam request context.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(UserContextKey)
			user, ok := val.(*AuthenticatedUser)
			if !ok || user == nil || user.Role != role {
				http.Error(w, `{"error":{"code":"FORBIDDEN","message":"Anda tidak memiliki izin untuk mengakses resource ini"}}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
