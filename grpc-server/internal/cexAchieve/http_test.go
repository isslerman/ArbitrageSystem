package cex

import (
	"net/http"
	"testing"
)

func TestHTTP_NewRequest(t *testing.T) {
	var tests = []struct {
		name               string
		url                string
		method             string
		expectedStatusCode int
	}{
		{"Get BITY Ticker", "https://api.bitpreco.com/btc-brl/ticker", http.MethodGet, http.StatusOK},
		{"Get BITY Private Unauth", "https://api.bitpreco.com/v1/trading/balance", http.MethodPost, http.StatusOK},
		{"Getting JSON from the URL", "https://reqbin.com", http.MethodGet, http.StatusOK},
	}

	for _, e := range tests {
		ro := RequestOptions{}
		hs := NewHTTPService()
		resp, _ := hs.NewRequest(e.url, e.method, ro)

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("GET returned wrong status code; expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}
}
