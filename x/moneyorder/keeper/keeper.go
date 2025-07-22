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
	
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// GramPensionHooks defines the interface expected by Gram Pension module
type GramPensionHooks interface {
	AfterSurakshaContribution(ctx sdk.Context, pensionAccountId string, contributor sdk.AccAddress, contribution sdk.Coin, villagePostalCode string) error
	AfterSurakshaMaturity(ctx sdk.Context, pensionAccountId string, beneficiary sdk.AccAddress, maturityAmount sdk.Coin) error
	MonthlyRevenueDistribution(ctx sdk.Context) error
}

// Keeper of the money order store
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	paramstore    paramtypes.Subspace
	
	// Keepers needed for cross-module interactions
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistributionKeeper
	authKeeper    types.AuthKeeper
	
	// Module account names
	feeCollectorName string
	
	// Hooks
	hooks types.MoneyOrderHooks
	
	// P2P Matching Engine
	matchingEngine *MatchingEngine
}

// NewKeeper creates a new money order Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistributionKeeper,
	authKeeper types.AuthKeeper,
	feeCollectorName string,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	
	keeper := &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		distrKeeper:      distrKeeper,
		authKeeper:       authKeeper,
		feeCollectorName: feeCollectorName,
	}
	
	// Initialize the P2P matching engine
	keeper.matchingEngine = keeper.NewMatchingEngine()
	
	return keeper
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the money order hooks
func (k *Keeper) SetHooks(hooks types.MoneyOrderHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set money order hooks twice")
	}
	k.hooks = hooks
	return k
}

// GetParams returns the total set of money order parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.MoneyOrderParams) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of money order parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.MoneyOrderParams) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMoneyOrderHooks returns the money order hooks for gram pension integration
func (k Keeper) GetMoneyOrderHooks() GramPensionHooks {
	return Hooks{k}
}

// GetNextPoolId returns and increments the global pool ID counter
func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefixSequence
	
	bz := store.Get(key)
	if bz == nil {
		// Initialize with 1 if not set
		store.Set(key, sdk.Uint64ToBigEndian(1))
		return 1
	}
	
	currentId := sdk.BigEndianToUint64(bz)
	nextId := currentId + 1
	store.Set(key, sdk.Uint64ToBigEndian(nextId))
	
	return currentId
}

// GetNextOrderId returns and increments the global order ID counter
func (k Keeper) GetNextOrderId(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixSequence, []byte("order")...)
	
	bz := store.Get(key)
	var currentId uint64
	if bz == nil {
		currentId = 1
	} else {
		currentId = sdk.BigEndianToUint64(bz)
	}
	
	nextId := currentId + 1
	store.Set(key, sdk.Uint64ToBigEndian(nextId))
	
	// Format as ORDER-YYYYMMDD-NNNNNN
	return fmt.Sprintf("ORDER-%s-%06d", ctx.BlockTime().Format("20060102"), currentId)
}

// TransferFunds transfers funds between accounts with proper error handling
func (k Keeper) TransferFunds(
	ctx sdk.Context,
	from sdk.AccAddress,
	to sdk.AccAddress,
	amount sdk.Coins,
) error {
	if err := k.bankKeeper.SendCoins(ctx, from, to, amount); err != nil {
		return err
	}
	
	return nil
}

// CollectFees collects fees and distributes them according to parameters
func (k Keeper) CollectFees(
	ctx sdk.Context,
	fees sdk.Coin,
	from sdk.AccAddress,
) error {
	// First, collect fees from sender
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, from, types.MoneyOrderFeeCollector, sdk.NewCoins(fees),
	); err != nil {
		return err
	}
	
	// Get distribution parameters
	params := k.GetParams(ctx)
	
	// Distribute fees according to parameters
	distribution := params.DistributeFees(fees)
	
	for poolName, amount := range distribution {
		if amount.IsZero() {
			continue
		}
		
		// Transfer from fee collector to specific pool
		if err := k.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			types.MoneyOrderFeeCollector,
			poolName,
			sdk.NewCoins(amount),
		); err != nil {
			return err
		}
	}
	
	// Emit fee collection event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMoneyOrder,
			sdk.NewAttribute(types.AttributeKeyFees, fees.String()),
			sdk.NewAttribute(types.AttributeKeySender, from.String()),
		),
	)
	
	return nil
}

// ResolveUPIAddress resolves a UPI-style address to an account address
func (k Keeper) ResolveUPIAddress(ctx sdk.Context, upiAddress string) (sdk.AccAddress, error) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("upi:"), []byte(upiAddress)...)
	
	bz := store.Get(key)
	if bz == nil {
		// Try to parse as regular address
		addr, err := sdk.AccAddressFromBech32(upiAddress)
		if err != nil {
			return nil, types.ErrInvalidReceiverUPI
		}
		return addr, nil
	}
	
	return sdk.AccAddress(bz), nil
}

// RegisterUPIAddress registers a UPI-style address for an account
func (k Keeper) RegisterUPIAddress(ctx sdk.Context, upiAddress string, accAddress sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("upi:"), []byte(upiAddress)...)
	
	// Check if already registered
	if store.Has(key) {
		return fmt.Errorf("UPI address already registered")
	}
	
	store.Set(key, accAddress.Bytes())
	return nil
}

// GetModuleAccountBalance gets the balance of a module account
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context, moduleName string) sdk.Coins {
	moduleAcc := k.accountKeeper.GetModuleAccount(ctx, moduleName)
	if moduleAcc == nil {
		return sdk.NewCoins()
	}
	return k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
}

// ValidateKYC validates if an address has completed KYC
func (k Keeper) ValidateKYC(ctx sdk.Context, address sdk.AccAddress) error {
	// Placeholder for KYC validation logic
	// In production, this would integrate with the KYC system
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("kyc:"), address.Bytes()...)
	
	if !store.Has(key) {
		return types.ErrKYCNotCompleted
	}
	
	return nil
}

// CheckDailyLimit checks if a user has exceeded their daily limit
func (k Keeper) CheckDailyLimit(ctx sdk.Context, address sdk.AccAddress, amount sdk.Int) error {
	params := k.GetParams(ctx)
	store := ctx.KVStore(k.storeKey)
	
	// Key format: daily_limit:<address>:<date>
	today := ctx.BlockTime().Format("20060102")
	key := append([]byte("daily_limit:"), append(address.Bytes(), []byte(today)...)...)
	
	var currentAmount sdk.Int
	bz := store.Get(key)
	if bz != nil {
		currentAmount = sdk.NewIntFromBigInt(sdk.BigIntFromBytes(bz))
	} else {
		currentAmount = sdk.ZeroInt()
	}
	
	newAmount := currentAmount.Add(amount)
	if newAmount.GT(params.MaxDailyUserLimit) {
		return types.ErrDailyLimitExceeded
	}
	
	// Update the daily amount
	store.Set(key, newAmount.BigInt().Bytes())
	
	return nil
}

// IsFestivalPeriod checks if current time is during a festival
func (k Keeper) IsFestivalPeriod(ctx sdk.Context) bool {
	// Placeholder for festival period logic
	// In production, this would check against configured festival dates
	store := ctx.KVStore(k.storeKey)
	key := []byte("festival:active")
	
	return store.Has(key)
}

// GetCulturalQuote returns a cultural quote for receipts
func (k Keeper) GetCulturalQuote(ctx sdk.Context, language string) string {
	// Placeholder for cultural quote system
	quotes := map[string]string{
		"en": "Where there is trust, there is happiness",
		"hi": "जहाँ भरोसा है, वहाँ खुशी है",
	}
	
	if quote, exists := quotes[language]; exists {
		return quote
	}
	
	return quotes["en"]
}

// AfterPoolCreated calls the hooks after a pool is created
func (k Keeper) AfterPoolCreated(ctx sdk.Context, poolId uint64) {
	if k.hooks != nil {
		k.hooks.AfterPoolCreated(ctx, poolId)
	}
}

// AfterSwap calls the hooks after a swap is executed
func (k Keeper) AfterSwap(ctx sdk.Context, poolId uint64, tokenIn, tokenOut sdk.Coin) {
	if k.hooks != nil {
		k.hooks.AfterSwap(ctx, poolId, tokenIn, tokenOut)
	}
}

// GetAllMoneyOrderReceipts returns all money order receipts
func (k Keeper) GetAllMoneyOrderReceipts(ctx sdk.Context) []*types.MoneyOrderReceipt {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixMoneyOrder)
	defer iterator.Close()
	
	var receipts []*types.MoneyOrderReceipt
	for ; iterator.Valid(); iterator.Next() {
		var receipt types.MoneyOrderReceipt
		k.cdc.MustUnmarshal(iterator.Value(), &receipt)
		receipts = append(receipts, &receipt)
	}
	
	return receipts
}

// ResetDailyLimits resets daily transaction limits
func (k Keeper) ResetDailyLimits(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("daily_limit:"))
	defer iterator.Close()
	
	// Delete all daily limit entries
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// UpdateFestivalStatus updates the festival period status
func (k Keeper) UpdateFestivalStatus(ctx sdk.Context) {
	// Placeholder for festival status update logic
	// In production, would check against configured festival dates
	store := ctx.KVStore(k.storeKey)
	
	// Example: Check if it's Diwali period (would be more sophisticated)
	currentMonth := ctx.BlockTime().Month()
	currentDay := ctx.BlockTime().Day()
	
	// Simplified festival check
	if (currentMonth == 10 || currentMonth == 11) && currentDay >= 1 && currentDay <= 5 {
		store.Set([]byte("festival:active"), []byte{1})
		store.Set([]byte("festival:name"), []byte(types.FestivalDiwali))
	} else {
		store.Delete([]byte("festival:active"))
		store.Delete([]byte("festival:name"))
	}
}

// GetDynamicPensionRate retrieves the current dynamic pension payout rate
// This integrates with the Gram Suraksha module's performance-based system
func (k Keeper) GetDynamicPensionRate(ctx sdk.Context) sdk.Dec {
	// Default to 30% if no rate is set
	defaultRate := sdk.NewDecWithPrec(30, 2)
	
	// Get the rate from storage (set by Gram Suraksha module)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyDynamicPensionRate)
	if bz == nil {
		return defaultRate
	}
	
	rate, err := sdk.NewDecFromStr(string(bz))
	if err != nil {
		return defaultRate
	}
	
	// Ensure rate is within bounds (8% to 50%)
	minRate := sdk.NewDecWithPrec(8, 2)
	maxRate := sdk.NewDecWithPrec(50, 2)
	
	if rate.LT(minRate) {
		return minRate
	}
	if rate.GT(maxRate) {
		return maxRate
	}
	
	return rate
}

// SetDynamicPensionRate updates the dynamic pension payout rate
func (k Keeper) SetDynamicPensionRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyDynamicPensionRate, []byte(rate.String()))
}