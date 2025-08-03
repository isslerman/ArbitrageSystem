package marketdata

type BookTrade struct {
	Timestamp  int64 // in nanoseconds
	LastPrice  float64
	LastSide   string // "buy" or "sell"
	LastVolume float64

	AskPrices  []float64
	AskVolumes []float64
	BidPrices  []float64
	BidVolumes []float64

	// Derived features (optional at runtime)
	BestAsk float64
	BestBid float64
	// Spread      float64
	// Skew        float64
	// Curtosis    float64
	// QuoteVolume float64
}
