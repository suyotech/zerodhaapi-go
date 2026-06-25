package instruments

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadOrDownloadAndFind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "instruments.csv")
	csv := `instrument_token,exchange_token,tradingsymbol,name,last_price,expiry,strike,tick_size,lot_size,instrument_type,segment,exchange
1,10,NIFTY26JUN24200CE,NIFTY,0,2026-06-25,24200,0.05,75,CE,NFO-OPT,NFO
2,11,NIFTY26JUN24200PE,NIFTY,0,2026-06-25,24200,0.05,75,PE,NFO-OPT,NFO
`
	if err := os.WriteFile(path, []byte(csv), 0o644); err != nil {
		t.Fatal(err)
	}

	store, err := LoadOrDownload(context.Background(), CacheConfig{
		FilePath: path,
		Now:      func() time.Time { return time.Date(2026, 6, 25, 8, 0, 0, 0, time.Local) },
	})
	if err != nil {
		t.Fatal(err)
	}
	name := "NIFTY"
	exchange := "NFO"
	items := Find(store, FindRequest{Name: &name, Exchange: &exchange})
	if len(items) != 2 {
		t.Fatalf("expected two instruments, got %#v", items)
	}
	strikes := OptionStrikes(store, FindRequest{Name: &name, Exchange: &exchange})
	if len(strikes) != 1 || strikes[0] != 24200 {
		t.Fatalf("unexpected strikes: %#v", strikes)
	}
	expiries := ExpiryDates(store, FindRequest{Name: &name, Exchange: &exchange})
	if len(expiries) != 1 {
		t.Fatalf("unexpected expiries: %#v", expiries)
	}
}
