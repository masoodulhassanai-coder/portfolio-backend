package httpapi

import (
	"net/http"
	"portfolio/internal/config"
)

func AdminAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// HARD BLOCK: If the server doesn't have a key configured, block everything
			if cfg.AdminAPIKey == "" {
				http.Error(w, `{"error":"admin key not configured on server"}`, http.StatusInternalServerError)
				return
			}

			authHeader := r.Header.Get("Authorization")
			expected := "Bearer " + cfg.AdminAPIKey
			
			if authHeader != expected {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}