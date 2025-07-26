package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/DeshChain/DeshChain-Ecosystem/x/revenue/types"
)

// Keeper of the revenue store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	
	// Module keepers for revenue tracking
	taxKeeper        types.TaxKeeper
	dinrKeeper       types.DINRKeeper
	dusdKeeper       types.DUSDKeeper
	lendingKeeper    types.LendingKeeper
	tradeKeeper      types.TradeKeeper
	remittanceKeeper types.RemittanceKeeper
	governanceKeeper types.GovernanceKeeper
}

// NewKeeper creates a new revenue Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetModuleKeepers sets the module keepers for revenue tracking
func (k *Keeper) SetModuleKeepers(
	taxKeeper types.TaxKeeper,
	dinrKeeper types.DINRKeeper,
	dusdKeeper types.DUSDKeeper,
	lendingKeeper types.LendingKeeper,
	tradeKeeper types.TradeKeeper,
	remittanceKeeper types.RemittanceKeeper,
	governanceKeeper types.GovernanceKeeper,
) {
	k.taxKeeper = taxKeeper
	k.dinrKeeper = dinrKeeper
	k.dusdKeeper = dusdKeeper
	k.lendingKeeper = lendingKeeper
	k.tradeKeeper = tradeKeeper
	k.remittanceKeeper = remittanceKeeper
	k.governanceKeeper = governanceKeeper
}

// RecordRevenueStream records a new revenue stream
func (k Keeper) RecordRevenueStream(ctx sdk.Context, stream types.RevenueStream) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.RevenueStreamKey, []byte(stream.StreamID)...)
	bz := k.cdc.MustMarshal(&stream)
	store.Set(key, bz)
	
	// Update module revenue
	k.updateModuleRevenue(ctx, stream.ModuleName, stream.Amount)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRevenueRecorded,
			sdk.NewAttribute(types.AttributeKeyStreamID, stream.StreamID),
			sdk.NewAttribute(types.AttributeKeyModule, stream.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAmount, stream.Amount.String()),
		),
	)
}

// updateModuleRevenue updates revenue tracking for a module
func (k Keeper) updateModuleRevenue(ctx sdk.Context, moduleName string, amount sdk.Coins) {
	// Get current period (daily)
	period := "daily"
	startTime := ctx.BlockTime().Truncate(24 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)
	
	// Generate key
	key := append(types.ModuleRevenueKey, []byte(fmt.Sprintf("%s-%s-%d", moduleName, period, startTime.Unix()))...)
	store := ctx.KVStore(k.storeKey)
	
	// Get existing or create new
	var moduleRev types.ModuleRevenue
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &moduleRev)
		moduleRev.Revenue = moduleRev.Revenue.Add(amount...)
		moduleRev.TransactionCount++
	} else {
		moduleRev = types.ModuleRevenue{
			ModuleName:       moduleName,
			Period:           period,
			Revenue:          amount,
			TransactionCount: 1,
			UniqueUsers:      0, // Would track unique users separately
			StartTime:        startTime,
			EndTime:          endTime,
		}
	}
	
	// Save updated revenue
	store.Set(key, k.cdc.MustMarshal(&moduleRev))
}

// GetCurrentPerformanceMetrics retrieves current platform performance metrics
func (k Keeper) GetCurrentPerformanceMetrics(ctx sdk.Context) types.PerformanceMetrics {
	// Aggregate metrics from various modules
	metrics := types.PerformanceMetrics{
		Timestamp: ctx.BlockTime(),
	}
	
	// Get platform revenue and expenses
	stats := k.GetPlatformStatistics(ctx)
	metrics.PlatformRevenue = sdk.NewDecFromInt(stats.TotalRevenue.AmountOf("namo"))
	metrics.PlatformExpenses = k.calculatePlatformExpenses(ctx)
	
	// Get trading volume
	metrics.TradingVolume = k.getTradingVolume(ctx)
	metrics.PreviousTradingVolume = k.getPreviousTradingVolume(ctx)
	
	// Get lending metrics
	if k.lendingKeeper != nil {
		metrics.LendingVolume = k.lendingKeeper.GetTotalLendingVolume(ctx)
		metrics.DefaultRate = k.lendingKeeper.GetDefaultRate(ctx)
	}
	
	// Get DUSD metrics
	if k.dusdKeeper != nil {
		metrics.DUSDRevenue = k.dusdKeeper.GetRevenue(ctx)
		metrics.DUSDVolume = k.dusdKeeper.GetVolume(ctx)
	}
	
	// Get user and transaction counts
	metrics.ActiveUsers = k.getActiveUserCount(ctx)
	metrics.TransactionCount = k.getTransactionCount(ctx)
	
	return metrics
}

// CalculateAndDistributeRevenue calculates and distributes revenue
func (k Keeper) CalculateAndDistributeRevenue(ctx sdk.Context) error {
	// Get total revenue for the period
	totalRevenue := k.getTotalRevenue(ctx)
	if totalRevenue.IsZero() {
		return nil
	}
	
	// Get charity percentage from governance
	charityPercent := k.governanceKeeper.GetCurrentCharityPercentage(ctx)
	
	// Calculate distributions
	charityAmount := sdk.Coins{}
	for _, coin := range totalRevenue {
		charityAmt := sdk.NewDecFromInt(coin.Amount).Mul(charityPercent).TruncateInt()
		if charityAmt.IsPositive() {
			charityAmount = charityAmount.Add(sdk.NewCoin(coin.Denom, charityAmt))
		}
	}
	
	// Calculate operations amount (remaining after charity)
	operationsAmount := totalRevenue.Sub(charityAmount)
	
	// Record distribution
	distribution := types.RevenueDistribution{
		DistributionID:   fmt.Sprintf("dist-%d", ctx.BlockTime().Unix()),
		Timestamp:        ctx.BlockTime(),
		TotalRevenue:     totalRevenue,
		CharityAmount:    charityAmount,
		CharityPercent:   charityPercent,
		OperationsAmount: operationsAmount,
		ReservesAmount:   sdk.Coins{}, // Can allocate to reserves
		YieldAmount:      sdk.Coins{}, // For DINR yield
	}
	
	k.recordRevenueDistribution(ctx, distribution)
	
	return nil
}

// recordRevenueDistribution records a revenue distribution
func (k Keeper) recordRevenueDistribution(ctx sdk.Context, dist types.RevenueDistribution) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.RevenueDistributionKey, []byte(dist.DistributionID)...)
	bz := k.cdc.MustMarshal(&dist)
	store.Set(key, bz)
	
	// Update platform statistics
	k.updatePlatformStatistics(ctx, dist)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRevenueDistributed,
			sdk.NewAttribute(types.AttributeKeyDistributionID, dist.DistributionID),
			sdk.NewAttribute(types.AttributeKeyTotalRevenue, dist.TotalRevenue.String()),
			sdk.NewAttribute(types.AttributeKeyCharityAmount, dist.CharityAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCharityPercent, dist.CharityPercent.String()),
		),
	)
}

// GetPlatformStatistics retrieves platform statistics
func (k Keeper) GetPlatformStatistics(ctx sdk.Context) types.PlatformStatistics {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlatformStatsKey)
	
	if bz == nil {
		return types.PlatformStatistics{
			LastUpdated:             ctx.BlockTime(),
			TotalRevenue:            sdk.Coins{},
			TotalCharityDistributed: sdk.Coins{},
			TotalYieldDistributed:   sdk.Coins{},
			AverageYieldRate:        sdk.ZeroDec(),
			TotalUsers:              0,
			TotalTransactions:       0,
			PlatformUptime:          sdk.NewDecWithPrec(999, 3), // 99.9%
		}
	}
	
	var stats types.PlatformStatistics
	k.cdc.MustUnmarshal(bz, &stats)
	return stats
}

// updatePlatformStatistics updates platform statistics
func (k Keeper) updatePlatformStatistics(ctx sdk.Context, dist types.RevenueDistribution) {
	stats := k.GetPlatformStatistics(ctx)
	
	stats.LastUpdated = ctx.BlockTime()
	stats.TotalRevenue = stats.TotalRevenue.Add(dist.TotalRevenue...)
	stats.TotalCharityDistributed = stats.TotalCharityDistributed.Add(dist.CharityAmount...)
	
	if !dist.YieldAmount.IsZero() {
		stats.TotalYieldDistributed = stats.TotalYieldDistributed.Add(dist.YieldAmount...)
	}
	
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PlatformStatsKey, k.cdc.MustMarshal(&stats))
}

// Helper functions (would be implemented based on actual module integration)
func (k Keeper) calculatePlatformExpenses(ctx sdk.Context) sdk.Dec {
	// Placeholder - would calculate actual expenses
	return sdk.NewDec(800000)
}

func (k Keeper) getTradingVolume(ctx sdk.Context) sdk.Dec {
	// Placeholder - would get from DEX module
	return sdk.NewDec(50000000)
}

func (k Keeper) getPreviousTradingVolume(ctx sdk.Context) sdk.Dec {
	// Placeholder - would get historical volume
	return sdk.NewDec(40000000)
}

func (k Keeper) getActiveUserCount(ctx sdk.Context) uint64 {
	// Placeholder - would count active users
	return 10000
}

func (k Keeper) getTransactionCount(ctx sdk.Context) uint64 {
	// Placeholder - would count transactions
	return 50000
}

func (k Keeper) getTotalRevenue(ctx sdk.Context) sdk.Coins {
	// Placeholder - would aggregate revenue from all sources
	return sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000000)))
}