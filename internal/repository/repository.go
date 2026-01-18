package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"time"
)

type Incident interface {
	CreateIncident(incident models.IncidentRequest) error
	GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error)
	GetIncidentById(id int) (models.IncidentResponse, error)
	UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error)
	DeleteIncident(id int) error
	GetDangerStats(ctx context.Context, window time.Duration) (int64, error)
	GetActiveIncidents(ctx context.Context) ([]models.IncidentResponse, error)
}

type LocationCheck interface {
	CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error)
	SaveCheck(userID int, lat, lon float64, hasDanger bool) error
}

type IncidentCache interface {
	GetActive(ctx context.Context) ([]models.IncidentResponse, error)
	SetActive(ctx context.Context, incidents []models.IncidentResponse, ttl time.Duration) error
	InvalidateActive(ctx context.Context) error
}

type DB interface {
	PingContext(ctx context.Context) error
}

type Redis interface {
	Ping(ctx context.Context) error
}

type Repository struct {
	Incident
	LocationCheck
	DB
	Redis
}

func NewRepository(db *sqlx.DB, redis RedisPinger) *Repository {
	return &Repository{
		Incident:      NewIncidentPostgres(db),
		LocationCheck: NewLocationCheckPostgres(db),
		DB:            db,
		Redis:         redis,
	}
}
