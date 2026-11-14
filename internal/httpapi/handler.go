package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"portfolio/internal/service"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	resp, err := h.svc.GetProductDetail(r.Context(), slug)
	if err != nil {
		http.Error(w, `{"error":"product not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to fetch products"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
	})
}

func (h *ProductHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.GetStats(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *ProductHandler) PostReview(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	
	product, err := h.svc.GetProductDetail(r.Context(), slug) 
	if err != nil {
		http.Error(w, `{"error":"product not found"}`, http.StatusNotFound)
		return
	}

	var req struct {
		Rating int `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	ipHash := fmt.Sprintf("%x", sha256.Sum256([]byte(r.RemoteAddr)))
	err = h.svc.AddRating(r.Context(), product.Product.ID, req.Rating, ipHash)
	if err != nil {
		http.Error(w, `{"error":"failed to save rating"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *ProductHandler) TrackDownload(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "link_id")
	
	// Increment in background so it doesn't block the user
	go func() {
		_ = h.svc.IncrementDownload(context.Background(), linkID)
	}()

	w.WriteHeader(http.StatusOK)
}