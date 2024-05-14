package cex

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type BITY struct {
	apiBaseURL string
	key        string
	Id         string
	Name       string
	endPoint   string
	FeeTaker   float64
	FeeMaker   float64
	Limits     limits
}

func NewBITY() *BITY {
	limits := &limits{
		OrdMinAmount: 10,
	}

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s := os.Getenv("APISIG_BITY")
	k := os.Getenv("APIKEY_BITY")
	key := fmt.Sprintf("%s%s", s, k)

	return &BITY{
		apiBaseURL: "https://api.bitpreco.com/v1/trading/",
		key:        key,
		Id:         "BITY",
		Name:       "Bitpreco",
		FeeTaker:   0.0050,
		FeeMaker:   0.0025,
		Limits:     *limits,
		endPoint:   "/orders",
	}
}

// curl --location 'https://api.bitpreco.com/v1/trading' \
// --form 'cmd="sell"' \
// --form 'auth_token="{auth_token}"' \
// --form 'market="{market}"' \
// --form 'price="{price}"' \
// --form 'volume="{volume}"' \
// --form 'amount="{amount}"' \
// --form 'limited="{limited}"'

// func (ripi *BITY) CreateOrder(o CreateOrderRequest) (id string, err error) {

// 	endpoint := "orders"
// 	endpoint = fmt.Sprintf("%s%s", ripi.apiBaseURL, endpoint)
// 	fmt.Println("RIPI createOrder")
// 	fmt.Printf("DEBUG: endpoint: %s\n", endpoint)

// 	reqMethod := http.MethodPost

// 	hs := NewHTTPService()
// 	res, err := hs.NewRequest(reqMethod, endpoint, o, ripi.key)
// 	if err != nil {
// 		fmt.Printf("Error Sending RIPI Order: %s", err)
// 	}

// 	cor := &bityOrderResponse{}
// 	err = json.Unmarshal(res, &cor)
// 	if err != nil {
// 		log.Printf("error unmarshaling response: %+v", err)
// 		return "", err
// 	}
// 	return cor.OrderID, err
// }

// {
// 	"success": true,
// 	"order_id": "AQ1eBGRkZmNjZN",
// 	"message_cod": "ORDER_CREATED"
// }

// CreateOrderResponse represents the structure of the JSON response.
// type bityOrderResponse struct {
// 	Sucess     bool   `json:"success"`
// 	OrderID    string `json:"order_id"`
// 	MessageCod string `json:"message_cod"`
// }

func (b *BITY) CancelAllOrders() (success bool, err error) {
	// the endpoint to call
	endpoint := "all_orders_cancel"
	endpoint = fmt.Sprintf("%s%s", b.apiBaseURL, endpoint)
	fmt.Println("BITY CancelAllOrders")
	fmt.Printf("DEBUG: endpoint: %s\n", endpoint)

	// the request method to use
	reqMethod := http.MethodPost

	// creating the http service instance
	hs := NewHTTPService()

	// setting the options of the request
	ro := RequestOptions{
		FormAuthKey: b.key,
		ContentType: FORM,
	}

	resp, err := hs.NewRequest(endpoint, reqMethod, ro)
	if err != nil {
		fmt.Printf("error CancelAllOrders: %s", err)
	}
	defer resp.Body.Close()
	fmt.Println("Request successful")

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	// var caor bityCancelAllOrdersResponse
	// err = json.Unmarshal(resp.Body, &caor)
	// if err != nil {
	// 	fmt.Printf("error unmarshal: %s", err)
	// }

	// fmt.Println(caor)
	return true, nil
}

// bityCancelAllOrdersResponse represents the structure of the JSON response.
type bityCancelAllOrdersResponse struct {
	Success        bool  `json:"success"`
	OrdersCanceled int16 `json:"orders_canceled"`
}

// {
//     "success": true,
//     "orders_canceled": 0
// }

// HERE TIPS for deal with both results of json types
// REF.: https://youtu.be/Tgg-ChT4IZE?si=5WT8O7IQj85e4I-g&t=759

// wrapper for both
type Response struct {
	Success
	Error
}

type Success struct {
	Results []string `json:"results"`
}

type Error struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
