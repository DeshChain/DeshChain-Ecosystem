package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// GetTradeFinanceStats returns the current trade finance statistics
func (k Keeper) GetTradeFinanceStats(ctx sdk.Context) types.TradeFinanceStats {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TradeFinanceStatsKey)
	
	if bz == nil {
		// Return default stats
		return types.TradeFinanceStats{
			TotalLcsIssued:         0,
			TotalTradeValue:        sdk.NewCoin("dinr", sdk.ZeroInt()),
			ActiveLcs:              0,
			CompletedTrades:        0,
			AverageProcessingHours: 0,
			TopTradeCorridor:       "",
			TotalFeesCollected:     sdk.NewCoin("dinr", sdk.ZeroInt()),
			DocumentsVerified:      0,
			LastUpdate:             ctx.BlockTime(),
		}
	}
	
	var stats types.TradeFinanceStats
	k.cdc.MustUnmarshal(bz, &stats)
	return stats
}

// SetTradeFinanceStats saves the trade finance statistics
func (k Keeper) SetTradeFinanceStats(ctx sdk.Context, stats types.TradeFinanceStats) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set(types.TradeFinanceStatsKey, bz)
}