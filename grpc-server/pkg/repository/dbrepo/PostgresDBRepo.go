package dbrepo

import (
	"database/sql"
	"grpc-server/pkg/data"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type PostgresDBRepo struct {
	db *sql.DB
}

func NewPostgresDBRepo(db *sql.DB) *PostgresDBRepo {
	return &PostgresDBRepo{db: db}
}

func (c *PostgresDBRepo) SaveOrderHistory(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error) {
	id := uuid.New().String()
	_, err := c.db.Exec(`INSERT INTO orderHistory (id, spread, aexcid, aprice, apricevet, avolume, bexcid, bprice, bpricevet, bvolume, created_at) 
    VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
    )`,
		id, spread, ask.ExcID, ask.Price, ask.PriceVET, ask.Volume,
		bid.ExcID, bid.Price, bid.PriceVET, bid.Volume, createdAt)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c *PostgresDBRepo) SaveLoggerInfo(log string) {
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

func (c *PostgresDBRepo) SaveLoggerErr(log string) {
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
