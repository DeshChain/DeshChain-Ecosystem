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

package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/deshchain/namo/x/governance/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the governance MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// VetoProposal implements the veto proposal handler
func (k msgServer) VetoProposal(goCtx context.Context, msg *types.MsgVetoProposal) (*types.MsgVetoProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify sender is the founder
	founderAddr := k.GetFounderAddress(ctx)
	if msg.Authority != founderAddr {
		return nil, sdkerrors.Wrapf(types.ErrNotFounder, "only founder can veto proposals")
	}

	// Get current phase
	phase := k.GetGovernancePhase(ctx)
	if phase != types.GovernancePhase_FOUNDER_CONTROL && phase != types.GovernancePhase_SHARED_GOVERNANCE {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPhase, "veto not allowed in phase %s", phase.String())
	}

	// Get the proposal
	proposal, found := k.govKeeper.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotFound, "proposal %d not found", msg.ProposalId)
	}

	// Check if proposal is still active
	if proposal.Status != govv1.StatusVotingPeriod {
		return nil, sdkerrors.Wrapf(types.ErrInvalidProposalStatus, "proposal %d is not in voting period", msg.ProposalId)
	}

	// Mark proposal as vetoed
	k.SetProposalVetoed(ctx, msg.ProposalId)

	// Update proposal status to failed
	proposal.Status = govv1.StatusFailed
	k.govKeeper.SetProposal(ctx, proposal)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"proposal_vetoed",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("founder", msg.Authority),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)

	k.Logger(ctx).Info("Proposal vetoed",
		"proposal_id", msg.ProposalId,
		"founder", msg.Authority,
		"reason", msg.Reason)

	return &types.MsgVetoProposalResponse{}, nil
}

// ApproveFounderConsentProposal implements the founder consent handler
func (k msgServer) ApproveFounderConsentProposal(goCtx context.Context, msg *types.MsgApproveFounderConsentProposal) (*types.MsgApproveFounderConsentProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify sender is the founder
	founderAddr := k.GetFounderAddress(ctx)
	if msg.Authority != founderAddr {
		return nil, sdkerrors.Wrapf(types.ErrNotFounder, "only founder can approve consent proposals")
	}

	// Get the proposal
	proposal, found := k.govKeeper.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotFound, "proposal %d not found", msg.ProposalId)
	}

	// Check if proposal requires founder consent
	requiresConsent := false
	for _, message := range proposal.Messages {
		if k.requiresFounderConsent(ctx, message) {
			requiresConsent = true
			break
		}
	}

	if !requiresConsent {
		return nil, sdkerrors.Wrapf(types.ErrNoConsentRequired, "proposal %d does not require founder consent", msg.ProposalId)
	}

	// Mark proposal as approved by founder
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetFounderApprovalKey(msg.ProposalId), []byte{1})

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"proposal_founder_approved",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("founder", msg.Authority),
		),
	)

	k.Logger(ctx).Info("Proposal approved by founder",
		"proposal_id", msg.ProposalId,
		"founder", msg.Authority)

	return &types.MsgApproveFounderConsentProposalResponse{}, nil
}

// UpdateProtectedParameter implements the protected parameter update handler
func (k msgServer) UpdateProtectedParameter(goCtx context.Context, msg *types.MsgUpdateProtectedParameter) (*types.MsgUpdateProtectedParameterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority (must be governance module)
	if msg.Authority != k.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Get current parameter
	param, found := k.GetProtectedParameter(ctx, msg.Name)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrParameterNotFound, "parameter %s not found", msg.Name)
	}

	// Check if parameter is immutable
	if param.Protection == types.ProtectionType_IMMUTABLE {
		return nil, sdkerrors.Wrapf(types.ErrImmutableParameter, "parameter %s is immutable", msg.Name)
	}

	// Update protection level
	param.Protection = msg.NewProtection
	k.SetProtectedParameter(ctx, param)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"protected_parameter_updated",
			sdk.NewAttribute("name", msg.Name),
			sdk.NewAttribute("new_protection", msg.NewProtection.String()),
		),
	)

	k.Logger(ctx).Info("Protected parameter updated",
		"name", msg.Name,
		"new_protection", msg.NewProtection.String())

	return &types.MsgUpdateProtectedParameterResponse{}, nil
}

// Helper function to check if a proposal message requires founder consent
func (k msgServer) requiresFounderConsent(ctx sdk.Context, msg *sdk.Any) bool {
	// Check if the message affects protected parameters
	protectedParams := []string{
		"chain_upgrade_handler",
		"crisis_module_permissions",
		"slashing_parameters",
		"consensus_parameters",
	}

	// Extract message type and check against protected operations
	msgTypeURL := msg.TypeUrl
	
	// Check for parameter change proposals
	if msgTypeURL == "/cosmos.params.v1beta1.ParameterChangeProposal" {
		return true // All parameter changes need review
	}

	// Check for software upgrade proposals
	if msgTypeURL == "/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal" {
		return true
	}

	// Check for other critical operations
	criticalOperations := []string{
		"/cosmos.slashing.v1beta1.MsgUpdateParams",
		"/cosmos.consensus.v1.MsgUpdateParams",
		"/cosmos.crisis.v1beta1.MsgUpdateParams",
	}

	for _, op := range criticalOperations {
		if msgTypeURL == op {
			return true
		}
	}

	return false
}