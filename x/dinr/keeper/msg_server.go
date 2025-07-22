package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/dinr/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the dinr MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// MintDINR implements types.MsgServer
func (k msgServer) MintDINR(goCtx context.Context, msg *types.MsgMintDINR) (*types.MsgMintDINRResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.MintDINR(ctx, minter, msg.Collateral, msg.DinrToMint); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintDINR,
			sdk.NewAttribute(types.AttributeKeyMinter, msg.Minter),
			sdk.NewAttribute(types.AttributeKeyCollateral, msg.Collateral.String()),
			sdk.NewAttribute(types.AttributeKeyDINRMinted, msg.DinrToMint.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Minter),
		),
	})

	return &types.MsgMintDINRResponse{}, nil
}

// BurnDINR implements types.MsgServer
func (k msgServer) BurnDINR(goCtx context.Context, msg *types.MsgBurnDINR) (*types.MsgBurnDINRResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	burner, err := sdk.AccAddressFromBech32(msg.Burner)
	if err != nil {
		return nil, err
	}

	collateralReturned, err := k.Keeper.BurnDINR(ctx, burner, msg.DinrToBurn, msg.CollateralDenom)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnDINR,
			sdk.NewAttribute(types.AttributeKeyBurner, msg.Burner),
			sdk.NewAttribute(types.AttributeKeyDINRBurned, msg.DinrToBurn.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralReturned, collateralReturned.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Burner),
		),
	})

	return &types.MsgBurnDINRResponse{
		CollateralReturned: collateralReturned,
	}, nil
}

// DepositCollateral implements types.MsgServer
func (k msgServer) DepositCollateral(goCtx context.Context, msg *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.DepositCollateral(ctx, depositor, msg.Collateral); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyDepositor, msg.Depositor),
			sdk.NewAttribute(types.AttributeKeyCollateral, msg.Collateral.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
		),
	})

	return &types.MsgDepositCollateralResponse{}, nil
}

// WithdrawCollateral implements types.MsgServer
func (k msgServer) WithdrawCollateral(goCtx context.Context, msg *types.MsgWithdrawCollateral) (*types.MsgWithdrawCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	withdrawer, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.WithdrawCollateral(ctx, withdrawer, msg.Collateral); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyWithdrawer, msg.Withdrawer),
			sdk.NewAttribute(types.AttributeKeyCollateral, msg.Collateral.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Withdrawer),
		),
	})

	return &types.MsgWithdrawCollateralResponse{}, nil
}

// Liquidate implements types.MsgServer
func (k msgServer) Liquidate(goCtx context.Context, msg *types.MsgLiquidate) (*types.MsgLiquidateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	liquidator, err := sdk.AccAddressFromBech32(msg.Liquidator)
	if err != nil {
		return nil, err
	}

	user, err := sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return nil, err
	}

	collateralReceived, err := k.Keeper.Liquidate(ctx, liquidator, user, msg.DinrToCover)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLiquidate,
			sdk.NewAttribute(types.AttributeKeyLiquidator, msg.Liquidator),
			sdk.NewAttribute(types.AttributeKeyUser, msg.User),
			sdk.NewAttribute(types.AttributeKeyDINRCovered, msg.DinrToCover.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralReceived, collateralReceived.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Liquidator),
		),
	})

	return &types.MsgLiquidateResponse{
		CollateralReceived: collateralReceived,
	}, nil
}

// UpdateParams implements types.MsgServer
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority - should be governance module account
	if msg.Authority != k.authority {
		return nil, types.ErrUnauthorized
	}

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateParams,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
		),
	})

	return &types.MsgUpdateParamsResponse{}, nil
}