package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/treasury/types"
)

// TreasuryManager handles comprehensive multi-pool treasury management
type TreasuryManager struct {
	keeper Keeper
}

// NewTreasuryManager creates a new treasury manager
func NewTreasuryManager(keeper Keeper) *TreasuryManager {
	return &TreasuryManager{
		keeper: keeper,
	}
}

// TreasuryPool represents a specific treasury pool with distinct purposes
type TreasuryPool struct {
	PoolID          string                    `json:"pool_id"`
	PoolName        string                    `json:"pool_name"`
	PoolType        string                    `json:"pool_type"` // OPERATIONAL, DEVELOPMENT, RESERVE, CHARITY, SECURITY
	Purpose         string                    `json:"purpose"`
	Allocation      types.PoolAllocation      `json:"allocation"`
	Balance         sdk.Coins                 `json:"balance"`
	ReserveRatio    sdk.Dec                   `json:"reserve_ratio"`
	TargetBalance   sdk.Coins                 `json:"target_balance"`
	MinBalance      sdk.Coins                 `json:"min_balance"`
	MaxBalance      sdk.Coins                 `json:"max_balance"`
	RebalanceConfig types.RebalanceConfig     `json:"rebalance_config"`
	AccessControl   types.AccessControl       `json:"access_control"`
	Performance     types.PoolPerformance     `json:"performance"`
	LastRebalance   time.Time                 `json:"last_rebalance"`
	Status          string                    `json:"status"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// InitializeTreasuryPools creates and initializes all treasury pools
func (tm *TreasuryManager) InitializeTreasuryPools(ctx sdk.Context) error {
	pools := []TreasuryPool{
		tm.createOperationalPool(ctx),
		tm.createDevelopmentPool(ctx),
		tm.createReservePool(ctx),
		tm.createCharityPool(ctx),
		tm.createSecurityPool(ctx),
		tm.createFounderPool(ctx),
		tm.createLiquidityPool(ctx),
		tm.createIncentivePool(ctx),
	}

	for _, pool := range pools {
		err := tm.keeper.SetTreasuryPool(ctx, pool)
		if err != nil {
			return fmt.Errorf("failed to initialize pool %s: %w", pool.PoolID, err)
		}
	}

	// Emit treasury initialization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryInitialized,
			sdk.NewAttribute(types.AttributeKeyPoolCount, fmt.Sprintf("%d", len(pools))),
		),
	)

	return nil
}

// ProcessTreasuryRevenue processes incoming revenue and distributes to pools
func (tm *TreasuryManager) ProcessTreasuryRevenue(ctx sdk.Context, revenue sdk.Coins, source string) error {
	// Get revenue distribution parameters
	params := tm.keeper.GetParams(ctx)
	
	// Calculate allocations based on revenue source and current pool states
	allocation := tm.calculateRevenueAllocation(ctx, revenue, source, params)

	// Distribute revenue to pools
	for poolID, amount := range allocation {
		err := tm.allocateToPool(ctx, poolID, amount, fmt.Sprintf("Revenue from %s", source))
		if err != nil {
			return fmt.Errorf("failed to allocate to pool %s: %w", poolID, err)
		}
	}

	// Record revenue transaction
	tm.recordRevenueTransaction(ctx, revenue, source, allocation)

	// Check if rebalancing is needed
	tm.scheduleRebalanceCheck(ctx)

	// Emit revenue processing event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryRevenueProcessed,
			sdk.NewAttribute(types.AttributeKeyRevenue, revenue.String()),
			sdk.NewAttribute(types.AttributeKeySource, source),
		),
	)

	return nil
}

// ExecuteAutomaticRebalancing performs automated treasury rebalancing
func (tm *TreasuryManager) ExecuteAutomaticRebalancing(ctx sdk.Context) error {
	// Get all treasury pools
	pools := tm.keeper.GetAllTreasuryPools(ctx)
	
	// Analyze current pool states
	rebalanceNeeded, rebalancePlan := tm.analyzeRebalanceNeed(ctx, pools)
	
	if !rebalanceNeeded {
		return nil
	}

	// Execute rebalancing plan
	for _, action := range rebalancePlan.Actions {
		err := tm.executeRebalanceAction(ctx, action)
		if err != nil {
			return fmt.Errorf("failed to execute rebalance action: %w", err)
		}
	}

	// Update pool states and performance metrics
	tm.updatePoolPerformanceMetrics(ctx, pools)

	// Record rebalancing transaction
	tm.recordRebalanceTransaction(ctx, rebalancePlan)

	// Emit rebalancing event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryRebalanced,
			sdk.NewAttribute(types.AttributeKeyRebalanceID, rebalancePlan.RebalanceID),
			sdk.NewAttribute(types.AttributeKeyActionsCount, fmt.Sprintf("%d", len(rebalancePlan.Actions))),
		),
	)

	return nil
}

// ProcessPoolWithdrawal handles withdrawals from treasury pools with governance
func (tm *TreasuryManager) ProcessPoolWithdrawal(ctx sdk.Context, request types.WithdrawalRequest) error {
	// Validate withdrawal request
	if err := tm.validateWithdrawalRequest(ctx, request); err != nil {
		return fmt.Errorf("withdrawal validation failed: %w", err)
	}

	// Get pool
	pool, found := tm.keeper.GetTreasuryPool(ctx, request.PoolID)
	if !found {
		return fmt.Errorf("treasury pool not found: %s", request.PoolID)
	}

	// Check access control
	if !tm.checkWithdrawalPermissions(ctx, request, pool) {
		return fmt.Errorf("insufficient permissions for withdrawal from pool %s", request.PoolID)
	}

	// Check pool balance and limits
	if err := tm.validatePoolWithdrawal(ctx, pool, request.Amount); err != nil {
		return fmt.Errorf("pool withdrawal validation failed: %w", err)
	}

	// Execute withdrawal
	err := tm.executePoolWithdrawal(ctx, pool, request)
	if err != nil {
		return fmt.Errorf("failed to execute pool withdrawal: %w", err)
	}

	// Update pool state
	pool.Balance = pool.Balance.Sub(request.Amount)
	pool.UpdatedAt = ctx.BlockTime()
	tm.keeper.SetTreasuryPool(ctx, pool)

	// Record withdrawal transaction
	tm.recordWithdrawalTransaction(ctx, request)

	// Check if rebalancing is needed after withdrawal
	tm.scheduleRebalanceCheck(ctx)

	// Emit withdrawal event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryWithdrawal,
			sdk.NewAttribute(types.AttributeKeyPoolID, request.PoolID),
			sdk.NewAttribute(types.AttributeKeyAmount, request.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, request.Recipient),
		),
	)

	return nil
}

// GenerateTreasuryReport generates comprehensive treasury performance report
func (tm *TreasuryManager) GenerateTreasuryReport(ctx sdk.Context, timeRange types.TimeRange) (*types.TreasuryReport, error) {
	report := &types.TreasuryReport{
		ReportID:    tm.generateReportID(ctx),
		TimeRange:   timeRange,
		GeneratedAt: ctx.BlockTime(),
	}

	// Get all pools
	pools := tm.keeper.GetAllTreasuryPools(ctx)
	report.TotalPools = int64(len(pools))

	// Calculate total treasury value
	var totalValue sdk.Coins
	for _, pool := range pools {
		totalValue = totalValue.Add(pool.Balance...)
	}
	report.TotalTreasuryValue = totalValue

	// Generate pool summaries
	for _, pool := range pools {
		summary := tm.generatePoolSummary(ctx, pool, timeRange)
		report.PoolSummaries = append(report.PoolSummaries, summary)
	}

	// Calculate treasury performance metrics
	report.PerformanceMetrics = tm.calculateTreasuryPerformance(ctx, pools, timeRange)

	// Analyze revenue streams
	report.RevenueAnalysis = tm.analyzeRevenueStreams(ctx, timeRange)

	// Calculate allocation efficiency
	report.AllocationEfficiency = tm.calculateAllocationEfficiency(ctx, pools, timeRange)

	// Generate risk assessment
	report.RiskAssessment = tm.assessTreasuryRisks(ctx, pools)

	// Provide recommendations
	report.Recommendations = tm.generateTreasuryRecommendations(ctx, pools, report)

	// Store report
	tm.keeper.SetTreasuryReport(ctx, *report)

	return report, nil
}

// Helper functions for treasury pool creation

func (tm *TreasuryManager) createOperationalPool(ctx sdk.Context) TreasuryPool {
	params := tm.keeper.GetParams(ctx)
	
	return TreasuryPool{
		PoolID:   "OPERATIONAL",
		PoolName: "Operational Treasury",
		PoolType: "OPERATIONAL",
		Purpose:  "Day-to-day operational expenses, validator rewards, and network maintenance",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(25, 2), // 25%
			MinPercentage:    sdk.NewDecWithPrec(20, 2), // 20%
			MaxPercentage:    sdk.NewDecWithPrec(35, 2), // 35%
		},
		ReserveRatio: sdk.NewDecWithPrec(10, 2), // 10% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(5, 2), // 5% deviation triggers rebalance
			MinRebalanceGap:   24 * time.Hour,           // Minimum 24 hours between rebalances
			MaxRebalanceGap:   7 * 24 * time.Hour,      // Maximum 7 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 2,
			AuthorizedRoles:    []string{"TREASURY_MANAGER", "OPERATIONAL_ADMIN"},
			GovernanceRequired: false,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createDevelopmentPool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "DEVELOPMENT",
		PoolName: "Development Treasury",
		PoolType: "DEVELOPMENT",
		Purpose:  "Platform development, research, innovation, and technical upgrades",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(20, 2), // 20%
			MinPercentage:    sdk.NewDecWithPrec(15, 2), // 15%
			MaxPercentage:    sdk.NewDecWithPrec(30, 2), // 30%
		},
		ReserveRatio: sdk.NewDecWithPrec(15, 2), // 15% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(10, 2), // 10% deviation triggers rebalance
			MinRebalanceGap:   48 * time.Hour,            // Minimum 48 hours between rebalances
			MaxRebalanceGap:   14 * 24 * time.Hour,       // Maximum 14 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 3,
			AuthorizedRoles:    []string{"DEVELOPMENT_LEAD", "TREASURY_MANAGER", "TECHNICAL_COMMITTEE"},
			GovernanceRequired: true,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createReservePool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "RESERVE",
		PoolName: "Strategic Reserve",
		PoolType: "RESERVE",
		Purpose:  "Emergency fund, market stability, and long-term strategic initiatives",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(30, 2), // 30%
			MinPercentage:    sdk.NewDecWithPrec(25, 2), // 25%
			MaxPercentage:    sdk.NewDecWithPrec(40, 2), // 40%
		},
		ReserveRatio: sdk.NewDecWithPrec(95, 2), // 95% reserve (minimal spending)
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(15, 2), // 15% deviation triggers rebalance
			MinRebalanceGap:   7 * 24 * time.Hour,        // Minimum 7 days between rebalances
			MaxRebalanceGap:   30 * 24 * time.Hour,       // Maximum 30 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 5,
			AuthorizedRoles:    []string{"TREASURY_MANAGER", "BOARD_MEMBER", "FOUNDER"},
			GovernanceRequired: true,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createCharityPool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "CHARITY",
		PoolName: "Social Impact Treasury",
		PoolType: "CHARITY",
		Purpose:  "40% of all fees dedicated to charitable and social impact initiatives",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(40, 2), // 40% of fees (not total treasury)
			MinPercentage:    sdk.NewDecWithPrec(35, 2), // 35%
			MaxPercentage:    sdk.NewDecWithPrec(45, 2), // 45%
		},
		ReserveRatio: sdk.NewDecWithPrec(5, 2), // 5% reserve (active distribution)
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(10, 2), // 10% deviation triggers rebalance
			MinRebalanceGap:   24 * time.Hour,            // Minimum 24 hours between rebalances
			MaxRebalanceGap:   7 * 24 * time.Hour,        // Maximum 7 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 2,
			AuthorizedRoles:    []string{"CHARITY_MANAGER", "SOCIAL_IMPACT_LEAD"},
			GovernanceRequired: false,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createSecurityPool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "SECURITY",
		PoolName: "Security Treasury",
		PoolType: "SECURITY",
		Purpose:  "Bug bounties, security audits, incident response, and validator slashing insurance",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(5, 2),  // 5%
			MinPercentage:    sdk.NewDecWithPrec(3, 2),  // 3%
			MaxPercentage:    sdk.NewDecWithPrec(10, 2), // 10%
		},
		ReserveRatio: sdk.NewDecWithPrec(80, 2), // 80% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(20, 2), // 20% deviation triggers rebalance
			MinRebalanceGap:   24 * time.Hour,            // Minimum 24 hours between rebalances
			MaxRebalanceGap:   30 * 24 * time.Hour,       // Maximum 30 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 3,
			AuthorizedRoles:    []string{"SECURITY_LEAD", "TREASURY_MANAGER", "TECHNICAL_LEAD"},
			GovernanceRequired: true,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createFounderPool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "FOUNDER",
		PoolName: "Founder Allocation",
		PoolType: "FOUNDER",
		Purpose:  "10% allocation for founder with perpetual 0.10% tax royalty and 5% platform royalty",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(10, 2), // 10% (immutable)
			MinPercentage:    sdk.NewDecWithPrec(10, 2), // 10%
			MaxPercentage:    sdk.NewDecWithPrec(10, 2), // 10%
		},
		ReserveRatio: sdk.NewDecWithPrec(50, 2), // 50% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           false, // No automatic rebalancing for founder pool
			ThresholdPercent:  sdk.ZeroDec(),
			MinRebalanceGap:   0,
			MaxRebalanceGap:   0,
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 1,
			AuthorizedRoles:    []string{"FOUNDER"},
			GovernanceRequired: false,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createLiquidityPool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "LIQUIDITY",
		PoolName: "Liquidity Support",
		PoolType: "LIQUIDITY",
		Purpose:  "20% allocation for market making, DEX liquidity, and price stability",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(20, 2), // 20%
			MinPercentage:    sdk.NewDecWithPrec(15, 2), // 15%
			MaxPercentage:    sdk.NewDecWithPrec(25, 2), // 25%
		},
		ReserveRatio: sdk.NewDecWithPrec(20, 2), // 20% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(5, 2), // 5% deviation triggers rebalance
			MinRebalanceGap:   12 * time.Hour,           // Minimum 12 hours between rebalances
			MaxRebalanceGap:   7 * 24 * time.Hour,      // Maximum 7 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 2,
			AuthorizedRoles:    []string{"LIQUIDITY_MANAGER", "TREASURY_MANAGER"},
			GovernanceRequired: false,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

func (tm *TreasuryManager) createIncentivePool(ctx sdk.Context) TreasuryPool {
	return TreasuryPool{
		PoolID:   "INCENTIVE",
		PoolName: "Community Incentives",
		PoolType: "INCENTIVE",
		Purpose:  "15% allocation for user rewards, staking incentives, and community programs",
		Allocation: types.PoolAllocation{
			TargetPercentage: sdk.NewDecWithPrec(15, 2), // 15%
			MinPercentage:    sdk.NewDecWithPrec(10, 2), // 10%
			MaxPercentage:    sdk.NewDecWithPrec(20, 2), // 20%
		},
		ReserveRatio: sdk.NewDecWithPrec(10, 2), // 10% reserve
		RebalanceConfig: types.RebalanceConfig{
			Enabled:           true,
			ThresholdPercent:  sdk.NewDecWithPrec(8, 2), // 8% deviation triggers rebalance
			MinRebalanceGap:   24 * time.Hour,           // Minimum 24 hours between rebalances
			MaxRebalanceGap:   14 * 24 * time.Hour,      // Maximum 14 days without rebalance
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 2,
			AuthorizedRoles:    []string{"COMMUNITY_MANAGER", "TREASURY_MANAGER"},
			GovernanceRequired: false,
		},
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
}

// Helper utility functions
func (tm *TreasuryManager) generateReportID(ctx sdk.Context) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("TREASURY_REPORT_%d", timestamp)
}