package oracle

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/oracle/keeper"
	"github.com/deshchain/deshchain/x/oracle/types"
)

// BeginBlocker processes oracle aggregation windows and updates prices
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Process price aggregation windows
	if err := k.ProcessAggregationWindow(ctx); err != nil {
		k.Logger(ctx).Error("failed to process aggregation window", "error", err)
	}

	// Clean up old submissions and historical data
	k.cleanupOldData(ctx)
}

// EndBlocker performs end of block processing for oracle module
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Check for stale prices and emit warnings
	k.checkStalePrices(ctx)

	// Update validator statistics
	k.updateValidatorStats(ctx)
}