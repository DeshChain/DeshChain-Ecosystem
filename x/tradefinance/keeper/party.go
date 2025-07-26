package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// RegisterParty registers a new trade party
func (k Keeper) RegisterParty(ctx sdk.Context, party types.TradeParty) (string, error) {
	// Check if party with same address already exists
	existingPartyID := k.GetPartyIDByAddress(ctx, party.DeshAddress)
	if existingPartyID != "" {
		return "", types.ErrPartyAlreadyExists
	}

	// Generate new party ID
	partyID := k.GetNextPartyID(ctx)
	party.PartyId = fmt.Sprintf("PARTY%06d", partyID)

	// Set verification timestamp
	party.VerifiedAt = ctx.BlockTime()
	party.IsVerified = true // Auto-verify for now, can add KYC later

	// Save party
	k.SetTradeParty(ctx, party)

	// Create index for address lookup
	k.SetPartyIDByAddress(ctx, party.DeshAddress, party.PartyId)

	// Increment counter
	k.SetNextPartyID(ctx, partyID+1)

	// Update stats
	k.IncrementPartyCount(ctx, party.PartyType)

	return party.PartyId, nil
}

// GetTradeParty returns a trade party by ID
func (k Keeper) GetTradeParty(ctx sdk.Context, partyID string) (types.TradeParty, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradePartyPrefix)
	
	bz := store.Get([]byte(partyID))
	if bz == nil {
		return types.TradeParty{}, false
	}

	var party types.TradeParty
	k.cdc.MustUnmarshal(bz, &party)
	return party, true
}

// SetTradeParty saves a trade party
func (k Keeper) SetTradeParty(ctx sdk.Context, party types.TradeParty) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradePartyPrefix)
	bz := k.cdc.MustMarshal(&party)
	store.Set([]byte(party.PartyId), bz)
}

// GetAllTradeParties returns all trade parties
func (k Keeper) GetAllTradeParties(ctx sdk.Context) []types.TradeParty {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradePartyPrefix)
	
	var parties []types.TradeParty
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var party types.TradeParty
		k.cdc.MustUnmarshal(iterator.Value(), &party)
		parties = append(parties, party)
	}
	
	return parties
}

// GetPartiesByType returns parties filtered by type
func (k Keeper) GetPartiesByType(ctx sdk.Context, partyType string) []types.TradeParty {
	var filteredParties []types.TradeParty
	
	parties := k.GetAllTradeParties(ctx)
	for _, party := range parties {
		if party.PartyType == partyType {
			filteredParties = append(filteredParties, party)
		}
	}
	
	return filteredParties
}

// GetPartyIDByAddress returns party ID by blockchain address
func (k Keeper) GetPartyIDByAddress(ctx sdk.Context, address string) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PartyByAddressPrefix)
	
	bz := store.Get([]byte(address))
	if bz == nil {
		return ""
	}
	
	return string(bz)
}

// SetPartyIDByAddress sets the party ID for an address
func (k Keeper) SetPartyIDByAddress(ctx sdk.Context, address, partyID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PartyByAddressPrefix)
	store.Set([]byte(address), []byte(partyID))
}

// GetNextPartyID returns the next party ID
func (k Keeper) GetNextPartyID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextPartyIDKey)
	
	if bz == nil {
		return 1
	}
	
	return sdk.BigEndianToUint64(bz)
}

// SetNextPartyID sets the next party ID
func (k Keeper) SetNextPartyID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextPartyIDKey, sdk.Uint64ToBigEndian(id))
}

// ValidatePartyRole validates if a party has the required role
func (k Keeper) ValidatePartyRole(ctx sdk.Context, partyID, requiredType string) error {
	party, found := k.GetTradeParty(ctx, partyID)
	if !found {
		return types.ErrPartyNotFound
	}
	
	if party.PartyType != requiredType {
		return fmt.Errorf("party %s is not a %s", partyID, requiredType)
	}
	
	if !party.IsVerified {
		return fmt.Errorf("party %s is not verified", partyID)
	}
	
	return nil
}

// IncrementPartyCount increments the party count in stats
func (k Keeper) IncrementPartyCount(ctx sdk.Context, partyType string) {
	stats := k.GetTradeFinanceStats(ctx)
	
	// This is a simplified version - you might want to track by type
	// For now, we'll just increment total count
	// stats.RegisteredParties++
	
	k.SetTradeFinanceStats(ctx, stats)
}