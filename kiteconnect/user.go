package kiteconnect

import (
	"context"
	"net/http"
)

type Profile struct {
	UserID        string         `json:"user_id"`
	UserName      string         `json:"user_name"`
	UserShortName string         `json:"user_shortname"`
	UserType      string         `json:"user_type"`
	Email         string         `json:"email"`
	Broker        string         `json:"broker"`
	Exchanges     []string       `json:"exchanges"`
	Products      []string       `json:"products"`
	OrderTypes    []string       `json:"order_types"`
	AvatarURL     string         `json:"avatar_url"`
	Meta          map[string]any `json:"meta"`
}

type Margins struct {
	Equity    MarginSegment `json:"equity"`
	Commodity MarginSegment `json:"commodity"`
}

type MarginSegment struct {
	Enabled   bool         `json:"enabled"`
	Net       float64      `json:"net"`
	Available MarginValues `json:"available"`
	Utilised  MarginValues `json:"utilised"`
}

type MarginValues struct {
	Cash             float64 `json:"cash"`
	OpeningBalance   float64 `json:"opening_balance"`
	LiveBalance      float64 `json:"live_balance"`
	Collateral       float64 `json:"collateral"`
	IntradayPayin    float64 `json:"intraday_payin"`
	AdhocMargin      float64 `json:"adhoc_margin"`
	UsedMargin       float64 `json:"used_margin"`
	Debits           float64 `json:"debits"`
	Exposure         float64 `json:"exposure"`
	M2MRealised      float64 `json:"m2m_realised"`
	M2MUnrealised    float64 `json:"m2m_unrealised"`
	OptionPremium    float64 `json:"option_premium"`
	Payout           float64 `json:"payout"`
	SPAN             float64 `json:"span"`
	HoldingSales     float64 `json:"holding_sales"`
	Turnover         float64 `json:"turnover"`
	LiquidCollateral float64 `json:"liquid_collateral"`
	StockCollateral  float64 `json:"stock_collateral"`
	Delivery         float64 `json:"delivery"`
}

func (c *Client) Profile(ctx context.Context) (*Profile, error) {
	var out Profile
	if err := c.do(ctx, http.MethodGet, "/user/profile", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Margins(ctx context.Context) (*Margins, error) {
	var out Margins
	if err := c.do(ctx, http.MethodGet, "/user/margins", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SegmentMargins(ctx context.Context, segment string) (*MarginSegment, error) {
	var out MarginSegment
	if err := c.do(ctx, http.MethodGet, "/user/margins/"+segment, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
