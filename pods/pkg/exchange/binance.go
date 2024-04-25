// Ref.: https://binance-docs.github.io/apidocs/spot/en/#change-log
// Taxas: https://www.binance.com/br/fee/trading

// PUBLIC:
// OrderBook: https://binance-docs.github.io/apidocs/spot/en/#order-book

package exchange

import (
	"encoding/json"
	"grpc-client/internal/data"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Layer: App layer - basic types
type Binance struct {
	apiURL    string
	Id        string
	Name      string
	FeeTaker  float64
	FeeMaker  float64
	CreatedAt int
}

func NewBinance() *Binance {
	return &Binance{
		apiURL:   "https://api.binance.com/api/v3/depth?limit=10&symbol=SOLBRL",
		Id:       "BINA",
		Name:     "Binance",
		FeeTaker: 0.01 * 0.10, // 0,075% with BNB active
		FeeMaker: 0.01 * 0.10,
	}
}

func (e *Binance) BestOrder() (*data.BestOrder, error) {
	apiData, err := e.fetchApiData()
	if err != nil {
		return nil, err
	}

	priceAsk, _ := strconv.ParseFloat(apiData.Asks[0][0], 64)
	priceAskVET := priceAsk * (1 + e.FeeTaker)
	volumeAsk, _ := strconv.ParseFloat(apiData.Asks[0][1], 64)

	priceBid, _ := strconv.ParseFloat(apiData.Bids[0][0], 64)
	priceBidVET := priceBid * (1 - e.FeeTaker)
	volumeBid, _ := strconv.ParseFloat(apiData.Bids[0][1], 64)

	createdAt := time.Now()

	bestAsk := &data.BestAsk{
		Price:     priceAsk,
		PriceVET:  priceAskVET,
		Volume:    volumeAsk,
		CreatedAt: createdAt,
	}

	bestBid := &data.BestBid{
		Price:     priceBid,
		PriceVET:  priceBidVET,
		Volume:    volumeBid,
		CreatedAt: createdAt,
	}

	return &data.BestOrder{
		BestAsk: bestAsk,
		BestBid: bestBid,
	}, nil
}

type apiDataBinance struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateId int        `json:"lastUpdateId"`
}

func (e *Binance) fetchApiData() (*apiDataBinance, error) {
	req, err := http.Get(e.apiURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var d apiDataBinance
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("DEBUG: [%s,%s]", d.Asks[0][0], d.Asks[0][1])
	return &d, nil
}
