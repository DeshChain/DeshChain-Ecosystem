package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
)

// Additional utility functions for business loan processing

// calculateProcessingFee calculates processing fee for business loans
func (blp *BusinessLoanProcessor) calculateProcessingFee(ctx sdk.Context, loanAmount sdk.Coin, loanType string) sdk.Coin {
	params := blp.keeper.GetParams(ctx)
	baseFeeRate := params.ProcessingFeeRate // e.g., 1%
	
	// Different processing fees for different loan types
	switch loanType {
	case "INVOICE_FINANCING":
		baseFeeRate = baseFeeRate.Mul(sdk.NewDecWithPrec(5, 1)) // 50% of base rate
	case "EQUIPMENT":
		baseFeeRate = baseFeeRate.Mul(sdk.NewDecWithPrec(15, 1)) // 150% of base rate
	}
	
	feeAmount := loanAmount.Amount.ToDec().Mul(baseFeeRate).TruncateInt()
	return sdk.NewCoin(loanAmount.Denom, feeAmount)
}

// calculateBusinessDisbursementFee calculates disbursement fee
func (blp *BusinessLoanProcessor) calculateBusinessDisbursementFee(ctx sdk.Context, loanAmount sdk.Coin) sdk.Coin {
	// Fixed disbursement fee
	feeAmount := sdk.NewInt(500) // ₹500 fixed fee
	return sdk.NewCoin(loanAmount.Denom, feeAmount)
}

// determineDisbursementMode determines how funds should be disbursed
func (blp *BusinessLoanProcessor) determineDisbursementMode(ctx sdk.Context, loanType string) string {
	switch loanType {
	case "EQUIPMENT":
		return "EQUIPMENT_VENDOR"
	case "INVOICE_FINANCING":
		return "INVOICE_FINANCING"
	default:
		return "DIRECT_TRANSFER"
	}
}

// determineApprovalLevel determines approval authority level
func (blp *BusinessLoanProcessor) determineApprovalLevel(ctx sdk.Context, loanAmount sdk.Coin) string {
	if loanAmount.Amount.GT(sdk.NewInt(10000000)) { // > 1 Cr
		return "SENIOR_MANAGEMENT"
	} else if loanAmount.Amount.GT(sdk.NewInt(5000000)) { // > 50 lakh
		return "MIDDLE_MANAGEMENT"
	} else {
		return "AUTOMATED"
	}
}

// generateApprovalConditions generates loan approval conditions
func (blp *BusinessLoanProcessor) generateApprovalConditions(ctx sdk.Context, collateralRequired, guaranteeRequired bool) []string {
	conditions := []string{
		"Submission of post-dated cheques",
		"Insurance coverage as per policy",
		"Compliance with loan covenants",
	}
	
	if collateralRequired {
		conditions = append(conditions, "Collateral verification and registration")
	}
	if guaranteeRequired {
		conditions = append(conditions, "Personal/Corporate guarantee execution")
	}
	
	return conditions
}

// calculateBusinessOutstandingAmount calculates total outstanding amount
func (blp *BusinessLoanProcessor) calculateBusinessOutstandingAmount(ctx sdk.Context, loan ProcessedBusinessLoan) sdk.Coin {
	totalRepaid := sdk.ZeroInt()
	for _, repayment := range loan.RepaymentHistory {
		totalRepaid = totalRepaid.Add(repayment.TotalAmount.Amount)
	}
	
	// Calculate total amount due (principal + interest)
	totalInterest := sdk.ZeroDec()
	for _, installment := range loan.RepaymentSchedule {
		totalInterest = totalInterest.Add(installment.InterestAmount.Amount.ToDec())
	}
	
	totalDue := loan.LoanAmount.Amount.Add(totalInterest.TruncateInt())
	outstanding := totalDue.Sub(totalRepaid)
	
	return sdk.NewCoin(loan.LoanAmount.Denom, outstanding)
}

// allocateBusinessRepayment allocates repayment between principal, interest, and fees
func (blp *BusinessLoanProcessor) allocateBusinessRepayment(ctx sdk.Context, loan ProcessedBusinessLoan, repaymentAmount sdk.Coin) (sdk.Coin, sdk.Coin, sdk.Coin) {
	// First pay any outstanding fees
	outstandingFees := blp.calculateOutstandingFees(ctx, loan)
	feesPortion := repaymentAmount
	if repaymentAmount.Amount.GT(outstandingFees.Amount) {
		feesPortion = outstandingFees
	}
	
	remainingAmount := repaymentAmount.Sub(feesPortion)
	
	// Then pay interest
	outstandingInterest := blp.calculateOutstandingInterest(ctx, loan)
	interestPortion := remainingAmount
	if remainingAmount.Amount.GT(outstandingInterest.Amount) {
		interestPortion = outstandingInterest
	}
	
	remainingForPrincipal := remainingAmount.Sub(interestPortion)
	
	// Rest goes to principal
	principalPortion := remainingForPrincipal
	
	return principalPortion, interestPortion, feesPortion
}

// generateBusinessRepaymentID generates unique repayment ID
func (blp *BusinessLoanProcessor) generateBusinessRepaymentID(ctx sdk.Context, loanID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("REP-%s-%d", loanID[4:12], timestamp)
}

// Utility functions for outstanding calculations
func (blp *BusinessLoanProcessor) calculateOutstandingFees(ctx sdk.Context, loan ProcessedBusinessLoan) sdk.Coin {
	// Calculate any late fees, penalty charges, etc.
	return sdk.NewCoin(loan.LoanAmount.Denom, sdk.ZeroInt())
}

func (blp *BusinessLoanProcessor) calculateOutstandingInterest(ctx sdk.Context, loan ProcessedBusinessLoan) sdk.Coin {
	totalInterestDue := sdk.ZeroInt()
	totalInterestPaid := sdk.ZeroInt()
	
	// Calculate total interest from schedule
	for _, installment := range loan.RepaymentSchedule {
		totalInterestDue = totalInterestDue.Add(installment.InterestAmount.Amount)
	}
	
	// Calculate total interest paid
	for _, repayment := range loan.RepaymentHistory {
		totalInterestPaid = totalInterestPaid.Add(repayment.InterestAmount.Amount)
	}
	
	outstanding := totalInterestDue.Sub(totalInterestPaid)
	return sdk.NewCoin(loan.LoanAmount.Denom, outstanding)
}

// Verification functions
func (blp *BusinessLoanProcessor) isCollateralVerified(ctx sdk.Context, loanID string) bool {
	// Check if collateral has been verified and registered
	return true // Simplified for now
}

func (blp *BusinessLoanProcessor) isGuaranteeVerified(ctx sdk.Context, loanID string) bool {
	// Check if guarantee documents have been executed
	return true // Simplified for now
}

func (blp *BusinessLoanProcessor) validateFinalCompliance(ctx sdk.Context, loan ProcessedBusinessLoan) bool {
	// Final compliance validation before disbursement
	return loan.ComplianceChecks.OverallStatus == "PASSED"
}

// Processing functions for different disbursement modes
func (blp *BusinessLoanProcessor) processInvoiceFinancingDisbursement(ctx sdk.Context, loan ProcessedBusinessLoan) error {
	// Process invoice financing specific disbursement
	return nil // Implementation depends on invoice verification system
}

func (blp *BusinessLoanProcessor) processVendorPayment(ctx sdk.Context, loan ProcessedBusinessLoan, amount sdk.Coin) error {
	// Process payment to equipment vendor
	return nil // Implementation depends on vendor payment system
}

// Additional business workflow functions
func (blp *BusinessLoanProcessor) distributeDisbursementFees(ctx sdk.Context, fees sdk.Coin, loanType string) error {
	// Distribute fees to revenue, operational pools
	return nil // Implementation based on fee distribution policy
}

func (blp *BusinessLoanProcessor) updateBusinessProfileAfterDisbursement(ctx sdk.Context, businessID string, loan ProcessedBusinessLoan) {
	// Update business profile with new active loan
	profile, found := blp.keeper.GetBusinessProfile(ctx, businessID)
	if found {
		profile.ActiveLoans++
		profile.TotalBorrowed = profile.TotalBorrowed.Add(loan.LoanAmount)
		blp.keeper.SetBusinessProfile(ctx, profile)
	}
}

func (blp *BusinessLoanProcessor) updateBusinessProfileAfterRepayment(ctx sdk.Context, businessID string, repayment types.BusinessRepayment) {
	// Update business profile after repayment
	profile, found := blp.keeper.GetBusinessProfile(ctx, businessID)
	if found {
		profile.TotalRepaid = profile.TotalRepaid.Add(repayment.TotalAmount)
		profile.LastRepaymentDate = repayment.RepaymentDate
		blp.keeper.SetBusinessProfile(ctx, profile)
	}
}

func (blp *BusinessLoanProcessor) processBusinessLoanClosure(ctx sdk.Context, loan ProcessedBusinessLoan) error {
	// Process loan closure activities
	profile, found := blp.keeper.GetBusinessProfile(ctx, loan.BusinessID)
	if found {
		profile.ActiveLoans--
		profile.CompletedLoans++
		blp.keeper.SetBusinessProfile(ctx, profile)
	}
	return nil
}

// Credit line management functions
func (blp *BusinessLoanProcessor) activateCreditLine(ctx sdk.Context, businessID string, creditLine types.CreditLineInfo) error {
	// Activate credit line for working capital access
	creditLine.IsActive = true
	blp.keeper.SetCreditLine(ctx, types.CreditLine{
		ID:           fmt.Sprintf("CL-%s-%d", businessID[:8], ctx.BlockTime().Unix()),
		BusinessID:   businessID,
		CreditLimit:  creditLine.CreditLimit,
		Utilized:     sdk.ZeroInt(),
		Available:    creditLine.CreditLimit.Amount,
		InterestRate: creditLine.InterestRate.String(),
		IsActive:     true,
		CreatedAt:    ctx.BlockTime(),
	})
	return nil
}

func (blp *BusinessLoanProcessor) processCreditLineRepayment(ctx sdk.Context, businessID string, amount sdk.Coin) error {
	// Process credit line repayment and increase available limit
	creditLine, found := blp.keeper.GetCreditLine(ctx, businessID)
	if found {
		creditLine.Utilized = creditLine.Utilized.Sub(amount.Amount)
		creditLine.Available = creditLine.Available.Add(amount.Amount)
		blp.keeper.SetCreditLine(ctx, creditLine)
	}
	return nil
}