package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/dinr/types"
	taxkeeper "github.com/deshchain/namo/x/tax/keeper"
)

// CalculateFeeInNAMO calculates the DINR operation fee and converts to NAMO
func (k Keeper) CalculateFeeInNAMO(
	ctx sdk.Context, 
	amount sdk.Int, 
	operation string,
) (sdk.Coin, error) {
	// Calculate fee in INR using tiered structure
	feeINR, err := k.CalculateTieredFee(ctx, amount, operation)
	if err != nil {
		return sdk.Coin{}, err
	}
	
	// Convert INR fee to NAMO (1 NAMO = 1 INR)
	// Fee is in micro units, so divide by 1000000
	feeNAMO := sdk.NewCoin("namo", feeINR)
	
	return feeNAMO, nil
}

// CollectDINRFeeInNAMO collects DINR operation fee in NAMO tokens
func (k Keeper) CollectDINRFeeInNAMO(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	amount sdk.Int,
	operation string,
	userPaymentToken sdk.Coin,
) error {
	// Calculate fee in NAMO
	feeNAMO, err := k.CalculateFeeInNAMO(ctx, amount, operation)
	if err != nil {
		return fmt.Errorf("failed to calculate fee: %w", err)
	}
	
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
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		userAddr,
		types.ModuleName,
		sdk.NewCoins(feeNAMO),
	)
	if err != nil {
		return fmt.Errorf("failed to collect fee: %w", err)
	}
	
	// Distribute collected fee
	err = k.DistributeDINRFee(ctx, feeNAMO, operation)
	if err != nil {
		return fmt.Errorf("failed to distribute fee: %w", err)
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDINRFeeCollected,
			sdk.NewAttribute("user", userAddr.String()),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("fee", feeNAMO.String()),
			sdk.NewAttribute("operation", operation),
		),
	)
	
	return nil
}

// DistributeDINRFee distributes collected DINR fees according to platform distribution
func (k Keeper) DistributeDINRFee(ctx sdk.Context, fee sdk.Coin, operation string) error {
	// Use tax keeper's platform revenue distribution if available
	if k.taxKeeper != nil {
		return k.taxKeeper.DistributePlatformRevenue(ctx, "dinr_"+operation, fee)
	}
	
	// Fallback: send to revenue module
	if k.revenueKeeper != nil {
		return k.revenueKeeper.RecordRevenue(ctx, "dinr_"+operation, fee)
	}
	
	return nil
}

// MintDINRWithNAMOFee mints DINR and collects fee in NAMO
func (k Keeper) MintDINRWithNAMOFee(
	ctx sdk.Context,
	minter sdk.AccAddress,
	collateral sdk.Coin,
	dinrAmount sdk.Int,
) error {
	// Calculate fee
	feeNAMO, err := k.CalculateFeeInNAMO(ctx, dinrAmount, "mint")
	if err != nil {
		return err
	}
	
	// Collect fee in NAMO
	err = k.CollectDINRFeeInNAMO(ctx, minter, dinrAmount, "mint", collateral)
	if err != nil {
		return err
	}
	
	// Proceed with minting (existing logic)
	// Lock collateral
	err = k.LockCollateral(ctx, minter, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}
	
	// Mint DINR
	dinrCoin := sdk.NewCoin("dinr", dinrAmount)
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(dinrCoin))
	if err != nil {
		return err
	}
	
	// Send to user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		minter,
		sdk.NewCoins(dinrCoin),
	)
	if err != nil {
		return err
	}
	
	return nil
}

// BurnDINRWithNAMOFee burns DINR and collects fee in NAMO
func (k Keeper) BurnDINRWithNAMOFee(
	ctx sdk.Context,
	burner sdk.AccAddress,
	dinrAmount sdk.Int,
) error {
	// Calculate fee
	feeNAMO, err := k.CalculateFeeInNAMO(ctx, dinrAmount, "burn")
	if err != nil {
		return err
	}
	
	// Check user has enough DINR
	dinrBalance := k.bankKeeper.GetBalance(ctx, burner, "dinr")
	if dinrBalance.Amount.LT(dinrAmount) {
		return fmt.Errorf("insufficient DINR balance")
	}
	
	// Collect fee in NAMO (user might pay in different token)
	userPaymentToken := k.bankKeeper.GetBalance(ctx, burner, "namo")
	if userPaymentToken.IsZero() {
		// Try other tokens
		balances := k.bankKeeper.GetAllBalances(ctx, burner)
		if !balances.IsZero() {
			userPaymentToken = balances[0] // Use first available token
		}
	}
	
	err = k.CollectDINRFeeInNAMO(ctx, burner, dinrAmount, "burn", userPaymentToken)
	if err != nil {
		return err
	}
	
	// Burn DINR
	dinrCoin := sdk.NewCoin("dinr", dinrAmount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		burner,
		types.ModuleName,
		sdk.NewCoins(dinrCoin),
	)
	if err != nil {
		return err
	}
	
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(dinrCoin))
	if err != nil {
		return err
	}
	
	// Release collateral (existing logic)
	collateral, err := k.GetUserCollateral(ctx, burner)
	if err == nil && !collateral.IsZero() {
		// Calculate proportional collateral to release
		userTotalDINR := k.GetUserMintedDINR(ctx, burner)
		if !userTotalDINR.IsZero() {
			releaseRatio := sdk.NewDecFromInt(dinrAmount).Quo(sdk.NewDecFromInt(userTotalDINR))
			releaseAmount := sdk.NewDecFromInt(collateral.AmountOf("namo")).Mul(releaseRatio).TruncateInt()
			
			if releaseAmount.IsPositive() {
				releaseCoin := sdk.NewCoin("namo", releaseAmount)
				err = k.ReleaseCollateral(ctx, burner, sdk.NewCoins(releaseCoin))
				if err != nil {
					// Log error but don't fail the burn
					ctx.Logger().Error("failed to release collateral", "error", err)
				}
			}
		}
	}
	
	return nil
}

// GetFeeEstimate returns fee estimate in NAMO for a DINR operation
func (k Keeper) GetFeeEstimate(ctx sdk.Context, amount sdk.Int, operation string) (sdk.Coin, error) {
	return k.CalculateFeeInNAMO(ctx, amount, operation)
}

// SetInclusiveFeeOption sets whether fees are inclusive or on-top
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