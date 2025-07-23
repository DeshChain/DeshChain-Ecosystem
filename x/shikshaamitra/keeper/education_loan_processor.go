package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/shikshaamitra/types"
)

// EducationLoanProcessor handles comprehensive education loan processing
type EducationLoanProcessor struct {
	keeper           Keeper
	creditAnalyzer   *StudentCreditAnalyzer
	institutionMgr   *InstitutionManager
}

// NewEducationLoanProcessor creates a new education loan processor
func NewEducationLoanProcessor(keeper Keeper) *EducationLoanProcessor {
	return &EducationLoanProcessor{
		keeper:         keeper,
		creditAnalyzer: NewStudentCreditAnalyzer(keeper),
		institutionMgr: NewInstitutionManager(keeper),
	}
}

// ProcessedEducationLoan represents a comprehensive education loan
type ProcessedEducationLoan struct {
	LoanID                string                          `json:"loan_id"`
	ApplicationID         string                          `json:"application_id"`
	StudentID             string                          `json:"student_id"`
	CoApplicantID         string                          `json:"co_applicant_id,omitempty"`
	LoanType              string                          `json:"loan_type"` // UNDERGRADUATE, POSTGRADUATE, PROFESSIONAL, ABROAD
	CourseDetails         types.CourseInfo                `json:"course_details"`
	InstitutionDetails    types.InstitutionInfo           `json:"institution_details"`
	LoanAmount            sdk.Coin                        `json:"loan_amount"`
	TuitionFees           sdk.Coin                        `json:"tuition_fees"`
	LivingExpenses        sdk.Coin                        `json:"living_expenses"`
	InterestRate          sdk.Dec                         `json:"interest_rate"`
	LoanTenure            int64                           `json:"loan_tenure"` // months
	MoratoriumPeriod      int64                           `json:"moratorium_period"` // months
	RepaymentType         string                          `json:"repayment_type"` // INCOME_DRIVEN, STANDARD, GRADUATED
	RepaymentSchedule     []types.EducationRepayment      `json:"repayment_schedule"`
	IncomeThresholds      types.IncomeThresholds          `json:"income_thresholds"`
	ColateralRequired     bool                            `json:"collateral_required"`
	CollateralDetails     types.EducationCollateral       `json:"collateral_details,omitempty"`
	GuarantorRequired     bool                            `json:"guarantor_required"`
	GuarantorDetails      types.GuarantorInfo             `json:"guarantor_details,omitempty"`
	Scholarships          []types.ScholarshipInfo         `json:"scholarships"`
	DisbursementSchedule  []types.EducationDisbursement   `json:"disbursement_schedule"`
	AcademicProgress      types.AcademicProgress          `json:"academic_progress"`
	EmploymentTracking    types.EmploymentTracking        `json:"employment_tracking"`
	Status                types.EducationLoanStatus       `json:"status"`
	StudentProfile        *StudentCreditProfile           `json:"student_profile"`
	RiskAssessment        types.StudentRiskAssessment     `json:"risk_assessment"`
	ApprovalWorkflow      types.EducationApprovalWorkflow `json:"approval_workflow"`
	ComplianceChecks      types.EducationComplianceChecks `json:"compliance_checks"`
	CreatedAt             time.Time                       `json:"created_at"`
	UpdatedAt             time.Time                       `json:"updated_at"`
	ExpectedGraduation    time.Time                       `json:"expected_graduation"`
}

// ProcessEducationLoanApplication processes a comprehensive education loan application
func (elp *EducationLoanProcessor) ProcessEducationLoanApplication(ctx sdk.Context, applicationID string) (*ProcessedEducationLoan, error) {
	// Get loan application
	application, found := elp.keeper.GetEducationLoanApplication(ctx, applicationID)
	if !found {
		return nil, fmt.Errorf("education loan application not found: %s", applicationID)
	}

	// Validate application status
	if application.Status != types.ApplicationStatusSubmitted {
		return nil, fmt.Errorf("application not in submitted status: %s", application.Status)
	}

	// Perform comprehensive student credit analysis
	studentProfile, err := elp.creditAnalyzer.AnalyzeStudentCredit(ctx, application.StudentID, application.CoApplicantID)
	if err != nil {
		return nil, fmt.Errorf("student credit analysis failed: %w", err)
	}

	// Verify institution partnership and accreditation
	institutionStatus, err := elp.institutionMgr.VerifyInstitution(ctx, application.InstitutionID)
	if err != nil {
		return nil, fmt.Errorf("institution verification failed: %w", err)
	}

	// Update application status
	application.Status = types.ApplicationStatusUnderReview
	application.CreditScore = studentProfile.CreditScore
	application.RiskCategory = studentProfile.RiskCategory
	elp.keeper.SetEducationLoanApplication(ctx, application)

	// Perform compliance checks
	complianceChecks, err := elp.performEducationComplianceChecks(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("compliance checks failed: %w", err)
	}

	// Check eligibility based on comprehensive analysis
	eligible, eligibilityReason := elp.checkEducationEligibility(ctx, application, studentProfile, institutionStatus)
	if !eligible {
		return elp.rejectEducationApplication(ctx, application, eligibilityReason)
	}

	// Determine loan parameters
	loanAmount := elp.calculateOptimalLoanAmount(ctx, application, studentProfile)
	interestRate := elp.calculateEducationInterestRate(ctx, application, studentProfile, institutionStatus)

	// Apply education-specific discounts and scholarships
	interestRate, scholarships := elp.applyEducationDiscounts(ctx, application.StudentID, interestRate, application.CourseType)

	// Determine loan tenure and moratorium
	tenure, moratorium := elp.calculateEducationTenure(ctx, application.CourseType, application.CourseDuration)

	// Determine repayment type based on course and income potential
	repaymentType := elp.determineRepaymentType(ctx, application.CourseType, studentProfile)

	// Create income thresholds for income-driven repayment
	incomeThresholds := elp.createIncomeThresholds(ctx, application.CourseType, loanAmount)

	// Generate repayment schedule
	repaymentSchedule := elp.generateEducationRepaymentSchedule(
		ctx, loanAmount, interestRate, tenure, moratorium, repaymentType, incomeThresholds,
	)

	// Create disbursement schedule tied to academic progress
	disbursementSchedule := elp.createDisbursementSchedule(ctx, loanAmount, application.CourseDuration)

	// Assess collateral and guarantor requirements
	collateralRequired, collateralDetails := elp.assessEducationCollateral(ctx, loanAmount, studentProfile)
	guarantorRequired, guarantorDetails := elp.assessGuarantorRequirements(ctx, loanAmount, studentProfile, application.CoApplicantID)

	// Generate loan ID
	loanID := elp.generateEducationLoanID(ctx, application.StudentID)

	// Create education loan
	educationLoan := &ProcessedEducationLoan{
		LoanID:            loanID,
		ApplicationID:     applicationID,
		StudentID:         application.StudentID,
		CoApplicantID:     application.CoApplicantID,
		LoanType:          application.LoanType,
		CourseDetails: types.CourseInfo{
			CourseName:     application.CourseName,
			CourseType:     application.CourseType,
			CourseDuration: application.CourseDuration,
			Specialization: application.Specialization,
			CourseLevel:    application.CourseLevel,
		},
		InstitutionDetails:   institutionStatus.InstitutionInfo,
		LoanAmount:          loanAmount,
		TuitionFees:         application.TuitionFees,
		LivingExpenses:      application.LivingExpenses,
		InterestRate:        interestRate,
		LoanTenure:          tenure,
		MoratoriumPeriod:    moratorium,
		RepaymentType:       repaymentType,
		RepaymentSchedule:   repaymentSchedule,
		IncomeThresholds:    incomeThresholds,
		ColateralRequired:   collateralRequired,
		CollateralDetails:   collateralDetails,
		GuarantorRequired:   guarantorRequired,
		GuarantorDetails:    guarantorDetails,
		Scholarships:        scholarships,
		DisbursementSchedule: disbursementSchedule,
		AcademicProgress: types.AcademicProgress{
			StudentID:           application.StudentID,
			CurrentSemester:     1,
			CGPA:                sdk.ZeroDec(),
			AttendanceRate:      sdk.ZeroDec(),
			AcademicStanding:    "GOOD",
			LastUpdated:         ctx.BlockTime(),
		},
		EmploymentTracking: types.EmploymentTracking{
			StudentID:           application.StudentID,
			EmploymentStatus:    "STUDENT",
			CurrentIncome:       sdk.ZeroDec(),
			EmploymentSector:    "",
			LastUpdated:         ctx.BlockTime(),
		},
		Status:         types.EducationLoanStatusApproved,
		StudentProfile: studentProfile,
		RiskAssessment: elp.createStudentRiskAssessment(ctx, studentProfile, institutionStatus),
		ApprovalWorkflow: types.EducationApprovalWorkflow{
			ApprovedBy:        "AUTOMATED_SYSTEM",
			ApprovedAmount:    loanAmount,
			ApprovedRate:      interestRate,
			ApprovalDate:      ctx.BlockTime(),
			ApprovalLevel:     elp.determineEducationApprovalLevel(ctx, loanAmount),
			Conditions:        elp.generateEducationApprovalConditions(ctx, collateralRequired, guarantorRequired),
			InstitutionPartnership: institutionStatus.PartnershipLevel,
		},
		ComplianceChecks:   complianceChecks,
		CreatedAt:          ctx.BlockTime(),
		UpdatedAt:          ctx.BlockTime(),
		ExpectedGraduation: ctx.BlockTime().AddDate(int(application.CourseDuration/12), int(application.CourseDuration%12), 0),
	}

	// Store education loan
	elp.keeper.SetProcessedEducationLoan(ctx, *educationLoan)

	// Update application status
	application.Status = types.ApplicationStatusApproved
	application.ApprovedLoanID = loanID
	elp.keeper.SetEducationLoanApplication(ctx, application)

	// Create academic progress tracking
	elp.initializeAcademicTracking(ctx, educationLoan)

	// Emit loan approval event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEducationLoanApproved,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyStudentID, application.StudentID),
			sdk.NewAttribute(types.AttributeKeyLoanAmount, loanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, interestRate.String()),
			sdk.NewAttribute(types.AttributeKeyCourseType, application.CourseType),
			sdk.NewAttribute(types.AttributeKeyInstitution, application.InstitutionID),
			sdk.NewAttribute(types.AttributeKeyRepaymentType, repaymentType),
		),
	)

	return educationLoan, nil
}

// DisburseEducationLoan handles education loan disbursement with academic progress verification
func (elp *EducationLoanProcessor) DisburseEducationLoan(ctx sdk.Context, loanID string, disbursementRequest types.DisbursementRequest) error {
	// Get loan details
	loan, found := elp.keeper.GetProcessedEducationLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("education loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.EducationLoanStatusApproved && loan.Status != types.EducationLoanStatusActive {
		return fmt.Errorf("loan not in disbursable status: %s", loan.Status)
	}

	// Verify academic progress for subsequent disbursements
	if disbursementRequest.DisbursementNumber > 1 {
		progressValid, err := elp.verifyAcademicProgress(ctx, loan.StudentID, disbursementRequest.RequiredSemester)
		if err != nil || !progressValid {
			return fmt.Errorf("academic progress verification failed for disbursement")
		}
	}

	// Check disbursement conditions
	if loan.ColateralRequired && !elp.isEducationCollateralVerified(ctx, loanID) {
		return fmt.Errorf("collateral verification pending for loan: %s", loanID)
	}

	if loan.GuarantorRequired && !elp.isGuarantorVerified(ctx, loanID) {
		return fmt.Errorf("guarantor verification pending for loan: %s", loanID)
	}

	// Find the disbursement in schedule
	var targetDisbursement *types.EducationDisbursement
	for i, disburse := range loan.DisbursementSchedule {
		if disburse.DisbursementNumber == disbursementRequest.DisbursementNumber {
			targetDisbursement = &loan.DisbursementSchedule[i]
			break
		}
	}

	if targetDisbursement == nil {
		return fmt.Errorf("disbursement number %d not found in schedule", disbursementRequest.DisbursementNumber)
	}

	// Get student address
	studentAddr, err := sdk.AccAddressFromBech32(loan.StudentID)
	if err != nil {
		return fmt.Errorf("invalid student address: %s", loan.StudentID)
	}

	// Calculate disbursement amount and fees
	disbursementAmount := targetDisbursement.Amount
	processingFee := elp.calculateEducationProcessingFee(ctx, disbursementAmount)
	netDisbursement := disbursementAmount.Sub(processingFee)

	// Handle different disbursement modes
	switch targetDisbursement.DisbursementMode {
	case "DIRECT_TO_INSTITUTION":
		err = elp.processInstitutionPayment(ctx, loan.InstitutionDetails.InstitutionID, netDisbursement)
	case "DIRECT_TO_STUDENT":
		err = elp.processDirectStudentTransfer(ctx, studentAddr, netDisbursement)
	case "SPLIT_PAYMENT":
		err = elp.processSplitPayment(ctx, loan, netDisbursement, targetDisbursement)
	default:
		err = elp.processDirectStudentTransfer(ctx, studentAddr, netDisbursement)
	}

	if err != nil {
		return fmt.Errorf("failed to disburse education loan: %w", err)
	}

	// Update disbursement status
	targetDisbursement.Status = types.DisbursementStatusCompleted
	targetDisbursement.DisbursedAmount = netDisbursement
	targetDisbursement.ProcessingFee = processingFee
	targetDisbursement.DisbursedDate = ctx.BlockTime()

	// Update loan status
	if loan.Status == types.EducationLoanStatusApproved {
		loan.Status = types.EducationLoanStatusActive
	}
	loan.UpdatedAt = ctx.BlockTime()

	elp.keeper.SetProcessedEducationLoan(ctx, loan)

	// Update student profile
	elp.updateStudentProfileAfterDisbursement(ctx, loan.StudentID, disbursementAmount)

	// Emit disbursement event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEducationLoanDisbursed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyDisbursedAmount, netDisbursement.String()),
			sdk.NewAttribute(types.AttributeKeyDisbursementNumber, fmt.Sprintf("%d", disbursementRequest.DisbursementNumber)),
			sdk.NewAttribute(types.AttributeKeyDisbursementMode, targetDisbursement.DisbursementMode),
		),
	)

	return nil
}

// ProcessIncomeBasedRepayment handles income-driven repayment calculations
func (elp *EducationLoanProcessor) ProcessIncomeBasedRepayment(ctx sdk.Context, loanID string, incomeVerification types.IncomeVerification) error {
	// Get loan details
	loan, found := elp.keeper.GetProcessedEducationLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("education loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.EducationLoanStatusActive {
		return fmt.Errorf("loan not in active status: %s", loan.Status)
	}

	// Update employment tracking
	loan.EmploymentTracking.CurrentIncome = incomeVerification.MonthlyIncome
	loan.EmploymentTracking.EmploymentStatus = incomeVerification.EmploymentStatus
	loan.EmploymentTracking.EmploymentSector = incomeVerification.EmploymentSector
	loan.EmploymentTracking.LastUpdated = ctx.BlockTime()

	// Calculate income-based repayment amount
	newRepaymentAmount := elp.calculateIncomeBasedRepayment(ctx, loan, incomeVerification.MonthlyIncome)

	// Update repayment schedule if income-driven
	if loan.RepaymentType == "INCOME_DRIVEN" {
		err := elp.updateIncomeBasedSchedule(ctx, &loan, newRepaymentAmount)
		if err != nil {
			return fmt.Errorf("failed to update income-based schedule: %w", err)
		}
	}

	// Check for payment caps and forgiveness eligibility
	elp.checkPaymentCapsAndForgiveness(ctx, &loan)

	elp.keeper.SetProcessedEducationLoan(ctx, loan)

	// Emit income update event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncomeUpdated,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyMonthlyIncome, incomeVerification.MonthlyIncome.String()),
			sdk.NewAttribute(types.AttributeKeyNewRepaymentAmount, newRepaymentAmount.String()),
		),
	)

	return nil
}

// ProcessEducationRepayment handles education loan repayment with income consideration
func (elp *EducationLoanProcessor) ProcessEducationRepayment(ctx sdk.Context, loanID string, repaymentAmount sdk.Coin) error {
	// Get loan details
	loan, found := elp.keeper.GetProcessedEducationLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("education loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.EducationLoanStatusActive {
		return fmt.Errorf("loan not in active status for repayment: %s", loan.Status)
	}

	// Check if still in moratorium period
	if elp.isInMoratoriumPeriod(ctx, loan) {
		return fmt.Errorf("loan is still in moratorium period")
	}

	// Get student address
	studentAddr, err := sdk.AccAddressFromBech32(loan.StudentID)
	if err != nil {
		return fmt.Errorf("invalid student address: %s", loan.StudentID)
	}

	// Calculate outstanding amount
	outstandingAmount := elp.calculateEducationOutstandingAmount(ctx, loan)
	if repaymentAmount.Amount.GT(outstandingAmount.Amount) {
		repaymentAmount = outstandingAmount // Cap at outstanding amount
	}

	// Process payment from student to module
	err = elp.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		studentAddr,
		types.ModuleName,
		sdk.NewCoins(repaymentAmount),
	)
	if err != nil {
		return fmt.Errorf("failed to process education loan repayment: %w", err)
	}

	// Allocate repayment (principal vs interest)
	principalPortion, interestPortion := elp.allocateEducationRepayment(ctx, loan, repaymentAmount)

	// Create repayment record
	repaymentRecord := types.EducationRepayment{
		RepaymentID:      elp.generateEducationRepaymentID(ctx, loanID),
		LoanID:           loanID,
		RepaymentDate:    ctx.BlockTime(),
		TotalAmount:      repaymentAmount,
		PrincipalAmount:  principalPortion,
		InterestAmount:   interestPortion,
		PaymentMethod:    "BLOCKCHAIN_TRANSFER",
		Status:           types.RepaymentStatusCompleted,
		OutstandingAfter: outstandingAmount.Sub(repaymentAmount),
	}

	// Update loan with repayment
	loan.RepaymentSchedule = append(loan.RepaymentSchedule, repaymentRecord)
	loan.UpdatedAt = ctx.BlockTime()

	// Check if loan is fully repaid
	newOutstanding := outstandingAmount.Sub(repaymentAmount)
	if newOutstanding.IsZero() {
		loan.Status = types.EducationLoanStatusClosed
		
		// Process loan closure and update employment tracking
		err = elp.processEducationLoanClosure(ctx, loan)
		if err != nil {
			elp.keeper.Logger(ctx).Error("Failed to process education loan closure", 
				"loan_id", loanID, "error", err)
		}
	}

	elp.keeper.SetProcessedEducationLoan(ctx, loan)

	// Update student profile
	elp.updateStudentProfileAfterRepayment(ctx, loan.StudentID, repaymentRecord)

	// Emit repayment event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEducationRepaymentProcessed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyRepaymentAmount, repaymentAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPrincipalAmount, principalPortion.String()),
			sdk.NewAttribute(types.AttributeKeyInterestAmount, interestPortion.String()),
			sdk.NewAttribute(types.AttributeKeyOutstandingAmount, newOutstanding.String()),
		),
	)

	return nil
}

// Helper functions for education loan processing

func (elp *EducationLoanProcessor) checkEducationEligibility(ctx sdk.Context, application types.EducationLoanApplication, studentProfile *StudentCreditProfile, institutionStatus types.InstitutionStatus) (bool, string) {
	// Check minimum academic qualifications
	if studentProfile.AcademicProfile.PreviousEducationScore < 60 {
		return false, "Minimum 60% marks required in previous education"
	}

	// Check institution accreditation
	if !institutionStatus.IsAccredited {
		return false, "Institution not accredited for education loans"
	}

	// Check age limits
	if studentProfile.PersonalInfo.Age < 18 || studentProfile.PersonalInfo.Age > 35 {
		return false, "Age must be between 18-35 years for education loans"
	}

	// Check co-applicant income for large loans
	if application.LoanAmount.Amount.GT(sdk.NewInt(1000000)) && studentProfile.FamilyProfile.MonthlyIncome.LT(sdk.NewDec(50000)) {
		return false, "Insufficient family income for requested loan amount"
	}

	// Check for previous defaults
	if len(studentProfile.CreditHistory.DefaultedLoans) > 0 {
		return false, "Previous loan defaults found"
	}

	return true, ""
}

func (elp *EducationLoanProcessor) generateEducationLoanID(ctx sdk.Context, studentID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("SM-%s-%d", studentID[:8], timestamp)
}

// Additional helper functions will be implemented in separate files
// - calculateOptimalLoanAmount
// - calculateEducationInterestRate
// - applyEducationDiscounts
// - calculateEducationTenure
// - determineRepaymentType
// - createIncomeThresholds
// - generateEducationRepaymentSchedule
// - createDisbursementSchedule
// - And other education-specific operations