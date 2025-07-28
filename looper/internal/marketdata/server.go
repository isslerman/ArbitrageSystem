package marketdata

import (
	"context"
	"fmt"
	"log/slog"
	"looper/internal/pb/orders"
	"looper/pkg/data"
	"net"

	"google.golang.org/grpc"
)

type MarketDataReceiver struct {
	// backwards compatibility
	orders.UnimplementedOrderServiceServer
	Models *data.Models
}

func NewOrderServer() *MarketDataReceiver {
	return &MarketDataReceiver{
		Models: data.NewModels(),
	}
}

// using the code generated - receving the request and the response
func (o *MarketDataReceiver) WriteOrder(ctx context.Context, req *orders.OrderRequest) (*orders.OrderResponse, error) {
	// get the input req
	inputAsk := req.GetOrderAsk()
	inputBid := req.GetOrderBid()

	// create a new orderEntry
	askEntry := &data.Order{
		ID:        inputAsk.Id,
		Price:     inputAsk.Price,
		PriceVET:  inputAsk.PriceVET,
		Volume:    inputAsk.Volume,
		CreatedAt: inputAsk.CreatedAt,
	}

	bidEntry := &data.Order{
		ID:        inputBid.Id,
		Price:     inputBid.Price,
		PriceVET:  inputBid.PriceVET,
		Volume:    inputBid.Volume,
		CreatedAt: inputBid.CreatedAt,
	}

	// add order to orderbook
	o.Models.OrderBook.AddOrUpdateAskOrder(askEntry)
	o.Models.OrderBook.AddOrUpdateBidOrder(bidEntry)
	fmt.Printf("order received:[%d] ID[%s], Price:%f, Volume:%f\n", askEntry.CreatedAt, askEntry.ID, askEntry.Price, askEntry.Volume)

	// return the res
	res := &orders.OrderResponse{
		Result: "sucess",
	}
	return res, nil
}

func GRPCListen(os *MarketDataReceiver, port string) error {
	// grpc server listeing on a tcp port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	s := grpc.NewServer()
	orders.RegisterOrderServiceServer(s, os)

	slog.Info("Market Data gRPC server started", "port", port)
	return s.Serve(lis)
}
