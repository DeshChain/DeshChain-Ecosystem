package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/oracle/types"
)

// NodeManager handles oracle node registration, management, and incentives
type NodeManager struct {
	keeper Keeper
}

// NewNodeManager creates a new oracle node manager
func NewNodeManager(keeper Keeper) *NodeManager {
	return &NodeManager{
		keeper: keeper,
	}
}

// OracleNode represents a registered oracle node
type OracleNode struct {
	NodeID          string                    `json:"node_id"`
	OperatorAddress string                    `json:"operator_address"`
	NodeInfo        types.NodeInfo            `json:"node_info"`
	NodeType        string                    `json:"node_type"` // PRICE_FEED, WEATHER, MARKET_DATA, CUSTOM
	DataSources     []string                  `json:"data_sources"`
	SupportedFeeds  []string                  `json:"supported_feeds"`
	StakeAmount     sdk.Coin                  `json:"stake_amount"`
	Performance     types.NodePerformance     `json:"performance"`
	ReputationScore sdk.Dec                   `json:"reputation_score"`
	Earnings        types.NodeEarnings        `json:"earnings"`
	Status          string                    `json:"status"` // ACTIVE, INACTIVE, SUSPENDED, SLASHED
	LastUpdate      time.Time                 `json:"last_update"`
	RegisteredAt    time.Time                 `json:"registered_at"`
	SlashingHistory []types.SlashingEvent     `json:"slashing_history"`
	IncentiveConfig types.NodeIncentiveConfig `json:"incentive_config"`
}

// RegisterOracleNode registers a new oracle node
func (nm *NodeManager) RegisterOracleNode(ctx sdk.Context, request types.NodeRegistrationRequest) (*OracleNode, error) {
	// Validate registration request
	if err := nm.validateRegistrationRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("registration validation failed: %w", err)
	}

	// Check minimum stake requirement
	params := nm.keeper.GetParams(ctx)
	if !request.StakeAmount.IsGTE(params.MinNodeStake) {
		return nil, fmt.Errorf("insufficient stake amount: required %s, provided %s", 
			params.MinNodeStake.String(), request.StakeAmount.String())
	}

	// Generate node ID
	nodeID := nm.generateNodeID(ctx, request.OperatorAddress)

	// Check if node already exists
	if _, found := nm.keeper.GetOracleNode(ctx, nodeID); found {
		return nil, fmt.Errorf("oracle node already registered for operator: %s", request.OperatorAddress)
	}

	// Lock stake tokens
	operatorAddr, err := sdk.AccAddressFromBech32(request.OperatorAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid operator address: %w", err)
	}

	err = nm.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx, operatorAddr, types.ModuleName, sdk.NewCoins(request.StakeAmount),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to lock stake: %w", err)
	}

	// Create oracle node
	node := &OracleNode{
		NodeID:          nodeID,
		OperatorAddress: request.OperatorAddress,
		NodeInfo: types.NodeInfo{
			Name:        request.NodeName,
			Description: request.Description,
			Website:     request.Website,
			Location:    request.Location,
			PublicKey:   request.PublicKey,
		},
		NodeType:       request.NodeType,
		DataSources:    request.DataSources,
		SupportedFeeds: request.SupportedFeeds,
		StakeAmount:    request.StakeAmount,
		Performance: types.NodePerformance{
			NodeID:            nodeID,
			TotalSubmissions:  0,
			AcceptedSubmissions: 0,
			RejectedSubmissions: 0,
			AccuracyScore:     sdk.ZeroDec(),
			UptimeScore:       sdk.NewDec(100), // Start with perfect uptime
			ResponseTime:      0,
			LastSubmission:    ctx.BlockTime(),
		},
		ReputationScore: sdk.NewDec(50), // Start with neutral reputation
		Earnings: types.NodeEarnings{
			NodeID:           nodeID,
			TotalEarned:      sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentPeriod:    sdk.NewCoin("namo", sdk.ZeroInt()),
			LastPayout:       ctx.BlockTime(),
			PendingRewards:   sdk.NewCoin("namo", sdk.ZeroInt()),
		},
		Status:          "ACTIVE",
		LastUpdate:      ctx.BlockTime(),
		RegisteredAt:    ctx.BlockTime(),
		SlashingHistory: []types.SlashingEvent{},
		IncentiveConfig: nm.createDefaultIncentiveConfig(request.NodeType),
	}

	// Store oracle node
	nm.keeper.SetOracleNode(ctx, *node)

	// Add to active nodes list
	nm.keeper.AddActiveNode(ctx, nodeID)

	// Emit registration event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleNodeRegistered,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyOperator, request.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyNodeType, request.NodeType),
			sdk.NewAttribute(types.AttributeKeyStakeAmount, request.StakeAmount.String()),
		),
	)

	return node, nil
}

// SubmitOracleData handles data submission from oracle nodes
func (nm *NodeManager) SubmitOracleData(ctx sdk.Context, nodeID string, submission types.DataSubmission) error {
	// Get oracle node
	node, found := nm.keeper.GetOracleNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("oracle node not found: %s", nodeID)
	}

	// Validate node status
	if node.Status != "ACTIVE" {
		return fmt.Errorf("node not active: %s", node.Status)
	}

	// Validate submission
	if err := nm.validateDataSubmission(ctx, node, submission); err != nil {
		return fmt.Errorf("data submission validation failed: %w", err)
	}

	// Process submission based on feed type
	accepted, err := nm.processDataSubmission(ctx, node, submission)
	if err != nil {
		return fmt.Errorf("failed to process data submission: %w", err)
	}

	// Update node performance
	node.Performance.TotalSubmissions++
	if accepted {
		node.Performance.AcceptedSubmissions++
		
		// Calculate and distribute rewards
		reward := nm.calculateSubmissionReward(ctx, node, submission)
		if reward.IsPositive() {
			nm.distributeNodeReward(ctx, &node, reward, "DATA_SUBMISSION")
		}
	} else {
		node.Performance.RejectedSubmissions++
		
		// Apply penalty for rejected submission
		penalty := nm.calculateRejectionPenalty(ctx, node, submission)
		if penalty.IsPositive() {
			nm.applyNodePenalty(ctx, &node, penalty, "REJECTED_SUBMISSION")
		}
	}

	// Update performance metrics
	nm.updateNodePerformanceMetrics(ctx, &node)

	// Update reputation score
	node.ReputationScore = nm.calculateReputationScore(ctx, node)

	// Update last activity
	node.LastUpdate = ctx.BlockTime()

	// Store updated node
	nm.keeper.SetOracleNode(ctx, node)

	// Emit submission event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleDataSubmitted,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyFeedType, submission.FeedType),
			sdk.NewAttribute(types.AttributeKeyAccepted, fmt.Sprintf("%t", accepted)),
		),
	)

	return nil
}

// ProcessPeriodicIncentives processes periodic incentives for oracle nodes
func (nm *NodeManager) ProcessPeriodicIncentives(ctx sdk.Context) error {
	// Get all active oracle nodes
	activeNodes := nm.keeper.GetActiveOracleNodes(ctx)

	for _, nodeID := range activeNodes {
		node, found := nm.keeper.GetOracleNode(ctx, nodeID)
		if !found {
			continue
		}

		// Calculate periodic incentives
		incentives := nm.calculatePeriodicIncentives(ctx, node)

		// Distribute incentives
		for incentiveType, amount := range incentives {
			if amount.IsPositive() {
				nm.distributeNodeReward(ctx, &node, amount, incentiveType)
			}
		}

		// Update node earnings
		nm.keeper.SetOracleNode(ctx, node)
	}

	// Emit incentive processing event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePeriodicIncentivesProcessed,
			sdk.NewAttribute(types.AttributeKeyProcessedNodes, fmt.Sprintf("%d", len(activeNodes))),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)

	return nil
}

// SlashNode applies slashing to misbehaving oracle nodes
func (nm *NodeManager) SlashNode(ctx sdk.Context, nodeID string, reason string, severity sdk.Dec) error {
	// Get oracle node
	node, found := nm.keeper.GetOracleNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("oracle node not found: %s", nodeID)
	}

	// Calculate slash amount
	slashAmount := node.StakeAmount.Amount.ToDec().Mul(severity).TruncateInt()
	slashedCoins := sdk.NewCoin(node.StakeAmount.Denom, slashAmount)

	// Apply slashing
	err := nm.keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(slashedCoins))
	if err != nil {
		return fmt.Errorf("failed to burn slashed tokens: %w", err)
	}

	// Update node stake
	node.StakeAmount = node.StakeAmount.Sub(slashedCoins)

	// Record slashing event
	slashingEvent := types.SlashingEvent{
		NodeID:      nodeID,
		Reason:      reason,
		Severity:    severity,
		SlashedAmount: slashedCoins,
		Timestamp:   ctx.BlockTime(),
		BlockHeight: ctx.BlockHeight(),
	}
	
	node.SlashingHistory = append(node.SlashingHistory, slashingEvent)

	// Update reputation score (significant penalty)
	reputationPenalty := severity.Mul(sdk.NewDec(50)) // Up to 50 point penalty
	node.ReputationScore = node.ReputationScore.Sub(reputationPenalty)
	if node.ReputationScore.LT(sdk.ZeroDec()) {
		node.ReputationScore = sdk.ZeroDec()
	}

	// Check if node should be suspended
	if node.ReputationScore.LT(sdk.NewDec(20)) || node.StakeAmount.Amount.LT(nm.keeper.GetParams(ctx).MinNodeStake.Amount) {
		node.Status = "SUSPENDED"
		nm.keeper.RemoveActiveNode(ctx, nodeID)
	}

	// Store updated node
	nm.keeper.SetOracleNode(ctx, node)

	// Emit slashing event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleNodeSlashed,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
			sdk.NewAttribute(types.AttributeKeySlashedAmount, slashedCoins.String()),
			sdk.NewAttribute(types.AttributeKeyNewStatus, node.Status),
		),
	)

	return nil
}

// UnregisterNode handles node unregistration and stake return
func (nm *NodeManager) UnregisterNode(ctx sdk.Context, nodeID string, operatorAddr string) error {
	// Get oracle node
	node, found := nm.keeper.GetOracleNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("oracle node not found: %s", nodeID)
	}

	// Validate operator
	if node.OperatorAddress != operatorAddr {
		return fmt.Errorf("unauthorized: only node operator can unregister")
	}

	// Check if node has pending obligations
	if nm.hasPendingObligations(ctx, node) {
		return fmt.Errorf("node has pending obligations, cannot unregister")
	}

	// Calculate final earnings
	finalEarnings := nm.calculateFinalEarnings(ctx, node)
	
	// Return stake and earnings
	operatorAddress, err := sdk.AccAddressFromBech32(node.OperatorAddress)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	totalReturn := sdk.NewCoins(node.StakeAmount).Add(finalEarnings...)
	err = nm.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, operatorAddress, totalReturn,
	)
	if err != nil {
		return fmt.Errorf("failed to return stake and earnings: %w", err)
	}

	// Update status
	node.Status = "INACTIVE"
	node.LastUpdate = ctx.BlockTime()

	// Remove from active nodes
	nm.keeper.RemoveActiveNode(ctx, nodeID)

	// Store final node state
	nm.keeper.SetOracleNode(ctx, node)

	// Emit unregistration event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleNodeUnregistered,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyOperator, operatorAddr),
			sdk.NewAttribute(types.AttributeKeyReturnedAmount, totalReturn.String()),
		),
	)

	return nil
}

// Helper functions for node management

func (nm *NodeManager) validateRegistrationRequest(ctx sdk.Context, request types.NodeRegistrationRequest) error {
	// Validate operator address
	_, err := sdk.AccAddressFromBech32(request.OperatorAddress)
	if err != nil {
		return fmt.Errorf("invalid operator address: %w", err)
	}

	// Validate node type
	validTypes := []string{"PRICE_FEED", "WEATHER", "MARKET_DATA", "CUSTOM"}
	validType := false
	for _, vt := range validTypes {
		if request.NodeType == vt {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid node type: %s", request.NodeType)
	}

	// Validate data sources
	if len(request.DataSources) == 0 {
		return fmt.Errorf("at least one data source must be specified")
	}

	// Validate supported feeds
	if len(request.SupportedFeeds) == 0 {
		return fmt.Errorf("at least one supported feed must be specified")
	}

	return nil
}

func (nm *NodeManager) calculateSubmissionReward(ctx sdk.Context, node OracleNode, submission types.DataSubmission) sdk.Coin {
	params := nm.keeper.GetParams(ctx)
	baseReward := params.BaseSubmissionReward

	// Apply multipliers based on node performance
	reputationMultiplier := node.ReputationScore.QuoInt64(100) // 0-1 based on reputation
	accuracyMultiplier := node.Performance.AccuracyScore       // 0-1 based on accuracy

	// Apply feed type multiplier
	feedMultiplier := nm.getFeedTypeMultiplier(submission.FeedType)

	// Calculate final reward
	finalReward := baseReward.Amount.ToDec().
		Mul(reputationMultiplier).
		Mul(accuracyMultiplier).
		Mul(feedMultiplier)

	return sdk.NewCoin(baseReward.Denom, finalReward.TruncateInt())
}

func (nm *NodeManager) calculatePeriodicIncentives(ctx sdk.Context, node OracleNode) map[string]sdk.Coin {
	incentives := make(map[string]sdk.Coin)
	params := nm.keeper.GetParams(ctx)

	// Uptime incentive
	if node.Performance.UptimeScore.GTE(sdk.NewDecWithPrec(95, 2)) { // >= 95% uptime
		uptimeReward := params.UptimeIncentive.Amount.ToDec().Mul(node.Performance.UptimeScore.QuoInt64(100))
		incentives["UPTIME"] = sdk.NewCoin(params.UptimeIncentive.Denom, uptimeReward.TruncateInt())
	}

	// Accuracy incentive
	if node.Performance.AccuracyScore.GTE(sdk.NewDecWithPrec(90, 2)) { // >= 90% accuracy
		accuracyReward := params.AccuracyIncentive.Amount.ToDec().Mul(node.Performance.AccuracyScore.QuoInt64(100))
		incentives["ACCURACY"] = sdk.NewCoin(params.AccuracyIncentive.Denom, accuracyReward.TruncateInt())
	}

	// Consistency incentive
	if node.Performance.TotalSubmissions > 100 {
		consistencyScore := nm.calculateConsistencyScore(ctx, node)
		if consistencyScore.GTE(sdk.NewDecWithPrec(80, 2)) { // >= 80% consistency
			consistencyReward := params.ConsistencyIncentive.Amount.ToDec().Mul(consistencyScore.QuoInt64(100))
			incentives["CONSISTENCY"] = sdk.NewCoin(params.ConsistencyIncentive.Denom, consistencyReward.TruncateInt())
		}
	}

	return incentives
}

func (nm *NodeManager) calculateReputationScore(ctx sdk.Context, node OracleNode) sdk.Dec {
	weights := map[string]sdk.Dec{
		"accuracy":    sdk.NewDecWithPrec(40, 2), // 40%
		"uptime":      sdk.NewDecWithPrec(25, 2), // 25%
		"consistency": sdk.NewDecWithPrec(20, 2), // 20%
		"longevity":   sdk.NewDecWithPrec(10, 2), // 10%
		"slashing":    sdk.NewDecWithPrec(5, 2),  // 5%
	}

	// Calculate component scores
	accuracyScore := node.Performance.AccuracyScore.Mul(sdk.NewDec(100))
	uptimeScore := node.Performance.UptimeScore
	consistencyScore := nm.calculateConsistencyScore(ctx, node).Mul(sdk.NewDec(100))
	longevityScore := nm.calculateLongevityScore(ctx, node).Mul(sdk.NewDec(100))
	slashingScore := nm.calculateSlashingScore(node).Mul(sdk.NewDec(100))

	// Calculate weighted average
	totalScore := sdk.ZeroDec()
	totalScore = totalScore.Add(accuracyScore.Mul(weights["accuracy"]))
	totalScore = totalScore.Add(uptimeScore.Mul(weights["uptime"]))
	totalScore = totalScore.Add(consistencyScore.Mul(weights["consistency"]))
	totalScore = totalScore.Add(longevityScore.Mul(weights["longevity"]))
	totalScore = totalScore.Add(slashingScore.Mul(weights["slashing"]))

	// Ensure score is within bounds [0, 100]
	if totalScore.GT(sdk.NewDec(100)) {
		totalScore = sdk.NewDec(100)
	}
	if totalScore.LT(sdk.ZeroDec()) {
		totalScore = sdk.ZeroDec()
	}

	return totalScore
}

func (nm *NodeManager) distributeNodeReward(ctx sdk.Context, node *OracleNode, amount sdk.Coin, rewardType string) {
	// Add to pending rewards
	node.Earnings.PendingRewards = node.Earnings.PendingRewards.Add(amount)
	node.Earnings.CurrentPeriod = node.Earnings.CurrentPeriod.Add(amount)

	// Record reward transaction
	reward := types.NodeRewardTransaction{
		NodeID:      node.NodeID,
		Amount:      amount,
		RewardType:  rewardType,
		Timestamp:   ctx.BlockTime(),
		BlockHeight: ctx.BlockHeight(),
	}
	
	nm.keeper.SetNodeRewardTransaction(ctx, reward)
}

func (nm *NodeManager) createDefaultIncentiveConfig(nodeType string) types.NodeIncentiveConfig {
	config := types.NodeIncentiveConfig{
		NodeType: nodeType,
	}

	switch nodeType {
	case "PRICE_FEED":
		config.BaseRewardMultiplier = sdk.NewDecWithPrec(12, 1) // 1.2x
		config.AccuracyBonus = sdk.NewDecWithPrec(5, 1)         // 0.5x bonus for high accuracy
		config.UptimeBonus = sdk.NewDecWithPrec(3, 1)           // 0.3x bonus for high uptime
	case "WEATHER":
		config.BaseRewardMultiplier = sdk.OneDec()              // 1.0x
		config.AccuracyBonus = sdk.NewDecWithPrec(3, 1)         // 0.3x bonus
		config.UptimeBonus = sdk.NewDecWithPrec(2, 1)           // 0.2x bonus
	case "MARKET_DATA":
		config.BaseRewardMultiplier = sdk.NewDecWithPrec(11, 1) // 1.1x
		config.AccuracyBonus = sdk.NewDecWithPrec(4, 1)         // 0.4x bonus
		config.UptimeBonus = sdk.NewDecWithPrec(3, 1)           // 0.3x bonus
	default:
		config.BaseRewardMultiplier = sdk.OneDec()              // 1.0x
		config.AccuracyBonus = sdk.NewDecWithPrec(2, 1)         // 0.2x bonus
		config.UptimeBonus = sdk.NewDecWithPrec(2, 1)           // 0.2x bonus
	}

	return config
}

// Helper utility functions
func (nm *NodeManager) generateNodeID(ctx sdk.Context, operatorAddr string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("NODE-%s-%d", operatorAddr[:8], timestamp)
}

func (nm *NodeManager) getFeedTypeMultiplier(feedType string) sdk.Dec {
	multipliers := map[string]sdk.Dec{
		"INR_USD":     sdk.NewDecWithPrec(15, 1), // 1.5x - critical feed
		"BTC_INR":     sdk.NewDecWithPrec(13, 1), // 1.3x
		"ETH_INR":     sdk.NewDecWithPrec(12, 1), // 1.2x
		"WEATHER":     sdk.OneDec(),              // 1.0x
		"MARKET_DATA": sdk.NewDecWithPrec(11, 1), // 1.1x
	}

	if multiplier, found := multipliers[feedType]; found {
		return multiplier
	}
	return sdk.OneDec() // Default 1.0x
}

// Additional helper methods would include:
// - validateDataSubmission
// - processDataSubmission
// - updateNodePerformanceMetrics
// - calculateRejectionPenalty
// - applyNodePenalty
// - calculateConsistencyScore
// - calculateLongevityScore
// - calculateSlashingScore
// - hasPendingObligations
// - calculateFinalEarnings