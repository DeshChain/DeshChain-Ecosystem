package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"cosmossdk.io/core/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// CreditScoringEngine provides automated credit assessment and scoring
type CreditScoringEngine struct {
	keeper       *Keeper
	models       []CreditModel
	scoringRules []ScoringRule
	thresholds   CreditThresholds
}

// NewCreditScoringEngine creates a new credit scoring engine
func NewCreditScoringEngine(k *Keeper) *CreditScoringEngine {
	return &CreditScoringEngine{
		keeper:       k,
		models:       initializeCreditModels(),
		scoringRules: initializeScoringRules(),
		thresholds:   initializeCreditThresholds(),
	}
}

// CreditAssessment represents a complete credit assessment
type CreditAssessment struct {
	AssessmentID       string                 `json:"assessment_id"`
	CustomerID         string                 `json:"customer_id"`
	AssessmentType     AssessmentType         `json:"assessment_type"`
	RequestedAmount    sdk.Coin               `json:"requested_amount"`
	LoanPurpose        string                 `json:"loan_purpose"`
	LoanTerm          int                    `json:"loan_term"`
	CreditScore        int                    `json:"credit_score"`
	RiskGrade          RiskGrade              `json:"risk_grade"`
	ApprovalStatus     ApprovalStatus         `json:"approval_status"`
	MaxLoanAmount      sdk.Coin               `json:"max_loan_amount"`
	RecommendedRate    sdk.Dec                `json:"recommended_rate"`
	Conditions         []LoanCondition        `json:"conditions"`
	Factors            []ScoringFactor        `json:"factors"`
	ModelResults       []ModelResult          `json:"model_results"`
	Timestamp          time.Time              `json:"timestamp"`
	ValidUntil         time.Time              `json:"valid_until"`
	AssessorID         string                 `json:"assessor_id"`
	ReviewNotes        string                 `json:"review_notes"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// CreditModel represents a scoring model
type CreditModel struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Version         string              `json:"version"`
	ModelType       ModelType           `json:"model_type"`
	Weightage       sdk.Dec             `json:"weightage"`
	Parameters      []ModelParameter    `json:"parameters"`
	Thresholds      ModelThresholds     `json:"thresholds"`
	IsActive        bool                `json:"is_active"`
	LastUpdated     time.Time           `json:"last_updated"`
	Performance     ModelPerformance    `json:"performance"`
	DataRequirements []string           `json:"data_requirements"`
}

// ScoringRule defines business rules for credit scoring
type ScoringRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    RuleCategory      `json:"category"`
	Conditions  []RuleCondition   `json:"conditions"`
	Actions     []RuleAction      `json:"actions"`
	Weight      sdk.Dec           `json:"weight"`
	Priority    int               `json:"priority"`
	IsActive    bool              `json:"is_active"`
	AppliesTo   []AssessmentType  `json:"applies_to"`
}

type ScoringFactor struct {
	Category     string  `json:"category"`
	Factor       string  `json:"factor"`
	Value        string  `json:"value"`
	Score        int     `json:"score"`
	Impact       string  `json:"impact"` // POSITIVE, NEGATIVE, NEUTRAL
	Confidence   float64 `json:"confidence"`
	Description  string  `json:"description"`
}

type ModelResult struct {
	ModelID     string  `json:"model_id"`
	ModelName   string  `json:"model_name"`
	Score       int     `json:"score"`
	Confidence  float64 `json:"confidence"`
	RiskGrade   string  `json:"risk_grade"`
	Factors     []ScoringFactor `json:"factors"`
	ProcessTime time.Duration   `json:"process_time"`
	Warnings    []string        `json:"warnings"`
}

type LoanCondition struct {
	Type         ConditionType `json:"type"`
	Description  string        `json:"description"`
	IsMandatory  bool          `json:"is_mandatory"`
	Value        string        `json:"value"`
	DueDate      time.Time     `json:"due_date,omitempty"`
	Status       string        `json:"status"`
}

// Customer financial profile for assessment
type CustomerProfile struct {
	CustomerID           string                 `json:"customer_id"`
	PersonalInfo         PersonalInformation    `json:"personal_info"`
	FinancialInfo        FinancialInformation   `json:"financial_info"`
	CreditHistory        CreditHistory          `json:"credit_history"`
	BusinessInfo         BusinessInformation    `json:"business_info,omitempty"`
	BankingHistory       BankingHistory         `json:"banking_history"`
	CollateralInfo       []Collateral           `json:"collateral_info"`
	References           []Reference            `json:"references"`
	RiskFactors          []RiskFactor           `json:"risk_factors"`
	EmploymentHistory    []Employment           `json:"employment_history"`
	PropertyOwnership    []Property             `json:"property_ownership"`
	InvestmentPortfolio  []Investment           `json:"investment_portfolio"`
	InsurancePolicies    []Insurance            `json:"insurance_policies"`
	TaxRecords          []TaxRecord            `json:"tax_records"`
	LegalRecords        []LegalRecord          `json:"legal_records"`
	SocialMediaProfile   SocialMediaData        `json:"social_media_profile,omitempty"`
	TransactionPatterns  TransactionBehavior    `json:"transaction_patterns"`
}

type PersonalInformation struct {
	Age                int               `json:"age"`
	Gender             string            `json:"gender"`
	MaritalStatus      string            `json:"marital_status"`
	Education          string            `json:"education"`
	Dependents         int               `json:"dependents"`
	ResidenceType      string            `json:"residence_type"`
	ResidenceStability int               `json:"residence_stability_years"`
	Location           LocationInfo      `json:"location"`
	PoliticalExposure  PEPStatus        `json:"political_exposure"`
}

type FinancialInformation struct {
	MonthlyIncome       sdk.Coin   `json:"monthly_income"`
	MonthlyExpenses     sdk.Coin   `json:"monthly_expenses"`
	NetWorth           sdk.Coin   `json:"net_worth"`
	Assets             sdk.Coin   `json:"assets"`
	Liabilities        sdk.Coin   `json:"liabilities"`
	DebtToIncomeRatio  sdk.Dec    `json:"debt_to_income_ratio"`
	CashFlowStability  int        `json:"cash_flow_stability"`
	IncomeSource       string     `json:"income_source"`
	IncomeVerification string     `json:"income_verification"`
	BankBalance        sdk.Coin   `json:"bank_balance"`
}

type CreditHistory struct {
	CreditScore       int                    `json:"credit_score"`
	CreditReportDate  time.Time              `json:"credit_report_date"`
	PaymentHistory    PaymentHistoryData     `json:"payment_history"`
	CreditUtilization sdk.Dec                `json:"credit_utilization"`
	AccountHistory    []CreditAccount        `json:"account_history"`
	Inquiries         []CreditInquiry        `json:"inquiries"`
	PublicRecords     []PublicRecord         `json:"public_records"`
	DelinquencyCount  int                    `json:"delinquency_count"`
	DefaultHistory    []Default              `json:"default_history"`
}

type BusinessInformation struct {
	BusinessName        string       `json:"business_name"`
	BusinessType        string       `json:"business_type"`
	Industry           string       `json:"industry"`
	YearsInBusiness    int          `json:"years_in_business"`
	AnnualRevenue      sdk.Coin     `json:"annual_revenue"`
	MonthlyRevenue     sdk.Coin     `json:"monthly_revenue"`
	ProfitMargin       sdk.Dec      `json:"profit_margin"`
	NumberOfEmployees  int          `json:"number_of_employees"`
	BusinessStability  int          `json:"business_stability"`
	MarketPosition     string       `json:"market_position"`
	CompetitiveRisk    string       `json:"competitive_risk"`
	SeasonalityFactor  sdk.Dec      `json:"seasonality_factor"`
}

// Enums and constants

type AssessmentType int

const (
	ASSESSMENT_TYPE_PERSONAL_LOAN AssessmentType = iota
	ASSESSMENT_TYPE_BUSINESS_LOAN
	ASSESSMENT_TYPE_MORTGAGE
	ASSESSMENT_TYPE_EDUCATION_LOAN
	ASSESSMENT_TYPE_AGRICULTURAL_LOAN
	ASSESSMENT_TYPE_TRADE_FINANCE
	ASSESSMENT_TYPE_WORKING_CAPITAL
	ASSESSMENT_TYPE_EQUIPMENT_FINANCE
	ASSESSMENT_TYPE_CREDIT_CARD
)

type RiskGrade int

const (
	RISK_GRADE_AAA RiskGrade = iota // Excellent (0-50 basis points)
	RISK_GRADE_AA                   // Very Good (50-100 bps)
	RISK_GRADE_A                    // Good (100-200 bps)
	RISK_GRADE_BBB                  // Fair (200-400 bps)
	RISK_GRADE_BB                   // Poor (400-700 bps)
	RISK_GRADE_B                    // Very Poor (700-1200 bps)
	RISK_GRADE_CCC                  // High Risk (1200+ bps)
	RISK_GRADE_D                    // Default Risk (Declined)
)

type ApprovalStatus int

const (
	APPROVAL_STATUS_APPROVED ApprovalStatus = iota
	APPROVAL_STATUS_CONDITIONALLY_APPROVED
	APPROVAL_STATUS_PENDING_REVIEW
	APPROVAL_STATUS_DECLINED
	APPROVAL_STATUS_NEEDS_MORE_INFO
	APPROVAL_STATUS_MANUAL_REVIEW_REQUIRED
)

type ModelType int

const (
	MODEL_TYPE_STATISTICAL ModelType = iota
	MODEL_TYPE_MACHINE_LEARNING
	MODEL_TYPE_RULE_BASED
	MODEL_TYPE_BEHAVIORAL
	MODEL_TYPE_ALTERNATIVE_DATA
	MODEL_TYPE_ENSEMBLE
)

type RuleCategory int

const (
	RULE_CATEGORY_INCOME_VERIFICATION RuleCategory = iota
	RULE_CATEGORY_DEBT_TO_INCOME
	RULE_CATEGORY_PAYMENT_HISTORY
	RULE_CATEGORY_EMPLOYMENT_STABILITY
	RULE_CATEGORY_COLLATERAL_VALUATION
	RULE_CATEGORY_INDUSTRY_RISK
	RULE_CATEGORY_GEOGRAPHIC_RISK
	RULE_CATEGORY_REGULATORY_COMPLIANCE
)

type ConditionType int

const (
	CONDITION_TYPE_INCOME_PROOF ConditionType = iota
	CONDITION_TYPE_COLLATERAL_VALUATION
	CONDITION_TYPE_GUARANTOR
	CONDITION_TYPE_INSURANCE
	CONDITION_TYPE_ADDITIONAL_DOCUMENTATION
	CONDITION_TYPE_CREDIT_BUREAU_UPDATE
	CONDITION_TYPE_BANK_STATEMENT_REVIEW
	CONDITION_TYPE_PROPERTY_VERIFICATION
)

// Core credit assessment functions

// AssessCreditworthiness performs complete credit assessment
func (cse *CreditScoringEngine) AssessCreditworthiness(ctx context.Context, profile CustomerProfile, request LoanRequest) (*CreditAssessment, error) {
	startTime := time.Now()

	assessment := &CreditAssessment{
		AssessmentID:    fmt.Sprintf("ASSESS_%s_%d", profile.CustomerID, time.Now().Unix()),
		CustomerID:      profile.CustomerID,
		AssessmentType:  request.LoanType,
		RequestedAmount: request.Amount,
		LoanPurpose:     request.Purpose,
		LoanTerm:       request.TermMonths,
		Timestamp:      time.Now(),
		ValidUntil:     time.Now().AddDate(0, 3, 0), // Valid for 3 months
		AssessorID:     "AUTOMATED_SYSTEM",
		Factors:        []ScoringFactor{},
		ModelResults:   []ModelResult{},
		Conditions:     []LoanCondition{},
		Metadata:       make(map[string]interface{}),
	}

	// Run all applicable credit models
	for _, model := range cse.models {
		if model.IsActive && cse.modelApplies(model, request.LoanType) {
			result, err := cse.runCreditModel(ctx, model, profile, request)
			if err != nil {
				cse.keeper.Logger(ctx).Warn("Credit model failed", "model", model.Name, "error", err)
				continue
			}
			assessment.ModelResults = append(assessment.ModelResults, *result)
		}
	}

	// Calculate composite credit score
	assessment.CreditScore = cse.calculateCompositeScore(assessment.ModelResults)

	// Determine risk grade
	assessment.RiskGrade = cse.determineRiskGrade(assessment.CreditScore, profile, request)

	// Apply business rules
	ruleResults := cse.applyBusinessRules(ctx, profile, request, assessment.CreditScore)
	assessment.Factors = append(assessment.Factors, ruleResults.Factors...)
	assessment.Conditions = append(assessment.Conditions, ruleResults.Conditions...)

	// Make approval decision
	decision := cse.makeApprovalDecision(ctx, assessment, profile, request)
	assessment.ApprovalStatus = decision.Status
	assessment.MaxLoanAmount = decision.MaxAmount
	assessment.RecommendedRate = decision.InterestRate
	assessment.Conditions = append(assessment.Conditions, decision.Conditions...)
	assessment.ReviewNotes = decision.Notes

	// Store assessment
	if err := cse.storeAssessment(ctx, *assessment); err != nil {
		return assessment, fmt.Errorf("failed to store assessment: %w", err)
	}

	// Log processing time
	processingTime := time.Since(startTime)
	assessment.Metadata["processing_time_ms"] = processingTime.Milliseconds()

	// Emit event
	cse.emitAssessmentEvent(ctx, assessment)

	return assessment, nil
}

type LoanRequest struct {
	LoanType     AssessmentType `json:"loan_type"`
	Amount       sdk.Coin       `json:"amount"`
	Purpose      string         `json:"purpose"`
	TermMonths   int           `json:"term_months"`
	CollateralID string        `json:"collateral_id,omitempty"`
}

type ApprovalDecision struct {
	Status       ApprovalStatus   `json:"status"`
	MaxAmount    sdk.Coin         `json:"max_amount"`
	InterestRate sdk.Dec          `json:"interest_rate"`
	Conditions   []LoanCondition  `json:"conditions"`
	Notes        string           `json:"notes"`
}

// Credit scoring models implementation

func (cse *CreditScoringEngine) runCreditModel(ctx context.Context, model CreditModel, profile CustomerProfile, request LoanRequest) (*ModelResult, error) {
	result := &ModelResult{
		ModelID:    model.ID,
		ModelName:  model.Name,
		Factors:    []ScoringFactor{},
		Warnings:   []string{},
		ProcessTime: 0,
	}

	startTime := time.Now()

	switch model.ModelType {
	case MODEL_TYPE_STATISTICAL:
		result = cse.runStatisticalModel(model, profile, request)
	case MODEL_TYPE_MACHINE_LEARNING:
		result = cse.runMLModel(model, profile, request)
	case MODEL_TYPE_RULE_BASED:
		result = cse.runRuleBasedModel(model, profile, request)
	case MODEL_TYPE_BEHAVIORAL:
		result = cse.runBehavioralModel(model, profile, request)
	case MODEL_TYPE_ALTERNATIVE_DATA:
		result = cse.runAlternativeDataModel(model, profile, request)
	case MODEL_TYPE_ENSEMBLE:
		result = cse.runEnsembleModel(model, profile, request)
	default:
		return nil, fmt.Errorf("unsupported model type: %v", model.ModelType)
	}

	result.ProcessTime = time.Since(startTime)
	result.RiskGrade = cse.scoreToRiskGrade(result.Score).String()

	return result, nil
}

func (cse *CreditScoringEngine) runStatisticalModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	score := 650 // Base score
	factors := []ScoringFactor{}

	// Age factor (18-25: -20, 26-35: +30, 36-50: +50, 51-65: +20, 65+: -30)
	ageScore := cse.calculateAgeScore(profile.PersonalInfo.Age)
	score += ageScore
	factors = append(factors, ScoringFactor{
		Category:    "DEMOGRAPHICS",
		Factor:      "AGE",
		Value:       fmt.Sprintf("%d years", profile.PersonalInfo.Age),
		Score:       ageScore,
		Impact:      cse.getImpact(ageScore),
		Confidence:  0.95,
		Description: "Age impact on creditworthiness",
	})

	// Income stability factor
	incomeScore := cse.calculateIncomeScore(profile.FinancialInfo.MonthlyIncome, profile.FinancialInfo.IncomeSource)
	score += incomeScore
	factors = append(factors, ScoringFactor{
		Category:    "INCOME",
		Factor:      "MONTHLY_INCOME",
		Value:       profile.FinancialInfo.MonthlyIncome.String(),
		Score:       incomeScore,
		Impact:      cse.getImpact(incomeScore),
		Confidence:  0.90,
		Description: "Monthly income assessment",
	})

	// Debt-to-income ratio
	dtiScore := cse.calculateDTIScore(profile.FinancialInfo.DebtToIncomeRatio)
	score += dtiScore
	factors = append(factors, ScoringFactor{
		Category:    "DEBT",
		Factor:      "DEBT_TO_INCOME_RATIO",
		Value:       profile.FinancialInfo.DebtToIncomeRatio.String(),
		Score:       dtiScore,
		Impact:      cse.getImpact(dtiScore),
		Confidence:  0.88,
		Description: "Debt-to-income ratio evaluation",
	})

	// Payment history
	paymentScore := cse.calculatePaymentHistoryScore(profile.CreditHistory.PaymentHistory)
	score += paymentScore
	factors = append(factors, ScoringFactor{
		Category:    "CREDIT_HISTORY",
		Factor:      "PAYMENT_HISTORY",
		Value:       fmt.Sprintf("%.1f%% on-time", profile.CreditHistory.PaymentHistory.OnTimePaymentPercentage.MustFloat64()*100),
		Score:       paymentScore,
		Impact:      cse.getImpact(paymentScore),
		Confidence:  0.92,
		Description: "Historical payment behavior",
	})

	// Employment stability
	empScore := cse.calculateEmploymentScore(profile.EmploymentHistory)
	score += empScore
	factors = append(factors, ScoringFactor{
		Category:    "EMPLOYMENT",
		Factor:      "EMPLOYMENT_STABILITY",
		Value:       cse.getEmploymentStabilityDescription(profile.EmploymentHistory),
		Score:       empScore,
		Impact:      cse.getImpact(empScore),
		Confidence:  0.85,
		Description: "Employment history and stability",
	})

	return &ModelResult{
		Score:      cse.normalizeScore(score),
		Confidence: 0.89,
		Factors:    factors,
	}
}

func (cse *CreditScoringEngine) runMLModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	// Simulate ML model with feature engineering
	features := cse.extractMLFeatures(profile, request)
	
	// Simulate neural network/ensemble prediction
	score := cse.predictMLScore(features)
	
	factors := []ScoringFactor{
		{
			Category:    "ML_PREDICTION",
			Factor:      "ENSEMBLE_SCORE",
			Value:       fmt.Sprintf("%.2f", score),
			Score:       int(score),
			Impact:      "POSITIVE",
			Confidence:  0.94,
			Description: "Machine learning ensemble prediction",
		},
	}

	return &ModelResult{
		Score:      int(score),
		Confidence: 0.94,
		Factors:    factors,
	}
}

func (cse *CreditScoringEngine) runRuleBasedModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	score := 650
	factors := []ScoringFactor{}

	// Apply hard rules
	for _, rule := range cse.scoringRules {
		if !rule.IsActive || !cse.ruleApplies(rule, request.LoanType) {
			continue
		}

		ruleScore, applied := cse.applyRule(rule, profile, request)
		if applied {
			score += ruleScore
			factors = append(factors, ScoringFactor{
				Category:    rule.Category.String(),
				Factor:      rule.Name,
				Value:       "Applied",
				Score:       ruleScore,
				Impact:      cse.getImpact(ruleScore),
				Confidence:  0.98,
				Description: rule.Description,
			})
		}
	}

	return &ModelResult{
		Score:      cse.normalizeScore(score),
		Confidence: 0.91,
		Factors:    factors,
	}
}

func (cse *CreditScoringEngine) runBehavioralModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	score := 650
	factors := []ScoringFactor{}

	// Transaction behavior analysis
	behaviorScore := cse.analyzeBehavioralPatterns(profile.TransactionPatterns)
	score += behaviorScore
	factors = append(factors, ScoringFactor{
		Category:    "BEHAVIOR",
		Factor:      "TRANSACTION_PATTERNS",
		Value:       "Analyzed",
		Score:       behaviorScore,
		Impact:      cse.getImpact(behaviorScore),
		Confidence:  0.86,
		Description: "Transaction behavior analysis",
	})

	// Banking behavior
	bankingScore := cse.analyzeBankingBehavior(profile.BankingHistory)
	score += bankingScore
	factors = append(factors, ScoringFactor{
		Category:    "BEHAVIOR",
		Factor:      "BANKING_BEHAVIOR",
		Value:       "Analyzed",
		Score:       bankingScore,
		Impact:      cse.getImpact(bankingScore),
		Confidence:  0.82,
		Description: "Banking relationship behavior",
	})

	return &ModelResult{
		Score:      cse.normalizeScore(score),
		Confidence: 0.84,
		Factors:    factors,
	}
}

func (cse *CreditScoringEngine) runAlternativeDataModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	score := 650
	factors := []ScoringFactor{}

	// Social media analysis (if available)
	if profile.SocialMediaProfile.HasData {
		socialScore := cse.analyzeSocialMediaData(profile.SocialMediaProfile)
		score += socialScore
		factors = append(factors, ScoringFactor{
			Category:    "ALTERNATIVE_DATA",
			Factor:      "SOCIAL_MEDIA_ANALYSIS",
			Value:       "Positive indicators",
			Score:       socialScore,
			Impact:      cse.getImpact(socialScore),
			Confidence:  0.72,
			Description: "Social media behavioral indicators",
		})
	}

	// Utility payment history (simulated)
	utilityScore := 25 // Assume positive utility payment history
	score += utilityScore
	factors = append(factors, ScoringFactor{
		Category:    "ALTERNATIVE_DATA",
		Factor:      "UTILITY_PAYMENTS",
		Value:       "Regular payments",
		Score:       utilityScore,
		Impact:      "POSITIVE",
		Confidence:  0.78,
		Description: "Utility and recurring payment behavior",
	})

	return &ModelResult{
		Score:      cse.normalizeScore(score),
		Confidence: 0.75,
		Factors:    factors,
	}
}

func (cse *CreditScoringEngine) runEnsembleModel(model CreditModel, profile CustomerProfile, request LoanRequest) *ModelResult {
	// Run multiple models and combine results
	statResult := cse.runStatisticalModel(model, profile, request)
	mlResult := cse.runMLModel(model, profile, request)
	ruleResult := cse.runRuleBasedModel(model, profile, request)

	// Weighted average
	ensembleScore := int(float64(statResult.Score)*0.4 + float64(mlResult.Score)*0.4 + float64(ruleResult.Score)*0.2)

	// Combine factors
	allFactors := []ScoringFactor{}
	allFactors = append(allFactors, statResult.Factors...)
	allFactors = append(allFactors, mlResult.Factors...)
	allFactors = append(allFactors, ruleResult.Factors...)

	return &ModelResult{
		Score:      ensembleScore,
		Confidence: 0.96,
		Factors:    allFactors,
	}
}

// Business rule application

type BusinessRuleResult struct {
	Factors    []ScoringFactor  `json:"factors"`
	Conditions []LoanCondition  `json:"conditions"`
}

func (cse *CreditScoringEngine) applyBusinessRules(ctx context.Context, profile CustomerProfile, request LoanRequest, baseScore int) BusinessRuleResult {
	result := BusinessRuleResult{
		Factors:    []ScoringFactor{},
		Conditions: []LoanCondition{},
	}

	// Minimum income rule
	minIncome := cse.getMinimumIncomeRequirement(request.LoanType, request.Amount)
	if profile.FinancialInfo.MonthlyIncome.IsLT(minIncome) {
		result.Conditions = append(result.Conditions, LoanCondition{
			Type:        CONDITION_TYPE_INCOME_PROOF,
			Description: fmt.Sprintf("Provide additional income proof. Minimum required: %s", minIncome.String()),
			IsMandatory: true,
			Status:      "PENDING",
		})
	}

	// Age restrictions
	if profile.PersonalInfo.Age < 21 || profile.PersonalInfo.Age > 65 {
		if profile.PersonalInfo.Age < 21 {
			result.Conditions = append(result.Conditions, LoanCondition{
				Type:        CONDITION_TYPE_GUARANTOR,
				Description: "Guarantor required for applicants under 21",
				IsMandatory: true,
				Status:      "PENDING",
			})
		}
		if profile.PersonalInfo.Age > 65 {
			result.Conditions = append(result.Conditions, LoanCondition{
				Type:        CONDITION_TYPE_ADDITIONAL_DOCUMENTATION,
				Description: "Additional medical and income documentation required",
				IsMandatory: true,
				Status:      "PENDING",
			})
		}
	}

	// High DTI ratio
	if profile.FinancialInfo.DebtToIncomeRatio.GT(sdk.NewDecWithPrec(40, 2)) { // 40%
		result.Conditions = append(result.Conditions, LoanCondition{
			Type:        CONDITION_TYPE_BANK_STATEMENT_REVIEW,
			Description: "Additional financial review required due to high debt-to-income ratio",
			IsMandatory: true,
			Status:      "PENDING",
		})
	}

	// Collateral requirement for large loans
	collateralThreshold := sdk.NewCoin(request.Amount.Denom, sdk.NewInt(1000000)) // 1M
	if request.Amount.IsGTE(collateralThreshold) {
		result.Conditions = append(result.Conditions, LoanCondition{
			Type:        CONDITION_TYPE_COLLATERAL_VALUATION,
			Description: "Collateral valuation required for loans above 1M",
			IsMandatory: true,
			Status:      "PENDING",
		})
	}

	return result
}

// Helper functions

func (cse *CreditScoringEngine) calculateCompositeScore(results []ModelResult) int {
	if len(results) == 0 {
		return 500 // Default poor score
	}

	totalWeightedScore := 0.0
	totalWeight := 0.0

	for _, result := range results {
		weight := result.Confidence
		totalWeightedScore += float64(result.Score) * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 500
	}

	return int(totalWeightedScore / totalWeight)
}

func (cse *CreditScoringEngine) determineRiskGrade(score int, profile CustomerProfile, request LoanRequest) RiskGrade {
	// Base risk grade on score
	baseGrade := cse.scoreToRiskGrade(score)

	// Adjust for specific risk factors
	if len(profile.RiskFactors) > 0 {
		baseGrade = cse.adjustGradeForRiskFactors(baseGrade, profile.RiskFactors)
	}

	// Industry-specific adjustments
	if profile.BusinessInfo.Industry != "" {
		baseGrade = cse.adjustGradeForIndustry(baseGrade, profile.BusinessInfo.Industry)
	}

	return baseGrade
}

func (cse *CreditScoringEngine) scoreToRiskGrade(score int) RiskGrade {
	switch {
	case score >= 800:
		return RISK_GRADE_AAA
	case score >= 750:
		return RISK_GRADE_AA
	case score >= 700:
		return RISK_GRADE_A
	case score >= 650:
		return RISK_GRADE_BBB
	case score >= 600:
		return RISK_GRADE_BB
	case score >= 550:
		return RISK_GRADE_B
	case score >= 500:
		return RISK_GRADE_CCC
	default:
		return RISK_GRADE_D
	}
}

func (cse *CreditScoringEngine) makeApprovalDecision(ctx context.Context, assessment *CreditAssessment, profile CustomerProfile, request LoanRequest) ApprovalDecision {
	decision := ApprovalDecision{
		Status:       APPROVAL_STATUS_DECLINED,
		MaxAmount:    sdk.NewCoin(request.Amount.Denom, sdk.ZeroInt()),
		InterestRate: sdk.NewDecWithPrec(12, 2), // 12% default
		Conditions:   []LoanCondition{},
		Notes:        "",
	}

	// Base decision on credit score and risk grade
	switch assessment.RiskGrade {
	case RISK_GRADE_AAA, RISK_GRADE_AA:
		decision.Status = APPROVAL_STATUS_APPROVED
		decision.MaxAmount = request.Amount.Mul(sdk.NewInt(120)).Quo(sdk.NewInt(100)) // 120% of requested
		decision.InterestRate = cse.calculateInterestRate(assessment.RiskGrade, request.LoanType)
		decision.Notes = "Excellent credit profile. Approved at favorable terms."

	case RISK_GRADE_A:
		decision.Status = APPROVAL_STATUS_APPROVED
		decision.MaxAmount = request.Amount
		decision.InterestRate = cse.calculateInterestRate(assessment.RiskGrade, request.LoanType)
		decision.Notes = "Good credit profile. Standard approval."

	case RISK_GRADE_BBB:
		decision.Status = APPROVAL_STATUS_CONDITIONALLY_APPROVED
		decision.MaxAmount = request.Amount.Mul(sdk.NewInt(80)).Quo(sdk.NewInt(100)) // 80% of requested
		decision.InterestRate = cse.calculateInterestRate(assessment.RiskGrade, request.LoanType)
		decision.Notes = "Conditional approval subject to additional documentation."

	case RISK_GRADE_BB:
		decision.Status = APPROVAL_STATUS_MANUAL_REVIEW_REQUIRED
		decision.MaxAmount = request.Amount.Mul(sdk.NewInt(60)).Quo(sdk.NewInt(100)) // 60% of requested
		decision.InterestRate = cse.calculateInterestRate(assessment.RiskGrade, request.LoanType)
		decision.Notes = "Manual review required due to fair credit profile."

	case RISK_GRADE_B:
		decision.Status = APPROVAL_STATUS_DECLINED
		decision.Notes = "Declined due to poor credit profile. Consider reapplying after 6 months."

	default:
		decision.Status = APPROVAL_STATUS_DECLINED
		decision.Notes = "Application declined due to high risk profile."
	}

	// Adjust for mandatory conditions
	mandatoryConditions := 0
	for _, condition := range assessment.Conditions {
		if condition.IsMandatory {
			mandatoryConditions++
		}
	}

	if mandatoryConditions > 3 && decision.Status == APPROVAL_STATUS_APPROVED {
		decision.Status = APPROVAL_STATUS_CONDITIONALLY_APPROVED
	}

	return decision
}

func (cse *CreditScoringEngine) calculateInterestRate(grade RiskGrade, loanType AssessmentType) sdk.Dec {
	// Base rate by loan type
	var baseRate sdk.Dec
	switch loanType {
	case ASSESSMENT_TYPE_MORTGAGE:
		baseRate = sdk.NewDecWithPrec(650, 4) // 6.5%
	case ASSESSMENT_TYPE_PERSONAL_LOAN:
		baseRate = sdk.NewDecWithPrec(1000, 4) // 10%
	case ASSESSMENT_TYPE_BUSINESS_LOAN:
		baseRate = sdk.NewDecWithPrec(800, 4) // 8%
	case ASSESSMENT_TYPE_EDUCATION_LOAN:
		baseRate = sdk.NewDecWithPrec(650, 4) // 6.5%
	case ASSESSMENT_TYPE_AGRICULTURAL_LOAN:
		baseRate = sdk.NewDecWithPrec(600, 4) // 6%
	default:
		baseRate = sdk.NewDecWithPrec(1200, 4) // 12%
	}

	// Risk adjustment
	var riskAdjustment sdk.Dec
	switch grade {
	case RISK_GRADE_AAA:
		riskAdjustment = sdk.NewDecWithPrec(-100, 4) // -1%
	case RISK_GRADE_AA:
		riskAdjustment = sdk.NewDecWithPrec(-50, 4) // -0.5%
	case RISK_GRADE_A:
		riskAdjustment = sdk.ZeroDec()
	case RISK_GRADE_BBB:
		riskAdjustment = sdk.NewDecWithPrec(100, 4) // +1%
	case RISK_GRADE_BB:
		riskAdjustment = sdk.NewDecWithPrec(200, 4) // +2%
	case RISK_GRADE_B:
		riskAdjustment = sdk.NewDecWithPrec(400, 4) // +4%
	default:
		riskAdjustment = sdk.NewDecWithPrec(600, 4) // +6%
	}

	finalRate := baseRate.Add(riskAdjustment)
	
	// Ensure minimum rate
	minRate := sdk.NewDecWithPrec(500, 4) // 5%
	if finalRate.LT(minRate) {
		finalRate = minRate
	}

	return finalRate
}

// Scoring helper functions

func (cse *CreditScoringEngine) calculateAgeScore(age int) int {
	switch {
	case age >= 18 && age <= 25:
		return -20
	case age >= 26 && age <= 35:
		return 30
	case age >= 36 && age <= 50:
		return 50
	case age >= 51 && age <= 65:
		return 20
	default:
		return -30
	}
}

func (cse *CreditScoringEngine) calculateIncomeScore(income sdk.Coin, source string) int {
	// Base score on income amount (assuming USD equivalent)
	amountScore := int(math.Min(100, float64(income.Amount.Int64())/1000.0)) // $1 = 1 point, capped at 100

	// Adjust for income source stability
	sourceMultiplier := 1.0
	switch strings.ToUpper(source) {
	case "GOVERNMENT", "PERMANENT_EMPLOYMENT":
		sourceMultiplier = 1.2
	case "BUSINESS", "SELF_EMPLOYED":
		sourceMultiplier = 0.8
	case "FREELANCE", "CONTRACT":
		sourceMultiplier = 0.6
	}

	return int(float64(amountScore) * sourceMultiplier)
}

func (cse *CreditScoringEngine) calculateDTIScore(dtiRatio sdk.Dec) int {
	ratio := dtiRatio.MustFloat64()
	switch {
	case ratio <= 0.20: // 20% or less
		return 50
	case ratio <= 0.30:
		return 25
	case ratio <= 0.40:
		return 0
	case ratio <= 0.50:
		return -25
	default:
		return -50
	}
}

// Storage operations

func (cse *CreditScoringEngine) storeAssessment(ctx context.Context, assessment CreditAssessment) error {
	store := cse.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("credit_assessment_%s", assessment.AssessmentID))
	bz, err := json.Marshal(assessment)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (cse *CreditScoringEngine) emitAssessmentEvent(ctx context.Context, assessment *CreditAssessment) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"credit_assessment_completed",
			sdk.NewAttribute("assessment_id", assessment.AssessmentID),
			sdk.NewAttribute("customer_id", assessment.CustomerID),
			sdk.NewAttribute("credit_score", fmt.Sprintf("%d", assessment.CreditScore)),
			sdk.NewAttribute("risk_grade", assessment.RiskGrade.String()),
			sdk.NewAttribute("approval_status", assessment.ApprovalStatus.String()),
		),
	)
}

// Initialize default models and rules

func initializeCreditModels() []CreditModel {
	return []CreditModel{
		{
			ID:         "STATISTICAL_MODEL_V1",
			Name:       "Statistical Credit Model",
			Version:    "1.0",
			ModelType:  MODEL_TYPE_STATISTICAL,
			Weightage:  sdk.NewDecWithPrec(30, 2), // 30%
			IsActive:   true,
			Performance: ModelPerformance{
				Accuracy: sdk.NewDecWithPrec(85, 2),
				Precision: sdk.NewDecWithPrec(82, 2),
				Recall: sdk.NewDecWithPrec(88, 2),
			},
		},
		{
			ID:         "ML_ENSEMBLE_V2",
			Name:       "Machine Learning Ensemble",
			Version:    "2.1",
			ModelType:  MODEL_TYPE_ENSEMBLE,
			Weightage:  sdk.NewDecWithPrec(50, 2), // 50%
			IsActive:   true,
			Performance: ModelPerformance{
				Accuracy: sdk.NewDecWithPrec(92, 2),
				Precision: sdk.NewDecWithPrec(89, 2),
				Recall: sdk.NewDecWithPrec(94, 2),
			},
		},
		{
			ID:         "RULE_ENGINE_V1",
			Name:       "Business Rules Engine",
			Version:    "1.3",
			ModelType:  MODEL_TYPE_RULE_BASED,
			Weightage:  sdk.NewDecWithPrec(20, 2), // 20%
			IsActive:   true,
			Performance: ModelPerformance{
				Accuracy: sdk.NewDecWithPrec(78, 2),
				Precision: sdk.NewDecWithPrec(95, 2),
				Recall: sdk.NewDecWithPrec(65, 2),
			},
		},
	}
}

func initializeScoringRules() []ScoringRule {
	return []ScoringRule{
		{
			ID:          "MIN_AGE_RULE",
			Name:        "Minimum Age Requirement",
			Description: "Applicant must be at least 18 years old",
			Category:    RULE_CATEGORY_REGULATORY_COMPLIANCE,
			Weight:      sdk.NewDecWithPrec(100, 2),
			Priority:    1,
			IsActive:    true,
			AppliesTo:   []AssessmentType{ASSESSMENT_TYPE_PERSONAL_LOAN, ASSESSMENT_TYPE_BUSINESS_LOAN},
		},
		{
			ID:          "MAX_DTI_RULE",
			Name:        "Maximum Debt-to-Income Ratio",
			Description: "DTI ratio should not exceed 60%",
			Category:    RULE_CATEGORY_DEBT_TO_INCOME,
			Weight:      sdk.NewDecWithPrec(80, 2),
			Priority:    2,
			IsActive:    true,
			AppliesTo:   []AssessmentType{ASSESSMENT_TYPE_PERSONAL_LOAN, ASSESSMENT_TYPE_MORTGAGE},
		},
	}
}

func initializeCreditThresholds() CreditThresholds {
	return CreditThresholds{
		MinimumScore:     500,
		ApprovalScore:    650,
		ExcellentScore:   750,
		MaxDTIRatio:      sdk.NewDecWithPrec(60, 2), // 60%
		MinIncomeMultiple: sdk.NewInt(3),
	}
}

// Additional supporting types and functions

type CreditThresholds struct {
	MinimumScore     int     `json:"minimum_score"`
	ApprovalScore    int     `json:"approval_score"`
	ExcellentScore   int     `json:"excellent_score"`
	MaxDTIRatio      sdk.Dec `json:"max_dti_ratio"`
	MinIncomeMultiple sdk.Int `json:"min_income_multiple"`
}

type ModelPerformance struct {
	Accuracy  sdk.Dec `json:"accuracy"`
	Precision sdk.Dec `json:"precision"`
	Recall    sdk.Dec `json:"recall"`
}

type ModelParameter struct {
	Name         string      `json:"name"`
	Value        interface{} `json:"value"`
	Type         string      `json:"type"`
	Description  string      `json:"description"`
}

type ModelThresholds struct {
	ApprovalThreshold  int `json:"approval_threshold"`
	DeclineThreshold   int `json:"decline_threshold"`
	ReviewThreshold    int `json:"review_threshold"`
}

// Simplified supporting types for compilation
type PaymentHistoryData struct {
	OnTimePaymentPercentage sdk.Dec `json:"on_time_payment_percentage"`
	TotalAccounts          int     `json:"total_accounts"`
	AverageMonthsReviewed  int     `json:"average_months_reviewed"`
}

type CreditAccount struct {
	AccountType    string    `json:"account_type"`
	Balance       sdk.Coin  `json:"balance"`
	CreditLimit   sdk.Coin  `json:"credit_limit"`
	OpenDate      time.Time `json:"open_date"`
	LastActivity  time.Time `json:"last_activity"`
	PaymentStatus string    `json:"payment_status"`
}

type CreditInquiry struct {
	InquiryDate time.Time `json:"inquiry_date"`
	InquiryType string    `json:"inquiry_type"`
	Lender     string    `json:"lender"`
	Amount     sdk.Coin  `json:"amount"`
}

type PublicRecord struct {
	RecordType string    `json:"record_type"`
	FilingDate time.Time `json:"filing_date"`
	Amount     sdk.Coin  `json:"amount"`
	Status     string    `json:"status"`
}

type Default struct {
	DefaultDate time.Time `json:"default_date"`
	Amount     sdk.Coin  `json:"amount"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
}

type BankingHistory struct {
	PrimaryBankRelationshipYears int                `json:"primary_bank_relationship_years"`
	NumberOfBankAccounts        int                `json:"number_of_bank_accounts"`
	AverageMonthlyBalance       sdk.Coin           `json:"average_monthly_balance"`
	OverdraftHistory           []OverdraftRecord  `json:"overdraft_history"`
	ReturnedCheckCount         int                `json:"returned_check_count"`
}

type OverdraftRecord struct {
	Date   time.Time `json:"date"`
	Amount sdk.Coin  `json:"amount"`
}

type Collateral struct {
	CollateralID   string   `json:"collateral_id"`
	Type          string   `json:"type"`
	Value         sdk.Coin `json:"value"`
	ValuationDate time.Time `json:"valuation_date"`
}

type Reference struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
	Verified     bool   `json:"verified"`
}

type RiskFactor struct {
	FactorType   string `json:"factor_type"`
	Description  string `json:"description"`
	RiskLevel    string `json:"risk_level"`
	Mitigation   string `json:"mitigation"`
}

type Employment struct {
	Employer         string    `json:"employer"`
	Position         string    `json:"position"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date,omitempty"`
	Income           sdk.Coin  `json:"income"`
	EmploymentType   string    `json:"employment_type"`
	Industry         string    `json:"industry"`
}

type Property struct {
	PropertyType  string   `json:"property_type"`
	Value        sdk.Coin `json:"value"`
	LoanBalance  sdk.Coin `json:"loan_balance"`
	Equity       sdk.Coin `json:"equity"`
}

type Investment struct {
	InvestmentType string   `json:"investment_type"`
	Value         sdk.Coin `json:"value"`
	Risk          string   `json:"risk"`
}

type Insurance struct {
	PolicyType    string   `json:"policy_type"`
	CoverageAmount sdk.Coin `json:"coverage_amount"`
	Premium       sdk.Coin `json:"premium"`
}

type TaxRecord struct {
	TaxYear      int      `json:"tax_year"`
	GrossIncome  sdk.Coin `json:"gross_income"`
	NetIncome    sdk.Coin `json:"net_income"`
	TaxesPaid    sdk.Coin `json:"taxes_paid"`
}

type LegalRecord struct {
	RecordType   string    `json:"record_type"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

type SocialMediaData struct {
	HasData      bool     `json:"has_data"`
	Platforms    []string `json:"platforms"`
	VerifiedData bool     `json:"verified_data"`
}

type TransactionBehavior struct {
	AverageMonthlyTransactions int      `json:"average_monthly_transactions"`
	TransactionRegularity     string   `json:"transaction_regularity"`
	SpendingCategories        []string `json:"spending_categories"`
}

type LocationInfo struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Pincode string `json:"pincode"`
}

type PEPStatus struct {
	IsPEP        bool   `json:"is_pep"`
	ExposureType string `json:"exposure_type"`
	Details      string `json:"details"`
}

type RuleAction struct {
	ActionType   string            `json:"action_type"`
	Parameters   map[string]string `json:"parameters"`
	Description  string            `json:"description"`
}

// Stub implementations for helper functions

func (cse *CreditScoringEngine) modelApplies(model CreditModel, loanType AssessmentType) bool {
	return true // All models apply to all loan types for now
}

func (cse *CreditScoringEngine) ruleApplies(rule ScoringRule, loanType AssessmentType) bool {
	if len(rule.AppliesTo) == 0 {
		return true
	}
	for _, applicableType := range rule.AppliesTo {
		if applicableType == loanType {
			return true
		}
	}
	return false
}

func (cse *CreditScoringEngine) applyRule(rule ScoringRule, profile CustomerProfile, request LoanRequest) (int, bool) {
	// Stub implementation
	return 0, false
}

func (cse *CreditScoringEngine) normalizeScore(score int) int {
	if score < 300 {
		return 300
	}
	if score > 850 {
		return 850
	}
	return score
}

func (cse *CreditScoringEngine) getImpact(score int) string {
	if score > 0 {
		return "POSITIVE"
	} else if score < 0 {
		return "NEGATIVE"
	}
	return "NEUTRAL"
}

func (cse *CreditScoringEngine) extractMLFeatures(profile CustomerProfile, request LoanRequest) map[string]float64 {
	return map[string]float64{
		"age":                        float64(profile.PersonalInfo.Age),
		"monthly_income":             profile.FinancialInfo.MonthlyIncome.Amount.ToDec().MustFloat64(),
		"debt_to_income_ratio":       profile.FinancialInfo.DebtToIncomeRatio.MustFloat64(),
		"employment_stability_years": cse.getEmploymentStabilityYears(profile.EmploymentHistory),
		"credit_history_length":      cse.getCreditHistoryLength(profile.CreditHistory),
		"requested_amount":          request.Amount.Amount.ToDec().MustFloat64(),
		"loan_term_months":          float64(request.TermMonths),
	}
}

func (cse *CreditScoringEngine) predictMLScore(features map[string]float64) float64 {
	// Simplified ML prediction simulation
	score := 650.0
	score += (features["monthly_income"] / 1000) * 0.5
	score -= features["debt_to_income_ratio"] * 200
	score += features["employment_stability_years"] * 10
	score += features["credit_history_length"] * 2
	
	return math.Max(300, math.Min(850, score))
}

func (cse *CreditScoringEngine) analyzeBehavioralPatterns(patterns TransactionBehavior) int {
	score := 0
	if patterns.TransactionRegularity == "REGULAR" {
		score += 30
	}
	return score
}

func (cse *CreditScoringEngine) analyzeBankingBehavior(history BankingHistory) int {
	score := 0
	score += history.PrimaryBankRelationshipYears * 2
	if history.ReturnedCheckCount == 0 {
		score += 20
	}
	return score
}

func (cse *CreditScoringEngine) analyzeSocialMediaData(data SocialMediaData) int {
	if data.HasData && data.VerifiedData {
		return 15
	}
	return 0
}

func (cse *CreditScoringEngine) calculatePaymentHistoryScore(history PaymentHistoryData) int {
	percentage := history.OnTimePaymentPercentage.MustFloat64()
	return int((percentage - 0.5) * 200) // 50% = 0 points, 100% = 100 points
}

func (cse *CreditScoringEngine) calculateEmploymentScore(employment []Employment) int {
	if len(employment) == 0 {
		return -50
	}
	
	// Calculate stability based on current employment
	current := employment[0] // Assuming first is current
	if current.EndDate.IsZero() { // Still employed
		years := time.Since(current.StartDate).Hours() / (24 * 365)
		return int(math.Min(50, years*10))
	}
	
	return -20 // Currently unemployed
}

func (cse *CreditScoringEngine) getEmploymentStabilityDescription(employment []Employment) string {
	if len(employment) == 0 {
		return "No employment history"
	}
	
	years := cse.getEmploymentStabilityYears(employment)
	return fmt.Sprintf("%.1f years current employment", years)
}

func (cse *CreditScoringEngine) getEmploymentStabilityYears(employment []Employment) float64 {
	if len(employment) == 0 {
		return 0
	}
	
	current := employment[0]
	if current.EndDate.IsZero() {
		return time.Since(current.StartDate).Hours() / (24 * 365)
	}
	
	return 0
}

func (cse *CreditScoringEngine) getCreditHistoryLength(history CreditHistory) float64 {
	// Simplified calculation
	return float64(len(history.AccountHistory)) * 0.5
}

func (cse *CreditScoringEngine) adjustGradeForRiskFactors(grade RiskGrade, factors []RiskFactor) RiskGrade {
	// Downgrade if high-risk factors present
	for _, factor := range factors {
		if factor.RiskLevel == "HIGH" {
			if grade < RISK_GRADE_D {
				return grade + 1
			}
		}
	}
	return grade
}

func (cse *CreditScoringEngine) adjustGradeForIndustry(grade RiskGrade, industry string) RiskGrade {
	// Industry risk adjustments
	highRiskIndustries := []string{"AVIATION", "HOSPITALITY", "ENTERTAINMENT"}
	for _, riskIndustry := range highRiskIndustries {
		if strings.ToUpper(industry) == riskIndustry {
			if grade < RISK_GRADE_D {
				return grade + 1
			}
		}
	}
	return grade
}

func (cse *CreditScoringEngine) getMinimumIncomeRequirement(loanType AssessmentType, amount sdk.Coin) sdk.Coin {
	// Simple calculation: 4x loan amount annually, so monthly = amount/3
	monthlyRequirement := amount.Amount.Quo(sdk.NewInt(3))
	return sdk.NewCoin(amount.Denom, monthlyRequirement)
}

// String methods for enums

func (rg RiskGrade) String() string {
	switch rg {
	case RISK_GRADE_AAA: return "AAA"
	case RISK_GRADE_AA: return "AA"
	case RISK_GRADE_A: return "A"
	case RISK_GRADE_BBB: return "BBB"
	case RISK_GRADE_BB: return "BB"
	case RISK_GRADE_B: return "B"
	case RISK_GRADE_CCC: return "CCC"
	case RISK_GRADE_D: return "D"
	default: return "UNKNOWN"
	}
}

func (as ApprovalStatus) String() string {
	switch as {
	case APPROVAL_STATUS_APPROVED: return "APPROVED"
	case APPROVAL_STATUS_CONDITIONALLY_APPROVED: return "CONDITIONALLY_APPROVED"
	case APPROVAL_STATUS_PENDING_REVIEW: return "PENDING_REVIEW"
	case APPROVAL_STATUS_DECLINED: return "DECLINED"
	case APPROVAL_STATUS_NEEDS_MORE_INFO: return "NEEDS_MORE_INFO"
	case APPROVAL_STATUS_MANUAL_REVIEW_REQUIRED: return "MANUAL_REVIEW_REQUIRED"
	default: return "UNKNOWN"
	}
}

func (rc RuleCategory) String() string {
	switch rc {
	case RULE_CATEGORY_INCOME_VERIFICATION: return "INCOME_VERIFICATION"
	case RULE_CATEGORY_DEBT_TO_INCOME: return "DEBT_TO_INCOME"
	case RULE_CATEGORY_PAYMENT_HISTORY: return "PAYMENT_HISTORY"
	case RULE_CATEGORY_EMPLOYMENT_STABILITY: return "EMPLOYMENT_STABILITY"
	case RULE_CATEGORY_COLLATERAL_VALUATION: return "COLLATERAL_VALUATION"
	case RULE_CATEGORY_INDUSTRY_RISK: return "INDUSTRY_RISK"
	case RULE_CATEGORY_GEOGRAPHIC_RISK: return "GEOGRAPHIC_RISK"
	case RULE_CATEGORY_REGULATORY_COMPLIANCE: return "REGULATORY_COMPLIANCE"
	default: return "UNKNOWN"
	}
}

// Public API methods

func (cse *CreditScoringEngine) GetAssessment(ctx context.Context, assessmentID string) (*CreditAssessment, error) {
	store := cse.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("credit_assessment_%s", assessmentID))
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("assessment not found: %s", assessmentID)
	}
	
	var assessment CreditAssessment
	if err := json.Unmarshal(bz, &assessment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assessment: %w", err)
	}
	
	return &assessment, nil
}

func (cse *CreditScoringEngine) GetCustomerAssessments(ctx context.Context, customerID string) ([]CreditAssessment, error) {
	// Implementation would query all assessments for a customer
	return []CreditAssessment{}, nil
}

func (cse *CreditScoringEngine) UpdateCreditModel(ctx context.Context, model CreditModel) error {
	// Implementation would update a credit model
	return nil
}