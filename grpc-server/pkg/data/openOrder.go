package data

import (
	"errors"
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

// OpenOrder is an order that can be sent to an exchange.
type AskOpenOrder OpenOrder
type BidOpenOrder OpenOrder

// OpenOrder - represents an open order that can be sent to an exchange
type OpenOrder struct {
	id        string  // id received from exchange
	Pair      string  // pair name
	Amount    float64 // amount to buy/sell
	Price     float64 // price to buy/sell
	OrderType string  // "market" | "limit"
	Side      string  // "buy" | "sell"
	createdAt int64   // time in unix the order is created
}

func NewAskOpenOrder(amount, price float64, pair, orderType string) (AskOpenOrder, error) {
	order := OpenOrder{
		id:        "",
		Pair:      pair,
		Amount:    amount,
		Price:     price,
		OrderType: orderType,
		Side:      "sell",
		createdAt: time.Now().Unix(),
	}

	err := order.validate()
	if err != nil {
		return AskOpenOrder{}, err
	}
	return AskOpenOrder(order), nil
}

func (o *OpenOrder) validate() error {
	if o.OrderType != "market" && o.OrderType != "limit" {
		return errors.New("invalid order type")
	}
	if o.Side != "buy" && o.Side != "sell" {
		return errors.New("invalid side")
	}
	return nil
}
