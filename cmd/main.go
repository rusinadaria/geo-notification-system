package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"fmt"
	_ "github.com/lib/pq"
	"github.com/rusinadaria/geo-notification-system/internal/config"
	"github.com/rusinadaria/geo-notification-system/internal/handlers"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	redisrepo "github.com/rusinadaria/geo-notification-system/internal/repository/redis"
	"github.com/rusinadaria/geo-notification-system/internal/services"
	"github.com/rusinadaria/geo-notification-system/internal/worker"
	"os/signal"
	"time"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := configLogger()
	cfg, err := config.GetConfig()
	if err != nil {
		return nil
	}

	db, err := repository.ConnectDatabase(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}
	defer db.Close()

	redisClient, err := redisrepo.NewClient(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	repo := repository.NewRepository(db, redisClient)
	queue := redisrepo.NewWebhookQueue(redisClient.Client())

	worker := worker.NewWebhookWorker(redisClient.Client(), cfg.WebhookURL)
	workerCtx, workerCancel := context.WithCancel(ctx)
	go worker.Run(workerCtx)

	srv := services.NewService(repo, cfg.WindowMin, queue)
	handler := handlers.NewHandler(srv)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: handler.InitRoutes(cfg, logger),
	}

	go func() {
		log.Println("listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("ListenAndServe failed", slog.Any("error", err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workerCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP shutdown failed: %w", err)
	}

	return nil
}

func configLogger() *slog.Logger {
	var logger *slog.Logger

	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("Unable to open a file for writing")
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger = slog.New(slog.NewJSONHandler(f, opts))
	return logger
}
