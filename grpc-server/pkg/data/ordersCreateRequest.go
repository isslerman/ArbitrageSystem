package data

// CreateOrderRequest represents the structure of the JSON request.
type OrdersCreateRequest struct {
	Amount float64 `json:"amount"`
	Pair   string  `json:"pair"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Type   string  `json:"type"`
}

// orderLimits is the validation for an order
type OrderLimits struct {
	OrdMinAmount float64
}
