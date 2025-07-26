package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "charitabletrust"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_charitabletrust"
)

// Store key prefixes
var (
	ParamsKey               = []byte{0x01}
	TrustFundBalanceKey     = []byte{0x02}
	CharitableAllocationKey = []byte{0x03}
	AllocationCountKey      = []byte{0x04}
	ImpactReportKey         = []byte{0x05}
	FraudAlertKey           = []byte{0x06}
	TrustGovernanceKey      = []byte{0x07}
	AllocationProposalKey   = []byte{0x08}
	ProposalCountKey        = []byte{0x09}
	
	// Index keys
	AllocationByCategoryKey     = []byte{0x10}
	AllocationByStatusKey       = []byte{0x11}
	AllocationByOrgKey          = []byte{0x12}
	ReportByAllocationKey       = []byte{0x13}
	AlertByAllocationKey        = []byte{0x14}
	ProposalByStatusKey         = []byte{0x15}
)

// GetCharitableAllocationKey returns the key for a charitable allocation
func GetCharitableAllocationKey(allocationID uint64) []byte {
	return append(CharitableAllocationKey, GetUint64Bytes(allocationID)...)
}

// GetImpactReportKey returns the key for an impact report
func GetImpactReportKey(reportID uint64) []byte {
	return append(ImpactReportKey, GetUint64Bytes(reportID)...)
}

// GetFraudAlertKey returns the key for a fraud alert
func GetFraudAlertKey(alertID uint64) []byte {
	return append(FraudAlertKey, GetUint64Bytes(alertID)...)
}

// GetAllocationProposalKey returns the key for an allocation proposal
func GetAllocationProposalKey(proposalID uint64) []byte {
	return append(AllocationProposalKey, GetUint64Bytes(proposalID)...)
}

// GetAllocationByCategoryKey returns the key for allocation by category index
func GetAllocationByCategoryKey(category string, allocationID uint64) []byte {
	return append(append(AllocationByCategoryKey, []byte(category)...), GetUint64Bytes(allocationID)...)
}

// GetAllocationByStatusKey returns the key for allocation by status index
func GetAllocationByStatusKey(status string, allocationID uint64) []byte {
	return append(append(AllocationByStatusKey, []byte(status)...), GetUint64Bytes(allocationID)...)
}

// GetAllocationByOrgKey returns the key for allocation by organization index
func GetAllocationByOrgKey(orgID uint64, allocationID uint64) []byte {
	key := append(AllocationByOrgKey, GetUint64Bytes(orgID)...)
	return append(key, GetUint64Bytes(allocationID)...)
}

// GetReportByAllocationKey returns the key for reports by allocation index
func GetReportByAllocationKey(allocationID uint64, reportID uint64) []byte {
	key := append(ReportByAllocationKey, GetUint64Bytes(allocationID)...)
	return append(key, GetUint64Bytes(reportID)...)
}

// GetAlertByAllocationKey returns the key for alerts by allocation index
func GetAlertByAllocationKey(allocationID uint64, alertID uint64) []byte {
	key := append(AlertByAllocationKey, GetUint64Bytes(allocationID)...)
	return append(key, GetUint64Bytes(alertID)...)
}

// GetProposalByStatusKey returns the key for proposals by status index
func GetProposalByStatusKey(status string, proposalID uint64) []byte {
	return append(append(ProposalByStatusKey, []byte(status)...), GetUint64Bytes(proposalID)...)
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