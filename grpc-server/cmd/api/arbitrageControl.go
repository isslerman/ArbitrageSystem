package main

import (
	"errors"
	"fmt"
	"grpc-server/internal/cex"
	"grpc-server/pkg/data"
	"grpc-server/pkg/repository"
	"log/slog"
	"time"
)

// ArbitrageControl is who control the arbitrage between two exchanges (CEX)
// Launching an ask limit order in exchange A and waiting this order to be executed
// to launch the bid order in exchange B.
type ArbitrageControl struct {
	AskOrder       data.AskOpenOrder // last order open at exchange A
	BidOrder       data.BidOpenOrder // initial order received
	AskOrderStatus data.OrderState   // actual state of an order
	BidOrderStatus data.OrderState   // actual state of an order
	createdAt      time.Time         // time the the arbitrage was created
	AskSymbol      string            // symbol of the ask side formatted for exchange A
	BidSymbol      string            // symbol of the bid side formatted for exchange B
	cexAsk         cex.Cex           // instance of exchange A
	cexBid         cex.Cex           // instance of exchange B
	Threshold      float64           // threshold value to recreate or not the sell order at exchange A Ask
	Profit         float64           // how much is the profit/distance from the bestask price received
	Dryrun         bool              // if true, it will not send the orders to the exchanges
	DB             repository.DatabaseRepo
}

func NewArbitrageControl(excA, excB cex.ID, aSymbol, bSymbol string, db repository.DatabaseRepo) (*ArbitrageControl, error) {
	ac := &ArbitrageControl{
		AskOrderStatus: data.StateWaiting, // state of the ask order
		BidOrderStatus: data.StateWaiting, // state of the bid order
		AskSymbol:      aSymbol,           // exchange that owns the ask order
		BidSymbol:      bSymbol,           // exchange that owns the bid order
		cexAsk:         cex.New(excA),     // exchange A (sell) - ask order
		cexBid:         cex.New(excB),     // exchange B (buy) - bid order
		Threshold:      0.2,               // 0.2 = 0.2%
		Profit:         0.4,               // 0.4 = 0.4% profit
		createdAt:      time.Now(),        // time the arbitrage was created
		Dryrun:         true,              // if true, it will not send the orders to the exchanges
		DB:             db,
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
func (ao *ArbitrageControl) AskOpenOrder(a data.AskOpenOrder) {
	// adjusting the price with the profit
	a.Price = a.Price * (1 + (ao.Profit / 100))
	// is there any ask order created?
	if !ao.hasAskOpenOrders() {
		// !!! JOGAR O DRYRUN PARA DENTRO DAS FUNCOES, LOG FUNCIONA?
		if !ao.Dryrun {
			_, err := ao.createLimitOrder(data.OpenOrder(a), "sell")
			if err != nil {
				// source, msg, context to help
				err := fmt.Sprintf("[error creating limit order], %s, %v", err, a)
				ao.DB.SaveLoggerErr(err)
				return
				// TODO: handle error
			}

		} else { // DRYRUN MODE
			info := fmt.Sprintf("[limitOrderCreated][dryrun], [%s] Price:%.2f Amount:%.2f Side:%s Pair:%s", ao.cexAsk.Id(), a.Price, a.Amount, a.Side, a.Pair)
			ao.DB.SaveLoggerInfo(info)
			ao.AskOrderStatus = data.StateCreated
			ao.AskOrder = a
			info = "[askorder status changed to StateCreated][dryrun]"
			ao.DB.SaveLoggerInfo(info)
		}
	} else {
		// check if the new price is inside the range of the threshold
		// here the threshold value is the newprice +- the threshold value
		// TODO: change this to a method
		if isInsideThreshold(ao.AskOrder.Price, ao.Threshold, a.Price) {
			info := fmt.Sprintf("[price inside threshold], %f, %f, %f, %f", a.Price, ao.Threshold, ao.AskOrder.Price, ao.AskOrder.Price/a.Price)
			ao.DB.SaveLoggerInfo(info)
			return
		} else {
			// Log info outside threshold
			info := fmt.Sprintf("[price outside threshold], %f, %f, %f, %f", a.Price, ao.Threshold, ao.AskOrder.Price, ao.AskOrder.Price/a.Price)
			ao.DB.SaveLoggerInfo(info)

			// if dryrun is TRUE we don't execute orders to cex
			if !ao.Dryrun {
				err := ao.cancelAllAskOrders()
				if err != nil {
					err := fmt.Sprintf("[error cancelling allaskorders], %s", err)
					ao.DB.SaveLoggerErr(err)
					return
				}

				_, err = ao.createLimitOrder(data.OpenOrder(a), "sell")
				if err != nil {
					err := fmt.Sprintf("[error creating limit order], %s, %v", err, a)
					ao.DB.SaveLoggerErr(err)
					return
					// TODO: handle error
				}
			} else { // DRYRUN MODE
				ao.AskOrderStatus = data.StateCancelled
				info := "[cancelling allaskorders][dryrun], "
				ao.DB.SaveLoggerInfo(info)
				info = "[askorder status changed to StateCancelled][dryrun], "
				ao.DB.SaveLoggerInfo(info)

				info = fmt.Sprintf("[limitOrderCreated][dryrun], [%s] Price:%.2f Amount:%.2f Side:%s Pair:%s", ao.cexAsk.Id(), a.Price, a.Amount, a.Side, a.Pair)
				ao.DB.SaveLoggerInfo(info)
				ao.AskOrderStatus = data.StateCreated
				ao.AskOrder = a
				info = "[askorder status changed to StateCreated][dryrun]"
				ao.DB.SaveLoggerInfo(info)

			}
		}

	}
}

// hasAskOpenOrders returns true if there is an ask order created and valid on the exchange
func (ao *ArbitrageControl) hasAskOpenOrders() bool {
	return ao.AskOrderStatus == data.StateCreated
}

// hasBidOpenOrders returns true if there is an ask order created and valid on the exchange
func (ao *ArbitrageControl) hasBidOpenOrders() bool {
	return ao.BidOrderStatus == data.StateCreated
}

// createLimitOrder creates a limit order on the exchange ask | bid
// o OpenOrder
// cexSide string - "ask" | "bid"
func (ao *ArbitrageControl) createLimitOrder(o data.OpenOrder, cexSide string) (string, error) {
	if ao.hasAskOpenOrders() {
		return "", errors.New("ask order already created and open")
	}
	if o.Amount == 0 {
		return "", errors.New("order not created. amount 0")
	}
	if o.Price == 0 {
		return "", errors.New("order not created. price 0")
	}
	if cexSide != "ask" && cexSide != "bid" {
		return "", errors.New("invalid cex side")
	}
	order := &data.OrdersCreateRequest{
		Amount: o.Amount,
		Pair:   o.Pair,
		Price:  o.Price,
		Side:   o.Side,
		Type:   o.OrderType,
	}

	if cexSide == "ask" {
		orderId, err := ao.cexAsk.CreateOrder(order)
		if err != nil {
			slog.Error("Error:", err)
			return "", err
		}
		// setting the status of the order to created
		ao.AskOrderStatus = data.StateCreated

		// update the AskOrder
		ao.AskOrder.Price = order.Price
		ao.AskOrder.Amount = order.Amount

		return orderId, nil
	} else if cexSide == "bid" {
		orderId, err := ao.cexBid.CreateOrder(order)
		if err != nil {
			slog.Error("Error:", err)
			return "", err
		}
		// setting the status of the order to created
		ao.BidOrderStatus = data.StateCreated
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
	ao.AskOrderStatus = data.StateCancelled
	return nil
}

// check if the valueToCheck is between the base with the threshold
// threshold must be a percentage. Ex. 0.4 = 0.4%
// range = [base - threshold, base + threshold]
func isInsideThreshold(base, threshold, valueToCheck float64) bool {
	if (valueToCheck >= ((1 - threshold/100) * base)) && (valueToCheck <= ((1 + threshold/100) * base)) {
		return true
	} else {
		return false
	}
}
