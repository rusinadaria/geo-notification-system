package services

import (
	"context"
	"time"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
)

type healthService struct {
    db    repository.DBPinger
    redis repository.RedisPinger
}

func NewHealthService(db repository.DBPinger, redis repository.RedisPinger) HealthService {
    return &healthService{
        db:    db,
        redis: redis,
    }
}

func (s *healthService) Check(ctx context.Context) models.HealthResponse {
    checks := map[string]string{}
    status := models.HealthOK

    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    if err := s.db.PingContext(ctx); err != nil {
        checks["postgres"] = string(models.HealthDown)
        status = models.HealthDegraded
    } else {
        checks["postgres"] = string(models.HealthOK)
    }

    if err := s.redis.Ping(ctx); err != nil {
        checks["redis"] = string(models.HealthDown)
        status = models.HealthDegraded
    } else {
        checks["redis"] = string(models.HealthOK)
    }

    return models.HealthResponse{
        Status:    status,
        Checks:    checks,
        Timestamp: time.Now().UTC(),
    }
}