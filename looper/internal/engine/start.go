package engine

import (
	"context"
	"looper/internal/marketdata"
	"looper/internal/repository"
	"looper/internal/strategy"
	"strconv"
)

func (e *Engine) Start(ctx context.Context) error {
	//////////////////////
	// Connect to Postgres
	conn, err := repository.ConnectToPostgres(e.cfg.Database.DSN)
	if err != nil {
		e.Logger.Errorw("failed to connect to Postgres",
			"db", e.cfg.Database.Name,
			"err", err,
		)
		return err
	}
	// Assign repository to Engine
	e.Repo = repository.NewPostgresDBRepo(conn)

	// Initial checkup
	e.Repo.LogInfo("---> SERVER STARTING <---")

	// Register the gRPC Server
	md := marketdata.NewMarketData(e.Repo, e.Logger)
	_, err = marketdata.GRPCListen(md, strconv.Itoa(e.cfg.Server.GrpcPort))
	if err != nil {
		e.Logger.Fatalf("failed to start gRPC server: %v", err)
	}
	e.MarketData = md

	// Register the OMS
	// oms := oms.NewBina(ctx, cfg)
	e.Strategy = strategy.NewMakerStrategy(ctx, e.cfg, e.Logger)

	go e.Strategy.Run(ctx, e.MarketData.Updates(), e.Strategy.GetOrderStream())
	// ob := srv.Models.OrderBook

	// Strategie config and load
	// cexAsk := cex.InstanceRipi
	// cexBid := cex.InstanceBina
	// aSymbol := fmt.Sprintf("%s_%s", a.cfg.Symbols.Base, a.cfg.Symbols.Quote) // exchange A SOL_BRL <base_quote>
	// bSymbol := fmt.Sprintf("%s%s", a.cfg.Symbols.Base, a.cfg.Symbols.Quote)  // exchange B SOLBRL <basequote>

	// // creating a new AC
	// a.ac, err = NewArbitrageControl(cexAsk, cexBid, aSymbol, bSymbol, a.DB, a.Notify, a.cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return nil
}
