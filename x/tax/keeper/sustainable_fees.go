package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

// SustainableFeeStructure defines the graduated fee structure for sustainability
type SustainableFeeStructure struct {
	BaseFeePct     sdk.Dec // Base fee percentage (0.50%)
	MinFeePct      sdk.Dec // Minimum fee percentage (0.25%)
	VolumeDiscount map[string]sdk.Dec // Volume-based discounts
}

// DefaultSustainableFeeStructure returns the default sustainable fee structure
func DefaultSustainableFeeStructure() SustainableFeeStructure {
	return SustainableFeeStructure{
		BaseFeePct: sdk.NewDecWithPrec(50, 4), // 0.50%
		MinFeePct:  sdk.NewDecWithPrec(25, 4), // 0.25% floor
		VolumeDiscount: map[string]sdk.Dec{
			"tier1": sdk.NewDecWithPrec(40, 4), // 0.40% for > ₹10 lakh monthly
			"tier2": sdk.NewDecWithPrec(30, 4), // 0.30% for > ₹1 Cr monthly
			"tier3": sdk.NewDecWithPrec(25, 4), // 0.25% for > ₹10 Cr monthly (minimum)
		},
	}
}

// CalculateSustainableFee calculates the fee based on sustainable model
func (k Keeper) CalculateSustainableFee(ctx sdk.Context, amount sdk.Coin, monthlyVolume sdk.Dec) (sdk.Coin, error) {
	feeStructure := k.GetSustainableFeeStructure(ctx)
	
	// Determine fee percentage based on volume
	var feePct sdk.Dec
	switch {
	case monthlyVolume.GTE(sdk.NewDec(100000000)): // ₹10 Cr
		feePct = feeStructure.VolumeDiscount["tier3"]
	case monthlyVolume.GTE(sdk.NewDec(10000000)): // ₹1 Cr
		feePct = feeStructure.VolumeDiscount["tier2"]
	case monthlyVolume.GTE(sdk.NewDec(1000000)): // ₹10 lakh
		feePct = feeStructure.VolumeDiscount["tier1"]
	default:
		feePct = feeStructure.BaseFeePct
	}
	
	// Ensure minimum fee percentage
	if feePct.LT(feeStructure.MinFeePct) {
		feePct = feeStructure.MinFeePct
	}
	
	// Calculate fee amount
	feeAmount := sdk.NewDecFromInt(amount.Amount).Mul(feePct).TruncateInt()
	
	// Apply minimum fee (₹1)
	minFee := sdk.NewInt(1)
	if feeAmount.LT(minFee) {
		feeAmount = minFee
	}
	
	return sdk.NewCoin(amount.Denom, feeAmount), nil
}

// GetSustainableFeeStructure retrieves the sustainable fee structure
func (k Keeper) GetSustainableFeeStructure(ctx sdk.Context) SustainableFeeStructure {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeySustainableFeeStructure)
	
	if bz == nil {
		return DefaultSustainableFeeStructure()
	}
	
	var structure SustainableFeeStructure
	k.cdc.MustUnmarshal(bz, &structure)
	return structure
}

// SetSustainableFeeStructure sets the sustainable fee structure
func (k Keeper) SetSustainableFeeStructure(ctx sdk.Context, structure SustainableFeeStructure) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&structure)
	store.Set(types.KeySustainableFeeStructure, bz)
}

// UpdateBaseTaxRateForSustainability updates the base tax rate to sustainable levels
func (k Keeper) UpdateBaseTaxRateForSustainability(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	
	// Update base tax rate to 0.50%
	params.BaseTaxRate = "0.005"
	
	// Remove artificial caps
	params.MaxTaxAmount = "10000" // ₹10,000 max
	params.MinTaxAmount = "1"     // ₹1 minimum
	
	k.SetParams(ctx, params)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSustainableFeeUpdate,
			sdk.NewAttribute("base_rate", params.BaseTaxRate),
			sdk.NewAttribute("min_fee", params.MinTaxAmount),
			sdk.NewAttribute("max_fee", params.MaxTaxAmount),
		),
	)
	
	return nil
}

// GetMonthlyVolume calculates the monthly transaction volume for an address
func (k Keeper) GetMonthlyVolume(ctx sdk.Context, address string) sdk.Dec {
	// This would be implemented with actual transaction history tracking
	// For now, return a placeholder
	return sdk.NewDec(0)
}

// ValidateSustainableFees ensures fees meet sustainability requirements
func (k Keeper) ValidateSustainableFees(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	
	// Parse base tax rate
	baseTaxRate, err := strconv.ParseFloat(params.BaseTaxRate, 64)
	if err != nil {
		return fmt.Errorf("invalid base tax rate: %w", err)
	}
	
	// Ensure minimum 0.25% fee
	if baseTaxRate < 0.0025 {
		return fmt.Errorf("base tax rate too low for sustainability: %f < 0.0025", baseTaxRate)
	}
	
	// Ensure minimum fee is at least ₹1
	minTaxAmount, err := strconv.ParseInt(params.MinTaxAmount, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid min tax amount: %w", err)
	}
	
	if minTaxAmount < 1 {
		return fmt.Errorf("minimum tax amount too low: %d < 1", minTaxAmount)
	}
	
	return nil
}