package instruments

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"zerodhaapi-go/kiteconnect"
)

type CacheConfig struct {
	Client       *kiteconnect.Client
	FilePath     string
	Exchange     string
	RefreshHour  int
	RefreshMin   int
	LocationName string
	Now          func() time.Time
}

type Store struct {
	items []Instrument
}

type Instrument struct {
	InstrumentToken uint32
	ExchangeToken   uint32
	TradingSymbol   string
	Name            string
	LastPrice       float64
	Expiry          time.Time
	Strike          float64
	TickSize        float64
	LotSize         int
	InstrumentType  string
	Segment         string
	Exchange        string
}

type FindRequest struct {
	InstrumentToken *uint32
	ExchangeToken   *uint32
	Name            *string
	TradingSymbol   *string
	InstrumentType  *string
	Strike          *float64
	Expiry          *time.Time
	Segment         *string
	Exchange        *string
}

func LoadOrDownload(ctx context.Context, cfg CacheConfig) (*Store, error) {
	if cfg.FilePath == "" {
		return nil, errors.New("instruments: file path is required")
	}
	if needsDownload(cfg) {
		if cfg.Client == nil {
			return nil, errors.New("instruments: client is required to download")
		}
		body, err := cfg.Client.Instruments(ctx, cfg.Exchange)
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(cfg.FilePath, body, 0o644); err != nil {
			return nil, err
		}
	}
	return load(cfg.FilePath)
}

func Find(store *Store, req FindRequest) []Instrument {
	if store == nil {
		return nil
	}
	out := make([]Instrument, 0)
	for _, inst := range store.items {
		if matches(inst, req) {
			out = append(out, inst)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Exchange == out[j].Exchange {
			return out[i].TradingSymbol < out[j].TradingSymbol
		}
		return out[i].Exchange < out[j].Exchange
	})
	return out
}

func ExpiryDates(store *Store, req FindRequest) []time.Time {
	seen := map[time.Time]bool{}
	for _, inst := range Find(store, req) {
		if !inst.Expiry.IsZero() {
			seen[inst.Expiry] = true
		}
	}
	out := make([]time.Time, 0, len(seen))
	for expiry := range seen {
		out = append(out, expiry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out
}

func OptionStrikes(store *Store, req FindRequest) []float64 {
	seen := map[float64]bool{}
	for _, inst := range Find(store, req) {
		if inst.InstrumentType == "CE" || inst.InstrumentType == "PE" {
			seen[inst.Strike] = true
		}
	}
	out := make([]float64, 0, len(seen))
	for strike := range seen {
		out = append(out, strike)
	}
	sort.Float64s(out)
	return out
}

func needsDownload(cfg CacheConfig) bool {
	info, err := os.Stat(cfg.FilePath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		return false
	}
	now := time.Now()
	if cfg.Now != nil {
		now = cfg.Now()
	}
	locName := cfg.LocationName
	if locName == "" {
		locName = "Asia/Kolkata"
	}
	loc, err := time.LoadLocation(locName)
	if err != nil {
		loc = time.Local
	}
	hour := cfg.RefreshHour
	if hour == 0 {
		hour = 8
	}
	min := cfg.RefreshMin
	if min == 0 {
		min = 30
	}
	now = now.In(loc)
	refreshAt := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, loc)
	return !now.Before(refreshAt) && info.ModTime().In(loc).Before(refreshAt)
}

func load(path string) (*Store, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	index := map[string]int{}
	for i, column := range header {
		index[strings.TrimSpace(column)] = i
	}

	store := &Store{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		store.items = append(store.items, instrument(row, index))
	}
	return store, nil
}

func instrument(row []string, index map[string]int) Instrument {
	value := func(name string) string {
		i, ok := index[name]
		if !ok || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}
	return Instrument{
		InstrumentToken: uint32Value(value("instrument_token")),
		ExchangeToken:   uint32Value(value("exchange_token")),
		TradingSymbol:   value("tradingsymbol"),
		Name:            value("name"),
		LastPrice:       floatValue(value("last_price")),
		Expiry:          dateValue(value("expiry")),
		Strike:          floatValue(value("strike")),
		TickSize:        floatValue(value("tick_size")),
		LotSize:         intValue(value("lot_size")),
		InstrumentType:  value("instrument_type"),
		Segment:         value("segment"),
		Exchange:        value("exchange"),
	}
}

func matches(inst Instrument, req FindRequest) bool {
	if req.InstrumentToken != nil && inst.InstrumentToken != *req.InstrumentToken {
		return false
	}
	if req.ExchangeToken != nil && inst.ExchangeToken != *req.ExchangeToken {
		return false
	}
	if req.Name != nil && !sameText(inst.Name, *req.Name) {
		return false
	}
	if req.TradingSymbol != nil && !sameText(inst.TradingSymbol, *req.TradingSymbol) {
		return false
	}
	if req.InstrumentType != nil && !sameText(inst.InstrumentType, *req.InstrumentType) {
		return false
	}
	if req.Strike != nil && inst.Strike != *req.Strike {
		return false
	}
	if req.Expiry != nil && !sameDate(inst.Expiry, *req.Expiry) {
		return false
	}
	if req.Segment != nil && !sameText(inst.Segment, *req.Segment) {
		return false
	}
	if req.Exchange != nil && !sameText(inst.Exchange, *req.Exchange) {
		return false
	}
	return true
}

func sameText(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func uint32Value(value string) uint32 {
	n, _ := strconv.ParseUint(value, 10, 32)
	return uint32(n)
}

func intValue(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}

func floatValue(value string) float64 {
	n, _ := strconv.ParseFloat(value, 64)
	return n
}

func dateValue(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	for _, layout := range []string{"2006-01-02", "02-Jan-2006"} {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
