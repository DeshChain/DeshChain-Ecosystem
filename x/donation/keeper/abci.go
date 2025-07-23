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
	"deshchain/x/donation/types"
)

// BeginBlocker is called at the beginning of each block
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Check if module is paused
	if k.IsModulePaused(ctx) {
		return
	}
	
	// Process recurring donations
	k.ProcessRecurringDonations(ctx)
	
	// Check and update campaign statuses
	currentTime := ctx.BlockTime().Unix()
	campaigns := k.GetAllCampaigns(ctx)
	
	for _, campaign := range campaigns {
		if campaign.IsActive && campaign.EndDate <= currentTime {
			// Campaign has ended
			campaign.IsActive = false
			campaign.Status = "ended"
			if campaign.CompletedAt == 0 {
				campaign.CompletedAt = currentTime
			}
			campaign.UpdatedAt = currentTime
			k.SetCampaign(ctx, campaign)
			
			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCampaignCompleted,
					sdk.NewAttribute(types.AttributeKeyCampaignID, sdk.FormatInvariant("%d", campaign.Id)),
					sdk.NewAttribute(types.AttributeKeyCampaignName, campaign.Name),
					sdk.NewAttribute(types.AttributeKeyTargetAmount, campaign.TargetAmount.String()),
					sdk.NewAttribute(types.AttributeKeyDonationAmount, campaign.RaisedAmount.String()),
				),
			)
		}
	}
	
	// Update module statistics
	k.UpdateModuleStatistics(ctx)
}

// EndBlocker is called at the end of each block
func EndBlocker(ctx sdk.Context, k Keeper) []sdk.ValidatorUpdate {
	// Check if module is paused
	if k.IsModulePaused(ctx) {
		return []sdk.ValidatorUpdate{}
	}
	
	// Process audit reminders
	ngos := k.GetAllNGOWallets(ctx)
	currentTime := ctx.BlockTime().Unix()
	
	for _, ngo := range ngos {
		if ngo.IsActive && ngo.IsVerified {
			// Check if audit is due
			if ngo.NextAuditDue > 0 && ngo.NextAuditDue <= currentTime {
				// Emit audit reminder event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeAuditReminder,
						sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", ngo.Id)),
						sdk.NewAttribute(types.AttributeKeyNGOName, ngo.Name),
						sdk.NewAttribute(types.AttributeKeyAuditDueDate, sdk.FormatInvariant("%d", ngo.NextAuditDue)),
					),
				)
			}
			
			// Recalculate transparency score if needed
			score, found := k.GetTransparencyScore(ctx, ngo.Id)
			if !found || score.NextCalculationDue <= currentTime {
				k.CalculateTransparencyScore(ctx, ngo.Id)
			}
		}
	}
	
	// Process verification queue
	verificationItems := k.GetAllVerificationQueueItems(ctx)
	for _, item := range verificationItems {
		if item.Status == types.VerificationStatusPending {
			// Check if item has been pending too long
			if currentTime-item.RequestedAt > 86400*30 { // 30 days
				// Mark as expired
				item.Status = types.VerificationStatusExpired
				item.CompletedAt = currentTime
				item.Result = "Expired due to timeout"
				k.AddToVerificationQueue(ctx, item)
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeVerificationExpired,
						sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", item.NgoWalletId)),
						sdk.NewAttribute(types.AttributeKeyVerificationType, item.VerificationType),
					),
				)
			}
		}
	}
	
	// Update average transparency score in statistics
	ngos = k.GetActiveNGOs(ctx)
	if len(ngos) > 0 {
		totalScore := int32(0)
		for _, ngo := range ngos {
			totalScore += ngo.TransparencyScore
		}
		avgScore := float64(totalScore) / float64(len(ngos))
		
		stats, _ := k.GetStatistics(ctx)
		stats.AverageTransparencyScore = avgScore
		stats.LastUpdated = currentTime
		k.SetStatistics(ctx, stats)
	}
	
	return []sdk.ValidatorUpdate{}
}

// ProcessDonationDistribution processes the distribution of donations to NGOs
func (k Keeper) ProcessDonationDistribution(ctx sdk.Context, amount sdk.Coins) error {
	// Get active NGOs
	activeNGOs := k.GetActiveNGOs(ctx)
	if len(activeNGOs) == 0 {
		return types.ErrNoActiveNGOs
	}
	
	// Distribute equally among active NGOs
	ngoCount := sdk.NewInt(int64(len(activeNGOs)))
	perNGOAmount := sdk.Coins{}
	
	for _, coin := range amount {
		if coin.Amount.GT(sdk.ZeroInt()) {
			share := coin.Amount.Quo(ngoCount)
			if share.GT(sdk.ZeroInt()) {
				perNGOAmount = perNGOAmount.Add(sdk.NewCoin(coin.Denom, share))
			}
		}
	}
	
	if perNGOAmount.IsZero() {
		return nil
	}
	
	// Distribute to each NGO
	for _, ngo := range activeNGOs {
		ngoAddr, err := sdk.AccAddressFromBech32(ngo.Address)
		if err != nil {
			continue
		}
		
		// Send funds from module to NGO
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ngoAddr, perNGOAmount); err != nil {
			return err
		}
		
		// Update NGO balance
		if err := k.UpdateNGOBalance(ctx, ngo.Id, perNGOAmount, true); err != nil {
			return err
		}
		
		// Record fund flow
		k.RecordFundFlow(ctx, "distribution", types.ModuleName, ngo.Address, perNGOAmount, "NGO donation distribution", 0, "auto_distribution")
		
		// Create distribution record
		count := k.GetDistributionRecordCount(ctx)
		distributionId := count + 1
		
		distribution := types.DistributionRecord{
			Id:               distributionId,
			NgoWalletId:      ngo.Id,
			Recipient:        ngo.Address,
			Amount:           perNGOAmount,
			Purpose:          "Automated donation distribution",
			Category:         "auto_distribution",
			ProjectName:      "Platform Distribution",
			BeneficiaryName:  ngo.Name,
			DistributedAt:    ctx.BlockTime().Unix(),
			TransactionHash:  sdk.FormatInvariant("dist_%d", distributionId),
			Status:           "completed",
			ExecutedBy:       types.ModuleName,
			BlockHeight:      ctx.BlockHeight(),
		}
		
		k.SetDistributionRecord(ctx, distribution)
		k.SetDistributionRecordCount(ctx, distributionId)
		k.AddDistributionByNGO(ctx, ngo.Id, distributionId)
		
		// Update statistics
		k.UpdateDistributionStatistics(ctx, perNGOAmount, 0)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDistributeFunds,
				sdk.NewAttribute(types.AttributeKeyDistributionID, sdk.FormatInvariant("%d", distributionId)),
				sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", ngo.Id)),
				sdk.NewAttribute(types.AttributeKeyDistributionAmount, perNGOAmount.String()),
				sdk.NewAttribute(types.AttributeKeyRecipient, ngo.Address),
			),
		)
	}
	
	return nil
}