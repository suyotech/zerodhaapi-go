package ticker

import "encoding/binary"

type Tick struct {
	InstrumentToken uint32
	Mode            Mode
	LastPrice       float64
	LastQuantity    uint32
	AveragePrice    float64
	Volume          uint32
	BuyQuantity     uint32
	SellQuantity    uint32
	Open            float64
	High            float64
	Low             float64
	Close           float64
	Change          float64
	Depth           Depth
}

type Depth struct {
	Buy  []DepthItem
	Sell []DepthItem
}

type DepthItem struct {
	Quantity uint32
	Price    float64
	Orders   uint16
}

func Parse(payload []byte) ([]Tick, error) {
	if len(payload) == 1 {
		return nil, nil
	}
	if len(payload) < 2 {
		return nil, nil
	}
	count := int(binary.BigEndian.Uint16(payload[0:2]))
	offset := 2
	ticks := make([]Tick, 0, count)
	for i := 0; i < count; i++ {
		if offset+2 > len(payload) {
			break
		}
		size := int(binary.BigEndian.Uint16(payload[offset : offset+2]))
		offset += 2
		if offset+size > len(payload) {
			break
		}
		tick := parsePacket(payload[offset : offset+size])
		offset += size
		if tick.InstrumentToken != 0 {
			ticks = append(ticks, tick)
		}
	}
	return ticks, nil
}

func parsePacket(packet []byte) Tick {
	if len(packet) < 8 {
		return Tick{}
	}
	tick := Tick{
		InstrumentToken: binary.BigEndian.Uint32(packet[0:4]),
		LastPrice:       price(packet[4:8]),
	}
	switch len(packet) {
	case 8:
		tick.Mode = ModeLTP
	case 28, 32, 44:
		tick.Mode = ModeQuote
		parseQuote(packet, &tick)
	default:
		tick.Mode = ModeFull
		parseQuote(packet, &tick)
		if len(packet) >= 184 {
			parseDepth(packet[64:184], &tick)
		}
	}
	return tick
}

func parseQuote(packet []byte, tick *Tick) {
	if len(packet) >= 28 {
		tick.LastQuantity = binary.BigEndian.Uint32(packet[8:12])
		tick.AveragePrice = price(packet[12:16])
		tick.Volume = binary.BigEndian.Uint32(packet[16:20])
		tick.BuyQuantity = binary.BigEndian.Uint32(packet[20:24])
		tick.SellQuantity = binary.BigEndian.Uint32(packet[24:28])
	}
	if len(packet) >= 44 {
		tick.Open = price(packet[28:32])
		tick.High = price(packet[32:36])
		tick.Low = price(packet[36:40])
		tick.Close = price(packet[40:44])
		if tick.Close != 0 {
			tick.Change = ((tick.LastPrice - tick.Close) / tick.Close) * 100
		}
	}
}

func parseDepth(packet []byte, tick *Tick) {
	for i := 0; i < 10 && i*12+12 <= len(packet); i++ {
		item := DepthItem{
			Quantity: binary.BigEndian.Uint32(packet[i*12 : i*12+4]),
			Price:    price(packet[i*12+4 : i*12+8]),
			Orders:   binary.BigEndian.Uint16(packet[i*12+8 : i*12+10]),
		}
		if i < 5 {
			tick.Depth.Buy = append(tick.Depth.Buy, item)
		} else {
			tick.Depth.Sell = append(tick.Depth.Sell, item)
		}
	}
}

func price(raw []byte) float64 {
	return float64(binary.BigEndian.Uint32(raw)) / 100
}
