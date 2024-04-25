package main

import (
	"flag"
	"grpc-client/infra/api"
	"grpc-client/pkg/exchanges"
	"log"
)

type application struct {
	Name    string
	Version string
	port    int
	api     api.ApiService
}

func main() {
	var app application
	flag.StringVar(&app.Name, "name", "Binance", "CEX Binance Exchange")
	flag.StringVar(&app.Version, "version", "1.0.0", "Pod Version")
	flag.IntVar(&app.port, "port", 15000, "Pod API port")
	flag.Parse()

	e := exchanges.NewBinance()
	log.Printf("Starting api on port %d\n", app.port)
	api := api.NewApiService(e, app.port)
	api.StartServer()
}
