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


func (h *Handler) CheckLocation (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// 1. Получаем координаты
	var checkReq models.LocationCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
        common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }
	// 2. Валидируем

	if checkReq.UserID <= 0 {
        // return errors.New("user_id is required")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Пустой user_id")
        return
    }
    if checkReq.Lat < -90 || checkReq.Lat > 90 {
        // return errors.New("invalid latitude")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный или пустой lat")
        return
    }
    if checkReq.Lon < -180 || checkReq.Lon > 180 {
        // return errors.New("invalid longitude")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный или пустой lan")
        return
    }

	// 3. Применяем алгоритм для нахождения ближайших зон
	incidents, err := h.services.CheckLocation(checkReq)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось получить информацию о ближайших опасных зонах")
		return
	}

	// 4. Возвращаем ответ

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incidents)
}

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
	// json.NewEncoder(w).Encode(info)
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