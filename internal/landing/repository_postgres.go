package landing

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetLandingData(ctx context.Context) (*LandingData, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetLandingData(ctx context.Context) (*LandingData, error) {
	// Ambil template featured dari database jika ada
	const qTemplates = `
		SELECT id, name, slug, category, thumbnail_url, is_featured
		FROM templates
		WHERE status = 'active' AND is_featured = TRUE
		LIMIT 6
	`
	var templates []TemplatePreviewItem
	rows, err := r.db.Query(ctx, qTemplates)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t TemplatePreviewItem
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Category, &t.ThumbnailURL, &t.IsFeatured); err == nil {
				templates = append(templates, t)
			}
		}
	}

	// Ambil reviews approved dari database jika ada
	const qReviews = `
		SELECT display_name, rating, message, to_char(created_at, 'DD Mon YYYY')
		FROM reviews
		WHERE status = 'approved'
		ORDER BY created_at DESC
		LIMIT 6
	`
	var reviews []ReviewItem
	rRows, err := r.db.Query(ctx, qReviews)
	if err == nil {
		defer rRows.Close()
		for rRows.Next() {
			var rev ReviewItem
			if err := rRows.Scan(&rev.Name, &rev.Rating, &rev.Message, &rev.Date); err == nil {
				rev.Role = "Pengantin GolanInvite"
				reviews = append(reviews, rev)
			}
		}
	}

	data := GetDefaultLandingData()
	if len(templates) > 0 {
		data.Templates = templates
	}
	if len(reviews) > 0 {
		data.Reviews = reviews
	}

	return data, nil
}

func GetDefaultLandingData() *LandingData {
	return &LandingData{
		Hero: HeroSection{
			Badge:        "Platform Undangan Digital Modern #1",
			Headline:     "Rayakan Momen Spesial dengan Undangan Digital Berkelas & Custom Domain",
			Subheadline:  "Buat undangan pernikahan eksklusif dengan nama tamu personal, RSVP instan, amplop digital, dan custom domain pribadi Anda sendiri tanpa ribet.",
			PrimaryCTA:   "Pilih Tema & Pesan Sekarang",
			SecondaryCTA: "Lihat Demo Undangan",
			RatingNote:   "Dipercaya oleh 1.500+ pasangan di seluruh Indonesia",
		},
		Features: []FeatureItem{
			{
				Icon:        "Globe",
				Title:       "Custom Domain Eksklusif",
				Description: "Gunakan nama domain Anda sendiri (misal: rani-dan-budi.com) dengan proteksi SSL/HTTPS aktif otomatis.",
			},
			{
				Icon:        "Send",
				Title:       "Generator Nama Tamu WhatsApp",
				Description: "Satu klik buat link undangan personal untuk setiap tamu dan kirim via WhatsApp secara praktis.",
			},
			{
				Icon:        "CheckCircle",
				Title:       "Manajemen RSVP & Ucapan Realtime",
				Description: "Pantau konfirmasi kehadiran tamu serta kelola ucapan dan doa restu langsung dari dashboard Anda.",
			},
			{
				Icon:        "CreditCard",
				Title:       "Amplop Digital & QRIS Terintegrasi",
				Description: "Kemudahan bagi tamu memberikan hadiah pernikahan secara cashless melalui nomor rekening atau QRIS.",
			},
			{
				Icon:        "Music",
				Title:       "Musik Latar & Galeri Foto HD",
				Description: "Sentuhan romantis dengan lagu pilihan dan album galeri kenangan berkualitas tinggi.",
			},
			{
				Icon:        "MapPin",
				Title:       "Integrasi Google Maps Akurat",
				Description: "Panduan navigasi langsung menuju lokasi resepsi dan akad nikah tanpa khawatir tamu tersesat.",
			},
		},
		HowItWorks: []HowItWorksStep{
			{
				Step:        1,
				Title:       "Pilih Desain & Paket",
				Description: "Jelajahi katalog tema pernikahan berkelas kami dan pilih paket yang sesuai kebutuhan acara Anda.",
			},
			{
				Step:        2,
				Title:       "Kirim Detail Acara",
				Description: "Isi form ringkas memuat nama mempelai, tanggal akad/resepsi, lokasi, foto galeri, dan cerita cinta Anda.",
			},
			{
				Step:        3,
				Title:       "Proses Kilat & Publikasi",
				Description: "Tim kami memproses undangan dalam hitungan jam, mengaktifkan custom domain, dan siap disebarkan.",
			},
		},
		Templates: []TemplatePreviewItem{
			{
				ID:           "1",
				Name:         "Elysian Bloom",
				Slug:         "elysian-bloom",
				Category:     "Botanical Modern",
				ThumbnailURL: "https://images.unsplash.com/photo-1519741497674-611481863552?w=800&q=80",
				IsFeatured:   true,
			},
			{
				ID:           "2",
				Name:         "Royal Serenade",
				Slug:         "royal-serenade",
				Category:     "Luxury Gold",
				ThumbnailURL: "https://images.unsplash.com/photo-1511285560929-80b456fea0bc?w=800&q=80",
				IsFeatured:   true,
			},
			{
				ID:           "3",
				Name:         "Minimalist Noir",
				Slug:         "minimalist-noir",
				Category:     "Clean Minimal",
				ThumbnailURL: "https://images.unsplash.com/photo-1465495976277-4387d4b0b4c6?w=800&q=80",
				IsFeatured:   true,
			},
		},
		Reviews: []ReviewItem{
			{
				Name:    "Rani & Dimas",
				Role:    "Pengantin Jakarta",
				Rating:  5,
				Message: "Tamu undangan kami sangat terpukau dengan domain pribadi rani-dimas.com! Fitur RSVP-nya sangat membantu kami mengatur catering resepsi dengan tepat.",
				Date:    "12 Feb 2026",
			},
			{
				Name:    "Sarah & Rizky",
				Role:    "Pengantin Surabaya",
				Rating:  5,
				Message: "Proses pembuatannya sangat cepat, desainnya mewah dan tidak norak. Amplop digital QRIS memudahkan keluarga yang berhalangan hadir untuk tetap berkirim kado.",
				Date:    "28 Jan 2026",
			},
			{
				Name:    "Anisa & Fikri",
				Role:    "Pengantin Bandung",
				Rating:  5,
				Message: "Pelayanan admin GolanInvite sangat responsif! Desainnya responsive dibuka di smartphone mana pun terasa smooth dan elegan.",
				Date:    "05 Jan 2026",
			},
		},
		FAQs: []FAQItem{
			{
				Question: "Berapa lama proses pembuatan undangan digital?",
				Answer:   "Setelah Anda melengkapi data acara dan menyelesaikan pemesanan, undangan akan selesai dalam waktu 1 - 24 jam kerja.",
			},
			{
				Question: "Apakah saya bisa menggunakan domain sendiri seperti nama-kami.com?",
				Answer:   "Tentu saja! Pada paket Platinum, kami menyediakan integrasi Custom Domain (.com / .id) lengkap dengan sertifikat SSL gratis tanpa biaya tersembunyi.",
			},
			{
				Question: "Apakah ada batasan jumlah tamu yang bisa dikirimkan undangan?",
				Answer:   "Tidak ada batasan (Unlimited Guest). Anda dapat membuat link dengan nama tamu sebanyak apa pun menggunakan fitur Generator WhatsApp kami.",
			},
			{
				Question: "Bisakah saya mengubah data atau jadwal jika terjadi perubahan acara?",
				Answer:   "Bisa, Anda memiliki akses dashboard untuk memperbarui data teks, waktu, maupun susunan acara kapan saja selama masa aktif undangan.",
			},
		},
		ContactInfo: map[string]string{
			"email":     "support@golaninvite.com",
			"whatsapp":  "+62 812-3456-7890",
			"instagram": "@golaninvite",
			"address":   "Jakarta Selatan, Indonesia",
		},
	}
}
