package main

import (
	"flag"
	"fmt"
	"grpc-server/infra/ntfy"
	"grpc-server/internal/cex"
	"grpc-server/pkg/data"
	"grpc-server/pkg/repository"
	"grpc-server/pkg/repository/dbrepo"
	"log"
	"log/slog"
	"time"
)

// Our app confiiguration.
// Errors and Infos for debug: All the erros and info msgs are saved to a DB.
type App struct {
	DB         repository.DatabaseRepo
	DSN        string     // DSN to connect to DB
	gRpcPort   string     // port to listen on for gRPC requests
	Notify     *ntfy.Ntfy // Notify sends notifications to mobile
	DryRunMode bool       // Send the orders or not
	ac         *data.ArbitrageControl
	baseToken  string // base token to use for arbitrage
	quoteToken string // quote token to use for arbitrage
}

func NewApp() *App {
	return &App{}
}

func main() {
	// starting the app
	app := App{}
	app.setupNotifyService()
	app.DryRunMode = true

	flag.StringVar(&app.DSN, "dsn", "host=aws-postgre-db.cjwmyk2c0zku.sa-east-1.rds.amazonaws.com port=5432 user=aws_postgre_db password=rfq6PlYM1NzFgmZm9QZ1 dbname=arbitrage_system sslmode=require timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.gRpcPort, "grpcPort", "50001", "gRPC Port Server Port")
	flag.StringVar(&app.baseToken, "base", "SOL", "base token to use for arbitrage (base/quote)")
	flag.StringVar(&app.quoteToken, "quote", "BRL", "quote token to use for arbitrage (base/quote)")
	flag.Parse()

	//postgre db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = dbrepo.NewPostgresDBRepo(conn)

	// Initial checkup
	app.DB.SaveLoggerInfo("Server is up - LoggerInfo Test")
	app.DB.SaveLoggerErr("Server is up - LoggerErr Test")

	// Register the gRPC Server
	srv := NewOrderServer()
	go app.gRPCListen(srv)
	ob := srv.Models.OrderBook

	// Strategie config and load
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina
	aSymbol := fmt.Sprintf("%s%s", app.baseToken, app.quoteToken)
	bSymbol := fmt.Sprintf("%s_%s", app.baseToken, app.quoteToken)
	// creating a new AC
	app.ac, err = data.NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol)
	if err != nil {
		log.Fatal(err)
	}

	pairInfo := fmt.Sprintf("Arbitrage pair is %s/%s", app.baseToken, app.quoteToken)
	slog.Info(pairInfo)
	slog.Info("Server is up and running.")
	app.DB.SaveLoggerInfo(pairInfo)
	app.DB.SaveLoggerInfo("Server is up and running.")
	app.run(ob)

}

// bestOrder receives the OrderBook over gRPC
// and find the best ask and bid orders every 2 seconds
func (app *App) run(ob *data.OrderBook) {
	// looping every 2 seconds to get the best ask and bid orders
	for {
		// filter only orders with volume >0
		if (ob.SizeAsk() > 0) && (ob.SizeBid() > 0) {
			ba := ob.BestAsk()
			bb := ob.BestBid()
			spread := (1 - (ba.PriceVET / bb.PriceVET)) * 100
			createdAt := time.Now().Unix()

			// Strategy ArbitrageControl Two Cex
			app.strategyArbitrageContol(ba, bb, spread)

			// aho and bho save all orders to the database
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

			_, err := app.DB.SaveOrderHistory(spread, aho, bho, createdAt)
			if err != nil {
				app.DB.SaveLoggerErr(fmt.Sprintf("app.bestOrder, error saving to DB, %v", err))
			}

			// if spread > 0.0 {
			fmt.Printf("Order, good,%.2f ,", spread)
			fmt.Printf("ASK[%s], %f, %f, %f, ", ba.ID, ba.Price, ba.PriceVET, ba.Volume)
			fmt.Printf("BID[%s], %f, %f, %f, ", bb.ID, bb.Price, bb.PriceVET, bb.Volume)
			fmt.Printf("time,%s\n", ba.CreatedAtHuman())

			// Strategy Buy and Sell at same time
			// if spread > 0.4 {
			// 	app.execOrder(ba, bb, spread)
			// }
			// app.execOrder(ba, bb, spread)

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

func (app *App) strategyArbitrageContol(ba, bb *data.Order, spread float64) {
	// Log info
	info := fmt.Sprintf("[SAC]New Price received [%s] - PriceVET, %f, Vol, %f", ba.ID, ba.PriceVET, ba.Volume)
	app.DB.SaveLoggerInfo(info)

	// creating the new askorder received
	a, err := data.NewAskOpenOrder(ba.Volume, ba.PriceVET, app.ac.AskSymbol, "sell")
	if err != nil {
		err := fmt.Sprintf("[SAC]Error creting NewAskOpenOrder [%s] - %f, %f, %s, sell", ba.ID, ba.Volume, ba.PriceVET, app.ac.AskSymbol)
		app.DB.SaveLoggerErr(err)
		return
	}
	app.ac.AskOpenOrder(a)
	_ = ba
	_ = bb
	_ = spread
}

// execOrder - receives the best ask and best bid and execute the order
// when the spread is higher than 0.4
// func (app *App) execOrder(ba, bb *data.Order, spread float64) {
// validations before exec the order

// getting the same volume (low) to use in both
// if ba.Volume <= bb.Volume {
// 	bb.Volume = ba.Volume
// } else {
// 	ba.Volume = bb.Volume
// }

// 1. check for open orders
// If there is already an open order, we need to check if we need to cancel it or not

// Logic to be done here
// 1. If we have any order already open, we need to check the threshol to cancel or not.
// 2. If the threshold is met, we cancel the order and set a new one
// 3. if the threshold is not met, we leave it as is
// 4. after the order is created, we need to check if the orders has been executed,
// maybe here inside the loop or outside of the loop in another goroutine?
// this goroutine will check if when the orden has been executed or canceled.
// if the order has been executed, we need to create the same order with the opposite side
// at binance.
// }

func (app *App) setupNotifyService() {
	app.Notify = ntfy.NewNtfy()
	app.Notify.SendMsg("Server is up - NTFY Test", "Server is up - LoggerErr Test", false)
}
