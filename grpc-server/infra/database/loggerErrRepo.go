package database

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type LoggerErrRepo struct {
	db *sql.DB
}

func NewLoggerErrRepo(db *sql.DB) *LoggerInfoRepo {
	return &LoggerInfoRepo{db: db}
}

func (c *LoggerErrRepo) Save(log string) {
	id := uuid.New().String()
	createdAt := time.Now()
	_, err := c.db.Exec(`INSERT INTO log_error (id, err, created_at) 
    VALUES (
        $1, $2, $3
    )`,
		id, log, createdAt)
	if err != nil {
		slog.Error("Error: LoggerErrRepo: can't insert Error Log - ", err)
	}
}
