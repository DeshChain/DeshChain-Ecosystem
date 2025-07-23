package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/treasury/types"
)

// DUSDReserveManager manages DUSD reserves and USD collateral within the treasury system
type DUSDReserveManager struct {
	treasuryManager *TreasuryManager
	keeper          *Keeper
}

// NewDUSDReserveManager creates a new DUSD reserve manager
func NewDUSDReserveManager(treasuryManager *TreasuryManager, keeper *Keeper) *DUSDReserveManager {
	return &DUSDReserveManager{
		treasuryManager: treasuryManager,
		keeper:          keeper,
	}
}

// Enhanced Treasury Pool for DUSD Operations
type DUSDTreasuryPool struct {
	// Base treasury pool
	TreasuryPool
	
	// DUSD-specific fields
	USDCollateralRatio    sdk.Dec                    `json:"usd_collateral_ratio"`    // Target USD collateral ratio
	DUSDSupplyBacked      sdk.Coin                   `json:"dusd_supply_backed"`      // DUSD supply backed by this pool
	USDReserveAssets      []types.ReserveAsset       `json:"usd_reserve_assets"`      // USD-denominated assets
	StabilityBuffer       sdk.Coin                   `json:"stability_buffer"`        // Buffer for stability operations
	CrossCurrencyExposure map[string]sdk.Dec         `json:"cross_currency_exposure"` // Exposure to other currencies
	RebalanceStrategy     types.USDRebalanceStrategy `json:"rebalance_strategy"`      // USD-specific rebalancing
}

// InitializeDUSDReservePools creates DUSD-specific treasury pools
func (drm *DUSDReserveManager) InitializeDUSDReservePools(ctx sdk.Context) error {
	// Create DUSD reserve pool as 9th treasury pool
	dusdPool := drm.createDUSDReservePool(ctx)
	
	// Create USD operations pool for collateral management
	usdOpsPool := drm.createUSDOperationsPool(ctx)
	
	// Create cross-currency pool for multi-currency operations
	crossCurrencyPool := drm.createCrossCurrencyPool(ctx)
	
	pools := []*DUSDTreasuryPool{dusdPool, usdOpsPool, crossCurrencyPool}
	
	for _, pool := range pools {
		if err := drm.storeDUSDPool(ctx, pool); err != nil {
			return fmt.Errorf("failed to store DUSD pool %s: %w", pool.PoolID, err)
		}
	}
	
	// Initialize USD collateral reserves
	if err := drm.initializeUSDCollateral(ctx); err != nil {
		return fmt.Errorf("failed to initialize USD collateral: %w", err)
	}
	
	drm.keeper.Logger(ctx).Info("initialized DUSD treasury pools",
		"pool_count", len(pools),
		"total_usd_backing", "initial",
	)
	
	return nil
}

// createDUSDReservePool creates the main DUSD reserve pool
func (drm *DUSDReserveManager) createDUSDReservePool(ctx sdk.Context) *DUSDTreasuryPool {
	basePool := TreasuryPool{
		PoolID:          "dusd_reserve",
		PoolName:        "DUSD Reserve Pool",
		PoolType:        "DUSD_RESERVE",
		Purpose:         "Backing DUSD stablecoin with USD reserves and collateral",
		Allocation: types.PoolAllocation{
			PercentageAllocation: sdk.NewDecWithPrec(15, 2), // 15% of total treasury
			FixedAllocation:     sdk.NewCoins(),
			Priority:           2, // High priority
		},
		ReserveRatio: sdk.NewDecWithPrec(150, 2), // 150% collateral ratio
		RebalanceConfig: types.RebalanceConfig{
			Enabled:             true,
			DeviationThreshold:  sdk.NewDecWithPrec(5, 2),  // 5% deviation trigger
			RebalanceFrequency:  24 * time.Hour,            // Daily rebalancing
			MinRebalanceAmount:  sdk.NewCoins(sdk.NewCoin("NAMO", sdk.NewInt(10000))),
			MaxRebalanceAmount:  sdk.NewCoins(sdk.NewCoin("NAMO", sdk.NewInt(10000000))),
			AutoRebalanceEnabled: true,
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 3,
			AuthorizedRoles:   []string{"treasury_manager", "dusd_operator", "governor"},
			EmergencyAccess:   true,
		},
		Status:    "active",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
	
	return &DUSDTreasuryPool{
		TreasuryPool:          basePool,
		USDCollateralRatio:    sdk.NewDecWithPrec(150, 2), // 150%
		DUSDSupplyBacked:      sdk.NewCoin("DUSD", sdk.ZeroInt()),
		USDReserveAssets:      drm.getInitialUSDAssets(),
		StabilityBuffer:       sdk.NewCoin("USD", sdk.NewInt(1000000)), // $1M buffer
		CrossCurrencyExposure: make(map[string]sdk.Dec),
		RebalanceStrategy:     drm.getDefaultUSDRebalanceStrategy(),
	}
}

// createUSDOperationsPool creates pool for USD operational expenses
func (drm *DUSDReserveManager) createUSDOperationsPool(ctx sdk.Context) *DUSDTreasuryPool {
	basePool := TreasuryPool{
		PoolID:          "usd_operations",
		PoolName:        "USD Operations Pool",
		PoolType:        "USD_OPERATIONS",
		Purpose:         "USD-denominated operational expenses and multi-currency operations",
		Allocation: types.PoolAllocation{
			PercentageAllocation: sdk.NewDecWithPrec(5, 2), // 5% of total treasury
			Priority:           3,
		},
		ReserveRatio: sdk.NewDecWithPrec(25, 2), // 25% reserve ratio
		RebalanceConfig: types.RebalanceConfig{
			Enabled:             true,
			DeviationThreshold:  sdk.NewDecWithPrec(10, 2), // 10% deviation trigger
			RebalanceFrequency:  7 * 24 * time.Hour,        // Weekly rebalancing
			AutoRebalanceEnabled: true,
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 2,
			AuthorizedRoles:   []string{"treasury_manager", "operations_manager"},
		},
		Status:    "active",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
	
	return &DUSDTreasuryPool{
		TreasuryPool:          basePool,
		USDCollateralRatio:    sdk.NewDecWithPrec(100, 2), // 100% (operational, not backing)
		USDReserveAssets:      []types.ReserveAsset{},
		StabilityBuffer:       sdk.NewCoin("USD", sdk.NewInt(100000)), // $100K buffer
		CrossCurrencyExposure: make(map[string]sdk.Dec),
		RebalanceStrategy:     drm.getOperationalUSDStrategy(),
	}
}

// createCrossCurrencyPool creates pool for cross-currency operations
func (drm *DUSDReserveManager) createCrossCurrencyPool(ctx sdk.Context) *DUSDTreasuryPool {
	basePool := TreasuryPool{
		PoolID:          "cross_currency",
		PoolName:        "Cross-Currency Pool",
		PoolType:        "CROSS_CURRENCY",
		Purpose:         "Managing exposure across multiple stablecoin currencies (DUSD, DEUR, DSGD)",
		Allocation: types.PoolAllocation{
			PercentageAllocation: sdk.NewDecWithPrec(10, 2), // 10% of total treasury
			Priority:           3,
		},
		ReserveRatio: sdk.NewDecWithPrec(120, 2), // 120% reserve ratio
		RebalanceConfig: types.RebalanceConfig{
			Enabled:             true,
			DeviationThreshold:  sdk.NewDecWithPrec(3, 2),   // 3% deviation trigger (more sensitive)
			RebalanceFrequency:  6 * time.Hour,              // 4x daily rebalancing
			AutoRebalanceEnabled: true,
		},
		AccessControl: types.AccessControl{
			RequiredSignatures: 3,
			AuthorizedRoles:   []string{"treasury_manager", "forex_manager", "risk_manager"},
		},
		Status:    "active",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}
	
	return &DUSDTreasuryPool{
		TreasuryPool:       basePool,
		USDCollateralRatio: sdk.NewDecWithPrec(120, 2), // 120%
		USDReserveAssets:   []types.ReserveAsset{},
		StabilityBuffer:    sdk.NewCoin("USD", sdk.NewInt(500000)), // $500K buffer
		CrossCurrencyExposure: map[string]sdk.Dec{
			"USD": sdk.NewDecWithPrec(40, 2), // 40%
			"EUR": sdk.NewDecWithPrec(30, 2), // 30%
			"SGD": sdk.NewDecWithPrec(20, 2), // 20%
			"GBP": sdk.NewDecWithPrec(10, 2), // 10%
		},
		RebalanceStrategy: drm.getCrossCurrencyStrategy(),
	}
}

// ManageDUSDCollateral handles DUSD collateral requirements
func (drm *DUSDReserveManager) ManageDUSDCollateral(ctx sdk.Context, dusdMintedAmount sdk.Coin) error {
	// Get DUSD reserve pool
	dusdPool, err := drm.getDUSDReservePool(ctx)
	if err != nil {
		return fmt.Errorf("failed to get DUSD reserve pool: %w", err)
	}
	
	// Calculate required USD collateral (150% ratio)
	requiredUSDCollateral := dusdMintedAmount.Amount.ToLegacyDec().Mul(dusdPool.USDCollateralRatio)
	
	// Check current USD reserves
	currentReserves := drm.calculateTotalUSDReserves(ctx, dusdPool)
	
	// If insufficient reserves, trigger rebalancing
	if currentReserves.LT(requiredUSDCollateral) {
		shortfall := requiredUSDCollateral.Sub(currentReserves)
		if err := drm.acquireAdditionalUSDReserves(ctx, shortfall); err != nil {
			return fmt.Errorf("failed to acquire additional USD reserves: %w", err)
		}
	}
	
	// Update pool backing statistics
	dusdPool.DUSDSupplyBacked = dusdPool.DUSDSupplyBacked.Add(dusdMintedAmount)
	dusdPool.UpdatedAt = ctx.BlockTime()
	
	// Store updated pool
	if err := drm.storeDUSDPool(ctx, dusdPool); err != nil {
		return fmt.Errorf("failed to update DUSD pool: %w", err)
	}
	
	// Emit collateral management event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"dusd_collateral_managed",
			sdk.NewAttribute("dusd_minted", dusdMintedAmount.String()),
			sdk.NewAttribute("required_collateral", requiredUSDCollateral.String()),
			sdk.NewAttribute("current_reserves", currentReserves.String()),
		),
	)
	
	return nil
}

// RebalanceCrossCurrencyExposure rebalances exposure across multiple currencies
func (drm *DUSDReserveManager) RebalanceCrossCurrencyExposure(ctx sdk.Context) error {
	// Get cross-currency pool
	crossPool, err := drm.getCrossCurrencyPool(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cross-currency pool: %w", err)
	}
	
	// Calculate current exposure vs target
	currentExposure := drm.calculateCurrentCurrencyExposure(ctx)
	targetExposure := crossPool.CrossCurrencyExposure
	
	rebalanceActions := []types.RebalanceAction{}
	
	for currency, targetPct := range targetExposure {
		currentPct := currentExposure[currency]
		deviation := currentPct.Sub(targetPct).Abs()
		
		// If deviation exceeds threshold, create rebalancing action
		if deviation.GT(crossPool.RebalanceConfig.DeviationThreshold) {
			action := types.RebalanceAction{
				Currency:        currency,
				CurrentExposure: currentPct,
				TargetExposure:  targetPct,
				Deviation:       deviation,
				Action:          drm.determineRebalanceAction(currentPct, targetPct),
				Timestamp:       ctx.BlockTime(),
			}
			rebalanceActions = append(rebalanceActions, action)
		}
	}
	
	// Execute rebalancing actions
	for _, action := range rebalanceActions {
		if err := drm.executeRebalanceAction(ctx, action); err != nil {
			drm.keeper.Logger(ctx).Error("failed to execute rebalance action",
				"currency", action.Currency,
				"action", action.Action,
				"error", err,
			)
			continue
		}
	}
	
	// Update last rebalance timestamp
	crossPool.LastRebalance = ctx.BlockTime()
	if err := drm.storeDUSDPool(ctx, crossPool); err != nil {
		return fmt.Errorf("failed to update cross-currency pool: %w", err)
	}
	
	drm.keeper.Logger(ctx).Info("completed cross-currency rebalancing",
		"actions_executed", len(rebalanceActions),
	)
	
	return nil
}

// calculateTotalUSDReserves calculates total USD reserves across all assets
func (drm *DUSDReserveManager) calculateTotalUSDReserves(ctx sdk.Context, pool *DUSDTreasuryPool) sdk.Dec {
	total := sdk.ZeroDec()
	
	for _, asset := range pool.USDReserveAssets {
		// Convert asset value to USD if needed
		assetValue := asset.Value.ToLegacyDec()
		if asset.Currency != "USD" {
			// Get exchange rate to USD
			if rate, err := drm.getUSDExchangeRate(ctx, asset.Currency); err == nil {
				assetValue = assetValue.Mul(rate)
			}
		}
		total = total.Add(assetValue)
	}
	
	return total
}

// Helper functions for initialization and operations

func (drm *DUSDReserveManager) getInitialUSDAssets() []types.ReserveAsset {
	return []types.ReserveAsset{
		{
			AssetID:    "usdc_reserves",
			AssetType:  "STABLECOIN",
			Currency:   "USDC",
			Value:      sdk.NewCoin("USDC", sdk.NewInt(5000000)), // $5M USDC
			Liquidity:  sdk.NewDecWithPrec(95, 2),               // 95% liquid
			Risk:       sdk.NewDecWithPrec(5, 2),                // 5% risk
			Yield:      sdk.NewDecWithPrec(3, 2),                // 3% yield
			MaturityDate: nil, // No maturity for stablecoins
		},
		{
			AssetID:    "treasury_bills",
			AssetType:  "GOVERNMENT_BOND",
			Currency:   "USD",
			Value:      sdk.NewCoin("USD", sdk.NewInt(10000000)), // $10M Treasury Bills
			Liquidity:  sdk.NewDecWithPrec(80, 2),               // 80% liquid
			Risk:       sdk.NewDecWithPrec(2, 2),                // 2% risk
			Yield:      sdk.NewDecWithPrec(45, 3),               // 4.5% yield
			MaturityDate: &types.MaturityDate{
				Date: time.Now().AddDate(0, 6, 0), // 6 months
			},
		},
	}
}

func (drm *DUSDReserveManager) getDefaultUSDRebalanceStrategy() types.USDRebalanceStrategy {
	return types.USDRebalanceStrategy{
		Strategy:           "CONSERVATIVE",
		MaxDeviation:       sdk.NewDecWithPrec(5, 2),   // 5%
		RebalanceFrequency: 24 * time.Hour,             // Daily
		AssetWeights: map[string]sdk.Dec{
			"USDC":          sdk.NewDecWithPrec(30, 2), // 30%
			"USDT":          sdk.NewDecWithPrec(20, 2), // 20%
			"Treasury_Bills": sdk.NewDecWithPrec(40, 2), // 40%
			"Cash":          sdk.NewDecWithPrec(10, 2), // 10%
		},
		RiskLimits: map[string]sdk.Dec{
			"single_asset": sdk.NewDecWithPrec(50, 2), // Max 50% in single asset
			"stablecoin":   sdk.NewDecWithPrec(60, 2), // Max 60% in stablecoins
			"bonds":        sdk.NewDecWithPrec(70, 2), // Max 70% in bonds
		},
	}
}

func (drm *DUSDReserveManager) getOperationalUSDStrategy() types.USDRebalanceStrategy {
	return types.USDRebalanceStrategy{
		Strategy:           "OPERATIONAL",
		MaxDeviation:       sdk.NewDecWithPrec(15, 2), // 15% (more flexible)
		RebalanceFrequency: 7 * 24 * time.Hour,        // Weekly
		AssetWeights: map[string]sdk.Dec{
			"Cash": sdk.NewDecWithPrec(80, 2),          // 80% cash for operations
			"USDC": sdk.NewDecWithPrec(20, 2),          // 20% USDC
		},
		RiskLimits: map[string]sdk.Dec{
			"cash_min": sdk.NewDecWithPrec(50, 2), // Min 50% cash
		},
	}
}

func (drm *DUSDReserveManager) getCrossCurrencyStrategy() types.USDRebalanceStrategy {
	return types.USDRebalanceStrategy{
		Strategy:           "CROSS_CURRENCY",
		MaxDeviation:       sdk.NewDecWithPrec(3, 2),  // 3% (very sensitive)
		RebalanceFrequency: 6 * time.Hour,             // 4x daily
		AssetWeights: map[string]sdk.Dec{
			"USD": sdk.NewDecWithPrec(40, 2), // 40%
			"EUR": sdk.NewDecWithPrec(30, 2), // 30%
			"SGD": sdk.NewDecWithPrec(20, 2), // 20%
			"GBP": sdk.NewDecWithPrec(10, 2), // 10%
		},
		RiskLimits: map[string]sdk.Dec{
			"single_currency": sdk.NewDecWithPrec(50, 2), // Max 50% in single currency
			"currency_pair":   sdk.NewDecWithPrec(25, 2), // Max 25% correlation risk
		},
	}
}

// Storage and retrieval functions

func (drm *DUSDReserveManager) storeDUSDPool(ctx sdk.Context, pool *DUSDTreasuryPool) error {
	store := ctx.KVStore(drm.keeper.storeKey)
	key := types.GetDUSDPoolKey(pool.PoolID)
	
	bz, err := drm.keeper.cdc.Marshal(pool)
	if err != nil {
		return err
	}
	
	store.Set(key, bz)
	return nil
}

func (drm *DUSDReserveManager) getDUSDReservePool(ctx sdk.Context) (*DUSDTreasuryPool, error) {
	store := ctx.KVStore(drm.keeper.storeKey)
	key := types.GetDUSDPoolKey("dusd_reserve")
	
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("DUSD reserve pool not found")
	}
	
	var pool DUSDTreasuryPool
	if err := drm.keeper.cdc.Unmarshal(bz, &pool); err != nil {
		return nil, err
	}
	
	return &pool, nil
}

func (drm *DUSDReserveManager) getCrossCurrencyPool(ctx sdk.Context) (*DUSDTreasuryPool, error) {
	store := ctx.KVStore(drm.keeper.storeKey)
	key := types.GetDUSDPoolKey("cross_currency")
	
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("cross-currency pool not found")
	}
	
	var pool DUSDTreasuryPool
	if err := drm.keeper.cdc.Unmarshal(bz, &pool); err != nil {
		return nil, err
	}
	
	return &pool, nil
}

// Placeholder functions for complex operations

func (drm *DUSDReserveManager) initializeUSDCollateral(ctx sdk.Context) error {
	// Initialize USD collateral reserves
	drm.keeper.Logger(ctx).Info("initializing USD collateral reserves")
	return nil
}

func (drm *DUSDReserveManager) acquireAdditionalUSDReserves(ctx sdk.Context, amount sdk.Dec) error {
	// Acquire additional USD reserves through various mechanisms
	drm.keeper.Logger(ctx).Info("acquiring additional USD reserves", "amount", amount.String())
	return nil
}

func (drm *DUSDReserveManager) calculateCurrentCurrencyExposure(ctx sdk.Context) map[string]sdk.Dec {
	// Calculate current currency exposure across all treasury pools
	return map[string]sdk.Dec{
		"USD": sdk.NewDecWithPrec(42, 2), // 42%
		"EUR": sdk.NewDecWithPrec(28, 2), // 28%
		"SGD": sdk.NewDecWithPrec(22, 2), // 22%
		"GBP": sdk.NewDecWithPrec(8, 2),  // 8%
	}
}

func (drm *DUSDReserveManager) determineRebalanceAction(current, target sdk.Dec) string {
	if current.GT(target) {
		return "REDUCE"
	}
	return "INCREASE"
}

func (drm *DUSDReserveManager) executeRebalanceAction(ctx sdk.Context, action types.RebalanceAction) error {
	// Execute the rebalancing action (buy/sell currencies, adjust positions)
	drm.keeper.Logger(ctx).Info("executing rebalance action",
		"currency", action.Currency,
		"action", action.Action,
		"deviation", action.Deviation.String(),
	)
	return nil
}

func (drm *DUSDReserveManager) getUSDExchangeRate(ctx sdk.Context, currency string) (sdk.Dec, error) {
	// Get USD exchange rate from oracle (placeholder)
	rates := map[string]sdk.Dec{
		"EUR": sdk.NewDecWithPrec(92, 2),  // 0.92 USD per EUR
		"SGD": sdk.NewDecWithPrec(74, 2),  // 0.74 USD per SGD
		"GBP": sdk.NewDecWithPrec(127, 2), // 1.27 USD per GBP
		"INR": sdk.NewDecWithPrec(12, 3),  // 0.012 USD per INR
	}
	
	if rate, found := rates[currency]; found {
		return rate, nil
	}
	
	return sdk.ZeroDec(), fmt.Errorf("exchange rate not found for %s", currency)
}