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

package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// UnifiedLiquidityPool manages liquidity across Pension, DEX, and Agricultural Lending
type UnifiedLiquidityPool struct {
	PoolId              uint64    `json:"pool_id"`
	VillagePoolId       uint64    `json:"village_pool_id"`
	TotalLiquidity      sdk.Coins `json:"total_liquidity"`
	
	// Allocations
	SurakshaReserve      sdk.Coins `json:"suraksha_reserve"`      // 20% - Reserved for pension payouts
	DexLiquidity        sdk.Coins `json:"dex_liquidity"`        // 30% - For Money Order trading
	AgriLendingPool     sdk.Coins `json:"agri_lending_pool"`   // 40% - For Kisaan Mitra loans
	EmergencyReserve    sdk.Coins `json:"emergency_reserve"`    // 10% - Emergency buffer
	
	// Tracking
	ActivePensionAccounts uint32    `json:"active_pension_accounts"`
	ActiveDexPairs        uint32    `json:"active_dex_pairs"`
	ActiveAgriLoans       uint32    `json:"active_agri_loans"`
	
	// Performance
	MonthlyDexRevenue     sdk.Coins `json:"monthly_dex_revenue"`
	MonthlyLendingRevenue sdk.Coins `json:"monthly_lending_revenue"`
	TotalRevenue          sdk.Coins `json:"total_revenue"`
	
	// Configuration
	AllocationConfig      AllocationConfig `json:"allocation_config"`
	LastRebalanceTime     time.Time        `json:"last_rebalance_time"`
	CreatedAt             time.Time        `json:"created_at"`
}

// AllocationConfig defines how liquidity is allocated across different uses
type AllocationConfig struct {
	SurakshaReserveRatio   sdk.Dec `json:"suraksha_reserve_ratio"`    // Default: 20%
	DexLiquidityRatio     sdk.Dec `json:"dex_liquidity_ratio"`      // Default: 30%
	AgriLendingRatio      sdk.Dec `json:"agri_lending_ratio"`       // Default: 40%
	EmergencyReserveRatio sdk.Dec `json:"emergency_reserve_ratio"`  // Default: 10%
	
	// Risk parameters
	MaxLendingExposure    sdk.Dec `json:"max_lending_exposure"`     // Max % that can be lent out
	MinDexLiquidity       sdk.Dec `json:"min_dex_liquidity"`        // Min % for DEX operations
	RebalanceFrequency    uint32  `json:"rebalance_frequency_days"` // How often to rebalance
}

// DefaultAllocationConfig returns the default allocation configuration
func DefaultAllocationConfig() AllocationConfig {
	return AllocationConfig{
		SurakshaReserveRatio:   sdk.NewDecWithPrec(20, 2), // 20%
		DexLiquidityRatio:     sdk.NewDecWithPrec(30, 2), // 30%
		AgriLendingRatio:      sdk.NewDecWithPrec(40, 2), // 40%
		EmergencyReserveRatio: sdk.NewDecWithPrec(10, 2), // 10%
		MaxLendingExposure:    sdk.NewDecWithPrec(70, 2), // Max 70% can be lent
		MinDexLiquidity:       sdk.NewDecWithPrec(20, 2), // Min 20% for DEX
		RebalanceFrequency:    7,                          // Weekly rebalancing
	}
}

// CreateUnifiedLiquidityPool creates a new unified pool for a village
func (k Keeper) CreateUnifiedLiquidityPool(
	ctx sdk.Context,
	villagePoolId uint64,
	initialLiquidity sdk.Coins,
) (*UnifiedLiquidityPool, error) {
	// Verify village pool exists
	villagePool, found := k.GetVillagePool(ctx, villagePoolId)
	if !found {
		return nil, types.ErrVillagePoolNotFound
	}

	if !villagePool.Active || !villagePool.Verified {
		return nil, sdkerrors.Wrap(types.ErrVillagePoolInactive, "pool must be active and verified")
	}

	// Get next pool ID
	poolId := k.GetNextUnifiedPoolId(ctx)
	
	// Create unified pool with default allocation
	config := DefaultAllocationConfig()
	pool := &UnifiedLiquidityPool{
		PoolId:            poolId,
		VillagePoolId:     villagePoolId,
		TotalLiquidity:    initialLiquidity,
		AllocationConfig:  config,
		CreatedAt:         ctx.BlockTime(),
		LastRebalanceTime: ctx.BlockTime(),
	}

	// Perform initial allocation
	if err := k.allocateLiquidity(ctx, pool); err != nil {
		return nil, err
	}

	// Save pool
	k.SetUnifiedLiquidityPool(ctx, *pool)
	k.SetNextUnifiedPoolId(ctx, poolId+1)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnifiedPoolCreated,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyVillagePoolId, fmt.Sprintf("%d", villagePoolId)),
			sdk.NewAttribute(types.AttributeKeyLiquidity, initialLiquidity.String()),
		),
	)

	return pool, nil
}

// AddSurakshaContribution adds pension contribution to unified pool
func (k Keeper) AddSurakshaContribution(
	ctx sdk.Context,
	unifiedPoolId uint64,
	contribution sdk.Coin,
	pensionAccountId string,
) error {
	pool, found := k.GetUnifiedLiquidityPool(ctx, unifiedPoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Add to total liquidity
	pool.TotalLiquidity = pool.TotalLiquidity.Add(contribution)
	pool.ActivePensionAccounts++

	// Rebalance allocations
	if err := k.rebalanceIfNeeded(ctx, &pool); err != nil {
		return err
	}

	// Track pension contribution
	k.SetSurakshaContribution(ctx, SurakshaContribution{
		PensionAccountId: pensionAccountId,
		UnifiedPoolId:    unifiedPoolId,
		Amount:           contribution,
		ContributionTime: ctx.BlockTime(),
		MaturityTime:     ctx.BlockTime().AddDate(0, 12, 0), // 12 months
	})

	// Update pool
	k.SetUnifiedLiquidityPool(ctx, pool)

	return nil
}

// ProcessAgriLoan processes an agricultural loan from the unified pool
func (k Keeper) ProcessAgriLoan(
	ctx sdk.Context,
	unifiedPoolId uint64,
	loanAmount sdk.Coin,
	borrower sdk.AccAddress,
	loanType string,
	duration uint32, // months
) error {
	pool, found := k.GetUnifiedLiquidityPool(ctx, unifiedPoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Check if sufficient lending capacity
	availableForLending := pool.AgriLendingPool.AmountOf(loanAmount.Denom)
	if availableForLending.LT(loanAmount.Amount) {
		return sdkerrors.Wrap(types.ErrInsufficientLiquidity, "insufficient agricultural lending capacity")
	}

	// Calculate interest rate based on loan type
	interestRate := k.getAgriLoanInterestRate(loanType)
	
	// Deduct from lending pool
	pool.AgriLendingPool = pool.AgriLendingPool.Sub(loanAmount)
	pool.ActiveAgriLoans++

	// Create loan record
	k.SetAgriLoan(ctx, AgriLoan{
		LoanId:           k.GetNextLoanId(ctx),
		UnifiedPoolId:    unifiedPoolId,
		Borrower:         borrower.String(),
		Amount:           loanAmount,
		InterestRate:     interestRate,
		LoanType:         loanType,
		Duration:         duration,
		DisbursementTime: ctx.BlockTime(),
		MaturityTime:     ctx.BlockTime().AddDate(0, int(duration), 0),
		Status:           "active",
	})

	// Update pool
	k.SetUnifiedLiquidityPool(ctx, pool)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAgriLoanDisbursed,
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", unifiedPoolId)),
			sdk.NewAttribute("borrower", borrower.String()),
			sdk.NewAttribute("amount", loanAmount.String()),
			sdk.NewAttribute("interest_rate", interestRate.String()),
			sdk.NewAttribute("duration_months", fmt.Sprintf("%d", duration)),
		),
	)

	return nil
}

// ProcessLoanRepayment handles agricultural loan repayment
func (k Keeper) ProcessLoanRepayment(
	ctx sdk.Context,
	loanId uint64,
	repaymentAmount sdk.Coin,
) error {
	loan, found := k.GetAgriLoan(ctx, loanId)
	if !found {
		return sdkerrors.Wrap(types.ErrNotFound, "loan not found")
	}

	pool, found := k.GetUnifiedLiquidityPool(ctx, loan.UnifiedPoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Calculate interest earned
	principal := loan.Amount
	interestEarned := repaymentAmount.Amount.Sub(principal.Amount)
	
	// Add principal back to lending pool
	pool.AgriLendingPool = pool.AgriLendingPool.Add(principal)
	
	// Distribute interest as revenue
	if interestEarned.GT(sdk.ZeroInt()) {
		interestCoin := sdk.NewCoin(repaymentAmount.Denom, interestEarned)
		pool.MonthlyLendingRevenue = pool.MonthlyLendingRevenue.Add(interestCoin)
		pool.TotalRevenue = pool.TotalRevenue.Add(interestCoin)
		
		// Add to total liquidity for redistribution
		pool.TotalLiquidity = pool.TotalLiquidity.Add(interestCoin)
	}

	// Update loan status
	loan.Status = "repaid"
	loan.RepaymentTime = ctx.BlockTime()
	loan.RepaymentAmount = repaymentAmount
	k.SetAgriLoan(ctx, loan)

	// Update pool
	pool.ActiveAgriLoans--
	k.SetUnifiedLiquidityPool(ctx, pool)

	return nil
}

// RecordDexRevenue records trading fee revenue from Money Order DEX
func (k Keeper) RecordDexRevenue(
	ctx sdk.Context,
	unifiedPoolId uint64,
	tradingFees sdk.Coins,
) error {
	pool, found := k.GetUnifiedLiquidityPool(ctx, unifiedPoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Add to revenue tracking
	pool.MonthlyDexRevenue = pool.MonthlyDexRevenue.Add(tradingFees...)
	pool.TotalRevenue = pool.TotalRevenue.Add(tradingFees...)
	
	// Add to total liquidity for redistribution
	pool.TotalLiquidity = pool.TotalLiquidity.Add(tradingFees...)

	// Update pool
	k.SetUnifiedLiquidityPool(ctx, pool)

	return nil
}

// CalculatePensionReturns calculates returns for pension accounts from unified revenue
func (k Keeper) CalculatePensionReturns(
	ctx sdk.Context,
	unifiedPoolId uint64,
) (sdk.Coins, error) {
	pool, found := k.GetUnifiedLiquidityPool(ctx, unifiedPoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	// Total revenue from DEX + Agricultural lending
	totalMonthlyRevenue := pool.MonthlyDexRevenue.Add(pool.MonthlyLendingRevenue...)
	
	// Calculate pension share (prioritized to ensure 50% returns)
	// Monthly contribution: 1000 NAMO * Active accounts
	// Required monthly return: 4.17% (to achieve 50% in 12 months)
	monthlyContributions := sdk.NewInt(1000).Mul(sdk.NewInt(int64(pool.ActivePensionAccounts)))
	requiredReturn := monthlyContributions.ToDec().Mul(sdk.NewDecWithPrec(417, 4)) // 4.17%
	
	pensionReturns := sdk.NewCoins()
	for _, revenue := range totalMonthlyRevenue {
		// Allocate revenue to meet pension return requirements
		allocation := sdk.MinInt(revenue.Amount, requiredReturn.TruncateInt())
		pensionReturns = pensionReturns.Add(sdk.NewCoin(revenue.Denom, allocation))
	}

	return pensionReturns, nil
}

// Helper functions

// allocateLiquidity distributes liquidity according to allocation config
func (k Keeper) allocateLiquidity(ctx sdk.Context, pool *UnifiedLiquidityPool) error {
	config := pool.AllocationConfig
	
	// Reset allocations
	pool.SurakshaReserve = sdk.NewCoins()
	pool.DexLiquidity = sdk.NewCoins()
	pool.AgriLendingPool = sdk.NewCoins()
	pool.EmergencyReserve = sdk.NewCoins()
	
	// Allocate each coin type
	for _, coin := range pool.TotalLiquidity {
		pensionAmount := coin.Amount.ToDec().Mul(config.SurakshaReserveRatio).TruncateInt()
		dexAmount := coin.Amount.ToDec().Mul(config.DexLiquidityRatio).TruncateInt()
		agriAmount := coin.Amount.ToDec().Mul(config.AgriLendingRatio).TruncateInt()
		emergencyAmount := coin.Amount.ToDec().Mul(config.EmergencyReserveRatio).TruncateInt()
		
		pool.SurakshaReserve = pool.SurakshaReserve.Add(sdk.NewCoin(coin.Denom, pensionAmount))
		pool.DexLiquidity = pool.DexLiquidity.Add(sdk.NewCoin(coin.Denom, dexAmount))
		pool.AgriLendingPool = pool.AgriLendingPool.Add(sdk.NewCoin(coin.Denom, agriAmount))
		pool.EmergencyReserve = pool.EmergencyReserve.Add(sdk.NewCoin(coin.Denom, emergencyAmount))
	}
	
	return nil
}

// rebalanceIfNeeded checks if rebalancing is needed and performs it
func (k Keeper) rebalanceIfNeeded(ctx sdk.Context, pool *UnifiedLiquidityPool) error {
	daysSinceRebalance := uint32(ctx.BlockTime().Sub(pool.LastRebalanceTime).Hours() / 24)
	
	if daysSinceRebalance >= pool.AllocationConfig.RebalanceFrequency {
		if err := k.allocateLiquidity(ctx, pool); err != nil {
			return err
		}
		pool.LastRebalanceTime = ctx.BlockTime()
	}
	
	return nil
}

// getAgriLoanInterestRate returns interest rate based on loan type
func (k Keeper) getAgriLoanInterestRate(loanType string) sdk.Dec {
	switch loanType {
	case "input": // Seeds, fertilizers
		return sdk.NewDecWithPrec(6, 2) // 6%
	case "equipment": // Tractors, tools
		return sdk.NewDecWithPrec(8, 2) // 8%
	case "emergency": // Crop failure, medical
		return sdk.NewDecWithPrec(9, 2) // 9%
	case "expansion": // Land purchase
		return sdk.NewDecWithPrec(10, 2) // 10%
	default:
		return sdk.NewDecWithPrec(8, 2) // Default 8%
	}
}

// Types

type SurakshaContribution struct {
	PensionAccountId string    `json:"pension_account_id"`
	UnifiedPoolId    uint64    `json:"unified_pool_id"`
	Amount           sdk.Coin  `json:"amount"`
	ContributionTime time.Time `json:"contribution_time"`
	MaturityTime     time.Time `json:"maturity_time"`
}

type AgriLoan struct {
	LoanId           uint64    `json:"loan_id"`
	UnifiedPoolId    uint64    `json:"unified_pool_id"`
	Borrower         string    `json:"borrower"`
	Amount           sdk.Coin  `json:"amount"`
	InterestRate     sdk.Dec   `json:"interest_rate"`
	LoanType         string    `json:"loan_type"`
	Duration         uint32    `json:"duration_months"`
	DisbursementTime time.Time `json:"disbursement_time"`
	MaturityTime     time.Time `json:"maturity_time"`
	RepaymentTime    time.Time `json:"repayment_time,omitempty"`
	RepaymentAmount  sdk.Coin  `json:"repayment_amount,omitempty"`
	Status           string    `json:"status"`
}

// Store functions

func (k Keeper) SetUnifiedLiquidityPool(ctx sdk.Context, pool UnifiedLiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	key := getUnifiedPoolKey(pool.PoolId)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

func (k Keeper) GetUnifiedLiquidityPool(ctx sdk.Context, poolId uint64) (UnifiedLiquidityPool, bool) {
	store := ctx.KVStore(k.storeKey)
	key := getUnifiedPoolKey(poolId)
	bz := store.Get(key)
	if bz == nil {
		return UnifiedLiquidityPool{}, false
	}
	
	var pool UnifiedLiquidityPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) GetNextUnifiedPoolId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyNextUnifiedPoolId)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetNextUnifiedPoolId(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyNextUnifiedPoolId, sdk.Uint64ToBigEndian(poolId))
}

func (k Keeper) SetSurakshaContribution(ctx sdk.Context, pc SurakshaContribution) {
	store := ctx.KVStore(k.storeKey)
	key := getSurakshaContributionKey(pc.PensionAccountId, pc.UnifiedPoolId)
	bz := k.cdc.MustMarshal(&pc)
	store.Set(key, bz)
}

func (k Keeper) SetAgriLoan(ctx sdk.Context, loan AgriLoan) {
	store := ctx.KVStore(k.storeKey)
	key := getAgriLoanKey(loan.LoanId)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(key, bz)
}

func (k Keeper) GetAgriLoan(ctx sdk.Context, loanId uint64) (AgriLoan, bool) {
	store := ctx.KVStore(k.storeKey)
	key := getAgriLoanKey(loanId)
	bz := store.Get(key)
	if bz == nil {
		return AgriLoan{}, false
	}
	
	var loan AgriLoan
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

func (k Keeper) GetNextLoanId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyNextLoanId)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// Key functions

func getUnifiedPoolKey(poolId uint64) []byte {
	return append(types.KeyPrefixUnifiedPool, sdk.Uint64ToBigEndian(poolId)...)
}

func getSurakshaContributionKey(pensionAccountId string, poolId uint64) []byte {
	return append(append(types.KeyPrefixSurakshaContribution, []byte(pensionAccountId)...), sdk.Uint64ToBigEndian(poolId)...)
}

func getAgriLoanKey(loanId uint64) []byte {
	return append(types.KeyPrefixAgriLoan, sdk.Uint64ToBigEndian(loanId)...)
}