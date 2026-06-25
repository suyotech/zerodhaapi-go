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

const (
	refreshHour   = 8
	refreshMinute = 30
	locationName  = "Asia/Kolkata"
)

const (
	InstrumentTypeEQ  = "EQ"
	InstrumentTypeFUT = "FUT"
	InstrumentTypeCE  = "CE"
	InstrumentTypePE  = "PE"
)

const (
	OptionTypeCE = InstrumentTypeCE
	OptionTypePE = InstrumentTypePE
)

const (
	ExchangeNSE = "NSE"
	ExchangeBSE = "BSE"
	ExchangeNFO = "NFO"
	ExchangeBFO = "BFO"
	ExchangeCDS = "CDS"
	ExchangeBCD = "BCD"
	ExchangeMCX = "MCX"
)

const (
	SegmentNSE     = "NSE"
	SegmentBSE     = "BSE"
	SegmentIndices = "INDICES"
	SegmentNFOFUT  = "NFO-FUT"
	SegmentNFOOPT  = "NFO-OPT"
	SegmentBFOFUT  = "BFO-FUT"
	SegmentBFOOPT  = "BFO-OPT"
	SegmentCDSFUT  = "CDS-FUT"
	SegmentCDSOPT  = "CDS-OPT"
	SegmentBCDFUT  = "BCD-FUT"
	SegmentBCDOPT  = "BCD-OPT"
	SegmentMCXFUT  = "MCX-FUT"
	SegmentMCXOPT  = "MCX-OPT"
)

var instrumentFilePath = "data/instruments.csv"

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

func CheckDownload(ctx context.Context, client *kiteconnect.Client, filePath string) error {
	if strings.TrimSpace(filePath) != "" {
		instrumentFilePath = filePath
	}
	info, err := os.Stat(instrumentFilePath)
	download := os.IsNotExist(err)
	if err != nil && !download {
		return nil
	}
	if !download {
		now := time.Now()
		loc, err := time.LoadLocation(locationName)
		if err != nil {
			loc = time.Local
		}
		now = now.In(loc)
		refreshAt := time.Date(now.Year(), now.Month(), now.Day(), refreshHour, refreshMinute, 0, 0, loc)
		download = !now.Before(refreshAt) && info.ModTime().In(loc).Before(refreshAt)
	}
	if !download {
		return nil
	}
	if client == nil {
		return errors.New("instruments: client is required to download")
	}
	body, err := client.Instruments(ctx, "")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(instrumentFilePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(instrumentFilePath, body, 0o644)
}

func Load() ([]Instrument, error) {
	return load(instrumentFilePath)
}

func Find(req FindRequest, instruments ...[]Instrument) ([]Instrument, error) {
	items := firstInstruments(instruments)
	if items == nil {
		loaded, err := Load()
		if err != nil {
			return nil, err
		}
		items = loaded
	}

	out := make([]Instrument, 0)
	for _, inst := range items {
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
	return out, nil
}

func ExpiryDates(req FindRequest, instruments ...[]Instrument) ([]time.Time, error) {
	items, err := Find(req, instruments...)
	if err != nil {
		return nil, err
	}
	seen := map[time.Time]bool{}
	for _, inst := range items {
		if !inst.Expiry.IsZero() {
			seen[inst.Expiry] = true
		}
	}
	out := make([]time.Time, 0, len(seen))
	for expiry := range seen {
		out = append(out, expiry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out, nil
}

func OptionStrikes(req FindRequest, instruments ...[]Instrument) ([]float64, error) {
	items, err := Find(req, instruments...)
	if err != nil {
		return nil, err
	}
	seen := map[float64]bool{}
	for _, inst := range items {
		if inst.InstrumentType == InstrumentTypeCE || inst.InstrumentType == InstrumentTypePE {
			seen[inst.Strike] = true
		}
	}
	out := make([]float64, 0, len(seen))
	for strike := range seen {
		out = append(out, strike)
	}
	sort.Float64s(out)
	return out, nil
}

func firstInstruments(instruments [][]Instrument) []Instrument {
	if len(instruments) == 0 {
		return nil
	}
	items := instruments[0]
	if items == nil {
		return nil
	}
	return items
}

func load(path string) ([]Instrument, error) {
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

	items := make([]Instrument, 0)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		items = append(items, instrument(row, index))
	}
	return items, nil
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
