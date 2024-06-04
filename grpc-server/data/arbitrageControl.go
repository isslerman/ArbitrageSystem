package data

import (
	"errors"
	"grpc-server/internal/cex"
	"grpc-server/internal/cex/data"
	"log/slog"
	"time"
)

// OrderStatus represents the state of a order
type OrderState int

const (
	// the order is waiting to be sent for execution (initial state)
	StateWaiting OrderState = iota
	// successful order that has been created
	StateCreated
	// successful order that has been partially filled.
	StatePartiallyFilled
	// cancelled order that has been cancelled.
	StateCancelled
	// failed order that has been failed.
	StateFailed
)

// ArbitrageControl is who control the arbitrage between two exchanges (CEX)
// Launching an ask limit order in exchange A and waiting this order to be executed
// to launch the bid order in exchange B.
type ArbitrageControl struct {
	InitialAskOrder   AskOrder   // initial order received
	InitialBidOrder   BidOrder   // initial order received
	ArbitrageAskOrder AskOrder   // order modified to be used to send to the exchange
	ArbitrageBidOrder BidOrder   // order modified to be used to send to the exchange
	AskOrderStatus    OrderState // actual state of an order
	BidOrderStatus    OrderState // actual state of an order
	createdAt         time.Time  // time the the arbitrage was created
	askSymbol         string
	bidSymbol         string
}

func NewArbitrageControl(a AskOrder, b BidOrder, askSymbol, bidSymbol string) *ArbitrageControl {
	return &ArbitrageControl{
		InitialAskOrder: a,
		InitialBidOrder: b,
		AskOrderStatus:  StateWaiting,
		BidOrderStatus:  StateWaiting,
		askSymbol:       askSymbol,
		bidSymbol:       bidSymbol,
		createdAt:       time.Now(),
	}
}

// hasAskOpenOrders returns true if there is an ask order created and valid on the exchange
func (ao *ArbitrageControl) hasAskOpenOrders() bool {
	return ao.AskOrderStatus == StateCreated
}

// hasBidOpenOrders returns true if there is an ask order created and valid on the exchange
func (ao *ArbitrageControl) hasBidOpenOrders() bool {
	return ao.BidOrderStatus == StateCreated
}

func (ao *ArbitrageControl) createAskOrder(id cex.ID) (string, error) {
	if ao.ArbitrageAskOrder.Volume == 0 {
		return "", errors.New("ask order not created. amount 0")
	}
	if ao.ArbitrageAskOrder.PriceVET == 0 {
		return "", errors.New("ask order not created. price 0")
	}

	// New instance of the exchange
	exc := cex.New(id)

	order := &data.OrdersCreateRequest{
		Amount: ao.ArbitrageAskOrder.Volume,
		Pair:   ao.askSymbol,
		Price:  ao.ArbitrageAskOrder.PriceVET,
		Side:   "sell",
		Type:   "limit",
	}

	// creating an order
	orderId, err := exc.CreateOrder(order)
	if err != nil {
		slog.Error("Error:", err)
		return "", err
	}

	// setting the status of the order to created
	ao.AskOrderStatus = StateCreated

	return orderId, nil
}
