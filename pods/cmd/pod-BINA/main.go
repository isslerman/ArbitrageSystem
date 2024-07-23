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

type config struct {
	Name        string
	Version     string
	Port        int
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

func newApplication() *application {
	app := &application{}
	app.loadConfig()

	return app
}

func main() {
	// initiate our application
	app := newApplication()

	symbol := fmt.Sprintf("%s%s", app.config.BaseSymbol, app.config.QuoteSymbol)
	// create the exchange instance to be injected into the application
	e := exchange.NewBinance(symbol, app.config.FeeTaker, app.config.FeeMaker)

	// create apiSrv instance
	apiSrv := api.NewServer(e, app.config.Port, app.l)
	app.apiSrv = apiSrv
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

	// init logger
	app.l = logger.Get(app.Name)

	return nil
}
