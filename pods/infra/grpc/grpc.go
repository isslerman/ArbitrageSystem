// GRPC package -
// Layer: App Layer - send the data thru grpc
// Data type: OrderRequest
package grpc

import (
	"context"
	"fmt"
	"grpc-client/internal/pb/orders"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Just a gRPC client that connects with our server and send a msg.
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func SendViaGRPC(orderAsk, orderBid *orders.Order) {
	app := Config{
		Host: "localhost",
		Port: 50001,
	}
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
