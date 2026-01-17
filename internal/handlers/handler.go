package handlers

import (
	"github.com/go-chi/chi"
	"github.com/rusinadaria/geo-notification-system/internal/config"
	"github.com/rusinadaria/geo-notification-system/internal/handlers/middleware"
	"github.com/rusinadaria/geo-notification-system/internal/services"
	"log/slog"
	"net/http"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(cfg *config.Config, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {

		// Публичный эндпоинт
		r.Post("/location/check", h.CheckLocation)
		r.Get("/system/health", h.HealthCheck)

		// CRUD для инцидентов
		r.Route("/incidents", func(r chi.Router) {

			r.Use(middleware.APIKeyAuth)

			r.Post("/", h.CreateIncidentHandler)
			r.Get("/", h.ListIncidents)
			r.Get("/{id}", h.GetIncident)
			r.Put("/{id}", h.UpdateIncident)
			r.Delete("/{id}", h.DeleteIncident)
			r.Get("/stats", h.GetStats)
		})
	})

	return r
}
