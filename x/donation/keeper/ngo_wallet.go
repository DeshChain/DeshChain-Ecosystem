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
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/donation/types"
)

// SetNGOWallet sets an NGO wallet in the store
func (k Keeper) SetNGOWallet(ctx sdk.Context, ngo types.NGOWallet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOWalletKey)
	bz := k.cdc.MustMarshal(&ngo)
	store.Set(sdk.Uint64ToBigEndian(ngo.Id), bz)
}

// GetNGOWallet retrieves an NGO wallet from the store
func (k Keeper) GetNGOWallet(ctx sdk.Context, id uint64) (types.NGOWallet, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOWalletKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.NGOWallet{}, false
	}
	var ngo types.NGOWallet
	k.cdc.MustUnmarshal(bz, &ngo)
	return ngo, true
}

// GetAllNGOWallets returns all NGO wallets
func (k Keeper) GetAllNGOWallets(ctx sdk.Context) []types.NGOWallet {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOWalletKey)
	var ngos []types.NGOWallet
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var ngo types.NGOWallet
		k.cdc.MustUnmarshal(iterator.Value(), &ngo)
		ngos = append(ngos, ngo)
	}
	
	return ngos
}

// GetActiveNGOs returns all active and verified NGOs
func (k Keeper) GetActiveNGOs(ctx sdk.Context) []types.NGOWallet {
	var activeNGOs []types.NGOWallet
	
	ngos := k.GetAllNGOWallets(ctx)
	for _, ngo := range ngos {
		if ngo.IsActive && ngo.IsVerified {
			activeNGOs = append(activeNGOs, ngo)
		}
	}
	
	return activeNGOs
}

// SetNGOWalletCount sets the total NGO wallet count
func (k Keeper) SetNGOWalletCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOWalletCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetNGOWalletCount gets the total NGO wallet count
func (k Keeper) GetNGOWalletCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOWalletCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNGOByAddress sets the NGO ID by address mapping
func (k Keeper) SetNGOByAddress(ctx sdk.Context, address string, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOByAddressKey)
	store.Set([]byte(address), sdk.Uint64ToBigEndian(id))
}

// GetNGOByAddress gets the NGO ID by address
func (k Keeper) GetNGOByAddress(ctx sdk.Context, address string) (uint64, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NGOByAddressKey)
	bz := store.Get([]byte(address))
	if bz == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

// UpdateNGOBalance updates the balance of an NGO wallet
func (k Keeper) UpdateNGOBalance(ctx sdk.Context, ngoId uint64, amount sdk.Coins, isIncoming bool) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	if isIncoming {
		ngo.TotalReceived = ngo.TotalReceived.Add(amount...)
		ngo.CurrentBalance = ngo.CurrentBalance.Add(amount...)
	} else {
		ngo.TotalDistributed = ngo.TotalDistributed.Add(amount...)
		ngo.CurrentBalance = ngo.CurrentBalance.Sub(amount...)
	}
	
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// IncrementNGOBeneficiaryCount increments the beneficiary count for an NGO
func (k Keeper) IncrementNGOBeneficiaryCount(ctx sdk.Context, ngoId uint64, count uint64) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.BeneficiaryCount += count
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// UpdateNGOTransparencyScore updates the transparency score for an NGO
func (k Keeper) UpdateNGOTransparencyScore(ctx sdk.Context, ngoId uint64, score int32) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.TransparencyScore = score
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// UpdateNGOAuditInfo updates audit-related information for an NGO
func (k Keeper) UpdateNGOAuditInfo(ctx sdk.Context, ngoId uint64, lastAuditDate, nextAuditDue int64) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.LastAuditDate = lastAuditDate
	ngo.NextAuditDue = nextAuditDue
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// UpdateNGOImpactMetrics updates impact metrics for an NGO
func (k Keeper) UpdateNGOImpactMetrics(ctx sdk.Context, ngoId uint64, metrics []types.ImpactMetric) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.ImpactMetrics = metrics
	ngo.UpdatedAt = ctx.BlockTime().Unix()
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// VerifyNGO marks an NGO as verified
func (k Keeper) VerifyNGO(ctx sdk.Context, ngoId uint64, verifier string) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.IsVerified = true
	ngo.VerifiedBy = verifier
	ngo.UpdatedAt = ctx.BlockTime().Unix()
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// DeactivateNGO marks an NGO as inactive
func (k Keeper) DeactivateNGO(ctx sdk.Context, ngoId uint64) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.IsActive = false
	ngo.UpdatedAt = ctx.BlockTime().Unix()
	k.SetNGOWallet(ctx, ngo)
	return nil
}

// ActivateNGO marks an NGO as active
func (k Keeper) ActivateNGO(ctx sdk.Context, ngoId uint64) error {
	ngo, found := k.GetNGOWallet(ctx, ngoId)
	if !found {
		return types.ErrNGONotFound
	}
	
	ngo.IsActive = true
	ngo.UpdatedAt = ctx.BlockTime().Unix()
	k.SetNGOWallet(ctx, ngo)
	return nil
}