package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types
const (
	TypeMsgApplyLoan          = "apply_loan"
	TypeMsgApproveLoan        = "approve_loan"
	TypeMsgRejectLoan         = "reject_loan"
	TypeMsgDisburseLoan       = "disburse_loan"
	TypeMsgRepayLoan          = "repay_loan"
	TypeMsgUpdateCropData     = "update_crop_data"
	TypeMsgClaimInsurance     = "claim_insurance"
	TypeMsgRegisterKCC        = "register_kcc"
)

// MsgApplyLoan - Apply for agricultural loan
type MsgApplyLoan struct {
	Applicant        string         `json:"applicant"`
	DhanPataID       string         `json:"dhanpata_id"`
	LoanAmount       sdk.Coin       `json:"loan_amount"`
	CropType         CropType       `json:"crop_type"`
	LandArea         string         `json:"land_area"` // in acres
	LandOwnershipDoc string         `json:"land_ownership_doc"`
	Pincode          string         `json:"pincode"`
	Purpose          string         `json:"purpose"`
	Duration         int64          `json:"duration"` // in months
	KCCNumber        string         `json:"kcc_number,omitempty"`
	CulturalQuote    string         `json:"cultural_quote"`
}

func NewMsgApplyLoan(
	applicant, dhanPataID string,
	loanAmount sdk.Coin,
	cropType CropType,
	landArea, landOwnershipDoc, pincode, purpose string,
	duration int64,
	kccNumber, culturalQuote string,
) *MsgApplyLoan {
	return &MsgApplyLoan{
		Applicant:        applicant,
		DhanPataID:       dhanPataID,
		LoanAmount:       loanAmount,
		CropType:         cropType,
		LandArea:         landArea,
		LandOwnershipDoc: landOwnershipDoc,
		Pincode:          pincode,
		Purpose:          purpose,
		Duration:         duration,
		KCCNumber:        kccNumber,
		CulturalQuote:    culturalQuote,
	}
}

func (msg *MsgApplyLoan) Route() string { return RouterKey }
func (msg *MsgApplyLoan) Type() string  { return TypeMsgApplyLoan }

func (msg *MsgApplyLoan) GetSigners() []sdk.AccAddress {
	applicant, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{applicant}
}

func (msg *MsgApplyLoan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant address (%s)", err)
	}

	if !msg.LoanAmount.IsValid() || msg.LoanAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "loan amount must be positive")
	}

	if msg.LandArea == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "land area cannot be empty")
	}

	if msg.Duration < 3 || msg.Duration > 60 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan duration must be between 3-60 months")
	}

	if len(msg.Pincode) != 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid PIN code")
	}

	return nil
}

// MsgApproveLoan - Approve agricultural loan
type MsgApproveLoan struct {
	Approver         string   `json:"approver"`
	LoanID           string   `json:"loan_id"`
	ApprovedAmount   sdk.Coin `json:"approved_amount"`
	InterestRate     sdk.Dec  `json:"interest_rate"`
	RepaymentSchedule string   `json:"repayment_schedule"`
	Remarks          string   `json:"remarks"`
}

func NewMsgApproveLoan(
	approver, loanID string,
	approvedAmount sdk.Coin,
	interestRate sdk.Dec,
	repaymentSchedule, remarks string,
) *MsgApproveLoan {
	return &MsgApproveLoan{
		Approver:          approver,
		LoanID:            loanID,
		ApprovedAmount:    approvedAmount,
		InterestRate:      interestRate,
		RepaymentSchedule: repaymentSchedule,
		Remarks:           remarks,
	}
}

func (msg *MsgApproveLoan) Route() string { return RouterKey }
func (msg *MsgApproveLoan) Type() string  { return TypeMsgApproveLoan }

func (msg *MsgApproveLoan) GetSigners() []sdk.AccAddress {
	approver, err := sdk.AccAddressFromBech32(msg.Approver)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{approver}
}

func (msg *MsgApproveLoan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Approver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid approver address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if !msg.ApprovedAmount.IsValid() || msg.ApprovedAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "approved amount must be positive")
	}

	if msg.InterestRate.IsNegative() || msg.InterestRate.GT(sdk.NewDecWithPrec(12, 2)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "interest rate must be between 0-12%")
	}

	return nil
}

// MsgRepayLoan - Repay agricultural loan
type MsgRepayLoan struct {
	Borrower      string   `json:"borrower"`
	LoanID        string   `json:"loan_id"`
	RepayAmount   sdk.Coin `json:"repay_amount"`
	CulturalQuote string   `json:"cultural_quote"`
}

func NewMsgRepayLoan(
	borrower, loanID string,
	repayAmount sdk.Coin,
	culturalQuote string,
) *MsgRepayLoan {
	return &MsgRepayLoan{
		Borrower:      borrower,
		LoanID:        loanID,
		RepayAmount:   repayAmount,
		CulturalQuote: culturalQuote,
	}
}

func (msg *MsgRepayLoan) Route() string { return RouterKey }
func (msg *MsgRepayLoan) Type() string  { return TypeMsgRepayLoan }

func (msg *MsgRepayLoan) GetSigners() []sdk.AccAddress {
	borrower, err := sdk.AccAddressFromBech32(msg.Borrower)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{borrower}
}

func (msg *MsgRepayLoan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRepayLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Borrower)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid borrower address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if !msg.RepayAmount.IsValid() || msg.RepayAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "repay amount must be positive")
	}

	return nil
}

// MsgUpdateCropData - Update crop yield and market data
type MsgUpdateCropData struct {
	Authority    string  `json:"authority"`
	LoanID       string  `json:"loan_id"`
	CropYield    string  `json:"crop_yield"` // in quintals
	MarketPrice  sdk.Dec `json:"market_price"`
	WeatherEvent string  `json:"weather_event,omitempty"`
	Remarks      string  `json:"remarks"`
}

func NewMsgUpdateCropData(
	authority, loanID, cropYield string,
	marketPrice sdk.Dec,
	weatherEvent, remarks string,
) *MsgUpdateCropData {
	return &MsgUpdateCropData{
		Authority:    authority,
		LoanID:       loanID,
		CropYield:    cropYield,
		MarketPrice:  marketPrice,
		WeatherEvent: weatherEvent,
		Remarks:      remarks,
	}
}

func (msg *MsgUpdateCropData) Route() string { return RouterKey }
func (msg *MsgUpdateCropData) Type() string  { return TypeMsgUpdateCropData }

func (msg *MsgUpdateCropData) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateCropData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCropData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if msg.MarketPrice.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "market price cannot be negative")
	}

	return nil
}