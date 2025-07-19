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
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Revenue distribution shares for all platform revenues
const (
	// Development Fund - 30% of platform revenues
	RevenueDevelopmentShare = "0.30"
	
	// Community Treasury - 25% of platform revenues
	RevenueCommunityShare = "0.25"
	
	// Liquidity Provision - 20% of platform revenues
	RevenueLiquidityShare = "0.20"
	
	// NGO Donations - 10% of platform revenues
	RevenueNGOShare = "0.10"
	
	// Emergency Reserve - 10% of platform revenues
	RevenueEmergencyShare = "0.10"
	
	// Founder Royalty - 5% of platform revenues
	RevenueFounderShare = "0.05"
)

// Module account names for revenue distribution
const (
	// DevelopmentFundPool for development activities
	DevelopmentFundPool = "revenue_development_fund"
	
	// CommunityTreasuryPool for community governance
	CommunityTreasuryPool = "revenue_community_treasury"
	
	// LiquidityProvisionPool for providing liquidity
	LiquidityProvisionPool = "revenue_liquidity_provision"
	
	// NGODonationPool for charitable donations
	NGODonationPool = "revenue_ngo_donation"
	
	// EmergencyReservePool for emergency situations
	EmergencyReservePool = "revenue_emergency_reserve"
	
	// FounderRoyaltyPool for founder's 5% share
	FounderRoyaltyPool = "revenue_founder_royalty"
)

// RevenueStream represents a source of platform revenue
type RevenueStream struct {
	// Unique identifier for the revenue stream
	ID string `json:"id"`
	
	// Name of the revenue stream (e.g., "DEX Trading Fees")
	Name string `json:"name"`
	
	// Type of revenue stream (e.g., "dex_trading")
	Type string `json:"type"`
	
	// Total revenue collected lifetime
	TotalCollected sdk.Coin `json:"total_collected"`
	
	// Revenue collected in current month
	MonthlyRevenue sdk.Coin `json:"monthly_revenue"`
	
	// Timestamp of last collection
	LastCollected time.Time `json:"last_collected"`
	
	// Active status
	Active bool `json:"active"`
}

// RevenueDistribution represents how platform revenue is distributed
type RevenueDistribution struct {
	// Development fund allocation
	Development sdk.Coin `json:"development"`
	
	// Community treasury allocation
	Community sdk.Coin `json:"community"`
	
	// Liquidity provision allocation
	Liquidity sdk.Coin `json:"liquidity"`
	
	// NGO donation allocation
	NGODonation sdk.Coin `json:"ngo_donation"`
	
	// Emergency reserve allocation
	Emergency sdk.Coin `json:"emergency"`
	
	// Founder royalty allocation
	FounderRoyalty sdk.Coin `json:"founder_royalty"`
}

// FounderRoyaltyConfig represents the founder's royalty configuration
type FounderRoyaltyConfig struct {
	// Current beneficiary address
	BeneficiaryAddress string `json:"beneficiary_address"`
	
	// Backup beneficiary addresses for inheritance
	BackupBeneficiaries []string `json:"backup_beneficiaries"`
	
	// Royalty percentage (5% by default)
	RoyaltyPercentage sdk.Dec `json:"royalty_percentage"`
	
	// Total royalties earned lifetime
	TotalEarned sdk.Coin `json:"total_earned"`
	
	// Unclaimed royalties
	UnclaimedAmount sdk.Coin `json:"unclaimed_amount"`
	
	// Last claim timestamp
	LastClaimTime time.Time `json:"last_claim_time"`
	
	// Inheritance enabled flag
	InheritanceEnabled bool `json:"inheritance_enabled"`
	
	// Lock period for inheritance (in days)
	InheritanceLockDays int64 `json:"inheritance_lock_days"`
}

// DistributionRecord represents a single revenue distribution event
type DistributionRecord struct {
	// Unique identifier
	ID uint64 `json:"id"`
	
	// Revenue stream that generated this
	StreamID string `json:"stream_id"`
	
	// Total amount distributed
	TotalAmount sdk.Coin `json:"total_amount"`
	
	// Distribution breakdown
	Distribution RevenueDistribution `json:"distribution"`
	
	// Timestamp of distribution
	Timestamp time.Time `json:"timestamp"`
	
	// Block height
	BlockHeight int64 `json:"block_height"`
	
	// Transaction hash
	TxHash string `json:"tx_hash"`
}

// NewRevenueStream creates a new revenue stream
func NewRevenueStream(id, name, streamType string) RevenueStream {
	return RevenueStream{
		ID:             id,
		Name:           name,
		Type:           streamType,
		TotalCollected: sdk.NewCoin("namo", sdk.ZeroInt()),
		MonthlyRevenue: sdk.NewCoin("namo", sdk.ZeroInt()),
		LastCollected:  time.Time{},
		Active:         true,
	}
}

// NewDefaultFounderRoyaltyConfig creates default founder royalty configuration
func NewDefaultFounderRoyaltyConfig(beneficiary string) FounderRoyaltyConfig {
	royaltyPct, _ := sdk.NewDecFromStr(RevenueFounderShare)
	
	return FounderRoyaltyConfig{
		BeneficiaryAddress:  beneficiary,
		BackupBeneficiaries: []string{},
		RoyaltyPercentage:   royaltyPct,
		TotalEarned:         sdk.NewCoin("namo", sdk.ZeroInt()),
		UnclaimedAmount:     sdk.NewCoin("namo", sdk.ZeroInt()),
		LastClaimTime:       time.Time{},
		InheritanceEnabled:  true,
		InheritanceLockDays: 90, // 90 days lock period for inheritance
	}
}

// CalculateDistribution calculates how to distribute platform revenue
func CalculateDistribution(totalRevenue sdk.Coin) (RevenueDistribution, error) {
	if totalRevenue.IsNegative() {
		return RevenueDistribution{}, fmt.Errorf("revenue amount cannot be negative")
	}
	
	// Parse distribution percentages
	devShare, _ := sdk.NewDecFromStr(RevenueDevelopmentShare)
	communityShare, _ := sdk.NewDecFromStr(RevenueCommunityShare)
	liquidityShare, _ := sdk.NewDecFromStr(RevenueLiquidityShare)
	ngoShare, _ := sdk.NewDecFromStr(RevenueNGOShare)
	emergencyShare, _ := sdk.NewDecFromStr(RevenueEmergencyShare)
	founderShare, _ := sdk.NewDecFromStr(RevenueFounderShare)
	
	// Calculate distributions
	distribution := RevenueDistribution{
		Development: sdk.NewCoin(totalRevenue.Denom, 
			devShare.MulInt(totalRevenue.Amount).TruncateInt()),
		Community: sdk.NewCoin(totalRevenue.Denom, 
			communityShare.MulInt(totalRevenue.Amount).TruncateInt()),
		Liquidity: sdk.NewCoin(totalRevenue.Denom, 
			liquidityShare.MulInt(totalRevenue.Amount).TruncateInt()),
		NGODonation: sdk.NewCoin(totalRevenue.Denom, 
			ngoShare.MulInt(totalRevenue.Amount).TruncateInt()),
		Emergency: sdk.NewCoin(totalRevenue.Denom, 
			emergencyShare.MulInt(totalRevenue.Amount).TruncateInt()),
		FounderRoyalty: sdk.NewCoin(totalRevenue.Denom, 
			founderShare.MulInt(totalRevenue.Amount).TruncateInt()),
	}
	
	// Verify total equals input (accounting for rounding)
	total := distribution.Development.
		Add(distribution.Community).
		Add(distribution.Liquidity).
		Add(distribution.NGODonation).
		Add(distribution.Emergency).
		Add(distribution.FounderRoyalty)
	
	// Handle any rounding difference by adding to development fund
	if !total.IsEqual(totalRevenue) {
		diff := totalRevenue.Sub(total)
		distribution.Development = distribution.Development.Add(diff)
	}
	
	return distribution, nil
}

// GetDistributionPools returns the mapping of distribution to pool names
func GetDistributionPools() map[string]string {
	return map[string]string{
		"development":     DevelopmentFundPool,
		"community":       CommunityTreasuryPool,
		"liquidity":       LiquidityProvisionPool,
		"ngo_donation":    NGODonationPool,
		"emergency":       EmergencyReservePool,
		"founder_royalty": FounderRoyaltyPool,
	}
}

// ValidateFounderRoyaltyConfig validates the founder royalty configuration
func (f FounderRoyaltyConfig) Validate() error {
	// Validate beneficiary address
	if _, err := sdk.AccAddressFromBech32(f.BeneficiaryAddress); err != nil {
		return fmt.Errorf("invalid beneficiary address: %w", err)
	}
	
	// Validate backup beneficiaries
	for i, backup := range f.BackupBeneficiaries {
		if _, err := sdk.AccAddressFromBech32(backup); err != nil {
			return fmt.Errorf("invalid backup beneficiary address at index %d: %w", i, err)
		}
	}
	
	// Validate royalty percentage
	if f.RoyaltyPercentage.IsNegative() || f.RoyaltyPercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("royalty percentage must be between 0 and 1")
	}
	
	// Validate inheritance lock days
	if f.InheritanceLockDays < 0 || f.InheritanceLockDays > 365 {
		return fmt.Errorf("inheritance lock days must be between 0 and 365")
	}
	
	return nil
}

// CanClaim checks if the beneficiary can claim royalties
func (f FounderRoyaltyConfig) CanClaim(currentTime time.Time) bool {
	// No minimum claim period for now, can claim anytime
	return f.UnclaimedAmount.IsPositive()
}

// AddRevenue adds revenue to a stream
func (r *RevenueStream) AddRevenue(amount sdk.Coin, timestamp time.Time) error {
	if amount.IsNegative() {
		return fmt.Errorf("revenue amount cannot be negative")
	}
	
	if amount.Denom != r.TotalCollected.Denom && !r.TotalCollected.IsZero() {
		return fmt.Errorf("revenue denom mismatch: expected %s, got %s", 
			r.TotalCollected.Denom, amount.Denom)
	}
	
	// Update total collected
	r.TotalCollected = r.TotalCollected.Add(amount)
	
	// Update monthly revenue (reset if new month)
	if r.LastCollected.Month() != timestamp.Month() || 
		r.LastCollected.Year() != timestamp.Year() {
		r.MonthlyRevenue = amount
	} else {
		r.MonthlyRevenue = r.MonthlyRevenue.Add(amount)
	}
	
	r.LastCollected = timestamp
	
	return nil
}

// GetTotalDistributed returns the sum of all distributions
func (d RevenueDistribution) GetTotalDistributed() sdk.Coin {
	return d.Development.
		Add(d.Community).
		Add(d.Liquidity).
		Add(d.NGODonation).
		Add(d.Emergency).
		Add(d.FounderRoyalty)
}

// Platform revenue sources with expected annual revenues (Year 5 projections)
var DefaultRevenueStreams = []struct {
	ID          string
	Name        string
	Type        string
	Description string
}{
	{
		ID:          "dex_001",
		Name:        "Money Order DEX Trading Fees",
		Type:        RevenueStreamDEX,
		Description: "0.3% trading fees from DEX",
	},
	{
		ID:          "nft_001",
		Name:        "Bharat Kala NFT Marketplace",
		Type:        RevenueStreamNFT,
		Description: "2.5% marketplace fees",
	},
	{
		ID:          "launch_001",
		Name:        "Sikkebaaz Launchpad Fees",
		Type:        RevenueStreamLaunchpad,
		Description: "100 NAMO + 2% of tokens",
	},
	{
		ID:          "pension_001",
		Name:        "Gram Pension Management Fees",
		Type:        RevenueStreamPension,
		Description: "80.6% profit margin on pension products",
	},
	{
		ID:          "lending_001",
		Name:        "Kisaan Mitra Interest Spread",
		Type:        RevenueStreamLending,
		Description: "Interest rate differential on agricultural loans",
	},
	{
		ID:          "privacy_001",
		Name:        "Privacy Transaction Fees",
		Type:        RevenueStreamPrivacy,
		Description: "₹50-150 per private transaction",
	},
	{
		ID:          "gov_001",
		Name:        "Governance Participation Fees",
		Type:        RevenueStreamGovernance,
		Description: "Proposal submission and voting fees",
	},
}