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
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tax distribution percentages (from calculated tax) - New Model
const (
	// NGO Donation Rate - 28% of calculated tax
	DefaultNGODonationRate = "0.28" // 28%
	
	// Validators Rate - 25% of calculated tax
	DefaultValidatorsRate = "0.25" // 25%
	
	// Community Rewards Rate - 18% of calculated tax
	DefaultCommunityRewardsRate = "0.18" // 18%
	
	// Tech Innovation Rate - 8% of calculated tax
	DefaultTechInnovationRate = "0.08" // 8%
	
	// Operations Rate - 6% of calculated tax
	DefaultOperationsRate = "0.06" // 6%
	
	// Founder Rate - 5% of calculated tax (consistent across all sources)
	DefaultFounderRate = "0.05" // 5%
	
	// Strategic Reserve Rate - 5% of calculated tax
	DefaultStrategicReserveRate = "0.05" // 5%
	
	// Co-Founders Rate - 3% of calculated tax
	DefaultCoFoundersRate = "0.03" // 3%
	
	// NAMO Burn Rate - 2% of calculated tax
	DefaultNAMOBurnRate = "0.02" // 2%
	
	// Module account names for NAMO burn
	NAMOBurnPoolName = "namo_burn_pool"
	FounderRoyaltyPoolName = "founder_royalty_pool"
)

// Platform revenue distribution percentages - New Model
const (
	// Development Fund - 25% of platform revenues
	PlatformDevelopmentShare = "0.25" // 25%
	
	// Community Treasury - 24% of platform revenues
	PlatformCommunityShare = "0.24" // 24%
	
	// Liquidity Provision - 18% of platform revenues
	PlatformLiquidityShare = "0.18" // 18%
	
	// NGO Donations - 10% of platform revenues
	PlatformNGOShare = "0.10" // 10%
	
	// Emergency Reserve - 8% of platform revenues
	PlatformEmergencyShare = "0.08" // 8%
	
	// Validators - 8% of platform revenues
	PlatformValidatorShare = "0.08" // 8%
	
	// Founder Royalty - 5% of platform revenues (consistent)
	PlatformFounderShare = "0.05" // 5%
	
	// NAMO Burn - 2% of platform revenues
	PlatformNAMOBurnShare = "0.02" // 2%
)

// Module account names for tax distribution - v2.0 Economic Model
const (
	// NGO and charity accounts
	NGOPoolName = "ngo_donation_pool"
	
	// Validator and network accounts
	ValidatorPoolName = "validator_pool"
	CommunityRewardsPoolName = "community_rewards_pool"
	
	// Development and operations accounts
	TechInnovationPoolName = "tech_innovation_pool"
	OperationsPoolName = "operations_pool"
	TalentAcquisitionPoolName = "talent_acquisition_pool"
	StrategicReservePoolName = "strategic_reserve_pool"
	
	// Founder and early supporter accounts
	FounderPoolName = "founder_pool"
	CoFoundersPoolName = "co_founders_pool"
	AngelInvestorsPoolName = "angel_investors_pool"
	
	// Platform accounts (for platform revenues)
	PlatformLiquidityPoolName = "platform_liquidity_pool"
	PlatformEmergencyPoolName = "platform_emergency_pool"
	DevelopmentPoolName = "development_pool"
)

// TaxDistribution represents the distribution of transaction tax - New Model
type TaxDistribution struct {
	NGODonations      sdk.Dec `json:"ngo_donations"`      // 28%
	Validators        sdk.Dec `json:"validators"`          // 25%
	CommunityRewards  sdk.Dec `json:"community_rewards"`  // 18%
	TechInnovation    sdk.Dec `json:"tech_innovation"`    // 8%
	Operations        sdk.Dec `json:"operations"`          // 6%
	Founder           sdk.Dec `json:"founder"`             // 5%
	StrategicReserve  sdk.Dec `json:"strategic_reserve"`  // 5%
	CoFounders        sdk.Dec `json:"co_founders"`         // 3%
	NAMOBurn          sdk.Dec `json:"namo_burn"`           // 2%
}

// PlatformRevenueDistribution represents the distribution of platform revenues
type PlatformRevenueDistribution struct {
	Development    sdk.Dec `json:"development"`     // 25%
	Community      sdk.Dec `json:"community"`       // 24%
	Liquidity      sdk.Dec `json:"liquidity"`       // 18%
	NGODonation    sdk.Dec `json:"ngo_donation"`    // 10%
	Emergency      sdk.Dec `json:"emergency"`       // 8%
	Validators     sdk.Dec `json:"validators"`      // 8%
	FounderRoyalty sdk.Dec `json:"founder_royalty"` // 5%
	NAMOBurn       sdk.Dec `json:"namo_burn"`       // 2%
}

// NewDefaultTaxDistribution creates default tax distribution - v2.0 Economic Model
func NewDefaultTaxDistribution() TaxDistribution {
	ngo, _ := sdk.NewDecFromStr(DefaultNGODonationRate)
	validators, _ := sdk.NewDecFromStr(DefaultValidatorsRate)
	community, _ := sdk.NewDecFromStr(DefaultCommunityRewardsRate)
	techInnovation, _ := sdk.NewDecFromStr(DefaultTechInnovationRate)
	ops, _ := sdk.NewDecFromStr(DefaultOperationsRate)
	founder, _ := sdk.NewDecFromStr(DefaultFounderRate)
	strategic, _ := sdk.NewDecFromStr(DefaultStrategicReserveRate)
	coFounders, _ := sdk.NewDecFromStr(DefaultCoFoundersRate)
	namoBurn, _ := sdk.NewDecFromStr(DefaultNAMOBurnRate)
	
	return TaxDistribution{
		NGODonations:     ngo,
		Validators:       validators,
		CommunityRewards: community,
		TechInnovation:   techInnovation,
		Operations:       ops,
		Founder:          founder,
		StrategicReserve: strategic,
		CoFounders:       coFounders,
		NAMOBurn:         namoBurn,
	}
}

// NewDefaultPlatformDistribution creates default platform revenue distribution
func NewDefaultPlatformDistribution() PlatformRevenueDistribution {
	dev, _ := sdk.NewDecFromStr(PlatformDevelopmentShare)
	community, _ := sdk.NewDecFromStr(PlatformCommunityShare)
	liquidity, _ := sdk.NewDecFromStr(PlatformLiquidityShare)
	ngo, _ := sdk.NewDecFromStr(PlatformNGOShare)
	emergency, _ := sdk.NewDecFromStr(PlatformEmergencyShare)
	validators, _ := sdk.NewDecFromStr(PlatformValidatorShare)
	founder, _ := sdk.NewDecFromStr(PlatformFounderShare)
	namoBurn, _ := sdk.NewDecFromStr(PlatformNAMOBurnShare)
	
	return PlatformRevenueDistribution{
		Development:    dev,
		Community:      community,
		Liquidity:      liquidity,
		NGODonation:    ngo,
		Emergency:      emergency,
		Validators:     validators,
		FounderRoyalty: founder,
		NAMOBurn:       namoBurn,
	}
}

// Progressive tax structure - rates determined by transaction amount
// No fixed base rate as tax is calculated based on transaction value

// Validate validates the tax distribution
func (td TaxDistribution) Validate() error {
	// Check that all rates are non-negative
	if td.NGODonations.IsNegative() || td.Validators.IsNegative() || 
		td.CommunityRewards.IsNegative() || td.TechInnovation.IsNegative() || 
		td.Operations.IsNegative() || td.StrategicReserve.IsNegative() || 
		td.Founder.IsNegative() || td.CoFounders.IsNegative() || 
		td.NAMOBurn.IsNegative() {
		return fmt.Errorf("all distribution rates must be non-negative")
	}
	
	// Calculate total rate
	total := td.NGODonations.Add(td.Validators).Add(td.CommunityRewards).
		Add(td.TechInnovation).Add(td.Operations).Add(td.StrategicReserve).
		Add(td.Founder).Add(td.CoFounders).Add(td.NAMOBurn)
	
	// Check that total equals 100%
	if !total.Equal(sdk.OneDec()) {
		return fmt.Errorf("total distribution (%s) must equal 1.0 (100%%)", total)
	}
	
	return nil
}

// Validate validates the platform revenue distribution
func (pd PlatformRevenueDistribution) Validate() error {
	// Check that all shares are non-negative
	if pd.Development.IsNegative() || pd.Community.IsNegative() || 
		pd.Liquidity.IsNegative() || pd.NGODonation.IsNegative() || 
		pd.Emergency.IsNegative() || pd.Validators.IsNegative() || 
		pd.FounderRoyalty.IsNegative() || pd.NAMOBurn.IsNegative() {
		return fmt.Errorf("all distribution shares must be non-negative")
	}
	
	// Calculate total share
	total := pd.Development.Add(pd.Community).
		Add(pd.Liquidity).Add(pd.NGODonation).
		Add(pd.Emergency).Add(pd.Validators).
		Add(pd.FounderRoyalty).Add(pd.NAMOBurn)
	
	// Check that total equals 100%
	if !total.Equal(sdk.OneDec()) {
		return fmt.Errorf("total distribution (%s) must equal 1.0 (100%%)", total)
	}
	
	return nil
}

// CalculateTaxAmounts calculates the actual amounts for each distribution category
func (td TaxDistribution) CalculateTaxAmounts(taxAmount sdk.Coin) map[string]sdk.Coin {
	amounts := make(map[string]sdk.Coin)
	
	// Calculate each distribution amount from the collected tax
	amounts[NGOPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.NGODonations.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[ValidatorPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.Validators.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[CommunityRewardsPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.CommunityRewards.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[TechInnovationPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.TechInnovation.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[OperationsPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.Operations.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[StrategicReservePoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.StrategicReserve.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[FounderPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.Founder.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[CoFoundersPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.CoFounders.MulInt(taxAmount.Amount).TruncateInt())
	
	amounts[NAMOBurnPoolName] = sdk.NewCoin(taxAmount.Denom, 
		td.NAMOBurn.MulInt(taxAmount.Amount).TruncateInt())
	
	return amounts
}

// CalculatePlatformAmounts calculates the actual amounts for platform revenue distribution
func (pd PlatformRevenueDistribution) CalculatePlatformAmounts(revenueAmount sdk.Coin) map[string]sdk.Coin {
	amounts := make(map[string]sdk.Coin)
	
	// Calculate each distribution amount
	amounts[DevelopmentPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.Development.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[CommunityRewardsPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.Community.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[PlatformLiquidityPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.Liquidity.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[NGOPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.NGODonation.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[PlatformEmergencyPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.Emergency.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[ValidatorPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.Validators.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[FounderRoyaltyPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.FounderRoyalty.MulInt(revenueAmount.Amount).TruncateInt())
	
	amounts[NAMOBurnPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.NAMOBurn.MulInt(revenueAmount.Amount).TruncateInt())
	
	return amounts
}

// GetTotalDistribution returns the total distribution (should be 100%)
func (td TaxDistribution) GetTotalDistribution() sdk.Dec {
	return td.NGODonations.Add(td.Validators).Add(td.CommunityRewards).
		Add(td.TechInnovation).Add(td.Operations).Add(td.StrategicReserve).
		Add(td.Founder).Add(td.CoFounders).Add(td.NAMOBurn)
}

// GetFounderTotalShare returns the total founder share from both tax and platform revenues
// This includes 5% from transaction tax + 5% from platform revenues
func GetFounderTotalShare(taxAmount sdk.Coin, platformRevenue sdk.Coin) sdk.Coin {
	taxDist := NewDefaultTaxDistribution()
	platformDist := NewDefaultPlatformDistribution()
	
	// Calculate founder royalty from transaction tax (5%)
	taxRoyalty := sdk.NewCoin(taxAmount.Denom, 
		taxDist.Founder.MulInt(taxAmount.Amount).TruncateInt())
	
	// Calculate founder royalty from platform revenue (5%)
	platformRoyalty := sdk.NewCoin(platformRevenue.Denom, 
		platformDist.FounderRoyalty.MulInt(platformRevenue.Amount).TruncateInt())
	
	// Return total
	return taxRoyalty.Add(platformRoyalty)
}

// GetNGOTotalShare returns the total NGO share from both tax and platform revenues
// This includes 28% from transaction tax + 10% from platform revenues
func GetNGOTotalShare(taxAmount sdk.Coin, platformRevenue sdk.Coin) sdk.Coin {
	taxDist := NewDefaultTaxDistribution()
	platformDist := NewDefaultPlatformDistribution()
	
	// Calculate NGO donation from transaction tax (28%)
	taxDonation := sdk.NewCoin(taxAmount.Denom, 
		taxDist.NGODonations.MulInt(taxAmount.Amount).TruncateInt())
	
	// Calculate NGO donation from platform revenue (10%)
	platformDonation := sdk.NewCoin(platformRevenue.Denom, 
		platformDist.NGODonation.MulInt(platformRevenue.Amount).TruncateInt())
	
	// Return total
	return taxDonation.Add(platformDonation)
}

// FounderRoyaltyInfo represents the founder royalty configuration
type FounderRoyaltyInfo struct {
	// Transaction tax royalty rate (0.10% of transaction)
	TransactionTaxRate sdk.Dec `json:"transaction_tax_rate"`
	
	// Platform revenue royalty rate (5% of platform revenues)
	PlatformRevenueRate sdk.Dec `json:"platform_revenue_rate"`
	
	// Inheritable flag - royalty can be passed to heirs
	Inheritable bool `json:"inheritable"`
	
	// Current beneficiary address
	BeneficiaryAddress string `json:"beneficiary_address"`
	
	// Backup beneficiary addresses (for inheritance)
	BackupBeneficiaries []string `json:"backup_beneficiaries"`
}

// NewDefaultFounderRoyaltyInfo creates default founder royalty configuration
func NewDefaultFounderRoyaltyInfo(beneficiary string) FounderRoyaltyInfo {
	txRate, _ := sdk.NewDecFromStr(DefaultFounderRate)
	platformRate, _ := sdk.NewDecFromStr(PlatformFounderShare)
	
	return FounderRoyaltyInfo{
		TransactionTaxRate:  txRate,
		PlatformRevenueRate: platformRate,
		Inheritable:         true, // Royalty is inheritable by default
		BeneficiaryAddress:  beneficiary,
		BackupBeneficiaries: []string{}, // To be set by founder
	}
}