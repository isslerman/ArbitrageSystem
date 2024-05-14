package main

import (
	"fmt"
	"grpc-server/internal/cex"
)

func main() {
	// Example usage
	bity := cex.NewBITY()

	// the body to pass
	// cor := cex.CreateOrderRequest{
	// 	Amount: 0.5,
	// 	Pair:   "SOL_BRL",
	// 	Price:  779.99,
	// 	Side:   "sell",
	// 	Type:   "limit",
	// }

	suc, err := bity.CancelAllOrders()
	if err != nil {
		fmt.Printf("cannot cancel all orders: %v", err)
	}

	fmt.Println(suc)
}
