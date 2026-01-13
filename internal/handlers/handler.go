package handlers

import (
	"github.com/rusinadaria/geo-notification-system/internal/services"
	"net/http"
	"github.com/go-chi/chi"
	"log/slog"
	"github.com/rusinadaria/geo-notification-system/internal/handlers/middleware"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {

		// Публичный эндпоинт
		r.Post("/location/check", h.CheckLocation)

		// CRUD для инцидентов
		r.Route("/incidents", func(r chi.Router) { // + middleware для API-key

			r.Use(middleware.APIKeyAuth)

			r.Post("/", h.CreateIncidentHandler)
			// r.Get("/", h.ListIncidents)
			r.Get("/{id}", h.GetIncident)
			r.Put("/{id}", h.UpdateIncident)
			r.Delete("/{id}", h.DeleteIncident)
		})
	})

	return r
}