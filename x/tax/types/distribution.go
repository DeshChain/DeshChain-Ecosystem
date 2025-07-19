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

// Tax distribution percentages (of the 2.5% base tax) - v2.0 Economic Model
const (
	// NGO Donation Rate - 30% of base tax (0.75% of transaction)
	DefaultNGODonationRate = "0.0075" // 0.75%
	
	// Validators Rate - 25% of base tax (0.625% of transaction)
	DefaultValidatorsRate = "0.00625" // 0.625%
	
	// Community Rewards Rate - 20% of base tax (0.50% of transaction)
	DefaultCommunityRewardsRate = "0.0050" // 0.50%
	
	// Tech Innovation Rate - 6% of base tax (0.15% of transaction)
	DefaultTechInnovationRate = "0.0015" // 0.15%
	
	// Operations Rate - 5% of base tax (0.125% of transaction)
	DefaultOperationsRate = "0.00125" // 0.125%
	
	// Talent Acquisition Rate - 4% of base tax (0.10% of transaction)
	DefaultTalentAcquisitionRate = "0.0010" // 0.10%
	
	// Strategic Reserve Rate - 4% of base tax (0.10% of transaction)
	DefaultStrategicReserveRate = "0.0010" // 0.10%
	
	// Founder Rate - 3.5% of base tax (0.0875% of transaction)
	DefaultFounderRate = "0.000875" // 0.0875%
	
	// Co-Founders Rate - 1.8% of base tax (0.045% of transaction)
	DefaultCoFoundersRate = "0.00045" // 0.045%
	
	// Angel Investors Rate - 0.7% of base tax (0.0175% of transaction)
	DefaultAngelInvestorsRate = "0.000175" // 0.0175%
)

// Platform revenue distribution percentages
const (
	// Development Fund - 30% of platform revenues
	PlatformDevelopmentShare = "0.30" // 30%
	
	// Community Treasury - 25% of platform revenues
	PlatformCommunityShare = "0.25" // 25%
	
	// Liquidity Provision - 20% of platform revenues
	PlatformLiquidityShare = "0.20" // 20%
	
	// NGO Donations - 10% of platform revenues
	PlatformNGOShare = "0.10" // 10%
	
	// Emergency Reserve - 10% of platform revenues
	PlatformEmergencyShare = "0.10" // 10%
	
	// Founder Royalty - 5% of platform revenues
	PlatformFounderShare = "0.05" // 5%
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

// TaxDistribution represents the distribution of transaction tax - v2.0 Economic Model
type TaxDistribution struct {
	NGODonations      sdk.Dec `json:"ngo_donations"`      // 30%
	Validators        sdk.Dec `json:"validators"`          // 25%
	CommunityRewards  sdk.Dec `json:"community_rewards"`  // 20%
	TechInnovation    sdk.Dec `json:"tech_innovation"`    // 6%
	Operations        sdk.Dec `json:"operations"`          // 5%
	TalentAcquisition sdk.Dec `json:"talent_acquisition"` // 4%
	StrategicReserve  sdk.Dec `json:"strategic_reserve"`  // 4%
	Founder           sdk.Dec `json:"founder"`             // 3.5%
	CoFounders        sdk.Dec `json:"co_founders"`         // 1.8%
	AngelInvestors    sdk.Dec `json:"angel_investors"`     // 0.7%
}

// PlatformRevenueDistribution represents the distribution of platform revenues
type PlatformRevenueDistribution struct {
	Development    sdk.Dec `json:"development"`
	Community      sdk.Dec `json:"community"`
	Liquidity      sdk.Dec `json:"liquidity"`
	NGODonation    sdk.Dec `json:"ngo_donation"`
	Emergency      sdk.Dec `json:"emergency"`
	FounderRoyalty sdk.Dec `json:"founder_royalty"`
}

// NewDefaultTaxDistribution creates default tax distribution - v2.0 Economic Model
func NewDefaultTaxDistribution() TaxDistribution {
	ngo, _ := sdk.NewDecFromStr(DefaultNGODonationRate)
	validators, _ := sdk.NewDecFromStr(DefaultValidatorsRate)
	community, _ := sdk.NewDecFromStr(DefaultCommunityRewardsRate)
	techInnovation, _ := sdk.NewDecFromStr(DefaultTechInnovationRate)
	ops, _ := sdk.NewDecFromStr(DefaultOperationsRate)
	talent, _ := sdk.NewDecFromStr(DefaultTalentAcquisitionRate)
	strategic, _ := sdk.NewDecFromStr(DefaultStrategicReserveRate)
	founder, _ := sdk.NewDecFromStr(DefaultFounderRate)
	coFounders, _ := sdk.NewDecFromStr(DefaultCoFoundersRate)
	angels, _ := sdk.NewDecFromStr(DefaultAngelInvestorsRate)
	
	return TaxDistribution{
		NGODonations:      ngo,
		Validators:        validators,
		CommunityRewards:  community,
		TechInnovation:    techInnovation,
		Operations:        ops,
		TalentAcquisition: talent,
		StrategicReserve:  strategic,
		Founder:           founder,
		CoFounders:        coFounders,
		AngelInvestors:    angels,
	}
}

// NewDefaultPlatformDistribution creates default platform revenue distribution
func NewDefaultPlatformDistribution() PlatformRevenueDistribution {
	dev, _ := sdk.NewDecFromStr(PlatformDevelopmentShare)
	community, _ := sdk.NewDecFromStr(PlatformCommunityShare)
	liquidity, _ := sdk.NewDecFromStr(PlatformLiquidityShare)
	ngo, _ := sdk.NewDecFromStr(PlatformNGOShare)
	emergency, _ := sdk.NewDecFromStr(PlatformEmergencyShare)
	founder, _ := sdk.NewDecFromStr(PlatformFounderShare)
	
	return PlatformRevenueDistribution{
		Development:    dev,
		Community:      community,
		Liquidity:      liquidity,
		NGODonation:    ngo,
		Emergency:      emergency,
		FounderRoyalty: founder,
	}
}

// Base tax rate constant
const DefaultBaseTaxRate = "0.025" // 2.5%

// Validate validates the tax distribution
func (td TaxDistribution) Validate() error {
	// Check that all rates are non-negative
	if td.NGODonations.IsNegative() || td.Validators.IsNegative() || 
		td.CommunityRewards.IsNegative() || td.TechInnovation.IsNegative() || 
		td.Operations.IsNegative() || td.TalentAcquisition.IsNegative() ||
		td.StrategicReserve.IsNegative() || td.Founder.IsNegative() ||
		td.CoFounders.IsNegative() || td.AngelInvestors.IsNegative() {
		return fmt.Errorf("all distribution rates must be non-negative")
	}
	
	// Calculate total rate
	total := td.NGODonations.Add(td.Validators).Add(td.CommunityRewards).
		Add(td.TechInnovation).Add(td.Operations).Add(td.TalentAcquisition).
		Add(td.StrategicReserve).Add(td.Founder).Add(td.CoFounders).Add(td.AngelInvestors)
	
	// Base tax rate
	baseTax, _ := sdk.NewDecFromStr(DefaultBaseTaxRate)
	
	// Check that total equals base tax rate (2.5%)
	if !total.Equal(baseTax) {
		return fmt.Errorf("total distribution (%s) must equal base tax rate (%s)", total, baseTax)
	}
	
	return nil
}

// Validate validates the platform revenue distribution
func (pd PlatformRevenueDistribution) Validate() error {
	// Check that all shares are non-negative
	if pd.Development.IsNegative() || pd.Community.IsNegative() || 
		pd.Liquidity.IsNegative() || pd.NGODonation.IsNegative() || 
		pd.Emergency.IsNegative() || pd.FounderRoyalty.IsNegative() {
		return fmt.Errorf("all distribution shares must be non-negative")
	}
	
	// Calculate total share
	total := pd.Development.Add(pd.Community).
		Add(pd.Liquidity).Add(pd.NGODonation).
		Add(pd.Emergency).Add(pd.FounderRoyalty)
	
	// Check that total equals 100%
	if !total.Equal(sdk.OneDec()) {
		return fmt.Errorf("total distribution (%s) must equal 1.0 (100%%)", total)
	}
	
	return nil
}

// CalculateTaxAmounts calculates the actual amounts for each distribution category
func (td TaxDistribution) CalculateTaxAmounts(transactionAmount sdk.Coin) map[string]sdk.Coin {
	amounts := make(map[string]sdk.Coin)
	
	// Calculate each distribution amount
	amounts[NGOPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.NGODonations.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[ValidatorPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.Validators.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[CommunityRewardsPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.CommunityRewards.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[TechInnovationPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.TechInnovation.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[OperationsPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.Operations.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[TalentAcquisitionPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.TalentAcquisition.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[StrategicReservePoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.StrategicReserve.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[FounderPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.Founder.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[CoFoundersPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.CoFounders.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[AngelInvestorsPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.AngelInvestors.MulInt(transactionAmount.Amount).TruncateInt())
	
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
	
	amounts[FounderRoyaltyPoolName] = sdk.NewCoin(revenueAmount.Denom, 
		pd.FounderRoyalty.MulInt(revenueAmount.Amount).TruncateInt())
	
	return amounts
}

// GetTotalTaxRate returns the total tax rate (should be 2.5%)
func (td TaxDistribution) GetTotalTaxRate() sdk.Dec {
	return td.NGODonations.Add(td.Validators).Add(td.CommunityRewards).
		Add(td.TechInnovation).Add(td.Operations).Add(td.TalentAcquisition).
		Add(td.StrategicReserve).Add(td.Founder).Add(td.CoFounders).Add(td.AngelInvestors)
}

// GetFounderTotalShare returns the total founder share from both tax and platform revenues
// This includes 0.0875% from transaction tax + 5% from platform revenues
func GetFounderTotalShare(taxAmount sdk.Coin, platformRevenue sdk.Coin) sdk.Coin {
	taxDist := NewDefaultTaxDistribution()
	platformDist := NewDefaultPlatformDistribution()
	
	// Calculate founder royalty from transaction tax
	taxRoyalty := sdk.NewCoin(taxAmount.Denom, 
		taxDist.Founder.MulInt(taxAmount.Amount).TruncateInt())
	
	// Calculate founder royalty from platform revenue
	platformRoyalty := sdk.NewCoin(platformRevenue.Denom, 
		platformDist.FounderRoyalty.MulInt(platformRevenue.Amount).TruncateInt())
	
	// Return total
	return taxRoyalty.Add(platformRoyalty)
}

// GetNGOTotalShare returns the total NGO share from both tax and platform revenues
// This includes 0.75% from transaction tax + 10% from platform revenues
func GetNGOTotalShare(taxAmount sdk.Coin, platformRevenue sdk.Coin) sdk.Coin {
	taxDist := NewDefaultTaxDistribution()
	platformDist := NewDefaultPlatformDistribution()
	
	// Calculate NGO donation from transaction tax
	taxDonation := sdk.NewCoin(taxAmount.Denom, 
		taxDist.NGODonations.MulInt(taxAmount.Amount).TruncateInt())
	
	// Calculate NGO donation from platform revenue
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