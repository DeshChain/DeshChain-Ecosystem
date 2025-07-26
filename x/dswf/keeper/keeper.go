package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/dswf/types"
)

// Keeper of the x/dswf store
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper      types.AccountKeeper
	bankKeeper         types.BankKeeper
	stakingKeeper      types.StakingKeeper
	distributionKeeper types.DistributionKeeper
	govKeeper          types.GovKeeper
	revenueKeeper      types.RevenueKeeper

	// the address capable of executing governance actions
	authority string
}

// NewKeeper creates a new DSWF Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distributionKeeper types.DistributionKeeper,
	govKeeper types.GovKeeper,
	revenueKeeper types.RevenueKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		paramstore:         ps,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		stakingKeeper:      stakingKeeper,
		distributionKeeper: distributionKeeper,
		govKeeper:          govKeeper,
		revenueKeeper:      revenueKeeper,
		authority:          authority,
	}
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetModuleAccountAddress returns the module account address
func (k Keeper) GetModuleAccountAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetFundBalance returns the current balance of the DSWF
func (k Keeper) GetFundBalance(ctx sdk.Context) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, k.GetModuleAccountAddress())
}

// GetFundAllocationCount returns the current allocation count
func (k Keeper) GetFundAllocationCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.FundAllocationCountKey)
	if bz == nil {
		return 0
	}
	return types.GetUint64FromBytes(bz)
}

// SetFundAllocationCount sets the allocation count
func (k Keeper) SetFundAllocationCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FundAllocationCountKey, types.GetUint64Bytes(count))
}

// IncrementFundAllocationCount increments and returns the allocation count
func (k Keeper) IncrementFundAllocationCount(ctx sdk.Context) uint64 {
	count := k.GetFundAllocationCount(ctx)
	count++
	k.SetFundAllocationCount(ctx, count)
	return count
}

// SetFundAllocation sets a fund allocation
func (k Keeper) SetFundAllocation(ctx sdk.Context, allocation types.FundAllocation) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&allocation)
	store.Set(types.GetFundAllocationKey(allocation.Id), bz)
	
	// Set indexes
	store.Set(types.GetAllocationByCategoryKey(allocation.Category, allocation.Id), []byte{})
	store.Set(types.GetAllocationByStatusKey(allocation.Status, allocation.Id), []byte{})
	store.Set(types.GetAllocationByRecipientKey(allocation.Recipient, allocation.Id), []byte{})
}

// GetFundAllocation returns a fund allocation by ID
func (k Keeper) GetFundAllocation(ctx sdk.Context, allocationID uint64) (types.FundAllocation, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFundAllocationKey(allocationID))
	if bz == nil {
		return types.FundAllocation{}, false
	}
	
	var allocation types.FundAllocation
	k.cdc.MustUnmarshal(bz, &allocation)
	return allocation, true
}

// ValidateMilestoneProof validates milestone proof for disbursement
func (k Keeper) ValidateMilestoneProof(ctx sdk.Context, proof string) error {
	if proof == "" {
		return types.ErrInvalidProof.Wrap("milestone proof cannot be empty")
	}
	
	// In production, this would:
	// 1. Verify digital signatures on the proof
	// 2. Check proof format and structure
	// 3. Validate against expected milestone criteria
	// 4. Ensure proof hasn't been tampered with
	
	// For now, basic validation
	if len(proof) < 32 {
		return types.ErrInvalidProof.Wrap("milestone proof too short")
	}
	
	return nil
}

// ExecutePortfolioRebalancing executes portfolio rebalancing based on strategy
func (k Keeper) ExecutePortfolioRebalancing(ctx sdk.Context, portfolio *types.InvestmentPortfolio) error {
	params := k.GetParams(ctx)
	
	// Get target allocations based on risk profile
	targetAllocations := k.calculateTargetAllocations(params.TargetReturnRate)
	
	// Calculate current allocations
	currentAllocations := k.calculateCurrentAllocations(portfolio)
	
	// Execute rebalancing transactions
	for assetType, targetPct := range targetAllocations {
		currentPct := currentAllocations[assetType]
		
		// If deviation is > 5%, rebalance
		deviation := targetPct - currentPct
		if deviation > sdk.NewDecWithPrec(5, 2) || deviation < sdk.NewDecWithPrec(-5, 2) {
			if err := k.executeAssetRebalancing(ctx, portfolio, assetType, targetPct); err != nil {
				return err
			}
		}
	}
	
	// Update portfolio metrics
	portfolio.RiskScore = k.calculatePortfolioRiskScore(portfolio)
	portfolio.AnnualReturnRate = k.calculateExpectedReturn(portfolio)
	
	return nil
}

// Helper methods for portfolio rebalancing
func (k Keeper) calculateTargetAllocations(targetReturn sdk.Dec) map[string]sdk.Dec {
	// Conservative strategy for sovereign wealth fund
	return map[string]sdk.Dec{
		"conservative": sdk.NewDecWithPrec(30, 2), // 30%
		"growth":       sdk.NewDecWithPrec(40, 2), // 40%
		"innovation":   sdk.NewDecWithPrec(20, 2), // 20%
		"reserve":      sdk.NewDecWithPrec(10, 2), // 10%
	}
}

func (k Keeper) calculateCurrentAllocations(portfolio *types.InvestmentPortfolio) map[string]sdk.Dec {
	if portfolio.TotalValue.Amount.IsZero() {
		return make(map[string]sdk.Dec)
	}
	
	allocations := make(map[string]sdk.Dec)
	totalValue := portfolio.TotalValue.Amount
	
	for _, component := range portfolio.Components {
		pct := sdk.NewDecFromInt(component.CurrentValue.Amount).Quo(sdk.NewDecFromInt(totalValue))
		allocations[component.AssetType] = pct
	}
	
	return allocations
}

func (k Keeper) executeAssetRebalancing(ctx sdk.Context, portfolio *types.InvestmentPortfolio, assetType string, targetPct sdk.Dec) error {
	// In production, this would execute actual trades/transfers
	// For now, we update the portfolio components
	
	targetAmount := targetPct.MulInt(portfolio.TotalValue.Amount)
	
	// Find or create component for this asset type
	found := false
	for i, component := range portfolio.Components {
		if component.AssetType == assetType {
			portfolio.Components[i].CurrentValue = sdk.NewCoin("unamo", targetAmount.TruncateInt())
			found = true
			break
		}
	}
	
	if !found {
		// Add new component
		newComponent := types.PortfolioComponent{
			AssetType:     assetType,
			Amount:        sdk.NewCoin("unamo", targetAmount.TruncateInt()),
			CurrentValue:  sdk.NewCoin("unamo", targetAmount.TruncateInt()),
			ReturnRate:    sdk.NewDecWithPrec(8, 2), // 8% default
			RiskRating:    "medium",
			LastUpdated:   ctx.BlockTime(),
		}
		portfolio.Components = append(portfolio.Components, newComponent)
	}
	
	return nil
}

func (k Keeper) calculatePortfolioRiskScore(portfolio *types.InvestmentPortfolio) int32 {
	// Calculate weighted risk score based on components
	totalWeight := sdk.ZeroDec()
	weightedRisk := sdk.ZeroDec()
	
	for _, component := range portfolio.Components {
		if !portfolio.TotalValue.Amount.IsZero() {
			weight := sdk.NewDecFromInt(component.CurrentValue.Amount).Quo(sdk.NewDecFromInt(portfolio.TotalValue.Amount))
			totalWeight = totalWeight.Add(weight)
			
			// Map risk rating to numeric score
			riskScore := sdk.NewDec(3) // medium default
			switch component.RiskRating {
			case "low":
				riskScore = sdk.NewDec(2)
			case "medium":
				riskScore = sdk.NewDec(3)
			case "high":
				riskScore = sdk.NewDec(4)
			}
			
			weightedRisk = weightedRisk.Add(weight.Mul(riskScore))
		}
	}
	
	if totalWeight.IsZero() {
		return 3 // Default medium risk
	}
	
	return int32(weightedRisk.Quo(totalWeight).TruncateInt64())
}

func (k Keeper) calculateExpectedReturn(portfolio *types.InvestmentPortfolio) sdk.Dec {
	// Calculate weighted expected return
	totalWeight := sdk.ZeroDec()
	weightedReturn := sdk.ZeroDec()
	
	for _, component := range portfolio.Components {
		if !portfolio.TotalValue.Amount.IsZero() {
			weight := sdk.NewDecFromInt(component.CurrentValue.Amount).Quo(sdk.NewDecFromInt(portfolio.TotalValue.Amount))
			totalWeight = totalWeight.Add(weight)
			weightedReturn = weightedReturn.Add(weight.Mul(component.ReturnRate))
		}
	}
	
	if totalWeight.IsZero() {
		return sdk.NewDecWithPrec(8, 2) // 8% default
	}
	
	return weightedReturn.Quo(totalWeight)
}

// ValidateMultiSignature validates multi-signature requirements for fund operations
func (k Keeper) ValidateMultiSignature(ctx sdk.Context, signers []string) bool {
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return false
	}
	
	// Check minimum signatures requirement
	if len(signers) < int(governance.RequiredSignatures) {
		return false
	}
	
	// Verify all signers are valid fund managers
	validSigners := 0
	signerMap := make(map[string]bool)
	
	for _, signer := range signers {
		// Prevent duplicate signers
		if signerMap[signer] {
			return false
		}
		signerMap[signer] = true
		
		// Check if signer is a fund manager
		for _, manager := range governance.FundManagers {
			if manager.Address == signer {
				validSigners++
				break
			}
		}
	}
	
	// Must have enough valid fund manager signatures
	return validSigners >= int(governance.RequiredSignatures)
}

// IsFundManager checks if an address is a fund manager
func (k Keeper) IsFundManager(ctx sdk.Context, address string) bool {
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return false
	}
	
	for _, manager := range governance.FundManagers {
		if manager.Address == address {
			return true
		}
	}
	
	return false
}

// ValidateRebalanceAuthority validates if the authority can perform portfolio rebalancing
func (k Keeper) ValidateRebalanceAuthority(ctx sdk.Context, authority string) error {
	// Check if authority is the governance module
	if authority == k.GetAuthority() {
		return nil
	}
	
	// Check if authority is an authorized fund manager with rebalancing privileges
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return types.ErrGovernanceNotFound
	}
	
	for _, manager := range governance.FundManagers {
		if manager.Address == authority {
			// Check if manager has rebalancing permissions
			// In production, this would check specific role permissions
			return nil
		}
	}
	
	return types.ErrUnauthorized.Wrapf("authority %s is not authorized for portfolio rebalancing", authority)
}

// ValidateGovernanceUpdate validates governance parameter updates
func (k Keeper) ValidateGovernanceUpdate(ctx sdk.Context, authority string) error {
	// Only governance module can update governance parameters
	if authority != k.GetAuthority() {
		return types.ErrUnauthorized.Wrapf("only governance module can update governance parameters, got %s", authority)
	}
	
	return nil
}

// UpdateAllocationStatus updates the status of an allocation
func (k Keeper) UpdateAllocationStatus(ctx sdk.Context, allocationID uint64, newStatus string) error {
	allocation, found := k.GetFundAllocation(ctx, allocationID)
	if !found {
		return types.ErrAllocationNotFound
	}
	
	// Remove old status index
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAllocationByStatusKey(allocation.Status, allocation.Id))
	
	// Update status
	allocation.Status = newStatus
	k.SetFundAllocation(ctx, allocation)
	
	return nil
}

// GetInvestmentPortfolio returns the current investment portfolio
func (k Keeper) GetInvestmentPortfolio(ctx sdk.Context) (types.InvestmentPortfolio, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InvestmentPortfolioKey)
	if bz == nil {
		return types.InvestmentPortfolio{}, false
	}
	
	var portfolio types.InvestmentPortfolio
	k.cdc.MustUnmarshal(bz, &portfolio)
	return portfolio, true
}

// SetInvestmentPortfolio sets the investment portfolio
func (k Keeper) SetInvestmentPortfolio(ctx sdk.Context, portfolio types.InvestmentPortfolio) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&portfolio)
	store.Set(types.InvestmentPortfolioKey, bz)
}

// GetFundGovernance returns the fund governance configuration
func (k Keeper) GetFundGovernance(ctx sdk.Context) (types.FundGovernance, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.FundGovernanceKey)
	if bz == nil {
		return types.FundGovernance{}, false
	}
	
	var governance types.FundGovernance
	k.cdc.MustUnmarshal(bz, &governance)
	return governance, true
}

// SetFundGovernance sets the fund governance configuration
func (k Keeper) SetFundGovernance(ctx sdk.Context, governance types.FundGovernance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&governance)
	store.Set(types.FundGovernanceKey, bz)
}

// IsFundManager checks if an address is a fund manager
func (k Keeper) IsFundManager(ctx sdk.Context, address string) bool {
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return false
	}
	
	for _, manager := range governance.FundManagers {
		if manager == address {
			return true
		}
	}
	return false
}

// SetMonthlyReport sets a monthly report
func (k Keeper) SetMonthlyReport(ctx sdk.Context, report types.MonthlyReport) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&report)
	store.Set(types.GetMonthlyReportKey(report.Period), bz)
}

// GetMonthlyReport returns a monthly report by period
func (k Keeper) GetMonthlyReport(ctx sdk.Context, period string) (types.MonthlyReport, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMonthlyReportKey(period))
	if bz == nil {
		return types.MonthlyReport{}, false
	}
	
	var report types.MonthlyReport
	k.cdc.MustUnmarshal(bz, &report)
	return report, true
}

// RecordDSWFRevenueActivity records DSWF-related revenue activities
func (k Keeper) RecordDSWFRevenueActivity(ctx sdk.Context, activityType string, amount sdk.Coin) {
	// Only record if revenue module is enabled
	if k.revenueKeeper != nil && k.revenueKeeper.IsRevenueEnabled(ctx) {
		// Create revenue stream record for DSWF activities
		stream := map[string]interface{}{
			"id":            fmt.Sprintf("dswf_%s_%d", activityType, ctx.BlockHeight()),
			"name":          fmt.Sprintf("DSWF %s", activityType),
			"type":          "dswf_activity",
			"amount":        amount,
			"source_module": "dswf",
			"activity_type": activityType,
			"block_height":  ctx.BlockHeight(),
			"timestamp":     ctx.BlockTime(),
		}
		
		k.revenueKeeper.RecordRevenueStream(ctx, stream)
	}
}

// ReceiveRevenueAllocation handles receiving revenue from the revenue distribution system
func (k Keeper) ReceiveRevenueAllocation(ctx sdk.Context, amount sdk.Coin) error {
	// Verify the amount is positive
	if amount.Amount.LTE(sdk.ZeroInt()) {
		return types.ErrInvalidAmount.Wrap("revenue allocation amount must be positive")
	}
	
	// Update fund balance and portfolio
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if !found {
		// Initialize portfolio if doesn't exist
		portfolio = types.InvestmentPortfolio{
			TotalValue:       amount,
			LiquidAssets:     amount,
			InvestedAssets:   sdk.NewCoin(amount.Denom, sdk.ZeroInt()),
			ReservedAssets:   sdk.NewCoin(amount.Denom, sdk.ZeroInt()),
			Components:       []types.PortfolioComponent{},
			TotalReturns:     sdk.NewCoin(amount.Denom, sdk.ZeroInt()),
			AnnualReturnRate: sdk.NewDecWithPrec(8, 2), // 8% initial target
			RiskScore:        3, // Medium risk
			LastRebalanced:   ctx.BlockTime(),
		}
	} else {
		// Add to existing portfolio
		portfolio.TotalValue = portfolio.TotalValue.Add(amount)
		portfolio.LiquidAssets = portfolio.LiquidAssets.Add(amount)
	}
	
	k.SetInvestmentPortfolio(ctx, portfolio)
	
	// Record this revenue receipt
	k.RecordDSWFRevenueActivity(ctx, "revenue_received", amount)
	
	return nil
}

// ValidateAllocationProposalWithRecovery validates a fund allocation proposal with comprehensive error handling and recovery
func (k Keeper) ValidateAllocationProposalWithRecovery(ctx sdk.Context, amount sdk.Coin, category string) error {
	// Comprehensive validation with detailed error context
	defer func() {
		if r := recover(); r != nil {
			k.Logger(ctx).Error("panic during allocation proposal validation", "panic", r, "amount", amount.String(), "category", category)
		}
	}()
	
	return k.ValidateAllocationProposal(ctx, amount, category)
}

// ValidateSystemIntegrity performs comprehensive cross-module integrity checks
func (k Keeper) ValidateSystemIntegrity(ctx sdk.Context) error {
	// Check module account integrity
	moduleAddr := k.GetModuleAccountAddress()
	if moduleAddr.Empty() {
		return types.ErrInvalidProposal.Wrap("DSWF module account address is empty")
	}
	
	// Validate account balance consistency with portfolio
	moduleBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddr)
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if found && !moduleBalance.IsZero() {
		totalExpected := portfolio.TotalValue
		actualBalance := moduleBalance.AmountOf(totalExpected.Denom)
		
		// Allow 1% tolerance for rounding errors
		tolerance := actualBalance.Quo(sdk.NewInt(100))
		diff := actualBalance.Sub(totalExpected.Amount).Abs()
		
		if diff.GT(tolerance) {
			return types.ErrInvalidAmount.Wrapf("portfolio balance mismatch: expected %s, actual %s, diff %s", totalExpected.String(), actualBalance.String(), diff.String())
		}
	}
	
	// Check governance integrity
	governance, found := k.GetFundGovernance(ctx)
	if found {
		if len(governance.FundManagers) == 0 {
			return types.ErrGovernanceNotFound.Wrap("no fund managers configured")
		}
		
		if governance.RequiredSignatures > int32(len(governance.FundManagers)) {
			return types.ErrInvalidProposal.Wrap("required signatures exceeds number of fund managers")
		}
	}
	
	return nil
}

// EmergencyPause pauses all fund operations in emergency situations
func (k Keeper) EmergencyPause(ctx sdk.Context, authority string, reason string) error {
	// Verify emergency pause authority
	params := k.GetParams(ctx)
	authorized := false
	
	// Check if authority is the governance module
	if authority == k.GetAuthority() {
		authorized = true
	}
	
	// Check if authority is in emergency pause authorities list
	for _, addr := range params.EmergencyPauseAuthorities {
		if addr == authority {
			authorized = true
			break
		}
	}
	
	if !authorized {
		return types.ErrUnauthorized.Wrapf("authority %s not authorized for emergency pause", authority)
	}
	
	// Set emergency pause parameters
	params.Enabled = false
	params.EmergencyPauseReason = reason
	params.EmergencyPauseTimestamp = ctx.BlockTime()
	k.SetParams(ctx, params)
	
	// Log emergency pause
	k.Logger(ctx).Error("DSWF emergency pause activated", "authority", authority, "reason", reason, "timestamp", ctx.BlockTime())
	
	// Emit emergency event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"dswf_emergency_pause",
			sdk.NewAttribute("authority", authority),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	)
	
	return nil
}

// ValidateOperationalSafety performs pre-operation safety checks
func (k Keeper) ValidateOperationalSafety(ctx sdk.Context, operationType string, amount sdk.Coin) error {
	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrInvalidProposal.Wrapf("DSWF operations paused: %s", params.EmergencyPauseReason)
	}
	
	// Check system integrity before major operations
	if err := k.ValidateSystemIntegrity(ctx); err != nil {
		return types.ErrInvalidProposal.Wrapf("system integrity check failed: %v", err)
	}
	
	// Circuit breaker: Check if amount exceeds safety thresholds
	if operationType == "disbursement" || operationType == "allocation" {
		portfolio, found := k.GetInvestmentPortfolio(ctx)
		if found {
			// No single operation should exceed 10% of total portfolio
			maxSafeAmount := portfolio.TotalValue.Amount.Quo(sdk.NewInt(10))
			if amount.Amount.GT(maxSafeAmount) {
				return types.ErrAllocationLimitExceeded.Wrapf("operation amount %s exceeds safety limit %s (10%% of portfolio)", amount.String(), maxSafeAmount.String())
			}
		}
	}
	
	return nil
}

// ValidateAllocationProposal validates a fund allocation proposal
func (k Keeper) ValidateAllocationProposal(ctx sdk.Context, amount sdk.Coin, category string) error {
	params := k.GetParams(ctx)
	
	// Check if DSWF is enabled
	if !params.Enabled {
		return types.ErrInvalidProposal.Wrap("DSWF is not enabled")
	}
	
	// Validate amount is positive
	if amount.Amount.LTE(sdk.ZeroInt()) {
		return types.ErrInvalidAmount.Wrap("allocation amount must be positive")
	}
	
	// Check minimum allocation amount
	if amount.Amount.LT(params.MinFundBalance.Amount.Quo(sdk.NewInt(1000))) { // 0.1% of min fund balance
		return types.ErrInvalidAmount.Wrap("allocation amount too small")
	}
	
	// Get current fund balance and portfolio
	fundBalance := k.GetFundBalance(ctx)
	if len(fundBalance) == 0 {
		return types.ErrInsufficientFunds.Wrap("no funds available in DSWF")
	}
	
	totalFunds := fundBalance.AmountOf(amount.Denom)
	
	// Check minimum fund balance after allocation
	if totalFunds.Sub(amount.Amount).LT(params.MinFundBalance.Amount) {
		return types.ErrMinimumBalanceRequired.Wrapf("allocation would violate minimum fund balance requirement")
	}
	
	// Check allocation percentage limit
	maxAllocation := totalFunds.ToDec().Mul(params.MaxAllocationPercentage).TruncateInt()
	if amount.Amount.GT(maxAllocation) {
		return types.ErrAllocationLimitExceeded.Wrapf("allocation %s exceeds maximum allowed %s", amount, maxAllocation)
	}
	
	// Check portfolio liquidity constraints
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if found {
		// Ensure adequate liquidity remains
		availableLiquidity := portfolio.LiquidAssets.Amount
		if availableLiquidity.LT(amount.Amount) {
			return types.ErrInsufficientLiquidity.Wrapf("insufficient liquid assets: need %s, have %s", amount, availableLiquidity)
		}
		
		// Check minimum liquidity ratio after allocation
		remainingLiquidity := availableLiquidity.Sub(amount.Amount)
		liquidityRatio := sdk.NewDecFromInt(remainingLiquidity).Quo(sdk.NewDecFromInt(portfolio.TotalValue.Amount))
		if liquidityRatio.LT(params.MinLiquidityRatio) {
			return types.ErrInsufficientLiquidity.Wrapf("allocation would violate minimum liquidity ratio: %s < %s", liquidityRatio, params.MinLiquidityRatio)
		}
	}
	
	// Validate category
	governance, found := k.GetFundGovernance(ctx)
	if !found {
		return types.ErrInvalidProposal.Wrap("fund governance not initialized")
	}
	
	validCategory := false
	for _, cat := range governance.Categories {
		if cat.Name == category {
			validCategory = true
			break
		}
	}
	
	if !validCategory {
		return types.ErrInvalidCategory.Wrapf("category %s not allowed", category)
	}
	
	return nil
}