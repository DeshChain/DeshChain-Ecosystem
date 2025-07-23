package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/oracle/types"
)

// FallbackManager handles oracle failure detection and fallback mechanisms
type FallbackManager struct {
	keeper Keeper
}

// NewFallbackManager creates a new fallback manager
func NewFallbackManager(keeper Keeper) *FallbackManager {
	return &FallbackManager{
		keeper: keeper,
	}
}

// FallbackConfig represents fallback configuration for different feed types
type FallbackConfig struct {
	FeedType            string                    `json:"feed_type"`
	PrimaryOracles      []string                  `json:"primary_oracles"`
	SecondaryOracles    []string                  `json:"secondary_oracles"`
	EmergencyOracles    []string                  `json:"emergency_oracles"`
	FallbackThresholds  types.FallbackThresholds  `json:"fallback_thresholds"`
	DataStaleness       time.Duration             `json:"data_staleness"`
	MinimumOracles      int                       `json:"minimum_oracles"`
	ConsensusThreshold  sdk.Dec                   `json:"consensus_threshold"`
	EmergencyProtocol   types.EmergencyProtocol   `json:"emergency_protocol"`
	RecoveryMechanism   types.RecoveryMechanism   `json:"recovery_mechanism"`
	LastUpdate          time.Time                 `json:"last_update"`
}

// OracleHealthStatus represents the health status of oracle system
type OracleHealthStatus struct {
	OverallHealth       string                          `json:"overall_health"` // HEALTHY, DEGRADED, CRITICAL, EMERGENCY
	FeedHealthStatus    map[string]types.FeedHealth     `json:"feed_health_status"`
	ActiveOracles       int                             `json:"active_oracles"`
	FailedOracles       int                             `json:"failed_oracles"`
	DataFreshness       map[string]time.Duration        `json:"data_freshness"`
	ConsensusStatus     map[string]types.ConsensusInfo  `json:"consensus_status"`
	FallbacksActive     []string                        `json:"fallbacks_active"`
	EmergencyMode       bool                            `json:"emergency_mode"`
	LastHealthCheck     time.Time                       `json:"last_health_check"`
	SystemAlerts        []types.SystemAlert             `json:"system_alerts"`
}

// InitializeFallbackSystem initializes fallback configurations for all feed types
func (fm *FallbackManager) InitializeFallbackSystem(ctx sdk.Context) error {
	feedTypes := []string{"INR_USD", "BTC_INR", "ETH_INR", "NAMO_INR", "WEATHER", "MARKET_DATA"}
	
	for _, feedType := range feedTypes {
		config := fm.createDefaultFallbackConfig(ctx, feedType)
		fm.keeper.SetFallbackConfig(ctx, feedType, config)
	}

	// Initialize system health monitoring
	fm.initializeHealthMonitoring(ctx)

	// Emit initialization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFallbackSystemInitialized,
			sdk.NewAttribute(types.AttributeKeyFeedCount, fmt.Sprintf("%d", len(feedTypes))),
		),
	)

	return nil
}

// MonitorOracleHealth continuously monitors oracle system health
func (fm *FallbackManager) MonitorOracleHealth(ctx sdk.Context) (*OracleHealthStatus, error) {
	status := &OracleHealthStatus{
		FeedHealthStatus: make(map[string]types.FeedHealth),
		DataFreshness:    make(map[string]time.Duration),
		ConsensusStatus:  make(map[string]types.ConsensusInfo),
		FallbacksActive:  []string{},
		SystemAlerts:     []types.SystemAlert{},
		LastHealthCheck:  ctx.BlockTime(),
	}

	// Get all feed types
	feedTypes := fm.keeper.GetAllFeedTypes(ctx)
	
	// Check health of each feed
	criticalIssues := 0
	degradedFeeds := 0
	
	for _, feedType := range feedTypes {
		feedHealth := fm.checkFeedHealth(ctx, feedType)
		status.FeedHealthStatus[feedType] = feedHealth
		
		// Check data freshness
		latestData := fm.keeper.GetLatestOracleData(ctx, feedType)
		if latestData != nil {
			freshness := ctx.BlockTime().Sub(latestData.Timestamp)
			status.DataFreshness[feedType] = freshness
			
			// Check if data is stale
			config, _ := fm.keeper.GetFallbackConfig(ctx, feedType)
			if freshness > config.DataStaleness {
				alert := types.SystemAlert{
					AlertType:   "STALE_DATA",
					FeedType:    feedType,
					Severity:    "HIGH",
					Message:     fmt.Sprintf("Data for %s is stale by %v", feedType, freshness),
					Timestamp:   ctx.BlockTime(),
				}
				status.SystemAlerts = append(status.SystemAlerts, alert)
				criticalIssues++
			}
		}
		
		// Check consensus status
		consensus := fm.checkConsensusStatus(ctx, feedType)
		status.ConsensusStatus[feedType] = consensus
		
		if consensus.ConsensusStrength.LT(sdk.NewDecWithPrec(7, 1)) { // < 70%
			degradedFeeds++
		}
		
		if consensus.ConsensusStrength.LT(sdk.NewDecWithPrec(5, 1)) { // < 50%
			criticalIssues++
		}
	}

	// Count active and failed oracles
	allNodes := fm.keeper.GetAllOracleNodes(ctx)
	for _, nodeID := range allNodes {
		node, found := fm.keeper.GetOracleNode(ctx, nodeID)
		if !found {
			continue
		}
		
		if node.Status == "ACTIVE" {
			status.ActiveOracles++
		} else {
			status.FailedOracles++
		}
	}

	// Determine overall health
	if criticalIssues > 0 || status.ActiveOracles < 3 {
		status.OverallHealth = "CRITICAL"
	} else if degradedFeeds > 0 || status.ActiveOracles < 5 {
		status.OverallHealth = "DEGRADED"
	} else {
		status.OverallHealth = "HEALTHY"
	}

	// Check if emergency mode should be activated
	if status.OverallHealth == "CRITICAL" {
		fm.activateEmergencyMode(ctx, status)
		status.EmergencyMode = true
	}

	// Store health status
	fm.keeper.SetOracleHealthStatus(ctx, *status)

	return status, nil
}

// HandleOracleFailure handles individual oracle node failures
func (fm *FallbackManager) HandleOracleFailure(ctx sdk.Context, nodeID string, failureType string) error {
	// Get oracle node
	node, found := fm.keeper.GetOracleNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("oracle node not found: %s", nodeID)
	}

	// Record failure
	failure := types.OracleFailure{
		NodeID:      nodeID,
		FailureType: failureType,
		Timestamp:   ctx.BlockTime(),
		Severity:    fm.assessFailureSeverity(failureType),
		Recovery:    false,
	}
	
	fm.keeper.SetOracleFailure(ctx, failure)

	// Update node status based on failure severity
	switch failure.Severity {
	case "LOW":
		// Apply minor penalty
		fm.applyFailurePenalty(ctx, &node, sdk.NewDecWithPrec(1, 2)) // 1% penalty
	case "MEDIUM":
		// Apply moderate penalty and reduce reputation
		fm.applyFailurePenalty(ctx, &node, sdk.NewDecWithPrec(5, 2)) // 5% penalty
		node.ReputationScore = node.ReputationScore.Sub(sdk.NewDec(10))
	case "HIGH":
		// Suspend node temporarily
		node.Status = "SUSPENDED"
		fm.keeper.RemoveActiveNode(ctx, nodeID)
		fm.applyFailurePenalty(ctx, &node, sdk.NewDecWithPrec(15, 2)) // 15% penalty
	case "CRITICAL":
		// Immediate slashing and suspension
		nodeManager := NewNodeManager(fm.keeper)
		nodeManager.SlashNode(ctx, nodeID, failureType, sdk.NewDecWithPrec(25, 2)) // 25% slash
	}

	// Check if fallback is needed for affected feeds
	for _, feedType := range node.SupportedFeeds {
		fm.checkAndActivateFallback(ctx, feedType, nodeID)
	}

	// Store updated node
	fm.keeper.SetOracleNode(ctx, node)

	// Emit failure event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleFailure,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyFailureType, failureType),
			sdk.NewAttribute(types.AttributeKeySeverity, failure.Severity),
		),
	)

	return nil
}

// ActivateFallback activates fallback mechanism for a specific feed
func (fm *FallbackManager) ActivateFallback(ctx sdk.Context, feedType string, reason string) error {
	// Get fallback configuration
	config, found := fm.keeper.GetFallbackConfig(ctx, feedType)
	if !found {
		return fmt.Errorf("fallback config not found for feed: %s", feedType)
	}

	// Determine fallback level based on current state
	level := fm.determineFallbackLevel(ctx, feedType)

	switch level {
	case "SECONDARY":
		// Activate secondary oracles
		err := fm.activateSecondaryOracles(ctx, feedType, config.SecondaryOracles)
		if err != nil {
			return fmt.Errorf("failed to activate secondary oracles: %w", err)
		}
	case "EMERGENCY":
		// Activate emergency oracles
		err := fm.activateEmergencyOracles(ctx, feedType, config.EmergencyOracles)
		if err != nil {
			return fmt.Errorf("failed to activate emergency oracles: %w", err)
		}
	case "PROTOCOL":
		// Execute emergency protocol
		err := fm.executeEmergencyProtocol(ctx, feedType, config.EmergencyProtocol)
		if err != nil {
			return fmt.Errorf("failed to execute emergency protocol: %w", err)
		}
	}

	// Record fallback activation
	fallback := types.FallbackActivation{
		FeedType:    feedType,
		Level:       level,
		Reason:      reason,
		ActivatedAt: ctx.BlockTime(),
		Status:      "ACTIVE",
	}
	
	fm.keeper.SetFallbackActivation(ctx, fallback)

	// Emit fallback activation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFallbackActivated,
			sdk.NewAttribute(types.AttributeKeyFeedType, feedType),
			sdk.NewAttribute(types.AttributeKeyFallbackLevel, level),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)

	return nil
}

// RecoverFromFailure attempts to recover from oracle failures
func (fm *FallbackManager) RecoverFromFailure(ctx sdk.Context, feedType string) error {
	// Check current fallback status
	fallback, found := fm.keeper.GetActiveFallback(ctx, feedType)
	if !found {
		return nil // No active fallback to recover from
	}

	// Check if primary oracles are healthy again
	primaryHealth := fm.checkPrimaryOracleHealth(ctx, feedType)
	
	if primaryHealth.HealthScore.GTE(sdk.NewDecWithPrec(8, 1)) { // >= 80% health
		// Gradually switch back to primary oracles
		err := fm.executeGradualRecovery(ctx, feedType, fallback)
		if err != nil {
			return fmt.Errorf("failed to execute gradual recovery: %w", err)
		}

		// Deactivate fallback
		fallback.Status = "RECOVERED"
		fallback.RecoveredAt = &ctx.BlockTime()
		fm.keeper.SetFallbackActivation(ctx, fallback)

		// Emit recovery event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeFallbackRecovered,
				sdk.NewAttribute(types.AttributeKeyFeedType, feedType),
				sdk.NewAttribute(types.AttributeKeyRecoveryTime, ctx.BlockTime().String()),
			),
		)
	}

	return nil
}

// Emergency protocol implementations

func (fm *FallbackManager) activateEmergencyMode(ctx sdk.Context, status *OracleHealthStatus) {
	// Implement circuit breaker mechanism
	for feedType := range status.FeedHealthStatus {
		// Pause non-critical operations
		fm.pauseNonCriticalOperations(ctx, feedType)
		
		// Activate emergency price feeds
		fm.activateEmergencyPriceFeeds(ctx, feedType)
		
		// Notify governance of emergency
		fm.notifyGovernanceEmergency(ctx, feedType, status)
	}

	// Store emergency activation
	emergency := types.EmergencyActivation{
		Reason:      "CRITICAL_ORACLE_HEALTH",
		ActivatedAt: ctx.BlockTime(),
		Status:      "ACTIVE",
		Triggers:    fm.extractEmergencyTriggers(status),
	}
	
	fm.keeper.SetEmergencyActivation(ctx, emergency)
}

func (fm *FallbackManager) createDefaultFallbackConfig(ctx sdk.Context, feedType string) FallbackConfig {
	config := FallbackConfig{
		FeedType:           feedType,
		DataStaleness:      5 * time.Minute,  // 5 minutes for most feeds
		MinimumOracles:     3,                // Minimum 3 oracles for consensus
		ConsensusThreshold: sdk.NewDecWithPrec(67, 2), // 67% consensus threshold
		LastUpdate:         ctx.BlockTime(),
	}

	// Feed-specific configurations
	switch feedType {
	case "INR_USD":
		config.DataStaleness = 2 * time.Minute // More frequent updates for critical feed
		config.MinimumOracles = 5
		config.ConsensusThreshold = sdk.NewDecWithPrec(75, 2) // 75% consensus
		config.FallbackThresholds = types.FallbackThresholds{
			PrimaryFailureThreshold:   sdk.NewDecWithPrec(3, 1), // 30%
			SecondaryFailureThreshold: sdk.NewDecWithPrec(5, 1), // 50%
			EmergencyThreshold:        sdk.NewDecWithPrec(7, 1), // 70%
		}
	case "BTC_INR", "ETH_INR":
		config.DataStaleness = 3 * time.Minute
		config.MinimumOracles = 4
		config.FallbackThresholds = types.FallbackThresholds{
			PrimaryFailureThreshold:   sdk.NewDecWithPrec(4, 1), // 40%
			SecondaryFailureThreshold: sdk.NewDecWithPrec(6, 1), // 60%
			EmergencyThreshold:        sdk.NewDecWithPrec(8, 1), // 80%
		}
	case "WEATHER":
		config.DataStaleness = 15 * time.Minute // Weather data can be less frequent
		config.MinimumOracles = 3
		config.FallbackThresholds = types.FallbackThresholds{
			PrimaryFailureThreshold:   sdk.NewDecWithPrec(5, 1), // 50%
			SecondaryFailureThreshold: sdk.NewDecWithPrec(7, 1), // 70%
			EmergencyThreshold:        sdk.NewDecWithPrec(9, 1), // 90%
		}
	default:
		config.FallbackThresholds = types.FallbackThresholds{
			PrimaryFailureThreshold:   sdk.NewDecWithPrec(4, 1), // 40%
			SecondaryFailureThreshold: sdk.NewDecWithPrec(6, 1), // 60%
			EmergencyThreshold:        sdk.NewDecWithPrec(8, 1), // 80%
		}
	}

	// Set up recovery mechanism
	config.RecoveryMechanism = types.RecoveryMechanism{
		AutoRecovery:        true,
		RecoveryDelay:       10 * time.Minute,
		GradualRecovery:     true,
		HealthCheckInterval: 1 * time.Minute,
		RecoveryThreshold:   sdk.NewDecWithPrec(8, 1), // 80% health for recovery
	}

	// Set up emergency protocol
	config.EmergencyProtocol = types.EmergencyProtocol{
		CircuitBreaker:      true,
		EmergencyGovernance: true,
		ExternalDataSources: fm.getEmergencyDataSources(feedType),
		ManualOverride:      true,
		NotificationChannels: []string{"GOVERNANCE", "VALIDATORS", "TREASURY"},
	}

	return config
}

func (fm *FallbackManager) checkFeedHealth(ctx sdk.Context, feedType string) types.FeedHealth {
	health := types.FeedHealth{
		FeedType: feedType,
	}

	// Get recent data points
	recentData := fm.keeper.GetRecentOracleData(ctx, feedType, 10)
	if len(recentData) == 0 {
		health.HealthScore = sdk.ZeroDec()
		health.Status = "NO_DATA"
		return health
	}

	// Calculate data consistency
	consistency := fm.calculateDataConsistency(recentData)
	health.ConsistencyScore = consistency

	// Calculate oracle participation
	participation := fm.calculateOracleParticipation(ctx, feedType)
	health.ParticipationRate = participation

	// Calculate overall health score
	healthScore := consistency.Mul(sdk.NewDecWithPrec(6, 1)).Add(participation.Mul(sdk.NewDecWithPrec(4, 1)))
	health.HealthScore = healthScore

	// Determine status
	if healthScore.GTE(sdk.NewDecWithPrec(8, 1)) {
		health.Status = "HEALTHY"
	} else if healthScore.GTE(sdk.NewDecWithPrec(6, 1)) {
		health.Status = "DEGRADED"
	} else {
		health.Status = "CRITICAL"
	}

	health.LastCheck = ctx.BlockTime()
	return health
}

func (fm *FallbackManager) assessFailureSeverity(failureType string) string {
	severityMap := map[string]string{
		"DATA_TIMEOUT":        "MEDIUM",
		"INVALID_SIGNATURE":   "HIGH",
		"PRICE_DEVIATION":     "MEDIUM",
		"CONSENSUS_FAILURE":   "HIGH",
		"NETWORK_TIMEOUT":     "LOW",
		"MALICIOUS_BEHAVIOR":  "CRITICAL",
		"DATA_CORRUPTION":     "HIGH",
		"REPEATED_FAILURES":   "CRITICAL",
	}

	if severity, found := severityMap[failureType]; found {
		return severity
	}
	return "MEDIUM" // Default severity
}

func (fm *FallbackManager) determineFallbackLevel(ctx sdk.Context, feedType string) string {
	// Get current oracle status
	activeOracles := fm.keeper.GetActiveOraclesForFeed(ctx, feedType)
	config, _ := fm.keeper.GetFallbackConfig(ctx, feedType)

	failureRate := fm.calculateCurrentFailureRate(ctx, feedType)

	if failureRate.GTE(config.FallbackThresholds.EmergencyThreshold) {
		return "PROTOCOL"
	} else if failureRate.GTE(config.FallbackThresholds.SecondaryFailureThreshold) {
		return "EMERGENCY"
	} else if failureRate.GTE(config.FallbackThresholds.PrimaryFailureThreshold) {
		return "SECONDARY"
	}

	return "PRIMARY"
}

func (fm *FallbackManager) getEmergencyDataSources(feedType string) []string {
	emergencySources := map[string][]string{
		"INR_USD": {"RBI_REFERENCE_RATE", "FOREX_BACKUP", "CENTRAL_BANK_FEED"},
		"BTC_INR": {"COINBASE_API", "BINANCE_API", "COINGECKO_API"},
		"ETH_INR": {"COINBASE_API", "BINANCE_API", "COINGECKO_API"},
		"WEATHER": {"IMD_API", "OPENWEATHER_API", "WEATHER_GOV_API"},
	}

	if sources, found := emergencySources[feedType]; found {
		return sources
	}
	return []string{"EXTERNAL_API_BACKUP"}
}

// Helper utility functions
func (fm *FallbackManager) initializeHealthMonitoring(ctx sdk.Context) {
	// Set up initial health check schedule
	fm.keeper.SetHealthCheckSchedule(ctx, types.HealthCheckSchedule{
		Interval:    1 * time.Minute,
		LastCheck:   ctx.BlockTime(),
		NextCheck:   ctx.BlockTime().Add(1 * time.Minute),
		Enabled:     true,
	})
}

// Additional helper methods would include:
// - checkConsensusStatus
// - checkAndActivateFallback
// - activateSecondaryOracles
// - activateEmergencyOracles
// - executeEmergencyProtocol
// - checkPrimaryOracleHealth
// - executeGradualRecovery
// - pauseNonCriticalOperations
// - activateEmergencyPriceFeeds
// - notifyGovernanceEmergency
// - extractEmergencyTriggers
// - calculateDataConsistency
// - calculateOracleParticipation
// - calculateCurrentFailureRate
// - applyFailurePenalty