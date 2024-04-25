package main

import (
	"fmt"
	"grpc-client/internal/exchanges"
	"grpc-client/orders"
	"log/slog"
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

	e := exchanges.NewRipio()
	for {
		order, err := e.BestOrder()
		if err != nil {
			slog.Warn("MBTC: BestOrder error", err)
		} else {
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

			if orderAsk.Volume != 0 {
				app.sendViaGRPC(orderAsk, orderBid)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
