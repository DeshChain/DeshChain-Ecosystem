package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
)

// BusinessLoanProcessor handles comprehensive business loan processing
type BusinessLoanProcessor struct {
	keeper         Keeper
	creditAnalyzer *BusinessCreditAnalyzer
}

// NewBusinessLoanProcessor creates a new business loan processor
func NewBusinessLoanProcessor(keeper Keeper) *BusinessLoanProcessor {
	return &BusinessLoanProcessor{
		keeper:         keeper,
		creditAnalyzer: NewBusinessCreditAnalyzer(keeper),
	}
}

// ProcessedBusinessLoan represents a comprehensive business loan
type ProcessedBusinessLoan struct {
	LoanID              string                        `json:"loan_id"`
	ApplicationID       string                        `json:"application_id"`
	BusinessID          string                        `json:"business_id"`
	LoanType            string                        `json:"loan_type"` // WORKING_CAPITAL, EQUIPMENT, EXPANSION, INVOICE_FINANCING
	BusinessCategory    string                        `json:"business_category"`
	LoanAmount          sdk.Coin                      `json:"loan_amount"`
	InterestRate        sdk.Dec                       `json:"interest_rate"`
	LoanTenure          int64                         `json:"loan_tenure"` // months
	RepaymentFrequency  string                        `json:"repayment_frequency"` // MONTHLY, QUARTERLY, BULLET
	RepaymentSchedule   []types.BusinessRepayment     `json:"repayment_schedule"`
	CollateralRequired  bool                          `json:"collateral_required"`
	CollateralDetails   types.BusinessCollateral      `json:"collateral_details,omitempty"`
	GuaranteeRequired   bool                          `json:"guarantee_required"`
	GuaranteeDetails    types.BusinessGuarantee       `json:"guarantee_details,omitempty"`
	CreditLineComponent *types.CreditLineInfo         `json:"credit_line,omitempty"`
	InvoiceFinancing    []types.InvoiceFinancingInfo  `json:"invoice_financing,omitempty"`
	Purpose             string                        `json:"purpose"`
	DisbursementInfo    types.BusinessDisbursement    `json:"disbursement_info"`
	RepaymentHistory    []types.BusinessRepayment     `json:"repayment_history"`
	Status              types.BusinessLoanStatus      `json:"status"`
	CreditProfile       *BusinessCreditProfile        `json:"credit_profile"`
	RiskAssessment      types.BusinessRiskAssessment  `json:"risk_assessment"`
	ApprovalWorkflow    types.ApprovalWorkflow        `json:"approval_workflow"`
	ComplianceChecks    types.ComplianceChecks        `json:"compliance_checks"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
	MaturityDate        time.Time                     `json:"maturity_date"`
}

// ProcessBusinessLoanApplication processes a comprehensive business loan application
func (blp *BusinessLoanProcessor) ProcessBusinessLoanApplication(ctx sdk.Context, applicationID string) (*ProcessedBusinessLoan, error) {
	// Get loan application
	application, found := blp.keeper.GetLoanApplication(ctx, applicationID)
	if !found {
		return nil, fmt.Errorf("business loan application not found: %s", applicationID)
	}

	// Validate application status
	if application.Status != types.ApplicationStatusSubmitted {
		return nil, fmt.Errorf("application not in submitted status: %s", application.Status)
	}

	// Perform comprehensive credit analysis
	creditProfile, err := blp.creditAnalyzer.AnalyzeBusinessCredit(ctx, application.BusinessID)
	if err != nil {
		return nil, fmt.Errorf("business credit analysis failed: %w", err)
	}

	// Update application with credit analysis results
	application.Status = types.ApplicationStatusUnderReview
	application.CreditScore = creditProfile.CreditScore
	application.RiskCategory = creditProfile.RiskCategory
	blp.keeper.SetLoanApplication(ctx, application)

	// Perform compliance checks
	complianceChecks, err := blp.performComplianceChecks(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("compliance checks failed: %w", err)
	}

	// Check eligibility based on comprehensive analysis
	eligible, eligibilityReason := blp.checkBusinessEligibility(ctx, application, creditProfile)
	if !eligible {
		return blp.rejectBusinessApplication(ctx, application, eligibilityReason)
	}

	// Determine loan parameters
	loanAmount := application.RequestedAmount
	if loanAmount.Amount.GT(creditProfile.MaxLoanEligibility.Amount) {
		loanAmount = creditProfile.MaxLoanEligibility
	}

	interestRate := creditProfile.RecommendedRate

	// Apply business-specific discounts and festival offers
	interestRate = blp.applyBusinessDiscounts(ctx, application.BusinessID, interestRate)

	// Determine loan tenure based on loan type and business profile
	tenure := blp.calculateOptimalTenure(ctx, application.LoanType, loanAmount, creditProfile)

	// Determine repayment frequency
	repaymentFreq := blp.determineRepaymentFrequency(ctx, application.LoanType, creditProfile)

	// Generate repayment schedule
	repaymentSchedule := blp.generateBusinessRepaymentSchedule(
		ctx, loanAmount, interestRate, tenure, repaymentFreq, application.LoanType,
	)

	// Assess collateral and guarantee requirements
	collateralRequired, collateralDetails := blp.assessCollateralRequirements(
		ctx, loanAmount, creditProfile.RiskCategory, application.LoanType,
	)
	guaranteeRequired, guaranteeDetails := blp.assessGuaranteeRequirements(
		ctx, loanAmount, creditProfile, application.LoanType,
	)

	// Create credit line component if applicable
	var creditLineInfo *types.CreditLineInfo
	if application.LoanType == "WORKING_CAPITAL" && creditProfile.CreditLineEligibility.IsPositive() {
		creditLineInfo = &types.CreditLineInfo{
			CreditLimit:     creditProfile.CreditLineEligibility,
			UtilizedAmount:  sdk.ZeroInt(),
			AvailableAmount: creditProfile.CreditLineEligibility.Amount,
			InterestRate:    interestRate.Add(sdk.NewDecWithPrec(5, 3)), // +0.5% for credit line
			IsActive:        true,
			CreatedAt:       ctx.BlockTime(),
		}
	}

	// Handle invoice financing if applicable
	var invoiceFinancing []types.InvoiceFinancingInfo
	if application.LoanType == "INVOICE_FINANCING" {
		invoiceFinancing = blp.setupInvoiceFinancing(ctx, application, creditProfile)
	}

	// Generate loan ID
	loanID := blp.generateBusinessLoanID(ctx, application.BusinessID)

	// Create business loan
	businessLoan := &ProcessedBusinessLoan{
		LoanID:             loanID,
		ApplicationID:      applicationID,
		BusinessID:         application.BusinessID,
		LoanType:          application.LoanType,
		BusinessCategory:  application.BusinessCategory,
		LoanAmount:        loanAmount,
		InterestRate:      interestRate,
		LoanTenure:        tenure,
		RepaymentFrequency: repaymentFreq,
		RepaymentSchedule: repaymentSchedule,
		CollateralRequired: collateralRequired,
		CollateralDetails:  collateralDetails,
		GuaranteeRequired:  guaranteeRequired,
		GuaranteeDetails:   guaranteeDetails,
		CreditLineComponent: creditLineInfo,
		InvoiceFinancing:   invoiceFinancing,
		Purpose:           application.Purpose,
		DisbursementInfo: types.BusinessDisbursement{
			RequestedDate:    application.PreferredDisbursementDate,
			ApprovedDate:     ctx.BlockTime(),
			Status:           types.DisbursementStatusPending,
			DisbursementMode: blp.determineDisbursementMode(ctx, application.LoanType),
		},
		Status:         types.BusinessLoanStatusApproved,
		CreditProfile:  creditProfile,
		RiskAssessment: blp.createRiskAssessment(ctx, creditProfile),
		ApprovalWorkflow: types.ApprovalWorkflow{
			ApprovedBy:      "AUTOMATED_SYSTEM",
			ApprovedAmount:  loanAmount,
			ApprovedRate:    interestRate,
			ApprovalDate:    ctx.BlockTime(),
			ApprovalLevel:   blp.determineApprovalLevel(ctx, loanAmount),
			Conditions:      blp.generateApprovalConditions(ctx, collateralRequired, guaranteeRequired),
		},
		ComplianceChecks: complianceChecks,
		CreatedAt:        ctx.BlockTime(),
		UpdatedAt:        ctx.BlockTime(),
		MaturityDate:     ctx.BlockTime().AddDate(0, int(tenure), 0),
	}

	// Store business loan
	blp.keeper.SetProcessedBusinessLoan(ctx, *businessLoan)

	// Update application status
	application.Status = types.ApplicationStatusApproved
	application.ApprovedLoanID = loanID
	blp.keeper.SetLoanApplication(ctx, application)

	// Emit loan approval event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBusinessLoanApproved,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyBusinessID, application.BusinessID),
			sdk.NewAttribute(types.AttributeKeyLoanAmount, loanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, interestRate.String()),
			sdk.NewAttribute(types.AttributeKeyCreditScore, fmt.Sprintf("%d", creditProfile.CreditScore)),
			sdk.NewAttribute(types.AttributeKeyLoanType, application.LoanType),
		),
	)

	return businessLoan, nil
}

// DisburseBusinesLoan handles business loan disbursement with comprehensive checks
func (blp *BusinessLoanProcessor) DisburseBusinessLoan(ctx sdk.Context, loanID string) error {
	// Get loan details
	loan, found := blp.keeper.GetProcessedBusinessLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("business loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.BusinessLoanStatusApproved {
		return fmt.Errorf("loan not in approved status: %s", loan.Status)
	}

	// Check disbursement conditions
	if loan.CollateralRequired && !blp.isCollateralVerified(ctx, loanID) {
		return fmt.Errorf("collateral verification pending for loan: %s", loanID)
	}

	if loan.GuaranteeRequired && !blp.isGuaranteeVerified(ctx, loanID) {
		return fmt.Errorf("guarantee verification pending for loan: %s", loanID)
	}

	// Perform final compliance checks
	if !blp.validateFinalCompliance(ctx, loan) {
		return fmt.Errorf("final compliance checks failed for loan: %s", loanID)
	}

	// Get business address
	businessAddr, err := sdk.AccAddressFromBech32(loan.BusinessID)
	if err != nil {
		return fmt.Errorf("invalid business address: %s", loan.BusinessID)
	}

	// Calculate disbursement fees and processing charges
	processingFee := blp.calculateProcessingFee(ctx, loan.LoanAmount, loan.LoanType)
	disbursementFee := blp.calculateBusinessDisbursementFee(ctx, loan.LoanAmount)
	totalFees := processingFee.Add(disbursementFee)
	netDisbursement := loan.LoanAmount.Sub(totalFees)

	// Handle different disbursement modes
	switch loan.DisbursementInfo.DisbursementMode {
	case "DIRECT_TRANSFER":
		err = blp.processDirectTransfer(ctx, businessAddr, netDisbursement)
	case "INVOICE_FINANCING":
		err = blp.processInvoiceFinancingDisbursement(ctx, loan)
	case "EQUIPMENT_VENDOR":
		err = blp.processVendorPayment(ctx, loan, netDisbursement)
	default:
		err = blp.processDirectTransfer(ctx, businessAddr, netDisbursement)
	}

	if err != nil {
		return fmt.Errorf("failed to disburse business loan: %w", err)
	}

	// Distribute fees to appropriate pools
	if totalFees.Amount.GT(sdk.ZeroInt()) {
		err = blp.distributeDisbursementFees(ctx, totalFees, loan.LoanType)
		if err != nil {
			blp.keeper.Logger(ctx).Error("Failed to distribute disbursement fees", 
				"fees", totalFees.String(), "error", err)
		}
	}

	// Update loan status and disbursement info
	loan.Status = types.BusinessLoanStatusActive
	loan.DisbursementInfo.Status = types.DisbursementStatusCompleted
	loan.DisbursementInfo.DisbursedAmount = netDisbursement
	loan.DisbursementInfo.ProcessingFee = processingFee
	loan.DisbursementInfo.DisbursementFee = disbursementFee
	loan.DisbursementInfo.DisbursedDate = ctx.BlockTime()
	loan.UpdatedAt = ctx.BlockTime()

	blp.keeper.SetProcessedBusinessLoan(ctx, loan)

	// Update business profile with active loan
	blp.updateBusinessProfileAfterDisbursement(ctx, loan.BusinessID, loan)

	// Activate credit line if applicable
	if loan.CreditLineComponent != nil {
		err = blp.activateCreditLine(ctx, loan.BusinessID, *loan.CreditLineComponent)
		if err != nil {
			blp.keeper.Logger(ctx).Error("Failed to activate credit line", 
				"business_id", loan.BusinessID, "error", err)
		}
	}

	// Emit disbursement event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBusinessLoanDisbursed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyDisbursedAmount, netDisbursement.String()),
			sdk.NewAttribute(types.AttributeKeyProcessingFee, processingFee.String()),
			sdk.NewAttribute(types.AttributeKeyDisbursementMode, loan.DisbursementInfo.DisbursementMode),
		),
	)

	return nil
}

// ProcessBusinessRepayment handles business loan repayment with flexible options
func (blp *BusinessLoanProcessor) ProcessBusinessRepayment(ctx sdk.Context, loanID string, repaymentAmount sdk.Coin, repaymentType string) error {
	// Get loan details
	loan, found := blp.keeper.GetProcessedBusinessLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("business loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.BusinessLoanStatusActive {
		return fmt.Errorf("loan not in active status: %s", loan.Status)
	}

	// Get business address
	businessAddr, err := sdk.AccAddressFromBech32(loan.BusinessID)
	if err != nil {
		return fmt.Errorf("invalid business address: %s", loan.BusinessID)
	}

	// Calculate outstanding amount
	outstandingAmount := blp.calculateBusinessOutstandingAmount(ctx, loan)
	if repaymentAmount.Amount.GT(outstandingAmount.Amount) {
		return fmt.Errorf("repayment amount exceeds outstanding balance")
	}

	// Process payment from business to module
	err = blp.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		businessAddr,
		types.ModuleName,
		sdk.NewCoins(repaymentAmount),
	)
	if err != nil {
		return fmt.Errorf("failed to process business repayment: %w", err)
	}

	// Allocate repayment (principal vs interest vs fees)
	principalPortion, interestPortion, feesPortion := blp.allocateBusinessRepayment(ctx, loan, repaymentAmount)

	// Create repayment record
	repaymentRecord := types.BusinessRepayment{
		RepaymentID:      blp.generateBusinessRepaymentID(ctx, loanID),
		LoanID:           loanID,
		RepaymentDate:    ctx.BlockTime(),
		TotalAmount:      repaymentAmount,
		PrincipalAmount:  principalPortion,
		InterestAmount:   interestPortion,
		FeesAmount:       feesPortion,
		RepaymentType:    repaymentType,
		Status:           types.RepaymentStatusCompleted,
		OutstandingAfter: outstandingAmount.Sub(repaymentAmount),
	}

	// Update loan with repayment
	loan.RepaymentHistory = append(loan.RepaymentHistory, repaymentRecord)
	loan.UpdatedAt = ctx.BlockTime()

	// Check if loan is fully repaid
	newOutstanding := outstandingAmount.Sub(repaymentAmount)
	if newOutstanding.IsZero() {
		loan.Status = types.BusinessLoanStatusClosed
		
		// Process loan closure
		err = blp.processBusinessLoanClosure(ctx, loan)
		if err != nil {
			blp.keeper.Logger(ctx).Error("Failed to process business loan closure", 
				"loan_id", loanID, "error", err)
		}
	}

	blp.keeper.SetProcessedBusinessLoan(ctx, loan)

	// Update business profile and merchant rating
	blp.updateBusinessProfileAfterRepayment(ctx, loan.BusinessID, repaymentRecord)
	blp.keeper.UpdateMerchantRating(ctx, loan.BusinessID, types.Repayment{
		LoanID:    loanID,
		Amount:    repaymentAmount,
		PaidAt:    ctx.BlockTime(),
		Type:      repaymentType,
	})

	// Handle credit line repayment if applicable
	if loan.CreditLineComponent != nil && repaymentType == "CREDIT_LINE_REPAYMENT" {
		err = blp.processCreditLineRepayment(ctx, loan.BusinessID, principalPortion)
		if err != nil {
			blp.keeper.Logger(ctx).Error("Failed to process credit line repayment", 
				"business_id", loan.BusinessID, "error", err)
		}
	}

	// Emit repayment event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBusinessRepaymentProcessed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyRepaymentAmount, repaymentAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPrincipalAmount, principalPortion.String()),
			sdk.NewAttribute(types.AttributeKeyInterestAmount, interestPortion.String()),
			sdk.NewAttribute(types.AttributeKeyRepaymentType, repaymentType),
			sdk.NewAttribute(types.AttributeKeyOutstandingAmount, newOutstanding.String()),
		),
	)

	return nil
}

// Helper functions for business loan processing
func (blp *BusinessLoanProcessor) checkBusinessEligibility(ctx sdk.Context, application types.LoanApplication, creditProfile *BusinessCreditProfile) (bool, string) {
	// Check minimum credit score
	if creditProfile.CreditScore < 600 {
		return false, "Credit score too low for business loan (minimum 600 required)"
	}

	// Check business age
	if creditProfile.BusinessInfo.YearsInBusiness < 1 {
		return false, "Business must be operational for at least 1 year"
	}

	// Check financial health
	if creditProfile.FinancialMetrics.CurrentRatio.LT(sdk.NewDecWithPrec(8, 1)) { // < 0.8
		return false, "Poor liquidity ratio - current ratio below acceptable threshold"
	}

	if creditProfile.FinancialMetrics.DebtToEquityRatio.GT(sdk.NewDec(3)) { // > 3.0
		return false, "Excessive leverage - debt to equity ratio too high"
	}

	// Check active loans limit
	if creditProfile.BusinessInfo.ActiveLoans >= 5 {
		return false, "Maximum 5 active business loans allowed"
	}

	// Check for red flags
	if len(creditProfile.RedFlags) > 0 {
		for _, flag := range creditProfile.RedFlags {
			if flag == "DEFAULTED_LOANS" || flag == "REGULATORY_VIOLATIONS" {
				return false, fmt.Sprintf("Critical red flag identified: %s", flag)
			}
		}
	}

	return true, ""
}

func (blp *BusinessLoanProcessor) calculateOptimalTenure(ctx sdk.Context, loanType string, loanAmount sdk.Coin, creditProfile *BusinessCreditProfile) int64 {
	// Loan type-specific tenure recommendations
	tenureMap := map[string]int64{
		"WORKING_CAPITAL":   12, // 12 months for working capital
		"EQUIPMENT":         36, // 36 months for equipment
		"EXPANSION":         60, // 60 months for expansion
		"INVOICE_FINANCING": 3,  // 3 months for invoice financing
		"EMERGENCY":         6,  // 6 months for emergency loans
	}

	baseTenure := tenureMap[loanType]
	if baseTenure == 0 {
		baseTenure = 24 // Default 24 months
	}

	// Adjust based on loan amount and business profile
	if loanAmount.Amount.GT(sdk.NewInt(5000000)) { // > 50 lakh
		baseTenure = baseTenure + 12 // Add 12 months for large loans
	}

	// Adjust based on cash flow stability
	if creditProfile.CashFlowAnalysis.CashFlowStability.LT(sdk.NewDecWithPrec(7, 1)) { // < 70% stability
		baseTenure = baseTenure + 6 // Add 6 months for unstable cash flows
	}

	// Cap the tenure
	if baseTenure > 84 { // Max 7 years
		baseTenure = 84
	}

	return baseTenure
}

func (blp *BusinessLoanProcessor) determineRepaymentFrequency(ctx sdk.Context, loanType string, creditProfile *BusinessCreditProfile) string {
	// Default frequencies by loan type
	switch loanType {
	case "WORKING_CAPITAL":
		return "MONTHLY"
	case "EQUIPMENT":
		return "MONTHLY"
	case "EXPANSION":
		return "QUARTERLY"
	case "INVOICE_FINANCING":
		return "BULLET" // Single payment at maturity
	default:
		return "MONTHLY"
	}
}

func (blp *BusinessLoanProcessor) applyBusinessDiscounts(ctx sdk.Context, businessID string, baseRate sdk.Dec) sdk.Dec {
	// Apply festival discounts
	festivalOffers := blp.keeper.GetActiveFestivalOffers(ctx)
	for _, offer := range festivalOffers {
		discount, _ := sdk.NewDecFromStr(offer.InterestReduction)
		baseRate = baseRate.Sub(discount)
	}

	// Apply business-specific discounts
	businessProfile, found := blp.keeper.GetBusinessProfile(ctx, businessID)
	if found {
		if businessProfile.IsWomenOwned {
			params := blp.keeper.GetParams(ctx)
			womenDiscount, _ := sdk.NewDecFromStr(params.WomenEntrepreneurDiscount)
			baseRate = baseRate.Sub(womenDiscount)
		}

		if businessProfile.IsStartup {
			params := blp.keeper.GetParams(ctx)
			startupDiscount, _ := sdk.NewDecFromStr(params.StartupDiscount)
			baseRate = baseRate.Sub(startupDiscount)
		}

		// MSME category discount
		if businessProfile.Category == "MSME" {
			msmeDiscount := sdk.NewDecWithPrec(5, 3) // 0.5%
			baseRate = baseRate.Sub(msmeDiscount)
		}
	}

	// Ensure rate doesn't go below minimum
	params := blp.keeper.GetParams(ctx)
	minRate, _ := sdk.NewDecFromStr(params.MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}

	return baseRate
}

// Additional helper functions would include:
// - generateBusinessRepaymentSchedule
// - assessGuaranteeRequirements
// - setupInvoiceFinancing
// - performComplianceChecks
// - createRiskAssessment
// - processDirectTransfer
// - And other business loan specific operations

func (blp *BusinessLoanProcessor) rejectBusinessApplication(ctx sdk.Context, application types.LoanApplication, reason string) (*ProcessedBusinessLoan, error) {
	application.Status = types.ApplicationStatusRejected
	application.RejectionReason = reason
	blp.keeper.SetLoanApplication(ctx, application)

	// Emit rejection event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBusinessLoanRejected,
			sdk.NewAttribute(types.AttributeKeyApplicationID, application.ID),
			sdk.NewAttribute(types.AttributeKeyBusinessID, application.BusinessID),
			sdk.NewAttribute(types.AttributeKeyRejectionReason, reason),
		),
	)

	return nil, fmt.Errorf("business loan application rejected: %s", reason)
}

func (blp *BusinessLoanProcessor) generateBusinessLoanID(ctx sdk.Context, businessID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("VYM-%s-%d", businessID[:8], timestamp)
}
