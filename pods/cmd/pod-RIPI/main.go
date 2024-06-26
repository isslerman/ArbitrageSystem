package main

import (
	"flag"
	"fmt"
	"pods/infra/api"
	"pods/pkg/exchange"
	"pods/pkg/logger"

	"go.uber.org/zap"
)

type config struct {
	port int
	env  string
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

	// create the exchange instance to be injected into the application
	e := exchange.NewRipio()

	// create apiSrv instance
	apiSrv := api.NewServer(e, app.config.port, app.l)
	app.apiSrv = apiSrv
	app.apiSrv.StartServer()
}

func (app *application) loadConfig() {
	app.config.env = "development"

	flag.StringVar(&app.Name, "name", "RIPI", "CEX Ripio Exchange")
	flag.StringVar(&app.Version, "version", "1.0.0", "Pod Version")
	flag.IntVar(&app.config.port, "port", 15004, "Pod API port")
	flag.Parse()

	app.l = logger.Get(app.Name)
	app.l.Info("logging files info",
		zap.String("dir", fmt.Sprintf("logs/%s.log", app.Name)))
}
