package data

import (
	"testing"
	"time"
)

// test orderbook shoud be empty
func Test_orderbook_beEmpty(t *testing.T) {
	ob := NewOrderBook()

	if ob.SizeAsk() != 0 {
		t.Errorf("Expected ask orderbook to be empty and have size of %d", ob.SizeAsk())
	}

	if ob.SizeBid() != 0 {
		t.Errorf("Expected bid orderbook to be empty and have size of %d", ob.SizeBid())
	}
}

func Test_orderBook_AddAskOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
	}

	ob.AddAskOrder(orders[0])
	if ob.SizeAsk() != 1 {
		t.Errorf("Expected SizeAsk to be %d but got %d", 1, ob.SizeAsk())
	}

	ob.AddAskOrder(orders[1])
	if ob.SizeAsk() != 2 {
		t.Errorf("Expected SizeAsk to be %d but got %d", 2, ob.SizeAsk())
	}

	ob.AddAskOrder(orders[2])
	if ob.SizeAsk() != 3 {
		t.Errorf("Expected SizeAsk to be %d but got %d", 3, ob.SizeAsk())
	}
}

func Test_orderBook_AddBidOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
	}

	ob.AddBidOrder(orders[0])
	if ob.SizeBid() != 1 {
		t.Errorf("Expected SizeBid to be %d but got %d", 1, ob.SizeBid())
	}

	ob.AddBidOrder(orders[1])
	if ob.SizeBid() != 2 {
		t.Errorf("Expected SizeBid to be %d but got %d", 2, ob.SizeBid())
	}

	ob.AddBidOrder(orders[2])
	if ob.SizeBid() != 3 {
		t.Errorf("Expected SizeBid to be %d but got %d", 3, ob.SizeBid())
	}
}

func Test_orderBook_removeAskOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
	}

	ob.AddAskOrder(orders[0])
	ob.AddAskOrder(orders[1])
	ob.AddAskOrder(orders[2])
	ob.RemoveAskOrder(orders[0])
	if ob.SizeAsk() != 2 {
		t.Errorf("Expected sizeAsk to be %d but got %d", 2, ob.SizeAsk())
	}
	ob.RemoveAskOrder(orders[1])
	ob.RemoveAskOrder(orders[2])

	if ob.SizeAsk() != 0 {
		t.Errorf("Expected sizeAsk to be %d but got %d", 0, ob.SizeAsk())
	}
}

func Test_orderBook_removeBidOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
	}

	ob.AddBidOrder(orders[0]) // 1 order
	ob.AddBidOrder(orders[1]) // 2 orders
	ob.AddBidOrder(orders[2]) // 3 orders
	ob.RemoveBidOrder(orders[0])
	if ob.SizeBid() != 2 {
		t.Errorf("Expected sizeBid to be %d but got %d", 2, ob.SizeBid())
	}
	ob.RemoveBidOrder(orders[1])
	ob.RemoveBidOrder(orders[2])

	if ob.SizeBid() != 0 {
		t.Errorf("Expected sizeBid to be %d but got %d", 0, ob.SizeBid())
	}
}

func Test_orderBook_addOrUpdateAskOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
	}
	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateAskOrder(orders[0])
	ob.AddOrUpdateAskOrder(orders[1])
	ob.AddOrUpdateAskOrder(orders[2]) // 2 orders

	if ob.SizeAsk() != 2 {
		t.Errorf("Expected sizeAsk to be %d but got %d", 0, ob.SizeAsk())
	}

	// order 0 must be updated
	ob.AddOrUpdateAskOrder(orders[3]) // same as order[0], time updated
	if ob.SizeAsk() != 2 {
		t.Errorf("Expected sizeAsk to be %d but got %d", 0, ob.SizeAsk())
	}

	// return order at idx 0 - order 0 updated with time of order 3
	idx := ob.OrderAskExist(orders[3]) // yes, idx[0]
	if idx != 0 {
		t.Errorf("Expected idx to be %d but got %d", 0, idx)
	}

	ob.RemoveAskOrder(orders[0])
	// return -1 - order not exists
	idx = ob.OrderAskExist(orders[1])
	if idx != -1 {
		t.Errorf("Expected idx to be %d but got %d", -1, idx)
	}

	// return 0
	idx = ob.OrderAskExist(orders[2])
	if idx != 0 {
		t.Errorf("Expected idx to be %d but got %d", 0, idx)
	}
}

func Test_orderBook_addOrUpdateBidOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
	}
	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateBidOrder(orders[0])
	ob.AddOrUpdateBidOrder(orders[1])
	ob.AddOrUpdateBidOrder(orders[2]) // 2 orders

	if ob.SizeBid() != 2 {
		t.Errorf("Expected sizeBid to be %d but got %d", 0, ob.SizeBid())
	}

	// order 0 must be updated
	ob.AddOrUpdateBidOrder(orders[3]) // same as order[0], time updated
	if ob.SizeBid() != 2 {
		t.Errorf("Expected sizeBid to be %d but got %d", 0, ob.SizeBid())
	}

	// return order at idx 0 - order 0 updated with time of order 3
	idx := ob.OrderBidExist(orders[3]) // yes, idx[0]
	if idx != 0 {
		t.Errorf("Expected idx to be %d but got %d", 0, idx)
	}

	ob.RemoveBidOrder(orders[0])
	// return -1 - order not exists
	idx = ob.OrderBidExist(orders[1])
	if idx != -1 {
		t.Errorf("Expected idx to be %d but got %d", -1, idx)
	}

	// return 0
	idx = ob.OrderBidExist(orders[2])
	if idx != 0 {
		t.Errorf("Expected idx to be %d but got %d", 0, idx)
	}
}

// test adding two orders with dif timestamp. Must remain one order with timestamp updated.
func Test_orderBook_AddOrUpdateAskOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
	}
	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateAskOrder(orders[0])
	ob.AddOrUpdateAskOrder(orders[1])
	ob.AddOrUpdateAskOrder(orders[2])
	ob.AddOrUpdateAskOrder(orders[3])

	if ob.SizeAsk() != 2 {
		t.Errorf("Expected sizeAsk to be %d but got %d", 0, ob.SizeAsk())
	}

	if ob.Asks[0].CreatedAt != orders[3].CreatedAt {
		t.Errorf("Expected CreatedAt to be %d but got %d", orders[0].CreatedAt, orders[3].CreatedAt)
	}
}

func Test_orderBook_AddOrUpdateBidOrder(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
	}
	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateBidOrder(orders[0])
	ob.AddOrUpdateBidOrder(orders[1])
	ob.AddOrUpdateBidOrder(orders[2])
	ob.AddOrUpdateBidOrder(orders[3])

	if ob.SizeBid() != 2 {
		t.Errorf("Expected sizeBid to be %d but got %d", 0, ob.SizeBid())
	}

	if ob.Bids[0].CreatedAt != orders[3].CreatedAt {
		t.Errorf("Expected CreatedAt to be %d but got %d", orders[0].CreatedAt, orders[3].CreatedAt)
	}
}

func Test_orderBook_bestAsk(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
		NewOrder("ABC", 90.00, 90.00, 0.5, time.Now().Unix()+int64(time.Second)),
	}

	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateAskOrder(orders[0])
	ob.AddOrUpdateAskOrder(orders[1])
	ob.AddOrUpdateAskOrder(orders[2]) // 2 orders
	ob.AddOrUpdateAskOrder(orders[3]) // 2 orders
	ob.AddOrUpdateAskOrder(orders[4]) // 2 orders

	if ob.BestAsk().Price != orders[4].Price {
		t.Errorf("Expected bestAsk to be %f but got %f", orders[4].Price, ob.BestAsk().Price)
	}
}

func Test_orderBook_bestBid(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("BCD", 120.00, 120.00, 2.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()+int64(time.Second)),
		NewOrder("ABC", 90.00, 90.00, 0.5, time.Now().Unix()+int64(time.Second)),
	}

	// adding two orders - 0 and 1 are equal, added 1
	ob.AddOrUpdateBidOrder(orders[0])
	ob.AddOrUpdateBidOrder(orders[1])
	ob.AddOrUpdateBidOrder(orders[2]) // 2 orders
	ob.AddOrUpdateBidOrder(orders[3]) // 2 orders
	ob.AddOrUpdateBidOrder(orders[4]) // 2 orders

	if ob.BestBid().Price != orders[2].Price {
		t.Errorf("Expected bestBid to be %f but got %f", orders[2].Price, ob.BestBid().Price)
	}
}

func Test_orderBook_removeExpiredAsks(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 110.00, 110.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()-TIME_TO_REMOVE_FROM_ASK+int64(time.Hour.Seconds())),
		NewOrder("ABC", 110.00, 110.00, 1.0, time.Now().Unix()-TIME_TO_REMOVE_FROM_ASK-int64(time.Hour.Seconds())),
	}

	// adding 2 orders
	ob.AddOrUpdateAskOrder(orders[0]) // 1 order
	ob.AddOrUpdateAskOrder(orders[1]) // 2 orders
	ob.AddOrUpdateAskOrder(orders[2]) // idx[0] updated and not expired
	ob.AddOrUpdateAskOrder(orders[3]) // idx[1] updated and expired

	// 2 orders
	if ob.SizeAsk() != 2 {
		t.Errorf("Expected Ask size to be %d but got %d", 1, ob.SizeAsk())
	}

	// 1 order
	size := ob.RemoveExpiredAsks()
	if ob.SizeAsk() != 1 {
		t.Errorf("Expected Ask size to be %d but got %d", 1, ob.SizeAsk())
	}

	if size != 1 {
		t.Errorf("Expected 1 size to be removed but got %d", size)
	}

}

func Test_orderBook_removeExpiredBids(t *testing.T) {
	ob := NewOrderBook()

	orders := []*Order{
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 110.00, 110.00, 1.0, time.Now().Unix()),
		NewOrder("ABC", 100.00, 100.00, 1.0, time.Now().Unix()-TIME_TO_REMOVE_FROM_ASK+int64(time.Hour.Seconds())),
		NewOrder("ABC", 110.00, 110.00, 1.0, time.Now().Unix()-TIME_TO_REMOVE_FROM_ASK-int64(time.Hour.Seconds())),
	}

	// adding 2 orders
	ob.AddOrUpdateBidOrder(orders[0]) // 1 order
	ob.AddOrUpdateBidOrder(orders[1]) // 2 orders
	ob.AddOrUpdateBidOrder(orders[2]) // idx[0] updated and not expired
	ob.AddOrUpdateBidOrder(orders[3]) // idx[1] updated and expired

	// 2 orders
	if ob.SizeBid() != 2 {
		t.Errorf("Expected Bid size to be %d but got %d", 1, ob.SizeBid())
	}

	// 1 order
	size := ob.RemoveExpiredBids()
	if ob.SizeBid() != 1 {
		t.Errorf("Expected Bid size to be %d but got %d", 1, ob.SizeBid())
	}

	if size != 1 {
		t.Errorf("Expected 1 size to be removed but got %d", size)
	}

}
