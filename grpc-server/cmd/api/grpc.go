package main

import (
	"context"
	"fmt"
	"grpc-server/data"
	"grpc-server/internal/pb/orders"
	"log"
	"net"

	"google.golang.org/grpc"
)

type OrderServer struct {
	// backwards compatibility
	orders.UnimplementedOrderServiceServer
	Models *data.Models
}

func NewOrderServer() *OrderServer {
	return &OrderServer{
		Models: data.NewModels(),
	}
}

// using the code generated - receving the request and the response
func (o *OrderServer) WriteOrder(ctx context.Context, req *orders.OrderRequest) (*orders.OrderResponse, error) {
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
	// fmt.Printf("order received:[%d] ID[%s], Price:%f, Volume:%f\n", orderEntry.CreatedAt, orderEntry.ID, orderEntry.Price, orderEntry.Volume)

	// return the res
	res := &orders.OrderResponse{
		Result: "sucess",
	}
	return res, nil
}

func (app *Config) gRPCListen(os *OrderServer) {
	// grpc server listeing on a tcp port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer(
	// grpc.MaxConcurrentStreams(20000),  // Limit the number of concurrent streams
	// grpc.MaxRecvMsgSize(1024*1024*10), // Set the maximum message size
	// grpc.MaxSendMsgSize(1024*1024*10), // Set the maximum message size
	// grpc.ReadBufferSize(1024*1024*60),
	// grpc.NumStreamWorkers(1),
	)

	orders.RegisterOrderServiceServer(s, os)

	// log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
