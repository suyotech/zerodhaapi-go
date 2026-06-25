package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"zerodhaapi-go/ticker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := ticker.NewClient("api_key", "access_token")
	client.OnTick(func(ticks []ticker.Tick) {
		for _, tick := range ticks {
			fmt.Printf("%d %.2f\n", tick.InstrumentToken, tick.LastPrice)
		}
	})
	client.OnError(func(err error) {
		log.Println("ticker error:", err)
	})

	if err := client.Subscribe(ticker.ModeFull, 408065); err != nil {
		log.Fatal(err)
	}
	if err := client.Connect(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
