package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/deshchain/deshchain/x/donation/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace
	
	// Dependencies
	bankKeeper types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: ps,
		bankKeeper: bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// DistributeToNGOs distributes funds to NGO wallets based on configured percentages
func (k Keeper) DistributeToNGOs(ctx sdk.Context, amount sdk.Coins) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return nil
	}
	
	// Get active NGO wallets
	ngos := k.GetActiveNGOs(ctx)
	if len(ngos) == 0 {
		return types.ErrNoActiveNGOs
	}
	
	// Calculate distribution per NGO (equal distribution for now)
	// In production, this could be weighted by transparency score, impact metrics, etc.
	perNGOAmount := sdk.Coins{}
	for _, coin := range amount {
		perCoinAmount := coin.Amount.Quo(sdk.NewInt(int64(len(ngos))))
		if perCoinAmount.IsPositive() {
			perNGOAmount = perNGOAmount.Add(sdk.NewCoin(coin.Denom, perCoinAmount))
		}
	}
	
	// Distribute to each NGO
	for _, ngo := range ngos {
		ngoAddr, err := sdk.AccAddressFromBech32(ngo.Address)
		if err != nil {
			k.Logger(ctx).Error("invalid NGO address", "ngo", ngo.Name, "error", err)
			continue
		}
		
		// Send funds from donation module to NGO address
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, ngoAddr, perNGOAmount,
		); err != nil {
			k.Logger(ctx).Error("failed to send to NGO", "ngo", ngo.Name, "error", err)
			continue
		}
		
		// Update NGO stats
		k.UpdateNGOStats(ctx, ngo.Id, perNGOAmount)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDonationDistributed,
				sdk.NewAttribute(types.AttributeKeyNGOId, fmt.Sprintf("%d", ngo.Id)),
				sdk.NewAttribute(types.AttributeKeyNGOName, ngo.Name),
				sdk.NewAttribute(types.AttributeKeyAmount, perNGOAmount.String()),
			),
		)
	}
	
	return nil
}

// GetActiveNGOs returns all active and verified NGOs
func (k Keeper) GetActiveNGOs(ctx sdk.Context) []types.NGOWallet {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NGOWalletKeyPrefix)
	defer iterator.Close()
	
	var ngos []types.NGOWallet
	for ; iterator.Valid(); iterator.Next() {
		var ngo types.NGOWallet
		k.cdc.MustUnmarshal(iterator.Value(), &ngo)
		
		if ngo.IsActive && ngo.IsVerified {
			ngos = append(ngos, ngo)
		}
	}
	
	// If no NGOs in store, use default NGOs
	if len(ngos) == 0 {
		defaultNGOs := types.GetActiveNGOWallets()
		for _, ngo := range defaultNGOs {
			if ngo.IsVerified {
				ngos = append(ngos, ngo)
			}
		}
	}
	
	return ngos
}

// SetNGOWallet stores an NGO wallet
func (k Keeper) SetNGOWallet(ctx sdk.Context, ngo types.NGOWallet) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetNGOWalletKey(ngo.Id)
	value := k.cdc.MustMarshal(&ngo)
	store.Set(key, value)
}

// GetNGOWallet retrieves an NGO wallet by ID
func (k Keeper) GetNGOWallet(ctx sdk.Context, id uint64) (types.NGOWallet, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetNGOWalletKey(id)
	
	if !store.Has(key) {
		// Check default NGOs
		wallets := types.DefaultNGOWallets()
		for _, wallet := range wallets {
			if wallet.Id == id {
				return wallet, true
			}
		}
		return types.NGOWallet{}, false
	}
	
	var ngo types.NGOWallet
	k.cdc.MustUnmarshal(store.Get(key), &ngo)
	return ngo, true
}

// UpdateNGOStats updates the statistics for an NGO after a donation
func (k Keeper) UpdateNGOStats(ctx sdk.Context, ngoId uint64, amount sdk.Coins) {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return
	}
	
	// Update total received
	ngo.TotalReceived = ngo.TotalReceived.Add(amount...)
	
	// Update current balance
	ngo.CurrentBalance = ngo.CurrentBalance.Add(amount...)
	
	// Update timestamp
	ngo.UpdatedAt = ctx.BlockTime().Unix()
	
	// Save updated NGO
	k.SetNGOWallet(ctx, ngo)
}

// GetAllNGOWallets returns all NGO wallets
func (k Keeper) GetAllNGOWallets(ctx sdk.Context) []types.NGOWallet {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NGOWalletKeyPrefix)
	defer iterator.Close()
	
	var ngos []types.NGOWallet
	for ; iterator.Valid(); iterator.Next() {
		var ngo types.NGOWallet
		k.cdc.MustUnmarshal(iterator.Value(), &ngo)
		ngos = append(ngos, ngo)
	}
	
	// If no NGOs in store, return default NGOs
	if len(ngos) == 0 {
		return types.DefaultNGOWallets()
	}
	
	return ngos
}

// DirectDonation handles direct donations to a specific NGO
func (k Keeper) DirectDonation(ctx sdk.Context, from sdk.AccAddress, ngoId uint64, amount sdk.Coins) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return types.ErrModuleDisabled
	}
	
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	if !ngo.IsActive {
		return types.ErrNGONotActive
	}
	
	if !ngo.IsVerified {
		return types.ErrNGONotVerified
	}
	
	ngoAddr, err := sdk.AccAddressFromBech32(ngo.Address)
	if err != nil {
		return types.ErrInvalidNGOAddress
	}
	
	// Transfer funds directly to NGO
	if err := k.bankKeeper.SendCoins(ctx, from, ngoAddr, amount); err != nil {
		return err
	}
	
	// Update NGO stats
	k.UpdateNGOStats(ctx, ngoId, amount)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDirectDonation,
			sdk.NewAttribute(types.AttributeKeyDonor, from.String()),
			sdk.NewAttribute(types.AttributeKeyNGOId, fmt.Sprintf("%d", ngoId)),
			sdk.NewAttribute(types.AttributeKeyNGOName, ngo.Name),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)
	
	return nil
}