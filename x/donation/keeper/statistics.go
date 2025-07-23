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

// SetStatistics sets the donation module statistics
func (k Keeper) SetStatistics(ctx sdk.Context, stats types.Statistics) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StatisticsKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set([]byte{0}, bz)
}

// GetStatistics retrieves the donation module statistics
func (k Keeper) GetStatistics(ctx sdk.Context) (types.Statistics, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StatisticsKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return types.Statistics{}, false
	}
	var stats types.Statistics
	k.cdc.MustUnmarshal(bz, &stats)
	return stats, true
}

// UpdateDonationStatistics updates statistics after a donation
func (k Keeper) UpdateDonationStatistics(ctx sdk.Context, amount sdk.Coins, donorCount uint64) {
	stats, found := k.GetStatistics(ctx)
	if !found {
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
		}
	}
	
	stats.TotalDonations = stats.TotalDonations.Add(amount...)
	stats.TotalDonors += donorCount
	stats.LastUpdated = ctx.BlockTime().Unix()
	
	// Recalculate utilization rate
	if !stats.TotalDonations.IsZero() {
		// Calculate utilization as distributed/donations ratio
		utilizationRate := 0.0
		for _, distributed := range stats.TotalDistributed {
			for _, donation := range stats.TotalDonations {
				if distributed.Denom == donation.Denom && !donation.Amount.IsZero() {
					rate := distributed.Amount.ToDec().Quo(donation.Amount.ToDec())
					utilizationRate = rate.MustFloat64()
					break
				}
			}
		}
		stats.UtilizationRate = utilizationRate
	}
	
	k.SetStatistics(ctx, stats)
}

// UpdateDistributionStatistics updates statistics after a distribution
func (k Keeper) UpdateDistributionStatistics(ctx sdk.Context, amount sdk.Coins, beneficiaryCount uint64) {
	stats, found := k.GetStatistics(ctx)
	if !found {
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
		}
	}
	
	stats.TotalDistributed = stats.TotalDistributed.Add(amount...)
	stats.TotalBeneficiaries += beneficiaryCount
	stats.LastUpdated = ctx.BlockTime().Unix()
	
	// Recalculate utilization rate
	if !stats.TotalDonations.IsZero() {
		utilizationRate := 0.0
		for _, distributed := range stats.TotalDistributed {
			for _, donation := range stats.TotalDonations {
				if distributed.Denom == donation.Denom && !donation.Amount.IsZero() {
					rate := distributed.Amount.ToDec().Quo(donation.Amount.ToDec())
					utilizationRate = rate.MustFloat64()
					break
				}
			}
		}
		stats.UtilizationRate = utilizationRate
	}
	
	k.SetStatistics(ctx, stats)
}

// UpdateNGOStatistics updates NGO-related statistics
func (k Keeper) UpdateNGOStatistics(ctx sdk.Context) {
	stats, found := k.GetStatistics(ctx)
	if !found {
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
		}
	}
	
	// Count total and active NGOs
	ngos := k.GetAllNGOWallets(ctx)
	activeCount := uint64(0)
	totalTransparencyScore := int32(0)
	
	for _, ngo := range ngos {
		if ngo.IsActive && ngo.IsVerified {
			activeCount++
		}
		totalTransparencyScore += ngo.TransparencyScore
	}
	
	stats.TotalNGOs = uint64(len(ngos))
	stats.ActiveNGOs = activeCount
	
	// Calculate average transparency score
	if len(ngos) > 0 {
		stats.AverageTransparencyScore = float64(totalTransparencyScore) / float64(len(ngos))
	}
	
	stats.LastUpdated = ctx.BlockTime().Unix()
	k.SetStatistics(ctx, stats)
}

// UpdateCampaignStatistics updates campaign-related statistics
func (k Keeper) UpdateCampaignStatistics(ctx sdk.Context) {
	stats, found := k.GetStatistics(ctx)
	if !found {
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
		}
	}
	
	// Count total and active campaigns
	campaigns := k.GetAllCampaigns(ctx)
	activeCount := uint64(0)
	
	for _, campaign := range campaigns {
		if campaign.IsActive {
			activeCount++
		}
	}
	
	stats.TotalCampaigns = uint64(len(campaigns))
	stats.ActiveCampaigns = activeCount
	stats.LastUpdated = ctx.BlockTime().Unix()
	
	k.SetStatistics(ctx, stats)
}

// UpdateRecurringDonationStatistics updates recurring donation statistics
func (k Keeper) UpdateRecurringDonationStatistics(ctx sdk.Context) {
	stats, found := k.GetStatistics(ctx)
	if !found {
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
		}
	}
	
	// Count total and active recurring donations
	recurringDonations := k.GetAllRecurringDonations(ctx)
	activeCount := uint64(0)
	
	for _, recurring := range recurringDonations {
		if recurring.IsActive {
			activeCount++
		}
	}
	
	stats.TotalRecurringDonations = uint64(len(recurringDonations))
	stats.ActiveRecurringDonations = activeCount
	stats.LastUpdated = ctx.BlockTime().Unix()
	
	k.SetStatistics(ctx, stats)
}

// GetModuleStatistics returns comprehensive module statistics
func (k Keeper) GetModuleStatistics(ctx sdk.Context) types.Statistics {
	stats, found := k.GetStatistics(ctx)
	if !found {
		// Initialize with current state
		stats = types.Statistics{
			TotalDonations:   sdk.Coins{},
			TotalDistributed: sdk.Coins{},
			LastUpdated:      ctx.BlockTime().Unix(),
		}
		
		// Update all statistics
		k.UpdateNGOStatistics(ctx)
		k.UpdateCampaignStatistics(ctx)
		k.UpdateRecurringDonationStatistics(ctx)
		
		// Retrieve updated stats
		stats, _ = k.GetStatistics(ctx)
	}
	
	return stats
}

// SetEmergencyPause sets the emergency pause state
func (k Keeper) SetEmergencyPause(ctx sdk.Context, pause types.EmergencyPause) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EmergencyPauseKey)
	bz := k.cdc.MustMarshal(&pause)
	store.Set([]byte{0}, bz)
}

// GetEmergencyPause retrieves the emergency pause state
func (k Keeper) GetEmergencyPause(ctx sdk.Context) (types.EmergencyPause, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.EmergencyPauseKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return types.EmergencyPause{}, false
	}
	var pause types.EmergencyPause
	k.cdc.MustUnmarshal(bz, &pause)
	return pause, true
}

// IsModulePaused checks if the module is paused
func (k Keeper) IsModulePaused(ctx sdk.Context) bool {
	pause, found := k.GetEmergencyPause(ctx)
	return found && pause.IsPaused
}

// SetFundFlow records a fund flow event
func (k Keeper) SetFundFlow(ctx sdk.Context, flow types.FundFlow) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundFlowKey)
	key := append(sdk.Uint64ToBigEndian(uint64(flow.Timestamp)), sdk.Uint64ToBigEndian(flow.Id)...)
	bz := k.cdc.MustMarshal(&flow)
	store.Set(key, bz)
}

// GetAllFundFlows retrieves all fund flow records
func (k Keeper) GetAllFundFlows(ctx sdk.Context) []types.FundFlow {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundFlowKey)
	var flows []types.FundFlow
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var flow types.FundFlow
		k.cdc.MustUnmarshal(iterator.Value(), &flow)
		flows = append(flows, flow)
	}
	
	return flows
}

// RecordFundFlow creates a new fund flow record
func (k Keeper) RecordFundFlow(ctx sdk.Context, flowType, from, to string, amount sdk.Coins, purpose string, relatedId uint64, relatedType string) {
	flow := types.FundFlow{
		Id:              uint64(ctx.BlockTime().UnixNano()),
		FlowType:        flowType,
		FromAddress:     from,
		ToAddress:       to,
		Amount:          amount,
		Purpose:         purpose,
		TransactionHash: sdk.FormatInvariant("flow_%d", ctx.BlockHeight()),
		Timestamp:       ctx.BlockTime().Unix(),
		BlockHeight:     ctx.BlockHeight(),
		RelatedId:       relatedId,
		RelatedType:     relatedType,
	}
	
	k.SetFundFlow(ctx, flow)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFundFlowUpdate,
			sdk.NewAttribute(types.AttributeKeyFundFlowType, flowType),
			sdk.NewAttribute(types.AttributeKeyFromAddress, from),
			sdk.NewAttribute(types.AttributeKeyToAddress, to),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyPurpose, purpose),
		),
	)
}

// SetTransparencyScore sets a transparency score for an NGO
func (k Keeper) SetTransparencyScore(ctx sdk.Context, score types.TransparencyScore) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransparencyScoreKey)
	bz := k.cdc.MustMarshal(&score)
	store.Set(sdk.Uint64ToBigEndian(score.NgoWalletId), bz)
}

// GetTransparencyScore retrieves a transparency score for an NGO
func (k Keeper) GetTransparencyScore(ctx sdk.Context, ngoWalletId uint64) (types.TransparencyScore, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransparencyScoreKey)
	bz := store.Get(sdk.Uint64ToBigEndian(ngoWalletId))
	if bz == nil {
		return types.TransparencyScore{}, false
	}
	var score types.TransparencyScore
	k.cdc.MustUnmarshal(bz, &score)
	return score, true
}

// GetAllTransparencyScores retrieves all transparency scores
func (k Keeper) GetAllTransparencyScores(ctx sdk.Context) []types.TransparencyScore {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TransparencyScoreKey)
	var scores []types.TransparencyScore
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var score types.TransparencyScore
		k.cdc.MustUnmarshal(iterator.Value(), &score)
		scores = append(scores, score)
	}
	
	return scores
}

// AddToVerificationQueue adds an item to the verification queue
func (k Keeper) AddToVerificationQueue(ctx sdk.Context, item types.VerificationQueueItem) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VerificationQueueKey)
	key := append([]byte(item.Priority), sdk.Uint64ToBigEndian(uint64(item.RequestedAt))...)
	bz := k.cdc.MustMarshal(&item)
	store.Set(key, bz)
}

// GetAllVerificationQueueItems retrieves all items from the verification queue
func (k Keeper) GetAllVerificationQueueItems(ctx sdk.Context) []types.VerificationQueueItem {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VerificationQueueKey)
	var items []types.VerificationQueueItem
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var item types.VerificationQueueItem
		k.cdc.MustUnmarshal(iterator.Value(), &item)
		items = append(items, item)
	}
	
	return items
}