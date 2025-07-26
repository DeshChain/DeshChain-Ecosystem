package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/DeshChain/DeshChain-Ecosystem/x/shikshamitra/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// ApplyEducationLoan handles education loan applications with REVOLUTIONARY member-only restrictions
func (k msgServer) ApplyEducationLoan(goCtx context.Context, msg *types.MsgApplyEducationLoan) (*types.MsgApplyEducationLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// REVOLUTIONARY RESTRICTION: Verify borrower is pool member before processing
	applicantAddr, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid applicant address")
	}

	// Check if liquidity is available and borrower is eligible
	amount, err := sdk.NewDecFromStr(msg.LoanAmount.Amount.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid loan amount")
	}

	canProcess, message := k.liquidityKeeper.CanProcessLoan(ctx, amount, "shikshamitra", applicantAddr)
	if !canProcess {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "🚫 REVOLUTIONARY EDUCATION LENDING RESTRICTION: %s", message)
	}

	// Verify DhanPata IDs
	if !k.accountKeeper.HasDhanPataID(ctx, msg.Applicant, msg.DhanPataID) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDhanPata, "DhanPata ID does not match applicant")
	}

	if !k.accountKeeper.HasDhanPataID(ctx, msg.CoApplicant, msg.CoApplicantDhanPataID) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDhanPata, "Co-applicant DhanPata ID does not match")
	}

	// Parse institution and course types
	institutionType := types.InstitutionType(types.InstitutionType_value[msg.InstitutionType])
	courseType := types.CourseType(types.CourseType_value[msg.CourseType])

	// Check eligibility
	eligible, reason := k.CheckEducationEligibility(ctx, msg, institutionType, courseType)
	if !eligible {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "eligibility check failed: %s", reason)
	}

	// Calculate base interest rate based on institution tier and course
	baseRate := k.CalculateEducationInterestRate(ctx, institutionType, courseType, msg.FamilyIncome)

	// Apply merit discount
	if msg.EntranceExamScore != "" {
		meritScore := k.ParseEntranceScore(msg.EntranceExamScore)
		if meritScore >= 80 {
			meritDiscount := sdk.MustNewDecFromStr(k.GetParams(ctx).MeritDiscount)
			baseRate = baseRate.Sub(meritDiscount)
		}
	}

	// Check for women student discount
	if k.IsWomenStudent(ctx, msg.StudentName) {
		womenDiscount := sdk.MustNewDecFromStr(k.GetParams(ctx).WomenStudentDiscount)
		baseRate = baseRate.Sub(womenDiscount)
	}

	// Check for reserved category discount
	if k.IsReservedCategory(ctx, msg.Applicant) {
		reservedDiscount := sdk.MustNewDecFromStr(k.GetParams(ctx).ReservedCategoryDiscount)
		baseRate = baseRate.Sub(reservedDiscount)
	}

	// Apply festival discount if active
	festivalOffer := k.GetActiveFestivalOffer(ctx)
	if festivalOffer != nil && k.IsEligibleForEducationFestivalOffer(msg, festivalOffer, courseType, institutionType) {
		festivalDiscount := sdk.MustNewDecFromStr(festivalOffer.InterestReduction)
		baseRate = baseRate.Sub(festivalDiscount)
	}

	// Ensure rate is within bounds
	minRate := sdk.MustNewDecFromStr(k.GetParams(ctx).MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}

	// Check if collateral is required
	collateralThreshold := sdk.MustNewDecFromStr(k.GetParams(ctx).CollateralThreshold)
	collateralRequired := msg.LoanAmount.Amount.ToDec().GT(collateralThreshold)

	// Convert loan components
	loanComponents := make([]*types.LoanComponent, len(msg.LoanComponents))
	for i, comp := range msg.LoanComponents {
		loanComponents[i] = &types.LoanComponent{
			ComponentType: comp.ComponentType,
			Amount:        comp.Amount,
			Description:   comp.Description,
		}
	}

	// Create loan application
	loanID := k.GenerateLoanID(ctx)
	loan := types.EducationLoan{
		LoanID:                   loanID,
		Borrower:                 msg.Applicant,
		DhanPataID:               msg.DhanPataID,
		StudentName:              msg.StudentName,
		CoApplicant:              msg.CoApplicant,
		CoApplicantDhanPataID:    msg.CoApplicantDhanPataID,
		InstitutionName:          msg.InstitutionName,
		InstitutionType:          institutionType,
		CourseType:               courseType,
		CourseName:               msg.CourseName,
		CourseDuration:           msg.CourseDuration,
		AcademicYear:             msg.AcademicYear,
		LoanAmount:               msg.LoanAmount,
		ApprovedAmount:           sdk.NewCoin(msg.LoanAmount.Denom, sdk.ZeroInt()),
		LoanComponents:           loanComponents,
		InterestRate:             baseRate.String(),
		Status:                   types.LoanStatus_LOAN_STATUS_PENDING,
		ApplicationDate:          ctx.BlockTime(),
		EntranceExamScore:        msg.EntranceExamScore,
		PreviousAcademicRecord:   msg.PreviousAcademicRecord,
		AdmissionLetter:          msg.AdmissionLetter,
		Pincode:                  msg.Pincode,
		FamilyIncome:             msg.FamilyIncome,
		CollateralOffered:        msg.CollateralOffered,
		CulturalQuote:            msg.CulturalQuote,
	}

	// Store loan application
	k.SetLoan(ctx, loan)

	// Update or create student profile
	k.UpdateStudentProfile(ctx, msg)

	// Calculate estimated EMI
	moratoriumPeriod := k.GetParams(ctx).MoratoriumPeriod
	totalMonths := msg.CourseDuration + int32(moratoriumPeriod) + 120 // 10 years repayment
	estimatedEMI := k.CalculateEMI(ctx, msg.LoanAmount, baseRate, int64(totalMonths))

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEducationLoanApplied,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyBorrower, msg.Applicant),
			sdk.NewAttribute(types.AttributeKeyStudentName, msg.StudentName),
			sdk.NewAttribute(types.AttributeKeyInstitution, msg.InstitutionName),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.LoanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, baseRate.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralRequired, fmt.Sprintf("%t", collateralRequired)),
		),
	})

	return &types.MsgApplyEducationLoanResponse{
		LoanID:             loanID,
		InterestRate:       baseRate.String(),
		Status:             "Application submitted successfully",
		CollateralRequired: collateralRequired,
		MoratoriumPeriod:   fmt.Sprintf("%d months", moratoriumPeriod),
		EstimatedEMI:       estimatedEMI.String(),
	}, nil
}

// UpdateAcademicProgress handles academic progress updates
func (k msgServer) UpdateAcademicProgress(goCtx context.Context, msg *types.MsgUpdateAcademicProgress) (*types.MsgUpdateAcademicProgressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority
	if !k.IsAuthorizedEducationVerifier(ctx, msg.Authority) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to update academic progress")
	}

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Parse GPA
	gpa, err := sdk.NewDecFromStr(msg.GradePointAvg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid GPA format")
	}

	// Parse attendance
	attendance, err := sdk.NewDecFromStr(msg.Attendance)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid attendance format")
	}

	// Create academic progress record
	progress := types.AcademicProgress{
		Semester:           msg.Semester,
		GradePointAvg:      msg.GradePointAvg,
		Attendance:         msg.Attendance,
		Marksheet:          msg.Marksheet,
		Remarks:            msg.Remarks,
		ContinuationStatus: msg.ContinuationStatus,
		UpdateDate:         ctx.BlockTime(),
	}

	// Add to loan's academic progress
	loan.AcademicProgress = append(loan.AcademicProgress, &progress)

	// Check if interest rate should be revised based on performance
	interestRevised := false
	newInterestRate := sdk.MustNewDecFromStr(loan.InterestRate)

	// Merit-based interest rate reduction
	if gpa.GTE(sdk.NewDecWithPrec(90, 1)) && attendance.GTE(sdk.NewDecWithPrec(85, 0)) {
		// Excellent performance: 0.25% reduction
		reduction := sdk.MustNewDecFromStr("0.0025")
		newInterestRate = newInterestRate.Sub(reduction)
		interestRevised = true
	} else if gpa.GTE(sdk.NewDecWithPrec(80, 1)) && attendance.GTE(sdk.NewDecWithPrec(75, 0)) {
		// Good performance: 0.1% reduction
		reduction := sdk.MustNewDecFromStr("0.001")
		newInterestRate = newInterestRate.Sub(reduction)
		interestRevised = true
	}

	// Ensure rate doesn't go below minimum
	minRate := sdk.MustNewDecFromStr(k.GetParams(ctx).MinInterestRate)
	if newInterestRate.LT(minRate) {
		newInterestRate = minRate
	}

	if interestRevised {
		loan.InterestRate = newInterestRate.String()
	}

	// Update loan
	k.SetLoan(ctx, loan)

	// Update student profile
	profile, _ := k.GetStudentProfile(ctx, loan.Borrower)
	k.UpdateStudentCGPA(ctx, profile, gpa)

	// Determine next disbursement status
	nextDisbursementStatus := "Ready for next semester disbursement"
	if msg.ContinuationStatus == "DROPPED" {
		nextDisbursementStatus = "Disbursements suspended - Student dropped"
		loan.Status = types.LoanStatus_LOAN_STATUS_IN_MORATORIUM
		k.SetLoan(ctx, loan)
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcademicProgressUpdated,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeySemester, fmt.Sprintf("%d", msg.Semester)),
			sdk.NewAttribute(types.AttributeKeyGPA, msg.GradePointAvg),
			sdk.NewAttribute(types.AttributeKeyAttendance, msg.Attendance),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.ContinuationStatus),
			sdk.NewAttribute(types.AttributeKeyInterestRevised, fmt.Sprintf("%t", interestRevised)),
		),
	})

	return &types.MsgUpdateAcademicProgressResponse{
		Success:                  true,
		InterestRateRevised:      interestRevised,
		NewInterestRate:          newInterestRate.String(),
		NextDisbursementStatus:   nextDisbursementStatus,
	}, nil
}

// UpdateEmploymentStatus handles employment status updates
func (k msgServer) UpdateEmploymentStatus(goCtx context.Context, msg *types.MsgUpdateEmploymentStatus) (*types.MsgUpdateEmploymentStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Verify student is the borrower
	if loan.Borrower != msg.Student {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the borrower of this loan")
	}

	// Update employment details
	loan.EmploymentDetails = &types.EmploymentDetails{
		EmploymentStatus: msg.EmploymentStatus,
		EmployerName:     msg.EmployerName,
		JobTitle:         msg.JobTitle,
		MonthlySalary:    msg.MonthlySalary,
		JoiningDate:      msg.JoiningDate,
		OfferLetter:      msg.OfferLetter,
		UpdateDate:       &ctx.BlockTime(),
	}

	// Calculate new EMI based on employment status and salary
	repaymentScheduleUpdated := false
	var newEMI sdk.Coin
	var repaymentStartDate time.Time

	if msg.EmploymentStatus == "EMPLOYED" && msg.MonthlySalary.IsPositive() {
		// Calculate affordable EMI (30% of monthly salary)
		affordableEMI := msg.MonthlySalary.Amount.MulRaw(30).QuoRaw(100)
		
		// Calculate remaining loan amount
		remainingAmount := loan.TotalRepayment.Sub(loan.RepaidAmount)
		
		// Calculate new repayment duration based on affordable EMI
		interestRate := sdk.MustNewDecFromStr(loan.InterestRate)
		newDuration := k.CalculateRepaymentDuration(ctx, remainingAmount, affordableEMI, interestRate)
		
		// Set repayment start date (1 month from employment)
		joiningDate, _ := time.Parse("2006-01-02", msg.JoiningDate)
		repaymentStartDate = joiningDate.AddDate(0, 1, 0)
		loan.RepaymentStartDate = &repaymentStartDate
		
		newEMI = sdk.NewCoin(msg.MonthlySalary.Denom, affordableEMI)
		loan.EMIAmount = newEMI
		repaymentScheduleUpdated = true
		
		// Update loan status
		loan.Status = types.LoanStatus_LOAN_STATUS_REPAYMENT
	}

	// Update loan
	k.SetLoan(ctx, loan)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEmploymentStatusUpdated,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeyEmploymentStatus, msg.EmploymentStatus),
			sdk.NewAttribute(types.AttributeKeyEmployer, msg.EmployerName),
			sdk.NewAttribute(types.AttributeKeySalary, msg.MonthlySalary.String()),
			sdk.NewAttribute(types.AttributeKeyRepaymentUpdated, fmt.Sprintf("%t", repaymentScheduleUpdated)),
		),
	})

	return &types.MsgUpdateEmploymentStatusResponse{
		Success:                   true,
		RepaymentScheduleUpdated:  repaymentScheduleUpdated,
		NewEMIAmount:              newEMI.String(),
		RepaymentStartDate:        repaymentStartDate.Format("2006-01-02"),
	}, nil
}

// ApplyScholarship handles scholarship applications
func (k msgServer) ApplyScholarship(goCtx context.Context, msg *types.MsgApplyScholarship) (*types.MsgApplyScholarshipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Verify applicant is the borrower
	if loan.Borrower != msg.Applicant {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the borrower of this loan")
	}

	// Parse academic score
	academicScore, err := sdk.NewDecFromStr(msg.AcademicScore)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid academic score format")
	}

	// Check eligibility based on scholarship type
	eligible := false
	evaluationCriteria := ""

	switch msg.ScholarshipType {
	case "MERIT":
		if academicScore.GTE(sdk.NewDecWithPrec(85, 0)) {
			eligible = true
			evaluationCriteria = "Academic excellence (85%+ required)"
		}
	case "NEED":
		incomeThreshold := sdk.NewCoin("NAMO", sdk.NewInt(500000)) // 5 lakh annual
		if loan.FamilyIncome.IsLT(incomeThreshold) {
			eligible = true
			evaluationCriteria = "Financial need (family income < 5 lakh)"
		}
	case "SPORTS":
		// Check achievements for sports-related keywords
		eligible = k.HasSportsAchievements(msg.Achievements)
		evaluationCriteria = "State/National level sports achievements"
	case "SPECIAL":
		// Special category scholarships
		eligible = true
		evaluationCriteria = "Special category evaluation"
	}

	if !eligible {
		return nil, sdkerrors.Wrap(types.ErrNotEligible, "not eligible for this scholarship type")
	}

	// Generate scholarship ID
	scholarshipID := k.GenerateScholarshipID(ctx)

	// Create scholarship record
	scholarship := types.ScholarshipRecord{
		ScholarshipID:   scholarshipID,
		ScholarshipType: msg.ScholarshipType,
		Amount:          msg.RequestedAmount,
		AcademicYear:    loan.AcademicYear,
		Status:          "APPLIED",
		AwardDate:       &ctx.BlockTime(),
	}

	// Add to loan's scholarship records
	loan.Scholarships = append(loan.Scholarships, &scholarship)
	k.SetLoan(ctx, loan)

	// Calculate expected decision date (15 days from application)
	expectedDecisionDate := ctx.BlockTime().AddDate(0, 0, 15)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeScholarshipApplied,
			sdk.NewAttribute(types.AttributeKeyScholarshipID, scholarshipID),
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeyScholarshipType, msg.ScholarshipType),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.RequestedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyAcademicScore, msg.AcademicScore),
		),
	})

	return &types.MsgApplyScholarshipResponse{
		ScholarshipID:        scholarshipID,
		Status:               "Application submitted successfully",
		EvaluationCriteria:   evaluationCriteria,
		ExpectedDecisionDate: expectedDecisionDate.Format("2006-01-02"),
	}, nil
}

// REVOLUTIONARY STAGED LOAN DISBURSEMENT

// RequestSemesterDisbursement handles semester-by-semester loan disbursement requests
func (k msgServer) RequestSemesterDisbursement(goCtx context.Context, msg *types.MsgRequestSemesterDisbursement) (*types.MsgRequestSemesterDisbursementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// REVOLUTIONARY MEMBER VERIFICATION: Only pool members can request disbursements
	studentAddr, err := sdk.AccAddressFromBech32(msg.Student)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid student address")
	}

	if !k.liquidityKeeper.IsPoolMember(ctx, studentAddr) {
		return nil, sdkerrors.Wrap(types.ErrNotEligible, "🚫 Only Suraksha Pool members can request education loan disbursements")
	}

	// Get loan
	loan, found := k.GetEducationLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "education loan not found")
	}

	// Verify student is the borrower
	if loan.Borrower != msg.Student {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the borrower of this loan")
	}

	// Verify loan is active
	if loan.Status != "ACTIVE" && loan.Status != "DISBURSING" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidLoanStatus, "loan status is %s, expected ACTIVE", loan.Status)
	}

	// Check if this is the next expected semester
	expectedSemester := int32(len(loan.SemesterDisbursements)) + 1
	if msg.Semester != expectedSemester {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "expected semester %d, got %d", expectedSemester, msg.Semester)
	}

	// Check if student has exceeded maximum semesters for course type
	maxSemesters := k.CalculateTotalSemesters(loan.CourseType.String(), loan.CourseDuration)
	if msg.Semester > maxSemesters {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "semester %d exceeds maximum %d for %s course", 
			msg.Semester, maxSemesters, loan.CourseType.String())
	}

	// REVOLUTIONARY VALIDATION: Verify student has paid 20% deposit
	studentDepositRequired := msg.SemesterFee.Amount.MulRaw(20).QuoRaw(100)
	if !msg.StudentDepositPaid || msg.StudentDepositAmount.Amount.LT(studentDepositRequired) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientDeposit, 
			"student must pay 20%% deposit (₹%s) before platform disbursement. Paid: ₹%s", 
			studentDepositRequired.String(), msg.StudentDepositAmount.Amount.String())
	}

	// REVOLUTIONARY VALIDATION: Verify college fee direct payment
	if !k.CheckCollegeFeeDirectPayment(ctx, msg.CollegeAddress, msg.SemesterFee) {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollege, "college not verified or fee structure mismatch")
	}

	// Check platform liquidity for this disbursement
	platformPortion := msg.SemesterFee.Amount.MulRaw(80).QuoRaw(100)
	platformAmount, err := sdk.NewDecFromStr(platformPortion.String())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid platform portion amount")
	}

	canProcess, message := k.liquidityKeeper.CanProcessLoan(ctx, platformAmount, "shikshamitra", studentAddr)
	if !canProcess {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientLiquidity, "🚫 %s", message)
	}

	// Create semester disbursement
	disbursement, err := k.CreateSemesterDisbursement(ctx, msg.LoanID, msg.Semester, msg.SemesterFee, msg.StudentDepositPaid)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create semester disbursement")
	}

	// Add disbursement to loan
	loan.SemesterDisbursements = append(loan.SemesterDisbursements, disbursement)
	
	// Update loan total disbursed amount
	loan.TotalDisbursed = loan.TotalDisbursed.Add(disbursement.DisbursedAmount)
	
	// Update loan status
	if msg.Semester >= maxSemesters {
		loan.Status = "COMPLETED_DISBURSEMENT"
		loan.RepaymentStartDate = &ctx.BlockTime()
	} else {
		loan.Status = "DISBURSING"
	}
	
	k.SetEducationLoan(ctx, loan)

	// REVOLUTIONARY: Direct payment to college (not student) to prevent misuse
	collegeAddr, err := sdk.AccAddressFromBech32(msg.CollegeAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid college address")
	}

	// Transfer funds directly to college
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, collegeAddr, 
		sdk.NewCoins(disbursement.DisbursedAmount))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to transfer funds to college")
	}

	// Update disbursement status
	disbursement.Status = "DISBURSED"
	disbursement.CollegeAddress = msg.CollegeAddress
	
	// Calculate risk reduction vs traditional model
	traditional, staged, riskReduction := k.CalculateMaximumExposureReduction(loan.CourseType.String(), 
		sdk.NewCoin(msg.SemesterFee.Denom, msg.SemesterFee.Amount.MulRaw(int64(maxSemesters))))

	// REVOLUTIONARY TRANSPARENCY: Emit comprehensive event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"semester_disbursement_completed",
			sdk.NewAttribute("loan_id", msg.LoanID),
			sdk.NewAttribute("semester", fmt.Sprintf("%d", msg.Semester)),
			sdk.NewAttribute("student", msg.Student),
			sdk.NewAttribute("college", msg.CollegeAddress),
			sdk.NewAttribute("total_fee", msg.SemesterFee.String()),
			sdk.NewAttribute("platform_portion", disbursement.PlatformPortion.String()),
			sdk.NewAttribute("student_deposit", disbursement.StudentDeposit.String()),
			sdk.NewAttribute("processing_fee", disbursement.ProcessingFee.String()),
			sdk.NewAttribute("disbursed_to_college", disbursement.DisbursedAmount.String()),
			sdk.NewAttribute("traditional_exposure", traditional.TruncateInt().String()),
			sdk.NewAttribute("staged_exposure", staged.TruncateInt().String()),
			sdk.NewAttribute("risk_reduction", fmt.Sprintf("%.1f%%", riskReduction.MustFloat64())),
			sdk.NewAttribute("next_semester_required", fmt.Sprintf("%t", msg.Semester < maxSemesters)),
		),
	})

	return &types.MsgRequestSemesterDisbursementResponse{
		Success:              true,
		DisbursedAmount:      disbursement.DisbursedAmount.String(),
		ProcessingFee:        disbursement.ProcessingFee.String(),
		NextSemesterNumber:   msg.Semester + 1,
		MaxSemesters:         maxSemesters,
		CompletionStatus:     msg.Semester >= maxSemesters,
		RiskReduction:        fmt.Sprintf("%.1f%% safer than traditional loans", riskReduction.MustFloat64()),
		PaymentMethod:        "Direct to college - prevents fund misuse",
		AcademicValidation:   "Previous semester performance verified",
	}, nil
}

// PayStudentDeposit handles student's 20% semester deposit payment
func (k msgServer) PayStudentDeposit(goCtx context.Context, msg *types.MsgPayStudentDeposit) (*types.MsgPayStudentDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get loan
	loan, found := k.GetEducationLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "education loan not found")
	}

	// Verify student is the borrower
	if loan.Borrower != msg.Student {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the borrower of this loan")
	}

	// Calculate required deposit (20% of semester fee)
	requiredDeposit := msg.SemesterFee.Amount.MulRaw(20).QuoRaw(100)
	if msg.DepositAmount.Amount.LT(requiredDeposit) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientDeposit, 
			"deposit amount ₹%s is less than required 20%% (₹%s)", 
			msg.DepositAmount.Amount.String(), requiredDeposit.String())
	}

	// Transfer deposit from student to module
	studentAddr, err := sdk.AccAddressFromBech32(msg.Student)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid student address")
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, studentAddr, types.ModuleName, 
		sdk.NewCoins(msg.DepositAmount))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to collect student deposit")
	}

	// Record deposit payment
	depositRecord := types.StudentDepositRecord{
		LoanID:        msg.LoanID,
		Semester:      msg.Semester,
		SemesterFee:   msg.SemesterFee,
		DepositAmount: msg.DepositAmount,
		PaymentDate:   ctx.BlockTime(),
		Status:        "PAID",
	}

	// Add to loan's deposit records
	loan.StudentDeposits = append(loan.StudentDeposits, &depositRecord)
	k.SetEducationLoan(ctx, loan)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"student_deposit_paid",
			sdk.NewAttribute("loan_id", msg.LoanID),
			sdk.NewAttribute("semester", fmt.Sprintf("%d", msg.Semester)),
			sdk.NewAttribute("deposit_amount", msg.DepositAmount.String()),
			sdk.NewAttribute("required_amount", requiredDeposit.String()),
			sdk.NewAttribute("commitment_shown", "Student has skin in the game"),
		),
	)

	return &types.MsgPayStudentDepositResponse{
		Success:                true,
		DepositPaid:           msg.DepositAmount.String(),
		RequiredDeposit:       requiredDeposit.String(),
		ReadyForDisbursement:  true,
		CommitmentMessage:     "20% deposit shows student commitment to education",
	}, nil
}