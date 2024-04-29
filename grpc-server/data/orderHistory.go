package data

type AskOrder OrderHistory
type BidOrder OrderHistory

type OrderHistory struct {
	ExcID    string
	Price    float64
	PriceVET float64
	Volume   float64
}
