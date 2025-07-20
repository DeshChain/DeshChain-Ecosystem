package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/vyavasayamitra/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Loan queries a specific business loan by ID
func (k Keeper) Loan(c context.Context, req *types.QueryLoanRequest) (*types.QueryLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.LoanId == "" {
		return nil, status.Error(codes.InvalidArgument, "loan ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	loan, found := k.GetLoan(ctx, req.LoanId)
	if !found {
		return nil, status.Error(codes.NotFound, "loan not found")
	}

	return &types.QueryLoanResponse{Loan: &loan}, nil
}

// LoansByBusiness queries all loans for a specific business
func (k Keeper) LoansByBusiness(c context.Context, req *types.QueryLoansByBusinessRequest) (*types.QueryLoansByBusinessResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Business == "" {
		return nil, status.Error(codes.InvalidArgument, "business address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.BusinessLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.BusinessLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.Borrower == req.Business {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLoansByBusinessResponse{
		Loans:      loans,
		Pagination: pageRes,
	}, nil
}

// LoansByIndustry queries all loans for a specific industry
func (k Keeper) LoansByIndustry(c context.Context, req *types.QueryLoansByIndustryRequest) (*types.QueryLoansByIndustryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.BusinessLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.BusinessLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.BusinessType.String() == req.Industry {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Calculate summary statistics
	summary := k.calculateIndustrySummary(loans)

	return &types.QueryLoansByIndustryResponse{
		Loans:      loans,
		Pagination: pageRes,
		Summary:    summary,
	}, nil
}

// BusinessProfile queries a business's profile
func (k Keeper) BusinessProfile(c context.Context, req *types.QueryBusinessProfileRequest) (*types.QueryBusinessProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Business == "" {
		return nil, status.Error(codes.InvalidArgument, "business address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	profile, found := k.GetBusinessProfile(ctx, req.Business)
	if !found {
		return nil, status.Error(codes.NotFound, "business profile not found")
	}

	// Get active loans
	var activeLoans []types.BusinessLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.BusinessLoan
		if err := k.cdc.Unmarshal(iterator.Value(), &loan); err != nil {
			continue
		}
		
		if loan.Borrower == req.Business && loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE {
			activeLoans = append(activeLoans, loan)
		}
	}

	// Calculate credit line utilization
	creditUtilization := k.CalculateCreditUtilization(activeLoans)

	return &types.QueryBusinessProfileResponse{
		Profile:           &profile,
		ActiveLoans:       activeLoans,
		TotalBorrowed:     profile.TotalBorrowed.String(),
		TotalRepaid:       profile.TotalRepaid.String(),
		CreditUtilization: creditUtilization,
	}, nil
}

// InvoiceFinancings queries invoice financing records
func (k Keeper) InvoiceFinancings(c context.Context, req *types.QueryInvoiceFinancingsRequest) (*types.QueryInvoiceFinancingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.LoanId == "" {
		return nil, status.Error(codes.InvalidArgument, "loan ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	loan, found := k.GetLoan(ctx, req.LoanId)
	if !found {
		return nil, status.Error(codes.NotFound, "loan not found")
	}

	// Filter by status if provided
	var financings []*types.InvoiceFinancing
	for _, inv := range loan.InvoiceFinancings {
		if req.Status == "" || inv.Status == req.Status {
			financings = append(financings, inv)
		}
	}

	// Calculate summary
	totalFinanced := sdk.ZeroInt()
	totalPending := sdk.ZeroInt()
	for _, inv := range financings {
		totalFinanced = totalFinanced.Add(inv.FinancedAmount.Amount)
		if inv.Status == "PENDING" {
			totalPending = totalPending.Add(inv.InvoiceAmount.Amount)
		}
	}

	return &types.QueryInvoiceFinancingsResponse{
		Financings:    financings,
		TotalFinanced: totalFinanced.String(),
		TotalPending:  totalPending.String(),
	}, nil
}

// LoanStatistics queries overall loan statistics
func (k Keeper) LoanStatistics(c context.Context, req *types.QueryLoanStatisticsRequest) (*types.QueryLoanStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	stats := k.CalculateLoanStatistics(ctx, req.Period)

	return &types.QueryLoanStatisticsResponse{
		TotalLoans:            stats.TotalLoans,
		TotalDisbursed:        stats.TotalDisbursed,
		TotalRepaid:           stats.TotalRepaid,
		ActiveLoans:           stats.ActiveLoans,
		TotalCreditLines:      stats.TotalCreditLines,
		AverageLoanAmount:     stats.AverageLoanAmount,
		AverageInterestRate:   stats.AverageInterestRate,
		LoansByBusinessType:   stats.LoansByBusinessType,
		LoansByPurpose:        stats.LoansByPurpose,
		DefaultRate:           stats.DefaultRate,
		InvoiceFinancingVolume: stats.InvoiceFinancingVolume,
		AverageBusinessAge:    stats.AverageBusinessAge,
	}, nil
}

// FestivalOffers queries active festival offers
func (k Keeper) FestivalOffers(c context.Context, req *types.QueryFestivalOffersRequest) (*types.QueryFestivalOffersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	activeOffers := k.GetActiveFestivalOffers(ctx)
	upcomingOffer := k.GetNextFestivalOffer(ctx)

	return &types.QueryFestivalOffersResponse{
		ActiveOffers:  activeOffers,
		UpcomingOffer: upcomingOffer,
	}, nil
}

// EligibilityCheck checks loan eligibility for a business
func (k Keeper) EligibilityCheck(c context.Context, req *types.QueryEligibilityCheckRequest) (*types.QueryEligibilityCheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Business == "" {
		return nil, status.Error(codes.InvalidArgument, "business address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Parse loan amount
	loanAmount, err := sdk.ParseCoinNormalized(req.LoanAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid loan amount")
	}

	// Parse annual revenue
	annualRevenue, err := sdk.ParseCoinNormalized(req.AnnualRevenue)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid annual revenue")
	}

	// Check eligibility
	eligible, maxAmount, reasons := k.CheckBusinessEligibility(ctx, req.Business, loanAmount, req.BusinessType, annualRevenue)

	// Calculate estimated interest rate
	estimatedRate := k.EstimateInterestRate(ctx, req.Business, req.BusinessType, loanAmount, req.CreditScore)

	// Get applicable discounts
	discounts := k.GetApplicableDiscounts(ctx, req.Business, req.BusinessType)

	// Calculate estimated EMI
	duration, _ := sdk.NewDecFromStr(req.Duration)
	estimatedEMI := k.CalculateEMI(ctx, loanAmount, estimatedRate, duration.TruncateInt64())

	// Check if eligible for credit line
	creditLineEligible := k.IsEligibleForCreditLine(ctx, req.Business, annualRevenue)

	// Get cultural tip
	culturalTip := k.culturalKeeper.GetRandomCulturalQuote(ctx)

	return &types.QueryEligibilityCheckResponse{
		Eligible:               eligible,
		MaxLoanAmount:          maxAmount.String(),
		EstimatedInterestRate:  estimatedRate.String(),
		CreditLineEligible:     creditLineEligible,
		MaxCreditLine:          k.CalculateMaxCreditLine(annualRevenue).String(),
		RequiredDocuments:      k.GetRequiredDocuments(req.BusinessType),
		ApplicableDiscounts:    discounts,
		EstimatedEmi:           estimatedEMI.String(),
		EligibilityReasons:     reasons,
		CulturalTip:            culturalTip,
	}, nil
}

// Helper functions

func (k Keeper) calculateIndustrySummary(loans []types.BusinessLoan) *types.LoanSummary {
	totalAmount := sdk.ZeroInt()
	totalInterest := sdk.ZeroDec()
	activeCount := int64(0)
	
	for _, loan := range loans {
		totalAmount = totalAmount.Add(loan.LoanAmount.Amount)
		rate, _ := sdk.NewDecFromStr(loan.InterestRate)
		totalInterest = totalInterest.Add(rate)
		
		if loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE {
			activeCount++
		}
	}
	
	avgInterestRate := sdk.ZeroDec()
	if len(loans) > 0 {
		avgInterestRate = totalInterest.Quo(sdk.NewDec(int64(len(loans))))
	}
	
	return &types.LoanSummary{
		TotalLoans:          int64(len(loans)),
		TotalAmount:         totalAmount.String(),
		AverageInterestRate: avgInterestRate.String(),
		ActiveBusinesses:    activeCount,
		RepaymentRate:       k.CalculateRepaymentRate(loans),
	}
}

func (k Keeper) CalculateCreditUtilization(loans []types.BusinessLoan) string {
	totalLimit := sdk.ZeroInt()
	totalUtilized := sdk.ZeroInt()
	
	for _, loan := range loans {
		if loan.CreditLine != nil && loan.CreditLine.IsActive {
			totalLimit = totalLimit.Add(loan.CreditLine.CreditLimit.Amount)
			totalUtilized = totalUtilized.Add(loan.CreditLine.UtilizedAmount.Amount)
		}
	}
	
	if totalLimit.IsZero() {
		return "0%"
	}
	
	utilization := totalUtilized.Mul(sdk.NewInt(100)).Quo(totalLimit)
	return utilization.String() + "%"
}

func (k Keeper) CalculateRepaymentRate(loans []types.BusinessLoan) string {
	if len(loans) == 0 {
		return "0%"
	}
	
	repaidCount := 0
	for _, loan := range loans {
		if loan.Status == types.LoanStatus_LOAN_STATUS_REPAID {
			repaidCount++
		}
	}
	
	rate := sdk.NewDec(int64(repaidCount)).Mul(sdk.NewDec(100)).Quo(sdk.NewDec(int64(len(loans))))
	return rate.String() + "%"
}