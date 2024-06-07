package data

import (
	"errors"
	"time"
)

// OpenOrder - represents an open order that can be sent to an exchange
type OpenOrder struct {
	id        string  // id received from exchange
	pair      string  // pair name
	amount    float64 // amount to buy/sell
	price     float64 // price to buy/sell
	orderType string  // "market" | "limit"
	side      string  // "buy" | "sell"
	createdAt int64   // time in unix the order is created
}

func NewAskOpenOrder(amount, price float64, pair, orderType string) (AskOpenOrder, error) {
	order := OpenOrder{
		id:        "",
		pair:      pair,
		amount:    amount,
		price:     price,
		orderType: orderType,
		side:      "sell",
		createdAt: time.Now().Unix(),
	}

	err := order.validate()
	if err != nil {
		return AskOpenOrder{}, err
	}
	return AskOpenOrder(order), nil
}

func (o *OpenOrder) validate() error {
	if o.orderType != "market" && o.orderType != "limit" {
		return errors.New("invalid order type")
	}
	if o.side != "buy" && o.side != "sell" {
		return errors.New("invalid side")
	}
	return nil
}
