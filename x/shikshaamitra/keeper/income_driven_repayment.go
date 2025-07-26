package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/shikshaamitra/types"
)

// IncomeBasedRepaymentProcessor handles income-driven repayment calculations
type IncomeBasedRepaymentProcessor struct {
	keeper Keeper
}

// NewIncomeBasedRepaymentProcessor creates a new income-based repayment processor
func NewIncomeBasedRepaymentProcessor(keeper Keeper) *IncomeBasedRepaymentProcessor {
	return &IncomeBasedRepaymentProcessor{
		keeper: keeper,
	}
}

// IncomeBasedRepaymentPlan represents a comprehensive income-driven repayment plan
type IncomeBasedRepaymentPlan struct {
	PlanID                string                     `json:"plan_id"`
	LoanID                string                     `json:"loan_id"`
	StudentID             string                     `json:"student_id"`
	PlanType              string                     `json:"plan_type"` // IBR, PAYE, REPAYE, ICR
	IncomeThresholds      types.IncomeThresholds     `json:"income_thresholds"`
	PaymentCapPercentage  sdk.Dec                    `json:"payment_cap_percentage"`
	ForgivenessEligible   bool                       `json:"forgiveness_eligible"`
	ForgivenessYears      int64                      `json:"forgiveness_years"`
	CurrentPaymentAmount  sdk.Coin                   `json:"current_payment_amount"`
	AnnualIncomeVerification types.IncomeVerification `json:"annual_income_verification"`
	PaymentHistory        []types.IncomeBasedPayment `json:"payment_history"`
	PaymentAdjustments    []types.PaymentAdjustment  `json:"payment_adjustments"`
	InterestCapitalization types.InterestCapitalization `json:"interest_capitalization"`
	SubsidyEligibility    types.SubsidyEligibility   `json:"subsidy_eligibility"`
	HardshipProvisions    types.HardshipProvisions   `json:"hardship_provisions"`
	Status                string                     `json:"status"`
	CreatedAt             time.Time                  `json:"created_at"`
	UpdatedAt             time.Time                  `json:"updated_at"`
	NextReviewDate        time.Time                  `json:"next_review_date"`
}

// CreateIncomeBasedRepaymentPlan creates a new income-driven repayment plan
func (idrp *IncomeBasedRepaymentProcessor) CreateIncomeBasedRepaymentPlan(ctx sdk.Context, loanID string, planType string, initialIncome sdk.Dec) (*IncomeBasedRepaymentPlan, error) {
	// Get loan details
	loan, found := idrp.keeper.GetProcessedEducationLoan(ctx, loanID)
	if !found {
		return nil, fmt.Errorf("education loan not found: %s", loanID)
	}

	// Validate loan eligibility for income-driven repayment
	if !idrp.isEligibleForIncomeBasedRepayment(ctx, loan) {
		return nil, fmt.Errorf("loan not eligible for income-driven repayment")
	}

	// Create income thresholds based on plan type
	incomeThresholds := idrp.createIncomeThresholdsForPlan(ctx, planType, loan.LoanAmount)

	// Calculate initial payment amount
	initialPayment := idrp.calculateInitialIncomeBasedPayment(ctx, loan, planType, initialIncome)

	// Determine forgiveness eligibility and timeline
	forgivenessEligible, forgivenessYears := idrp.determineForgiveness(ctx, loan, planType)

	// Generate plan ID
	planID := idrp.generateRepaymentPlanID(ctx, loanID)

	// Create repayment plan
	plan := &IncomeBasedRepaymentPlan{
		PlanID:               planID,
		LoanID:               loanID,
		StudentID:            loan.StudentID,
		PlanType:             planType,
		IncomeThresholds:     incomeThresholds,
		PaymentCapPercentage: idrp.getPaymentCapPercentage(planType),
		ForgivenessEligible:  forgivenessEligible,
		ForgivenessYears:     forgivenessYears,
		CurrentPaymentAmount: initialPayment,
		AnnualIncomeVerification: types.IncomeVerification{
			StudentID:        loan.StudentID,
			MonthlyIncome:    initialIncome,
			VerificationDate: ctx.BlockTime(),
			VerificationMethod: "INITIAL_APPLICATION",
			EmploymentStatus: "EMPLOYED",
			IncomeSource:     "SALARY",
		},
		PaymentHistory:       []types.IncomeBasedPayment{},
		PaymentAdjustments:   []types.PaymentAdjustment{},
		InterestCapitalization: idrp.createInterestCapitalizationRules(planType),
		SubsidyEligibility:   idrp.assessSubsidyEligibility(ctx, loan, initialIncome),
		HardshipProvisions:   idrp.createHardshipProvisions(planType),
		Status:               "ACTIVE",
		CreatedAt:            ctx.BlockTime(),
		UpdatedAt:            ctx.BlockTime(),
		NextReviewDate:       ctx.BlockTime().AddDate(1, 0, 0), // Annual review
	}

	// Store repayment plan
	idrp.keeper.SetIncomeBasedRepaymentPlan(ctx, *plan)

	// Update loan with income-based repayment plan
	loan.RepaymentType = "INCOME_DRIVEN"
	loan.IncomeThresholds = incomeThresholds
	idrp.keeper.SetProcessedEducationLoan(ctx, loan)

	// Emit plan creation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncomeBasedPlanCreated,
			sdk.NewAttribute(types.AttributeKeyPlanID, planID),
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyPlanType, planType),
			sdk.NewAttribute(types.AttributeKeyInitialPayment, initialPayment.String()),
		),
	)

	return plan, nil
}

// UpdateIncomeBasedPayment updates payment amount based on new income verification
func (idrp *IncomeBasedRepaymentProcessor) UpdateIncomeBasedPayment(ctx sdk.Context, planID string, incomeVerification types.IncomeVerification) error {
	// Get repayment plan
	plan, found := idrp.keeper.GetIncomeBasedRepaymentPlan(ctx, planID)
	if !found {
		return fmt.Errorf("income-based repayment plan not found: %s", planID)
	}

	// Validate income verification
	if err := idrp.validateIncomeVerification(ctx, incomeVerification); err != nil {
		return fmt.Errorf("income verification failed: %w", err)
	}

	// Calculate new payment amount
	newPaymentAmount := idrp.calculateIncomeBasedPayment(ctx, plan, incomeVerification.MonthlyIncome)

	// Check for significant payment changes
	oldPayment := plan.CurrentPaymentAmount
	paymentChange := newPaymentAmount.Amount.ToDec().Sub(oldPayment.Amount.ToDec()).Quo(oldPayment.Amount.ToDec())

	// Create payment adjustment record
	adjustment := types.PaymentAdjustment{
		AdjustmentID:      idrp.generateAdjustmentID(ctx, planID),
		PlanID:            planID,
		PreviousPayment:   oldPayment,
		NewPayment:        newPaymentAmount,
		PaymentChange:     paymentChange,
		IncomeChange:      incomeVerification.MonthlyIncome.Sub(plan.AnnualIncomeVerification.MonthlyIncome),
		EffectiveDate:     ctx.BlockTime(),
		Reason:            "ANNUAL_INCOME_VERIFICATION",
		AdjustmentType:    idrp.determineAdjustmentType(paymentChange),
	}

	// Update plan
	plan.CurrentPaymentAmount = newPaymentAmount
	plan.AnnualIncomeVerification = incomeVerification
	plan.PaymentAdjustments = append(plan.PaymentAdjustments, adjustment)
	plan.UpdatedAt = ctx.BlockTime()
	plan.NextReviewDate = ctx.BlockTime().AddDate(1, 0, 0)

	// Handle significant payment reductions (potential hardship)
	if paymentChange.LT(sdk.NewDecWithPrec(-3, 1)) { // > 30% reduction
		idrp.evaluateHardshipEligibility(ctx, &plan, incomeVerification)
	}

	// Update subsidy eligibility
	plan.SubsidyEligibility = idrp.assessSubsidyEligibility(ctx, types.ProcessedEducationLoan{
		LoanID:    plan.LoanID,
		StudentID: plan.StudentID,
	}, incomeVerification.MonthlyIncome)

	idrp.keeper.SetIncomeBasedRepaymentPlan(ctx, plan)

	// Emit payment update event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncomeBasedPaymentUpdated,
			sdk.NewAttribute(types.AttributeKeyPlanID, planID),
			sdk.NewAttribute(types.AttributeKeyOldPayment, oldPayment.String()),
			sdk.NewAttribute(types.AttributeKeyNewPayment, newPaymentAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPaymentChange, paymentChange.String()),
		),
	)

	return nil
}

// ProcessIncomeBasedPayment processes a payment under income-driven plan
func (idrp *IncomeBasedRepaymentProcessor) ProcessIncomeBasedPayment(ctx sdk.Context, planID string, paymentAmount sdk.Coin) error {
	// Get repayment plan
	plan, found := idrp.keeper.GetIncomeBasedRepaymentPlan(ctx, planID)
	if !found {
		return fmt.Errorf("income-based repayment plan not found: %s", planID)
	}

	// Get associated loan
	loan, found := idrp.keeper.GetProcessedEducationLoan(ctx, plan.LoanID)
	if !found {
		return fmt.Errorf("associated loan not found: %s", plan.LoanID)
	}

	// Validate payment amount
	if paymentAmount.Amount.GT(plan.CurrentPaymentAmount.Amount.Mul(sdk.NewInt(2))) {
		return fmt.Errorf("payment amount exceeds maximum allowed (2x current payment)")
	}

	// Process payment allocation with income-based rules
	principalPortion, interestPortion, subsidyPortion := idrp.allocateIncomeBasedPayment(ctx, plan, loan, paymentAmount)

	// Handle interest subsidy if eligible
	var subsidyAmount sdk.Coin
	if plan.SubsidyEligibility.IsEligible {
		subsidyAmount = idrp.calculateInterestSubsidy(ctx, plan, loan)
		if subsidyAmount.IsPositive() {
			err := idrp.processInterestSubsidy(ctx, plan.StudentID, subsidyAmount)
			if err != nil {
				idrp.keeper.Logger(ctx).Error("Failed to process interest subsidy", "error", err)
			}
		}
	}

	// Create payment record
	paymentRecord := types.IncomeBasedPayment{
		PaymentID:       idrp.generateIncomePaymentID(ctx, planID),
		PlanID:          planID,
		LoanID:          plan.LoanID,
		PaymentDate:     ctx.BlockTime(),
		TotalAmount:     paymentAmount,
		PrincipalAmount: principalPortion,
		InterestAmount:  interestPortion,
		SubsidyAmount:   subsidyAmount,
		PaymentMethod:   "INCOME_BASED",
		IncomeAtPayment: plan.AnnualIncomeVerification.MonthlyIncome,
		Status:          "COMPLETED",
	}

	// Update plan payment history
	plan.PaymentHistory = append(plan.PaymentHistory, paymentRecord)
	plan.UpdatedAt = ctx.BlockTime()

	// Check for interest capitalization events
	idrp.checkInterestCapitalization(ctx, &plan, loan)

	// Check forgiveness progress
	idrp.updateForgivenessProgress(ctx, &plan)

	idrp.keeper.SetIncomeBasedRepaymentPlan(ctx, plan)

	// Emit payment processed event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncomeBasedPaymentProcessed,
			sdk.NewAttribute(types.AttributeKeyPlanID, planID),
			sdk.NewAttribute(types.AttributeKeyPaymentAmount, paymentAmount.String()),
			sdk.NewAttribute(types.AttributeKeySubsidyAmount, subsidyAmount.String()),
		),
	)

	return nil
}

// EvaluateLoanForgiveness evaluates eligibility for loan forgiveness
func (idrp *IncomeBasedRepaymentProcessor) EvaluateLoanForgiveness(ctx sdk.Context, planID string) error {
	// Get repayment plan
	plan, found := idrp.keeper.GetIncomeBasedRepaymentPlan(ctx, planID)
	if !found {
		return fmt.Errorf("income-based repayment plan not found: %s", planID)
	}

	// Check forgiveness eligibility
	if !plan.ForgivenessEligible {
		return fmt.Errorf("loan not eligible for forgiveness under current plan")
	}

	// Calculate qualifying payments
	qualifyingPayments := idrp.countQualifyingPayments(plan.PaymentHistory)

	// Calculate required payments for forgiveness
	requiredPayments := plan.ForgivenessYears * 12 // Monthly payments

	// Check if forgiveness criteria met
	if qualifyingPayments < requiredPayments {
		return fmt.Errorf("insufficient qualifying payments: %d/%d", qualifyingPayments, requiredPayments)
	}

	// Get loan for forgiveness calculation
	loan, found := idrp.keeper.GetProcessedEducationLoan(ctx, plan.LoanID)
	if !found {
		return fmt.Errorf("loan not found for forgiveness: %s", plan.LoanID)
	}

	// Calculate forgiveness amount
	outstandingAmount := idrp.calculateOutstandingAmount(ctx, loan)
	forgivenessAmount := outstandingAmount

	// Process loan forgiveness
	err := idrp.processLoanForgiveness(ctx, plan, forgivenessAmount)
	if err != nil {
		return fmt.Errorf("failed to process loan forgiveness: %w", err)
	}

	// Update loan status
	loan.Status = types.EducationLoanStatusForgiven
	loan.UpdatedAt = ctx.BlockTime()
	idrp.keeper.SetProcessedEducationLoan(ctx, loan)

	// Update plan status
	plan.Status = "FORGIVEN"
	plan.UpdatedAt = ctx.BlockTime()
	idrp.keeper.SetIncomeBasedRepaymentPlan(ctx, plan)

	// Emit forgiveness event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLoanForgiven,
			sdk.NewAttribute(types.AttributeKeyLoanID, plan.LoanID),
			sdk.NewAttribute(types.AttributeKeyPlanID, planID),
			sdk.NewAttribute(types.AttributeKeyForgivenessAmount, forgivenessAmount.String()),
			sdk.NewAttribute(types.AttributeKeyQualifyingPayments, fmt.Sprintf("%d", qualifyingPayments)),
		),
	)

	return nil
}

// Helper functions for income-based repayment

func (idrp *IncomeBasedRepaymentProcessor) createIncomeThresholdsForPlan(ctx sdk.Context, planType string, loanAmount sdk.Coin) types.IncomeThresholds {
	params := idrp.keeper.GetParams(ctx)
	
	thresholds := types.IncomeThresholds{
		PlanType: planType,
	}

	switch planType {
	case "IBR": // Income-Based Repayment
		thresholds.PovertyGuideline = params.PovertyGuideline
		thresholds.PaymentCapPercentage = sdk.NewDecWithPrec(10, 2) // 10% of discretionary income
		thresholds.DiscretionaryIncomeThreshold = sdk.NewDecWithPrec(15, 1) // 150% of poverty guideline
		thresholds.MaxPaymentCap = loanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(12, 2)) // 12% of original loan
		
	case "PAYE": // Pay As You Earn
		thresholds.PovertyGuideline = params.PovertyGuideline
		thresholds.PaymentCapPercentage = sdk.NewDecWithPrec(10, 2) // 10% of discretionary income
		thresholds.DiscretionaryIncomeThreshold = sdk.NewDecWithPrec(15, 1) // 150% of poverty guideline
		thresholds.MaxPaymentCap = loanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(10, 2)) // 10% of original loan
		
	case "REPAYE": // Revised Pay As You Earn
		thresholds.PovertyGuideline = params.PovertyGuideline
		thresholds.PaymentCapPercentage = sdk.NewDecWithPrec(10, 2) // 10% of discretionary income
		thresholds.DiscretionaryIncomeThreshold = sdk.NewDecWithPrec(15, 1) // 150% of poverty guideline
		thresholds.MaxPaymentCap = sdk.ZeroDec() // No payment cap
		
	case "ICR": // Income-Contingent Repayment
		thresholds.PovertyGuideline = params.PovertyGuideline
		thresholds.PaymentCapPercentage = sdk.NewDecWithPrec(20, 2) // 20% of discretionary income
		thresholds.DiscretionaryIncomeThreshold = sdk.NewDecWithPrec(10, 1) // 100% of poverty guideline
		thresholds.MaxPaymentCap = loanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(12, 2)) // 12% of original loan
	}

	return thresholds
}

func (idrp *IncomeBasedRepaymentProcessor) calculateInitialIncomeBasedPayment(ctx sdk.Context, loan ProcessedEducationLoan, planType string, income sdk.Dec) sdk.Coin {
	thresholds := idrp.createIncomeThresholdsForPlan(ctx, planType, loan.LoanAmount)
	
	// Calculate annual income
	annualIncome := income.Mul(sdk.NewInt(12))
	
	// Calculate discretionary income
	povertyThreshold := thresholds.PovertyGuideline.Mul(thresholds.DiscretionaryIncomeThreshold)
	discretionaryIncome := annualIncome.Sub(povertyThreshold)
	
	// Ensure discretionary income is not negative
	if discretionaryIncome.IsNegative() {
		discretionaryIncome = sdk.ZeroDec()
	}
	
	// Calculate payment as percentage of discretionary income
	annualPayment := discretionaryIncome.Mul(thresholds.PaymentCapPercentage)
	monthlyPayment := annualPayment.QuoInt64(12)
	
	// Apply payment cap if applicable
	if !thresholds.MaxPaymentCap.IsZero() && monthlyPayment.GT(thresholds.MaxPaymentCap.QuoInt64(12)) {
		monthlyPayment = thresholds.MaxPaymentCap.QuoInt64(12)
	}
	
	// Minimum payment protection
	if monthlyPayment.LT(sdk.NewDec(100)) {
		monthlyPayment = sdk.NewDec(100) // Minimum ₹100 payment
	}
	
	return sdk.NewCoin(loan.LoanAmount.Denom, monthlyPayment.TruncateInt())
}

func (idrp *IncomeBasedRepaymentProcessor) determineForgiveness(ctx sdk.Context, loan ProcessedEducationLoan, planType string) (bool, int64) {
	switch planType {
	case "IBR":
		return true, 25 // 25 years
	case "PAYE":
		return true, 20 // 20 years
	case "REPAYE":
		if loan.LoanType == "UNDERGRADUATE" {
			return true, 20 // 20 years for undergraduate
		} else {
			return true, 25 // 25 years for graduate
		}
	case "ICR":
		return true, 25 // 25 years
	default:
		return false, 0
	}
}

func (idrp *IncomeBasedRepaymentProcessor) getPaymentCapPercentage(planType string) sdk.Dec {
	switch planType {
	case "IBR", "PAYE", "REPAYE":
		return sdk.NewDecWithPrec(10, 2) // 10%
	case "ICR":
		return sdk.NewDecWithPrec(20, 2) // 20%
	default:
		return sdk.NewDecWithPrec(10, 2) // Default 10%
	}
}

func (idrp *IncomeBasedRepaymentProcessor) createInterestCapitalizationRules(planType string) types.InterestCapitalization {
	rules := types.InterestCapitalization{
		PlanType: planType,
	}

	switch planType {
	case "IBR":
		rules.CapitalizationEvents = []string{"Plan_Exit", "Payment_Default", "Income_Increase_Above_Cap"}
		rules.MaxCapitalizationAmount = sdk.NewDecWithPrec(10, 2) // 10% of original principal
		rules.SubsidizedInterestCap = true
		
	case "PAYE":
		rules.CapitalizationEvents = []string{"Plan_Exit", "Payment_Default"}
		rules.MaxCapitalizationAmount = sdk.NewDecWithPrec(10, 2) // 10% of original principal
		rules.SubsidizedInterestCap = true
		
	case "REPAYE":
		rules.CapitalizationEvents = []string{"Plan_Exit"}
		rules.MaxCapitalizationAmount = sdk.NewDecWithPrec(10, 2) // 10% of original principal
		rules.SubsidizedInterestCap = true
		rules.UnsubsidizedInterestCap = true // 50% subsidy on unsubsidized loans
		
	case "ICR":
		rules.CapitalizationEvents = []string{"Plan_Exit", "Payment_Default", "Annual_Review"}
		rules.MaxCapitalizationAmount = sdk.ZeroDec() // No cap
		rules.SubsidizedInterestCap = false
	}

	return rules
}

func (idrp *IncomeBasedRepaymentProcessor) assessSubsidyEligibility(ctx sdk.Context, loan ProcessedEducationLoan, income sdk.Dec) types.SubsidyEligibility {
	eligibility := types.SubsidyEligibility{
		LoanID: loan.LoanID,
	}

	params := idrp.keeper.GetParams(ctx)
	
	// Check income eligibility for subsidy
	povertyThreshold := params.PovertyGuideline.Mul(sdk.NewDecWithPrec(15, 1)) // 150% of poverty guideline
	annualIncome := income.Mul(sdk.NewInt(12))
	
	if annualIncome.LTE(povertyThreshold) {
		eligibility.IsEligible = true
		eligibility.SubsidyType = "INTEREST_SUBSIDY"
		eligibility.SubsidyPercentage = sdk.NewDecWithPrec(50, 2) // 50% subsidy
		eligibility.MaxSubsidyAmount = loan.LoanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(5, 2)) // 5% of loan amount per year
	}

	return eligibility
}

func (idrp *IncomeBasedRepaymentProcessor) createHardshipProvisions(planType string) types.HardshipProvisions {
	provisions := types.HardshipProvisions{
		PlanType: planType,
	}

	// Common hardship provisions for all plan types
	provisions.UnemploymentDeferment = 36 // 36 months maximum
	provisions.EconomicHardshipDeferment = 36 // 36 months maximum
	provisions.ForbearanceOptions = []string{"General", "Mandatory", "Administrative"}
	provisions.ZeroPaymentEligibility = true
	provisions.PaymentPostponement = 12 // 12 months maximum

	// Plan-specific provisions
	switch planType {
	case "REPAYE":
		provisions.InterestSubsidyDuringHardship = true
	}

	return provisions
}

// Additional utility functions
func (idrp *IncomeBasedRepaymentProcessor) generateRepaymentPlanID(ctx sdk.Context, loanID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("IRP-%s-%d", loanID[:8], timestamp)
}

func (idrp *IncomeBasedRepaymentProcessor) generateAdjustmentID(ctx sdk.Context, planID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("ADJ-%s-%d", planID[:8], timestamp)
}

func (idrp *IncomeBasedRepaymentProcessor) generateIncomePaymentID(ctx sdk.Context, planID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("IPY-%s-%d", planID[:8], timestamp)
}

// Additional helper functions would include:
// - validateIncomeVerification
// - determineAdjustmentType
// - evaluateHardshipEligibility
// - allocateIncomeBasedPayment
// - calculateInterestSubsidy
// - processInterestSubsidy
// - checkInterestCapitalization
// - updateForgivenessProgress
// - countQualifyingPayments
// - processLoanForgiveness
// - calculateOutstandingAmount