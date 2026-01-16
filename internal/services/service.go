package services

import (
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	// "github.com/rusinadaria/geo-notification-system/internal/queue"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Incident interface {
	CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
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

// type WebhookQueue interface {
// 	Enqueue(payload models.WebhookPayload) error
// }

type Service struct {
	Incident
	HealthService
	// WebhookQueue
}

func NewService(repos *repository.Repository, windowMin int) *Service {
	return &Service{
		Incident: NewIncidentService(repos.Incident, windowMin),
		HealthService: NewHealthService(repos.DB, repos.Redis),
	}
}

// func NewService(repos *repository.Repository,  webhookQueue queue.WebhookQueue) *Service {
// 	return &Service{
// 		Incident: NewIncidentService(repos.Incident),
// 		WebhookQueue: &webhookQueue,
// 	}
// }