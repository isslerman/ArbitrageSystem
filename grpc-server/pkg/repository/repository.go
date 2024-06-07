package repository

import "grpc-server/pkg/data"

type DatabaseRepo interface {
	SaveOrderHistory(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error)
	SaveLoggerErr(log string)
	SaveLoggerInfo(log string)
}
