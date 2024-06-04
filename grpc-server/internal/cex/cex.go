package cex

import (
	"grpc-server/internal/cex/bina"
	"grpc-server/internal/cex/data"
	"grpc-server/internal/cex/ripi"
)

// ID represents different exchanges
type ID int

const (
	InstanceBina ID = iota
	InstanceBity
	InstanceRipi
	InstanceFoxb
)

// Cex is the interface that all exchanges must implement
type Cex interface {
	Balance(asset string) (amount float64, err error)
	CancelAllOrders() error
	Id() string
	CreateOrder(o *data.OrdersCreateRequest) (string, error)
}

// New is the factory method to create instances of different exchanges
func New(id ID) Cex {
	factoryFunctions := map[ID]func() Cex{
		InstanceBina: func() Cex { return bina.New() },
		// IDBity: func() Cex { return &Binance{} },
		InstanceRipi: func() Cex { return ripi.New() },
		// IDFoxb: func() Cex { return &Binance{} },
	}

	if factoryFunc, ok := factoryFunctions[id]; ok {
		return factoryFunc()
	}
	return nil
}
