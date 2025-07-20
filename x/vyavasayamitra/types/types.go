package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LoanStatus represents the status of a business loan
type LoanStatus int32

const (
	LoanStatus_PENDING LoanStatus = iota
	LoanStatus_APPROVED
	LoanStatus_DISBURSED
	LoanStatus_REPAYING
	LoanStatus_COMPLETED
	LoanStatus_DEFAULTED
	LoanStatus_REJECTED
	LoanStatus_RESTRUCTURED
)

// BusinessType represents different business categories
type BusinessType int32

const (
	BusinessType_SOLE_PROPRIETORSHIP BusinessType = iota
	BusinessType_PARTNERSHIP
	BusinessType_PRIVATE_LIMITED
	BusinessType_LLP
	BusinessType_MSME
	BusinessType_STARTUP
	BusinessType_EXPORT_UNIT
)

// BusinessLoan represents a business loan
type BusinessLoan struct {
	ID                string        `json:"id"`
	BusinessID        string        `json:"business_id"`
	Borrower          string        `json:"borrower"`
	DhanPataAddress   string        `json:"dhanpata_address"`
	Amount            sdk.Coin      `json:"amount"`
	InterestRate      sdk.Dec       `json:"interest_rate"`
	Term              int64         `json:"term"` // in months
	Status            LoanStatus    `json:"status"`
	Purpose           string        `json:"purpose"`
	BusinessCategory  string        `json:"business_category"`
	CreatedAt         time.Time     `json:"created_at"`
	DisbursedAt       *time.Time    `json:"disbursed_at,omitempty"`
	MaturityDate      *time.Time    `json:"maturity_date,omitempty"`
	RepaidAmount      sdk.Coin      `json:"repaid_amount"`
	LastRepaymentDate *time.Time    `json:"last_repayment_date,omitempty"`
	PINCode           string        `json:"pin_code"`
	FestivalBonus     bool          `json:"festival_bonus"`
	CreditLineID      string        `json:"credit_line_id,omitempty"`
	InvoiceFinancing  bool          `json:"invoice_financing"`
	InsuranceRequired bool          `json:"insurance_required"`
	InsurancePremium  sdk.Coin      `json:"insurance_premium"`
}

// BusinessProfile represents a business entity
type BusinessProfile struct {
	ID                  string         `json:"id"`
	BusinessName        string         `json:"business_name"`
	Owner               string         `json:"owner"`
	DhanPataAddress     string         `json:"dhanpata_address"`
	BusinessType        BusinessType   `json:"business_type"`
	Category            string         `json:"category"`
	RegistrationNumber  string         `json:"registration_number"`
	GSTNumber           string         `json:"gst_number"`
	PanNumber           string         `json:"pan_number"`
	EstablishedDate     time.Time      `json:"established_date"`
	AnnualRevenue       sdk.Coin       `json:"annual_revenue"`
	EmployeeCount       int32          `json:"employee_count"`
	PINCode             string         `json:"pin_code"`
	BankAccount         string         `json:"bank_account"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	CreditScore         int32          `json:"credit_score"`
	VerificationStatus  string         `json:"verification_status"`
	MerchantRating      sdk.Dec        `json:"merchant_rating"`
	TotalLoansAvailed   int32          `json:"total_loans_availed"`
	ActiveLoans         int32          `json:"active_loans"`
}

// LoanApplication represents a business loan application
type LoanApplication struct {
	ID                   string              `json:"id"`
	BusinessID           string              `json:"business_id"`
	Applicant            string              `json:"applicant"`
	DhanPataAddress      string              `json:"dhanpata_address"`
	RequestedAmount      sdk.Coin            `json:"requested_amount"`
	Purpose              string              `json:"purpose"`
	BusinessPlan         string              `json:"business_plan"`
	ProjectedRevenue     sdk.Coin            `json:"projected_revenue"`
	CollateralOffered    string              `json:"collateral_offered"`
	FinancialStatements  []FinancialDocument `json:"financial_statements"`
	BankStatements       []string            `json:"bank_statements"`
	AppliedAt            time.Time           `json:"applied_at"`
	ReviewedBy           string              `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time          `json:"reviewed_at,omitempty"`
	Status               LoanStatus          `json:"status"`
	RejectionReason      string              `json:"rejection_reason,omitempty"`
	RequestedTerm        int64               `json:"requested_term"`
	ProposedInterestRate sdk.Dec             `json:"proposed_interest_rate"`
}

// CreditLine represents a pre-approved credit facility
type CreditLine struct {
	ID              string    `json:"id"`
	BusinessID      string    `json:"business_id"`
	ApprovedAmount  sdk.Coin  `json:"approved_amount"`
	UtilizedAmount  sdk.Coin  `json:"utilized_amount"`
	AvailableAmount sdk.Coin  `json:"available_amount"`
	InterestRate    sdk.Dec   `json:"interest_rate"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	Status          string    `json:"status"`
	ReviewPeriod    int32     `json:"review_period"` // in days
	LastReviewDate  time.Time `json:"last_review_date"`
}

// Invoice for invoice financing
type Invoice struct {
	ID             string    `json:"id"`
	BusinessID     string    `json:"business_id"`
	InvoiceNumber  string    `json:"invoice_number"`
	CustomerName   string    `json:"customer_name"`
	InvoiceAmount  sdk.Coin  `json:"invoice_amount"`
	InvoiceDate    time.Time `json:"invoice_date"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `json:"status"`
	FinancedAmount sdk.Coin  `json:"financed_amount"`
	FinancedAt     *time.Time `json:"financed_at,omitempty"`
	CollectedAt    *time.Time `json:"collected_at,omitempty"`
}

// Collateral for business loans
type Collateral struct {
	LoanID           string    `json:"loan_id"`
	Type             string    `json:"type"` // property, equipment, inventory, receivables
	Description      string    `json:"description"`
	ValuationAmount  sdk.Coin  `json:"valuation_amount"`
	ValuationDate    time.Time `json:"valuation_date"`
	ValuationAgency  string    `json:"valuation_agency"`
	DocumentHash     string    `json:"document_hash"`
	InsuranceStatus  string    `json:"insurance_status"`
	LienMarked       bool      `json:"lien_marked"`
}

// Repayment represents a loan repayment
type Repayment struct {
	ID              string    `json:"id"`
	LoanID          string    `json:"loan_id"`
	Amount          sdk.Coin  `json:"amount"`
	Principal       sdk.Coin  `json:"principal"`
	Interest        sdk.Coin  `json:"interest"`
	PenaltyAmount   sdk.Coin  `json:"penalty_amount,omitempty"`
	PaidBy          string    `json:"paid_by"`
	PaidAt          time.Time `json:"paid_at"`
	TransactionID   string    `json:"transaction_id"`
	PaymentMethod   string    `json:"payment_method"`
	ReceiptNumber   string    `json:"receipt_number"`
}

// PINCodeEligibility for business loans
type PINCodeEligibility struct {
	PINCode            string   `json:"pin_code"`
	CityName           string   `json:"city_name"`
	StateName          string   `json:"state_name"`
	IsEligible         bool     `json:"is_eligible"`
	MaxLoanAmount      sdk.Coin `json:"max_loan_amount"`
	BaseInterestRate   sdk.Dec  `json:"base_interest_rate"`
	IndustrialArea     bool     `json:"industrial_area"`
	ExportHub          bool     `json:"export_hub"`
	StartupHub         bool     `json:"startup_hub"`
	EligibleCategories []string `json:"eligible_categories"`
	LocalFestivals     []string `json:"local_festivals"`
}

// MerchantRating for creditworthiness
type MerchantRating struct {
	BusinessID          string    `json:"business_id"`
	OverallRating       sdk.Dec   `json:"overall_rating"`
	PaymentScore        sdk.Dec   `json:"payment_score"`
	BusinessScore       sdk.Dec   `json:"business_score"`
	FinancialScore      sdk.Dec   `json:"financial_score"`
	TotalTransactions   int32     `json:"total_transactions"`
	OnTimePayments      int32     `json:"on_time_payments"`
	DelayedPayments     int32     `json:"delayed_payments"`
	DefaultedPayments   int32     `json:"defaulted_payments"`
	LastUpdated         time.Time `json:"last_updated"`
	RecommendedCreditLimit sdk.Coin `json:"recommended_credit_limit"`
}

// FinancialDocument for loan applications
type FinancialDocument struct {
	Type         string    `json:"type"` // balance_sheet, p&l, cash_flow
	Year         int32     `json:"year"`
	DocumentHash string    `json:"document_hash"`
	AuditedBy    string    `json:"audited_by,omitempty"`
	UploadedAt   time.Time `json:"uploaded_at"`
}