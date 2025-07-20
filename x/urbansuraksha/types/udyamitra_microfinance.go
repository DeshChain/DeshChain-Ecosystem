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

// UdyamitraSMELoan represents micro-finance loans for Small and Medium Enterprises
// Udyamitra (उद्यमित्र) = "Friend of Enterprise" - culturally rooted in supporting Indian entrepreneurs
type UdyamitraSMELoan struct {
	LoanID              string         `json:"loan_id"`
	BusinessName        string         `json:"business_name"`
	OwnerName           string         `json:"owner_name"`
	OwnerAddress        sdk.AccAddress `json:"owner_address"`
	BusinessType        string         `json:"business_type"`        // "manufacturing", "trading", "services", "digital"
	BusinessCategory    string         `json:"business_category"`    // "retail", "wholesale", "export", "handicrafts", etc.
	
	// Loan details
	LoanAmount          sdk.Coin       `json:"loan_amount"`
	InterestRate        sdk.Dec        `json:"interest_rate"`        // 8-15% based on risk and business type
	LoanTerm            uint32         `json:"loan_term"`            // 6-36 months
	RepaymentFrequency  string         `json:"repayment_frequency"`  // "weekly", "monthly", "quarterly"
	CollateralRequired  bool           `json:"collateral_required"`
	CollateralValue     sdk.Coin       `json:"collateral_value,omitempty"`
	
	// Business specifics
	MonthlyRevenue      sdk.Coin       `json:"monthly_revenue"`
	MonthlyExpenses     sdk.Coin       `json:"monthly_expenses"`
	YearsInBusiness     uint32         `json:"years_in_business"`
	EmployeeCount       uint32         `json:"employee_count"`
	BusinessLocation    string         `json:"business_location"`
	BusinessPAN         string         `json:"business_pan"`
	GSTNumber           string         `json:"gst_number,omitempty"`
	
	// Purpose and utilization
	LoanPurpose         string         `json:"loan_purpose"`         // "working_capital", "equipment", "expansion", "inventory"
	ExpectedROI         sdk.Dec        `json:"expected_roi"`         // Expected return on investment
	UtilizationPlan     []string       `json:"utilization_plan"`     // Detailed breakdown of fund usage
	
	// Risk assessment
	CreditScore         uint32         `json:"credit_score"`
	BusinessRiskLevel   string         `json:"business_risk_level"`  // "low", "medium", "high"
	MarketConditions    string         `json:"market_conditions"`    // Assessment of market for the business
	SeasonalityFactor   sdk.Dec        `json:"seasonality_factor"`   // Business seasonality impact
	
	// Community validation
	CommunityEndorsed   bool           `json:"community_endorsed"`
	LocalValidators     []string       `json:"local_validators,omitempty"`
	ValidationComments  string         `json:"validation_comments,omitempty"`
	
	// Repayment tracking
	TotalRepaid         sdk.Coin       `json:"total_repaid"`
	InterestPaid        sdk.Coin       `json:"interest_paid"`
	PrincipalPaid       sdk.Coin       `json:"principal_paid"`
	OutstandingAmount   sdk.Coin       `json:"outstanding_amount"`
	NextPaymentDue      time.Time      `json:"next_payment_due"`
	PaymentsMade        uint32         `json:"payments_made"`
	PaymentsRemaining   uint32         `json:"payments_remaining"`
	DaysOverdue         uint32         `json:"days_overdue"`
	
	// Performance metrics
	BusinessGrowth      sdk.Dec        `json:"business_growth"`      // Revenue growth since loan
	JobsCreated         uint32         `json:"jobs_created"`         // New employment generated
	DigitalAdoption     bool           `json:"digital_adoption"`     // Using digital payment systems
	TaxCompliance       bool           `json:"tax_compliance"`       // Regular tax filing
	
	// Cultural and social impact
	WomenEnterprise     bool           `json:"women_enterprise"`     // Women-owned business
	TribalEnterprise    bool           `json:"tribal_enterprise"`    // Tribal entrepreneur
	HandicraftBusiness  bool           `json:"handicraft_business"`  // Traditional handicrafts
	RuralEnterprise     bool           `json:"rural_enterprise"`     // Rural location
	SocialImpactScore   sdk.Dec        `json:"social_impact_score"`  // Overall social contribution
	
	// Loan status and timeline
	Status              string         `json:"status"`
	ApplicationDate     time.Time      `json:"application_date"`
	ApprovalDate        time.Time      `json:"approval_date,omitempty"`
	DisbursementDate    time.Time      `json:"disbursement_date,omitempty"`
	MaturityDate        time.Time      `json:"maturity_date,omitempty"`
	ClosureDate         time.Time      `json:"closure_date,omitempty"`
	
	// Udyamitra specific features
	MentorshipProvided  bool           `json:"mentorship_provided"`
	BusinessTraining    []string       `json:"business_training,omitempty"`
	NetworkingEvents    uint32         `json:"networking_events"`
	MarketLinkages      []string       `json:"market_linkages,omitempty"`
	TechnologySupport   bool           `json:"technology_support"`
}

// UdyamitraBusinessCategory represents different categories of businesses supported
type UdyamitraBusinessCategory struct {
	CategoryID          string         `json:"category_id"`
	CategoryName        string         `json:"category_name"`
	MinLoanAmount       sdk.Coin       `json:"min_loan_amount"`
	MaxLoanAmount       sdk.Coin       `json:"max_loan_amount"`
	BaseInterestRate    sdk.Dec        `json:"base_interest_rate"`
	RiskMultiplier      sdk.Dec        `json:"risk_multiplier"`
	
	// Category-specific features
	SeasonalBusiness    bool           `json:"seasonal_business"`
	CollateralMandatory bool           `json:"collateral_mandatory"`
	CommunityValidation bool           `json:"community_validation"`
	MentorshipRequired  bool           `json:"mentorship_required"`
	
	// Support programs
	TrainingPrograms    []string       `json:"training_programs"`
	MarketSupport       []string       `json:"market_support"`
	TechnologyTools     []string       `json:"technology_tools"`
	
	// Cultural significance
	TraditionalCraft    bool           `json:"traditional_craft"`
	CulturalImportance  string         `json:"cultural_importance"`
	GovernmentSchemes   []string       `json:"government_schemes,omitempty"`
}

// UdyamitraLiquidity represents the SME liquidity pool management
type UdyamitraLiquidity struct {
	PoolID              uint64         `json:"pool_id"`
	PoolName            string         `json:"pool_name"`
	RegionCode          string         `json:"region_code"`
	
	// Pool composition from Urban Pension
	SurakshaContribution sdk.Coin       `json:"suraksha_contribution"`  // 15% from Urban Pension Pool
	DirectInvestment    sdk.Coin       `json:"direct_investment"`     // Direct SME investors
	GovernmentSubsidy   sdk.Coin       `json:"government_subsidy"`    // MUDRA/MSME scheme integration
	CorporateCSR        sdk.Coin       `json:"corporate_csr"`         // CSR fund allocation
	
	// Pool utilization
	TotalLiquidity      sdk.Coin       `json:"total_liquidity"`
	ActiveLoans         sdk.Coin       `json:"active_loans"`
	AvailableLiquidity  sdk.Coin       `json:"available_liquidity"`
	ReserveAmount       sdk.Coin       `json:"reserve_amount"`        // 10% safety reserve
	
	// Performance metrics
	PortfolioPerformance sdk.Dec       `json:"portfolio_performance"`
	DefaultRate         sdk.Dec        `json:"default_rate"`
	AverageTicketSize   sdk.Coin       `json:"average_ticket_size"`
	AverageInterestRate sdk.Dec        `json:"average_interest_rate"`
	
	// Impact metrics
	BusinessesSupported uint32         `json:"businesses_supported"`
	JobsCreated         uint32         `json:"jobs_created"`
	WomenBorrowers      uint32         `json:"women_borrowers"`
	RuralBorrowers      uint32         `json:"rural_borrowers"`
	
	// Revenue generation
	MonthlyInterestIncome sdk.Coin     `json:"monthly_interest_income"`
	ProcessingFees       sdk.Coin      `json:"processing_fees"`
	LateFees            sdk.Coin       `json:"late_fees"`
	TotalRevenue        sdk.Coin       `json:"total_revenue"`
	
	// Cultural programs
	HandicraftLoans     uint32         `json:"handicraft_loans"`
	TraditionalArts     uint32         `json:"traditional_arts"`
	CulturalEvents      uint32         `json:"cultural_events"`
	HeritageBusinesses  uint32         `json:"heritage_businesses"`
}

// Revenue-Based Financing for high-growth SMEs
type UdyamitraRBF struct {
	RBFID               string         `json:"rbf_id"`
	BusinessID          string         `json:"business_id"`
	OwnerAddress        sdk.AccAddress `json:"owner_address"`
	
	// RBF terms
	InvestmentAmount    sdk.Coin       `json:"investment_amount"`
	RevenuePercentage   sdk.Dec        `json:"revenue_percentage"`    // 2-8% of monthly revenue
	CapMultiple         sdk.Dec        `json:"cap_multiple"`          // 1.2x to 3x investment amount
	PaybackCap          sdk.Coin       `json:"payback_cap"`           // Maximum total payback
	MinMonthlyPayment   sdk.Coin       `json:"min_monthly_payment"`
	MaxMonthlyPayment   sdk.Coin       `json:"max_monthly_payment"`
	
	// Performance tracking
	MonthlyRevenue      []sdk.Coin     `json:"monthly_revenue"`       // Historical revenue
	MonthlyPayments     []sdk.Coin     `json:"monthly_payments"`      // Payments made
	TotalPaid           sdk.Coin       `json:"total_paid"`
	RemainingCap        sdk.Coin       `json:"remaining_cap"`
	
	// Business growth metrics
	RevenueGrowthRate   sdk.Dec        `json:"revenue_growth_rate"`
	ProfitMargin        sdk.Dec        `json:"profit_margin"`
	MarketExpansion     bool           `json:"market_expansion"`
	ProductDiversification bool        `json:"product_diversification"`
	
	Status              string         `json:"status"`
	StartDate           time.Time      `json:"start_date"`
	ExpectedCompletion  time.Time      `json:"expected_completion"`
}

// Supply Chain Financing for traders and manufacturers
type UdyamitraSupplyChainFinance struct {
	SCFID               string         `json:"scf_id"`
	SupplierAddress     sdk.AccAddress `json:"supplier_address"`
	BuyerAddress        sdk.AccAddress `json:"buyer_address"`
	
	// Transaction details
	InvoiceAmount       sdk.Coin       `json:"invoice_amount"`
	AdvancePercentage   sdk.Dec        `json:"advance_percentage"`    // 70-90% of invoice
	AdvanceAmount       sdk.Coin       `json:"advance_amount"`
	FinancingFee        sdk.Dec        `json:"financing_fee"`         // 1-3% of invoice amount
	PaymentTerms        uint32         `json:"payment_terms"`         // Days to payment
	
	// Supply chain details
	ProductCategory     string         `json:"product_category"`
	DeliveryDate        time.Time      `json:"delivery_date"`
	PaymentDueDate      time.Time      `json:"payment_due_date"`
	InvoiceNumber       string         `json:"invoice_number"`
	DeliveryProof       string         `json:"delivery_proof"`
	
	// Risk mitigation
	BuyerCreditScore    uint32         `json:"buyer_credit_score"`
	SupplierHistory     uint32         `json:"supplier_history"`      // Transaction count
	InsuranceCoverage   sdk.Coin       `json:"insurance_coverage"`
	
	Status              string         `json:"status"`
	CreatedDate         time.Time      `json:"created_date"`
	SettlementDate      time.Time      `json:"settlement_date,omitempty"`
}

// Digital Gold Micro-Lending for jewelry businesses
type UdyamitraGoldLending struct {
	GoldLoanID          string         `json:"gold_loan_id"`
	BorrowerAddress     sdk.AccAddress `json:"borrower_address"`
	
	// Gold collateral
	GoldWeight          sdk.Dec        `json:"gold_weight"`           // In grams
	GoldPurity          uint32         `json:"gold_purity"`           // 22K, 24K etc
	CurrentGoldPrice    sdk.Dec        `json:"current_gold_price"`    // Per gram
	CollateralValue     sdk.Coin       `json:"collateral_value"`
	LoanToValue         sdk.Dec        `json:"loan_to_value"`         // 75-85% of gold value
	
	// Loan terms
	LoanAmount          sdk.Coin       `json:"loan_amount"`
	InterestRate        sdk.Dec        `json:"interest_rate"`         // 12-18% for gold loans
	LoanTerm            uint32         `json:"loan_term"`             // 6-12 months typically
	
	// Gold storage and tracking
	StorageLocation     string         `json:"storage_location"`
	StorageProvider     string         `json:"storage_provider"`
	DigitalCertificate  string         `json:"digital_certificate"`
	QRCode              string         `json:"qr_code"`               // For tracking
	
	// Market risk management
	PriceAlerts         []sdk.Dec      `json:"price_alerts"`          // Alert thresholds
	MarginCalls         []time.Time    `json:"margin_calls"`
	TopUpRequirements   sdk.Coin       `json:"top_up_requirements"`
	
	Status              string         `json:"status"`
	PledgeDate          time.Time      `json:"pledge_date"`
	MaturityDate        time.Time      `json:"maturity_date"`
	ReleaseDate         time.Time      `json:"release_date,omitempty"`
}

// Constants for Udyamitra SME Loans
const (
	// Loan categories
	UdyamitraManufacturing = "manufacturing"
	UdyamitraTrading      = "trading"
	UdyamitraServices     = "services"
	UdyamitraDigital      = "digital"
	UdyamitraHandicrafts  = "handicrafts"
	UdyamitraAgriculture  = "agriculture"
	UdyamitraTextiles     = "textiles"
	UdyamitraJewelry      = "jewelry"
	
	// Interest rate ranges
	MinSMEInterestRate = 8.0  // 8% for low-risk, collateralized loans
	MaxSMEInterestRate = 15.0 // 15% for high-risk, unsecured loans
	
	// Loan amount limits
	MinSMELoan = 25000    // ₹25,000 minimum
	MaxSMELoan = 5000000  // ₹50 lakhs maximum
	
	// Special category bonuses
	WomenEntrepreneurDiscount = 1.0 // 1% interest rate discount
	TribalEntrepreneurDiscount = 1.5 // 1.5% discount
	HandicraftBonus = 2.0 // 2% discount for traditional crafts
	RuralBusinessBonus = 0.5 // 0.5% discount for rural enterprises
	
	// Pool allocation from Urban Pension
	SMEPoolAllocation = 15 // 15% of Urban Pension pool for SME loans
)

// Loan status constants
const (
	SMEStatusPending     = "pending"
	SMEStatusUnderReview = "under_review"
	SMEStatusApproved    = "approved"
	SMEStatusDisbursed   = "disbursed"
	SMEStatusActive      = "active"
	SMEStatusCompleted   = "completed"
	SMEStatusDefaulted   = "defaulted"
	SMEStatusRestructured = "restructured"
)

// Validation functions
func (sme UdyamitraSMELoan) Validate() error {
	if sme.LoanID == "" {
		return ErrInvalidLoanID
	}
	if sme.OwnerAddress.Empty() {
		return ErrInvalidBorrower
	}
	if !sme.LoanAmount.IsPositive() {
		return ErrInvalidLoanAmount
	}
	if sme.InterestRate.IsNegative() {
		return ErrInvalidInterestRate
	}
	if sme.BusinessName == "" || sme.OwnerName == "" {
		return ErrInvalidBusinessDetails
	}
	return nil
}

// CalculateOptimalInterestRate calculates interest rate based on multiple factors
func (sme UdyamitraSMELoan) CalculateOptimalInterestRate(baseRate sdk.Dec) sdk.Dec {
	finalRate := baseRate
	
	// Business type adjustments
	switch sme.BusinessType {
	case UdyamitraManufacturing:
		finalRate = finalRate.Add(sdk.NewDecWithPrec(5, 1)) // +0.5%
	case UdyamitraTrading:
		finalRate = finalRate.Add(sdk.NewDec(1)) // +1%
	case UdyamitraServices:
		finalRate = finalRate.Add(sdk.NewDecWithPrec(15, 1)) // +1.5%
	case UdyamitraDigital:
		finalRate = finalRate.Sub(sdk.NewDecWithPrec(5, 1)) // -0.5% (lower risk)
	case UdyamitraHandicrafts:
		finalRate = finalRate.Sub(sdk.NewDec(2)) // -2% (cultural support)
	}
	
	// Experience bonus
	if sme.YearsInBusiness >= 5 {
		finalRate = finalRate.Sub(sdk.NewDecWithPrec(5, 1)) // -0.5% for 5+ years
	}
	if sme.YearsInBusiness >= 10 {
		finalRate = finalRate.Sub(sdk.NewDecWithPrec(5, 1)) // Additional -0.5% for 10+ years
	}
	
	// Social impact discounts
	if sme.WomenEnterprise {
		finalRate = finalRate.Sub(sdk.NewDec(1)) // -1% for women entrepreneurs
	}
	if sme.TribalEnterprise {
		finalRate = finalRate.Sub(sdk.NewDecWithPrec(15, 1)) // -1.5% for tribal entrepreneurs
	}
	if sme.RuralEnterprise {
		finalRate = finalRate.Sub(sdk.NewDecWithPrec(5, 1)) // -0.5% for rural enterprises
	}
	
	// Risk adjustments
	switch sme.BusinessRiskLevel {
	case "low":
		finalRate = finalRate.Sub(sdk.NewDec(1)) // -1% for low risk
	case "high":
		finalRate = finalRate.Add(sdk.NewDec(2)) // +2% for high risk
	}
	
	// Ensure within bounds
	minRate := sdk.NewDecWithPrec(int64(MinSMEInterestRate), 2)
	maxRate := sdk.NewDecWithPrec(int64(MaxSMEInterestRate), 2)
	
	if finalRate.LT(minRate) {
		finalRate = minRate
	}
	if finalRate.GT(maxRate) {
		finalRate = maxRate
	}
	
	return finalRate
}

// CalculateSocialImpactScore calculates the social impact score
func (sme UdyamitraSMELoan) CalculateSocialImpactScore() sdk.Dec {
	score := sdk.ZeroDec()
	
	// Base score for job creation
	if sme.EmployeeCount > 0 {
		score = score.Add(sdk.NewDec(int64(sme.EmployeeCount)).Mul(sdk.NewDecWithPrec(1, 1))) // 0.1 per employee
	}
	
	// Jobs created since loan
	if sme.JobsCreated > 0 {
		score = score.Add(sdk.NewDec(int64(sme.JobsCreated)).Mul(sdk.NewDecWithPrec(2, 1))) // 0.2 per new job
	}
	
	// Special category bonuses
	if sme.WomenEnterprise {
		score = score.Add(sdk.NewDec(5)) // +5 points for women entrepreneurship
	}
	if sme.TribalEnterprise {
		score = score.Add(sdk.NewDec(7)) // +7 points for tribal entrepreneurship
	}
	if sme.HandicraftBusiness {
		score = score.Add(sdk.NewDec(8)) // +8 points for handicraft preservation
	}
	if sme.RuralEnterprise {
		score = score.Add(sdk.NewDec(3)) // +3 points for rural development
	}
	
	// Digital adoption bonus
	if sme.DigitalAdoption {
		score = score.Add(sdk.NewDec(2)) // +2 points for digital payments
	}
	
	// Tax compliance bonus
	if sme.TaxCompliance {
		score = score.Add(sdk.NewDec(2)) // +2 points for tax compliance
	}
	
	// Business growth bonus
	if sme.BusinessGrowth.GT(sdk.NewDecWithPrec(10, 2)) { // >10% growth
		score = score.Add(sdk.NewDec(3)) // +3 points for good growth
	}
	
	return score
}

// IsEligibleForRBF checks if business is eligible for Revenue-Based Financing
func (sme UdyamitraSMELoan) IsEligibleForRBF() bool {
	// Must have consistent revenue
	if sme.MonthlyRevenue.Amount.LT(sdk.NewInt(100000)) { // < ₹1 lakh monthly revenue
		return false
	}
	
	// Must be in business for at least 1 year
	if sme.YearsInBusiness < 1 {
		return false
	}
	
	// Must have digital payment adoption
	if !sme.DigitalAdoption {
		return false
	}
	
	// Good performance on existing loan (if any)
	if sme.DaysOverdue > 30 {
		return false
	}
	
	return true
}