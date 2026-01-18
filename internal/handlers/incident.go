package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rusinadaria/geo-notification-system/internal/common"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

func (h *Handler) CreateIncidentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var incidentData models.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&incidentData); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	err := h.services.CreateIncident(incidentData)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось добавить инцидент")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListIncidents(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	list, err := h.services.GetAllIncidents(limit, offset)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось получить инциденты")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	fmt.Println(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	incident, err := h.services.GetIncidentById(id)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось найти инцидент по id")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incident)
}

func (h *Handler) UpdateIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var newIncident models.IncidentRequest

	if err := json.NewDecoder(r.Body).Decode(&newIncident); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if newIncident.RadiusMeters <= 0 {
		http.Error(w, "radius_meters must be > 0", http.StatusBadRequest)
		return
	}

	incident, err := h.services.UpdateIncident(id, newIncident)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось изменить инцидент по id")
		return
	}

	json.NewEncoder(w).Encode(incident)
}

func (h *Handler) DeleteIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.services.DeleteIncident(id)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось удалить инцидент")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.services.GetIncidentStats(ctx)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
