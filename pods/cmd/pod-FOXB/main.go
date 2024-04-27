package main

// Just a gRPC client that connects with our server and send a msg.

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	// app := Config{
	// 	Host: "localhost",
	// 	Port: 50001,
	// }

	// e := exchanges.NewFoxbit()

}
