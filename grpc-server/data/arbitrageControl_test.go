package data

import (
	"grpc-server/internal/cex"
	"testing"
)

func Test_CreateAskOrder(t *testing.T) {
	// Binance uses SOL_BRL
	// RIPI uses SOLBRL
	ao := AskOrder{
		ExcID:    "SOL_BRL",
		Price:    9999.00,
		PriceVET: 9999.00,
		Volume:   0.10,
	}

	bo := BidOrder{
		ExcID:    "SOLBRL",
		Price:    9999.00,
		PriceVET: 9999.00,
		Volume:   0.10,
	}

	ac := NewArbitrageControl(ao, bo, ao.ExcID, bo.ExcID)
	_, err := ac.createAskOrder(cex.InstanceRipi)
	if err != nil {
		t.Errorf("error creating ask order: %d", err)
	}
}
