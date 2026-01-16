package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/rusinadaria/geo-notification-system/internal/config"
	"github.com/rusinadaria/geo-notification-system/internal/handlers"

	// "github.com/rusinadaria/geo-notification-system/internal/queue"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	redisrepo "github.com/rusinadaria/geo-notification-system/internal/repository/redis"
	"github.com/rusinadaria/geo-notification-system/internal/services"

	// "github.com/ilyakaznacheev/cleanenv"

	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	logger := configLogger()
	cfg := config.GetConfig()

	db, err := repository.ConnectDatabase(cfg, logger)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	ctx := context.Background()
	redisClient, err := redisrepo.NewClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatal("Не удалось подключиться к redis:", err)
	}
	// client := queue.NewWebhookQueue("redisAddr")

	windowMin, err := strconv.Atoi(os.Getenv("STATS_TIME_WINDOW_MINUTES"))
	if err != nil {
		windowMin = 15
	}

	repo := repository.NewRepository(db, redisClient)
	srv := services.NewService(repo, windowMin)
	handler := handlers.NewHandler(srv)

	err = http.ListenAndServe(":8080", handler.InitRoutes(cfg, logger))
	if err != nil {
		log.Fatal("Не удалось запустить сервер:", err)
	}
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