package oms

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type OrderState struct {
	ID     string
	Price  float64
	Amount float64
	Side   string
	Time   time.Time
}

type OrderTracker struct {
	mu    sync.RWMutex
	state map[string]OrderState // "ask" or "bid"
}

func NewOrderTracker() *OrderTracker {
	return &OrderTracker{
		state: make(map[string]OrderState),
	}
}

func (ot *OrderTracker) Set(side, id string, price, amount float64) {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if id == "" {
		id = fmt.Sprintf("%s%d", strings.ToLower(side), time.Now().UnixNano())
	}

	ot.state[strings.ToLower(side)] = OrderState{
		ID:     id,
		Price:  price,
		Amount: amount,
		Side:   strings.ToLower(side),
		Time:   time.Now(),
	}
}

func (ot *OrderTracker) RemoveByID(id string) {
	ot.mu.Lock()
	defer ot.mu.Unlock()
	for side, ord := range ot.state {
		if ord.ID == id {
			delete(ot.state, side)
		}
	}
}

func (ot *OrderTracker) IsOpen(side string) bool {
	ot.mu.RLock()
	defer ot.mu.RUnlock()
	_, ok := ot.state[strings.ToLower(side)]
	return ok
}

func (ot *OrderTracker) GetPrice(side string) (float64, bool) {
	ot.mu.RLock()
	defer ot.mu.RUnlock()
	ord, ok := ot.state[strings.ToLower(side)]
	return ord.Price, ok
}

func (ot *OrderTracker) GetID(side string) (string, bool) {
	ot.mu.RLock()
	defer ot.mu.RUnlock()
	ord, ok := ot.state[strings.ToLower(side)]
	return ord.ID, ok
}

func (ot *OrderTracker) GetOrder(side string) (OrderState, bool) {
	ot.mu.RLock()
	defer ot.mu.RUnlock()
	ord, ok := ot.state[strings.ToLower(side)]
	return ord, ok
}

func (ot *OrderTracker) GetOpenPrice(side string) (float64, bool) {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	order, ok := ot.state[strings.ToLower(side)]
	if !ok {
		return 0.0, false
	}
	return order.Price, true
}
