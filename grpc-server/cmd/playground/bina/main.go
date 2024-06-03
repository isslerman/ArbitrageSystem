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
	err := bina.OrdersCreate(order)
	if err != nil {
		slog.Error("Error:", err)
	}

	// List Open orders
	err = bina.CancelAllOrders()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Print Exchange Name
	fmt.Println("Exchange ID:", bina.Id())
}
