package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/deshchain/deshchain/x/charitabletrust/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateAllocationProposal creates a new allocation proposal
func (k msgServer) CreateAllocationProposal(goCtx context.Context, msg *types.MsgCreateAllocationProposal) (*types.MsgCreateAllocationProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify proposer is a trustee
	if !k.IsTrustee(ctx, msg.Proposer) {
		return nil, types.ErrNotTrustee.Wrapf("proposer %s is not a trustee", msg.Proposer)
	}
	
	// Validate total amount matches allocations
	totalAllocated := sdk.NewCoin(msg.TotalAmount.Denom, sdk.ZeroInt())
	for _, alloc := range msg.Allocations {
		// Validate organization
		if err := k.ValidateCharitableOrganization(ctx, alloc.CharitableOrgWalletId); err != nil {
			return nil, err
		}
		
		totalAllocated = totalAllocated.Add(alloc.Amount)
	}
	
	if !totalAllocated.IsEqual(msg.TotalAmount) {
		return nil, types.ErrInvalidAmount.Wrap("sum of allocations doesn't match total amount")
	}
	
	// Check fund balance
	balance, found := k.GetTrustFundBalance(ctx)
	if !found {
		return nil, types.ErrInsufficientFunds.Wrap("trust fund not initialized")
	}
	
	if balance.AvailableAmount.IsLT(msg.TotalAmount) {
		return nil, types.ErrInsufficientFunds
	}
	
	// Create proposal
	proposalID := k.IncrementProposalCount(ctx)
	proposal := types.AllocationProposal{
		Id:             proposalID,
		Proposer:       msg.Proposer,
		Title:          msg.Title,
		Description:    msg.Description,
		TotalAmount:    msg.TotalAmount,
		Allocations:    msg.Allocations,
		Justification:  msg.Justification,
		ExpectedImpact: msg.ExpectedImpact,
		Documents:      msg.Documents,
		VotingStart:    ctx.BlockTime(),
		VotingEnd:      ctx.BlockTime().Add(time.Duration(k.GetParams(ctx).ProposalVotingPeriod) * time.Second),
		Votes:          []types.Vote{},
		Status:         "pending",
	}
	
	k.SetAllocationProposal(ctx, proposal)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalCreated,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.TotalAmount.String()),
		),
	)
	
	return &types.MsgCreateAllocationProposalResponse{
		ProposalId: proposalID,
	}, nil
}

// VoteOnProposal handles trustee votes
func (k msgServer) VoteOnProposal(goCtx context.Context, msg *types.MsgVoteOnProposal) (*types.MsgVoteOnProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify voter is a trustee
	if !k.IsTrustee(ctx, msg.Voter) {
		return nil, types.ErrNotTrustee.Wrapf("voter %s is not a trustee", msg.Voter)
	}
	
	// Get proposal
	proposal, found := k.GetAllocationProposal(ctx, msg.ProposalId)
	if !found {
		return nil, types.ErrInvalidProposal.Wrap("proposal not found")
	}
	
	// Check voting period
	if ctx.BlockTime().Before(proposal.VotingStart) {
		return nil, types.ErrVotingPeriodActive.Wrap("voting hasn't started")
	}
	
	if ctx.BlockTime().After(proposal.VotingEnd) {
		return nil, types.ErrVotingPeriodExpired
	}
	
	// Check for duplicate vote
	for _, vote := range proposal.Votes {
		if vote.Voter == msg.Voter {
			return nil, types.ErrDuplicateVote
		}
	}
	
	// Add vote
	proposal.Votes = append(proposal.Votes, types.Vote{
		Voter:   msg.Voter,
		Vote:    msg.Vote,
		Reason:  msg.Reason,
		VotedAt: ctx.BlockTime(),
	})
	
	// Check if we have enough votes
	governance, _ := k.GetTrustGovernance(ctx)
	yesVotes := 0
	noVotes := 0
	
	for _, vote := range proposal.Votes {
		if vote.Vote == "yes" {
			yesVotes++
		} else if vote.Vote == "no" {
			noVotes++
		}
	}
	
	totalTrustees := len(governance.Trustees)
	requiredVotes := int(float64(totalTrustees) * governance.ApprovalThreshold.MustFloat64())
	
	if yesVotes >= requiredVotes {
		proposal.Status = "approved"
		
		// Update fund balance
		balance, _ := k.GetTrustFundBalance(ctx)
		balance.AllocatedAmount = balance.AllocatedAmount.Add(proposal.TotalAmount)
		balance.AvailableAmount = balance.AvailableAmount.Sub(proposal.TotalAmount)
		k.SetTrustFundBalance(ctx, balance)
		
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalApproved,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	} else if noVotes > (totalTrustees - requiredVotes) {
		proposal.Status = "rejected"
		
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalRejected,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	}
	
	k.SetAllocationProposal(ctx, proposal)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVoted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Voter),
			sdk.NewAttribute(types.AttributeKeyVote, msg.Vote),
		),
	)
	
	return &types.MsgVoteOnProposalResponse{
		Success: true,
		Status:  proposal.Status,
	}, nil
}

// ExecuteAllocation executes an approved proposal
func (k msgServer) ExecuteAllocation(goCtx context.Context, msg *types.MsgExecuteAllocation) (*types.MsgExecuteAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get proposal
	proposal, found := k.GetAllocationProposal(ctx, msg.ProposalId)
	if !found {
		return nil, types.ErrInvalidProposal.Wrap("proposal not found")
	}
	
	// Check proposal is approved
	if proposal.Status != "approved" {
		return nil, types.ErrInvalidProposal.Wrapf("proposal status is %s", proposal.Status)
	}
	
	// Check voting period has ended
	if ctx.BlockTime().Before(proposal.VotingEnd) {
		return nil, types.ErrVotingPeriodActive
	}
	
	// Execute allocations
	allocationIDs := []uint64{}
	
	for _, proposedAlloc := range proposal.Allocations {
		// Create allocation
		allocationID := k.IncrementAllocationCount(ctx)
		allocation := types.CharitableAllocation{
			Id:                    allocationID,
			CharitableOrgWalletId: proposedAlloc.CharitableOrgWalletId,
			OrganizationName:      proposedAlloc.OrganizationName,
			Amount:                proposedAlloc.Amount,
			Purpose:               proposedAlloc.Purpose,
			Category:              proposedAlloc.Category,
			ProposalId:            proposal.Id,
			ApprovedBy:            []string{}, // Collect from votes
			AllocatedAt:           ctx.BlockTime(),
			ExpectedImpact:        proposal.ExpectedImpact,
			Monitoring: types.MonitoringRequirements{
				ReportingFrequency:    k.GetParams(ctx).ImpactReportFrequency,
				RequiredReports:       []string{"impact", "financial", "beneficiary"},
				Kpis:                  []string{"beneficiaries_reached", "funds_utilized", "outcomes_achieved"},
				MonitoringDuration:    180, // 6 months
				SiteVisitsRequired:    true,
				FinancialAuditRequired: proposedAlloc.Amount.Amount.GT(sdk.NewInt(1000000000)), // >1000 NAMO
			},
			Status: "distributed",
		}
		
		// Collect approvers from votes
		for _, vote := range proposal.Votes {
			if vote.Vote == "yes" {
				allocation.ApprovedBy = append(allocation.ApprovedBy, vote.Voter)
			}
		}
		
		// TODO: Get actual organization wallet address from donation module
		// For now, using a placeholder
		recipientAddr := sdk.AccAddress([]byte(fmt.Sprintf("org%d", proposedAlloc.CharitableOrgWalletId)))
		
		// Transfer funds
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, sdk.NewCoins(proposedAlloc.Amount)); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to transfer funds")
		}
		
		// Set distribution details
		allocation.Distribution = &types.DistributionDetails{
			TxHash:        fmt.Sprintf("%X", ctx.TxBytes()),
			DistributedAt: ctx.BlockTime(),
			DistributedBy: msg.Executor,
		}
		
		k.SetCharitableAllocation(ctx, allocation)
		allocationIDs = append(allocationIDs, allocationID)
		
		// Update fund balance
		balance, _ := k.GetTrustFundBalance(ctx)
		balance.AllocatedAmount = balance.AllocatedAmount.Sub(proposedAlloc.Amount)
		balance.TotalDistributed = balance.TotalDistributed.Add(proposedAlloc.Amount)
		k.SetTrustFundBalance(ctx, balance)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeFundsDistributed,
				sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", allocationID)),
				sdk.NewAttribute(types.AttributeKeyOrganizationID, fmt.Sprintf("%d", proposedAlloc.CharitableOrgWalletId)),
				sdk.NewAttribute(types.AttributeKeyAmount, proposedAlloc.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyCategory, proposedAlloc.Category),
			),
		)
	}
	
	// Update proposal status
	proposal.Status = "executed"
	k.SetAllocationProposal(ctx, proposal)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAllocationExecuted,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("allocation_count", fmt.Sprintf("%d", len(allocationIDs))),
		),
	)
	
	return &types.MsgExecuteAllocationResponse{
		Success:       true,
		AllocationIds: allocationIDs,
	}, nil
}

// SubmitImpactReport handles impact report submissions
func (k msgServer) SubmitImpactReport(goCtx context.Context, msg *types.MsgSubmitImpactReport) (*types.MsgSubmitImpactReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get allocation
	allocation, found := k.GetCharitableAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	// TODO: Verify submitter is authorized (org representative)
	
	// Create impact report
	reportID := k.GetAllocationCount(ctx) + 1000000 // Simple ID generation
	report := types.ImpactReport{
		Id:                  reportID,
		AllocationId:        msg.AllocationId,
		Period:              msg.Period,
		BeneficiariesReached: msg.BeneficiariesReached,
		FundsUtilized:       msg.FundsUtilized,
		Metrics:             msg.Metrics,
		Documents:           msg.Documents,
		Media:               msg.Media,
		Challenges:          msg.Challenges,
		SubmittedBy:         msg.Submitter,
		SubmittedAt:         ctx.BlockTime(),
		Verification: &types.VerificationStatus{
			IsVerified: false,
		},
	}
	
	k.SetImpactReport(ctx, report)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeImpactReportSubmitted,
			sdk.NewAttribute(types.AttributeKeyReportID, fmt.Sprintf("%d", reportID)),
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", msg.AllocationId)),
			sdk.NewAttribute(types.AttributeKeyBeneficiaries, fmt.Sprintf("%d", msg.BeneficiariesReached)),
		),
	)
	
	return &types.MsgSubmitImpactReportResponse{
		ReportId: reportID,
	}, nil
}

// VerifyImpactReport handles report verification
func (k msgServer) VerifyImpactReport(goCtx context.Context, msg *types.MsgVerifyImpactReport) (*types.MsgVerifyImpactReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if verifier is authorized
	params := k.GetParams(ctx)
	isAuthorized := false
	for _, investigator := range params.AuthorizedInvestigators {
		if investigator == msg.Verifier {
			isAuthorized = true
			break
		}
	}
	
	if !isAuthorized && !k.IsTrustee(ctx, msg.Verifier) {
		return nil, types.ErrUnauthorized
	}
	
	// Get report
	report, found := k.GetImpactReport(ctx, msg.ReportId)
	if !found {
		return nil, types.ErrReportNotFound
	}
	
	// Update verification status
	report.Verification = &types.VerificationStatus{
		IsVerified:             msg.Verified,
		VerifiedBy:             msg.Verifier,
		VerifiedAt:             ctx.BlockTime(),
		Notes:                  msg.Notes,
		SiteVisitConducted:     msg.SiteVisitConducted,
		FinancialAuditConducted: msg.FinancialAuditConducted,
	}
	
	k.SetImpactReport(ctx, report)
	
	// Calculate impact score
	impactScore := 0
	if msg.Verified {
		for _, metric := range report.Metrics {
			if metric.AchievementPercentage.GTE(sdk.NewDecWithPrec(8, 1)) { // 80%+
				impactScore += 20
			} else if metric.AchievementPercentage.GTE(sdk.NewDecWithPrec(5, 1)) { // 50%+
				impactScore += 10
			}
		}
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeImpactReportVerified,
			sdk.NewAttribute(types.AttributeKeyReportID, fmt.Sprintf("%d", msg.ReportId)),
			sdk.NewAttribute(types.AttributeKeyVerifier, msg.Verifier),
			sdk.NewAttribute("verified", fmt.Sprintf("%t", msg.Verified)),
			sdk.NewAttribute(types.AttributeKeyImpactScore, fmt.Sprintf("%d", impactScore)),
		),
	)
	
	return &types.MsgVerifyImpactReportResponse{Success: true}, nil
}

// ReportFraud handles fraud reporting
func (k msgServer) ReportFraud(goCtx context.Context, msg *types.MsgReportFraud) (*types.MsgReportFraudResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get allocation
	_, found := k.GetCharitableAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	// Create fraud alert
	alertID := k.GetAllocationCount(ctx) + 2000000 // Simple ID generation
	alert := types.FraudAlert{
		Id:           alertID,
		AllocationId: msg.AllocationId,
		AlertType:    msg.AlertType,
		Severity:     msg.Severity,
		Description:  msg.Description,
		Evidence:     msg.Evidence,
		ReportedBy:   msg.Reporter,
		ReportedAt:   ctx.BlockTime(),
		Status:       "reported",
	}
	
	k.SetFraudAlert(ctx, alert)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFraudAlertCreated,
			sdk.NewAttribute(types.AttributeKeyAlertID, fmt.Sprintf("%d", alertID)),
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", msg.AllocationId)),
			sdk.NewAttribute(types.AttributeKeySeverity, msg.Severity),
		),
	)
	
	return &types.MsgReportFraudResponse{
		AlertId: alertID,
	}, nil
}

// InvestigateFraud handles fraud investigation
func (k msgServer) InvestigateFraud(goCtx context.Context, msg *types.MsgInvestigateFraud) (*types.MsgInvestigateFraudResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if investigator is authorized
	params := k.GetParams(ctx)
	isAuthorized := false
	for _, investigator := range params.AuthorizedInvestigators {
		if investigator == msg.Investigator {
			isAuthorized = true
			break
		}
	}
	
	if !isAuthorized {
		return nil, types.ErrUnauthorized
	}
	
	// Get alert
	alert, found := k.GetFraudAlert(ctx, msg.AlertId)
	if !found {
		return nil, types.ErrAlertNotFound
	}
	
	// Update investigation
	if alert.Investigation == nil {
		alert.Investigation = &types.Investigation{
			Investigator: msg.Investigator,
			StartedAt:    ctx.BlockTime(),
		}
	}
	
	alert.Investigation.Findings = msg.Findings
	alert.Investigation.Recommendation = msg.Recommendation
	alert.Investigation.Report = msg.Report
	
	if msg.InvestigationComplete {
		alert.Investigation.CompletedAt = ctx.BlockTime()
		alert.Status = "investigated"
	} else {
		alert.Status = "investigating"
	}
	
	k.SetFraudAlert(ctx, alert)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFraudInvestigated,
			sdk.NewAttribute(types.AttributeKeyAlertID, fmt.Sprintf("%d", msg.AlertId)),
			sdk.NewAttribute(types.AttributeKeyInvestigator, msg.Investigator),
			sdk.NewAttribute(types.AttributeKeyStatus, alert.Status),
		),
	)
	
	return &types.MsgInvestigateFraudResponse{
		Success: true,
		Status:  alert.Status,
	}, nil
}

// UpdateTrustees updates the board of trustees
func (k msgServer) UpdateTrustees(goCtx context.Context, msg *types.MsgUpdateTrustees) (*types.MsgUpdateTrusteesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify authority
	if msg.Authority != k.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}
	
	// Get current governance
	governance, found := k.GetTrustGovernance(ctx)
	if !found {
		governance = types.TrustGovernance{}
	}
	
	// Update trustees
	governance.Trustees = msg.NewTrustees
	
	// Validate governance
	if err := governance.Validate(); err != nil {
		return nil, types.ErrInvalidProposal.Wrap(err.Error())
	}
	
	k.SetTrustGovernance(ctx, governance)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrusteesUpdated,
			sdk.NewAttribute("trustee_count", fmt.Sprintf("%d", len(msg.NewTrustees))),
		),
	)
	
	return &types.MsgUpdateTrusteesResponse{Success: true}, nil
}

// UpdateParams updates the module parameters
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	if msg.Authority != k.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}
	
	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}
	
	k.SetParams(ctx, msg.Params)
	
	return &types.MsgUpdateParamsResponse{}, nil
}