package kiteconnect

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type MarginOrder struct {
	Exchange        string  `json:"exchange"`                // Example: ExchangeNSE, ExchangeNFO.
	TradingSymbol   string  `json:"tradingsymbol"`           // Example: "INFY".
	TransactionType string  `json:"transaction_type"`        // Example: TransactionBuy, TransactionSell.
	Variety         string  `json:"variety,omitempty"`       // Example: VarietyRegular.
	Product         string  `json:"product"`                 // Example: ProductMIS, ProductNRML.
	OrderType       string  `json:"order_type"`              // Example: OrderTypeMarket, OrderTypeLimit.
	Quantity        int     `json:"quantity"`                // Example: 1.
	Price           float64 `json:"price,omitempty"`         // Example: 1525.50 for limit orders.
	TriggerPrice    float64 `json:"trigger_price,omitempty"` // Example: 1520.00 for SL/SL-M orders.
}

type OrderMargin struct {
	Type          string  `json:"type"`
	TradingSymbol string  `json:"tradingsymbol"`
	Exchange      string  `json:"exchange"`
	SPAN          float64 `json:"span"`
	Exposure      float64 `json:"exposure"`
	OptionPremium float64 `json:"option_premium"`
	Additional    float64 `json:"additional"`
	BO            float64 `json:"bo"`
	Cash          float64 `json:"cash"`
	VAR           float64 `json:"var"`
	Total         float64 `json:"total"`
}

type BasketMargin struct {
	Initial OrderMargin   `json:"initial"`
	Final   OrderMargin   `json:"final"`
	Orders  []OrderMargin `json:"orders"`
	Charges any           `json:"charges"`
}

func (c *Client) OrderMargins(ctx context.Context, orders []MarginOrder) ([]OrderMargin, error) {
	body, _ := json.Marshal(orders)
	form := url.Values{}
	form.Set("orders", string(body))
	var out []OrderMargin
	if err := c.do(ctx, http.MethodPost, "/margins/orders", nil, form, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) BasketMargins(ctx context.Context, orders []MarginOrder, considerPositions bool) (*BasketMargin, error) {
	body, _ := json.Marshal(orders)
	form := url.Values{}
	form.Set("orders", string(body))
	if considerPositions {
		form.Set("consider_positions", "true")
	}
	var out BasketMargin
	if err := c.do(ctx, http.MethodPost, "/margins/basket", nil, form, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
