// Reference to ArdanLabs project template
// https://github.com/ardanlabs/service/blob/master/business/domain/userbus/model.go
package data

import (
	"time"
)

// Layer: Business Layer
// Complex types, parsed from Application layer to guarantee data integrity
type Ask Order

type Bid Order

type Orderbook struct {
	Ask       *Ask
	Bid       *Bid
	CreatedAt time.Time
}

func NewOrderBook(a *Ask, b *Bid) *Orderbook {
	return &Orderbook{
		Ask: a,
		Bid: b,
	}
}
