package repo

import (
	"context"

	"portfolio/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	query := `
		SELECT id, slug, name, coalesce(tagline,''), coalesce(description,''), 
		       kind, coalesce(cover_image_url,''), screenshots, tags, status,
		       coalesce(homepage_url,''), coalesce(repo_url,''), coalesce(demo_url,''), 
		       coalesce(video_url,''), coalesce(guide_url,''), is_open_source, created_at
		FROM products
		WHERE slug = $1 AND status IN ('published', 'coming_soon')
	`
	var p domain.Product
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Name, &p.Tagline, &p.Description,
		&p.Kind, &p.CoverImageURL, &p.Screenshots, &p.Tags, &p.Status,
		&p.HomepageURL, &p.RepoURL, &p.DemoURL, &p.VideoURL, &p.GuideURL, &p.IsOpenSource, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) GetPlatformsForProduct(ctx context.Context, productID string) ([]domain.ProductPlatform, error) {
	query := `
		SELECT pp.id, pp.product_id, 
		       pp.platform, pl.display_name,
		       pp.price_cents, pp.currency,
		       pp.external_url, 
		       pp.provider, pv.display_name,
		       coalesce(pp.label,''), coalesce(pp.version,''), pp.is_primary,
		       pp.download_count, coalesce(pp.min_os_version, '')
		FROM product_platforms pp
		JOIN platforms pl ON pl.code = pp.platform
		JOIN providers pv ON pv.code = pp.provider
		WHERE pp.product_id = $1 AND pp.is_available = true
		ORDER BY pl.sort_order, pp.sort_order
	`
	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []domain.ProductPlatform
	for rows.Next() {
		var pp domain.ProductPlatform
		err := rows.Scan(
			&pp.ID, &pp.ProductID, 
			&pp.Platform, &pp.PlatformDisplayName,
			&pp.PriceCents, &pp.Currency,
			&pp.ExternalURL, 
			&pp.Provider, &pp.ProviderDisplayName,
			&pp.Label, &pp.Version, &pp.IsPrimary,
			&pp.DownloadCount, &pp.MinOSVersion,
		)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, pp)
	}
	return platforms, rows.Err()
}

func (r *ProductRepo) ListPublished(ctx context.Context) ([]domain.Product, error) {
	query := `
		SELECT id, slug, name, coalesce(tagline,''), kind, coalesce(cover_image_url,''), status, tags, created_at
		FROM products
		WHERE status IN ('published', 'coming_soon')
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Name, &p.Tagline,
			&p.Kind, &p.CoverImageURL, &p.Status, &p.Tags, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepo) AddRating(ctx context.Context, productID string, rating int, ipHash string) error {
	query := `
		INSERT INTO reviews (product_id, rating, ip_hash)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, productID, rating, ipHash)
	return err
}

func (r *ProductRepo) GetReviews(ctx context.Context, productID string) ([]domain.Review, error) {
	query := `
		SELECT id, product_id, rating, created_at
		FROM reviews
		WHERE product_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		err := rows.Scan(&rev.ID, &rev.ProductID, &rev.Rating, &rev.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

func (r *ProductRepo) IncrementDownload(ctx context.Context, linkID string) error {
	query := `UPDATE product_platforms SET download_count = download_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, linkID)
	return err
}

func (r *ProductRepo) GetTotalDownloadsForProduct(ctx context.Context, productID string) (int, error) {
	var total int
	query := `SELECT coalesce(sum(download_count), 0) FROM product_platforms WHERE product_id = $1`
	err := r.db.QueryRow(ctx, query, productID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *ProductRepo) GetGlobalStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalProducts int
	_ = r.db.QueryRow(ctx, `SELECT count(*) FROM products WHERE status = 'published'`).Scan(&totalProducts)
	stats["total_products"] = totalProducts

	var totalReviews int
	_ = r.db.QueryRow(ctx, `SELECT count(*) FROM reviews`).Scan(&totalReviews)
	stats["total_reviews"] = totalReviews

	var avgRating float64
	_ = r.db.QueryRow(ctx, `SELECT coalesce(avg(rating), 0) FROM reviews`).Scan(&avgRating)
	stats["avg_rating"] = avgRating

	var totalDownloads int
	_ = r.db.QueryRow(ctx, `SELECT coalesce(sum(download_count), 0) FROM product_platforms`).Scan(&totalDownloads)
	stats["total_downloads"] = totalDownloads

	stats["free_apps"] = totalProducts

	return stats, nil
}