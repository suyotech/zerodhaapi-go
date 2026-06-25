package kiteconnect

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type PlaceOrderParams struct {
	Variety           string
	Exchange          string
	TradingSymbol     string
	TransactionType   string
	Quantity          int
	Product           string
	OrderType         string
	Price             float64
	TriggerPrice      float64
	DisclosedQuantity int
	Validity          string
	ValidityTTL       int
	IcebergLegs       int
	IcebergQuantity   int
	AuctionNumber     string
	MarketProtection  int
	AutoSlice         bool
	Tag               string
}

type ModifyOrderParams struct {
	Quantity          int
	Price             float64
	OrderType         string
	TriggerPrice      float64
	DisclosedQuantity int
	Validity          string
	ValidityTTL       int
	MarketProtection  int
}

type OrderResponse struct {
	OrderID string `json:"order_id"`
}

type Order struct {
	OrderID                 string         `json:"order_id"`
	ExchangeOrderID         string         `json:"exchange_order_id"`
	ParentOrderID           string         `json:"parent_order_id"`
	Status                  string         `json:"status"`
	StatusMessage           string         `json:"status_message"`
	OrderTimestamp          string         `json:"order_timestamp"`
	ExchangeUpdateTimestamp string         `json:"exchange_update_timestamp"`
	ExchangeTimestamp       string         `json:"exchange_timestamp"`
	Variety                 string         `json:"variety"`
	Exchange                string         `json:"exchange"`
	TradingSymbol           string         `json:"tradingsymbol"`
	InstrumentToken         uint32         `json:"instrument_token"`
	OrderType               string         `json:"order_type"`
	TransactionType         string         `json:"transaction_type"`
	Validity                string         `json:"validity"`
	Product                 string         `json:"product"`
	Quantity                int            `json:"quantity"`
	DisclosedQuantity       int            `json:"disclosed_quantity"`
	Price                   float64        `json:"price"`
	TriggerPrice            float64        `json:"trigger_price"`
	AveragePrice            float64        `json:"average_price"`
	FilledQuantity          int            `json:"filled_quantity"`
	PendingQuantity         int            `json:"pending_quantity"`
	CancelledQuantity       int            `json:"cancelled_quantity"`
	MarketProtection        int            `json:"market_protection"`
	Meta                    map[string]any `json:"meta"`
	Tag                     string         `json:"tag"`
	GUID                    string         `json:"guid"`
}

type Trade struct {
	TradeID           string  `json:"trade_id"`
	OrderID           string  `json:"order_id"`
	Exchange          string  `json:"exchange"`
	TradingSymbol     string  `json:"tradingsymbol"`
	InstrumentToken   uint32  `json:"instrument_token"`
	Product           string  `json:"product"`
	AveragePrice      float64 `json:"average_price"`
	Quantity          int     `json:"quantity"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	TransactionType   string  `json:"transaction_type"`
	FillTimestamp     string  `json:"fill_timestamp"`
	OrderTimestamp    string  `json:"order_timestamp"`
	ExchangeTimestamp string  `json:"exchange_timestamp"`
}

func (c *Client) PlaceOrder(ctx context.Context, p PlaceOrderParams) (*OrderResponse, error) {
	variety := p.Variety
	if variety == "" {
		variety = VarietyRegular
	}
	var out OrderResponse
	if err := c.do(ctx, http.MethodPost, "/orders/"+variety, nil, placeOrderValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyOrder(ctx context.Context, variety, orderID string, p ModifyOrderParams) (*OrderResponse, error) {
	if variety == "" {
		variety = VarietyRegular
	}
	var out OrderResponse
	if err := c.do(ctx, http.MethodPut, "/orders/"+variety+"/"+orderID, nil, modifyOrderValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CancelOrder(ctx context.Context, variety, orderID string) (*OrderResponse, error) {
	if variety == "" {
		variety = VarietyRegular
	}
	var out OrderResponse
	if err := c.do(ctx, http.MethodDelete, "/orders/"+variety+"/"+orderID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Orders(ctx context.Context) ([]Order, error) {
	var out []Order
	if err := c.do(ctx, http.MethodGet, "/orders", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) OrderHistory(ctx context.Context, orderID string) ([]Order, error) {
	var out []Order
	if err := c.do(ctx, http.MethodGet, "/orders/"+orderID, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Trades(ctx context.Context) ([]Trade, error) {
	var out []Trade
	if err := c.do(ctx, http.MethodGet, "/trades", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) OrderTrades(ctx context.Context, orderID string) ([]Trade, error) {
	var out []Trade
	if err := c.do(ctx, http.MethodGet, "/orders/"+orderID+"/trades", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func placeOrderValues(p PlaceOrderParams) url.Values {
	v := url.Values{}
	setString(v, "exchange", p.Exchange)
	setString(v, "tradingsymbol", p.TradingSymbol)
	setString(v, "transaction_type", p.TransactionType)
	setInt(v, "quantity", p.Quantity)
	setString(v, "product", p.Product)
	setString(v, "order_type", p.OrderType)
	setFloat(v, "price", p.Price)
	setFloat(v, "trigger_price", p.TriggerPrice)
	setInt(v, "disclosed_quantity", p.DisclosedQuantity)
	setString(v, "validity", p.Validity)
	setInt(v, "validity_ttl", p.ValidityTTL)
	setInt(v, "iceberg_legs", p.IcebergLegs)
	setInt(v, "iceberg_quantity", p.IcebergQuantity)
	setString(v, "auction_number", p.AuctionNumber)
	setInt(v, "market_protection", p.MarketProtection)
	if p.AutoSlice {
		v.Set("autoslice", "true")
	}
	setString(v, "tag", p.Tag)
	return v
}

func modifyOrderValues(p ModifyOrderParams) url.Values {
	v := url.Values{}
	setInt(v, "quantity", p.Quantity)
	setFloat(v, "price", p.Price)
	setString(v, "order_type", p.OrderType)
	setFloat(v, "trigger_price", p.TriggerPrice)
	setInt(v, "disclosed_quantity", p.DisclosedQuantity)
	setString(v, "validity", p.Validity)
	setInt(v, "validity_ttl", p.ValidityTTL)
	setInt(v, "market_protection", p.MarketProtection)
	return v
}

func setString(v url.Values, key, value string) {
	if value != "" {
		v.Set(key, value)
	}
}

func setInt(v url.Values, key string, value int) {
	if value != 0 {
		v.Set(key, strconv.Itoa(value))
	}
}

func setFloat(v url.Values, key string, value float64) {
	if value != 0 {
		v.Set(key, strconv.FormatFloat(value, 'f', -1, 64))
	}
}
