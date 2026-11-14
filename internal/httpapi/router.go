package httpapi

import (
	"net/http"
	"time"

	"portfolio/internal/config"
	"portfolio/internal/repo"
	"portfolio/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func NewRouter(cfg *config.Config, productSvc *service.ProductService, emailSvc *service.EmailService, db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(SecurityHeaders)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://yourdomain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	productHandler := NewProductHandler(productSvc)
	contactHandler := NewContactHandler(emailSvc)
	
	adminRepo := repo.NewAdminRepo(db)
	adminHandler := NewAdminHandler(adminRepo)

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Minute))
		r.Get("/api/products", productHandler.ListProducts)
		r.Get("/api/products/{slug}", productHandler.GetProduct)
		r.Get("/api/stats", productHandler.GetStats)
		r.Post("/api/downloads/{link_id}", productHandler.TrackDownload)
	})

	// Strict Public Routes
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, 1*time.Minute))
		r.Post("/api/products/{slug}/reviews", productHandler.PostReview)
		r.Post("/api/contact", contactHandler.SendMessage)
	})

	// ADMIN ROUTES (Protected)
	r.Route("/api/admin", func(r chi.Router) {
		r.Use(AdminAuth(cfg))
		
		// Products
		r.Post("/products", adminHandler.CreateProduct)
		r.Put("/products/{slug}", adminHandler.UpdateProduct)
		r.Delete("/products/{slug}", adminHandler.DeleteProduct)

		// Platforms (Download Links)
		r.Post("/platforms", adminHandler.AddPlatform)
		r.Put("/platforms/{id}", adminHandler.UpdatePlatform)
		r.Delete("/platforms/{id}", adminHandler.DeletePlatform)

		// Reviews
		r.Delete("/reviews/{id}", adminHandler.DeleteReview)

		// Lookup Tables
		r.Post("/lookup/platforms", adminHandler.CreateLookupPlatform)
		r.Post("/lookup/providers", adminHandler.CreateLookupProvider)
	})

	return r
}