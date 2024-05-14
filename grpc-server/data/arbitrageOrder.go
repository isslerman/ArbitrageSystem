package data

import (
	"time"

	"github.com/google/uuid"
)

type ArbitrageOrder struct {
	AskOrder    AskOrder
	BidOrder    BidOrder
	ID          string
	HasBeenSent bool
	HasFailed   bool
	createdAt   time.Time
}

func NewArbitrageOrder(a AskOrder, b BidOrder) *ArbitrageOrder {
	return &ArbitrageOrder{
		AskOrder:    a,
		BidOrder:    b,
		ID:          uuid.New().String(),
		HasBeenSent: false,
		HasFailed:   false,
		createdAt:   time.Now(),
	}

}
