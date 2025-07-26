package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// ApplyBusinessLoan handles business loan applications with REVOLUTIONARY member-only restrictions
func (k msgServer) ApplyBusinessLoan(goCtx context.Context, msg *types.MsgApplyBusinessLoan) (*types.MsgApplyBusinessLoanResponse, error) {
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

	canProcess, message := k.liquidityKeeper.CanProcessLoan(ctx, amount, "vyavasayamitra", applicantAddr)
	if !canProcess {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "🚫 REVOLUTIONARY LENDING RESTRICTION: %s", message)
	}

	// Check if business already has too many active loans
	activeLoans := k.GetActiveLoansCount(ctx, msg.Applicant)
	if activeLoans >= 3 {
		return nil, sdkerrors.Wrap(types.ErrTooManyActiveLoans, "maximum 3 active loans allowed")
	}

	// Verify DhanPata ID
	if !k.accountKeeper.HasDhanPataID(ctx, msg.Applicant, msg.DhanPataID) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDhanPata, "DhanPata ID does not match applicant")
	}

	// Verify GST number
	if !k.IsValidGSTNumber(msg.GSTNumber) {
		return nil, sdkerrors.Wrap(types.ErrInvalidGST, "invalid GST number format")
	}

	// Check eligibility based on business metrics
	eligible, reason := k.CheckBusinessEligibility(ctx, msg)
	if !eligible {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "eligibility check failed: %s", reason)
	}

	// Calculate credit score based on various factors
	creditScore := k.CalculateBusinessCreditScore(ctx, msg)

	// Calculate interest rate based on business type, credit score, and other factors
	baseRate := k.CalculateBaseInterestRate(ctx, msg.BusinessType, creditScore, msg.Duration)

	// Apply discounts
	if msg.BusinessType == types.BusinessType_BUSINESS_TYPE_STARTUP {
		startupDiscount := sdk.MustNewDecFromStr(k.GetParams(ctx).StartupDiscount)
		baseRate = baseRate.Sub(startupDiscount)
	}

	// Check for women entrepreneur discount
	if k.IsWomenOwnedBusiness(ctx, msg.Applicant) {
		womenDiscount := sdk.MustNewDecFromStr(k.GetParams(ctx).WomenEntrepreneurDiscount)
		baseRate = baseRate.Sub(womenDiscount)
	}

	// Apply festival discount if active
	festivalOffer := k.GetActiveFestivalOffer(ctx)
	if festivalOffer != nil && k.IsEligibleForFestivalOffer(msg, festivalOffer) {
		festivalDiscount := sdk.MustNewDecFromStr(festivalOffer.InterestReduction)
		baseRate = baseRate.Sub(festivalDiscount)
	}

	// Ensure rate is within bounds
	minRate := sdk.MustNewDecFromStr(k.GetParams(ctx).MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}

	// Create loan application
	loanID := k.GenerateLoanID(ctx)
	loan := types.BusinessLoan{
		LoanID:                  loanID,
		Borrower:                msg.Applicant,
		DhanPataID:              msg.DhanPataID,
		BusinessName:            msg.BusinessName,
		BusinessType:            msg.BusinessType,
		GSTNumber:               msg.GSTNumber,
		LoanAmount:              msg.LoanAmount,
		ApprovedAmount:          sdk.NewCoin(msg.LoanAmount.Denom, sdk.ZeroInt()),
		LoanPurpose:             msg.LoanPurpose,
		InterestRate:            baseRate.String(),
		Duration:                msg.Duration,
		Status:                  types.LoanStatus_LOAN_STATUS_PENDING,
		ApplicationDate:         ctx.BlockTime(),
		AnnualRevenue:           msg.AnnualRevenue,
		EmployeeCount:           msg.EmployeeCount,
		BusinessAge:             msg.BusinessAge,
		Pincode:                 msg.Pincode,
		CollateralDetails:       msg.CollateralOffered,
		BusinessPlanHash:        msg.BusinessPlan,
		FinancialStatementsHash: msg.FinancialStatements,
		CulturalQuote:           msg.CulturalQuote,
	}

	// Store loan application
	k.SetLoan(ctx, loan)

	// Update or create business profile
	k.UpdateBusinessProfile(ctx, msg)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBusinessLoanApplied,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyBorrower, msg.Applicant),
			sdk.NewAttribute(types.AttributeKeyBusinessName, msg.BusinessName),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.LoanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, baseRate.String()),
			sdk.NewAttribute(types.AttributeKeyCreditScore, fmt.Sprintf("%d", creditScore)),
		),
	})

	return &types.MsgApplyBusinessLoanResponse{
		LoanID:       loanID,
		InterestRate: baseRate.String(),
		Status:       "Application submitted successfully",
		CreditScore:  fmt.Sprintf("%d", creditScore),
	}, nil
}

// ApplyInvoiceFinancing handles invoice financing applications with REVOLUTIONARY member-only restrictions
func (k msgServer) ApplyInvoiceFinancing(goCtx context.Context, msg *types.MsgApplyInvoiceFinancing) (*types.MsgApplyInvoiceFinancingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// REVOLUTIONARY RESTRICTION: Verify borrower is pool member before processing
	applicantAddr, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid applicant address")
	}

	// Check if member is eligible for invoice financing
	if !k.liquidityKeeper.IsPoolMember(ctx, applicantAddr) {
		return nil, sdkerrors.Wrap(types.ErrNotEligible, "🚫 REVOLUTIONARY FINANCING RESTRICTION: Only Suraksha Pool members can access invoice financing!")
	}

	// Get the associated business loan
	loan, found := k.GetLoan(ctx, msg.BusinessLoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "business loan not found")
	}

	// Verify applicant owns the loan
	if loan.Borrower != msg.Applicant {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the owner of this business loan")
	}

	// Check if loan has active credit line
	if loan.CreditLine == nil || !loan.CreditLine.IsActive {
		return nil, sdkerrors.Wrap(types.ErrNoCreditLine, "no active credit line available")
	}

	// Validate invoice details
	if !k.IsValidInvoice(ctx, msg) {
		return nil, sdkerrors.Wrap(types.ErrInvalidInvoice, "invoice validation failed")
	}

	// Calculate financing amount based on LTV
	ltv := sdk.MustNewDecFromStr(k.GetParams(ctx).InvoiceFinancingLTV)
	financingPercent := sdk.MustNewDecFromStr(msg.FinancingPercent)
	actualLTV := sdk.MinDec(ltv, financingPercent)
	financedAmount := msg.InvoiceAmount.Amount.ToDec().Mul(actualLTV).TruncateInt()

	// Check credit line availability
	if loan.CreditLine.AvailableCredit.IsLT(sdk.NewCoin(msg.InvoiceAmount.Denom, financedAmount)) {
		return nil, sdkerrors.Wrap(types.ErrInsufficientCreditLine, "insufficient credit line available")
	}

	// Calculate interest rate for invoice financing (usually lower than term loan)
	invoiceRate := k.CalculateInvoiceFinancingRate(ctx, loan, msg)

	// Calculate expected repayment date
	dueDate, _ := time.Parse("2006-01-02", msg.DueDate)
	expectedRepayment := dueDate.Add(-5 * 24 * time.Hour) // 5 days before due date

	// Create invoice financing record
	financingID := k.GenerateFinancingID(ctx)
	invoiceFinancing := types.InvoiceFinancing{
		FinancingID:          financingID,
		LoanID:               msg.BusinessLoanID,
		InvoiceNumber:        msg.InvoiceNumber,
		InvoiceAmount:        msg.InvoiceAmount,
		FinancedAmount:       sdk.NewCoin(msg.InvoiceAmount.Denom, financedAmount),
		InvoiceDate:          ctx.BlockTime(),
		DueDate:              dueDate,
		BuyerGSTNumber:       msg.BuyerGSTNumber,
		BuyerName:            msg.BuyerName,
		InvoiceDocumentHash:  msg.InvoiceDocument,
		Status:               "APPROVED",
		FinancingDate:        &ctx.BlockTime(),
	}

	// Update credit line utilization
	loan.CreditLine.UtilizedAmount = loan.CreditLine.UtilizedAmount.Add(sdk.NewCoin(msg.InvoiceAmount.Denom, financedAmount))
	loan.CreditLine.AvailableCredit = loan.CreditLine.CreditLimit.Sub(loan.CreditLine.UtilizedAmount)
	loan.InvoiceFinancings = append(loan.InvoiceFinancings, &invoiceFinancing)
	k.SetLoan(ctx, loan)

	// Disburse funds immediately for invoice financing
	borrowerAddr, _ := sdk.AccAddressFromBech32(msg.Applicant)
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrowerAddr, sdk.NewCoins(invoiceFinancing.FinancedAmount))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to disburse invoice financing")
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInvoiceFinancingApproved,
			sdk.NewAttribute(types.AttributeKeyFinancingID, financingID),
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.BusinessLoanID),
			sdk.NewAttribute(types.AttributeKeyInvoiceNumber, msg.InvoiceNumber),
			sdk.NewAttribute(types.AttributeKeyAmount, invoiceFinancing.FinancedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyInterestRate, invoiceRate.String()),
		),
	})

	return &types.MsgApplyInvoiceFinancingResponse{
		FinancingID:           financingID,
		ApprovedAmount:        invoiceFinancing.FinancedAmount,
		InterestRate:          invoiceRate.String(),
		ExpectedRepaymentDate: expectedRepayment.Format("2006-01-02"),
	}, nil
}

// UpdateBusinessMetrics handles business performance updates
func (k msgServer) UpdateBusinessMetrics(goCtx context.Context, msg *types.MsgUpdateBusinessMetrics) (*types.MsgUpdateBusinessMetricsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority
	if !k.IsAuthorizedMetricsProvider(ctx, msg.Authority) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to update business metrics")
	}

	// Get loan
	loan, found := k.GetLoan(ctx, msg.BusinessLoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Create business metrics record
	netProfit := msg.MonthlyRevenue.Sub(msg.MonthlyExpenses)
	metrics := types.BusinessMetrics{
		LoanID:            msg.BusinessLoanID,
		ReportingPeriod:   msg.ReportingPeriod,
		MonthlyRevenue:    msg.MonthlyRevenue,
		MonthlyExpenses:   msg.MonthlyExpenses,
		NetProfit:         netProfit,
		EmployeeCount:     msg.NewEmployees,
		GSTFilingStatus:   msg.GSTFilingStatus,
		CreditUtilization: msg.CreditUtilization.String(),
		UpdateDate:        ctx.BlockTime(),
	}

	// Calculate growth rate
	previousMetrics := k.GetPreviousMonthMetrics(ctx, msg.BusinessLoanID)
	if previousMetrics != nil {
		growthRate := msg.MonthlyRevenue.Amount.Sub(previousMetrics.MonthlyRevenue.Amount).ToDec().
			Quo(previousMetrics.MonthlyRevenue.Amount.ToDec()).Mul(sdk.NewDec(100))
		metrics.GrowthRate = growthRate.String() + "%"
	}

	// Store metrics
	k.SetBusinessMetrics(ctx, metrics)

	// Update credit score based on performance
	profile, _ := k.GetBusinessProfile(ctx, loan.Borrower)
	newCreditScore := k.RecalculateCreditScore(ctx, profile, metrics)
	profile.CreditScore = newCreditScore
	k.SetBusinessProfile(ctx, profile)

	// Check if eligible for credit line increase
	eligibleForIncrease := k.IsEligibleForCreditLineIncrease(ctx, loan, metrics)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBusinessMetricsUpdated,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.BusinessLoanID),
			sdk.NewAttribute(types.AttributeKeyReportingPeriod, msg.ReportingPeriod),
			sdk.NewAttribute("revenue", msg.MonthlyRevenue.String()),
			sdk.NewAttribute("profit", netProfit.String()),
			sdk.NewAttribute("growth_rate", metrics.GrowthRate),
			sdk.NewAttribute("credit_score", fmt.Sprintf("%d", newCreditScore)),
		),
	})

	return &types.MsgUpdateBusinessMetricsResponse{
		Success:                      true,
		UpdatedCreditScore:           fmt.Sprintf("%d", newCreditScore),
		CreditLineIncreaseEligible:   eligibleForIncrease,
	}, nil
}