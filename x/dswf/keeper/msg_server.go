package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/deshchain/deshchain/x/dswf/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// ProposeAllocation handles fund allocation proposals
func (k msgServer) ProposeAllocation(goCtx context.Context, msg *types.MsgProposeAllocation) (*types.MsgProposeAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Operational safety check
	if err := k.ValidateOperationalSafety(ctx, "allocation", msg.Amount); err != nil {
		return nil, err
	}
	
	// Validate the proposal
	if err := k.ValidateAllocationProposal(ctx, msg.Amount, msg.Category); err != nil {
		return nil, types.ErrInvalidProposal.Wrapf("allocation proposal validation failed: %v", err)
	}
	
	// Verify multi-signature requirement for fund managers
	if !k.ValidateMultiSignature(ctx, msg.Proposers) {
		return nil, types.ErrInsufficientSignatures.Wrap("allocation requires multiple fund manager signatures")
	}
	
	// Create new allocation
	allocationID := k.IncrementFundAllocationCount(ctx)
	allocation := types.FundAllocation{
		Id:               allocationID,
		Purpose:          msg.Purpose,
		Category:         msg.Category,
		Amount:           msg.Amount,
		Recipient:        msg.Recipient,
		ApprovedBy:       msg.Proposers, // All proposers who signed
		ProposalId:       0, // Will be set when governance integration is complete
		AllocatedAt:      ctx.BlockTime(),
		ExpectedOutcomes: msg.ExpectedOutcomes,
		Disbursements:    msg.DisbursementSchedule,
		Status:           "pending",
		Metrics:          []types.PerformanceMetric{},
		Roi:              sdk.ZeroDec(),
	}
	
	// Validate disbursement schedule
	totalDisbursement := sdk.NewCoin(msg.Amount.Denom, sdk.ZeroInt())
	for i, disbursement := range allocation.Disbursements {
		totalDisbursement = totalDisbursement.Add(disbursement.Amount)
		allocation.Disbursements[i].Status = "pending"
	}
	
	if !totalDisbursement.IsEqual(msg.Amount) {
		return nil, types.ErrInvalidAmount.Wrap("disbursement schedule total doesn't match allocation amount")
	}
	
	// Save allocation
	k.SetFundAllocation(ctx, allocation)
	
	// Record revenue activity for DSWF allocation proposal
	k.RecordDSWFRevenueActivity(ctx, "allocation_proposed", msg.Amount)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAllocationProposed,
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", allocationID)),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCategory, msg.Category),
		),
	)
	
	return &types.MsgProposeAllocationResponse{
		AllocationId: allocationID,
		ProposalId:   0, // Governance proposal ID will be implemented in governance integration
	}, nil
}

// ApproveAllocation handles allocation approvals
func (k msgServer) ApproveAllocation(goCtx context.Context, msg *types.MsgApproveAllocation) (*types.MsgApproveAllocationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify approver is a fund manager
	if !k.IsFundManager(ctx, msg.Approver) {
		return nil, types.ErrUnauthorized.Wrapf("approver %s is not a fund manager", msg.Approver)
	}
	
	// Additional validation: Check if allocation still needs approval
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return nil, types.ErrGovernanceNotFound
	}
	
	// Get allocation
	allocation, found := k.GetFundAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	// Check if already approved by this manager
	for _, approver := range allocation.ApprovedBy {
		if approver == msg.Approver {
			return nil, types.ErrDuplicateApproval
		}
	}
	
	// Add approval
	if msg.Approved {
		allocation.ApprovedBy = append(allocation.ApprovedBy, msg.Approver)
		
		// Check if we have enough approvals
		governance, _ := k.GetFundGovernance(ctx)
		requiredApprovals := int(float64(len(governance.FundManagers)) * governance.Threshold.MustFloat64())
		
		if len(allocation.ApprovedBy) >= requiredApprovals {
			allocation.Status = "approved"
			
			// Reserve funds for allocation
			moduleAddr := k.GetModuleAccountAddress()
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ModuleName+"_reserved", sdk.NewCoins(allocation.Amount)); err != nil {
				return nil, sdkerrors.Wrap(err, "failed to reserve funds")
			}
		}
	}
	
	k.SetFundAllocation(ctx, allocation)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAllocationApproved,
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", msg.AllocationId)),
			sdk.NewAttribute(types.AttributeKeyApprover, msg.Approver),
			sdk.NewAttribute(types.AttributeKeyApproved, fmt.Sprintf("%t", msg.Approved)),
			sdk.NewAttribute(types.AttributeKeyStatus, allocation.Status),
		),
	)
	
	return &types.MsgApproveAllocationResponse{
		Success: true,
		Status:  allocation.Status,
	}, nil
}

// ExecuteDisbursement handles disbursement execution
func (k msgServer) ExecuteDisbursement(goCtx context.Context, msg *types.MsgExecuteDisbursement) (*types.MsgExecuteDisbursementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get allocation to check disbursement amount for safety validation
	allocation, found := k.GetFundAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	if int(msg.DisbursementIndex) >= len(allocation.Disbursements) {
		return nil, types.ErrInvalidProposal.Wrap("invalid disbursement index")
	}
	
	disbursementAmount := allocation.Disbursements[msg.DisbursementIndex].Amount
	
	// Operational safety check
	if err := k.ValidateOperationalSafety(ctx, "disbursement", disbursementAmount); err != nil {
		return nil, err
	}
	
	// Get allocation (again for consistency)
	allocation, found := k.GetFundAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	// Check allocation is approved
	if allocation.Status != "approved" && allocation.Status != "active" {
		return nil, types.ErrInvalidStatus.Wrapf("allocation status is %s", allocation.Status)
	}
	
	// Get disbursement
	if int(msg.DisbursementIndex) >= len(allocation.Disbursements) {
		return nil, types.ErrInvalidProposal.Wrap("invalid disbursement index")
	}
	
	disbursement := &allocation.Disbursements[msg.DisbursementIndex]
	
	// Check if already disbursed
	if disbursement.Status == "disbursed" {
		return nil, types.ErrInvalidStatus.Wrap("disbursement already executed")
	}
	
	// Check if scheduled date has passed
	if disbursement.ScheduledDate.After(ctx.BlockTime()) {
		return nil, types.ErrDisbursementNotReady
	}
	
	// Verify milestone proof if required
	if disbursement.MilestoneProof != "" {
		if err := k.ValidateMilestoneProof(ctx, disbursement.MilestoneProof); err != nil {
			return nil, types.ErrInvalidProof.Wrap("milestone proof validation failed")
		}
	}
	
	// Execute disbursement
	recipientAddr, err := sdk.AccAddressFromBech32(allocation.Recipient)
	if err != nil {
		return nil, err
	}
	
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName+"_reserved", recipientAddr, sdk.NewCoins(disbursement.Amount)); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to execute disbursement")
	}
	
	// Record revenue activity for DSWF disbursement
	k.RecordDSWFRevenueActivity(ctx, "funds_disbursed", disbursement.Amount)
	
	// Update disbursement status
	disbursement.Status = "disbursed"
	disbursement.DisbursedAt = ctx.BlockTime()
	disbursement.TxHash = fmt.Sprintf("%X", ctx.TxBytes()) // Simplified, use actual tx hash
	
	// Update allocation status
	allocation.Status = "active"
	allDisbursed := true
	for _, d := range allocation.Disbursements {
		if d.Status != "disbursed" {
			allDisbursed = false
			break
		}
	}
	if allDisbursed {
		allocation.Status = "completed"
	}
	
	k.SetFundAllocation(ctx, allocation)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDisbursementExecuted,
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", msg.AllocationId)),
			sdk.NewAttribute(types.AttributeKeyDisbursementIndex, fmt.Sprintf("%d", msg.DisbursementIndex)),
			sdk.NewAttribute(types.AttributeKeyAmount, disbursement.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, allocation.Recipient),
		),
	)
	
	return &types.MsgExecuteDisbursementResponse{
		Success:         true,
		TxHash:          disbursement.TxHash,
		AmountDisbursed: disbursement.Amount,
	}, nil
}

// UpdateInvestmentStrategy handles investment strategy updates
func (k msgServer) UpdateInvestmentStrategy(goCtx context.Context, msg *types.MsgUpdateInvestmentStrategy) (*types.MsgUpdateInvestmentStrategyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify authority
	if msg.Authority != k.GetAuthority() {
		return nil, types.ErrUnauthorized.Wrapf("invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}
	
	// Validate new strategy
	if err := msg.NewStrategy.Validate(); err != nil {
		return nil, types.ErrInvalidInvestmentStrategy.Wrap(err.Error())
	}
	
	// Update params
	params := k.GetParams(ctx)
	params.InvestmentStrategy = *msg.NewStrategy
	k.SetParams(ctx, params)
	
	// TODO: Trigger portfolio rebalancing
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStrategyUpdated,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute("conservative", msg.NewStrategy.ConservativePercentage.String()),
			sdk.NewAttribute("moderate", msg.NewStrategy.ModeratePercentage.String()),
			sdk.NewAttribute("aggressive", msg.NewStrategy.AggressivePercentage.String()),
		),
	)
	
	return &types.MsgUpdateInvestmentStrategyResponse{Success: true}, nil
}

// RebalancePortfolio handles portfolio rebalancing
func (k msgServer) RebalancePortfolio(goCtx context.Context, msg *types.MsgRebalancePortfolio) (*types.MsgRebalancePortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify authority has permission to rebalance
	if err := k.ValidateRebalanceAuthority(ctx, msg.Authority); err != nil {
		return nil, types.ErrUnauthorized.Wrapf("rebalance authority validation failed: %v", err)
	}
	
	// Get current portfolio
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if !found {
		// Initialize new portfolio
		portfolio = types.InvestmentPortfolio{
			TotalValue:       k.GetFundBalance(ctx)[0], // Assuming single denom
			LiquidAssets:     k.GetFundBalance(ctx)[0],
			InvestedAssets:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			ReservedAssets:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			Components:       []types.PortfolioComponent{},
			TotalReturns:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			AnnualReturnRate: sdk.ZeroDec(),
			RiskScore:        5,
			LastRebalanced:   ctx.BlockTime(),
		}
	}
	
	// Implement portfolio rebalancing logic
	if err := k.ExecutePortfolioRebalancing(ctx, &portfolio); err != nil {
		return nil, types.ErrInvalidInvestmentStrategy.Wrapf("portfolio rebalancing execution failed: %v", err)
	}
	
	portfolio.LastRebalanced = ctx.BlockTime()
	k.SetInvestmentPortfolio(ctx, portfolio)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePortfolioRebalanced,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute("total_value", portfolio.TotalValue.String()),
			sdk.NewAttribute("risk_score", fmt.Sprintf("%d", portfolio.RiskScore)),
		),
	)
	
	return &types.MsgRebalancePortfolioResponse{
		Success:      true,
		NewPortfolio: &portfolio,
	}, nil
}

// SubmitPerformanceMetrics handles performance metric submissions
func (k msgServer) SubmitPerformanceMetrics(goCtx context.Context, msg *types.MsgSubmitPerformanceMetrics) (*types.MsgSubmitPerformanceMetricsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get allocation
	allocation, found := k.GetFundAllocation(ctx, msg.AllocationId)
	if !found {
		return nil, types.ErrAllocationNotFound
	}
	
	// Verify submitter is the recipient or an authorized auditor
	params := k.GetParams(ctx)
	isAuthorized := allocation.Recipient == msg.Submitter
	for _, auditor := range params.AuthorizedAuditors {
		if auditor == msg.Submitter {
			isAuthorized = true
			break
		}
	}
	
	if !isAuthorized {
		return nil, types.ErrUnauthorized
	}
	
	// Update metrics
	allocation.Metrics = append(allocation.Metrics, msg.Metrics...)
	
	// Calculate achievement percentage
	totalAchievement := sdk.ZeroDec()
	for _, metric := range allocation.Metrics {
		if metric.Status == "on_track" {
			totalAchievement = totalAchievement.Add(sdk.OneDec())
		} else if metric.Status == "at_risk" {
			totalAchievement = totalAchievement.Add(sdk.NewDecWithPrec(5, 1)) // 0.5
		}
	}
	
	if len(allocation.Metrics) > 0 {
		avgAchievement := totalAchievement.Quo(sdk.NewDec(int64(len(allocation.Metrics))))
		// Simple ROI calculation based on achievement
		allocation.Roi = avgAchievement.Mul(sdk.NewDecWithPrec(2, 1)) // 20% max ROI
	}
	
	k.SetFundAllocation(ctx, allocation)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMetricsSubmitted,
			sdk.NewAttribute(types.AttributeKeyAllocationID, fmt.Sprintf("%d", msg.AllocationId)),
			sdk.NewAttribute(types.AttributeKeySubmitter, msg.Submitter),
			sdk.NewAttribute("metrics_count", fmt.Sprintf("%d", len(msg.Metrics))),
			sdk.NewAttribute("roi", allocation.Roi.String()),
		),
	)
	
	return &types.MsgSubmitPerformanceMetricsResponse{
		Success:    true,
		Evaluation: fmt.Sprintf("ROI: %s", allocation.Roi.String()),
	}, nil
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