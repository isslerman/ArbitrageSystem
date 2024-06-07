package main

import (
	"flag"
	"fmt"
	"grpc-server/infra/grpc"
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
	ac         *ArbitrageControl
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
	app.DB.SaveLoggerInfo("---> SERVER STARTING <---")
	app.DB.SaveLoggerErr("---> SERVER STARTING <---")

	// Register the gRPC Server
	srv := grpc.NewOrderServer()
	go grpc.GRPCListen(srv, app.gRpcPort)
	ob := srv.Models.OrderBook

	// Strategie config and load
	cexAsk := cex.InstanceRipi
	cexBid := cex.InstanceBina
	aSymbol := fmt.Sprintf("%s%s", app.baseToken, app.quoteToken)
	bSymbol := fmt.Sprintf("%s_%s", app.baseToken, app.quoteToken)
	// creating a new AC
	app.ac, err = NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, app.DB)
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
	info := fmt.Sprintf("[New Price received] [%s] - PriceVET, %f, Vol, %f", ba.ID, ba.PriceVET, ba.Volume)
	app.DB.SaveLoggerInfo(info)

	// creating the new askorder received
	a, err := data.NewAskOpenOrder(ba.Volume, ba.PriceVET, app.ac.AskSymbol, "limit")
	if err != nil {
		err := fmt.Sprintf("[Error creating NewAskOpenOrder] [%s] - %f, %f, %s, sell, %s", ba.ID, ba.Volume, ba.PriceVET, app.ac.AskSymbol, err)
		app.DB.SaveLoggerErr(err)
		return
	}
	app.ac.AskOpenOrder(a)
	_ = bb
	_ = spread
}

// setupNotifyService send a msg to ntfy mobile app
func (app *App) setupNotifyService() {
	// TEMPORARILY DISABLED - REMOVE THE TWO COMMENTS TO ENABLE

	// app.Notify = ntfy.NewNtfy()
	// app.Notify.SendMsg("Server is up - NTFY Test", "Server is up - LoggerErr Test", false)
}
