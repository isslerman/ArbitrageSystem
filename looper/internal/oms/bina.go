package oms

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"looper/config"
	"looper/pkg/data"
	"os"
	"strings"
	"sync"

	binance "github.com/adshao/go-binance/v2"
)

type Bina struct {
	cfg          *config.Config
	apiBaseURL   string
	id           string
	name         string
	FeeTaker     float64
	FeeMaker     float64
	Limits       data.OrderLimits
	orderTracker *OrderTracker
	client       *binance.Client
	mu           sync.RWMutex
	// cancelFn     context.CancelFunc
}

func NewBina(ctx context.Context, cfg *config.Config) (*Bina, <-chan *ExecutionReport) {
	apiKey := os.Getenv("APIKEY_BINA")
	secretKey := os.Getenv("APISEC_BINA")
	apiKeyER := os.Getenv("APIKEYER_BINA")
	secretKeyER := os.Getenv("APISECER_BINA")

	if apiKey == "" || secretKey == "" || apiKeyER == "" || secretKeyER == "" {
		log.Fatal("APIXXX_BINA not set")
	}

	limits := &data.OrderLimits{
		OrdMinAmount: 10,
	}

	client := binance.NewClient(apiKey, secretKey)

	_, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatalf("Error, NewBina: APIConn, apikey or ip changes, err:%v", err)
	}

	orderTracker := NewOrderTracker()

	// Starting WebSocket
	log.Println("Starting Binance WebSocket")
	ws, err := NewBinaWSConsumer(apiKeyER, secretKeyER)
	if err != nil {
		log.Fatalf("Failed to create WS consumer: %v", err)
	}
	stream, err := ws.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start WS: %v", err)
	}
	log.Println("🛰 Binance WS started successfully")

	// // Start a goroutine to handle incoming messages
	// go func() {
	// 	for msg := range stream {
	// 		log.Printf("executionReport: %+v", msg)
	// 		log.Printf("🔥 Order %d filled at %s", msg.OrderID, msg.LastExecutedPrice)

	// 		// TODO: atualize order tracker, estratégias, etc.
	// 	}
	// 	log.Println("WebSocket handler done")
	// }()

	return &Bina{
		cfg:          cfg,
		apiBaseURL:   "https://api.binance.com/api/v3/",
		id:           "BINA",
		name:         "Binance",
		FeeTaker:     cfg.Fees.Taker,
		FeeMaker:     cfg.Fees.Maker,
		Limits:       *limits,
		client:       client,
		orderTracker: orderTracker,
	}, stream
}

func NewTestBina() *Bina {
	return &Bina{
		orderTracker: NewOrderTracker(),
		mu:           sync.RWMutex{},
	}
}

// func (e *Bina) CloseWS() {
// 	e.cancelFn()
// }

func (e *Bina) Balance(asset string) (amount float64, err error) {
	return 0.0, nil
}

func (e *Bina) Id() string {
	return e.id
}

// CancelOrder -
func (e *Bina) CancelOrder(id string) error {
	log.Fatal("not implemented")
	return nil
}

// CancelAllOrders -
func (e *Bina) CancelAllOrders() error {
	openOrders, err := e.client.NewListOpenOrdersService().
		Symbol("USDTBRL").
		Do(context.Background())
	if err != nil {
		return err
	}

	var errs []error

	for _, o := range openOrders {
		_, err := e.client.NewCancelOrderService().
			Symbol("USDTBRL").
			OrderID(o.OrderID).
			Do(context.Background())
		if err != nil {
			errs = append(errs, fmt.Errorf("cancel failed: id=%d err=%w", o.OrderID, err))
			continue
		}

		e.orderTracker.RemoveByID(fmt.Sprintf("%d", o.OrderID))
	}

	return combineErrors(errs)
}

// OrdersCreate -
func (e *Bina) CreateOrder(o *data.OrdersCreateRequest) (string, error) {

	var side binance.SideType
	if o.Side == "sell" {
		side = binance.SideTypeSell
	} else {
		if o.Side == "buy" {
			side = binance.SideTypeBuy
		} else {
			log.Fatal("bina.go|CreateOrder()| Invalid side order.")
		}
	}

	qty := fmt.Sprint(o.Amount)
	price := fmt.Sprintf("%v", o.Price)

	// validation - amount less than min allowed
	if o.Amount < 0.001 {
		return "", fmt.Errorf("bina.go|CreateOrder()| Bina amount less than min allowed. order: %v", o)
	}
	if o.Amount > 20 {
		log.Fatalf("bina.go|CreateOrder()| WARNING: order with too high amount: 3")
	}

	if e.cfg.Trade.DryRun {
		err := e.client.NewCreateOrderService().
			Symbol(o.Pair).
			Side(side).Type(binance.OrderTypeLimit).
			TimeInForce(binance.TimeInForceTypeGTC).
			Quantity(qty).
			Price(price).
			Test(context.Background())

		if err != nil {
			return "", fmt.Errorf("dry run failed: %w", err)
		}
		e.orderTracker.Set(OrderSideAsk, "", o.Price, o.Amount)

		slog.Info("Dry run successful (no real order placed)")
		return "dry-run-id", nil
	} else {
		order, err := e.client.NewCreateOrderService().
			Symbol(o.Pair).
			Side(side).Type(binance.OrderTypeLimit).
			TimeInForce(binance.TimeInForceTypeGTC).
			Quantity(qty).
			Price(price).
			Do(context.Background())

		if err != nil {
			return "", fmt.Errorf("real order failed: %w", err)
		}
		e.orderTracker.Set(OrderSideAsk, fmt.Sprintf("%d", order.OrderID), o.Price, o.Amount)

		slog.Info("Real order placed", "order", order)
		return fmt.Sprintf("%d", order.OrderID), nil
	}

}

// OrdersUpdate -
func (e *Bina) UpdateOrder(o *data.OrdersUpdateRequest) error {
	log.Fatal("not implemented")
	return nil
}

func (e *Bina) IsBidOpen() bool {
	return e.orderTracker.IsOpen(OrderSideBid)
}

func (e *Bina) IsAskOpen() bool {
	return e.orderTracker.IsOpen(OrderSideAsk)
}

func (e *Bina) GetAskPrice() (float64, bool) {
	return e.orderTracker.GetOpenPrice(OrderSideAsk)
}

func (e *Bina) GetBidPrice() (float64, bool) {
	return e.orderTracker.GetOpenPrice(OrderSideBid)
}

func combineErrors(errs []error) error {
	var sb strings.Builder
	for _, err := range errs {
		if err != nil {
			sb.WriteString(err.Error() + "; ")
		}
	}
	if sb.Len() == 0 {
		return nil
	}
	return errors.New(sb.String())
}
