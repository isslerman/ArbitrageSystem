package bina

import (
	"grpc-server/internal/cex/data"
	"testing"
)

func Test_CancelAllOrders(t *testing.T) {
	// exchange instance
	ripi := New()
	err := ripi.CancelAllOrders()

	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}
}

func Test_OrdersCreate(t *testing.T) {
	// New instance of the exchange RIPI
	ripi := New()

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOLBRL",
		Price:  700.00,
		Side:   "buy",
		Type:   "limit",
	}

	// creating an order
	err := ripi.OrdersCreate(order)
	if err != nil {
		t.Errorf("error creating an order: %d", err)
	}

	// cancel all orders
	err = ripi.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}
}
