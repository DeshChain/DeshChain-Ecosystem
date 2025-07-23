package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tax/types"
)

// GetParams returns the total set of tax parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of tax parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetTaxRate returns the current tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)
	return params.TaxRate
}

// IsTaxEnabled returns whether tax collection is enabled
func (k Keeper) IsTaxEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.Enabled
}