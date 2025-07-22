package types

const (
	// ModuleName defines the module name
	ModuleName = "tradefinance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for trade finance
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tradefinance"
)

// Store key prefixes
var (
	ParamsKey                 = []byte{0x01}
	TradePartyPrefix          = []byte{0x02}
	LetterOfCreditPrefix      = []byte{0x03}
	TradeDocumentPrefix       = []byte{0x04}
	InsurancePolicyPrefix     = []byte{0x05}
	ShipmentTrackingPrefix    = []byte{0x06}
	PaymentInstructionPrefix  = []byte{0x07}
	TradeFinanceStatsKey      = []byte{0x08}
	
	// Counter keys
	NextPartyIDKey        = []byte{0x10}
	NextLcIDKey           = []byte{0x11}
	NextDocumentIDKey     = []byte{0x12}
	NextPolicyIDKey       = []byte{0x13}
	NextInstructionIDKey  = []byte{0x14}
	
	// Index keys
	PartyByAddressPrefix  = []byte{0x20}
	LcByPartyPrefix       = []byte{0x21}
	DocumentByLcPrefix    = []byte{0x22}
	PolicyByLcPrefix      = []byte{0x23}
	PaymentByLcPrefix     = []byte{0x24}
)

// GetTradePartyKey returns the store key for a trade party
func GetTradePartyKey(partyID string) []byte {
	return append(TradePartyPrefix, []byte(partyID)...)
}

// GetLetterOfCreditKey returns the store key for an LC
func GetLetterOfCreditKey(lcID string) []byte {
	return append(LetterOfCreditPrefix, []byte(lcID)...)
}

// GetTradeDocumentKey returns the store key for a document
func GetTradeDocumentKey(documentID string) []byte {
	return append(TradeDocumentPrefix, []byte(documentID)...)
}

// GetInsurancePolicyKey returns the store key for an insurance policy
func GetInsurancePolicyKey(policyID string) []byte {
	return append(InsurancePolicyPrefix, []byte(policyID)...)
}

// GetShipmentTrackingKey returns the store key for shipment tracking
func GetShipmentTrackingKey(trackingID string) []byte {
	return append(ShipmentTrackingPrefix, []byte(trackingID)...)
}

// GetPaymentInstructionKey returns the store key for a payment instruction
func GetPaymentInstructionKey(instructionID string) []byte {
	return append(PaymentInstructionPrefix, []byte(instructionID)...)
}

// GetPartyByAddressKey returns the index key for party by address
func GetPartyByAddressKey(address string) []byte {
	return append(PartyByAddressPrefix, []byte(address)...)
}

// GetLcByPartyKey returns the index key for LC by party
func GetLcByPartyKey(partyID, lcID string) []byte {
	return append(append(LcByPartyPrefix, []byte(partyID)...), []byte(lcID)...)
}

// GetDocumentByLcKey returns the index key for documents by LC
func GetDocumentByLcKey(lcID, documentID string) []byte {
	return append(append(DocumentByLcPrefix, []byte(lcID)...), []byte(documentID)...)
}