package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/kisaanmitra/types"
)

// Keeper of the kisaanmitra store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	memKey   sdk.StoreKey
}

// NewKeeper creates a new kisaanmitra Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetBorrower stores a borrower in the store
func (k Keeper) SetBorrower(ctx sdk.Context, borrower types.Borrower) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&borrower)
	store.Set(types.GetBorrowerKey(borrower.BorrowerID), bz)
}

// GetBorrower retrieves a borrower from the store
func (k Keeper) GetBorrower(ctx sdk.Context, borrowerID string) (borrower types.Borrower, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBorrowerKey(borrowerID))
	if bz == nil {
		return borrower, false
	}
	k.cdc.MustUnmarshal(bz, &borrower)
	return borrower, true
}

// SetLoan stores a loan in the store
func (k Keeper) SetLoan(ctx sdk.Context, loan types.Loan) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(types.GetLoanKey(loan.LoanID), bz)
}

// GetLoan retrieves a loan from the store
func (k Keeper) GetLoan(ctx sdk.Context, loanID string) (loan types.Loan, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLoanKey(loanID))
	if bz == nil {
		return loan, false
	}
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

// GetAllBorrowers retrieves all borrowers
func (k Keeper) GetAllBorrowers(ctx sdk.Context) []types.Borrower {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.BorrowerPrefix.Bytes(), nil)
	defer iterator.Close()

	var borrowers []types.Borrower
	for ; iterator.Valid(); iterator.Next() {
		var borrower types.Borrower
		k.cdc.MustUnmarshal(iterator.Value(), &borrower)
		borrowers = append(borrowers, borrower)
	}
	return borrowers
}

// GetAllLoans retrieves all loans
func (k Keeper) GetAllLoans(ctx sdk.Context) []types.Loan {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.LoanPrefix.Bytes(), nil)
	defer iterator.Close()

	var loans []types.Loan
	for ; iterator.Valid(); iterator.Next() {
		var loan types.Loan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		loans = append(loans, loan)
	}
	return loans
}