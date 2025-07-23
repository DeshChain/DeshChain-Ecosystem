package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/dinr/types"
)

// TieredFeeStructure defines the tiered fee structure for DINR operations
type TieredFeeStructure struct {
	Tiers []FeeTier
}

// FeeTier represents a single fee tier
type FeeTier struct {
	MinAmount sdk.Int  // Minimum amount for this tier
	MaxAmount sdk.Int  // Maximum amount for this tier (empty for no limit)
	FeeRate   sdk.Dec  // Fee rate as decimal (e.g., 0.001 for 0.1%)
	MinFee    sdk.Int  // Minimum fee for this tier
}

// GetDefaultTieredFeeStructure returns the default tiered fee structure
func GetDefaultTieredFeeStructure() TieredFeeStructure {
	return TieredFeeStructure{
		Tiers: []FeeTier{
			{
				MinAmount: sdk.NewInt(0),
				MaxAmount: sdk.NewInt(10000),      // Up to ₹10,000
				FeeRate:   sdk.NewDecWithPrec(10, 4), // 0.10%
				MinFee:    sdk.NewInt(10),          // Min ₹10
			},
			{
				MinAmount: sdk.NewInt(10000),
				MaxAmount: sdk.NewInt(100000),     // ₹10K - ₹1L
				FeeRate:   sdk.NewDecWithPrec(8, 4),  // 0.08%
				MinFee:    sdk.NewInt(10),          // Min ₹10
			},
			{
				MinAmount: sdk.NewInt(100000),
				MaxAmount: sdk.NewInt(1000000),    // ₹1L - ₹10L
				FeeRate:   sdk.NewDecWithPrec(6, 4),  // 0.06%
				MinFee:    sdk.NewInt(80),          // Min ₹80
			},
			{
				MinAmount: sdk.NewInt(1000000),
				MaxAmount: sdk.NewInt(10000000),   // ₹10L - ₹1Cr
				FeeRate:   sdk.NewDecWithPrec(4, 4),  // 0.04%
				MinFee:    sdk.NewInt(600),         // Min ₹600
			},
			{
				MinAmount: sdk.NewInt(10000000),
				MaxAmount: sdk.Int{},               // ₹1 Cr+
				FeeRate:   sdk.NewDecWithPrec(2, 4),  // 0.02% (no cap)
				MinFee:    sdk.NewInt(4000),        // Min ₹4,000
			},
		},
	}
}

// CalculateTieredFee calculates the fee based on tiered structure
func (k Keeper) CalculateTieredFee(ctx sdk.Context, amount sdk.Int, operation string) (sdk.Int, error) {
	feeStructure := k.GetTieredFeeStructure(ctx)
	
	// Find the appropriate tier
	var applicableTier *FeeTier
	for _, tier := range feeStructure.Tiers {
		if amount.GTE(tier.MinAmount) && (tier.MaxAmount.IsNil() || amount.LT(tier.MaxAmount)) {
			applicableTier = &tier
			break
		}
	}
	
	if applicableTier == nil {
		// Use highest tier by default
		applicableTier = &feeStructure.Tiers[len(feeStructure.Tiers)-1]
	}
	
	// Calculate fee
	fee := sdk.NewDecFromInt(amount).Mul(applicableTier.FeeRate).TruncateInt()
	
	// Apply minimum fee
	if fee.LT(applicableTier.MinFee) {
		fee = applicableTier.MinFee
	}
	
	// No maximum cap for sustainability
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTieredFeeCalculated,
			sdk.NewAttribute("operation", operation),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("fee_rate", applicableTier.FeeRate.String()),
			sdk.NewAttribute("fee", fee.String()),
		),
	)
	
	return fee, nil
}

// GetTieredFeeStructure retrieves the tiered fee structure
func (k Keeper) GetTieredFeeStructure(ctx sdk.Context) TieredFeeStructure {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyTieredFeeStructure)
	
	if bz == nil {
		return GetDefaultTieredFeeStructure()
	}
	
	var structure TieredFeeStructure
	k.cdc.MustUnmarshal(bz, &structure)
	return structure
}

// SetTieredFeeStructure sets the tiered fee structure
func (k Keeper) SetTieredFeeStructure(ctx sdk.Context, structure TieredFeeStructure) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&structure)
	store.Set(types.KeyTieredFeeStructure, bz)
}

// CalculateMintFeeTiered calculates the mint fee using tiered structure
func (k Keeper) CalculateMintFeeTiered(ctx sdk.Context, amount sdk.Int) (sdk.Int, error) {
	return k.CalculateTieredFee(ctx, amount, "mint")
}

// CalculateBurnFeeTiered calculates the burn fee using tiered structure
func (k Keeper) CalculateBurnFeeTiered(ctx sdk.Context, amount sdk.Int) (sdk.Int, error) {
	return k.CalculateTieredFee(ctx, amount, "burn")
}

// GetFeeInfo returns fee information for a given amount
func (k Keeper) GetFeeInfo(ctx sdk.Context, amount sdk.Int) types.FeeInfo {
	feeStructure := k.GetTieredFeeStructure(ctx)
	
	// Find applicable tier
	var applicableTier *FeeTier
	var tierIndex int
	for i, tier := range feeStructure.Tiers {
		if amount.GTE(tier.MinAmount) && (tier.MaxAmount.IsNil() || amount.LT(tier.MaxAmount)) {
			applicableTier = &tier
			tierIndex = i
			break
		}
	}
	
	if applicableTier == nil {
		tierIndex = len(feeStructure.Tiers) - 1
		applicableTier = &feeStructure.Tiers[tierIndex]
	}
	
	// Calculate fee
	fee := sdk.NewDecFromInt(amount).Mul(applicableTier.FeeRate).TruncateInt()
	if fee.LT(applicableTier.MinFee) {
		fee = applicableTier.MinFee
	}
	
	return types.FeeInfo{
		Amount:     amount,
		FeeRate:    applicableTier.FeeRate,
		Fee:        fee,
		TierIndex:  tierIndex,
		MinFee:     applicableTier.MinFee,
		HasCap:     false, // No cap for sustainability
	}
}