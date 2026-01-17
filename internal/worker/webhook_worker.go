package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

type WebhookWorker struct {
	client     *redis.Client
	webhookURL string
}

func NewWebhookWorker(client *redis.Client, url string) *WebhookWorker {
	return &WebhookWorker{
		client:     client,
		webhookURL: url,
	}
}

func (w *WebhookWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := w.client.BRPop(ctx, 0, "webhooks:queue").Result()
			if err != nil {
				log.Println(err)
				continue
			}

			payload := res[1]

			var job models.WebhookPayload
			if err := json.Unmarshal([]byte(payload), &job); err != nil {
				log.Println(err)
				continue
			}
			w.send(job)
		}
	}
}

func (w *WebhookWorker) send(job models.WebhookPayload) {
	body, _ := json.Marshal(job)

	req, _ := http.NewRequest("POST", w.webhookURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("webhook http error:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("webhook response status:", resp.Status)

	if resp.StatusCode >= 300 {
		log.Println("webhook failed, retry later")
		return
	}

	log.Println("webhook delivered successfully")
}
