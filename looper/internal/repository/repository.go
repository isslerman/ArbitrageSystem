package repository

import "looper/pkg/data"

type Repository interface {
	StoreOrderHistory(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error)
	LogInfo(log string)
	LogError(log string)
}
