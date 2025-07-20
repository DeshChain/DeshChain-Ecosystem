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

	"github.com/deshchain/deshchain/x/cultural/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace

	// Cached data
	cachedQuotes          []types.Quote
	cachedHistoricalEvents []types.HistoricalEvent
	cachedCulturalWisdom  []types.CulturalWisdom
}

// NewKeeper creates new instances of the cultural Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
) Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	keeper := Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
	}

	// Initialize with default data
	keeper.InitializeDefaultData()

	return keeper
}

// InitializeDefaultData loads default cultural data into memory
func (k *Keeper) InitializeDefaultData() {
	k.cachedQuotes = types.DefaultQuotes()
	k.cachedHistoricalEvents = types.DefaultHistoricalEvents()
	k.cachedCulturalWisdom = types.DefaultCulturalWisdom()
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetQuote retrieves a quote by ID
func (k Keeper) GetQuote(ctx sdk.Context, id uint64) (types.Quote, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetQuoteKey(id)
	
	// First check cached quotes
	for _, quote := range k.cachedQuotes {
		if quote.Id == id {
			return quote, true
		}
	}
	
	// Then check store
	bz := store.Get(key)
	if bz == nil {
		return types.Quote{}, false
	}
	
	var quote types.Quote
	k.cdc.MustUnmarshal(bz, &quote)
	return quote, true
}

// SetQuote stores a quote
func (k Keeper) SetQuote(ctx sdk.Context, quote types.Quote) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetQuoteKey(quote.Id)
	bz := k.cdc.MustMarshal(&quote)
	store.Set(key, bz)
}

// GetAllQuotes returns all quotes (cached + stored)
func (k Keeper) GetAllQuotes(ctx sdk.Context) []types.Quote {
	quotes := make([]types.Quote, 0, len(k.cachedQuotes))
	quotes = append(quotes, k.cachedQuotes...)
	
	// Add any additional quotes from store
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.QuotePrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var quote types.Quote
		k.cdc.MustUnmarshal(iterator.Value(), &quote)
		
		// Check if quote is already in cached list
		found := false
		for _, cachedQuote := range k.cachedQuotes {
			if cachedQuote.Id == quote.Id {
				found = true
				break
			}
		}
		
		if !found {
			quotes = append(quotes, quote)
		}
	}
	
	return quotes
}

// GetHistoricalEvent retrieves a historical event by ID
func (k Keeper) GetHistoricalEvent(ctx sdk.Context, id uint64) (types.HistoricalEvent, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalEventKey(id)
	
	// First check cached events
	for _, event := range k.cachedHistoricalEvents {
		if event.Id == id {
			return event, true
		}
	}
	
	// Then check store
	bz := store.Get(key)
	if bz == nil {
		return types.HistoricalEvent{}, false
	}
	
	var event types.HistoricalEvent
	k.cdc.MustUnmarshal(bz, &event)
	return event, true
}

// SetHistoricalEvent stores a historical event
func (k Keeper) SetHistoricalEvent(ctx sdk.Context, event types.HistoricalEvent) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalEventKey(event.Id)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// GetAllHistoricalEvents returns all historical events
func (k Keeper) GetAllHistoricalEvents(ctx sdk.Context) []types.HistoricalEvent {
	events := make([]types.HistoricalEvent, 0, len(k.cachedHistoricalEvents))
	events = append(events, k.cachedHistoricalEvents...)
	
	// Add any additional events from store
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.HistoricalEventPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var event types.HistoricalEvent
		k.cdc.MustUnmarshal(iterator.Value(), &event)
		
		// Check if event is already in cached list
		found := false
		for _, cachedEvent := range k.cachedHistoricalEvents {
			if cachedEvent.Id == event.Id {
				found = true
				break
			}
		}
		
		if !found {
			events = append(events, event)
		}
	}
	
	return events
}

// GetCulturalWisdom retrieves cultural wisdom by ID
func (k Keeper) GetCulturalWisdom(ctx sdk.Context, id uint64) (types.CulturalWisdom, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCulturalWisdomKey(id)
	
	// First check cached wisdom
	for _, wisdom := range k.cachedCulturalWisdom {
		if wisdom.Id == id {
			return wisdom, true
		}
	}
	
	// Then check store
	bz := store.Get(key)
	if bz == nil {
		return types.CulturalWisdom{}, false
	}
	
	var wisdom types.CulturalWisdom
	k.cdc.MustUnmarshal(bz, &wisdom)
	return wisdom, true
}

// SetCulturalWisdom stores cultural wisdom
func (k Keeper) SetCulturalWisdom(ctx sdk.Context, wisdom types.CulturalWisdom) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCulturalWisdomKey(wisdom.Id)
	bz := k.cdc.MustMarshal(&wisdom)
	store.Set(key, bz)
}

// GetAllCulturalWisdom returns all cultural wisdom
func (k Keeper) GetAllCulturalWisdom(ctx sdk.Context) []types.CulturalWisdom {
	wisdom := make([]types.CulturalWisdom, 0, len(k.cachedCulturalWisdom))
	wisdom = append(wisdom, k.cachedCulturalWisdom...)
	
	// Add any additional wisdom from store
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CulturalWisdomPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var w types.CulturalWisdom
		k.cdc.MustUnmarshal(iterator.Value(), &w)
		
		// Check if wisdom is already in cached list
		found := false
		for _, cachedWisdom := range k.cachedCulturalWisdom {
			if cachedWisdom.Id == w.Id {
				found = true
				break
			}
		}
		
		if !found {
			wisdom = append(wisdom, w)
		}
	}
	
	return wisdom
}

// StoreTransactionQuote stores a quote associated with a transaction
func (k Keeper) StoreTransactionQuote(ctx sdk.Context, txQuote types.TransactionQuote) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransactionQuoteKey(txQuote.TxHash)
	bz := k.cdc.MustMarshal(&txQuote)
	store.Set(key, bz)
	
	// Update quote usage count
	if quote, found := k.GetQuote(ctx, txQuote.QuoteId); found {
		quote.UsageCount++
		k.SetQuote(ctx, quote)
	}
}

// GetTransactionQuote retrieves the quote associated with a transaction
func (k Keeper) GetTransactionQuote(ctx sdk.Context, txHash string) (types.TransactionQuote, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTransactionQuoteKey(txHash)
	bz := store.Get(key)
	if bz == nil {
		return types.TransactionQuote{}, false
	}
	
	var txQuote types.TransactionQuote
	k.cdc.MustUnmarshal(bz, &txQuote)
	return txQuote, true
}

// GetNextQuoteID returns the next available quote ID
func (k Keeper) GetNextQuoteID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.QuoteCountKey)
	if bz == nil {
		// Start from a high number to avoid conflicts with default quotes
		return 10000
	}
	
	return sdk.BigEndianToUint64(bz) + 1
}

// SetQuoteCount sets the total quote count
func (k Keeper) SetQuoteCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.QuoteCountKey, sdk.Uint64ToBigEndian(count))
}

// GetNextHistoricalEventID returns the next available historical event ID
func (k Keeper) GetNextHistoricalEventID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HistoricalEventCountKey)
	if bz == nil {
		// Start from a high number to avoid conflicts with default events
		return 10000
	}
	
	return sdk.BigEndianToUint64(bz) + 1
}

// SetHistoricalEventCount sets the total historical event count
func (k Keeper) SetHistoricalEventCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HistoricalEventCountKey, sdk.Uint64ToBigEndian(count))
}

// GetNextCulturalWisdomID returns the next available cultural wisdom ID
func (k Keeper) GetNextCulturalWisdomID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CulturalWisdomCountKey)
	if bz == nil {
		// Start from a high number to avoid conflicts with default wisdom
		return 10000
	}
	
	return sdk.BigEndianToUint64(bz) + 1
}

// SetCulturalWisdomCount sets the total cultural wisdom count
func (k Keeper) SetCulturalWisdomCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CulturalWisdomCountKey, sdk.Uint64ToBigEndian(count))
}