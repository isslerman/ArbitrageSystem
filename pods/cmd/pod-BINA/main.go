package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pods/infra/api"
	"pods/pkg/exchange"
	"pods/pkg/logger"

	"go.uber.org/zap"
)

// Pod is an API service that collects data from the exchange/source, proccess and send
// to the grpc server

// Application config
type config struct {
	Name        string
	Version     string
	Port        int     `json:"api_port"`
	BaseSymbol  string  `json:"base_symbol"`
	QuoteSymbol string  `json:"quote_symbol"`
	FeeTaker    float64 `json:"fee_taker"`
	FeeMaker    float64 `json:"fee_maker"`
}

type application struct {
	Name    string
	Version string
	config  config
	apiSrv  api.Server
	l       *zap.Logger
}

// Start the application and load the config
func newApplication() *application {
	app := &application{}
	app.loadConfig()

	// init logger
	app.l = logger.Get(app.Name)

	return app
}

func main() {
	// initiate our application
	app := newApplication()

	symbol := fmt.Sprintf("%s%s", app.config.BaseSymbol, app.config.QuoteSymbol)
	// create the exchange instance to be injected into the application
	e := exchange.NewBinance(symbol, app.config.FeeTaker, app.config.FeeMaker)

	// create server api and start
	app.apiSrv = api.NewServer(e, app.config.Port, app.l)
	app.apiSrv.StartServer()
}

func (app *application) loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error loading config file, %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app.config)
	if err != nil {
		log.Fatalf("Error loading config file, %s", err)
	}

	return nil
}
