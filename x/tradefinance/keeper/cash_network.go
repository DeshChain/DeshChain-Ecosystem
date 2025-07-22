package keeper

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CashNetworkManager manages cash-in/cash-out operations
type CashNetworkManager struct {
	keeper               Keeper
	agentNetwork         *CashAgentNetwork
	atmNetwork           *ATMNetworkManager
	bankingPartners      *BankingPartnerManager
	retailNetwork        *RetailNetworkManager
	settlementEngine     *CashSettlementEngine
	complianceManager    *CashComplianceManager
	mu                   sync.RWMutex
}

// CashAgentNetwork manages cash agents
type CashAgentNetwork struct {
	agents               map[string]*CashAgent
	locationIndex        *LocationIndex
	capacityManager      *AgentCapacityManager
	commissionCalculator *CommissionCalculator
	performanceTracker   *AgentPerformanceTracker
	trainingSystem       *AgentTrainingSystem
}

// CashAgent represents a cash-in/cash-out agent
type CashAgent struct {
	AgentID             string
	AgentType           AgentType
	BusinessName        string
	Location            Location
	ServiceHours        []ServiceHour
	Capabilities        AgentCapabilities
	LiquidityStatus     LiquidityStatus
	CurrentCashBalance  map[string]sdk.Coin
	DailyLimit          map[string]sdk.Coin
	MonthlyLimit        map[string]sdk.Coin
	UtilizedToday       map[string]sdk.Coin
	UtilizedThisMonth   map[string]sdk.Coin
	CommissionRate      CommissionStructure
	Rating              float64
	Status              AgentStatus
	ComplianceStatus    ComplianceStatus
	LastAudit           time.Time
}

// ATMNetworkManager manages ATM integrations
type ATMNetworkManager struct {
	atmProviders         map[string]*ATMProvider
	atmLocations         map[string]*ATMLocation
	transactionProcessor *ATMTransactionProcessor
	cardlessWithdrawal   *CardlessWithdrawalManager
	qrATMManager         *QRBasedATMManager
	maintenanceTracker   *ATMMaintenanceTracker
}

// BankingPartnerManager manages banking partnerships
type BankingPartnerManager struct {
	partners             map[string]*BankingPartner
	branchNetwork        map[string]*BankBranch
	accountManager       *PartnerAccountManager
	settlementScheduler  *BankSettlementScheduler
	reconciliationEngine *BankReconciliationEngine
}

// RetailNetworkManager manages retail partners
type RetailNetworkManager struct {
	retailers            map[string]*RetailPartner
	posIntegration       *POSIntegrationManager
	barcodeManager       *BarcodeTransactionManager
	inventoryTracker     *CashInventoryTracker
	incentiveProgram     *RetailerIncentiveProgram
}

// CashSettlementEngine handles settlements
type CashSettlementEngine struct {
	settlementQueue      *SettlementQueue
	batchProcessor       *BatchSettlementProcessor
	realTimeSettlement   *RealTimeSettler
	disputeManager       *DisputeResolutionManager
	auditTrail           *SettlementAuditTrail
}

// CashComplianceManager ensures regulatory compliance
type CashComplianceManager struct {
	kycVerifier          *CashKYCVerifier
	transactionMonitor   *CashTransactionMonitor
	reportGenerator      *ComplianceReportGenerator
	limitEnforcer        *TransactionLimitEnforcer
	suspiciousActivity   *SuspiciousActivityDetector
}

// Types and enums
type AgentType int
type AgentStatus int
type ComplianceStatus int
type TransactionType int
type SettlementStatus int

const (
	// Agent Types
	IndividualAgent AgentType = iota
	RetailStore
	BankBranch
	PostOffice
	PetrolPump
	MobileAgent
	
	// Agent Status
	AgentActive AgentStatus = iota
	AgentInactive
	AgentSuspended
	AgentUnderReview
	
	// Transaction Types
	CashIn TransactionType = iota
	CashOut
	BalanceInquiry
	MiniStatement
)

// Core cash network methods

// RegisterCashAgent registers a new cash agent
func (k Keeper) RegisterCashAgent(ctx context.Context, agentRequest CashAgentRequest) (*CashAgent, error) {
	cnm := k.getCashNetworkManager()
	
	// Validate agent request
	if err := cnm.validateAgentRequest(agentRequest); err != nil {
		return nil, fmt.Errorf("invalid agent request: %w", err)
	}
	
	// Perform KYC verification
	kycResult, err := cnm.complianceManager.kycVerifier.verifyAgent(agentRequest.KYCDocuments)
	if err != nil || !kycResult.Approved {
		return nil, fmt.Errorf("KYC verification failed: %w", err)
	}
	
	// Create cash agent
	agent := &CashAgent{
		AgentID:            generateID("AGENT"),
		AgentType:          agentRequest.Type,
		BusinessName:       agentRequest.BusinessName,
		Location:           agentRequest.Location,
		ServiceHours:       agentRequest.ServiceHours,
		Capabilities:       agentRequest.Capabilities,
		LiquidityStatus:    LiquidityNormal,
		CurrentCashBalance: make(map[string]sdk.Coin),
		DailyLimit:         agentRequest.DailyLimits,
		MonthlyLimit:       agentRequest.MonthlyLimits,
		UtilizedToday:      make(map[string]sdk.Coin),
		UtilizedThisMonth:  make(map[string]sdk.Coin),
		CommissionRate:     cnm.agentNetwork.commissionCalculator.calculateCommission(agentRequest),
		Rating:             5.0, // Start with perfect rating
		Status:             AgentActive,
		ComplianceStatus:   CompliantStatus,
		LastAudit:          time.Now(),
	}
	
	// Initialize cash balances
	for currency, amount := range agentRequest.InitialDeposit {
		agent.CurrentCashBalance[currency] = amount
	}
	
	// Add to network
	cnm.agentNetwork.agents[agent.AgentID] = agent
	
	// Update location index
	cnm.agentNetwork.locationIndex.addAgent(agent)
	
	// Provide training access
	if agentRequest.RequiresTraining {
		cnm.agentNetwork.trainingSystem.enrollAgent(agent.AgentID)
	}
	
	// Store agent
	if err := k.storeCashAgent(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to store agent: %w", err)
	}
	
	return agent, nil
}

// ProcessCashIn processes a cash-in transaction
func (k Keeper) ProcessCashIn(ctx context.Context, cashInRequest CashInRequest) (*CashTransactionResult, error) {
	cnm := k.getCashNetworkManager()
	
	// Find nearest available agent
	agent, err := cnm.findAvailableAgent(cashInRequest.Location, cashInRequest.Amount)
	if err != nil {
		return nil, fmt.Errorf("no available agent found: %w", err)
	}
	
	// Validate transaction limits
	if err := cnm.complianceManager.limitEnforcer.validateCashIn(agent, cashInRequest); err != nil {
		return nil, fmt.Errorf("limit validation failed: %w", err)
	}
	
	// Create transaction
	transaction := &CashTransaction{
		TransactionID:   generateID("CASHIN"),
		Type:            CashIn,
		AgentID:         agent.AgentID,
		CustomerID:      cashInRequest.CustomerID,
		Amount:          cashInRequest.Amount,
		Currency:        cashInRequest.Currency,
		InitiatedTime:   time.Now(),
		Status:          TransactionPending,
		VerificationCode: generateVerificationCode(),
	}
	
	// Generate QR code for verification
	qrCode, err := cnm.generateTransactionQR(transaction)
	if err != nil {
		return nil, fmt.Errorf("QR generation failed: %w", err)
	}
	
	// Lock agent capacity
	if err := cnm.agentNetwork.capacityManager.reserveCapacity(agent.AgentID, cashInRequest.Amount); err != nil {
		return nil, fmt.Errorf("capacity reservation failed: %w", err)
	}
	
	// Process biometric verification if required
	if agent.Capabilities.BiometricEnabled && cashInRequest.BiometricData != nil {
		verified, err := k.verifyBiometric(ctx, cashInRequest.CustomerID, cashInRequest.BiometricData)
		if err != nil || !verified {
			cnm.agentNetwork.capacityManager.releaseCapacity(agent.AgentID, cashInRequest.Amount)
			return nil, fmt.Errorf("biometric verification failed")
		}
	}
	
	// Update agent balance
	agent.CurrentCashBalance[cashInRequest.Currency] = agent.CurrentCashBalance[cashInRequest.Currency].Add(cashInRequest.Amount)
	agent.UtilizedToday[cashInRequest.Currency] = agent.UtilizedToday[cashInRequest.Currency].Add(cashInRequest.Amount)
	agent.UtilizedThisMonth[cashInRequest.Currency] = agent.UtilizedThisMonth[cashInRequest.Currency].Add(cashInRequest.Amount)
	
	// Credit customer account
	if err := k.creditCustomerAccount(ctx, cashInRequest.CustomerID, cashInRequest.Amount); err != nil {
		// Rollback agent balance
		agent.CurrentCashBalance[cashInRequest.Currency] = agent.CurrentCashBalance[cashInRequest.Currency].Sub(cashInRequest.Amount)
		cnm.agentNetwork.capacityManager.releaseCapacity(agent.AgentID, cashInRequest.Amount)
		return nil, fmt.Errorf("account credit failed: %w", err)
	}
	
	// Calculate commission
	commission := cnm.agentNetwork.commissionCalculator.calculateCashInCommission(agent, cashInRequest.Amount)
	
	// Complete transaction
	transaction.Status = TransactionCompleted
	transaction.CompletedTime = timePtr(time.Now())
	transaction.Commission = commission
	
	// Store transaction
	if err := k.storeCashTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to store transaction: %w", err)
	}
	
	// Update agent metrics
	cnm.agentNetwork.performanceTracker.recordTransaction(agent.AgentID, transaction)
	
	// Check for suspicious activity
	if suspicious := cnm.complianceManager.suspiciousActivity.analyze(transaction); suspicious {
		cnm.complianceManager.reportGenerator.generateSAR(transaction)
	}
	
	result := &CashTransactionResult{
		TransactionID:    transaction.TransactionID,
		Status:           transaction.Status,
		Agent:            agent,
		Amount:           cashInRequest.Amount,
		Commission:       commission,
		VerificationCode: transaction.VerificationCode,
		QRCode:           qrCode,
		CompletionTime:   transaction.CompletedTime,
		Receipt:          cnm.generateReceipt(transaction),
	}
	
	// Send confirmation
	go cnm.sendConfirmation(cashInRequest.CustomerID, result)
	
	return result, nil
}

// ProcessCashOut processes a cash-out transaction
func (k Keeper) ProcessCashOut(ctx context.Context, cashOutRequest CashOutRequest) (*CashTransactionResult, error) {
	cnm := k.getCashNetworkManager()
	
	// Validate customer balance
	balance, err := k.getCustomerBalance(ctx, cashOutRequest.CustomerID, cashOutRequest.Currency)
	if err != nil || balance.IsLT(cashOutRequest.Amount) {
		return nil, fmt.Errorf("insufficient balance")
	}
	
	// Find agent with sufficient cash
	agent, err := cnm.findAgentWithCash(cashOutRequest.Location, cashOutRequest.Amount, cashOutRequest.Currency)
	if err != nil {
		return nil, fmt.Errorf("no agent with sufficient cash found: %w", err)
	}
	
	// Validate transaction limits
	if err := cnm.complianceManager.limitEnforcer.validateCashOut(agent, cashOutRequest); err != nil {
		return nil, fmt.Errorf("limit validation failed: %w", err)
	}
	
	// Create transaction
	transaction := &CashTransaction{
		TransactionID:    generateID("CASHOUT"),
		Type:             CashOut,
		AgentID:          agent.AgentID,
		CustomerID:       cashOutRequest.CustomerID,
		Amount:           cashOutRequest.Amount,
		Currency:         cashOutRequest.Currency,
		InitiatedTime:    time.Now(),
		Status:           TransactionPending,
		VerificationCode: generateVerificationCode(),
		OTP:              generateOTP(),
	}
	
	// Send OTP to customer
	if err := cnm.sendOTP(cashOutRequest.CustomerID, transaction.OTP); err != nil {
		return nil, fmt.Errorf("OTP send failed: %w", err)
	}
	
	// Wait for OTP verification (with timeout)
	otpVerified := make(chan bool)
	go cnm.waitForOTPVerification(transaction.TransactionID, transaction.OTP, otpVerified)
	
	select {
	case verified := <-otpVerified:
		if !verified {
			transaction.Status = TransactionFailed
			transaction.FailureReason = "OTP verification failed"
			k.storeCashTransaction(ctx, transaction)
			return nil, fmt.Errorf("OTP verification failed")
		}
	case <-time.After(5 * time.Minute):
		transaction.Status = TransactionExpired
		k.storeCashTransaction(ctx, transaction)
		return nil, fmt.Errorf("transaction expired")
	}
	
	// Debit customer account
	if err := k.debitCustomerAccount(ctx, cashOutRequest.CustomerID, cashOutRequest.Amount); err != nil {
		return nil, fmt.Errorf("account debit failed: %w", err)
	}
	
	// Update agent balance
	agent.CurrentCashBalance[cashOutRequest.Currency] = agent.CurrentCashBalance[cashOutRequest.Currency].Sub(cashOutRequest.Amount)
	agent.UtilizedToday[cashOutRequest.Currency] = agent.UtilizedToday[cashOutRequest.Currency].Add(cashOutRequest.Amount)
	agent.UtilizedThisMonth[cashOutRequest.Currency] = agent.UtilizedThisMonth[cashOutRequest.Currency].Add(cashOutRequest.Amount)
	
	// Calculate commission
	commission := cnm.agentNetwork.commissionCalculator.calculateCashOutCommission(agent, cashOutRequest.Amount)
	
	// Complete transaction
	transaction.Status = TransactionCompleted
	transaction.CompletedTime = timePtr(time.Now())
	transaction.Commission = commission
	
	// Store transaction
	if err := k.storeCashTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to store transaction: %w", err)
	}
	
	// Schedule settlement
	cnm.settlementEngine.scheduleSettlement(agent.AgentID, transaction)
	
	result := &CashTransactionResult{
		TransactionID:    transaction.TransactionID,
		Status:           transaction.Status,
		Agent:            agent,
		Amount:           cashOutRequest.Amount,
		Commission:       commission,
		VerificationCode: transaction.VerificationCode,
		CompletionTime:   transaction.CompletedTime,
		Receipt:          cnm.generateReceipt(transaction),
	}
	
	return result, nil
}

// ATM Network methods

func (anm *ATMNetworkManager) processCardlessWithdrawal(ctx context.Context, request CardlessWithdrawalRequest) (*ATMTransactionResult, error) {
	// Find nearest ATM
	atm, err := anm.findNearestATM(request.Location, request.Features)
	if err != nil {
		return nil, fmt.Errorf("no suitable ATM found: %w", err)
	}
	
	// Generate withdrawal code
	withdrawalCode := anm.cardlessWithdrawal.generateCode(request)
	
	// Create ATM transaction
	atmTx := &ATMTransaction{
		TransactionID:   generateID("ATM"),
		ATMID:           atm.ATMID,
		CustomerID:      request.CustomerID,
		Amount:          request.Amount,
		Type:            CardlessWithdrawalType,
		WithdrawalCode:  withdrawalCode,
		ExpiryTime:      time.Now().Add(30 * time.Minute),
		Status:          ATMTransactionPending,
	}
	
	// Reserve amount at ATM
	if err := anm.reserveATMCash(atm, request.Amount); err != nil {
		return nil, fmt.Errorf("ATM cash reservation failed: %w", err)
	}
	
	// Send withdrawal code to customer
	if err := anm.sendWithdrawalCode(request.CustomerID, withdrawalCode, atm); err != nil {
		anm.releaseATMCash(atm, request.Amount)
		return nil, fmt.Errorf("code delivery failed: %w", err)
	}
	
	result := &ATMTransactionResult{
		TransactionID:   atmTx.TransactionID,
		ATM:             atm,
		WithdrawalCode:  withdrawalCode,
		ExpiryTime:      atmTx.ExpiryTime,
		Status:          atmTx.Status,
	}
	
	return result, nil
}

// Settlement methods

func (cse *CashSettlementEngine) processAgentSettlement(ctx context.Context, agentID string) (*SettlementResult, error) {
	// Get pending transactions
	pendingTxns := cse.getPendingTransactions(agentID)
	if len(pendingTxns) == 0 {
		return nil, fmt.Errorf("no pending transactions")
	}
	
	// Calculate net position
	netPosition := cse.calculateNetPosition(pendingTxns)
	
	// Create settlement batch
	batch := &SettlementBatch{
		BatchID:       generateID("SETTLE"),
		AgentID:       agentID,
		Transactions:  pendingTxns,
		NetAmount:     netPosition,
		CreatedAt:     time.Now(),
		Status:        SettlementPending,
	}
	
	// Process settlement based on amount
	if netPosition.Amount.Abs().LT(sdk.NewInt(10000)) {
		// Real-time settlement for small amounts
		result, err := cse.realTimeSettlement.process(ctx, batch)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		// Batch settlement for larger amounts
		cse.settlementQueue.enqueue(batch)
		return &SettlementResult{
			BatchID:        batch.BatchID,
			Status:         SettlementQueued,
			EstimatedTime:  cse.getNextSettlementTime(),
		}, nil
	}
}

// Helper types

type CashAgentRequest struct {
	Type             AgentType
	BusinessName     string
	Location         Location
	ServiceHours     []ServiceHour
	Capabilities     AgentCapabilities
	KYCDocuments     []Document
	InitialDeposit   map[string]sdk.Coin
	DailyLimits      map[string]sdk.Coin
	MonthlyLimits    map[string]sdk.Coin
	RequiresTraining bool
}

type CashInRequest struct {
	CustomerID    string
	Amount        sdk.Coin
	Currency      string
	Location      Location
	BiometricData []byte
	Purpose       string
}

type CashOutRequest struct {
	CustomerID string
	Amount     sdk.Coin
	Currency   string
	Location   Location
	Purpose    string
}

type CashTransaction struct {
	TransactionID    string
	Type             TransactionType
	AgentID          string
	CustomerID       string
	Amount           sdk.Coin
	Currency         string
	InitiatedTime    time.Time
	CompletedTime    *time.Time
	Status           TransactionStatus
	VerificationCode string
	OTP              string
	Commission       sdk.Coin
	FailureReason    string
}

type CashTransactionResult struct {
	TransactionID    string
	Status           TransactionStatus
	Agent            *CashAgent
	Amount           sdk.Coin
	Commission       sdk.Coin
	VerificationCode string
	QRCode           string
	CompletionTime   *time.Time
	Receipt          *TransactionReceipt
}

type Location struct {
	Latitude     float64
	Longitude    float64
	Address      string
	City         string
	State        string
	Country      string
	PostalCode   string
	Landmark     string
}

type ServiceHour struct {
	DayOfWeek  int
	OpenTime   string
	CloseTime  string
	IsHoliday  bool
}

type AgentCapabilities struct {
	MaxCashInAmount      sdk.Coin
	MaxCashOutAmount     sdk.Coin
	SupportedCurrencies  []string
	BiometricEnabled     bool
	QRCodeEnabled        bool
	ReceiptPrinter       bool
	InternetConnectivity bool
	BackupPower          bool
}

type LiquidityStatus int

const (
	LiquidityNormal LiquidityStatus = iota
	LiquidityLow
	LiquidityCritical
	LiquidityExcess
)

type CommissionStructure struct {
	CashInRate       sdk.Dec
	CashOutRate      sdk.Dec
	MinCommission    sdk.Coin
	MaxCommission    sdk.Coin
	VolumeIncentives []VolumeIncentive
}

type VolumeIncentive struct {
	MinVolume        sdk.Coin
	MaxVolume        sdk.Coin
	BonusPercentage  sdk.Dec
}

type ATMLocation struct {
	ATMID            string
	Provider         string
	Location         Location
	Features         []string
	CashAvailable    map[string]sdk.Coin
	Status           ATMStatus
	LastRefill       time.Time
	MaintenanceDue   time.Time
}

type CardlessWithdrawalRequest struct {
	CustomerID string
	Amount     sdk.Coin
	Location   Location
	Features   []string
}

type ATMTransaction struct {
	TransactionID   string
	ATMID           string
	CustomerID      string
	Amount          sdk.Coin
	Type            ATMTransactionType
	WithdrawalCode  string
	ExpiryTime      time.Time
	Status          ATMTransactionStatus
}

type SettlementBatch struct {
	BatchID      string
	AgentID      string
	Transactions []*CashTransaction
	NetAmount    sdk.Coin
	CreatedAt    time.Time
	SettledAt    *time.Time
	Status       SettlementStatus
}

type SettlementResult struct {
	BatchID         string
	Status          SettlementStatus
	SettledAmount   sdk.Coin
	TransactionRef  string
	EstimatedTime   time.Time
	ActualTime      *time.Time
}

type ATMStatus int
type ATMTransactionType int
type ATMTransactionStatus int

const (
	ATMOperational ATMStatus = iota
	ATMOffline
	ATMMaintenance
	ATMOutOfCash
	
	CardlessWithdrawalType ATMTransactionType = iota
	QRWithdrawalType
	BiometricWithdrawalType
	
	ATMTransactionPending ATMTransactionStatus = iota
	ATMTransactionCompleted
	ATMTransactionExpired
	ATMTransactionFailed
)

// Utility functions

func (cnm *CashNetworkManager) findAvailableAgent(location Location, amount sdk.Coin) (*CashAgent, error) {
	// Find agents within radius
	nearbyAgents := cnm.agentNetwork.locationIndex.findAgentsWithinRadius(location, 5.0) // 5km radius
	
	var bestAgent *CashAgent
	minDistance := math.MaxFloat64
	
	for _, agent := range nearbyAgents {
		// Check if agent is active and has capacity
		if agent.Status != AgentActive {
			continue
		}
		
		// Check daily limit
		utilized := agent.UtilizedToday[amount.Denom]
		limit := agent.DailyLimit[amount.Denom]
		if utilized.Add(amount).GT(limit) {
			continue
		}
		
		// Calculate distance
		distance := calculateDistance(location, agent.Location)
		if distance < minDistance {
			minDistance = distance
			bestAgent = agent
		}
	}
	
	if bestAgent == nil {
		return nil, fmt.Errorf("no available agent found")
	}
	
	return bestAgent, nil
}

func calculateDistance(loc1, loc2 Location) float64 {
	// Haversine formula for distance calculation
	const earthRadius = 6371.0 // km
	
	lat1Rad := loc1.Latitude * math.Pi / 180
	lat2Rad := loc2.Latitude * math.Pi / 180
	deltaLat := (loc2.Latitude - loc1.Latitude) * math.Pi / 180
	deltaLon := (loc2.Longitude - loc1.Longitude) * math.Pi / 180
	
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

func generateVerificationCode() string {
	return generateID("VER")[:8]
}

func generateOTP() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}