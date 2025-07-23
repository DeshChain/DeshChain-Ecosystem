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

// SetDonationRecord sets a donation record in the store
func (k Keeper) SetDonationRecord(ctx sdk.Context, record types.DonationRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationRecordKey)
	bz := k.cdc.MustMarshal(&record)
	store.Set(sdk.Uint64ToBigEndian(record.Id), bz)
}

// GetDonationRecord retrieves a donation record from the store
func (k Keeper) GetDonationRecord(ctx sdk.Context, id uint64) (types.DonationRecord, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationRecordKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.DonationRecord{}, false
	}
	var record types.DonationRecord
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// GetAllDonationRecords returns all donation records
func (k Keeper) GetAllDonationRecords(ctx sdk.Context) []types.DonationRecord {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationRecordKey)
	var records []types.DonationRecord
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var record types.DonationRecord
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}
	
	return records
}

// SetDonationRecordCount sets the total donation record count
func (k Keeper) SetDonationRecordCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationRecordCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetDonationRecordCount gets the total donation record count
func (k Keeper) GetDonationRecordCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationRecordCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// AddDonationByDonor adds a donation ID to the donor's donation list
func (k Keeper) AddDonationByDonor(ctx sdk.Context, donor string, donationId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationByDonorKey)
	key := append([]byte(donor), sdk.Uint64ToBigEndian(donationId)...)
	store.Set(key, []byte{1})
}

// GetDonationsByDonor retrieves all donation IDs for a donor
func (k Keeper) GetDonationsByDonor(ctx sdk.Context, donor string) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationByDonorKey)
	var donationIds []uint64
	
	iterator := store.Iterator([]byte(donor), sdk.PrefixEndBytes([]byte(donor)))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract donation ID from key (donor address + donation ID)
		if len(key) >= len(donor)+8 {
			donationId := sdk.BigEndianToUint64(key[len(donor):])
			donationIds = append(donationIds, donationId)
		}
	}
	
	return donationIds
}

// AddDonationByNGO adds a donation ID to the NGO's donation list
func (k Keeper) AddDonationByNGO(ctx sdk.Context, ngoWalletId, donationId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationByNGOKey)
	key := append(sdk.Uint64ToBigEndian(ngoWalletId), sdk.Uint64ToBigEndian(donationId)...)
	store.Set(key, []byte{1})
}

// GetDonationsByNGO retrieves all donation IDs for an NGO
func (k Keeper) GetDonationsByNGO(ctx sdk.Context, ngoWalletId uint64) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DonationByNGOKey)
	var donationIds []uint64
	
	ngoKey := sdk.Uint64ToBigEndian(ngoWalletId)
	iterator := store.Iterator(ngoKey, sdk.PrefixEndBytes(ngoKey))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract donation ID from key (NGO ID + donation ID)
		if len(key) >= 16 {
			donationId := sdk.BigEndianToUint64(key[8:])
			donationIds = append(donationIds, donationId)
		}
	}
	
	return donationIds
}

// CreateDonationRecord creates a new donation record
func (k Keeper) CreateDonationRecord(ctx sdk.Context, donor string, ngoWalletId uint64, amount sdk.Coins, purpose string, isAnonymous bool) (uint64, error) {
	// Verify NGO exists and is active
	ngo, found := k.GetNGOWallet(ctx, ngoWalletId)
	if !found {
		return 0, types.ErrNGONotFound
	}
	if !ngo.IsActive {
		return 0, types.ErrNGONotActive
	}
	if !ngo.IsVerified {
		return 0, types.ErrNGONotVerified
	}
	
	// Get next donation ID
	count := k.GetDonationRecordCount(ctx)
	donationId := count + 1
	
	// Create donation record
	record := types.DonationRecord{
		Id:              donationId,
		Donor:           donor,
		NgoWalletId:     ngoWalletId,
		Amount:          amount,
		Purpose:         purpose,
		Category:        ngo.Category,
		IsAnonymous:     isAnonymous,
		DonatedAt:       ctx.BlockTime().Unix(),
		TransactionHash: sdk.FormatInvariant("donation_%d", donationId),
		BlockHeight:     ctx.BlockHeight(),
	}
	
	// Calculate tax benefit (50% of donation amount)
	taxBenefit := sdk.Coins{}
	for _, coin := range amount {
		benefitAmount := coin.Amount.Quo(sdk.NewInt(2)) // 50%
		taxBenefit = taxBenefit.Add(sdk.NewCoin(coin.Denom, benefitAmount))
	}
	record.TaxBenefitAmount = taxBenefit
	
	// Store the record
	k.SetDonationRecord(ctx, record)
	k.SetDonationRecordCount(ctx, donationId)
	
	// Add to indices
	k.AddDonationByDonor(ctx, donor, donationId)
	k.AddDonationByNGO(ctx, ngoWalletId, donationId)
	
	// Update NGO balance
	if err := k.UpdateNGOBalance(ctx, ngoWalletId, amount, true); err != nil {
		return 0, err
	}
	
	// Update statistics
	k.UpdateDonationStatistics(ctx, amount, 1)
	
	// Emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDonate,
			sdk.NewAttribute(types.AttributeKeyDonor, donor),
			sdk.NewAttribute(types.AttributeKeyDonationID, sdk.FormatInvariant("%d", donationId)),
			sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", ngoWalletId)),
			sdk.NewAttribute(types.AttributeKeyDonationAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyDonationPurpose, purpose),
			sdk.NewAttribute(types.AttributeKeyIsAnonymous, sdk.FormatInvariant("%t", isAnonymous)),
			sdk.NewAttribute(types.AttributeKeyTaxBenefitAmount, taxBenefit.String()),
		),
	})
	
	return donationId, nil
}

// SetDistributionRecord sets a distribution record in the store
func (k Keeper) SetDistributionRecord(ctx sdk.Context, record types.DistributionRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionRecordKey)
	bz := k.cdc.MustMarshal(&record)
	store.Set(sdk.Uint64ToBigEndian(record.Id), bz)
}

// GetDistributionRecord retrieves a distribution record from the store
func (k Keeper) GetDistributionRecord(ctx sdk.Context, id uint64) (types.DistributionRecord, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionRecordKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.DistributionRecord{}, false
	}
	var record types.DistributionRecord
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// GetAllDistributionRecords returns all distribution records
func (k Keeper) GetAllDistributionRecords(ctx sdk.Context) []types.DistributionRecord {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionRecordKey)
	var records []types.DistributionRecord
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var record types.DistributionRecord
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}
	
	return records
}

// SetDistributionRecordCount sets the total distribution record count
func (k Keeper) SetDistributionRecordCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionRecordCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetDistributionRecordCount gets the total distribution record count
func (k Keeper) GetDistributionRecordCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionRecordCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// AddDistributionByNGO adds a distribution ID to the NGO's distribution list
func (k Keeper) AddDistributionByNGO(ctx sdk.Context, ngoWalletId, distributionId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionByNGOKey)
	key := append(sdk.Uint64ToBigEndian(ngoWalletId), sdk.Uint64ToBigEndian(distributionId)...)
	store.Set(key, []byte{1})
}

// GetDistributionsByNGO retrieves all distribution IDs for an NGO
func (k Keeper) GetDistributionsByNGO(ctx sdk.Context, ngoWalletId uint64) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DistributionByNGOKey)
	var distributionIds []uint64
	
	ngoKey := sdk.Uint64ToBigEndian(ngoWalletId)
	iterator := store.Iterator(ngoKey, sdk.PrefixEndBytes(ngoKey))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract distribution ID from key (NGO ID + distribution ID)
		if len(key) >= 16 {
			distributionId := sdk.BigEndianToUint64(key[8:])
			distributionIds = append(distributionIds, distributionId)
		}
	}
	
	return distributionIds
}