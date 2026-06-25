package main

import (
	"context"
	"fmt"
	"log"

	"zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("api_key", "api_secret")
	client.SetAccessToken("access_token")

	order, err := client.PlaceOrder(context.Background(), kiteconnect.PlaceOrderParams{
		Exchange:        kiteconnect.ExchangeNSE,
		TradingSymbol:   "INFY",
		TransactionType: kiteconnect.TransactionBuy,
		Quantity:        1,
		Product:         kiteconnect.ProductMIS,
		OrderType:       kiteconnect.OrderTypeMarket,
		Validity:        kiteconnect.ValidityDay,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(order.OrderID)
}
