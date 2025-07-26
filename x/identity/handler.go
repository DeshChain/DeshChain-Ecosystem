package identity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/deshchain/x/identity/keeper"
	"github.com/deshchain/x/identity/types"
)

// NewHandler creates an sdk.Handler for all the identity module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateIdentity:
			res, err := msgServer.CreateIdentity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgUpdateIdentity:
			res, err := msgServer.UpdateIdentity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgRevokeIdentity:
			res, err := msgServer.RevokeIdentity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgRegisterDID:
			res, err := msgServer.RegisterDID(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgUpdateDID:
			res, err := msgServer.UpdateDID(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgDeactivateDID:
			res, err := msgServer.DeactivateDID(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgIssueCredential:
			res, err := msgServer.IssueCredential(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgRevokeCredential:
			res, err := msgServer.RevokeCredential(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgPresentCredential:
			res, err := msgServer.PresentCredential(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgCreateZKProof:
			res, err := msgServer.CreateZKProof(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgVerifyZKProof:
			res, err := msgServer.VerifyZKProof(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgLinkAadhaar:
			res, err := msgServer.LinkAadhaar(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgConnectDigiLocker:
			res, err := msgServer.ConnectDigiLocker(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgLinkUPI:
			res, err := msgServer.LinkUPI(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgGiveConsent:
			res, err := msgServer.GiveConsent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgWithdrawConsent:
			res, err := msgServer.WithdrawConsent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgAddRecoveryMethod:
			res, err := msgServer.AddRecoveryMethod(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgInitiateRecovery:
			res, err := msgServer.InitiateRecovery(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgCompleteRecovery:
			res, err := msgServer.CompleteRecovery(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			
		case *types.MsgUpdatePrivacySettings:
			res, err := msgServer.UpdatePrivacySettings(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}