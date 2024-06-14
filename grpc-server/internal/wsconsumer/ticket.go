package wsconsumer

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

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
