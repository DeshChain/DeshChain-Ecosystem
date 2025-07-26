package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/deshchain/deshchain/x/charitabletrust/types"
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
	
	// Set index
	store.Set(types.GetAlertByAllocationKey(alert.AllocationId, alert.Id), []byte{})
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
	
	// TODO: Check if organization is under investigation
	
	return nil
}

// CalculateOrganizationMonthlyAllocation calculates total monthly allocation for an org
func (k Keeper) CalculateOrganizationMonthlyAllocation(ctx sdk.Context, orgWalletID uint64) sdk.Coin {
	// TODO: Implement monthly calculation logic
	// This would iterate through allocations for the current month
	return sdk.NewCoin("unamo", sdk.ZeroInt())
}