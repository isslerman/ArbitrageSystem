package bina

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"grpc-server/internal/cex/data"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
)

type Bina struct {
	apiBaseURL string
	id         string
	name       string
	FeeTaker   float64
	FeeMaker   float64
	Limits     data.OrderLimits
	client     *binance.Client
}

func New() *Bina {
	limits := &data.OrderLimits{
		OrdMinAmount: 10,
	}

	// Get the current working directory
	// abs, _ := os.Getwd()
	// temp var for tests run
	abs := "/Users/marcosissler/projects/202404-ArbitrageSystem/grpc-server"
	envFile := fmt.Sprintf("%s/.env", abs)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("APIKEY_BINA")
	secretKey := os.Getenv("APISEC_BINA")

	client := binance.NewClient(apiKey, secretKey)

	// Return Ripi with default values
	return &Bina{
		apiBaseURL: "https://api.binance.com/api/v3/",
		id:         "BINA",
		name:       "Binance",
		FeeTaker:   0.0050,
		FeeMaker:   0.0025,
		Limits:     *limits,
		client:     client,
	}
}

func (e *Bina) Balance(asset string) (amount float64, err error) {
	return 0.0, nil
}

func (e *Bina) Id() string {
	return e.id
}

// CancelAllOrders -
func (e *Bina) CancelAllOrders() error {
	openOrders, err := e.client.NewListOpenOrdersService().Symbol("SOLBRL").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}
	ids := make([]int64, len(openOrders))
	errors := make([]error, len(ids))
	for i, o := range openOrders {
		ids = append(ids, o.OrderID)
		slog.Info(fmt.Sprintf("OrderID: %d", o.OrderID))
		slog.Info(fmt.Sprintf("openOrders [%d]: %v", i, o))

		// cancel
		res, err := e.client.NewCancelOrderService().Symbol("SOLBRL").
			OrderID(o.OrderID).
			Do(context.Background())
		if err != nil {
			errors = append(errors, err)
		}
		slog.Info(fmt.Sprintf("CancelOrderResponse [%d]: %v", i, res))
	}

	return combineErrors(errors)
}

// OrdersCreate -
func (e *Bina) OrdersCreate(o *data.OrdersCreateRequest) error {

	var side binance.SideType
	if o.Side == "sell" {
		side = binance.SideTypeSell
	} else {
		if o.Side == "buy" {
			side = binance.SideTypeBuy
		} else {
			log.Fatal("Invalid side order")
		}
	}

	qty := fmt.Sprint(o.Amount)
	price := fmt.Sprint(o.Price)

	order, err := e.client.NewCreateOrderService().
		Symbol(o.Pair).
		Side(side).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(qty).
		Price(price).
		Do(context.Background())
	if err != nil {
		log.Fatal("error creating order", err)
		return err
	}
	slog.Info(fmt.Sprintf("OrdersCreate: %v", order))

	// Use Test() instead of Do() for testing.
	return nil
}

func combineErrors(errs []error) error {
	var sb strings.Builder
	for _, err := range errs {
		if err != nil {
			sb.WriteString(err.Error() + "; ")
		}
	}
	if sb.Len() == 0 {
		return nil
	}
	return errors.New(sb.String())
}
