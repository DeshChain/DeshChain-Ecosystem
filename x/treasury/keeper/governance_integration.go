package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/treasury/types"
)

// GovernanceIntegration handles multi-signature governance for treasury operations
type GovernanceIntegration struct {
	keeper Keeper
}

// NewGovernanceIntegration creates a new governance integration
func NewGovernanceIntegration(keeper Keeper) *GovernanceIntegration {
	return &GovernanceIntegration{
		keeper: keeper,
	}
}

// GovernanceProposal represents a treasury governance proposal
type GovernanceProposal struct {
	ProposalID        string                    `json:"proposal_id"`
	ProposalType      string                    `json:"proposal_type"` // WITHDRAWAL, REBALANCE, PARAMETER_CHANGE, POOL_CREATION
	Title             string                    `json:"title"`
	Description       string                    `json:"description"`
	Proposer          string                    `json:"proposer"`
	ProposalDetails   types.ProposalDetails     `json:"proposal_details"`
	RequiredSignatures int                      `json:"required_signatures"`
	Signatures        []types.GovernanceSignature `json:"signatures"`
	VotingPeriod      types.VotingPeriod        `json:"voting_period"`
	ExecutionWindow   types.ExecutionWindow     `json:"execution_window"`
	Status            string                    `json:"status"` // PENDING, VOTING, APPROVED, REJECTED, EXECUTED, EXPIRED
	CreatedAt         time.Time                 `json:"created_at"`
	VotingStartTime   *time.Time                `json:"voting_start_time,omitempty"`
	VotingEndTime     *time.Time                `json:"voting_end_time,omitempty"`
	ExecutedAt        *time.Time                `json:"executed_at,omitempty"`
	RiskAssessment    types.ProposalRiskAssessment `json:"risk_assessment"`
}

// SubmitTreasuryProposal submits a new treasury governance proposal
func (gi *GovernanceIntegration) SubmitTreasuryProposal(ctx sdk.Context, request types.ProposalRequest) (*GovernanceProposal, error) {
	// Validate proposer permissions
	if err := gi.validateProposerPermissions(ctx, request.Proposer, request.ProposalType); err != nil {
		return nil, fmt.Errorf("proposer validation failed: %w", err)
	}

	// Validate proposal details
	if err := gi.validateProposalDetails(ctx, request); err != nil {
		return nil, fmt.Errorf("proposal validation failed: %w", err)
	}

	// Generate proposal ID
	proposalID := gi.generateProposalID(ctx, request.ProposalType)

	// Determine required signatures based on proposal type and amount
	requiredSignatures := gi.calculateRequiredSignatures(ctx, request)

	// Create governance proposal
	proposal := &GovernanceProposal{
		ProposalID:         proposalID,
		ProposalType:       request.ProposalType,
		Title:              request.Title,
		Description:        request.Description,
		Proposer:           request.Proposer,
		ProposalDetails:    request.Details,
		RequiredSignatures: requiredSignatures,
		Signatures:         []types.GovernanceSignature{},
		VotingPeriod:       gi.calculateVotingPeriod(ctx, request),
		ExecutionWindow:    gi.calculateExecutionWindow(ctx, request),
		Status:             "PENDING",
		CreatedAt:          ctx.BlockTime(),
		RiskAssessment:     gi.assessProposalRisk(ctx, request),
	}

	// If proposal requires immediate voting, start voting period
	if gi.requiresVoting(request.ProposalType) {
		gi.startVotingPeriod(ctx, proposal)
	}

	// Store proposal
	gi.keeper.SetGovernanceProposal(ctx, *proposal)

	// Emit proposal submission event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryProposalSubmitted,
			sdk.NewAttribute(types.AttributeKeyProposalID, proposalID),
			sdk.NewAttribute(types.AttributeKeyProposalType, request.ProposalType),
			sdk.NewAttribute(types.AttributeKeyProposer, request.Proposer),
			sdk.NewAttribute(types.AttributeKeyRequiredSignatures, fmt.Sprintf("%d", requiredSignatures)),
		),
	)

	return proposal, nil
}

// SignTreasuryProposal adds a signature to a treasury proposal
func (gi *GovernanceIntegration) SignTreasuryProposal(ctx sdk.Context, proposalID string, signer string, signatureData types.SignatureData) error {
	// Get proposal
	proposal, found := gi.keeper.GetGovernanceProposal(ctx, proposalID)
	if !found {
		return fmt.Errorf("proposal not found: %s", proposalID)
	}

	// Validate proposal status
	if proposal.Status != "PENDING" && proposal.Status != "VOTING" {
		return fmt.Errorf("proposal not in signable state: %s", proposal.Status)
	}

	// Validate signer permissions
	if err := gi.validateSignerPermissions(ctx, signer, proposal); err != nil {
		return fmt.Errorf("signer validation failed: %w", err)
	}

	// Check if signer already signed
	for _, sig := range proposal.Signatures {
		if sig.Signer == signer {
			return fmt.Errorf("signer already signed this proposal: %s", signer)
		}
	}

	// Verify signature
	if err := gi.verifySignature(ctx, proposal, signer, signatureData); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	// Add signature
	signature := types.GovernanceSignature{
		Signer:      signer,
		SignedAt:    ctx.BlockTime(),
		Signature:   signatureData.Signature,
		PublicKey:   signatureData.PublicKey,
		SignatureType: signatureData.SignatureType,
	}
	
	proposal.Signatures = append(proposal.Signatures, signature)

	// Check if proposal has enough signatures
	if len(proposal.Signatures) >= proposal.RequiredSignatures {
		proposal.Status = "APPROVED"
		
		// Start execution window for approved proposals
		gi.startExecutionWindow(ctx, &proposal)
	}

	// Update proposal
	gi.keeper.SetGovernanceProposal(ctx, proposal)

	// Emit signature event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryProposalSigned,
			sdk.NewAttribute(types.AttributeKeyProposalID, proposalID),
			sdk.NewAttribute(types.AttributeKeySigner, signer),
			sdk.NewAttribute(types.AttributeKeySignatureCount, fmt.Sprintf("%d", len(proposal.Signatures))),
			sdk.NewAttribute(types.AttributeKeyRequiredSignatures, fmt.Sprintf("%d", proposal.RequiredSignatures)),
		),
	)

	return nil
}

// ExecuteTreasuryProposal executes an approved treasury proposal
func (gi *GovernanceIntegration) ExecuteTreasuryProposal(ctx sdk.Context, proposalID string, executor string) error {
	// Get proposal
	proposal, found := gi.keeper.GetGovernanceProposal(ctx, proposalID)
	if !found {
		return fmt.Errorf("proposal not found: %s", proposalID)
	}

	// Validate proposal status
	if proposal.Status != "APPROVED" {
		return fmt.Errorf("proposal not approved for execution: %s", proposal.Status)
	}

	// Check execution window
	if proposal.ExecutionWindow.EndTime != nil && ctx.BlockTime().After(*proposal.ExecutionWindow.EndTime) {
		proposal.Status = "EXPIRED"
		gi.keeper.SetGovernanceProposal(ctx, proposal)
		return fmt.Errorf("proposal execution window expired")
	}

	// Validate executor permissions
	if err := gi.validateExecutorPermissions(ctx, executor, proposal); err != nil {
		return fmt.Errorf("executor validation failed: %w", err)
	}

	// Execute proposal based on type
	err := gi.executeProposalByType(ctx, proposal, executor)
	if err != nil {
		return fmt.Errorf("proposal execution failed: %w", err)
	}

	// Update proposal status
	now := ctx.BlockTime()
	proposal.Status = "EXECUTED"
	proposal.ExecutedAt = &now
	gi.keeper.SetGovernanceProposal(ctx, proposal)

	// Emit execution event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryProposalExecuted,
			sdk.NewAttribute(types.AttributeKeyProposalID, proposalID),
			sdk.NewAttribute(types.AttributeKeyExecutor, executor),
			sdk.NewAttribute(types.AttributeKeyExecutedAt, now.String()),
		),
	)

	return nil
}

// executeProposalByType executes proposal based on its type
func (gi *GovernanceIntegration) executeProposalByType(ctx sdk.Context, proposal GovernanceProposal, executor string) error {
	switch proposal.ProposalType {
	case "WITHDRAWAL":
		return gi.executeWithdrawalProposal(ctx, proposal, executor)
	case "REBALANCE":
		return gi.executeRebalanceProposal(ctx, proposal, executor)
	case "PARAMETER_CHANGE":
		return gi.executeParameterChangeProposal(ctx, proposal, executor)
	case "POOL_CREATION":
		return gi.executePoolCreationProposal(ctx, proposal, executor)
	case "POOL_MODIFICATION":
		return gi.executePoolModificationProposal(ctx, proposal, executor)
	case "EMERGENCY_ACTION":
		return gi.executeEmergencyActionProposal(ctx, proposal, executor)
	default:
		return fmt.Errorf("unknown proposal type: %s", proposal.ProposalType)
	}
}

// executeWithdrawalProposal executes a withdrawal proposal
func (gi *GovernanceIntegration) executeWithdrawalProposal(ctx sdk.Context, proposal GovernanceProposal, executor string) error {
	withdrawalDetails := proposal.ProposalDetails.WithdrawalDetails
	if withdrawalDetails == nil {
		return fmt.Errorf("withdrawal details not found in proposal")
	}

	// Create withdrawal request
	request := types.WithdrawalRequest{
		PoolID:     withdrawalDetails.PoolID,
		Amount:     withdrawalDetails.Amount,
		Recipient:  withdrawalDetails.Recipient,
		Purpose:    withdrawalDetails.Purpose,
		Requester:  executor,
		Timestamp:  ctx.BlockTime(),
		Approved:   true, // Already approved through governance
		ProposalID: proposal.ProposalID,
	}

	// Get treasury manager and execute withdrawal
	treasuryManager := NewTreasuryManager(gi.keeper)
	return treasuryManager.ProcessPoolWithdrawal(ctx, request)
}

// executeRebalanceProposal executes a rebalance proposal
func (gi *GovernanceIntegration) executeRebalanceProposal(ctx sdk.Context, proposal GovernanceProposal, executor string) error {
	rebalanceDetails := proposal.ProposalDetails.RebalanceDetails
	if rebalanceDetails == nil {
		return fmt.Errorf("rebalance details not found in proposal")
	}

	// Get rebalance engine and execute rebalancing
	rebalanceEngine := NewRebalanceEngine(gi.keeper)
	
	// Create rebalance plan from proposal details
	plan := &RebalancePlan{
		RebalanceID:   fmt.Sprintf("GOV_%s", proposal.ProposalID),
		TriggerReason: fmt.Sprintf("Governance proposal: %s", proposal.Title),
		Actions:       rebalanceDetails.Actions,
		CreatedAt:     ctx.BlockTime(),
		Status:        "APPROVED",
	}

	return rebalanceEngine.ExecuteRebalancePlan(ctx, plan)
}

// validateProposerPermissions validates if proposer can submit this type of proposal
func (gi *GovernanceIntegration) validateProposerPermissions(ctx sdk.Context, proposer string, proposalType string) error {
	// Get proposer's roles
	roles := gi.keeper.GetUserRoles(ctx, proposer)
	
	// Check if proposer has required role for proposal type
	requiredRoles := gi.getRequiredRolesForProposal(proposalType)
	
	for _, requiredRole := range requiredRoles {
		for _, userRole := range roles {
			if userRole == requiredRole {
				return nil // Found matching role
			}
		}
	}
	
	return fmt.Errorf("proposer %s lacks required roles %v for proposal type %s", proposer, requiredRoles, proposalType)
}

// calculateRequiredSignatures calculates required signatures based on proposal type and amount
func (gi *GovernanceIntegration) calculateRequiredSignatures(ctx sdk.Context, request types.ProposalRequest) int {
	baseSignatures := 2 // Minimum signatures
	
	switch request.ProposalType {
	case "WITHDRAWAL":
		if request.Details.WithdrawalDetails != nil {
			amount := request.Details.WithdrawalDetails.Amount.AmountOf("namo")
			
			// Higher amounts require more signatures
			if amount.GT(sdk.NewInt(1000000000000)) { // > 1M NAMO
				return 5
			} else if amount.GT(sdk.NewInt(100000000000)) { // > 100K NAMO
				return 4
			} else if amount.GT(sdk.NewInt(10000000000)) { // > 10K NAMO
				return 3
			}
		}
		return baseSignatures
		
	case "REBALANCE":
		return 3 // Rebalancing requires 3 signatures
		
	case "PARAMETER_CHANGE":
		return 4 // Parameter changes require 4 signatures
		
	case "POOL_CREATION", "POOL_MODIFICATION":
		return 4 // Pool operations require 4 signatures
		
	case "EMERGENCY_ACTION":
		return 2 // Emergency actions can be executed faster
		
	default:
		return baseSignatures
	}
}

// assessProposalRisk assesses the risk level of a proposal
func (gi *GovernanceIntegration) assessProposalRisk(ctx sdk.Context, request types.ProposalRequest) types.ProposalRiskAssessment {
	assessment := types.ProposalRiskAssessment{
		RiskLevel: "LOW",
		RiskFactors: []string{},
	}
	
	switch request.ProposalType {
	case "WITHDRAWAL":
		if request.Details.WithdrawalDetails != nil {
			amount := request.Details.WithdrawalDetails.Amount.AmountOf("namo")
			
			// Check withdrawal amount risk
			if amount.GT(sdk.NewInt(1000000000000)) { // > 1M NAMO
				assessment.RiskLevel = "HIGH"
				assessment.RiskFactors = append(assessment.RiskFactors, "Large withdrawal amount")
			} else if amount.GT(sdk.NewInt(100000000000)) { // > 100K NAMO
				assessment.RiskLevel = "MEDIUM"
				assessment.RiskFactors = append(assessment.RiskFactors, "Moderate withdrawal amount")
			}
			
			// Check pool impact
			pool, found := gi.keeper.GetTreasuryPool(ctx, request.Details.WithdrawalDetails.PoolID)
			if found {
				withdrawalPercentage := amount.ToDec().Quo(pool.Balance.AmountOf("namo").ToDec())
				if withdrawalPercentage.GT(sdk.NewDecWithPrec(50, 2)) { // > 50% of pool
					assessment.RiskLevel = "HIGH"
					assessment.RiskFactors = append(assessment.RiskFactors, "High percentage of pool withdrawal")
				}
			}
		}
		
	case "REBALANCE":
		assessment.RiskLevel = "MEDIUM"
		assessment.RiskFactors = append(assessment.RiskFactors, "Treasury rebalancing operation")
		
	case "PARAMETER_CHANGE":
		assessment.RiskLevel = "HIGH"
		assessment.RiskFactors = append(assessment.RiskFactors, "System parameter modification")
		
	case "EMERGENCY_ACTION":
		assessment.RiskLevel = "HIGH"
		assessment.RiskFactors = append(assessment.RiskFactors, "Emergency action with potential system impact")
	}
	
	return assessment
}

// getRequiredRolesForProposal returns required roles for each proposal type
func (gi *GovernanceIntegration) getRequiredRolesForProposal(proposalType string) []string {
	roleMap := map[string][]string{
		"WITHDRAWAL":        {"TREASURY_MANAGER", "POOL_MANAGER", "BOARD_MEMBER"},
		"REBALANCE":         {"TREASURY_MANAGER", "FINANCIAL_OFFICER"},
		"PARAMETER_CHANGE":  {"TECHNICAL_LEAD", "TREASURY_MANAGER", "BOARD_MEMBER"},
		"POOL_CREATION":     {"TREASURY_MANAGER", "BOARD_MEMBER"},
		"POOL_MODIFICATION": {"TREASURY_MANAGER", "POOL_MANAGER"},
		"EMERGENCY_ACTION":  {"EMERGENCY_RESPONDER", "TREASURY_MANAGER", "FOUNDER"},
	}
	
	if roles, found := roleMap[proposalType]; found {
		return roles
	}
	
	return []string{"TREASURY_MANAGER"} // Default required role
}

// Helper utility functions
func (gi *GovernanceIntegration) generateProposalID(ctx sdk.Context, proposalType string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("%s_PROP_%d_%d", proposalType, ctx.BlockHeight(), timestamp)
}

func (gi *GovernanceIntegration) requiresVoting(proposalType string) bool {
	votingTypes := map[string]bool{
		"PARAMETER_CHANGE": true,
		"POOL_CREATION":    true,
		"EMERGENCY_ACTION": false, // Can be executed with signatures only
		"WITHDRAWAL":       false, // Can be executed with signatures only
		"REBALANCE":        false, // Can be executed with signatures only
	}
	
	return votingTypes[proposalType]
}

func (gi *GovernanceIntegration) startVotingPeriod(ctx sdk.Context, proposal *GovernanceProposal) {
	now := ctx.BlockTime()
	proposal.Status = "VOTING"
	proposal.VotingStartTime = &now
	
	endTime := now.Add(proposal.VotingPeriod.Duration)
	proposal.VotingEndTime = &endTime
}

func (gi *GovernanceIntegration) startExecutionWindow(ctx sdk.Context, proposal *GovernanceProposal) {
	now := ctx.BlockTime()
	proposal.ExecutionWindow.StartTime = now
	
	endTime := now.Add(proposal.ExecutionWindow.Duration)
	proposal.ExecutionWindow.EndTime = &endTime
}

// Additional helper methods would include:
// - validateProposalDetails
// - calculateVotingPeriod
// - calculateExecutionWindow
// - validateSignerPermissions
// - verifySignature
// - validateExecutorPermissions
// - executeParameterChangeProposal
// - executePoolCreationProposal
// - executePoolModificationProposal
// - executeEmergencyActionProposal