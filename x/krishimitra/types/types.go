package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LoanStatus represents the status of a loan
type LoanStatus int32

const (
	LoanStatus_PENDING LoanStatus = iota
	LoanStatus_APPROVED
	LoanStatus_DISBURSED
	LoanStatus_REPAYING
	LoanStatus_COMPLETED
	LoanStatus_DEFAULTED
	LoanStatus_REJECTED
)

// CropType represents different crop categories
type CropType int32

const (
	CropType_KHARIF CropType = iota
	CropType_RABI
	CropType_ZAID
	CropType_PERENNIAL
	CropType_HORTICULTURE
)

// Loan represents an agricultural loan
type Loan struct {
	ID                string        `json:"id"`
	Borrower          string        `json:"borrower"`
	DhanPataAddress   string        `json:"dhanpata_address"`
	Amount            sdk.Coin      `json:"amount"`
	InterestRate      sdk.Dec       `json:"interest_rate"`
	Term              int64         `json:"term"` // in months
	Status            LoanStatus    `json:"status"`
	Purpose           string        `json:"purpose"`
	CropType          CropType      `json:"crop_type"`
	LandArea          sdk.Dec       `json:"land_area"` // in acres
	ExpectedYield     sdk.Dec       `json:"expected_yield"`
	CreatedAt         time.Time     `json:"created_at"`
	DisbursedAt       *time.Time    `json:"disbursed_at,omitempty"`
	MaturityDate      *time.Time    `json:"maturity_date,omitempty"`
	RepaidAmount      sdk.Coin      `json:"repaid_amount"`
	LastRepaymentDate *time.Time    `json:"last_repayment_date,omitempty"`
	PINCode           string        `json:"pin_code"`
	VillageCode       string        `json:"village_code"`
	FestivalBonus     bool          `json:"festival_bonus"`
	InsuranceRequired bool          `json:"insurance_required"`
	InsurancePremium  sdk.Coin      `json:"insurance_premium"`
}

// LoanApplication represents a loan application
type LoanApplication struct {
	ID                   string          `json:"id"`
	Applicant            string          `json:"applicant"`
	DhanPataAddress      string          `json:"dhanpata_address"`
	RequestedAmount      sdk.Coin        `json:"requested_amount"`
	Purpose              string          `json:"purpose"`
	CropType             CropType        `json:"crop_type"`
	LandOwnershipProof   string          `json:"land_ownership_proof"`
	LandArea             sdk.Dec         `json:"land_area"`
	PreviousYield        sdk.Dec         `json:"previous_yield"`
	ExpectedYield        sdk.Dec         `json:"expected_yield"`
	PINCode              string          `json:"pin_code"`
	VillageCode          string          `json:"village_code"`
	AadhaarHash          string          `json:"aadhaar_hash"`
	KisanCreditCardNo    string          `json:"kisan_credit_card_no,omitempty"`
	BankAccount          string          `json:"bank_account"`
	AppliedAt            time.Time       `json:"applied_at"`
	ReviewedBy           string          `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time      `json:"reviewed_at,omitempty"`
	Status               LoanStatus      `json:"status"`
	RejectionReason      string          `json:"rejection_reason,omitempty"`
	CreditScore          int32           `json:"credit_score"`
	RepaymentHistory     []LoanReference `json:"repayment_history"`
}

// Collateral represents loan collateral
type Collateral struct {
	LoanID          string   `json:"loan_id"`
	Type            string   `json:"type"` // land, crop, equipment
	Description     string   `json:"description"`
	ValuationAmount sdk.Coin `json:"valuation_amount"`
	DocumentHash    string   `json:"document_hash"`
	VerifiedBy      string   `json:"verified_by"`
	VerifiedAt      time.Time `json:"verified_at"`
}

// Repayment represents a loan repayment
type Repayment struct {
	ID            string    `json:"id"`
	LoanID        string    `json:"loan_id"`
	Amount        sdk.Coin  `json:"amount"`
	Principal     sdk.Coin  `json:"principal"`
	Interest      sdk.Coin  `json:"interest"`
	PaidBy        string    `json:"paid_by"`
	PaidAt        time.Time `json:"paid_at"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"` // UPI, bank_transfer, cash
}

// PINCodeEligibility represents eligibility criteria for a PIN code
type PINCodeEligibility struct {
	PINCode            string   `json:"pin_code"`
	DistrictName       string   `json:"district_name"`
	StateName          string   `json:"state_name"`
	IsEligible         bool     `json:"is_eligible"`
	MaxLoanAmount      sdk.Coin `json:"max_loan_amount"`
	BaseInterestRate   sdk.Dec  `json:"base_interest_rate"`
	PriorityDistrict   bool     `json:"priority_district"`
	DroughtProne       bool     `json:"drought_prone"`
	EligibleCrops      []string `json:"eligible_crops"`
	LocalFestivalDates []string `json:"local_festival_dates"`
}

// FestivalPeriod represents a festival period with bonus benefits
type FestivalPeriod struct {
	Name              string    `json:"name"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	InterestReduction sdk.Dec   `json:"interest_reduction"`
	Regions           []string  `json:"regions"` // PIN codes or states
}

// LoanReference for tracking loan history
type LoanReference struct {
	LoanID          string     `json:"loan_id"`
	Amount          sdk.Coin   `json:"amount"`
	Status          LoanStatus `json:"status"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	DelayedPayments int32      `json:"delayed_payments"`
}