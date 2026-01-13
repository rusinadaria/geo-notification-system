package repository

import (
	// "database/sql"
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
	"log"
	"log/slog"
)

func ConnectDatabase(storage_path string, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", "user=postgres password=root dbname=geo sslmode=disable")
	if err != nil {
		log.Fatal("Failed connect database")
		return nil, err
	}
	logger.Info("Connect database")
	return db, nil
}