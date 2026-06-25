package kiteconnect

import "testing"

func TestConstants(t *testing.T) {
	tests := map[string]string{
		"ExchangeNSE":                  ExchangeNSE,
		"SegmentNFOOPT":                SegmentNFOOPT,
		"VarietyRegular":               VarietyRegular,
		"TransactionBuy":               TransactionBuy,
		"PositionDay":                  PositionDay,
		"OrderTypeMarket":              OrderTypeMarket,
		"ProductMIS":                   ProductMIS,
		"ValidityDay":                  ValidityDay,
		"OrderStatusComplete":          OrderStatusComplete,
		"MarginTypeEquity":             MarginTypeEquity,
		"IntervalMinute":               IntervalMinute,
		"GTTTypeSingle":                GTTTypeSingle,
		"GTTStatusActive":              GTTStatusActive,
		"MFTransactionBuy":             MFTransactionBuy,
		"MFOrderVarietyRegular":        MFOrderVarietyRegular,
		"MFPurchaseTypeFresh":          MFPurchaseTypeFresh,
		"MFDividendTypeReinvestment":   MFDividendTypeReinvestment,
		"MFSIPFrequencyMonthly":        MFSIPFrequencyMonthly,
		"MFStatusComplete":             MFStatusComplete,
		"OrderStatusValidationPending": OrderStatusValidationPending,
	}

	for name, value := range tests {
		if value == "" {
			t.Fatalf("%s is empty", name)
		}
	}
}
