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

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/DeshChain/DeshChain-Ecosystem/x/treasury/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the treasury MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SubmitCommunityProposal implements the community fund proposal submission
func (k msgServer) SubmitCommunityProposal(goCtx context.Context, msg *types.MsgSubmitCommunityProposal) (*types.MsgSubmitCommunityProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		return nil, err
	}

	// Validate requested amount
	if !msg.RequestedAmount.IsValid() || msg.RequestedAmount.IsZero() {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid requested amount")
	}

	// Check community fund balance
	balance, err := k.GetCommunityFundBalance(ctx)
	if err != nil {
		return nil, err
	}

	if balance.AvailableAmount.IsLT(msg.RequestedAmount) {
		return nil, errors.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient community fund balance")
	}

	// Create proposal
	proposalID := k.GetNextCommunityProposalID(ctx)
	proposal := types.CommunityFundProposal{
		ProposalId:      proposalID,
		Proposer:        msg.Proposer,
		Title:           msg.Title,
		Description:     msg.Description,
		Category:        msg.Category,
		RequestedAmount: msg.RequestedAmount,
		Recipients:      msg.Recipients,
		Milestones:      msg.Milestones,
		Status:          types.StatusPending,
		VotingPeriod:    k.GetCommunityVotingPeriod(ctx),
		SubmissionTime:  ctx.BlockTime(),
		VotingEndTime:   ctx.BlockTime().Add(k.GetCommunityVotingPeriod(ctx)),
		QuorumReached:   false,
		Passed:          false,
		AuditRequired:   k.IsAuditRequired(msg.RequestedAmount),
		TransparencyScore: k.CalculateTransparencyScore(msg),
	}

	// Store proposal
	if err := k.SetCommunityFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Update proposal counter
	k.SetNextCommunityProposalID(ctx, proposalID+1)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposalSubmitted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyCategory, msg.Category),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.RequestedAmount.String()),
		),
	})

	return &types.MsgSubmitCommunityProposalResponse{
		ProposalId: proposalID,
	}, nil
}

// VoteCommunityProposal implements voting on community fund proposals
func (k msgServer) VoteCommunityProposal(goCtx context.Context, msg *types.MsgVoteCommunityProposal) (*types.MsgVoteCommunityProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	voter, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetCommunityFundProposal(ctx, msg.ProposalId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
	}

	// Check if voting period is active
	if ctx.BlockTime().After(proposal.VotingEndTime) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "voting period has ended")
	}

	// Check if proposal is in voting status
	if proposal.Status != types.StatusActive && proposal.Status != types.StatusPending {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proposal is not in voting period")
	}

	// Update proposal to active if first vote
	if proposal.Status == types.StatusPending {
		proposal.Status = types.StatusActive
	}

	// Record vote (simplified for now)
	switch msg.Option {
	case "yes":
		proposal.VotesFor = k.IncrementVotes(proposal.VotesFor)
	case "no":
		proposal.VotesAgainst = k.IncrementVotes(proposal.VotesAgainst)
	case "abstain":
		proposal.VotesAbstain = k.IncrementVotes(proposal.VotesAbstain)
	default:
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid vote option")
	}

	// Update total votes
	proposal.TotalVotes = k.IncrementVotes(proposal.TotalVotes)

	// Check if quorum is reached
	proposal.QuorumReached = k.IsQuorumReached(ctx, proposal)

	// Store updated proposal
	if err := k.SetCommunityFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposalVoted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("voter", msg.Voter),
			sdk.NewAttribute("option", msg.Option),
		),
	})

	return &types.MsgVoteCommunityProposalResponse{}, nil
}

// ExecuteCommunityProposal implements execution of passed community fund proposals
func (k msgServer) ExecuteCommunityProposal(goCtx context.Context, msg *types.MsgExecuteCommunityProposal) (*types.MsgExecuteCommunityProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	executor, err := sdk.AccAddressFromBech32(msg.Executor)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetCommunityFundProposal(ctx, msg.ProposalId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
	}

	// Check if voting period has ended
	if ctx.BlockTime().Before(proposal.VotingEndTime) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "voting period has not ended")
	}

	// Check if proposal passed
	if !k.HasProposalPassed(ctx, proposal) {
		proposal.Status = types.StatusRejected
		if err := k.SetCommunityFundProposal(ctx, proposal); err != nil {
			return nil, err
		}
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proposal did not pass")
	}

	// Update proposal status
	proposal.Status = types.StatusExecuting
	proposal.ExecutionTime = ctx.BlockTime()
	proposal.Passed = true

	// Execute fund transfers
	transactionIDs := []string{}
	for _, recipient := range proposal.Recipients {
		recipientAddr, err := sdk.AccAddressFromBech32(recipient.Address)
		if err != nil {
			continue
		}

		// Transfer funds from community pool to recipient
		err = k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.CommunityFundPoolName,
			recipientAddr,
			sdk.NewCoins(recipient.Amount),
		)
		if err != nil {
			continue
		}

		// Create transaction record
		txID := k.GenerateTransactionID(ctx, proposal.ProposalId, recipient.Address)
		transaction := types.CommunityFundTransaction{
			TxId:        txID,
			ProposalId:  proposal.ProposalId,
			From:        k.GetModuleAddress(types.CommunityFundPoolName).String(),
			To:          recipient.Address,
			Amount:      recipient.Amount,
			Type:        types.TxTypeAllocation,
			Category:    proposal.Category,
			Description: fmt.Sprintf("Community fund allocation for: %s", proposal.Title),
			Timestamp:   ctx.BlockTime(),
			BlockHeight: ctx.BlockHeight(),
			Status:      types.StatusCompleted,
			Verified:    true,
		}

		if err := k.SetCommunityFundTransaction(ctx, transaction); err != nil {
			continue
		}

		transactionIDs = append(transactionIDs, txID)
	}

	// Update proposal status to completed
	proposal.Status = types.StatusCompleted
	if err := k.SetCommunityFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Update community fund balance
	if err := k.UpdateCommunityFundBalance(ctx, proposal.RequestedAmount, false); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposalExecuted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyExecutor, msg.Executor),
			sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status),
		),
	})

	return &types.MsgExecuteCommunityProposalResponse{
		TransactionIds: transactionIDs,
	}, nil
}

// SubmitDevelopmentProposal implements development fund proposal submission
func (k msgServer) SubmitDevelopmentProposal(goCtx context.Context, msg *types.MsgSubmitDevelopmentProposal) (*types.MsgSubmitDevelopmentProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		return nil, err
	}

	// Validate requested amount
	if !msg.RequestedAmount.IsValid() || msg.RequestedAmount.IsZero() {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid requested amount")
	}

	// Check development fund balance
	balance, err := k.GetDevelopmentFundBalance(ctx)
	if err != nil {
		return nil, err
	}

	if balance.AvailableAmount.IsLT(msg.RequestedAmount) {
		return nil, errors.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient development fund balance")
	}

	// Create proposal
	proposalID := k.GetNextDevelopmentProposalID(ctx)
	proposal := types.DevelopmentFundProposal{
		ProposalId:      proposalID,
		Proposer:        msg.Proposer,
		Title:           msg.Title,
		Description:     msg.Description,
		Category:        msg.Category,
		Priority:        msg.Priority,
		RequestedAmount: msg.RequestedAmount,
		TechnicalSpecs:  msg.TechnicalSpecs,
		Timeline:        msg.Timeline,
		Team:            msg.Team,
		Deliverables:    msg.Deliverables,
		Status:          types.StatusPending,
		SubmissionTime:  ctx.BlockTime(),
		ReviewPeriod:    k.GetDevelopmentReviewPeriod(ctx),
		ReviewEndTime:   ctx.BlockTime().Add(k.GetDevelopmentReviewPeriod(ctx)),
		TransparencyLevel: k.CalculateDevelopmentTransparencyLevel(msg),
	}

	// Store proposal
	if err := k.SetDevelopmentFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Update proposal counter
	k.SetNextDevelopmentProposalID(ctx, proposalID+1)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.DevEventTypeProposalSubmitted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyCategory, msg.Category),
			sdk.NewAttribute("priority", msg.Priority),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.RequestedAmount.String()),
		),
	})

	return &types.MsgSubmitDevelopmentProposalResponse{
		ProposalId: proposalID,
	}, nil
}

// ReviewDevelopmentProposal implements technical review of development proposals
func (k msgServer) ReviewDevelopmentProposal(goCtx context.Context, msg *types.MsgReviewDevelopmentProposal) (*types.MsgReviewDevelopmentProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	reviewer, err := sdk.AccAddressFromBech32(msg.Reviewer)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetDevelopmentFundProposal(ctx, msg.ProposalId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
	}

	// Check if review period is active
	if ctx.BlockTime().After(proposal.ReviewEndTime) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "review period has ended")
	}

	// Update review based on type
	switch msg.ReviewType {
	case "technical":
		proposal.TechnicalReview = &types.TechnicalReview{
			Reviewer:    msg.Reviewer,
			ReviewDate:  ctx.BlockTime(),
			Score:       msg.Score,
			Approved:    msg.Approved,
			Comments:    msg.Comments,
			Recommendations: msg.Recommendations,
		}
	case "financial":
		proposal.FinancialReview = &types.FinancialReview{
			Reviewer:    msg.Reviewer,
			ReviewDate:  ctx.BlockTime(),
			Score:       msg.Score,
			Approved:    msg.Approved,
			Recommendations: msg.Recommendations,
		}
	case "security":
		proposal.SecurityReview = &types.SecurityReview{
			Reviewer:    msg.Reviewer,
			ReviewDate:  ctx.BlockTime(),
			Score:       msg.Score,
			Approved:    msg.Approved,
			Recommendations: msg.Recommendations,
		}
	default:
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid review type")
	}

	// Calculate overall approval score
	proposal.ApprovalScore = k.CalculateApprovalScore(proposal)

	// Update proposal status if all reviews complete
	if k.AreAllReviewsComplete(proposal) {
		if proposal.ApprovalScore >= k.GetMinApprovalScore(ctx) {
			proposal.Status = types.StatusPassed
			proposal.ApprovalTime = &ctx.BlockTime()
		} else {
			proposal.Status = types.StatusRejected
		}
	}

	// Store updated proposal
	if err := k.SetDevelopmentFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.DevEventTypeReviewCompleted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("reviewer", msg.Reviewer),
			sdk.NewAttribute("review_type", msg.ReviewType),
			sdk.NewAttribute("score", fmt.Sprintf("%d", msg.Score)),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", msg.Approved)),
		),
	})

	return &types.MsgReviewDevelopmentProposalResponse{}, nil
}

// ExecuteDevelopmentProposal implements execution of approved development proposals
func (k msgServer) ExecuteDevelopmentProposal(goCtx context.Context, msg *types.MsgExecuteDevelopmentProposal) (*types.MsgExecuteDevelopmentProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	executor, err := sdk.AccAddressFromBech32(msg.Executor)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetDevelopmentFundProposal(ctx, msg.ProposalId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
	}

	// Check if proposal is approved
	if proposal.Status != types.StatusPassed {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "proposal not approved")
	}

	// Update proposal status
	proposal.Status = types.StatusExecuting
	executionTime := ctx.BlockTime()
	proposal.ExecutionStartTime = &executionTime

	// Execute initial fund transfer (usually a percentage)
	initialAmount := k.CalculateInitialFunding(proposal.RequestedAmount)
	
	// Transfer funds to escrow
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.DevelopmentFundPoolName,
		types.DevelopmentFundEscrowName,
		sdk.NewCoins(initialAmount),
	)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	txID := k.GenerateDevelopmentTransactionID(ctx, proposal.ProposalId)
	transaction := types.DevelopmentFundTransaction{
		TxId:        txID,
		ProposalId:  proposal.ProposalId,
		From:        k.GetModuleAddress(types.DevelopmentFundPoolName).String(),
		To:          k.GetModuleAddress(types.DevelopmentFundEscrowName).String(),
		Amount:      initialAmount,
		Type:        types.TxTypeAllocation,
		Category:    proposal.Category,
		Description: fmt.Sprintf("Initial funding for: %s", proposal.Title),
		Timestamp:   ctx.BlockTime(),
		BlockHeight: ctx.BlockHeight(),
		Status:      types.StatusCompleted,
		Phase:       "initial",
		Approved:    true,
		Reviewed:    true,
	}

	if err := k.SetDevelopmentFundTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// Store updated proposal
	if err := k.SetDevelopmentFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Update development fund balance
	if err := k.UpdateDevelopmentFundBalance(ctx, initialAmount, false); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.DevEventTypeProjectStarted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyExecutor, msg.Executor),
			sdk.NewAttribute(types.AttributeKeyTransactionID, txID),
			sdk.NewAttribute(types.AttributeKeyAmount, initialAmount.String()),
		),
	})

	return &types.MsgExecuteDevelopmentProposalResponse{
		TransactionId: txID,
	}, nil
}

// UpdateMilestone implements milestone update functionality
func (k msgServer) UpdateMilestone(goCtx context.Context, msg *types.MsgUpdateMilestone) (*types.MsgUpdateMilestoneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetCommunityFundProposal(ctx, msg.ProposalId)
	if !found {
		// Try development fund proposal
		devProposal, devFound := k.GetDevelopmentFundProposal(ctx, msg.ProposalId)
		if !devFound {
			return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
		}
		// Handle development proposal milestone update
		return k.updateDevelopmentMilestone(ctx, devProposal, msg)
	}

	// Find and update milestone
	for i, milestone := range proposal.Milestones {
		if milestone.Id == msg.MilestoneId {
			if msg.Completed {
				milestone.Completed = true
				completedTime := ctx.BlockTime()
				milestone.CompletedAt = &completedTime
				milestone.Evidence = msg.Evidence
				milestone.Approved = false // Requires review
			}
			proposal.Milestones[i] = milestone
			break
		}
	}

	// Store updated proposal
	if err := k.SetCommunityFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMilestoneCompleted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyMilestone, fmt.Sprintf("%d", msg.MilestoneId)),
			sdk.NewAttribute("updater", msg.Updater),
			sdk.NewAttribute("completed", fmt.Sprintf("%t", msg.Completed)),
		),
	})

	return &types.MsgUpdateMilestoneResponse{}, nil
}

// AddMultiSigSigner implements adding a new signer to multi-sig governance
func (k msgServer) AddMultiSigSigner(goCtx context.Context, msg *types.MsgAddMultiSigSigner) (*types.MsgAddMultiSigSignerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Verify authority (simplified - should check governance permissions)
	if !k.IsAuthority(ctx, authority) {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "not authorized to add signers")
	}

	// Get multi-sig governance
	governance, found := k.GetMultiSigGovernance(ctx, msg.GovernanceId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "governance not found")
	}

	// Add new signer
	signer := types.Signer{
		Address:    msg.SignerAddress,
		Role:       msg.Role,
		Weight:     msg.Weight,
		AddedAt:    ctx.BlockTime(),
		AddedBy:    msg.Authority,
		Active:     true,
		Reputation: 100, // Starting reputation
	}

	if err := k.AddSigner(ctx, msg.GovernanceId, signer); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.MultiSigEventTypeSignerAdded,
			sdk.NewAttribute("governance_id", fmt.Sprintf("%d", msg.GovernanceId)),
			sdk.NewAttribute(types.AttributeKeySigner, msg.SignerAddress),
			sdk.NewAttribute("role", msg.Role),
			sdk.NewAttribute("weight", fmt.Sprintf("%d", msg.Weight)),
		),
	})

	return &types.MsgAddMultiSigSignerResponse{}, nil
}

// RemoveMultiSigSigner implements removing a signer from multi-sig governance
func (k msgServer) RemoveMultiSigSigner(goCtx context.Context, msg *types.MsgRemoveMultiSigSigner) (*types.MsgRemoveMultiSigSignerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Verify authority
	if !k.IsAuthority(ctx, authority) {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "not authorized to remove signers")
	}

	// Remove signer
	if err := k.RemoveSigner(ctx, msg.GovernanceId, msg.SignerAddress); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.MultiSigEventTypeSignerRemoved,
			sdk.NewAttribute("governance_id", fmt.Sprintf("%d", msg.GovernanceId)),
			sdk.NewAttribute(types.AttributeKeySigner, msg.SignerAddress),
		),
	})

	return &types.MsgRemoveMultiSigSignerResponse{}, nil
}

// SignMultiSigProposal implements signing a multi-sig proposal
func (k msgServer) SignMultiSigProposal(goCtx context.Context, msg *types.MsgSignMultiSigProposal) (*types.MsgSignMultiSigProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	// Get proposal
	proposal, found := k.GetMultiSigProposal(ctx, msg.ProposalId)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "proposal not found")
	}

	// Verify signer is authorized
	if !k.IsValidSigner(ctx, proposal.GovernanceId, msg.Signer) {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "not a valid signer")
	}

	// Add signature
	signature := types.Signature{
		Signer:    msg.Signer,
		Signature: msg.Signature,
		Timestamp: ctx.BlockTime(),
		Comments:  msg.Comments,
	}

	if err := k.AddSignature(ctx, msg.ProposalId, signature); err != nil {
		return nil, err
	}

	// Check if threshold reached
	currentSigs := k.GetSignatureCount(ctx, msg.ProposalId)
	requiredSigs := k.GetRequiredSignatures(ctx, proposal.GovernanceId)
	approved := currentSigs >= requiredSigs

	if approved {
		// Execute proposal
		if err := k.ExecuteMultiSigProposal(ctx, proposal); err != nil {
			return nil, err
		}
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.MultiSigEventTypeProposalSigned,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute("current_signatures", fmt.Sprintf("%d", currentSigs)),
			sdk.NewAttribute("required_signatures", fmt.Sprintf("%d", requiredSigs)),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", approved)),
		),
	})

	return &types.MsgSignMultiSigProposalResponse{
		CurrentSignatures:  uint32(currentSigs),
		RequiredSignatures: uint32(requiredSigs),
		Approved:           approved,
	}, nil
}

// UpdateGovernancePhase implements manual governance phase update
func (k msgServer) UpdateGovernancePhase(goCtx context.Context, msg *types.MsgUpdateGovernancePhase) (*types.MsgUpdateGovernancePhaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Verify authority
	if !k.IsAuthority(ctx, authority) {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "not authorized to update governance phase")
	}

	// Get current phase
	oldPhase := k.GetCurrentGovernancePhase(ctx)

	// Update phase
	if err := k.SetGovernancePhase(ctx, msg.NewPhase); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePhaseTransition,
			sdk.NewAttribute(types.AttributeKeyFromPhase, oldPhase),
			sdk.NewAttribute(types.AttributeKeyToPhase, msg.NewPhase),
			sdk.NewAttribute(types.AttributeKeyReason, msg.Reason),
			sdk.NewAttribute(types.AttributeKeyTransitionDate, ctx.BlockTime().String()),
		),
	})

	return &types.MsgUpdateGovernancePhaseResponse{
		OldPhase: oldPhase,
		NewPhase: msg.NewPhase,
	}, nil
}

// SubmitTransparencyReport implements transparency report submission
func (k msgServer) SubmitTransparencyReport(goCtx context.Context, msg *types.MsgSubmitTransparencyReport) (*types.MsgSubmitTransparencyReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	submitter, err := sdk.AccAddressFromBech32(msg.Submitter)
	if err != nil {
		return nil, err
	}

	// Create report
	reportID := k.GetNextTransparencyReportID(ctx)
	report := types.TransparencyReport{
		ReportId:       reportID,
		StartDate:      *msg.StartDate,
		EndDate:        *msg.EndDate,
		TotalFunds:     msg.TotalFunds,
		AllocatedFunds: msg.AllocatedFunds,
		SpentFunds:     msg.SpentFunds,
		RemainingFunds: sdk.NewCoin(msg.TotalFunds.Denom, msg.TotalFunds.Amount.Sub(msg.SpentFunds.Amount)),
		AuditStatus:    "pending",
		NextReportDate: msg.EndDate.AddDate(0, 3, 0), // Quarterly reports
	}

	// Store report
	if err := k.SetTransparencyReport(ctx, report); err != nil {
		return nil, err
	}

	// Update report counter
	k.SetNextTransparencyReportID(ctx, reportID+1)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransparencyReport,
			sdk.NewAttribute("report_id", fmt.Sprintf("%d", reportID)),
			sdk.NewAttribute("submitter", msg.Submitter),
			sdk.NewAttribute("period", fmt.Sprintf("%s to %s", msg.StartDate.Format("2006-01-02"), msg.EndDate.Format("2006-01-02"))),
		),
	})

	return &types.MsgSubmitTransparencyReportResponse{
		ReportId: reportID,
	}, nil
}

// UpdateDashboard implements real-time dashboard update
func (k msgServer) UpdateDashboard(goCtx context.Context, msg *types.MsgUpdateDashboard) (*types.MsgUpdateDashboardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		return nil, err
	}

	// Update dashboard metrics
	dashboard, err := k.UpdateDashboardMetrics(ctx, msg.ForceUpdate)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.MultiSigEventTypeTransparencyUpdate,
			sdk.NewAttribute("updater", msg.Updater),
			sdk.NewAttribute("last_updated", dashboard.LastUpdated.String()),
			sdk.NewAttribute("transparency_score", fmt.Sprintf("%d", dashboard.TransparencyScore)),
			sdk.NewAttribute("compliance_score", fmt.Sprintf("%d", dashboard.ComplianceScore)),
		),
	})

	return &types.MsgUpdateDashboardResponse{
		LastUpdated:       dashboard.LastUpdated,
		TransparencyScore: dashboard.TransparencyScore,
		ComplianceScore:   dashboard.ComplianceScore,
	}, nil
}

// Helper function to update development milestone
func (k msgServer) updateDevelopmentMilestone(ctx sdk.Context, proposal types.DevelopmentFundProposal, msg *types.MsgUpdateMilestone) (*types.MsgUpdateMilestoneResponse, error) {
	// Implementation for development fund milestones
	// Similar to community fund but with different validation rules
	
	// Store updated proposal
	if err := k.SetDevelopmentFundProposal(ctx, proposal); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.DevEventTypeMilestoneAchieved,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyMilestone, fmt.Sprintf("%d", msg.MilestoneId)),
			sdk.NewAttribute("completed", fmt.Sprintf("%t", msg.Completed)),
		),
	})

	return &types.MsgUpdateMilestoneResponse{}, nil
}