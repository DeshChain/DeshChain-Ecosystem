package tradefinance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/deshchain/deshchain/x/tradefinance/keeper"
)

// BeginBlocker processes automatic LC expiry and status updates
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Process LC expiries
	k.ProcessLcExpiries(ctx)
	
	// Process payment due dates
	k.ProcessPaymentDueDates(ctx)
	
	// Update statistics
	k.UpdateStatistics(ctx)
}