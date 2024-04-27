// Package api is responsible for handling incoming requests.
// Routing is a concern of the infrastructure
package api

import (
	"fmt"
	"log"
	"net/http"

	"pods/internal/pod"
	"pods/pkg/toolbox"

	"pods/pkg/exchange"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	tb        toolbox.Tools
	stopCh    chan bool
	isRunning bool
	exchange  exchange.IExchange
	port      int
	pod       *pod.Pod
}

// NewServer initializes a new Server
func NewServer(e exchange.IExchange, port int) Server {
	return Server{
		tb:       toolbox.New(),   // toolbox utilities
		stopCh:   make(chan bool), // channel to help with start/stop goroutine
		exchange: e,
		port:     port,
		pod:      nil,
		// app config?
		// log is part of the app or the api service? R: App
	}
}

func (api *Server) StartServer() {

	err := http.ListenAndServe(fmt.Sprintf(":%d", api.port), api.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

func (api *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	// remove one or leave both?
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/status", api.Status)

	mux.HandleFunc("/start", api.StartPod)
	mux.HandleFunc("/stop", api.StopPod)

	return mux
}
