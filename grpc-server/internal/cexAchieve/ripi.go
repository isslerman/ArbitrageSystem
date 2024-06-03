package cexAchieve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type RIPI struct {
	apiBaseURL string
	key        string
	Id         string
	Name       string
	endPoint   string
	FeeTaker   float64
	FeeMaker   float64
	Limits     limits
}

func NewRIPI() *RIPI {
	limits := &limits{
		OrdMinAmount: 10,
	}

	key := os.Getenv("APIKEY_RIPI")
	return &RIPI{
		apiBaseURL: "https://api.ripiotrade.co/v4/",
		key:        key,
		Id:         "RIPI",
		Name:       "Ripio",
		FeeTaker:   0.0050,
		FeeMaker:   0.0025,
		Limits:     *limits,
		endPoint:   "/orders",
	}
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

func (r *RIPI) CreateOrder(o CreateOrderRequest) (id string, err error) {
	endpoint := "orders"
	endpoint = fmt.Sprintf("%s%s", r.apiBaseURL, endpoint)
	fmt.Println("RIPI createOrder")
	fmt.Printf("DEBUG: endpoint: %s\n", endpoint)

	reqMethod := http.MethodPost

	hs := NewHTTPService()
	res, err := hs.NewRequest(reqMethod, endpoint, "json", o, r.key, nil)
	if err != nil {
		fmt.Printf("Error Sending RIPI Order: %s", err)
	}

	or := &ripiOrderResponse{}
	err = json.Unmarshal(res, &or)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return "", err
	}
	return or.Data.ID, err
}

// CreateOrderResponse represents the structure of the JSON response.
type ripiOrderResponse struct {
	Data      ripidata `json:"data"`
	ErrorCode *string  `json:"error_code"`
	Message   *string  `json:"message"`
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

func (r *RIPI) CancelAllOrders() (id string, err error) {
	endpoint := "orders/all"
	endpoint = fmt.Sprintf("%s%s", r.apiBaseURL, endpoint)
	fmt.Println("RIPI CancelAllOrders")
	fmt.Printf("DEBUG: endpoint: %s\n", endpoint)

	reqMethod := http.MethodDelete

	hs := NewHTTPService()
	res, err := hs.NewRequest(reqMethod, endpoint, "json", nil, r.key, nil)
	if err != nil {
		fmt.Printf("Error Sending RIPI Order: %s", err)
	}

	or := &ripiOrderResponse{}
	err = json.Unmarshal(res, &or)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return "", err
	}
	return or.Data.ID, err

	return "", nil
}

// {
//     "data": [
//         {
//             "order_id": "D3D9E29D-B080-4211-974D-BD5B74335C5F",
//             "pair_code": "BRLSOL",
//             "success": true
//         }
//     ],
//     "error_code": null,
//     "message": null
// }

// RequestOption define option type for request
// type RequestOption func(*request)

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
