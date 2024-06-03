package exchange

import (
	"strings"
	"testing"
)

func Test_name(t *testing.T) {
	e := NewMercadoBitcoin()

	if len(strings.TrimSpace(e.Name)) == 0 {
		t.Error("missing exchange name")
	}

	if len(strings.TrimSpace(e.Id)) == 0 {
		t.Error("missing exchange id")
	}

	if len(strings.TrimSpace(e.Id)) == 0 {
		t.Error("missing exchange apiURL")
	}
}

func Test_bestAsk(t *testing.T) {
	e := NewMercadoBitcoin()

	ask, bid, err := e.BestOrder()
	if err != nil {
		t.Errorf("Error getting bestAsk. Got error %s", err)
	}

	if ask.Price < 0.0 || bid.Price < 0.0 {
		t.Error("expected bestAsk to be zero or positive number.")
	}
}

func Test_fetchApiData(t *testing.T) {
	e := NewMercadoBitcoin()

	d, err := e.fetchApiData()
	if err != nil {
		t.Errorf("Error fetchApiData from %s, %v", e.Id, err)
	}

	if len(d.Asks) == 0 {
		t.Errorf("No asks found in %s", e.Id)
	}
}
