package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/deshchain/deshchain/x/liquiditymanager/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	bankKeeper     types.BankKeeper
	stakingKeeper  types.StakingKeeper
	surakshaKeeper types.SurakshaKeeper
	dexKeeper      types.DEXKeeper
}

// NewKeeper creates a new liquidity manager keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	surakshaKeeper types.SurakshaKeeper,
	dexKeeper types.DEXKeeper,
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
		bankKeeper:     bankKeeper,
		stakingKeeper:  stakingKeeper,
		surakshaKeeper: surakshaKeeper,
		dexKeeper:      dexKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) sdk.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetParams get all parameters as types.LiquidityParams
func (k Keeper) GetParams(ctx sdk.Context) types.LiquidityParams {
	var params types.LiquidityParams
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.LiquidityParams) {
	k.paramstore.SetParamSet(ctx, &params)
}

// REVOLUTIONARY LENDING SYSTEM FUNCTIONS

// ProcessLoan processes a loan request with revolutionary restrictions
func (k Keeper) ProcessLoan(ctx sdk.Context, borrower sdk.AccAddress, amount sdk.Dec, module string) error {
	// Check if loan can be processed with member-only restrictions
	canProcess, message := k.CanProcessLoan(ctx, amount, module, borrower)
	if !canProcess {
		return types.ErrLoanRejected.Wrap(message)
	}

	// Record the loan processing
	k.recordLoanProcessing(ctx, borrower, amount, module)

	return nil
}

// ProcessCollateralLoan processes a NAMO collateral loan with 70% LTV
func (k Keeper) ProcessCollateralLoan(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) error {
	// Check if collateral loan can be processed
	canProcess, message := k.CanProcessCollateralLoan(ctx, loanAmount, collateralAmount, borrower)
	if !canProcess {
		return types.ErrCollateralLoanRejected.Wrap(message)
	}

	// Lock the collateral
	if err := k.LockCollateral(ctx, borrower, collateralAmount); err != nil {
		return types.ErrCollateralLockFailed.Wrap(err.Error())
	}

	// Record the collateral loan
	k.recordCollateralLoan(ctx, borrower, loanAmount, collateralAmount)

	return nil
}

// RepayCollateralLoan repays a collateral loan and unlocks NAMO
func (k Keeper) RepayCollateralLoan(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) error {
	// Unlock the collateral
	if err := k.UnlockCollateral(ctx, borrower, collateralAmount); err != nil {
		return types.ErrCollateralUnlockFailed.Wrap(err.Error())
	}

	// Record the repayment
	k.recordLoanRepayment(ctx, borrower, loanAmount, collateralAmount)

	return nil
}

// Helper functions for loan processing
func (k Keeper) recordLoanProcessing(ctx sdk.Context, borrower sdk.AccAddress, amount sdk.Dec, module string) {
	// Record loan in store for tracking
	store := ctx.KVStore(k.storeKey)
	
	// Update daily lending used
	currentUsed := k.getDailyLendingUsed(ctx)
	newUsed := currentUsed.Add(amount)
	
	key := []byte("daily_lending_used")
	bz := k.cdc.MustMarshal(&newUsed)
	store.Set(key, bz)
}

func (k Keeper) recordCollateralLoan(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) {
	// Record collateral loan details
	store := ctx.KVStore(k.storeKey)
	
	loan := types.CollateralLoan{
		Borrower:         borrower.String(),
		LoanAmount:       loanAmount,
		CollateralAmount: collateralAmount,
		Timestamp:        ctx.BlockTime(),
		IsActive:         true,
	}
	
	key := append([]byte("collateral_loan_"), borrower.Bytes()...)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(key, bz)
}

func (k Keeper) recordLoanRepayment(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) {
	// Mark loan as repaid
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("collateral_loan_"), borrower.Bytes()...)
	
	bz := store.Get(key)
	if bz != nil {
		var loan types.CollateralLoan
		k.cdc.MustUnmarshal(bz, &loan)
		loan.IsActive = false
		
		bz = k.cdc.MustMarshal(&loan)
		store.Set(key, bz)
	}
}