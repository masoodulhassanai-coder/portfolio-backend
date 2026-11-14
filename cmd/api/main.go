package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"portfolio/internal/config"
	"portfolio/internal/httpapi"
	"portfolio/internal/repo"
	"portfolio/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Connect to Database
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}
	fmt.Println("✅ Successfully connected to PostgreSQL!")

	// 3. Wire up the layers
	productRepo := repo.NewProductRepo(dbpool)
	productService := service.NewProductService(productRepo)
	emailService := service.NewEmailService(cfg)
	
	// 4. Setup Router (PASS CFG AND DBPOOL HERE)
	router := httpapi.NewRouter(cfg, productService, emailService, dbpool)

	// 5. Start Server
	server := &http.Server{
		Addr:    "0.0.0.0" + cfg.Port,
		Handler: router,
	}

	fmt.Printf("🚀 Server starting on http://localhost:%s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v\n", err)
	}
}