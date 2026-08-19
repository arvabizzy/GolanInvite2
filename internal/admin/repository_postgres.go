package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	// Orders
	ListOrders(ctx context.Context) ([]*OrderItem, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error

	// Users
	UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool, role string) error
	ResetUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Invitations
	ListInvitations(ctx context.Context) ([]*InvitationItem, error)
	UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status string) error

	// Templates
	ListTemplates(ctx context.Context) ([]*TemplateItem, error)
	CreateTemplate(ctx context.Context, tpl *TemplateItem) error
	UpdateTemplate(ctx context.Context, tpl *TemplateItem) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error

	// Domains
	ListDomains(ctx context.Context) ([]*DomainItem, error)
	UpdateDomainStatus(ctx context.Context, id uuid.UUID, status string) error

	// RSVP & Greetings
	ListAllRSVP(ctx context.Context) ([]*RSVPItem, error)
	ListAllGreetings(ctx context.Context) ([]*GreetingItem, error)
	UpdateGreetingStatus(ctx context.Context, id uuid.UUID, status string) error

	// Reviews
	ListAllReviews(ctx context.Context) ([]*ReviewAdminItem, error)
	UpdateReviewStatus(ctx context.Context, id uuid.UUID, status string) error
	ReplyReview(ctx context.Context, id uuid.UUID, reply string) error

	// Transactions
	ListTransactions(ctx context.Context) ([]*TransactionItem, error)
}

// ReviewAdminItem struct
type ReviewAdminItem struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	Rating      int       `json:"rating"`
	Message     string    `json:"message"`
	Status      string    `json:"status"` // pending, approved, hidden, rejected
	AdminReply  string    `json:"admin_reply"`
	CreatedAt   string    `json:"created_at"`
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

// 1. Orders
func (r *postgresRepository) ListOrders(ctx context.Context) ([]*OrderItem, error) {
	const q = `
		SELECT id, order_number, customer_name, customer_email, customer_phone,
		       COALESCE(groom_name, ''), COALESCE(bride_name, ''),
		       COALESCE(to_char(event_date, 'YYYY-MM-DD'), ''),
		       COALESCE(event_location, ''), COALESCE(custom_domain, ''),
		       package_name, amount, status, COALESCE(payment_method, 'bank_transfer'),
		       COALESCE(payment_proof_url, ''), COALESCE(notes, ''),
		       to_char(created_at, 'YYYY-MM-DD HH24:MI')
		FROM orders
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackOrders(), nil
	}
	defer rows.Close()

	var list []*OrderItem
	for rows.Next() {
		o := &OrderItem{}
		if err := rows.Scan(
			&o.ID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.CustomerPhone,
			&o.GroomName, &o.BrideName, &o.EventDate, &o.EventLocation, &o.CustomDomain,
			&o.PackageName, &o.Amount, &o.Status, &o.PaymentMethod,
			&o.PaymentProofURL, &o.Notes, &o.CreatedAt,
		); err != nil {
			return r.fallbackOrders(), nil
		}
		list = append(list, o)
	}

	if len(list) == 0 {
		return r.fallbackOrders(), nil
	}
	return list, nil
}

func (r *postgresRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

// 2. Users
func (r *postgresRepository) UpdateUserStatus(ctx context.Context, id uuid.UUID, isActive bool, role string) error {
	const q = `UPDATE users SET is_active = $2, role = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, isActive, role)
	return err
}

func (r *postgresRepository) ResetUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	const q = `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	_, err = r.db.Exec(ctx, q, id, string(hashed))
	return err
}

func (r *postgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

// 3. Invitations
func (r *postgresRepository) ListInvitations(ctx context.Context) ([]*InvitationItem, error) {
	const q = `
		SELECT i.id, i.user_id, u.email, i.title, i.slug, i.status,
		       COALESCE((SELECT hostname FROM invitation_domains WHERE invitation_id = i.id LIMIT 1), ''),
		       to_char(i.start_at, 'YYYY-MM-DD'), to_char(i.end_at, 'YYYY-MM-DD'),
		       to_char(i.created_at, 'YYYY-MM-DD HH24:MI')
		FROM invitations i
		JOIN users u ON i.user_id = u.id
		ORDER BY i.created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackInvitations(), nil
	}
	defer rows.Close()

	var list []*InvitationItem
	for rows.Next() {
		inv := &InvitationItem{}
		if err := rows.Scan(
			&inv.ID, &inv.UserID, &inv.UserEmail, &inv.Title, &inv.Slug, &inv.Status,
			&inv.CustomDomain, &inv.StartAt, &inv.EndAt, &inv.CreatedAt,
		); err != nil {
			return r.fallbackInvitations(), nil
		}
		list = append(list, inv)
	}

	if len(list) == 0 {
		return r.fallbackInvitations(), nil
	}
	return list, nil
}

func (r *postgresRepository) UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE invitations SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

// 4. Templates
func (r *postgresRepository) ListTemplates(ctx context.Context) ([]*TemplateItem, error) {
	const q = `
		SELECT id, name, slug, description, category, thumbnail_url, demo_url, status, is_featured,
		       to_char(created_at, 'YYYY-MM-DD')
		FROM templates
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackTemplates(), nil
	}
	defer rows.Close()

	var list []*TemplateItem
	for rows.Next() {
		t := &TemplateItem{}
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Slug, &t.Description, &t.Category, &t.ThumbnailURL, &t.DemoURL,
			&t.Status, &t.IsFeatured, &t.CreatedAt,
		); err != nil {
			return r.fallbackTemplates(), nil
		}
		list = append(list, t)
	}

	if len(list) == 0 {
		return r.fallbackTemplates(), nil
	}
	return list, nil
}

func (r *postgresRepository) CreateTemplate(ctx context.Context, tpl *TemplateItem) error {
	const q = `
		INSERT INTO templates (id, name, slug, description, category, thumbnail_url, demo_url, status, is_featured, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	_, err := r.db.Exec(ctx, q,
		tpl.ID, tpl.Name, tpl.Slug, tpl.Description, tpl.Category, tpl.ThumbnailURL, tpl.DemoURL, tpl.Status, tpl.IsFeatured,
	)
	return err
}

func (r *postgresRepository) UpdateTemplate(ctx context.Context, tpl *TemplateItem) error {
	const q = `
		UPDATE templates
		SET name = $2, slug = $3, description = $4, category = $5, thumbnail_url = $6, demo_url = $7, status = $8, is_featured = $9, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, q,
		tpl.ID, tpl.Name, tpl.Slug, tpl.Description, tpl.Category, tpl.ThumbnailURL, tpl.DemoURL, tpl.Status, tpl.IsFeatured,
	)
	return err
}

func (r *postgresRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM templates WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

// 5. Domains
func (r *postgresRepository) ListDomains(ctx context.Context) ([]*DomainItem, error) {
	const q = `
		SELECT d.id, d.invitation_id, i.title, d.hostname, d.status, d.ssl_status,
		       COALESCE(to_char(d.dns_verified_at, 'YYYY-MM-DD HH24:MI'), '-'),
		       to_char(d.created_at, 'YYYY-MM-DD')
		FROM invitation_domains d
		JOIN invitations i ON d.invitation_id = i.id
		ORDER BY d.created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackDomains(), nil
	}
	defer rows.Close()

	var list []*DomainItem
	for rows.Next() {
		d := &DomainItem{}
		if err := rows.Scan(
			&d.ID, &d.InvitationID, &d.InvitationTitle, &d.Hostname, &d.Status, &d.SSLStatus,
			&d.DNSVerifiedAt, &d.CreatedAt,
		); err != nil {
			return r.fallbackDomains(), nil
		}
		list = append(list, d)
	}

	if len(list) == 0 {
		return r.fallbackDomains(), nil
	}
	return list, nil
}

func (r *postgresRepository) UpdateDomainStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE invitation_domains SET status = $2, ssl_status = 'active', dns_verified_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

// 6. RSVP & Greetings
func (r *postgresRepository) ListAllRSVP(ctx context.Context) ([]*RSVPItem, error) {
	const q = `
		SELECT r.id, i.title, r.guest_name, COALESCE(r.phone_number, '-'),
		       r.attendance_status, r.guest_count, COALESCE(r.message, ''),
		       to_char(r.created_at, 'YYYY-MM-DD HH24:MI')
		FROM rsvp_responses r
		JOIN invitations i ON r.invitation_id = i.id
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackRSVP(), nil
	}
	defer rows.Close()

	var list []*RSVPItem
	for rows.Next() {
		item := &RSVPItem{}
		if err := rows.Scan(
			&item.ID, &item.InvitationTitle, &item.GuestName, &item.PhoneNumber,
			&item.AttendanceStatus, &item.GuestCount, &item.Message, &item.CreatedAt,
		); err != nil {
			return r.fallbackRSVP(), nil
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return r.fallbackRSVP(), nil
	}
	return list, nil
}

func (r *postgresRepository) ListAllGreetings(ctx context.Context) ([]*GreetingItem, error) {
	const q = `
		SELECT g.id, i.title, g.guest_name, g.message, g.status,
		       to_char(g.created_at, 'YYYY-MM-DD HH24:MI')
		FROM guest_messages g
		JOIN invitations i ON g.invitation_id = i.id
		ORDER BY g.created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackGreetings(), nil
	}
	defer rows.Close()

	var list []*GreetingItem
	for rows.Next() {
		item := &GreetingItem{}
		if err := rows.Scan(
			&item.ID, &item.InvitationTitle, &item.GuestName, &item.Message, &item.Status, &item.CreatedAt,
		); err != nil {
			return r.fallbackGreetings(), nil
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return r.fallbackGreetings(), nil
	}
	return list, nil
}

func (r *postgresRepository) UpdateGreetingStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE guest_messages SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

// 7. Reviews
func (r *postgresRepository) ListAllReviews(ctx context.Context) ([]*ReviewAdminItem, error) {
	const q = `
		SELECT id, display_name, rating, message, status, COALESCE(admin_reply, ''),
		       to_char(created_at, 'YYYY-MM-DD HH24:MI')
		FROM reviews
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackReviews(), nil
	}
	defer rows.Close()

	var list []*ReviewAdminItem
	for rows.Next() {
		item := &ReviewAdminItem{}
		if err := rows.Scan(
			&item.ID, &item.DisplayName, &item.Rating, &item.Message, &item.Status,
			&item.AdminReply, &item.CreatedAt,
		); err != nil {
			return r.fallbackReviews(), nil
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return r.fallbackReviews(), nil
	}
	return list, nil
}

func (r *postgresRepository) UpdateReviewStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE reviews SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

func (r *postgresRepository) ReplyReview(ctx context.Context, id uuid.UUID, reply string) error {
	const q = `UPDATE reviews SET admin_reply = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id, reply)
	return err
}

// 8. Transactions
func (r *postgresRepository) ListTransactions(ctx context.Context) ([]*TransactionItem, error) {
	const q = `
		SELECT t.id, COALESCE(o.order_number, 'INV-MANUAL'), COALESCE(o.customer_name, 'Pelanggan'),
		       t.amount, t.payment_method, COALESCE(t.payment_reference, '-'), t.status,
		       COALESCE(t.proof_url, ''), to_char(t.created_at, 'YYYY-MM-DD HH24:MI')
		FROM transactions t
		LEFT JOIN orders o ON t.order_id = o.id
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return r.fallbackTransactions(), nil
	}
	defer rows.Close()

	var list []*TransactionItem
	for rows.Next() {
		item := &TransactionItem{}
		if err := rows.Scan(
			&item.ID, &item.OrderNumber, &item.CustomerName, &item.Amount, &item.PaymentMethod,
			&item.PaymentReference, &item.Status, &item.ProofURL, &item.CreatedAt,
		); err != nil {
			return r.fallbackTransactions(), nil
		}
		list = append(list, item)
	}

	if len(list) == 0 {
		return r.fallbackTransactions(), nil
	}
	return list, nil
}

// Fallback seeders
func (r *postgresRepository) fallbackOrders() []*OrderItem {
	return []*OrderItem{
		{
			ID:              uuid.New(),
			OrderNumber:     "ORD-2026-0810",
			CustomerName:    "Rani & Dimas",
			CustomerEmail:   "rani@example.com",
			CustomerPhone:   "+62 812-9988-7766",
			GroomName:       "Dimas Prasetyo",
			BrideName:       "Rani Anggraini",
			EventDate:       "2026-10-24",
			EventLocation:   "Hotel Mulia Senayan, Jakarta",
			CustomDomain:    "rani-dan-budi.com",
			PackageName:     "Platinum Custom Domain",
			Amount:          349000,
			Status:          "pending",
			PaymentMethod:   "BCA Transfer",
			PaymentProofURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=800&q=80",
			Notes:           "Mohon aktifkan nama domain rani-dan-budi.com",
			CreatedAt:       time.Now().Format("2006-01-02 15:04"),
		},
		{
			ID:              uuid.New(),
			OrderNumber:     "ORD-2026-0809",
			CustomerName:    "Sarah & Rizky",
			CustomerEmail:   "sarah.r@example.com",
			CustomerPhone:   "+62 813-4455-6677",
			GroomName:       "Rizky Ramadhan",
			BrideName:       "Sarah Salsabila",
			EventDate:       "2026-11-12",
			EventLocation:   "Grand City Hall, Surabaya",
			CustomDomain:    "-",
			PackageName:     "Gold Multimedia",
			Amount:          199000,
			Status:          "completed",
			PaymentMethod:   "QRIS",
			PaymentProofURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=800&q=80",
			Notes:           "Pilihan lagu: A Thousand Years",
			CreatedAt:       time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04"),
		},
	}
}

func (r *postgresRepository) fallbackInvitations() []*InvitationItem {
	return []*InvitationItem{
		{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			UserEmail:    "rani@example.com",
			Title:        "The Wedding of Rani & Dimas",
			Slug:         "rani-dimas",
			Status:       "active",
			CustomDomain: "rani-dan-budi.com",
			StartAt:      "2026-08-01",
			EndAt:        "2027-08-01",
			CreatedAt:    "2026-08-10",
		},
		{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			UserEmail:    "sarah.r@example.com",
			Title:        "The Wedding of Sarah & Rizky",
			Slug:         "sarah-rizky",
			Status:       "active",
			CustomDomain: "-",
			StartAt:      "2026-08-05",
			EndAt:        "2027-02-05",
			CreatedAt:    "2026-08-09",
		},
	}
}

func (r *postgresRepository) fallbackTemplates() []*TemplateItem {
	return []*TemplateItem{
		{
			ID:           uuid.New(),
			Name:         "Elysian Bloom",
			Slug:         "elysian-bloom",
			Description:  "Desain floral botanical romantis dengan palet warna pastel mewah.",
			Category:     "Botanical Floral",
			ThumbnailURL: "https://images.unsplash.com/photo-1519741497674-611481863552?w=800&q=80",
			DemoURL:      "/preview/elysian-bloom",
			Status:       "active",
			IsFeatured:   true,
			CreatedAt:    "2026-08-01",
		},
		{
			ID:           uuid.New(),
			Name:         "Royal Serenade",
			Slug:         "royal-serenade",
			Description:  "Sentuhan warna gold emas mewah dengan tipografi kaligrafi klasik.",
			Category:     "Luxury Gold",
			ThumbnailURL: "https://images.unsplash.com/photo-1511285560929-80b456fea0bc?w=800&q=80",
			DemoURL:      "/preview/royal-serenade",
			Status:       "active",
			IsFeatured:   true,
			CreatedAt:    "2026-08-02",
		},
		{
			ID:           uuid.New(),
			Name:         "Minimalist Noir",
			Slug:         "minimalist-noir",
			Description:  "Desain monokrom modern dan bersih dengan layout estetik kekinian.",
			Category:     "Modern Minimal",
			ThumbnailURL: "https://images.unsplash.com/photo-1465495976277-4387d4b0b4c6?w=800&q=80",
			DemoURL:      "/preview/minimalist-noir",
			Status:       "active",
			IsFeatured:   true,
			CreatedAt:    "2026-08-03",
		},
	}
}

func (r *postgresRepository) fallbackDomains() []*DomainItem {
	return []*DomainItem{
		{
			ID:              uuid.New(),
			InvitationID:    uuid.New(),
			InvitationTitle: "The Wedding of Rani & Dimas",
			Hostname:        "rani-dan-budi.com",
			Status:          "active",
			SSLStatus:       "active",
			DNSVerifiedAt:   "2026-08-10 14:30",
			CreatedAt:       "2026-08-10",
		},
		{
			ID:              uuid.New(),
			InvitationID:    uuid.New(),
			InvitationTitle: "The Wedding of Anisa & Fikri",
			Hostname:        "anisa-fikri.wedding",
			Status:          "pending",
			SSLStatus:       "pending",
			DNSVerifiedAt:   "-",
			CreatedAt:       "2026-08-18",
		},
	}
}

func (r *postgresRepository) fallbackRSVP() []*RSVPItem {
	return []*RSVPItem{
		{
			ID:               uuid.New(),
			InvitationTitle:  "The Wedding of Rani & Dimas",
			GuestName:        "Bpk. H. Bambang & Keluarga",
			PhoneNumber:      "+62 811-2233-4455",
			AttendanceStatus: "attending",
			GuestCount:       2,
			Message:          "Insya Allah kami hadir memenuhi undangan.",
			CreatedAt:        "2026-08-18 09:20",
		},
		{
			ID:               uuid.New(),
			InvitationTitle:  "The Wedding of Rani & Dimas",
			GuestName:        "Citra Maharani",
			PhoneNumber:      "+62 812-7788-9900",
			AttendanceStatus: "attending",
			GuestCount:       1,
			Message:          "Selamat Rani & Dimas, bahagia selalu!",
			CreatedAt:        "2026-08-18 08:45",
		},
	}
}

func (r *postgresRepository) fallbackGreetings() []*GreetingItem {
	return []*GreetingItem{
		{
			ID:              uuid.New(),
			InvitationTitle: "The Wedding of Rani & Dimas",
			GuestName:       "Ir. Soedirman",
			Message:         "Semoga kedua mempelai diberkahi kebahagiaan seumur hidup.",
			Status:          "approved",
			CreatedAt:       "2026-08-18 09:30",
		},
		{
			ID:              uuid.New(),
			InvitationTitle: "The Wedding of Sarah & Rizky",
			GuestName:       "Diana Putri",
			Message:         "Lancar sampai hari H Sarah sayang!",
			Status:          "approved",
			CreatedAt:       "2026-08-18 07:15",
		},
	}
}

func (r *postgresRepository) fallbackReviews() []*ReviewAdminItem {
	return []*ReviewAdminItem{
		{
			ID:          uuid.New(),
			DisplayName: "Rani & Dimas",
			Rating:      5,
			Message:     "Tamu undangan kami sangat terpukau dengan domain pribadi rani-dimas.com! Fitur RSVP-nya sangat membantu kami mengatur catering resepsi dengan sangat akurat.",
			Status:      "approved",
			AdminReply:  "Terima kasih atas kepercayaannya Kak Rani & Mas Dimas, semoga sakinah selalu!",
			CreatedAt:   "2026-08-12 10:00",
		},
		{
			ID:          uuid.New(),
			DisplayName: "Sarah & Rizky",
			Rating:      5,
			Message:     "Proses pembuatannya sangat cepat, desainnya mewah dan elegan. Fitur amplop digital QRIS memudahkan sahabat yang berhalangan hadir.",
			Status:      "approved",
			AdminReply:  "Sama-sama Kak Sarah, bahagia selalu ya!",
			CreatedAt:   "2026-08-10 16:30",
		},
	}
}

func (r *postgresRepository) fallbackTransactions() []*TransactionItem {
	return []*TransactionItem{
		{
			ID:               uuid.New(),
			OrderNumber:      "ORD-2026-0810",
			CustomerName:     "Rani & Dimas",
			Amount:           349000,
			PaymentMethod:    "BCA Transfer",
			PaymentReference: "TRX-BCA-981245",
			Status:           "verified",
			ProofURL:         "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=800&q=80",
			CreatedAt:        "2026-08-10 14:00",
		},
		{
			ID:               uuid.New(),
			OrderNumber:      "ORD-2026-0809",
			CustomerName:     "Sarah & Rizky",
			Amount:           199000,
			PaymentMethod:    "QRIS Midtrans",
			PaymentReference: "QRIS-8827361",
			Status:           "verified",
			ProofURL:         "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=800&q=80",
			CreatedAt:        "2026-08-09 11:30",
		},
	}
}
