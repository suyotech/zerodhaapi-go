# zerodhaapi-go

Small Go SDK for Zerodha Kite Connect v3.

Packages:

- `kiteconnect`: REST API client.
- `ticker`: WebSocket ticker client and binary tick parser.
- `instruments`: instrument CSV cache and lookup helpers.

## Structure

```text
kiteconnect/    REST API client
ticker/         WebSocket ticker client and tick parser
instruments/    instrument dump cache and lookup
examples/       small runnable examples
```

## Login

```go
client := kiteconnect.NewClient("api_key", "api_secret")
fmt.Println(client.LoginURL())

session, err := client.GenerateSession(ctx, "request_token")
if err != nil {
    log.Fatal(err)
}
client.SetAccessToken(session.AccessToken)
```

## Trading

```go
order, err := client.PlaceOrder(ctx, kiteconnect.PlaceOrderParams{
    Exchange:        kiteconnect.ExchangeNSE,
    TradingSymbol:   "INFY",
    TransactionType: kiteconnect.TransactionBuy,
    Quantity:        1,
    Product:         kiteconnect.ProductMIS,
    OrderType:       kiteconnect.OrderTypeMarket,
    Validity:        kiteconnect.ValidityDay,
})
```

## Market Data

```go
ltp, err := client.LTP(ctx, "NSE:INFY")
candles, err := client.HistoricalData(ctx, kiteconnect.HistoricalDataParams{
    InstrumentToken: 408065,
    Interval:        kiteconnect.IntervalDay,
    From:            time.Now().AddDate(0, -1, 0),
    To:              time.Now(),
})
```

## Ticker

```go
t := ticker.NewClient("api_key", session.AccessToken)
t.OnTick(func(ticks []ticker.Tick) {
    for _, tick := range ticks {
        fmt.Println(tick.InstrumentToken, tick.LastPrice)
    }
})
t.Subscribe(ticker.ModeFull, 408065)
err := t.Connect(ctx)
```

Available ticker modes are `ticker.ModeLTP`, `ticker.ModeQuote`, and `ticker.ModeFull`.

## Instruments

```go
store, err := instruments.LoadOrDownload(ctx, instruments.CacheConfig{
    Client:   client,
    FilePath: "data/instruments.csv",
})

name := "NIFTY"
exchange := "NFO"
expiries := instruments.ExpiryDates(store, instruments.FindRequest{
    Name:     &name,
    Exchange: &exchange,
})
strikes := instruments.OptionStrikes(store, instruments.FindRequest{
    Name:     &name,
    Exchange: &exchange,
    Expiry:   &expiries[0],
})
```
