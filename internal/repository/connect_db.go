package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rusinadaria/geo-notification-system/internal/config"
	"log"
	"log/slog"
)

func ConnectDatabase(cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DBPath)
	if err != nil {
		log.Fatal("Failed connect database")
		return nil, err
	}
	logger.Info("Connect database")
	return db, nil
}
