package main

import (
	"flag"
	"pods/infra/api"
	"pods/pkg/exchange"
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
	// infoLog  *log.Logger
	// errorLog *log.Logger
	// DB        repository.DatabaseRepo
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
	e := exchange.NewFoxbit()
	// e := exchange.NewFoxbit()

	// create apiSrv instance
	apiSrv := api.NewServer(e, app.config.port)
	app.apiSrv = apiSrv
	app.apiSrv.StartServer()
}

func (app *application) loadConfig() {
	app.config.env = "development"

	flag.StringVar(&app.Name, "name", "FOXB", "CEX FoxBit Exchange")
	flag.StringVar(&app.Version, "version", "1.0.0", "Pod Version")
	flag.IntVar(&app.config.port, "port", 15002, "Pod API port")
	flag.Parse()
}
