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
	"github.com/deshchain/deshchain/x/gramsuraksha/types"
)

// MakeContribution processes a monthly contribution
func (k Keeper) MakeContribution(ctx sdk.Context, participantID string, amount sdk.Coin, month uint32) error {
	// Get participant
	participant, found := k.GetParticipant(ctx, participantID)
	if !found {
		return types.ErrParticipantNotFound
	}

	// Check participant status
	if participant.Status != types.StatusActive {
		return types.ErrNotEnrolled
	}

	// Get scheme
	scheme, found := k.GetScheme(ctx, participant.SchemeID)
	if !found {
		return types.ErrSchemeNotFound
	}

	// Validate contribution month
	if month > scheme.ContributionPeriod {
		return types.ErrInvalidContributionMonth
	}

	// Check if contribution already made for this month
	if k.HasContribution(ctx, participantID, month) {
		return types.ErrContributionAlreadyMade
	}

	// Validate amount
	expectedAmount := scheme.MonthlyContribution
	if !amount.Equal(expectedAmount) {
		return types.ErrInvalidContribution
	}

	// Calculate admin fee (0.1% of contribution)
	adminFeeRate := sdk.MustNewDecFromStr(types.DefaultSchemeAdminFee)
	adminFee := amount.Amount.ToDec().Mul(adminFeeRate).TruncateInt()
	adminFeeCoin := sdk.NewCoin(amount.Denom, adminFee)
	
	// Net contribution after admin fee
	netContribution := sdk.NewCoin(amount.Denom, amount.Amount.Sub(adminFee))
	
	// Transfer funds from contributor to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, participant.Address, types.ModuleName, sdk.NewCoins(amount),
	); err != nil {
		return err
	}
	
	// Collect admin fee using revenue keeper
	if adminFee.IsPositive() && k.revenueKeeper != nil {
		feeCoins := sdk.NewCoins(adminFeeCoin)
		if err := k.revenueKeeper.CollectServiceFee(ctx, types.ModuleName, participant.Address, feeCoins, "pension_contribution"); err != nil {
			k.Logger(ctx).Error("Failed to collect admin fee", "error", err, "participant", participantID, "fee", adminFee)
			// Continue even if fee collection fails
		}
	}

	// Create contribution record
	contribution := types.SurakshaContribution{
		ContributionID:   k.GenerateContributionID(ctx, participantID, month),
		ParticipantID:    participantID,
		SchemeID:         participant.SchemeID,
		Address:          participant.Address,
		Amount:           amount,
		Month:            month,
		ContributionDate: ctx.BlockTime(),
		TransactionHash:  fmt.Sprintf("%s-%d", ctx.TxBytes(), ctx.BlockHeight()),
		Status:           types.StatusCompleted,
		LiquidityProvided: scheme.LiquidityProvision,
	}

	// Apply on-time bonus if applicable
	daysSinceExpected := k.CalculateDaysSinceExpectedPayment(ctx, participant, month)
	if daysSinceExpected <= int(scheme.GracePeriodDays) && scheme.OnTimeBonusPercent.IsPositive() {
		bonusAmount := amount.Amount.ToDec().Mul(scheme.OnTimeBonusPercent).TruncateInt()
		contribution.BonusApplied = sdk.NewCoin(amount.Denom, bonusAmount)
		participant.BonusEarned = participant.BonusEarned.Add(contribution.BonusApplied)
		participant.OnTimePayments++
	} else if daysSinceExpected > int(scheme.GracePeriodDays) && scheme.LatePaymentPenalty.IsPositive() {
		// Apply late payment penalty
		penaltyAmount := amount.Amount.ToDec().Mul(scheme.LatePaymentPenalty).TruncateInt()
		contribution.PenaltyApplied = sdk.NewCoin(amount.Denom, penaltyAmount)
		participant.PenaltyIncurred = participant.PenaltyIncurred.Add(contribution.PenaltyApplied)
		participant.MissedPayments++
	} else {
		participant.OnTimePayments++
	}

	// Get cultural quote if enabled
	if scheme.CulturalIntegration && k.culturalKeeper != nil {
		quote, _ := k.culturalKeeper.GetRandomQuote(ctx, "contribution")
		contribution.CulturalQuote = quote
	}

	// Update participant totals (use net contribution after admin fee)
	participant.TotalContributed = participant.TotalContributed.Add(netContribution)
	participant.LastContribution = ctx.BlockTime()

	// Handle liquidity provision if enabled
	if scheme.LiquidityProvision && k.hooks != nil {
		liquidityAmount := scheme.GetLiquidityAmount(amount)
		if liquidityAmount.IsPositive() {
			// Call money order hook to add liquidity
			pensionAccountID := fmt.Sprintf("%s-%s", participant.SchemeID, participant.ParticipantID)
			if err := k.CallHookAfterContribution(ctx, pensionAccountID, participant.Address, liquidityAmount, participant.VillagePostalCode); err != nil {
				k.Logger(ctx).Error("failed to provide liquidity", "error", err)
				// Don't fail contribution, just log error
			} else {
				participant.LiquidityContributed = participant.LiquidityContributed.Add(liquidityAmount)
			}
		}
	}

	// Process referral rewards for first contribution
	if month == 1 && !participant.ReferrerAddress.Empty() && scheme.ReferralRewardPercent.IsPositive() {
		k.ProcessReferralReward(ctx, participant, amount, scheme)
	}

	// Update patriotism score
	if k.culturalKeeper != nil {
		k.culturalKeeper.UpdatePatriotismScore(ctx, participant.Address, 10) // 10 points per contribution
	}

	// Store contribution
	k.SetContribution(ctx, contribution)
	
	// Update participant
	k.SetParticipant(ctx, participant)

	// Update scheme statistics
	k.UpdateSchemeStatistics(ctx, participant.SchemeID)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContributionMade,
			sdk.NewAttribute(types.AttributeKeyParticipantID, participantID),
			sdk.NewAttribute(types.AttributeKeyContributionAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyContributionMonth, fmt.Sprintf("%d", month)),
			sdk.NewAttribute(types.AttributeKeyContributionStatus, contribution.Status),
			sdk.NewAttribute(types.AttributeKeyCulturalQuoteText, contribution.CulturalQuote),
		),
	)

	return nil
}

// SetContribution stores a contribution
func (k Keeper) SetContribution(ctx sdk.Context, contribution types.SurakshaContribution) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&contribution)
	store.Set(types.ContributionPrefix.Bytes(append([]byte(contribution.ContributionID)), bz)
}

// GetContribution retrieves a contribution
func (k Keeper) GetContribution(ctx sdk.Context, contributionID string) (types.SurakshaContribution, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ContributionPrefix.Bytes(append([]byte(contributionID)))
	if bz == nil {
		return types.SurakshaContribution{}, false
	}

	var contribution types.SurakshaContribution
	k.cdc.MustUnmarshal(bz, &contribution)
	return contribution, true
}

// HasContribution checks if a contribution exists for a participant and month
func (k Keeper) HasContribution(ctx sdk.Context, participantID string, month uint32) bool {
	contributionID := fmt.Sprintf("%s-M%02d", participantID, month)
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ContributionPrefix.Bytes(append([]byte(contributionID)))
}

// GetParticipantContributions gets all contributions for a participant
func (k Keeper) GetParticipantContributions(ctx sdk.Context, participantID string) []types.SurakshaContribution {
	var contributions []types.SurakshaContribution
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ContributionPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var contribution types.SurakshaContribution
		k.cdc.MustUnmarshal(iterator.Value(), &contribution)
		if contribution.ParticipantID == participantID {
			contributions = append(contributions, contribution)
		}
	}

	return contributions
}

// ProcessReferralReward processes referral rewards
func (k Keeper) ProcessReferralReward(ctx sdk.Context, participant types.SurakshaParticipant, contributionAmount sdk.Coin, scheme types.SurakshaScheme) {
	if participant.ReferrerAddress.Empty() || scheme.ReferralRewardPercent.IsZero() {
		return
	}

	// Calculate referral reward
	rewardAmount := contributionAmount.Amount.ToDec().Mul(scheme.ReferralRewardPercent).TruncateInt()
	if rewardAmount.IsZero() {
		return
	}

	reward := sdk.NewCoin(contributionAmount.Denom, rewardAmount)

	// Get referrer participant
	referrer, found := k.GetParticipantByAddress(ctx, participant.ReferrerAddress, participant.SchemeID)
	if !found {
		return
	}

	// Update referrer rewards
	referrer.ReferralRewards = referrer.ReferralRewards.Add(reward)
	k.SetParticipant(ctx, referrer)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeReferralRewardProcessed,
			sdk.NewAttribute(types.AttributeKeyReferrerAddress, participant.ReferrerAddress.String()),
			sdk.NewAttribute(types.AttributeKeyReferredParticipant, participant.ParticipantID),
			sdk.NewAttribute(types.AttributeKeyReferralReward, reward.String()),
		),
	)
}

// CalculateDaysSinceExpectedPayment calculates days since expected payment date
func (k Keeper) CalculateDaysSinceExpectedPayment(ctx sdk.Context, participant types.SurakshaParticipant, month uint32) int {
	// Calculate expected payment date (month * 30 days from enrollment)
	expectedDate := participant.EnrollmentDate.AddDate(0, 0, int((month-1)*30))
	currentDate := ctx.BlockTime()
	
	// If current date is before expected date, return negative days
	if currentDate.Before(expectedDate) {
		return -int(expectedDate.Sub(currentDate).Hours() / 24)
	}
	
	// Return positive days if late
	return int(currentDate.Sub(expectedDate).Hours() / 24)
}

// GenerateContributionID generates a unique contribution ID
func (k Keeper) GenerateContributionID(ctx sdk.Context, participantID string, month uint32) string {
	return fmt.Sprintf("%s-M%02d-%d", participantID, month, ctx.BlockTime().Unix())
}

// GetMonthlyContributions gets contributions for a specific month across all schemes
func (k Keeper) GetMonthlyContributions(ctx sdk.Context, year int, month int) sdk.Coins {
	totalContributions := sdk.NewCoins()
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ContributionPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var contribution types.SurakshaContribution
		k.cdc.MustUnmarshal(iterator.Value(), &contribution)
		
		// Check if contribution is from the specified month
		contribYear := contribution.ContributionDate.Year()
		contribMonth := int(contribution.ContributionDate.Month())
		
		if contribYear == year && contribMonth == month {
			totalContributions = totalContributions.Add(contribution.Amount)
		}
	}

	return totalContributions
}