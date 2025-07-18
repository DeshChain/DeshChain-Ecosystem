package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tax distribution percentages (of the 2.5% base tax)
const (
	// NGO Donation Rate - 30% of base tax (0.75% of transaction)
	DefaultNGODonationRate = "0.0075" // 0.75%
	
	// Community Rewards Rate - 20% of base tax (0.50% of transaction)
	DefaultCommunityRewardsRate = "0.0050" // 0.50%
	
	// Development Rate - 18% of base tax (0.45% of transaction)
	DefaultDevelopmentRate = "0.0045" // 0.45%
	
	// Operations Rate - 18% of base tax (0.45% of transaction)
	DefaultOperationsRate = "0.0045" // 0.45%
	
	// Token Burn Rate - 10% of base tax (0.25% of transaction)
	DefaultBurnRate = "0.0025" // 0.25%
	
	// Founder Royalty Rate - 4% of base tax (0.10% of transaction)
	DefaultFounderRoyaltyRate = "0.0010" // 0.10%
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

// Module account names for tax distribution
const (
	// NGOPoolName is the module account for NGO donations
	NGOPoolName = "ngo_donation_pool"
	
	// CommunityRewardsPoolName is the module account for community rewards
	CommunityRewardsPoolName = "community_rewards_pool"
	
	// DevelopmentPoolName is the module account for development fund
	DevelopmentPoolName = "development_pool"
	
	// OperationsPoolName is the module account for operations
	OperationsPoolName = "operations_pool"
	
	// BurnPoolName is the module account for token burning
	BurnPoolName = "burn_pool"
	
	// FounderRoyaltyPoolName is the module account for founder royalty
	FounderRoyaltyPoolName = "founder_royalty_pool"
	
	// PlatformLiquidityPoolName is the module account for platform liquidity
	PlatformLiquidityPoolName = "platform_liquidity_pool"
	
	// PlatformEmergencyPoolName is the module account for emergency reserves
	PlatformEmergencyPoolName = "platform_emergency_pool"
)

// TaxDistribution represents the distribution of transaction tax
type TaxDistribution struct {
	NGODonation      sdk.Dec `json:"ngo_donation"`
	CommunityRewards sdk.Dec `json:"community_rewards"`
	Development      sdk.Dec `json:"development"`
	Operations       sdk.Dec `json:"operations"`
	TokenBurn        sdk.Dec `json:"token_burn"`
	FounderRoyalty   sdk.Dec `json:"founder_royalty"`
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

// NewDefaultTaxDistribution creates default tax distribution
func NewDefaultTaxDistribution() TaxDistribution {
	ngo, _ := sdk.NewDecFromStr(DefaultNGODonationRate)
	community, _ := sdk.NewDecFromStr(DefaultCommunityRewardsRate)
	dev, _ := sdk.NewDecFromStr(DefaultDevelopmentRate)
	ops, _ := sdk.NewDecFromStr(DefaultOperationsRate)
	burn, _ := sdk.NewDecFromStr(DefaultBurnRate)
	founder, _ := sdk.NewDecFromStr(DefaultFounderRoyaltyRate)
	
	return TaxDistribution{
		NGODonation:      ngo,
		CommunityRewards: community,
		Development:      dev,
		Operations:       ops,
		TokenBurn:        burn,
		FounderRoyalty:   founder,
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

// Validate validates the tax distribution
func (td TaxDistribution) Validate() error {
	// Check that all rates are non-negative
	if td.NGODonation.IsNegative() || td.CommunityRewards.IsNegative() || 
		td.Development.IsNegative() || td.Operations.IsNegative() || 
		td.TokenBurn.IsNegative() || td.FounderRoyalty.IsNegative() {
		return fmt.Errorf("all distribution rates must be non-negative")
	}
	
	// Calculate total rate
	total := td.NGODonation.Add(td.CommunityRewards).
		Add(td.Development).Add(td.Operations).
		Add(td.TokenBurn).Add(td.FounderRoyalty)
	
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
		td.NGODonation.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[CommunityRewardsPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.CommunityRewards.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[DevelopmentPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.Development.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[OperationsPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.Operations.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[BurnPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.TokenBurn.MulInt(transactionAmount.Amount).TruncateInt())
	
	amounts[FounderRoyaltyPoolName] = sdk.NewCoin(transactionAmount.Denom, 
		td.FounderRoyalty.MulInt(transactionAmount.Amount).TruncateInt())
	
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
	return td.NGODonation.Add(td.CommunityRewards).
		Add(td.Development).Add(td.Operations).
		Add(td.TokenBurn).Add(td.FounderRoyalty)
}

// GetFounderTotalShare returns the total founder share from both tax and platform revenues
// This includes 0.10% from transaction tax + variable amount from platform revenues
func GetFounderTotalShare(taxAmount sdk.Coin, platformRevenue sdk.Coin) sdk.Coin {
	taxDist := NewDefaultTaxDistribution()
	platformDist := NewDefaultPlatformDistribution()
	
	// Calculate founder royalty from transaction tax
	taxRoyalty := sdk.NewCoin(taxAmount.Denom, 
		taxDist.FounderRoyalty.MulInt(taxAmount.Amount).TruncateInt())
	
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
		taxDist.NGODonation.MulInt(taxAmount.Amount).TruncateInt())
	
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
	txRate, _ := sdk.NewDecFromStr(DefaultFounderRoyaltyRate)
	platformRate, _ := sdk.NewDecFromStr(PlatformFounderShare)
	
	return FounderRoyaltyInfo{
		TransactionTaxRate:  txRate,
		PlatformRevenueRate: platformRate,
		Inheritable:         true, // Royalty is inheritable by default
		BeneficiaryAddress:  beneficiary,
		BackupBeneficiaries: []string{}, // To be set by founder
	}
}