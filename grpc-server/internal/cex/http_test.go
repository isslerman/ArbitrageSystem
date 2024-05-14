package cex

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
	}

	for _, e := range tests {
		ro := RequestOptions{}
		hs := NewHTTPService()
		resp, _ := hs.NewRequest(e.url, e.method, ro)

		body, _ := io.ReadAll(resp.Body)
		spew.Dump(body)
		fmt.Println(body)

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("GET returned wrong status code; expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}
}
