package exchange

import "grpc-client/internal/data"

type IExchange interface {
	BestOrder() (*data.BestOrder, error)
}
