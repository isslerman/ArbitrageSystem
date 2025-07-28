package trading

import (
	"testing"
)

func TestOrdersManager_CreateOrder(t *testing.T) {
	om := &TradingManager{
		openOrders: make(map[string]Order),
	}

	if err := om.CreateOrder("ORDER-1"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, exists := om.openOrders["ORDER-1"]; !exists {
		t.Errorf("expected order to exist in map")
	}
}

func TestOrdersManager_ListOrders(t *testing.T) {
	om := &TradingManager{
		openOrders: make(map[string]Order),
	}

	err := om.CreateOrder("ORDER-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	orderList := om.ListOrders()
	if len(orderList) != 1 {
		t.Errorf("expected list to have 1 order, got %d", len(orderList))
	}
	for _, order := range orderList {
		if order.ID != "ORDER-1" {
			t.Errorf("expected order ID to be 'ORDER-1', got '%s'", order.ID)
		}
	}
}

func TestOrdersManager_RemoveOrder(t *testing.T) {
	om := &TradingManager{
		openOrders: make(map[string]Order),
	}

	err := om.CreateOrder("ORDER-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = om.RemoveOrder("ORDER-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, exists := om.openOrders["ORDER-1"]; exists {
		t.Errorf("expected order to be removed from map")
	}
}

func TestOrdersManager_RemoveNonExistentOrder(t *testing.T) {
	om := &TradingManager{
		openOrders: make(map[string]Order),
	}

	err := om.RemoveOrder("ORDER-1")
	if err == nil {
		t.Errorf("expected error, got no error")
	}
}

func TestOrdersManager_ListAfterRemoval(t *testing.T) {
	om := &TradingManager{
		openOrders: make(map[string]Order),
	}

	err := om.CreateOrder("ORDER-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = om.RemoveOrder("ORDER-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	orderList := om.ListOrders()
	if len(orderList) != 0 {
		t.Errorf("expected list to be empty, got %d openOrders", len(orderList))
	}
}
