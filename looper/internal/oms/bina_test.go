package oms

import (
	"errors"
	"log"
	"looper/config"
	"looper/pkg/data"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	root, err := FindProjectRootWithEnv()
	if err != nil {
		log.Fatalf("Failed to find project root: %v", err)
	}

	envPath := filepath.Join(root, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Failed to load .env file at %s: %v", envPath, err)
	}

	cfg = loadCfg()

	os.Exit(m.Run())
}

func Test_IsAskOpen_IsBidOpen(t *testing.T) {
	tracker := NewOrderTracker()

	if tracker.IsOpen("ask") {
		t.Errorf("Expected ask to be false initially")
	}
	if tracker.IsOpen("bid") {
		t.Errorf("Expected bid to be false initially")
	}

	tracker.Set("ask", "ask-123", 5.5, 1.0)
	if !tracker.IsOpen("ask") {
		t.Errorf("Expected ask to be open after Set()")
	}
	if tracker.IsOpen("bid") {
		t.Errorf("Expected bid to still be false")
	}

	tracker.Set("bid", "bid-456", 5.4, 1.0)
	if !tracker.IsOpen("bid") {
		t.Errorf("Expected bid to be open after Set()")
	}

	tracker.RemoveByID("bid-123")
	if tracker.IsOpen("ask") {
		t.Errorf("Expected ask to be false after RemoveByID()")
	}
	if !tracker.IsOpen("bid") {
		t.Errorf("Expected bid to still be open")
	}
}

func Test_OrderTracker_GetValues(t *testing.T) {
	tracker := NewOrderTracker()
	tracker.Set("ask", "id-abc", 10.0, 0.1)

	price, ok := tracker.GetPrice("ask")
	if !ok || price != 10.0 {
		t.Errorf("Expected price 10.0, got %f", price)
	}

	id, ok := tracker.GetID("ask")
	if !ok || id != "id-abc" {
		t.Errorf("Expected id 'id-abc', got %s", id)
	}
}

func Test_CancelAllOrders(t *testing.T) {
	bina := NewBina(cfg)

	err := bina.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %v", err)
	}
}

func Test_OrdersCreate(t *testing.T) {
	bina := NewBina(cfg)

	order := &data.OrdersCreateRequest{
		Amount: 3,
		Pair:   "USDTBRL",
		Price:  5.0,
		Side:   "buy",
		Type:   "limit",
	}

	_, err := bina.CreateOrder(order)
	if err != nil {
		t.Errorf("error creating an order: %v", err)
	}

	err = bina.CancelAllOrders()
	if err != nil {
		t.Errorf("CancelAllOrders returned an error: %v", err)
	}
}

func loadCfg() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config.")
	}
	return cfg
}

func FindProjectRootWithEnv() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root without finding .env
		}
		dir = parent
	}

	return "", errors.New(".env not found in any parent directory")
}
