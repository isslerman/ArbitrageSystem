package cex

// func TestRIPI_CreateOrder(t *testing.T) {

// }

// func TestRIPI_CancelAllOrders(t *testing.T) {
// 	// Define the request body
// 	order := map[string]interface{}{
// 		"symbol":   "SOLBRL",
// 		"side":     "sell",
// 		"type":     "limit",
// 		"price":    9999,
// 		"quantity": 0.01,
// 	}

// 	// Convert order to JSON
// 	orderJSON, err := json.Marshal(order)
// 	if err != nil {
// 		t.Fatalf("failed to marshal order: %v", err)
// 	}

// 	// Create a request with the order JSON as the body
// 	req, err := http.NewRequest("POST", "https://api.binance.com/api/v3/order", bytes.NewBuffer(orderJSON))
// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}

// 	// Set Binance API key and secret if necessary
// 	// req.Header.Set("X-MBX-APIKEY", "YOUR_API_KEY")

// 	// Create a test HTTP server to mock Binance API response
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if request body matches expected order
// 		var receivedOrder map[string]interface{}
// 		err := json.NewDecoder(r.Body).Decode(&receivedOrder)
// 		if err != nil {
// 			http.Error(w, "failed to decode request body", http.StatusBadRequest)
// 			return
// 		}

// 		// Compare received order with expected order
// 		if !compareOrders(order, receivedOrder) {
// 			http.Error(w, "received order does not match expected order", http.StatusBadRequest)
// 			return
// 		}

// 		// Simulate a successful response from Binance API
// 		w.WriteHeader(http.StatusOK)
// 	}))

// 	defer server.Close()

// 	// Send request to the test server
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Fatalf("request failed: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check if the response was successful (status code 200)
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
// 	}
// }

// // compareOrders checks if two orders are identical
// func compareOrders(expected, received map[string]interface{}) bool {
// 	// Check each field of the order
// 	for key, val := range expected {
// 		if receivedVal, ok := received[key]; !ok || receivedVal != val {
// 			return false
// 		}
// 	}
// 	return true
// }
