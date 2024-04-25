// Reference to ArdanLabs project template
// https://github.com/ardanlabs/service/blob/master/business/domain/userbus/model.go
package bestOrder

import (
	"fmt"
	"grpc-client/internal/services"
	"time"
)

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

func NewBestOrder() *BestOrder {
	return &BestOrder{}
}

func (bo *BestOrder) CreateBestOrder(ask *BestAsk, bid *BestBid) {
	// validate ask and bid
	bo.BestAsk = ask
	bo.BestBid = bid
}

// Ref to Samvcodes
type Service struct{}

func NewService() *Service {
	return &Service{}
}

// req is requestInput
func (s *Service) CreateBestOrder(req CreateBestOrderInput, order *types.Order) (*types.Profile, error) {
	// map of validations
	// var validationErrs errsx.Map
	bestAsk, err := types.NewBestAsk(req.BestAsk)
	if err != nil {
		validationErrs.Set("bestAsk", err)
	}
	bestBid, err := types.NewBestBisd(req.BestBid)
	if err != nil {
		validationErrs.Set("bestBid", err)
	}

	if validationErrs != nil {
		return nil, fmt.Errorf("%w: validationErrs", services.ErrInvalidInput)
	}

	// return &types.BestOrder{
	return &BestOrder{
		BestAsk:   bestAsk,
		BestBid:   bestBid,
		CreatedAt: time.Now().UTC(),
	}, nil

}

// where we create the BestOrder
func fetchsomething(){
	//get the orders
	var createProfileReq profile.CreateProfileInput
	var createBestOrderReq bestOrder.CreateBestOrderInput
}

// Other methods
//func (s *Service) DeleteOrders(req CreateBestOrderInput, order *types.Order) (*types.Profile, error)