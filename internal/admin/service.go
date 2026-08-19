package admin

import (
	"context"

	"github.com/google/uuid"
	"golaninvite/internal/packages"
)

type Service interface {
	// Orders
	GetOrders(ctx context.Context) ([]*OrderItem, error)
	SetOrderStatus(ctx context.Context, id uuid.UUID, status string) error

	// Users
	SetUserStatus(ctx context.Context, id uuid.UUID, isActive bool, role string) error
	ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error
	RemoveUser(ctx context.Context, id uuid.UUID) error

	// Invitations
	GetInvitations(ctx context.Context) ([]*InvitationItem, error)
	SetInvitationStatus(ctx context.Context, id uuid.UUID, status string) error

	// Templates
	GetTemplates(ctx context.Context) ([]*TemplateItem, error)
	AddTemplate(ctx context.Context, tpl *TemplateItem) error
	EditTemplate(ctx context.Context, tpl *TemplateItem) error
	RemoveTemplate(ctx context.Context, id uuid.UUID) error

	// Domains
	GetDomains(ctx context.Context) ([]*DomainItem, error)
	SetDomainStatus(ctx context.Context, id uuid.UUID, status string) error

	// RSVP & Greetings
	GetRSVP(ctx context.Context) ([]*RSVPItem, error)
	GetGreetings(ctx context.Context) ([]*GreetingItem, error)
	SetGreetingStatus(ctx context.Context, id uuid.UUID, status string) error

	// Transactions
	GetTransactions(ctx context.Context) ([]*TransactionItem, error)

	// Packages
	GetPackages(ctx context.Context) ([]*packages.Package, error)
	AddPackage(ctx context.Context, pkg *packages.Package) error
	EditPackage(ctx context.Context, pkg *packages.Package) error
	RemovePackage(ctx context.Context, id uuid.UUID) error

	// Reviews
	GetReviews(ctx context.Context) ([]*ReviewAdminItem, error)
	SetReviewStatus(ctx context.Context, id uuid.UUID, status string) error
	AddReviewReply(ctx context.Context, id uuid.UUID, reply string) error
}

type service struct {
	repo         Repository
	packagesRepo packages.Repository
}

func NewService(repo Repository, packagesRepo packages.Repository) Service {
	return &service{
		repo:         repo,
		packagesRepo: packagesRepo,
	}
}

func (s *service) GetOrders(ctx context.Context) ([]*OrderItem, error) {
	return s.repo.ListOrders(ctx)
}

func (s *service) SetOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateOrderStatus(ctx, id, status)
}

func (s *service) SetUserStatus(ctx context.Context, id uuid.UUID, isActive bool, role string) error {
	return s.repo.UpdateUserStatus(ctx, id, isActive, role)
}

func (s *service) ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	return s.repo.ResetUserPassword(ctx, id, newPassword)
}

func (s *service) RemoveUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *service) GetInvitations(ctx context.Context) ([]*InvitationItem, error) {
	return s.repo.ListInvitations(ctx)
}

func (s *service) SetInvitationStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateInvitationStatus(ctx, id, status)
}

func (s *service) GetTemplates(ctx context.Context) ([]*TemplateItem, error) {
	return s.repo.ListTemplates(ctx)
}

func (s *service) AddTemplate(ctx context.Context, tpl *TemplateItem) error {
	if tpl.ID == uuid.Nil {
		tpl.ID = uuid.New()
	}
	return s.repo.CreateTemplate(ctx, tpl)
}

func (s *service) EditTemplate(ctx context.Context, tpl *TemplateItem) error {
	return s.repo.UpdateTemplate(ctx, tpl)
}

func (s *service) RemoveTemplate(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, id)
}

func (s *service) GetDomains(ctx context.Context) ([]*DomainItem, error) {
	return s.repo.ListDomains(ctx)
}

func (s *service) SetDomainStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateDomainStatus(ctx, id, status)
}

func (s *service) GetRSVP(ctx context.Context) ([]*RSVPItem, error) {
	return s.repo.ListAllRSVP(ctx)
}

func (s *service) GetGreetings(ctx context.Context) ([]*GreetingItem, error) {
	return s.repo.ListAllGreetings(ctx)
}

func (s *service) SetGreetingStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateGreetingStatus(ctx, id, status)
}

func (s *service) GetTransactions(ctx context.Context) ([]*TransactionItem, error) {
	return s.repo.ListTransactions(ctx)
}

func (s *service) GetPackages(ctx context.Context) ([]*packages.Package, error) {
	return s.packagesRepo.FindAll(ctx)
}

func (s *service) AddPackage(ctx context.Context, pkg *packages.Package) error {
	if pkg.ID == uuid.Nil {
		pkg.ID = uuid.New()
	}
	return s.packagesRepo.Create(ctx, pkg)
}

func (s *service) EditPackage(ctx context.Context, pkg *packages.Package) error {
	return s.packagesRepo.Update(ctx, pkg)
}

func (s *service) RemovePackage(ctx context.Context, id uuid.UUID) error {
	return s.packagesRepo.Delete(ctx, id)
}

func (s *service) GetReviews(ctx context.Context) ([]*ReviewAdminItem, error) {
	return s.repo.ListAllReviews(ctx)
}

func (s *service) SetReviewStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateReviewStatus(ctx, id, status)
}

func (s *service) AddReviewReply(ctx context.Context, id uuid.UUID, reply string) error {
	return s.repo.ReplyReview(ctx, id, reply)
}
