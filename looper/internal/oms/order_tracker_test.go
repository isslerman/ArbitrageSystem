package oms

import (
	"testing"
)

func TestOrderTracker_SetAndGet(t *testing.T) {
	tracker := NewOrderTracker()

	tracker.Set("ask", "test-id", 5.20, 100.0)

	if !tracker.IsOpen("ask") {
		t.Error("Expected ask order to be open")
	}

	price, ok := tracker.GetPrice("ask")
	if !ok || price != 5.20 {
		t.Errorf("Expected price 5.20, got %.2f", price)
	}

	id, ok := tracker.GetID("ask")
	if !ok || id != "test-id" {
		t.Errorf("Expected id 'test-id', got %s", id)
	}

	ord, ok := tracker.GetOrder("ask")
	if !ok || ord.Amount != 100.0 || ord.Price != 5.20 || ord.ID != "test-id" {
		t.Error("OrderState not matching expected values")
	}
}

func TestOrderTracker_AutoID(t *testing.T) {
	tracker := NewOrderTracker()
	tracker.Set("bid", "", 4.98, 50.0)

	id, ok := tracker.GetID("bid")
	if !ok || id == "" {
		t.Error("Expected auto-generated ID for bid")
	}

	if !tracker.IsOpen("bid") {
		t.Error("Expected bid order to be open")
	}
}

func TestOrderTracker_Remove(t *testing.T) {
	tracker := NewOrderTracker()
	tracker.Set("ask", "removable-id", 5.10, 70.0)

	if !tracker.IsOpen("ask") {
		t.Fatal("Expected ask to be open before removal")
	}

	tracker.RemoveByID("removable-id")

	if tracker.IsOpen("ask") {
		t.Error("Expected ask to be removed")
	}
}
