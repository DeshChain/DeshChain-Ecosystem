package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// FraudDetectionEngine implements ML-based fraud detection for trade finance
type FraudDetectionEngine struct {
	keeper *Keeper
}

// NewFraudDetectionEngine creates a new fraud detection engine
func NewFraudDetectionEngine(k *Keeper) *FraudDetectionEngine {
	return &FraudDetectionEngine{keeper: k}
}

// FraudDetectionResult represents the result of fraud detection analysis
type FraudDetectionResult struct {
	TransactionID       string                 `json:"transaction_id"`
	OverallRiskScore    float64                `json:"overall_risk_score"`
	RiskLevel          FraudRiskLevel         `json:"risk_level"`
	FraudProbability   float64                `json:"fraud_probability"`
	RiskFactors        []RiskFactor           `json:"risk_factors"`
	AnomalyScores      map[string]float64     `json:"anomaly_scores"`
	ModelPredictions   []ModelPrediction      `json:"model_predictions"`
	RecommendedAction  FraudAction            `json:"recommended_action"`
	RequiredApprovals  []string               `json:"required_approvals"`
	MonitoringFlags    []MonitoringFlag       `json:"monitoring_flags"`
	ComplianceFlags    []ComplianceAlert      `json:"compliance_flags"`
	Timestamp          time.Time              `json:"timestamp"`
	ProcessingTimeMs   int64                  `json:"processing_time_ms"`
}

// FraudAnalysisInput contains all data needed for fraud analysis
type FraudAnalysisInput struct {
	TransactionData    TransactionData        `json:"transaction_data"`
	CustomerProfile    CustomerProfile        `json:"customer_profile"`
	HistoricalData     HistoricalData         `json:"historical_data"`
	ExternalData       ExternalDataSources    `json:"external_data"`
	ContextualData     ContextualInformation  `json:"contextual_data"`
}

type TransactionData struct {
	TransactionID      string                 `json:"transaction_id"`
	Type              string                 `json:"type"` // LC, SBLC, Guarantee, etc.
	Amount            sdk.Coin               `json:"amount"`
	Currency          string                 `json:"currency"`
	BeneficiaryData   PartyData              `json:"beneficiary_data"`
	ApplicantData     PartyData              `json:"applicant_data"`
	DocumentData      []DocumentInfo         `json:"document_data"`
	PaymentTerms      PaymentTerms           `json:"payment_terms"`
	TradeGoods        []TradeGood            `json:"trade_goods"`
	ShippingData      ShippingInformation    `json:"shipping_data"`
	Timestamp         time.Time              `json:"timestamp"`
}

type CustomerProfile struct {
	CustomerID        string                 `json:"customer_id"`
	KYCLevel          int                    `json:"kyc_level"`
	CustomerType      string                 `json:"customer_type"` // Individual, Corporate, SME
	Industry          string                 `json:"industry"`
	CountryOfOrigin   string                 `json:"country_of_origin"`
	BusinessAge       int                    `json:"business_age"` // Years
	AnnualTurnover    sdk.Coin               `json:"annual_turnover"`
	CreditRating      string                 `json:"credit_rating"`
	RiskRating        string                 `json:"risk_rating"`
	RelationshipAge   int                    `json:"relationship_age"` // Months with bank
	PreviousDefaults  int                    `json:"previous_defaults"`
	ComplianceHistory []ComplianceIncident   `json:"compliance_history"`
}

type HistoricalData struct {
	PreviousTransactions []HistoricalTransaction `json:"previous_transactions"`
	BehaviorPatterns     CustomerBehavior        `json:"behavior_patterns"`
	SeasonalPatterns     []SeasonalPattern       `json:"seasonal_patterns"`
	PeerComparisons      PeerBenchmark           `json:"peer_comparisons"`
}

type ExternalDataSources struct {
	SanctionsLists     []SanctionsMatch       `json:"sanctions_lists"`
	CreditBureauData   CreditBureauInfo       `json:"credit_bureau_data"`
	NewsAndMedia       []NewsItem             `json:"news_and_media"`
	RegulatoryWarnings []RegulatoryWarning    `json:"regulatory_warnings"`
	IndustryData       IndustryMetrics        `json:"industry_data"`
}

type ContextualInformation struct {
	GeopoliticalRisk   float64               `json:"geopolitical_risk"`
	EconomicIndicators EconomicContext       `json:"economic_indicators"`
	MarketConditions   MarketContext         `json:"market_conditions"`
	SeasonalFactors    SeasonalContext       `json:"seasonal_factors"`
	RegulatoryChanges  []RegulatoryChange    `json:"regulatory_changes"`
}

// Supporting data structures

type PartyData struct {
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	Country         string    `json:"country"`
	TaxID           string    `json:"tax_id"`
	BankDetails     string    `json:"bank_details"`
	ContactInfo     string    `json:"contact_info"`
	BusinessType    string    `json:"business_type"`
	RegistrationDate time.Time `json:"registration_date"`
}

type DocumentInfo struct {
	DocumentType    string    `json:"document_type"`
	DocumentID      string    `json:"document_id"`
	IssueDate       time.Time `json:"issue_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	Issuer          string    `json:"issuer"`
	Authenticity    float64   `json:"authenticity"` // ML-based authenticity score
	Anomalies       []string  `json:"anomalies"`
}

type PaymentTerms struct {
	PaymentMethod   string    `json:"payment_method"`
	CreditPeriod    int       `json:"credit_period"` // Days
	InterestRate    sdk.Dec   `json:"interest_rate"`
	Guarantees      []string  `json:"guarantees"`
	Collateral      string    `json:"collateral"`
}

type TradeGood struct {
	Description     string   `json:"description"`
	HSCode          string   `json:"hs_code"`
	Quantity        sdk.Dec  `json:"quantity"`
	UnitPrice       sdk.Coin `json:"unit_price"`
	TotalValue      sdk.Coin `json:"total_value"`
	Origin          string   `json:"origin"`
	Destination     string   `json:"destination"`
	RiskCategory    string   `json:"risk_category"`
}

type ShippingInformation struct {
	ShippingLine    string    `json:"shipping_line"`
	PortOfLoading   string    `json:"port_of_loading"`
	PortOfDischarge string    `json:"port_of_discharge"`
	ShipmentDate    time.Time `json:"shipment_date"`
	ArrivalDate     time.Time `json:"arrival_date"`
	TrackingNumber  string    `json:"tracking_number"`
	RouteRisk       float64   `json:"route_risk"`
}

type HistoricalTransaction struct {
	TransactionID   string    `json:"transaction_id"`
	Amount          sdk.Coin  `json:"amount"`
	Type            string    `json:"type"`
	Counterparty    string    `json:"counterparty"`
	Date            time.Time `json:"date"`
	Status          string    `json:"status"`
	FraudScore      float64   `json:"fraud_score"`
}

type CustomerBehavior struct {
	AverageTransactionAmount sdk.Coin  `json:"average_transaction_amount"`
	TransactionFrequency     float64   `json:"transaction_frequency"` // Per month
	PreferredCounterparties  []string  `json:"preferred_counterparties"`
	SeasonalityScore         float64   `json:"seasonality_score"`
	BehaviorStability        float64   `json:"behavior_stability"`
	RiskAppetite            float64   `json:"risk_appetite"`
}

type SeasonalPattern struct {
	Period          string   `json:"period"` // Monthly, Quarterly, etc.
	AverageAmount   sdk.Coin `json:"average_amount"`
	TransactionCount int     `json:"transaction_count"`
	Volatility      float64  `json:"volatility"`
}

type PeerBenchmark struct {
	IndustryAverage    sdk.Coin `json:"industry_average"`
	SizeClassAverage   sdk.Coin `json:"size_class_average"`
	RegionAverage      sdk.Coin `json:"region_average"`
	PercentileRank     float64  `json:"percentile_rank"`
}

// Risk assessment models

type RiskFactor struct {
	Factor          string  `json:"factor"`
	Score           float64 `json:"score"`
	Weight          float64 `json:"weight"`
	ContributionPct float64 `json:"contribution_pct"`
	Description     string  `json:"description"`
	Severity        string  `json:"severity"`
}

type ModelPrediction struct {
	ModelName       string  `json:"model_name"`
	ModelVersion    string  `json:"model_version"`
	FraudProbability float64 `json:"fraud_probability"`
	Confidence      float64 `json:"confidence"`
	ModelType       string  `json:"model_type"` // RandomForest, XGBoost, Neural Network
	Features        []FeatureImportance `json:"features"`
}

type FeatureImportance struct {
	FeatureName string  `json:"feature_name"`
	Importance  float64 `json:"importance"`
	Value       float64 `json:"value"`
}

type MonitoringFlag struct {
	FlagType    string `json:"flag_type"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Action      string `json:"action"`
}

type ComplianceAlert struct {
	AlertType   string    `json:"alert_type"`
	Message     string    `json:"message"`
	RiskLevel   string    `json:"risk_level"`
	RequiredBy  time.Time `json:"required_by"`
	Authority   string    `json:"authority"`
}

// Enums
type FraudRiskLevel int

const (
	RISK_LEVEL_VERY_LOW FraudRiskLevel = iota
	RISK_LEVEL_LOW
	RISK_LEVEL_MEDIUM
	RISK_LEVEL_HIGH
	RISK_LEVEL_VERY_HIGH
	RISK_LEVEL_CRITICAL
)

type FraudAction int

const (
	ACTION_APPROVE FraudAction = iota
	ACTION_APPROVE_WITH_MONITORING
	ACTION_MANUAL_REVIEW
	ACTION_ENHANCED_DUE_DILIGENCE
	ACTION_REJECT
	ACTION_ESCALATE
)

// Core fraud detection methods

// AnalyzeTransaction performs comprehensive fraud analysis on a transaction
func (fde *FraudDetectionEngine) AnalyzeTransaction(ctx context.Context, input FraudAnalysisInput) (*FraudDetectionResult, error) {
	startTime := time.Now()

	result := &FraudDetectionResult{
		TransactionID:    input.TransactionData.TransactionID,
		RiskFactors:      []RiskFactor{},
		AnomalyScores:    make(map[string]float64),
		ModelPredictions: []ModelPrediction{},
		MonitoringFlags:  []MonitoringFlag{},
		ComplianceFlags:  []ComplianceAlert{},
		Timestamp:        time.Now(),
	}

	// 1. Basic validation and data preprocessing
	if err := fde.validateInput(input); err != nil {
		return result, fmt.Errorf("input validation failed: %w", err)
	}

	// 2. Rule-based risk assessment
	ruleBasedScore, riskFactors := fde.performRuleBasedAnalysis(ctx, input)
	result.RiskFactors = append(result.RiskFactors, riskFactors...)

	// 3. Anomaly detection
	anomalyScores := fde.detectAnomalies(ctx, input)
	result.AnomalyScores = anomalyScores

	// 4. ML model predictions
	modelPredictions, err := fde.runMLModels(ctx, input)
	if err != nil {
		fde.keeper.Logger(ctx).Error("ML model prediction failed", "error", err)
	} else {
		result.ModelPredictions = modelPredictions
	}

	// 5. Behavioral analysis
	behaviorScore := fde.analyzeBehavioralPatterns(ctx, input)

	// 6. Document authenticity analysis
	documentScore := fde.analyzeDocumentAuthenticity(ctx, input.TransactionData.DocumentData)

	// 7. Network analysis (counterparty risk)
	networkScore := fde.analyzeCounterpartyNetwork(ctx, input)

	// 8. Geopolitical risk assessment
	geopoliticalScore := fde.assessGeopoliticalRisk(ctx, input)

	// 9. Combine all scores using weighted ensemble
	weights := map[string]float64{
		"rule_based":      0.25,
		"anomaly":         0.15,
		"ml_models":       0.30,
		"behavioral":      0.15,
		"document":        0.10,
		"network":         0.03,
		"geopolitical":    0.02,
	}

	// Calculate weighted average of ML model predictions
	avgMLScore := 0.0
	if len(modelPredictions) > 0 {
		for _, pred := range modelPredictions {
			avgMLScore += pred.FraudProbability
		}
		avgMLScore /= float64(len(modelPredictions))
	}

	// Calculate overall risk score
	result.OverallRiskScore = weights["rule_based"]*ruleBasedScore +
		weights["anomaly"]*fde.aggregateAnomalyScores(anomalyScores) +
		weights["ml_models"]*avgMLScore +
		weights["behavioral"]*behaviorScore +
		weights["document"]*documentScore +
		weights["network"]*networkScore +
		weights["geopolitical"]*geopoliticalScore

	result.FraudProbability = result.OverallRiskScore

	// 10. Determine risk level and recommended action
	result.RiskLevel = fde.determineRiskLevel(result.OverallRiskScore)
	result.RecommendedAction = fde.determineRecommendedAction(result.RiskLevel, riskFactors)

	// 11. Generate monitoring flags
	result.MonitoringFlags = fde.generateMonitoringFlags(ctx, input, result)

	// 12. Generate compliance alerts
	result.ComplianceFlags = fde.generateComplianceAlerts(ctx, input, result)

	// 13. Determine required approvals
	result.RequiredApprovals = fde.determineRequiredApprovals(result.RiskLevel, result.RecommendedAction)

	// Record processing time
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	// Store analysis result for audit and learning
	if err := fde.storeAnalysisResult(ctx, *result); err != nil {
		fde.keeper.Logger(ctx).Error("Failed to store analysis result", "error", err)
	}

	// Emit fraud detection event
	fde.emitFraudDetectionEvent(ctx, result)

	return result, nil
}

// Rule-based analysis using predefined business rules
func (fde *FraudDetectionEngine) performRuleBasedAnalysis(ctx context.Context, input FraudAnalysisInput) (float64, []RiskFactor) {
	var riskFactors []RiskFactor
	totalScore := 0.0

	// Amount-based rules
	if input.TransactionData.Amount.Amount.GT(sdk.NewInt(1000000)) { // > $1M
		factor := RiskFactor{
			Factor:          "LARGE_AMOUNT",
			Score:           0.3,
			Weight:          1.0,
			ContributionPct: 30.0,
			Description:     "Transaction amount exceeds $1M threshold",
			Severity:        "MEDIUM",
		}
		riskFactors = append(riskFactors, factor)
		totalScore += factor.Score * factor.Weight
	}

	// Country risk rules
	highRiskCountries := []string{"AF", "KP", "IR", "SY", "MM", "BY"}
	for _, country := range highRiskCountries {
		if input.TransactionData.BeneficiaryData.Country == country {
			factor := RiskFactor{
				Factor:          "HIGH_RISK_COUNTRY",
				Score:           0.4,
				Weight:          1.0,
				ContributionPct: 40.0,
				Description:     fmt.Sprintf("Beneficiary in high-risk country: %s", country),
				Severity:        "HIGH",
			}
			riskFactors = append(riskFactors, factor)
			totalScore += factor.Score * factor.Weight
			break
		}
	}

	// New customer risk
	if input.CustomerProfile.RelationshipAge < 6 { // Less than 6 months
		factor := RiskFactor{
			Factor:          "NEW_CUSTOMER",
			Score:           0.2,
			Weight:          1.0,
			ContributionPct: 20.0,
			Description:     "Customer relationship less than 6 months old",
			Severity:        "MEDIUM",
		}
		riskFactors = append(riskFactors, factor)
		totalScore += factor.Score * factor.Weight
	}

	// Document inconsistency rules
	for _, doc := range input.TransactionData.DocumentData {
		if doc.Authenticity < 0.7 {
			factor := RiskFactor{
				Factor:          "DOCUMENT_AUTHENTICITY",
				Score:           0.35,
				Weight:          1.0,
				ContributionPct: 35.0,
				Description:     fmt.Sprintf("Low authenticity score for document: %s", doc.DocumentType),
				Severity:        "HIGH",
			}
			riskFactors = append(riskFactors, factor)
			totalScore += factor.Score * factor.Weight
		}
	}

	// Velocity rules - checking transaction frequency
	if len(input.HistoricalData.PreviousTransactions) > 0 {
		recentTxns := 0
		for _, tx := range input.HistoricalData.PreviousTransactions {
			if time.Since(tx.Date).Hours() < 24 { // Last 24 hours
				recentTxns++
			}
		}
		if recentTxns > 5 { // More than 5 transactions in 24 hours
			factor := RiskFactor{
				Factor:          "HIGH_VELOCITY",
				Score:           0.25,
				Weight:          1.0,
				ContributionPct: 25.0,
				Description:     fmt.Sprintf("High transaction velocity: %d transactions in 24h", recentTxns),
				Severity:        "MEDIUM",
			}
			riskFactors = append(riskFactors, factor)
			totalScore += factor.Score * factor.Weight
		}
	}

	// Sanctions screening results
	for _, match := range input.ExternalData.SanctionsLists {
		if match.MatchScore > 0.8 {
			factor := RiskFactor{
				Factor:          "SANCTIONS_MATCH",
				Score:           0.9,
				Weight:          1.0,
				ContributionPct: 90.0,
				Description:     fmt.Sprintf("High sanctions match: %s", match.ListName),
				Severity:        "CRITICAL",
			}
			riskFactors = append(riskFactors, factor)
			totalScore += factor.Score * factor.Weight
		}
	}

	// Normalize score to 0-1 range
	if totalScore > 1.0 {
		totalScore = 1.0
	}

	return totalScore, riskFactors
}

// Anomaly detection using statistical methods
func (fde *FraudDetectionEngine) detectAnomalies(ctx context.Context, input FraudAnalysisInput) map[string]float64 {
	anomalies := make(map[string]float64)

	// Amount anomaly detection
	if len(input.HistoricalData.PreviousTransactions) > 0 {
		amounts := make([]float64, len(input.HistoricalData.PreviousTransactions))
		for i, tx := range input.HistoricalData.PreviousTransactions {
			amounts[i] = float64(tx.Amount.Amount.Int64())
		}
		
		mean := calculateMean(amounts)
		stdDev := calculateStdDev(amounts, mean)
		currentAmount := float64(input.TransactionData.Amount.Amount.Int64())
		
		if stdDev > 0 {
			zScore := math.Abs((currentAmount - mean) / stdDev)
			anomalies["amount_zscore"] = math.Min(zScore/3.0, 1.0) // Normalize to 0-1
		}
	}

	// Time-based anomaly detection
	if input.CustomerProfile.RelationshipAge > 0 {
		// Check if transaction time is unusual for this customer
		currentHour := input.TransactionData.Timestamp.Hour()
		hourlyPattern := fde.calculateHourlyPattern(input.HistoricalData.PreviousTransactions)
		
		if normalizedHourlyFreq, exists := hourlyPattern[currentHour]; exists {
			anomalies["temporal"] = 1.0 - normalizedHourlyFreq // Lower frequency = higher anomaly
		}
	}

	// Counterparty anomaly detection
	isNewCounterparty := true
	for _, preferred := range input.HistoricalData.BehaviorPatterns.PreferredCounterparties {
		if preferred == input.TransactionData.BeneficiaryData.Name {
			isNewCounterparty = false
			break
		}
	}
	if isNewCounterparty {
		anomalies["new_counterparty"] = 0.3 // Moderate anomaly for new counterparty
	}

	// Geographic anomaly detection
	beneficiaryCountry := input.TransactionData.BeneficiaryData.Country
	countryFreq := fde.calculateCountryFrequency(input.HistoricalData.PreviousTransactions, beneficiaryCountry)
	if countryFreq == 0 {
		anomalies["geographic"] = 0.4 // New country is moderately anomalous
	} else if countryFreq < 0.1 {
		anomalies["geographic"] = 0.2 // Rarely used country
	}

	return anomalies
}

// ML model predictions (simplified implementation)
func (fde *FraudDetectionEngine) runMLModels(ctx context.Context, input FraudAnalysisInput) ([]ModelPrediction, error) {
	var predictions []ModelPrediction

	// Random Forest Model (simulated)
	rfPrediction := ModelPrediction{
		ModelName:        "RandomForest_v2.1",
		ModelVersion:     "2.1.0",
		FraudProbability: fde.simulateRandomForestPrediction(input),
		Confidence:       0.85,
		ModelType:        "RandomForest",
		Features:         fde.extractFeatureImportances("RandomForest", input),
	}
	predictions = append(predictions, rfPrediction)

	// XGBoost Model (simulated)
	xgbPrediction := ModelPrediction{
		ModelName:        "XGBoost_v1.6",
		ModelVersion:     "1.6.0",
		FraudProbability: fde.simulateXGBoostPrediction(input),
		Confidence:       0.82,
		ModelType:        "XGBoost",
		Features:         fde.extractFeatureImportances("XGBoost", input),
	}
	predictions = append(predictions, xgbPrediction)

	// Neural Network Model (simulated)
	nnPrediction := ModelPrediction{
		ModelName:        "NeuralNet_v3.0",
		ModelVersion:     "3.0.0",
		FraudProbability: fde.simulateNeuralNetPrediction(input),
		Confidence:       0.78,
		ModelType:        "NeuralNetwork",
		Features:         fde.extractFeatureImportances("NeuralNetwork", input),
	}
	predictions = append(predictions, nnPrediction)

	return predictions, nil
}

// Additional analysis methods

func (fde *FraudDetectionEngine) analyzeBehavioralPatterns(ctx context.Context, input FraudAnalysisInput) float64 {
	score := 0.0

	// Analyze deviation from typical behavior
	if input.HistoricalData.BehaviorPatterns.BehaviorStability < 0.5 {
		score += 0.3 // High variability in behavior
	}

	currentAmount := input.TransactionData.Amount
	avgAmount := input.HistoricalData.BehaviorPatterns.AverageTransactionAmount
	
	if avgAmount.IsPositive() {
		ratio := float64(currentAmount.Amount.Int64()) / float64(avgAmount.Amount.Int64())
		if ratio > 5.0 || ratio < 0.2 {
			score += 0.4 // Significantly different from typical amount
		}
	}

	return math.Min(score, 1.0)
}

func (fde *FraudDetectionEngine) analyzeDocumentAuthenticity(ctx context.Context, documents []DocumentInfo) float64 {
	if len(documents) == 0 {
		return 0.5 // Moderate risk for no documents
	}

	totalAuthenticity := 0.0
	for _, doc := range documents {
		totalAuthenticity += doc.Authenticity
	}
	avgAuthenticity := totalAuthenticity / float64(len(documents))
	
	return 1.0 - avgAuthenticity // Convert authenticity to risk score
}

func (fde *FraudDetectionEngine) analyzeCounterpartyNetwork(ctx context.Context, input FraudAnalysisInput) float64 {
	// Simplified network analysis - would be more complex in production
	beneficiaryName := input.TransactionData.BeneficiaryData.Name
	
	// Check if beneficiary appears in any sanctions or warning lists
	for _, sanctions := range input.ExternalData.SanctionsLists {
		if sanctions.EntityName == beneficiaryName && sanctions.MatchScore > 0.6 {
			return 0.8 // High network risk
		}
	}

	return 0.1 // Low baseline network risk
}

func (fde *FraudDetectionEngine) assessGeopoliticalRisk(ctx context.Context, input FraudAnalysisInput) float64 {
	return input.ContextualData.GeopoliticalRisk
}

// Helper methods for calculations

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(values)-1))
}

func (fde *FraudDetectionEngine) calculateHourlyPattern(transactions []HistoricalTransaction) map[int]float64 {
	hourCounts := make(map[int]int)
	totalTxns := len(transactions)
	
	for _, tx := range transactions {
		hour := tx.Date.Hour()
		hourCounts[hour]++
	}
	
	// Normalize to frequencies
	hourPattern := make(map[int]float64)
	for hour, count := range hourCounts {
		hourPattern[hour] = float64(count) / float64(totalTxns)
	}
	
	return hourPattern
}

func (fde *FraudDetectionEngine) calculateCountryFrequency(transactions []HistoricalTransaction, country string) float64 {
	countryCount := 0
	for _, tx := range transactions {
		// This would need to be enhanced to extract country from transaction data
		if tx.Counterparty == country { // Simplified assumption
			countryCount++
		}
	}
	
	if len(transactions) == 0 {
		return 0
	}
	
	return float64(countryCount) / float64(len(transactions))
}

func (fde *FraudDetectionEngine) aggregateAnomalyScores(anomalies map[string]float64) float64 {
	if len(anomalies) == 0 {
		return 0
	}
	
	total := 0.0
	for _, score := range anomalies {
		total += score
	}
	
	return math.Min(total/float64(len(anomalies)), 1.0)
}

func (fde *FraudDetectionEngine) determineRiskLevel(score float64) FraudRiskLevel {
	switch {
	case score >= 0.9:
		return RISK_LEVEL_CRITICAL
	case score >= 0.7:
		return RISK_LEVEL_VERY_HIGH
	case score >= 0.5:
		return RISK_LEVEL_HIGH
	case score >= 0.3:
		return RISK_LEVEL_MEDIUM
	case score >= 0.1:
		return RISK_LEVEL_LOW
	default:
		return RISK_LEVEL_VERY_LOW
	}
}

func (fde *FraudDetectionEngine) determineRecommendedAction(riskLevel FraudRiskLevel, riskFactors []RiskFactor) FraudAction {
	// Check for critical risk factors
	for _, factor := range riskFactors {
		if factor.Factor == "SANCTIONS_MATCH" || factor.Severity == "CRITICAL" {
			return ACTION_REJECT
		}
	}

	switch riskLevel {
	case RISK_LEVEL_CRITICAL:
		return ACTION_REJECT
	case RISK_LEVEL_VERY_HIGH:
		return ACTION_ESCALATE
	case RISK_LEVEL_HIGH:
		return ACTION_ENHANCED_DUE_DILIGENCE
	case RISK_LEVEL_MEDIUM:
		return ACTION_MANUAL_REVIEW
	case RISK_LEVEL_LOW:
		return ACTION_APPROVE_WITH_MONITORING
	default:
		return ACTION_APPROVE
	}
}

// Simulation methods for ML models (in production, these would call actual ML services)

func (fde *FraudDetectionEngine) simulateRandomForestPrediction(input FraudAnalysisInput) float64 {
	// Simplified simulation based on key features
	score := 0.0
	
	// Amount factor
	if input.TransactionData.Amount.Amount.GT(sdk.NewInt(500000)) {
		score += 0.2
	}
	
	// Country risk
	highRiskCountries := []string{"AF", "KP", "IR", "SY"}
	for _, country := range highRiskCountries {
		if input.TransactionData.BeneficiaryData.Country == country {
			score += 0.4
			break
		}
	}
	
	// Customer age factor
	if input.CustomerProfile.RelationshipAge < 12 {
		score += 0.1
	}
	
	return math.Min(score, 1.0)
}

func (fde *FraudDetectionEngine) simulateXGBoostPrediction(input FraudAnalysisInput) float64 {
	// Different weighting than RandomForest
	score := 0.0
	
	// Focus more on behavioral patterns
	if input.HistoricalData.BehaviorPatterns.BehaviorStability < 0.6 {
		score += 0.3
	}
	
	// Document authenticity focus
	avgAuth := 0.0
	if len(input.TransactionData.DocumentData) > 0 {
		for _, doc := range input.TransactionData.DocumentData {
			avgAuth += doc.Authenticity
		}
		avgAuth /= float64(len(input.TransactionData.DocumentData))
		if avgAuth < 0.7 {
			score += 0.4
		}
	}
	
	return math.Min(score, 1.0)
}

func (fde *FraudDetectionEngine) simulateNeuralNetPrediction(input FraudAnalysisInput) float64 {
	// Neural network might capture complex interactions
	score := 0.0
	
	// Complex interaction between amount and customer profile
	amountFactor := float64(input.TransactionData.Amount.Amount.Int64()) / 1000000.0 // Scale to millions
	ageFactor := float64(input.CustomerProfile.RelationshipAge) / 12.0 // Scale to years
	
	// Non-linear interaction
	interaction := amountFactor * (1.0 - ageFactor) // High amount + new customer = higher risk
	if interaction > 0.5 {
		score += 0.3
	}
	
	return math.Min(score, 1.0)
}

func (fde *FraudDetectionEngine) extractFeatureImportances(modelType string, input FraudAnalysisInput) []FeatureImportance {
	// Simulated feature importances - would come from actual models in production
	switch modelType {
	case "RandomForest":
		return []FeatureImportance{
			{FeatureName: "transaction_amount", Importance: 0.25, Value: float64(input.TransactionData.Amount.Amount.Int64())},
			{FeatureName: "beneficiary_country", Importance: 0.20, Value: 1.0}, // Categorical encoded
			{FeatureName: "customer_age", Importance: 0.15, Value: float64(input.CustomerProfile.RelationshipAge)},
			{FeatureName: "document_authenticity", Importance: 0.12, Value: 0.8}, // Average authenticity
		}
	case "XGBoost":
		return []FeatureImportance{
			{FeatureName: "behavior_stability", Importance: 0.30, Value: input.HistoricalData.BehaviorPatterns.BehaviorStability},
			{FeatureName: "transaction_amount", Importance: 0.22, Value: float64(input.TransactionData.Amount.Amount.Int64())},
			{FeatureName: "document_authenticity", Importance: 0.18, Value: 0.8},
		}
	case "NeuralNetwork":
		return []FeatureImportance{
			{FeatureName: "amount_customer_interaction", Importance: 0.35, Value: 0.6},
			{FeatureName: "geopolitical_risk", Importance: 0.25, Value: input.ContextualData.GeopoliticalRisk},
			{FeatureName: "network_risk", Importance: 0.20, Value: 0.1},
		}
	}
	return []FeatureImportance{}
}

// Monitoring and compliance methods

func (fde *FraudDetectionEngine) generateMonitoringFlags(ctx context.Context, input FraudAnalysisInput, result *FraudDetectionResult) []MonitoringFlag {
	var flags []MonitoringFlag

	if result.OverallRiskScore > 0.3 {
		flags = append(flags, MonitoringFlag{
			FlagType:    "ENHANCED_MONITORING",
			Description: "Transaction requires enhanced monitoring",
			Severity:    "MEDIUM",
			Action:      "Monitor subsequent transactions for 90 days",
		})
	}

	if result.OverallRiskScore > 0.7 {
		flags = append(flags, MonitoringFlag{
			FlagType:    "HIGH_RISK_CUSTOMER",
			Description: "Customer profile shows high risk indicators",
			Severity:    "HIGH",
			Action:      "Conduct periodic risk assessment every 30 days",
		})
	}

	return flags
}

func (fde *FraudDetectionEngine) generateComplianceAlerts(ctx context.Context, input FraudAnalysisInput, result *FraudDetectionResult) []ComplianceAlert {
	var alerts []ComplianceAlert

	// Large transaction reporting
	if input.TransactionData.Amount.Amount.GT(sdk.NewInt(10000)) {
		alerts = append(alerts, ComplianceAlert{
			AlertType:  "CTR_REPORTING",
			Message:    "Transaction exceeds CTR reporting threshold",
			RiskLevel:  "INFO",
			RequiredBy: time.Now().Add(24 * time.Hour),
			Authority:  "FinCEN",
		})
	}

	// Suspicious activity reporting
	if result.OverallRiskScore > 0.8 {
		alerts = append(alerts, ComplianceAlert{
			AlertType:  "SAR_FILING",
			Message:    "Suspicious activity detected - SAR filing may be required",
			RiskLevel:  "HIGH",
			RequiredBy: time.Now().Add(30 * 24 * time.Hour), // 30 days
			Authority:  "FinCEN",
		})
	}

	return alerts
}

func (fde *FraudDetectionEngine) determineRequiredApprovals(riskLevel FraudRiskLevel, action FraudAction) []string {
	var approvals []string

	switch riskLevel {
	case RISK_LEVEL_HIGH, RISK_LEVEL_VERY_HIGH:
		approvals = append(approvals, "senior_officer", "compliance_manager")
	case RISK_LEVEL_CRITICAL:
		approvals = append(approvals, "senior_officer", "compliance_manager", "chief_risk_officer")
	}

	if action == ACTION_ENHANCED_DUE_DILIGENCE || action == ACTION_ESCALATE {
		approvals = append(approvals, "aml_specialist")
	}

	return approvals
}

// Storage and event methods

func (fde *FraudDetectionEngine) validateInput(input FraudAnalysisInput) error {
	if input.TransactionData.TransactionID == "" {
		return fmt.Errorf("transaction ID is required")
	}
	if !input.TransactionData.Amount.IsValid() {
		return fmt.Errorf("invalid transaction amount")
	}
	if input.CustomerProfile.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	return nil
}

func (fde *FraudDetectionEngine) storeAnalysisResult(ctx context.Context, result FraudDetectionResult) error {
	store := fde.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("fraud_analysis_%s", result.TransactionID))
	bz := fde.keeper.cdc.MustMarshal(&result)
	store.Set(key, bz)
	return nil
}

func (fde *FraudDetectionEngine) emitFraudDetectionEvent(ctx context.Context, result *FraudDetectionResult) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fraud_detection_analysis",
			sdk.NewAttribute("transaction_id", result.TransactionID),
			sdk.NewAttribute("risk_score", fmt.Sprintf("%.3f", result.OverallRiskScore)),
			sdk.NewAttribute("risk_level", fde.riskLevelToString(result.RiskLevel)),
			sdk.NewAttribute("recommended_action", fde.actionToString(result.RecommendedAction)),
			sdk.NewAttribute("processing_time_ms", fmt.Sprintf("%d", result.ProcessingTimeMs)),
		),
	)
}

func (fde *FraudDetectionEngine) riskLevelToString(level FraudRiskLevel) string {
	switch level {
	case RISK_LEVEL_VERY_LOW:
		return "VERY_LOW"
	case RISK_LEVEL_LOW:
		return "LOW"
	case RISK_LEVEL_MEDIUM:
		return "MEDIUM"
	case RISK_LEVEL_HIGH:
		return "HIGH"
	case RISK_LEVEL_VERY_HIGH:
		return "VERY_HIGH"
	case RISK_LEVEL_CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func (fde *FraudDetectionEngine) actionToString(action FraudAction) string {
	switch action {
	case ACTION_APPROVE:
		return "APPROVE"
	case ACTION_APPROVE_WITH_MONITORING:
		return "APPROVE_WITH_MONITORING"
	case ACTION_MANUAL_REVIEW:
		return "MANUAL_REVIEW"
	case ACTION_ENHANCED_DUE_DILIGENCE:
		return "ENHANCED_DUE_DILIGENCE"
	case ACTION_REJECT:
		return "REJECT"
	case ACTION_ESCALATE:
		return "ESCALATE"
	default:
		return "UNKNOWN"
	}
}

// Supporting data structures for external data

type SanctionsMatch struct {
	EntityName  string  `json:"entity_name"`
	ListName    string  `json:"list_name"`
	MatchScore  float64 `json:"match_score"`
	MatchType   string  `json:"match_type"`
	ListType    string  `json:"list_type"`
}

type CreditBureauInfo struct {
	CreditScore     int     `json:"credit_score"`
	PaymentHistory  string  `json:"payment_history"`
	CreditUtilization float64 `json:"credit_utilization"`
	AccountAge      int     `json:"account_age"`
}

type NewsItem struct {
	Headline    string    `json:"headline"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Sentiment   float64   `json:"sentiment"`
	Relevance   float64   `json:"relevance"`
}

type RegulatoryWarning struct {
	Authority   string    `json:"authority"`
	WarningType string    `json:"warning_type"`
	Description string    `json:"description"`
	IssuedAt    time.Time `json:"issued_at"`
	Severity    string    `json:"severity"`
}

type IndustryMetrics struct {
	IndustryCode    string  `json:"industry_code"`
	AverageRisk     float64 `json:"average_risk"`
	FraudRate       float64 `json:"fraud_rate"`
	MarketCondition string  `json:"market_condition"`
}

type EconomicContext struct {
	GDPGrowthRate     float64 `json:"gdp_growth_rate"`
	InflationRate     float64 `json:"inflation_rate"`
	UnemploymentRate  float64 `json:"unemployment_rate"`
	CurrencyVolatility float64 `json:"currency_volatility"`
}

type MarketContext struct {
	MarketVolatility  float64 `json:"market_volatility"`
	TradingVolume     float64 `json:"trading_volume"`
	SectorPerformance float64 `json:"sector_performance"`
	LiquidityConditions string `json:"liquidity_conditions"`
}

type SeasonalContext struct {
	SeasonalFactor    float64 `json:"seasonal_factor"`
	HolidayEffect     float64 `json:"holiday_effect"`
	MonthlyPattern    float64 `json:"monthly_pattern"`
	WeeklyPattern     float64 `json:"weekly_pattern"`
}

type RegulatoryChange struct {
	ChangeType      string    `json:"change_type"`
	EffectiveDate   time.Time `json:"effective_date"`
	ImpactLevel     string    `json:"impact_level"`
	Description     string    `json:"description"`
	ComplianceReqd  bool      `json:"compliance_required"`
}

type ComplianceIncident struct {
	IncidentType    string    `json:"incident_type"`
	Date            time.Time `json:"date"`
	Severity        string    `json:"severity"`
	Resolution      string    `json:"resolution"`
	ActionTaken     string    `json:"action_taken"`
}