package httpapi

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/domain"
	"portfolio/internal/repo"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	repo *repo.AdminRepo
}

func NewAdminHandler(repo *repo.AdminRepo) *AdminHandler {
	return &AdminHandler{repo: repo}
}

// --- Products ---
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateProduct(r.Context(), &p); err != nil {
		http.Error(w, `{"error":"failed to create product (slug might already exist)"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success", "message":"product created"}`))
}

func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	p.Slug = slug
	if err := h.repo.UpdateProduct(r.Context(), &p); err != nil {
		http.Error(w, `{"error":"failed to update product"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"success", "message":"product updated"}`))
}

func (h *AdminHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if err := h.repo.DeleteProduct(r.Context(), slug); err != nil {
		http.Error(w, `{"error":"failed to delete product"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"success", "message":"product deleted"}`))
}

// --- Platforms ---
func (h *AdminHandler) AddPlatform(w http.ResponseWriter, r *http.Request) {
	var p domain.ProductPlatform
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.AddPlatform(r.Context(), &p); err != nil {
		http.Error(w, `{"error":"failed to add platform"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success", "message":"platform added"}`))
}

func (h *AdminHandler) UpdatePlatform(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p domain.ProductPlatform
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.repo.UpdatePlatform(r.Context(), &p); err != nil {
		http.Error(w, `{"error":"failed to update platform"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"success", "message":"platform updated"}`))
}

func (h *AdminHandler) DeletePlatform(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeletePlatform(r.Context(), id); err != nil {
		http.Error(w, `{"error":"failed to delete platform"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"success", "message":"platform deleted"}`))
}

// --- Reviews ---
func (h *AdminHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteReview(r.Context(), id); err != nil {
		http.Error(w, `{"error":"failed to delete review"}`, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status":"success", "message":"review deleted"}`))
}

// --- Lookup Tables ---
func (h *AdminHandler) CreateLookupPlatform(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code        string `json:"code"`
		DisplayName string `json:"display_name"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateLookupPlatform(r.Context(), req.Code, req.DisplayName, req.SortOrder); err != nil {
		http.Error(w, `{"error":"failed to create platform lookup"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success", "message":"lookup platform created"}`))
}

func (h *AdminHandler) CreateLookupProvider(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code        string `json:"code"`
		DisplayName string `json:"display_name"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateLookupProvider(r.Context(), req.Code, req.DisplayName, req.SortOrder); err != nil {
		http.Error(w, `{"error":"failed to create provider lookup"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"success", "message":"lookup provider created"}`))
}