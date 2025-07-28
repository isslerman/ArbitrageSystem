package trading

import "fmt"

// trading package

// manage open openOrders maker side
// create order
// list all open openOrders: ids
// create open openOrders
type Order struct {
	ID string
}

// func main() {
// 	tm := createTradingManager()

// 	// Create openOrders
// 	if err := tm.CreateOrder("4237BB6F-CE0A-4E6A-86CB-99AF0FF9D341"); err != nil {
// 		fmt.Println(err)
// 	}
// 	if err := tm.CreateOrder("86CB"); err != nil {
// 		fmt.Println(err)
// 	}
// 	if err := tm.CreateOrder("4237BB6F-CE0A-4E6A-86CB-99AF0FF9D343"); err != nil {
// 		fmt.Println(err)
// 	}

// 	// List openOrders
// 	orderList := tm.ListOrders()
// 	for _, order := range orderList {
// 		fmt.Printf("ID: %s\n", order.ID)
// 	}

// 	// Remove an order
// 	err := tm.RemoveOrder("86CB")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// List again to see if the removed order is gone
// 	orderList = tm.ListOrders()
// 	for _, order := range orderList {
// 		fmt.Printf("ID: %s\n", order.ID)
// 	}
// }

func NewTradingManager() *TradingManager {
	return &TradingManager{
		openOrders: make(map[string]Order),
	}
}

// TradingManager is a struct that holds a map of openOrders.
// It provides methods to create, list and remove openOrders.
type TradingManager struct {
	openOrders map[string]Order
}

func (m *TradingManager) CreateOrder(id string) error {
	if _, exists := m.openOrders[id]; exists {
		return fmt.Errorf("order with ID %s already exists", id)
	}
	m.openOrders[id] = Order{ID: id}
	return nil
}

func (m *TradingManager) ListOrders() []Order {
	orderList := make([]Order, 0, len(m.openOrders))
	for _, order := range m.openOrders {
		orderList = append(orderList, order)
	}
	return orderList
}

func (m *TradingManager) RemoveOrder(id string) error {
	if _, exists := m.openOrders[id]; !exists {
		return fmt.Errorf("order with ID %s does not exist", id)
	}
	delete(m.openOrders, id)
	return nil
}
