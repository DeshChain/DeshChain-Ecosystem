package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
)

// CollateralManager handles all collateral-related operations
type CollateralManager struct {
	keeper Keeper
}

// NewCollateralManager creates a new collateral manager
func NewCollateralManager(keeper Keeper) *CollateralManager {
	return &CollateralManager{
		keeper: keeper,
	}
}

// ValidateCollateral checks if a collateral asset is valid and active
func (cm *CollateralManager) ValidateCollateral(ctx sdk.Context, denom string) (*types.CollateralAsset, error) {
	params := cm.keeper.GetParams(ctx)
	
	for _, asset := range params.CollateralAssets {
		if asset.Denom == denom && asset.IsActive {
			return &asset, nil
		}
	}
	
	return nil, types.ErrInvalidCollateral
}

// CalculateCollateralValue calculates the total value of collateral in INR
func (cm *CollateralManager) CalculateCollateralValue(ctx sdk.Context, collateral sdk.Coins) (sdk.Int, error) {
	totalValue := sdk.ZeroInt()
	
	for _, coin := range collateral {
		// Validate collateral asset
		asset, err := cm.ValidateCollateral(ctx, coin.Denom)
		if err != nil {
			continue // Skip invalid collateral
		}
		
		// Get price from oracle
		price, err := cm.keeper.oracleKeeper.GetPrice(ctx, asset.OracleScriptId)
		if err != nil {
			return sdk.ZeroInt(), fmt.Errorf("failed to get price for %s: %w", coin.Denom, err)
		}
		
		// Calculate value in INR
		coinValue := coin.Amount.ToDec().Mul(price).TruncateInt()
		totalValue = totalValue.Add(coinValue)
	}
	
	return totalValue, nil
}

// CalculateHealthFactor calculates the health factor for a position
func (cm *CollateralManager) CalculateHealthFactor(ctx sdk.Context, collateralValue, debtValue sdk.Int) sdk.Dec {
	params := cm.keeper.GetParams(ctx)
	
	if debtValue.IsZero() {
		return sdk.NewDec(1000) // Very high health factor if no debt
	}
	
	// Health Factor = (Collateral Value * Liquidation Threshold) / Debt Value
	liquidationThreshold := sdk.NewDec(int64(params.LiquidationThreshold)).QuoInt64(10000)
	collateralValueDec := collateralValue.ToDec()
	debtValueDec := debtValue.ToDec()
	
	healthFactor := collateralValueDec.Mul(liquidationThreshold).Quo(debtValueDec)
	return healthFactor
}

// CheckLiquidationEligibility determines if a position can be liquidated
func (cm *CollateralManager) CheckLiquidationEligibility(ctx sdk.Context, userAddr string) (bool, *types.UserPosition, error) {
	position, found := cm.keeper.GetUserPosition(ctx, userAddr)
	if !found {
		return false, nil, types.ErrPositionNotFound
	}
	
	// Calculate current collateral value
	collateralValue, err := cm.CalculateCollateralValue(ctx, position.Collateral)
	if err != nil {
		return false, nil, err
	}
	
	// Calculate health factor
	healthFactor := cm.CalculateHealthFactor(ctx, collateralValue, position.DinrMinted.Amount)
	
	// Position is liquidatable if health factor < 1.0
	isLiquidatable := healthFactor.LT(sdk.OneDec())
	
	return isLiquidatable, &position, nil
}

// CalculateLiquidationAmount determines how much debt can be liquidated
func (cm *CollateralManager) CalculateLiquidationAmount(ctx sdk.Context, position *types.UserPosition, maxLiquidationAmount sdk.Int) (sdk.Int, error) {
	params := cm.keeper.GetParams(ctx)
	
	// Maximum liquidation is 50% of the debt or the specified amount, whichever is smaller
	maxAllowedLiquidation := position.DinrMinted.Amount.QuoInt64(2)
	
	if maxLiquidationAmount.IsZero() || maxLiquidationAmount.GT(maxAllowedLiquidation) {
		return maxAllowedLiquidation, nil
	}
	
	return maxLiquidationAmount, nil
}

// CalculateLiquidationReward calculates collateral reward for liquidator
func (cm *CollateralManager) CalculateLiquidationReward(ctx sdk.Context, liquidationAmount sdk.Int, position *types.UserPosition) (sdk.Coins, error) {
	params := cm.keeper.GetParams(ctx)
	
	// Calculate base collateral value to seize (1:1 with debt)
	baseValue := liquidationAmount
	
	// Add liquidation penalty (bonus for liquidator)
	penaltyBps := params.Fees.LiquidationPenalty
	penalty := baseValue.Mul(sdk.NewInt(int64(penaltyBps))).QuoInt64(10000)
	totalValue := baseValue.Add(penalty)
	
	// Select collateral proportionally
	return cm.selectCollateralForLiquidation(ctx, position.Collateral, totalValue)
}

// selectCollateralForLiquidation selects collateral coins for liquidation
func (cm *CollateralManager) selectCollateralForLiquidation(ctx sdk.Context, collateral sdk.Coins, targetValue sdk.Int) (sdk.Coins, error) {
	totalCollateralValue, err := cm.CalculateCollateralValue(ctx, collateral)
	if err != nil {
		return nil, err
	}
	
	if totalCollateralValue.IsZero() {
		return sdk.NewCoins(), nil
	}
	
	selectedCollateral := sdk.NewCoins()
	
	// Calculate proportion to liquidate
	proportion := targetValue.ToDec().Quo(totalCollateralValue.ToDec())
	if proportion.GT(sdk.OneDec()) {
		proportion = sdk.OneDec() // Cap at 100%
	}
	
	for _, coin := range collateral {
		amountToSeize := coin.Amount.ToDec().Mul(proportion).TruncateInt()
		if amountToSeize.GT(sdk.ZeroInt()) {
			selectedCollateral = selectedCollateral.Add(sdk.NewCoin(coin.Denom, amountToSeize))
		}
	}
	
	return selectedCollateral, nil
}

// ProcessLiquidation executes a liquidation
func (cm *CollateralManager) ProcessLiquidation(ctx sdk.Context, liquidator sdk.AccAddress, userAddr string, dinrAmount sdk.Int) error {
	// Check liquidation eligibility
	isEligible, position, err := cm.CheckLiquidationEligibility(ctx, userAddr)
	if err != nil {
		return err
	}
	
	if !isEligible {
		return types.ErrPositionNotLiquidatable
	}
	
	// Calculate liquidation amount
	liquidationAmount, err := cm.CalculateLiquidationAmount(ctx, position, dinrAmount)
	if err != nil {
		return err
	}
	
	// Calculate collateral reward
	collateralReward, err := cm.CalculateLiquidationReward(ctx, liquidationAmount, position)
	if err != nil {
		return err
	}
	
	// Execute liquidation
	dinrCoin := sdk.NewCoin("dinr", liquidationAmount)
	
	// Transfer DINR from liquidator to module and burn
	err = cm.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, liquidator, types.ModuleName, sdk.NewCoins(dinrCoin))
	if err != nil {
		return fmt.Errorf("failed to transfer DINR from liquidator: %w", err)
	}
	
	err = cm.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(dinrCoin))
	if err != nil {
		return fmt.Errorf("failed to burn DINR: %w", err)
	}
	
	// Transfer collateral reward to liquidator
	err = cm.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidator, collateralReward)
	if err != nil {
		return fmt.Errorf("failed to transfer collateral to liquidator: %w", err)
	}
	
	// Update user position
	userAddress, _ := sdk.AccAddressFromBech32(userAddr)
	err = cm.updatePositionAfterLiquidation(ctx, userAddress, liquidationAmount, collateralReward)
	if err != nil {
		return err
	}
	
	// Emit liquidation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLiquidate,
			sdk.NewAttribute(types.AttributeKeyLiquidator, liquidator.String()),
			sdk.NewAttribute(types.AttributeKeyUser, userAddr),
			sdk.NewAttribute(types.AttributeKeyDINRCovered, dinrCoin.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralReceived, collateralReward.String()),
		),
	)
	
	return nil
}

// updatePositionAfterLiquidation updates user position after liquidation
func (cm *CollateralManager) updatePositionAfterLiquidation(ctx sdk.Context, userAddr sdk.AccAddress, liquidatedDebt sdk.Int, seizedCollateral sdk.Coins) error {
	position, found := cm.keeper.GetUserPosition(ctx, userAddr.String())
	if !found {
		return types.ErrPositionNotFound
	}
	
	// Update debt
	position.DinrMinted = position.DinrMinted.Sub(sdk.NewCoin("dinr", liquidatedDebt))
	
	// Update collateral
	position.Collateral = position.Collateral.Sub(seizedCollateral)
	
	// If position is fully liquidated, remove it
	if position.DinrMinted.IsZero() || position.Collateral.IsZero() {
		cm.keeper.RemoveUserPosition(ctx, userAddr.String())
		return nil
	}
	
	// Recalculate health factor
	collateralValue, err := cm.CalculateCollateralValue(ctx, position.Collateral)
	if err != nil {
		return err
	}
	
	healthFactor := cm.CalculateHealthFactor(ctx, collateralValue, position.DinrMinted.Amount)
	position.HealthFactor = healthFactor.String()
	position.LastUpdate = ctx.BlockTime()
	
	cm.keeper.SetUserPosition(ctx, position)
	return nil
}

// GetCollateralRatio calculates the current collateral ratio for a position
func (cm *CollateralManager) GetCollateralRatio(ctx sdk.Context, userAddr string) (sdk.Dec, error) {
	position, found := cm.keeper.GetUserPosition(ctx, userAddr)
	if !found {
		return sdk.ZeroDec(), types.ErrPositionNotFound
	}
	
	if position.DinrMinted.IsZero() {
		return sdk.NewDec(10000), nil // Very high ratio if no debt
	}
	
	collateralValue, err := cm.CalculateCollateralValue(ctx, position.Collateral)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Collateral Ratio = (Collateral Value / Debt Value) * 100
	ratio := collateralValue.ToDec().Quo(position.DinrMinted.Amount.ToDec()).MulInt64(100)
	return ratio, nil
}

// GetMaxMintableAmount calculates maximum DINR that can be minted with given collateral
func (cm *CollateralManager) GetMaxMintableAmount(ctx sdk.Context, collateral sdk.Coins) (sdk.Int, error) {
	params := cm.keeper.GetParams(ctx)
	
	collateralValue, err := cm.CalculateCollateralValue(ctx, collateral)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	
	// Max mintable = (Collateral Value * 100) / Min Collateral Ratio
	minRatio := sdk.NewDec(int64(params.MinCollateralRatio)).QuoInt64(100)
	maxMintable := collateralValue.ToDec().Quo(minRatio).TruncateInt()
	
	return maxMintable, nil
}

// ValidateCollateralSufficiency checks if collateral is sufficient for minting
func (cm *CollateralManager) ValidateCollateralSufficiency(ctx sdk.Context, collateral sdk.Coins, mintAmount sdk.Int) error {
	maxMintable, err := cm.GetMaxMintableAmount(ctx, collateral)
	if err != nil {
		return err
	}
	
	if mintAmount.GT(maxMintable) {
		return types.ErrInsufficientCollateral
	}
	
	return nil
}