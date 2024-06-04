package main

import (
	"fmt"
	"grpc-server/internal/cex"
	"grpc-server/internal/cex/data"
	"log/slog"
)

func main() {

	// New instance of the exchange RIPI
	bina := cex.New(cex.InstanceBina)

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOLBRL",
		Price:  1200.00,
		Side:   "sell",
		Type:   "limit",
	}

	// creating an order
	id, err := bina.CreateOrder(order)
	if err != nil {
		slog.Error("Error:", err)
	}

	// List Open orders
	err = bina.CancelAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print the exchange ID and the order ID
	fmt.Printf("Exchange ID: %s OrderID: %s", bina.Id(), id)
}
