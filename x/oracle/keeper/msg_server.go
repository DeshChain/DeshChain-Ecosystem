package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SubmitPrice(goCtx context.Context, req *types.MsgSubmitPrice) (*types.MsgSubmitPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Submit the price using the keeper
	err := k.Keeper.SubmitPrice(ctx, req.Validator, req.Symbol, req.Price, req.Source, req.Timestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitPriceResponse{}, nil
}

func (k msgServer) SubmitExchangeRate(goCtx context.Context, req *types.MsgSubmitExchangeRate) (*types.MsgSubmitExchangeRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Submit the exchange rate using the keeper
	err := k.Keeper.SubmitExchangeRate(ctx, req.Validator, req.Base, req.Target, req.Rate, req.Source, req.Timestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitExchangeRateResponse{}, nil
}

func (k msgServer) RegisterOracleValidator(goCtx context.Context, req *types.MsgRegisterOracleValidator) (*types.MsgRegisterOracleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Register the oracle validator using the keeper
	err := k.Keeper.RegisterOracleValidator(ctx, req.Validator, req.Power, req.Description)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterOracleValidatorResponse{}, nil
}

func (k msgServer) UpdateOracleValidator(goCtx context.Context, req *types.MsgUpdateOracleValidator) (*types.MsgUpdateOracleValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get existing oracle validator
	oracleValidator, found := k.Keeper.GetOracleValidator(ctx, req.Validator)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	// Update the oracle validator
	oracleValidator.Power = req.Power
	oracleValidator.Active = req.Active

	k.Keeper.SetOracleValidator(ctx, oracleValidator)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorUpdated,
			sdk.NewAttribute(types.AttributeKeyValidator, req.Validator),
			sdk.NewAttribute(types.AttributeKeyPower, sdk.NewInt(int64(req.Power)).String()),
			sdk.NewAttribute(types.AttributeKeyActive, sdk.FormatBool(req.Active)),
		),
	)

	return &types.MsgUpdateOracleValidatorResponse{}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate authority (typically governance module account)
	if k.GetAuthority() != req.Authority {
		return nil, types.ErrInvalidSigner
	}

	// Update parameters
	k.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}