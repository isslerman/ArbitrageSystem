package main

import (
	"fmt"
	"grpc-server/data"
	"time"
)

const (
	gRpcPort = "50001"
)

type Config struct {
	Models data.Models
}

func main() {
	app := Config{}

	// Register the gRPC Server
	srv := NewOrderServer()
	go app.gRPCListen(srv)
	ob := srv.Models.OrderBook

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// go removeOldAsks(ob, ctx)
	// go removeOldBids(ob, ctx)
	app.bestOrder(ob)

	// TBD clean interrupt
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
}

func (app *Config) bestOrder(ob *data.OrderBook) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	for {
		if (ob.SizeAsk() > 0) && (ob.SizeBid() > 0) {
			ba := ob.BestAsk()
			bb := ob.BestBid()
			spread := (1 - (ba.PriceVET / bb.PriceVET)) * 100
			// if spread > 0.0 {
			fmt.Printf("Order, good,%.2f ,", spread)
			fmt.Printf("ASK[%s], %f, %f, %f, ", ba.ID, ba.Price, ba.PriceVET, ba.Volume)
			fmt.Printf("BID[%s], %f, %f, %f, ", bb.ID, bb.Price, bb.PriceVET, bb.Volume)
			fmt.Printf("time,%s\n", ba.CreatedAtTime())
			// }
			// } else {
			// 	// fmt.Printf("OB, low,%.2f , ASK[%s], %f, %f, BID[%s], %f, %f, time,%s\n", spread, ba.ID, ba.Price, ba.Volume, bb.ID, bb.Price, bb.Volume, ba.CreatedAtTime())
			// }
		}
		// Housekeeping
		ob.RemoveExpiredAsks()
		ob.RemoveExpiredBids()

		time.Sleep(time.Second * 2)

		// debug
		// fmt.Printf("OB: [%d] orders\n", ob.SizeAsk())
		// fmt.Printf("OB: [%d] orders\n", ob.SizeAsk())
	}
}
