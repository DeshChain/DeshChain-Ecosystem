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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// SetSevaMitra stores an sevaMitra
func (k Keeper) SetSevaMitra(ctx sdk.Context, sevaMitra *types.SevaMitra) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraKey(sevaMitra.MitraId)
	value := k.cdc.MustMarshal(sevaMitra)
	store.Set(key, value)
	
	// Index by address
	k.SetSevaMitraAddressIndex(ctx, sevaMitra.Address, sevaMitra.MitraId)
}

// GetSevaMitra retrieves an sevaMitra by ID
func (k Keeper) GetSevaMitra(ctx sdk.Context, sevaMitraID string) (*types.SevaMitra, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraKey(sevaMitraID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	var sevaMitra types.SevaMitra
	k.cdc.MustUnmarshal(value, &sevaMitra)
	return &sevaMitra, true
}

// IsSevaMitraRegistered checks if an address is already registered as sevaMitra
func (k Keeper) IsSevaMitraRegistered(ctx sdk.Context, address string) bool {
	sevaMitraID := k.GetSevaMitraAddressIndex(ctx, address)
	return sevaMitraID != ""
}

// SetSevaMitraAddressIndex stores sevaMitra ID by address
func (k Keeper) SetSevaMitraAddressIndex(ctx sdk.Context, address, sevaMitraID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraAddressIndexKey(address)
	store.Set(key, []byte(sevaMitraID))
}

// GetSevaMitraAddressIndex retrieves sevaMitra ID by address
func (k Keeper) GetSevaMitraAddressIndex(ctx sdk.Context, address string) string {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraAddressIndexKey(address)
	value := store.Get(key)
	if value == nil {
		return ""
	}
	return string(value)
}

// AddSevaMitraToDistrictIndex adds sevaMitra to district index
func (k Keeper) AddSevaMitraToDistrictIndex(ctx sdk.Context, district, sevaMitraID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraDistrictIndexKey(district, sevaMitraID)
	store.Set(key, []byte{1})
}

// GetSevaMitrasByDistrict returns all sevaMitras in a district
func (k Keeper) GetSevaMitrasByDistrict(ctx sdk.Context, district string) []*types.SevaMitra {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetSevaMitraDistrictIndexPrefix(district)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var sevaMitras []*types.SevaMitra
	for ; iterator.Valid(); iterator.Next() {
		// Extract sevaMitra ID from key
		sevaMitraID := string(iterator.Key()[len(prefix):])
		sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
		if found {
			sevaMitras = append(sevaMitras, sevaMitra)
		}
	}
	
	return sevaMitras
}

// sevaMitraProvidesServices checks if sevaMitra provides required services
func (k Keeper) sevaMitraProvidesServices(sevaMitra *types.SevaMitra, requiredServices []types.SevaMitraService) bool {
	if len(requiredServices) == 0 {
		return true
	}
	
	serviceMap := make(map[types.SevaMitraService]bool)
	for _, service := range sevaMitra.Services {
		serviceMap[service] = true
	}
	
	for _, required := range requiredServices {
		if !serviceMap[required] {
			return false
		}
	}
	
	return true
}

// isWithinDistance checks if two postal codes are within distance
func (k Keeper) isWithinDistance(postalCode1, postalCode2 string, maxDistanceKm int32) bool {
	// Simple distance check based on postal code similarity
	// In production, would use actual geographic data
	
	if postalCode1 == postalCode2 {
		return true
	}
	
	// Same first 3 digits = same district (approximately)
	if postalCode1[:3] == postalCode2[:3] {
		return maxDistanceKm >= 10
	}
	
	// Same first 2 digits = same region
	if postalCode1[:2] == postalCode2[:2] {
		return maxDistanceKm >= 50
	}
	
	// Same state
	if k.GetStateName(postalCode1) == k.GetStateName(postalCode2) {
		return maxDistanceKm >= 200
	}
	
	return false
}

// GetSevaMitraDailyVolume retrieves sevaMitra's daily transaction volume
func (k Keeper) GetSevaMitraDailyVolume(ctx sdk.Context, sevaMitraID, date string) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraDailyVolumeKey(sevaMitraID, date)
	value := store.Get(key)
	if value == nil {
		return sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt())
	}
	
	var volume sdk.Coin
	k.cdc.MustUnmarshal(value, &volume)
	return volume
}

// SetSevaMitraDailyVolume stores sevaMitra's daily transaction volume
func (k Keeper) SetSevaMitraDailyVolume(ctx sdk.Context, sevaMitraID, date string, volume sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraDailyVolumeKey(sevaMitraID, date)
	value := k.cdc.MustMarshal(&volume)
	store.Set(key, value)
}

// StoreSevaMitraRating stores a rating for an sevaMitra
func (k Keeper) StoreSevaMitraRating(ctx sdk.Context, sevaMitraID string, rater sdk.AccAddress, rating int32, comment string) {
	ratingData := &types.SevaMitraRating{
		MitraId:   sevaMitraID,
		Rater:     rater.String(),
		Rating:    rating,
		Comment:   comment,
		Timestamp: ctx.BlockTime(),
	}
	
	store := ctx.KVStore(k.storeKey)
	key := types.GetSevaMitraRatingKey(sevaMitraID, rater)
	value := k.cdc.MustMarshal(ratingData)
	store.Set(key, value)
}

// GetSevaMitraRatings retrieves all ratings for an sevaMitra
func (k Keeper) GetSevaMitraRatings(ctx sdk.Context, sevaMitraID string) []*types.SevaMitraRating {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetSevaMitraRatingsPrefix(sevaMitraID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var ratings []*types.SevaMitraRating
	for ; iterator.Valid(); iterator.Next() {
		var rating types.SevaMitraRating
		k.cdc.MustUnmarshal(iterator.Value(), &rating)
		ratings = append(ratings, &rating)
	}
	
	return ratings
}

// HasRole checks if an address has a specific role
func (k Keeper) HasRole(ctx sdk.Context, address sdk.AccAddress, role string) bool {
	// In production, would integrate with role-based access control
	// For now, simplified implementation
	return true
}

// AddToKYCQueue adds sevaMitra to KYC verification queue
func (k Keeper) AddToKYCQueue(ctx sdk.Context, sevaMitraID string, scheduledAt time.Time) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetKYCQueueKey(scheduledAt, sevaMitraID)
	store.Set(key, []byte{1})
}

// GetActiveSevaMitras returns all active sevaMitras
func (k Keeper) GetActiveSevaMitras(ctx sdk.Context) []*types.SevaMitra {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SevaMitraPrefix)
	defer iterator.Close()
	
	var sevaMitras []*types.SevaMitra
	for ; iterator.Valid(); iterator.Next() {
		var sevaMitra types.SevaMitra
		k.cdc.MustUnmarshal(iterator.Value(), &sevaMitra)
		if sevaMitra.Status == types.SevaMitraStatus_MITRA_STATUS_ACTIVE {
			sevaMitras = append(sevaMitras, &sevaMitra)
		}
	}
	
	return sevaMitras
}

// GetSevaMitraStats returns aggregated statistics for all sevaMitras
func (k Keeper) GetSevaMitraStats(ctx sdk.Context) *types.SevaMitraSystemStats {
	stats := &types.SevaMitraSystemStats{
		TotalSevaMitras:       0,
		ActiveSevaMitras:      0,
		SuspendedSevaMitras:   0,
		TotalVolume:       sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
		TotalTransactions: 0,
		AverageRating:     0,
		ServiceCoverage:   make(map[string]int32),
		DistrictCoverage:  make(map[string]int32),
	}
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SevaMitraPrefix)
	defer iterator.Close()
	
	totalRating := float64(0)
	ratedSevaMitras := int32(0)
	
	for ; iterator.Valid(); iterator.Next() {
		var sevaMitra types.SevaMitra
		k.cdc.MustUnmarshal(iterator.Value(), &sevaMitra)
		
		stats.TotalSevaMitras++
		
		switch sevaMitra.Status {
		case types.SevaMitraStatus_MITRA_STATUS_ACTIVE:
			stats.ActiveSevaMitras++
		case types.SevaMitraStatus_MITRA_STATUS_SUSPENDED:
			stats.SuspendedSevaMitras++
		}
		
		// Aggregate stats
		if sevaMitra.Stats != nil {
			stats.TotalVolume = stats.TotalVolume.Add(sevaMitra.Stats.TotalVolume)
			stats.TotalTransactions += sevaMitra.Stats.TotalTransactions
			
			if sevaMitra.Stats.AverageRating > 0 {
				totalRating += sevaMitra.Stats.AverageRating
				ratedSevaMitras++
			}
		}
		
		// Service coverage
		for _, service := range sevaMitra.Services {
			stats.ServiceCoverage[service.String()]++
		}
		
		// District coverage
		stats.DistrictCoverage[sevaMitra.District]++
	}
	
	// Calculate average rating
	if ratedSevaMitras > 0 {
		stats.AverageRating = totalRating / float64(ratedSevaMitras)
	}
	
	return stats
}