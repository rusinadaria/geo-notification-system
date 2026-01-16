package queue

import (
	// "context"
	// "time"
	"encoding/json"
	"github.com/hibiken/asynq"
	// "github.com/redis/go-redis/v9"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)
const TaskSendWebhook = "webhook:send"

type WebhookQueue struct {
	client *asynq.Client
}

func NewWebhookQueue(redisAddr string) *WebhookQueue {
	return &WebhookQueue{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr: redisAddr,
		}),
	}
}

func (q *WebhookQueue) Enqueue(payload models.WebhookPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(
		TaskSendWebhook,
		data,
		asynq.MaxRetry(5),
		asynq.Queue("webhooks"),
	)

	_, err = q.client.Enqueue(task)
	return err
}

// type Queue struct {
// 	client *redis.Client
// 	name   string
// }

// func NewQueue(client *redis.Client, name string) *Queue {
// 	return &Queue{
// 		client: client,
// 		name:   name,
// 	}
// }

// func (q *Queue) Push(ctx context.Context, payload models.WebhookPayload) error {
// 	data, err := json.Marshal(paylod)
// 	if err != nil {
// 		return err
// 	}

// 	return q.client.RPush(ctx, q.name, data).Err()
// }

// func (q *Queue) PopBlocking(ctx context.Context, timeout time.Duration) (string, error) {
// 	result, err := q.client.BLPop(ctx, timeout, q.name).Result()
// 	if err == redis.Nil {
// 		return "", nil
// 	}
// 	if err != nil {
// 		return "", err
// 	}

// 	// result = [queue_name, value]
// 	return result[1], nil
// }