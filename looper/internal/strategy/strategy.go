package strategy

import (
	"context"
	"looper/config"
	"looper/internal/marketdata"
	"looper/internal/oms"
	"looper/pkg/data"
	"math"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type MakerStrategy struct {
	cfg         *config.Config
	Logger      *zap.SugaredLogger
	OMS         oms.OMS
	OrderStream <-chan *oms.ExecutionReport
	SCfg        *StrategyConfig
	lastUpdate  time.Time
	bestAsk     float64
	bestBid     float64
}

type StrategyConfig struct {
	base                string
	quote               string
	OrderDistanceBps    float64 // Distance below best ask to place order
	ProfitThresholdBps  float64 // Minimum bps profit required after fees
	RepriceThresholdBps float64 // When price moves enough to warrant cancel & replace
}

func NewMakerStrategy(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) *MakerStrategy {
	ms := MakerStrategy{
		cfg:    cfg,
		Logger: logger,
	}
	ms.OMS, ms.OrderStream = oms.NewBina(ctx, cfg)

	ms.SCfg = &StrategyConfig{
		quote:               cfg.Symbols.Quote,
		base:                cfg.Symbols.Base,
		OrderDistanceBps:    20.0, // Profit Brute - This is the profit + minfee (~19bps)
		RepriceThresholdBps: 30.0,
		ProfitThresholdBps:  20.0,
	}

	return &ms
}

func (s *MakerStrategy) Run(ctx context.Context, marketData <-chan *marketdata.BookTrade, orderEvents <-chan *oms.ExecutionReport) error {
	// Stop the strategy if there is a inactivity timeout for the marketdata
	inactivityTimeout := 30 * time.Second
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Since(s.lastUpdate) > inactivityTimeout {
				s.Logger.Errorw("No marketdata received for 30s. Stopping strategy.")
				return ctx.Err()
			}

		case <-ctx.Done():
			s.Logger.Info("MakerStrategy shutting down")
			return ctx.Err()

		case book := <-marketData:
			s.lastUpdate = time.Now()
			s.bestAsk = book.BestAsk
			s.bestBid = book.BestBid
			minBps := s.minBpsToAskProfit(s.bestAsk, s.bestBid)
			s.Logger.Infow("Current minBps", "minBps", minBps, "price", s.bestAsk*(1.0+minBps/10_000))
			if s.bestAsk <= 0 || s.bestBid <= 0 {
				s.Logger.Warn("Invalid market data received", "ask", s.bestAsk, "bid", s.bestBid)
				continue
			}

			// Is Ask order already open?
			if s.OMS.IsAskOpen() {
				s.Logger.Info("Strategy -> IsAskOpen:True")
				if price, ok := s.OMS.GetAskPrice(); ok {
					s.Logger.Infow("Current ask order price", "price", price)
					drift := bpsBetween(s.bestAsk, price)
					s.Logger.Infow("Current drift", "drift", drift)

					if drift < minBps || drift > s.SCfg.RepriceThresholdBps {
						s.Logger.Infow("Order too far from current ask, cancelling",
							"current_ask", s.bestAsk,
							"order_price", price,
							"bps_drift", drift,
						)

						if err := s.OMS.CancelAllOrders(); err != nil {
							s.Logger.Errorw("Failed to cancel all orders", "err", err)
						}
					} else {
						s.Logger.Infow("Ask order within range. Skipping update")
						continue
					}

				}

			} else {
				s.Logger.Info("Strategy -> IsAskOpen:False")

				// No ask order open. Placing a new one.

				s.Logger.Infow("minBPL calculated", "minBPS", minBps)
				bps := math.Max(minBps, s.SCfg.OrderDistanceBps) // 40bps
				s.Logger.Infow("minBPL adjusted", "minBPS", bps)

				price := s.bestAsk * (1.0 + bps/10_000)
				s.Logger.Infow("price order", "price", price)
				price = truncateTo4Decimals(price)
				s.Logger.Infow("price order truncated", "price", price)

				drift := bpsBetween(s.bestAsk, price)
				s.Logger.Infow("Current drift", "drift", drift)
				if drift < minBps || drift > s.SCfg.RepriceThresholdBps {
					s.Logger.Infow("Order with neg pnl, cancelling",
						"current_ask", s.bestAsk,
						"order_price", price,
						"bps_drift", drift,
					)
				} else {
					order := &data.OrdersCreateRequest{
						Pair:   s.SCfg.Pair(),
						Price:  price,
						Amount: s.cfg.Trade.Amount,
						Side:   oms.OrderSideSell,
						Type:   "limit",
					}

					orderID, err := s.OMS.CreateOrder(order)
					if err != nil {
						s.Logger.Errorw("Failed to create ask order", "err", err)
						return err
					}

					s.Logger.Infow("New ask order created",
						"price", price,
						"bps", bps,
						"order_id", orderID,
					)
				}
			}
		case evt := <-orderEvents:
			s.handleOrderEvent(evt)
		}
	}
}

func (s *MakerStrategy) handleOrderEvent(evt *oms.ExecutionReport) {
	// Log all events
	s.Logger.Infow("Received order event",
		"type", evt.ExecutionType,
		"status", evt.OrderStatus,
		"order_id", evt.OrderID,
		"symbol", evt.Symbol,
	)

	switch evt.OrderStatus {
	case "FILLED":

		s.Logger.Infow("Order filled",
			"id", evt.OrderID,
			"price", evt.LastExecutedPrice,
			"qty", evt.CumulativeFilledQty,
		)

		// 1. Parse fill price and quantity
		fillPrice, err := strconv.ParseFloat(evt.LastExecutedPrice, 64)
		if err != nil {
			s.Logger.Errorw("invalid fill price", "raw", evt.LastExecutedPrice, "err", err)
			return
		}
		fillQty, err := strconv.ParseFloat(evt.CumulativeFilledQty, 64)
		if err != nil {
			s.Logger.Errorw("invalid fill quantity", "raw", evt.CumulativeFilledQty, "err", err)
			return
		}

		hedgeBps := s.SCfg.ProfitThresholdBps
		targetBuyPrice := fillPrice * (1.0 - hedgeBps/10_000)

		s.Logger.Infow("Calculated hedge buy price",
			"fill_price", fillPrice,
			"bps_target", hedgeBps,
			"target_price", targetBuyPrice,
		)

		// 4. Place hedge MARKET order using margin
		order := &data.OrdersCreateRequest{
			Pair:   s.SCfg.Pair(),
			Price:  targetBuyPrice, // optional if it's a true market order
			Amount: fillQty,
			Side:   oms.OrderSideBuy,
			Type:   "market_margin", // custom: your OMS interprets this
		}

		orderID, err := s.OMS.CreateOrder(order)
		if err != nil {
			s.Logger.Errorw("Failed to create hedge order", "err", err)
			return
		}

		s.Logger.Infow("Hedge BID order created",
			"order_id", orderID,
			"price", targetBuyPrice,
			"amount", fillQty,
		)

	case "CANCELED":
		s.Logger.Infow("Order canceled",
			"id", evt.OrderID,
			"client_id", evt.ClientOrderID,
		)

	// optionally: handle others
	case "REJECTED":
		s.Logger.Warnw("Order rejected",
			"reason", evt.OrderRejectReason,
			"id", evt.OrderID,
		)

	case "NEW":
		s.Logger.Infow("Order accepted",
			"id", evt.OrderID,
			"price", evt.Price,
			"side", evt.Side,
		)
	}
}

func (s *MakerStrategy) StopAllOrders() error {
	err := s.OMS.CancelAllOrders()
	if err != nil {
		return err
	}
	return nil
}

func (s *MakerStrategy) GetOrderStream() <-chan *oms.ExecutionReport {
	return s.OrderStream
}

func (c *StrategyConfig) Pair() string {
	return strings.ToUpper(c.base + c.quote)
}

func bpsBetween(current, reference float64) float64 {
	return math.Abs((current - reference) / reference * 10000)
}

// MinBpsToProfit calcs the min BPS above the bestAsk
// to have profit
func (s *MakerStrategy) minBpsToAskProfit(bestAsk, bestBid float64) float64 {
	// Min price to have profit
	costToRebuy := bestBid * (1.0 + s.cfg.Fees.Taker)
	requiredSellPrice := costToRebuy / (1.0 - s.cfg.Fees.Maker)

	// diff to bestAsk converted to BPS
	bps := (requiredSellPrice - bestAsk) / bestAsk * 10_000

	// return positive BPS
	return math.Max(0, bps)
}

func (s *MakerStrategy) minBpsToBidProfit(bestBid, bestAsk float64) float64 {
	sellPrice := bestAsk * (1.0 - s.cfg.Fees.Taker)
	requiredBuyPrice := sellPrice / (1.0 + s.cfg.Fees.Maker)
	bps := (bestBid - requiredBuyPrice) / bestBid * 10_000
	return math.Max(0, bps)
}

func truncateTo4Decimals(f float64) float64 {
	return math.Floor(f*1e4) / 1e4
}
