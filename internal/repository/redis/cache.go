package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

const activeIncidentsKey = "incidents:active"

type IncidentCache struct {
	client *redis.Client
}

func NewIncidentCache(client *redis.Client) *IncidentCache {
	return &IncidentCache{client: client}
}

func (c *IncidentCache) GetActive(ctx context.Context) ([]models.IncidentResponse, error) {
	data, err := c.client.Get(ctx, activeIncidentsKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var incidents []models.IncidentResponse
	if err := json.Unmarshal(data, &incidents); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (c *IncidentCache) SetActive(
	ctx context.Context,
	incidents []models.IncidentResponse,
	ttl time.Duration,
) error {
	data, err := json.Marshal(incidents)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, activeIncidentsKey, data, ttl).Err()
}

func (c *IncidentCache) InvalidateActive(ctx context.Context) error {
	return c.client.Del(ctx, activeIncidentsKey).Err()
}
