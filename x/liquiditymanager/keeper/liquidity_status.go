package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/liquiditymanager/types"
)

// LiquidityStatus represents the current state of lending liquidity
type LiquidityStatus string

const (
	LiquidityStatusBuilding LiquidityStatus = "BUILDING"
	LiquidityStatusBasic    LiquidityStatus = "BASIC"
	LiquidityStatusMedium   LiquidityStatus = "MEDIUM"
	LiquidityStatusFull     LiquidityStatus = "FULL"
	LiquidityStatusPaused   LiquidityStatus = "PAUSED"
)

// LiquidityInfo contains comprehensive information about pool liquidity
type LiquidityInfo struct {
	TotalPoolValue      sdk.Dec         `json:"total_pool_value"`
	AvailableForLending sdk.Dec         `json:"available_for_lending"`
	ReserveAmount       sdk.Dec         `json:"reserve_amount"`
	EmergencyReserve    sdk.Dec         `json:"emergency_reserve"`
	Status              LiquidityStatus `json:"status"`
	MaxLoanAmount       sdk.Dec         `json:"max_loan_amount"`
	DailyLendingLimit   sdk.Dec         `json:"daily_lending_limit"`
	AvailableModules    []string        `json:"available_modules"`
	NextThreshold       sdk.Dec         `json:"next_threshold"`
	ProgressToNext      sdk.Dec         `json:"progress_to_next"`
	EstimatedDaysToNext int64           `json:"estimated_days_to_next"`
}

// GetLiquidityInfo returns comprehensive liquidity information with enhanced risk management
func (k Keeper) GetLiquidityInfo(ctx sdk.Context) LiquidityInfo {
	params := k.GetParams(ctx)
	totalPool := k.GetTotalPoolValue(ctx)
	
	// REVOLUTIONARY: Calculate all required reserves for financial safety
	reserveAmount := totalPool.Mul(params.MinimumReserveRatio)              // 50% minimum reserve
	emergencyReserve := totalPool.Mul(params.EmergencyReserveRatio)         // 15% emergency reserve
	loanLossProvision := totalPool.Mul(params.LoanLossProvisionRatio)       // 10% loan loss provision
	
	// Total reserves = 75% of pool (industry-leading safety)
	totalReserves := reserveAmount.Add(emergencyReserve).Add(loanLossProvision)
	availableForLending := totalPool.Sub(totalReserves)
	
	// Ensure available amount is not negative (critical safety check)
	if availableForLending.IsNegative() {
		availableForLending = sdk.ZeroDec()
	}

	// Determine current status and capabilities
	status := k.calculateLiquidityStatus(totalPool, params)
	maxLoanAmount := k.getMaxLoanAmount(status)
	dailyLimit := k.getDailyLendingLimit(status, availableForLending)
	modules := k.getAvailableModules(status)
	
	// Calculate progress to next threshold
	nextThreshold, progress := k.calculateProgressToNext(totalPool, params)
	estimatedDays := k.estimateDaysToNextThreshold(totalPool, nextThreshold, ctx)

	return LiquidityInfo{
		TotalPoolValue:      totalPool,
		AvailableForLending: availableForLending,
		ReserveAmount:       reserveAmount,
		EmergencyReserve:    emergencyReserve,
		Status:              status,
		MaxLoanAmount:       maxLoanAmount,
		DailyLendingLimit:   dailyLimit,
		AvailableModules:    modules,
		NextThreshold:       nextThreshold,
		ProgressToNext:      progress,
		EstimatedDaysToNext: estimatedDays,
	}
}

// GetTotalPoolValue calculates the total value across all liquidity pools
func (k Keeper) GetTotalPoolValue(ctx sdk.Context) sdk.Dec {
	// Get values from Suraksha pools (village + urban)
	surakshaValue := k.getSurakshaPoolValue(ctx)
	
	// Get values from Money Order DEX liquidity
	dexValue := k.getDEXLiquidityValue(ctx)
	
	// Get values from lending pools
	lendingValue := k.getLendingPoolValue(ctx)
	
	return surakshaValue.Add(dexValue).Add(lendingValue)
}

// calculateLiquidityStatus determines the current lending status
func (k Keeper) calculateLiquidityStatus(totalPool sdk.Dec, params types.LiquidityParams) LiquidityStatus {
	// Check if we're below emergency reserve threshold
	emergencyThreshold := totalPool.Mul(params.EmergencyReserveRatio)
	if totalPool.LTE(emergencyThreshold) {
		return LiquidityStatusPaused
	}

	// Check threshold levels
	fullThreshold := sdk.NewDecFromInt(params.LendingFullThreshold)
	mediumThreshold := sdk.NewDecFromInt(params.LendingMediumThreshold)
	basicThreshold := sdk.NewDecFromInt(params.LendingBasicThreshold)

	if totalPool.GTE(fullThreshold) {
		return LiquidityStatusFull
	} else if totalPool.GTE(mediumThreshold) {
		return LiquidityStatusMedium
	} else if totalPool.GTE(basicThreshold) {
		return LiquidityStatusBasic
	} else {
		return LiquidityStatusBuilding
	}
}

// getMaxLoanAmount returns the maximum loan amount for current status
func (k Keeper) getMaxLoanAmount(status LiquidityStatus) sdk.Dec {
	switch status {
	case LiquidityStatusBasic:
		return sdk.NewDec(50000000000)  // ₹50K in micro-NAMO
	case LiquidityStatusMedium:
		return sdk.NewDec(200000000000) // ₹2L in micro-NAMO
	case LiquidityStatusFull:
		return sdk.NewDec(500000000000) // ₹5L in micro-NAMO
	default:
		return sdk.ZeroDec()
	}
}

// getDailyLendingLimit returns FINANCIALLY SOUND daily lending limits to prevent liquidity crisis
func (k Keeper) getDailyLendingLimit(status LiquidityStatus, available sdk.Dec) sdk.Dec {
	switch status {
	case LiquidityStatusBasic:
		// CONSERVATIVE: Max 1% of available liquidity per day (prevents pool exhaustion)
		return available.Mul(sdk.NewDecWithPrec(1, 2))
	case LiquidityStatusMedium:
		// CONSERVATIVE: Max 2% of available liquidity per day
		return available.Mul(sdk.NewDecWithPrec(2, 2))
	case LiquidityStatusFull:
		// CONSERVATIVE: Max 3% of available liquidity per day (vs previous 15%!)
		return available.Mul(sdk.NewDecWithPrec(3, 2))
	default:
		return sdk.ZeroDec()
	}
}

// getAvailableModules returns which lending modules are available
func (k Keeper) getAvailableModules(status LiquidityStatus) []string {
	switch status {
	case LiquidityStatusBasic:
		return []string{"krishimitra"}
	case LiquidityStatusMedium:
		return []string{"krishimitra", "vyavasayamitra"}
	case LiquidityStatusFull:
		return []string{"krishimitra", "vyavasayamitra", "shikshamitra"}
	default:
		return []string{}
	}
}

// IsLendingAvailable checks if lending is currently available
func (k Keeper) IsLendingAvailable(ctx sdk.Context) bool {
	info := k.GetLiquidityInfo(ctx)
	return info.Status != LiquidityStatusBuilding && info.Status != LiquidityStatusPaused
}

// CanProcessLoan checks if a specific loan can be processed with member-only restrictions
func (k Keeper) CanProcessLoan(ctx sdk.Context, amount sdk.Dec, module string, borrower sdk.AccAddress) (bool, string) {
	info := k.GetLiquidityInfo(ctx)
	
	// Check if lending is available at all
	if !k.IsLendingAvailable(ctx) {
		return false, k.getLendingUnavailableMessage(info.Status, info.EstimatedDaysToNext)
	}
	
	// REVOLUTIONARY RESTRICTION 1: Only pool members can get loans
	if !k.IsPoolMember(ctx, borrower) {
		return false, "🚫 Lending access restricted to Suraksha Pool members only. Join Village or Urban Suraksha Pool to access revolutionary lending rates!"
	}
	
	// Check if module is available
	moduleAvailable := false
	for _, availableModule := range info.AvailableModules {
		if availableModule == module {
			moduleAvailable = true
			break
		}
	}
	if !moduleAvailable {
		return false, fmt.Sprintf("Module %s not available at current liquidity level", module)
	}
	
	// Check loan amount limits
	if amount.GT(info.MaxLoanAmount) {
		return false, fmt.Sprintf("Loan amount exceeds maximum of ₹%.0f", info.MaxLoanAmount.Quo(sdk.NewDec(1000000)).TruncateInt64())
	}
	
	// Check available liquidity
	if amount.GT(info.AvailableForLending) {
		return false, "Insufficient liquidity available for lending"
	}
	
	// Check daily lending limit
	dailyUsed := k.getDailyLendingUsed(ctx)
	if dailyUsed.Add(amount).GT(info.DailyLendingLimit) {
		return false, "Daily lending limit reached. Please try again tomorrow"
	}
	
	return true, ""
}

// CanProcessCollateralLoan checks if a NAMO collateral loan can be processed
func (k Keeper) CanProcessCollateralLoan(ctx sdk.Context, loanAmount sdk.Dec, collateralAmount sdk.Dec, borrower sdk.AccAddress) (bool, string) {
	// REVOLUTIONARY RESTRICTION 2: 70% LTV against staked NAMO
	maxLoanAmount := collateralAmount.Mul(sdk.NewDecWithPrec(70, 2)) // 70% LTV
	
	if loanAmount.GT(maxLoanAmount) {
		return false, fmt.Sprintf("🏛️ Loan amount exceeds 70%% collateral limit. Max loan: ₹%.0f against ₹%.0f collateral", 
			maxLoanAmount.Quo(sdk.NewDec(1000000)).TruncateInt64(),
			collateralAmount.Quo(sdk.NewDec(1000000)).TruncateInt64())
	}
	
	// Verify borrower has sufficient staked NAMO
	stakedAmount := k.GetStakedNAMO(ctx, borrower)
	if stakedAmount.LT(collateralAmount) {
		return false, fmt.Sprintf("🔒 Insufficient staked NAMO. Required: ₹%.0f, Available: ₹%.0f", 
			collateralAmount.Quo(sdk.NewDec(1000000)).TruncateInt64(),
			stakedAmount.Quo(sdk.NewDec(1000000)).TruncateInt64())
	}
	
	// Check if collateral is already used
	if k.IsCollateralLocked(ctx, borrower, collateralAmount) {
		return false, "💎 Part of your staked NAMO is already used as collateral. Stake more NAMO or repay existing loans"
	}
	
	return true, ""
}

// getLendingUnavailableMessage returns user-friendly message when lending is unavailable
func (k Keeper) getLendingUnavailableMessage(status LiquidityStatus, estimatedDays int64) string {
	switch status {
	case LiquidityStatusBuilding:
		if estimatedDays > 0 {
			return fmt.Sprintf("Lending will be available when our liquidity pool reaches ₹10 Cr. Estimated in %d days. Join Suraksha Pool to help build liquidity faster!", estimatedDays)
		}
		return "Lending will be available when our liquidity pool reaches ₹10 Cr. Join Suraksha Pool to help build liquidity faster!"
	case LiquidityStatusPaused:
		return "Lending is temporarily paused due to low liquidity. It will resume automatically when reserves are restored."
	default:
		return "Lending is currently unavailable. Please check back later."
	}
}

// Helper functions for pool value calculations
func (k Keeper) getSurakshaPoolValue(ctx sdk.Context) sdk.Dec {
	// This would integrate with the Suraksha module to get total pool value
	// For now, return a placeholder implementation
	return sdk.ZeroDec()
}

func (k Keeper) getDEXLiquidityValue(ctx sdk.Context) sdk.Dec {
	// This would integrate with the MoneyOrder DEX to get liquidity value
	// For now, return a placeholder implementation
	return sdk.ZeroDec()
}

func (k Keeper) getLendingPoolValue(ctx sdk.Context) sdk.Dec {
	// This would get the value currently deployed in active loans
	// For now, return a placeholder implementation
	return sdk.ZeroDec()
}

func (k Keeper) calculateProgressToNext(totalPool sdk.Dec, params types.LiquidityParams) (sdk.Dec, sdk.Dec) {
	fullThreshold := sdk.NewDecFromInt(params.LendingFullThreshold)
	mediumThreshold := sdk.NewDecFromInt(params.LendingMediumThreshold)
	basicThreshold := sdk.NewDecFromInt(params.LendingBasicThreshold)

	var nextThreshold sdk.Dec
	var progress sdk.Dec

	if totalPool.LT(basicThreshold) {
		nextThreshold = basicThreshold
		progress = totalPool.Quo(basicThreshold)
	} else if totalPool.LT(mediumThreshold) {
		nextThreshold = mediumThreshold
		progress = totalPool.Sub(basicThreshold).Quo(mediumThreshold.Sub(basicThreshold))
	} else if totalPool.LT(fullThreshold) {
		nextThreshold = fullThreshold
		progress = totalPool.Sub(mediumThreshold).Quo(fullThreshold.Sub(mediumThreshold))
	} else {
		nextThreshold = fullThreshold
		progress = sdk.OneDec()
	}

	return nextThreshold, progress
}

func (k Keeper) estimateDaysToNextThreshold(currentPool, nextThreshold sdk.Dec, ctx sdk.Context) int64 {
	// This would calculate based on historical pool growth rates
	// For now, return a simple estimation based on daily contribution patterns
	remaining := nextThreshold.Sub(currentPool)
	if remaining.LTE(sdk.ZeroDec()) {
		return 0
	}
	
	// Assume average daily contribution of ₹50L (placeholder)
	averageDailyGrowth := sdk.NewDec(50000000000) // ₹50L in micro-NAMO
	days := remaining.Quo(averageDailyGrowth).TruncateInt64()
	
	if days < 0 {
		return 0
	}
	return days
}

func (k Keeper) getDailyLendingUsed(ctx sdk.Context) sdk.Dec {
	// This would track daily lending usage
	// For now, return zero as placeholder
	return sdk.ZeroDec()
}

// REVOLUTIONARY FINANCIAL PROTECTION FUNCTIONS

// IsPoolMember checks if user is a member of Village or Urban Suraksha Pool
func (k Keeper) IsPoolMember(ctx sdk.Context, user sdk.AccAddress) bool {
	// Check Village Suraksha Pool membership
	if k.isVillageSurakshaPoolMember(ctx, user) {
		return true
	}
	
	// Check Urban Suraksha Pool membership  
	if k.isUrbanSurakshaPoolMember(ctx, user) {
		return true
	}
	
	// Check if user holds minimum NAMO staking requirement
	stakedAmount := k.GetStakedNAMO(ctx, user)
	minimumStakeForLending := sdk.NewDec(100000000000) // ₹1L minimum stake for lending access
	
	return stakedAmount.GTE(minimumStakeForLending)
}

// isVillageSurakshaPoolMember checks Village pool membership
func (k Keeper) isVillageSurakshaPoolMember(ctx sdk.Context, user sdk.AccAddress) bool {
	// This would integrate with the GramSuraksha module
	// Check if user has active village pool participation
	// For now, placeholder implementation
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("village_member_"), user.Bytes()...)
	return store.Has(key)
}

// isUrbanSurakshaPoolMember checks Urban pool membership
func (k Keeper) isUrbanSurakshaPoolMember(ctx sdk.Context, user sdk.AccAddress) bool {
	// This would integrate with the UrbanSuraksha module
	// Check if user has active urban pool participation
	// For now, placeholder implementation
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("urban_member_"), user.Bytes()...)
	return store.Has(key)
}

// GetStakedNAMO returns the amount of NAMO staked by user
func (k Keeper) GetStakedNAMO(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	// This would integrate with the staking module or NAMO module
	// Get total staked NAMO for the user
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("staked_namo_"), user.Bytes()...)
	
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	
	var amount sdk.Dec
	k.cdc.MustUnmarshal(bz, &amount)
	return amount
}

// IsCollateralLocked checks if NAMO collateral is already being used
func (k Keeper) IsCollateralLocked(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) bool {
	// Get total locked collateral for user
	lockedAmount := k.getLockedCollateral(ctx, user)
	availableCollateral := k.GetStakedNAMO(ctx, user).Sub(lockedAmount)
	
	return amount.GT(availableCollateral)
}

// getLockedCollateral returns total NAMO locked as collateral
func (k Keeper) getLockedCollateral(ctx sdk.Context, user sdk.AccAddress) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("locked_collateral_"), user.Bytes()...)
	
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	
	var amount sdk.Dec
	k.cdc.MustUnmarshal(bz, &amount)
	return amount
}

// LockCollateral locks NAMO as collateral for a loan
func (k Keeper) LockCollateral(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) error {
	currentLocked := k.getLockedCollateral(ctx, user)
	newLocked := currentLocked.Add(amount)
	
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("locked_collateral_"), user.Bytes()...)
	
	bz := k.cdc.MustMarshal(&newLocked)
	store.Set(key, bz)
	
	return nil
}

// UnlockCollateral unlocks NAMO collateral when loan is repaid
func (k Keeper) UnlockCollateral(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) error {
	currentLocked := k.getLockedCollateral(ctx, user)
	newLocked := currentLocked.Sub(amount)
	
	if newLocked.IsNegative() {
		newLocked = sdk.ZeroDec()
	}
	
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("locked_collateral_"), user.Bytes()...)
	
	bz := k.cdc.MustMarshal(&newLocked)
	store.Set(key, bz)
	
	return nil
}

// SetPoolMembership sets pool membership for a user
func (k Keeper) SetPoolMembership(ctx sdk.Context, user sdk.AccAddress, poolType string, active bool) {
	store := ctx.KVStore(k.storeKey)
	var key []byte
	
	switch poolType {
	case "village":
		key = append([]byte("village_member_"), user.Bytes()...)
	case "urban":
		key = append([]byte("urban_member_"), user.Bytes()...)
	default:
		return
	}
	
	if active {
		store.Set(key, []byte{1})
	} else {
		store.Delete(key)
	}
}

// SetStakedNAMO sets the staked NAMO amount for a user
func (k Keeper) SetStakedNAMO(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("staked_namo_"), user.Bytes()...)
	
	bz := k.cdc.MustMarshal(&amount)
	store.Set(key, bz)
}

// REVOLUTIONARY FEE CALCULATION FUNCTIONS

// CalculateProcessingFee calculates the processing fee with ₹2500 cap (ultra-competitive)
func (k Keeper) CalculateProcessingFee(ctx sdk.Context, loanAmount sdk.Dec) sdk.Dec {
	params := k.GetParams(ctx)
	
	// Calculate 1% processing fee
	processingFee := loanAmount.Mul(params.ProcessingFeeRate)
	
	// Apply ₹2500 cap (borrower protection)
	feeCap := sdk.NewDecFromInt(params.ProcessingFeeCap)
	if processingFee.GT(feeCap) {
		processingFee = feeCap
	}
	
	return processingFee
}

// CalculateEarlySettlementFee calculates early settlement fee (0.5% on remaining principal)
func (k Keeper) CalculateEarlySettlementFee(ctx sdk.Context, remainingPrincipal sdk.Dec) sdk.Dec {
	params := k.GetParams(ctx)
	
	// Calculate 0.5% early settlement fee on remaining principal only
	earlySettlementFee := remainingPrincipal.Mul(params.EarlySettlementRate)
	
	return earlySettlementFee
}

// CalculateDisbursementAmount calculates actual disbursed amount (99% after fees)
func (k Keeper) CalculateDisbursementAmount(ctx sdk.Context, approvedAmount sdk.Dec) (sdk.Dec, sdk.Dec) {
	processingFee := k.CalculateProcessingFee(ctx, approvedAmount)
	disbursedAmount := approvedAmount.Sub(processingFee)
	
	return disbursedAmount, processingFee
}

// GetLoanBreakdown provides comprehensive loan cost breakdown for transparency
func (k Keeper) GetLoanBreakdown(ctx sdk.Context, loanAmount sdk.Dec, interestRate sdk.Dec, termMonths int64) LoanBreakdown {
	processingFee := k.CalculateProcessingFee(ctx, loanAmount)
	disbursedAmount := loanAmount.Sub(processingFee)
	
	// Calculate total interest over loan term
	monthlyRate := interestRate.Quo(sdk.NewDec(12))
	totalInterest := loanAmount.Mul(monthlyRate).Mul(sdk.NewDec(termMonths))
	
	// Calculate monthly EMI (simplified calculation)
	totalRepayment := loanAmount.Add(totalInterest)
	monthlyEMI := totalRepayment.Quo(sdk.NewDec(termMonths))
	
	return LoanBreakdown{
		ApprovedAmount:   loanAmount,
		ProcessingFee:    processingFee,
		DisbursedAmount:  disbursedAmount,
		InterestRate:     interestRate,
		TotalInterest:    totalInterest,
		TotalRepayment:   totalRepayment,
		MonthlyEMI:       monthlyEMI,
		TermMonths:       termMonths,
		EffectiveFeeRate: processingFee.Quo(loanAmount), // Actual fee percentage
	}
}

// LoanBreakdown provides transparent loan cost information
type LoanBreakdown struct {
	ApprovedAmount   sdk.Dec `json:"approved_amount"`
	ProcessingFee    sdk.Dec `json:"processing_fee"`
	DisbursedAmount  sdk.Dec `json:"disbursed_amount"`
	InterestRate     sdk.Dec `json:"interest_rate"`
	TotalInterest    sdk.Dec `json:"total_interest"`
	TotalRepayment   sdk.Dec `json:"total_repayment"`
	MonthlyEMI       sdk.Dec `json:"monthly_emi"`
	TermMonths       int64   `json:"term_months"`
	EffectiveFeeRate sdk.Dec `json:"effective_fee_rate"`
}