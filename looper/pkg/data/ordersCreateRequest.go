package data

// CreateOrderRequest represents the structure of the JSON request.
type OrdersCreateRequest struct {
	Amount float64 `json:"amount"`
	Pair   string  `json:"pair"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Type   string  `json:"type"`
}

type OrdersUpdateRequest struct {
	Amount  float64 `json:"amount"`
	Price   float64 `json:"uni_price"`
	OrderID string  `json:"order_id"`
}

// orderLimits is the validation for an order
type OrderLimits struct {
	OrdMinAmount float64
}

// type CancelOrderRequest struct {
// 	ID string `json:"id"`
// }
