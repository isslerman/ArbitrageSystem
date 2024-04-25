package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// midPrice := float64(50)

	bCh := make(chan float64, 5)
	aCh := make(chan float64, 5)
	go bidMaker(bCh)
	go askMaker(aCh)

	defer close(bCh)
	defer close(aCh)

	for {
		checkDeal(bCh, aCh)
	}
}

func checkDeal(bCh, aCh chan float64) {
	bidPrice := 0.0
	askPrice := 0.0

	select {
	case bid := <-bCh:
		if bid > bidPrice {
			fmt.Printf("Good! %f is greater than %f\n", bid, bidPrice)
			bidPrice = bid
		}
	case ask := <-aCh:
		askPrice = ask
		if askPrice < bidPrice {
			fmt.Printf("Good! %f is lower than %f\n", askPrice, bidPrice)
		}
	}
}

func bidMaker(bm chan float64) {
	for {
		number := rand.Float64() * 100 // [0,100]
		fmt.Println("Bid:", number)
		bm <- number
		time.Sleep(time.Second * 1)
	}
}

func askMaker(am chan float64) {
	for {
		number := rand.Float64() * 100 // [0,100]
		fmt.Println("Ask:", number)
		am <- number
		time.Sleep(time.Second * 1)
	}
}
