package oms

import "looper/pkg/data"

type OMS interface {
	CreateOrder(o *data.OrdersCreateRequest) (string, error)
	CancelOrder(id string) error
	CancelAllOrders() error
	UpdateOrder(o *data.OrdersUpdateRequest) error
	IsAskOpen() bool
	IsBidOpen() bool
	GetAskPrice() (float64, bool)
	GetBidPrice() (float64, bool)
}
