package landing

import (
	"context"
	"encoding/json"
	"net/http"

	"golaninvite/internal/validation"
)

type Service interface {
	GetLandingData(ctx context.Context) (*LandingData, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetLandingData(ctx context.Context) (*LandingData, error) {
	return s.repo.GetLandingData(ctx)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// HandleGetPublicLanding handles GET /api/v1/public/landing (SSOT §66)
func (h *Handler) HandleGetPublicLanding(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetLandingData(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat konten landing page", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
	})
}
