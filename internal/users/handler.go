package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golaninvite/internal/validation"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// HandleAdminListUsers menangani GET /api/v1/admin/users (SSOT §62)
func (h *Handler) HandleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	userList, total, err := h.svc.ListUsers(r.Context(), limit, offset)
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil daftar pengguna", nil)
		return
	}

	type UserDTO struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Role      Role   `json:"role"`
		IsActive  bool   `json:"is_active"`
		CreatedAt string `json:"created_at"`
	}

	var dtos []UserDTO
	for _, u := range userList {
		dtos = append(dtos, UserDTO{
			ID:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": dtos,
		"meta": map[string]interface{}{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// HandleAdminCreateUser menangani POST /api/v1/admin/users (SSOT §62)
func (h *Handler) HandleAdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format data tidak valid", nil)
		return
	}

	fields := make(map[string]string)
	if req.Name == "" {
		fields["name"] = "Nama pengguna wajib diisi"
	}
	if req.Email == "" {
		fields["email"] = "Email pengguna wajib diisi"
	}
	if req.Password == "" || len(req.Password) < 6 {
		fields["password"] = "Password minimal 6 karakter"
	}
	if req.Role == "" {
		req.Role = RoleUser
	}
	if !req.Role.IsValid() {
		fields["role"] = "Role harus 'admin' atau 'user'"
	}

	if len(fields) > 0 {
		validation.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Data pengguna tidak lengkap/valid", fields)
		return
	}

	newUser, err := h.svc.CreateUser(r.Context(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			validation.WriteError(w, http.StatusConflict, "EMAIL_EXISTS", "Email sudah terdaftar dalam sistem", nil)
			return
		}
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat pengguna baru", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pengguna berhasil dibuat",
		"data": map[string]interface{}{
			"id":         newUser.ID,
			"name":       newUser.Name,
			"email":      newUser.Email,
			"role":       newUser.Role,
			"is_active":  newUser.IsActive,
			"created_at": newUser.CreatedAt,
		},
	})
}

// HandleAdminOverview menangani GET /api/v1/admin/overview
func (h *Handler) HandleAdminOverview(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetAdminStats(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat statistik admin", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}
