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

	ltp, err := client.LTP(context.Background(), "NSE:INFY")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ltp["NSE:INFY"].LastPrice)
}
