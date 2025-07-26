package keeper

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/identity/types"
)

// Cross-Module Identity Sharing Message Server Implementation

// CreateShareRequest handles creating a new identity sharing request
func (k msgServer) CreateShareRequest(goCtx context.Context, msg *types.MsgCreateShareRequest) (*types.MsgCreateShareRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert TTL from hours to duration
	ttl := time.Duration(msg.TTLHours) * time.Hour

	// Create the share request
	request, err := k.CreateShareRequest(
		ctx,
		msg.RequesterModule,
		msg.ProviderModule,
		msg.HolderDID,
		msg.RequestedData,
		msg.Purpose,
		msg.Justification,
		ttl,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_share_request_created",
			sdk.NewAttribute("request_id", request.RequestID),
			sdk.NewAttribute("requester_module", request.RequesterModule),
			sdk.NewAttribute("provider_module", request.ProviderModule),
			sdk.NewAttribute("holder_did", request.HolderDID),
			sdk.NewAttribute("purpose", request.Purpose),
			sdk.NewAttribute("status", request.Status.String()),
		),
	)

	return &types.MsgCreateShareRequestResponse{
		RequestID: request.RequestID,
	}, nil
}

// ApproveShareRequest handles approving a pending share request
func (k msgServer) ApproveShareRequest(goCtx context.Context, msg *types.MsgApproveShareRequest) (*types.MsgApproveShareRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the authority address
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Generate access token if not provided
	accessToken := msg.AccessToken
	if accessToken == "" {
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			return nil, err
		}
		accessToken = hex.EncodeToString(tokenBytes)
	}

	// Approve the request
	err = k.ApproveShareRequest(ctx, authority, msg.RequestID, accessToken)
	if err != nil {
		return nil, err
	}

	// Get the updated request to determine expiry
	request, found := k.GetShareRequest(ctx, msg.RequestID)
	if !found {
		return nil, types.ErrShareRequestNotFound
	}

	expiresAt := ctx.BlockTime().Add(request.TTL)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_share_request_approved",
			sdk.NewAttribute("request_id", msg.RequestID),
			sdk.NewAttribute("holder_did", request.HolderDID),
			sdk.NewAttribute("approved_by", msg.Authority),
			sdk.NewAttribute("expires_at", expiresAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgApproveShareRequestResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

// DenyShareRequest handles denying a pending share request
func (k msgServer) DenyShareRequest(goCtx context.Context, msg *types.MsgDenyShareRequest) (*types.MsgDenyShareRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the authority address
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Deny the request
	err = k.DenyShareRequest(ctx, authority, msg.RequestID, msg.DenialReason)
	if err != nil {
		return nil, err
	}

	// Get the request for event details
	request, found := k.GetShareRequest(ctx, msg.RequestID)
	if !found {
		return nil, types.ErrShareRequestNotFound
	}

	deniedAt := ctx.BlockTime()

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_share_request_denied",
			sdk.NewAttribute("request_id", msg.RequestID),
			sdk.NewAttribute("holder_did", request.HolderDID),
			sdk.NewAttribute("denied_by", msg.Authority),
			sdk.NewAttribute("denial_reason", msg.DenialReason),
			sdk.NewAttribute("denied_at", deniedAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgDenyShareRequestResponse{
		DeniedAt: deniedAt,
	}, nil
}

// CreateSharingAgreement handles creating a standing agreement between modules
func (k msgServer) CreateSharingAgreement(goCtx context.Context, msg *types.MsgCreateSharingAgreement) (*types.MsgCreateSharingAgreementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the authority address
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Convert durations
	maxTTL := time.Duration(msg.MaxTTLHours) * time.Hour
	validityPeriod := time.Duration(msg.ValidityDays) * 24 * time.Hour

	// Create the agreement
	agreement, err := k.CreateSharingAgreement(
		ctx,
		authority,
		msg.RequesterModule,
		msg.ProviderModule,
		msg.AllowedDataTypes,
		msg.Purposes,
		msg.TrustLevel,
		msg.AutoApprove,
		maxTTL,
		validityPeriod,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_sharing_agreement_created",
			sdk.NewAttribute("agreement_id", agreement.AgreementID),
			sdk.NewAttribute("requester_module", agreement.RequesterModule),
			sdk.NewAttribute("provider_module", agreement.ProviderModule),
			sdk.NewAttribute("trust_level", agreement.TrustLevel),
			sdk.NewAttribute("auto_approve", fmt.Sprintf("%t", agreement.AutoApprove)),
			sdk.NewAttribute("expires_at", agreement.ExpiresAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgCreateSharingAgreementResponse{
		AgreementID: agreement.AgreementID,
	}, nil
}

// CreateAccessPolicy handles creating an access policy for identity sharing
func (k msgServer) CreateAccessPolicy(goCtx context.Context, msg *types.MsgCreateAccessPolicy) (*types.MsgCreateAccessPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create the access policy
	policy, err := k.CreateAccessPolicy(
		ctx,
		msg.HolderDID,
		msg.AllowedModules,
		msg.DeniedModules,
		msg.DataRestrictions,
		msg.PurposeRestrictions,
		msg.TimeRestrictions,
		msg.MaxSharesPerDay,
		msg.RequireExplicitConsent,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_access_policy_created",
			sdk.NewAttribute("policy_id", policy.PolicyID),
			sdk.NewAttribute("holder_did", policy.HolderDID),
			sdk.NewAttribute("max_shares_per_day", fmt.Sprintf("%d", policy.MaxSharesPerDay)),
			sdk.NewAttribute("require_explicit_consent", fmt.Sprintf("%t", policy.RequireExplicitConsent)),
		),
	)

	return &types.MsgCreateAccessPolicyResponse{
		PolicyID: policy.PolicyID,
	}, nil
}