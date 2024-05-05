// package pod is responsable for receive a connect within an exchage via RESP API and get the orderbook price for a token
// it will clean, transform the orderbook and return the best bid order and the best ask order with a minimum of setup
// configured
// Layer: Business Layer. But the RESP API data must come from the App layer. How to do it?
// Input: Exchange data orderbook
// Output: Formated data to be sent by gRPC to the server.
package pod

import (
	"errors"
	"fmt"
	"pods/infra/grpc"
	"pods/internal/data"
	"pods/internal/pb/orders"
	"pods/pkg/exchange"

	"go.uber.org/zap"
)

// type OrderMsg *orders.Order

type Pod struct {
	exchange  exchange.IExchange
	orderbook *data.Orderbook
	// msgService *infra.IMsgService
	conn *grpc.Config // create grpc instance

}

func NewPod(e exchange.IExchange, l *zap.Logger) *Pod {
	return &Pod{
		exchange: e,
		conn:     grpc.NewServiceGRPC(l),
	}
}

var (
	ErrAskZeroVolume     = errors.New("ask volume zero")
	ErrBidZeroVolume     = errors.New("bid volume zero")
	ErrfetchBestOrderNil = errors.New("fetching best order returning nil")
)

// Run() fetch the data from the Iexchange one time
func (pod *Pod) Run() error {
	// check if exchange if nil
	// fetchExchange
	// here we need to have all data, before filter
	ask, bid, err := pod.fetchBestOrder()

	// validate - pod is the man with the bus rules
	// has volume?
	// has the min that attend?
	// here we now that the data comming from input is validated
	// insert this into validation
	if err != nil {
		fmt.Printf("Error: %v", ErrfetchBestOrderNil)
		return nil
	}

	if (ask == nil) || (bid == nil) {
		fmt.Printf("Error: %v", ErrfetchBestOrderNil)
		return nil
	}
	// we are in the bus/first class
	pod.orderbook = data.NewOrderBook(ask, bid)

	if pod.orderbook.Ask.Volume == 0 {
		return ErrAskZeroVolume
	}

	if pod.orderbook.Bid.Volume == 0 {
		return ErrBidZeroVolume
	}

	// orderbook.validate.rules
	// if something... only 3 orders
	// now we are just passing withour rules

	// send msg over gRPC
	pod.sendBestOrderViaGRPC()

	return nil
}

func (pod *Pod) fetchBestOrder() (*data.Ask, *data.Bid, error) {
	// fmt.Println("Pod.fetchBestOrder()")
	// we are using the interface, the exchange was choosen in the main.go
	// fetch the data from the exchange
	ask, bid, err := pod.exchange.BestOrder()
	if err != nil {
		return nil, nil, err
	}

	// insert into order data bus layer
	return ask, bid, nil
}

func (pod *Pod) sendBestOrderViaGRPC() {
	// fmt.Println("Pod.sendBestOrderViaGRPC")

	oa := &orders.Order{
		Id:        pod.exchange.ExchangeID(),
		Price:     pod.orderbook.Ask.Price,
		PriceVET:  pod.orderbook.Ask.PriceVET,
		Volume:    pod.orderbook.Ask.Volume,
		CreatedAt: pod.orderbook.Ask.CreatedAt.Unix(),
	}

	ob := &orders.Order{
		Id:        pod.exchange.ExchangeID(),
		Price:     pod.orderbook.Bid.Price,
		PriceVET:  pod.orderbook.Bid.PriceVET,
		Volume:    pod.orderbook.Bid.Volume,
		CreatedAt: pod.orderbook.Bid.CreatedAt.Unix(),
	}
	// send msg over grpc
	fmt.Printf("send: ask %v, %v\n", oa, ob)
	pod.conn.SendViaGRPC(oa, ob)
}

// type inputData struct {
// }

// type outputData struct {
// }

// value object - bestOrder - muda a cada pesquisa e fetch
// entidade - não muda.
