package main

import (
	"context"
	"fmt"
	"log"

	"github.com/suyotech/zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("api_key", "api_secret")
	client.SetAccessToken("access_token")

	positions, err := client.Positions(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(positions.Net), len(positions.Day))
}
