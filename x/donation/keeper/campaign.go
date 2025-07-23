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

// SetCampaign sets a campaign in the store
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKey)
	bz := k.cdc.MustMarshal(&campaign)
	store.Set(sdk.Uint64ToBigEndian(campaign.Id), bz)
}

// GetCampaign retrieves a campaign from the store
func (k Keeper) GetCampaign(ctx sdk.Context, id uint64) (types.Campaign, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.Campaign{}, false
	}
	var campaign types.Campaign
	k.cdc.MustUnmarshal(bz, &campaign)
	return campaign, true
}

// GetAllCampaigns returns all campaigns
func (k Keeper) GetAllCampaigns(ctx sdk.Context) []types.Campaign {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKey)
	var campaigns []types.Campaign
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var campaign types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &campaign)
		campaigns = append(campaigns, campaign)
	}
	
	return campaigns
}

// SetCampaignCount sets the total campaign count
func (k Keeper) SetCampaignCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetCampaignCount gets the total campaign count
func (k Keeper) GetCampaignCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// GetActiveCampaigns returns all active campaigns
func (k Keeper) GetActiveCampaigns(ctx sdk.Context) []types.Campaign {
	var activeCampaigns []types.Campaign
	currentTime := ctx.BlockTime().Unix()
	
	campaigns := k.GetAllCampaigns(ctx)
	for _, campaign := range campaigns {
		if campaign.IsActive && campaign.StartDate <= currentTime && campaign.EndDate > currentTime {
			activeCampaigns = append(activeCampaigns, campaign)
		}
	}
	
	return activeCampaigns
}

// UpdateCampaignProgress updates the raised amount for a campaign
func (k Keeper) UpdateCampaignProgress(ctx sdk.Context, campaignId uint64, amount sdk.Coins) error {
	campaign, found := k.GetCampaign(ctx, campaignId)
	if !found {
		return nil // Campaign donation not required, proceed without error
	}
	
	campaign.RaisedAmount = campaign.RaisedAmount.Add(amount...)
	campaign.DonorCount++
	campaign.UpdatedAt = ctx.BlockTime().Unix()
	
	// Check if campaign goal is reached
	if campaign.RaisedAmount.IsAllGTE(campaign.TargetAmount) && campaign.CompletedAt == 0 {
		campaign.CompletedAt = ctx.BlockTime().Unix()
		campaign.Status = "completed"
		campaign.IsActive = false
		
		// Emit campaign completed event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCampaignCompleted,
				sdk.NewAttribute(types.AttributeKeyCampaignID, sdk.FormatInvariant("%d", campaignId)),
				sdk.NewAttribute(types.AttributeKeyCampaignName, campaign.Name),
				sdk.NewAttribute(types.AttributeKeyTargetAmount, campaign.TargetAmount.String()),
				sdk.NewAttribute(types.AttributeKeyDonationAmount, campaign.RaisedAmount.String()),
			),
		)
	}
	
	k.SetCampaign(ctx, campaign)
	k.UpdateCampaignStatistics(ctx)
	
	return nil
}

// SetRecurringDonation sets a recurring donation in the store
func (k Keeper) SetRecurringDonation(ctx sdk.Context, recurring types.RecurringDonation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecurringDonationKey)
	bz := k.cdc.MustMarshal(&recurring)
	store.Set(sdk.Uint64ToBigEndian(recurring.Id), bz)
}

// GetRecurringDonation retrieves a recurring donation from the store
func (k Keeper) GetRecurringDonation(ctx sdk.Context, id uint64) (types.RecurringDonation, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecurringDonationKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.RecurringDonation{}, false
	}
	var recurring types.RecurringDonation
	k.cdc.MustUnmarshal(bz, &recurring)
	return recurring, true
}

// GetAllRecurringDonations returns all recurring donations
func (k Keeper) GetAllRecurringDonations(ctx sdk.Context) []types.RecurringDonation {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecurringDonationKey)
	var recurringDonations []types.RecurringDonation
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var recurring types.RecurringDonation
		k.cdc.MustUnmarshal(iterator.Value(), &recurring)
		recurringDonations = append(recurringDonations, recurring)
	}
	
	return recurringDonations
}

// SetRecurringDonationCount sets the total recurring donation count
func (k Keeper) SetRecurringDonationCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecurringDonationCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetRecurringDonationCount gets the total recurring donation count
func (k Keeper) GetRecurringDonationCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecurringDonationCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// GetActiveRecurringDonations returns all active recurring donations
func (k Keeper) GetActiveRecurringDonations(ctx sdk.Context) []types.RecurringDonation {
	var activeRecurring []types.RecurringDonation
	currentTime := ctx.BlockTime().Unix()
	
	recurringDonations := k.GetAllRecurringDonations(ctx)
	for _, recurring := range recurringDonations {
		if recurring.IsActive && recurring.StartDate <= currentTime && (recurring.EndDate == 0 || recurring.EndDate > currentTime) {
			activeRecurring = append(activeRecurring, recurring)
		}
	}
	
	return activeRecurring
}

// ProcessRecurringDonations processes all due recurring donations
func (k Keeper) ProcessRecurringDonations(ctx sdk.Context) {
	currentTime := ctx.BlockTime().Unix()
	recurringDonations := k.GetActiveRecurringDonations(ctx)
	
	for _, recurring := range recurringDonations {
		if recurring.NextExecutionDate <= currentTime {
			// Process the recurring donation
			donorAddr, err := sdk.AccAddressFromBech32(recurring.Donor)
			if err != nil {
				continue
			}
			
			// Check if donor has sufficient balance
			balance := k.bankKeeper.GetAllBalances(ctx, donorAddr)
			if !balance.IsAllGTE(recurring.Amount) {
				// Skip this recurring donation due to insufficient funds
				continue
			}
			
			// Transfer funds from donor to module
			if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, donorAddr, types.ModuleName, recurring.Amount); err != nil {
				continue
			}
			
			// Create donation record
			donationId, err := k.CreateDonationRecord(ctx, recurring.Donor, recurring.NgoWalletId, recurring.Amount, recurring.Purpose, false)
			if err != nil {
				// Refund the amount back to donor
				k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, donorAddr, recurring.Amount)
				continue
			}
			
			// Update recurring donation record
			donation, _ := k.GetDonationRecord(ctx, donationId)
			donation.IsRecurring = true
			donation.RecurringId = recurring.Id
			k.SetDonationRecord(ctx, donation)
			
			// Update recurring donation statistics
			recurring.TotalDonated = recurring.TotalDonated.Add(recurring.Amount...)
			recurring.ExecutionCount++
			recurring.LastExecutionDate = currentTime
			
			// Calculate next execution date based on frequency
			switch recurring.Frequency {
			case "daily":
				recurring.NextExecutionDate = currentTime + 86400 // 24 hours
			case "weekly":
				recurring.NextExecutionDate = currentTime + 604800 // 7 days
			case "monthly":
				recurring.NextExecutionDate = currentTime + 2592000 // 30 days
			case "quarterly":
				recurring.NextExecutionDate = currentTime + 7776000 // 90 days
			case "yearly":
				recurring.NextExecutionDate = currentTime + 31536000 // 365 days
			}
			
			// Check if recurring donation has ended
			if recurring.EndDate > 0 && recurring.NextExecutionDate > recurring.EndDate {
				recurring.IsActive = false
			}
			
			recurring.UpdatedAt = currentTime
			k.SetRecurringDonation(ctx, recurring)
			
			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeRecurringDonation,
					sdk.NewAttribute(types.AttributeKeyRecurringDonationID, sdk.FormatInvariant("%d", recurring.Id)),
					sdk.NewAttribute(types.AttributeKeyDonor, recurring.Donor),
					sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", recurring.NgoWalletId)),
					sdk.NewAttribute(types.AttributeKeyDonationAmount, recurring.Amount.String()),
					sdk.NewAttribute(types.AttributeKeyDonationID, sdk.FormatInvariant("%d", donationId)),
				),
			)
		}
	}
	
	// Update statistics
	k.UpdateRecurringDonationStatistics(ctx)
}

// CancelRecurringDonation cancels a recurring donation
func (k Keeper) CancelRecurringDonation(ctx sdk.Context, id uint64, reason string) error {
	recurring, found := k.GetRecurringDonation(ctx, id)
	if !found {
		return nil // Not found, no error
	}
	
	recurring.IsActive = false
	recurring.CancelledAt = ctx.BlockTime().Unix()
	recurring.CancelReason = reason
	recurring.UpdatedAt = ctx.BlockTime().Unix()
	
	k.SetRecurringDonation(ctx, recurring)
	k.UpdateRecurringDonationStatistics(ctx)
	
	return nil
}