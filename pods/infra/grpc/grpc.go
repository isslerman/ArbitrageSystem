// GRPC package -
// Layer: App Layer - send the data thru grpc
// Data type: OrderRequest
package grpc

import (
	"context"
	"fmt"
	"log"
	"pods/internal/pb/orders"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Just a gRPC client that connects with our server and send a msg.
type Config struct {
	host string
	port int
	Conn *grpc.ClientConn
	l    *zap.Logger
}

func NewServiceGRPC(l *zap.Logger) *Config {
	app := Config{
		host: "localhost",
		port: 50001,
	}

	target := fmt.Sprintf("%s:%d", app.host, app.port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("Error Conn")
	}
	// defer conn.Close()

	return &Config{
		host: app.host,
		port: app.port,
		Conn: conn,
	}
}

func (grpc *Config) SendViaGRPC(orderAsk, orderBid *orders.Order) {
	c := orders.NewOrderServiceClient(grpc.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.WriteOrder(ctx, &orders.OrderRequest{
		OrderAsk: orderAsk,
		OrderBid: orderBid,
	})

	if err != nil {
		log.Fatal("Error WriteOrder: ", err)
	}
}
