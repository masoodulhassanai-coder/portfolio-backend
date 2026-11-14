package domain

import (
	"time"
)

type Product struct {
	ID            string    `json:"id"`
	Slug          string    `json:"slug"`
	Name          string    `json:"name"`
	Tagline       string    `json:"tagline"`
	Description   string    `json:"description"`
	Kind          string    `json:"kind"`
	CoverImageURL string    `json:"cover_image_url"`
	Screenshots   []string  `json:"screenshots"`
	Tags          []string  `json:"tags"`
	Status        string    `json:"status"`
	HomepageURL   string    `json:"homepage_url,omitempty"`
	RepoURL       string    `json:"repo_url,omitempty"`
	DemoURL       string    `json:"demo_url,omitempty"`
	VideoURL      string    `json:"video_url,omitempty"`
	GuideURL      string    `json:"guide_url,omitempty"`
	IsOpenSource  bool      `json:"is_open_source"` // <-- NEW
	CreatedAt     time.Time `json:"created_at"`
}

type ProductPlatform struct {
	ID                  string `json:"id"`
	ProductID           string `json:"product_id"`
	Platform            string `json:"platform"`
	PlatformDisplayName string `json:"platform_display_name"`
	PriceCents          int    `json:"price_cents"`
	Currency            string `json:"currency"`
	ExternalURL         string `json:"external_url"`
	Provider            string `json:"provider"`
	ProviderDisplayName string `json:"provider_display_name"`
	Label               string `json:"label"`
	Version             string `json:"version,omitempty"`
	IsPrimary           bool   `json:"is_primary"`
	DownloadCount       int    `json:"download_count"`
	MinOSVersion        string `json:"min_os_version,omitempty"`
}

type Review struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}