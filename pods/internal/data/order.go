package data

import "time"

type Order struct {
	Price     float64
	PriceVET  float64
	Volume    float64
	CreatedAt time.Time
}
