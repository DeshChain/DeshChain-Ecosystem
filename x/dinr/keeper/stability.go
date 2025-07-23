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
	collateralManager := k.GetCollateralManager()

	k.IterateAllUserPositions(ctx, func(position types.UserPosition) bool {
		// Check if position is liquidatable using CollateralManager
		isLiquidatable, _, err := collateralManager.CheckLiquidationEligibility(ctx, position.Address)
		if err != nil {
			return false // Continue to next position
		}

		if isLiquidatable {
			// Calculate current collateral ratio for event
			collateralValue, err := collateralManager.CalculateCollateralValue(ctx, position.Collateral)
			if err != nil {
				return false
			}
			
			healthFactor := collateralManager.CalculateHealthFactor(ctx, collateralValue, position.DinrMinted.Amount)
			
			// Emit liquidatable position event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"position_liquidatable",
					sdk.NewAttribute(types.AttributeKeyUser, position.Address),
					sdk.NewAttribute("health_factor", healthFactor.String()),
					sdk.NewAttribute("collateral_value", collateralValue.String()),
					sdk.NewAttribute("debt_amount", position.DinrMinted.String()),
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

// StabilityController manages DINR price stability mechanisms
type StabilityController struct {
	keeper      Keeper
	targetPrice sdk.Dec
	tolerance   sdk.Dec // Allowed deviation percentage (e.g., 0.01 = 1%)
}

// NewStabilityController creates a new stability controller
func NewStabilityController(keeper Keeper) *StabilityController {
	return &StabilityController{
		keeper:      keeper,
		targetPrice: sdk.OneDec(), // $1.00 target
		tolerance:   sdk.NewDecWithPrec(1, 2), // 1% tolerance
	}
}

// MaintainPeg checks current price and executes stability actions if needed
func (sc *StabilityController) MaintainPeg(ctx sdk.Context) error {
	// Get current DINR price from oracle
	currentPrice, err := sc.keeper.GetCurrentPrice(ctx, "DINR")
	if err != nil {
		return fmt.Errorf("failed to get DINR price: %w", err)
	}

	// Calculate deviation from target
	deviation := sc.calculateDeviation(currentPrice)
	
	// Check if intervention is needed
	if deviation.Abs().GT(sc.tolerance) {
		return sc.executeStabilityAction(ctx, currentPrice, deviation)
	}

	return nil
}

// calculateDeviation returns the percentage deviation from target price
func (sc *StabilityController) calculateDeviation(currentPrice sdk.Dec) sdk.Dec {
	if sc.targetPrice.IsZero() {
		return sdk.ZeroDec()
	}
	
	// (currentPrice - targetPrice) / targetPrice
	return currentPrice.Sub(sc.targetPrice).Quo(sc.targetPrice)
}

// executeStabilityAction performs minting/burning to stabilize price
func (sc *StabilityController) executeStabilityAction(ctx sdk.Context, currentPrice, deviation sdk.Dec) error {
	params := sc.keeper.GetParams(ctx)
	
	// Calculate intervention amount based on deviation magnitude
	interventionAmount := sc.calculateInterventionAmount(ctx, deviation)
	
	if deviation.IsPositive() {
		// Price above target - increase supply (mint DINR)
		return sc.executeMinting(ctx, interventionAmount, params)
	} else {
		// Price below target - decrease supply (burn DINR)
		return sc.executeBurning(ctx, interventionAmount, params)
	}
}

// calculateInterventionAmount determines how much to mint/burn
func (sc *StabilityController) calculateInterventionAmount(ctx sdk.Context, deviation sdk.Dec) sdk.Int {
	// Get current DINR supply
	supply := sc.keeper.bankKeeper.GetSupply(ctx, "dinr")
	
	// Base intervention: 1% of supply per 1% deviation (with max cap)
	interventionRate := deviation.Abs()
	maxIntervention := sdk.NewDecWithPrec(5, 2) // 5% max intervention
	
	if interventionRate.GT(maxIntervention) {
		interventionRate = maxIntervention
	}
	
	interventionAmount := supply.Amount.ToDec().Mul(interventionRate).TruncateInt()
	
	// Minimum intervention threshold
	minIntervention := sdk.NewInt(1000000) // 1 DINR minimum
	if interventionAmount.LT(minIntervention) {
		interventionAmount = minIntervention
	}
	
	return interventionAmount
}

// executeMinting mints new DINR to increase supply
func (sc *StabilityController) executeMinting(ctx sdk.Context, amount sdk.Int, params types.Params) error {
	// Check if minting is enabled and within limits
	if !params.MintingEnabled {
		return types.ErrMintingDisabled
	}
	
	// Check daily minting limit
	dailyMinted := sc.keeper.GetDailyMintedAmount(ctx)
	if dailyMinted.Add(amount).GT(sdk.NewInt(int64(params.MaxDailyMinting))) {
		// Reduce to maximum allowed
		amount = sdk.NewInt(int64(params.MaxDailyMinting)).Sub(dailyMinted)
		if amount.IsZero() || amount.IsNegative() {
			return types.ErrDailyMintingLimitExceeded
		}
	}
	
	// Mint DINR to stability pool
	coins := sdk.NewCoins(sdk.NewCoin("dinr", amount))
	stabilityPoolAddr := sc.keeper.accountKeeper.GetModuleAddress(types.StabilityPoolName)
	
	err := sc.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return fmt.Errorf("failed to mint DINR: %w", err)
	}
	
	err = sc.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, stabilityPoolAddr, coins)
	if err != nil {
		return fmt.Errorf("failed to send minted DINR to stability pool: %w", err)
	}
	
	// Update daily minting tracking
	sc.keeper.SetDailyMintedAmount(ctx, dailyMinted.Add(amount))
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStabilityMint,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyReason, "price_above_target"),
		),
	)
	
	return nil
}

// executeBurning burns DINR to decrease supply
func (sc *StabilityController) executeBurning(ctx sdk.Context, amount sdk.Int, params types.Params) error {
	// Check if burning is enabled
	if !params.BurningEnabled {
		return types.ErrBurningDisabled
	}
	
	// Check stability pool balance
	stabilityPoolAddr := sc.keeper.accountKeeper.GetModuleAddress(types.StabilityPoolName)
	balance := sc.keeper.bankKeeper.GetBalance(ctx, stabilityPoolAddr, "dinr")
	
	if balance.Amount.LT(amount) {
		// Burn only what's available
		amount = balance.Amount
		if amount.IsZero() {
			return types.ErrInsufficientStabilityPoolBalance
		}
	}
	
	// Check daily burning limit
	dailyBurned := sc.keeper.GetDailyBurnedAmount(ctx)
	if dailyBurned.Add(amount).GT(sdk.NewInt(int64(params.MaxDailyBurning))) {
		// Reduce to maximum allowed
		amount = sdk.NewInt(int64(params.MaxDailyBurning)).Sub(dailyBurned)
		if amount.IsZero() || amount.IsNegative() {
			return types.ErrDailyBurningLimitExceeded
		}
	}
	
	// Burn DINR from stability pool
	coins := sdk.NewCoins(sdk.NewCoin("dinr", amount))
	
	err := sc.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, stabilityPoolAddr, types.ModuleName, coins)
	if err != nil {
		return fmt.Errorf("failed to send DINR from stability pool: %w", err)
	}
	
	err = sc.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return fmt.Errorf("failed to burn DINR: %w", err)
	}
	
	// Update daily burning tracking
	sc.keeper.SetDailyBurnedAmount(ctx, dailyBurned.Add(amount))
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStabilityBurn,
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyReason, "price_below_target"),
		),
	)
	
	return nil
}

// CheckStabilityHealth returns health metrics for the stability system
func (sc *StabilityController) CheckStabilityHealth(ctx sdk.Context) (*types.StabilityHealth, error) {
	currentPrice, err := sc.keeper.GetCurrentPrice(ctx, "DINR")
	if err != nil {
		return nil, err
	}
	
	deviation := sc.calculateDeviation(currentPrice)
	stabilityPoolAddr := sc.keeper.accountKeeper.GetModuleAddress(types.StabilityPoolName)
	poolBalance := sc.keeper.bankKeeper.GetBalance(ctx, stabilityPoolAddr, "dinr")
	totalSupply := sc.keeper.bankKeeper.GetSupply(ctx, "dinr")
	
	health := &types.StabilityHealth{
		CurrentPrice:    currentPrice.String(),
		TargetPrice:     sc.targetPrice.String(),
		Deviation:       deviation.String(),
		PoolBalance:     poolBalance.Amount.String(),
		TotalSupply:     totalSupply.Amount.String(),
		WithinTolerance: deviation.Abs().LTE(sc.tolerance),
		DailyMinted:     sc.keeper.GetDailyMintedAmount(ctx).String(),
		DailyBurned:     sc.keeper.GetDailyBurnedAmount(ctx).String(),
		LastUpdate:      ctx.BlockTime(),
	}
	
	return health, nil
}