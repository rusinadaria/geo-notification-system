package services

import (
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Incident interface {
	CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
	CreateIncident(incidentData models.IncidentRequest) error
	GetIncidentById(id int) (models.Incident, error)
	UpdateIncident(id int, req models.IncidentRequest) (models.Incident, error)
	DeleteIncident(id int) error
}

type Service struct {
	Incident
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Incident: NewIncidentService(repos.Incident),
	}
}