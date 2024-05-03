package cex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RIPI struct {
	apiBaseURL string
	apiURL     string
	key        string
	Id         string
	Name       string
	endPoint   string
	FeeTaker   float64
	FeeMaker   float64
	Limits     limits
}

type limits struct {
	OrdMinAmount float64
}

func NewRIPI() *RIPI {
	limits := &limits{
		OrdMinAmount: 10,
	}

	return &RIPI{
		apiURL:     "https://api.ripiotrade.co/v4/public/orders/level-2?pair=SOL_BRL",
		apiBaseURL: "https://api.ripiotrade.co/v4/",
		key:        "U2FsdGVkX18sLp6CsmTh7evKaqcz44ZHHnc4Qh2tF73iepGSKxm4ZBFNPecDP1Bm",
		Id:         "RIPI",
		Name:       "Ripio",
		FeeTaker:   0.0050,
		FeeMaker:   0.0025,
		Limits:     *limits,
		endPoint:   "/orders",
	}
}

// CreateOrderResponse represents the structure of the JSON response.
type createOrderResponse struct {
	Data      data    `json:"data"`
	ErrorCode *string `json:"error_code"`
	Message   *string `json:"message"`
}

// data represents the structure of the "data" object within the JSON.
type data struct {
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

// {
//     "data": {
//         "create_date": "2024-05-02T13:00:29.497Z",
//         "executed_amount": 0,
//         "id": "47A9EFDC-CD5B-4AC2-9845-B64B5F2140C8",
//         "pair": "SOL_BRL",
//         "remaining_amount": 0.014,
//         "remaining_value": 10.28,
//         "requested_amount": 0.014,
//         "side": "sell",
//         "status": "pending_creation"
//     },
//     "error_code": null,
//     "message": null
// }

// CreateOrderRequest represents the structure of the JSON request.
type createOrderRequest struct {
	Amount float64 `json:"amount"`
	Pair   string  `json:"pair"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Type   string  `json:"type"`
}

func (ripi *RIPI) CreateOrder() (id string, err error) {

	endpoint := "orders"
	endpoint = fmt.Sprintf("%s%s", ripi.apiBaseURL, endpoint)
	fmt.Println("RIPI createOrder")
	fmt.Printf("DEBUG: endpoint: %s", endpoint)

	reqMethod := http.MethodPost

	// the body to pass
	cor := &createOrderRequest{
		Amount: 0.5,
		Pair:   "SOL_BRL",
		Price:  779.99,
		Side:   "sell",
		Type:   "limit",
	}

	// Convert the User object to JSON.
	jsonData, err := json.Marshal(cor)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(reqMethod, endpoint, body)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	req = req.WithContext(ctx)

	// Set the content type to application/json.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ripi.key)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return the response body as a string
	return string(respBody), nil
}

// request define an API request
type request struct {
	method     string
	endpoint   string
	query      url.Values
	form       url.Values
	recvWindow int64
	// secType    secType
	header  http.Header
	body    io.Reader
	fullURL string
}

// RequestOption define option type for request
type RequestOption func(*request)

// baseURL - url base
// endpoint - /orders
// fullURL - base + url
// method - get post
// auth - header
// body - json request

// func (ripi *RIPI) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
// 	err = ripi.parseRequest(r, opts...)
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	req, err := http.NewRequest(r.method, r.fullURL, r.body)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	req = req.WithContext(ctx)
// 	req.Header = r.header
// 	c.debug("request: %#v", req)
// 	f := c.do
// 	if f == nil {
// 		f = c.HTTPClient.Do
// 	}
// 	res, err := f(req)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	data, err = io.ReadAll(res.Body)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	defer func() {
// 		cerr := res.Body.Close()
// 		// Only overwrite the retured error if the original error was nil and an
// 		// error occurred while closing the body.
// 		if err == nil && cerr != nil {
// 			err = cerr
// 		}
// 	}()
// 	c.debug("response: %#v", res)
// 	c.debug("response body: %s", string(data))
// 	c.debug("response status code: %d", res.StatusCode)

// 	if res.StatusCode >= http.StatusBadRequest {
// 		apiErr := new(common.APIError)
// 		e := json.Unmarshal(data, apiErr)
// 		if e != nil {
// 			c.debug("failed to unmarshal json: %s", e)
// 		}
// 		return nil, apiErr
// 	}
// 	return data, nil
// }
