package webhook

// import (
// 	"errors"
// 	"net/http"
// 	"encoding/json"
// 	"context"
// 	"github.com/rusinadaria/geo-notification-system/internal/config"
// 	"github.com/rusinadaria/geo-notification-system/internal/models"
// 	"github.com/hibiken/asynq"
// )
// func HandleSendWebhook(cfg *config.Config) asynq.HandlerFunc {
// 	return func(ctx context.Context, t *asynq.Task) error {
// 		var payload models.WebhookPayload
// 		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
// 			return err
// 		}

// 		body, _ := json.Marshal(payload)

// 		req, _ := http.NewRequest(
// 			http.MethodPost,
// 			cfg.WebhookURL,
// 			bytes.NewBuffer(body),
// 		)
// 		req.Header.Set("Content-Type", "application/json")

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil || resp.StatusCode >= 300 {
// 			return errors.New("webhook failed")
// 		}

// 		return nil
// 	}
// }