package data

type Models struct {
	Order     Order
	OrderBook *OrderBook
}

func NewModels() *Models {
	return &Models{
		Order:     Order{},
		OrderBook: &OrderBook{},
	}
}
