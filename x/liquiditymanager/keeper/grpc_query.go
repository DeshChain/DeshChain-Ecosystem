package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/liquiditymanager/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier implements the Query gRPC service
var _ types.QueryServer = Keeper{}

// LiquidityStatus returns the current liquidity status
func (k Keeper) LiquidityStatus(c context.Context, req *types.QueryLiquidityStatusRequest) (*types.QueryLiquidityStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	info := k.GetLiquidityInfo(ctx)

	return &types.QueryLiquidityStatusResponse{
		TotalPoolValue:      info.TotalPoolValue.String(),
		AvailableForLending: info.AvailableForLending.String(),
		ReserveAmount:       info.ReserveAmount.String(),
		Status:              string(info.Status),
		MaxLoanAmount:       info.MaxLoanAmount.String(),
		DailyLendingLimit:   info.DailyLendingLimit.String(),
		AvailableModules:    info.AvailableModules,
		NextThreshold:       info.NextThreshold.String(),
		ProgressToNext:      info.ProgressToNext.String(),
		EstimatedDaysToNext: info.EstimatedDaysToNext,
	}, nil
}

// LendingAvailability checks if lending is available for a specific amount and module
func (k Keeper) LendingAvailability(c context.Context, req *types.QueryLendingAvailabilityRequest) (*types.QueryLendingAvailabilityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Amount == "" {
		return nil, status.Error(codes.InvalidArgument, "amount is required")
	}

	if req.Module == "" {
		return nil, status.Error(codes.InvalidArgument, "module is required")
	}

	if req.Borrower == "" {
		return nil, status.Error(codes.InvalidArgument, "borrower address is required")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	amount, err := sdk.NewDecFromStr(req.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount format")
	}

	borrower, err := sdk.AccAddressFromBech32(req.Borrower)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid borrower address")
	}

	available, message := k.CanProcessLoan(ctx, amount, req.Module, borrower)

	return &types.QueryLendingAvailabilityResponse{
		Available: available,
		Message:   message,
	}, nil
}

// PoolProgress returns detailed progress information towards next threshold
func (k Keeper) PoolProgress(c context.Context, req *types.QueryPoolProgressRequest) (*types.QueryPoolProgressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	info := k.GetLiquidityInfo(ctx)
	params := k.GetParams(ctx)

	// Calculate milestone progress
	basicThreshold := sdk.NewDecFromInt(params.LendingBasicThreshold)
	mediumThreshold := sdk.NewDecFromInt(params.LendingMediumThreshold)
	fullThreshold := sdk.NewDecFromInt(params.LendingFullThreshold)

	milestones := []*types.PoolMilestone{
		{
			Name:        "🏛️ Basic Lending (BANK-GRADE SAFETY)",
			Description: "Conservative loans up to ₹50K with 75% reserves (Krishi Mitra only)",
			Threshold:   basicThreshold.String(),
			Achieved:    info.TotalPoolValue.GTE(basicThreshold),
			Benefits:    []string{"₹50K max loans", "Krishi Mitra lending", "1% daily limit", "75% total reserves", "10% default provision"},
		},
		{
			Name:        "🚀 Medium Lending (ENHANCED SECURITY)",
			Description: "Secure loans up to ₹2L with comprehensive risk management (Krishi + Vyavasaya)",
			Threshold:   mediumThreshold.String(),
			Achieved:    info.TotalPoolValue.GTE(mediumThreshold),
			Benefits:    []string{"₹2L max loans", "Business lending", "2% daily limit", "Advanced risk controls", "Pool member protection"},
		},
		{
			Name:        "💎 Full Lending (REVOLUTIONARY PLATFORM)",
			Description: "Maximum loans up to ₹5L with world-class financial safety (All modules)",
			Threshold:   fullThreshold.String(),
			Achieved:    info.TotalPoolValue.GTE(fullThreshold),
			Benefits:    []string{"₹5L max loans", "Education lending", "3% daily limit", "Industry-leading reserves", "Complete DeFi ecosystem"},
		},
	}

	return &types.QueryPoolProgressResponse{
		CurrentValue:        info.TotalPoolValue.String(),
		NextThreshold:       info.NextThreshold.String(),
		ProgressPercentage:  info.ProgressToNext.Mul(sdk.NewDec(100)).String(),
		EstimatedDaysToNext: info.EstimatedDaysToNext,
		Milestones:          milestones,
		Status:              string(info.Status),
	}, nil
}

// DailyLendingStats returns daily lending statistics
func (k Keeper) DailyLendingStats(c context.Context, req *types.QueryDailyLendingStatsRequest) (*types.QueryDailyLendingStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	info := k.GetLiquidityInfo(ctx)
	dailyUsed := k.getDailyLendingUsed(ctx)
	
	remaining := info.DailyLendingLimit.Sub(dailyUsed)
	if remaining.IsNegative() {
		remaining = sdk.ZeroDec()
	}

	utilizationPercentage := sdk.ZeroDec()
	if info.DailyLendingLimit.GT(sdk.ZeroDec()) {
		utilizationPercentage = dailyUsed.Quo(info.DailyLendingLimit).Mul(sdk.NewDec(100))
	}

	return &types.QueryDailyLendingStatsResponse{
		DailyLimit:            info.DailyLendingLimit.String(),
		DailyUsed:             dailyUsed.String(),
		DailyRemaining:        remaining.String(),
		UtilizationPercentage: utilizationPercentage.String(),
		Status:                string(info.Status),
	}, nil
}

// LendingQueue returns information about pending loan applications
func (k Keeper) LendingQueue(c context.Context, req *types.QueryLendingQueueRequest) (*types.QueryLendingQueueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Get queue information (placeholder implementation)
	queueSize := k.getQueueSize(ctx)
	estimatedWaitTime := k.getEstimatedWaitTime(ctx)
	
	return &types.QueryLendingQueueResponse{
		QueueSize:           queueSize,
		EstimatedWaitTimeHours: estimatedWaitTime,
		AcceptingNewApplications: k.IsLendingAvailable(ctx),
	}, nil
}

// CollateralLoanAvailability checks if a NAMO collateral loan can be processed
func (k Keeper) CollateralLoanAvailability(c context.Context, req *types.QueryCollateralLoanAvailabilityRequest) (*types.QueryCollateralLoanAvailabilityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.LoanAmount == "" {
		return nil, status.Error(codes.InvalidArgument, "loan amount is required")
	}

	if req.CollateralAmount == "" {
		return nil, status.Error(codes.InvalidArgument, "collateral amount is required")
	}

	if req.Borrower == "" {
		return nil, status.Error(codes.InvalidArgument, "borrower address is required")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	loanAmount, err := sdk.NewDecFromStr(req.LoanAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid loan amount format")
	}

	collateralAmount, err := sdk.NewDecFromStr(req.CollateralAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collateral amount format")
	}

	borrower, err := sdk.AccAddressFromBech32(req.Borrower)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid borrower address")
	}

	available, message := k.CanProcessCollateralLoan(ctx, loanAmount, collateralAmount, borrower)

	// Calculate max loan amount based on collateral
	maxLoanAmount := collateralAmount.Mul(sdk.NewDecWithPrec(70, 2)) // 70% LTV
	
	// Get user's staking info
	stakedAmount := k.GetStakedNAMO(ctx, borrower)
	lockedAmount := k.getLockedCollateral(ctx, borrower)
	availableStake := stakedAmount.Sub(lockedAmount)

	return &types.QueryCollateralLoanAvailabilityResponse{
		Available:           available,
		Message:             message,
		MaxLoanAmount:       maxLoanAmount.String(),
		RequiredCollateral:  collateralAmount.String(),
		UserStakedAmount:    stakedAmount.String(),
		UserLockedAmount:    lockedAmount.String(),
		UserAvailableStake:  availableStake.String(),
		LoanToValueRatio:    "70.00", // Fixed 70% LTV
	}, nil
}

// UserStakeInfo returns detailed staking information for a user
func (k Keeper) UserStakeInfo(c context.Context, req *types.QueryUserStakeInfoRequest) (*types.QueryUserStakeInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.User == "" {
		return nil, status.Error(codes.InvalidArgument, "user address is required")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	user, err := sdk.AccAddressFromBech32(req.User)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user address")
	}

	stakedAmount := k.GetStakedNAMO(ctx, user)
	lockedAmount := k.getLockedCollateral(ctx, user)
	availableStake := stakedAmount.Sub(lockedAmount)
	
	// Check pool memberships
	isVillageMember := k.isVillageSurakshaPoolMember(ctx, user)
	isUrbanMember := k.isUrbanSurakshaPoolMember(ctx, user)
	isPoolMember := k.IsPoolMember(ctx, user)

	// Calculate max borrowing capacity
	maxBorrowCapacity := availableStake.Mul(sdk.NewDecWithPrec(70, 2)) // 70% of available stake

	return &types.QueryUserStakeInfoResponse{
		StakedAmount:        stakedAmount.String(),
		LockedCollateral:    lockedAmount.String(),
		AvailableStake:      availableStake.String(),
		MaxBorrowCapacity:   maxBorrowCapacity.String(),
		IsVillagePoolMember: isVillageMember,
		IsUrbanPoolMember:   isUrbanMember,
		IsEligibleForLending: isPoolMember,
		LoanToValueRatio:    "70.00",
	}, nil
}

// LoanBreakdown returns comprehensive loan cost breakdown for transparency
func (k Keeper) LoanBreakdown(c context.Context, req *types.QueryLoanBreakdownRequest) (*types.QueryLoanBreakdownResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.LoanAmount == "" {
		return nil, status.Error(codes.InvalidArgument, "loan amount is required")
	}

	if req.InterestRate == "" {
		return nil, status.Error(codes.InvalidArgument, "interest rate is required")
	}

	if req.TermMonths == 0 {
		return nil, status.Error(codes.InvalidArgument, "term months is required")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	loanAmount, err := sdk.NewDecFromStr(req.LoanAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid loan amount format")
	}

	interestRate, err := sdk.NewDecFromStr(req.InterestRate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid interest rate format")
	}

	breakdown := k.GetLoanBreakdown(ctx, loanAmount, interestRate, req.TermMonths)

	return &types.QueryLoanBreakdownResponse{
		ApprovedAmount:   breakdown.ApprovedAmount.String(),
		ProcessingFee:    breakdown.ProcessingFee.String(),
		DisbursedAmount:  breakdown.DisbursedAmount.String(),
		InterestRate:     breakdown.InterestRate.String(),
		TotalInterest:    breakdown.TotalInterest.String(),
		TotalRepayment:   breakdown.TotalRepayment.String(),
		MonthlyEmi:       breakdown.MonthlyEMI.String(),
		TermMonths:       breakdown.TermMonths,
		EffectiveFeeRate: breakdown.EffectiveFeeRate.String(),
		
		// REVOLUTIONARY TRANSPARENCY
		ProcessingFeeDetails: fmt.Sprintf("1%% processing fee (capped at ₹2,500). Your fee: ₹%.0f", 
			breakdown.ProcessingFee.Quo(sdk.NewDec(1000000)).TruncateInt64()),
		CompetitiveComparison: "🏆 50-70% cheaper than traditional banks (who charge 2-5% + hidden fees)",
		TransparencyNote: "💎 Complete transparency - no hidden charges, no surprises",
	}, nil
}

// FeeCalculator returns processing and early settlement fee calculations
func (k Keeper) FeeCalculator(c context.Context, req *types.QueryFeeCalculatorRequest) (*types.QueryFeeCalculatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var processingFee sdk.Dec
	var disbursedAmount sdk.Dec
	var earlySettlementFee sdk.Dec

	// Calculate processing fee if loan amount provided
	if req.LoanAmount != "" {
		loanAmount, err := sdk.NewDecFromStr(req.LoanAmount)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid loan amount format")
		}
		
		processingFee = k.CalculateProcessingFee(ctx, loanAmount)
		disbursedAmount, _ = k.CalculateDisbursementAmount(ctx, loanAmount)
	}

	// Calculate early settlement fee if remaining principal provided
	if req.RemainingPrincipal != "" {
		remainingPrincipal, err := sdk.NewDecFromStr(req.RemainingPrincipal)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid remaining principal format")
		}
		
		earlySettlementFee = k.CalculateEarlySettlementFee(ctx, remainingPrincipal)
	}

	return &types.QueryFeeCalculatorResponse{
		ProcessingFee:       processingFee.String(),
		DisbursedAmount:     disbursedAmount.String(),
		EarlySettlementFee:  earlySettlementFee.String(),
		ProcessingFeeRate:   "1.00%",
		ProcessingFeeCap:    "₹2,500",
		EarlySettlementRate: "0.50%",
		BorrowerAdvantage:   "🎯 Ultra-competitive rates vs banks (2-5% processing + 2-4% early settlement)",
	}, nil
}

// Helper functions for queue management
func (k Keeper) getQueueSize(ctx sdk.Context) int64 {
	// This would return the actual queue size from store
	// For now, return placeholder
	return 0
}

func (k Keeper) getEstimatedWaitTime(ctx sdk.Context) int64 {
	// This would calculate wait time based on queue and lending capacity
	// For now, return placeholder
	return 0
}