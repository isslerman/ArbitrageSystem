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
	err := ripi.OrdersCreate(order)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("orders cancelled")
	}

	// Print Exchange Name
	fmt.Println("Exchange ID:", ripi.Id())

	// cancel all orders
	err = ripi.CancelAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("orders cancelled")
	}

	// Check if there is any order open

	// Create an order

	// Cancel all orders

	// the body to pass
	// cor := cex.CreateOrderRequest{
	// 	Amount: 0.5,
	// 	Pair:   "SOL_BRL",
	// 	Price:  779.99,
	// 	Side:   "sell",
	// 	Type:   "limit",
	// }

	// pass a copy, not a pointer
	// id, err := ripi.CreateOrder(cor)
	// if err != nil {
	// 	fmt.Println("Error creating order:", err)
	// 	return
	// }
	// fmt.Println("Order created: ", id)
}
