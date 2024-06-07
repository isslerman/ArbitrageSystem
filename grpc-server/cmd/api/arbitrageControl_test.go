package main

import (
	"grpc-server/internal/cex"
	"grpc-server/pkg/data"
	"testing"
)

func Test_NewArbitrageControl(t *testing.T) {
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina

	aSymbol := "SOLBRL"
	bSymbol := "SOL_BRL"

	// creating a new AC
	ac, err := NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, nil)
	if err != nil {
		t.Errorf("error creating ArbitrageControl: %d", err)
	}
	if ac.AskSymbol != "SOLBRL" {
		t.Errorf("error creating ArbitrageControl. expected SOL_BRL, got: %s", ac.AskSymbol)
	}
}

func Test_hasAskOpenOrders(t *testing.T) {
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina

	aSymbol := "SOL_BRL"
	bSymbol := "SOLBRL"

	// creating a new AC
	ac, err := NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, nil)
	if err != nil {
		t.Errorf("error creating ArbitrageControl: %d", err)
	}

	// create a new askopenorder
	ao, err := data.NewAskOpenOrder(0.1, 9999.00, aSymbol, "limit")
	if err != nil {
		t.Errorf("error creating newaskopenorder, %s", err)
	}

	_, err = ac.createLimitOrder(data.OpenOrder(ao), "ask")
	if err != nil {
		t.Errorf("error creating limit order, %s", err)
	}
	// set a new askopenorder to the control
	// ac.AskOpenOrder(ao)
	// test
	if ac.hasAskOpenOrders() != true {
		t.Errorf("error Test_hasAskOpenOrders. Expected false, got %t", ac.hasAskOpenOrders())
	}
	ac.cancelAllAskOrders()
}

func Test_createLimitOrder(t *testing.T) {
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina

	aSymbol := "SOL_BRL"
	bSymbol := "SOLBRL"

	ao, err := data.NewAskOpenOrder(0.1, 9999.00, aSymbol, "limit")
	if err != nil {
		t.Errorf("error creating newaskopenorder, %s", err)
	}

	// creating a new AC
	ac, err := NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, nil)
	if err != nil {
		t.Errorf("error creating ArbitrageControl: %d", err)
	}

	// create a limit order on the ask side
	_, err = ac.createLimitOrder(data.OpenOrder(ao), "ask")
	if err != nil {
		t.Errorf("error creating ask order: %d", err)
	}

	// check is there is one order open
	if ac.hasAskOpenOrders() != true {
		t.Errorf("error creating ask order. Expected to hasAskOpenOrders true but got %t", ac.hasAskOpenOrders())
	}

	// cancell all orders
	err = ac.cancelAllAskOrders()
	if err != nil {
		t.Errorf("error Test_createLimitOrder, cancelAllAskOrders: %d", err)
	}

	// check if all orders are cancelled
	if ac.hasAskOpenOrders() != false {
		t.Errorf("error Test_createLimitOrder. Expected to hasAskOpenOrders false but got %t", ac.hasAskOpenOrders())
	}
	ac.cancelAllAskOrders()
}

func Test_TryToCreateTwoLimitOrder(t *testing.T) {
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina

	aSymbol := "SOL_BRL"
	bSymbol := "SOLBRL"

	ao, err := data.NewAskOpenOrder(0.1, 9999.00, aSymbol, "limit")
	if err != nil {
		t.Errorf("error creating newaskopenorder, %s", err)
	}

	// creating a new AC
	ac, err := NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, nil)
	if err != nil {
		t.Errorf("error creating ArbitrageControl: %d", err)
	}

	_, err = ac.createLimitOrder(data.OpenOrder(ao), "ask")
	if err != nil {
		t.Errorf("error creating ask order: %d", err)
	}
	_, err = ac.createLimitOrder(data.OpenOrder(ao), "ask")
	if err.Error() != "ask order already created and open" {
		t.Errorf("error creating second ask order: %d", err)
	}
	ac.cancelAllAskOrders()
}
