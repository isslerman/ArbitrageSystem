package bestOrder

type CreateBestOrderInput struct {
	
}

// Layer: Business Layer
// Complex types, parsed from Application layer to guarantee data integrity
type BestAsk struct {
	Price     float64
	PriceVET  float64
	Volume    float64
	CreatedAt time.Time
}

type BestBid struct {
	Price     float64
	PriceVET  float64
	Volume    float64
	CreatedAt time.Time
}

type BestOrder struct {
	BestAsk   *BestAsk
	BestBid   *BestBid
	CreatedAt time.Time
}