package database

import (
	"database/sql"
	"grpc-server/data"

	"github.com/google/uuid"
)

type OrderHistoryRepo struct {
	db *sql.DB
	// id       string
	// spread   float64
	// ask      AskOrder
	// bid      BidOrder
	// createAt time.Time
}

func NewOrderHistory(db *sql.DB) *OrderHistoryRepo {
	return &OrderHistoryRepo{db: db}
}

func (c *OrderHistoryRepo) Save(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error) {
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
