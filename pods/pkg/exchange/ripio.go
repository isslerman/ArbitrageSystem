// Ref.: https://apidocs.ripiotrade.co/v4#
// Taxas:
// Grafico: https://trade.ripio.com/market/market-out?pairCode=SOLBRL

// PUBLIC:
// https://api.ripiotrade.co/v4
// OrderBook: https://api.ripiotrade.co/v4/public/orders/level-2
// Markets: https://api.ripiotrade.co/v4/public/pairs

package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"pods/internal/data"
	"time"
)

type Ripio struct {
	apiURL   string
	Id       string
	Name     string
	FeeTaker float64
	FeeMaker float64
	// vol min to trade? how to make aggregated orders?
}

func NewRipio() *Ripio {
	return &Ripio{
		apiURL:   "https://api.ripiotrade.co/v4/public/orders/level-2?pair=SOL_BRL",
		Id:       "RIPI",
		Name:     "Ripio",
		FeeTaker: 0.0050,
		FeeMaker: 0.0025,
	}
}

func (e *Ripio) BestOrder() (*data.Ask, *data.Bid, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, nil, err
	}
	var priceAsk, priceBid float64
	var priceAskVET, priceBidVET float64
	var volumeAsk, volumeBid float64

	// try 5 orders with the min volume
	for i := 0; i < 5; i++ {
		priceAsk = apiData.Data.Asks[i].Price
		priceAskVET = priceAsk * (1 + e.FeeTaker)
		volumeAsk = apiData.Data.Asks[i].Amount

		// test if order < minOrder
		// better to implement a medium price of agg orders
		if priceAsk*volumeAsk > 10 {
			break
		}
	}

	// try 5 orders with the min volume
	for i := 0; i < 5; i++ {
		priceBid = apiData.Data.Bids[i].Price
		priceBidVET = priceBid * (1 - e.FeeTaker)
		volumeBid = apiData.Data.Bids[i].Amount

		// test if order < minOrder
		// better to implement a medium price of agg orders
		if priceAsk*volumeBid > 10 {
			break
		}
	}

	createdAt := time.Now()

	ask := &data.Ask{
		Price:     priceAsk,
		PriceVET:  priceAskVET,
		Volume:    volumeAsk,
		CreatedAt: createdAt,
	}

	bid := &data.Bid{
		Price:     priceBid,
		PriceVET:  priceBidVET,
		Volume:    volumeBid,
		CreatedAt: createdAt,
	}

	return ask, bid, nil
}

func (e *Ripio) ExchangeID() string {
	return e.Id
}

type orderRipio struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type dataRipio struct {
	Asks      []orderRipio `json:"asks"`
	Bids      []orderRipio `json:"bids"`
	Hash      string       `json:"hash"`
	Timestamp int          `json:"timestamp"`
}

type apiDataRipio struct {
	Data      dataRipio `json:"data"`
	ErrorCode string    `json:"error_code"`
	Message   string    `json:"message"`
}

func (e *Ripio) fetchApiData() (*apiDataRipio, error) {
	req, err := http.Get(e.apiURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var d apiDataRipio
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: [%f]", d.Data.Asks)
	return &d, nil
}
