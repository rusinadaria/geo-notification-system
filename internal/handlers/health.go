package handlers

import (
	"encoding/json"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"net/http"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := h.services.HealthService.Check(r.Context())

	code := http.StatusOK
	if resp.Status != models.HealthOK {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
