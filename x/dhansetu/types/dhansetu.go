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

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DhanPataAddress represents a virtual address in the DhanSetu ecosystem
type DhanPataAddress struct {
	Name           string    `json:"name" yaml:"name"`                       // username@dhan
	Owner          string    `json:"owner" yaml:"owner"`                     // blockchain address
	BlockchainAddr string    `json:"blockchain_addr" yaml:"blockchain_addr"` // mapped address
	AddressType    string    `json:"address_type" yaml:"address_type"`       // personal, business, service
	Metadata       *DhanPataMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" yaml:"updated_at"`
	IsActive       bool      `json:"is_active" yaml:"is_active"`
}

// DhanPataMetadata contains additional information for DhanPata addresses
type DhanPataMetadata struct {
	DisplayName     string            `json:"display_name" yaml:"display_name"`
	Description     string            `json:"description" yaml:"description"`
	ProfileImageURL string            `json:"profile_image_url" yaml:"profile_image_url"`
	QRCodeData      string            `json:"qr_code_data" yaml:"qr_code_data"`
	Tags            []string          `json:"tags" yaml:"tags"`
	SocialLinks     map[string]string `json:"social_links" yaml:"social_links"`
	BusinessInfo    *BusinessInfo     `json:"business_info,omitempty" yaml:"business_info,omitempty"`
	Verified        bool              `json:"verified" yaml:"verified"`
}

// BusinessInfo contains business-specific metadata
type BusinessInfo struct {
	BusinessType    string   `json:"business_type" yaml:"business_type"`
	Category        string   `json:"category" yaml:"category"`
	Location        string   `json:"location" yaml:"location"`
	Pincode         string   `json:"pincode" yaml:"pincode"`
	GST             string   `json:"gst" yaml:"gst"`
	Licenses        []string `json:"licenses" yaml:"licenses"`
	OperatingHours  string   `json:"operating_hours" yaml:"operating_hours"`
	ContactNumber   string   `json:"contact_number" yaml:"contact_number"`
}

// EnhancedMitraProfile extends the existing SevaMitra with DhanSetu features
type EnhancedMitraProfile struct {
	MitraId         string         `json:"mitra_id" yaml:"mitra_id"`
	DhanPataName    string         `json:"dhanpata_name" yaml:"dhanpata_name"`    // @dhan address
	MitraType       string         `json:"mitra_type" yaml:"mitra_type"`          // individual, business, global
	TrustScore      int64          `json:"trust_score" yaml:"trust_score"`        // 0-100
	DailyLimit      sdk.Int        `json:"daily_limit" yaml:"daily_limit"`
	MonthlyLimit    sdk.Int        `json:"monthly_limit" yaml:"monthly_limit"`
	DailyVolume     sdk.Int        `json:"daily_volume" yaml:"daily_volume"`
	MonthlyVolume   sdk.Int        `json:"monthly_volume" yaml:"monthly_volume"`
	TotalTrades     uint64         `json:"total_trades" yaml:"total_trades"`
	SuccessfulTrades uint64        `json:"successful_trades" yaml:"successful_trades"`
	ActiveEscrows   []string       `json:"active_escrows" yaml:"active_escrows"`
	Specializations []string       `json:"specializations" yaml:"specializations"` // crypto-to-fiat, bulk-orders, etc.
	OperatingRegions []string      `json:"operating_regions" yaml:"operating_regions"` // pincodes or districts
	PaymentMethods  []PaymentMethod `json:"payment_methods" yaml:"payment_methods"`
	KYCStatus       string         `json:"kyc_status" yaml:"kyc_status"`
	LastActiveAt    time.Time      `json:"last_active_at" yaml:"last_active_at"`
	CreatedAt       time.Time      `json:"created_at" yaml:"created_at"`
	IsActive        bool           `json:"is_active" yaml:"is_active"`
}

// PaymentMethod represents supported payment options
type PaymentMethod struct {
	Type        string `json:"type" yaml:"type"`               // UPI, IMPS, NEFT, Bank
	Provider    string `json:"provider" yaml:"provider"`       // GPay, PhonePe, Paytm, etc.
	Identifier  string `json:"identifier" yaml:"identifier"`   // UPI ID, Account number
	IsPreferred bool   `json:"is_preferred" yaml:"is_preferred"`
	IsVerified  bool   `json:"is_verified" yaml:"is_verified"`
}

// KshetraCoin represents a pincode-based community memecoin
type KshetraCoin struct {
	Pincode         string    `json:"pincode" yaml:"pincode"`
	CoinName        string    `json:"coin_name" yaml:"coin_name"`         // "Connaught Place Coin"
	CoinSymbol      string    `json:"coin_symbol" yaml:"coin_symbol"`     // "CP110001"
	Creator         string    `json:"creator" yaml:"creator"`             // DhanPata address
	TotalSupply     sdk.Int   `json:"total_supply" yaml:"total_supply"`
	CirculatingSupply sdk.Int `json:"circulating_supply" yaml:"circulating_supply"`
	MarketCap       sdk.Int   `json:"market_cap" yaml:"market_cap"`
	HolderCount     uint64    `json:"holder_count" yaml:"holder_count"`
	CommunityFund   sdk.Int   `json:"community_fund" yaml:"community_fund"` // 1% of trades
	NGOBeneficiary  string    `json:"ngo_beneficiary" yaml:"ngo_beneficiary"`
	Description     string    `json:"description" yaml:"description"`
	LocalLandmarks  []string  `json:"local_landmarks" yaml:"local_landmarks"`
	CreatedAt       time.Time `json:"created_at" yaml:"created_at"`
	IsActive        bool      `json:"is_active" yaml:"is_active"`
}

// CrossModuleBridge represents integration between DhanSetu and other modules
type CrossModuleBridge struct {
	BridgeId      string                 `json:"bridge_id" yaml:"bridge_id"`
	SourceModule  string                 `json:"source_module" yaml:"source_module"`  // moneyorder, namo, cultural
	TargetModule  string                 `json:"target_module" yaml:"target_module"`  // dhansetu
	SourceEntity  string                 `json:"source_entity" yaml:"source_entity"`  // order_id, token_id, etc.
	TargetEntity  string                 `json:"target_entity" yaml:"target_entity"`  // dhanpata_name, mitra_id
	BridgeType    string                 `json:"bridge_type" yaml:"bridge_type"`      // order_mapping, fee_sharing, etc.
	Metadata      map[string]interface{} `json:"metadata" yaml:"metadata"`
	CreatedAt     time.Time              `json:"created_at" yaml:"created_at"`
	IsActive      bool                   `json:"is_active" yaml:"is_active"`
}

// TradeHistoryEntry represents unified trade history across all DhanSetu products
type TradeHistoryEntry struct {
	TradeId       string                 `json:"trade_id" yaml:"trade_id"`
	UserDhanPata  string                 `json:"user_dhanpata" yaml:"user_dhanpata"`
	TradeType     string                 `json:"trade_type" yaml:"trade_type"`     // money_order, token_swap, memecoin, etc.
	SourceProduct string                 `json:"source_product" yaml:"source_product"` // moneyorder, sikkebaaz, kshetra
	Amount        sdk.Coin               `json:"amount" yaml:"amount"`
	Fee           sdk.Coin               `json:"fee" yaml:"fee"`
	Counterparty  string                 `json:"counterparty" yaml:"counterparty"` // other DhanPata or address
	Status        string                 `json:"status" yaml:"status"`
	Metadata      map[string]interface{} `json:"metadata" yaml:"metadata"`
	Timestamp     time.Time              `json:"timestamp" yaml:"timestamp"`
}

// DhanSetuFeeSummary tracks fee distribution across the ecosystem
type DhanSetuFeeSummary struct {
	Period        string   `json:"period" yaml:"period"`         // daily, weekly, monthly
	TotalFees     sdk.Coin `json:"total_fees" yaml:"total_fees"`
	PlatformShare sdk.Coin `json:"platform_share" yaml:"platform_share"` // 40% to platform
	NGOShare      sdk.Coin `json:"ngo_share" yaml:"ngo_share"`           // 40% to charity
	FounderShare  sdk.Coin `json:"founder_share" yaml:"founder_share"`   // 20% to founders
	CreatedAt     time.Time `json:"created_at" yaml:"created_at"`
}

// Validation functions

// ValidateDhanPataName validates DhanPata address format
func ValidateDhanPataName(name string) error {
	// Basic validation - implement comprehensive validation
	if len(name) < DhanPataMinLength || len(name) > DhanPataMaxLength {
		return ErrInvalidDhanPataName
	}
	// Add more validation logic here
	return nil
}

// ValidatePincode validates Indian PIN code format
func ValidatePincode(pincode string) error {
	if len(pincode) != 6 {
		return ErrInvalidPincode
	}
	// Add PIN code format validation
	return nil
}

// CalculateMitraLimits returns daily/monthly limits based on mitra type
func CalculateMitraLimits(mitraType string, trustScore int64) (daily, monthly sdk.Int) {
	baseDaily := sdk.ZeroInt()
	
	switch mitraType {
	case MitraTypeIndividual:
		baseDaily = sdk.NewInt(MitraIndividualLimit)
	case MitraTypeBusiness:
		baseDaily = sdk.NewInt(MitraBusinessLimit)
	case MitraTypeGlobal:
		return sdk.ZeroInt(), sdk.ZeroInt() // No limits for global mitras
	}
	
	// Apply trust score multiplier (50-150% of base limit)
	multiplier := sdk.NewDecWithPrec(int64(50+trustScore), 2) // 0.50 to 1.50
	adjustedDaily := multiplier.MulInt(baseDaily).TruncateInt()
	adjustedMonthly := adjustedDaily.MulRaw(30) // 30 days
	
	return adjustedDaily, adjustedMonthly
}