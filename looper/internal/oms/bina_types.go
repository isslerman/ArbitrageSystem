package oms

import (
	"encoding/json"
	"fmt"
)

// ExecutionReport struct define exchange websocket transactions
// https://developers.binance.com/docs/binance-spot-api-docs/user-data-stream#conditional-fields-in-execution-report
type ExecutionReport struct {
	EventType               string  `json:"e"` // Event type
	EventTime               int64   `json:"E"` // Event time (epoch ms)
	Symbol                  string  `json:"s"` // Symbol
	ClientOrderID           string  `json:"c"` // Client order ID
	Side                    string  `json:"S"` // Side (BUY/SELL)
	OrderType               string  `json:"o"` // Order type (LIMIT/MARKET/etc.)
	TimeInForce             string  `json:"f"` // Time in force (GTC/IOC/etc.)
	Quantity                string  `json:"q"` // Order quantity
	Price                   string  `json:"p"` // Order price
	StopPrice               string  `json:"P"` // Stop price
	IcebergQuantity         string  `json:"F"` // Iceberg quantity
	OrderListID             int64   `json:"g"` // OrderListId
	OrigClientOrderID       string  `json:"C"` // Original client order ID (if cancel)
	ExecutionType           string  `json:"x"` // Execution type (NEW/TRADE/CANCELED/etc.)
	OrderStatus             string  `json:"X"` // Current order status
	OrderRejectReason       string  `json:"r"` // Order reject reason
	OrderID                 int64   `json:"i"` // Order ID
	LastExecutedQuantity    string  `json:"l"` // Last executed quantity
	CumulativeFilledQty     string  `json:"z"` // Cumulative filled quantity
	LastExecutedPrice       string  `json:"L"` // Last executed price
	CommissionAmount        string  `json:"n"` // Commission amount
	CommissionAsset         *string `json:"N"` // Commission asset (can be null)
	TransactionTime         int64   `json:"T"` // Transaction time
	TradeID                 int64   `json:"t"` // Trade ID
	PreventedMatchID        *int64  `json:"v"` // Prevented Match ID (only if expired due to STP)
	ExecutionID             int64   `json:"I"` // Execution ID
	IsOrderOnBook           bool    `json:"w"` // Is order on the book?
	IsMakerSide             bool    `json:"m"` // Is trade the maker side?
	IgnoreField             bool    `json:"M"` // Ignore
	OrderCreationTime       int64   `json:"O"` // Order creation time
	CumulativeQuoteQty      string  `json:"Z"` // Cumulative quote asset transacted
	LastQuoteQty            string  `json:"Y"` // Last quote asset transacted
	QuoteOrderQty           string  `json:"Q"` // Quote order quantity
	WorkingTime             int64   `json:"W"` // Working time
	SelfTradePreventionMode string  `json:"V"` // STP mode (e.g. NONE)
}

func parseExecReport(msg []byte) (*ExecutionReport, bool, error) {
	// Step 1: Unmarshal the top-level message and extract "event"
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(msg, &envelope); err != nil {
		return nil, false, fmt.Errorf("failed to decode envelope: %w", err)
	}

	eventRaw, ok := envelope["event"]
	if !ok {
		return nil, false, fmt.Errorf(`missing "event" key`)
	}

	// Step 2: Extract event type "e" safely to check if it's an executionReport
	var inner map[string]json.RawMessage
	if err := json.Unmarshal(eventRaw, &inner); err != nil {
		return nil, false, fmt.Errorf("failed to decode inner event JSON: %w", err)
	}

	var eventType string
	if err := json.Unmarshal(inner["e"], &eventType); err != nil {
		return nil, false, fmt.Errorf("failed to decode event type: %w", err)
	}
	if eventType != "executionReport" {
		return nil, false, nil // not an executionReport, ignore it
	}

	// Step 3: Parse the full event as an ExecutionReport
	var er ExecutionReport
	if err := json.Unmarshal(eventRaw, &er); err != nil {
		return nil, false, fmt.Errorf("failed to decode executionReport: %w", err)
	}

	return &er, true, nil
}
