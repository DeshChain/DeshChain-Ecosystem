package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/revenue/types"
)

// GetParams returns the total set of revenue parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of revenue parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// IsRevenueEnabled returns whether revenue collection is enabled
func (k Keeper) IsRevenueEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.Enabled
}