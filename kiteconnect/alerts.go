package kiteconnect

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type AlertParams struct {
	Name             string
	Type             string
	LHSExchange      string
	LHSTradingSymbol string
	LHSAttribute     string
	Operator         string
	RHSType          string
	RHSConstant      float64
	RHSExchange      string
	RHSTradingSymbol string
	RHSAttribute     string
	Basket           string
}

type Alert struct {
	UUID             string         `json:"uuid"`
	Type             string         `json:"type"`
	UserID           string         `json:"user_id"`
	Name             string         `json:"name"`
	Status           string         `json:"status"`
	DisabledReason   string         `json:"disabled_reason"`
	LHSAttribute     string         `json:"lhs_attribute"`
	LHSExchange      string         `json:"lhs_exchange"`
	LHSTradingSymbol string         `json:"lhs_tradingsymbol"`
	Operator         string         `json:"operator"`
	RHSType          string         `json:"rhs_type"`
	RHSAttribute     string         `json:"rhs_attribute"`
	RHSExchange      string         `json:"rhs_exchange"`
	RHSTradingSymbol string         `json:"rhs_tradingsymbol"`
	RHSConstant      float64        `json:"rhs_constant"`
	AlertCount       int            `json:"alert_count"`
	Basket           map[string]any `json:"basket"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

func (c *Client) CreateAlert(ctx context.Context, p AlertParams) (*Alert, error) {
	var out Alert
	if err := c.do(ctx, http.MethodPost, "/alerts", nil, alertValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ModifyAlert(ctx context.Context, uuid string, p AlertParams) (*Alert, error) {
	var out Alert
	if err := c.do(ctx, http.MethodPut, "/alerts/"+uuid, nil, alertValues(p), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Alerts(ctx context.Context, status string, page, pageSize int) ([]Alert, error) {
	q := url.Values{}
	setString(q, "status", status)
	setInt(q, "page", page)
	setInt(q, "page_size", pageSize)
	var out []Alert
	if err := c.do(ctx, http.MethodGet, "/alerts", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) Alert(ctx context.Context, uuid string) (*Alert, error) {
	var out Alert
	if err := c.do(ctx, http.MethodGet, "/alerts/"+uuid, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteAlerts(ctx context.Context, uuids ...string) error {
	q := url.Values{}
	for _, uuid := range uuids {
		q.Add("uuid", uuid)
	}
	return c.do(ctx, http.MethodDelete, "/alerts", q, nil, nil)
}

func (c *Client) AlertHistory(ctx context.Context, uuid string) ([]map[string]any, error) {
	var out []map[string]any
	if err := c.do(ctx, http.MethodGet, "/alerts/"+uuid+"/history", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func alertValues(p AlertParams) url.Values {
	v := url.Values{}
	setString(v, "name", p.Name)
	setString(v, "type", p.Type)
	setString(v, "lhs_exchange", p.LHSExchange)
	setString(v, "lhs_tradingsymbol", p.LHSTradingSymbol)
	setString(v, "lhs_attribute", p.LHSAttribute)
	setString(v, "operator", p.Operator)
	setString(v, "rhs_type", p.RHSType)
	if p.RHSConstant != 0 {
		v.Set("rhs_constant", strconv.FormatFloat(p.RHSConstant, 'f', -1, 64))
	}
	setString(v, "rhs_exchange", p.RHSExchange)
	setString(v, "rhs_tradingsymbol", p.RHSTradingSymbol)
	setString(v, "rhs_attribute", p.RHSAttribute)
	setString(v, "basket", p.Basket)
	return v
}
