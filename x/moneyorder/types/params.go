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

// Parameter store keys
var (
	KeyTradingFeeRate            = []byte("TradingFeeRate")
	KeyMoneyOrderBaseFee         = []byte("MoneyOrderBaseFee")
	KeyExpressDeliveryFee        = []byte("ExpressDeliveryFee")
	KeyBulkOrderDiscount         = []byte("BulkOrderDiscount")
	KeyMinMoneyOrderAmount       = []byte("MinMoneyOrderAmount")
	KeyMaxMoneyOrderAmount       = []byte("MaxMoneyOrderAmount")
	KeyMaxDailyUserLimit         = []byte("MaxDailyUserLimit")
	KeyKYCRequiredThreshold      = []byte("KYCRequiredThreshold")
	
	// Fee distribution keys
	KeyValidatorFeeShare         = []byte("ValidatorFeeShare")
	KeyPlatformFeeShare          = []byte("PlatformFeeShare")
	KeyLiquidityProviderFeeShare = []byte("LiquidityProviderFeeShare")
	KeyDevelopmentFeeShare       = []byte("DevelopmentFeeShare")
	KeyNGODonationFeeShare       = []byte("NGODonationFeeShare")
	KeyFounderRoyaltyFeeShare    = []byte("FounderRoyaltyFeeShare")
	
	// Cultural discount keys
	KeyFestivalDiscount          = []byte("FestivalDiscount")
	KeyVillagePoolDiscount       = []byte("VillagePoolDiscount")
	KeySeniorCitizenDiscount     = []byte("SeniorCitizenDiscount")
	KeyCulturalTokenDiscount     = []byte("CulturalTokenDiscount")
	
	// Operational keys
	KeyEnableCulturalFeatures    = []byte("EnableCulturalFeatures")
	KeyEnableFestivalBonuses     = []byte("EnableFestivalBonuses")
	KeyRequireKYCForLargeOrders  = []byte("RequireKYCForLargeOrders")
	KeyEnableVillagePools        = []byte("EnableVillagePools")
	KeyEnableFixedRatePools      = []byte("EnableFixedRatePools")
)

// MoneyOrderParams defines the parameters for the Money Order module
type MoneyOrderParams struct {
	// Trading fee configuration
	TradingFeeRate       sdk.Dec `json:"trading_fee_rate" yaml:"trading_fee_rate"`
	MoneyOrderBaseFee    sdk.Dec `json:"money_order_base_fee" yaml:"money_order_base_fee"`
	ExpressDeliveryFee   sdk.Dec `json:"express_delivery_fee" yaml:"express_delivery_fee"`
	BulkOrderDiscount    sdk.Dec `json:"bulk_order_discount" yaml:"bulk_order_discount"`
	
	// Order limits
	MinMoneyOrderAmount  sdk.Int `json:"min_money_order_amount" yaml:"min_money_order_amount"`
	MaxMoneyOrderAmount  sdk.Int `json:"max_money_order_amount" yaml:"max_money_order_amount"`
	MaxDailyUserLimit    sdk.Int `json:"max_daily_user_limit" yaml:"max_daily_user_limit"`
	KYCRequiredThreshold sdk.Int `json:"kyc_required_threshold" yaml:"kyc_required_threshold"`
	
	// Fee distribution (from x/dex)
	ValidatorFeeShare         sdk.Dec `json:"validator_fee_share" yaml:"validator_fee_share"`
	PlatformFeeShare          sdk.Dec `json:"platform_fee_share" yaml:"platform_fee_share"`
	LiquidityProviderFeeShare sdk.Dec `json:"liquidity_provider_fee_share" yaml:"liquidity_provider_fee_share"`
	DevelopmentFeeShare       sdk.Dec `json:"development_fee_share" yaml:"development_fee_share"`
	NGODonationFeeShare       sdk.Dec `json:"ngo_donation_fee_share" yaml:"ngo_donation_fee_share"`
	FounderRoyaltyFeeShare    sdk.Dec `json:"founder_royalty_fee_share" yaml:"founder_royalty_fee_share"`
	
	// Cultural discounts
	FestivalDiscount        sdk.Dec `json:"festival_discount" yaml:"festival_discount"`
	VillagePoolDiscount     sdk.Dec `json:"village_pool_discount" yaml:"village_pool_discount"`
	SeniorCitizenDiscount   sdk.Dec `json:"senior_citizen_discount" yaml:"senior_citizen_discount"`
	CulturalTokenDiscount   sdk.Dec `json:"cultural_token_discount" yaml:"cultural_token_discount"`
	
	// Operational flags
	EnableCulturalFeatures   bool `json:"enable_cultural_features" yaml:"enable_cultural_features"`
	EnableFestivalBonuses    bool `json:"enable_festival_bonuses" yaml:"enable_festival_bonuses"`
	RequireKYCForLargeOrders bool `json:"require_kyc_for_large_orders" yaml:"require_kyc_for_large_orders"`
	EnableVillagePools       bool `json:"enable_village_pools" yaml:"enable_village_pools"`
	EnableFixedRatePools     bool `json:"enable_fixed_rate_pools" yaml:"enable_fixed_rate_pools"`
}

// DefaultParams returns default Money Order parameters
func DefaultParams() MoneyOrderParams {
	// Parse decimal values
	tradingFee, _ := sdk.NewDecFromStr(DefaultTradingFeeRate)
	moneyOrderFee, _ := sdk.NewDecFromStr(MoneyOrderBaseFee)
	expressFee, _ := sdk.NewDecFromStr(ExpressDeliveryFee)
	bulkDiscount, _ := sdk.NewDecFromStr(BulkOrderDiscount)
	
	// Fee distribution (from x/dex/types/params.go)
	validatorShare, _ := sdk.NewDecFromStr("0.45")    // 45%
	platformShare, _ := sdk.NewDecFromStr("0.25")     // 25%
	liquidityShare, _ := sdk.NewDecFromStr("0.20")    // 20%
	developmentShare, _ := sdk.NewDecFromStr("0.05")  // 5%
	ngoShare, _ := sdk.NewDecFromStr("0.03")          // 3%
	founderShare, _ := sdk.NewDecFromStr("0.02")      // 2%
	
	// Cultural discounts
	festivalDisc, _ := sdk.NewDecFromStr(FestivalDiscount)
	villageDisc, _ := sdk.NewDecFromStr(VillagePoolDiscount)
	seniorDisc, _ := sdk.NewDecFromStr(SeniorCitizenDiscount)
	culturalDisc, _ := sdk.NewDecFromStr(CulturalTokenDiscount)
	
	return MoneyOrderParams{
		// Trading fees
		TradingFeeRate:       tradingFee,
		MoneyOrderBaseFee:    moneyOrderFee,
		ExpressDeliveryFee:   expressFee,
		BulkOrderDiscount:    bulkDiscount,
		
		// Order limits
		MinMoneyOrderAmount:  MinMoneyOrderAmount,
		MaxMoneyOrderAmount:  MaxMoneyOrderAmount,
		MaxDailyUserLimit:    MaxDailyUserLimit,
		KYCRequiredThreshold: KYCRequiredThreshold,
		
		// Fee distribution
		ValidatorFeeShare:         validatorShare,
		PlatformFeeShare:          platformShare,
		LiquidityProviderFeeShare: liquidityShare,
		DevelopmentFeeShare:       developmentShare,
		NGODonationFeeShare:       ngoShare,
		FounderRoyaltyFeeShare:    founderShare,
		
		// Cultural discounts
		FestivalDiscount:        festivalDisc,
		VillagePoolDiscount:     villageDisc,
		SeniorCitizenDiscount:   seniorDisc,
		CulturalTokenDiscount:   culturalDisc,
		
		// Operational flags (all enabled by default)
		EnableCulturalFeatures:   true,
		EnableFestivalBonuses:    true,
		RequireKYCForLargeOrders: true,
		EnableVillagePools:       true,
		EnableFixedRatePools:     true,
	}
}

// ParamSetPairs returns the parameter set pairs for Money Order module
func (p *MoneyOrderParams) ParamSetPairs() []ParamSetPair {
	return []ParamSetPair{
		{KeyTradingFeeRate, &p.TradingFeeRate, validateFeeRate},
		{KeyMoneyOrderBaseFee, &p.MoneyOrderBaseFee, validateFeeRate},
		{KeyExpressDeliveryFee, &p.ExpressDeliveryFee, validateFeeRate},
		{KeyBulkOrderDiscount, &p.BulkOrderDiscount, validateDiscount},
		{KeyMinMoneyOrderAmount, &p.MinMoneyOrderAmount, validateAmount},
		{KeyMaxMoneyOrderAmount, &p.MaxMoneyOrderAmount, validateAmount},
		{KeyMaxDailyUserLimit, &p.MaxDailyUserLimit, validateAmount},
		{KeyKYCRequiredThreshold, &p.KYCRequiredThreshold, validateAmount},
		{KeyValidatorFeeShare, &p.ValidatorFeeShare, validateFeeShare},
		{KeyPlatformFeeShare, &p.PlatformFeeShare, validateFeeShare},
		{KeyLiquidityProviderFeeShare, &p.LiquidityProviderFeeShare, validateFeeShare},
		{KeyDevelopmentFeeShare, &p.DevelopmentFeeShare, validateFeeShare},
		{KeyNGODonationFeeShare, &p.NGODonationFeeShare, validateFeeShare},
		{KeyFounderRoyaltyFeeShare, &p.FounderRoyaltyFeeShare, validateFeeShare},
		{KeyFestivalDiscount, &p.FestivalDiscount, validateDiscount},
		{KeyVillagePoolDiscount, &p.VillagePoolDiscount, validateDiscount},
		{KeySeniorCitizenDiscount, &p.SeniorCitizenDiscount, validateDiscount},
		{KeyCulturalTokenDiscount, &p.CulturalTokenDiscount, validateDiscount},
		{KeyEnableCulturalFeatures, &p.EnableCulturalFeatures, validateBool},
		{KeyEnableFestivalBonuses, &p.EnableFestivalBonuses, validateBool},
		{KeyRequireKYCForLargeOrders, &p.RequireKYCForLargeOrders, validateBool},
		{KeyEnableVillagePools, &p.EnableVillagePools, validateBool},
		{KeyEnableFixedRatePools, &p.EnableFixedRatePools, validateBool},
	}
}

// Validate validates the Money Order parameters
func (p MoneyOrderParams) Validate() error {
	// Validate fee rates
	if err := validateFeeRate(p.TradingFeeRate); err != nil {
		return fmt.Errorf("invalid trading fee rate: %w", err)
	}
	if err := validateFeeRate(p.MoneyOrderBaseFee); err != nil {
		return fmt.Errorf("invalid money order base fee: %w", err)
	}
	if err := validateFeeRate(p.ExpressDeliveryFee); err != nil {
		return fmt.Errorf("invalid express delivery fee: %w", err)
	}
	
	// Validate discounts
	if err := validateDiscount(p.BulkOrderDiscount); err != nil {
		return fmt.Errorf("invalid bulk order discount: %w", err)
	}
	if err := validateDiscount(p.FestivalDiscount); err != nil {
		return fmt.Errorf("invalid festival discount: %w", err)
	}
	if err := validateDiscount(p.VillagePoolDiscount); err != nil {
		return fmt.Errorf("invalid village pool discount: %w", err)
	}
	if err := validateDiscount(p.SeniorCitizenDiscount); err != nil {
		return fmt.Errorf("invalid senior citizen discount: %w", err)
	}
	if err := validateDiscount(p.CulturalTokenDiscount); err != nil {
		return fmt.Errorf("invalid cultural token discount: %w", err)
	}
	
	// Validate amounts
	if err := validateAmount(p.MinMoneyOrderAmount); err != nil {
		return fmt.Errorf("invalid min money order amount: %w", err)
	}
	if err := validateAmount(p.MaxMoneyOrderAmount); err != nil {
		return fmt.Errorf("invalid max money order amount: %w", err)
	}
	if p.MinMoneyOrderAmount.GT(p.MaxMoneyOrderAmount) {
		return fmt.Errorf("min money order amount cannot exceed max amount")
	}
	if err := validateAmount(p.MaxDailyUserLimit); err != nil {
		return fmt.Errorf("invalid max daily user limit: %w", err)
	}
	if err := validateAmount(p.KYCRequiredThreshold); err != nil {
		return fmt.Errorf("invalid KYC required threshold: %w", err)
	}
	
	// Validate fee distribution (must sum to 100%)
	totalFeeShare := p.ValidatorFeeShare.Add(p.PlatformFeeShare).
		Add(p.LiquidityProviderFeeShare).Add(p.DevelopmentFeeShare).
		Add(p.NGODonationFeeShare).Add(p.FounderRoyaltyFeeShare)
	
	if !totalFeeShare.Equal(sdk.OneDec()) {
		return fmt.Errorf("total fee distribution must equal 100%%, got %s", totalFeeShare.Mul(sdk.NewDec(100)))
	}
	
	return nil
}

// CalculateTradingFee calculates the trading fee for a given amount
func (p MoneyOrderParams) CalculateTradingFee(
	amount sdk.Int,
	orderType string,
	applicableDiscounts sdk.Dec,
) sdk.Int {
	// Start with base fee
	var baseFee sdk.Dec
	switch orderType {
	case "express":
		baseFee = p.MoneyOrderBaseFee.Add(p.ExpressDeliveryFee)
	case "bulk":
		discount := p.MoneyOrderBaseFee.Mul(p.BulkOrderDiscount)
		baseFee = p.MoneyOrderBaseFee.Sub(discount)
	default:
		baseFee = p.MoneyOrderBaseFee
	}
	
	// Apply additional discounts
	if !applicableDiscounts.IsZero() {
		discount := baseFee.Mul(applicableDiscounts)
		baseFee = baseFee.Sub(discount)
	}
	
	// Ensure fee is not negative
	if baseFee.IsNegative() {
		baseFee = sdk.ZeroDec()
	}
	
	// Calculate fee amount
	feeAmount := baseFee.MulInt(amount)
	return feeAmount.TruncateInt()
}

// DistributeFees distributes collected fees according to configured shares
func (p MoneyOrderParams) DistributeFees(totalFees sdk.Coin) map[string]sdk.Coin {
	distribution := make(map[string]sdk.Coin)
	
	distribution[ValidatorPoolName] = sdk.NewCoin(totalFees.Denom,
		p.ValidatorFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution[PlatformPoolName] = sdk.NewCoin(totalFees.Denom,
		p.PlatformFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution[LiquidityProviderPoolName] = sdk.NewCoin(totalFees.Denom,
		p.LiquidityProviderFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution[DevelopmentPoolName] = sdk.NewCoin(totalFees.Denom,
		p.DevelopmentFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution[NGODonationPoolName] = sdk.NewCoin(totalFees.Denom,
		p.NGODonationFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	distribution[FounderRoyaltyPoolName] = sdk.NewCoin(totalFees.Denom,
		p.FounderRoyaltyFeeShare.MulInt(totalFees.Amount).TruncateInt())
	
	return distribution
}

// GetCulturalDiscount calculates applicable cultural discounts
func (p MoneyOrderParams) GetCulturalDiscount(
	isFestival bool,
	isVillagePool bool,
	isSeniorCitizen bool,
	isCulturalToken bool,
) sdk.Dec {
	if !p.EnableCulturalFeatures {
		return sdk.ZeroDec()
	}
	
	discount := sdk.ZeroDec()
	
	if isFestival && p.EnableFestivalBonuses {
		discount = discount.Add(p.FestivalDiscount)
	}
	
	if isVillagePool && p.EnableVillagePools {
		discount = discount.Add(p.VillagePoolDiscount)
	}
	
	if isSeniorCitizen {
		discount = discount.Add(p.SeniorCitizenDiscount)
	}
	
	if isCulturalToken {
		discount = discount.Add(p.CulturalTokenDiscount)
	}
	
	// Cap maximum discount at 50%
	maxDiscount := sdk.NewDecWithPrec(50, 2)
	if discount.GT(maxDiscount) {
		discount = maxDiscount
	}
	
	return discount
}

// IsLargeOrder checks if an order requires KYC
func (p MoneyOrderParams) IsLargeOrder(amount sdk.Int) bool {
	return p.RequireKYCForLargeOrders && amount.GTE(p.KYCRequiredThreshold)
}

// Validation functions
func validateFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNegative() {
		return fmt.Errorf("fee rate cannot be negative: %s", v)
	}
	
	if v.GT(sdk.NewDecWithPrec(10, 2)) { // Max 10%
		return fmt.Errorf("fee rate cannot exceed 10%%: %s", v)
	}
	
	return nil
}

func validateFeeShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNegative() {
		return fmt.Errorf("fee share cannot be negative: %s", v)
	}
	
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("fee share cannot exceed 100%%: %s", v)
	}
	
	return nil
}

func validateDiscount(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNegative() {
		return fmt.Errorf("discount cannot be negative: %s", v)
	}
	
	if v.GT(sdk.NewDecWithPrec(50, 2)) { // Max 50% discount
		return fmt.Errorf("discount cannot exceed 50%%: %s", v)
	}
	
	return nil
}

func validateAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNegative() {
		return fmt.Errorf("amount cannot be negative: %s", v)
	}
	
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}