package main

import (
	"context"
	"log"
	"log/slog"
	"looper/config"
	"looper/internal/engine"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		return
	}

	e := engine.New(cfg)

	if err := e.Start(ctx, stop); err != nil {
		slog.Error("service failed", "err", err)
		return
	}

	slog.Info("Server is up and running.")

	// Run the main process
	// go e.run(ctx, app.orderSrv)

	<-ctx.Done()
	e.GracefullyShutdown()
}

func initEnv() {
	projectRoot, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Failed to get working directory")
	}

	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Failed to load env file at %s: %v", envPath, err)
	}
}
