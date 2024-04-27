// Ref.: https://docs.foxbit.com.br/
// Taxas: https://foxbit.com.br/taxas/
// Go SDK: https://github.com/foxbit-group/foxbit-api-samples/tree/main/rest-v3/go

// PUBLIC:
// https://docs.foxbit.com.br/rest/v3/
// OrderBook: https://docs.foxbit.com.br/rest/v3/#tag/Market-Data/operation/MarketsController_findOrderbook
// Markets: https://api.foxbit.com.br/rest/v3/markets

package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"pods/internal/data"
	"strconv"
	"time"
)

type Foxbit struct {
	apiURL    string
	Id        string
	Name      string
	FeeTaker  float64
	FeeMaker  float64
	CreatedAt int
}

func NewFoxbit() *Foxbit {
	return &Foxbit{
		apiURL:   "https://api.foxbit.com.br/rest/v3/markets/solbrl/orderbook",
		Id:       "FOXB",
		Name:     "Foxbit",
		FeeTaker: 0.0050,
		FeeMaker: 0.0025,
	}
}

func (e *Foxbit) BestAsk() (*data.Ask, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, err
	}
	// price with fee
	price, _ := strconv.ParseFloat(apiData.Asks[0][0], 64)
	price = price * (1 + e.FeeTaker)
	volume, _ := strconv.ParseFloat(apiData.Asks[0][1], 64)
	return &data.Ask{
		Price:     price,
		Volume:    volume,
		CreatedAt: time.Now(),
	}, nil
}

func (e *Foxbit) BestBid() (*data.Bid, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, err
	}
	price, _ := strconv.ParseFloat(apiData.Bids[0][0], 64)
	volume, _ := strconv.ParseFloat(apiData.Bids[0][1], 64)
	return &data.Bid{
		Price:     price,
		Volume:    volume,
		CreatedAt: time.Now(),
	}, nil
}

func (e *Foxbit) BestOrder() (*data.Ask, *data.Bid, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, nil, err
	}

	priceAsk, _ := strconv.ParseFloat(apiData.Asks[0][0], 64)
	priceAsk = priceAsk * (1 + e.FeeTaker)
	volumeAsk, _ := strconv.ParseFloat(apiData.Asks[0][1], 64)

	priceBid, _ := strconv.ParseFloat(apiData.Bids[0][0], 64)
	priceBid = priceBid * (1 - e.FeeTaker)
	volumeBid, _ := strconv.ParseFloat(apiData.Bids[0][1], 64)

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

type apiDataFoxBit struct {
	Asks       [][]string `json:"asks"`
	Bids       [][]string `json:"bids"`
	SequenceID int        `json:"sequence_id"`
}

func (e *Foxbit) fetchApiData() (*apiDataFoxBit, error) {
	req, err := http.Get(e.apiURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var d apiDataFoxBit
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: [%s,%s]", d.Asks[0][0], d.Asks[0][1])
	return &d, nil
}
