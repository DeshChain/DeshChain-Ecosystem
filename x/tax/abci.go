package tax

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Tax module doesn't need begin block processing
	// Tax collection happens during transaction execution
}