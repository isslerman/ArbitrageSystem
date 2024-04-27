// Ref.: https://bitypreco.com/api#publica
// Taxas: https://suporte.bity.com.br/pt-BR/articles/6967815-taxas

// PUBLIC:
// https://docs.BitPreco.com.br/rest/v3/
// OrderBook: https://docs.BitPreco.com.br/rest/v3/#tag/Market-Data/operation/MarketsController_findOrderbook
// Markets: https://api.BitPreco.com.br/rest/v3/markets

package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"pods/internal/data"
	"time"
)

type BitPreco struct {
	apiURL    string
	Id        string
	Name      string
	FeeTaker  float64
	FeeMaker  float64
	CreatedAt int
}

func NewBitPreco() *BitPreco {
	return &BitPreco{
		apiURL:   "https://api.bitpreco.com/sol-brl/orderbook",
		Id:       "BITP",
		Name:     "BitPreco",
		FeeTaker: 0.0020,
		FeeMaker: 0.0020,
	}
}

func (e *BitPreco) BestOrder() (*data.Ask, *data.Bid, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, nil, err
	}

	var priceAsk float64
	var volumeAsk float64
	for i := 0; i < 5; i++ {
		priceAsk = apiData.Asks[0].Price
		priceAsk = priceAsk * (1 + e.FeeTaker)
		volumeAsk = apiData.Asks[0].Amount
		// test if order < minOrder
		// better to implement a medium price of agg orders
		if priceAsk*volumeAsk > 10 {
			break
		}
	}
	// fmt.Printf("DEBUG: priceAsk %f, vol %.2f\n", priceAsk, volumeAsk)

	priceBid := apiData.Bids[0].Price
	priceBid = priceBid * (1 - e.FeeTaker)
	volumeBid := apiData.Bids[0].Amount

	createdAt := time.Now()

	ask := &data.Ask{
		Price:     priceAsk,
		Volume:    volumeAsk,
		CreatedAt: createdAt,
	}

	bid := &data.Bid{
		Price:     priceBid,
		Volume:    volumeBid,
		CreatedAt: createdAt,
	}

	return ask, bid, nil
}

type apiData struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
	Id     string  `json:"id"`
}

type apiDataBitPreco struct {
	Asks      []apiData `json:"asks"`
	Bids      []apiData `json:"bids"`
	Success   bool      `json:"success"`
	Timestamp string    `json:"timestamp"`
}

func (e *BitPreco) fetchApiData() (*apiDataBitPreco, error) {
	req, err := http.Get(e.apiURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var d apiDataBitPreco
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: [%v]", d.Asks)
	return &d, nil
}
