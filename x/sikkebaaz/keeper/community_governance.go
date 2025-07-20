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
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/sikkebaaz/types"
)

// InitiateCommunityVeto creates a new community veto proposal
func (k Keeper) InitiateCommunityVeto(ctx sdk.Context, initiator, launchID, reason string) error {
	// Validate initiator has sufficient stake
	initiatorAddr, err := sdk.AccAddressFromBech32(initiator)
	if err != nil {
		return err
	}

	namoBalance := k.bankKeeper.GetBalance(ctx, initiatorAddr, types.DefaultDenom)
	minStake := sdk.NewInt(10000).MulRaw(1e12) // 10,000 NAMO required
	if namoBalance.Amount.LT(minStake) {
		return types.ErrInsufficientVotingPower
	}

	// Check if launch exists and is active
	launch, found := k.GetTokenLaunch(ctx, launchID)
	if !found {
		return types.ErrLaunchNotFound
	}

	if launch.Status != types.LaunchStatusActive {
		return types.ErrLaunchNotActive
	}

	// Check if veto already exists
	existingVeto := k.getCommunityVeto(ctx, launchID)
	if existingVeto != nil {
		return types.ErrCommunityVetoActive
	}

	// Create community veto
	veto := types.CommunityVeto{
		LaunchID:         launchID,
		InitiatedBy:      initiator,
		VoteStartTime:    ctx.BlockTime(),
		VoteEndTime:      ctx.BlockTime().Add(72 * time.Hour), // 72 hours voting period
		Votes:            make(map[string]bool),
		VotingPower:      make(map[string]sdk.Int),
		TotalVotingPower: sdk.ZeroInt(),
		VetoThreshold:    sdk.MustNewDecFromStr(types.CommunityVetoThreshold),
		Status:           "active",
		Reason:           reason,
	}

	k.setCommunityVeto(ctx, veto)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommunityVeto,
			sdk.NewAttribute("action", "initiated"),
			sdk.NewAttribute("launch_id", launchID),
			sdk.NewAttribute("initiator", initiator),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("voting_end_time", veto.VoteEndTime.String()),
		),
	)

	k.Logger(ctx).Info("Community veto initiated",
		"launch_id", launchID,
		"initiator", initiator,
		"reason", reason,
		"voting_ends", veto.VoteEndTime,
	)

	return nil
}

// VoteOnCommunityVeto allows community members to vote on a veto proposal
func (k Keeper) VoteOnCommunityVeto(ctx sdk.Context, voter, launchID string, vote bool) error {
	// Get community veto
	veto := k.getCommunityVeto(ctx, launchID)
	if veto == nil {
		return types.ErrCommunityVetoActive
	}

	// Check if voting period is active
	if ctx.BlockTime().After(veto.VoteEndTime) {
		return types.ErrVotingPeriodExpired
	}

	// Check if already voted
	if _, exists := veto.Votes[voter]; exists {
		return types.ErrAlreadyVoted
	}

	// Calculate voting power based on multiple factors
	votingPower := k.calculateVotingPower(ctx, voter, launchID)
	if votingPower.IsZero() {
		return types.ErrInsufficientVotingPower
	}

	// Record vote
	veto.Votes[voter] = vote
	veto.VotingPower[voter] = votingPower
	veto.TotalVotingPower = veto.TotalVotingPower.Add(votingPower)

	// Check if threshold reached
	k.checkVetoThreshold(ctx, veto)

	k.setCommunityVeto(ctx, *veto)

	// Emit voting event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommunityVeto,
			sdk.NewAttribute("action", "voted"),
			sdk.NewAttribute("launch_id", launchID),
			sdk.NewAttribute("voter", voter),
			sdk.NewAttribute("vote", func() string {
				if vote {
					return "yes"
				}
				return "no"
			}()),
			sdk.NewAttribute("voting_power", votingPower.String()),
			sdk.NewAttribute("total_power", veto.TotalVotingPower.String()),
			sdk.NewAttribute("status", veto.Status),
		),
	)

	return nil
}

// ProcessExpiredVetos processes expired veto proposals
func (k Keeper) ProcessExpiredVetos(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCommunityVeto)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var veto types.CommunityVeto
		k.cdc.MustUnmarshal(iterator.Value(), &veto)

		// Skip already processed vetos
		if veto.Status != "active" {
			continue
		}

		// Check if voting period expired
		if ctx.BlockTime().After(veto.VoteEndTime) {
			k.finalizeVeto(ctx, &veto)
		}
	}
}

// calculateVotingPower calculates voting power based on multiple factors
func (k Keeper) calculateVotingPower(ctx sdk.Context, voter, launchID string) sdk.Int {
	voterAddr, err := sdk.AccAddressFromBech32(voter)
	if err != nil {
		return sdk.ZeroInt()
	}

	// Base voting power from NAMO balance
	namoBalance := k.bankKeeper.GetBalance(ctx, voterAddr, types.DefaultDenom)
	basePower := namoBalance.Amount

	// Get launch details for regional weighting
	launch, found := k.GetTokenLaunch(ctx, launchID)
	if !found {
		return basePower
	}

	// Regional multiplier: higher power for same pincode
	regionalMultiplier := sdk.NewDec(1)
	voterPincode := k.getVoterPincode(ctx, voter)
	if voterPincode == launch.CreatorPincode {
		regionalMultiplier = sdk.MustNewDecFromStr("1.5") // 50% bonus for local community
	} else if k.isSameRegion(voterPincode, launch.CreatorPincode) {
		regionalMultiplier = sdk.MustNewDecFromStr("1.2") // 20% bonus for same region
	}

	// Stake duration multiplier: longer-term holders get more power
	stakeDuration := k.getStakeDuration(ctx, voter)
	durationMultiplier := k.calculateDurationMultiplier(stakeDuration)

	// Participation history multiplier
	participationMultiplier := k.calculateParticipationMultiplier(ctx, voter)

	// Apply all multipliers
	finalPower := basePower.ToDec().
		Mul(regionalMultiplier).
		Mul(durationMultiplier).
		Mul(participationMultiplier).
		TruncateInt()

	// Minimum voting power threshold
	minPower := sdk.NewInt(1000).MulRaw(1e12) // 1,000 NAMO minimum
	if finalPower.LT(minPower) {
		return sdk.ZeroInt()
	}

	return finalPower
}

// checkVetoThreshold checks if veto threshold has been reached
func (k Keeper) checkVetoThreshold(ctx sdk.Context, veto *types.CommunityVeto) {
	if veto.TotalVotingPower.IsZero() {
		return
	}

	// Calculate yes votes
	yesVotes := sdk.ZeroInt()
	for voter, vote := range veto.Votes {
		if vote {
			yesVotes = yesVotes.Add(veto.VotingPower[voter])
		}
	}

	// Calculate percentage
	yesPercentage := yesVotes.ToDec().Quo(veto.TotalVotingPower.ToDec())

	// Check if threshold reached
	if yesPercentage.GTE(veto.VetoThreshold) {
		veto.Status = "passed"
		k.executeVeto(ctx, veto)
	}
}

// finalizeVeto finalizes an expired veto proposal
func (k Keeper) finalizeVeto(ctx sdk.Context, veto *types.CommunityVeto) {
	if veto.Status != "active" {
		return
	}

	// Calculate final result
	yesVotes := sdk.ZeroInt()
	noVotes := sdk.ZeroInt()

	for voter, vote := range veto.Votes {
		power := veto.VotingPower[voter]
		if vote {
			yesVotes = yesVotes.Add(power)
		} else {
			noVotes = noVotes.Add(power)
		}
	}

	// Determine result
	if veto.TotalVotingPower.IsPositive() {
		yesPercentage := yesVotes.ToDec().Quo(veto.TotalVotingPower.ToDec())
		if yesPercentage.GTE(veto.VetoThreshold) {
			veto.Status = "passed"
			k.executeVeto(ctx, veto)
		} else {
			veto.Status = "failed"
		}
	} else {
		veto.Status = "failed" // No votes cast
	}

	k.setCommunityVeto(ctx, *veto)

	// Emit finalization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommunityVeto,
			sdk.NewAttribute("action", "finalized"),
			sdk.NewAttribute("launch_id", veto.LaunchID),
			sdk.NewAttribute("status", veto.Status),
			sdk.NewAttribute("yes_votes", yesVotes.String()),
			sdk.NewAttribute("no_votes", noVotes.String()),
			sdk.NewAttribute("total_power", veto.TotalVotingPower.String()),
		),
	)

	k.Logger(ctx).Info("Community veto finalized",
		"launch_id", veto.LaunchID,
		"status", veto.Status,
		"yes_votes", yesVotes,
		"no_votes", noVotes,
		"total_power", veto.TotalVotingPower,
	)
}

// executeVeto executes a passed veto proposal
func (k Keeper) executeVeto(ctx sdk.Context, veto *types.CommunityVeto) {
	// Get launch and cancel it
	launch, found := k.GetTokenLaunch(ctx, veto.LaunchID)
	if !found {
		k.Logger(ctx).Error("Launch not found for veto execution", "launch_id", veto.LaunchID)
		return
	}

	// Update launch status
	launch.Status = types.LaunchStatusVetoed
	launch.UpdatedAt = ctx.BlockTime()
	k.SetTokenLaunch(ctx, launch)

	// Process refunds to participants
	if err := k.processRefunds(ctx, &launch); err != nil {
		k.Logger(ctx).Error("Failed to process refunds during veto execution",
			"launch_id", veto.LaunchID,
			"error", err,
		)
	}

	// Distribute veto rewards to participants
	k.distributeVetoRewards(ctx, veto)

	// Emit execution event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommunityVeto,
			sdk.NewAttribute("action", "executed"),
			sdk.NewAttribute("launch_id", veto.LaunchID),
			sdk.NewAttribute("refund_amount", launch.RaisedAmount.String()),
		),
	)

	k.Logger(ctx).Info("Community veto executed",
		"launch_id", veto.LaunchID,
		"refund_amount", launch.RaisedAmount,
	)
}

// distributeVetoRewards distributes rewards to veto participants
func (k Keeper) distributeVetoRewards(ctx sdk.Context, veto *types.CommunityVeto) {
	// Calculate total reward pool (1% of launch target amount)
	launch, found := k.GetTokenLaunch(ctx, veto.LaunchID)
	if !found {
		return
	}

	rewardPool := launch.TargetAmount.MulRaw(1).QuoRaw(100) // 1% reward
	
	// Distribute rewards proportionally to yes voters
	yesVotes := sdk.ZeroInt()
	for voter, vote := range veto.Votes {
		if vote {
			yesVotes = yesVotes.Add(veto.VotingPower[voter])
		}
	}

	if yesVotes.IsZero() {
		return
	}

	// Distribute rewards
	for voter, vote := range veto.Votes {
		if !vote {
			continue // Only reward yes voters
		}

		voterPower := veto.VotingPower[voter]
		voterReward := rewardPool.Mul(voterPower).Quo(yesVotes)

		if voterReward.IsPositive() {
			voterAddr, err := sdk.AccAddressFromBech32(voter)
			if err != nil {
				continue
			}

			rewardCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, voterReward))
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CommunityIncentivePool, voterAddr, rewardCoins); err != nil {
				k.Logger(ctx).Error("Failed to distribute veto reward",
					"voter", voter,
					"reward", voterReward,
					"error", err,
				)
			}
		}
	}
}

// Helper functions for voting power calculation

func (k Keeper) getVoterPincode(ctx sdk.Context, voter string) string {
	// This would integrate with user profile system
	// For now, return empty string
	return ""
}

func (k Keeper) isSameRegion(pincode1, pincode2 string) bool {
	if len(pincode1) < 1 || len(pincode2) < 1 {
		return false
	}
	// Compare first digit for region matching
	return pincode1[0] == pincode2[0]
}

func (k Keeper) getStakeDuration(ctx sdk.Context, voter string) time.Duration {
	// This would track how long user has held NAMO tokens
	// For now, return default duration
	return 30 * 24 * time.Hour // 30 days default
}

func (k Keeper) calculateDurationMultiplier(duration time.Duration) sdk.Dec {
	days := duration.Hours() / 24
	
	switch {
	case days >= 365: // 1+ years
		return sdk.MustNewDecFromStr("2.0") // 2x multiplier
	case days >= 180: // 6+ months
		return sdk.MustNewDecFromStr("1.5") // 1.5x multiplier
	case days >= 90:  // 3+ months
		return sdk.MustNewDecFromStr("1.2") // 1.2x multiplier
	case days >= 30:  // 1+ month
		return sdk.MustNewDecFromStr("1.1") // 1.1x multiplier
	default:
		return sdk.OneDec() // No bonus
	}
}

func (k Keeper) calculateParticipationMultiplier(ctx sdk.Context, voter string) sdk.Dec {
	// Count previous governance participation
	participationCount := k.getGovernanceParticipationCount(ctx, voter)
	
	switch {
	case participationCount >= 10:
		return sdk.MustNewDecFromStr("1.3") // 30% bonus for active participants
	case participationCount >= 5:
		return sdk.MustNewDecFromStr("1.2") // 20% bonus
	case participationCount >= 1:
		return sdk.MustNewDecFromStr("1.1") // 10% bonus
	default:
		return sdk.OneDec() // No bonus for first-time voters
	}
}

func (k Keeper) getGovernanceParticipationCount(ctx sdk.Context, voter string) int {
	// This would track user's governance participation history
	// For now, return default count
	return 1
}

// Community governance helper functions

func (k Keeper) getCommunityVeto(ctx sdk.Context, launchID string) *types.CommunityVeto {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCommunityVetoKey(launchID))
	if bz == nil {
		return nil
	}

	var veto types.CommunityVeto
	k.cdc.MustUnmarshal(bz, &veto)
	return &veto
}

func (k Keeper) setCommunityVeto(ctx sdk.Context, veto types.CommunityVeto) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&veto)
	store.Set(types.GetCommunityVetoKey(veto.LaunchID), bz)
}

func (k Keeper) deleteCommunityVeto(ctx sdk.Context, launchID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetCommunityVetoKey(launchID))
}

// GetAllCommunityVetos returns all community vetos
func (k Keeper) GetAllCommunityVetos(ctx sdk.Context) []types.CommunityVeto {
	var vetos []types.CommunityVeto
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCommunityVeto)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var veto types.CommunityVeto
		k.cdc.MustUnmarshal(iterator.Value(), &veto)
		vetos = append(vetos, veto)
	}

	return vetos
}

// GetCommunityVetoByLaunchID returns community veto for a specific launch
func (k Keeper) GetCommunityVetoByLaunchID(ctx sdk.Context, launchID string) (types.CommunityVeto, bool) {
	veto := k.getCommunityVeto(ctx, launchID)
	if veto == nil {
		return types.CommunityVeto{}, false
	}
	return *veto, true
}