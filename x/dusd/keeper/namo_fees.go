package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/dusd/types"
)

// CollectDUSDFeeInNAMO collects DUSD operation fee in NAMO tokens
func (k Keeper) CollectDUSDFeeInNAMO(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	amount sdk.Dec,
	operation string,
	userPaymentToken sdk.Coin,
) error {
	// Get user's monthly volume for tiered pricing
	monthlyVolume := k.GetMonthlyVolume(ctx, userAddr.String())
	
	// Calculate fee in NAMO
	feeNAMO := k.CalculateSustainableDUSDFee(ctx, amount, monthlyVolume)
	
	// If user is not paying in NAMO, use swap router
	if userPaymentToken.Denom != "namo" && k.taxKeeper != nil {
		swapRouter := k.taxKeeper.GetNAMOSwapRouter()
		if swapRouter != nil {
			// Swap user's token for NAMO to pay fee
			swappedNAMO, err := swapRouter.SwapForNAMOFee(
				ctx,
				userAddr,
				feeNAMO,
				userPaymentToken,
				false, // not inclusive
			)
			if err != nil {
				return fmt.Errorf("failed to swap for NAMO fee: %w", err)
			}
			feeNAMO = swappedNAMO
		}
	}
	
	// Collect fee from user
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		userAddr,
		types.ModuleName,
		sdk.NewCoins(feeNAMO),
	)
	if err != nil {
		return fmt.Errorf("failed to collect fee: %w", err)
	}
	
	// Distribute collected fee
	err = k.DistributeDUSDFee(ctx, feeNAMO, operation)
	if err != nil {
		return fmt.Errorf("failed to distribute fee: %w", err)
	}
	
	// Update user's monthly volume
	k.UpdateMonthlyVolume(ctx, userAddr.String(), amount)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDUSDFeeCollected,
			sdk.NewAttribute("user", userAddr.String()),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("fee", feeNAMO.String()),
			sdk.NewAttribute("operation", operation),
		),
	)
	
	return nil
}

// DistributeDUSDFee distributes collected DUSD fees according to platform distribution
func (k Keeper) DistributeDUSDFee(ctx sdk.Context, fee sdk.Coin, operation string) error {
	// Use tax keeper's platform revenue distribution if available
	if k.taxKeeper != nil {
		return k.taxKeeper.DistributePlatformRevenue(ctx, "dusd_"+operation, fee)
	}
	
	// Fallback: send to revenue module
	if k.revenueKeeper != nil {
		return k.revenueKeeper.RecordRevenue(ctx, "dusd_"+operation, fee)
	}
	
	return nil
}

// MintDUSDWithNAMOFee mints DUSD and collects fee in NAMO
func (k Keeper) MintDUSDWithNAMOFee(
	ctx sdk.Context,
	minter sdk.AccAddress,
	amount sdk.Dec,
	collateral sdk.Coin,
) error {
	// Collect fee in NAMO
	err := k.CollectDUSDFeeInNAMO(ctx, minter, amount, "mint", collateral)
	if err != nil {
		return err
	}
	
	// Proceed with minting (existing logic)
	// For DUSD, typically backed by USD stable assets
	dusdAmount := amount.TruncateInt()
	dusdCoin := sdk.NewCoin("dusd", dusdAmount)
	
	// Lock collateral if required
	if !collateral.IsZero() {
		err = k.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			minter,
			types.ModuleName,
			sdk.NewCoins(collateral),
		)
		if err != nil {
			return err
		}
	}
	
	// Mint DUSD
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(dusdCoin))
	if err != nil {
		return err
	}
	
	// Send to user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		minter,
		sdk.NewCoins(dusdCoin),
	)
	if err != nil {
		return err
	}
	
	return nil
}

// BurnDUSDWithNAMOFee burns DUSD and collects fee in NAMO
func (k Keeper) BurnDUSDWithNAMOFee(
	ctx sdk.Context,
	burner sdk.AccAddress,
	amount sdk.Dec,
) error {
	// Check user has enough DUSD
	dusdBalance := k.bankKeeper.GetBalance(ctx, burner, "dusd")
	dusdAmount := amount.TruncateInt()
	if dusdBalance.Amount.LT(dusdAmount) {
		return fmt.Errorf("insufficient DUSD balance")
	}
	
	// Collect fee in NAMO
	userPaymentToken := k.bankKeeper.GetBalance(ctx, burner, "namo")
	if userPaymentToken.IsZero() {
		// Try other tokens
		balances := k.bankKeeper.GetAllBalances(ctx, burner)
		if !balances.IsZero() {
			userPaymentToken = balances[0] // Use first available token
		}
	}
	
	err := k.CollectDUSDFeeInNAMO(ctx, burner, amount, "burn", userPaymentToken)
	if err != nil {
		return err
	}
	
	// Burn DUSD
	dusdCoin := sdk.NewCoin("dusd", dusdAmount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		burner,
		types.ModuleName,
		sdk.NewCoins(dusdCoin),
	)
	if err != nil {
		return err
	}
	
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(dusdCoin))
	if err != nil {
		return err
	}
	
	// Release collateral if any (existing logic)
	// DUSD may have backing collateral to release
	
	return nil
}

// ConvertDUSDWithNAMOFee converts between DUSD and other currencies with NAMO fee
func (k Keeper) ConvertDUSDWithNAMOFee(
	ctx sdk.Context,
	converter sdk.AccAddress,
	fromAmount sdk.Dec,
	fromCurrency string,
	toCurrency string,
) error {
	// Calculate cross-currency fee
	fee := k.CalculateCrossCurrencyFee(ctx, fromAmount, fromCurrency, toCurrency)
	
	// Collect fee in NAMO
	userPaymentToken := k.bankKeeper.GetBalance(ctx, converter, "namo")
	if userPaymentToken.IsZero() {
		balances := k.bankKeeper.GetAllBalances(ctx, converter)
		if !balances.IsZero() {
			userPaymentToken = balances[0]
		}
	}
	
	// If user is not paying in NAMO, use swap router
	if userPaymentToken.Denom != "namo" && k.taxKeeper != nil {
		swapRouter := k.taxKeeper.GetNAMOSwapRouter()
		if swapRouter != nil {
			swappedNAMO, err := swapRouter.SwapForNAMOFee(
				ctx,
				converter,
				fee,
				userPaymentToken,
				false,
			)
			if err != nil {
				return fmt.Errorf("failed to swap for NAMO fee: %w", err)
			}
			fee = swappedNAMO
		}
	}
	
	// Collect fee
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		converter,
		types.ModuleName,
		sdk.NewCoins(fee),
	)
	if err != nil {
		return err
	}
	
	// Distribute fee
	err = k.DistributeDUSDFee(ctx, fee, "convert")
	if err != nil {
		return err
	}
	
	// Proceed with conversion logic
	// This would involve oracle rates and actual token swaps
	
	return nil
}

// GetFeeEstimate returns fee estimate in NAMO for a DUSD operation
func (k Keeper) GetFeeEstimate(ctx sdk.Context, amount sdk.Dec, operation string, userAddr sdk.AccAddress) sdk.Coin {
	monthlyVolume := sdk.NewDec(0)
	if userAddr != nil {
		monthlyVolume = k.GetMonthlyVolume(ctx, userAddr.String())
	}
	
	if operation == "convert" {
		return k.CalculateCrossCurrencyFee(ctx, amount, "dusd", "other")
	}
	
	return k.CalculateSustainableDUSDFee(ctx, amount, monthlyVolume)
}

// SetInclusiveFeeOption sets whether fees are inclusive or on-top for a user
func (k Keeper) SetInclusiveFeeOption(ctx sdk.Context, user sdk.AccAddress, inclusive bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("fee_option/"), user.Bytes()...)
	
	if inclusive {
		store.Set(key, []byte{1})
	} else {
		store.Set(key, []byte{0})
	}
}

// GetInclusiveFeeOption gets user's fee preference
func (k Keeper) GetInclusiveFeeOption(ctx sdk.Context, user sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("fee_option/"), user.Bytes()...)
	
	bz := store.Get(key)
	if bz == nil {
		return false // Default to on-top
	}
	
	return bz[0] == 1
}