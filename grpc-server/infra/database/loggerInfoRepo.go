package database

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type LoggerInfoRepo struct {
	db *sql.DB
}

func NewLoggerInfoRepo(db *sql.DB) *LoggerInfoRepo {
	return &LoggerInfoRepo{db: db}
}

func (c *LoggerInfoRepo) Save(log string) {
	id := uuid.New().String()
	createdAt := time.Now()
	_, err := c.db.Exec(`INSERT INTO log_info (id, info, created_at) 
    VALUES (
        $1, $2, $3
    )`,
		id, log, createdAt)
	if err != nil {
		slog.Error("Error: LoggerInfoRepo: can't insert Info Log - ", err)
	}

}
