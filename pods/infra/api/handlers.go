package api

import (
	"fmt"
	"pods/internal/pod"

	"net/http"
	"time"
)

func (api *Server) StartPod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting Pod")
	if api.isRunning {
		fmt.Println("Pod already running")
		w.Write([]byte("Pod already running"))
		return
	}
	go api.runPod(api.stopCh)
	api.isRunning = true

	w.Write([]byte("Pod started"))
	fmt.Println("leaving StartPod Func")
}

func (api *Server) StopPod(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stoping Pod")
	if !api.isRunning {
		fmt.Println("Pod already stopped")
		w.Write([]byte("Pod already stopped"))
		return
	}
	// Send a signal to stop the worker goroutine.
	api.stopCh <- true
	api.isRunning = false

	w.Write([]byte("Pod stopped"))
	fmt.Println("Pod stopped")
}

func (api *Server) runPod(stopChan <-chan bool) {
	fmt.Println("I am doing somethin here.")

	pod := pod.NewPod(api.exchange)

	for {
		select {
		case <-stopChan:
			fmt.Println("We must stop...")
			return
		default:
			// fmt.Println("Running...")
			// run the pod
			err := pod.Run()
			if err != nil {
				fmt.Println("error:", err)
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
	Status string `json:"status"`
}

func (api *Server) Status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status ping")
	status := &status{
		Status: "ok",
	}
	_ = api.tb.WriteJSON(w, http.StatusOK, status)
}
