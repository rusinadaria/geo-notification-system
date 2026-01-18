package handlers

import (
	"context"
	"encoding/json"
	"github.com/rusinadaria/geo-notification-system/internal/common"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"net/http"
)

func (h *Handler) CheckLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var checkReq models.LocationCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	if checkReq.UserID <= 0 {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Пустой user_id")
		return
	}
	if checkReq.Lat < -90 || checkReq.Lat > 90 {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный или пустой lat")
		return
	}
	if checkReq.Lon < -180 || checkReq.Lon > 180 {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный или пустой lan")
		return
	}

	incidents, err := h.services.CheckLocation(context.Background(), checkReq)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось получить информацию о ближайших опасных зонах")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incidents)
}
