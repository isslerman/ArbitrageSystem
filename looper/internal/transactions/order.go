package transactions

import (
	"fmt"
)

// this order implements transaction orders to be sent to cex
type order struct {
}

type CreateOrderService struct {
	amount float64
	pair   string // BTC_BRL
	price  float64
	side   string // buy || sell
	Type   string // market || limit
}

type IOrderTransaction interface {
}

func NewCreateOrderService() *CreateOrderService {
	fmt.Println("OrderService created.")
	return &CreateOrderService{}
}

func (os *CreateOrderService) createOrder() (err error) {
	return nil
}

func (os *CreateOrderService) Do() (err error) {
	return nil
}

func (os *CreateOrderService) Test() (err error) {
	return nil
}

type CreateOrderResponse struct {
	data       struct{}
	error_code int64
	message    string
}

// dados recebidos para gerarmos uma ordem
// ordem criada com dados que vieram do sistema de arbitragem
// ordem enviada a cex
// ordem executada
// envio de aviso de ordem executada
// ordem não executada
// envio de aviso de ordem não executada
// tentamos executar novamente?
// ordem executada // não executada
