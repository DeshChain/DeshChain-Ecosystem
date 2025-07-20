/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LoanScheme represents an agricultural loan scheme
type LoanScheme struct {
	SchemeID            string    `json:"scheme_id"`
	SchemeName          string    `json:"scheme_name"`
	Description         string    `json:"description"`
	LoanType            string    `json:"loan_type"`
	MinAmount           sdk.Coin  `json:"min_amount"`
	MaxAmount           sdk.Coin  `json:"max_amount"`
	MinInterestRate     sdk.Dec   `json:"min_interest_rate"`
	MaxInterestRate     sdk.Dec   `json:"max_interest_rate"`
	MinTerm             uint32    `json:"min_term"`           // in months
	MaxTerm             uint32    `json:"max_term"`           // in months
	CollateralRequired  bool      `json:"collateral_required"`
	CollateralRatio     sdk.Dec   `json:"collateral_ratio"`   // e.g., 1.25 = 125%
	InsuranceRequired   bool      `json:"insurance_required"`
	InsurancePremium    sdk.Dec   `json:"insurance_premium"`  // percentage
	ProcessingFee       sdk.Dec   `json:"processing_fee"`     // percentage
	GracePeriodDays     uint32    `json:"grace_period_days"`
	LatePaymentPenalty  sdk.Dec   `json:"late_payment_penalty"` // percentage
	MinCreditScore      uint32    `json:"min_credit_score"`
	MaxCreditScore      uint32    `json:"max_credit_score"`
	CropSeasons         []string  `json:"crop_seasons,omitempty"` // applicable seasons
	EligibleCrops       []string  `json:"eligible_crops,omitempty"` // specific crops
	VillageRestrictions []string  `json:"village_restrictions,omitempty"` // postal codes
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ValidFrom           time.Time `json:"valid_from"`
	ValidTo             time.Time `json:"valid_to"`
	
	// Community features
	CommunityValidation bool      `json:"community_validation"`
	MinValidators       uint32    `json:"min_validators"`
	
	// Liquidity pool integration
	LiquiditySourceType string    `json:"liquidity_source_type"` // "unified_pool", "dedicated", "mixed"
	LiquidityAllocation sdk.Dec   `json:"liquidity_allocation"`  // percentage of pool
	
	// Risk management
	RiskBasedPricing    bool      `json:"risk_based_pricing"`
	MaxRiskLevel        string    `json:"max_risk_level"`
	WeatherInsurance    bool      `json:"weather_insurance"`
	CropInsurance       bool      `json:"crop_insurance"`
}

// Borrower represents a farmer/borrower
type Borrower struct {
	BorrowerID        string         `json:"borrower_id"`
	Address           sdk.AccAddress `json:"address"`
	Name              string         `json:"name"`
	Age               uint32         `json:"age"`
	Phone             string         `json:"phone"`
	VillageCode       string         `json:"village_code"`
	VillageName       string         `json:"village_name"`
	LandSize          sdk.Dec        `json:"land_size"`        // in acres
	LandOwnership     string         `json:"land_ownership"`   // "owned", "leased", "shared"
	PrimaryMobilization string       `json:"primary_crops"`
	SecondaryMobilization string     `json:"secondary_crops,omitempty"`
	BankAccount       string         `json:"bank_account"`
	IFSCCode          string         `json:"ifsc_code"`
	AadhaarNumber     string         `json:"aadhaar_number"`
	PANNumber         string         `json:"pan_number,omitempty"`
	CreditScore       uint32         `json:"credit_score"`
	Status            string         `json:"status"`
	KYCStatus         string         `json:"kyc_status"`
	RegisteredAt      time.Time      `json:"registered_at"`
	LastUpdated       time.Time      `json:"last_updated"`
	
	// Performance metrics
	TotalLoansCount      uint32    `json:"total_loans_count"`
	ActiveLoansCount     uint32    `json:"active_loans_count"`
	TotalBorrowed        sdk.Coin  `json:"total_borrowed"`
	TotalRepaid          sdk.Coin  `json:"total_repaid"`
	DefaultedLoansCount  uint32    `json:"defaulted_loans_count"`
	OnTimePaymentsCount  uint32    `json:"on_time_payments_count"`
	LatePaymentsCount    uint32    `json:"late_payments_count"`
	
	// Community standing
	CommunityScore       uint32    `json:"community_score"`
	Recommendations      uint32    `json:"recommendations"`
	CommunityWarnings    uint32    `json:"community_warnings"`
	
	// Risk assessment
	RiskScore            sdk.Dec   `json:"risk_score"`
	RiskLevel            string    `json:"risk_level"`
	LastRiskAssessment   time.Time `json:"last_risk_assessment"`
}

// LoanApplication represents a loan application
type LoanApplication struct {
	ApplicationID    string         `json:"application_id"`
	LoanSchemeID     string         `json:"loan_scheme_id"`
	BorrowerID       string         `json:"borrower_id"`
	BorrowerAddress  sdk.AccAddress `json:"borrower_address"`
	RequestedAmount  sdk.Coin       `json:"requested_amount"`
	RequestedTerm    uint32         `json:"requested_term"`     // in months
	Purpose          string         `json:"purpose"`
	CropType         string         `json:"crop_type,omitempty"`
	CropSeason       string         `json:"crop_season,omitempty"`
	ExpectedHarvest  time.Time      `json:"expected_harvest,omitempty"`
	MarketPrice      sdk.Dec        `json:"market_price,omitempty"` // per unit
	ExpectedYield    sdk.Dec        `json:"expected_yield,omitempty"` // in units
	
	// Collateral information
	CollateralType   string         `json:"collateral_type,omitempty"`
	CollateralValue  sdk.Coin       `json:"collateral_value,omitempty"`
	CollateralDocs   []string       `json:"collateral_docs,omitempty"`
	
	// Application status
	Status           string         `json:"status"`
	SubmittedAt      time.Time      `json:"submitted_at"`
	ReviewedAt       time.Time      `json:"reviewed_at,omitempty"`
	ApprovedAt       time.Time      `json:"approved_at,omitempty"`
	RejectedAt       time.Time      `json:"rejected_at,omitempty"`
	RejectionReason  string         `json:"rejection_reason,omitempty"`
	
	// Approval details
	ApprovedAmount   sdk.Coin       `json:"approved_amount,omitempty"`
	ApprovedTerm     uint32         `json:"approved_term,omitempty"`
	InterestRate     sdk.Dec        `json:"interest_rate,omitempty"`
	ProcessingFee    sdk.Coin       `json:"processing_fee,omitempty"`
	
	// Community validation
	CommunityValidators []string    `json:"community_validators,omitempty"`
	ValidationStatus    string      `json:"validation_status"`
	ValidationComments  string      `json:"validation_comments,omitempty"`
	
	// Risk assessment
	RiskScore        sdk.Dec        `json:"risk_score,omitempty"`
	RiskFactors      []string       `json:"risk_factors,omitempty"`
	MitigationSteps  []string       `json:"mitigation_steps,omitempty"`
}

// Loan represents an active loan
type Loan struct {
	LoanID           string         `json:"loan_id"`
	ApplicationID    string         `json:"application_id"`
	LoanSchemeID     string         `json:"loan_scheme_id"`
	BorrowerID       string         `json:"borrower_id"`
	BorrowerAddress  sdk.AccAddress `json:"borrower_address"`
	Principal        sdk.Coin       `json:"principal"`
	InterestRate     sdk.Dec        `json:"interest_rate"`
	Term             uint32         `json:"term"`              // in months
	MonthlyPayment   sdk.Coin       `json:"monthly_payment"`
	TotalInterest    sdk.Coin       `json:"total_interest"`
	TotalPayable     sdk.Coin       `json:"total_payable"`
	
	// Loan timeline
	DisbursedAt      time.Time      `json:"disbursed_at"`
	FirstPaymentDue  time.Time      `json:"first_payment_due"`
	LastPaymentDue   time.Time      `json:"last_payment_due"`
	MaturityDate     time.Time      `json:"maturity_date"`
	
	// Repayment tracking
	AmountPaid       sdk.Coin       `json:"amount_paid"`
	InterestPaid     sdk.Coin       `json:"interest_paid"`
	PrincipalPaid    sdk.Coin       `json:"principal_paid"`
	OutstandingAmount sdk.Coin      `json:"outstanding_amount"`
	PaymentsMade     uint32         `json:"payments_made"`
	PaymentsRemaining uint32        `json:"payments_remaining"`
	
	// Status and performance
	Status           string         `json:"status"`
	NextPaymentDue   time.Time      `json:"next_payment_due"`
	DaysOverdue      uint32         `json:"days_overdue"`
	LatePaymentPenalty sdk.Coin     `json:"late_payment_penalty"`
	
	// Collateral
	CollateralID     string         `json:"collateral_id,omitempty"`
	CollateralValue  sdk.Coin       `json:"collateral_value,omitempty"`
	CollateralStatus string         `json:"collateral_status,omitempty"`
	
	// Insurance
	InsuranceID      string         `json:"insurance_id,omitempty"`
	InsurancePremium sdk.Coin       `json:"insurance_premium,omitempty"`
	InsuranceStatus  string         `json:"insurance_status,omitempty"`
	
	// Community and performance
	CommunityEndorsed bool          `json:"community_endorsed"`
	PerformanceRating sdk.Dec       `json:"performance_rating"`
	
	// Liquidity source
	LiquidityPoolID  string         `json:"liquidity_pool_id,omitempty"`
	LiquiditySource  string         `json:"liquidity_source"` // "unified_pool", "dedicated", etc.
}

// LoanRepayment represents a loan repayment
type LoanRepayment struct {
	RepaymentID      string         `json:"repayment_id"`
	LoanID           string         `json:"loan_id"`
	BorrowerAddress  sdk.AccAddress `json:"borrower_address"`
	Amount           sdk.Coin       `json:"amount"`
	PrincipalPortion sdk.Coin       `json:"principal_portion"`
	InterestPortion  sdk.Coin       `json:"interest_portion"`
	PenaltyPortion   sdk.Coin       `json:"penalty_portion,omitempty"`
	PaymentDate      time.Time      `json:"payment_date"`
	DueDate          time.Time      `json:"due_date"`
	DaysLate         uint32         `json:"days_late"`
	PaymentMethod    string         `json:"payment_method"`
	TransactionHash  string         `json:"transaction_hash"`
	Status           string         `json:"status"`
	
	// Performance metrics
	OnTime           bool           `json:"on_time"`
	EarlyPayment     bool           `json:"early_payment"`
	PartialPayment   bool           `json:"partial_payment"`
}

// VillageLendingPool represents a village-level lending pool
type VillageLendingPool struct {
	PoolID           string         `json:"pool_id"`
	VillageCode      string         `json:"village_code"`
	VillageName      string         `json:"village_name"`
	CoordinatorAddress sdk.AccAddress `json:"coordinator_address"`
	TotalLiquidity   sdk.Coin       `json:"total_liquidity"`
	AvailableLiquidity sdk.Coin     `json:"available_liquidity"`
	LoansOutstanding sdk.Coin       `json:"loans_outstanding"`
	TotalLoansCount  uint32         `json:"total_loans_count"`
	ActiveLoansCount uint32         `json:"active_loans_count"`
	DefaultRate      sdk.Dec        `json:"default_rate"`
	AverageInterestRate sdk.Dec     `json:"average_interest_rate"`
	PerformanceScore sdk.Dec        `json:"performance_score"`
	Status           string         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	LastUpdated      time.Time      `json:"last_updated"`
	
	// Unified pool integration
	UnifiedPoolAllocation sdk.Coin  `json:"unified_pool_allocation"`
	UnifiedPoolUtilization sdk.Dec  `json:"unified_pool_utilization"`
	MonthlyRevenue       sdk.Coin   `json:"monthly_revenue"`
	
	// Community metrics
	BorrowerCount        uint32     `json:"borrower_count"`
	CommunityScore       sdk.Dec    `json:"community_score"`
	SocialImpactScore    sdk.Dec    `json:"social_impact_score"`
}

// Validate performs basic validation on LoanScheme
func (ls LoanScheme) Validate() error {
	if ls.SchemeID == "" {
		return ErrInvalidLoanScheme
	}
	if ls.SchemeName == "" {
		return ErrInvalidLoanScheme
	}
	if !ls.MinAmount.IsPositive() || !ls.MaxAmount.IsPositive() {
		return ErrInvalidLoanAmount
	}
	if ls.MinAmount.Amount.GT(ls.MaxAmount.Amount) {
		return ErrInvalidLoanAmount
	}
	if ls.MinInterestRate.IsNegative() || ls.MaxInterestRate.IsNegative() {
		return ErrInvalidInterestRate
	}
	if ls.MinInterestRate.GT(ls.MaxInterestRate) {
		return ErrInvalidInterestRate
	}
	if ls.MinTerm == 0 || ls.MaxTerm == 0 || ls.MinTerm > ls.MaxTerm {
		return ErrInvalidLoanTerm
	}
	return nil
}

// Validate performs basic validation on Borrower
func (b Borrower) Validate() error {
	if b.BorrowerID == "" {
		return ErrInvalidBorrower
	}
	if b.Address.Empty() {
		return ErrInvalidBorrower
	}
	if b.Name == "" {
		return ErrInvalidBorrower
	}
	if b.Age == 0 || b.Age > 100 {
		return ErrInvalidBorrower
	}
	if len(b.Phone) < MinPhoneLength || len(b.Phone) > MaxPhoneLength {
		return ErrInvalidBorrower
	}
	if len(b.AadhaarNumber) != MaxAadhaarLength {
		return ErrInvalidBorrower
	}
	if b.LandSize.IsNegative() {
		return ErrInvalidBorrower
	}
	return nil
}

// IsEligible checks if a borrower is eligible for a loan scheme
func (ls LoanScheme) IsEligible(borrower Borrower, amount sdk.Coin) bool {
	// Check amount limits
	if amount.Amount.LT(ls.MinAmount.Amount) || amount.Amount.GT(ls.MaxAmount.Amount) {
		return false
	}
	
	// Check credit score
	if borrower.CreditScore < ls.MinCreditScore || borrower.CreditScore > ls.MaxCreditScore {
		return false
	}
	
	// Check scheme status and validity
	if ls.Status != StatusActive {
		return false
	}
	
	// Check village restrictions
	if len(ls.VillageRestrictions) > 0 {
		found := false
		for _, code := range ls.VillageRestrictions {
			if code == borrower.VillageCode {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

// CalculateInterestRate calculates interest rate based on risk and scheme
func (ls LoanScheme) CalculateInterestRate(riskScore sdk.Dec) sdk.Dec {
	if !ls.RiskBasedPricing {
		return ls.MinInterestRate
	}
	
	// Linear interpolation based on risk score (0.0 to 1.0)
	rateRange := ls.MaxInterestRate.Sub(ls.MinInterestRate)
	adjustedRate := ls.MinInterestRate.Add(rateRange.Mul(riskScore))
	
	// Ensure within bounds
	if adjustedRate.LT(ls.MinInterestRate) {
		return ls.MinInterestRate
	}
	if adjustedRate.GT(ls.MaxInterestRate) {
		return ls.MaxInterestRate
	}
	
	return adjustedRate
}

// CalculateMonthlyPayment calculates monthly payment amount
func (l Loan) CalculateMonthlyPayment() sdk.Coin {
	if l.Term == 0 {
		return sdk.NewCoin(l.Principal.Denom, sdk.ZeroInt())
	}
	
	// Simple calculation for now; can be enhanced with compound interest
	monthlyInterest := l.Principal.Amount.ToDec().Mul(l.InterestRate).Quo(sdk.NewDec(12))
	monthlyPrincipal := l.Principal.Amount.ToDec().Quo(sdk.NewDec(int64(l.Term)))
	monthlyPayment := monthlyPrincipal.Add(monthlyInterest).TruncateInt()
	
	return sdk.NewCoin(l.Principal.Denom, monthlyPayment)
}

// IsOverdue checks if loan payment is overdue
func (l Loan) IsOverdue(currentTime time.Time) bool {
	return currentTime.After(l.NextPaymentDue)
}

// GetOutstandingAmount calculates current outstanding amount
func (l Loan) GetOutstandingAmount() sdk.Coin {
	return l.TotalPayable.Sub(l.AmountPaid)
}