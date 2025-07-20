package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types
const (
	TypeMsgApplyBusinessLoan     = "apply_business_loan"
	TypeMsgApproveLoan          = "approve_loan"
	TypeMsgRejectLoan           = "reject_loan"
	TypeMsgDisburseLoan         = "disburse_loan"
	TypeMsgRepayLoan            = "repay_loan"
	TypeMsgUpdateBusinessMetrics = "update_business_metrics"
	TypeMsgApplyInvoiceFinancing = "apply_invoice_financing"
	TypeMsgUpdateCreditLine     = "update_credit_line"
)

// MsgApplyBusinessLoan - Apply for business loan
type MsgApplyBusinessLoan struct {
	Applicant           string         `json:"applicant"`
	DhanPataID          string         `json:"dhanpata_id"`
	BusinessName        string         `json:"business_name"`
	BusinessType        BusinessType   `json:"business_type"`
	GSTNumber           string         `json:"gst_number"`
	LoanAmount          sdk.Coin       `json:"loan_amount"`
	LoanPurpose         LoanPurpose    `json:"loan_purpose"`
	Duration            int64          `json:"duration"` // in months
	AnnualRevenue       sdk.Coin       `json:"annual_revenue"`
	EmployeeCount       int32          `json:"employee_count"`
	BusinessAge         int32          `json:"business_age"` // in months
	Pincode             string         `json:"pincode"`
	CollateralOffered   string         `json:"collateral_offered"`
	BusinessPlan        string         `json:"business_plan"` // IPFS hash
	FinancialStatements string         `json:"financial_statements"` // IPFS hash
	CulturalQuote       string         `json:"cultural_quote"`
}

func NewMsgApplyBusinessLoan(
	applicant, dhanPataID, businessName string,
	businessType BusinessType,
	gstNumber string,
	loanAmount sdk.Coin,
	loanPurpose LoanPurpose,
	duration int64,
	annualRevenue sdk.Coin,
	employeeCount, businessAge int32,
	pincode, collateralOffered, businessPlan, financialStatements, culturalQuote string,
) *MsgApplyBusinessLoan {
	return &MsgApplyBusinessLoan{
		Applicant:           applicant,
		DhanPataID:          dhanPataID,
		BusinessName:        businessName,
		BusinessType:        businessType,
		GSTNumber:           gstNumber,
		LoanAmount:          loanAmount,
		LoanPurpose:         loanPurpose,
		Duration:            duration,
		AnnualRevenue:       annualRevenue,
		EmployeeCount:       employeeCount,
		BusinessAge:         businessAge,
		Pincode:             pincode,
		CollateralOffered:   collateralOffered,
		BusinessPlan:        businessPlan,
		FinancialStatements: financialStatements,
		CulturalQuote:       culturalQuote,
	}
}

func (msg *MsgApplyBusinessLoan) Route() string { return RouterKey }
func (msg *MsgApplyBusinessLoan) Type() string  { return TypeMsgApplyBusinessLoan }

func (msg *MsgApplyBusinessLoan) GetSigners() []sdk.AccAddress {
	applicant, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{applicant}
}

func (msg *MsgApplyBusinessLoan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyBusinessLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant address (%s)", err)
	}

	if msg.BusinessName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business name cannot be empty")
	}

	if !msg.LoanAmount.IsValid() || msg.LoanAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "loan amount must be positive")
	}

	if msg.Duration < 6 || msg.Duration > 84 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan duration must be between 6-84 months")
	}

	if len(msg.Pincode) != 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid PIN code")
	}

	if msg.EmployeeCount < 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "employee count cannot be negative")
	}

	return nil
}

// MsgApplyInvoiceFinancing - Apply for invoice financing
type MsgApplyInvoiceFinancing struct {
	Applicant         string   `json:"applicant"`
	BusinessLoanID    string   `json:"business_loan_id"`
	InvoiceNumber     string   `json:"invoice_number"`
	InvoiceAmount     sdk.Coin `json:"invoice_amount"`
	InvoiceDate       string   `json:"invoice_date"`
	DueDate           string   `json:"due_date"`
	BuyerGSTNumber    string   `json:"buyer_gst_number"`
	BuyerName         string   `json:"buyer_name"`
	InvoiceDocument   string   `json:"invoice_document"` // IPFS hash
	FinancingPercent  sdk.Dec  `json:"financing_percent"` // 70-90% of invoice
}

func NewMsgApplyInvoiceFinancing(
	applicant, businessLoanID, invoiceNumber string,
	invoiceAmount sdk.Coin,
	invoiceDate, dueDate, buyerGSTNumber, buyerName, invoiceDocument string,
	financingPercent sdk.Dec,
) *MsgApplyInvoiceFinancing {
	return &MsgApplyInvoiceFinancing{
		Applicant:        applicant,
		BusinessLoanID:   businessLoanID,
		InvoiceNumber:    invoiceNumber,
		InvoiceAmount:    invoiceAmount,
		InvoiceDate:      invoiceDate,
		DueDate:          dueDate,
		BuyerGSTNumber:   buyerGSTNumber,
		BuyerName:        buyerName,
		InvoiceDocument:  invoiceDocument,
		FinancingPercent: financingPercent,
	}
}

func (msg *MsgApplyInvoiceFinancing) Route() string { return RouterKey }
func (msg *MsgApplyInvoiceFinancing) Type() string  { return TypeMsgApplyInvoiceFinancing }

func (msg *MsgApplyInvoiceFinancing) GetSigners() []sdk.AccAddress {
	applicant, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{applicant}
}

func (msg *MsgApplyInvoiceFinancing) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyInvoiceFinancing) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant address (%s)", err)
	}

	if msg.InvoiceNumber == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invoice number cannot be empty")
	}

	if !msg.InvoiceAmount.IsValid() || msg.InvoiceAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invoice amount must be positive")
	}

	if msg.FinancingPercent.LT(sdk.NewDecWithPrec(70, 2)) || msg.FinancingPercent.GT(sdk.NewDecWithPrec(90, 2)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "financing percent must be between 70-90%")
	}

	return nil
}

// MsgUpdateBusinessMetrics - Update business performance metrics
type MsgUpdateBusinessMetrics struct {
	Authority          string   `json:"authority"`
	BusinessLoanID     string   `json:"business_loan_id"`
	MonthlyRevenue     sdk.Coin `json:"monthly_revenue"`
	MonthlyExpenses    sdk.Coin `json:"monthly_expenses"`
	NewEmployees       int32    `json:"new_employees"`
	GSTFilingStatus    string   `json:"gst_filing_status"`
	CreditUtilization  sdk.Dec  `json:"credit_utilization"`
	ReportingPeriod    string   `json:"reporting_period"`
}

func NewMsgUpdateBusinessMetrics(
	authority, businessLoanID string,
	monthlyRevenue, monthlyExpenses sdk.Coin,
	newEmployees int32,
	gstFilingStatus string,
	creditUtilization sdk.Dec,
	reportingPeriod string,
) *MsgUpdateBusinessMetrics {
	return &MsgUpdateBusinessMetrics{
		Authority:         authority,
		BusinessLoanID:    businessLoanID,
		MonthlyRevenue:    monthlyRevenue,
		MonthlyExpenses:   monthlyExpenses,
		NewEmployees:      newEmployees,
		GSTFilingStatus:   gstFilingStatus,
		CreditUtilization: creditUtilization,
		ReportingPeriod:   reportingPeriod,
	}
}

func (msg *MsgUpdateBusinessMetrics) Route() string { return RouterKey }
func (msg *MsgUpdateBusinessMetrics) Type() string  { return TypeMsgUpdateBusinessMetrics }

func (msg *MsgUpdateBusinessMetrics) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateBusinessMetrics) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBusinessMetrics) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.BusinessLoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business loan ID cannot be empty")
	}

	if !msg.MonthlyRevenue.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid monthly revenue")
	}

	if !msg.MonthlyExpenses.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid monthly expenses")
	}

	return nil
}

// MsgUpdateCreditLine - Update credit line for a business
type MsgUpdateCreditLine struct {
	Authority       string   `json:"authority"`
	BusinessLoanID  string   `json:"business_loan_id"`
	NewCreditLimit  sdk.Coin `json:"new_credit_limit"`
	Reason          string   `json:"reason"`
	ValidityPeriod  int64    `json:"validity_period"` // in days
}

func NewMsgUpdateCreditLine(
	authority, businessLoanID string,
	newCreditLimit sdk.Coin,
	reason string,
	validityPeriod int64,
) *MsgUpdateCreditLine {
	return &MsgUpdateCreditLine{
		Authority:      authority,
		BusinessLoanID: businessLoanID,
		NewCreditLimit: newCreditLimit,
		Reason:         reason,
		ValidityPeriod: validityPeriod,
	}
}

func (msg *MsgUpdateCreditLine) Route() string { return RouterKey }
func (msg *MsgUpdateCreditLine) Type() string  { return TypeMsgUpdateCreditLine }

func (msg *MsgUpdateCreditLine) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateCreditLine) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCreditLine) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.BusinessLoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business loan ID cannot be empty")
	}

	if !msg.NewCreditLimit.IsValid() || msg.NewCreditLimit.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "new credit limit must be positive")
	}

	if msg.ValidityPeriod < 30 || msg.ValidityPeriod > 365 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validity period must be between 30-365 days")
	}

	return nil
}