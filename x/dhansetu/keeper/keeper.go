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
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/deshchain/deshchain/x/dhansetu/types"
	moneyordertypes "github.com/deshchain/deshchain/x/moneyorder/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for DhanSetu
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	// Cross-module keepers for integration
	accountKeeper     types.AccountKeeper
	bankKeeper        types.BankKeeper
	moneyOrderKeeper  types.MoneyOrderKeeper
	culturalKeeper    types.CulturalKeeper
	namoKeeper        types.NAMOKeeper
}

// NewKeeper creates a new DhanSetu keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	moneyOrderKeeper types.MoneyOrderKeeper,
	culturalKeeper types.CulturalKeeper,
	namoKeeper types.NAMOKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		memKey:            memKey,
		paramstore:        ps,
		accountKeeper:     accountKeeper,
		bankKeeper:        bankKeeper,
		moneyOrderKeeper:  moneyOrderKeeper,
		culturalKeeper:    culturalKeeper,
		namoKeeper:        namoKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// DhanPata Virtual Address System

// RegisterDhanPataAddress registers a new DhanPata virtual address
func (k Keeper) RegisterDhanPataAddress(ctx sdk.Context, name, owner, addressType string, metadata *types.DhanPataMetadata) error {
	// Validate DhanPata name format
	if err := types.ValidateDhanPataName(name); err != nil {
		return err
	}

	// Check if name already exists
	if k.HasDhanPataAddress(ctx, name) {
		return types.ErrDhanPataAlreadyExists
	}

	// Get owner's blockchain address
	ownerAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return err
	}

	// Create DhanPata address entry
	dhanpataAddr := types.DhanPataAddress{
		Name:           name,
		Owner:          owner,
		BlockchainAddr: ownerAddr.String(),
		AddressType:    addressType,
		Metadata:       metadata,
		CreatedAt:      ctx.BlockTime(),
		UpdatedAt:      ctx.BlockTime(),
		IsActive:       true,
	}

	// Store the mapping
	k.SetDhanPataAddress(ctx, dhanpataAddr)

	// Create reverse mapping
	k.SetAddressToDhanPata(ctx, owner, name)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDhanPataRegistered,
			sdk.NewAttribute(types.AttributeKeyDhanPataName, name),
			sdk.NewAttribute("owner", owner),
			sdk.NewAttribute("address_type", addressType),
		),
	)

	return nil
}

// SetDhanPataAddress stores a DhanPata address
func (k Keeper) SetDhanPataAddress(ctx sdk.Context, address types.DhanPataAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&address)
	store.Set(types.GetDhanPataAddressKey(address.Name), bz)
}

// GetDhanPataAddress retrieves a DhanPata address by name
func (k Keeper) GetDhanPataAddress(ctx sdk.Context, name string) (types.DhanPataAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDhanPataAddressKey(name))
	if bz == nil {
		return types.DhanPataAddress{}, false
	}

	var address types.DhanPataAddress
	k.cdc.MustUnmarshal(bz, &address)
	return address, true
}

// HasDhanPataAddress checks if a DhanPata name exists
func (k Keeper) HasDhanPataAddress(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDhanPataAddressKey(name))
}

// SetAddressToDhanPata creates reverse mapping from blockchain address to DhanPata name
func (k Keeper) SetAddressToDhanPata(ctx sdk.Context, address, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAddressToDhanPataKey(address), []byte(name))
}

// GetDhanPataByAddress retrieves DhanPata name by blockchain address
func (k Keeper) GetDhanPataByAddress(ctx sdk.Context, address string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAddressToDhanPataKey(address))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// Mitra Exchange Integration

// RegisterEnhancedMitra creates an enhanced mitra profile for DhanSetu
func (k Keeper) RegisterEnhancedMitra(ctx sdk.Context, profile types.EnhancedMitraProfile) error {
	// Validate mitra type
	if profile.MitraType != types.MitraTypeIndividual &&
		profile.MitraType != types.MitraTypeBusiness &&
		profile.MitraType != types.MitraTypeGlobal {
		return types.ErrInvalidMitraType
	}

	// Calculate limits based on type and trust score
	daily, monthly := types.CalculateMitraLimits(profile.MitraType, profile.TrustScore)
	profile.DailyLimit = daily
	profile.MonthlyLimit = monthly

	// Set creation time
	profile.CreatedAt = ctx.BlockTime()
	profile.LastActiveAt = ctx.BlockTime()
	profile.IsActive = true

	// Store the profile
	k.SetEnhancedMitraProfile(ctx, profile)

	return nil
}

// SetEnhancedMitraProfile stores an enhanced mitra profile
func (k Keeper) SetEnhancedMitraProfile(ctx sdk.Context, profile types.EnhancedMitraProfile) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&profile)
	store.Set(types.GetMitraProfileKey(profile.MitraId), bz)
}

// GetEnhancedMitraProfile retrieves an enhanced mitra profile
func (k Keeper) GetEnhancedMitraProfile(ctx sdk.Context, mitraId string) (types.EnhancedMitraProfile, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMitraProfileKey(mitraId))
	if bz == nil {
		return types.EnhancedMitraProfile{}, false
	}

	var profile types.EnhancedMitraProfile
	k.cdc.MustUnmarshal(bz, &profile)
	return profile, true
}

// Kshetra Coins (Pincode-based memecoins)

// CreateKshetraCoin creates a new pincode-based community memecoin
func (k Keeper) CreateKshetraCoin(ctx sdk.Context, coin types.KshetraCoin) error {
	// Validate pincode
	if err := types.ValidatePincode(coin.Pincode); err != nil {
		return err
	}

	// Check if coin already exists for this pincode
	if k.HasKshetraCoin(ctx, coin.Pincode) {
		return types.ErrKshetraCoinExists
	}

	// Set creation time
	coin.CreatedAt = ctx.BlockTime()
	coin.IsActive = true

	// Store the coin
	k.SetKshetraCoin(ctx, coin)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKshetraCoinCreated,
			sdk.NewAttribute(types.AttributeKeyPincode, coin.Pincode),
			sdk.NewAttribute("coin_name", coin.CoinName),
			sdk.NewAttribute("coin_symbol", coin.CoinSymbol),
			sdk.NewAttribute("creator", coin.Creator),
		),
	)

	return nil
}

// SetKshetraCoin stores a Kshetra coin
func (k Keeper) SetKshetraCoin(ctx sdk.Context, coin types.KshetraCoin) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&coin)
	store.Set(types.GetKshetraCoinKey(coin.Pincode), bz)
}

// GetKshetraCoin retrieves a Kshetra coin by pincode
func (k Keeper) GetKshetraCoin(ctx sdk.Context, pincode string) (types.KshetraCoin, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetKshetraCoinKey(pincode))
	if bz == nil {
		return types.KshetraCoin{}, false
	}

	var coin types.KshetraCoin
	k.cdc.MustUnmarshal(bz, &coin)
	return coin, true
}

// HasKshetraCoin checks if a Kshetra coin exists for a pincode
func (k Keeper) HasKshetraCoin(ctx sdk.Context, pincode string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetKshetraCoinKey(pincode))
}

// Cross-Module Integration

// CreateOrderBridge creates a bridge between money order and DhanSetu
func (k Keeper) CreateOrderBridge(ctx sdk.Context, orderId, dhanpataName string) error {
	// Verify DhanPata exists
	if !k.HasDhanPataAddress(ctx, dhanpataName) {
		return types.ErrDhanPataNotFound
	}

	// Create bridge mapping
	bridge := types.CrossModuleBridge{
		BridgeId:     fmt.Sprintf("order-%s", orderId),
		SourceModule: moneyordertypes.ModuleName,
		TargetModule: types.ModuleName,
		SourceEntity: orderId,
		TargetEntity: dhanpataName,
		BridgeType:   "order_mapping",
		Metadata:     make(map[string]interface{}),
		CreatedAt:    ctx.BlockTime(),
		IsActive:     true,
	}

	// Store the bridge
	k.SetCrossModuleBridge(ctx, bridge)

	return nil
}

// SetCrossModuleBridge stores a cross-module bridge
func (k Keeper) SetCrossModuleBridge(ctx sdk.Context, bridge types.CrossModuleBridge) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bridge)
	store.Set(types.GetOrderBridgeKey(bridge.BridgeId), bz)
}

// ProcessMoneyOrderWithDhanPata processes a money order with DhanPata integration
func (k Keeper) ProcessMoneyOrderWithDhanPata(ctx sdk.Context, senderAddr, receiverDhanPata string, amount sdk.Coin, note string) error {
	// Resolve DhanPata to blockchain address
	dhanpataAddr, found := k.GetDhanPataAddress(ctx, receiverDhanPata)
	if !found {
		return types.ErrDhanPataNotFound
	}

	// Create money order through money order module
	// This would call the money order keeper's CreateMoneyOrder function
	// with the resolved blockchain address

	// Create order bridge for tracking
	orderId := fmt.Sprintf("mo-%d", ctx.BlockHeight())
	err := k.CreateOrderBridge(ctx, orderId, receiverDhanPata)
	if err != nil {
		return err
	}

	// Record trade history
	trade := types.TradeHistoryEntry{
		TradeId:       orderId,
		UserDhanPata:  receiverDhanPata,
		TradeType:     "money_order",
		SourceProduct: "moneyorder",
		Amount:        amount,
		Fee:           sdk.NewCoin(amount.Denom, amount.Amount.MulRaw(5).QuoRaw(1000)), // 0.5% fee
		Counterparty:  senderAddr,
		Status:        "completed",
		Metadata:      map[string]interface{}{"note": note},
		Timestamp:     ctx.BlockTime(),
	}

	k.RecordTradeHistory(ctx, trade)

	// Emit cross-module event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCrossModuleTransfer,
			sdk.NewAttribute(types.AttributeKeySourceModule, moneyordertypes.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTargetModule, types.ModuleName),
			sdk.NewAttribute("sender", senderAddr),
			sdk.NewAttribute("receiver_dhanpata", receiverDhanPata),
			sdk.NewAttribute("receiver_address", dhanpataAddr.BlockchainAddr),
			sdk.NewAttribute("amount", amount.String()),
		),
	)

	return nil
}

// RecordTradeHistory records a trade in the unified history
func (k Keeper) RecordTradeHistory(ctx sdk.Context, trade types.TradeHistoryEntry) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&trade)
	key := append(types.KeyPrefixTradeHistory, []byte(trade.TradeId)...)
	store.Set(key, bz)
}

// GetTradeHistory retrieves trade history for a DhanPata user
func (k Keeper) GetTradeHistory(ctx sdk.Context, dhanpataName string) []types.TradeHistoryEntry {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixTradeHistory)
	defer iterator.Close()

	var trades []types.TradeHistoryEntry
	for ; iterator.Valid(); iterator.Next() {
		var trade types.TradeHistoryEntry
		k.cdc.MustUnmarshal(iterator.Value(), &trade)
		
		if trade.UserDhanPata == dhanpataName {
			trades = append(trades, trade)
		}
	}

	return trades
}

// GetAllDhanPataAddresses returns all registered DhanPata addresses
func (k Keeper) GetAllDhanPataAddresses(ctx sdk.Context) []types.DhanPataAddress {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixDhanPataAddress)
	defer iterator.Close()

	var addresses []types.DhanPataAddress
	for ; iterator.Valid(); iterator.Next() {
		var address types.DhanPataAddress
		k.cdc.MustUnmarshal(iterator.Value(), &address)
		addresses = append(addresses, address)
	}

	return addresses
}

// GetAllKshetraCoins returns all Kshetra coins
func (k Keeper) GetAllKshetraCoins(ctx sdk.Context) []types.KshetraCoin {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixKshetraCoin)
	defer iterator.Close()

	var coins []types.KshetraCoin
	for ; iterator.Valid(); iterator.Next() {
		var coin types.KshetraCoin
		k.cdc.MustUnmarshal(iterator.Value(), &coin)
		coins = append(coins, coin)
	}

	return coins
}