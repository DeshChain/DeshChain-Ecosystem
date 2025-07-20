package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/krishimitra/types"
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

// Loan queries a specific agricultural loan by ID
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

// LoansByFarmer queries all loans for a specific farmer
func (k Keeper) LoansByFarmer(c context.Context, req *types.QueryLoansByFarmerRequest) (*types.QueryLoansByFarmerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Farmer == "" {
		return nil, status.Error(codes.InvalidArgument, "farmer address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.AgriculturalLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.AgriculturalLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.Borrower == req.Farmer {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLoansByFarmerResponse{
		Loans:      loans,
		Pagination: pageRes,
	}, nil
}

// LoansByPincode queries all loans for a specific pincode
func (k Keeper) LoansByPincode(c context.Context, req *types.QueryLoansByPincodeRequest) (*types.QueryLoansByPincodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Pincode) != 6 {
		return nil, status.Error(codes.InvalidArgument, "invalid pincode")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.AgriculturalLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.AgriculturalLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.Pincode == req.Pincode {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Calculate summary statistics
	summary := k.calculateLoanSummary(loans)

	return &types.QueryLoansByPincodeResponse{
		Loans:      loans,
		Pagination: pageRes,
		Summary:    summary,
	}, nil
}

// FarmerProfile queries a farmer's profile
func (k Keeper) FarmerProfile(c context.Context, req *types.QueryFarmerProfileRequest) (*types.QueryFarmerProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Farmer == "" {
		return nil, status.Error(codes.InvalidArgument, "farmer address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	profile, found := k.GetFarmerProfile(ctx, req.Farmer)
	if !found {
		return nil, status.Error(codes.NotFound, "farmer profile not found")
	}

	// Get active loans
	var activeLoans []types.AgriculturalLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.AgriculturalLoan
		if err := k.cdc.Unmarshal(iterator.Value(), &loan); err != nil {
			continue
		}
		
		if loan.Borrower == req.Farmer && loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE {
			activeLoans = append(activeLoans, loan)
		}
	}

	// Calculate totals
	totalBorrowed := sdk.ZeroInt()
	totalRepaid := sdk.ZeroInt()
	for _, loan := range activeLoans {
		totalBorrowed = totalBorrowed.Add(loan.LoanAmount.Amount)
		totalRepaid = totalRepaid.Add(loan.RepaidAmount.Amount)
	}

	return &types.QueryFarmerProfileResponse{
		Profile:       &profile,
		ActiveLoans:   activeLoans,
		TotalBorrowed: totalBorrowed.String(),
		TotalRepaid:   totalRepaid.String(),
	}, nil
}

// WeatherData queries weather data for agricultural planning
func (k Keeper) WeatherData(c context.Context, req *types.QueryWeatherDataRequest) (*types.QueryWeatherDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Pincode) != 6 {
		return nil, status.Error(codes.InvalidArgument, "invalid pincode")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Get weather data from store
	weatherData, found := k.GetWeatherData(ctx, req.Pincode)
	if !found {
		// Return mock data for demonstration
		weatherData = types.WeatherInfo{
			Pincode:     req.Pincode,
			Temperature: "28°C",
			Humidity:    "65%",
			Rainfall:    "12mm",
			WindSpeed:   "10 km/h",
			Forecast:    "Partly cloudy with chance of rain",
			LastUpdated: ctx.BlockTime(),
		}
	}

	// Get crop recommendations based on weather
	recommendations := k.GetCropRecommendations(ctx, req.Pincode, weatherData)

	return &types.QueryWeatherDataResponse{
		Weather:         &weatherData,
		Recommendations: recommendations,
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
		TotalLoans:          stats.TotalLoans,
		TotalDisbursed:      stats.TotalDisbursed,
		TotalRepaid:         stats.TotalRepaid,
		ActiveLoans:         stats.ActiveLoans,
		DefaultedLoans:      stats.DefaultedLoans,
		AverageLoanAmount:   stats.AverageLoanAmount,
		AverageInterestRate: stats.AverageInterestRate,
		LoansByCrop:         stats.LoansByCrop,
		LoansByPincode:      stats.LoansByPincode,
		DefaultRate:         stats.DefaultRate,
		SuccessfulHarvests:  stats.SuccessfulHarvests,
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

// EligibilityCheck checks loan eligibility for a farmer
func (k Keeper) EligibilityCheck(c context.Context, req *types.QueryEligibilityCheckRequest) (*types.QueryEligibilityCheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Farmer == "" {
		return nil, status.Error(codes.InvalidArgument, "farmer address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Parse loan amount
	loanAmount, err := sdk.ParseCoinNormalized(req.LoanAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid loan amount")
	}

	// Check eligibility
	eligible, maxAmount, reasons := k.CheckFarmerEligibility(ctx, req.Farmer, loanAmount, req.CropType, req.LandSize)

	// Calculate estimated interest rate
	estimatedRate := k.EstimateInterestRate(ctx, req.Farmer, req.CropType, loanAmount)

	// Get applicable discounts
	discounts := k.GetApplicableDiscounts(ctx, req.Farmer, req.CropType)

	// Calculate estimated EMI
	duration := int64(6) // 6 months default
	estimatedEMI := k.CalculateEMI(ctx, loanAmount, estimatedRate, duration)

	// Get cultural tip
	culturalTip := k.culturalKeeper.GetRandomCulturalQuote(ctx)

	return &types.QueryEligibilityCheckResponse{
		Eligible:              eligible,
		MaxLoanAmount:         maxAmount.String(),
		EstimatedInterestRate: estimatedRate.String(),
		RequiredDocuments:     k.GetRequiredDocuments(req.CropType),
		ApplicableDiscounts:   discounts,
		EstimatedEmi:          estimatedEMI.String(),
		EligibilityReasons:    reasons,
		CulturalTip:           culturalTip,
	}, nil
}

// Helper functions

func (k Keeper) calculateLoanSummary(loans []types.AgriculturalLoan) *types.LoanSummary {
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
		ActiveFarmers:       activeCount,
		RepaymentRate:       k.CalculateRepaymentRate(loans),
	}
}

func (k Keeper) CalculateRepaymentRate(loans []types.AgriculturalLoan) string {
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