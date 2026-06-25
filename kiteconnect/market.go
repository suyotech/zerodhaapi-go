package kiteconnect

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Quote struct {
	InstrumentToken   uint32      `json:"instrument_token"`
	Timestamp         string      `json:"timestamp"`
	LastTradeTime     string      `json:"last_trade_time"`
	LastPrice         float64     `json:"last_price"`
	LastQuantity      int         `json:"last_quantity"`
	BuyQuantity       int         `json:"buy_quantity"`
	SellQuantity      int         `json:"sell_quantity"`
	Volume            int         `json:"volume"`
	AveragePrice      float64     `json:"average_price"`
	OI                int         `json:"oi"`
	OIDayHigh         int         `json:"oi_day_high"`
	OIDayLow          int         `json:"oi_day_low"`
	NetChange         float64     `json:"net_change"`
	LowerCircuitLimit float64     `json:"lower_circuit_limit"`
	UpperCircuitLimit float64     `json:"upper_circuit_limit"`
	OHLC              OHLC        `json:"ohlc"`
	Depth             MarketDepth `json:"depth"`
}

type OHLC struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

type MarketDepth struct {
	Buy  []DepthItem `json:"buy"`
	Sell []DepthItem `json:"sell"`
}

type DepthItem struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Orders   int     `json:"orders"`
}

type LTP struct {
	InstrumentToken uint32  `json:"instrument_token"`
	LastPrice       float64 `json:"last_price"`
}

func (c *Client) Quote(ctx context.Context, instruments ...string) (map[string]Quote, error) {
	var out map[string]Quote
	if err := c.do(ctx, http.MethodGet, "/quote", instrumentQuery(instruments), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) OHLC(ctx context.Context, instruments ...string) (map[string]Quote, error) {
	var out map[string]Quote
	if err := c.do(ctx, http.MethodGet, "/quote/ohlc", instrumentQuery(instruments), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) LTP(ctx context.Context, instruments ...string) (map[string]LTP, error) {
	var out map[string]LTP
	if err := c.do(ctx, http.MethodGet, "/quote/ltp", instrumentQuery(instruments), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Instruments(ctx context.Context, exchange string) ([]byte, error) {
	path := "/instruments"
	if exchange != "" {
		path += "/" + exchange
	}
	body, header, err := c.raw(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(header.Get("Content-Encoding"), "gzip") && !(len(body) >= 2 && body[0] == 0x1f && body[1] == 0x8b) {
		return body, nil
	}
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func instrumentQuery(instruments []string) url.Values {
	q := url.Values{}
	for _, instrument := range instruments {
		q.Add("i", instrument)
	}
	return q
}
