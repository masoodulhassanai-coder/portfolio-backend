package service

import (
	"context"

	"portfolio/internal/domain"
	"portfolio/internal/repo"
)

type ProductService struct {
	repo *repo.ProductRepo
}

func NewProductService(repo *repo.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

type PlatformGroup struct {
	Platform            string                   `json:"platform"`
	PlatformDisplayName string                   `json:"platform_display_name"`
	PriceCents          int                      `json:"price_cents"`
	Currency            string                   `json:"currency"`
	Links               []domain.ProductPlatform `json:"links"`
}

type ProductDetailResponse struct {
	Product        *domain.Product `json:"product"`
	Platforms      []PlatformGroup `json:"platforms"`
	Reviews        []domain.Review `json:"reviews"`
	AvgRating      float64         `json:"avg_rating"`
	TotalDownloads int             `json:"total_downloads"`
}

func (s *ProductService) GetProductDetail(ctx context.Context, slug string) (*ProductDetailResponse, error) {
	product, err := s.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if product.Status == "coming_soon" {
		return &ProductDetailResponse{
			Product:        product,
			Platforms:      []PlatformGroup{},
			Reviews:        []domain.Review{},
			AvgRating:      0,
			TotalDownloads: 0,
		}, nil
	}

	platformRows, err := s.repo.GetPlatformsForProduct(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	totalDownloads, _ := s.repo.GetTotalDownloadsForProduct(ctx, product.ID)

	reviews, _ := s.repo.GetReviews(ctx, product.ID)
	
	avg := 0.0
	if len(reviews) > 0 {
		total := 0
		for _, r := range reviews {
			total += r.Rating
		}
		avg = float64(total) / float64(len(reviews))
	}

	return &ProductDetailResponse{
		Product:        product,
		Platforms:      s.groupByPlatform(platformRows),
		Reviews:        reviews,
		AvgRating:      avg,
		TotalDownloads: totalDownloads,
	}, nil
}

func (s *ProductService) AddRating(ctx context.Context, productID string, rating int, ip string) error {
	return s.repo.AddRating(ctx, productID, rating, ip)
}

func (s *ProductService) IncrementDownload(ctx context.Context, linkID string) error {
	return s.repo.IncrementDownload(ctx, linkID)
}

func (s *ProductService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetGlobalStats(ctx)
}

func (s *ProductService) groupByPlatform(rows []domain.ProductPlatform) []PlatformGroup {
	var result []PlatformGroup
	platformMap := make(map[string]*PlatformGroup)

	for _, row := range rows {
		if _, exists := platformMap[row.Platform]; !exists {
			group := &PlatformGroup{
				Platform:            row.Platform,
				PlatformDisplayName: row.PlatformDisplayName,
				PriceCents:          row.PriceCents,
				Currency:            row.Currency,
				Links:               []domain.ProductPlatform{},
			}
			platformMap[row.Platform] = group
			result = append(result, PlatformGroup{})
			result[len(result)-1] = *group
			platformMap[row.Platform] = &result[len(result)-1]
		}
		group := platformMap[row.Platform]
		group.Links = append(group.Links, row)
	}
	return result
}

type ProductSummary struct {
	Slug           string   `json:"slug"`
	Name           string   `json:"name"`
	Tagline        string   `json:"tagline"`
	Kind           string   `json:"kind"`
	Status         string   `json:"status"`
	CoverImage     string   `json:"cover_image_url"`
	Tags           []string `json:"tags"`
	MinPrice       int      `json:"min_price_cents"`
	Currency       string   `json:"currency"`
	TotalDownloads int      `json:"total_downloads"`
}

func (s *ProductService) ListProducts(ctx context.Context) ([]ProductSummary, error) {
	products, err := s.repo.ListPublished(ctx)
	if err != nil {
		return nil, err
	}

	var summaries []ProductSummary
	for _, p := range products {
		minPrice := 0
		currency := "USD"
		totalDownloads := 0

		if p.Status != "coming_soon" {
			platforms, _ := s.repo.GetPlatformsForProduct(ctx, p.ID)
			minPrice = -1
			for _, plat := range platforms {
				if minPrice == -1 || plat.PriceCents < minPrice {
					minPrice = plat.PriceCents
					currency = plat.Currency
				}
			}
			if minPrice == -1 {
				minPrice = 0
			}
			totalDownloads, _ = s.repo.GetTotalDownloadsForProduct(ctx, p.ID)
		}

		summaries = append(summaries, ProductSummary{
			Slug:           p.Slug,
			Name:           p.Name,
			Tagline:        p.Tagline,
			Kind:           p.Kind,
			Status:         p.Status,
			CoverImage:     p.CoverImageURL,
			Tags:           p.Tags,
			MinPrice:       minPrice,
			Currency:       currency,
			TotalDownloads: totalDownloads,
		})
	}
	return summaries, nil
}