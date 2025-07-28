package engine

import (
	"context"
	"looper/internal/marketdata"
	"looper/internal/repository"
	"strconv"

	"github.com/gookit/slog"
)

func (e *Engine) Start(ctx context.Context) error {
	// Connect to Postgres
	conn, err := repository.ConnectToPostgres(e.cfg.Database.DSN)
	if err != nil {
		slog.Error("failed to connect to Postgres", "db", e.cfg.Database.Name, "err", err)
		return err
	}

	// Assign repository to Engine
	e.Repo = repository.NewPostgresDBRepo(conn)

	// Ensure DB is closed on shutdown
	go func() {
		<-ctx.Done()
		slog.Info("Closing database connection")
		_ = conn.Close()
	}()

	// Initial checkup
	e.Repo.LogInfo("---> SERVER STARTING <---")

	// Register the gRPC Server
	srv := marketdata.NewOrderServer()
	slog.Debug("config", "cfg", e.cfg)
	go marketdata.GRPCListen(srv, strconv.Itoa(e.cfg.Server.GrpcPort))

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
