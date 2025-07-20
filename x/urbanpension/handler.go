/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package urbanpension

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/deshchain/deshchain/x/urbanpension/keeper"
	"github.com/deshchain/deshchain/x/urbanpension/types"
)

// NewHandler creates a new urban pension handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateUrbanPensionScheme:
			res, err := msgServer.CreateUrbanPensionScheme(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgContributeToPension:
			res, err := msgServer.ContributeToPension(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgApplyForEducationLoan:
			res, err := msgServer.ApplyForEducationLoan(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCreateInsurancePolicy:
			res, err := msgServer.CreateInsurancePolicy(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgApplyForSMELoan:
			res, err := msgServer.ApplyForSMELoan(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCreateRBFInvestment:
			res, err := msgServer.CreateRBFInvestment(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgProcessPensionMaturity:
			res, err := msgServer.ProcessPensionMaturity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRepayLoan:
			res, err := msgServer.RepayLoan(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgClaimInsurance:
			res, err := msgServer.ClaimInsurance(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgUpdateParams:
			res, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// BeginBlocker handles logic at the beginning of each block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Process daily operations
	k.ProcessDailyOperations(ctx)
}

// EndBlocker handles logic at the end of each block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Process end of block operations
	k.ProcessEndBlockOperations(ctx)
}