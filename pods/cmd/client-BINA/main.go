package main

import (
	"fmt"
	"grpc-client/internal/data"
	"grpc-client/internal/pb/orders"
	"grpc-client/pkg/exchanges"
	"log"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Just a gRPC client that connects with our server and send a msg.

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	GRPC *grpc.ClientConn
}

func main() {
	app := Config{
		Host: "localhost",
		Port: 50001,
	}

	grpcConn, err := app.connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	e := exchanges.NewBinance()

	// loop checking best order from exchange every 1/3 seconds
	for {
		order, err := e.BestOrder()
		if err != nil {
			slog.Warn("BINA: BestOrder error", err)
		} else {
			app.sendViaGRPC(order)
		}

		time.Sleep(time.Millisecond * 1000 / 3)
	}
}

func (app *Config) sendViaGRPC(order *data.BestOrder) {
	fmt.Printf("%s BestAsk: %f, Volume: %f\n", e.Id, order.BestAsk.Price, order.BestAsk.Volume)
	fmt.Printf("%s BestBid: %f, Volume: %f\n", e.Id, order.BestBid.Price, order.BestBid.Volume)

	orderAsk := &orders.Order{
		Id:        e.Id,
		Price:     order.BestAsk.Price,
		PriceVET:  order.BestAsk.PriceVET,
		Volume:    order.BestAsk.Volume,
		CreatedAt: order.BestAsk.CreatedAt.Unix(),
	}

	orderBid := &orders.Order{
		Id:        e.Id,
		Price:     order.BestBid.Price,
		PriceVET:  order.BestBid.PriceVET,
		Volume:    order.BestBid.Volume,
		CreatedAt: order.BestBid.CreatedAt.Unix(),
	}

	service.sendViaGRPC(order)
}

func (app *Config) connect() (*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:%d", app.Host, app.Port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("Error Conn")
		return nil, err
	}
	defer conn.Close()

	return conn, nil
}
