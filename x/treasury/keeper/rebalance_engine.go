package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/treasury/types"
)

// RebalanceEngine handles automated treasury rebalancing logic
type RebalanceEngine struct {
	keeper Keeper
}

// NewRebalanceEngine creates a new rebalance engine
func NewRebalanceEngine(keeper Keeper) *RebalanceEngine {
	return &RebalanceEngine{
		keeper: keeper,
	}
}

// RebalancePlan represents a comprehensive rebalancing plan
type RebalancePlan struct {
	RebalanceID       string                    `json:"rebalance_id"`
	TriggerReason     string                    `json:"trigger_reason"`
	Actions           []RebalanceAction         `json:"actions"`
	TotalRebalanced   sdk.Coins                 `json:"total_rebalanced"`
	TargetDeviations  map[string]sdk.Dec        `json:"target_deviations"`
	ExpectedOutcome   types.RebalanceOutcome    `json:"expected_outcome"`
	RiskAssessment    types.RebalanceRisk       `json:"risk_assessment"`
	ExecutionStrategy types.ExecutionStrategy   `json:"execution_strategy"`
	Timeline          types.RebalanceTimeline   `json:"timeline"`
	CreatedAt         time.Time                 `json:"created_at"`
	ExecutedAt        *time.Time                `json:"executed_at,omitempty"`
	Status            string                    `json:"status"`
}

// RebalanceAction represents a single rebalancing action
type RebalanceAction struct {
	ActionID        string    `json:"action_id"`
	ActionType      string    `json:"action_type"` // TRANSFER, MINT, BURN, SWAP
	SourcePool      string    `json:"source_pool"`
	DestinationPool string    `json:"destination_pool"`
	Amount          sdk.Coins `json:"amount"`
	Priority        int       `json:"priority"`
	Reason          string    `json:"reason"`
	ExpectedImpact  sdk.Dec   `json:"expected_impact"`
	RiskLevel       string    `json:"risk_level"`
	Prerequisites   []string  `json:"prerequisites"`
	ExecutedAt      *time.Time `json:"executed_at,omitempty"`
	Status          string    `json:"status"`
}

// AnalyzeRebalanceNeed analyzes if treasury rebalancing is needed
func (re *RebalanceEngine) AnalyzeRebalanceNeed(ctx sdk.Context, pools []TreasuryPool) (bool, *RebalancePlan) {
	// Calculate current total treasury value
	totalValue := re.calculateTotalTreasuryValue(pools)
	
	// Analyze deviations from target allocations
	deviations := re.calculateAllocationDeviations(pools, totalValue)
	
	// Check if any pool exceeds rebalance threshold
	rebalanceNeeded := false
	var triggerReasons []string
	
	for poolID, deviation := range deviations {
		pool := re.findPoolByID(pools, poolID)
		if pool == nil {
			continue
		}
		
		// Check if deviation exceeds threshold
		if deviation.Abs().GTE(pool.RebalanceConfig.ThresholdPercent) {
			rebalanceNeeded = true
			if deviation.IsPositive() {
				triggerReasons = append(triggerReasons, fmt.Sprintf("Pool %s over-allocated by %s", poolID, deviation.String()))
			} else {
				triggerReasons = append(triggerReasons, fmt.Sprintf("Pool %s under-allocated by %s", poolID, deviation.Abs().String()))
			}
		}
	}
	
	// Check time-based rebalancing requirements
	for _, pool := range pools {
		if pool.RebalanceConfig.Enabled {
			timeSinceLastRebalance := ctx.BlockTime().Sub(pool.LastRebalance)
			if timeSinceLastRebalance > pool.RebalanceConfig.MaxRebalanceGap {
				rebalanceNeeded = true
				triggerReasons = append(triggerReasons, fmt.Sprintf("Pool %s exceeds maximum rebalance gap", pool.PoolID))
			}
		}
	}
	
	if !rebalanceNeeded {
		return false, nil
	}
	
	// Create rebalancing plan
	plan := &RebalancePlan{
		RebalanceID:      re.generateRebalanceID(ctx),
		TriggerReason:    fmt.Sprintf("Multiple triggers: %v", triggerReasons),
		TargetDeviations: deviations,
		CreatedAt:        ctx.BlockTime(),
		Status:          "PLANNED",
	}
	
	// Generate rebalancing actions
	plan.Actions = re.generateRebalanceActions(ctx, pools, deviations, totalValue)
	
	// Calculate expected outcome
	plan.ExpectedOutcome = re.calculateExpectedOutcome(ctx, pools, plan.Actions)
	
	// Assess rebalancing risks
	plan.RiskAssessment = re.assessRebalanceRisks(ctx, plan.Actions)
	
	// Determine execution strategy
	plan.ExecutionStrategy = re.determineExecutionStrategy(ctx, plan.Actions, plan.RiskAssessment)
	
	// Create execution timeline
	plan.Timeline = re.createExecutionTimeline(ctx, plan.Actions)
	
	// Calculate total rebalanced amount
	var totalRebalanced sdk.Coins
	for _, action := range plan.Actions {
		totalRebalanced = totalRebalanced.Add(action.Amount...)
	}
	plan.TotalRebalanced = totalRebalanced
	
	return true, plan
}

// ExecuteRebalancePlan executes a rebalancing plan
func (re *RebalanceEngine) ExecuteRebalancePlan(ctx sdk.Context, plan *RebalancePlan) error {
	// Validate plan before execution
	if err := re.validateRebalancePlan(ctx, plan); err != nil {
		return fmt.Errorf("rebalance plan validation failed: %w", err)
	}
	
	// Sort actions by priority
	sort.Slice(plan.Actions, func(i, j int) bool {
		return plan.Actions[i].Priority < plan.Actions[j].Priority
	})
	
	// Execute actions in order
	for i := range plan.Actions {
		action := &plan.Actions[i]
		
		// Check prerequisites
		if err := re.checkActionPrerequisites(ctx, *action); err != nil {
			return fmt.Errorf("action %s prerequisites not met: %w", action.ActionID, err)
		}
		
		// Execute action
		err := re.executeRebalanceAction(ctx, action)
		if err != nil {
			return fmt.Errorf("failed to execute action %s: %w", action.ActionID, err)
		}
		
		// Mark action as executed
		now := ctx.BlockTime()
		action.ExecutedAt = &now
		action.Status = "EXECUTED"
		
		// Emit action execution event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRebalanceActionExecuted,
				sdk.NewAttribute(types.AttributeKeyActionID, action.ActionID),
				sdk.NewAttribute(types.AttributeKeyActionType, action.ActionType),
				sdk.NewAttribute(types.AttributeKeyAmount, action.Amount.String()),
			),
		)
	}
	
	// Update plan status
	now := ctx.BlockTime()
	plan.ExecutedAt = &now
	plan.Status = "EXECUTED"
	
	// Store executed plan
	re.keeper.SetRebalancePlan(ctx, *plan)
	
	return nil
}

// generateRebalanceActions generates specific rebalancing actions
func (re *RebalanceEngine) generateRebalanceActions(ctx sdk.Context, pools []TreasuryPool, deviations map[string]sdk.Dec, totalValue sdk.Coins) []RebalanceAction {
	var actions []RebalanceAction
	actionCounter := 0
	
	// Separate over-allocated and under-allocated pools
	overAllocated := make(map[string]sdk.Dec)
	underAllocated := make(map[string]sdk.Dec)
	
	for poolID, deviation := range deviations {
		if deviation.IsPositive() {
			overAllocated[poolID] = deviation
		} else if deviation.IsNegative() {
			underAllocated[poolID] = deviation.Abs()
		}
	}
	
	// Generate transfer actions from over-allocated to under-allocated pools
	for overPoolID, overDeviation := range overAllocated {
		overPool := re.findPoolByID(pools, overPoolID)
		if overPool == nil {
			continue
		}
		
		// Calculate excess amount
		targetValue := totalValue.AmountOf("namo").ToDec().Mul(overPool.Allocation.TargetPercentage)
		currentValue := overPool.Balance.AmountOf("namo").ToDec()
		excessAmount := currentValue.Sub(targetValue)
		
		if excessAmount.LTE(sdk.ZeroDec()) {
			continue
		}
		
		// Find best recipient pools
		for underPoolID, underDeviation := range underAllocated {
			if excessAmount.LTE(sdk.ZeroDec()) {
				break
			}
			
			underPool := re.findPoolByID(pools, underPoolID)
			if underPool == nil {
				continue
			}
			
			// Calculate needed amount for under-allocated pool
			targetValue := totalValue.AmountOf("namo").ToDec().Mul(underPool.Allocation.TargetPercentage)
			currentValue := underPool.Balance.AmountOf("namo").ToDec()
			neededAmount := targetValue.Sub(currentValue)
			
			if neededAmount.LTE(sdk.ZeroDec()) {
				continue
			}
			
			// Calculate transfer amount (minimum of excess and needed)
			transferAmount := sdk.MinDec(excessAmount, neededAmount)
			
			// Create transfer action
			action := RebalanceAction{
				ActionID:        fmt.Sprintf("REBAL_%d_%d", ctx.BlockHeight(), actionCounter),
				ActionType:      "TRANSFER",
				SourcePool:      overPoolID,
				DestinationPool: underPoolID,
				Amount:          sdk.NewCoins(sdk.NewCoin("namo", transferAmount.TruncateInt())),
				Priority:        re.calculateActionPriority(overPool, underPool),
				Reason:          fmt.Sprintf("Rebalance from over-allocated %s to under-allocated %s", overPoolID, underPoolID),
				ExpectedImpact:  transferAmount.Quo(totalValue.AmountOf("namo").ToDec()),
				RiskLevel:       re.assessActionRisk(overPool, underPool, transferAmount),
				Prerequisites:   []string{},
				Status:          "PENDING",
			}
			
			actions = append(actions, action)
			actionCounter++
			
			// Update remaining amounts
			excessAmount = excessAmount.Sub(transferAmount)
			underAllocated[underPoolID] = underAllocated[underPoolID].Sub(transferAmount.Quo(totalValue.AmountOf("namo").ToDec()))
		}
	}
	
	// Add emergency actions if needed
	emergencyActions := re.generateEmergencyActions(ctx, pools, totalValue)
	actions = append(actions, emergencyActions...)
	
	return actions
}

// executeRebalanceAction executes a specific rebalance action
func (re *RebalanceEngine) executeRebalanceAction(ctx sdk.Context, action *RebalanceAction) error {
	switch action.ActionType {
	case "TRANSFER":
		return re.executeTransferAction(ctx, action)
	case "MINT":
		return re.executeMintAction(ctx, action)
	case "BURN":
		return re.executeBurnAction(ctx, action)
	case "SWAP":
		return re.executeSwapAction(ctx, action)
	default:
		return fmt.Errorf("unknown action type: %s", action.ActionType)
	}
}

// executeTransferAction executes a transfer between pools
func (re *RebalanceEngine) executeTransferAction(ctx sdk.Context, action *RebalanceAction) error {
	// Get source pool
	sourcePool, found := re.keeper.GetTreasuryPool(ctx, action.SourcePool)
	if !found {
		return fmt.Errorf("source pool not found: %s", action.SourcePool)
	}
	
	// Get destination pool
	destPool, found := re.keeper.GetTreasuryPool(ctx, action.DestinationPool)
	if !found {
		return fmt.Errorf("destination pool not found: %s", action.DestinationPool)
	}
	
	// Validate transfer amount
	if !sourcePool.Balance.IsAllGTE(action.Amount) {
		return fmt.Errorf("insufficient balance in source pool %s", action.SourcePool)
	}
	
	// Check minimum balance constraints
	newSourceBalance := sourcePool.Balance.Sub(action.Amount)
	if !newSourceBalance.IsAllGTE(sourcePool.MinBalance) {
		return fmt.Errorf("transfer would violate minimum balance for pool %s", action.SourcePool)
	}
	
	// Check maximum balance constraints
	newDestBalance := destPool.Balance.Add(action.Amount...)
	if destPool.MaxBalance.IsAllPositive() && !destPool.MaxBalance.IsAllGTE(newDestBalance) {
		return fmt.Errorf("transfer would exceed maximum balance for pool %s", action.DestinationPool)
	}
	
	// Execute the transfer
	sourcePool.Balance = newSourceBalance
	destPool.Balance = newDestBalance
	
	// Update timestamps
	now := ctx.BlockTime()
	sourcePool.LastRebalance = now
	sourcePool.UpdatedAt = now
	destPool.LastRebalance = now
	destPool.UpdatedAt = now
	
	// Save updated pools
	re.keeper.SetTreasuryPool(ctx, sourcePool)
	re.keeper.SetTreasuryPool(ctx, destPool)
	
	// Record transaction
	transaction := types.TreasuryTransaction{
		TransactionID: re.generateTransactionID(ctx, action.ActionID),
		Type:          "REBALANCE_TRANSFER",
		SourcePool:    action.SourcePool,
		DestPool:      action.DestinationPool,
		Amount:        action.Amount,
		Reason:        action.Reason,
		ExecutedBy:    "REBALANCE_ENGINE",
		Timestamp:     now,
		BlockHeight:   ctx.BlockHeight(),
	}
	
	re.keeper.SetTreasuryTransaction(ctx, transaction)
	
	return nil
}

// calculateAllocationDeviations calculates how much each pool deviates from target
func (re *RebalanceEngine) calculateAllocationDeviations(pools []TreasuryPool, totalValue sdk.Coins) map[string]sdk.Dec {
	deviations := make(map[string]sdk.Dec)
	totalNamo := totalValue.AmountOf("namo").ToDec()
	
	if totalNamo.IsZero() {
		return deviations
	}
	
	for _, pool := range pools {
		currentValue := pool.Balance.AmountOf("namo").ToDec()
		currentPercentage := currentValue.Quo(totalNamo)
		targetPercentage := pool.Allocation.TargetPercentage
		
		deviation := currentPercentage.Sub(targetPercentage)
		deviations[pool.PoolID] = deviation
	}
	
	return deviations
}

// calculateTotalTreasuryValue calculates total value across all pools
func (re *RebalanceEngine) calculateTotalTreasuryValue(pools []TreasuryPool) sdk.Coins {
	var total sdk.Coins
	for _, pool := range pools {
		total = total.Add(pool.Balance...)
	}
	return total
}

// findPoolByID finds a pool by its ID
func (re *RebalanceEngine) findPoolByID(pools []TreasuryPool, poolID string) *TreasuryPool {
	for i := range pools {
		if pools[i].PoolID == poolID {
			return &pools[i]
		}
	}
	return nil
}

// calculateActionPriority calculates priority for rebalance action
func (re *RebalanceEngine) calculateActionPriority(sourcePool, destPool *TreasuryPool) int {
	priority := 100 // Base priority
	
	// Higher priority for critical pools
	if sourcePool.PoolType == "RESERVE" || destPool.PoolType == "RESERVE" {
		priority -= 50
	}
	
	if sourcePool.PoolType == "OPERATIONAL" || destPool.PoolType == "OPERATIONAL" {
		priority -= 30
	}
	
	if sourcePool.PoolType == "SECURITY" || destPool.PoolType == "SECURITY" {
		priority -= 20
	}
	
	return priority
}

// assessActionRisk assesses risk level for rebalance action
func (re *RebalanceEngine) assessActionRisk(sourcePool, destPool *TreasuryPool, amount sdk.Dec) string {
	// Calculate percentage of source pool being transferred
	sourcePercentage := amount.Quo(sourcePool.Balance.AmountOf("namo").ToDec())
	
	if sourcePercentage.GTE(sdk.NewDecWithPrec(50, 2)) { // >= 50%
		return "HIGH"
	} else if sourcePercentage.GTE(sdk.NewDecWithPrec(25, 2)) { // >= 25%
		return "MEDIUM"
	} else {
		return "LOW"
	}
}

// generateEmergencyActions generates emergency rebalancing actions
func (re *RebalanceEngine) generateEmergencyActions(ctx sdk.Context, pools []TreasuryPool, totalValue sdk.Coins) []RebalanceAction {
	var actions []RebalanceAction
	
	// Check for pools below minimum balance
	for _, pool := range pools {
		if pool.MinBalance.IsAllPositive() && !pool.Balance.IsAllGTE(pool.MinBalance) {
			// Generate emergency funding action
			needed := pool.MinBalance.Sub(pool.Balance)
			
			// Find source pool (prefer reserve pool for emergencies)
			var sourcePool *TreasuryPool
			for i := range pools {
				if pools[i].PoolType == "RESERVE" && pools[i].Balance.IsAllGTE(needed) {
					sourcePool = &pools[i]
					break
				}
			}
			
			if sourcePool != nil {
				action := RebalanceAction{
					ActionID:        fmt.Sprintf("EMERGENCY_%d_%s", ctx.BlockHeight(), pool.PoolID),
					ActionType:      "TRANSFER",
					SourcePool:      sourcePool.PoolID,
					DestinationPool: pool.PoolID,
					Amount:          needed,
					Priority:        1, // Highest priority
					Reason:          fmt.Sprintf("Emergency funding for %s below minimum balance", pool.PoolID),
					ExpectedImpact:  needed.AmountOf("namo").ToDec().Quo(totalValue.AmountOf("namo").ToDec()),
					RiskLevel:       "HIGH",
					Prerequisites:   []string{},
					Status:          "PENDING",
				}
				
				actions = append(actions, action)
			}
		}
	}
	
	return actions
}

// Helper utility functions
func (re *RebalanceEngine) generateRebalanceID(ctx sdk.Context) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("REBAL_%d_%d", ctx.BlockHeight(), timestamp)
}

func (re *RebalanceEngine) generateTransactionID(ctx sdk.Context, actionID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("TX_%s_%d", actionID, timestamp)
}

// Additional helper methods would include:
// - validateRebalancePlan
// - checkActionPrerequisites
// - executeMintAction
// - executeBurnAction
// - executeSwapAction
// - calculateExpectedOutcome
// - assessRebalanceRisks
// - determineExecutionStrategy
// - createExecutionTimeline