package dbrepo

import (
	"database/sql"
	"grpc-server/pkg/data"
	"log/slog"
	"time"
)

type PostgresDBRepo struct {
	db *sql.DB
}

func NewPostgresDBRepo(db *sql.DB) *PostgresDBRepo {
	return &PostgresDBRepo{db: db}
}

func (c *PostgresDBRepo) SaveOrderHistory(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error) {
	_, err := c.db.Exec(`INSERT INTO orderHistory (spread, aexcid, aprice, apricevet, avolume, bexcid, bprice, bpricevet, bvolume, created_at) 
    VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
    )`,
		spread, ask.ExcID, ask.Price, ask.PriceVET, ask.Volume,
		bid.ExcID, bid.Price, bid.PriceVET, bid.Volume, createdAt)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (c *PostgresDBRepo) SaveLoggerInfo(log string) {
	log = truncate(log, 250)
	createdAt := time.Now()

	_, err := c.db.Exec(`INSERT INTO log_info (info, created_at) 
    VALUES (
        $1, $2
    )`,
		log, createdAt)
	if err != nil {
		slog.Error("Error: LoggerInfoRepo: can't insert Info Log - ", err)
	}

}

func (c *PostgresDBRepo) SaveLoggerErr(log string) {
	log = truncate(log, 250)
	createdAt := time.Now()

	_, err := c.db.Exec(`INSERT INTO log_error (err, created_at) 
    VALUES (
        $1, $2
    )`,
		log, createdAt)
	if err != nil {
		slog.Error("Error: LoggerErrRepo: can't insert Error Log - ", err)
	}
}

func truncate(s string, n int) string {
	// truncating to 250 chars
	if len(s) > n {
		return s[:250]
	} else {
		return s
	}
}
