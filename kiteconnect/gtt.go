package kiteconnect

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GTTTypeSingle = "single"
	GTTTypeOCO    = "two-leg"
)

type GTTParams struct {
	Type          string
	Exchange      string
	TradingSymbol string
	TriggerValues []float64
	LastPrice     float64
	Orders        []GTTOrder
}

type GTTOrder struct {
	TransactionType string  `json:"transaction_type"`
	Quantity        int     `json:"quantity"`
	OrderType       string  `json:"order_type"`
	Product         string  `json:"product"`
	Price           float64 `json:"price"`
}

type GTT struct {
	ID            int64      `json:"id"`
	UserID        string     `json:"user_id"`
	ParentTrigger string     `json:"parent_trigger"`
	Type          string     `json:"type"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	ExpiresAt     string     `json:"expires_at"`
	Status        string     `json:"status"`
	Condition     any        `json:"condition"`
	Orders        []GTTOrder `json:"orders"`
}

type GTTResponse struct {
	TriggerID int64 `json:"trigger_id"`
}

func (c *Client) PlaceGTT(ctx context.Context, p GTTParams) (*GTTResponse, error) {
	var out GTTResponse
	if err := c.do(ctx, http.MethodPost, "/gtt/triggers", nil, gttValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyGTT(ctx context.Context, triggerID int64, p GTTParams) (*GTTResponse, error) {
	var out GTTResponse
	if err := c.do(ctx, http.MethodPut, "/gtt/triggers/"+strconv.FormatInt(triggerID, 10), nil, gttValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteGTT(ctx context.Context, triggerID int64) (*GTTResponse, error) {
	var out GTTResponse
	if err := c.do(ctx, http.MethodDelete, "/gtt/triggers/"+strconv.FormatInt(triggerID, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GTTs(ctx context.Context) ([]GTT, error) {
	var out []GTT
	if err := c.do(ctx, http.MethodGet, "/gtt/triggers", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GTT(ctx context.Context, triggerID int64) (*GTT, error) {
	var out GTT
	if err := c.do(ctx, http.MethodGet, "/gtt/triggers/"+strconv.FormatInt(triggerID, 10), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func gttValues(p GTTParams) url.Values {
	condition, _ := json.Marshal(map[string]any{
		"exchange":       p.Exchange,
		"tradingsymbol":  p.TradingSymbol,
		"trigger_values": p.TriggerValues,
		"last_price":     p.LastPrice,
	})
	orders, _ := json.Marshal(p.Orders)
	v := url.Values{}
	setString(v, "type", p.Type)
	v.Set("condition", string(condition))
	v.Set("orders", string(orders))
	return v
}
