package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/deshchain/x/krishimitra/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// ApplyLoan handles loan application from farmers with REVOLUTIONARY member-only restrictions
func (k msgServer) ApplyLoan(goCtx context.Context, msg *types.MsgApplyLoan) (*types.MsgApplyLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// REVOLUTIONARY RESTRICTION: Verify borrower is pool member before processing
	applicantAddr, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid applicant address")
	}

	// Check if liquidity is available and borrower is eligible
	amount, err := sdk.NewDecFromStr(msg.LoanAmount.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid loan amount")
	}

	canProcess, message := k.liquidityKeeper.CanProcessLoan(ctx, amount, "krishimitra", applicantAddr)
	if !canProcess {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "🚫 REVOLUTIONARY LENDING RESTRICTION: %s", message)
	}

	// Check if applicant already has an active loan
	if k.HasActiveLoan(ctx, msg.Applicant) {
		return nil, sdkerrors.Wrap(types.ErrActiveLoanExists, "applicant already has an active loan")
	}

	// Verify DhanPata ID
	if !k.accountKeeper.HasDhanPataID(ctx, msg.Applicant, msg.DhanPataID) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDhanPata, "DhanPata ID does not match applicant")
	}

	// Check eligibility based on PIN code
	eligible, reason := k.CheckPINCodeEligibility(ctx, msg.Pincode, msg.LoanAmount)
	if !eligible {
		return nil, sdkerrors.Wrapf(types.ErrNotEligible, "PIN code eligibility: %s", reason)
	}

	// Calculate interest rate based on various factors
	interestRate := k.CalculateInterestRate(ctx, msg.CropType, msg.Duration, msg.KCCNumber != "")

	// Check if festival bonus applies
	festivalBonus := k.GetActiveFestivalBonus(ctx)
	if festivalBonus != nil {
		interestRate = interestRate.Sub(festivalBonus.InterestReduction)
		if interestRate.LT(sdk.NewDecWithPrec(6, 2)) {
			interestRate = sdk.NewDecWithPrec(6, 2) // Minimum 6%
		}
	}

	// Create loan application
	loanID := k.GenerateLoanID(ctx)
	loan := types.AgriculturalLoan{
		LoanID:           loanID,
		Borrower:         msg.Applicant,
		DhanPataID:       msg.DhanPataID,
		Amount:           msg.LoanAmount,
		RequestedAmount:  msg.LoanAmount,
		InterestRate:     interestRate,
		Duration:         msg.Duration,
		Status:           types.LoanStatusPending,
		ApplicationDate:  ctx.BlockTime(),
		CropType:         msg.CropType,
		LandArea:         msg.LandArea,
		LandOwnership:    msg.LandOwnershipDoc,
		Pincode:          msg.Pincode,
		Purpose:          msg.Purpose,
		KCCNumber:        msg.KCCNumber,
		CulturalQuote:    msg.CulturalQuote,
	}

	// Store loan application
	k.SetLoan(ctx, loan)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanApplied,
			sdk.NewAttribute(types.AttributeKeyLoanID, loanID),
			sdk.NewAttribute(types.AttributeKeyBorrower, msg.Applicant),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.LoanAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCropType, string(msg.CropType)),
			sdk.NewAttribute(types.AttributeKeyInterestRate, interestRate.String()),
		),
	})

	return &types.MsgApplyLoanResponse{
		LoanID:       loanID,
		InterestRate: interestRate,
		Status:       "Application submitted successfully",
	}, nil
}

// ApproveLoan handles loan approval by authorized validators
func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if approver is authorized
	if !k.IsAuthorizedApprover(ctx, msg.Approver) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to approve loans")
	}

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan application not found")
	}

	if loan.Status != types.LoanStatusPending {
		return nil, sdkerrors.Wrapf(types.ErrInvalidLoanStatus, "loan status is %s, expected pending", loan.Status)
	}

	// Update loan details
	loan.Amount = msg.ApprovedAmount
	loan.InterestRate = msg.InterestRate
	loan.Status = types.LoanStatusApproved
	loan.ApprovalDate = &ctx.BlockTime()
	loan.RepaymentSchedule = msg.RepaymentSchedule
	loan.Remarks = msg.Remarks

	// Calculate total repayment amount
	interest := msg.ApprovedAmount.Amount.ToDec().Mul(msg.InterestRate).Mul(sdk.NewDec(loan.Duration)).Quo(sdk.NewDec(12))
	totalRepayment := msg.ApprovedAmount.Amount.Add(interest.TruncateInt())
	loan.TotalRepayment = sdk.NewCoin(msg.ApprovedAmount.Denom, totalRepayment)

	// Update loan
	k.SetLoan(ctx, loan)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanApproved,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeyApprover, msg.Approver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.ApprovedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTotalRepayment, loan.TotalRepayment.String()),
		),
	})

	return &types.MsgApproveLoanResponse{
		Success:        true,
		TotalRepayment: loan.TotalRepayment,
	}, nil
}

// DisburseLoan handles loan disbursement
func (k msgServer) DisburseLoan(goCtx context.Context, msg *types.MsgDisburseLoanRequest) (*types.MsgDisburseLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if disburser is authorized
	if !k.IsAuthorizedDisburser(ctx, msg.Disburser) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to disburse loans")
	}

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	if loan.Status != types.LoanStatusApproved {
		return nil, sdkerrors.Wrapf(types.ErrInvalidLoanStatus, "loan status is %s, expected approved", loan.Status)
	}

	// Get borrower address
	borrowerAddr, err := sdk.AccAddressFromBech32(loan.Borrower)
	if err != nil {
		return nil, err
	}

	// REVOLUTIONARY FEE STRUCTURE: Calculate processing fee and disbursement amount
	approvedAmount, err := sdk.NewDecFromStr(loan.Amount.Amount.String())
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid loan amount")
	}

	disbursedAmount, processingFee := k.liquidityKeeper.CalculateDisbursementAmount(ctx, approvedAmount)
	
	// Convert back to sdk.Coin for transfers
	disbursedCoin := sdk.NewCoin(loan.Amount.Denom, disbursedAmount.TruncateInt())
	processingFeeCoin := sdk.NewCoin(loan.Amount.Denom, processingFee.TruncateInt())

	// Transfer 99% of approved amount to borrower (after processing fee)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrowerAddr, sdk.NewCoins(disbursedCoin))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to disburse loan")
	}

	// Keep processing fee in module for platform revenue (automatic)
	// (No transfer needed - fee stays in lending pool for platform sustainability)

	// Update loan status with fee breakdown
	loan.Status = types.LoanStatusActive
	loan.DisbursedAmount = disbursedCoin
	loan.ProcessingFee = processingFeeCoin
	loan.DisbursementDate = &ctx.BlockTime()
	loan.NextPaymentDate = k.CalculateNextPaymentDate(ctx, loan.RepaymentSchedule)
	k.SetLoan(ctx, loan)

	// Create farmer profile if not exists
	profile, found := k.GetFarmerProfile(ctx, loan.Borrower)
	if !found {
		profile = types.FarmerProfile{
			Address:      loan.Borrower,
			DhanPataID:   loan.DhanPataID,
			Pincode:      loan.Pincode,
			TotalLoans:   1,
			ActiveLoans:  1,
			CreditScore:  750, // Default credit score
			JoinedDate:   ctx.BlockTime(),
		}
	} else {
		profile.TotalLoans++
		profile.ActiveLoans++
	}
	k.SetFarmerProfile(ctx, profile)

	// REVOLUTIONARY TRANSPARENCY: Emit comprehensive event with fee breakdown
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanDisbursed,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeyBorrower, loan.Borrower),
			sdk.NewAttribute(types.AttributeKeyApprovedAmount, loan.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDisbursedAmount, disbursedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyProcessingFee, processingFeeCoin.String()),
			sdk.NewAttribute("effective_fee_rate", fmt.Sprintf("%.2f%%", processingFee.Quo(approvedAmount).Mul(sdk.NewDec(100)).MustFloat64())),
			sdk.NewAttribute("savings_vs_banks", "🎯 50-70% cheaper processing fee than traditional banks"),
			sdk.NewAttribute("transparency", "💎 Complete fee breakdown - no hidden charges"),
		),
	})

	return &types.MsgDisburseLoanResponse{
		Success:           true,
		ApprovedAmount:    loan.Amount.String(),
		DisbursedAmount:   disbursedCoin.String(),
		ProcessingFee:     processingFeeCoin.String(),
		EffectiveFeeRate:  fmt.Sprintf("%.2f%%", processingFee.Quo(approvedAmount).Mul(sdk.NewDec(100)).MustFloat64()),
		Transaction:       ctx.TxBytes().String(),
		BorrowerAdvantage: "🏆 Revolutionary 1% processing fee (capped ₹2500) vs banks 2-5%",
	}, nil
}

// RepayLoan handles loan repayment
func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Verify borrower
	if loan.Borrower != msg.Borrower {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the borrower of this loan")
	}

	if loan.Status != types.LoanStatusActive {
		return nil, sdkerrors.Wrapf(types.ErrInvalidLoanStatus, "loan status is %s, expected active", loan.Status)
	}

	// Get borrower address
	borrowerAddr, err := sdk.AccAddressFromBech32(msg.Borrower)
	if err != nil {
		return nil, err
	}

	// Check if early repayment bonus applies
	earlyBonus := sdk.ZeroDec()
	if k.IsEarlyRepayment(ctx, loan) {
		earlyBonus = sdk.NewDecWithPrec(2, 2) // 2% discount
	}

	// Calculate actual repayment amount after bonuses
	repayAmount := msg.RepayAmount
	if !earlyBonus.IsZero() {
		discountAmount := repayAmount.Amount.ToDec().Mul(earlyBonus).TruncateInt()
		repayAmount.Amount = repayAmount.Amount.Sub(discountAmount)
	}

	// Transfer repayment to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrowerAddr, types.ModuleName, sdk.NewCoins(repayAmount))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to process repayment")
	}

	// Update loan
	loan.RepaidAmount = loan.RepaidAmount.Add(repayAmount)
	loan.LastPaymentDate = &ctx.BlockTime()

	// Check if fully repaid
	if loan.RepaidAmount.IsGTE(loan.TotalRepayment) {
		loan.Status = types.LoanStatusRepaid
		loan.CompletionDate = &ctx.BlockTime()
		
		// Update farmer profile
		profile, _ := k.GetFarmerProfile(ctx, loan.Borrower)
		profile.ActiveLoans--
		profile.RepaidLoans++
		profile.CreditScore += 10 // Improve credit score
		k.SetFarmerProfile(ctx, profile)
	} else {
		// Update next payment date
		loan.NextPaymentDate = k.CalculateNextPaymentDate(ctx, loan.RepaymentSchedule)
	}

	k.SetLoan(ctx, loan)

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanRepaid,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute(types.AttributeKeyBorrower, msg.Borrower),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.RepayAmount.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingAmount, loan.TotalRepayment.Sub(loan.RepaidAmount).String()),
			sdk.NewAttribute("early_bonus", earlyBonus.String()),
		),
	})

	return &types.MsgRepayLoanResponse{
		Success:         true,
		RemainingAmount: loan.TotalRepayment.Sub(loan.RepaidAmount),
		LoanStatus:      string(loan.Status),
	}, nil
}

// UpdateCropData handles crop yield and market data updates
func (k msgServer) UpdateCropData(goCtx context.Context, msg *types.MsgUpdateCropData) (*types.MsgUpdateCropDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if authority is authorized
	if !k.IsAuthorizedDataProvider(ctx, msg.Authority) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to update crop data")
	}

	// Get loan
	loan, found := k.GetLoan(ctx, msg.LoanID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrLoanNotFound, "loan not found")
	}

	// Create crop data entry
	cropData := types.CropData{
		LoanID:       msg.LoanID,
		CropYield:    msg.CropYield,
		MarketPrice:  msg.MarketPrice,
		WeatherEvent: msg.WeatherEvent,
		UpdateDate:   ctx.BlockTime(),
		Remarks:      msg.Remarks,
	}

	// Store crop data
	k.SetCropData(ctx, loan.Borrower, cropData)

	// Check if weather event qualifies for insurance
	if msg.WeatherEvent != "" && k.IsInsurableEvent(ctx, msg.WeatherEvent) {
		// Auto-trigger insurance claim process
		k.CreateInsuranceClaim(ctx, loan, msg.WeatherEvent)
	}

	// Emit event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCropDataUpdated,
			sdk.NewAttribute(types.AttributeKeyLoanID, msg.LoanID),
			sdk.NewAttribute("crop_yield", msg.CropYield),
			sdk.NewAttribute("market_price", msg.MarketPrice.String()),
			sdk.NewAttribute("weather_event", msg.WeatherEvent),
		),
	})

	return &types.MsgUpdateCropDataResponse{
		Success:        true,
		InsuranceClaim: msg.WeatherEvent != "",
	}, nil
}