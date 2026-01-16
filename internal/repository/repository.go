package repository

import (
	// "database/sql"
	"context"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	// "github.com/redis/go-redis/v9"
)


type Incident interface {
	// GetIncident(req models.IncidentRequest) (int, error)
	CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
	SaveCheck(userID int, lat, lon float64, hasDanger bool,) error
	CreateIncident(incident models.IncidentRequest) error
	GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error)
	GetIncidentById(id int) (models.IncidentResponse, error)
	UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error)
	DeleteIncident(id int) error
	GetDangerStats(ctx context.Context, window time.Duration) (int64, error)
}

type DB interface {
	PingContext(ctx context.Context) error
}

type Redis interface {
	Ping(ctx context.Context) error
}


type Repository struct {
	Incident
	DB
	Redis
}

func NewRepository(db *sqlx.DB, redis RedisPinger) *Repository {
	return &Repository{
		Incident: NewIncidentPostgres(db),
		DB:       db,
		Redis:    redis,
	}
}