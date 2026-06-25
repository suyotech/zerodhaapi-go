package ticker

import (
	"encoding/binary"
	"testing"
)

func TestParseLTP(t *testing.T) {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint16(payload[0:2], 1)
	binary.BigEndian.PutUint16(payload[2:4], 8)
	binary.BigEndian.PutUint32(payload[4:8], 408065)
	binary.BigEndian.PutUint32(payload[8:12], 152345)

	ticks, err := Parse(payload)
	if err != nil {
		t.Fatal(err)
	}
	if len(ticks) != 1 || ticks[0].InstrumentToken != 408065 || ticks[0].LastPrice != 1523.45 {
		t.Fatalf("unexpected ticks: %#v", ticks)
	}
}
