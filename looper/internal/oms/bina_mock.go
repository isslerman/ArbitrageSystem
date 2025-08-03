package oms

import (
	"github.com/adshao/go-binance/v2"
)

// EXAMPLE OF BINA MOCK - NOT USING YET

type MockBinanceClient struct{}

func (m *MockBinanceClient) NewCreateOrderService() *binance.CreateOrderService {
	return &binance.CreateOrderService{} // stub with test response manually if needed
}

func (m *MockBinanceClient) NewListOpenOrdersService() *binance.ListOpenOrdersService {
	return &binance.ListOpenOrdersService{}
}

func (m *MockBinanceClient) NewCancelOrderService() *binance.CancelOrderService {
	return &binance.CancelOrderService{}
}

// You can mock these methods further using testify/mock or gomock if needed
