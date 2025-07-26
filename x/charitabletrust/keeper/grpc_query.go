package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DeshChain/DeshChain-Ecosystem/x/charitabletrust/types"
)

var _ types.QueryServer = Keeper{}

// TrustFundBalance returns the current trust fund balance
func (k Keeper) TrustFundBalance(goCtx context.Context, req *types.QueryTrustFundBalanceRequest) (*types.QueryTrustFundBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Update balance before returning
	k.UpdateTrustFundBalance(ctx)
	
	balance, found := k.GetTrustFundBalance(ctx)
	if !found {
		balance = types.TrustFundBalance{
			TotalBalance:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			AllocatedAmount:  sdk.NewCoin("unamo", sdk.ZeroInt()),
			AvailableAmount:  sdk.NewCoin("unamo", sdk.ZeroInt()),
			TotalDistributed: sdk.NewCoin("unamo", sdk.ZeroInt()),
			LastUpdated:      ctx.BlockTime(),
		}
	}

	return &types.QueryTrustFundBalanceResponse{
		Balance: balance,
	}, nil
}

// Allocation returns details of a specific allocation
func (k Keeper) Allocation(goCtx context.Context, req *types.QueryAllocationRequest) (*types.QueryAllocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allocation, found := k.GetCharitableAllocation(ctx, req.AllocationId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("allocation %d not found", req.AllocationId))
	}

	return &types.QueryAllocationResponse{
		Allocation: allocation,
	}, nil
}

// Allocations returns all allocations with pagination
func (k Keeper) Allocations(goCtx context.Context, req *types.QueryAllocationsRequest) (*types.QueryAllocationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var allocations []types.CharitableAllocation
	store := ctx.KVStore(k.storeKey)

	var allocStore prefix.Store
	if req.Status != "" {
		// Filter by status using index
		allocStore = prefix.NewStore(store, types.GetAllocationByStatusKey(req.Status, 0)[:len(types.GetAllocationByStatusKey(req.Status, 0))-8])
	} else {
		// All allocations
		allocStore = prefix.NewStore(store, types.CharitableAllocationKey)
	}

	pageRes, err := query.Paginate(allocStore, req.Pagination, func(key []byte, value []byte) error {
		if req.Status != "" {
			// Extract ID from index key
			allocationID := types.GetUint64FromBytes(key[len(req.Status):])
			// Get actual allocation
			allocation, found := k.GetCharitableAllocation(ctx, allocationID)
			if found {
				// Apply category filter if provided
				if req.Category == "" || allocation.Category == req.Category {
					allocations = append(allocations, allocation)
				}
			}
		} else {
			// Direct allocation storage
			var allocation types.CharitableAllocation
			k.cdc.MustUnmarshal(value, &allocation)
			// Apply category filter if provided
			if req.Category == "" || allocation.Category == req.Category {
				allocations = append(allocations, allocation)
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllocationsResponse{
		Allocations: allocations,
		Pagination:  pageRes,
	}, nil
}

// AllocationProposal returns a specific proposal
func (k Keeper) AllocationProposal(goCtx context.Context, req *types.QueryAllocationProposalRequest) (*types.QueryAllocationProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetAllocationProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("proposal %d not found", req.ProposalId))
	}

	return &types.QueryAllocationProposalResponse{
		Proposal: proposal,
	}, nil
}

// AllocationProposals returns all proposals
func (k Keeper) AllocationProposals(goCtx context.Context, req *types.QueryAllocationProposalsRequest) (*types.QueryAllocationProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var proposals []types.AllocationProposal
	store := ctx.KVStore(k.storeKey)

	var proposalStore prefix.Store
	if req.Status != "" {
		// Filter by status using index
		proposalStore = prefix.NewStore(store, types.GetProposalByStatusKey(req.Status, 0)[:len(types.GetProposalByStatusKey(req.Status, 0))-8])
	} else {
		// All proposals
		proposalStore = prefix.NewStore(store, types.AllocationProposalKey)
	}

	pageRes, err := query.Paginate(proposalStore, req.Pagination, func(key []byte, value []byte) error {
		if req.Status != "" {
			// Extract ID from index key
			proposalID := types.GetUint64FromBytes(key[len(req.Status):])
			// Get actual proposal
			proposal, found := k.GetAllocationProposal(ctx, proposalID)
			if found {
				proposals = append(proposals, proposal)
			}
		} else {
			// Direct proposal storage
			var proposal types.AllocationProposal
			k.cdc.MustUnmarshal(value, &proposal)
			proposals = append(proposals, proposal)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllocationProposalsResponse{
		Proposals:  proposals,
		Pagination: pageRes,
	}, nil
}

// ImpactReport returns a specific impact report
func (k Keeper) ImpactReport(goCtx context.Context, req *types.QueryImpactReportRequest) (*types.QueryImpactReportResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	report, found := k.GetImpactReport(ctx, req.ReportId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("report %d not found", req.ReportId))
	}

	return &types.QueryImpactReportResponse{
		Report: report,
	}, nil
}

// ImpactReports returns impact reports for an allocation
func (k Keeper) ImpactReports(goCtx context.Context, req *types.QueryImpactReportsRequest) (*types.QueryImpactReportsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var reports []types.ImpactReport
	store := ctx.KVStore(k.storeKey)
	
	// Use allocation index
	reportStore := prefix.NewStore(store, types.GetReportByAllocationKey(req.AllocationId, 0)[:len(types.GetReportByAllocationKey(req.AllocationId, 0))-8])

	pageRes, err := query.Paginate(reportStore, req.Pagination, func(key []byte, value []byte) error {
		// Extract report ID from index key
		reportID := types.GetUint64FromBytes(key)
		
		// Get actual report
		report, found := k.GetImpactReport(ctx, reportID)
		if found {
			reports = append(reports, report)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryImpactReportsResponse{
		Reports:    reports,
		Pagination: pageRes,
	}, nil
}

// FraudAlert returns a specific fraud alert
func (k Keeper) FraudAlert(goCtx context.Context, req *types.QueryFraudAlertRequest) (*types.QueryFraudAlertResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	alert, found := k.GetFraudAlert(ctx, req.AlertId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("alert %d not found", req.AlertId))
	}

	return &types.QueryFraudAlertResponse{
		Alert: alert,
	}, nil
}

// FraudAlerts returns all fraud alerts
func (k Keeper) FraudAlerts(goCtx context.Context, req *types.QueryFraudAlertsRequest) (*types.QueryFraudAlertsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var alerts []types.FraudAlert
	store := ctx.KVStore(k.storeKey)
	alertStore := prefix.NewStore(store, types.FraudAlertKey)

	pageRes, err := query.Paginate(alertStore, req.Pagination, func(key []byte, value []byte) error {
		var alert types.FraudAlert
		k.cdc.MustUnmarshal(value, &alert)
		
		// Apply filters
		if req.Status != "" && alert.Status != req.Status {
			return nil
		}
		if req.Severity != "" && alert.Severity != req.Severity {
			return nil
		}
		
		alerts = append(alerts, alert)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFraudAlertsResponse{
		Alerts:     alerts,
		Pagination: pageRes,
	}, nil
}

// TrustGovernance returns the governance configuration
func (k Keeper) TrustGovernance(goCtx context.Context, req *types.QueryTrustGovernanceRequest) (*types.QueryTrustGovernanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	governance, found := k.GetTrustGovernance(ctx)
	if !found {
		// Return default governance if not set
		governance = types.TrustGovernance{
			Trustees:            []types.Trustee{},
			Quorum:              4,
			ApprovalThreshold:   sdk.MustNewDecFromStr("0.571"),
			AdvisoryCommittee:   []types.AdvisoryMember{},
			TransparencyOfficer: "",
			NextElection:        ctx.BlockTime().Add(365 * 24 * 3600 * 1e9), // 1 year
		}
	}

	return &types.QueryTrustGovernanceResponse{
		Governance: governance,
	}, nil
}

// AllocationsByOrganization returns allocations for a specific organization
func (k Keeper) AllocationsByOrganization(goCtx context.Context, req *types.QueryAllocationsByOrganizationRequest) (*types.QueryAllocationsByOrganizationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var allocations []types.CharitableAllocation
	totalAllocated := sdk.NewCoin("unamo", sdk.ZeroInt())
	monthlyAllocated := k.CalculateOrganizationMonthlyAllocation(ctx, req.OrgWalletId)
	
	store := ctx.KVStore(k.storeKey)
	orgStore := prefix.NewStore(store, types.GetAllocationByOrgKey(req.OrgWalletId, 0)[:len(types.GetAllocationByOrgKey(req.OrgWalletId, 0))-8])

	pageRes, err := query.Paginate(orgStore, req.Pagination, func(key []byte, value []byte) error {
		// Extract allocation ID from index key
		allocationID := types.GetUint64FromBytes(key)
		
		// Get actual allocation
		allocation, found := k.GetCharitableAllocation(ctx, allocationID)
		if !found {
			return nil
		}

		// Apply status filter if provided
		if req.Status != "" && allocation.Status != req.Status {
			return nil
		}

		allocations = append(allocations, allocation)
		totalAllocated = totalAllocated.Add(allocation.Amount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllocationsByOrganizationResponse{
		Allocations:      allocations,
		TotalAllocated:   totalAllocated,
		MonthlyAllocated: monthlyAllocated,
		Pagination:       pageRes,
	}, nil
}

// Params queries all parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}