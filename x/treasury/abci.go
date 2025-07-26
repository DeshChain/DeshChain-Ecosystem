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

package treasury

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/treasury/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/treasury/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker checks for automatic phase transitions and updates dashboard metrics
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Check for governance phase transition
	currentPhase := k.GetCurrentGovernancePhase(ctx)
	system, err := k.GetCommunityProposalSystem(ctx)
	if err == nil {
		nextPhase := checkPhaseTransition(ctx, system, currentPhase)
		if nextPhase != currentPhase {
			// Transition to next phase
			k.SetGovernancePhase(ctx, nextPhase)
			
			// Emit phase transition event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypePhaseTransition,
					sdk.NewAttribute(types.AttributeKeyFromPhase, currentPhase),
					sdk.NewAttribute(types.AttributeKeyToPhase, nextPhase),
					sdk.NewAttribute(types.AttributeKeyTransitionDate, ctx.BlockTime().String()),
				),
			)
		}
	}

	// Update dashboard metrics if needed (every 5 minutes)
	dashboard, err := k.GetRealTimeDashboard(ctx)
	if err == nil {
		if ctx.BlockTime().Sub(dashboard.LastUpdated) >= types.DashboardUpdateInterval {
			k.UpdateDashboardMetrics(ctx, false)
		}
	}
}

// EndBlocker processes completed proposals and expired proposals
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Process community fund proposals
	processCommunityProposals(ctx, k)

	// Process development fund proposals
	processDevelopmentProposals(ctx, k)

	// Process multi-sig proposals
	processMultiSigProposals(ctx, k)

	// Generate quarterly transparency report if needed
	generateTransparencyReport(ctx, k)

	return []abci.ValidatorUpdate{}
}

// processCommunityProposals checks for proposals that have ended voting period
func processCommunityProposals(ctx sdk.Context, k keeper.Keeper) {
	// Iterate through all active proposals
	iter, err := k.CommunityFundProposals.Iterate(ctx, nil)
	if err != nil {
		return
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		proposal, err := iter.Value()
		if err != nil {
			continue
		}

		// Skip if not in voting period
		if proposal.Status != types.StatusActive && proposal.Status != types.StatusPending {
			continue
		}

		// Check if voting period has ended
		if ctx.BlockTime().After(proposal.VotingEndTime) {
			// Check if proposal passed
			if k.HasProposalPassed(ctx, proposal) {
				proposal.Status = types.StatusPassed
				proposal.Passed = true
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeProposalPassed,
						sdk.NewAttribute(types.AttributeKeyProposalID, sdk.NewInt(int64(proposal.ProposalId)).String()),
						sdk.NewAttribute(types.AttributeKeyCategory, proposal.Category),
					),
				)
			} else {
				proposal.Status = types.StatusRejected
				proposal.Passed = false
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeProposalRejected,
						sdk.NewAttribute(types.AttributeKeyProposalID, sdk.NewInt(int64(proposal.ProposalId)).String()),
						sdk.NewAttribute(types.AttributeKeyCategory, proposal.Category),
					),
				)
			}

			// Update proposal
			k.SetCommunityFundProposal(ctx, proposal)
		}
	}
}

// processDevelopmentProposals checks for proposals that have ended review period
func processDevelopmentProposals(ctx sdk.Context, k keeper.Keeper) {
	// Iterate through all pending proposals
	iter, err := k.DevelopmentFundProposals.Iterate(ctx, nil)
	if err != nil {
		return
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		proposal, err := iter.Value()
		if err != nil {
			continue
		}

		// Skip if not pending
		if proposal.Status != types.StatusPending {
			continue
		}

		// Check if review period has ended
		if ctx.BlockTime().After(proposal.ReviewEndTime) {
			// Auto-reject if no reviews received
			if !k.AreAllReviewsComplete(proposal) {
				proposal.Status = types.StatusRejected
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.DevEventTypeProposalRejected,
						sdk.NewAttribute(types.AttributeKeyProposalID, sdk.NewInt(int64(proposal.ProposalId)).String()),
						sdk.NewAttribute(types.AttributeKeyReason, "review_timeout"),
					),
				)
				
				// Update proposal
				k.SetDevelopmentFundProposal(ctx, proposal)
			}
		}

		// Check milestone deadlines for executing proposals
		if proposal.Status == types.StatusExecuting {
			checkMilestoneDeadlines(ctx, k, proposal)
		}
	}
}

// processMultiSigProposals checks for expired multi-sig proposals
func processMultiSigProposals(ctx sdk.Context, k keeper.Keeper) {
	// Iterate through all multi-sig proposals
	iter, err := k.MultiSigProposals.Iterate(ctx, nil)
	if err != nil {
		return
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		proposal, err := iter.Value()
		if err != nil {
			continue
		}

		// Skip if already executed or cancelled
		if proposal.Status != types.StatusPending && proposal.Status != types.StatusActive {
			continue
		}

		// Check if proposal has expired
		if ctx.BlockTime().After(proposal.ExpiryTime) {
			proposal.Status = types.StatusExpired
			
			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.MultiSigEventTypeProposalCancelled,
					sdk.NewAttribute(types.AttributeKeyProposalID, sdk.NewInt(int64(proposal.ProposalId)).String()),
					sdk.NewAttribute(types.AttributeKeyReason, "expired"),
				),
			)
			
			// Update proposal
			k.SetMultiSigProposal(ctx, proposal)
		}
	}
}

// generateTransparencyReport generates quarterly transparency reports
func generateTransparencyReport(ctx sdk.Context, k keeper.Keeper) {
	// Get last report
	lastReportID := k.GetNextTransparencyReportID(ctx) - 1
	if lastReportID > 0 {
		lastReport, found := k.GetTransparencyReport(ctx, lastReportID)
		if found {
			// Check if it's time for next quarterly report
			if ctx.BlockTime().After(lastReport.NextReportDate) {
				generateNewTransparencyReport(ctx, k, lastReport.EndDate)
			}
		}
	} else {
		// Generate first report if none exists and system has been running for 3 months
		system, err := k.GetCommunityProposalSystem(ctx)
		if err == nil && ctx.BlockTime().Sub(system.LaunchDate) >= 90*24*time.Hour {
			generateNewTransparencyReport(ctx, k, system.LaunchDate)
		}
	}
}

// generateNewTransparencyReport creates a new transparency report
func generateNewTransparencyReport(ctx sdk.Context, k keeper.Keeper, startDate time.Time) {
	endDate := ctx.BlockTime()
	
	// Get fund balances
	communityBalance, _ := k.GetCommunityFundBalance(ctx)
	developmentBalance, _ := k.GetDevelopmentFundBalance(ctx)
	
	// Calculate totals
	totalFunds := sdk.NewCoin("namo", 
		communityBalance.TotalBalance.Amount.Add(developmentBalance.TotalBalance.Amount))
	allocatedFunds := sdk.NewCoin("namo",
		communityBalance.AllocatedAmount.Amount.Add(developmentBalance.AllocatedAmount.Amount))
	spentFunds := sdk.NewCoin("namo",
		communityBalance.TotalBalance.Amount.Sub(communityBalance.AvailableAmount.Amount).
			Add(developmentBalance.TotalBalance.Amount.Sub(developmentBalance.AvailableAmount.Amount)))
	remainingFunds := sdk.NewCoin("namo",
		communityBalance.AvailableAmount.Amount.Add(developmentBalance.AvailableAmount.Amount))
	
	// Create report
	reportID := k.GetNextTransparencyReportID(ctx)
	report := types.TransparencyReport{
		ReportId:       reportID,
		StartDate:      startDate,
		EndDate:        endDate,
		TotalFunds:     totalFunds,
		AllocatedFunds: allocatedFunds,
		SpentFunds:     spentFunds,
		RemainingFunds: remainingFunds,
		AuditStatus:    "pending",
		NextReportDate: endDate.AddDate(0, 3, 0), // Next quarterly report
	}
	
	// Store report
	k.SetTransparencyReport(ctx, report)
	k.SetNextTransparencyReportID(ctx, reportID+1)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransparencyReport,
			sdk.NewAttribute("report_id", sdk.NewInt(int64(reportID)).String()),
			sdk.NewAttribute("period", "quarterly"),
			sdk.NewAttribute("total_funds", totalFunds.String()),
		),
	)
}

// checkPhaseTransition determines if governance phase should transition
func checkPhaseTransition(ctx sdk.Context, system types.CommunityProposalSystem, currentPhase string) string {
	now := ctx.BlockTime()
	
	switch currentPhase {
	case types.PhaseFounderDriven:
		if now.After(system.PhaseSchedule.FounderDrivenEnd) {
			return types.PhaseTransitional
		}
	case types.PhaseTransitional:
		if now.After(system.PhaseSchedule.TransitionalEnd) {
			return types.PhaseCommunityProposal
		}
	case types.PhaseCommunityProposal:
		if now.After(system.PhaseSchedule.CommunityProposalEnd) {
			return types.PhaseFullGovernance
		}
	}
	
	return currentPhase
}

// checkMilestoneDeadlines checks if any milestones are overdue
func checkMilestoneDeadlines(ctx sdk.Context, k keeper.Keeper, proposal types.DevelopmentFundProposal) {
	updated := false
	
	for i, milestone := range proposal.Timeline.Milestones {
		if !milestone.Completed && ctx.BlockTime().After(milestone.DueDate) {
			// Mark milestone as overdue
			// This would trigger alerts and potentially affect future funding
			updated = true
			
			// Emit event for overdue milestone
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.DevEventTypeDeliverableSubmitted,
					sdk.NewAttribute(types.AttributeKeyProposalID, sdk.NewInt(int64(proposal.ProposalId)).String()),
					sdk.NewAttribute(types.AttributeKeyMilestone, sdk.NewInt(int64(milestone.Id)).String()),
					sdk.NewAttribute(types.AttributeKeyStatus, "overdue"),
				),
			)
		}
	}
	
	if updated {
		k.SetDevelopmentFundProposal(ctx, proposal)
	}
}