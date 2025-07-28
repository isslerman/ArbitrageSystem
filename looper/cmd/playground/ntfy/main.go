package main

import (
	"fmt"
	"grpc-server/infra/ntfy"
)

func main() {
	ntfy := ntfy.NewNtfy()

	isUrgent := true
	title := "This is the Title"
	msg := "This is the msg"

	res := ntfy.SendMsg(title, msg, isUrgent)
	if res {
		fmt.Printf("Mensage sent. %v", res)
	} else {
		fmt.Println("Mensage error, not sent.")
	}
}
