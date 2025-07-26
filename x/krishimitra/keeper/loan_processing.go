package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/krishimitra/types"
)

// LoanProcessor handles comprehensive agricultural loan processing
type LoanProcessor struct {
	keeper         Keeper
	creditEngine   *CreditScoringEngine
	insuranceEngine *InsuranceEngine
}

// NewLoanProcessor creates a new loan processor
func NewLoanProcessor(keeper Keeper) *LoanProcessor {
	return &LoanProcessor{
		keeper:          keeper,
		creditEngine:    NewCreditScoringEngine(keeper),
		insuranceEngine: NewInsuranceEngine(keeper),
	}
}

// AgriculturalLoan represents a comprehensive agricultural loan
type AgriculturalLoan struct {
	LoanID              string                        `json:"loan_id"`
	ApplicationID       string                        `json:"application_id"`
	FarmerID            string                        `json:"farmer_id"`
	LoanType            string                        `json:"loan_type"` // CROP, EQUIPMENT, INFRASTRUCTURE
	CropType            string                        `json:"crop_type,omitempty"`
	LoanAmount          sdk.Coin                      `json:"loan_amount"`
	InterestRate        sdk.Dec                       `json:"interest_rate"`
	LoanTenure          int64                         `json:"loan_tenure"` // months
	RepaymentSchedule   []types.RepaymentInstallment  `json:"repayment_schedule"`
	CollateralRequired  bool                          `json:"collateral_required"`
	CollateralDetails   types.CollateralInfo          `json:"collateral_details,omitempty"`
	InsuranceRequired   bool                          `json:"insurance_required"`
	InsurancePolicyID   string                        `json:"insurance_policy_id,omitempty"`
	Purpose             string                        `json:"purpose"`
	SeasonType          string                        `json:"season_type"` // KHARIF, RABI, ZAID
	DisbursementDetails types.DisbursementInfo        `json:"disbursement_details"`
	RepaymentHistory    []types.RepaymentRecord       `json:"repayment_history"`
	Status              types.LoanStatus              `json:"status"`
	CreditScore         int64                         `json:"credit_score"`
	RiskCategory        string                        `json:"risk_category"`
	ApprovalDetails     types.ApprovalInfo            `json:"approval_details"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
	MaturityDate        time.Time                     `json:"maturity_date"`
}

// ProcessLoanApplication processes a comprehensive loan application
func (lp *LoanProcessor) ProcessLoanApplication(ctx sdk.Context, applicationID string) (*AgriculturalLoan, error) {
	// Get loan application
	application, found := lp.keeper.GetLoanApplication(ctx, applicationID)
	if !found {
		return nil, fmt.Errorf("loan application not found: %s", applicationID)
	}

	// Validate application status
	if application.Status != types.ApplicationStatusSubmitted {
		return nil, fmt.Errorf("application not in submitted status: %s", application.Status)
	}

	// Perform comprehensive credit assessment
	eligibilityAssessment, err := lp.creditEngine.AssessLoanEligibility(
		ctx,
		application.FarmerID,
		application.RequestedAmount,
		application.CropType,
	)
	if err != nil {
		return nil, fmt.Errorf("credit assessment failed: %w", err)
	}

	// Update application status
	application.Status = types.ApplicationStatusUnderReview
	application.CreditScore = eligibilityAssessment.CreditScore
	application.RiskCategory = eligibilityAssessment.RiskCategory
	lp.keeper.SetLoanApplication(ctx, application)

	// Check eligibility
	if !eligibilityAssessment.IsEligible {
		return lp.rejectLoanApplication(ctx, application, eligibilityAssessment.EligibilityReason)
	}

	// Determine loan parameters
	loanAmount := application.RequestedAmount
	if loanAmount.Amount.GT(eligibilityAssessment.MaxEligibleAmount.Amount) {
		loanAmount = eligibilityAssessment.MaxEligibleAmount
	}

	interestRate := eligibilityAssessment.RecommendedRate

	// Apply any festival or special discounts
	interestRate = lp.applyDiscounts(ctx, application.FarmerID, interestRate)

	// Calculate loan tenure
	tenure := lp.calculateOptimalTenure(ctx, application.CropType, loanAmount)

	// Generate repayment schedule
	repaymentSchedule := lp.generateRepaymentSchedule(ctx, loanAmount, interestRate, tenure, application.CropType)

	// Determine collateral requirements
	collateralRequired, collateralDetails := lp.assessCollateralRequirements(ctx, loanAmount, eligibilityAssessment.RiskCategory)

	// Determine insurance requirements
	insuranceRequired := lp.isInsuranceRequired(ctx, application.CropType, loanAmount)

	// Generate loan ID
	loanID := lp.generateLoanID(ctx, application.FarmerID)

	// Create agricultural loan
	loan := &AgriculturalLoan{
		LoanID:            loanID,
		ApplicationID:     applicationID,
		FarmerID:          application.FarmerID,
		LoanType:          "CROP_LOAN",
		CropType:          application.CropType,
		LoanAmount:        loanAmount,
		InterestRate:      interestRate,
		LoanTenure:        tenure,
		RepaymentSchedule: repaymentSchedule,
		CollateralRequired: collateralRequired,
		CollateralDetails:  collateralDetails,
		InsuranceRequired:  insuranceRequired,
		Purpose:           application.Purpose,
		SeasonType:        application.SeasonType,
		DisbursementDetails: types.DisbursementInfo{
			RequestedDate:    application.PreferredDisbursementDate,
			ApprovedDate:     ctx.BlockTime(),
			Status:           types.DisbursementStatusPending,
		},
		Status:       types.LoanStatusApproved,
		CreditScore:  eligibilityAssessment.CreditScore,
		RiskCategory: eligibilityAssessment.RiskCategory,
		ApprovalDetails: types.ApprovalInfo{
			ApprovedBy:      "AUTOMATED_SYSTEM",
			ApprovedAmount:  loanAmount,
			ApprovedRate:    interestRate,
			ApprovalDate:    ctx.BlockTime(),
			Conditions:      lp.generateApprovalConditions(ctx, collateralRequired, insuranceRequired),
		},
		CreatedAt:    ctx.BlockTime(),
		UpdatedAt:    ctx.BlockTime(),
		MaturityDate: ctx.BlockTime().AddDate(0, int(tenure), 0),
	}

	// Store loan
	lp.keeper.SetAgriculturalLoan(ctx, *loan)

	// Update application status
	application.Status = types.ApplicationStatusApproved
	application.ApprovedLoanID = loanID
	lp.keeper.SetLoanApplication(ctx, application)

	// Create insurance policy if required
	if insuranceRequired {
		insurancePolicyID, err := lp.createMandatoryInsurance(ctx, loan)
		if err != nil {
			lp.keeper.Logger(ctx).Error("Failed to create mandatory insurance", "loan_id", loanID, "error", err)
		} else {
			loan.InsurancePolicyID = insurancePolicyID
			lp.keeper.SetAgriculturalLoan(ctx, *loan)
		}
	}

	// Emit loan approval event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLoanApproved,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyFarmerID, application.FarmerID),
			sdk.NewAttribute(types.AttributeKeyLoanAmount, loanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, interestRate.String()),
			sdk.NewAttribute(types.AttributeKeyCreditScore, fmt.Sprintf("%d", eligibilityAssessment.CreditScore)),
		),
	)

	return loan, nil
}

// DisburseLoan handles loan disbursement
func (lp *LoanProcessor) DisburseLoan(ctx sdk.Context, loanID string) error {
	// Get loan details
	loan, found := lp.keeper.GetAgriculturalLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.LoanStatusApproved {
		return fmt.Errorf("loan not in approved status: %s", loan.Status)
	}

	// Check disbursement conditions
	if loan.CollateralRequired && !lp.isCollateralVerified(ctx, loanID) {
		return fmt.Errorf("collateral verification pending for loan: %s", loanID)
	}

	if loan.InsuranceRequired && loan.InsurancePolicyID == "" {
		return fmt.Errorf("insurance policy not created for loan: %s", loanID)
	}

	// Transfer loan amount to farmer
	farmerAddr, err := sdk.AccAddressFromBech32(loan.FarmerID)
	if err != nil {
		return fmt.Errorf("invalid farmer address: %s", loan.FarmerID)
	}

	// Calculate disbursement fees
	disbursementFee := lp.calculateDisbursementFee(ctx, loan.LoanAmount)
	netDisbursement := loan.LoanAmount.Sub(disbursementFee)

	// Transfer from KrishiMitra module to farmer
	err = lp.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		farmerAddr,
		sdk.NewCoins(netDisbursement),
	)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %w", err)
	}

	// Distribute disbursement fee
	if disbursementFee.Amount.GT(sdk.ZeroInt()) {
		err = lp.distributeDisbursementFee(ctx, disbursementFee)
		if err != nil {
			lp.keeper.Logger(ctx).Error("Failed to distribute disbursement fee", "fee", disbursementFee.String(), "error", err)
		}
	}

	// Update loan status
	loan.Status = types.LoanStatusActive
	loan.DisbursementDetails.Status = types.DisbursementStatusCompleted
	loan.DisbursementDetails.DisbursedAmount = netDisbursement
	loan.DisbursementDetails.DisbursementFee = disbursementFee
	loan.DisbursementDetails.DisbursedDate = ctx.BlockTime()
	loan.UpdatedAt = ctx.BlockTime()

	lp.keeper.SetAgriculturalLoan(ctx, loan)

	// Update farmer profile
	lp.updateFarmerProfileAfterDisbursement(ctx, loan.FarmerID, loan)

	// Emit disbursement event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLoanDisbursed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyDisbursedAmount, netDisbursement.String()),
			sdk.NewAttribute(types.AttributeKeyDisbursementFee, disbursementFee.String()),
		),
	)

	return nil
}

// ProcessRepayment handles loan repayment
func (lp *LoanProcessor) ProcessRepayment(ctx sdk.Context, loanID string, repaymentAmount sdk.Coin) error {
	// Get loan details
	loan, found := lp.keeper.GetAgriculturalLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("loan not found: %s", loanID)
	}

	// Validate loan status
	if loan.Status != types.LoanStatusActive {
		return fmt.Errorf("loan not in active status: %s", loan.Status)
	}

	// Get farmer address
	farmerAddr, err := sdk.AccAddressFromBech32(loan.FarmerID)
	if err != nil {
		return fmt.Errorf("invalid farmer address: %s", loan.FarmerID)
	}

	// Validate repayment amount
	outstandingAmount := lp.calculateOutstandingAmount(ctx, loan)
	if repaymentAmount.Amount.GT(outstandingAmount.Amount) {
		return fmt.Errorf("repayment amount exceeds outstanding balance")
	}

	// Transfer repayment from farmer to module
	err = lp.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		farmerAddr,
		types.ModuleName,
		sdk.NewCoins(repaymentAmount),
	)
	if err != nil {
		return fmt.Errorf("failed to process repayment: %w", err)
	}

	// Allocate repayment (principal vs interest)
	principalPortion, interestPortion := lp.allocateRepayment(ctx, loan, repaymentAmount)

	// Create repayment record
	repaymentRecord := types.RepaymentRecord{
		RepaymentID:     lp.generateRepaymentID(ctx, loanID),
		LoanID:          loanID,
		RepaymentDate:   ctx.BlockTime(),
		TotalAmount:     repaymentAmount,
		PrincipalAmount: principalPortion,
		InterestAmount:  interestPortion,
		RepaymentType:   lp.determineRepaymentType(ctx, loan, repaymentAmount),
		Status:          types.RepaymentStatusCompleted,
	}

	// Update loan with repayment
	loan.RepaymentHistory = append(loan.RepaymentHistory, repaymentRecord)
	loan.UpdatedAt = ctx.BlockTime()

	// Check if loan is fully repaid
	newOutstanding := outstandingAmount.Sub(repaymentAmount)
	if newOutstanding.IsZero() {
		loan.Status = types.LoanStatusClosed
		
		// Process loan closure
		err = lp.processLoanClosure(ctx, loan)
		if err != nil {
			lp.keeper.Logger(ctx).Error("Failed to process loan closure", "loan_id", loanID, "error", err)
		}
	}

	lp.keeper.SetAgriculturalLoan(ctx, loan)

	// Update farmer profile
	lp.updateFarmerProfileAfterRepayment(ctx, loan.FarmerID, repaymentRecord)

	// Emit repayment event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRepaymentProcessed,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyRepaymentAmount, repaymentAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPrincipalAmount, principalPortion.String()),
			sdk.NewAttribute(types.AttributeKeyInterestAmount, interestPortion.String()),
		),
	)

	return nil
}

// calculateOptimalTenure determines the best loan tenure for a crop type
func (lp *LoanProcessor) calculateOptimalTenure(ctx sdk.Context, cropType string, loanAmount sdk.Coin) int64 {
	// Crop-specific tenure recommendations
	tenureMap := map[string]int64{
		"RICE":       6,  // 6 months for rice
		"WHEAT":      6,  // 6 months for wheat
		"COTTON":     8,  // 8 months for cotton
		"SUGARCANE":  12, // 12 months for sugarcane
		"VEGETABLES": 4,  // 4 months for vegetables
		"PULSES":     5,  // 5 months for pulses
	}

	if tenure, found := tenureMap[cropType]; found {
		// Adjust based on loan amount (larger loans may need longer tenure)
		if loanAmount.Amount.GT(sdk.NewInt(100000)) { // > 1 lakh
			return tenure + 2 // Add 2 months for large loans
		}
		return tenure
	}

	return 6 // Default 6 months
}

// generateRepaymentSchedule creates a detailed repayment schedule
func (lp *LoanProcessor) generateRepaymentSchedule(ctx sdk.Context, loanAmount sdk.Coin, interestRate sdk.Dec, tenure int64, cropType string) []types.RepaymentInstallment {
	schedule := []types.RepaymentInstallment{}

	// Calculate total interest
	monthlyRate := interestRate.QuoInt64(12) // Convert annual rate to monthly
	totalInterest := loanAmount.Amount.ToDec().Mul(monthlyRate).Mul(sdk.NewDec(tenure))
	totalAmount := loanAmount.Amount.ToDec().Add(totalInterest)

	// For agricultural loans, typically structured as:
	// - Interest-only payments during crop growth
	// - Principal + interest at harvest time
	
	if lp.isCropLoan(cropType) {
		// Interest-only payments for most of the tenure
		interestOnlyMonths := tenure - 1
		monthlyInterest := totalInterest.QuoInt64(tenure)
		
		// Interest-only installments
		for i := int64(1); i <= interestOnlyMonths; i++ {
			installment := types.RepaymentInstallment{
				InstallmentNumber: i,
				DueDate:          ctx.BlockTime().AddDate(0, int(i), 0),
				TotalAmount:      sdk.NewCoin(loanAmount.Denom, monthlyInterest.TruncateInt()),
				PrincipalAmount:  sdk.NewCoin(loanAmount.Denom, sdk.ZeroInt()),
				InterestAmount:   sdk.NewCoin(loanAmount.Denom, monthlyInterest.TruncateInt()),
				Status:           types.InstallmentStatusPending,
			}
			schedule = append(schedule, installment)
		}

		// Final payment with principal
		finalPayment := loanAmount.Amount.Add(monthlyInterest.TruncateInt())
		finalInstallment := types.RepaymentInstallment{
			InstallmentNumber: tenure,
			DueDate:          ctx.BlockTime().AddDate(0, int(tenure), 0),
			TotalAmount:      sdk.NewCoin(loanAmount.Denom, finalPayment),
			PrincipalAmount:  loanAmount,
			InterestAmount:   sdk.NewCoin(loanAmount.Denom, monthlyInterest.TruncateInt()),
			Status:           types.InstallmentStatusPending,
		}
		schedule = append(schedule, finalInstallment)
	} else {
		// Equal monthly installments for non-crop loans
		monthlyPayment := totalAmount.QuoInt64(tenure)
		
		for i := int64(1); i <= tenure; i++ {
			// Calculate principal and interest portions
			outstandingPrincipal := loanAmount.Amount.ToDec().Sub(
				loanAmount.Amount.ToDec().Mul(sdk.NewDec(i-1)).QuoInt64(tenure),
			)
			monthlyInterestAmount := outstandingPrincipal.Mul(monthlyRate)
			monthlyPrincipalAmount := monthlyPayment.Sub(monthlyInterestAmount)

			installment := types.RepaymentInstallment{
				InstallmentNumber: i,
				DueDate:          ctx.BlockTime().AddDate(0, int(i), 0),
				TotalAmount:      sdk.NewCoin(loanAmount.Denom, monthlyPayment.TruncateInt()),
				PrincipalAmount:  sdk.NewCoin(loanAmount.Denom, monthlyPrincipalAmount.TruncateInt()),
				InterestAmount:   sdk.NewCoin(loanAmount.Denom, monthlyInterestAmount.TruncateInt()),
				Status:           types.InstallmentStatusPending,
			}
			schedule = append(schedule, installment)
		}
	}

	return schedule
}

// Helper functions for loan processing
func (lp *LoanProcessor) rejectLoanApplication(ctx sdk.Context, application types.LoanApplication, reason string) (*AgriculturalLoan, error) {
	application.Status = types.ApplicationStatusRejected
	application.RejectionReason = reason
	lp.keeper.SetLoanApplication(ctx, application)

	// Emit rejection event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLoanRejected,
			sdk.NewAttribute(types.AttributeKeyApplicationID, application.ID),
			sdk.NewAttribute(types.AttributeKeyFarmerID, application.FarmerID),
			sdk.NewAttribute(types.AttributeKeyRejectionReason, reason),
		),
	)

	return nil, fmt.Errorf("loan application rejected: %s", reason)
}

func (lp *LoanProcessor) applyDiscounts(ctx sdk.Context, farmerID string, baseRate sdk.Dec) sdk.Dec {
	// Apply festival discounts
	festivalOffers := lp.keeper.GetActiveFestivalOffers(ctx)
	for _, offer := range festivalOffers {
		discount, _ := sdk.NewDecFromStr(offer.InterestReduction)
		baseRate = baseRate.Sub(discount)
	}

	// Apply farmer-specific discounts
	farmerProfile, found := lp.keeper.GetFarmerProfile(ctx, farmerID)
	if found {
		if farmerProfile.IsWomenFarmer {
			params := lp.keeper.GetParams(ctx)
			womenDiscount, _ := sdk.NewDecFromStr(params.WomenFarmerDiscount)
			baseRate = baseRate.Sub(womenDiscount)
		}

		if farmerProfile.TotalLandArea.LT(sdk.NewDec(5)) { // Small farmer
			params := lp.keeper.GetParams(ctx)
			smallFarmerDiscount, _ := sdk.NewDecFromStr(params.SmallFarmerDiscount)
			baseRate = baseRate.Sub(smallFarmerDiscount)
		}
	}

	// Ensure rate doesn't go below minimum
	params := lp.keeper.GetParams(ctx)
	minRate, _ := sdk.NewDecFromStr(params.MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}

	return baseRate
}

func (lp *LoanProcessor) assessCollateralRequirements(ctx sdk.Context, loanAmount sdk.Coin, riskCategory string) (bool, types.CollateralInfo) {
	params := lp.keeper.GetParams(ctx)
	
	// High-risk borrowers may need collateral
	if riskCategory == "HIGH_RISK" || riskCategory == "VERY_HIGH_RISK" {
		return true, types.CollateralInfo{
			CollateralType:  "LAND_DOCUMENTS",
			CollateralValue: loanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(12, 1)), // 120% of loan
			Required:        true,
		}
	}

	// Large loans may need collateral
	if loanAmount.Amount.GT(params.CollateralThreshold.Amount) {
		return true, types.CollateralInfo{
			CollateralType:  "LAND_DOCUMENTS_OR_GOLD",
			CollateralValue: loanAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(11, 1)), // 110% of loan
			Required:        true,
		}
	}

	return false, types.CollateralInfo{}
}

func (lp *LoanProcessor) isInsuranceRequired(ctx sdk.Context, cropType string, loanAmount sdk.Coin) bool {
	params := lp.keeper.GetParams(ctx)
	
	// Insurance mandatory for loans above threshold
	if loanAmount.Amount.GT(params.InsuranceThreshold.Amount) {
		return true
	}

	// Insurance mandatory for weather-sensitive crops
	weatherSensitiveCrops := []string{"COTTON", "SUGARCANE", "VEGETABLES", "WHEAT"}
	for _, crop := range weatherSensitiveCrops {
		if crop == cropType {
			return true
		}
	}

	return false
}

func (lp *LoanProcessor) generateLoanID(ctx sdk.Context, farmerID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("KM-%s-%d", farmerID[:8], timestamp)
}

func (lp *LoanProcessor) isCropLoan(cropType string) bool {
	return cropType != "" && cropType != "EQUIPMENT" && cropType != "INFRASTRUCTURE"
}

// Additional helper methods would include:
// - calculateOutstandingAmount
// - allocateRepayment
// - determineRepaymentType
// - processLoanClosure
// - updateFarmerProfileAfterDisbursement
// - updateFarmerProfileAfterRepayment
// - createMandatoryInsurance
// - generateApprovalConditions
// - calculateDisbursementFee
// - distributeDisbursementFee
// - isCollateralVerified