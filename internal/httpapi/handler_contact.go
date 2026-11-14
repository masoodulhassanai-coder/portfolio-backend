package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"portfolio/internal/service"
)

type ContactHandler struct {
	emailSvc *service.EmailService
}

func NewContactHandler(emailSvc *service.EmailService) *ContactHandler {
	return &ContactHandler{emailSvc: emailSvc}
}

type ContactRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	QueryType string `json:"query_type"`
	Message   string `json:"message"`
}

func (h *ContactHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	// NEW: Input Validation
	if len(req.Name) > 100 || len(req.Email) > 100 || len(req.Message) > 5000 {
		http.Error(w, `{"error":"input too long"}`, http.StatusBadRequest)
		return
	}
	
	if !strings.Contains(req.Email, "@") {
		http.Error(w, `{"error":"invalid email"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		http.Error(w, `{"error":"missing required fields"}`, http.StatusBadRequest)
		return
	}

	formattedMessage := fmt.Sprintf("Query Type: %s\n\n%s", req.QueryType, req.Message)

	err := h.emailSvc.SendContactEmail(req.Name, req.Email, req.QueryType, formattedMessage)
	if err != nil {
		http.Error(w, `{"error":"failed to send email"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}