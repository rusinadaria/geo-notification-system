package services

import (
	"context"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
)

//go:generate mockgen -destination=./mocks/mock.go -source=service.go -package=mocks

type Incident interface {
	CheckLocation(ctx context.Context, checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
	CreateIncident(incidentData models.IncidentRequest) error
	GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error)
	GetIncidentById(id int) (models.IncidentResponse, error)
	UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error)
	DeleteIncident(id int) error
	GetIncidentStats(ctx context.Context) (models.IncidentStatsResponse, error)
}

type HealthService interface {
	Check(ctx context.Context) models.HealthResponse
}

type WebhookQueue interface {
	Enqueue(ctx context.Context, job models.WebhookPayload) error
}

type Service struct {
	Incident
	HealthService
	WebhookQueue
}

func NewService(repos *repository.Repository, windowMin int, webhookQueue WebhookQueue) *Service {
	return &Service{
		Incident:      NewIncidentService(repos.Incident, windowMin, webhookQueue),
		HealthService: NewHealthService(repos.DB, repos.Redis),
	}
}
