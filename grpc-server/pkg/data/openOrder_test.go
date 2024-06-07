package data

import (
	"testing"
)

func Test_validate(t *testing.T) {
	var orders = []struct {
		name      string
		orderType string
		side      string
		expected  string
	}{
		{name: "valid order", orderType: "market", side: "sell", expected: ""},
		{name: "valid order", orderType: "limit", side: "buy", expected: ""},
		{name: "invalid order type", orderType: "test", side: "sell", expected: "invalid order type"},
	}

	for _, e := range orders {
		_, err := NewAskOpenOrder(10.0, 10.0, "SOLBRL", e.orderType)

		// valid orders
		if err != nil && e.expected == "" {
			t.Errorf("%s: did not expect error, but got one - %s", e.name, err.Error())
		}

		// invalid orders
		if err == nil && e.expected != "" {
			t.Errorf("%s: expected error, but did not get one", e.name)
		}
	}

}
