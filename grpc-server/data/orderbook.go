package data

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	TIME_TO_REMOVE_FROM_ASK = 2
	TIME_TO_REMOVE_FROM_BID = 2
)

type OrderBook struct {
	Asks []Order
	Bids []Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: []Order{},
		Bids: []Order{},
	}
}

func (o *OrderBook) AddAskOrder(order *Order) {
	o.Asks = append(o.Asks, *order)
}

func (o *OrderBook) AddBidOrder(order *Order) {
	o.Bids = append(o.Bids, *order)
}

func (o *OrderBook) UpdateIfExistAskOrder(order *Order) {
	idx := o.OrderAskExist(order)
	if idx != -1 {
		o.UpdateAskOrder(idx, order)
	}
}

func (o *OrderBook) UpdateIfExistBidOrder(order *Order) {
	idx := o.OrderBidExist(order)
	if idx != -1 {
		o.UpdateBidOrder(idx, order)
	}
}

func (o *OrderBook) AddOrUpdateAskOrder(order *Order) {
	idx := o.OrderAskExist(order)
	if idx == -1 {
		o.AddAskOrder(order)
	} else {
		o.UpdateAskOrder(idx, order)
	}
}

func (o *OrderBook) AddOrUpdateBidOrder(order *Order) {
	idx := o.OrderBidExist(order)
	if idx == -1 {
		o.AddBidOrder(order)
	} else {
		o.UpdateBidOrder(idx, order)
	}
}

func (o *OrderBook) RemoveAskOrder(order *Order) {
	for i, v := range o.Asks {
		if (v.ID == order.ID) && (v.PriceVET == order.PriceVET) && (v.Volume == order.Volume) {
			// fmt.Printf("removing order %s\n", order.ID)
			o.Asks = append(o.Asks[:i], o.Asks[i+1:]...)
			break
		}
	}
}

func (o *OrderBook) RemoveBidOrder(order *Order) {
	for i, v := range o.Bids {
		if (v.ID == order.ID) && (v.PriceVET == order.PriceVET) && (v.Volume == order.Volume) {
			o.Bids = append(o.Bids[:i], o.Bids[i+1:]...)
			break
		}
	}
}

func (o *OrderBook) OrderAskExist(order *Order) int {
	for i, v := range o.Asks {
		if order.ID == v.ID && order.PriceVET == v.PriceVET && order.Volume == v.Volume {
			return i
		}
	}
	return -1
}

func (o *OrderBook) OrderBidExist(order *Order) int {
	for i, v := range o.Bids {
		if order.ID == v.ID && order.PriceVET == v.PriceVET && order.Volume == v.Volume {
			return i
		}
	}
	return -1
}

func (o *OrderBook) UpdateAskOrder(idx int, order *Order) {
	if idx >= 0 && idx < int(o.SizeAsk()) {
		o.Asks[idx].CreatedAt = order.CreatedAt
	} else {
		fmt.Printf("index ask out of bounds. idx: %d of size: %d\n", idx, len(o.Asks))
	}
}

func (o *OrderBook) UpdateBidOrder(idx int, order *Order) {
	if idx >= 0 && idx < int(o.SizeBid()) {
		o.Bids[idx].CreatedAt = order.CreatedAt
	} else {
		slog.Info("index bid out of bounds\n")
	}
}

func (o *OrderBook) BestAsk() *Order {
	bestOrder := o.Asks[0]
	for _, order := range o.Asks {
		if order.PriceVET < bestOrder.PriceVET {
			bestOrder = order
		}
	}
	return &bestOrder
}

func (o *OrderBook) BestBid() *Order {
	bestOrder := o.Bids[0]
	for _, order := range o.Bids {
		if order.PriceVET > bestOrder.PriceVET {
			bestOrder = order
		}
	}
	return &bestOrder
}

// return the number of ask orders in the OrderBook
func (o *OrderBook) SizeAsk() int32 {
	return int32(len(o.Asks))
}

// return the number of bid orders in the OrderBook
func (o *OrderBook) SizeBid() int32 {
	return int32(len(o.Bids))
}

// remove ask orders that are expired
func (o *OrderBook) RemoveExpiredAsks() int {
	expireTime := time.Second * TIME_TO_REMOVE_FROM_ASK
	counter := 0

	for _, order := range o.Asks {
		if time.Since(time.Unix(order.CreatedAt, 0)) > expireTime {
			o.RemoveAskOrder(&order)
			counter++
		}
	}
	return counter
}

// remove bid orders that are expired
func (o *OrderBook) RemoveExpiredBids() int {
	expireTime := time.Second * TIME_TO_REMOVE_FROM_BID
	counter := 0

	for _, order := range o.Bids {
		if time.Since(time.Unix(order.CreatedAt, 0)) > expireTime {
			o.RemoveBidOrder(&order)
			counter++
		}
	}
	return counter
}
