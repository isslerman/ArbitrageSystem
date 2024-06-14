package bina

import (
	"grpc-server/pkg/data"
	"testing"
)

func Test_CancelAllOrders(t *testing.T) {
	// exchange instance
	bina := New()
	err := bina.CancelAllOrders()

	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}
}

func Test_OrdersCreate(t *testing.T) {
	// New instance of the exchange RIPI
	bina := New()

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOLBRL",
		Price:  700.0,
		Side:   "buy",
		Type:   "limit",
	}

	// creating an order
	_, err := bina.CreateOrder(order)
	if err != nil {
		t.Errorf("error creating an order: %d", err)
	}

	// cancel all orders
	err = bina.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}
}
