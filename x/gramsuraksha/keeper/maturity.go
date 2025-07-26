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
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// ProcessMaturity processes a matured pension account
func (k Keeper) ProcessMaturity(ctx sdk.Context, participantID string) error {
	// Get participant
	participant, found := k.GetParticipant(ctx, participantID)
	if !found {
		return types.ErrParticipantNotFound
	}

	// Check if already matured
	if participant.Status == types.StatusMatured {
		return types.ErrAlreadyMatured
	}

	// Check if maturity date reached
	if ctx.BlockTime().Before(participant.MaturityDate) {
		return types.ErrNotMatured
	}

	// Get scheme
	scheme, found := k.GetScheme(ctx, participant.SchemeID)
	if !found {
		return types.ErrSchemeNotFound
	}

	// Get current chain performance metrics
	metrics := k.GetChainPerformanceMetrics(ctx)
	params := k.GetDynamicPayoutParams(ctx)
	
	// Calculate dynamic payout rate based on chain performance
	dynamicPayoutRate := types.CalculateDynamicPayout(metrics, params)
	
	// Calculate maturity amount using dynamic rate
	maturityAmount := scheme.CalculateDynamicMaturityAmount(participant.TotalContributed, dynamicPayoutRate)
	
	// Add any bonuses earned
	maturityAmount = maturityAmount.Add(participant.BonusEarned)
	
	// Subtract any penalties
	if participant.PenaltyIncurred.IsPositive() {
		maturityAmount = maturityAmount.Sub(participant.PenaltyIncurred)
		if maturityAmount.IsNegative() {
			maturityAmount = sdk.NewCoin(maturityAmount.Denom, sdk.ZeroInt())
		}
	}
	
	// Log payout rate for transparency
	ctx.Logger().Info("Processing pension maturity with dynamic payout",
		"participant_id", participantID,
		"payout_rate", dynamicPayoutRate.String(),
		"health_score", metrics.HealthScore.String(),
		"risk_score", metrics.RiskScore.String(),
	)

	// Create maturity record
	maturity := types.SurakshaMaturity{
		MaturityID:       k.GenerateMaturityID(ctx, participantID),
		ParticipantID:    participantID,
		SchemeID:         participant.SchemeID,
		Address:          participant.Address,
		MaturityDate:     participant.MaturityDate,
		TotalContributed: participant.TotalContributed,
		MaturityBonus:    sdk.NewCoin(participant.TotalContributed.Denom, maturityAmount.Amount.Sub(participant.TotalContributed.Amount)),
		TotalPayout:      maturityAmount,
		Status:           types.StatusProcessing,
		ProcessedDate:    ctx.BlockTime(),
	}

	// Handle liquidity returns if applicable
	if scheme.LiquidityProvision && participant.LiquidityContributed.IsPositive() {
		// Call hook to process liquidity returns
		if k.hooks != nil {
			pensionAccountID := fmt.Sprintf("%s-%s", participant.SchemeID, participant.ParticipantID)
			if err := k.CallHookAfterMaturity(ctx, pensionAccountID, participant.Address, maturityAmount); err != nil {
				k.Logger(ctx).Error("failed to process liquidity returns", "error", err)
			}
		}
	}

	// Transfer maturity amount to participant
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, participant.Address, sdk.NewCoins(maturityAmount),
	); err != nil {
		return err
	}

	// Update maturity record
	maturity.Status = types.StatusCompleted
	maturity.TransactionHash = fmt.Sprintf("%s-%d", ctx.TxBytes(), ctx.BlockHeight())

	// Store maturity record
	k.SetMaturity(ctx, maturity)

	// Update participant status
	participant.Status = types.StatusMatured
	k.SetParticipant(ctx, participant)

	// Update scheme statistics
	stats, found := k.GetSchemeStatistics(ctx, participant.SchemeID)
	if found {
		stats.MaturedParticipants++
		stats.TotalMaturityPaid = stats.TotalMaturityPaid.Add(maturity.TotalPayout)
		stats.TotalBonusPaid = stats.TotalBonusPaid.Add(maturity.MaturityBonus)
		k.SetSchemeStatistics(ctx, stats)
	}

	// Get cultural quote
	culturalQuote := ""
	if scheme.CulturalIntegration && k.culturalKeeper != nil {
		quote, _ := k.culturalKeeper.GetRandomQuote(ctx, "maturity")
		culturalQuote = quote
	}

	// Update patriotism score for completing scheme
	if k.culturalKeeper != nil {
		k.culturalKeeper.UpdatePatriotismScore(ctx, participant.Address, 100) // 100 points for completion
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMaturityProcessed,
			sdk.NewAttribute(types.AttributeKeyParticipantID, participantID),
			sdk.NewAttribute(types.AttributeKeyMaturityAmount, maturityAmount.String()),
			sdk.NewAttribute(types.AttributeKeyMaturityDate, participant.MaturityDate.String()),
			sdk.NewAttribute(types.AttributeKeyMaturityBonus, maturity.MaturityBonus.String()),
			sdk.NewAttribute(types.AttributeKeyCulturalQuoteText, culturalQuote),
		),
	)

	return nil
}

// SetMaturity stores a maturity record
func (k Keeper) SetMaturity(ctx sdk.Context, maturity types.SurakshaMaturity) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&maturity)
	store.Set(types.MaturityProjectionPrefix.Bytes(append([]byte(maturity.MaturityID)), bz)
}

// GetMaturity retrieves a maturity record
func (k Keeper) GetMaturity(ctx sdk.Context, maturityID string) (types.SurakshaMaturity, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MaturityProjectionPrefix.Bytes(append([]byte(maturityID)))
	if bz == nil {
		return types.SurakshaMaturity{}, false
	}

	var maturity types.SurakshaMaturity
	k.cdc.MustUnmarshal(bz, &maturity)
	return maturity, true
}

// GetUpcomingMaturities gets participants maturing in the next N days
func (k Keeper) GetUpcomingMaturities(ctx sdk.Context, days int) []types.SurakshaParticipant {
	var participants []types.SurakshaParticipant
	
	currentTime := ctx.BlockTime()
	futureTime := currentTime.AddDate(0, 0, days)

	// Iterate all active participants
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ParticipantPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var participant types.SurakshaParticipant
		k.cdc.MustUnmarshal(iterator.Value(), &participant)
		
		if participant.Status == types.StatusActive &&
			participant.MaturityDate.After(currentTime) &&
			participant.MaturityDate.Before(futureTime) {
			participants = append(participants, participant)
		}
	}

	return participants
}

// ProcessAllMaturities processes all matured accounts
func (k Keeper) ProcessAllMaturities(ctx sdk.Context) {
	currentTime := ctx.BlockTime()
	
	// Iterate all active participants
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ParticipantPrefix)
	defer iterator.Close()

	processedCount := 0
	for ; iterator.Valid(); iterator.Next() {
		var participant types.SurakshaParticipant
		k.cdc.MustUnmarshal(iterator.Value(), &participant)
		
		// Check if eligible for maturity processing
		if participant.Status == types.StatusActive && 
			participant.MaturityDate.Before(currentTime) {
			
			if err := k.ProcessMaturity(ctx, participant.ParticipantID); err != nil {
				k.Logger(ctx).Error("failed to process maturity",
					"participant_id", participant.ParticipantID,
					"error", err)
			} else {
				processedCount++
			}
		}
	}

	if processedCount > 0 {
		k.Logger(ctx).Info("processed maturities", "count", processedCount)
	}
}

// CalculateProjectedMaturity calculates projected maturity for a participant
func (k Keeper) CalculateProjectedMaturity(ctx sdk.Context, participantID string) (sdk.Coin, error) {
	participant, found := k.GetParticipant(ctx, participantID)
	if !found {
		return sdk.Coin{}, types.ErrParticipantNotFound
	}

	scheme, found := k.GetScheme(ctx, participant.SchemeID)
	if !found {
		return sdk.Coin{}, types.ErrSchemeNotFound
	}

	// Calculate remaining contributions
	contributionsMade := participant.OnTimePayments + participant.MissedPayments
	remainingContributions := scheme.ContributionPeriod - contributionsMade
	
	// Project total contributions
	projectedContributions := participant.TotalContributed.Add(
		sdk.NewCoin(scheme.MonthlyContribution.Denom,
			scheme.MonthlyContribution.Amount.MulRaw(int64(remainingContributions))),
	)

	// Calculate maturity with bonus
	projectedMaturity := scheme.CalculateMaturityAmount(projectedContributions)
	
	// Add current bonuses
	projectedMaturity = projectedMaturity.Add(participant.BonusEarned)
	
	// Subtract penalties
	if participant.PenaltyIncurred.IsPositive() {
		projectedMaturity = projectedMaturity.Sub(participant.PenaltyIncurred)
		if projectedMaturity.IsNegative() {
			projectedMaturity = sdk.NewCoin(projectedMaturity.Denom, sdk.ZeroInt())
		}
	}

	return projectedMaturity, nil
}

// GenerateMaturityID generates a unique maturity ID
func (k Keeper) GenerateMaturityID(ctx sdk.Context, participantID string) string {
	return fmt.Sprintf("MAT-%s-%d", participantID, ctx.BlockTime().Unix())
}

// NotifyUpcomingMaturities sends notifications for upcoming maturities
func (k Keeper) NotifyUpcomingMaturities(ctx sdk.Context) {
	// Get participants maturing in next 30 days
	upcomingMaturities := k.GetUpcomingMaturities(ctx, 30)
	
	for _, participant := range upcomingMaturities {
		// Calculate days until maturity
		daysUntilMaturity := int(participant.MaturityDate.Sub(ctx.BlockTime()).Hours() / 24)
		
		// Emit notification event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"upcoming_maturity_notification",
				sdk.NewAttribute(types.AttributeKeyParticipantID, participant.ParticipantID),
				sdk.NewAttribute(types.AttributeKeyParticipantAddress, participant.Address.String()),
				sdk.NewAttribute("days_until_maturity", fmt.Sprintf("%d", daysUntilMaturity)),
				sdk.NewAttribute(types.AttributeKeyNotificationType, types.NotificationTypeMaturity),
			),
		)
	}
}