package revenue

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/deshchain/deshchain/x/revenue/keeper"
	"github.com/deshchain/deshchain/x/revenue/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Revenue module doesn't need begin block processing
	// Revenue collection and distribution happens during transaction execution
}