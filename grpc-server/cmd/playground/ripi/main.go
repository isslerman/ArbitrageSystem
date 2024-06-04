package main

import (
	"fmt"
	"grpc-server/internal/cex"
	"grpc-server/internal/cex/data"
)

func main() {

	// New instance of the exchange RIPI
	ripi := cex.New(cex.InstanceRipi)

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOL_BRL",
		Price:  999.00,
		Side:   "sell",
		Type:   "limit",
	}

	// creating an order
	id, err := ripi.CreateOrder(order)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("orders cancelled")
	}

	// Print the exchange ID and the order ID
	fmt.Printf("Exchange ID: %s OrderID: %s", ripi.Id(), id)

	// cancel all orders
	err = ripi.CancelAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("orders cancelled")
	}
}
