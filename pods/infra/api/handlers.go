package api

import (
	"pods/internal/pod"

	"net/http"
	"time"

	"go.uber.org/zap"
)

func (api *Server) StartPod(w http.ResponseWriter, r *http.Request) {
	api.l.Info("Starting Pod")

	if api.isRunning {
		w.Write([]byte("Pod already running"))
		return
	}
	go api.runPod(api.stopCh)
	api.isRunning = true

	w.Write([]byte("Pod started"))
	api.l.Info("Pod started")
}

func (api *Server) StopPod(w http.ResponseWriter, r *http.Request) {
	api.l.Info("Stoping Pod")
	if !api.isRunning {
		w.Write([]byte("Pod already stopped"))
		return
	}
	// Send a signal to stop the worker goroutine.
	api.stopCh <- true
	api.isRunning = false

	w.Write([]byte("Pod stopped"))
	api.l.Info("Pod stopped")
}

func (api *Server) runPod(stopChan <-chan bool) {
	pod := pod.NewPod(api.exchange, api.l)

	for {
		select {
		case <-stopChan:
			api.l.Info("We must stop")
			return
		default:
			// run the pod
			err := pod.Run()
			if err != nil {
				api.l.Error("error running pod", zap.Error(err))
			}

			time.Sleep(time.Second * 1)
		}
	}
}

type OrderPayload struct {
	ID     string
	Price  float64
	Volume float64
}
type status struct {
	Status bool `json:"status"`
}

func (api *Server) Status(w http.ResponseWriter, r *http.Request) {
	api.l.Info("Handler Status")
	status := &status{
		Status: api.isRunning,
	}
	_ = api.tb.WriteJSON(w, http.StatusOK, status)
}
