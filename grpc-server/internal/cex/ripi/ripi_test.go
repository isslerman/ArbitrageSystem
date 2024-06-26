package ripi

import (
	"grpc-server/pkg/data"
	"testing"
)

func Test_CancelAllOrders(t *testing.T) {
	// exchange instance
	ripi := New()
	err := ripi.CancelAllOrders()

	// no orders to cancel
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOL_BRL",
		Price:  9999.00,
		Side:   "sell",
		Type:   "limit",
	}

	// creating an order
	_, err = ripi.CreateOrder(order)
	if err != nil {
		t.Errorf("error creating an order: %d", err)
	}

	// cancel all orders
	err = ripi.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}

}

// to force withou cache
// go test -count=1 ./...
func Test_CreateOrder(t *testing.T) {
	// New instance of the exchange RIPI
	ripi := New()

	order := &data.OrdersCreateRequest{
		Amount: 0.10,
		Pair:   "SOL_BRL",
		Price:  9999.00,
		Side:   "sell",
		Type:   "limit",
	}

	// creating an order
	_, err := ripi.CreateOrder(order)
	if err != nil {
		t.Errorf("error creating an order: %d", err)
	}

	// cancel all orders
	err = ripi.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %d", err)
	}
}
