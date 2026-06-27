package instruments

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/suyotech/zerodhaapi-go/kiteconnect"
)

func TestInstrumentConstants(t *testing.T) {
	instrumentTypes := []string{
		InstrumentTypeEQ,
		InstrumentTypeFUT,
		InstrumentTypeCE,
		InstrumentTypePE,
	}
	if got, want := len(instrumentTypes), 4; got != want {
		t.Fatalf("unexpected instrument type count: %d", got)
	}

	optionTypes := []string{
		OptionTypeCE,
		OptionTypePE,
	}
	if got, want := len(optionTypes), 2; got != want {
		t.Fatalf("unexpected option type count: %d", got)
	}

	exchanges := []string{
		ExchangeNSE,
		ExchangeBSE,
		ExchangeNFO,
		ExchangeBFO,
		ExchangeCDS,
		ExchangeBCD,
		ExchangeMCX,
	}
	if got, want := len(exchanges), 7; got != want {
		t.Fatalf("unexpected exchange count: %d", got)
	}

	segments := []string{
		SegmentNSE,
		SegmentBSE,
		SegmentIndices,
		SegmentNFOFUT,
		SegmentNFOOPT,
		SegmentBFOFUT,
		SegmentBFOOPT,
		SegmentCDSFUT,
		SegmentCDSOPT,
		SegmentBCDFUT,
		SegmentBCDOPT,
		SegmentMCXFUT,
		SegmentMCXOPT,
	}
	if got, want := len(segments), 13; got != want {
		t.Fatalf("unexpected segment count: %d", got)
	}
}

func TestLoadAndFind(t *testing.T) {
	t.Chdir(t.TempDir())
	filePath := "data/instruments.csv"
	instrumentFilePath = filePath
	csv := `instrument_token,exchange_token,tradingsymbol,name,last_price,expiry,strike,tick_size,lot_size,instrument_type,segment,exchange
1,10,NIFTY26JUN24200CE,NIFTY,0,2026-06-25,24200,0.05,75,CE,NFO-OPT,NFO
2,11,NIFTY26JUN24200PE,NIFTY,0,2026-06-25,24200,0.05,75,PE,NFO-OPT,NFO
`
	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte(csv), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(filePath, time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}

	items, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	name := "NIFTY"
	exchange := "NFO"
	found, err := Find(FindRequest{Name: &name, Exchange: &exchange}, items)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 2 {
		t.Fatalf("expected two instruments, got %#v", found)
	}
	strikes, err := OptionStrikes(FindRequest{Name: &name, Exchange: &exchange}, items)
	if err != nil {
		t.Fatal(err)
	}
	if len(strikes) != 1 || strikes[0] != 24200 {
		t.Fatalf("unexpected strikes: %#v", strikes)
	}
	expiries, err := ExpiryDates(FindRequest{Name: &name, Exchange: &exchange}, items)
	if err != nil {
		t.Fatal(err)
	}
	if len(expiries) != 1 {
		t.Fatalf("unexpected expiries: %#v", expiries)
	}
	if len(items) != 2 {
		t.Fatalf("expected all instruments in memory, got %d", len(items))
	}
}

func TestCheckDownloadDownloadsAllInstruments(t *testing.T) {
	t.Chdir(t.TempDir())
	filePath := "cache/all-instruments.csv"
	csv := `instrument_token,exchange_token,tradingsymbol,name,last_price,expiry,strike,tick_size,lot_size,instrument_type,segment,exchange
1,10,INFY,INFY,0,,0,0.05,1,EQ,NSE,NSE
2,20,NIFTY26JUN24200CE,NIFTY,0,2026-06-25,24200,0.05,75,CE,NFO-OPT,NFO
`
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(csv))
	}))
	defer server.Close()

	client := kiteconnect.NewClient("key", "secret", kiteconnect.WithBaseURL(server.URL))
	client.SetAccessToken("access")

	if err := CheckDownload(context.Background(), client, filePath); err != nil {
		t.Fatal(err)
	}
	items, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/instruments" {
		t.Fatalf("expected full instruments endpoint, got %q", gotPath)
	}
	if len(items) != 2 {
		t.Fatalf("expected all downloaded instruments, got %#v", items)
	}

	token := uint32(1)
	found, err := Find(FindRequest{InstrumentToken: &token})
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 || found[0].TradingSymbol != "INFY" {
		t.Fatalf("unexpected token lookup: %#v", found)
	}
}
