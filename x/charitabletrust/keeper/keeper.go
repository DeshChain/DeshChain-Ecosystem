package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/charitabletrust/types"
)

// Keeper of the x/charitabletrust store
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	stakingKeeper  types.StakingKeeper
	donationKeeper types.DonationKeeper
	govKeeper      types.GovKeeper
	revenueKeeper  types.RevenueKeeper

	// the address capable of executing governance actions
	authority string
}

// NewKeeper creates a new CharitableTrust Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	donationKeeper types.DonationKeeper,
	govKeeper types.GovKeeper,
	revenueKeeper types.RevenueKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		stakingKeeper:  stakingKeeper,
		donationKeeper: donationKeeper,
		govKeeper:      govKeeper,
		revenueKeeper:  revenueKeeper,
		authority:      authority,
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

// GetTrustFundBalance returns the current trust fund balance
func (k Keeper) GetTrustFundBalance(ctx sdk.Context) (types.TrustFundBalance, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TrustFundBalanceKey)
	if bz == nil {
		return types.TrustFundBalance{}, false
	}
	
	var balance types.TrustFundBalance
	k.cdc.MustUnmarshal(bz, &balance)
	return balance, true
}

// SetTrustFundBalance sets the trust fund balance
func (k Keeper) SetTrustFundBalance(ctx sdk.Context, balance types.TrustFundBalance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&balance)
	store.Set(types.TrustFundBalanceKey, bz)
}

// UpdateTrustFundBalance updates the trust fund balance
func (k Keeper) UpdateTrustFundBalance(ctx sdk.Context) {
	moduleBalance := k.bankKeeper.GetAllBalances(ctx, k.GetModuleAccountAddress())
	
	balance, _ := k.GetTrustFundBalance(ctx)
	if len(moduleBalance) > 0 {
		balance.TotalBalance = moduleBalance[0] // Assuming single denom
		balance.AvailableAmount = balance.TotalBalance.Sub(balance.AllocatedAmount)
	}
	balance.LastUpdated = ctx.BlockTime()
	
	k.SetTrustFundBalance(ctx, balance)
}

// GetAllocationCount returns the current allocation count
func (k Keeper) GetAllocationCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AllocationCountKey)
	if bz == nil {
		return 0
	}
	return types.GetUint64FromBytes(bz)
}

// SetAllocationCount sets the allocation count
func (k Keeper) SetAllocationCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AllocationCountKey, types.GetUint64Bytes(count))
}

// IncrementAllocationCount increments and returns the allocation count
func (k Keeper) IncrementAllocationCount(ctx sdk.Context) uint64 {
	count := k.GetAllocationCount(ctx)
	count++
	k.SetAllocationCount(ctx, count)
	return count
}

// GetProposalCount returns the current proposal count
func (k Keeper) GetProposalCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalCountKey)
	if bz == nil {
		return 0
	}
	return types.GetUint64FromBytes(bz)
}

// SetProposalCount sets the proposal count
func (k Keeper) SetProposalCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalCountKey, types.GetUint64Bytes(count))
}

// IncrementProposalCount increments and returns the proposal count
func (k Keeper) IncrementProposalCount(ctx sdk.Context) uint64 {
	count := k.GetProposalCount(ctx)
	count++
	k.SetProposalCount(ctx, count)
	return count
}

// SetCharitableAllocation sets a charitable allocation
func (k Keeper) SetCharitableAllocation(ctx sdk.Context, allocation types.CharitableAllocation) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&allocation)
	store.Set(types.GetCharitableAllocationKey(allocation.Id), bz)
	
	// Set indexes
	store.Set(types.GetAllocationByCategoryKey(allocation.Category, allocation.Id), []byte{})
	store.Set(types.GetAllocationByStatusKey(allocation.Status, allocation.Id), []byte{})
	store.Set(types.GetAllocationByOrgKey(allocation.CharitableOrgWalletId, allocation.Id), []byte{})
}

// GetCharitableAllocation returns a charitable allocation by ID
func (k Keeper) GetCharitableAllocation(ctx sdk.Context, allocationID uint64) (types.CharitableAllocation, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCharitableAllocationKey(allocationID))
	if bz == nil {
		return types.CharitableAllocation{}, false
	}
	
	var allocation types.CharitableAllocation
	k.cdc.MustUnmarshal(bz, &allocation)
	return allocation, true
}

// SetAllocationProposal sets an allocation proposal
func (k Keeper) SetAllocationProposal(ctx sdk.Context, proposal types.AllocationProposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&proposal)
	store.Set(types.GetAllocationProposalKey(proposal.Id), bz)
	
	// Set index
	store.Set(types.GetProposalByStatusKey(proposal.Status, proposal.Id), []byte{})
}

// GetAllocationProposal returns an allocation proposal by ID
func (k Keeper) GetAllocationProposal(ctx sdk.Context, proposalID uint64) (types.AllocationProposal, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAllocationProposalKey(proposalID))
	if bz == nil {
		return types.AllocationProposal{}, false
	}
	
	var proposal types.AllocationProposal
	k.cdc.MustUnmarshal(bz, &proposal)
	return proposal, true
}

// GetTrustGovernance returns the trust governance configuration
func (k Keeper) GetTrustGovernance(ctx sdk.Context) (types.TrustGovernance, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TrustGovernanceKey)
	if bz == nil {
		return types.TrustGovernance{}, false
	}
	
	var governance types.TrustGovernance
	k.cdc.MustUnmarshal(bz, &governance)
	return governance, true
}

// SetTrustGovernance sets the trust governance configuration
func (k Keeper) SetTrustGovernance(ctx sdk.Context, governance types.TrustGovernance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&governance)
	store.Set(types.TrustGovernanceKey, bz)
}

// IsTrustee checks if an address is a trustee
func (k Keeper) IsTrustee(ctx sdk.Context, address string) bool {
	governance, found := k.GetTrustGovernance(ctx)
	if !found {
		return false
	}
	
	for _, trustee := range governance.Trustees {
		if trustee.Address == address && trustee.IsActive {
			return true
		}
	}
	return false
}

// SetImpactReport sets an impact report
func (k Keeper) SetImpactReport(ctx sdk.Context, report types.ImpactReport) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&report)
	store.Set(types.GetImpactReportKey(report.Id), bz)
	
	// Set index
	store.Set(types.GetReportByAllocationKey(report.AllocationId, report.Id), []byte{})
}

// GetImpactReport returns an impact report by ID
func (k Keeper) GetImpactReport(ctx sdk.Context, reportID uint64) (types.ImpactReport, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetImpactReportKey(reportID))
	if bz == nil {
		return types.ImpactReport{}, false
	}
	
	var report types.ImpactReport
	k.cdc.MustUnmarshal(bz, &report)
	return report, true
}

// SetFraudAlert sets a fraud alert
func (k Keeper) SetFraudAlert(ctx sdk.Context, alert types.FraudAlert) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&alert)
	store.Set(types.GetFraudAlertKey(alert.Id), bz)
	
	// Set allocation index
	store.Set(types.GetAlertByAllocationKey(alert.AllocationId, alert.Id), []byte{})
	
	// Set organization index for faster fraud status lookups
	// Get allocation to determine organization ID
	allocation, found := k.GetCharitableAllocation(ctx, alert.AllocationId)
	if found {
		store.Set(types.GetAlertByOrgKey(allocation.CharitableOrgWalletId, alert.Id), bz)
	}
}

// GetFraudAlert returns a fraud alert by ID
func (k Keeper) GetFraudAlert(ctx sdk.Context, alertID uint64) (types.FraudAlert, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFraudAlertKey(alertID))
	if bz == nil {
		return types.FraudAlert{}, false
	}
	
	var alert types.FraudAlert
	k.cdc.MustUnmarshal(bz, &alert)
	return alert, true
}

// ValidateCharitableOrganization validates if an organization can receive funds
func (k Keeper) ValidateCharitableOrganization(ctx sdk.Context, orgWalletID uint64) error {
	// Check if organization exists and is verified
	if !k.donationKeeper.IsWalletVerified(ctx, orgWalletID) {
		return types.ErrOrganizationNotVerified
	}
	
	// Check if organization is active
	if !k.donationKeeper.IsWalletActive(ctx, orgWalletID) {
		return types.ErrOrganizationInactive
	}
	
	// Check if organization is under investigation
	alerts := k.GetFraudAlertsByOrganization(ctx, orgWalletID)
	for _, alert := range alerts {
		if alert.Status == "investigating" || alert.Status == "pending" {
			return types.ErrOrganizationUnderInvestigation.Wrapf("organization %d is under fraud investigation", orgWalletID)
		}
	}
	
	return nil
}

// ValidateOrganizationFraudStatus validates if organization is clear of fraud issues
func (k Keeper) ValidateOrganizationFraudStatus(ctx sdk.Context, orgWalletID uint64) error {
	// Get all fraud alerts for this organization
	alerts := k.GetFraudAlertsByOrganization(ctx, orgWalletID)
	
	for _, alert := range alerts {
		switch alert.Status {
		case "investigating", "pending":
			return types.ErrOrganizationUnderInvestigation.Wrapf("organization %d has active fraud investigation (Alert ID: %d)", orgWalletID, alert.Id)
		case "confirmed":
			return types.ErrOrganizationFraudulent.Wrapf("organization %d has confirmed fraud case (Alert ID: %d)", orgWalletID, alert.Id)
		case "suspended":
			return types.ErrOrganizationSuspended.Wrapf("organization %d is suspended due to fraud concerns (Alert ID: %d)", orgWalletID, alert.Id)
		}
	}
	
	// Check recent investigation history
	recentInvestigations := k.GetRecentInvestigationsByOrganization(ctx, orgWalletID, 90) // Last 90 days
	if len(recentInvestigations) >= 3 {
		return types.ErrOrganizationHighRisk.Wrapf("organization %d has too many recent investigations (%d in last 90 days)", orgWalletID, len(recentInvestigations))
	}
	
	return nil
}

// GetFraudAlertsByOrganization returns all fraud alerts for an organization
func (k Keeper) GetFraudAlertsByOrganization(ctx sdk.Context, orgWalletID uint64) []types.FraudAlert {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetAlertByOrgPrefix(orgWalletID))
	defer iterator.Close()
	
	var alerts []types.FraudAlert
	for ; iterator.Valid(); iterator.Next() {
		var alert types.FraudAlert
		k.cdc.MustUnmarshal(iterator.Value(), &alert)
		alerts = append(alerts, alert)
	}
	
	// If no alerts found via direct organization index, search through all allocations for this organization
	if len(alerts) == 0 {
		allocations := k.GetAllCharitableAllocations(ctx)
		for _, allocation := range allocations {
			if allocation.CharitableOrgWalletId == orgWalletID {
				allocationAlerts := k.GetFraudAlertsByAllocation(ctx, allocation.Id)
				alerts = append(alerts, allocationAlerts...)
			}
		}
	}
	
	return alerts
}

// GetFraudAlertsByAllocation returns all fraud alerts for a specific allocation
func (k Keeper) GetFraudAlertsByAllocation(ctx sdk.Context, allocationID uint64) []types.FraudAlert {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.AlertByAllocationKey, types.GetUint64Bytes(allocationID)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var alerts []types.FraudAlert
	for ; iterator.Valid(); iterator.Next() {
		// Extract alert ID from the key and get the full alert
		alertIDBytes := iterator.Key()[len(prefix):]
		if len(alertIDBytes) == 8 {
			alertID := types.GetUint64FromBytes(alertIDBytes)
			alert, found := k.GetFraudAlert(ctx, alertID)
			if found {
				alerts = append(alerts, alert)
			}
		}
	}
	
	return alerts
}

// GetAllCharitableAllocations returns all charitable allocations (needed for fallback search)
func (k Keeper) GetAllCharitableAllocations(ctx sdk.Context) []types.CharitableAllocation {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CharitableAllocationKey)
	defer iterator.Close()
	
	var allocations []types.CharitableAllocation
	for ; iterator.Valid(); iterator.Next() {
		var allocation types.CharitableAllocation
		k.cdc.MustUnmarshal(iterator.Value(), &allocation)
		allocations = append(allocations, allocation)
	}
	
	return allocations
}

// GetRecentInvestigationsByOrganization returns recent investigations for an organization
func (k Keeper) GetRecentInvestigationsByOrganization(ctx sdk.Context, orgWalletID uint64, days int32) []types.FraudAlert {
	alerts := k.GetFraudAlertsByOrganization(ctx, orgWalletID)
	cutoffDate := ctx.BlockTime().AddDate(0, 0, -int(days))
	
	var recentInvestigations []types.FraudAlert
	for _, alert := range alerts {
		if alert.ReportedAt.After(cutoffDate) && 
		   (alert.Status == "investigating" || alert.Status == "investigated" || alert.Status == "confirmed") {
			recentInvestigations = append(recentInvestigations, alert)
		}
	}
	
	return recentInvestigations
}

// CalculateOrganizationMonthlyAllocation calculates total monthly allocation for an org
func (k Keeper) CalculateOrganizationMonthlyAllocation(ctx sdk.Context, orgWalletID uint64) sdk.Coin {
	// Calculate monthly allocation for organization
	currentMonth := ctx.BlockTime().Month()
	currentYear := ctx.BlockTime().Year()
	
	totalAllocation := sdk.ZeroInt()
	allocations := k.GetAllCharitableAllocations(ctx)
	
	for _, allocation := range allocations {
		// Check if allocation is from current month
		if allocation.AllocatedAt.Month() == currentMonth && 
		   allocation.AllocatedAt.Year() == currentYear && 
		   allocation.CharitableOrgWalletId == orgWalletID {
			totalAllocation = totalAllocation.Add(allocation.Amount.Amount)
		}
	}
	
	return sdk.NewCoin("unamo", totalAllocation)
}

// GetOrganizationWalletAddress gets wallet address for organization from donation module
func (k Keeper) GetOrganizationWalletAddress(ctx sdk.Context, orgWalletID uint64) (sdk.AccAddress, error) {
	// In production, this would query the donation module
	// For now, return a deterministic address based on ID
	if orgWalletID == 0 {
		return nil, types.ErrInvalidOrganization.Wrap("organization wallet ID cannot be zero")
	}
	
	// Generate deterministic address based on organization ID
	orgBytes := sdk.Uint64ToBigEndian(orgWalletID)
	prefix := []byte("deshchain-charity-org-")
	addrBytes := append(prefix, orgBytes...)
	
	// Truncate to 20 bytes for address
	if len(addrBytes) > 20 {
		addrBytes = addrBytes[:20]
	} else if len(addrBytes) < 20 {
		// Pad with zeros
		padding := make([]byte, 20-len(addrBytes))
		addrBytes = append(addrBytes, padding...)
	}
	
	return sdk.AccAddress(addrBytes), nil
}

// ValidateOrganizationRepresentative validates if submitter is authorized for organization
func (k Keeper) ValidateOrganizationRepresentative(ctx sdk.Context, orgID uint64, submitter string) error {
	submitterAddr, err := sdk.AccAddressFromBech32(submitter)
	if err != nil {
		return types.ErrInvalidAddress.Wrap("invalid submitter address")
	}
	
	// In production, this would:
	// 1. Query the donation module for organization representatives
	// 2. Check digital signatures and authorization
	// 3. Verify role-based permissions
	
	// For now, basic validation
	if len(submitterAddr) != 20 {
		return types.ErrUnauthorized.Wrap("invalid submitter address format")
	}
	
	// In production, integrate with donation module:
	// return k.donationKeeper.IsAuthorizedRepresentative(ctx, orgID, submitterAddr)
	
	return nil
}

// RecordCharitableTrustRevenueActivity records CharitableTrust-related revenue activities
func (k Keeper) RecordCharitableTrustRevenueActivity(ctx sdk.Context, activityType string, amount sdk.Coin) {
	// Only record if revenue module is enabled
	if k.revenueKeeper != nil && k.revenueKeeper.IsRevenueEnabled(ctx) {
		// Create revenue stream record for CharitableTrust activities
		stream := map[string]interface{}{
			"id":            fmt.Sprintf("charitabletrust_%s_%d", activityType, ctx.BlockHeight()),
			"name":          fmt.Sprintf("CharitableTrust %s", activityType),
			"type":          "charitabletrust_activity",
			"amount":        amount,
			"source_module": "charitabletrust",
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
	
	// Update trust fund balance
	balance, found := k.GetTrustFundBalance(ctx)
	if !found {
		// Initialize balance if doesn't exist
		balance = types.TrustFundBalance{
			TotalBalance:     amount,
			AllocatedAmount:  sdk.NewCoin(amount.Denom, sdk.ZeroInt()),
			AvailableAmount:  amount,
			TotalDistributed: sdk.NewCoin(amount.Denom, sdk.ZeroInt()),
			LastUpdated:      ctx.BlockTime(),
		}
	} else {
		// Add to existing balance
		balance.TotalBalance = balance.TotalBalance.Add(amount)
		balance.AvailableAmount = balance.AvailableAmount.Add(amount)
		balance.LastUpdated = ctx.BlockTime()
	}
	
	k.SetTrustFundBalance(ctx, balance)
	
	// Record this revenue receipt
	k.RecordCharitableTrustRevenueActivity(ctx, "revenue_received", amount)
	
	return nil
}

// ValidateAllocationAmountWithRecovery validates allocation amount with comprehensive error handling
func (k Keeper) ValidateAllocationAmountWithRecovery(ctx sdk.Context, amount sdk.Coin) error {
	// Comprehensive validation with panic recovery
	defer func() {
		if r := recover(); r != nil {
			k.Logger(ctx).Error("panic during allocation amount validation", "panic", r, "amount", amount.String())
		}
	}()
	
	return k.ValidateAllocationAmount(ctx, amount)
}

// ValidateSystemHealth performs comprehensive system health checks
func (k Keeper) ValidateSystemHealth(ctx sdk.Context) error {
	// Check module account health
	moduleAddr := k.GetModuleAccountAddress()
	moduleBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddr)
	if moduleBalance.IsZero() {
		k.Logger(ctx).Warn("charitable trust module account has zero balance", "address", moduleAddr.String())
	}
	
	// Check trust fund balance consistency
	balance, found := k.GetTrustFundBalance(ctx)
	if found {
		if balance.TotalBalance.Amount.LT(balance.AllocatedAmount.Amount.Add(balance.AvailableAmount.Amount)) {
			return types.ErrTrustFundNotInitialized.Wrap("trust fund balance inconsistency detected")
		}
	}
	
	// Check governance configuration
	governance, found := k.GetTrustGovernance(ctx)
	if found {
		if err := governance.Validate(); err != nil {
			return types.ErrInvalidProposal.Wrapf("governance configuration validation failed: %v", err)
		}
	}
	
	return nil
}

// EmergencyPause pauses all charitable trust operations in emergency situations
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
	
	// Set emergency pause
	params.Enabled = false
	k.SetParams(ctx, params)
	
	// Log emergency pause
	k.Logger(ctx).Error("CharitableTrust emergency pause activated", "authority", authority, "reason", reason, "timestamp", ctx.BlockTime())
	
	// Emit emergency event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"charitabletrust_emergency_pause",
			sdk.NewAttribute("authority", authority),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	)
	
	return nil
}

// ValidateOperationalSafety performs pre-operation safety checks for charitable operations
func (k Keeper) ValidateOperationalSafety(ctx sdk.Context, operationType string, amount sdk.Coin) error {
	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled.Wrap("CharitableTrust operations paused due to emergency")
	}
	
	// Check system health before major operations
	if err := k.ValidateSystemHealth(ctx); err != nil {
		return types.ErrTrustFundNotInitialized.Wrapf("system health check failed: %v", err)
	}
	
	// Circuit breaker: Check if amount exceeds safety thresholds
	if operationType == "distribution" || operationType == "allocation" {
		balance, found := k.GetTrustFundBalance(ctx)
		if found {
			// No single operation should exceed 25% of available funds
			maxSafeAmount := balance.AvailableAmount.Amount.Quo(sdk.NewInt(4))
			if amount.Amount.GT(maxSafeAmount) {
				return types.ErrInsufficientFunds.Wrapf("operation amount %s exceeds safety limit %s (25%% of available funds)", amount.String(), maxSafeAmount.String())
			}
		}
	}
	
	return nil
}

// ValidateAllocationAmount validates allocation amount against available funds
func (k Keeper) ValidateAllocationAmount(ctx sdk.Context, amount sdk.Coin) error {
	params := k.GetParams(ctx)
	
	// Check if module is enabled
	if !params.Enabled {
		return types.ErrModuleDisabled.Wrap("CharitableTrust module is disabled")
	}
	
	// Validate amount is positive
	if amount.Amount.LTE(sdk.ZeroInt()) {
		return types.ErrInvalidAmount.Wrap("allocation amount must be positive")
	}
	
	// Check minimum allocation amount
	if amount.Amount.LT(params.MinAllocationAmount.Amount) {
		return types.ErrInvalidAmount.Wrapf("allocation %s below minimum %s", amount, params.MinAllocationAmount)
	}
	
	// Get trust fund balance
	balance, found := k.GetTrustFundBalance(ctx)
	if !found {
		return types.ErrTrustFundNotInitialized.Wrap("trust fund balance not initialized")
	}
	
	// Check available funds
	if balance.AvailableAmount.Amount.LT(amount.Amount) {
		return types.ErrInsufficientFunds.Wrapf("insufficient available funds: requested %s, available %s", amount, balance.AvailableAmount)
	}
	
	// Check if allocation would leave minimum operational balance
	minOperationalBalance := balance.TotalBalance.Amount.Quo(sdk.NewInt(100)) // 1% of total balance
	if balance.AvailableAmount.Amount.Sub(amount.Amount).LT(minOperationalBalance) {
		return types.ErrInsufficientFunds.Wrapf("allocation would violate minimum operational balance requirement")
	}
	
	return nil
}

// ValidateTrusteeQuorum validates if enough trustees have voted on a proposal
func (k Keeper) ValidateTrusteeQuorum(ctx sdk.Context, votes []types.Vote) (bool, bool) {
	governance, found := k.GetTrustGovernance(ctx)
	if !found {
		return false, false
	}
	
	// Count valid votes from active trustees
	validVotes := 0
	yesVotes := 0
	totalVotingPower := int32(0)
	yesVotingPower := int32(0)
	
	// Get total active trustees voting power
	for _, trustee := range governance.Trustees {
		if trustee.Status == "active" && ctx.BlockTime().Before(trustee.TermEndDate) {
			totalVotingPower += trustee.VotingPower
		}
	}
	
	// Count votes from valid trustees
	for _, vote := range votes {
		for _, trustee := range governance.Trustees {
			if trustee.Address == vote.Voter && 
			   trustee.Status == "active" && 
			   ctx.BlockTime().Before(trustee.TermEndDate) {
				validVotes++
				if vote.VoteType == "yes" {
					yesVotes++
					yesVotingPower += trustee.VotingPower
				}
				break
			}
		}
	}
	
	// Check quorum
	hasQuorum := validVotes >= int(governance.Quorum)
	
	// Check approval threshold
	approved := false
	if hasQuorum && totalVotingPower > 0 {
		approvalPct := sdk.NewDec(int64(yesVotingPower)).Quo(sdk.NewDec(int64(totalVotingPower)))
		approved = approvalPct.GTE(governance.ApprovalThreshold)
	}
	
	return hasQuorum, approved
}

// ValidateGovernanceAuthority validates if authority can update governance settings
func (k Keeper) ValidateGovernanceAuthority(ctx sdk.Context, authority string) error {
	// Only the governance module can update trust governance
	if authority != k.GetAuthority() {
		return types.ErrUnauthorized.Wrapf("only governance module can update trust governance, got %s", authority)
	}
	
	return nil
}

// ValidateTrusteeAuthority validates if address is an active trustee
func (k Keeper) ValidateTrusteeAuthority(ctx sdk.Context, address string) error {
	if !k.IsTrustee(ctx, address) {
		return types.ErrNotTrustee.Wrapf("address %s is not an active trustee", address)
	}
	
	return nil
}