package kiteconnect

import (
	"context"
	"net/http"
)

type MFOrder struct {
	OrderID           string  `json:"order_id"`
	ExchangeOrderID   string  `json:"exchange_order_id"`
	TradingSymbol     string  `json:"tradingsymbol"`
	Status            string  `json:"status"`
	StatusMessage     string  `json:"status_message"`
	Folio             string  `json:"folio"`
	Fund              string  `json:"fund"`
	OrderTimestamp    string  `json:"order_timestamp"`
	ExchangeTimestamp string  `json:"exchange_timestamp"`
	SettlementID      string  `json:"settlement_id"`
	TransactionType   string  `json:"transaction_type"`
	Amount            float64 `json:"amount"`
	Variety           string  `json:"variety"`
	PurchaseType      string  `json:"purchase_type"`
	Quantity          float64 `json:"quantity"`
	LastPrice         float64 `json:"last_price"`
	AveragePrice      float64 `json:"average_price"`
	PlacedBy          string  `json:"placed_by"`
	LastPriceDate     string  `json:"last_price_date"`
	Tag               string  `json:"tag"`
}

type MFSIP struct {
	SIPID                string         `json:"sip_id"`
	TradingSymbol        string         `json:"tradingsymbol"`
	Fund                 string         `json:"fund"`
	DividendType         string         `json:"dividend_type"`
	TransactionType      string         `json:"transaction_type"`
	Status               string         `json:"status"`
	Created              string         `json:"created"`
	Frequency            string         `json:"frequency"`
	NextInstalment       string         `json:"next_instalment"`
	InstalmentAmount     float64        `json:"instalment_amount"`
	Instalments          int            `json:"instalments"`
	LastInstalment       string         `json:"last_instalment"`
	PendingInstalments   int            `json:"pending_instalments"`
	InstalmentDay        int            `json:"instalment_day"`
	CompletedInstalments int            `json:"completed_instalments"`
	Tag                  string         `json:"tag"`
	StepUp               map[string]any `json:"step_up"`
}

type MFHolding struct {
	Folio           string  `json:"folio"`
	Fund            string  `json:"fund"`
	TradingSymbol   string  `json:"tradingsymbol"`
	AveragePrice    float64 `json:"average_price"`
	LastPrice       float64 `json:"last_price"`
	LastPriceDate   string  `json:"last_price_date"`
	PledgedQuantity float64 `json:"pledged_quantity"`
	PNL             float64 `json:"pnl"`
	Quantity        float64 `json:"quantity"`
}

func (c *Client) MFOrders(ctx context.Context) ([]MFOrder, error) {
	var out []MFOrder
	if err := c.do(ctx, http.MethodGet, "/mf/orders", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MFOrder(ctx context.Context, orderID string) (*MFOrder, error) {
	var out MFOrder
	if err := c.do(ctx, http.MethodGet, "/mf/orders/"+orderID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) MFSIPs(ctx context.Context) ([]MFSIP, error) {
	var out []MFSIP
	if err := c.do(ctx, http.MethodGet, "/mf/sips", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MFHoldings(ctx context.Context) ([]MFHolding, error) {
	var out []MFHolding
	if err := c.do(ctx, http.MethodGet, "/mf/holdings", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) MFInstruments(ctx context.Context) ([]byte, error) {
	body, _, err := c.raw(ctx, http.MethodGet, "/mf/instruments", nil)
	return body, err
}
