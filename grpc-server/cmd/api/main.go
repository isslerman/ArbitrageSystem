package main

import (
	"database/sql"
	"fmt"
	"grpc-server/data"
	"grpc-server/infra/database"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	gRpcPort = "50001"
)

type Config struct {
	IorderHistoryRepo IOrderHistoryRepo
}

type IOrderHistoryRepo interface {
	Save(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt string) (string, error)
}

func main() {
	app := Config{}

	//sqlite db
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app.IorderHistoryRepo = database.NewOrderHistory(db)

	// Register the gRPC Server
	srv := NewOrderServer()
	go app.gRPCListen(srv)
	ob := srv.Models.OrderBook

	app.bestOrder(ob)
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

			aho := &data.AskOrder{
				ExcID:    ba.ID,
				Price:    ba.Price,
				PriceVET: ba.PriceVET,
				Volume:   ba.Volume,
			}
			bho := &data.BidOrder{
				ExcID:    bb.ID,
				Price:    bb.Price,
				PriceVET: bb.PriceVET,
				Volume:   bb.Volume,
			}
			id, err := app.IorderHistoryRepo.Save(spread, aho, bho, ba.CreatedAtTime())
			if err != nil {
				fmt.Printf("Error saving to DB: %v", err)
			}

			fmt.Printf("Saved to DB: [%s]\n", id)

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
