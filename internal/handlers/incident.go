package handlers

import (
	"github.com/rusinadaria/geo-notification-system/internal/common"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/go-chi/chi"
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

	if checkReq.UserID == "" {
        // return errors.New("user_id is required")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }
    if checkReq.Lat < -90 || checkReq.Lat > 90 {
        // return errors.New("invalid latitude")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }
    if checkReq.Lon < -180 || checkReq.Lon > 180 {
        // return errors.New("invalid longitude")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	// 3. Применяем алгоритм дял нахождения ближайших зон
	incidents, err := h.services.CheckLocation(checkReq)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось получить информацию о пользователе")
		return
	}

	// 4. Возвращаем ответ

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incidents)
}












func (h *Handler) CreateIncidentHandler (w http.ResponseWriter, r *http.Request) { // Получить информацию о монетках, инвентаре и истории транзакций.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var incidentData models.IncidentRequest
    if err := json.NewDecoder(r.Body).Decode(&incidentData); err != nil {
        common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	err := h.services.CreateIncident(incidentData)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось получить информацию о пользователе")
		return
	}

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(info)
}

// func (h *Handler) ListIncidents(w http.ResponseWriter, r *http.Request) {
// 	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
// 	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	list := make([]Incident, 0)
// 	for _, inc := range incidents {
// 		if inc.Active {
// 			list = append(list, inc)
// 		}
// 	}

// 	end := offset + limit
// 	if offset > len(list) {
// 		offset = len(list)
// 	}
// 	if end > len(list) {
// 		end = len(list)
// 	}

// 	json.NewEncoder(w).Encode(list[offset:end])
// }

func (h *Handler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	incident, err := h.services.GetIncidentById(id)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось найти инцидент по id")
		return
	}

	// if !ok || !incident.Active {
	// 	http.NotFound(w, r)
	// 	return
	// }
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

	incident, err := h.services.UpdateIncident(id, newIncident)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось изменить инцидент по id")
		return
	}

	// incident, ok := incidents[id]
	// if !ok || !incident.Active {
	// 	http.NotFound(w, r)
	// 	return
	// }

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