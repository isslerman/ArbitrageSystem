// Ref.: https://api.mercadobitcoin.net/api/v4/docs
// Taxas: https://www.mercadobitcoin.com.br/taxas-contas-limites#tabela-taxas

// PUBLIC:
// OrderBook: https://api.mercadobitcoin.net/api/v4/{symbol}/orderbook
// Symbols: https://api.mercadobitcoin.net/api/v4/symbols

package exchange

import (
	"encoding/json"
	"grpc-client/internal/data"
	"io"
	"net/http"
	"strconv"
	"time"
)

type MercadoBitcoin struct {
	apiURL   string
	Id       string
	Name     string
	FeeTaker float64
	FeeMaker float64
}

func NewMercadoBitcoin() *MercadoBitcoin {
	return &MercadoBitcoin{
		apiURL:   "https://api.mercadobitcoin.net/api/v4/SOL-BRL/orderbook",
		Id:       "MBTC",
		Name:     "Mercado Bitcoin",
		FeeTaker: 0.007,
		FeeMaker: 0.003,
	}
}

func (e *MercadoBitcoin) BestOrder() (*data.BestOrder, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, err
	}

	if len(apiData.Asks) == 0 || len(apiData.Bids) == 0 {
		return nil, nil
	}

	priceAsk, _ := strconv.ParseFloat(apiData.Asks[0][0], 64)
	priceAsk = priceAsk * (1 + e.FeeTaker)
	volumeAsk, _ := strconv.ParseFloat(apiData.Asks[0][1], 64)

	priceBid, _ := strconv.ParseFloat(apiData.Bids[0][0], 64)
	priceBid = priceBid * (1 - e.FeeTaker)
	volumeBid, _ := strconv.ParseFloat(apiData.Bids[0][1], 64)

	createdAt := time.Now()

	bestAsk := &data.BestAsk{
		Price:     priceAsk,
		Volume:    volumeAsk,
		CreatedAt: createdAt,
	}

	bestBid := &data.BestBid{
		Price:     priceBid,
		Volume:    volumeBid,
		CreatedAt: createdAt,
	}

	return &data.BestOrder{
		BestAsk: bestAsk,
		BestBid: bestBid,
	}, nil
}

type apiDataMB struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp int        `json:"timestamp"`
}

func (e *MercadoBitcoin) fetchApiData() (*apiDataMB, error) {
	req, err := http.Get(e.apiURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var d apiDataMB
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("[%s,%s]", d.Asks[0][0], d.Asks[0][1])
	return &d, nil
}
