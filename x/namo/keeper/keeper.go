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
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper creates new instances of the NAMO Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetTokenSupply returns the current token supply
func (k Keeper) GetTokenSupply(ctx sdk.Context) (types.TokenSupply, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenSupplyKey)
	if bz == nil {
		return types.TokenSupply{}, false
	}

	var supply types.TokenSupply
	k.cdc.MustUnmarshal(bz, &supply)
	return supply, true
}

// SetTokenSupply sets the token supply
func (k Keeper) SetTokenSupply(ctx sdk.Context, supply types.TokenSupply) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&supply)
	store.Set(types.TokenSupplyKey, bz)
}

// GetVestingSchedule retrieves a vesting schedule by beneficiary address
func (k Keeper) GetVestingSchedule(ctx sdk.Context, beneficiary string) (types.VestingSchedule, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVestingScheduleKey(beneficiary)
	bz := store.Get(key)
	if bz == nil {
		return types.VestingSchedule{}, false
	}

	var schedule types.VestingSchedule
	k.cdc.MustUnmarshal(bz, &schedule)
	return schedule, true
}

// SetVestingSchedule sets a vesting schedule
func (k Keeper) SetVestingSchedule(ctx sdk.Context, schedule types.VestingSchedule) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVestingScheduleKey(schedule.Beneficiary)
	bz := k.cdc.MustMarshal(&schedule)
	store.Set(key, bz)
}

// GetAllVestingSchedules returns all vesting schedules
func (k Keeper) GetAllVestingSchedules(ctx sdk.Context) []types.VestingSchedule {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VestingSchedulePrefix)
	defer iterator.Close()

	var schedules []types.VestingSchedule
	for ; iterator.Valid(); iterator.Next() {
		var schedule types.VestingSchedule
		k.cdc.MustUnmarshal(iterator.Value(), &schedule)
		schedules = append(schedules, schedule)
	}

	return schedules
}

// CalculateVestedTokens calculates the amount of tokens that have vested for a beneficiary
func (k Keeper) CalculateVestedTokens(ctx sdk.Context, beneficiary string) (sdk.Int, error) {
	schedule, found := k.GetVestingSchedule(ctx, beneficiary)
	if !found {
		return sdk.ZeroInt(), types.ErrVestingScheduleNotFound
	}

	currentTime := ctx.BlockTime()
	cliffTime := time.Unix(schedule.CliffTime, 0)
	endTime := time.Unix(schedule.EndTime, 0)

	// Before cliff, no tokens are vested
	if currentTime.Before(cliffTime) {
		return sdk.ZeroInt(), nil
	}

	totalAmount, ok := sdk.NewIntFromString(schedule.TotalAmount)
	if !ok {
		return sdk.ZeroInt(), types.ErrInvalidTokenAmount
	}

	// After end time, all tokens are vested
	if currentTime.After(endTime) || currentTime.Equal(endTime) {
		claimedAmount, ok := sdk.NewIntFromString(schedule.ClaimedAmount)
		if !ok {
			return sdk.ZeroInt(), types.ErrInvalidTokenAmount
		}
		return totalAmount.Sub(claimedAmount), nil
	}

	// Calculate proportional vesting
	vestingDuration := endTime.Unix() - cliffTime.Unix()
	elapsedTime := currentTime.Unix() - cliffTime.Unix()

	if vestingDuration <= 0 {
		return sdk.ZeroInt(), types.ErrInvalidVestingPeriod
	}

	// Proportional vesting calculation
	vestedAmount := totalAmount.MulRaw(elapsedTime).QuoRaw(vestingDuration)

	claimedAmount, ok := sdk.NewIntFromString(schedule.ClaimedAmount)
	if !ok {
		return sdk.ZeroInt(), types.ErrInvalidTokenAmount
	}

	claimableAmount := vestedAmount.Sub(claimedAmount)
	if claimableAmount.IsNegative() {
		return sdk.ZeroInt(), nil
	}

	return claimableAmount, nil
}

// ClaimVestedTokens allows a beneficiary to claim their vested tokens
func (k Keeper) ClaimVestedTokens(ctx sdk.Context, beneficiary string) error {
	claimableAmount, err := k.CalculateVestedTokens(ctx, beneficiary)
	if err != nil {
		return err
	}

	if claimableAmount.IsZero() {
		return types.ErrNoTokensToClaim
	}

	// Get beneficiary address
	beneficiaryAddr, err := sdk.AccAddressFromBech32(beneficiary)
	if err != nil {
		return err
	}

	// Get vesting schedule to update claimed amount
	schedule, found := k.GetVestingSchedule(ctx, beneficiary)
	if !found {
		return types.ErrVestingScheduleNotFound
	}

	// Transfer tokens from module account to beneficiary
	coins := sdk.NewCoins(sdk.NewCoin(types.TokenDenom, claimableAmount))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.VestingPoolName, beneficiaryAddr, coins)
	if err != nil {
		return err
	}

	// Update claimed amount
	currentClaimed, ok := sdk.NewIntFromString(schedule.ClaimedAmount)
	if !ok {
		return types.ErrInvalidTokenAmount
	}
	newClaimed := currentClaimed.Add(claimableAmount)
	schedule.ClaimedAmount = newClaimed.String()

	// Save updated schedule
	k.SetVestingSchedule(ctx, schedule)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimVestedTokens,
			sdk.NewAttribute(types.AttributeKeyBeneficiary, beneficiary),
			sdk.NewAttribute(types.AttributeKeyAmount, claimableAmount.String()),
		),
	)

	return nil
}

// BurnTokens burns tokens from the specified account
func (k Keeper) BurnTokens(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int) error {
	if amount.IsZero() || amount.IsNegative() {
		return types.ErrInvalidTokenAmount
	}

	// Burn tokens
	coins := sdk.NewCoins(sdk.NewCoin(types.TokenDenom, amount))
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.BurnPoolName, coins)
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.BurnPoolName, coins)
	if err != nil {
		return err
	}

	// Update token supply
	supply, found := k.GetTokenSupply(ctx)
	if found {
		currentTotal, ok := sdk.NewIntFromString(supply.TotalSupply)
		if ok {
			newTotal := currentTotal.Sub(amount)
			supply.TotalSupply = newTotal.String()
			k.SetTokenSupply(ctx, supply)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnTokens,
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)

	return nil
}

// CreateTokenDistributionEvent creates a new token distribution event
func (k Keeper) CreateTokenDistributionEvent(ctx sdk.Context, event types.TokenDistributionEvent) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenDistributionEventKey(event.Id)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)

	// Update next event ID
	k.SetNextDistributionEventID(ctx, event.Id+1)
}

// GetTokenDistributionEvent retrieves a token distribution event by ID
func (k Keeper) GetTokenDistributionEvent(ctx sdk.Context, id uint64) (types.TokenDistributionEvent, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTokenDistributionEventKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.TokenDistributionEvent{}, false
	}

	var event types.TokenDistributionEvent
	k.cdc.MustUnmarshal(bz, &event)
	return event, true
}

// GetAllTokenDistributionEvents returns all token distribution events
func (k Keeper) GetAllTokenDistributionEvents(ctx sdk.Context) []types.TokenDistributionEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TokenDistributionEventPrefix)
	defer iterator.Close()

	var events []types.TokenDistributionEvent
	for ; iterator.Valid(); iterator.Next() {
		var event types.TokenDistributionEvent
		k.cdc.MustUnmarshal(iterator.Value(), &event)
		events = append(events, event)
	}

	return events
}

// GetNextDistributionEventID returns the next available distribution event ID
func (k Keeper) GetNextDistributionEventID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextDistributionEventIDKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextDistributionEventID sets the next distribution event ID
func (k Keeper) SetNextDistributionEventID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextDistributionEventIDKey, sdk.Uint64ToBigEndian(id))
}

// GetTotalBurnedTokens returns the total amount of burned tokens
func (k Keeper) GetTotalBurnedTokens(ctx sdk.Context) sdk.Int {
	// Get burned tokens from burn pool account
	burnPoolAddr := k.accountKeeper.GetModuleAddress(types.BurnPoolName)
	if burnPoolAddr == nil {
		return sdk.ZeroInt()
	}

	burnedCoins := k.bankKeeper.GetAllBalances(ctx, burnPoolAddr)
	namoCoins := burnedCoins.AmountOf(types.TokenDenom)
	return namoCoins
}

// GetCirculatingSupply returns the circulating supply (total - vested but unclaimed)
func (k Keeper) GetCirculatingSupply(ctx sdk.Context) sdk.Int {
	supply, found := k.GetTokenSupply(ctx)
	if !found {
		return sdk.ZeroInt()
	}

	totalSupply, ok := sdk.NewIntFromString(supply.TotalSupply)
	if !ok {
		return sdk.ZeroInt()
	}

	// Subtract tokens still locked in vesting
	vestingPoolAddr := k.accountKeeper.GetModuleAddress(types.VestingPoolName)
	if vestingPoolAddr != nil {
		vestingCoins := k.bankKeeper.GetBalance(ctx, vestingPoolAddr, types.TokenDenom)
		totalSupply = totalSupply.Sub(vestingCoins.Amount)
	}

	return totalSupply
}

// ValidateTokenOperation validates a token operation
func (k Keeper) ValidateTokenOperation(ctx sdk.Context, operation string, amount sdk.Int) error {
	if amount.IsZero() || amount.IsNegative() {
		return types.ErrInvalidTokenAmount
	}

	params := k.GetParams(ctx)
	if !params.EnableTokenOperations {
		return types.ErrTokenOperationsDisabled
	}

	return nil
}