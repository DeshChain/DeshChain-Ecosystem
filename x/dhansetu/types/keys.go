/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

const (
	// ModuleName defines the module name
	ModuleName = "dhansetu"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dhansetu"
	
	// DefaultDenom is the default denomination for DhanSetu
	DefaultDenom = "namo"
)

// Store key prefixes for DhanSetu integration
var (
	// DhanPata Virtual Address System
	KeyPrefixDhanPataAddress     = []byte{0x01} // name@dhan -> blockchain address mapping
	KeyPrefixAddressToDhanPata   = []byte{0x02} // reverse mapping
	KeyPrefixDhanPataMetadata    = []byte{0x03} // profile info, QR codes, etc.
	
	// Mitra Exchange Integration
	KeyPrefixMitraProfile        = []byte{0x10} // Enhanced SevaMitra profiles
	KeyPrefixMitraRating         = []byte{0x11} // Trust scores and ratings
	KeyPrefixMitraLimits         = []byte{0x12} // Daily/monthly limits
	KeyPrefixMitraEscrow         = []byte{0x13} // Active escrows
	
	// Money Order DEX Bridge
	KeyPrefixOrderBridge         = []byte{0x20} // Money order -> DhanSetu mapping
	KeyPrefixTradeHistory        = []byte{0x21} // Cross-platform trade history
	KeyPrefixFeeDistribution     = []byte{0x22} // DhanSetu fee sharing
	
	// Kshetra Coins (Pincode-based memecoins)
	KeyPrefixKshetraCoin         = []byte{0x30} // pincode -> coin mapping
	KeyPrefixKshetraHolder       = []byte{0x31} // coin holder data
	KeyPrefixKshetraMetrics      = []byte{0x32} // community metrics
	
	// Cultural Integration Bridge
	KeyPrefixCulturalBonus       = []byte{0x40} // Festival bonuses
	KeyPrefixLanguagePreference  = []byte{0x41} // User language settings
	KeyPrefixRegionalOffers      = []byte{0x42} // Location-based offers
	
	// Cross-Module Integration
	KeyPrefixSurakshaIntegration = []byte{0x50} // Pension system hooks
	KeyPrefixSikkebaazBridge     = []byte{0x51} // Memecoin platform bridge
	KeyPrefixLendingIntegration  = []byte{0x52} // Krishi/Vyavasaya integration
)

// DhanPata address formats
const (
	DhanPataPersonalPrefix  = "@dhan"
	DhanPataBusinessPrefix  = ".biz@dhan"
	DhanPataServicePrefix   = ".svc@dhan"
	
	// Address validation patterns
	DhanPataMinLength = 4   // a@dhan
	DhanPataMaxLength = 32  // username.category@dhan
)

// Mitra types and limits
const (
	MitraTypeIndividual = "individual"
	MitraTypeBusiness   = "business"
	MitraTypeGlobal     = "global"
	
	// Daily limits (in base NAMO units)
	MitraIndividualLimit = 100_000_000_000   // ₹1L equivalent
	MitraBusinessLimit   = 1_000_000_000_000 // ₹10L equivalent
	// Global mitras have no limit (licensed entities)
)

// Fee structure for DhanSetu ecosystem
const (
	DhanSetuBaseFee        = "0.005"  // 0.5% platform fee
	DhanPataNamingFee      = "100"    // 100 NAMO for custom names
	KshetraCoinCreationFee = "1000"   // 1000 NAMO per pincode coin
	MitraRegistrationFee   = "500"    // 500 NAMO for verified mitra
)

// Event types for DhanSetu integration
const (
	EventTypeDhanPataRegistered   = "dhanpata_registered"
	EventTypeMitraTradeCompleted  = "mitra_trade_completed"
	EventTypeKshetraCoinCreated   = "kshetra_coin_created"
	EventTypeCrossModuleTransfer  = "cross_module_transfer"
	
	// Attribute keys
	AttributeKeyDhanPataName = "dhanpata_name"
	AttributeKeyMitraType    = "mitra_type"
	AttributeKeyPincode      = "pincode"
	AttributeKeySourceModule = "source_module"
	AttributeKeyTargetModule = "target_module"
)

// GetDhanPataAddressKey returns the store key for DhanPata name mapping
func GetDhanPataAddressKey(name string) []byte {
	return append(KeyPrefixDhanPataAddress, []byte(name)...)
}

// GetAddressToDhanPataKey returns the store key for reverse DhanPata mapping
func GetAddressToDhanPataKey(address string) []byte {
	return append(KeyPrefixAddressToDhanPata, []byte(address)...)
}

// GetMitraProfileKey returns the store key for Mitra profile
func GetMitraProfileKey(mitraId string) []byte {
	return append(KeyPrefixMitraProfile, []byte(mitraId)...)
}

// GetKshetraCoinKey returns the store key for Kshetra coin by pincode
func GetKshetraCoinKey(pincode string) []byte {
	return append(KeyPrefixKshetraCoin, []byte(pincode)...)
}

// GetOrderBridgeKey returns the store key for order bridge mapping
func GetOrderBridgeKey(orderId string) []byte {
	return append(KeyPrefixOrderBridge, []byte(orderId)...)
}