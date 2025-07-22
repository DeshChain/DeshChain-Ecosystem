package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/dinr/types"
)

// UpdateStabilityData updates the global stability metrics
func (k Keeper) UpdateStabilityData(ctx sdk.Context) {
	// Get current stability data
	stabilityData := k.GetStabilityData(ctx)

	// Calculate total DINR supply
	totalSupply := k.bankKeeper.GetSupply(ctx, types.DINRDenom)
	stabilityData.TotalSupply = totalSupply

	// Calculate total collateral value across all positions
	totalCollateralValue := sdk.ZeroInt()
	k.IterateAllUserPositions(ctx, func(position types.UserPosition) bool {
		collateralValue := k.calculateTotalCollateralValue(ctx, position.Collateral)
		totalCollateralValue = totalCollateralValue.Add(collateralValue)
		return false
	})

	stabilityData.TotalCollateralValue = sdk.NewCoin("inr", totalCollateralValue)

	// Calculate global collateral ratio
	if !totalSupply.Amount.IsZero() {
		stabilityData.GlobalCollateralRatio = k.calculateCollateralRatio(totalCollateralValue, totalSupply.Amount)
	} else {
		stabilityData.GlobalCollateralRatio = 0
	}

	// Get current price from oracle (for now, use target price)
	// In production, this would fetch from oracle
	stabilityData.CurrentPrice = stabilityData.TargetPrice

	// Calculate price deviation
	currentPriceDec, _ := sdk.NewDecFromStr(stabilityData.CurrentPrice)
	targetPriceDec, _ := sdk.NewDecFromStr(stabilityData.TargetPrice)
	deviation := currentPriceDec.Sub(targetPriceDec).Abs().Mul(sdk.NewDec(10000)).Quo(targetPriceDec)
	stabilityData.PriceDeviation = deviation.TruncateInt64()

	// Update timestamp
	stabilityData.LastUpdate = ctx.BlockTime()

	// Save updated stability data
	k.SetStabilityData(ctx, stabilityData)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStabilityUpdate,
			sdk.NewAttribute(types.AttributeKeyCurrentPrice, stabilityData.CurrentPrice),
			sdk.NewAttribute(types.AttributeKeyTargetPrice, stabilityData.TargetPrice),
			sdk.NewAttribute(types.AttributeKeyPriceDeviation, fmt.Sprintf("%d", stabilityData.PriceDeviation)),
			sdk.NewAttribute("total_supply", stabilityData.TotalSupply.String()),
			sdk.NewAttribute("global_collateral_ratio", fmt.Sprintf("%d", stabilityData.GlobalCollateralRatio)),
		),
	)
}

// ProcessYieldStrategies processes yield generation from various strategies
func (k Keeper) ProcessYieldStrategies(ctx sdk.Context) {
	params := k.GetParams(ctx)
	strategies := k.GetAllYieldStrategies(ctx)

	totalYield := sdk.ZeroInt()

	for _, strategy := range strategies {
		if !strategy.IsActive {
			continue
		}

		// Calculate yield based on deployed amount and APY
		deployedAmount := strategy.DeployedAmount.Amount
		if deployedAmount.IsZero() {
			continue
		}

		// Calculate hourly yield (APY / 365 / 24)
		apy, _ := sdk.NewDecFromStr(strategy.ExpectedApy)
		hourlyRate := apy.Quo(sdk.NewDec(365 * 24 * 100)) // Convert percentage to decimal
		yieldAmount := sdk.NewDecFromInt(deployedAmount).Mul(hourlyRate).TruncateInt()

		if yieldAmount.GT(sdk.ZeroInt()) {
			// Mint yield as DINR
			yieldCoin := sdk.NewCoin(types.DINRDenom, yieldAmount)
			err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(yieldCoin))
			if err != nil {
				continue
			}

			// Distribute yield to insurance fund and platform revenue
			insuranceAmount := yieldAmount.Mul(sdk.NewInt(int64(params.YieldToInsuranceRatio))).Quo(sdk.NewInt(10000))
			platformAmount := yieldAmount.Sub(insuranceAmount)

			// Add to insurance fund
			if insuranceAmount.GT(sdk.ZeroInt()) {
				insuranceCoin := sdk.NewCoin(types.DINRDenom, insuranceAmount)
				k.AddToInsuranceFund(ctx, insuranceCoin)
			}

			// Distribute to platform revenue
			if platformAmount.GT(sdk.ZeroInt()) {
				platformCoin := sdk.NewCoin(types.DINRDenom, platformAmount)
				k.revenueKeeper.DistributePlatformRevenue(ctx, sdk.NewCoins(platformCoin))
			}

			totalYield = totalYield.Add(yieldAmount)

			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeYieldDistribution,
					sdk.NewAttribute(types.AttributeKeyStrategy, strategy.Id),
					sdk.NewAttribute(types.AttributeKeyYieldAmount, yieldAmount.String()),
				),
			)
		}
	}

	// Update last yield processing time
	k.SetLastYieldProcessingTime(ctx, ctx.BlockTime())

	// Update total fees collected
	k.IncrementTotalFeesCollected(ctx, totalYield)
}

// ProcessLiquidations checks and processes liquidatable positions
func (k Keeper) ProcessLiquidations(ctx sdk.Context) {
	params := k.GetParams(ctx)

	k.IterateAllUserPositions(ctx, func(position types.UserPosition) bool {
		// Calculate current collateral ratio
		collateralValue := k.calculateTotalCollateralValue(ctx, position.Collateral)
		collateralRatio := k.calculateCollateralRatio(collateralValue, position.DinrMinted.Amount)

		// Check if position is liquidatable
		if collateralRatio < uint64(params.LiquidationThreshold) {
			// Mark position for liquidation
			// In a real implementation, this would trigger a liquidation auction
			// For now, we just emit an event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"position_liquidatable",
					sdk.NewAttribute(types.AttributeKeyUser, position.Address),
					sdk.NewAttribute("collateral_ratio", fmt.Sprintf("%d", collateralRatio)),
					sdk.NewAttribute("liquidation_threshold", fmt.Sprintf("%d", params.LiquidationThreshold)),
				),
			)
		}

		return false
	})
}

// AddToInsuranceFund adds funds to the insurance fund
func (k Keeper) AddToInsuranceFund(ctx sdk.Context, amount sdk.Coin) {
	insuranceFund := k.GetInsuranceFund(ctx)
	insuranceFund.Balance = insuranceFund.Balance.Add(amount)
	insuranceFund.Assets = insuranceFund.Assets.Add(amount)
	k.SetInsuranceFund(ctx, insuranceFund)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInsuranceFundUpdate,
			sdk.NewAttribute(types.AttributeKeyInsuranceBalance, insuranceFund.Balance.String()),
		),
	)
}

// IncrementTotalFeesCollected updates the total fees collected
func (k Keeper) IncrementTotalFeesCollected(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := types.TotalFeesCollectedKey

	// Get current total
	var total sdk.Int
	bz := store.Get(key)
	if bz != nil {
		total, _ = sdk.NewIntFromString(string(bz))
	} else {
		total = sdk.ZeroInt()
	}

	// Increment and save
	total = total.Add(amount)
	store.Set(key, []byte(total.String()))
}