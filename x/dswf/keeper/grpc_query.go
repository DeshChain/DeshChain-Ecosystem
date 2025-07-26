package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deshchain/deshchain/x/dswf/types"
)

var _ types.QueryServer = Keeper{}

// FundStatus returns the current status of the DSWF
func (k Keeper) FundStatus(goCtx context.Context, req *types.QueryFundStatusRequest) (*types.QueryFundStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get fund balance
	fundBalance := k.GetFundBalance(ctx)
	totalBalance := sdk.NewCoin("unamo", sdk.ZeroInt())
	if len(fundBalance) > 0 {
		totalBalance = fundBalance[0]
	}

	// Get portfolio
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if !found {
		portfolio = types.InvestmentPortfolio{
			TotalValue:       totalBalance,
			AllocatedAmount:  sdk.NewCoin("unamo", sdk.ZeroInt()),
			AvailableAmount:  totalBalance,
			InvestedAmount:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			TotalReturns:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			AnnualReturnRate: sdk.ZeroDec(),
		}
	}

	// Count allocations
	store := ctx.KVStore(k.storeKey)
	activeCount := 0
	completedCount := 0

	activeIterator := sdk.KVStorePrefixIterator(store, types.GetAllocationByStatusKey("active", 0)[:len(types.GetAllocationByStatusKey("active", 0))-8])
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		activeCount++
	}

	completedIterator := sdk.KVStorePrefixIterator(store, types.GetAllocationByStatusKey("completed", 0)[:len(types.GetAllocationByStatusKey("completed", 0))-8])
	defer completedIterator.Close()
	for ; completedIterator.Valid(); completedIterator.Next() {
		completedCount++
	}

	// Calculate available amount
	availableAmount := totalBalance.Sub(portfolio.AllocatedAmount)
	if availableAmount.IsNegative() {
		availableAmount = sdk.NewCoin("unamo", sdk.ZeroInt())
	}

	return &types.QueryFundStatusResponse{
		TotalBalance:        totalBalance,
		AllocatedAmount:     portfolio.AllocatedAmount,
		AvailableAmount:     availableAmount,
		InvestedAmount:      portfolio.InvestedAmount,
		TotalReturns:        portfolio.TotalReturns,
		AnnualReturnRate:    portfolio.AnnualReturnRate.String(),
		ActiveAllocations:   int32(activeCount),
		CompletedAllocations: int32(completedCount),
		LastUpdated:         portfolio.LastRebalanced,
	}, nil
}

// Allocation returns details of a specific allocation
func (k Keeper) Allocation(goCtx context.Context, req *types.QueryAllocationRequest) (*types.QueryAllocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allocation, found := k.GetFundAllocation(ctx, req.AllocationId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("allocation %d not found", req.AllocationId))
	}

	return &types.QueryAllocationResponse{
		Allocation: &allocation,
	}, nil
}

// Allocations returns all allocations with pagination
func (k Keeper) Allocations(goCtx context.Context, req *types.QueryAllocationsRequest) (*types.QueryAllocationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var allocations []types.FundAllocation
	store := ctx.KVStore(k.storeKey)

	var allocStore prefix.Store
	if req.Status != "" {
		// Filter by status using index
		allocStore = prefix.NewStore(store, types.GetAllocationByStatusKey(req.Status, 0)[:len(types.GetAllocationByStatusKey(req.Status, 0))-8])
	} else {
		// All allocations
		allocStore = prefix.NewStore(store, types.FundAllocationKey)
	}

	pageRes, err := query.Paginate(allocStore, req.Pagination, func(key []byte, value []byte) error {
		var allocationID uint64
		if req.Status != "" {
			// Extract ID from index key
			allocationID = types.GetUint64FromBytes(key[len(req.Status):])
			// Get actual allocation
			allocationBz := store.Get(types.GetFundAllocationKey(allocationID))
			if allocationBz != nil {
				var allocation types.FundAllocation
				k.cdc.MustUnmarshal(allocationBz, &allocation)
				allocations = append(allocations, allocation)
			}
		} else {
			// Direct allocation storage
			var allocation types.FundAllocation
			k.cdc.MustUnmarshal(value, &allocation)
			allocations = append(allocations, allocation)
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

// Portfolio returns the current investment portfolio
func (k Keeper) Portfolio(goCtx context.Context, req *types.QueryPortfolioRequest) (*types.QueryPortfolioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if !found {
		// Return empty portfolio
		portfolio = types.InvestmentPortfolio{
			TotalValue:       sdk.NewCoin("unamo", sdk.ZeroInt()),
			LiquidAssets:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			InvestedAssets:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			ReservedAssets:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			Components:       []types.PortfolioComponent{},
			TotalReturns:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			AnnualReturnRate: sdk.ZeroDec(),
			RiskScore:        5,
		}
	}

	return &types.QueryPortfolioResponse{
		Portfolio: &portfolio,
	}, nil
}

// MonthlyReports returns monthly reports
func (k Keeper) MonthlyReports(goCtx context.Context, req *types.QueryMonthlyReportsRequest) (*types.QueryMonthlyReportsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var reports []types.MonthlyReport
	store := ctx.KVStore(k.storeKey)
	reportStore := prefix.NewStore(store, types.MonthlyReportKey)

	// Filter by period range if provided
	pageRes, err := query.Paginate(reportStore, req.Pagination, func(key []byte, value []byte) error {
		period := string(key)
		
		// Apply period filters if provided
		if req.FromPeriod != "" && period < req.FromPeriod {
			return nil
		}
		if req.ToPeriod != "" && period > req.ToPeriod {
			return nil
		}

		var report types.MonthlyReport
		k.cdc.MustUnmarshal(value, &report)
		reports = append(reports, report)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMonthlyReportsResponse{
		Reports:    reports,
		Pagination: pageRes,
	}, nil
}

// Governance returns governance parameters
func (k Keeper) Governance(goCtx context.Context, req *types.QueryGovernanceRequest) (*types.QueryGovernanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "governance configuration not found")
	}

	return &types.QueryGovernanceResponse{
		Governance: &governance,
	}, nil
}

// AllocationsByCategory returns allocations filtered by category
func (k Keeper) AllocationsByCategory(goCtx context.Context, req *types.QueryAllocationsByCategoryRequest) (*types.QueryAllocationsByCategoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var allocations []types.FundAllocation
	totalAllocated := sdk.NewCoin("unamo", sdk.ZeroInt())
	
	store := ctx.KVStore(k.storeKey)
	categoryStore := prefix.NewStore(store, types.GetAllocationByCategoryKey(req.Category, 0)[:len(types.GetAllocationByCategoryKey(req.Category, 0))-8])

	pageRes, err := query.Paginate(categoryStore, req.Pagination, func(key []byte, value []byte) error {
		// Extract allocation ID from index key
		allocationID := types.GetUint64FromBytes(key[len(req.Category):])
		
		// Get actual allocation
		allocation, found := k.GetFundAllocation(ctx, allocationID)
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

	return &types.QueryAllocationsByCategoryResponse{
		Allocations:    allocations,
		TotalAllocated: totalAllocated,
		Pagination:     pageRes,
	}, nil
}

// PendingDisbursements returns pending disbursements
func (k Keeper) PendingDisbursements(goCtx context.Context, req *types.QueryPendingDisbursementsRequest) (*types.QueryPendingDisbursementsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var disbursements []types.QueryPendingDisbursementsResponse_PendingDisbursement
	totalPending := sdk.NewCoin("unamo", sdk.ZeroInt())

	store := ctx.KVStore(k.storeKey)
	
	// Iterate through all allocations to find pending disbursements
	allocIterator := sdk.KVStorePrefixIterator(store, types.FundAllocationKey)
	defer allocIterator.Close()

	currentTime := ctx.BlockTime()
	
	for ; allocIterator.Valid(); allocIterator.Next() {
		var allocation types.FundAllocation
		k.cdc.MustUnmarshal(allocIterator.Value(), &allocation)

		// Check each disbursement in the allocation
		for i, disbursement := range allocation.Disbursements {
			if disbursement.Status == "pending" && disbursement.ScheduledDate.Before(currentTime.Add(30*24*3600*1e9)) { // Next 30 days
				disbursements = append(disbursements, types.QueryPendingDisbursementsResponse_PendingDisbursement{
					AllocationId:       allocation.Id,
					AllocationPurpose:  allocation.Purpose,
					DisbursementIndex:  uint32(i),
					Amount:             disbursement.Amount,
					ScheduledDate:      disbursement.ScheduledDate,
					Milestone:          disbursement.Milestone,
					Recipient:          allocation.Recipient,
				})
				totalPending = totalPending.Add(disbursement.Amount)
			}
		}
	}

	// Apply pagination manually
	start := 0
	end := len(disbursements)
	
	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			start = int(req.Pagination.Offset)
			if start >= len(disbursements) {
				start = len(disbursements)
			}
		}
		if req.Pagination.Limit > 0 {
			end = start + int(req.Pagination.Limit)
			if end > len(disbursements) {
				end = len(disbursements)
			}
		}
	}

	pageRes := &query.PageResponse{
		Total: uint64(len(disbursements)),
	}

	return &types.QueryPendingDisbursementsResponse{
		Disbursements: disbursements[start:end],
		TotalPending:  totalPending,
		Pagination:    pageRes,
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