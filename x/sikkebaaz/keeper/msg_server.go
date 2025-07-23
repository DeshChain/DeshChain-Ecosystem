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
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/sikkebaaz/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateLaunch handles creating a new token launch
func (k msgServer) CreateLaunch(goCtx context.Context, msg *types.MsgCreateLaunch) (*types.MsgCreateLaunchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate unique launch ID
	launchID := k.generateLaunchID(ctx, msg.Creator, msg.TokenSymbol)

	// Create TokenLaunch struct
	launch := types.TokenLaunch{
		LaunchID:         launchID,
		Creator:          msg.Creator,
		CreatorPincode:   msg.CreatorPincode,
		TokenName:        msg.TokenName,
		TokenSymbol:      msg.TokenSymbol,
		TokenDescription: msg.TokenDescription,
		TotalSupply:      msg.TotalSupply,
		Decimals:         msg.Decimals,
		LaunchType:       msg.LaunchType,
		TargetAmount:     msg.TargetAmount,
		RaisedAmount:     sdk.ZeroInt(),
		MinContribution:  msg.MinContribution,
		MaxContribution:  msg.MaxContribution,
		StartTime:        msg.StartTime,
		EndTime:          msg.EndTime,
		TradingDelay:     msg.TradingDelay,
		AntiPumpConfig:   msg.AntiPumpConfig,
		CulturalQuote:    msg.CulturalQuote,
		FestivalBonus:    false,
		PatriotismScore:  0,
		Status:           types.LaunchStatusPending,
		ParticipantCount: 0,
		Whitelist:        msg.Whitelist,
		Metadata:         make(map[string]string),
	}

	// Validate cultural requirements and calculate patriotism score
	if err := k.ValidateCulturalRequirements(ctx, &launch); err != nil {
		return nil, err
	}

	// Check seasonal restrictions
	if err := k.ValidateSeasonalRestrictions(ctx, &launch); err != nil {
		return nil, err
	}

	// Apply festival bonuses if applicable
	k.ApplyFestivalBonuses(ctx, &launch)

	// Create the launch
	if err := k.CreateTokenLaunch(ctx, msg.Creator, &launch); err != nil {
		return nil, err
	}

	// Emit detailed event with cultural information
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenLaunched,
			sdk.NewAttribute(types.AttributeKeyLaunchID, launchID),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, msg.TokenSymbol),
			sdk.NewAttribute(types.AttributeKeyPincode, msg.CreatorPincode),
			sdk.NewAttribute(types.AttributeKeyLaunchType, msg.LaunchType),
			sdk.NewAttribute(types.AttributeKeyTargetAmount, msg.TargetAmount.String()),
			sdk.NewAttribute("patriotism_score", strconv.FormatInt(launch.PatriotismScore, 10)),
			sdk.NewAttribute("festival_bonus", strconv.FormatBool(launch.FestivalBonus)),
			sdk.NewAttribute("cultural_quote", launch.CulturalQuote),
		),
	)

	return &types.MsgCreateLaunchResponse{
		LaunchID: launchID,
		Status:   launch.Status,
	}, nil
}

// ParticipateLaunch handles participation in a token launch
func (k msgServer) ParticipateLaunch(goCtx context.Context, msg *types.MsgParticipateLaunch) (*types.MsgParticipateLaunchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Participate in the launch
	if err := k.Keeper.ParticipateInLaunch(ctx, msg.Participant, msg.LaunchID, msg.Amount); err != nil {
		return nil, err
	}

	// Get updated launch details
	launch, found := k.GetTokenLaunch(ctx, msg.LaunchID)
	if !found {
		return nil, types.ErrLaunchNotFound
	}

	// Calculate allocated tokens
	tokenAllocation := launch.CalculateTokenAllocation(msg.Amount)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"launch_participation",
			sdk.NewAttribute("launch_id", msg.LaunchID),
			sdk.NewAttribute("participant", msg.Participant),
			sdk.NewAttribute("contribution", msg.Amount.String()),
			sdk.NewAttribute("token_allocation", tokenAllocation.String()),
			sdk.NewAttribute("total_raised", launch.RaisedAmount.String()),
			sdk.NewAttribute("participant_count", strconv.FormatUint(launch.ParticipantCount, 10)),
		),
	)

	return &types.MsgParticipateLaunchResponse{
		TokensAllocated: tokenAllocation,
		TotalRaised:     launch.RaisedAmount,
	}, nil
}

// CompleteLaunch handles manual completion of a launch
func (k msgServer) CompleteLaunch(goCtx context.Context, msg *types.MsgCompleteLaunch) (*types.MsgCompleteLaunchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get launch
	launch, found := k.GetTokenLaunch(ctx, msg.LaunchID)
	if !found {
		return nil, types.ErrLaunchNotFound
	}

	// Verify creator
	if launch.Creator != msg.Creator {
		return nil, types.ErrNotCreator
	}

	// Check if launch can be completed
	if launch.Status != types.LaunchStatusActive {
		return nil, types.ErrLaunchNotActive
	}

	// Check if minimum target reached (flexible completion)
	minTarget := launch.TargetAmount.MulRaw(50).QuoRaw(100) // 50% minimum
	if launch.RaisedAmount.LT(minTarget) {
		return nil, types.ErrTargetReached
	}

	// Update launch status
	launch.Status = types.LaunchStatusSuccessful
	completedAt := ctx.BlockTime()
	launch.CompletedAt = &completedAt
	launch.UpdatedAt = ctx.BlockTime()

	// Deploy token with anti-pump features
	if err := k.deployToken(ctx, &launch); err != nil {
		return nil, err
	}

	k.SetTokenLaunch(ctx, launch)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLaunchCompleted,
			sdk.NewAttribute(types.AttributeKeyLaunchID, msg.LaunchID),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute("final_raised", launch.RaisedAmount.String()),
			sdk.NewAttribute("completion_type", "manual"),
		),
	)

	return &types.MsgCompleteLaunchResponse{
		Status:      launch.Status,
		TokenSymbol: launch.TokenSymbol,
		FinalRaised: launch.RaisedAmount,
	}, nil
}

// CancelLaunch handles cancellation of a launch
func (k msgServer) CancelLaunch(goCtx context.Context, msg *types.MsgCancelLaunch) (*types.MsgCancelLaunchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get launch
	launch, found := k.GetTokenLaunch(ctx, msg.LaunchID)
	if !found {
		return nil, types.ErrLaunchNotFound
	}

	// Verify creator
	if launch.Creator != msg.Creator {
		return nil, types.ErrNotCreator
	}

	// Check if launch can be cancelled
	if launch.Status == types.LaunchStatusSuccessful || launch.Status == types.LaunchStatusCancelled {
		return nil, types.ErrLaunchAlreadyCompleted
	}

	// Update launch status
	launch.Status = types.LaunchStatusCancelled
	launch.UpdatedAt = ctx.BlockTime()

	// Process refunds to participants
	if err := k.processRefunds(ctx, &launch); err != nil {
		k.Logger(ctx).Error("Failed to process refunds", "launch_id", msg.LaunchID, "error", err)
	}

	k.SetTokenLaunch(ctx, launch)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLaunchFailed,
			sdk.NewAttribute(types.AttributeKeyLaunchID, msg.LaunchID),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("refund_amount", launch.RaisedAmount.String()),
		),
	)

	return &types.MsgCancelLaunchResponse{
		Status:       launch.Status,
		RefundAmount: launch.RaisedAmount,
	}, nil
}

// ClaimTokens handles claiming allocated tokens after successful launch
func (k msgServer) ClaimTokens(goCtx context.Context, msg *types.MsgClaimTokens) (*types.MsgClaimTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get launch
	launch, found := k.GetTokenLaunch(ctx, msg.LaunchID)
	if !found {
		return nil, types.ErrLaunchNotFound
	}

	// Check if launch is successful
	if launch.Status != types.LaunchStatusSuccessful {
		return nil, types.ErrLaunchNotActive
	}

	// Check if trading delay has passed
	if err := k.CheckTradingDelay(ctx, launch.TokenSymbol); err != nil {
		return nil, err
	}

	// Get participation record
	participation, found := k.getLaunchParticipation(ctx, msg.Participant, msg.LaunchID)
	if !found {
		return nil, types.ErrNotWhitelisted
	}

	// Check if already claimed
	if participation.TokensClaimed.Equal(participation.TokensAllocated) {
		return nil, types.ErrRewardAlreadyClaimed
	}

	// Calculate claimable amount
	claimableAmount := participation.TokensAllocated.Sub(participation.TokensClaimed)

	// Transfer tokens from module to participant
	participantAddr, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		return nil, err
	}

	tokenCoins := sdk.NewCoins(sdk.NewCoin(launch.TokenSymbol, claimableAmount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, participantAddr, tokenCoins); err != nil {
		return nil, err
	}

	// Update participation record
	participation.TokensClaimed = participation.TokensAllocated
	claimedAt := ctx.BlockTime()
	participation.ClaimedAt = &claimedAt

	k.setLaunchParticipation(ctx, participation)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"tokens_claimed",
			sdk.NewAttribute("launch_id", msg.LaunchID),
			sdk.NewAttribute("participant", msg.Participant),
			sdk.NewAttribute("claimed_amount", claimableAmount.String()),
			sdk.NewAttribute("token_symbol", launch.TokenSymbol),
		),
	)

	return &types.MsgClaimTokensResponse{
		ClaimedAmount: claimableAmount,
		TokenSymbol:   launch.TokenSymbol,
	}, nil
}

// InitiateCommunityVeto handles initiation of community veto process
func (k msgServer) InitiateCommunityVeto(goCtx context.Context, msg *types.MsgInitiateCommunityVeto) (*types.MsgInitiateCommunityVetoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get launch
	launch, found := k.GetTokenLaunch(ctx, msg.LaunchID)
	if !found {
		return nil, types.ErrLaunchNotFound
	}

	// Check if launch is active
	if launch.Status != types.LaunchStatusActive {
		return nil, types.ErrLaunchNotActive
	}

	// Check if initiator has sufficient stake (simplified)
	initiatorAddr, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return nil, err
	}

	namoBalance := k.bankKeeper.GetBalance(ctx, initiatorAddr, types.DefaultDenom)
	minStake := sdk.NewInt(10000).MulRaw(1e12) // 10,000 NAMO
	if namoBalance.Amount.LT(minStake) {
		return nil, types.ErrInsufficientVotingPower
	}

	// Check if veto already exists
	if k.getCommunityVeto(ctx, msg.LaunchID) != nil {
		return nil, types.ErrCommunityVetoActive
	}

	// Create community veto
	veto := types.CommunityVeto{
		LaunchID:      msg.LaunchID,
		InitiatedBy:   msg.Initiator,
		VoteStartTime: ctx.BlockTime(),
		VoteEndTime:   ctx.BlockTime().Add(72 * time.Hour), // 72 hours voting period
		Votes:         make(map[string]bool),
		VotingPower:   make(map[string]sdk.Int),
		TotalVotingPower: sdk.ZeroInt(),
		VetoThreshold: sdk.MustNewDecFromStr(types.CommunityVetoThreshold),
		Status:        "active",
		Reason:        msg.Reason,
	}

	k.setCommunityVeto(ctx, veto)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommunityVeto,
			sdk.NewAttribute("launch_id", msg.LaunchID),
			sdk.NewAttribute("initiator", msg.Initiator),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("voting_end_time", veto.VoteEndTime.String()),
		),
	)

	return &types.MsgInitiateCommunityVetoResponse{
		VetoID:      msg.LaunchID,
		VotingEnds:  veto.VoteEndTime,
		Status:      veto.Status,
	}, nil
}

// VoteCommunityVeto handles voting on community veto
func (k msgServer) VoteCommunityVeto(goCtx context.Context, msg *types.MsgVoteCommunityVeto) (*types.MsgVoteCommunityVetoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get community veto
	veto := k.getCommunityVeto(ctx, msg.LaunchID)
	if veto == nil {
		return nil, types.ErrCommunityVetoActive
	}

	// Check if voting period is active
	if ctx.BlockTime().After(veto.VoteEndTime) {
		return nil, types.ErrVotingPeriodExpired
	}

	// Check if already voted
	if _, exists := veto.Votes[msg.Voter]; exists {
		return nil, types.ErrAlreadyVoted
	}

	// Calculate voting power based on NAMO balance
	voterAddr, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	namoBalance := k.bankKeeper.GetBalance(ctx, voterAddr, types.DefaultDenom)
	votingPower := namoBalance.Amount

	// Record vote
	veto.Votes[msg.Voter] = msg.Vote
	veto.VotingPower[msg.Voter] = votingPower
	veto.TotalVotingPower = veto.TotalVotingPower.Add(votingPower)

	// Check if threshold reached
	yesVotes := sdk.ZeroInt()
	for voter, vote := range veto.Votes {
		if vote {
			yesVotes = yesVotes.Add(veto.VotingPower[voter])
		}
	}

	// Calculate yes percentage
	if veto.TotalVotingPower.IsPositive() {
		yesPercentage := yesVotes.ToDec().Quo(veto.TotalVotingPower.ToDec())
		if yesPercentage.GTE(veto.VetoThreshold) {
			veto.Status = "passed"
			
			// Cancel the launch
			launch, _ := k.GetTokenLaunch(ctx, msg.LaunchID)
			launch.Status = types.LaunchStatusVetoed
			k.SetTokenLaunch(ctx, launch)
		}
	}

	k.setCommunityVeto(ctx, *veto)

	return &types.MsgVoteCommunityVetoResponse{
		VoteRecorded: true,
		VotingPower:  votingPower,
		Status:       veto.Status,
	}, nil
}

// ClaimCreatorReward handles claiming creator trading rewards
func (k msgServer) ClaimCreatorReward(goCtx context.Context, msg *types.MsgClaimCreatorReward) (*types.MsgClaimCreatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get creator reward
	reward, found := k.getCreatorReward(ctx, msg.Creator, msg.TokenAddress)
	if !found {
		return nil, types.ErrCreatorRewardNotFound
	}

	// Check if there are rewards to claim
	if reward.AccumulatedReward.IsZero() {
		return nil, types.ErrNoRewardsToClaim
	}

	// Transfer rewards to creator
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, reward.AccumulatedReward))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CreatorRewardsPool, creatorAddr, rewardCoins); err != nil {
		return nil, err
	}

	// Update reward record
	reward.TotalClaimed = reward.TotalClaimed.Add(reward.AccumulatedReward)
	claimedAmount := reward.AccumulatedReward
	reward.AccumulatedReward = sdk.ZeroInt()
	reward.LastClaimedAt = ctx.BlockTime()

	k.setCreatorReward(ctx, reward)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatorReward,
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("token_address", msg.TokenAddress),
			sdk.NewAttribute("claimed_amount", claimedAmount.String()),
			sdk.NewAttribute("total_claimed", reward.TotalClaimed.String()),
		),
	)

	return &types.MsgClaimCreatorRewardResponse{
		ClaimedAmount: claimedAmount,
		TotalClaimed:  reward.TotalClaimed,
	}, nil
}

// UpdateAntiPumpConfig handles updating anti-pump configuration
func (k msgServer) UpdateAntiPumpConfig(goCtx context.Context, msg *types.MsgUpdateAntiPumpConfig) (*types.MsgUpdateAntiPumpConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get launch by token address
	launch, found := k.getTokenLaunchByAddress(ctx, msg.TokenAddress)
	if !found {
		return nil, types.ErrTokenAlreadyExists
	}

	// Verify creator
	if launch.Creator != msg.Creator {
		return nil, types.ErrNotCreator
	}

	// Update anti-pump config (only allow more restrictive settings)
	oldConfig := launch.AntiPumpConfig
	newConfig := msg.AntiPumpConfig

	// Validate that new config is more restrictive
	if newConfig.MaxWalletPercent24h > oldConfig.MaxWalletPercent24h ||
		newConfig.MaxWalletPercentAfter > oldConfig.MaxWalletPercentAfter ||
		newConfig.LiquidityLockDays < oldConfig.LiquidityLockDays {
		return nil, types.ErrIncompatibleSettings
	}

	launch.AntiPumpConfig = newConfig
	launch.UpdatedAt = ctx.BlockTime()
	k.SetTokenLaunch(ctx, launch)

	return &types.MsgUpdateAntiPumpConfigResponse{
		Success: true,
		UpdatedConfig: newConfig,
	}, nil
}

// EmergencyStop handles emergency stop for problematic tokens
func (k msgServer) EmergencyStop(goCtx context.Context, msg *types.MsgEmergencyStop) (*types.MsgEmergencyStopResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority (this should be governance or specific authority address)
	// For now, simplified to allow any address with sufficient NAMO balance
	authorityAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	namoBalance := k.bankKeeper.GetBalance(ctx, authorityAddr, types.DefaultDenom)
	minAuthority := sdk.NewInt(100000).MulRaw(1e12) // 100,000 NAMO
	if namoBalance.Amount.LT(minAuthority) {
		return nil, types.ErrUnauthorizedAccess
	}

	// Create emergency control
	emergencyControl := types.EmergencyControl{
		TokenAddress: msg.TokenAddress,
		ControlType:  "emergency_stop",
		InitiatedBy:  msg.Authority,
		Reason:       msg.Reason,
		ActivatedAt:  ctx.BlockTime(),
		ExpiresAt:    nil, // Indefinite until resolved
		IsActive:     true,
		Metadata:     make(map[string]string),
	}

	k.setEmergencyControl(ctx, emergencyControl)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"emergency_stop",
			sdk.NewAttribute("token_address", msg.TokenAddress),
			sdk.NewAttribute("authority", msg.Authority),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)

	return &types.MsgEmergencyStopResponse{
		Success:   true,
		StoppedAt: ctx.BlockTime(),
	}, nil
}

// Helper functions

func (k msgServer) generateLaunchID(ctx sdk.Context, creator, symbol string) string {
	return fmt.Sprintf("%s-%s-%d", creator[:8], symbol, ctx.BlockHeight())
}

func (k msgServer) processRefunds(ctx sdk.Context, launch *types.TokenLaunch) error {
	// Implementation would iterate through participants and process refunds
	// Simplified for now
	return nil
}

func (k msgServer) getLaunchParticipation(ctx sdk.Context, participant, launchID string) (types.LaunchParticipation, bool) {
	return k.Keeper.getLaunchParticipation(ctx, participant, launchID)
}

func (k msgServer) getCommunityVeto(ctx sdk.Context, launchID string) *types.CommunityVeto {
	// Implementation would retrieve community veto
	// Simplified for now
	return nil
}

func (k msgServer) setCommunityVeto(ctx sdk.Context, veto types.CommunityVeto) {
	// Implementation would store community veto
	// Simplified for now
}

func (k msgServer) getCreatorReward(ctx sdk.Context, creator, tokenAddress string) (types.CreatorReward, bool) {
	return k.Keeper.getCreatorReward(ctx, creator, tokenAddress)
}

func (k msgServer) setEmergencyControl(ctx sdk.Context, control types.EmergencyControl) {
	// Implementation would store emergency control
	// Simplified for now
}