package main

import (
	"fmt"

	"zerodhaapi-go/ticker"
)

func main() {
	client := ticker.NewClient("api_key", "access_token")

	fmt.Println(client.Endpoint())
	fmt.Println(string(ticker.Subscribe(408065)))
	fmt.Println(string(ticker.SetMode(ticker.ModeFull, 408065)))
}
