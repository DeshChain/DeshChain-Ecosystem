package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_oracle"
)

var (
	// ParamsKey defines the key to store the module params
	ParamsKey = []byte{0x01}

	// PriceDataKeyPrefix defines the prefix key for price data
	PriceDataKeyPrefix = []byte{0x02}

	// OracleValidatorKeyPrefix defines the prefix key for oracle validators
	OracleValidatorKeyPrefix = []byte{0x03}

	// ExchangeRateKeyPrefix defines the prefix key for exchange rates
	ExchangeRateKeyPrefix = []byte{0x04}

	// ValidatorSubmissionKeyPrefix defines the prefix key for validator submissions
	ValidatorSubmissionKeyPrefix = []byte{0x05}

	// PriceHistoryKeyPrefix defines the prefix key for price history
	PriceHistoryKeyPrefix = []byte{0x06}

	// AggregationWindowKeyPrefix defines the prefix key for aggregation windows
	AggregationWindowKeyPrefix = []byte{0x07}

	// ValidatorMissCounterKeyPrefix defines the prefix key for validator miss counters
	ValidatorMissCounterKeyPrefix = []byte{0x08}
)

// PriceDataKey returns the store key to retrieve a PriceData from the symbol
func PriceDataKey(symbol string) []byte {
	return append(PriceDataKeyPrefix, []byte(symbol)...)
}

// OracleValidatorKey returns the store key to retrieve an OracleValidator from the validator address
func OracleValidatorKey(validator string) []byte {
	return append(OracleValidatorKeyPrefix, []byte(validator)...)
}

// ExchangeRateKey returns the store key to retrieve an ExchangeRate from base and target currencies
func ExchangeRateKey(base, target string) []byte {
	key := append(ExchangeRateKeyPrefix, []byte(base)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(target)...)
}

// ValidatorSubmissionKey returns the store key for a validator's price submission
func ValidatorSubmissionKey(validator, symbol string, blockHeight uint64) []byte {
	key := append(ValidatorSubmissionKeyPrefix, []byte(validator)...)
	key = append(key, []byte("/")...)
	key = append(key, []byte(symbol)...)
	key = append(key, []byte("/")...)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, blockHeight)
	return append(key, heightBytes...)
}

// ValidatorSubmissionKeyByValidator returns the prefix key for all submissions by a validator
func ValidatorSubmissionKeyByValidator(validator string) []byte {
	key := append(ValidatorSubmissionKeyPrefix, []byte(validator)...)
	return append(key, []byte("/")...)
}

// ValidatorSubmissionKeyBySymbol returns the prefix key for all submissions for a symbol
func ValidatorSubmissionKeyBySymbol(symbol string) []byte {
	return append(ValidatorSubmissionKeyPrefix, []byte(symbol)...)
}

// PriceHistoryKey returns the store key for price history
func PriceHistoryKey(symbol string, blockHeight uint64) []byte {
	key := append(PriceHistoryKeyPrefix, []byte(symbol)...)
	key = append(key, []byte("/")...)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, blockHeight)
	return append(key, heightBytes...)
}

// PriceHistoryKeyBySymbol returns the prefix key for all price history of a symbol
func PriceHistoryKeyBySymbol(symbol string) []byte {
	key := append(PriceHistoryKeyPrefix, []byte(symbol)...)
	return append(key, []byte("/")...)
}

// AggregationWindowKey returns the store key for aggregation window data
func AggregationWindowKey(symbol string, windowStart uint64) []byte {
	key := append(AggregationWindowKeyPrefix, []byte(symbol)...)
	key = append(key, []byte("/")...)
	windowBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(windowBytes, windowStart)
	return append(key, windowBytes...)
}

// ValidatorMissCounterKey returns the store key for validator miss counter
func ValidatorMissCounterKey(validator string) []byte {
	return append(ValidatorMissCounterKeyPrefix, []byte(validator)...)
}