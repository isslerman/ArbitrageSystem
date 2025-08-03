package marketdata

import (
	"context"
	"fmt"
	"looper/internal/pb/orders"
	"looper/internal/repository"
	"looper/pkg/data"
	"sync"
	"time"

	"go.uber.org/zap"
)

type MarketData struct {
	// backwards compatibility
	orders.UnimplementedOrderServiceServer
	Repo    repository.Repository
	Models  *data.Models
	Book    *BookTrade
	mu      sync.RWMutex
	Logger  *zap.SugaredLogger
	updates chan *BookTrade
}

func NewMarketData(repo repository.Repository, logger *zap.SugaredLogger) *MarketData {
	return &MarketData{
		Repo:    repo,
		Logger:  logger,
		updates: make(chan *BookTrade, 100),
	}
}

func (m *MarketData) Updates() <-chan *BookTrade {
	return m.updates
}

func (m *MarketData) WriteOrder(ctx context.Context, req *orders.OrderRequest) (*orders.OrderResponse, error) {
	inputAsk := req.GetOrderAsk()
	inputBid := req.GetOrderBid()

	book := &BookTrade{
		Timestamp:  inputAsk.CreatedAt, // or max(inputAsk.CreatedAt, inputBid.CreatedAt)
		LastPrice:  inputAsk.Price,     // optional logic here
		LastSide:   "ask",              // or "bid", or ""
		LastVolume: inputAsk.Volume,

		BestAsk: inputAsk.Price,
		BestBid: inputBid.Price,
	}

	m.mu.Lock()
	m.Book = book
	m.mu.Unlock()

	m.Repo.LogInfo(fmt.Sprintf("Book updated: Ask %.2f | Bid %.2f",
		book.BestAsk,
		book.BestBid,
	))

	now := time.Now().UnixMilli()
	lagMs := now - inputAsk.CreatedAt

	m.Logger.Debugw("MarketData received:",
		"ask_price", book.BestAsk,
		"bid_price", book.BestBid,
		"lag_ms", lagMs,
	)

	select {
	case m.updates <- book:
		// ok
	default:
		m.Logger.Warn("marketdata update channel is full")
	}

	return &orders.OrderResponse{Result: "success"}, nil
}
