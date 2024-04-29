// Ref.: https://api.mercadobitcoin.net/api/v4/docs
// Taxas: https://www.mercadobitcoin.com.br/taxas-contas-limites#tabela-taxas

// PUBLIC:
// OrderBook: https://api.mercadobitcoin.net/api/v4/{symbol}/orderbook
// Symbols: https://api.mercadobitcoin.net/api/v4/symbols

package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"pods/internal/data"
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
		apiURL: "https://api.mercadobitcoin.net/api/v4/SOL-BRL/orderbook",
		Id:     "MBTC",
		Name:   "Mercado Bitcoin",
		// FeeTaker: 0.007,
		FeeTaker: 0.002,
		// FeeMaker: 0.003,
		FeeMaker: 0.001,
	}
}

func (e *MercadoBitcoin) BestOrder() (*data.Ask, *data.Bid, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, nil, err
	}

	if len(apiData.Asks) == 0 || len(apiData.Bids) == 0 {
		return nil, nil, nil
	}

	priceAsk, _ := strconv.ParseFloat(apiData.Asks[0][0], 64)
	priceAskVET := priceAsk * (1 + e.FeeTaker)
	volumeAsk, _ := strconv.ParseFloat(apiData.Asks[0][1], 64)

	priceBid, _ := strconv.ParseFloat(apiData.Bids[0][0], 64)
	priceBidVET := priceBid * (1 - e.FeeTaker)
	volumeBid, _ := strconv.ParseFloat(apiData.Bids[0][1], 64)

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

func (e *MercadoBitcoin) ExchangeID() string {
	return e.Id
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
