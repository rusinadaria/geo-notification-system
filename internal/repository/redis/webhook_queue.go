package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

type WebhookRedisQueue struct {
	client *redis.Client
	key    string
}

func NewWebhookQueue(client *redis.Client) *WebhookRedisQueue {
	return &WebhookRedisQueue{
		client: client,
		key:    "webhooks:queue",
	}
}

func (q *WebhookRedisQueue) Enqueue(
	ctx context.Context,
	job models.WebhookPayload,
) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.LPush(ctx, q.key, data).Err()
}
