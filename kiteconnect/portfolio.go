package kiteconnect

import (
	"context"
	"net/http"
)

type Holding struct {
	TradingSymbol       string  `json:"tradingsymbol"`
	Exchange            string  `json:"exchange"`
	InstrumentToken     uint32  `json:"instrument_token"`
	ISIN                string  `json:"isin"`
	Product             string  `json:"product"`
	Price               float64 `json:"price"`
	Quantity            int     `json:"quantity"`
	UsedQuantity        int     `json:"used_quantity"`
	T1Quantity          int     `json:"t1_quantity"`
	AveragePrice        float64 `json:"average_price"`
	LastPrice           float64 `json:"last_price"`
	ClosePrice          float64 `json:"close_price"`
	PNL                 float64 `json:"pnl"`
	DayChange           float64 `json:"day_change"`
	DayChangePercentage float64 `json:"day_change_percentage"`
}

type Positions struct {
	Net []Position `json:"net"`
	Day []Position `json:"day"`
}

type Position struct {
	TradingSymbol     string  `json:"tradingsymbol"`
	Exchange          string  `json:"exchange"`
	InstrumentToken   uint32  `json:"instrument_token"`
	Product           string  `json:"product"`
	Quantity          int     `json:"quantity"`
	OvernightQuantity int     `json:"overnight_quantity"`
	AveragePrice      float64 `json:"average_price"`
	ClosePrice        float64 `json:"close_price"`
	LastPrice         float64 `json:"last_price"`
	Value             float64 `json:"value"`
	PNL               float64 `json:"pnl"`
	M2M               float64 `json:"m2m"`
	Unrealised        float64 `json:"unrealised"`
	Realised          float64 `json:"realised"`
	BuyQuantity       int     `json:"buy_quantity"`
	BuyPrice          float64 `json:"buy_price"`
	BuyValue          float64 `json:"buy_value"`
	SellQuantity      int     `json:"sell_quantity"`
	SellPrice         float64 `json:"sell_price"`
	SellValue         float64 `json:"sell_value"`
}

func (c *Client) Holdings(ctx context.Context) ([]Holding, error) {
	var out []Holding
	if err := c.do(ctx, http.MethodGet, "/portfolio/holdings", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Positions(ctx context.Context) (*Positions, error) {
	var out Positions
	if err := c.do(ctx, http.MethodGet, "/portfolio/positions", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
