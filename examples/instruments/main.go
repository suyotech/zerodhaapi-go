package main

import (
	"context"
	"fmt"
	"log"

	"zerodhaapi-go/instruments"
	"zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("api_key", "api_secret")
	client.SetAccessToken("access_token")

	store, err := instruments.LoadOrDownload(context.Background(), instruments.CacheConfig{
		Client:   client,
		FilePath: "data/instruments.csv",
	})
	if err != nil {
		log.Fatal(err)
	}

	name := "NIFTY"
	exchange := "NFO"

	expiries := instruments.ExpiryDates(store, instruments.FindRequest{
		Name:     &name,
		Exchange: &exchange,
	})
	if len(expiries) == 0 {
		log.Fatal("no expiries found")
	}

	strikes := instruments.OptionStrikes(store, instruments.FindRequest{
		Name:     &name,
		Exchange: &exchange,
		Expiry:   &expiries[0],
	})

	fmt.Println(expiries[0].Format("2006-01-02"), strikes)
}
