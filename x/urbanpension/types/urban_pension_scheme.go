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

// UrbanPensionScheme represents the Urban Pension Scheme with integrated services
type UrbanPensionScheme struct {
	SchemeID            string    `json:"scheme_id"`
	AccountID           string    `json:"account_id"`
	ContributorAddress  sdk.AccAddress `json:"contributor_address"`
	MonthlyContribution sdk.Coin `json:"monthly_contribution"` // ₹2500 worth of NAMO
	ContributionPeriod  uint32    `json:"contribution_period"`  // 18 months
	TotalContributions  sdk.Coin  `json:"total_contributions"`  // ₹45,000 worth
	StartDate           time.Time `json:"start_date"`
	MaturityDate        time.Time `json:"maturity_date"`        // 20th month
	Status              string    `json:"status"`
	
	// Contribution tracking
	PaidContributions   uint32    `json:"paid_contributions"`
	MissedContributions uint32    `json:"missed_contributions"`
	LastContribution    time.Time `json:"last_contribution"`
	NextContribution    time.Time `json:"next_contribution"`
	
	// Maturity and returns
	ExpectedReturn      sdk.Coin  `json:"expected_return"`      // Calculated return
	ReturnPercentage    sdk.Dec   `json:"return_percentage"`    // Dynamic based on pool performance
	MaturityPaid        bool      `json:"maturity_paid"`
	MaturityAmount      sdk.Coin  `json:"maturity_amount"`
	
	// Integrated services utilized
	EducationLoanTaken  bool      `json:"education_loan_taken"`
	EducationLoanAmount sdk.Coin  `json:"education_loan_amount,omitempty"`
	LifeInsuranceCover  sdk.Coin  `json:"life_insurance_cover"`
	HealthInsuranceCover sdk.Coin `json:"health_insurance_cover"`
	InsurancePremiumPaid sdk.Coin `json:"insurance_premium_paid"`
	
	// Referral system
	ReferrerAddress     sdk.AccAddress `json:"referrer_address,omitempty"`
	ReferralBonus       sdk.Coin       `json:"referral_bonus"`
	ReferralsGenerated  uint32         `json:"referrals_generated"`
	ReferralRewards     sdk.Coin       `json:"referral_rewards"`
	
	// Pool integration
	UrbanPoolID         uint64    `json:"urban_pool_id"`
	PoolContribution    sdk.Coin  `json:"pool_contribution"`    // 70% goes to unified pool
	ReserveContribution sdk.Coin  `json:"reserve_contribution"` // 30% stays in reserve
	
	// Performance metrics
	PerformanceScore    sdk.Dec   `json:"performance_score"`
	ConsistencyBonus    sdk.Coin  `json:"consistency_bonus"`
	CommunityRating     sdk.Dec   `json:"community_rating"`
}

// UrbanEducationLoan represents education loans for urban pension contributors
type UrbanEducationLoan struct {
	LoanID              string         `json:"loan_id"`
	PensionAccountID    string         `json:"pension_account_id"`
	BorrowerAddress     sdk.AccAddress `json:"borrower_address"`
	GuarantorAddress    sdk.AccAddress `json:"guarantor_address,omitempty"`
	
	// Loan details
	LoanAmount          sdk.Coin       `json:"loan_amount"`
	InterestRate        sdk.Dec        `json:"interest_rate"`        // 4-7% based on pension history
	LoanTerm            uint32         `json:"loan_term"`            // in months
	MonthlyEMI          sdk.Coin       `json:"monthly_emi"`
	
	// Education specifics
	InstitutionName     string         `json:"institution_name"`
	CourseType          string         `json:"course_type"`
	CourseDuration      uint32         `json:"course_duration"`      // in months
	ExpectedCompletion  time.Time      `json:"expected_completion"`
	FeeStructure        []sdk.Coin     `json:"fee_structure"`
	
	// Disbursement
	DisbursedAmount     sdk.Coin       `json:"disbursed_amount"`
	DisbursementSchedule []Disbursement `json:"disbursement_schedule"`
	
	// Repayment
	RepaymentStartDate  time.Time      `json:"repayment_start_date"` // After course completion + 6 months
	TotalRepaid         sdk.Coin       `json:"total_repaid"`
	OutstandingAmount   sdk.Coin       `json:"outstanding_amount"`
	
	// Performance incentives
	AcademicPerformance string         `json:"academic_performance"`  // A, B, C grades
	InterestWaiver      sdk.Coin       `json:"interest_waiver"`       // For excellent performance
	
	Status              string         `json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	LastUpdated         time.Time      `json:"last_updated"`
}

// Disbursement represents a loan disbursement schedule
type Disbursement struct {
	Amount      sdk.Coin  `json:"amount"`
	DueDate     time.Time `json:"due_date"`
	PaidDate    time.Time `json:"paid_date,omitempty"`
	Status      string    `json:"status"`
	Purpose     string    `json:"purpose"`     // "tuition", "hostel", "books", etc.
}

// UrbanInsurancePolicy represents life and health insurance for urban pension members
type UrbanInsurancePolicy struct {
	PolicyID            string         `json:"policy_id"`
	PensionAccountID    string         `json:"pension_account_id"`
	PolicyHolderAddress sdk.AccAddress `json:"policy_holder_address"`
	
	// Coverage details
	LifeCoverAmount     sdk.Coin       `json:"life_cover_amount"`     // ₹10-50 lakhs
	HealthCoverAmount   sdk.Coin       `json:"health_cover_amount"`   // ₹5-25 lakhs
	AccidentCover       sdk.Coin       `json:"accident_cover"`
	CriticalIllnessCover sdk.Coin      `json:"critical_illness_cover"`
	
	// Premium details
	MonthlyPremium      sdk.Coin       `json:"monthly_premium"`       // Auto-deducted from pension
	AnnualPremium       sdk.Coin       `json:"annual_premium"`
	PremiumPaid         sdk.Coin       `json:"premium_paid"`
	PremiumWaiver       bool           `json:"premium_waiver"`        // For consistent contributors
	
	// Policy period
	PolicyStartDate     time.Time      `json:"policy_start_date"`
	PolicyEndDate       time.Time      `json:"policy_end_date"`
	RenewalDate         time.Time      `json:"renewal_date"`
	
	// Beneficiaries
	LifeBeneficiaries   []Beneficiary  `json:"life_beneficiaries"`
	HealthBeneficiaries []string       `json:"health_beneficiaries"` // Family members
	
	// Claims
	ClaimsHistory       []InsuranceClaim `json:"claims_history"`
	TotalClaimsAmount   sdk.Coin         `json:"total_claims_amount"`
	
	// Wellness programs
	WellnessScore       sdk.Dec          `json:"wellness_score"`
	HealthCheckups      uint32           `json:"health_checkups"`
	WellnessRewards     sdk.Coin         `json:"wellness_rewards"`
	
	Status              string           `json:"status"`
	CreatedAt           time.Time        `json:"created_at"`
	LastUpdated         time.Time        `json:"last_updated"`
}

// Beneficiary represents insurance beneficiary
type Beneficiary struct {
	Name         string         `json:"name"`
	Relationship string         `json:"relationship"`
	Address      sdk.AccAddress `json:"address,omitempty"`
	Share        sdk.Dec        `json:"share"`         // percentage
	ContactInfo  string         `json:"contact_info"`
}

// InsuranceClaim represents an insurance claim
type InsuranceClaim struct {
	ClaimID         string    `json:"claim_id"`
	ClaimType       string    `json:"claim_type"`       // "life", "health", "accident", "critical"
	ClaimAmount     sdk.Coin  `json:"claim_amount"`
	ClaimedAmount   sdk.Coin  `json:"claimed_amount"`
	ApprovedAmount  sdk.Coin  `json:"approved_amount"`
	ClaimDate       time.Time `json:"claim_date"`
	ProcessedDate   time.Time `json:"processed_date,omitempty"`
	Status          string    `json:"status"`
	Documents       []string  `json:"documents"`
	ProcessingNotes string    `json:"processing_notes,omitempty"`
}

// UrbanUnifiedPool represents the urban unified liquidity pool
type UrbanUnifiedPool struct {
	PoolID              uint64    `json:"pool_id"`
	PoolName            string    `json:"pool_name"`
	CityCode            string    `json:"city_code"`
	CityName            string    `json:"city_name"`
	
	// Pool composition
	TotalLiquidity      sdk.Coin  `json:"total_liquidity"`
	PensionReserve      sdk.Coin  `json:"pension_reserve"`       // 25% - For maturity payments
	EducationLoanPool   sdk.Coin  `json:"education_loan_pool"`   // 35% - For education loans
	InsuranceReserve    sdk.Coin  `json:"insurance_reserve"`     // 15% - For insurance claims
	InvestmentPool      sdk.Coin  `json:"investment_pool"`       // 20% - For yield generation
	EmergencyReserve    sdk.Coin  `json:"emergency_reserve"`     // 5% - Emergency fund
	
	// Performance metrics
	MonthlyInflow       sdk.Coin  `json:"monthly_inflow"`
	MonthlyOutflow      sdk.Coin  `json:"monthly_outflow"`
	LoanRevenue         sdk.Coin  `json:"loan_revenue"`          // Interest from education loans
	InvestmentReturns   sdk.Coin  `json:"investment_returns"`    // From investment pool
	InsurancePremiums   sdk.Coin  `json:"insurance_premiums"`    // Premium collections
	
	// Pool health
	UtilizationRatio    sdk.Dec   `json:"utilization_ratio"`
	DefaultRate         sdk.Dec   `json:"default_rate"`
	LiquidityRatio      sdk.Dec   `json:"liquidity_ratio"`
	PerformanceScore    sdk.Dec   `json:"performance_score"`
	
	// Members and activity
	ActiveMembers       uint32    `json:"active_members"`
	NewMembersThisMonth uint32    `json:"new_members_this_month"`
	EducationLoansActive uint32   `json:"education_loans_active"`
	InsurancePoliciesActive uint32 `json:"insurance_policies_active"`
	
	// Referral program
	ReferralRewardPool  sdk.Coin  `json:"referral_reward_pool"`
	MonthlyReferrals    uint32    `json:"monthly_referrals"`
	
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	LastUpdated         time.Time `json:"last_updated"`
}

// ReferralReward represents referral rewards structure
type ReferralReward struct {
	ReferrerAddress    sdk.AccAddress `json:"referrer_address"`
	RefereeAddress     sdk.AccAddress `json:"referee_address"`
	RewardType         string         `json:"reward_type"`        // "pension", "loan", "insurance"
	RewardAmount       sdk.Coin       `json:"reward_amount"`
	RewardPercentage   sdk.Dec        `json:"reward_percentage"`
	RewardDate         time.Time      `json:"reward_date"`
	RewardStatus       string         `json:"reward_status"`
	
	// Milestone rewards
	MilestoneLevel     uint32         `json:"milestone_level"`    // 1, 5, 10, 25, 50, 100 referrals
	MilestoneBonus     sdk.Coin       `json:"milestone_bonus"`
	
	// Performance-based rewards
	RefereePerformance sdk.Dec        `json:"referee_performance"` // Referee's consistency score
	PerformanceBonus   sdk.Coin       `json:"performance_bonus"`
}

// UrbanPensionConstants for scheme configuration
const (
	UrbanPensionMonthlyContribution = 2500 // ₹2500 equivalent in NAMO
	UrbanPensionContributionPeriod  = 18   // 18 months
	UrbanPensionMaturityMonth       = 20   // Returns on 20th month
	
	// Default coverage amounts (in INR equivalent)
	DefaultLifeCover     = 1000000  // ₹10 lakhs
	DefaultHealthCover   = 500000   // ₹5 lakhs
	DefaultAccidentCover = 500000   // ₹5 lakhs
	DefaultCriticalCover = 250000   // ₹2.5 lakhs
	
	// Interest rates
	MinEducationLoanRate = 4.0  // 4% for excellent pension history
	MaxEducationLoanRate = 7.0  // 7% for new members
	
	// Pool allocation percentages
	PensionReserveAllocation  = 25 // 25%
	EducationLoanAllocation   = 35 // 35%
	InsuranceReserveAllocation = 15 // 15%
	InvestmentPoolAllocation  = 20 // 20%
	EmergencyReserveAllocation = 5  // 5%
	
	// Referral rewards
	PensionReferralReward = 1000 // ₹1000 for pension referral
	LoanReferralReward    = 2000 // ₹2000 for education loan referral
	InsuranceReferralReward = 500 // ₹500 for insurance referral
)

// Scheme status constants
const (
	StatusActive     = "active"
	StatusMatured    = "matured"
	StatusDefaulted  = "defaulted"
	StatusSuspended  = "suspended"
	StatusClosed     = "closed"
)

// Education loan status
const (
	LoanStatusPending      = "pending"
	LoanStatusApproved     = "approved"
	LoanStatusDisbursing   = "disbursing"
	LoanStatusActive       = "active"
	LoanStatusCompleted    = "completed"
	LoanStatusDefaulted    = "defaulted"
)

// Insurance status
const (
	PolicyStatusActive    = "active"
	PolicyStatusLapsed    = "lapsed"
	PolicyStatusClaimed   = "claimed"
	PolicyStatusMatured   = "matured"
	PolicyStatusCancelled = "cancelled"
)

// Validate validates the UrbanPensionScheme
func (ups UrbanPensionScheme) Validate() error {
	if ups.SchemeID == "" {
		return ErrInvalidSchemeID
	}
	if ups.ContributorAddress.Empty() {
		return ErrInvalidContributor
	}
	if !ups.MonthlyContribution.IsPositive() {
		return ErrInvalidContribution
	}
	if ups.ContributionPeriod != UrbanPensionContributionPeriod {
		return ErrInvalidContributionPeriod
	}
	return nil
}

// CalculateExpectedReturn calculates expected return based on pool performance
func (ups UrbanPensionScheme) CalculateExpectedReturn(poolPerformance sdk.Dec) sdk.Coin {
	baseReturn := sdk.NewDecWithPrec(30, 2) // 30% base return
	performanceBonus := poolPerformance.Mul(sdk.NewDecWithPrec(10, 2)) // Up to 10% performance bonus
	totalReturnRate := baseReturn.Add(performanceBonus)
	
	// Cap at 45% maximum return
	maxReturn := sdk.NewDecWithPrec(45, 2)
	if totalReturnRate.GT(maxReturn) {
		totalReturnRate = maxReturn
	}
	
	// Calculate return amount
	returnAmount := ups.TotalContributions.Amount.ToDec().Mul(totalReturnRate).TruncateInt()
	return ups.TotalContributions.Add(sdk.NewCoin(ups.TotalContributions.Denom, returnAmount))
}

// CalculateEducationLoanRate calculates interest rate based on pension performance
func (ups UrbanPensionScheme) CalculateEducationLoanRate() sdk.Dec {
	// Better pension performance = lower interest rate
	baseRate := sdk.NewDecWithPrec(7, 2) // 7% base rate
	
	if ups.MissedContributions == 0 && ups.PaidContributions >= 6 {
		// Excellent track record - 4% rate
		return sdk.NewDecWithPrec(4, 2)
	} else if ups.MissedContributions <= 1 && ups.PaidContributions >= 3 {
		// Good track record - 5% rate
		return sdk.NewDecWithPrec(5, 2)
	} else if ups.MissedContributions <= 2 {
		// Average track record - 6% rate
		return sdk.NewDecWithPrec(6, 2)
	}
	
	// New member or poor track record - 7% rate
	return baseRate
}

// CalculateInsurancePremium calculates monthly insurance premium
func (ups UrbanPensionScheme) CalculateInsurancePremium(lifeCover, healthCover sdk.Coin) sdk.Coin {
	// Life insurance: 0.1% of cover amount annually
	lifePremium := lifeCover.Amount.ToDec().Mul(sdk.NewDecWithPrec(1, 3)).Quo(sdk.NewDec(12))
	
	// Health insurance: 3% of cover amount annually
	healthPremium := healthCover.Amount.ToDec().Mul(sdk.NewDecWithPrec(3, 2)).Quo(sdk.NewDec(12))
	
	totalPremium := lifePremium.Add(healthPremium).TruncateInt()
	return sdk.NewCoin(ups.MonthlyContribution.Denom, totalPremium)
}

// IsEligibleForEducationLoan checks eligibility for education loan
func (ups UrbanPensionScheme) IsEligibleForEducationLoan() bool {
	// Must have at least 3 successful contributions
	if ups.PaidContributions < 3 {
		return false
	}
	
	// Must not have more than 2 missed contributions
	if ups.MissedContributions > 2 {
		return false
	}
	
	// Must be active
	if ups.Status != StatusActive {
		return false
	}
	
	// Cannot already have an education loan
	if ups.EducationLoanTaken {
		return false
	}
	
	return true
}

// CalculateReferralReward calculates referral reward amount
func CalculateReferralReward(rewardType string, milestoneLevel uint32) sdk.Coin {
	var baseReward int64
	
	switch rewardType {
	case "pension":
		baseReward = PensionReferralReward
	case "loan":
		baseReward = LoanReferralReward
	case "insurance":
		baseReward = InsuranceReferralReward
	default:
		baseReward = PensionReferralReward
	}
	
	// Milestone multipliers
	multiplier := sdk.NewDec(1)
	if milestoneLevel >= 100 {
		multiplier = sdk.NewDec(3) // 3x for 100+ referrals
	} else if milestoneLevel >= 50 {
		multiplier = sdk.NewDecWithPrec(25, 1) // 2.5x for 50+ referrals
	} else if milestoneLevel >= 25 {
		multiplier = sdk.NewDec(2) // 2x for 25+ referrals
	} else if milestoneLevel >= 10 {
		multiplier = sdk.NewDecWithPrec(15, 1) // 1.5x for 10+ referrals
	} else if milestoneLevel >= 5 {
		multiplier = sdk.NewDecWithPrec(12, 1) // 1.2x for 5+ referrals
	}
	
	finalReward := sdk.NewDec(baseReward).Mul(multiplier).TruncateInt()
	return sdk.NewCoin("unamo", finalReward)
}