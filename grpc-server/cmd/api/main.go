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
	"time"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotfound   = errors.New("not found")
)

// Our app confiiguration.
// Errors and Infos for debug: All the erros and info msgs are saved to a DB.
type App struct {
	IorderHistoryRepo IOrderHistoryRepo
	DSN               string          // DSN to connect to DB
	gRpcPort          string          // port to listen on for gRPC requests
	LoggerErr         ILoggerErrRepo  // Repo to save errors to DB
	LoggerInfo        ILoggerInfoRepo // Repo to save info to DB
	Notify            *ntfy.Ntfy      // Notify sends notifications to mobile
	DryRunMode        bool            // Send the orders or not
	ArbitrageOrder    *data.ArbitrageOrder
}

func NewApp() *App {
	return &App{}
}

type IOrderHistoryRepo interface {
	Save(spread float64, ask *data.AskOrder, bid *data.BidOrder, createdAt int64) (string, error)
}

type ILoggerErrRepo interface {
	Save(log string)
}

type ILoggerInfoRepo interface {
	Save(log string)
}

func main() {
	// configure slog

	app := App{}
	app.setupNotifyService()
	app.DryRunMode = true

	flag.StringVar(&app.DSN, "dsn", "host=aws-postgre-db.cjwmyk2c0zku.sa-east-1.rds.amazonaws.com port=5432 user=aws_postgre_db password=rfq6PlYM1NzFgmZm9QZ1 dbname=arbitrage_system sslmode=require timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.gRpcPort, "grpcPort", "50001", "gRPC Port Server Port")
	flag.Parse()

	//postgre db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.IorderHistoryRepo = database.NewOrderHistory(conn)
	app.LoggerErr = database.NewLoggerErrRepo(conn)
	app.LoggerInfo = database.NewLoggerInfoRepo(conn)

	// Initial checkup
	app.LoggerErr.Save("Server is up - LoggerErr Test")
	app.LoggerInfo.Save("Server is up - LoggerInfo Test")

	// Register the gRPC Server
	srv := NewOrderServer()
	go app.gRPCListen(srv)
	ob := srv.Models.OrderBook

	slog.Info("Server is up and running.")
	app.bestOrder(ob)

}

// bestOrder receives the OrderBook and gets the best ask and bid orders.
func (app *App) bestOrder(ob *data.OrderBook) {
	// looping every 2 seconds to get the best ask and bid orders
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

			// our filter to execute the order
			// what will be our rules?
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
				app.LoggerErr.Save(fmt.Sprintf("app.bestOrder, error saving to DB, %v", err))
				// slog.Error("error saving to DB:", err)
			}
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

// execOrder - receives the best ask and best bid and execute the order
func (app *App) execOrder(ba, bb *data.Order, spread float64) {
	// validations before exec the order

	// getting the same volume (low) to use in both
	if ba.Volume <= bb.Volume {
		bb.Volume = ba.Volume
	} else {
		ba.Volume = bb.Volume
	}

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
	if app.hasArbitrageOrder() {
		return
	} else {
		app.LoggerInfo.Save("ArbitrageOrder Exec")
		app.LoggerInfo.Save(fmt.Sprintf("%f", spread))
		// ao := data.NewArbitrageOrder(ba, bo)
		// res := ao.Exec()
	}
}

// func (app *Config) setupLogger() {
// 	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
// 	// slog.SetDefault(logger) // Updates slogs default instance of slog with our own handler.
// }

func (app *App) setupNotifyService() {
	app.Notify = ntfy.NewNtfy()
	app.Notify.SendMsg("Server is up - NTFY Test", "Server is up - LoggerErr Test", false)
}

func (app *App) hasArbitrageOrder() bool {
	return (app.ArbitrageOrder == nil)
}
