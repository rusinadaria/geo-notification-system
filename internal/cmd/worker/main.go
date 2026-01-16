package worker

// import (
// 	"github.com/hibiken/asynq"
// 	"github.com/rusinadaria/geo-notification-system/internal/webhook"
// 	"github.com/rusinadaria/geo-notification-system/internal/queue"
// )

//
// func NewWebhookServer(redisAddr string) *asynq.Server {
// 	return asynq.NewServer(
// 		asynq.RedisClientOpt{Addr: redisAddr},
// 		asynq.Config{
// 			Concurrency: 10,
// 			Queues: map[string]int{
// 				"webhooks": 10,
// 			},
// 			RetryDelayFunc: asynq.DefaultRetryDelayFunc,
// 		},
// 	)
// }

//
// func RunWorker() error {
// 	srv := NewWebhookServer("localhost:6379")

// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc(queue.TaskSendWebhook, webhook.HandleTask)

// 	return srv.Run(mux)
// }