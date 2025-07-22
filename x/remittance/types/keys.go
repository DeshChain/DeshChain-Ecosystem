package types

import (
	"encoding/binary"
	"strings"
)

const (
	// ModuleName defines the module name
	ModuleName = "remittance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_remittance"

	// RemittancePortID defines the port ID for IBC remittance transfers
	RemittancePortID = "remittance"

	// Version defines the current version the IBC module supports
	Version = "remittance-1"
)

// KVStore keys
var (
	// ParamsKey defines the key to store module parameters
	ParamsKey = []byte{0x01}

	// RemittanceTransferKeyPrefix defines the prefix for remittance transfer keys
	RemittanceTransferKeyPrefix = []byte{0x10}

	// LiquidityPoolKeyPrefix defines the prefix for liquidity pool keys
	LiquidityPoolKeyPrefix = []byte{0x20}

	// ExchangeRateKeyPrefix defines the prefix for exchange rate keys
	ExchangeRateKeyPrefix = []byte{0x30}

	// CorridorKeyPrefix defines the prefix for corridor keys
	CorridorKeyPrefix = []byte{0x40}

	// PartnerKeyPrefix defines the prefix for partner keys
	PartnerKeyPrefix = []byte{0x50}

	// SettlementKeyPrefix defines the prefix for settlement keys
	SettlementKeyPrefix = []byte{0x60}

	// CountersKey defines the key for module counters
	CountersKey = []byte{0x70}

	// CountryConfigKeyPrefix defines the prefix for country configuration keys
	CountryConfigKeyPrefix = []byte{0x80}

	// CurrencyConfigKeyPrefix defines the prefix for currency configuration keys
	CurrencyConfigKeyPrefix = []byte{0x90}

	// TransferBySenderKeyPrefix defines the prefix for transfer by sender index
	TransferBySenderKeyPrefix = []byte{0xA0}

	// TransferByRecipientKeyPrefix defines the prefix for transfer by recipient index
	TransferByRecipientKeyPrefix = []byte{0xB0}

	// TransferByStatusKeyPrefix defines the prefix for transfer by status index
	TransferByStatusKeyPrefix = []byte{0xC0}

	// TransferByCorridorKeyPrefix defines the prefix for transfer by corridor index
	TransferByCorridorKeyPrefix = []byte{0xD0}

	// LiquidityProviderKeyPrefix defines the prefix for liquidity provider keys
	LiquidityProviderKeyPrefix = []byte{0xE0}

	// ComplianceDataKeyPrefix defines the prefix for compliance data keys
	ComplianceDataKeyPrefix = []byte{0xF0}
	
	// Sewa Mitra key prefixes
	SewaMitraAgentKeyPrefix = []byte{0x100}
	SewaMitraCommissionKeyPrefix = []byte{0x110}
	AgentByCountryKeyPrefix = []byte{0x120}
	AgentByCityKeyPrefix = []byte{0x130}
	AgentByStatusKeyPrefix = []byte{0x140}
	AgentByCurrencyKeyPrefix = []byte{0x150}
	CommissionByAgentKeyPrefix = []byte{0x160}
	CommissionByTransferKeyPrefix = []byte{0x170}
	CommissionByStatusKeyPrefix = []byte{0x180}
)

// RemittanceTransferKey returns the store key for a remittance transfer
func RemittanceTransferKey(transferID string) []byte {
	return append(RemittanceTransferKeyPrefix, []byte(transferID)...)
}

// LiquidityPoolKey returns the store key for a liquidity pool
func LiquidityPoolKey(poolID string) []byte {
	return append(LiquidityPoolKeyPrefix, []byte(poolID)...)
}

// ExchangeRateKey returns the store key for an exchange rate
func ExchangeRateKey(baseCurrency, quoteCurrency string) []byte {
	key := append(ExchangeRateKeyPrefix, []byte(baseCurrency)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(quoteCurrency)...)
}

// CorridorKey returns the store key for a corridor
func CorridorKey(corridorID string) []byte {
	return append(CorridorKeyPrefix, []byte(corridorID)...)
}

// PartnerKey returns the store key for a partner
func PartnerKey(partnerID string) []byte {
	return append(PartnerKeyPrefix, []byte(partnerID)...)
}

// SettlementKey returns the store key for a settlement
func SettlementKey(transferID string) []byte {
	return append(SettlementKeyPrefix, []byte(transferID)...)
}

// CountryConfigKey returns the store key for a country configuration
func CountryConfigKey(countryCode string) []byte {
	return append(CountryConfigKeyPrefix, []byte(countryCode)...)
}

// CurrencyConfigKey returns the store key for a currency configuration
func CurrencyConfigKey(currencyCode string) []byte {
	return append(CurrencyConfigKeyPrefix, []byte(currencyCode)...)
}

// TransferBySenderKey returns the store key for transfer by sender index
func TransferBySenderKey(sender string, transferID string) []byte {
	key := append(TransferBySenderKeyPrefix, []byte(sender)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(transferID)...)
}

// TransferByRecipientKey returns the store key for transfer by recipient index
func TransferByRecipientKey(recipient string, transferID string) []byte {
	key := append(TransferByRecipientKeyPrefix, []byte(recipient)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(transferID)...)
}

// TransferByStatusKey returns the store key for transfer by status index
func TransferByStatusKey(status TransferStatus, transferID string) []byte {
	statusBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(statusBytes, uint32(status))
	key := append(TransferByStatusKeyPrefix, statusBytes...)
	key = append(key, []byte("/")...)
	return append(key, []byte(transferID)...)
}

// TransferByCorridorKey returns the store key for transfer by corridor index
func TransferByCorridorKey(corridorID string, transferID string) []byte {
	key := append(TransferByCorridorKeyPrefix, []byte(corridorID)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(transferID)...)
}

// LiquidityProviderKey returns the store key for a liquidity provider
func LiquidityProviderKey(poolID string, providerAddress string) []byte {
	key := append(LiquidityProviderKeyPrefix, []byte(poolID)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(providerAddress)...)
}

// ComplianceDataKey returns the store key for compliance data
func ComplianceDataKey(transferID string) []byte {
	return append(ComplianceDataKeyPrefix, []byte(transferID)...)
}

// ParseTransferBySenderKey parses a transfer by sender key
func ParseTransferBySenderKey(key []byte) (sender string, transferID string) {
	keyStr := string(key[len(TransferBySenderKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		sender = parts[0]
		transferID = parts[1]
	}
	return
}

// ParseTransferByRecipientKey parses a transfer by recipient key
func ParseTransferByRecipientKey(key []byte) (recipient string, transferID string) {
	keyStr := string(key[len(TransferByRecipientKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		recipient = parts[0]
		transferID = parts[1]
	}
	return
}

// ParseTransferByStatusKey parses a transfer by status key
func ParseTransferByStatusKey(key []byte) (status TransferStatus, transferID string) {
	if len(key) > len(TransferByStatusKeyPrefix)+4 {
		statusBytes := key[len(TransferByStatusKeyPrefix) : len(TransferByStatusKeyPrefix)+4]
		status = TransferStatus(binary.BigEndian.Uint32(statusBytes))
		
		remaining := key[len(TransferByStatusKeyPrefix)+4:]
		if len(remaining) > 1 && remaining[0] == '/' {
			transferID = string(remaining[1:])
		}
	}
	return
}

// ParseTransferByCorridorKey parses a transfer by corridor key
func ParseTransferByCorridorKey(key []byte) (corridorID string, transferID string) {
	keyStr := string(key[len(TransferByCorridorKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		corridorID = parts[0]
		transferID = parts[1]
	}
	return
}

// ParseLiquidityProviderKey parses a liquidity provider key
func ParseLiquidityProviderKey(key []byte) (poolID string, providerAddress string) {
	keyStr := string(key[len(LiquidityProviderKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		poolID = parts[0]
		providerAddress = parts[1]
	}
	return
}

// GetTransferIDFromKey extracts transfer ID from various key formats
func GetTransferIDFromKey(key []byte, prefix []byte) string {
	if len(key) <= len(prefix) {
		return ""
	}
	return string(key[len(prefix):])
}

// Event types
const (
	EventTypeInitiateTransfer   = "initiate_transfer"
	EventTypeConfirmTransfer    = "confirm_transfer"
	EventTypeCancelTransfer     = "cancel_transfer"
	EventTypeRefundTransfer     = "refund_transfer"
	EventTypeCompleteTransfer   = "complete_transfer"
	EventTypeAddLiquidity       = "add_liquidity"
	EventTypeRemoveLiquidity    = "remove_liquidity"
	EventTypeUpdateExchangeRate = "update_exchange_rate"
	EventTypeCreateCorridor     = "create_corridor"
	EventTypeRegisterPartner    = "register_partner"
	EventTypeProcessSettlement  = "process_settlement"
)

// Event attribute keys
const (
	AttributeKeyTransferID        = "transfer_id"
	AttributeKeySender            = "sender"
	AttributeKeyRecipient         = "recipient"
	AttributeKeyAmount            = "amount"
	AttributeKeySourceCountry     = "source_country"
	AttributeKeyDestinationCountry = "destination_country"
	AttributeKeySourceCurrency    = "source_currency"
	AttributeKeyDestinationCurrency = "destination_currency"
	AttributeKeyExchangeRate      = "exchange_rate"
	AttributeKeyFees              = "fees"
	AttributeKeySettlementMethod  = "settlement_method"
	AttributeKeyStatus            = "status"
	AttributeKeyCorridorID        = "corridor_id"
	AttributeKeyPartnerID         = "partner_id"
	AttributeKeyPoolID            = "pool_id"
	AttributeKeyProvider          = "provider"
	AttributeKeyLPTokens          = "lp_tokens"
	AttributeKeyBaseCurrency      = "base_currency"
	AttributeKeyQuoteCurrency     = "quote_currency"
	AttributeKeyRate              = "rate"
	AttributeKeyReason            = "reason"
	AttributeKeySettlementTime    = "settlement_time"
	
	// Sewa Mitra attribute keys
	AttributeKeySewaMitraAgentID = "seva_mitra_agent_id"
	AttributeKeySewaMitraCommission = "seva_mitra_commission"
	AttributeKeyAgentName = "agent_name"
	AttributeKeyAgentLocation = "agent_location"
	AttributeKeyCommissionAmount = "commission_amount"
)

// ============== Sewa Mitra Key Functions ==============

// SewaMitraAgentKey returns the store key for a Sewa Mitra agent
func SewaMitraAgentKey(agentID string) []byte {
	return append(SewaMitraAgentKeyPrefix, []byte(agentID)...)
}

// SewaMitraCommissionKey returns the store key for a Sewa Mitra commission
func SewaMitraCommissionKey(commissionID string) []byte {
	return append(SewaMitraCommissionKeyPrefix, []byte(commissionID)...)
}

// AgentByCountryKey returns the store key for agent by country index
func AgentByCountryKey(country string, agentID string) []byte {
	key := append(AgentByCountryKeyPrefix, []byte(country)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(agentID)...)
}

// AgentByCityKey returns the store key for agent by city index
func AgentByCityKey(country, city string, agentID string) []byte {
	key := append(AgentByCityKeyPrefix, []byte(country)...)
	key = append(key, []byte("/")...)
	key = append(key, []byte(city)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(agentID)...)
}

// AgentByStatusKey returns the store key for agent by status index
func AgentByStatusKey(status AgentStatus, agentID string) []byte {
	statusBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(statusBytes, uint32(status))
	key := append(AgentByStatusKeyPrefix, statusBytes...)
	key = append(key, []byte("/")...)
	return append(key, []byte(agentID)...)
}

// AgentByCurrencyKey returns the store key for agent by currency index
func AgentByCurrencyKey(currency string, agentID string) []byte {
	key := append(AgentByCurrencyKeyPrefix, []byte(currency)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(agentID)...)
}

// CommissionByAgentKey returns the store key for commission by agent index
func CommissionByAgentKey(agentID string, commissionID string) []byte {
	key := append(CommissionByAgentKeyPrefix, []byte(agentID)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(commissionID)...)
}

// CommissionByTransferKey returns the store key for commission by transfer index
func CommissionByTransferKey(transferID string, commissionID string) []byte {
	key := append(CommissionByTransferKeyPrefix, []byte(transferID)...)
	key = append(key, []byte("/")...)
	return append(key, []byte(commissionID)...)
}

// CommissionByStatusKey returns the store key for commission by status index
func CommissionByStatusKey(status CommissionStatus, commissionID string) []byte {
	statusBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(statusBytes, uint32(status))
	key := append(CommissionByStatusKeyPrefix, statusBytes...)
	key = append(key, []byte("/")...)
	return append(key, []byte(commissionID)...)
}

// ============== Sewa Mitra Parse Functions ==============

// ParseAgentByCountryKey parses an agent by country key
func ParseAgentByCountryKey(key []byte) (country string, agentID string) {
	keyStr := string(key[len(AgentByCountryKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		country = parts[0]
		agentID = parts[1]
	}
	return
}

// ParseAgentByCityKey parses an agent by city key
func ParseAgentByCityKey(key []byte) (country, city, agentID string) {
	keyStr := string(key[len(AgentByCityKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 3 {
		country = parts[0]
		city = parts[1]
		agentID = parts[2]
	}
	return
}

// ParseAgentByStatusKey parses an agent by status key
func ParseAgentByStatusKey(key []byte) (status AgentStatus, agentID string) {
	if len(key) > len(AgentByStatusKeyPrefix)+4 {
		statusBytes := key[len(AgentByStatusKeyPrefix) : len(AgentByStatusKeyPrefix)+4]
		status = AgentStatus(binary.BigEndian.Uint32(statusBytes))
		
		remaining := key[len(AgentByStatusKeyPrefix)+4:]
		if len(remaining) > 1 && remaining[0] == '/' {
			agentID = string(remaining[1:])
		}
	}
	return
}

// ParseAgentByCurrencyKey parses an agent by currency key
func ParseAgentByCurrencyKey(key []byte) (currency string, agentID string) {
	keyStr := string(key[len(AgentByCurrencyKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		currency = parts[0]
		agentID = parts[1]
	}
	return
}

// ParseCommissionByAgentKey parses a commission by agent key
func ParseCommissionByAgentKey(key []byte) (agentID string, commissionID string) {
	keyStr := string(key[len(CommissionByAgentKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		agentID = parts[0]
		commissionID = parts[1]
	}
	return
}

// ParseCommissionByTransferKey parses a commission by transfer key
func ParseCommissionByTransferKey(key []byte) (transferID string, commissionID string) {
	keyStr := string(key[len(CommissionByTransferKeyPrefix):])
	parts := strings.Split(keyStr, "/")
	if len(parts) >= 2 {
		transferID = parts[0]
		commissionID = parts[1]
	}
	return
}

// ParseCommissionByStatusKey parses a commission by status key
func ParseCommissionByStatusKey(key []byte) (status CommissionStatus, commissionID string) {
	if len(key) > len(CommissionByStatusKeyPrefix)+4 {
		statusBytes := key[len(CommissionByStatusKeyPrefix) : len(CommissionByStatusKeyPrefix)+4]
		status = CommissionStatus(binary.BigEndian.Uint32(statusBytes))
		
		remaining := key[len(CommissionByStatusKeyPrefix)+4:]
		if len(remaining) > 1 && remaining[0] == '/' {
			commissionID = string(remaining[1:])
		}
	}
	return
)