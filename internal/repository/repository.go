package repository

import (
	// "database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)


type Incident interface {
	// GetIncident(req models.IncidentRequest) (int, error)
	CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
	SaveCheck(userID string, lat, lon float64, hasDanger bool,) error
	CreateIncident(incident models.IncidentRequest) error
	GetIncidentById(id int) (models.Incident, error)
	UpdateIncident(id int, req models.IncidentRequest) (models.Incident, error)
	DeleteIncident(id int) error
}

type Repository struct {
	Incident
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Incident: NewIncidentPostgres(db),
	}
}