package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "dswf"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dswf"
)

// Store key prefixes
var (
	ParamsKey                = []byte{0x01}
	FundAllocationKey        = []byte{0x02}
	FundAllocationCountKey   = []byte{0x03}
	InvestmentPortfolioKey   = []byte{0x04}
	MonthlyReportKey         = []byte{0x05}
	FundGovernanceKey        = []byte{0x06}
	DisbursementScheduleKey  = []byte{0x07}
	PerformanceMetricKey     = []byte{0x08}
	
	// Index keys
	AllocationByCategoryKey  = []byte{0x10}
	AllocationByStatusKey    = []byte{0x11}
	AllocationByRecipientKey = []byte{0x12}
	PendingDisbursementKey   = []byte{0x13}
)

// GetFundAllocationKey returns the key for a fund allocation
func GetFundAllocationKey(allocationID uint64) []byte {
	return append(FundAllocationKey, GetUint64Bytes(allocationID)...)
}

// GetMonthlyReportKey returns the key for a monthly report
func GetMonthlyReportKey(period string) []byte {
	return append(MonthlyReportKey, []byte(period)...)
}

// GetAllocationByCategoryKey returns the key for allocation by category index
func GetAllocationByCategoryKey(category string, allocationID uint64) []byte {
	return append(append(AllocationByCategoryKey, []byte(category)...), GetUint64Bytes(allocationID)...)
}

// GetAllocationByStatusKey returns the key for allocation by status index
func GetAllocationByStatusKey(status string, allocationID uint64) []byte {
	return append(append(AllocationByStatusKey, []byte(status)...), GetUint64Bytes(allocationID)...)
}

// GetAllocationByRecipientKey returns the key for allocation by recipient index
func GetAllocationByRecipientKey(recipient string, allocationID uint64) []byte {
	return append(append(AllocationByRecipientKey, []byte(recipient)...), GetUint64Bytes(allocationID)...)
}

// GetDisbursementScheduleKey returns the key for a disbursement schedule
func GetDisbursementScheduleKey(allocationID uint64, index uint32) []byte {
	key := append(DisbursementScheduleKey, GetUint64Bytes(allocationID)...)
	return append(key, GetUint32Bytes(index)...)
}

// GetPerformanceMetricKey returns the key for performance metrics
func GetPerformanceMetricKey(allocationID uint64, metricName string) []byte {
	key := append(PerformanceMetricKey, GetUint64Bytes(allocationID)...)
	return append(key, []byte(metricName)...)
}

// GetPendingDisbursementKey returns the key for pending disbursements
func GetPendingDisbursementKey(scheduledDate int64, allocationID uint64, index uint32) []byte {
	key := append(PendingDisbursementKey, GetInt64Bytes(scheduledDate)...)
	key = append(key, GetUint64Bytes(allocationID)...)
	return append(key, GetUint32Bytes(index)...)
}

// GetUint64Bytes returns the byte representation of a uint64
func GetUint64Bytes(value uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, value)
	return bz
}

// GetUint64FromBytes returns uint64 from bytes
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// GetUint32Bytes returns the byte representation of a uint32
func GetUint32Bytes(value uint32) []byte {
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, value)
	return bz
}

// GetInt64Bytes returns the byte representation of an int64
func GetInt64Bytes(value int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(value))
	return bz
}