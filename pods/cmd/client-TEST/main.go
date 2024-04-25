package main

import (
	"grpc-client/orders"
	"time"
)

// Just a gRPC client that connects with our server and send a msg.

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	app := Config{
		Host: "localhost",
		Port: 50001,
	}

	for {
		orderAsk := []*orders.Order{
			{Id: "FOXB", Price: 100.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "MBTC", Price: 110.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "BITP", Price: 115.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "RIPI", Price: 118.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
		}

		orderBid := []*orders.Order{
			{Id: "FOXB", Price: 80.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "MBTC", Price: 90.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "BITP", Price: 88.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
			{Id: "RIPI", Price: 60.00, Volume: 1.0, CreatedAt: time.Now().Unix()},
		}

		app.sendViaGRPC(orderAsk[0], orderBid[0])
		app.sendViaGRPC(orderAsk[1], orderBid[1])
		app.sendViaGRPC(orderAsk[2], orderBid[2])
		app.sendViaGRPC(orderAsk[3], orderBid[3])
		time.Sleep(time.Millisecond * 1000 / 3)
	}
}
