package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
)

// LockCollateral locks user's collateral for DINR minting
func (k Keeper) LockCollateral(ctx sdk.Context, user sdk.AccAddress, collateral sdk.Coins) error {
	// Transfer collateral from user to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		user,
		types.ModuleName,
		collateral,
	)
	if err != nil {
		return err
	}
	
	// Update user's collateral record
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("collateral/"), user.Bytes()...)
	
	// Get existing collateral
	var existingCollateral sdk.Coins
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &existingCollateral)
	}
	
	// Add new collateral
	totalCollateral := existingCollateral.Add(collateral...)
	
	// Store updated collateral
	store.Set(key, k.cdc.MustMarshal(&totalCollateral))
	
	return nil
}

// ReleaseCollateral releases user's locked collateral
func (k Keeper) ReleaseCollateral(ctx sdk.Context, user sdk.AccAddress, amount sdk.Coins) error {
	// Get user's collateral
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("collateral/"), user.Bytes()...)
	
	var userCollateral sdk.Coins
	bz := store.Get(key)
	if bz == nil {
		return types.ErrInsufficientCollateral
	}
	k.cdc.MustUnmarshal(bz, &userCollateral)
	
	// Check sufficient collateral
	if !userCollateral.IsAllGTE(amount) {
		return types.ErrInsufficientCollateral
	}
	
	// Transfer collateral from module to user
	err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		user,
		amount,
	)
	if err != nil {
		return err
	}
	
	// Update user's collateral record
	remainingCollateral := userCollateral.Sub(amount)
	if remainingCollateral.IsZero() {
		store.Delete(key)
	} else {
		store.Set(key, k.cdc.MustMarshal(&remainingCollateral))
	}
	
	return nil
}

// GetUserCollateral returns user's locked collateral
func (k Keeper) GetUserCollateral(ctx sdk.Context, user sdk.AccAddress) (sdk.Coins, error) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("collateral/"), user.Bytes()...)
	
	bz := store.Get(key)
	if bz == nil {
		return sdk.Coins{}, nil
	}
	
	var collateral sdk.Coins
	k.cdc.MustUnmarshal(bz, &collateral)
	return collateral, nil
}

// GetUserMintedDINR returns total DINR minted by user
func (k Keeper) GetUserMintedDINR(ctx sdk.Context, user sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("minted/"), user.Bytes()...)
	
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroInt()
	}
	
	var amount sdk.Int
	k.cdc.MustUnmarshal(bz, &amount)
	return amount
}

// SetUserMintedDINR updates user's total minted DINR
func (k Keeper) SetUserMintedDINR(ctx sdk.Context, user sdk.AccAddress, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("minted/"), user.Bytes()...)
	
	if amount.IsZero() {
		store.Delete(key)
	} else {
		store.Set(key, k.cdc.MustMarshal(&amount))
	}
}

// UpdateUserMintedDINR adds to user's total minted DINR
func (k Keeper) UpdateUserMintedDINR(ctx sdk.Context, user sdk.AccAddress, additionalAmount sdk.Int) {
	currentAmount := k.GetUserMintedDINR(ctx, user)
	newAmount := currentAmount.Add(additionalAmount)
	k.SetUserMintedDINR(ctx, user, newAmount)
}