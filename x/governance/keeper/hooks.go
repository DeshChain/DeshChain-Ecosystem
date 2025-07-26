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

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/DeshChain/DeshChain-Ecosystem/x/governance/types"
)

// Hooks is the governance hooks wrapper
type Hooks struct {
	k Keeper
}

var _ govtypes.GovHooks = Hooks{}

// NewHooks creates new governance hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterProposalSubmission is called after a proposal is submitted
func (h Hooks) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
	// Get the proposal
	proposal, found := h.k.govKeeper.GetProposal(ctx, proposalID)
	if !found {
		return
	}

	// Check if any messages are trying to modify immutable parameters
	for _, msg := range proposal.Messages {
		if h.affectsImmutableParameters(ctx, msg) {
			// Immediately veto the proposal
			h.k.SetProposalVetoed(ctx, proposalID)
			proposal.Status = govv1.StatusFailed
			h.k.govKeeper.SetProposal(ctx, proposal)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"proposal_auto_vetoed",
					sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalID)),
					sdk.NewAttribute("reason", "affects_immutable_parameters"),
				),
			)

			h.k.Logger(ctx).Info("Proposal auto-vetoed for affecting immutable parameters",
				"proposal_id", proposalID)
			return
		}
	}

	// Log proposal submission
	h.k.Logger(ctx).Info("Proposal submitted",
		"proposal_id", proposalID,
		"messages", len(proposal.Messages))
}

// AfterProposalDeposit is called after a deposit is made
func (h Hooks) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
	// No special handling needed
}

// AfterProposalVote is called after a vote is cast
func (h Hooks) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	// Check if proposal has been vetoed
	if h.k.IsProposalVetoed(ctx, proposalID) {
		// Log that vote was cast on vetoed proposal
		h.k.Logger(ctx).Info("Vote cast on vetoed proposal",
			"proposal_id", proposalID,
			"voter", voterAddr.String())
	}
}

// AfterProposalFailedMinDeposit is called after a proposal fails to meet min deposit
func (h Hooks) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64) {
	// No special handling needed
}

// AfterProposalVotingPeriodEnded is called after the voting period ends
func (h Hooks) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	// Check if proposal was vetoed
	if h.k.IsProposalVetoed(ctx, proposalID) {
		h.k.Logger(ctx).Info("Vetoed proposal voting period ended",
			"proposal_id", proposalID)
		return
	}

	// Get the proposal
	proposal, found := h.k.govKeeper.GetProposal(ctx, proposalID)
	if !found {
		return
	}

	// Check if proposal requires founder consent
	requiresConsent := false
	for _, msg := range proposal.Messages {
		if h.requiresFounderConsent(ctx, msg) {
			requiresConsent = true
			break
		}
	}

	if requiresConsent {
		// Check if founder has approved
		store := ctx.KVStore(h.k.storeKey)
		if !store.Has(types.GetFounderApprovalKey(proposalID)) {
			// Fail the proposal
			proposal.Status = govv1.StatusFailed
			h.k.govKeeper.SetProposal(ctx, proposal)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"proposal_failed_no_founder_consent",
					sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalID)),
				),
			)

			h.k.Logger(ctx).Info("Proposal failed due to missing founder consent",
				"proposal_id", proposalID)
		}
	}

	// Check if proposal requires supermajority
	if h.requiresSupermajority(ctx, proposal) {
		tally := h.k.govKeeper.GetTallyResult(ctx, proposal)
		totalVotes := tally.YesCount.Add(tally.NoCount).Add(tally.NoWithVetoCount).Add(tally.AbstainCount)
		
		if totalVotes.IsZero() {
			return
		}

		// Calculate yes percentage
		yesPercentage := tally.YesCount.MulRaw(100).Quo(totalVotes)
		
		// Require 80% supermajority
		if yesPercentage.LT(sdk.NewInt(80)) {
			proposal.Status = govv1.StatusFailed
			h.k.govKeeper.SetProposal(ctx, proposal)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"proposal_failed_supermajority",
					sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalID)),
					sdk.NewAttribute("yes_percentage", yesPercentage.String()),
				),
			)

			h.k.Logger(ctx).Info("Proposal failed to meet supermajority requirement",
				"proposal_id", proposalID,
				"yes_percentage", yesPercentage.String())
		}
	}
}

// Helper function to check if a message affects immutable parameters
func (h Hooks) affectsImmutableParameters(ctx sdk.Context, msg *sdk.Any) bool {
	immutableParams := []string{
		"founder_token_allocation",
		"founder_tax_royalty",
		"founder_platform_royalty",
		"founder_inheritance_mechanism",
		"founder_minimum_voting_power",
	}

	// Check message type and content
	msgTypeURL := msg.TypeUrl
	
	// For parameter change proposals, check each parameter
	if msgTypeURL == "/cosmos.params.v1beta1.ParameterChangeProposal" {
		// Would need to unmarshal and check each parameter
		// For now, we'll be conservative and check manually
		return true // Requires manual review
	}

	return false
}

// Helper function to check if a proposal message requires founder consent
func (h Hooks) requiresFounderConsent(ctx sdk.Context, msg *sdk.Any) bool {
	// Check if the message affects protected parameters
	founderConsentParams := []string{
		"chain_upgrade_handler",
		"crisis_module_permissions",
		"slashing_parameters",
		"consensus_parameters",
	}

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

// Helper function to check if a proposal requires supermajority
func (h Hooks) requiresSupermajority(ctx sdk.Context, proposal govv1.Proposal) bool {
	supermajorityParams := []string{
		"governance_voting_period",
		"governance_deposit_amount",
		"distribution_community_tax",
		"mint_inflation_rate",
	}

	// Check each message in the proposal
	for _, msg := range proposal.Messages {
		msgTypeURL := msg.TypeUrl
		
		// Check for parameter changes that require supermajority
		if msgTypeURL == "/cosmos.params.v1beta1.ParameterChangeProposal" {
			// Would need to unmarshal and check against supermajorityParams
			return true // For now, be conservative
		}

		// Check for distribution parameter updates
		if msgTypeURL == "/cosmos.distribution.v1beta1.MsgCommunityPoolSpend" {
			return true
		}

		// Check for mint parameter updates
		if msgTypeURL == "/cosmos.mint.v1beta1.MsgUpdateParams" {
			return true
		}
	}

	return false
}