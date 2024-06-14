package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// {
// 	"method": "subscribe",
// 	"topics": ["order_status"],
// 	"ticket": "d90a9a10-06af-44af-8592-baf866dd1503"
// }

// MESSAGE RECEIVED AFTER EXECUTED -
// Received: {
// 	"body":{
// 		"amount":0.015,
// 		"average_execution_price":875.92,
// 		"id":"5E0D7702-1B42-4A12-993E-AFCC1B171071",
// 		"created_at":"2024-06-09T21:34:58.197Z",
// 		"executed_amount":0.015,"external_id":null,
// 		"pair":"SOL_BRL","price":875.92,
// 		"remaining_amount":0,"side":"sell",
// 		"status":"executed_completely",
// 		"type":"limit","updated_at":
// 		"2024-06-09T21:43:59.917Z",
// 		"user_id":"9EB837A5-6A81-4ADD-8B97-16672B13267A"
// 		},
// 	"timestamp":1717969441140,
// 	"topic":"order_status",
// 	"id":1036494597
// }

// status = {"executed_completely","executed_partially"}

func main() {

	ticket, err := generateTicket()
	if err != nil {
		log.Fatal(err)
	}
	// ticket := "2BDAFBA9-9DEC-4632-9078-C9732B0F3088"

	subscribeMsg := fmt.Sprintf(`{
		"method": "subscribe",
		"topics": [
		  "order_status"
		],
		"ticket": "%s"
	  }`, ticket)

	url := "wss://ws.ripiotrade.co"

	// Connect WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Channel to handle OS signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel to handle incoming messages
	done := make(chan struct{})

	// read the messages from ws
	readMessagesFromWebSocket(conn, done)

	err = conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg))
	if err != nil {
		log.Fatalf("Failed to send subscribe message: %v", err)
	}

	// Main loop to handle OS signals
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Received interrupt signal. Closing connection...")

			// Close the WebSocket connection gracefully
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Error during closing connection: %v", err)
				return
			}
			return
		}
	}
}

func readMessagesFromWebSocket(conn *websocket.Conn, done chan struct{}) {
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			var respUserTrade responseUserTrades
			err = json.Unmarshal(message, &respUserTrade)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			fmt.Printf("Received: %s\n", message)
			fmt.Println("Status", respUserTrade.Body.Status)
		}
	}()
}

// https://apidocs.ripiotrade.co/v4#tag/Orders
type responseUserTrades struct {
	ID        int    `json:"id"`
	Topic     string `json:"topic"`
	Timestamp int64  `json:"timestamp"`
	Body      struct {
		Amount                float64     `json:"amount"`
		AverageExecutionPrice float64     `json:"average_execution_price"`
		ID                    string      `json:"id"`
		CreatedAt             time.Time   `json:"created_at"`
		ExecutedAmount        float64     `json:"executed_amount"`
		ExternalID            interface{} `json:"external_id"`
		Pair                  string      `json:"pair"`
		Price                 float64     `json:"price"`
		RemainingAmount       float64     `json:"remaining_amount"`
		Side                  string      `json:"side"`
		Status                string      `json:"status"` // Array of strings (OrderStatusParamV4) "executed_completely" "executed_partially" "open" "canceled" "pending_creation"
		Type                  string      `json:"type"`
		UpdatedAt             time.Time   `json:"updated_at"`
		UserID                string      `json:"user_id"`
	} `json:"body"`
}

// TICKETS
//////////

// ticketResponse represents the top-level structure of the JSON data.
type ticketResponse struct {
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

// Data represents the nested data structure within the JSON data.
type Data struct {
	Ticket string `json:"ticket"`
}

// generateTicket uses the API post endpoint to generate a ticket number to use it with the websocket.
func generateTicket() (string, error) {
	endpoint := "https://api.ripiotrade.co/v4/ticket"
	key := getKey()
	client := resty.New()

	var res ticketResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", key).
		SetResult(&res).
		SetError(&res).
		EnableTrace().
		Post(endpoint)

	if err != nil || resp.StatusCode() != 200 {
		return "", errors.New(res.Message)
	}

	return res.Data.Ticket, nil
}

func getKey() string {
	// Get the current working directory
	// abs, _ := os.Getwd()
	// temp var for tests run
	abs := "/Users/marcosissler/projects/202404-ArbitrageSystem/grpc-server"
	envFile := fmt.Sprintf("%s/.env", abs)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("APIKEY_RIPI")
}
