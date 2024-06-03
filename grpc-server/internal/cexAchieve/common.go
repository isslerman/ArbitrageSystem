package cex

import (
	"io"
	"net/http"
	"net/url"
)

type limits struct {
	OrdMinAmount float64
}

// CreateOrderRequest represents the structure of the JSON request.
type CreateOrderRequest struct {
	Amount float64 `json:"amount"`
	Pair   string  `json:"pair"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Type   string  `json:"type"`
}

// request define an API request
type request struct {
	method     string
	endpoint   string
	query      url.Values
	form       url.Values
	recvWindow int64
	// secType    secType
	header  http.Header
	body    io.Reader
	fullURL string
}
