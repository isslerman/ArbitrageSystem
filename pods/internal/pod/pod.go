// package pod is responsable for receive a connect within an exchage via RESP API and get the orderbook price for a token
// it will clean, transform the orderbook and return the best bid order and the best ask order with a minimum of setup
// configured
// Layer: Business Layer. But the RESP API data must come from the App layer. How to do it?
// Input: Exchange data orderbook
// Output: Formated data to be sent by gRPC to the server.
package pod

import (
	"fmt"
	"pods/infra/grpc"
	"pods/internal/data"
	"pods/internal/pb/orders"
	"pods/pkg/exchange"
)

// type OrderMsg *orders.Order

type Pod struct {
	exchange  exchange.IExchange
	orderbook *data.Orderbook
	// msgService *infra.IMsgService
	conn *grpc.Config // create grpc instance

}

func NewPod(e exchange.IExchange) *Pod {
	return &Pod{
		exchange: e,
		conn:     grpc.NewServiceGRPC(),
	}
}

// Run() fetch the data from the Iexchange one time
func (pod *Pod) Run() error {
	// check if exchange if nil
	// fetchExchange
	// here we need to have all data, before filter
	ask, bid := pod.fetchBestOrder()
	// validate - pod is the man with the bus rules
	// has volume?
	// has the min that attend?
	// here we now that the data comming from input is validated
	// we are in the bus/first class
	pod.orderbook = data.NewOrderBook(ask, bid)

	// orderbook.validate.rules
	// if something... only 3 orders
	// now we are just passing withour rules

	// send msg over gRPC
	pod.sendBestOrderViaGRPC()

	return nil
}

func (pod *Pod) fetchBestOrder() (*data.Ask, *data.Bid) {
	// fmt.Println("Pod.fetchBestOrder()")
	// we are using the interface, the exchange was choosen in the main.go
	// fetch the data from the exchange
	ask, bid, _ := pod.exchange.BestOrder()
	// insert into order data bus layer
	return ask, bid
}

func (pod *Pod) sendBestOrderViaGRPC() {
	// fmt.Println("Pod.sendBestOrderViaGRPC")

	oa := &orders.Order{
		Id:        pod.exchange.ExchangeID(),
		Price:     pod.orderbook.Ask.Price,
		PriceVET:  pod.orderbook.Ask.PriceVET,
		CreatedAt: pod.orderbook.Ask.CreatedAt.Unix(),
	}

	ob := &orders.Order{
		Id:        pod.exchange.ExchangeID(),
		Price:     pod.orderbook.Bid.Price,
		PriceVET:  pod.orderbook.Bid.PriceVET,
		CreatedAt: pod.orderbook.Bid.CreatedAt.Unix(),
	}
	// send msg over grpc
	fmt.Printf("send: %v, %v\n", oa, ob)
	pod.conn.SendViaGRPC(oa, ob)
}

// type inputData struct {
// }

// type outputData struct {
// }

// value object - bestOrder - muda a cada pesquisa e fetch
// entidade - não muda.
