package packages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindAllActive(ctx context.Context) ([]*Package, error) {
	const q = `
		SELECT id, name, description, price, benefits, COALESCE(badge, ''), is_active, display_order, created_at, updated_at
		FROM packages
		WHERE is_active = TRUE
		ORDER BY display_order ASC, price ASC
	`
	return r.queryList(ctx, q)
}

func (r *postgresRepository) FindAll(ctx context.Context) ([]*Package, error) {
	const q = `
		SELECT id, name, description, price, benefits, COALESCE(badge, ''), is_active, display_order, created_at, updated_at
		FROM packages
		ORDER BY display_order ASC, price ASC
	`
	return r.queryList(ctx, q)
}

func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*Package, error) {
	const q = `
		SELECT id, name, description, price, benefits, COALESCE(badge, ''), is_active, display_order, created_at, updated_at
		FROM packages
		WHERE id = $1
	`
	pkg := &Package{}
	var benefitsRaw []byte
	err := r.db.QueryRow(ctx, q, id).Scan(
		&pkg.ID, &pkg.Name, &pkg.Description, &pkg.Price, &benefitsRaw,
		&pkg.Badge, &pkg.IsActive, &pkg.DisplayOrder, &pkg.CreatedAt, &pkg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("packages: paket tidak ditemukan")
		}
		return nil, fmt.Errorf("packages.repo.FindByID: %w", err)
	}

	if len(benefitsRaw) > 0 {
		_ = json.Unmarshal(benefitsRaw, &pkg.Benefits)
	}
	return pkg, nil
}

func (r *postgresRepository) Create(ctx context.Context, pkg *Package) error {
	const q = `
		INSERT INTO packages (id, name, description, price, benefits, badge, is_active, display_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	benefitsJSON, err := json.Marshal(pkg.Benefits)
	if err != nil {
		benefitsJSON = []byte("[]")
	}

	_, err = r.db.Exec(ctx, q,
		pkg.ID, pkg.Name, pkg.Description, pkg.Price, benefitsJSON,
		pkg.Badge, pkg.IsActive, pkg.DisplayOrder,
	)
	if err != nil {
		return fmt.Errorf("packages.repo.Create: %w", err)
	}
	return nil
}

func (r *postgresRepository) Update(ctx context.Context, pkg *Package) error {
	const q = `
		UPDATE packages
		SET name = $2, description = $3, price = $4, benefits = $5, badge = $6, is_active = $7, display_order = $8, updated_at = NOW()
		WHERE id = $1
	`
	benefitsJSON, err := json.Marshal(pkg.Benefits)
	if err != nil {
		benefitsJSON = []byte("[]")
	}

	tag, err := r.db.Exec(ctx, q,
		pkg.ID, pkg.Name, pkg.Description, pkg.Price, benefitsJSON,
		pkg.Badge, pkg.IsActive, pkg.DisplayOrder,
	)
	if err != nil {
		return fmt.Errorf("packages.repo.Update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("packages: paket tidak ditemukan")
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM packages WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("packages.repo.Delete: %w", err)
	}
	return nil
}

func (r *postgresRepository) queryList(ctx context.Context, query string, args ...interface{}) ([]*Package, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("packages.repo.queryList: %w", err)
	}
	defer rows.Close()

	var list []*Package
	for rows.Next() {
		pkg := &Package{}
		var benefitsRaw []byte
		if err := rows.Scan(
			&pkg.ID, &pkg.Name, &pkg.Description, &pkg.Price, &benefitsRaw,
			&pkg.Badge, &pkg.IsActive, &pkg.DisplayOrder, &pkg.CreatedAt, &pkg.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("packages.repo scan: %w", err)
		}
		if len(benefitsRaw) > 0 {
			_ = json.Unmarshal(benefitsRaw, &pkg.Benefits)
		}
		list = append(list, pkg)
	}
	return list, nil
}
