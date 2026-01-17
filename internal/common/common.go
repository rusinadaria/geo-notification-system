package common

import (
	"encoding/json"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{Errors: message}
	json.NewEncoder(w).Encode(errorResponse)
}
