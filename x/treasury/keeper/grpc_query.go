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

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/treasury/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// CommunityFundBalance queries the current community fund balance
func (k Keeper) CommunityFundBalance(goCtx context.Context, req *types.QueryCommunityFundBalanceRequest) (*types.QueryCommunityFundBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	balance, err := k.GetCommunityFundBalance(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCommunityFundBalanceResponse{
		Balance: balance,
	}, nil
}

// DevelopmentFundBalance queries the current development fund balance
func (k Keeper) DevelopmentFundBalance(goCtx context.Context, req *types.QueryDevelopmentFundBalanceRequest) (*types.QueryDevelopmentFundBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	balance, err := k.GetDevelopmentFundBalance(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDevelopmentFundBalanceResponse{
		Balance: balance,
	}, nil
}

// CommunityProposal queries a specific community fund proposal
func (k Keeper) CommunityProposal(goCtx context.Context, req *types.QueryCommunityProposalRequest) (*types.QueryCommunityProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetCommunityFundProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}

	return &types.QueryCommunityProposalResponse{
		Proposal: proposal,
	}, nil
}

// CommunityProposals queries all community fund proposals with pagination
func (k Keeper) CommunityProposals(goCtx context.Context, req *types.QueryCommunityProposalsRequest) (*types.QueryCommunityProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposals := []types.CommunityFundProposal{}
	
	// Get all proposals
	iter, err := k.CommunityFundProposals.Iterate(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		proposal, err := iter.Value()
		if err != nil {
			continue
		}

		// Filter by status if provided
		if req.Status != "" && proposal.Status != req.Status {
			continue
		}

		// Filter by category if provided
		if req.Category != "" && proposal.Category != req.Category {
			continue
		}

		proposals = append(proposals, proposal)
	}

	// Apply pagination
	start, end := client.Paginate(len(proposals), req.Pagination)
	if start < 0 || end < 0 {
		return &types.QueryCommunityProposalsResponse{}, nil
	}
	
	paginatedProposals := proposals[start:end]

	return &types.QueryCommunityProposalsResponse{
		Proposals: paginatedProposals,
		Pagination: &query.PageResponse{
			Total: uint64(len(proposals)),
		},
	}, nil
}

// DevelopmentProposal queries a specific development fund proposal
func (k Keeper) DevelopmentProposal(goCtx context.Context, req *types.QueryDevelopmentProposalRequest) (*types.QueryDevelopmentProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetDevelopmentFundProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}

	return &types.QueryDevelopmentProposalResponse{
		Proposal: proposal,
	}, nil
}

// DevelopmentProposals queries all development fund proposals with pagination
func (k Keeper) DevelopmentProposals(goCtx context.Context, req *types.QueryDevelopmentProposalsRequest) (*types.QueryDevelopmentProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposals := []types.DevelopmentFundProposal{}
	
	// Get all proposals
	iter, err := k.DevelopmentFundProposals.Iterate(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		proposal, err := iter.Value()
		if err != nil {
			continue
		}

		// Filter by status if provided
		if req.Status != "" && proposal.Status != req.Status {
			continue
		}

		// Filter by category if provided
		if req.Category != "" && proposal.Category != req.Category {
			continue
		}

		// Filter by priority if provided
		if req.Priority != "" && proposal.Priority != req.Priority {
			continue
		}

		proposals = append(proposals, proposal)
	}

	// Apply pagination
	start, end := client.Paginate(len(proposals), req.Pagination)
	if start < 0 || end < 0 {
		return &types.QueryDevelopmentProposalsResponse{}, nil
	}
	
	paginatedProposals := proposals[start:end]

	return &types.QueryDevelopmentProposalsResponse{
		Proposals: paginatedProposals,
		Pagination: &query.PageResponse{
			Total: uint64(len(proposals)),
		},
	}, nil
}

// MultiSigGovernance queries multi-sig governance configuration
func (k Keeper) MultiSigGovernance(goCtx context.Context, req *types.QueryMultiSigGovernanceRequest) (*types.QueryMultiSigGovernanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	governance, found := k.GetMultiSigGovernance(ctx, req.GovernanceId)
	if !found {
		return nil, status.Error(codes.NotFound, "governance not found")
	}

	// Convert to proto message format
	protoGovernance := types.MultiSigGovernance{
		Id:                 governance.Id,
		Name:               governance.Name,
		Description:        governance.Description,
		Type:               string(governance.Type),
		Threshold:          uint32(governance.Threshold),
		TotalSigners:       uint32(len(governance.Signers)),
		ActiveProposals:    uint64(len(governance.Proposals)),
		CompletedProposals: governance.CompletedProposals,
		CreatedAt:          governance.CreatedAt,
		UpdatedAt:          governance.LastUpdated,
		Status:             governance.Status,
	}

	return &types.QueryMultiSigGovernanceResponse{
		Governance: protoGovernance,
	}, nil
}

// ProposalSystem queries the community proposal system configuration
func (k Keeper) ProposalSystem(goCtx context.Context, req *types.QueryProposalSystemRequest) (*types.QueryProposalSystemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	system, err := k.GetCommunityProposalSystem(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to proto message format
	protoSystem := types.CommunityProposalSystem{
		Id:             system.Id,
		Name:           system.Name,
		Description:    system.Description,
		LaunchDate:     system.LaunchDate,
		ActivationDate: system.ActivationDate,
		CurrentPhase:   string(k.GetCurrentGovernancePhase(ctx)),
		Status:         system.Status,
	}

	return &types.QueryProposalSystemResponse{
		System: protoSystem,
	}, nil
}

// Dashboard queries the real-time dashboard data
func (k Keeper) Dashboard(goCtx context.Context, req *types.QueryDashboardRequest) (*types.QueryDashboardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	dashboard, err := k.GetRealTimeDashboard(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get current balances
	communityBalance, _ := k.GetCommunityFundBalance(ctx)
	developmentBalance, _ := k.GetDevelopmentFundBalance(ctx)

	// Build community fund metrics
	communityMetrics := types.CommunityFundMetrics{
		TotalAllocation:    sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
		CurrentBalance:     communityBalance.TotalBalance,
		AllocatedAmount:    communityBalance.AllocatedAmount,
		SpentAmount:        sdk.NewCoin("namo", communityBalance.TotalBalance.Amount.Sub(communityBalance.AvailableAmount.Amount)),
		RemainingAmount:    communityBalance.AvailableAmount,
		ActiveProposals:    k.GetActiveProposalCount(ctx, "community"),
		CompletedProposals: k.GetCompletedProposalCount(ctx, "community"),
		LastUpdated:        ctx.BlockTime(),
	}

	// Build development fund metrics
	developmentMetrics := types.DevelopmentFundMetrics{
		TotalAllocation:    sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
		CurrentBalance:     developmentBalance.TotalBalance,
		AllocatedAmount:    developmentBalance.AllocatedAmount,
		SpentAmount:        sdk.NewCoin("namo", developmentBalance.TotalBalance.Amount.Sub(developmentBalance.AvailableAmount.Amount)),
		RemainingAmount:    developmentBalance.AvailableAmount,
		ActiveProjects:     k.GetActiveProposalCount(ctx, "development"),
		CompletedProjects:  k.GetCompletedProposalCount(ctx, "development"),
		LastUpdated:        ctx.BlockTime(),
	}

	// Build overall metrics
	overallMetrics := types.OverallFundMetrics{
		TotalFunds:          sdk.NewCoin("namo", sdk.NewInt(types.TotalTreasuryAllocation)),
		TotalAllocated:      sdk.NewCoin("namo", communityBalance.AllocatedAmount.Amount.Add(developmentBalance.AllocatedAmount.Amount)),
		TotalSpent:          sdk.NewCoin("namo", communityMetrics.SpentAmount.Amount.Add(developmentMetrics.SpentAmount.Amount)),
		TotalRemaining:      sdk.NewCoin("namo", communityBalance.AvailableAmount.Amount.Add(developmentBalance.AvailableAmount.Amount)),
		OverallTransparency: dashboard.TransparencyScore,
		OverallCompliance:   dashboard.ComplianceScore,
		LastUpdated:         ctx.BlockTime(),
	}

	// Build proto dashboard
	protoDashboard := types.RealTimeDashboard{
		LastUpdated:       dashboard.LastUpdated,
		CommunityFund:     communityMetrics,
		DevelopmentFund:   developmentMetrics,
		OverallMetrics:    overallMetrics,
		TransparencyScore: dashboard.TransparencyScore,
		ComplianceScore:   dashboard.ComplianceScore,
		Status:            dashboard.Status,
	}

	return &types.QueryDashboardResponse{
		Dashboard: protoDashboard,
	}, nil
}

// CommunityTransactions queries community fund transactions with pagination
func (k Keeper) CommunityTransactions(goCtx context.Context, req *types.QueryCommunityTransactionsRequest) (*types.QueryCommunityTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	transactions := []types.CommunityFundTransaction{}
	
	// Get all transactions
	iter, err := k.CommunityFundTransactions.Iterate(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		tx, err := iter.Value()
		if err != nil {
			continue
		}

		// Filter by status if provided
		if req.Status != "" && tx.Status != req.Status {
			continue
		}

		// Filter by type if provided
		if req.Type != "" && tx.Type != req.Type {
			continue
		}

		// Filter by category if provided
		if req.Category != "" && tx.Category != req.Category {
			continue
		}

		transactions = append(transactions, tx)
	}

	// Apply pagination
	start, end := client.Paginate(len(transactions), req.Pagination)
	if start < 0 || end < 0 {
		return &types.QueryCommunityTransactionsResponse{}, nil
	}
	
	paginatedTransactions := transactions[start:end]

	return &types.QueryCommunityTransactionsResponse{
		Transactions: paginatedTransactions,
		Pagination: &query.PageResponse{
			Total: uint64(len(transactions)),
		},
	}, nil
}

// DevelopmentTransactions queries development fund transactions with pagination
func (k Keeper) DevelopmentTransactions(goCtx context.Context, req *types.QueryDevelopmentTransactionsRequest) (*types.QueryDevelopmentTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	transactions := []types.DevelopmentFundTransaction{}
	
	// Get all transactions
	iter, err := k.DevelopmentFundTransactions.Iterate(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		tx, err := iter.Value()
		if err != nil {
			continue
		}

		// Filter by status if provided
		if req.Status != "" && tx.Status != req.Status {
			continue
		}

		// Filter by type if provided
		if req.Type != "" && tx.Type != req.Type {
			continue
		}

		// Filter by category if provided
		if req.Category != "" && tx.Category != req.Category {
			continue
		}

		transactions = append(transactions, tx)
	}

	// Apply pagination
	start, end := client.Paginate(len(transactions), req.Pagination)
	if start < 0 || end < 0 {
		return &types.QueryDevelopmentTransactionsResponse{}, nil
	}
	
	paginatedTransactions := transactions[start:end]

	return &types.QueryDevelopmentTransactionsResponse{
		Transactions: paginatedTransactions,
		Pagination: &query.PageResponse{
			Total: uint64(len(transactions)),
		},
	}, nil
}

// GovernancePhase queries the current governance phase
func (k Keeper) GovernancePhase(goCtx context.Context, req *types.QueryGovernancePhaseRequest) (*types.QueryGovernancePhaseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	currentPhase := k.GetCurrentGovernancePhase(ctx)
	phaseConfig := k.GetPhaseConfiguration(ctx, currentPhase)
	
	// Calculate phase dates
	phaseStartDate := k.GetPhaseStartDate(ctx, currentPhase)
	nextPhaseDate := k.GetNextPhaseTransitionDate(ctx)
	nextPhase := k.GetNextPhase(currentPhase)

	// Build phase details
	phaseDetails := types.PhaseDetails{
		Description:            phaseConfig.Description,
		FounderAllocationPower: uint32(phaseConfig.FounderAllocationPower),
		CommunityProposalPower: uint32(phaseConfig.CommunityProposalPower),
		CommunityVotingEnabled: phaseConfig.CommunityVotingEnabled,
		FounderVetoEnabled:     phaseConfig.FounderVetoEnabled,
		AllowedProposalTypes:   phaseConfig.AllowedProposalTypes,
	}

	return &types.QueryGovernancePhaseResponse{
		CurrentPhase:   currentPhase,
		PhaseStartDate: phaseStartDate,
		NextPhaseDate:  nextPhaseDate,
		NextPhase:      nextPhase,
		PhaseDetails:   phaseDetails,
	}, nil
}

// TransparencyReport queries transparency reports
func (k Keeper) TransparencyReport(goCtx context.Context, req *types.QueryTransparencyReportRequest) (*types.QueryTransparencyReportResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	report, found := k.GetTransparencyReport(ctx, req.ReportId)
	if !found {
		return nil, status.Error(codes.NotFound, "report not found")
	}

	return &types.QueryTransparencyReportResponse{
		Report: report,
	}, nil
}

// Metrics queries treasury metrics
func (k Keeper) Metrics(goCtx context.Context, req *types.QueryMetricsRequest) (*types.QueryMetricsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Calculate metrics based on request parameters
	metrics := k.CalculateTreasuryMetrics(ctx, req.MetricType, req.Period)

	// Build proto metrics
	protoMetrics := types.TreasuryMetrics{
		Period:              req.Period,
		TotalInflow:         metrics.TotalInflow,
		TotalOutflow:        metrics.TotalOutflow,
		ProposalsSubmitted:  metrics.ProposalsSubmitted,
		ProposalsApproved:   metrics.ProposalsApproved,
		ProposalsRejected:   metrics.ProposalsRejected,
		ProposalsExecuted:   metrics.ProposalsExecuted,
		AllocationEfficiency: metrics.AllocationEfficiency,
		SpendingEfficiency:  metrics.SpendingEfficiency,
		TransparencyScore:   metrics.TransparencyScore,
		ComplianceScore:     metrics.ComplianceScore,
		CalculatedAt:        ctx.BlockTime(),
	}

	return &types.QueryMetricsResponse{
		Metrics: protoMetrics,
	}, nil
}

// Helper function for pagination
var client = struct {
	Paginate func(int, *query.PageRequest) (int, int)
}{
	Paginate: func(total int, pageReq *query.PageRequest) (int, int) {
		if pageReq == nil {
			return 0, total
		}
		
		offset := int(pageReq.Offset)
		limit := int(pageReq.Limit)
		
		if offset >= total {
			return -1, -1
		}
		
		if limit == 0 {
			limit = 100 // Default limit
		}
		
		end := offset + limit
		if end > total {
			end = total
		}
		
		return offset, end
	},
}