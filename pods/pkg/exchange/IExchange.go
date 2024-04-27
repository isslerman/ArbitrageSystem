package exchange

import "pods/internal/data"

type IExchange interface {
	BestOrder() (*data.Ask, *data.Bid, error)
	ExchangeID() string
}
