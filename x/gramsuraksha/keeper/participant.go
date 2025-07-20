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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/gramsuraksha/types"
)

// EnrollParticipant enrolls a new participant in a pension scheme
func (k Keeper) EnrollParticipant(ctx sdk.Context, participant types.SurakshaParticipant) error {
	// Get scheme
	scheme, found := k.GetScheme(ctx, participant.SchemeID)
	if !found {
		return types.ErrSchemeNotFound
	}

	// Check if scheme is active
	if scheme.Status != types.StatusActive {
		return types.ErrSchemeInactive
	}

	// Check eligibility
	if !scheme.IsEligible(participant.Age) {
		return types.ErrAgeNotEligible
	}

	// Check KYC if required
	if scheme.KYCRequired {
		if k.kycKeeper != nil && !k.kycKeeper.IsKYCVerified(ctx, participant.Address) {
			return types.ErrKYCNotVerified
		}
	}

	// Check if already enrolled
	if k.IsParticipantEnrolled(ctx, participant.Address, participant.SchemeID) {
		return types.ErrAlreadyEnrolled
	}

	// Validate participant
	if err := participant.Validate(); err != nil {
		return err
	}

	// Set enrollment details
	participant.EnrollmentDate = ctx.BlockTime()
	participant.MaturityDate = participant.EnrollmentDate.AddDate(0, int(scheme.ContributionPeriod), 0)
	participant.Status = types.StatusActive
	participant.TotalContributed = sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())
	participant.BonusEarned = sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())
	participant.PenaltyIncurred = sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())
	participant.ReferralRewards = sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())
	participant.LiquidityContributed = sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())

	// Handle referral
	if !participant.ReferrerAddress.Empty() {
		referrer, found := k.GetParticipantByAddress(ctx, participant.ReferrerAddress, participant.SchemeID)
		if found && referrer.Status == types.StatusActive {
			// Will process referral rewards when first contribution is made
		}
	}

	// Update scheme participant count
	scheme.CurrentParticipants++
	k.SetScheme(ctx, scheme)

	// Store participant
	k.SetParticipant(ctx, participant)

	// Create indexes
	k.SetParticipantByAddress(ctx, participant.Address, participant.SchemeID, participant.ParticipantID)
	k.SetParticipantByScheme(ctx, participant.SchemeID, participant.ParticipantID)

	// Get cultural quote if enabled
	culturalQuote := ""
	if scheme.CulturalIntegration && k.culturalKeeper != nil {
		quote, _ := k.culturalKeeper.GetRandomQuote(ctx, "enrollment")
		culturalQuote = quote
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeParticipantEnrolled,
			sdk.NewAttribute(types.AttributeKeyParticipantID, participant.ParticipantID),
			sdk.NewAttribute(types.AttributeKeySchemeID, participant.SchemeID),
			sdk.NewAttribute(types.AttributeKeyParticipantAddress, participant.Address.String()),
			sdk.NewAttribute(types.AttributeKeyCulturalQuoteText, culturalQuote),
		),
	)

	return nil
}

// SetParticipant stores a participant
func (k Keeper) SetParticipant(ctx sdk.Context, participant types.SurakshaParticipant) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&participant)
	store.Set(types.ParticipantPrefix.Bytes(append([]byte(participant.ParticipantID)), bz)
}

// GetParticipant retrieves a participant by ID
func (k Keeper) GetParticipant(ctx sdk.Context, participantID string) (types.SurakshaParticipant, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParticipantPrefix.Bytes(append([]byte(participantID)))
	if bz == nil {
		return types.SurakshaParticipant{}, false
	}

	var participant types.SurakshaParticipant
	k.cdc.MustUnmarshal(bz, &participant)
	return participant, true
}

// GetParticipantByAddress retrieves a participant by address and scheme
func (k Keeper) GetParticipantByAddress(ctx sdk.Context, address sdk.AccAddress, schemeID string) (types.SurakshaParticipant, bool) {
	participantID, found := k.GetParticipantIDByAddress(ctx, address, schemeID)
	if !found {
		return types.SurakshaParticipant{}, false
	}
	return k.GetParticipant(ctx, participantID)
}

// IsParticipantEnrolled checks if an address is enrolled in a scheme
func (k Keeper) IsParticipantEnrolled(ctx sdk.Context, address sdk.AccAddress, schemeID string) bool {
	_, found := k.GetParticipantIDByAddress(ctx, address, schemeID)
	return found
}

// SetParticipantByAddress creates an index by address
func (k Keeper) SetParticipantByAddress(ctx sdk.Context, address sdk.AccAddress, schemeID, participantID string) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.ParticipantByAddressPrefix, address.Bytes()...), []byte(schemeID)...)
	store.Set(key, []byte(participantID))
}

// GetParticipantIDByAddress gets participant ID by address
func (k Keeper) GetParticipantIDByAddress(ctx sdk.Context, address sdk.AccAddress, schemeID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.ParticipantByAddressPrefix, address.Bytes()...), []byte(schemeID)...)
	bz := store.Get(key)
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// SetParticipantByScheme creates an index by scheme
func (k Keeper) SetParticipantByScheme(ctx sdk.Context, schemeID, participantID string) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.ParticipantBySchemePrefix, []byte(schemeID)...), []byte(participantID)...)
	store.Set(key, []byte{1})
}

// IterateParticipantsByScheme iterates participants in a scheme
func (k Keeper) IterateParticipantsByScheme(ctx sdk.Context, schemeID string, cb func(types.SurakshaParticipant) bool) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.ParticipantBySchemePrefix, []byte(schemeID)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Extract participant ID from key
		fullKey := iterator.Key()
		participantID := string(fullKey[len(prefix):])
		
		participant, found := k.GetParticipant(ctx, participantID)
		if found && cb(participant) {
			break
		}
	}
}

// UpdateParticipantStatus updates participant status
func (k Keeper) UpdateParticipantStatus(ctx sdk.Context, participantID string, status string) error {
	participant, found := k.GetParticipant(ctx, participantID)
	if !found {
		return types.ErrParticipantNotFound
	}

	oldStatus := participant.Status
	participant.Status = status

	k.SetParticipant(ctx, participant)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeParticipantUpdated,
			sdk.NewAttribute(types.AttributeKeyParticipantID, participantID),
			sdk.NewAttribute("old_status", oldStatus),
			sdk.NewAttribute(types.AttributeKeyParticipantStatus, status),
		),
	)

	return nil
}

// GetActiveParticipants returns all active participants in a scheme
func (k Keeper) GetActiveParticipants(ctx sdk.Context, schemeID string) []types.SurakshaParticipant {
	var participants []types.SurakshaParticipant
	k.IterateParticipantsByScheme(ctx, schemeID, func(participant types.SurakshaParticipant) bool {
		if participant.Status == types.StatusActive {
			participants = append(participants, participant)
		}
		return false
	})
	return participants
}

// CheckParticipantDefaults checks for defaulted participants
func (k Keeper) CheckParticipantDefaults(ctx sdk.Context) {
	currentTime := ctx.BlockTime()
	
	// Iterate all active participants
	k.IterateSchemes(ctx, func(scheme types.SurakshaScheme) bool {
		if scheme.Status != types.StatusActive {
			return false
		}

		k.IterateParticipantsByScheme(ctx, scheme.SchemeID, func(participant types.SurakshaParticipant) bool {
			if participant.Status != types.StatusActive {
				return false
			}

			// Calculate expected contributions
			monthsSinceEnrollment := int(currentTime.Sub(participant.EnrollmentDate).Hours() / 24 / 30)
			expectedContributions := uint32(monthsSinceEnrollment)
			if expectedContributions > scheme.ContributionPeriod {
				expectedContributions = scheme.ContributionPeriod
			}

			actualContributions := participant.OnTimePayments + participant.MissedPayments

			// Check if too many payments missed
			if actualContributions < expectedContributions {
				missedPayments := expectedContributions - actualContributions
				if missedPayments > 3 { // More than 3 missed payments = default
					k.UpdateParticipantStatus(ctx, participant.ParticipantID, types.StatusDefaulted)
					
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(
							types.EventTypeParticipantSuspended,
							sdk.NewAttribute(types.AttributeKeyParticipantID, participant.ParticipantID),
							sdk.NewAttribute("missed_payments", fmt.Sprintf("%d", missedPayments)),
						),
					)
				}
			}
			
			return false
		})
		
		return false
	})
}

// GenerateParticipantID generates a unique participant ID
func (k Keeper) GenerateParticipantID(ctx sdk.Context, address sdk.AccAddress, schemeID string) string {
	timestamp := ctx.BlockTime().Unix()
	addrStr := address.String()
	if len(addrStr) > 8 {
		addrStr = addrStr[:8]
	}
	return fmt.Sprintf("PART-%s-%s-%d", schemeID[:min(8, len(schemeID))], addrStr, timestamp)
}