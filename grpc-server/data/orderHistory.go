package data

type AskOrder OrderHistory
type BidOrder OrderHistory

// OrderHistory contains information about an order that was offerred as the best bid or ask on a particular exchange.
// It is used to manage all the best orders received from each exchange and to calculate where is the best price for buy or sell.
type OrderHistory struct {
	ExcID    string  // Exchange ID. Ex. BINA
	Price    float64 // Price of the order
	PriceVET float64 // Price of the order with fees included. VET is calculated.
	Volume   float64 // Volume of the order
}
