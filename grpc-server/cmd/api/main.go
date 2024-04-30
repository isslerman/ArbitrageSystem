package main

import (
	"errors"
	"flag"
	"fmt"
	"grpc-server/data"
	"grpc-server/infra/database"
	"grpc-server/infra/ntfy"
	"log"
	"log/slog"
	"os"
	"time"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotfound   = errors.New("not found")
)

type Config struct {
	IorderHistoryRepo IOrderHistoryRepo
	DSN               string
	gRpcPort          string
	Logger            slog.Logger
	Notify            *ntfy.Ntfy
}

type IOrderHistoryRepo interface {
	Save(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error)
}

func main() {
	// configure slog

	app := Config{}
	app.setupLogger()
	app.setupNotifyService()

	flag.StringVar(&app.DSN, "dsn", "host=192.168.15.14 port=5432 user=root password=root dbname=arbitrage-system sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.gRpcPort, "grpcPort", "50001", "gRPC Port Server Port")
	flag.Parse()

	//postgre db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.IorderHistoryRepo = database.NewOrderHistory(conn)

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
			createdAt := time.Now().Unix()
			// if spread > 0.0 {
			fmt.Printf("Order, good,%.2f ,", spread)
			fmt.Printf("ASK[%s], %f, %f, %f, ", ba.ID, ba.Price, ba.PriceVET, ba.Volume)
			fmt.Printf("BID[%s], %f, %f, %f, ", bb.ID, bb.Price, bb.PriceVET, bb.Volume)
			fmt.Printf("time,%s\n", ba.CreatedAtHuman())

			if spread > 0.4 {
				app.execOrder(ba, bb, spread)
			}
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

			_, err := app.IorderHistoryRepo.Save(spread, aho, bho, createdAt)
			if err != nil {
				slog.Error("error saving to DB:", err)
			}

			// fmt.Printf("Saved to DB: [%s]\n", id)

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

func (app *Config) execOrder(ba, bb *data.Order, spread float64) {
	// How much volume have avaiable
	// if ba.Volume < bb.Volume {
	// 	slog.Info("Buy [%f]-[%f] at [%f] from %s, and Sell [%f]-[%f] at [%f] from %s\n", ba.Volume, (ba.Volume * ba.Price), ba.Price, ba.ID, bb.Volume, (bb.Volume * bb.Price), bb.Price, bb.ID)
	// } else {
	// 	slog.Info("Buy [%f]-[%f] at [%f] from %s, and Sell [%f]-[%f] at [%f] from %s\n", bb.Volume, (bb.Volume * bb.Price), bb.Price, bb.ID, ba.Volume, (ba.Volume * ba.Price), ba.Price, ba.ID)
	// }

}

func (app *Config) setupLogger() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Updates slogs default instance of slog with our own handler.

}

func (app *Config) setupNotifyService() {
	app.Notify = ntfy.NewNtfy()
}
