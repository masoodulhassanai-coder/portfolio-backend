package repo

import (
	"context"
	"portfolio/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{db: db}
}

// ==========================================
// PRODUCTS
// ==========================================
func (r *AdminRepo) CreateProduct(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (slug, name, tagline, description, kind, cover_image_url, tags, status, is_open_source, repo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'published', $8, $9)
	`
	_, err := r.db.Exec(ctx, query, p.Slug, p.Name, p.Tagline, p.Description, p.Kind, p.CoverImageURL, p.Tags, p.IsOpenSource, p.RepoURL)
	return err
}

func (r *AdminRepo) UpdateProduct(ctx context.Context, p *domain.Product) error {
	// Uses COALESCE and NULLIF to only update fields that are not empty
	query := `
		UPDATE products 
		SET 
			name = COALESCE(NULLIF($2, ''), name), 
			tagline = COALESCE(NULLIF($3, ''), tagline), 
			description = COALESCE(NULLIF($4, ''), description), 
			kind = COALESCE(NULLIF($5, ''), kind), 
			cover_image_url = COALESCE(NULLIF($6, ''), cover_image_url), 
			tags = CASE WHEN cardinality($7) > 0 THEN $7 ELSE tags END,
			is_open_source = $8
		WHERE slug = $1
	`
	_, err := r.db.Exec(ctx, query, p.Slug, p.Name, p.Tagline, p.Description, p.Kind, p.CoverImageURL, p.Tags, p.IsOpenSource)
	return err
}

func (r *AdminRepo) DeleteProduct(ctx context.Context, slug string) error {
	query := `DELETE FROM products WHERE slug = $1`
	_, err := r.db.Exec(ctx, query, slug)
	return err
}

// ==========================================
// PLATFORMS (Download Links)
// ==========================================
func (r *AdminRepo) AddPlatform(ctx context.Context, p *domain.ProductPlatform) error {
	query := `
		INSERT INTO product_platforms (product_id, platform, price_cents, currency, external_url, provider, label, version, is_primary, min_os_version)
		VALUES ((SELECT id FROM products WHERE slug = $1), $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query, p.ProductID, p.Platform, p.PriceCents, p.Currency, p.ExternalURL, p.Provider, p.Label, p.Version, p.IsPrimary, p.MinOSVersion)
	return err
}

func (r *AdminRepo) UpdatePlatform(ctx context.Context, p *domain.ProductPlatform) error {
	query := `
		UPDATE product_platforms 
		SET price_cents = $2, external_url = $3, provider = $4, label = $5, version = $6, is_primary = $7, min_os_version = $8
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, p.ID, p.PriceCents, p.ExternalURL, p.Provider, p.Label, p.Version, p.IsPrimary, p.MinOSVersion)
	return err
}

func (r *AdminRepo) DeletePlatform(ctx context.Context, linkID string) error {
	query := `DELETE FROM product_platforms WHERE id = $1`
	_, err := r.db.Exec(ctx, query, linkID)
	return err
}

// ==========================================
// REVIEWS
// ==========================================
func (r *AdminRepo) DeleteReview(ctx context.Context, reviewID string) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.Exec(ctx, query, reviewID)
	return err
}

// ==========================================
// LOOKUP TABLES
// ==========================================
func (r *AdminRepo) CreateLookupPlatform(ctx context.Context, code, displayName string, sortOrder int) error {
	query := `INSERT INTO platforms (code, display_name, sort_order) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, code, displayName, sortOrder)
	return err
}

func (r *AdminRepo) CreateLookupProvider(ctx context.Context, code, displayName string, sortOrder int) error {
	query := `INSERT INTO providers (code, display_name, sort_order) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, code, displayName, sortOrder)
	return err
}