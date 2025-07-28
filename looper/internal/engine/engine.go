package engine

import (
	"context"
	"fmt"
	"looper/config"
	"looper/infra/ntfy"
	"looper/internal/marketdata"
	"looper/internal/repository"
	"looper/pkg/data"
	"time"

	"github.com/gookit/slog"
)

// Our app confiiguration.
// Errors and Infos for debug: All the erros and info msgs are saved to a DB.
type Engine struct {
	cfg        *config.Config // Configuration for the app
	Repo       repository.Repository
	Notify     *ntfy.Ntfy // Notify sends notifications to mobile
	MarketData *marketdata.MarketDataReceiver
	// ac       *ArbitrageControld
}

func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

// bestOrder receives the OrderBook over gRPC
// and find the best ask and bid orders every 2 seconds
func (e *Engine) run(ctx context.Context, ob *data.OrderBook) {
	ticker := time.NewTicker(1500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Warn("run() stopped due to context cancellation")
			return

		case <-ticker.C:
			if ob.SizeAsk() > 0 && ob.SizeBid() > 0 {
				ba := ob.BestAsk()
				bb := ob.BestBid()
				spread := (1 - (ba.PriceVET / bb.PriceVET)) * 100
				createdAt := time.Now().Unix()

				// e.strategyArbitrageContol(ba, bb)

				aho := &data.AskOrder{
					ExcID:    ba.ID,
					Price:    ba.Price,
					PriceVET: ba.PriceVET,
					Volume:   ba.Volume,
				}
				bho := &data.BidOrder{
					ExcID:    bb.ID,
					Price:    bb.Price,
					PriceVET: bb.PriceVET,
					Volume:   bb.Volume,
				}

				if _, err := e.Repo.StoreOrderHistory(spread, aho, bho, createdAt); err != nil {
					e.Repo.LogError(fmt.Sprintf("app.run: error saving to DB: %v", err))
					slog.Error("failed to save order history", "err", err)
				}

				slog.Info("Best Ask", "id", ba.ID, "price", ba.Price, "priceVET", ba.PriceVET, "volume", ba.Volume)
				slog.Info("Best Bid", "id", bb.ID, "price", bb.Price, "priceVET", bb.PriceVET, "volume", bb.Volume)
			}

			ob.RemoveExpiredAsks()
			ob.RemoveExpiredBids()
		}
	}
}

func (e *Engine) GracefullyShutdown() {
	slog.Warn("Shutting down server...")
	e.Repo.LogInfo("ALERT: Shutting down server...")

	slog.Warn("Closing all open orders.")
	e.Repo.LogInfo("ALERT: Closing all open orders.")

	// TODO Change here to OMS
	// if err := e.ac.cancelAllOrders(); err != nil {
	// 	slog.Error("Error canceling all orders", "err", err)
	// 	e.DB.SaveLoggerErr("WARNING: Error canceling all orders, check the exchange for open orders.")
	// } else {
	// 	slog.Info("All orders cancelled successfully.")
	// 	app.DB.SaveLoggerInfo("No error, all orders cancelled.")
	// }

	finalMsg := "---> SERVER STOPPED <---"
	slog.Warn(finalMsg)
	e.Repo.LogInfo(finalMsg)
	e.Repo.LogError(finalMsg)
}

// func (e *Engine) strategyArbitrageContol(ba, bb *data.Order) {
// 	// Log info
// 	info := fmt.Sprintf("[New Price received] [%s] - PriceVET, %f, Vol, %f", ba.ID, ba.PriceVET, ba.Volume)
// 	e.DB.SaveLoggerInfo(info)

// 	// creating ask order
// 	a, err := data.NewAskOpenOrder(ba.Volume, ba.PriceVET, e.ac.ASymbol, "limit")
// 	if err != nil {
// 		err := fmt.Sprintf("[Error creating NewAskOpenOrder] [%s] - %f, %f, %s, sell, %s", ba.ID, ba.Volume, ba.PriceVET, app.ac.ASymbol, err)
// 		e.DB.SaveLoggerErr(err)
// 		return
// 	}

// 	// creating bid order
// 	b, err := data.NewBidOpenOrder(bb.Volume, bb.PriceVET, e.ac.ASymbol, "limit")
// 	if err != nil {
// 		err := fmt.Sprintf("[Error creating NewBidOpenOrder] [%s] - %f, %f, %s, buy, %s", bb.ID, bb.Volume, bb.PriceVET, app.ac.ASymbol, err)
// 		e.DB.SaveLoggerErr(err)
// 		return
// 	}

// 	e.ac.AskOpenOrder(a)
// 	e.ac.BidOpenOrder(b)
// }

// setupNotifyService send a msg to ntfy mobile app
func (e *Engine) setupNotifyService() {
	// TEMPORARILY DISABLED - REMOVE THE TWO COMMENTS TO ENABLE

	// app.Notify = ntfy.NewNtfy()
	// app.Notify.SendMsg("Server is up - NTFY Test", "Server is up - LoggerErr Test", false)
}
