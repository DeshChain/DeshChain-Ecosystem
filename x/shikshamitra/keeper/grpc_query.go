package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/shikshamitra/types"
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

// Loan queries a specific education loan by ID
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

// LoansByStudent queries all loans for a specific student
func (k Keeper) LoansByStudent(c context.Context, req *types.QueryLoansByStudentRequest) (*types.QueryLoansByStudentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Student == "" {
		return nil, status.Error(codes.InvalidArgument, "student address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.EducationLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.EducationLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.Borrower == req.Student {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLoansByStudentResponse{
		Loans:      loans,
		Pagination: pageRes,
	}, nil
}

// LoansByInstitution queries all loans for a specific institution
func (k Keeper) LoansByInstitution(c context.Context, req *types.QueryLoansByInstitutionRequest) (*types.QueryLoansByInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Institution == "" {
		return nil, status.Error(codes.InvalidArgument, "institution name cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var loans []types.EducationLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)

	pageRes, err := query.Paginate(loanStore, req.Pagination, func(key []byte, value []byte) error {
		var loan types.EducationLoan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}
		
		if loan.InstitutionName == req.Institution {
			loans = append(loans, loan)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Calculate summary statistics
	summary := k.calculateInstitutionSummary(loans)

	return &types.QueryLoansByInstitutionResponse{
		Loans:      loans,
		Pagination: pageRes,
		Summary:    summary,
	}, nil
}

// StudentProfile queries a student's profile
func (k Keeper) StudentProfile(c context.Context, req *types.QueryStudentProfileRequest) (*types.QueryStudentProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Student == "" {
		return nil, status.Error(codes.InvalidArgument, "student address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	profile, found := k.GetStudentProfile(ctx, req.Student)
	if !found {
		return nil, status.Error(codes.NotFound, "student profile not found")
	}

	// Get active loans
	var activeLoans []types.EducationLoan
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	totalBorrowed := sdk.ZeroInt()
	totalRepaid := sdk.ZeroInt()
	totalScholarship := sdk.ZeroInt()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		if err := k.cdc.Unmarshal(iterator.Value(), &loan); err != nil {
			continue
		}
		
		if loan.Borrower == req.Student {
			if loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE || 
			   loan.Status == types.LoanStatus_LOAN_STATUS_IN_MORATORIUM ||
			   loan.Status == types.LoanStatus_LOAN_STATUS_REPAYMENT {
				activeLoans = append(activeLoans, loan)
			}
			
			totalBorrowed = totalBorrowed.Add(loan.LoanAmount.Amount)
			totalRepaid = totalRepaid.Add(loan.RepaidAmount.Amount)
			
			// Sum scholarships
			for _, scholarship := range loan.Scholarships {
				if scholarship.Status == "DISBURSED" {
					totalScholarship = totalScholarship.Add(scholarship.Amount.Amount)
				}
			}
		}
	}

	return &types.QueryStudentProfileResponse{
		Profile:           &profile,
		ActiveLoans:       activeLoans,
		TotalBorrowed:     totalBorrowed.String(),
		TotalRepaid:       totalRepaid.String(),
		ScholarshipAmount: totalScholarship.String(),
	}, nil
}

// ActiveScholarships queries active scholarship opportunities
func (k Keeper) ActiveScholarships(c context.Context, req *types.QueryActiveScholarshipsRequest) (*types.QueryActiveScholarshipsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	var scholarships []types.ScholarshipOpportunity
	
	// Get all active scholarships from store
	store := ctx.KVStore(k.storeKey)
	scholarshipStore := prefix.NewStore(store, types.ScholarshipKeyPrefix)

	pageRes, err := query.Paginate(scholarshipStore, req.Pagination, func(key []byte, value []byte) error {
		var scholarship types.ScholarshipOpportunity
		if err := k.cdc.Unmarshal(value, &scholarship); err != nil {
			return err
		}
		
		// Filter by course type if provided
		if req.CourseType != "" && !k.isEligibleCourse(scholarship.EligibleCourses, req.CourseType) {
			return nil
		}
		
		// Filter by institution type if provided
		if req.InstitutionType != "" && !k.isEligibleInstitution(scholarship.EligibleInstitutions, req.InstitutionType) {
			return nil
		}
		
		// Check if still active
		if ctx.BlockTime().Before(scholarship.ApplicationDeadline) {
			scholarships = append(scholarships, scholarship)
		}
		
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryActiveScholarshipsResponse{
		Scholarships: scholarships,
		Pagination:   pageRes,
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
		LoansInMoratorium:   stats.LoansInMoratorium,
		AverageLoanAmount:   stats.AverageLoanAmount,
		AverageInterestRate: stats.AverageInterestRate,
		LoansByCourse:       stats.LoansByCourse,
		LoansByInstitution:  stats.LoansByInstitution,
		DefaultRate:         stats.DefaultRate,
		AverageCgpa:         stats.AverageCGPA,
		EmploymentRate:      stats.EmploymentRate,
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

// EligibilityCheck checks loan eligibility for a student
func (k Keeper) EligibilityCheck(c context.Context, req *types.QueryEligibilityCheckRequest) (*types.QueryEligibilityCheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Student == "" {
		return nil, status.Error(codes.InvalidArgument, "student address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Parse loan amount
	loanAmount, err := sdk.ParseCoinNormalized(req.LoanAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid loan amount")
	}

	// Parse family income
	familyIncome, err := sdk.ParseCoinNormalized(req.FamilyIncome)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid family income")
	}

	// Check eligibility
	eligible, maxAmount, collateralRequired, reasons := k.CheckStudentEligibility(
		ctx, req.Student, loanAmount, req.InstitutionType, req.CourseType, familyIncome,
	)

	// Calculate estimated interest rate
	academicScore, _ := sdk.NewDecFromStr(req.AcademicScore)
	estimatedRate := k.EstimateInterestRate(ctx, req.Student, req.InstitutionType, req.CourseType, academicScore, familyIncome)

	// Get applicable discounts
	discounts := k.GetApplicableDiscounts(ctx, req.Student, req.InstitutionType, req.CourseType, academicScore)

	// Calculate estimated EMI (after moratorium)
	courseDuration := k.GetCourseDuration(req.CourseType)
	moratoriumPeriod := k.GetParams(ctx).MoratoriumPeriod
	repaymentMonths := int64(120) // 10 years default
	estimatedEMI := k.CalculateEMI(ctx, loanAmount, estimatedRate, repaymentMonths)

	// Get required documents
	documents := k.GetRequiredDocuments(req.CourseType, collateralRequired)

	// Get cultural tip
	culturalTip := k.culturalKeeper.GetRandomCulturalQuote(ctx)

	return &types.QueryEligibilityCheckResponse{
		Eligible:              eligible,
		EstimatedInterestRate: estimatedRate.String(),
		MaxLoanAmount:         maxAmount.String(),
		CollateralRequired:    collateralRequired,
		ApplicableDiscounts:   discounts,
		EstimatedEmi:          estimatedEMI.String(),
		Requirements:          reasons,
		CulturalTip:           culturalTip,
	}, nil
}

// Helper functions

func (k Keeper) calculateInstitutionSummary(loans []types.EducationLoan) *types.LoanSummary {
	totalAmount := sdk.ZeroInt()
	totalInterest := sdk.ZeroDec()
	activeCount := int64(0)
	
	for _, loan := range loans {
		totalAmount = totalAmount.Add(loan.LoanAmount.Amount)
		rate, _ := sdk.NewDecFromStr(loan.InterestRate)
		totalInterest = totalInterest.Add(rate)
		
		if loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE ||
		   loan.Status == types.LoanStatus_LOAN_STATUS_IN_MORATORIUM ||
		   loan.Status == types.LoanStatus_LOAN_STATUS_REPAYMENT {
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
		ActiveStudents:      activeCount,
		RepaymentRate:       k.CalculateRepaymentRate(loans),
	}
}

func (k Keeper) CalculateRepaymentRate(loans []types.EducationLoan) string {
	if len(loans) == 0 {
		return "0%"
	}
	
	repaidCount := 0
	totalEligible := 0
	
	for _, loan := range loans {
		// Only count loans that have passed moratorium
		if loan.Status == types.LoanStatus_LOAN_STATUS_REPAID {
			repaidCount++
			totalEligible++
		} else if loan.Status == types.LoanStatus_LOAN_STATUS_REPAYMENT {
			totalEligible++
		}
	}
	
	if totalEligible == 0 {
		return "N/A"
	}
	
	rate := sdk.NewDec(int64(repaidCount)).Mul(sdk.NewDec(100)).Quo(sdk.NewDec(int64(totalEligible)))
	return rate.String() + "%"
}

func (k Keeper) isEligibleCourse(eligibleCourses []string, courseType string) bool {
	if len(eligibleCourses) == 0 {
		return true
	}
	
	for _, course := range eligibleCourses {
		if course == courseType {
			return true
		}
	}
	return false
}

func (k Keeper) isEligibleInstitution(eligibleInstitutions []string, institutionType string) bool {
	if len(eligibleInstitutions) == 0 {
		return true
	}
	
	for _, inst := range eligibleInstitutions {
		if inst == institutionType {
			return true
		}
	}
	return false
}

func (k Keeper) GetCourseDuration(courseType string) int64 {
	// Return typical course durations in months
	switch courseType {
	case "ENGINEERING":
		return 48 // 4 years
	case "MEDICAL":
		return 66 // 5.5 years
	case "MANAGEMENT":
		return 24 // 2 years
	case "LAW":
		return 36 // 3 years
	case "PHD":
		return 60 // 5 years
	default:
		return 36 // 3 years default
	}
}