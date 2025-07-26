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
	"fmt"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/governance/types"
)

// EnforceFounderProtection ensures all founder protections are enforced
func (k Keeper) EnforceFounderProtection(ctx sdk.Context, proposal types.Proposal) error {
	// Get founder address from params
	params := k.GetParams(ctx)
	founderAddr := params.FounderAddress
	
	// Check if proposal affects protected parameters
	for _, param := range proposal.AffectsProtectedParams {
		protectionLevel := types.GetProtectionLevel(param)
		
		switch protectionLevel {
		case types.ProtectionType_PROTECTION_TYPE_IMMUTABLE:
			// These parameters can NEVER be changed
			return errors.Wrapf(types.ErrImmutableParameter, 
				"parameter %s is immutable and cannot be changed by any governance proposal", param)
			
		case types.ProtectionType_PROTECTION_TYPE_FOUNDER_CONSENT:
			// These require founder approval
			if proposal.FounderApprovalStatus != types.FounderApprovalStatus_FOUNDER_APPROVAL_STATUS_APPROVED {
				return errors.Wrapf(types.ErrFounderConsentRequired,
					"parameter %s requires founder approval", param)
			}
			
		case types.ProtectionType_PROTECTION_TYPE_SUPERMAJORITY:
			// These require 80% supermajority
			// This check happens during vote tallying
			proposal.RequiresFounderApproval = false // Mark for supermajority check
		}
	}
	
	// Check if founder can veto this proposal
	genesisTime := k.GetGenesisTime(ctx)
	currentTime := ctx.BlockTime()
	
	if types.CanFounderVeto(proposal.ProposalType, currentTime, genesisTime) {
		// Founder has veto power for this proposal type
		k.Logger(ctx).Info("Proposal can be vetoed by founder", 
			"proposal_id", proposal.ProposalId,
			"type", proposal.ProposalType.String())
	}
	
	return nil
}

// ProcessFounderVeto handles founder veto action
func (k Keeper) ProcessFounderVeto(ctx sdk.Context, proposalID uint64, vetoAddr string) error {
	params := k.GetParams(ctx)
	
	// Only founder can veto
	if vetoAddr != params.FounderAddress {
		return errors.Wrapf(types.ErrUnauthorized, 
			"only founder (%s) can veto proposals", params.FounderAddress)
	}
	
	// Get the proposal
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		return errors.Wrapf(types.ErrProposalNotFound, "proposal %d not found", proposalID)
	}
	
	// Check if veto period is still active
	genesisTime := k.GetGenesisTime(ctx)
	if !types.CanFounderVeto(proposal.ProposalType, ctx.BlockTime(), genesisTime) {
		return errors.Wrap(types.ErrVetoPeriodExpired, "founder veto period has expired")
	}
	
	// Check proposal status
	if proposal.Status != types.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD {
		return errors.Wrap(types.ErrInvalidProposalStatus, "can only veto proposals in voting period")
	}
	
	// Apply veto
	proposal.Status = types.ProposalStatus_PROPOSAL_STATUS_VETOED
	proposal.FounderVetoUsed = true
	k.SetProposal(ctx, proposal)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFounderVeto,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyVetoBy, vetoAddr),
		),
	)
	
	k.Logger(ctx).Info("Founder vetoed proposal", 
		"proposal_id", proposalID,
		"vetoed_by", vetoAddr)
	
	return nil
}

// CalculateVotingPower calculates voting power ensuring founder minimum
func (k Keeper) CalculateVotingPower(ctx sdk.Context, voterAddr string) sdk.Int {
	params := k.GetParams(ctx)
	
	// Get actual token-based voting power
	actualPower := k.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(voterAddr))
	if actualPower == nil {
		actualPower = k.bankKeeper.GetBalance(ctx, sdk.AccAddress(voterAddr), "namo").Amount
	}
	
	// If this is the founder, ensure minimum voting power
	if voterAddr == params.FounderAddress {
		totalPower := k.GetTotalVotingPower(ctx)
		minPower := types.CalculateFounderVotingPower(totalPower, actualPower)
		return minPower
	}
	
	return actualPower
}

// CheckSupermajorityRequirement checks if proposal meets supermajority requirement
func (k Keeper) CheckSupermajorityRequirement(ctx sdk.Context, proposal types.Proposal, tally types.TallyResult) bool {
	// Check if this proposal requires supermajority
	requiresSupermajority := types.RequiresSupermajority(proposal.ProposalType, proposal.AffectsProtectedParams)
	
	if !requiresSupermajority {
		// Normal majority (>50%) is sufficient
		totalVotes := tally.YesCount.Add(tally.NoCount).Add(tally.AbstainCount).Add(tally.NoWithVetoCount)
		return tally.YesCount.GT(totalVotes.QuoRaw(2))
	}
	
	// Supermajority required (80%)
	totalVotes := tally.YesCount.Add(tally.NoCount).Add(tally.AbstainCount).Add(tally.NoWithVetoCount)
	requiredVotes := totalVotes.MulRaw(types.SupermajorityThresholdPercent).QuoRaw(100)
	
	return tally.YesCount.GTE(requiredVotes)
}

// ExecuteEmergencyAction allows founder to take emergency actions
func (k Keeper) ExecuteEmergencyAction(ctx sdk.Context, actionType types.EmergencyActionType, executorAddr string, description string) error {
	params := k.GetParams(ctx)
	
	// Validate the action
	err := types.ValidateEmergencyAction(actionType, params.FounderAddress, executorAddr)
	if err != nil {
		return err
	}
	
	// Create emergency action record
	actionID := k.GetNextEmergencyActionID(ctx)
	action := types.EmergencyAction{
		ActionId:      actionID,
		ActionType:    actionType,
		Description:   description,
		ExecutedBy:    executorAddr,
		ExecutionTime: ctx.BlockTime(),
	}
	
	// Execute the action
	switch actionType {
	case types.EmergencyActionType_EMERGENCY_ACTION_TYPE_HALT_CHAIN:
		// Halt the chain
		k.crisisKeeper.AssertInvariants(ctx)
		
	case types.EmergencyActionType_EMERGENCY_ACTION_TYPE_FREEZE_MODULE:
		// Freeze specified module
		// Implementation depends on module architecture
		
	case types.EmergencyActionType_EMERGENCY_ACTION_TYPE_ROLLBACK_UPGRADE:
		// Rollback recent upgrade
		// Implementation depends on upgrade module
		
	case types.EmergencyActionType_EMERGENCY_ACTION_TYPE_PATCH_VULNERABILITY:
		// Apply security patch
		// Implementation depends on specific vulnerability
	}
	
	// Store the action
	k.SetEmergencyAction(ctx, action)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmergencyAction,
			sdk.NewAttribute(types.AttributeKeyActionID, fmt.Sprintf("%d", actionID)),
			sdk.NewAttribute(types.AttributeKeyActionType, actionType.String()),
			sdk.NewAttribute(types.AttributeKeyExecutedBy, executorAddr),
		),
	)
	
	k.Logger(ctx).Info("Emergency action executed by founder",
		"action_id", actionID,
		"type", actionType.String(),
		"executor", executorAddr)
	
	return nil
}

// ValidateProposalSubmission validates proposals don't violate protections
func (k Keeper) ValidateProposalSubmission(ctx sdk.Context, proposal types.Proposal) error {
	// Check each affected parameter
	for _, param := range proposal.AffectsProtectedParams {
		// Get protection level
		protection := types.GetProtectionLevel(param)
		
		// Immutable parameters cannot have proposals
		if protection == types.ProtectionType_PROTECTION_TYPE_IMMUTABLE {
			return errors.Wrapf(types.ErrImmutableParameter,
				"cannot create proposal affecting immutable parameter: %s", param)
		}
		
		// Mark if founder consent needed
		if protection == types.ProtectionType_PROTECTION_TYPE_FOUNDER_CONSENT {
			proposal.RequiresFounderApproval = true
		}
	}
	
	// Special validation for founder-related proposals
	if proposal.ProposalType == types.ProposalType_PROPOSAL_TYPE_FOUNDER_RELATED {
		// These always require founder approval
		proposal.RequiresFounderApproval = true
		
		// Check notice period
		params := k.GetParams(ctx)
		noticeEnd := ctx.BlockTime().Add(params.GovernanceChangeNoticePeriod)
		if proposal.VotingStartTime.Before(noticeEnd) {
			return errors.Wrapf(types.ErrNoticePeriodViolation,
				"governance changes affecting founder require %d day notice period",
				types.GovernanceChangeNoticeDays)
		}
	}
	
	return nil
}

// ProcessFounderApproval handles founder's decision on proposals requiring consent
func (k Keeper) ProcessFounderApproval(ctx sdk.Context, proposalID uint64, approverAddr string, decision types.FounderApprovalStatus) error {
	params := k.GetParams(ctx)
	
	// Only founder can approve/reject
	if approverAddr != params.FounderAddress {
		return errors.Wrapf(types.ErrUnauthorized,
			"only founder (%s) can approve/reject proposals requiring consent", params.FounderAddress)
	}
	
	// Get the proposal
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		return errors.Wrapf(types.ErrProposalNotFound, "proposal %d not found", proposalID)
	}
	
	// Check if approval is required
	if !proposal.RequiresFounderApproval {
		return errors.Wrap(types.ErrInvalidProposal, "proposal does not require founder approval")
	}
	
	// Update approval status
	proposal.FounderApprovalStatus = decision
	k.SetProposal(ctx, proposal)
	
	// If rejected, end the proposal
	if decision == types.FounderApprovalStatus_FOUNDER_APPROVAL_STATUS_REJECTED {
		proposal.Status = types.ProposalStatus_PROPOSAL_STATUS_REJECTED
		k.SetProposal(ctx, proposal)
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFounderApproval,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyApprovalStatus, decision.String()),
			sdk.NewAttribute(types.AttributeKeyApprovedBy, approverAddr),
		),
	)
	
	return nil
}

// GetFounderRights returns the current founder rights status
func (k Keeper) GetFounderRights(ctx sdk.Context) types.FounderRights {
	params := k.GetParams(ctx)
	genesisTime := k.GetGenesisTime(ctx)
	currentTime := ctx.BlockTime()
	
	// Calculate veto expiry
	vetoExpiry := genesisTime.Add(time.Duration(types.FounderVetoDurationYears) * 365 * 24 * time.Hour)
	vetoPower := currentTime.Before(vetoExpiry)
	
	// Technical authority for 3 years
	technicalAuthority := currentTime.Before(vetoExpiry)
	
	return types.FounderRights{
		GuaranteedVotingPower:     fmt.Sprintf("%d%%", types.FounderMinimumVotingPowerPercent),
		VetoPower:                 vetoPower,
		VetoExpiry:                vetoExpiry,
		ProtectedAllocations:      []string{"10% token allocation", "48-month vesting"},
		ProtectedRoyalties:        []string{"0.10% tax royalty", "5% platform royalty"},
		EmergencyPowers:           true, // Always available
		TechnicalDecisionAuthority: technicalAuthority,
	}
}

// GetImmutableProtections returns the immutable protection status
func (k Keeper) GetImmutableProtections(ctx sdk.Context) types.ImmutableFounderProtections {
	// These are hardcoded as true - they can NEVER be changed
	return types.ImmutableFounderProtections{
		FounderTokenAllocationProtected:  true,
		FounderTaxRoyaltyProtected:      true,
		FounderPlatformRoyaltyProtected: true,
		InheritanceMechanismProtected:   true,
		MinimumVotingPowerProtected:     true,
		ProtectionRemovalForbidden:      true, // This protection itself cannot be removed
	}
}

// Override50YearPlan allows founder to override decisions that deviate from 50-year plan
func (k Keeper) Override50YearPlan(ctx sdk.Context, overriderAddr string, proposalID uint64, reason string) error {
	params := k.GetParams(ctx)
	
	// Only founder can override
	if overriderAddr != params.FounderAddress {
		return errors.Wrapf(types.ErrUnauthorized,
			"only founder can override decisions deviating from 50-year plan")
	}
	
	// Get the proposal
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		return errors.Wrapf(types.ErrProposalNotFound, "proposal %d not found", proposalID)
	}
	
	// Override the proposal
	proposal.Status = types.ProposalStatus_PROPOSAL_STATUS_VETOED
	proposal.FounderVetoUsed = true
	k.SetProposal(ctx, proposal)
	
	// Log the override reason
	k.Logger(ctx).Info("Founder overrode proposal to maintain 50-year plan",
		"proposal_id", proposalID,
		"reason", reason)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventType50YearPlanOverride,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
			sdk.NewAttribute(types.AttributeKeyOverrideBy, overriderAddr),
		),
	)
	
	return nil
}