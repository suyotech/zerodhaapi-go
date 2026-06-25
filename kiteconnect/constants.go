package kiteconnect

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

const (
	VarietyRegular = "regular"
	VarietyAMO     = "amo"
	VarietyCO      = "co"
	VarietyIceberg = "iceberg"
	VarietyAuction = "auction"
)

const (
	TransactionBuy  = "BUY"
	TransactionSell = "SELL"
)

const (
	PositionDay      = "day"
	PositionNet      = "net"
	HoldingTypeDemat = "demat"
	HoldingTypeT1    = "t1"
)

const (
	OrderTypeMarket = "MARKET"
	OrderTypeLimit  = "LIMIT"
	OrderTypeSL     = "SL"
	OrderTypeSLM    = "SL-M"
)

const (
	ProductCNC  = "CNC"
	ProductNRML = "NRML"
	ProductMIS  = "MIS"
	ProductMTF  = "MTF"
)

const (
	ValidityDay = "DAY"
	ValidityIOC = "IOC"
	ValidityTTL = "TTL"
)

const (
	OrderStatusComplete              = "COMPLETE"
	OrderStatusRejected              = "REJECTED"
	OrderStatusCancelled             = "CANCELLED"
	OrderStatusOpen                  = "OPEN"
	OrderStatusOpenPending           = "OPEN PENDING"
	OrderStatusValidationPending     = "VALIDATION PENDING"
	OrderStatusPutOrderRequest       = "PUT ORDER REQ RECEIVED"
	OrderStatusTriggerPending        = "TRIGGER PENDING"
	OrderStatusCancelPending         = "CANCEL PENDING"
	OrderStatusModifyPending         = "MODIFY PENDING"
	OrderStatusModifyValidation      = "MODIFY VALIDATION PENDING"
	OrderStatusAMORequest            = "AMO REQ RECEIVED"
	OrderStatusMarketProtectionError = "MARKET PROTECTION ERROR"
)

const (
	MarginTypeEquity    = "equity"
	MarginTypeCommodity = "commodity"
)

const (
	IntervalMinute   = "minute"
	IntervalDay      = "day"
	Interval3Minute  = "3minute"
	Interval5Minute  = "5minute"
	Interval10Minute = "10minute"
	Interval15Minute = "15minute"
	Interval30Minute = "30minute"
	Interval60Minute = "60minute"
)

const (
	GTTTypeSingle = "single"
	GTTTypeOCO    = "two-leg"
)

const (
	GTTStatusActive    = "active"
	GTTStatusTriggered = "triggered"
	GTTStatusDisabled  = "disabled"
	GTTStatusExpired   = "expired"
	GTTStatusCancelled = "cancelled"
	GTTStatusRejected  = "rejected"
	GTTStatusDeleted   = "deleted"
)

const (
	MFTransactionBuy  = "BUY"
	MFTransactionSell = "SELL"
)

const (
	MFOrderVarietyRegular = "regular"
	MFOrderVarietySIP     = "sip"
)

const (
	MFPurchaseTypeFresh      = "FRESH"
	MFPurchaseTypeAdditional = "ADDITIONAL"
)

const (
	MFDividendTypePayout       = "payout"
	MFDividendTypeReinvestment = "reinvestment"
	MFDividendTypeGrowth       = "growth"
)

const (
	MFSIPFrequencyWeekly     = "weekly"
	MFSIPFrequencyMonthly    = "monthly"
	MFSIPFrequencyQuarterly  = "quarterly"
	MFSIPFrequencyHalfYearly = "half-yearly"
)

const (
	MFStatusComplete  = "COMPLETE"
	MFStatusPending   = "PENDING"
	MFStatusCancelled = "CANCELLED"
	MFStatusRejected  = "REJECTED"
)
