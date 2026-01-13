package cmd

import (
	"log"
	"log/slog"
	"github.com/rusinadaria/geo-notification-system/internal/config"
	"github.com/rusinadaria/geo-notification-system/internal/handlers"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"github.com/rusinadaria/geo-notification-system/internal/services"
	"net/http"
	"os"

	// "github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := configLogger()
	cfg := config.GetConfig()

	db, err := repository.ConnectDatabase(cfg.DBPath, logger)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	repo := repository.NewRepository(db)
	srv := services.NewService(repo)
	handler := handlers.NewHandler(srv)

	err = http.ListenAndServe(cfg.Port, handler.InitRoutes(logger))
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