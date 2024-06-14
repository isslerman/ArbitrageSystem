package ripi

import (
	"errors"
	"fmt"
	"log"
	"os"

	"grpc-server/pkg/data"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type Ripi struct {
	apiBaseURL string
	key        string
	id         string
	name       string
	FeeTaker   float64
	FeeMaker   float64
	Limits     data.OrderLimits
	client     *resty.Client
}

func New() *Ripi {
	limits := &data.OrderLimits{
		OrdMinAmount: 10,
	}

	// Get the current working directory
	// abs, _ := os.Getwd()
	// temp var for tests run
	abs := "/Users/marcosissler/projects/202404-ArbitrageSystem/grpc-server"
	envFile := fmt.Sprintf("%s/.env", abs)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("APIKEY_RIPI")

	// Return Ripi with default values
	return &Ripi{
		apiBaseURL: "https://api.ripiotrade.co/v4/",
		key:        key,
		id:         "RIPI",
		name:       "Ripio",
		FeeTaker:   0.0050,
		FeeMaker:   0.0025,
		Limits:     *limits,
		client:     resty.New(),
	}
}

func (e *Ripi) Balance(asset string) (amount float64, err error) {
	endpoint := "orders/all"
	endpoint = fmt.Sprintf("%s%s", e.apiBaseURL, endpoint)

	var res ripiResponse
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", e.key).
		SetResult(&res).
		SetError(&res).
		EnableTrace().
		Get(endpoint)

	if resp.StatusCode() != 200 {
		return 0.0, errors.New(res.Message)
	}

	if err != nil {
		return 0.0, err
	}

	return 0.0, nil
	// Iterate over the Data slice to find the assets balances
	// for _, currency := range res. {
	// 	switch currency.CurrencyCode {
	// 	case "BTC":
	// 		fmt.Printf("BTC Balance: %f\n", currency.AvailableAmount)
	// 	case "BRL":
	// 		fmt.Printf("BRL Balance: %f\n", currency.AvailableAmount)
	// 	}
	// }

}

func (e *Ripi) Id() string {
	return e.id
}

// CancelAllOrders -
func (e *Ripi) CancelAllOrders() error {
	endpoint := "orders/all"
	endpoint = fmt.Sprintf("%s%s", e.apiBaseURL, endpoint)
	// fmt.Println("RIPI CancelAllOrders")
	// fmt.Printf("DEBUG: endpoint: %s\n", endpoint)
	// fmt.Printf("DEBUG: key: %s\n", e.key)

	body := `{
		"pair": "SOL_BRL"
	}`

	// response: {"data":[{"order_id":"30F9964D-CA07-44CF-8A7A-D06138A073FE","pair_code":"BRLSOL","success":true}],"error_code":null,"message":null}
	var res cancelAllOrdersResponse
	resp, err := e.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", e.key).
		SetBody(body).
		SetResult(&res).
		SetError(&res).
		EnableTrace().
		Delete(endpoint)

	// debug
	fmt.Println("DEBUG| CancelAllOrders() err:", err)
	printResp(resp, err)

	// error - not executed
	if err != nil {
		return err
	}

	// Success:
	// return error nil, status code 400, body: {"error_code":40022,"message":"No orders to cancel"}
	// return error nil, status code 200, body: {"data":[{"order_id":"042BDF94-7073-4DD9-9106-E567628C99F9","pair_code":"BRLSOL","success":true}],"error_code":null,"message":null}

	// OK - No orders to cancel
	if resp.StatusCode() == 400 && res.ErrorCode == 40022 {
		return nil
	}

	if resp.StatusCode() == 200 && res.ErrorCode == 0 {
		return nil
	}

	return fmt.Errorf("error: StatusCode: %d - %d - %s", resp.StatusCode(), res.ErrorCode, res.Message)
}

// {"data":[{"order_id":"30F9964D-CA07-44CF-8A7A-D06138A073FE","pair_code":"BRLSOL","success":true}],"error_code":null,"message":null}
type cancelAllOrdersResponse struct {
	Data []struct {
		OrderID  string `json:"order_id"`
		PairCode string `json:"pair_code"`
		Success  bool   `json:"success"`
	} `json:"data"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// ripiResponse represents the structure of the JSON response.
type ripiResponse struct {
	Data      ripidata `json:"data"`
	ErrorCode int      `json:"error_code"`
	Message   string   `json:"message"`
}

// data represents the structure of the "data" object within the JSON.
type ripidata struct {
	CreateDate      string  `json:"create_date"`
	ExecutedAmount  float64 `json:"executed_amount"`
	ID              string  `json:"id"`
	Pair            string  `json:"pair"`
	RemainingAmount float64 `json:"remaining_amount"`
	RemainingValue  float64 `json:"remaining_value"`
	RequestedAmount float64 `json:"requested_amount"`
	Side            string  `json:"side"`
	Status          string  `json:"status"`
}

// OrdersCreate -
func (e *Ripi) CreateOrder(o *data.OrdersCreateRequest) (string, error) {
	endpoint := "orders"
	endpoint = fmt.Sprintf("%s%s", e.apiBaseURL, endpoint)
	// fmt.Println("RIPI OrdersCreate")
	// fmt.Printf("DEBUG: endpoint: %s\n", endpoint)

	var res ripiResponse
	resp, err := e.client.R().
		SetHeader("Authorization", e.key).
		SetBody(o).
		SetResult(&res).
		SetError(&res).
		EnableTrace().
		Post(endpoint)

	// debug
	fmt.Println("DEBUG| CreateOrder() order:", o)
	printResp(resp, err)
	// printRespTraceInfo(resp)
	if resp.StatusCode() == 200 {
		return res.Data.ID, nil
	}

	if resp.StatusCode() == 400 {
		return "", errors.New(res.Message)
	}

	if resp.StatusCode() == 401 {
		return "", errors.New(res.Message)
	}

	if err != nil {
		return "", err
	}

	return "", errors.New("unknown error")
}

func printResp(resp *resty.Response, err error) {
	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}

func printRespTraceInfo(resp *resty.Response) {
	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}
