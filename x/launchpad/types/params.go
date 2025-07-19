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

// Sikkebaaz Launchpad fee configuration - v2.0 Economic Model
const (
	// Platform fee - 5% of total funds raised (increased from 2%)
	DefaultPlatformFeeRate = "0.05" // 5%
	
	// Listing fee - 1000 NAMO tokens (increased from 100 NAMO)
	DefaultListingFeeAmount = "1000000000000" // 1000 NAMO in base units
	
	// Anti-dump protection minimum lock period (30 days)
	DefaultMinLockPeriod = 2592000 // 30 days in seconds
	
	// Maximum lock period (2 years)
	DefaultMaxLockPeriod = 63072000 // 2 years in seconds
	
	// Minimum raise amount (10,000 NAMO)
	DefaultMinRaiseAmount = "10000000000000" // 10,000 NAMO in base units
	
	// Maximum raise amount (100 million NAMO)
	DefaultMaxRaiseAmount = "100000000000000000000" // 100M NAMO in base units
)

// Sikkebaaz platform fee distribution (of the 5% platform fee)
const (
	// Validators get 40% of platform fees
	LaunchpadValidatorFeeShare = "0.40" // 40%
	
	// Development fund gets 25% of platform fees
	LaunchpadDevelopmentFeeShare = "0.25" // 25%
	
	// Operations get 15% of platform fees
	LaunchpadOperationsFeeShare = "0.15" // 15%
	
	// NGO donations get 10% of platform fees
	LaunchpadNGOFeeShare = "0.10" // 10%
	
	// Community rewards get 5% of platform fees
	LaunchpadCommunityFeeShare = "0.05" // 5%
	
	// Founder royalty gets 5% of platform fees
	LaunchpadFounderFeeShare = "0.05" // 5%
)

// Anti-pump & dump protection parameters
const (
	// Gradual release percentage per month (10% per month after lock period)
	DefaultGradualReleaseRate = "0.10" // 10%
	
	// Maximum single sell amount (5% of total tokens)
	DefaultMaxSellPercentage = "0.05" // 5%
	
	// Whale protection threshold (1% of total supply)
	DefaultWhaleThreshold = "0.01" // 1%
	
	// Community verification requirement threshold (1 million NAMO raise)
	DefaultCommunityVerificationThreshold = "1000000000000000" // 1M NAMO
)

// Cultural features for Desi memecoin launches
const (
	// Festival launch bonus (20% fee discount during festivals)
	FestivalLaunchDiscount = "0.20" // 20%
	
	// Regional project bonus (15% fee discount for India-based projects)
	RegionalProjectDiscount = "0.15" // 15%
	
	// Cultural theme bonus (10% fee discount for cultural meme projects)
	CulturalThemeDiscount = "0.10" // 10%
	
	// Community-backed project bonus (25% fee discount for high community support)
	CommunityBackedDiscount = "0.25" // 25%
)

// SikkebaazParams defines the parameters for the Sikkebaaz launchpad
type SikkebaazParams struct {
	// Fee configuration
	PlatformFeeRate     sdk.Dec `json:"platform_fee_rate" yaml:"platform_fee_rate"`
	ListingFeeAmount    sdk.Int `json:"listing_fee_amount" yaml:"listing_fee_amount"`
	
	// Platform fee distribution
	ValidatorFeeShare   sdk.Dec `json:"validator_fee_share" yaml:"validator_fee_share"`
	DevelopmentFeeShare sdk.Dec `json:"development_fee_share" yaml:"development_fee_share"`
	OperationsFeeShare  sdk.Dec `json:"operations_fee_share" yaml:"operations_fee_share"`
	NGOFeeShare         sdk.Dec `json:"ngo_fee_share" yaml:"ngo_fee_share"`
	CommunityFeeShare   sdk.Dec `json:"community_fee_share" yaml:"community_fee_share"`
	FounderFeeShare     sdk.Dec `json:"founder_fee_share" yaml:"founder_fee_share"`
	
	// Launch constraints
	MinRaiseAmount      sdk.Int `json:"min_raise_amount" yaml:"min_raise_amount"`
	MaxRaiseAmount      sdk.Int `json:"max_raise_amount" yaml:"max_raise_amount"`
	MinLockPeriod       uint64  `json:"min_lock_period" yaml:"min_lock_period"`
	MaxLockPeriod       uint64  `json:"max_lock_period" yaml:"max_lock_period"`
	
	// Anti-dump protection
	GradualReleaseRate           sdk.Dec `json:"gradual_release_rate" yaml:"gradual_release_rate"`
	MaxSellPercentage            sdk.Dec `json:"max_sell_percentage" yaml:"max_sell_percentage"`
	WhaleThreshold               sdk.Dec `json:"whale_threshold" yaml:"whale_threshold"`
	CommunityVerificationThreshold sdk.Int `json:"community_verification_threshold" yaml:"community_verification_threshold"`
	
	// Cultural features
	FestivalDiscount         sdk.Dec `json:"festival_discount" yaml:"festival_discount"`
	RegionalDiscount         sdk.Dec `json:"regional_discount" yaml:"regional_discount"`
	CulturalThemeDiscount    sdk.Dec `json:"cultural_theme_discount" yaml:"cultural_theme_discount"`
	CommunityBackedDiscount  sdk.Dec `json:"community_backed_discount" yaml:"community_backed_discount"`
	
	// Feature toggles
	EnableAntiDumpProtection  bool `json:"enable_anti_dump_protection" yaml:"enable_anti_dump_protection"`
	EnableCulturalFeatures    bool `json:"enable_cultural_features" yaml:"enable_cultural_features"`
	EnableCommunityVerification bool `json:"enable_community_verification" yaml:"enable_community_verification"`
	RequireKYCForLaunchers    bool `json:"require_kyc_for_launchers" yaml:"require_kyc_for_launchers"`
}

// NewDefaultSikkebaazParams creates default Sikkebaaz parameters
func NewDefaultSikkebaazParams() SikkebaazParams {
	platformFee, _ := sdk.NewDecFromStr(DefaultPlatformFeeRate)
	listingFee, _ := sdk.NewIntFromString(DefaultListingFeeAmount)
	minRaise, _ := sdk.NewIntFromString(DefaultMinRaiseAmount)
	maxRaise, _ := sdk.NewIntFromString(DefaultMaxRaiseAmount)
	communityThreshold, _ := sdk.NewIntFromString(DefaultCommunityVerificationThreshold)
	
	validatorShare, _ := sdk.NewDecFromStr(LaunchpadValidatorFeeShare)
	developmentShare, _ := sdk.NewDecFromStr(LaunchpadDevelopmentFeeShare)
	operationsShare, _ := sdk.NewDecFromStr(LaunchpadOperationsFeeShare)
	ngoShare, _ := sdk.NewDecFromStr(LaunchpadNGOFeeShare)
	communityShare, _ := sdk.NewDecFromStr(LaunchpadCommunityFeeShare)
	founderShare, _ := sdk.NewDecFromStr(LaunchpadFounderFeeShare)
	
	gradualRelease, _ := sdk.NewDecFromStr(DefaultGradualReleaseRate)
	maxSell, _ := sdk.NewDecFromStr(DefaultMaxSellPercentage)
	whaleThreshold, _ := sdk.NewDecFromStr(DefaultWhaleThreshold)
	
	festivalDiscount, _ := sdk.NewDecFromStr(FestivalLaunchDiscount)
	regionalDiscount, _ := sdk.NewDecFromStr(RegionalProjectDiscount)
	culturalDiscount, _ := sdk.NewDecFromStr(CulturalThemeDiscount)
	communityDiscount, _ := sdk.NewDecFromStr(CommunityBackedDiscount)
	
	return SikkebaazParams{
		PlatformFeeRate:                platformFee,
		ListingFeeAmount:               listingFee,
		ValidatorFeeShare:              validatorShare,
		DevelopmentFeeShare:            developmentShare,
		OperationsFeeShare:             operationsShare,
		NGOFeeShare:                    ngoShare,
		CommunityFeeShare:              communityShare,
		FounderFeeShare:                founderShare,
		MinRaiseAmount:                 minRaise,
		MaxRaiseAmount:                 maxRaise,
		MinLockPeriod:                  DefaultMinLockPeriod,
		MaxLockPeriod:                  DefaultMaxLockPeriod,
		GradualReleaseRate:             gradualRelease,
		MaxSellPercentage:              maxSell,
		WhaleThreshold:                 whaleThreshold,
		CommunityVerificationThreshold: communityThreshold,
		FestivalDiscount:               festivalDiscount,
		RegionalDiscount:               regionalDiscount,
		CulturalThemeDiscount:          culturalDiscount,
		CommunityBackedDiscount:        communityDiscount,
		EnableAntiDumpProtection:       true,  // Anti-dump protection enabled by default
		EnableCulturalFeatures:         true,  // Cultural features enabled by default
		EnableCommunityVerification:    true,  // Community verification enabled by default
		RequireKYCForLaunchers:         true,  // KYC required for launchers
	}
}

// Validate validates the Sikkebaaz parameters
func (sp SikkebaazParams) Validate() error {
	// Validate platform fee rate
	if sp.PlatformFeeRate.IsNegative() {
		return fmt.Errorf("platform fee rate cannot be negative: %s", sp.PlatformFeeRate)
	}
	if sp.PlatformFeeRate.GT(sdk.NewDecWithPrec(10, 2)) { // Max 10%
		return fmt.Errorf("platform fee rate cannot exceed 10%%: %s", sp.PlatformFeeRate.Mul(sdk.NewDec(100)))
	}
	
	// Validate listing fee amount
	if sp.ListingFeeAmount.IsNegative() {
		return fmt.Errorf("listing fee amount cannot be negative: %s", sp.ListingFeeAmount)
	}
	
	// Validate fee distribution (must sum to 100%)
	totalFeeShare := sp.ValidatorFeeShare.Add(sp.DevelopmentFeeShare).
		Add(sp.OperationsFeeShare).Add(sp.NGOFeeShare).
		Add(sp.CommunityFeeShare).Add(sp.FounderFeeShare)
	
	if !totalFeeShare.Equal(sdk.OneDec()) {
		return fmt.Errorf("total fee distribution must equal 100%%: %s", totalFeeShare.Mul(sdk.NewDec(100)))
	}
	
	// Validate individual fee shares
	feeShares := []struct {
		name  string
		value sdk.Dec
	}{
		{"validator_fee_share", sp.ValidatorFeeShare},
		{"development_fee_share", sp.DevelopmentFeeShare},
		{"operations_fee_share", sp.OperationsFeeShare},
		{"ngo_fee_share", sp.NGOFeeShare},
		{"community_fee_share", sp.CommunityFeeShare},
		{"founder_fee_share", sp.FounderFeeShare},
	}
	
	for _, share := range feeShares {
		if share.value.IsNegative() {
			return fmt.Errorf("%s cannot be negative: %s", share.name, share.value)
		}
		if share.value.GT(sdk.OneDec()) {
			return fmt.Errorf("%s cannot exceed 100%%: %s", share.name, share.value.Mul(sdk.NewDec(100)))
		}
	}
	
	// Validate raise amounts
	if sp.MinRaiseAmount.IsNegative() {
		return fmt.Errorf("minimum raise amount cannot be negative: %s", sp.MinRaiseAmount)
	}
	if sp.MaxRaiseAmount.IsNegative() {
		return fmt.Errorf("maximum raise amount cannot be negative: %s", sp.MaxRaiseAmount)
	}
	if sp.MinRaiseAmount.GT(sp.MaxRaiseAmount) {
		return fmt.Errorf("minimum raise amount cannot exceed maximum raise amount")
	}
	
	// Validate lock periods
	if sp.MinLockPeriod > sp.MaxLockPeriod {
		return fmt.Errorf("minimum lock period cannot exceed maximum lock period")
	}
	
	// Validate anti-dump parameters
	if sp.GradualReleaseRate.IsNegative() || sp.GradualReleaseRate.GT(sdk.OneDec()) {
		return fmt.Errorf("gradual release rate must be between 0 and 100%%: %s", sp.GradualReleaseRate.Mul(sdk.NewDec(100)))
	}
	if sp.MaxSellPercentage.IsNegative() || sp.MaxSellPercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("max sell percentage must be between 0 and 100%%: %s", sp.MaxSellPercentage.Mul(sdk.NewDec(100)))
	}
	if sp.WhaleThreshold.IsNegative() || sp.WhaleThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("whale threshold must be between 0 and 100%%: %s", sp.WhaleThreshold.Mul(sdk.NewDec(100)))
	}
	
	// Validate cultural discounts
	discounts := []struct {
		name  string
		value sdk.Dec
	}{
		{"festival_discount", sp.FestivalDiscount},
		{"regional_discount", sp.RegionalDiscount},
		{"cultural_theme_discount", sp.CulturalThemeDiscount},
		{"community_backed_discount", sp.CommunityBackedDiscount},
	}
	
	for _, discount := range discounts {
		if discount.value.IsNegative() {
			return fmt.Errorf("%s cannot be negative: %s", discount.name, discount.value)
		}
		if discount.value.GT(sdk.NewDecWithPrec(50, 2)) { // Max 50% discount
			return fmt.Errorf("%s cannot exceed 50%%: %s", discount.name, discount.value.Mul(sdk.NewDec(100)))
		}
	}
	
	return nil
}

// CalculatePlatformFee calculates the platform fee for a given raise amount
func (sp SikkebaazParams) CalculatePlatformFee(raiseAmount sdk.Int, discountRate sdk.Dec) sdk.Int {
	// Base platform fee calculation
	baseFee := sp.PlatformFeeRate.MulInt(raiseAmount)
	
	// Apply any applicable discounts
	if !discountRate.IsZero() {
		discount := baseFee.Mul(discountRate)
		baseFee = baseFee.Sub(discount)
	}
	
	// Ensure fee is not negative
	if baseFee.IsNegative() {
		baseFee = sdk.ZeroDec()
	}
	
	return baseFee.TruncateInt()
}

// DistributePlatformFees distributes platform fees according to the configured percentages
func (sp SikkebaazParams) DistributePlatformFees(totalFees sdk.Coin) map[string]sdk.Coin {
	distribution := make(map[string]sdk.Coin)
	
	// Calculate each share
	distribution["validator_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.ValidatorFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution["development_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.DevelopmentFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution["operations_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.OperationsFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution["ngo_donation_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.NGOFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution["community_rewards_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.CommunityFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution["founder_royalty_pool"] = sdk.NewCoin(totalFees.Denom,
		sp.FounderFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	return distribution
}

// GetCulturalDiscount calculates the applicable cultural discount for a project launch
func (sp SikkebaazParams) GetCulturalDiscount(isFestivalPeriod bool, isIndianProject bool,
	isCulturalTheme bool, hasCommunitySupport bool) sdk.Dec {
	
	if !sp.EnableCulturalFeatures {
		return sdk.ZeroDec()
	}
	
	discount := sdk.ZeroDec()
	
	// Festival discount (during Indian festivals)
	if isFestivalPeriod {
		discount = discount.Add(sp.FestivalDiscount)
	}
	
	// Regional discount for India-based projects
	if isIndianProject {
		discount = discount.Add(sp.RegionalDiscount)
	}
	
	// Cultural theme discount
	if isCulturalTheme {
		discount = discount.Add(sp.CulturalThemeDiscount)
	}
	
	// Community-backed project discount
	if hasCommunitySupport {
		discount = discount.Add(sp.CommunityBackedDiscount)
	}
	
	// Cap maximum discount at 50%
	maxDiscount := sdk.NewDecWithPrec(50, 2) // 50%
	if discount.GT(maxDiscount) {
		discount = maxDiscount
	}
	
	return discount
}

// RequiresCommunityVerification checks if a project requires community verification
func (sp SikkebaazParams) RequiresCommunityVerification(raiseAmount sdk.Int) bool {
	if !sp.EnableCommunityVerification {
		return false
	}
	
	return raiseAmount.GTE(sp.CommunityVerificationThreshold)
}

// CalculateAntiDumpRelease calculates token release schedule for anti-dump protection
func (sp SikkebaazParams) CalculateAntiDumpRelease(totalTokens sdk.Int, monthsAfterLock uint64) sdk.Int {
	if !sp.EnableAntiDumpProtection || monthsAfterLock == 0 {
		return sdk.ZeroInt()
	}
	
	// Calculate total releasable amount (gradual release rate per month)
	releasePercentage := sp.GradualReleaseRate.Mul(sdk.NewDec(int64(monthsAfterLock)))
	
	// Cap at 100%
	if releasePercentage.GT(sdk.OneDec()) {
		releasePercentage = sdk.OneDec()
	}
	
	return releasePercentage.MulInt(totalTokens).TruncateInt()
}

// IsWhaleTransaction checks if a transaction exceeds whale threshold
func (sp SikkebaazParams) IsWhaleTransaction(sellAmount sdk.Int, totalSupply sdk.Int) bool {
	if !sp.EnableAntiDumpProtection {
		return false
	}
	
	whaleAmount := sp.WhaleThreshold.MulInt(totalSupply)
	return sellAmount.GT(whaleAmount.TruncateInt())
}

// TokenLaunch represents a token launch on Sikkebaaz
type TokenLaunch struct {
	// Basic launch information
	ProjectName     string    `json:"project_name" yaml:"project_name"`
	TokenSymbol     string    `json:"token_symbol" yaml:"token_symbol"`
	LauncherAddress string    `json:"launcher_address" yaml:"launcher_address"`
	
	// Financial details
	TargetRaise     sdk.Int   `json:"target_raise" yaml:"target_raise"`
	TokenPrice      sdk.Dec   `json:"token_price" yaml:"token_price"`
	TotalSupply     sdk.Int   `json:"total_supply" yaml:"total_supply"`
	
	// Anti-dump protection
	LockPeriod      uint64    `json:"lock_period" yaml:"lock_period"`
	GradualRelease  bool      `json:"gradual_release" yaml:"gradual_release"`
	
	// Cultural features
	IsCulturalTheme   bool    `json:"is_cultural_theme" yaml:"is_cultural_theme"`
	IsIndianProject   bool    `json:"is_indian_project" yaml:"is_indian_project"`
	CulturalCategory  string  `json:"cultural_category" yaml:"cultural_category"`
	
	// Community verification
	CommunitySupport      bool    `json:"community_support" yaml:"community_support"`
	CommunityVotes        uint64  `json:"community_votes" yaml:"community_votes"`
	RequiresVerification  bool    `json:"requires_verification" yaml:"requires_verification"`
	
	// Status
	LaunchStatus    string  `json:"launch_status" yaml:"launch_status"`
	FundsRaised     sdk.Int `json:"funds_raised" yaml:"funds_raised"`
	ParticipantCount uint64 `json:"participant_count" yaml:"participant_count"`
}

// Launch status constants
const (
	LaunchStatusPending     = "pending"
	LaunchStatusActive      = "active"
	LaunchStatusSuccessful  = "successful"
	LaunchStatusFailed      = "failed"
	LaunchStatusCancelled   = "cancelled"
)

// Cultural categories for Desi meme projects
const (
	CulturalCategoryBollywood   = "bollywood"
	CulturalCategoryFestivals   = "festivals"
	CulturalCategoryFood        = "food"
	CulturalCategoryLanguage    = "language"
	CulturalCategoryTradition   = "tradition"
	CulturalCategoryHistory     = "history"
	CulturalCategoryReligion    = "religion"
	CulturalCategoryRegional    = "regional"
)

// Module account names for Sikkebaaz fee distribution
const (
	SikkebaazValidatorPoolName     = "sikkebaaz_validator_pool"
	SikkebaazDevelopmentPoolName   = "sikkebaaz_development_pool"
	SikkebaazOperationsPoolName    = "sikkebaaz_operations_pool"
	SikkebaazNGOPoolName           = "sikkebaaz_ngo_pool"
	SikkebaazCommunityPoolName     = "sikkebaaz_community_pool"
	SikkebaazFounderRoyaltyPoolName = "sikkebaaz_founder_royalty_pool"
	SikkebaazListingFeePoolName    = "sikkebaaz_listing_fee_pool"
)

// GetCulturalCategoryBonus returns additional fee discount for specific cultural categories
func GetCulturalCategoryBonus(category string) sdk.Dec {
	switch category {
	case CulturalCategoryBollywood:
		return sdk.NewDecWithPrec(5, 2) // 5% bonus for Bollywood themes
	case CulturalCategoryFestivals:
		return sdk.NewDecWithPrec(10, 2) // 10% bonus for festival themes
	case CulturalCategoryFood:
		return sdk.NewDecWithPrec(3, 2) // 3% bonus for food themes
	case CulturalCategoryLanguage:
		return sdk.NewDecWithPrec(8, 2) // 8% bonus for language preservation
	case CulturalCategoryTradition:
		return sdk.NewDecWithPrec(12, 2) // 12% bonus for tradition preservation
	case CulturalCategoryHistory:
		return sdk.NewDecWithPrec(15, 2) // 15% bonus for historical themes
	case CulturalCategoryReligion:
		return sdk.NewDecWithPrec(7, 2) // 7% bonus for religious themes
	case CulturalCategoryRegional:
		return sdk.NewDecWithPrec(6, 2) // 6% bonus for regional themes
	default:
		return sdk.ZeroDec() // No bonus for non-cultural themes
	}
}

// ValidateTokenLaunch validates a token launch
func (tl TokenLaunch) Validate() error {
	if tl.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if tl.TokenSymbol == "" {
		return fmt.Errorf("token symbol cannot be empty")
	}
	if tl.LauncherAddress == "" {
		return fmt.Errorf("launcher address cannot be empty")
	}
	if tl.TargetRaise.IsNegative() || tl.TargetRaise.IsZero() {
		return fmt.Errorf("target raise must be positive: %s", tl.TargetRaise)
	}
	if tl.TokenPrice.IsNegative() || tl.TokenPrice.IsZero() {
		return fmt.Errorf("token price must be positive: %s", tl.TokenPrice)
	}
	if tl.TotalSupply.IsNegative() || tl.TotalSupply.IsZero() {
		return fmt.Errorf("total supply must be positive: %s", tl.TotalSupply)
	}
	
	return nil
}