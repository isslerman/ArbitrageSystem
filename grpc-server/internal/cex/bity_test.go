package cex

import (
	"testing"
)

func TestBITY_CancelAllOrders(t *testing.T) {
	// Example usage
	bity := NewBITY()

	// pass a copy, not a pointer
	suc, err := bity.CancelAllOrders()
	if err != nil {
		t.Errorf("cannot cancel all orders: %v", err)
	}
	if !suc {
		t.Errorf("cannot cancel all orders")
	}
}
