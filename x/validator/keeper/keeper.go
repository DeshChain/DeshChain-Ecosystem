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
	
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	
	"github.com/deshchain/namo/x/validator/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	
	// Expected keepers
	bankKeeper    BankKeeper
	stakingKeeper StakingKeeper
	accountKeeper AccountKeeper
	
	// Staking manager for new USD-based staking
	stakingManager *StakingManager
	
	// Referral keeper for referral system
	referralKeeper *ReferralKeeper
	
	// Oracle keeper for price feeds
	oracleKeeper OracleKeeper
}

// NewKeeper creates a new validator keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	accountKeeper AccountKeeper,
	oracleKeeper OracleKeeper,
) *Keeper {
	k := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		accountKeeper: accountKeeper,
		oracleKeeper:  oracleKeeper,
	}
	
	// Initialize staking manager
	k.stakingManager = NewStakingManager(*k)
	
	// Initialize referral keeper
	k.referralKeeper = NewReferralKeeper(*k)
	
	return k
}

// Expected keepers interfaces
type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	PowerReduction(ctx sdk.Context) sdk.Int
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

type OracleKeeper interface {
	GetNAMOPriceUSD(ctx sdk.Context) (sdk.Dec, error)
	GetExchangeRate(ctx sdk.Context, denom string) (sdk.Dec, error)
}

// GetStakingManager returns the staking manager
func (k Keeper) GetStakingManager() *StakingManager {
	return k.stakingManager
}

// GetReferralKeeper returns the referral keeper
func (k Keeper) GetReferralKeeper() *ReferralKeeper {
	return k.referralKeeper
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) sdk.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetInsurancePoolAddress returns the insurance pool address
func (k Keeper) GetInsurancePoolAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.InsurancePoolName)
}

// Storage functions for validator stakes

// SetValidatorStake stores a validator stake
func (k Keeper) SetValidatorStake(ctx sdk.Context, stake types.ValidatorStake) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stake)
	store.Set(types.GetValidatorStakeKey(stake.ValidatorAddr), bz)
}

// GetValidatorStake retrieves a validator stake
func (k Keeper) GetValidatorStake(ctx sdk.Context, validatorAddr string) (stake types.ValidatorStake, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorStakeKey(validatorAddr))
	if bz == nil {
		return stake, false
	}
	k.cdc.MustUnmarshal(bz, &stake)
	return stake, true
}

// GetAllValidatorStakes returns all validator stakes
func (k Keeper) GetAllValidatorStakes(ctx sdk.Context) []types.ValidatorStake {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorStakeKeyPrefix)
	defer iterator.Close()
	
	var stakes []types.ValidatorStake
	for ; iterator.Valid(); iterator.Next() {
		var stake types.ValidatorStake
		k.cdc.MustUnmarshal(iterator.Value(), &stake)
		stakes = append(stakes, stake)
	}
	
	return stakes
}

// GetGenesisNFT retrieves a genesis NFT by token ID
func (k Keeper) GetGenesisNFT(ctx sdk.Context, tokenID uint64) (nft types.GenesisValidatorNFT, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGenesisNFTKey(tokenID))
	if bz == nil {
		return nft, false
	}
	k.cdc.MustUnmarshal(bz, &nft)
	return nft, true
}

// SetGenesisNFT stores a genesis NFT
func (k Keeper) SetGenesisNFT(ctx sdk.Context, nft types.GenesisValidatorNFT) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&nft)
	store.Set(types.GetGenesisNFTKey(nft.TokenID), bz)
}

// GetNextNFTID returns the next available NFT ID
func (k Keeper) GetNextNFTID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextNFTIDKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextNFTID sets the next NFT ID
func (k Keeper) SetNextNFTID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextNFTIDKey, sdk.Uint64ToBigEndian(id))
}

// TransferValidatorRights transfers validator rights when NFT is traded
func (k Keeper) TransferValidatorRights(ctx sdk.Context, from, to string) error {
	// This would integrate with the staking module to transfer validator rights
	// For now, we'll emit an event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorRightsTransferred,
			sdk.NewAttribute("from", from),
			sdk.NewAttribute("to", to),
		),
	)
	return nil
}

// SendCoinsToValidator sends coins to a validator address
func (k Keeper) SendCoinsToValidator(ctx sdk.Context, validatorAddr string, amount sdk.Coins) error {
	addr, err := sdk.AccAddressFromBech32(validatorAddr)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}
	
	return k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		addr,
		amount,
	)
}

// GetAllActiveValidators returns all active validators
func (k Keeper) GetAllActiveValidators(ctx sdk.Context) []types.Validator {
	// For backward compatibility, we'll use the staking keeper to get validators
	// and convert them to our validator type
	stakedValidators := k.stakingKeeper.GetAllValidators(ctx)
	
	var validators []types.Validator
	for _, sv := range stakedValidators {
		if sv.IsBonded() {
			validators = append(validators, types.Validator{
				OperatorAddress: sv.OperatorAddress,
				Tokens:          sv.Tokens,
				Status:          sv.Status,
			})
		}
	}
	
	return validators
}

// Referral storage functions

// SetReferral stores a referral
func (k Keeper) SetReferral(ctx sdk.Context, referral types.Referral) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&referral)
	store.Set(types.GetReferralKey(referral.ReferralID), bz)
}

// GetReferral retrieves a referral by ID
func (k Keeper) GetReferral(ctx sdk.Context, referralID uint64) (types.Referral, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReferralKey(referralID))
	if bz == nil {
		return types.Referral{}, false
	}
	var referral types.Referral
	k.cdc.MustUnmarshal(bz, &referral)
	return referral, true
}

// GetAllReferrals returns all referrals
func (k Keeper) GetAllReferrals(ctx sdk.Context) []types.Referral {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReferralKeyPrefix)
	defer iterator.Close()
	
	var referrals []types.Referral
	for ; iterator.Valid(); iterator.Next() {
		var referral types.Referral
		k.cdc.MustUnmarshal(iterator.Value(), &referral)
		referrals = append(referrals, referral)
	}
	
	return referrals
}

// GetReferralsByReferrer returns all referrals by a referrer
func (k Keeper) GetReferralsByReferrer(ctx sdk.Context, referrerAddr string) []types.Referral {
	allReferrals := k.GetAllReferrals(ctx)
	var referrals []types.Referral
	
	for _, ref := range allReferrals {
		if ref.ReferrerAddr == referrerAddr {
			referrals = append(referrals, ref)
		}
	}
	
	return referrals
}

// SetReferralStats stores referral statistics
func (k Keeper) SetReferralStats(ctx sdk.Context, stats types.ReferralStats) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set(types.GetReferralStatsKey(stats.ValidatorAddr), bz)
}

// GetReferralStats retrieves referral statistics
func (k Keeper) GetReferralStats(ctx sdk.Context, validatorAddr string) (types.ReferralStats, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReferralStatsKey(validatorAddr))
	if bz == nil {
		return types.ReferralStats{}, false
	}
	var stats types.ReferralStats
	k.cdc.MustUnmarshal(bz, &stats)
	return stats, true
}

// SetValidatorToken stores a validator token
func (k Keeper) SetValidatorToken(ctx sdk.Context, token types.ValidatorToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&token)
	store.Set(types.GetValidatorTokenKey(token.TokenID), bz)
}

// GetValidatorToken retrieves a validator token
func (k Keeper) GetValidatorToken(ctx sdk.Context, tokenID uint64) (types.ValidatorToken, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorTokenKey(tokenID))
	if bz == nil {
		return types.ValidatorToken{}, false
	}
	var token types.ValidatorToken
	k.cdc.MustUnmarshal(bz, &token)
	return token, true
}

// SetCommissionPayout stores a commission payout
func (k Keeper) SetCommissionPayout(ctx sdk.Context, payout types.CommissionPayout) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&payout)
	store.Set(types.GetCommissionPayoutKey(payout.PayoutID), bz)
}

// GetCommissionPayout retrieves a commission payout
func (k Keeper) GetCommissionPayout(ctx sdk.Context, payoutID uint64) (types.CommissionPayout, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCommissionPayoutKey(payoutID))
	if bz == nil {
		return types.CommissionPayout{}, false
	}
	var payout types.CommissionPayout
	k.cdc.MustUnmarshal(bz, &payout)
	return payout, true
}

// ID management functions

// GetNextReferralID returns the next referral ID
func (k Keeper) GetNextReferralID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextReferralIDKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextReferralID sets the next referral ID
func (k Keeper) SetNextReferralID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextReferralIDKey, sdk.Uint64ToBigEndian(id))
}

// GetNextTokenID returns the next token ID
func (k Keeper) GetNextTokenID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextTokenIDKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextTokenID sets the next token ID
func (k Keeper) SetNextTokenID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextTokenIDKey, sdk.Uint64ToBigEndian(id))
}

// GetNextPayoutID returns the next payout ID
func (k Keeper) GetNextPayoutID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextPayoutIDKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextPayoutID sets the next payout ID
func (k Keeper) SetNextPayoutID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextPayoutIDKey, sdk.Uint64ToBigEndian(id))
}