package main

import (
	"fmt"
	"grpc-server/internal/cex"
)

func main() {
	// Example usage
	ripi := cex.NewRIPI()
	response, err := ripi.CreateOrder()
	if err != nil {
		fmt.Println("Error creating order:", err)
		return
	}
	fmt.Println("Response:", response)
}
