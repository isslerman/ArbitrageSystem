package main

import (
	"errors"
	"grpc-server/internal/cex"
	"grpc-server/internal/cex/data"
	"grpc-server/pkg/data"
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

// OpenOrder is an order that can be sent to an exchange.
type AskOpenOrder data.OpenOrder
type BidOpenOrder OpenOrder

// ArbitrageControl is who control the arbitrage between two exchanges (CEX)
// Launching an ask limit order in exchange A and waiting this order to be executed
// to launch the bid order in exchange B.
type ArbitrageControl struct {
	AskOrder       AskOpenOrder // initial order received
	BidOrder       BidOpenOrder // initial order received
	AskOrderStatus OrderState   // actual state of an order
	BidOrderStatus OrderState   // actual state of an order
	createdAt      time.Time    // time the the arbitrage was created
	AskSymbol      string       // symbol of the ask side formatted for exchange A
	BidSymbol      string       // symbol of the bid side formatted for exchange B
	cexAsk         cex.Cex      // instance of exchange A
	cexBid         cex.Cex      // instance of exchange B
	Threshold      float64      // threshold value to recreate or not the sell order at exchange A Ask
	Dryrun         bool         // if true, it will not send the orders to the exchanges
}

func NewArbitrageControl(excA, excB cex.ID, aSymbol, bSymbol string) (*ArbitrageControl, error) {
	ac := &ArbitrageControl{
		AskOrderStatus: StateWaiting,  // state of the ask order
		BidOrderStatus: StateWaiting,  // state of the bid order
		AskSymbol:      aSymbol,       // exchange that owns the ask order
		BidSymbol:      bSymbol,       // exchange that owns the bid order
		cexAsk:         cex.New(excA), // exchange A (sell) - ask order
		cexBid:         cex.New(excB), // exchange B (buy) - bid order
		createdAt:      time.Now(),    // time the arbitrage was created
		Dryrun:         true,          // if true, it will not send the orders to the exchanges
	}

	err := ac.validate()
	if err != nil {
		return nil, err
	}
	return ac, nil
}

// validates the ArbitrageControl created
func (ao *ArbitrageControl) validate() error {
	// TODO: validate
	return nil
}

// set a new value received for the AskOpenOrder
func (ao *ArbitrageControl) AskOpenOrder(a AskOpenOrder) {
	// run the new order received ?
	// is there any ask order created?
	if !ao.hasAskOpenOrders() {
		_, err := ao.createLimitOrder(OpenOrder(a), "sell")
		if err != nil {
			slog.Error("error creating limit order %s", err)
			return
			// TODO: handle error
		}
		ao.AskOrder = a
	} else {
		// check if the new price is inside the range of the threshold
		// TODO: change this to a method
		if a.price <= ao.AskOrder.price*ao.Threshold {
			// TODO: handle info log
			return
		} else {
			// if dryrun is TRUE we don't execute orders to cex
			if !ao.Dryrun {
				err := ao.cancelAllAskOrders()
				if err != nil {
					slog.Error("error cancelling all ask orders %s", err)
				}
				_, err = ao.createLimitOrder(OpenOrder(a), "sell")
				if err != nil {
					slog.Error("error creating limit order %s", err)
					return
					// TODO: handle error
				}
			}
			// TODO: handle info log
		}

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

// createLimitOrder creates a limit order on the exchange ask | bid
// o OpenOrder
// cexSide string - "ask" | "bid"
func (ao *ArbitrageControl) createLimitOrder(o OpenOrder, cexSide string) (string, error) {
	if ao.hasAskOpenOrders() {
		return "", errors.New("ask order already created and open")
	}
	if o.amount == 0 {
		return "", errors.New("order not created. amount 0")
	}
	if o.price == 0 {
		return "", errors.New("order not created. price 0")
	}
	if cexSide != "ask" && cexSide != "bid" {
		return "", errors.New("invalid cex side")
	}

	order := &data.OrdersCreateRequest{
		Amount: o.amount,
		Pair:   o.pair,
		Price:  o.price,
		Side:   o.side,
		Type:   o.orderType,
	}

	if cexSide == "ask" {
		orderId, err := ao.cexAsk.CreateOrder(order)
		if err != nil {
			slog.Error("Error:", err)
			return "", err
		}
		// setting the status of the order to created
		ao.AskOrderStatus = StateCreated
		return orderId, nil
	} else if cexSide == "bid" {
		orderId, err := ao.cexBid.CreateOrder(order)
		if err != nil {
			slog.Error("Error:", err)
			return "", err
		}
		// setting the status of the order to created
		ao.BidOrderStatus = StateCreated
		return orderId, nil
	}
	return "", errors.New("invalid cex side")
}

// cancelAllAskOrders cancels all open askorders on the exchange A - ask side
func (ao *ArbitrageControl) cancelAllAskOrders() error {
	// cancel all orders
	err := ao.cexAsk.CancelAllOrders()
	if err != nil {
		return err
	}
	ao.AskOrderStatus = StateCancelled
	return nil
}
