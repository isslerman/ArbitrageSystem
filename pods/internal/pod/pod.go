// package pod is responsable for receive a connect within an exchage via RESP API and get the orderbook price for a token
// it will clean, transform the orderbook and return the best bid order and the best ask order with a minimum of setup
// configured
// Layer: Business Layer. But the RESP API data must come from the App layer. How to do it?
// Input: Exchange data orderbook
// Output: Formated data to be sent by gRPC to the server.
package pod

import (
	"fmt"
	"grpc-client/pkg/exchange"
)

type Pod struct {
	exchange exchange.IExchange
}

func NewPod(e exchange.IExchange) *Pod {
	return &Pod{
		exchange: e,
	}
}

func (pod *Pod) fetchBestOrder() {
	fmt.Println("Pod.fetchBestOrder()")
	// fetch the data from the exchange, input data basic types
	// insert into order data bus layer
}

func (pod *Pod) sendBestOrderViaGRPC() {
	fmt.Println("Pod.sendBestOrderViaGRPC")
}

// type inputData struct {
// }

// type outputData struct {
// }

// value object - bestOrder - muda a cada pesquisa e fetch
// entidade - não muda.
