package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/tax/types"
)

// NAMOBurnManager handles the 2% NAMO burn mechanism
type NAMOBurnManager struct {
	keeper *Keeper
}

// NewNAMOBurnManager creates a new NAMO burn manager
func NewNAMOBurnManager(k *Keeper) *NAMOBurnManager {
	return &NAMOBurnManager{
		keeper: k,
	}
}

// BurnFromDistribution burns 2% NAMO from tax/platform distributions
func (nbm *NAMOBurnManager) BurnFromDistribution(
	ctx sdk.Context,
	amounts map[string]sdk.Coin,
) error {
	// Get the NAMO burn amount from distribution
	burnAmount, exists := amounts[types.NAMOBurnPoolName]
	if !exists || burnAmount.Amount.IsZero() {
		return nil // No burn amount specified
	}
	
	// Ensure it's NAMO token
	if burnAmount.Denom != "namo" {
		// If not NAMO, try to swap first
		swapRouter := NewNAMOSwapRouter(nbm.keeper)
		namoAmount, err := swapRouter.SwapForNAMOFee(
			ctx,
			nbm.keeper.GetModuleAddress(types.ModuleName),
			burnAmount,
			burnAmount,
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to swap for NAMO burn: %w", err)
		}
		burnAmount = namoAmount
	}
	
	// Send to module account first
	err := nbm.keeper.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		types.NAMOBurnPoolName,
		sdk.NewCoins(burnAmount),
	)
	if err != nil {
		return fmt.Errorf("failed to send to burn pool: %w", err)
	}
	
	// Burn the coins
	err = nbm.keeper.bankKeeper.BurnCoins(ctx, types.NAMOBurnPoolName, sdk.NewCoins(burnAmount))
	if err != nil {
		return fmt.Errorf("failed to burn NAMO: %w", err)
	}
	
	// Emit burn event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNAMOBurn,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, burnAmount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)
	
	// Update total burned counter
	nbm.updateTotalBurned(ctx, burnAmount)
	
	return nil
}

// BurnDirectNAMO burns NAMO directly from a user's transaction
func (nbm *NAMOBurnManager) BurnDirectNAMO(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	amount sdk.Coin,
) error {
	// Ensure it's NAMO
	if amount.Denom != "namo" {
		return fmt.Errorf("can only burn NAMO tokens, got %s", amount.Denom)
	}
	
	// Send to burn module account
	err := nbm.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		userAddr,
		types.NAMOBurnPoolName,
		sdk.NewCoins(amount),
	)
	if err != nil {
		return fmt.Errorf("failed to send to burn pool: %w", err)
	}
	
	// Burn the coins
	err = nbm.keeper.bankKeeper.BurnCoins(ctx, types.NAMOBurnPoolName, sdk.NewCoins(amount))
	if err != nil {
		return fmt.Errorf("failed to burn NAMO: %w", err)
	}
	
	// Emit burn event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNAMOBurn,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyUser, userAddr.String()),
		),
	)
	
	// Update total burned counter
	nbm.updateTotalBurned(ctx, amount)
	
	return nil
}

// GetTotalBurned returns the total amount of NAMO burned
func (nbm *NAMOBurnManager) GetTotalBurned(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(nbm.keeper.storeKey)
	key := []byte("total_namo_burned")
	
	bz := store.Get(key)
	if bz == nil {
		return sdk.NewCoin("namo", sdk.ZeroInt())
	}
	
	var amount sdk.Int
	err := amount.Unmarshal(bz)
	if err != nil {
		return sdk.NewCoin("namo", sdk.ZeroInt())
	}
	
	return sdk.NewCoin("namo", amount)
}

// updateTotalBurned updates the total burned counter
func (nbm *NAMOBurnManager) updateTotalBurned(ctx sdk.Context, burnAmount sdk.Coin) {
	if burnAmount.Denom != "namo" {
		return
	}
	
	store := ctx.KVStore(nbm.keeper.storeKey)
	key := []byte("total_namo_burned")
	
	// Get current total
	currentTotal := nbm.GetTotalBurned(ctx)
	
	// Add new burn amount
	newTotal := currentTotal.Amount.Add(burnAmount.Amount)
	
	// Store updated total
	bz, err := newTotal.Marshal()
	if err != nil {
		return
	}
	
	store.Set(key, bz)
}

// CalculateBurnAmount calculates 2% burn amount from any revenue
func (nbm *NAMOBurnManager) CalculateBurnAmount(revenue sdk.Coin) sdk.Coin {
	// 2% burn rate
	burnRate := sdk.NewDecWithPrec(2, 2) // 0.02
	burnAmount := sdk.NewDecFromInt(revenue.Amount).Mul(burnRate).TruncateInt()
	
	return sdk.NewCoin("namo", burnAmount)
}

// BurnFromPlatformRevenue burns 2% from any platform revenue source
func (nbm *NAMOBurnManager) BurnFromPlatformRevenue(
	ctx sdk.Context,
	revenueSource string,
	revenue sdk.Coin,
) error {
	// Calculate 2% burn
	burnAmount := nbm.CalculateBurnAmount(revenue)
	
	// If revenue is not in NAMO, swap first
	if revenue.Denom != "namo" {
		swapRouter := NewNAMOSwapRouter(nbm.keeper)
		namoAmount, err := swapRouter.SwapForNAMOFee(
			ctx,
			nbm.keeper.GetModuleAddress(types.ModuleName),
			burnAmount,
			revenue,
			false,
		)
		if err != nil {
			return fmt.Errorf("failed to swap for NAMO burn: %w", err)
		}
		burnAmount = namoAmount
	}
	
	// Burn the NAMO
	err := nbm.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnAmount))
	if err != nil {
		return fmt.Errorf("failed to burn NAMO from %s revenue: %w", revenueSource, err)
	}
	
	// Emit burn event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNAMOBurn,
			sdk.NewAttribute(types.AttributeKeyBurnAmount, burnAmount.String()),
			sdk.NewAttribute("revenue_source", revenueSource),
		),
	)
	
	// Update total burned counter
	nbm.updateTotalBurned(ctx, burnAmount)
	
	return nil
}

// GetBurnRate returns the current burn rate (2%)
func (nbm *NAMOBurnManager) GetBurnRate() sdk.Dec {
	return sdk.NewDecWithPrec(2, 2) // 0.02 = 2%
}

// ValidateBurnAmount validates that burn amount is correct 2% of revenue
func (nbm *NAMOBurnManager) ValidateBurnAmount(revenue sdk.Coin, burnAmount sdk.Coin) error {
	expectedBurn := nbm.CalculateBurnAmount(revenue)
	
	if !burnAmount.Equal(expectedBurn) {
		return fmt.Errorf("invalid burn amount: expected %s, got %s", expectedBurn, burnAmount)
	}
	
	return nil
}