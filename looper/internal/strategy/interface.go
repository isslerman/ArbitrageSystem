package strategy

import (
	"context"
	"looper/internal/marketdata"
	"looper/internal/oms"
)

type Strategy interface {
	Run(ctx context.Context, marketData <-chan *marketdata.BookTrade, orderEvents <-chan *oms.ExecutionReport) error
	StopAllOrders() error
	GetOrderStream() <-chan *oms.ExecutionReport
}
