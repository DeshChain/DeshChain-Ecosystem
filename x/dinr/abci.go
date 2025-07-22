package dinr

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/deshchain/deshchain/x/dinr/keeper"
	"github.com/deshchain/deshchain/x/dinr/types"
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

// EndBlocker currently does nothing but is here for potential future use
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {
	// Placeholder for future end-block processing
	return []abci.ValidatorUpdate{}
}

// shouldProcessYield checks if it's time to process yield strategies
func shouldProcessYield(ctx sdk.Context, k keeper.Keeper) bool {
	lastYieldTime := k.GetLastYieldProcessingTime(ctx)
	currentTime := ctx.BlockTime()
	
	// Process yield every hour
	return currentTime.Sub(lastYieldTime) >= time.Hour
}