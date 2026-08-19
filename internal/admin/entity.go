package admin

import (
	"github.com/google/uuid"
)

// OrderItem struct
type OrderItem struct {
	ID             uuid.UUID `json:"id"`
	OrderNumber    string    `json:"order_number"`
	CustomerName   string    `json:"customer_name"`
	CustomerEmail  string    `json:"customer_email"`
	CustomerPhone  string    `json:"customer_phone"`
	GroomName      string    `json:"groom_name"`
	BrideName      string    `json:"bride_name"`
	EventDate      string    `json:"event_date"`
	EventLocation  string    `json:"event_location"`
	CustomDomain   string    `json:"custom_domain"`
	PackageName    string    `json:"package_name"`
	Amount         int64     `json:"amount"`
	Status         string    `json:"status"` // pending, processing, accepted, rejected, completed
	PaymentMethod  string    `json:"payment_method"`
	PaymentProofURL string   `json:"payment_proof_url"`
	Notes          string    `json:"notes"`
	CreatedAt      string    `json:"created_at"`
}

// InvitationItem struct
type InvitationItem struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Status       string    `json:"status"` // draft, scheduled, active, expired, disabled
	CustomDomain string    `json:"custom_domain"`
	StartAt      string    `json:"start_at"`
	EndAt        string    `json:"end_at"`
	CreatedAt    string    `json:"created_at"`
}

// TemplateItem struct
type TemplateItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	ThumbnailURL string    `json:"thumbnail_url"`
	DemoURL      string    `json:"demo_url"`
	Status       string    `json:"status"` // draft, active, inactive, archived
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    string    `json:"created_at"`
}

// DomainItem struct
type DomainItem struct {
	ID             uuid.UUID `json:"id"`
	InvitationID   uuid.UUID `json:"invitation_id"`
	InvitationTitle string   `json:"invitation_title"`
	Hostname       string    `json:"hostname"`
	Status         string    `json:"status"` // pending, dns_verified, ssl_pending, active, failed, disabled
	SSLStatus      string    `json:"ssl_status"`
	DNSVerifiedAt  string    `json:"dns_verified_at"`
	CreatedAt      string    `json:"created_at"`
}

// RSVPItem struct
type RSVPItem struct {
	ID               uuid.UUID `json:"id"`
	InvitationTitle  string    `json:"invitation_title"`
	GuestName        string    `json:"guest_name"`
	PhoneNumber      string    `json:"phone_number"`
	AttendanceStatus string    `json:"attendance_status"`
	GuestCount       int       `json:"guest_count"`
	Message          string    `json:"message"`
	CreatedAt        string    `json:"created_at"`
}

// GreetingItem struct
type GreetingItem struct {
	ID              uuid.UUID `json:"id"`
	InvitationTitle string    `json:"invitation_title"`
	GuestName       string    `json:"guest_name"`
	Message         string    `json:"message"`
	Status          string    `json:"status"` // pending, approved, hidden, rejected
	CreatedAt       string    `json:"created_at"`
}

// TransactionItem struct
type TransactionItem struct {
	ID               uuid.UUID `json:"id"`
	OrderNumber      string    `json:"order_number"`
	CustomerName     string    `json:"customer_name"`
	Amount           int64     `json:"amount"`
	PaymentMethod    string    `json:"payment_method"`
	PaymentReference string    `json:"payment_reference"`
	Status           string    `json:"status"`
	ProofURL         string    `json:"proof_url"`
	CreatedAt        string    `json:"created_at"`
}
