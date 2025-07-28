package repository

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func ConnectToPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Optionally test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	slog.Info("Connected to Postgres")
	return db, nil
}
