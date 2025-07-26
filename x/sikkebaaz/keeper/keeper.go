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

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/DeshChain/DeshChain-Ecosystem/x/sikkebaaz/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	moneyordertypes "github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	bankKeeper       types.BankKeeper
	accountKeeper    types.AccountKeeper
	culturalKeeper   types.CulturalKeeper
	treasuryKeeper   types.TreasuryKeeper
	revenueKeeper    types.RevenueKeeper
	moneyOrderKeeper types.MoneyOrderKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	culturalKeeper types.CulturalKeeper,
	treasuryKeeper types.TreasuryKeeper,
	revenueKeeper types.RevenueKeeper,
	moneyOrderKeeper types.MoneyOrderKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		culturalKeeper:   culturalKeeper,
		treasuryKeeper:   treasuryKeeper,
		revenueKeeper:    revenueKeeper,
		moneyOrderKeeper: moneyOrderKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CreateTokenLaunch creates a new token launch with anti-pump features
func (k Keeper) CreateTokenLaunch(ctx sdk.Context, creator string, msg *types.TokenLaunch) error {
	// Validate token launch
	if err := types.ValidateTokenLaunch(*msg); err != nil {
		return err
	}

	// Validate anti-pump configuration
	if err := types.ValidateAntiPumpConfig(msg.AntiPumpConfig); err != nil {
		return err
	}

	// Validate creator's patriotism score
	if msg.PatriotismScore < types.MinPatriotismScore {
		return types.ErrInsufficientPatriotismScore
	}

	// Validate PIN code format
	if len(msg.CreatorPincode) != types.MaxPincodeLength {
		return types.ErrInvalidPincode
	}

	// Calculate launch fee
	launchFee := types.CalculateLaunchFee(msg.TargetAmount)
	msg.LaunchFee = launchFee

	// Calculate charity allocation
	charityAllocation := types.CalculateCharityAllocation(msg.TargetAmount)
	msg.CharityAllocation = charityAllocation

	// Charge launch fee from creator
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return err
	}

	feeCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, launchFee))
	// Use revenue keeper to collect launch fee - handles collection and distribution
	if err := k.revenueKeeper.CollectLaunchFee(ctx, types.ModuleName, creatorAddr, feeCoins, msg.TokenSymbol); err != nil {
		return types.ErrInsufficientFees
	}

	// Apply festival bonus if active
	k.applyFestivalBonus(ctx, msg)

	// Set timestamps
	msg.CreatedAt = ctx.BlockTime()
	msg.UpdatedAt = ctx.BlockTime()
	msg.Status = types.LaunchStatusPending

	// Store the launch
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(msg)
	store.Set(types.GetTokenLaunchKey(msg.LaunchID), bz)

	// Index by creator
	k.addCreatorLaunch(ctx, creator, msg.LaunchID)

	// Index by pincode
	k.addPincodeLaunch(ctx, msg.CreatorPincode, msg.LaunchID)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenLaunched,
			sdk.NewAttribute(types.AttributeKeyLaunchID, msg.LaunchID),
			sdk.NewAttribute(types.AttributeKeyCreator, creator),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, msg.TokenSymbol),
			sdk.NewAttribute(types.AttributeKeyPincode, msg.CreatorPincode),
			sdk.NewAttribute(types.AttributeKeyLaunchType, msg.LaunchType),
			sdk.NewAttribute(types.AttributeKeyTargetAmount, msg.TargetAmount.String()),
		),
	)

	return nil
}

// GetTokenLaunch retrieves a token launch by ID
func (k Keeper) GetTokenLaunch(ctx sdk.Context, launchID string) (types.TokenLaunch, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTokenLaunchKey(launchID))
	if bz == nil {
		return types.TokenLaunch{}, false
	}

	var launch types.TokenLaunch
	k.cdc.MustUnmarshal(bz, &launch)
	return launch, true
}

// SetTokenLaunch stores a token launch
func (k Keeper) SetTokenLaunch(ctx sdk.Context, launch types.TokenLaunch) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&launch)
	store.Set(types.GetTokenLaunchKey(launch.LaunchID), bz)
}

// ParticipateInLaunch allows users to participate in a token launch
func (k Keeper) ParticipateInLaunch(ctx sdk.Context, participant, launchID string, contribution sdk.Int) error {
	// Get launch
	launch, found := k.GetTokenLaunch(ctx, launchID)
	if !found {
		return types.ErrLaunchNotFound
	}

	// Check if launch is active
	if !launch.IsLaunchActive(ctx.BlockTime()) {
		return types.ErrLaunchNotActive
	}

	// Validate contribution amount
	if contribution.LT(launch.MinContribution) {
		return types.ErrContributionTooLow
	}
	if contribution.GT(launch.MaxContribution) {
		return types.ErrContributionTooHigh
	}

	// Check if target would be exceeded
	newTotal := launch.RaisedAmount.Add(contribution)
	if newTotal.GT(launch.TargetAmount) {
		return types.ErrTargetReached
	}

	// Check whitelist if applicable
	if launch.LaunchType == types.LaunchTypeWhitelist {
		if !k.isWhitelisted(participant, launch.Whitelist) {
			return types.ErrNotWhitelisted
		}
	}

	// Check if already participated
	if k.hasParticipated(ctx, participant, launchID) {
		return types.ErrAlreadyParticipated
	}

	// Transfer contribution to escrow
	participantAddr, err := sdk.AccAddressFromBech32(participant)
	if err != nil {
		return err
	}

	contributionCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, contribution))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, participantAddr, types.LaunchEscrowAccount, contributionCoins); err != nil {
		return types.ErrInsufficientBalance
	}

	// Calculate token allocation
	tokenAllocation := launch.CalculateTokenAllocation(contribution)

	// Create participation record
	participation := types.LaunchParticipation{
		LaunchID:          launchID,
		Participant:       participant,
		ContributedAmount: contribution,
		TokensAllocated:   tokenAllocation,
		TokensClaimed:     sdk.ZeroInt(),
		ParticipatedAt:    ctx.BlockTime(),
		IsRefunded:        false,
	}

	// Store participation
	k.setLaunchParticipation(ctx, participation)

	// Update launch totals
	launch.RaisedAmount = launch.RaisedAmount.Add(contribution)
	launch.ParticipantCount++
	launch.UpdatedAt = ctx.BlockTime()

	// Check if target reached
	if launch.RaisedAmount.GTE(launch.TargetAmount) {
		launch.Status = types.LaunchStatusSuccessful
		launch.CompletedAt = &ctx.BlockTime()
		
		// Deploy token with anti-pump features
		if err := k.deployToken(ctx, &launch); err != nil {
			return err
		}
	}

	k.SetTokenLaunch(ctx, launch)

	return nil
}

// deployToken deploys the actual token with anti-pump protection
func (k Keeper) deployToken(ctx sdk.Context, launch *types.TokenLaunch) error {
	// Create token metadata for bank module
	metadata := banktypes.Metadata{
		Description: launch.TokenDescription,
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    launch.TokenSymbol,
				Exponent: 0,
				Aliases:  []string{},
			},
			{
				Denom:    launch.TokenSymbol,
				Exponent: launch.Decimals,
				Aliases:  []string{},
			},
		},
		Base:    launch.TokenSymbol,
		Display: launch.TokenSymbol,
		Name:    launch.TokenName,
		Symbol:  launch.TokenSymbol,
	}

	// Set token metadata
	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	// Mint initial supply to module account for distribution
	tokenCoins := sdk.NewCoins(sdk.NewCoin(launch.TokenSymbol, launch.TotalSupply))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, tokenCoins); err != nil {
		return err
	}

	// Lock liquidity (80% minimum)
	liquidityAmount := launch.RaisedAmount.MulRaw(int64(launch.AntiPumpConfig.MinLiquidityPercent)).QuoRaw(100)
	liquidityLock := types.LiquidityLock{
		TokenAddress:   launch.TokenSymbol,
		LockOwner:      launch.Creator,
		LPTokenAddress: launch.TokenSymbol + "-NAMO-LP", // LP token address
		LockedAmount:   liquidityAmount,
		LockDate:       ctx.BlockTime(),
		UnlockDate:     ctx.BlockTime().AddDate(0, 0, int(launch.AntiPumpConfig.LiquidityLockDays)),
		IsWithdrawn:    false,
	}

	k.setLiquidityLock(ctx, liquidityLock)

	// Set wallet limits for all participants
	k.initializeWalletLimits(ctx, launch)

	// Initialize creator rewards
	creatorReward := types.CreatorReward{
		Creator:           launch.Creator,
		TokenAddress:      launch.TokenSymbol,
		RewardRate:        sdk.MustNewDecFromStr(types.CreatorTradingReward),
		AccumulatedReward: sdk.ZeroInt(),
		LastClaimedAt:     ctx.BlockTime(),
		TotalClaimed:      sdk.ZeroInt(),
		IsActive:          true,
	}

	k.setCreatorReward(ctx, creatorReward)

	// Create liquidity pool on MoneyOrder DEX
	if err := k.createLiquidityPool(ctx, launch); err != nil {
		k.Logger(ctx).Error("Failed to create liquidity pool", "error", err, "token", launch.TokenSymbol)
		// Don't fail deployment, pool can be created manually later
	}

	// Distribute fees according to DeshChain Platform Revenue Model
	return k.distributeLaunchFees(ctx, launch)
}

// initializeWalletLimits sets up wallet limits for anti-pump protection
func (k Keeper) initializeWalletLimits(ctx sdk.Context, launch *types.TokenLaunch) {
	// Get current wallet limit percentage
	currentLimit := launch.GetCurrentWalletLimit(ctx.BlockTime())
	maxAmount := launch.TotalSupply.MulRaw(int64(currentLimit)).QuoRaw(10000) // Convert from basis points

	// Set limits for creator
	creatorLimits := types.WalletLimits{
		TokenAddress:   launch.TokenSymbol,
		WalletAddress:  launch.Creator,
		MaxAmount:      maxAmount,
		CurrentAmount:  sdk.ZeroInt(),
		LastTxTime:     ctx.BlockTime(),
		LastTxBlock:    ctx.BlockHeight(),
		ViolationCount: 0,
		IsRestricted:   false,
	}

	k.setWalletLimits(ctx, creatorLimits)

	// Initialize trading metrics
	metrics := types.TradingMetrics{
		TokenAddress:   launch.TokenSymbol,
		TotalVolume:    sdk.ZeroInt(),
		DailyVolume:    sdk.ZeroInt(),
		TotalTrades:    0,
		DailyTrades:    0,
		UniqueTraders:  0,
		CurrentPrice:   sdk.ZeroDec(),
		PriceChange24h: sdk.ZeroDec(),
		MarketCap:      sdk.ZeroInt(),
		Liquidity:      sdk.ZeroInt(),
		LastUpdated:    ctx.BlockTime(),
	}

	k.setTradingMetrics(ctx, metrics)
}

// createLiquidityPool creates an AMM pool on MoneyOrder DEX for the newly deployed token
func (k Keeper) createLiquidityPool(ctx sdk.Context, launch *types.TokenLaunch) error {
	// Skip if MoneyOrderKeeper is not available
	if k.moneyOrderKeeper == nil {
		return nil
	}

	// Calculate liquidity amounts
	// Use 80% of raised NAMO for liquidity (as per anti-pump config)
	liquidityPercent := launch.AntiPumpConfig.MinLiquidityPercent
	namoLiquidity := launch.RaisedAmount.MulRaw(int64(liquidityPercent)).QuoRaw(100)
	
	// Calculate proportional token amount for initial price
	// Initial price = RaisedAmount / TotalSupply
	tokenLiquidity := launch.TotalSupply.MulRaw(int64(liquidityPercent)).QuoRaw(100)

	// Get creator address
	creatorAddr, err := sdk.AccAddressFromBech32(launch.Creator)
	if err != nil {
		return err
	}

	// Create pool assets
	poolAssets := []moneyordertypes.PoolAsset{
		{
			Denom:  types.DefaultDenom, // NAMO
			Amount: namoLiquidity,
			Weight: sdk.NewInt(50), // 50% weight
		},
		{
			Denom:  launch.TokenSymbol,
			Amount: tokenLiquidity,
			Weight: sdk.NewInt(50), // 50% weight
		},
	}

	// Default swap fee (0.3%)
	swapFee := sdk.MustNewDecFromStr("0.003")

	// Create AMM pool
	poolId, err := k.moneyOrderKeeper.CreateAMMPool(ctx, creatorAddr, poolAssets, swapFee)
	if err != nil {
		return err
	}

	// Store pool ID with launch for reference
	launch.LiquidityPoolId = poolId
	k.SetTokenLaunch(ctx, launch)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"liquidity_pool_created",
			sdk.NewAttribute("launch_id", launch.LaunchID),
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("token", launch.TokenSymbol),
			sdk.NewAttribute("namo_liquidity", namoLiquidity.String()),
			sdk.NewAttribute("token_liquidity", tokenLiquidity.String()),
		),
	)

	return nil
}

// distributeLaunchFees distributes launch fees according to DeshChain revenue model
func (k Keeper) distributeLaunchFees(ctx sdk.Context, launch *types.TokenLaunch) error {
	// Platform fees are already distributed by revenue keeper when collected
	// Just handle the charity allocation which is separate from platform fees
	if launch.CharityAllocation.IsPositive() {
		charityCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, launch.CharityAllocation))
		// This would integrate with treasury module for local NGO allocation
		if err := k.treasuryKeeper.AllocateToLocalNGO(ctx, launch.CreatorPincode, charityCoins); err != nil {
			k.Logger(ctx).Error("Failed to allocate charity", "pincode", launch.CreatorPincode, "amount", launch.CharityAllocation)
		}
	}

	return nil
}

// applyFestivalBonus applies festival bonus if currently active
func (k Keeper) applyFestivalBonus(ctx sdk.Context, launch *types.TokenLaunch) {
	if k.culturalKeeper.IsActiveFestival(ctx) {
		festivalName := k.culturalKeeper.GetCurrentFestival(ctx)
		bonusRate := sdk.MustNewDecFromStr(types.FestivalBonusRate) // 10%
		
		bonusAmount := launch.TargetAmount.ToDec().Mul(bonusRate).TruncateInt()
		launch.TargetAmount = launch.TargetAmount.Add(bonusAmount)
		launch.FestivalBonus = true
		
		// Store festival bonus record
		festivalBonus := types.FestivalBonus{
			LaunchID:      launch.LaunchID,
			FestivalName:  festivalName,
			BonusRate:     bonusRate,
			BonusAmount:   bonusAmount,
			AppliedAt:     ctx.BlockTime(),
			CulturalQuote: launch.CulturalQuote,
		}
		
		k.setFestivalBonus(ctx, festivalBonus)
		
		k.Logger(ctx).Info("Applied festival bonus", 
			"launch_id", launch.LaunchID, 
			"festival", festivalName, 
			"bonus_amount", bonusAmount,
		)
	}
}

// Helper functions for data management

func (k Keeper) addCreatorLaunch(ctx sdk.Context, creator, launchID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCreatorLaunches)
	key := []byte(creator)
	
	var launches []string
	if bz := store.Get(key); bz != nil {
		k.cdc.MustUnmarshal(bz, &launches)
	}
	
	launches = append(launches, launchID)
	bz := k.cdc.MustMarshal(&launches)
	store.Set(key, bz)
}

func (k Keeper) addPincodeLaunch(ctx sdk.Context, pincode, launchID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPincodeLaunches)
	key := []byte(pincode)
	
	var launches []string
	if bz := store.Get(key); bz != nil {
		k.cdc.MustUnmarshal(bz, &launches)
	}
	
	launches = append(launches, launchID)
	bz := k.cdc.MustMarshal(&launches)
	store.Set(key, bz)
}

func (k Keeper) isWhitelisted(participant string, whitelist []string) bool {
	for _, addr := range whitelist {
		if addr == participant {
			return true
		}
	}
	return false
}

func (k Keeper) hasParticipated(ctx sdk.Context, participant, launchID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenLaunch)
	// This would check participation records - simplified for now
	return false
}

func (k Keeper) setLaunchParticipation(ctx sdk.Context, participation types.LaunchParticipation) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLaunchParticipationKey(participation.LaunchID, participation.Participant)
	bz := k.cdc.MustMarshal(&participation)
	store.Set(key, bz)
}

func (k Keeper) getLaunchParticipation(ctx sdk.Context, participant, launchID string) (types.LaunchParticipation, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLaunchParticipationKey(launchID, participant)
	bz := store.Get(key)
	if bz == nil {
		return types.LaunchParticipation{}, false
	}
	
	var participation types.LaunchParticipation
	k.cdc.MustUnmarshal(bz, &participation)
	return participation, true
}

func (k Keeper) setLiquidityLock(ctx sdk.Context, lock types.LiquidityLock) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&lock)
	store.Set(types.GetLiquidityLockKey(lock.TokenAddress), bz)
}

func (k Keeper) setCreatorReward(ctx sdk.Context, reward types.CreatorReward) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&reward)
	store.Set(types.GetCreatorRewardsKey(reward.Creator, reward.TokenAddress), bz)
}

func (k Keeper) getCreatorReward(ctx sdk.Context, creator, tokenAddress string) (types.CreatorReward, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCreatorRewardsKey(creator, tokenAddress))
	if bz == nil {
		return types.CreatorReward{}, false
	}
	
	var reward types.CreatorReward
	k.cdc.MustUnmarshal(bz, &reward)
	return reward, true
}

func (k Keeper) setWalletLimits(ctx sdk.Context, limits types.WalletLimits) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&limits)
	store.Set(types.GetWalletLimitsKey(limits.TokenAddress, limits.WalletAddress), bz)
}

func (k Keeper) setTradingMetrics(ctx sdk.Context, metrics types.TradingMetrics) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTradingMetrics)
	key := []byte(metrics.TokenAddress)
	bz := k.cdc.MustMarshal(&metrics)
	store.Set(key, bz)
}

func (k Keeper) setFestivalBonus(ctx sdk.Context, bonus types.FestivalBonus) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFestivalBonuses)
	key := []byte(bonus.LaunchID)
	bz := k.cdc.MustMarshal(&bonus)
	store.Set(key, bz)
}