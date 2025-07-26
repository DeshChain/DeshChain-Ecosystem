package dinr

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
)

// BeginBlocker updates stability metrics and processes yield strategies
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Update stability metrics every block
	k.UpdateStabilityData(ctx)

	// Process yield strategies every hour
	if shouldProcessYield(ctx, k) {
		k.ProcessYieldStrategies(ctx)
	}

	// Check and liquidate unhealthy positions
	k.ProcessLiquidations(ctx)
}

// EndBlocker processes stability mechanisms and daily resets
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {
	// Run stability controller to maintain DINR peg
	stabilityController := k.GetStabilityController()
	err := stabilityController.MaintainPeg(ctx)
	if err != nil {
		// Log error but don't fail the block
		k.Logger(ctx).Error("stability controller error", "error", err)
	}
	
	// Reset daily limits at the start of each day (UTC midnight)
	if shouldResetDailyLimits(ctx) {
		k.ResetDailyAmounts(ctx)
	}
	
	return []abci.ValidatorUpdate{}
}

// shouldProcessYield checks if it's time to process yield strategies
func shouldProcessYield(ctx sdk.Context, k keeper.Keeper) bool {
	lastYieldTime := k.GetLastYieldProcessingTime(ctx)
	currentTime := ctx.BlockTime()
	
	// Process yield every hour
	return currentTime.Sub(lastYieldTime) >= time.Hour
}

// shouldResetDailyLimits checks if it's a new day and limits should be reset
func shouldResetDailyLimits(ctx sdk.Context) bool {
	currentTime := ctx.BlockTime()
	
	// Reset at UTC midnight
	_, _, day := currentTime.Date()
	year, month, _ := currentTime.AddDate(0, 0, -1).Date()
	previousMidnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	
	// Check if we've crossed midnight since the last block
	return currentTime.After(previousMidnight) && currentTime.Hour() == 0 && currentTime.Minute() < 10
}