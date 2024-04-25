package service

import (
	"context"
	"fmt"
	"grpc-client/internal/pb/orders"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderPayload struct {
	ID     string
	Price  float64
	Volume float64
}

func (app *Config) sendViaGRPC(orderAsk, orderBid *orders.Order) {
	target := fmt.Sprintf("%s:%d", app.Host, app.Port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("Error Conn")
	}
	defer conn.Close()

	// client
	c := orders.NewOrderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.WriteOrder(ctx, &orders.OrderRequest{
		OrderAsk: orderAsk,
		OrderBid: orderBid,
	})
	if err != nil {
		log.Fatal("Error WriteOrder: ", err)
	}
}
