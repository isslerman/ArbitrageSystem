package repository

import (
	"database/sql"
	"log/slog"
	"looper/pkg/data"
	"sync/atomic"
	"time"
)

type PostgresDBRepo struct {
	db     *sql.DB
	closed atomic.Bool
}

func NewPostgresDBRepo(db *sql.DB) *PostgresDBRepo {
	return &PostgresDBRepo{db: db}
}

func (c *PostgresDBRepo) Stop() {
	c.closed.Store(true)
	_ = c.db.Close()
}

func (c *PostgresDBRepo) StoreOrderHistory(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error) {
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

func (c *PostgresDBRepo) LogInfo(log string) {
	if c.closed.Load() {
		return
	}

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

func (c *PostgresDBRepo) LogError(log string) {
	if c.closed.Load() {
		return
	}

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
