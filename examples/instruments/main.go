package main

import (
	"context"
	"fmt"
	"log"

	"github.com/suyotech/zerodhaapi-go/instruments"
	"github.com/suyotech/zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("api_key", "api_secret")
	client.SetAccessToken("access_token")

	if err := instruments.CheckDownload(context.Background(), client, "data/instruments.csv"); err != nil {
		log.Fatal(err)
	}
	allInstruments, err := instruments.Load()
	if err != nil {
		log.Fatal(err)
	}

	name := "NIFTY"
	exchange := "NFO"

	expiries, err := instruments.ExpiryDates(instruments.FindRequest{
		Name:     &name,
		Exchange: &exchange,
	}, allInstruments)
	if err != nil {
		log.Fatal(err)
	}
	if len(expiries) == 0 {
		log.Fatal("no expiries found")
	}

	strikes, err := instruments.OptionStrikes(instruments.FindRequest{
		Name:     &name,
		Exchange: &exchange,
		Expiry:   &expiries[0],
	}, allInstruments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(expiries[0].Format("2006-01-02"), strikes)
}
