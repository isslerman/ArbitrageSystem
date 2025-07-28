package wsconsumer

import (
	"encoding/json"
	"fmt"
	"grpc-server/infra/ntfy"
	"grpc-server/internal/cex"
	"grpc-server/pkg/data"
	"grpc-server/pkg/repository"
	"log"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WSConsumer struct {
	ticket         string
	url            string
	conn           *websocket.Conn
	db             repository.DatabaseRepo
	notify         *ntfy.Ntfy // Notify sends notifications to mobile
	cexBid         cex.Cex    // instance of exchange B
	askOrderFilled chan bool
	bidOrderFilled chan bool
}

func NewWSConsumer(db repository.DatabaseRepo, notify *ntfy.Ntfy, cexBid cex.Cex, aof, bof chan bool) *WSConsumer {
	ticket, err := generateTicket()
	if err != nil {
		log.Fatalf("FAILURE: Unable to generate ticket for WebSocket: %s", err)
	}

	url := "wss://ws.ripiotrade.co"

	// Connect WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		msg := fmt.Sprintf("FAILURE: Unable to connect to WebSocket: %v", err)
		db.SaveLoggerErr(msg)
		db.SaveLoggerInfo(msg)
		log.Fatal(msg)
	}

	db.SaveLoggerInfo("Sucessfully created wsconsumer.")
	return &WSConsumer{
		ticket:         ticket,
		url:            url,
		conn:           conn,
		db:             db,
		notify:         notify,
		cexBid:         cexBid,
		askOrderFilled: aof,
		bidOrderFilled: bof,
	}
}

func (ws *WSConsumer) Start() {
	defer ws.conn.Close()

	// Channel to handle OS signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel to handle incoming messages
	done := make(chan struct{})

	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg()))
	if err != nil {
		msg := fmt.Sprintf("FAILURE: fail to send subscribe message: %v", err)
		ws.db.SaveLoggerErr(msg)
		ws.db.SaveLoggerInfo(msg)
		log.Fatal(msg)
	}

	// read the messages from ws
	ws.readMessagesFromWebSocket(done)

	// Ping/pong mechanism to keep the connection alive
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	// Main loop to handle OS signals
	for {
		select {
		case <-ticker.C:
			err := ws.conn.WriteMessage(websocket.PingMessage, []byte(pingMsg()))
			if err != nil {
				ws.db.SaveLoggerErr(fmt.Sprintf("[WebSocket Ping] - Ping error, %s", err))
				log.Println("write:", err)
				return
			}
			ws.db.SaveLoggerInfo("[WebSocket Ping] - Ping ok.")
		case <-done:
			// Cleanly close the WebSocket connection by sending a close message and then
			return
		case <-interrupt:
			log.Println("Received interrupt signal. Closing connection...")

			// Close the WebSocket connection gracefully
			err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Error during closing connection: %v", err)
				return
			}
			return
		}
	}

}

func (ws *WSConsumer) readMessagesFromWebSocket(done chan struct{}) {
	go func() {
		defer close(done)
		for {
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				msg := fmt.Sprintf("Error reading message: %s - %s", err, string(message[:]))
				ws.db.SaveLoggerErr(msg)
				break
				// log.Fatalf("Error reading message: %v", err)
			}

			var respUserTrade responseUserTrades
			err = json.Unmarshal(message, &respUserTrade)
			if err != nil {
				msg := fmt.Sprintf("Error unmarshaling message: %s - %v", err, message)
				ws.db.SaveLoggerErr(msg)
				ws.notify.SendMsg("ERROR BINA CHECK LOG", err.Error(), true)
				log.Fatalf("Error unmarshaling message: %s - %s", err, string(message))
			}

			_, err = ws.handleMsg(respUserTrade)
			if err != nil {
				ws.db.SaveLoggerErr(fmt.Sprintf("Error handling message: %s, %v, %v ", err, message, respUserTrade))
				ws.notify.SendMsg("ERROR BINA CHECK LOG", err.Error(), true)
				log.Fatalf("Error handling message: %s, %s, %v ", err, string(message), respUserTrade)
			}
		}
	}()
}

// https://apidocs.ripiotrade.co/v4#tag/Orders
// {id: 1544541741
//
//	 topic: order_status
//	 timestamp: 1722566262653
//			{
//				amount 0.1
//				average_execution_price 0
//				id AB48E61C-E960-4C47-8EDF-817DBDEE6837
//				created_at 2024-08-02 02:37:42.59 +0000 UTC
//				executed_amount 0
//				external_id <nil>
//				pair SOL_BRL
//				price 944.66
//				remaining_amount 0.1
//				side sell
//				status waiting
//				type limit
//				updated_at 2024-08-02 02:37:42.613 +0000 UTC
//				user_id 9EB837A5-6A81-4ADD-8B97-16672B13267A
//	}}
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
		Side                  string      `json:"side"`   // sell
		Status                string      `json:"status"` // Array of strings (OrderStatusParamV4) "executed_completely" "executed_partially" "open" "canceled" "pending_creation"
		Type                  string      `json:"type"`
		UpdatedAt             time.Time   `json:"updated_at"`
		UserID                string      `json:"user_id"`
	} `json:"body"`
}

func subscribeMsg() string {
	ticket, err := generateTicket()
	if err != nil {
		log.Fatalf("FAILURE: Unable to generate ticket for WebSocket: %s", err)
	}

	return fmt.Sprintf(`{
		"method": "subscribe",
		"topics": [
		  "order_status"
		],
		"ticket": "%s"
	  }`, ticket)
}

func pingMsg() string {
	ticket, err := generateTicket()
	if err != nil {
		log.Fatalf("FAILURE: Unable to generate ticket for WebSocket: %s", err)
	}

	return fmt.Sprintf(`{
		"method": "ping",
		"topics": [
		  "order_status"
		],
		"ticket": "%s"
	  }`, ticket)
}

func (ws *WSConsumer) handleMsg(msg responseUserTrades) (string, error) {
	// DEBUG
	ws.db.SaveLoggerInfo(fmt.Sprintf("[ws msg received]: %v", msg))
	fmt.Println(msg)
	// [ws msg received]: {1541489153 order_status 1722547427528 {0.4 0 BD65722A-41E1-4F11-8355-C19B17552F30 2024-08-01 21:23:02.487 +0000 UTC 0.4 <nil> SOL_BRL 945.26 0 buy executed_completely limit 2024-08-01 21:23:02.51 +0000 UTC 9EB837A5-6A81-4ADD-8B97-
	if msg.Body.Status == "executed_completely" || msg.Body.Status == "executed_partially" {
		// { SOLBRL 0.015 867.254388 limit sell 1717899767}

		var orderSide string // buy or sell for the order - if we receive a sell ws msg, we create a buy order
		var price float64
		side := msg.Body.Side
		// ask side - if we sell on maker ripi we need to buy on bina taker
		if side == "sell" {
			orderSide = "buy"
			price = (msg.Body.Price * 1.01)
			price = math.Round(price*10) / 10
		}
		// bid side - if we buy on maker ripi we need to sell on bina taker
		if side == "buy" {
			orderSide = "sell"
			// price = (msg.Body.Price * 0.99)
			price = msg.Body.Price // trying the same price
			price = math.Round(price*10) / 10
		}

		// amount need to have max LOT_SIZE 0.001 - 3 decimals
		amount := math.Round(msg.Body.ExecutedAmount*1000) / 1000
		// NOTIONAL ERROR - we need to check if price * amount >= R$10.00
		if price*amount < 10 {
			amount = math.Round((10/price)*1000)/1000 + 0.001 // get the mininum amount to be R$10.00
			amount = math.Round(amount*1000) / 1000           // round to 3 decimals
		}
		if amount > 3 {
			log.Fatalf("wsconsumer.go|handleMsg()| WARNING: order with too high amount: 3")
		}

		order := &data.OrdersCreateRequest{
			Amount: amount,
			Pair:   "SOLBRL",
			Price:  price, // for SOLBRL needs to be with one decimal. Ex. 715.4
			Side:   orderSide,
			Type:   "limit",
		}

		// creating an order
		id, err := ws.cexBid.CreateOrder(order)
		msgInfo := fmt.Sprintf("wsconsumer.go|handleMsg()| side:%s - price:%v - id:%v - order:%v", side, price, id, order)
		ws.db.SaveLoggerInfo(msgInfo)
		if err != nil {
			ws.db.SaveLoggerErr(fmt.Sprintf("[bina order error]: %s - %v", err, order))
			return "", fmt.Errorf("wsconsumer.go|CreateOrder()| err: %s, order: %v", err, order)
		}
		msgInfo = fmt.Sprintf("BINA %s order created| id:%s - order:%v", side, id, order)
		ws.db.SaveLoggerInfo(msgInfo)
		ws.notify.SendMsg(fmt.Sprintf("BINA %s order created", side), msgInfo, false)

		if msg.Body.Status == "executed_completely" {
			if side == "sell" {
				ws.askOrderFilled <- true // channel to notify arbitrage control that the order was completely executed
			}
			if side == "buy" {
				ws.bidOrderFilled <- true // channel to notify arbitrage control that the order was completely executed
			}
		}

		return id, nil
	} else {
		ws.db.SaveLoggerInfo(fmt.Sprintf("wsconsumer.go|handleMsg()|Msg received ignored: %s", msg.Body.Status))
		return "", nil
	}
}

// order created - status waiting {1058457078 order_status 1718160319431 {1 0 26073BAF-E962-4769-9D76-70F4464E0C99 2024-06-12 02:45:18.37 +0000 UTC 0 <nil> SOL_BRL 817.8 1 sell waiting limit 2024-06-12 02:45:18.37 +0000 UTC 9EB837A5-6A81-4ADD-8B97-16672B13267A}}
