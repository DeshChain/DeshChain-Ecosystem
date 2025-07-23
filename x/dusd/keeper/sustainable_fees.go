package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/dusd/types"
)

// SustainableDUSDFeeStructure defines the sustainable fee structure for DUSD
type SustainableDUSDFeeStructure struct {
	BaseFeeUSD sdk.Dec // $0.10 minimum
	MaxFeeUSD  sdk.Dec // $1.00 maximum
	FeeRate    sdk.Dec // 0.30% average for institutions
	VolumeTiers []DUSDVolumeTier
}

// DUSDVolumeTier represents a volume-based fee tier
type DUSDVolumeTier struct {
	MinVolume  sdk.Dec // Minimum monthly volume in USD
	FeeRate    sdk.Dec // Fee rate for this tier
	Description string
}

// GetDefaultSustainableDUSDFeeStructure returns sustainable DUSD fee structure
func GetDefaultSustainableDUSDFeeStructure() SustainableDUSDFeeStructure {
	return SustainableDUSDFeeStructure{
		BaseFeeUSD: sdk.NewDecWithPrec(10, 2),  // $0.10
		MaxFeeUSD:  sdk.NewDecWithPrec(100, 2), // $1.00
		FeeRate:    sdk.NewDecWithPrec(30, 4),  // 0.30% default
		VolumeTiers: []DUSDVolumeTier{
			{
				MinVolume:   sdk.NewDec(0),
				FeeRate:     sdk.NewDecWithPrec(30, 4), // 0.30% retail
				Description: "Retail tier",
			},
			{
				MinVolume:   sdk.NewDec(100000),        // $100K monthly
				FeeRate:     sdk.NewDecWithPrec(25, 4), // 0.25%
				Description: "Small business tier",
			},
			{
				MinVolume:   sdk.NewDec(1000000),       // $1M monthly
				FeeRate:     sdk.NewDecWithPrec(20, 4), // 0.20%
				Description: "Enterprise tier",
			},
			{
				MinVolume:   sdk.NewDec(10000000),      // $10M monthly
				FeeRate:     sdk.NewDecWithPrec(15, 4), // 0.15%
				Description: "Institutional tier",
			},
			{
				MinVolume:   sdk.NewDec(100000000),     // $100M monthly
				FeeRate:     sdk.NewDecWithPrec(10, 4), // 0.10%
				Description: "Market maker tier",
			},
		},
	}
}

// CalculateSustainableDUSDFee calculates DUSD fee with sustainable model
func (k Keeper) CalculateSustainableDUSDFee(ctx sdk.Context, amount sdk.Dec, monthlyVolume sdk.Dec) sdk.Coin {
	feeStructure := k.GetSustainableDUSDFeeStructure(ctx)
	
	// Find applicable tier based on monthly volume
	var applicableRate sdk.Dec
	for i := len(feeStructure.VolumeTiers) - 1; i >= 0; i-- {
		tier := feeStructure.VolumeTiers[i]
		if monthlyVolume.GTE(tier.MinVolume) {
			applicableRate = tier.FeeRate
			break
		}
	}
	
	// Calculate percentage-based fee
	percentageFee := amount.Mul(applicableRate)
	
	// Apply minimum fee
	if percentageFee.LT(feeStructure.BaseFeeUSD) {
		percentageFee = feeStructure.BaseFeeUSD
	}
	
	// Apply maximum fee
	if percentageFee.GT(feeStructure.MaxFeeUSD) {
		percentageFee = feeStructure.MaxFeeUSD
	}
	
	// Convert to integer
	feeAmount := percentageFee.TruncateInt()
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDUSDFeeCalculated,
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("monthly_volume", monthlyVolume.String()),
			sdk.NewAttribute("fee_rate", applicableRate.String()),
			sdk.NewAttribute("fee", feeAmount.String()),
		),
	)
	
	return sdk.NewCoin(types.DUSDDenom, feeAmount)
}

// GetSustainableDUSDFeeStructure retrieves the sustainable DUSD fee structure
func (k Keeper) GetSustainableDUSDFeeStructure(ctx sdk.Context) SustainableDUSDFeeStructure {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeySustainableDUSDFeeStructure)
	
	if bz == nil {
		return GetDefaultSustainableDUSDFeeStructure()
	}
	
	var structure SustainableDUSDFeeStructure
	k.cdc.MustUnmarshal(bz, &structure)
	return structure
}

// SetSustainableDUSDFeeStructure sets the sustainable DUSD fee structure
func (k Keeper) SetSustainableDUSDFeeStructure(ctx sdk.Context, structure SustainableDUSDFeeStructure) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&structure)
	store.Set(types.KeySustainableDUSDFeeStructure, bz)
}

// CalculateCrossCurrencyFee calculates fee for cross-currency operations
func (k Keeper) CalculateCrossCurrencyFee(ctx sdk.Context, amount sdk.Dec, fromCurrency, toCurrency string) sdk.Coin {
	// Standard 0.25% for cross-currency operations
	crossCurrencyRate := sdk.NewDecWithPrec(25, 4)
	
	fee := amount.Mul(crossCurrencyRate)
	
	// Minimum $0.50 for cross-currency
	minCrossCurrencyFee := sdk.NewDecWithPrec(50, 2)
	if fee.LT(minCrossCurrencyFee) {
		fee = minCrossCurrencyFee
	}
	
	return sdk.NewCoin(types.DUSDDenom, fee.TruncateInt())
}

// GetMonthlyVolume retrieves monthly volume for an address
func (k Keeper) GetMonthlyVolume(ctx sdk.Context, address string) sdk.Dec {
	// This would track actual monthly volume
	// For now, return placeholder
	return sdk.NewDec(0)
}

// UpdateMonthlyVolume updates the monthly volume for an address
func (k Keeper) UpdateMonthlyVolume(ctx sdk.Context, address string, amount sdk.Dec) {
	currentVolume := k.GetMonthlyVolume(ctx, address)
	newVolume := currentVolume.Add(amount)
	
	// Store updated volume
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyMonthlyVolume, []byte(address)...)
	store.Set(key, k.cdc.MustMarshal(&newVolume))
}