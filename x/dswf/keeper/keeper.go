package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/deshchain/deshchain/x/dswf/types"
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

// ValidateAllocationProposal validates a fund allocation proposal
func (k Keeper) ValidateAllocationProposal(ctx sdk.Context, amount sdk.Coin, category string) error {
	params := k.GetParams(ctx)
	
	// Check if DSWF is enabled
	if !params.Enabled {
		return types.ErrInvalidProposal.Wrap("DSWF is not enabled")
	}
	
	// Check minimum fund balance
	fundBalance := k.GetFundBalance(ctx)
	if fundBalance.AmountOf(amount.Denom).LT(params.MinimumFundBalance.Amount) {
		return types.ErrMinimumBalanceRequired
	}
	
	// Check allocation percentage limit
	maxAllocation := fundBalance.AmountOf(amount.Denom).
		ToDec().
		Mul(params.MaxAllocationPercentage).
		TruncateInt()
	
	if amount.Amount.GT(maxAllocation) {
		return types.ErrAllocationLimitExceeded
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
		return types.ErrInvalidCategory
	}
	
	return nil
}