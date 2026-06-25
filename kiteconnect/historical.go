package kiteconnect

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type HistoricalDataParams struct {
	InstrumentToken uint32
	Interval        string
	From            time.Time
	To              time.Time
	Continuous      bool
	OI              bool
}

type HistoricalCandle struct {
	Time   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
	OI     int
}

func (c *Client) HistoricalData(ctx context.Context, p HistoricalDataParams) ([]HistoricalCandle, error) {
	q := url.Values{}
	q.Set("from", p.From.Format("2006-01-02 15:04:05"))
	q.Set("to", p.To.Format("2006-01-02 15:04:05"))
	if p.Continuous {
		q.Set("continuous", "1")
	}
	if p.OI {
		q.Set("oi", "1")
	}

	var raw struct {
		Candles [][]any `json:"candles"`
	}
	path := "/instruments/historical/" + strconv.FormatUint(uint64(p.InstrumentToken), 10) + "/" + p.Interval
	if err := c.do(ctx, http.MethodGet, path, q, nil, &raw); err != nil {
		return nil, err
	}

	candles := make([]HistoricalCandle, 0, len(raw.Candles))
	for _, row := range raw.Candles {
		if len(row) < 6 {
			continue
		}
		candle := HistoricalCandle{
			Time:   stringValue(row[0]),
			Open:   numberValue(row[1]),
			High:   numberValue(row[2]),
			Low:    numberValue(row[3]),
			Close:  numberValue(row[4]),
			Volume: int(numberValue(row[5])),
		}
		if len(row) > 6 {
			candle.OI = int(numberValue(row[6]))
		}
		candles = append(candles, candle)
	}
	return candles, nil
}

func stringValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func numberValue(v any) float64 {
	switch value := v.(type) {
	case float64:
		return value
	case json.Number:
		f, _ := value.Float64()
		return f
	default:
		return 0
	}
}
