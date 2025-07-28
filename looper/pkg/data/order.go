package data

import (
	"time"
)

type Order struct {
	ID        string
	Price     float64
	PriceVET  float64
	Volume    float64
	CreatedAt int64
}

func NewOrder(id string, price, priceVET, volume float64, createdAt int64) *Order {
	return &Order{
		ID:        id,
		Price:     price,
		PriceVET:  priceVET,
		Volume:    volume,
		CreatedAt: createdAt,
	}
}

func (o *Order) CreatedAtHuman() string {
	t := time.Unix(o.CreatedAt, 0)
	strDate := t.Format(time.UnixDate)
	return strDate
}
