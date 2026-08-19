package packages

import (
	"context"

	"github.com/google/uuid"
)

// Repository mendefinisikan interface akses data paket.
type Repository interface {
	FindAllActive(ctx context.Context) ([]*Package, error)
	FindAll(ctx context.Context) ([]*Package, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Package, error)
	Create(ctx context.Context, pkg *Package) error
	Update(ctx context.Context, pkg *Package) error
	Delete(ctx context.Context, id uuid.UUID) error
}
