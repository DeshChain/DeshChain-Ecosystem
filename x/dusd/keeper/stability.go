package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dusd/types"
)

// StabilityEngine manages DUSD price stability using same algorithms as DINR
type StabilityEngine struct {
	keeper *Keeper
}

// NewStabilityEngine creates a new stability engine
func NewStabilityEngine(keeper *Keeper) *StabilityEngine {
	return &StabilityEngine{
		keeper: keeper,
	}
}

// CheckPriceStability monitors DUSD price and triggers rebalancing if needed
func (se *StabilityEngine) CheckPriceStability(ctx sdk.Context) error {
	// Get current USD price
	currentPrice, err := se.keeper.GetUSDPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}
	
	// Get parameters
	params, err := se.keeper.GetParams(ctx)
	if err != nil {
		return err
	}
	
	// Parse target price and thresholds
	targetPrice, err := sdk.NewDecFromStr(params.TargetPrice)
	if err != nil {
		return err
	}
	
	rebalanceThreshold, err := sdk.NewDecFromStr(params.RebalanceThreshold)
	if err != nil {
		return err
	}
	
	emergencyThreshold, err := sdk.NewDecFromStr(params.EmergencyThreshold)
	if err != nil {
		return err
	}
	
	// Calculate price deviation
	deviation := currentPrice.Sub(targetPrice).Quo(targetPrice).Abs()
	
	// Log price information
	se.keeper.Logger(ctx).Info("DUSD price stability check",
		"current_price", currentPrice.String(),
		"target_price", targetPrice.String(),
		"deviation", deviation.String(),
		"rebalance_threshold", rebalanceThreshold.String(),
	)
	
	// Check if emergency action is needed
	if deviation.GTE(emergencyThreshold) {
		return se.ExecuteEmergencyAction(ctx, currentPrice, targetPrice, deviation)
	}
	
	// Check if rebalancing is needed (same logic as DINR)
	if deviation.GTE(rebalanceThreshold) {
		return se.ExecuteRebalanceAction(ctx, currentPrice, targetPrice, deviation)
	}
	
	return nil
}

// ExecuteRebalanceAction performs algorithmic rebalancing (same as DINR)
func (se *StabilityEngine) ExecuteRebalanceAction(ctx sdk.Context, currentPrice, targetPrice, deviation sdk.Dec) error {
	// Determine action type based on price direction
	var actionType string
	if currentPrice.GT(targetPrice) {
		actionType = "expand_supply" // Price too high, increase supply
	} else {
		actionType = "contract_supply" // Price too low, decrease supply
	}
	
	// Calculate adjustment amount based on deviation and total supply
	totalSupply := se.keeper.bankKeeper.GetSupply(ctx, types.DUSDDenom)
	adjustmentPercentage := deviation.Mul(sdk.NewDecWithPrec(5, 1)) // 0.5x deviation
	adjustmentAmount := totalSupply.Amount.ToLegacyDec().Mul(adjustmentPercentage).TruncateInt()
	
	// Cap adjustment to reasonable limits (max 5% of supply)
	maxAdjustment := totalSupply.Amount.ToLegacyDec().Mul(sdk.NewDecWithPrec(5, 2)).TruncateInt()
	if adjustmentAmount.GT(maxAdjustment) {
		adjustmentAmount = maxAdjustment
	}
	
	// Execute the adjustment
	adjustmentCoin := sdk.NewCoin(types.DUSDDenom, adjustmentAmount)
	
	var err error
	switch actionType {
	case "expand_supply":
		err = se.ExpandSupply(ctx, adjustmentCoin)
	case "contract_supply":
		err = se.ContractSupply(ctx, adjustmentCoin)
	}
	
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", actionType, err)
	}
	
	// Record stability action
	actionID := types.GenerateActionID(actionType, ctx.BlockTime())
	stabilityAction := types.StabilityAction{
		Id:            actionID,
		ActionType:    actionType,
		TriggerPrice:  currentPrice.String(),
		TargetPrice:   targetPrice.String(),
		AmountAdjusted: adjustmentCoin,
		ExecutedAt:    ctx.BlockTime(),
		Result:        "success",
	}
	
	se.RecordStabilityAction(ctx, stabilityAction)
	
	se.keeper.Logger(ctx).Info("executed DUSD rebalance action",
		"action_type", actionType,
		"adjustment_amount", adjustmentCoin.String(),
		"trigger_price", currentPrice.String(),
		"target_price", targetPrice.String(),
	)
	
	return nil
}

// ExpandSupply increases DUSD supply to reduce price
func (se *StabilityEngine) ExpandSupply(ctx sdk.Context, amount sdk.Coin) error {
	// Mint new DUSD tokens to module account
	if err := se.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("failed to mint DUSD for supply expansion: %w", err)
	}
	
	// Distribute new tokens through treasury system for market operations
	// This could involve:
	// 1. Adding liquidity to DEX pools
	// 2. Strategic market operations
	// 3. Incentive programs
	
	return nil
}

// ContractSupply decreases DUSD supply to increase price
func (se *StabilityEngine) ContractSupply(ctx sdk.Context, amount sdk.Coin) error {
	// Check if module has enough DUSD to burn
	moduleAddr := se.keeper.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := se.keeper.bankKeeper.GetBalance(ctx, moduleAddr, types.DUSDDenom)
	
	if moduleBalance.Amount.LT(amount.Amount) {
		// If not enough in module, use treasury reserves
		if err := se.keeper.treasuryKeeper.WithdrawFromPool(ctx, "reserve", sdk.NewCoins(amount)); err != nil {
			return fmt.Errorf("insufficient DUSD for supply contraction: %w", err)
		}
	}
	
	// Burn DUSD tokens
	if err := se.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("failed to burn DUSD for supply contraction: %w", err)
	}
	
	return nil
}

// ExecuteEmergencyAction handles emergency price deviations
func (se *StabilityEngine) ExecuteEmergencyAction(ctx sdk.Context, currentPrice, targetPrice, deviation sdk.Dec) error {
	// Get parameters
	params, err := se.keeper.GetParams(ctx)
	if err != nil {
		return err
	}
	
	// Check if circuit breaker should be activated
	if params.CircuitBreakerEnabled && deviation.GTE(sdk.NewDecWithPrec(5, 2)) { // 5% deviation
		return se.ActivateCircuitBreaker(ctx, "emergency_price_deviation")
	}
	
	// Execute aggressive rebalancing
	totalSupply := se.keeper.bankKeeper.GetSupply(ctx, types.DUSDDenom)
	adjustmentPercentage := deviation.Mul(sdk.NewDecWithPrec(2, 1)) // 2x deviation for emergency
	adjustmentAmount := totalSupply.Amount.ToLegacyDec().Mul(adjustmentPercentage).TruncateInt()
	
	// Cap emergency adjustment to 10% of supply
	maxEmergencyAdjustment := totalSupply.Amount.ToLegacyDec().Mul(sdk.NewDecWithPrec(10, 2)).TruncateInt()
	if adjustmentAmount.GT(maxEmergencyAdjustment) {
		adjustmentAmount = maxEmergencyAdjustment
	}
	
	adjustmentCoin := sdk.NewCoin(types.DUSDDenom, adjustmentAmount)
	
	// Execute emergency adjustment
	var actionType string
	if currentPrice.GT(targetPrice) {
		actionType = "emergency_expand"
		err = se.ExpandSupply(ctx, adjustmentCoin)
	} else {
		actionType = "emergency_contract"
		err = se.ContractSupply(ctx, adjustmentCoin)
	}
	
	if err != nil {
		return fmt.Errorf("failed to execute emergency action: %w", err)
	}
	
	// Record emergency action
	actionID := types.GenerateActionID(actionType, ctx.BlockTime())
	stabilityAction := types.StabilityAction{
		Id:            actionID,
		ActionType:    actionType,
		TriggerPrice:  currentPrice.String(),
		TargetPrice:   targetPrice.String(),
		AmountAdjusted: adjustmentCoin,
		ExecutedAt:    ctx.BlockTime(),
		Result:        "emergency_executed",
	}
	
	se.RecordStabilityAction(ctx, stabilityAction)
	
	// Emit emergency event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"dusd_emergency_action",
			sdk.NewAttribute("action_type", actionType),
			sdk.NewAttribute("deviation", deviation.String()),
			sdk.NewAttribute("adjustment_amount", adjustmentCoin.String()),
		),
	)
	
	se.keeper.Logger(ctx).Error("executed DUSD emergency action",
		"action_type", actionType,
		"deviation", deviation.String(),
		"adjustment_amount", adjustmentCoin.String(),
	)
	
	return nil
}

// ActivateCircuitBreaker halts DUSD operations during extreme conditions
func (se *StabilityEngine) ActivateCircuitBreaker(ctx sdk.Context, reason string) error {
	// Record circuit breaker activation
	actionID := types.GenerateActionID("circuit_breaker", ctx.BlockTime())
	stabilityAction := types.StabilityAction{
		Id:         actionID,
		ActionType: "circuit_breaker",
		ExecutedAt: ctx.BlockTime(),
		Result:     fmt.Sprintf("activated: %s", reason),
	}
	
	se.RecordStabilityAction(ctx, stabilityAction)
	
	// Emit circuit breaker event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"dusd_circuit_breaker",
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	)
	
	se.keeper.Logger(ctx).Error("DUSD circuit breaker activated",
		"reason", reason,
		"timestamp", ctx.BlockTime().String(),
	)
	
	// In a real implementation, this would:
	// 1. Pause minting/burning operations
	// 2. Halt liquidations
	// 3. Send governance notifications
	// 4. Require manual intervention to resume
	
	return nil
}

// RecordStabilityAction stores a stability action record
func (se *StabilityEngine) RecordStabilityAction(ctx sdk.Context, action types.StabilityAction) {
	store := ctx.KVStore(se.keeper.storeKey)
	key := types.GetStabilityActionStoreKey(action.Id)
	bz := se.keeper.cdc.MustMarshal(&action)
	store.Set(key, bz)
}

// GetStabilityActions returns recent stability actions
func (se *StabilityEngine) GetStabilityActions(ctx sdk.Context, hours int64) []types.StabilityAction {
	store := ctx.KVStore(se.keeper.storeKey)
	
	// Calculate time threshold
	threshold := ctx.BlockTime().Add(-time.Duration(hours) * time.Hour)
	
	// Iterate through stability actions
	prefix := types.StabilityActionKey
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var actions []types.StabilityAction
	for ; iterator.Valid(); iterator.Next() {
		var action types.StabilityAction
		se.keeper.cdc.MustUnmarshal(iterator.Value(), &action)
		
		// Filter by time threshold
		if action.ExecutedAt.After(threshold) {
			actions = append(actions, action)
		}
	}
	
	return actions
}