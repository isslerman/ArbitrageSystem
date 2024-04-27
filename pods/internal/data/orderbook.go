// Reference to ArdanLabs project template
// https://github.com/ardanlabs/service/blob/master/business/domain/userbus/model.go
package data

import (
	"time"
)

// Layer: Business Layer
// Complex types, parsed from Application layer to guarantee data integrity
type Ask Order

type Bid Order

type Orderbook struct {
	Ask       *Ask
	Bid       *Bid
	CreatedAt time.Time
}

func NewOrderBook(a *Ask, b *Bid) *Orderbook {
	return &Orderbook{
		Ask: a,
		Bid: b,
	}
}

// func (bo *Orderbook) CreateBestOrder(ask *Ask, bid *Bid) {
// 	// validate ask and bid
// 	bo.BestAsk = ask
// 	bo.BestBid = bid
// }

// // Ref to Samvcodes
// type Service struct{}

// func NewService() *Service {
// 	return &Service{}
// }

// // req is requestInput
// func (s *Service) CreateBestOrder(req CreateBestOrderInput, order *types.Order) (*types.Profile, error) {
// 	// map of validations
// 	// var validationErrs errsx.Map
// 	bestAsk, err := types.NewBestAsk(req.BestAsk)
// 	if err != nil {
// 		validationErrs.Set("bestAsk", err)
// 	}
// 	bestBid, err := types.NewBestBisd(req.BestBid)
// 	if err != nil {
// 		validationErrs.Set("bestBid", err)
// 	}

// 	if validationErrs != nil {
// 		return nil, fmt.Errorf("%w: validationErrs", services.ErrInvalidInput)
// 	}

// 	// return &types.BestOrder{
// 	return &BestOrder{
// 		BestAsk:   bestAsk,
// 		BestBid:   bestBid,
// 		CreatedAt: time.Now().UTC(),
// 	}, nil

// }

// where we create the BestOrder
// func fetchsomething() {
// 	//get the orders
// 	var createProfileReq profile.CreateProfileInput
// 	var createBestOrderReq bestOrder.CreateBestOrderInput
// }

// Other methods
//func (s *Service) DeleteOrders(req CreateBestOrderInput, order *types.Order) (*types.Profile, error)
