package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tax/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the tax MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}