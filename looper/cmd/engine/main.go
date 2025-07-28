package main

import (
	"context"
	"log/slog"
	"looper/config"
	"looper/internal/engine"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		return
	}

	e := engine.New(cfg)

	if err := e.Start(ctx); err != nil {
		slog.Error("service failed", "err", err)
		return
	}

	slog.Info("Server is up and running.")

	// Run the main process
	// go e.run(ctx, app.orderSrv)

	<-ctx.Done()
	e.GracefullyShutdown()
}
