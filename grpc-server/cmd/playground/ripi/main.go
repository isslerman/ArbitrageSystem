package main

import (
	"fmt"
	"grpc-server/internal/cex"
)

func main() {
	// Example usage
	ripi := cex.NewRIPI()

	// the body to pass
	cor := cex.CreateOrderRequest{
		Amount: 0.5,
		Pair:   "SOL_BRL",
		Price:  779.99,
		Side:   "sell",
		Type:   "limit",
	}

	// pass a copy, not a pointer
	id, err := ripi.CreateOrder(cor)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}
	fmt.Println("Order created: ", id)
}
