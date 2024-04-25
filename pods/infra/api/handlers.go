package api

import (
	"context"
	"fmt"
	"grpc-client/internal/pb/orders"
	"log"

	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// 
func (api *Server) runPod(stopChan <-chan bool) {
	fmt.Println("I am doing somethin here.")
	for {
		select {
		case <-stopChan:
			fmt.Println("We must stop...")
			return
		default:
			fmt.Println("Running...")
			fmt.Printf("Pod loaded with %s config", api.exchange)

			
			SEND GRPC? 

			// here we need to execute the application service, that is?
			// feth data from exchange, send to Pod to process the orderbook and return the result.
			// the result is sent over gRPC
			// input data is only basic types

			fmt.Println("working here")
			// order, err := e.BestOrder()
			// if err != nil {
			// 	slog.Warn("RIPI: BestOrder error", err)
			// } else {
			// 	fmt.Printf("%s BestAsk: %f, Volume: %f\n", e.Id, order.BestAsk.Price, order.BestAsk.Volume)
			// 	fmt.Printf("%s BestBid: %f, Volume: %f\n", e.Id, order.BestBid.Price, order.BestBid.Volume)

			// 	orderAsk := &orders.Order{
			// 		Id:        "api.exchange",
			// 		Price:     order.BestAsk.Price,
			// 		PriceVET:  order.BestAsk.PriceVET,
			// 		Volume:    order.BestAsk.Volume,
			// 		CreatedAt: order.BestAsk.CreatedAt.Unix(),
			// 	}

			// 	orderBid := &orders.Order{
			// 		Id:        e.Id,
			// 		Price:     order.BestBid.Price,
			// 		PriceVET:  order.BestBid.PriceVET,
			// 		Volume:    order.BestBid.Volume,
			// 		CreatedAt: order.BestBid.CreatedAt.Unix(),
			// 	}

			// 	if orderAsk.Volume != 0 {
			// 		grpc.SendViaGRPC(orderAsk, orderBid)
			// 	}
			// }
			time.Sleep(time.Second * 1)
		}
	}
}

type OrderPayload struct {
	ID     string
	Price  float64
	Volume float64
}

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	GRPC *grpc.ClientConn
}

func (api *Server) sendViaGRPC(orderAsk, orderBid *orders.Order) {
	app := Config{
		Host: "localhost",
		Port: 50001,
	}
	target := fmt.Sprintf("%s:%d", app.Host, app.Port)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("Error Conn")
	}
	defer conn.Close()

	// client
	c := orders.NewOrderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.WriteOrder(ctx, &orders.OrderRequest{
		OrderAsk: orderAsk,
		OrderBid: orderBid,
	})
	if err != nil {
		log.Fatal("Error WriteOrder: ", err)
	}
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
