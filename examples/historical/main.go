package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("api_key", "api_secret")
	client.SetAccessToken("access_token")

	candles, err := client.HistoricalData(context.Background(), kiteconnect.HistoricalDataParams{
		InstrumentToken: 408065,
		Interval:        kiteconnect.IntervalDay,
		From:            time.Now().AddDate(0, -1, 0),
		To:              time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(candles))
}
