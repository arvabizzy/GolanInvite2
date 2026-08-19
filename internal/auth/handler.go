package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"golaninvite/internal/middleware"
	"golaninvite/internal/users"
	"golaninvite/internal/validation"
)

type Handler struct {
	authSvc Service
}

func NewHandler(authSvc Service) *Handler {
	return &Handler{authSvc: authSvc}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User struct {
		ID    uuid.UUID  `json:"id"`
		Name  string     `json:"name"`
		Email string     `json:"email"`
		Role  users.Role `json:"role"`
	} `json:"user"`
	RedirectURL string `json:"redirect_url"`
}

// HandleLogin menangani POST /api/v1/auth/login (SSOT §60)
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format JSON tidak valid", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		fields := make(map[string]string)
		if req.Email == "" {
			fields["email"] = "Email wajib diisi"
		}
		if req.Password == "" {
			fields["password"] = "Password wajib diisi"
		}
		validation.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email dan password wajib diisi", fields)
		return
	}

	sessionUser, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			validation.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau password yang Anda masukkan salah", nil)
			return
		}
		if errors.Is(err, users.ErrAccountInactive) {
			validation.WriteError(w, http.StatusForbidden, "ACCOUNT_INACTIVE", "Akun Anda telah dinonaktifkan. Silakan hubungi admin.", nil)
			return
		}
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan saat memproses login", nil)
		return
	}

	// Pasang Session Cookie sesuai SSOT §17
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionUser.SessionID.String(),
		Path:     "/",
		Expires:  sessionUser.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	var redirectURL string
	if sessionUser.Role == users.RoleAdmin {
		redirectURL = "/admin"
	} else {
		redirectURL = "/dashboard"
	}

	resp := LoginResponse{
		RedirectURL: redirectURL,
	}
	resp.User.ID = sessionUser.UserID
	resp.User.Name = sessionUser.Name
	resp.User.Email = sessionUser.Email
	resp.User.Role = sessionUser.Role

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": resp,
	})
}

// HandleLogout menangani POST /api/v1/auth/logout (SSOT §60)
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil && cookie.Value != "" {
		if sessionID, err := uuid.Parse(cookie.Value); err == nil {
			_ = h.authSvc.Logout(r.Context(), sessionID)
		}
	}

	// Hapus cookie di browser
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Berhasil keluar (logout)",
	})
}

// HandleGetSession menangani GET /api/v1/auth/session (SSOT §60)
func (h *Handler) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(middleware.UserContextKey)
	user, ok := val.(*middleware.AuthenticatedUser)
	if !ok || user == nil {
		validation.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Tidak ada sesi aktif", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"id":         user.UserID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"session_id": user.SessionID,
			"expires_at": user.ExpiresAt,
		},
	})
}
