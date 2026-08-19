package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golaninvite/internal/packages"
	"golaninvite/internal/validation"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// 1. Orders
func (h *Handler) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.GetOrders(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat daftar pesanan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": orders})
}

func (h *Handler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/orders/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Pesanan tidak valid", nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Status pesanan wajib diisi", nil)
		return
	}

	if err := h.svc.SetOrderStatus(r.Context(), id, req.Status); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui status pesanan", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Status pesanan berhasil diperbarui"})
}

// 2. Users Update/Reset/Delete
func (h *Handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/users/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Pengguna tidak valid", nil)
		return
	}

	var req struct {
		IsActive bool   `json:"is_active"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format JSON tidak valid", nil)
		return
	}

	if err := h.svc.SetUserStatus(r.Context(), id, req.IsActive, req.Role); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui pengguna", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Data pengguna berhasil diperbarui"})
}

func (h *Handler) HandleResetUserPassword(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/users/")
	idStr = strings.TrimSuffix(idStr, "/reset-password")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Pengguna tidak valid", nil)
		return
	}

	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.NewPassword) < 6 {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Password baru minimal 6 karakter", nil)
		return
	}

	if err := h.svc.ResetPassword(r.Context(), id, req.NewPassword); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mereset password", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Password berhasil direset"})
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/users/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Pengguna tidak valid", nil)
		return
	}

	if err := h.svc.RemoveUser(r.Context(), id); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus pengguna", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Pengguna berhasil dihapus"})
}

// 3. Invitations
func (h *Handler) HandleListInvitations(w http.ResponseWriter, r *http.Request) {
	invs, err := h.svc.GetInvitations(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat daftar undangan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": invs})
}

func (h *Handler) HandleUpdateInvitationStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/invitations/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Undangan tidak valid", nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Status wajib diisi", nil)
		return
	}

	if err := h.svc.SetInvitationStatus(r.Context(), id, req.Status); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengubah status undangan", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Status undangan berhasil diperbarui"})
}

// 4. Templates
func (h *Handler) HandleListTemplates(w http.ResponseWriter, r *http.Request) {
	tpls, err := h.svc.GetTemplates(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat template", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": tpls})
}

func (h *Handler) HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	var tpl TemplateItem
	if err := json.NewDecoder(r.Body).Decode(&tpl); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format JSON tidak valid", nil)
		return
	}
	if tpl.Name == "" || tpl.Slug == "" {
		validation.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Nama dan slug template wajib diisi", nil)
		return
	}

	if err := h.svc.AddTemplate(r.Context(), &tpl); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menyimpan template baru", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Template berhasil ditambahkan"})
}

func (h *Handler) HandleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/templates/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Template tidak valid", nil)
		return
	}
	if err := h.svc.RemoveTemplate(r.Context(), id); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus template", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Template berhasil dihapus"})
}

// 5. Packages
func (h *Handler) HandleListPackages(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.svc.GetPackages(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat paket", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": pkgs})
}

func (h *Handler) HandleCreatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg packages.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format JSON tidak valid", nil)
		return
	}
	if pkg.Name == "" || pkg.Price < 0 {
		validation.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Nama dan harga paket wajib diisi dengan benar", nil)
		return
	}

	if err := h.svc.AddPackage(r.Context(), &pkg); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat paket baru", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Paket berhasil dibuat"})
}

func (h *Handler) HandleDeletePackage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/packages/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Paket tidak valid", nil)
		return
	}
	if err := h.svc.RemovePackage(r.Context(), id); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus paket", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Paket berhasil dihapus"})
}

func (h *Handler) HandleUpdatePackage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/packages/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Paket tidak valid", nil)
		return
	}

	var pkg packages.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Format JSON tidak valid", nil)
		return
	}
	pkg.ID = id

	if err := h.svc.EditPackage(r.Context(), &pkg); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui paket", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Paket berhasil diperbarui"})
}

// 9. Reviews
func (h *Handler) HandleListReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.svc.GetReviews(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat ulasan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": reviews})
}

func (h *Handler) HandleModerateReview(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/reviews/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Ulasan tidak valid", nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Status wajib diisi", nil)
		return
	}

	if err := h.svc.SetReviewStatus(r.Context(), id, req.Status); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memoderasi ulasan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Status ulasan berhasil diperbarui"})
}

func (h *Handler) HandleReplyReview(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/reviews/")
	idStr = strings.TrimSuffix(idStr, "/reply")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Ulasan tidak valid", nil)
		return
	}

	var req struct {
		Reply string `json:"reply"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Reply == "" {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Balasan wajib diisi", nil)
		return
	}

	if err := h.svc.AddReviewReply(r.Context(), id, req.Reply); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menyimpan balasan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Balasan berhasil disimpan"})
}

// 6. Domains
func (h *Handler) HandleListDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := h.svc.GetDomains(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat domain", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": domains})
}

func (h *Handler) HandleActivateDomain(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/domains/")
	idStr = strings.TrimSuffix(idStr, "/activate")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Domain tidak valid", nil)
		return
	}

	if err := h.svc.SetDomainStatus(r.Context(), id, "active"); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengaktifkan domain", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Domain berhasil diaktifkan dan diverifikasi SSL"})
}

// 7. RSVP & Greetings
func (h *Handler) HandleListAllRSVP(w http.ResponseWriter, r *http.Request) {
	rsvp, err := h.svc.GetRSVP(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat RSVP", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": rsvp})
}

func (h *Handler) HandleListAllGreetings(w http.ResponseWriter, r *http.Request) {
	greetings, err := h.svc.GetGreetings(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat ucapan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": greetings})
}

func (h *Handler) HandleModerateGreeting(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/greetings/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := uuid.Parse(idStr)
	if err != nil {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_ID", "ID Ucapan tidak valid", nil)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		validation.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Status wajib diisi", nil)
		return
	}

	if err := h.svc.SetGreetingStatus(r.Context(), id, req.Status); err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memoderasi ucapan", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Status ucapan berhasil diperbarui"})
}

// 8. Transactions
func (h *Handler) HandleListTransactions(w http.ResponseWriter, r *http.Request) {
	txs, err := h.svc.GetTransactions(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat transaksi", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": txs})
}
