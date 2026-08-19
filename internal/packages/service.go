package packages

import (
	"context"
	"encoding/json"
	"net/http"

	"golaninvite/internal/validation"
)

type Service interface {
	GetPublicPackages(ctx context.Context) ([]*Package, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPublicPackages(ctx context.Context) ([]*Package, error) {
	pkgs, err := s.repo.FindAllActive(ctx)
	if err != nil {
		return nil, err
	}
	// Fallback data jika DB masih kosong untuk UX yang baik
	if len(pkgs) == 0 {
		return getDefaultPackages(), nil
	}
	return pkgs, nil
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// HandleGetPublicPackages handles GET /api/v1/public/packages (SSOT §66)
func (h *Handler) HandleGetPublicPackages(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.svc.GetPublicPackages(r.Context())
	if err != nil {
		validation.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memuat paket", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": pkgs,
	})
}

func getDefaultPackages() []*Package {
	return []*Package{
		{
			Name:        "Silver",
			Description: "Paket ekonomis dengan fitur esensial untuk acara pernikahan yang intim.",
			Price:       99000,
			Benefits: []string{
				"1 Pilihan Desain Elegan",
				"Masa Aktif 3 Bulan",
				"Hitung Mundur Waktu Acara",
				"Galeri Foto (Hingga 5 Foto)",
				"Google Maps Navigasi Lokasi",
				"RSVP & Kolom Ucapan Tamu",
			},
			IsActive:     true,
			DisplayOrder: 1,
		},
		{
			Name:        "Gold",
			Description: "Paket paling populer dengan kustomisasi lengkap dan fitur multimedia interaktif.",
			Price:       199000,
			Benefits: []string{
				"Semua Fitur Paket Silver",
				"Masa Aktif 6 Bulan",
				"Galeri Foto (Hingga 15 Foto) & Background Music",
				"Amplop Digital & QRIS Pembayaran",
				"Kirim Undangan Nama Tamu Otomatis (WhatsApp Generator)",
				"Cerita Cinta (Love Story Timeline)",
				"Live Streaming Link Integration",
			},
			IsActive:     true,
			DisplayOrder: 2,
		},
		{
			Name:        "Platinum Custom Domain",
			Description: "Paket eksklusif dengan Custom Domain pribadi (.com / .id) untuk pengalaman tak terlupakan.",
			Price:       349000,
			Benefits: []string{
				"Semua Fitur Paket Gold",
				"Custom Domain Pribadi (contoh: rani-dan-budi.com)",
				"Gratis SSL / HTTPS Otomatis",
				"Masa Aktif 1 Tahun Penuh",
				"Galeri Foto Tanpa Batas & Video Embed",
				"Prioritas Dukungan VIP & Revisi Bebas",
				"Statistik Tamu & Export Data RSVP Excel",
			},
			IsActive:     true,
			DisplayOrder: 3,
		},
	}
}
