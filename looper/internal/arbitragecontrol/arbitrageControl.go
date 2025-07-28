package main

import (
	"errors"
	"fmt"
	"looper/config"
	"looper/infra/ntfy"
	"looper/internal/cex"
	"looper/internal/trading"
	"looper/internal/wsconsumer"
	"looper/pkg/data"
	"looper/pkg/repository"

	"log/slog"
	"time"
)

// ArbitrageControl is who control the arbitrage between two exchanges (CEX)
// Launching an ask limit order in exchange A and waiting this order to be executed
// to launch the bid order in exchange B.
type ArbitrageControl struct {
	AskOrder       data.AskOpenOrder // last order open at exchange maker - ask side // confirm to remove this
	BidOrder       data.BidOpenOrder // last order open at exchange maker - bid side // confirm to remove this
	AskOrderStatus data.OrderState   // actual state of an order
	BidOrderStatus data.OrderState   // actual state of an order
	createdAt      time.Time         // time the the arbitrage was created
	ASymbol        string            // symbol of the ask side formatted for exchange A
	BSymbol        string            // symbol of the bid side formatted for exchange B
	cexA           cex.Cex           // instance of exchange A
	cexB           cex.Cex           // instance of exchange B
	// Threshold      float64
	// Profit      float64 // how much is the profit/distance from the bestask price received
	Dryrun         bool // if true, it will not send the orders to the exchanges
	DB             repository.DatabaseRepo
	Notify         *ntfy.Ntfy // Notify sends notifications to mobile
	consumer       *wsconsumer.WSConsumer
	askOrderFilled chan bool
	bidOrderFilled chan bool
	config         *config.Config
	tmA            *trading.TradingManager
	tmB            *trading.TradingManager
}

func NewArbitrageControl(excA, excB cex.ID, aSymbol, bSymbol string, db repository.DatabaseRepo, ntfy *ntfy.Ntfy, cfg *config) (*ArbitrageControl, error) {
	ac := &ArbitrageControl{
		AskOrderStatus: data.StateWaiting, // state of the ask order
		BidOrderStatus: data.StateWaiting, // state of the bid order
		ASymbol:        aSymbol,           // symbol formatted for exchange A
		BSymbol:        bSymbol,           // symbol formatted for exchange B
		cexA:           cex.New(excA),     // exchange A (sell) - ask order
		cexB:           cex.New(excB),     // exchange B (buy) - bid order
		config:         cfg,
		tmA:            trading.NewTradingManager(),
		tmB:            trading.NewTradingManager(),
		// Threshold:      0.10,              // 0.2 = 0.2%
		// This profit need sum the fee of the selling exchange exchange A plus the profit that you want.
		// 1.30 - tive poucas ordens - mas chegamos no nivel 3.
		// Profit:      0.90,       // 0.4 = 0.4% profit - Let's use 0.4 + 0.50 of fee = 0.9 - let's try 0.95, 1.00, 1.10
		createdAt:      time.Now(), // time the arbitrage was created
		Dryrun:         true,       // if true, it will not send the orders to the exchanges
		DB:             db,
		Notify:         ntfy,            // Notify sends notifications to mobile
		askOrderFilled: make(chan bool), // channel to handle when the consumer have a full filled ask order.
		bidOrderFilled: make(chan bool), // channel to handle when the consumer have a full filled ask order.
	}

	err := ac.validate()
	if err != nil {
		return nil, err
	}

	// websocket consumer to watch monitor when the ask order is filled.
	ac.consumer = wsconsumer.NewWSConsumer(db, ntfy, ac.cexB, ac.askOrderFilled, ac.bidOrderFilled)
	// start the consumer
	go ac.consumer.Start()

	return ac, nil
}

// validates the ArbitrageControl created
func (ao *ArbitrageControl) validate() error {
	// TODO: validate
	return nil
}

// set a new value received for the AskOpenOrder
func (ao *ArbitrageControl) AskOpenOrder(a data.AskOpenOrder) {
	// the price received here is the priceVET ( with fees )
	// adjusting the price with the profit
	slog.Info(fmt.Sprintf("arbitrageControl.go|AskOpenOrder())|ask price: %f, price with fee: %f, profit: %f", a.Price, a.Price*(1+(ao.config.profit/100)), ao.config.profit))

	a.Price = a.Price * (1 + (ao.config.profit / 100) + ao.config.feeMaker)
	// amount of the order to be open on maker side
	a.Amount = ao.config.amount

	// is there any ask order created?
	if !ao.hasAskOpenOrders() {
		_, err := ao.createLimitOrder(data.OpenOrder(a), "ask")
		if err != nil {
			// source, msg, context to help
			err := fmt.Sprintf("[error creating limit order], %s, %v", err, a)
			ao.DB.SaveLoggerErr(err)
			return
			// TODO: handle error
		}
	} else {
		// check if the new price is inside the range of the threshold
		// here the threshold value is the newprice +- the threshold value
		// TODO: change this to a method
		if isInsideThreshold(ao.AskOrder.Price, ao.config.thresholdReset, a.Price) {
			info := fmt.Sprintf("[price inside threshold], %f, %f, %f, %f", a.Price, ao.config.thresholdReset, ao.AskOrder.Price, ao.AskOrder.Price/a.Price)
			ao.DB.SaveLoggerInfo(info)
			return
		} else {
			// Log info outside threshold
			info := fmt.Sprintf("[price outside threshold], %f, %f, %f, %f", a.Price, ao.config.thresholdReset, ao.AskOrder.Price, ao.AskOrder.Price/a.Price)
			ao.DB.SaveLoggerInfo(info)

			ordersOpen := ao.tmA.ListOrders()
			orderOpen := ordersOpen[0]
			slog.Info(fmt.Sprintf("arbitrageControl.go|AskOpenOrder())|cancelOrder orderId: %s", orderOpen.ID))
			err := ao.cancelOrder(orderOpen.ID, "ask")
			if err != nil {
				err := fmt.Sprintf("arbitrageControl.go|AskOpenOrder())|cancelOrder orderId: %s", orderOpen.ID)
				ao.DB.SaveLoggerErr(err)
				return
			}

			_, err = ao.createLimitOrder(data.OpenOrder(a), "ask")
			if err != nil {
				err := fmt.Sprintf("[error creating limit order], %s, %v", err, a)
				ao.DB.SaveLoggerErr(err)
				return
				// TODO: handle error
			}

		}

	}
}

// set a new value received for the AskOpenOrder
func (ao *ArbitrageControl) BidOpenOrder(b data.BidOpenOrder) {
	// the price received here is the priceVET ( with fees )
	// adjusting the price with the profit
	slog.Info(fmt.Sprintf("arbitrageControl.go|BidOpenOrder())|bid price: %f, price with fee: %f, profit: %f", b.Price, b.Price*(1-(ao.config.profit/100)), ao.config.profit))
	b.Price = b.Price / (1 + (ao.config.profit / 100) + ao.config.feeMaker)

	// amount of the order to be open on maker side
	b.Amount = ao.config.amount

	// is there any ask order created?
	if !ao.hasBidOpenOrders() {
		_, err := ao.createLimitOrder(data.OpenOrder(b), "bid")
		if err != nil {
			// source, msg, context to help
			err := fmt.Sprintf("[error creating limit order], %s, %v", err, b)
			ao.DB.SaveLoggerErr(err)
			return
			// TODO: handle error
		}
	} else {
		// check if the new price is inside the range of the threshold
		// here the threshold value is the newprice +- the threshold value
		// TODO: change this to a method
		if isInsideThreshold(ao.BidOrder.Price, ao.config.thresholdReset, b.Price) {
			info := fmt.Sprintf("[price inside threshold], %f, %f, %f, %f", b.Price, ao.config.thresholdReset, ao.BidOrder.Price, ao.BidOrder.Price/b.Price)
			ao.DB.SaveLoggerInfo(info)
			return
		} else {
			// Log info outside threshold
			info := fmt.Sprintf("[price outside threshold], %f, %f, %f, %f", b.Price, ao.config.thresholdReset, ao.BidOrder.Price, ao.AskOrder.Price/b.Price)
			ao.DB.SaveLoggerInfo(info)

			ordersOpen := ao.tmB.ListOrders()
			orderOpen := ordersOpen[0]
			slog.Info(fmt.Sprintf("arbitrageControl.go|BidOpenOrder())|cancelOrder orderId: %s", orderOpen.ID))

			err := ao.cancelOrder(orderOpen.ID, "bid")
			if err != nil {
				err := fmt.Sprintf("arbitrageControl.go|BidOpenOrder())|cancelOrder orderId: %s", orderOpen.ID)
				ao.DB.SaveLoggerErr(err)
				return
			}

			_, err = ao.createLimitOrder(data.OpenOrder(b), "bid")
			if err != nil {
				err := fmt.Sprintf("[error creating limit order], %s, %v", err, b)
				ao.DB.SaveLoggerErr(err)
				return
				// TODO: handle error
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
	if cexSide == "ask" {
		if ao.hasAskOpenOrders() {
			return "", errors.New("ask order already created and open")
		}
	}
	if cexSide == "bid" {
		if ao.hasBidOpenOrders() {
			return "", errors.New("bid order already created and open")
		}
	}
	if o.Amount == 0 {
		return "", errors.New("order not created. amount 0")
	}
	if o.Price == 0 {
		return "", errors.New("order not created. price 0")
	}
	if cexSide != "ask" && cexSide != "bid" {
		return "", fmt.Errorf("invalid cex side: %s", cexSide)
	}

	// { SOLBRL 0.015 867.254388 limit sell 1717899767}
	order := &data.OrdersCreateRequest{
		Amount: o.Amount,
		Pair:   o.Pair,
		Price:  o.Price,
		Side:   o.Side,
		Type:   o.OrderType,
	}

	if cexSide == "ask" && ao.config.askSide {
		orderId, err := ao.cexA.CreateOrder(order)
		if err != nil {
			slog.Error(fmt.Sprintf("arbitrageControl.go|createLimitOrder())|ask - orderId: %s", orderId))
			return "", fmt.Errorf("arbitrageControl.go|createLimitOrder())|ask - orderId: %s", orderId)
		}
		// setting the status of the order to created
		ao.AskOrderStatus = data.StateCreated
		ao.tmA.CreateOrder(orderId)
		slog.Info(fmt.Sprintf("arbitrageControl.go|createLimitOrder())|ask - orderId: %s", orderId))

		// update the AskOrder
		ao.AskOrder.Price = order.Price
		ao.AskOrder.Amount = order.Amount
		ao.AskOrder.Id = orderId

		// start go routine to watch when order was filled
		// not sure if can be the same for ask and bid
		go ao.limitOrderExecuted()

		// notify the user that the order was created
		// ao.Notify.SendMsg("Ask order created", "Ask order created", false)
		return orderId, nil
	}
	if cexSide == "bid" && ao.config.bidSide {
		orderId, err := ao.cexA.CreateOrder(order)
		if err != nil {
			slog.Error(fmt.Sprintf("arbitrageControl.go|createLimitOrder())|bid - orderId: %s", orderId))
			return "", fmt.Errorf("arbitrageControl.go|createLimitOrder())|bid - orderId: %s", orderId)
		}
		// setting the status of the order to created
		ao.BidOrderStatus = data.StateCreated
		ao.tmB.CreateOrder(orderId)
		slog.Info(fmt.Sprintf("arbitrageControl.go|createLimitOrder())|bid - orderId: %s", orderId))

		// update the BidOrder
		ao.BidOrder.Price = order.Price
		ao.BidOrder.Amount = order.Amount
		ao.AskOrder.Id = orderId

		// start go routine to watch when order was filled
		// not sure if can be the same for ask and bid
		go ao.limitOrderExecuted()

		return orderId, nil
	}
	return "", nil
}

// cancelOrder cancel an order given an orderID
func (ao *ArbitrageControl) cancelOrder(orderId, side string) error {
	err := ao.cexA.CancelOrder(orderId)
	if err != nil {
		err := fmt.Errorf("arbitrageControl.go|cancelOrder())| orderId: %s", orderId)
		ao.DB.SaveLoggerErr(err.Error())
		return err
	}
	if side == "ask" {
		ao.AskOrderStatus = data.StateCancelled
		ao.tmA.RemoveOrder(orderId)
	} else if side == "bid" {
		ao.BidOrderStatus = data.StateCancelled
		ao.tmB.RemoveOrder(orderId)
	}
	return nil
}

// cancelAllAskOrders cancels all open askorders on the exchange A - ask side
func (ao *ArbitrageControl) cancelAllOrders() error {
	// cancel all orders
	err := ao.cexA.CancelAllOrders()
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

// go routine to watch when the limit order was fully executed using the channel orderFilled
func (ao *ArbitrageControl) limitOrderExecuted() {
	ao.DB.SaveLoggerInfo("[limitOrderExecuted waiting limit order to be filled]")
	// waiting for the channel
	select {
	case <-ao.askOrderFilled:
		ao.AskOrderStatus = data.StateWaiting
		// we only have one order, so we can get the first one
		orderId := ao.tmA.ListOrders()[0].ID
		ao.tmA.RemoveOrder(orderId)
		// log
		slog.Info(fmt.Sprintf("arbitrageControl.go|limitOrderExecuted())|askOrderFilled orderId: %s", orderId))
		msg := fmt.Sprintf("Order executed. ASK RIPI - %v", ao.AskOrder)
		ao.Notify.SendMsg("Order RIPI-ASK executed", msg, false)
		ao.DB.SaveLoggerInfo("[limitOrderExecuted limit order filled]")
	case <-ao.bidOrderFilled:
		ao.BidOrderStatus = data.StateWaiting
		// we only have one order, so we can get the first one
		orderId := ao.tmB.ListOrders()[0].ID
		ao.tmB.RemoveOrder(orderId)
		// log
		slog.Info(fmt.Sprintf("arbitrageControl.go|limitOrderExecuted())|bidOrderFilled orderId: %s", orderId))
		msg := fmt.Sprintf("Order executed. BID RIPI - %v", ao.BidOrder)
		ao.Notify.SendMsg("Order RIPI-BID executed", msg, false)
		ao.DB.SaveLoggerInfo("[limitOrderExecuted limit order filled]")
	}
}
