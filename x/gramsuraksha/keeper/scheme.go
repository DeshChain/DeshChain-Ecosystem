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
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// CreateScheme creates a new pension scheme
func (k Keeper) CreateScheme(ctx sdk.Context, scheme types.SurakshaScheme) error {
	// Validate scheme
	if err := scheme.Validate(); err != nil {
		return err
	}

	// Check if scheme ID already exists
	if k.HasScheme(ctx, scheme.SchemeID) {
		return types.ErrInvalidSchemeID
	}

	// Set timestamps
	scheme.CreatedAt = ctx.BlockTime()
	scheme.UpdatedAt = ctx.BlockTime()
	scheme.Status = types.StatusActive

	// Store scheme
	k.SetScheme(ctx, scheme)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSchemeCreated,
			sdk.NewAttribute(types.AttributeKeySchemeID, scheme.SchemeID),
			sdk.NewAttribute(types.AttributeKeySchemeName, scheme.SchemeName),
			sdk.NewAttribute(types.AttributeKeySchemeStatus, scheme.Status),
		),
	)

	return nil
}

// UpdateScheme updates an existing pension scheme
func (k Keeper) UpdateScheme(ctx sdk.Context, schemeID string, updates types.MsgUpdateScheme) error {
	scheme, found := k.GetScheme(ctx, schemeID)
	if !found {
		return types.ErrSchemeNotFound
	}

	// Update fields
	if updates.Description != "" {
		scheme.Description = updates.Description
	}
	if updates.Status != "" {
		scheme.Status = updates.Status
	}

	scheme.UpdatedAt = ctx.BlockTime()

	// Store updated scheme
	k.SetScheme(ctx, scheme)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSchemeUpdated,
			sdk.NewAttribute(types.AttributeKeySchemeID, scheme.SchemeID),
			sdk.NewAttribute(types.AttributeKeySchemeStatus, scheme.Status),
		),
	)

	return nil
}

// SetScheme stores a pension scheme
func (k Keeper) SetScheme(ctx sdk.Context, scheme types.SurakshaScheme) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&scheme)
	store.Set(types.SurakshaSchemePrefix.Bytes(append([]byte(scheme.SchemeID)), bz)
}

// GetScheme retrieves a pension scheme by ID
func (k Keeper) GetScheme(ctx sdk.Context, schemeID string) (types.SurakshaScheme, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SurakshaSchemePrefix.Bytes(append([]byte(schemeID)))
	if bz == nil {
		return types.SurakshaScheme{}, false
	}

	var scheme types.SurakshaScheme
	k.cdc.MustUnmarshal(bz, &scheme)
	return scheme, true
}

// HasScheme checks if a scheme exists
func (k Keeper) HasScheme(ctx sdk.Context, schemeID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.SurakshaSchemePrefix.Bytes(append([]byte(schemeID)))
}

// GetAllSchemes returns all pension schemes
func (k Keeper) GetAllSchemes(ctx sdk.Context) []types.SurakshaScheme {
	var schemes []types.SurakshaScheme
	k.IterateSchemes(ctx, func(scheme types.SurakshaScheme) bool {
		schemes = append(schemes, scheme)
		return false
	})
	return schemes
}

// IterateSchemes iterates over all schemes
func (k Keeper) IterateSchemes(ctx sdk.Context, cb func(types.SurakshaScheme) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SurakshaSchemePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var scheme types.SurakshaScheme
		k.cdc.MustUnmarshal(iterator.Value(), &scheme)
		if cb(scheme) {
			break
		}
	}
}

// GetActiveSchemes returns all active pension schemes
func (k Keeper) GetActiveSchemes(ctx sdk.Context) []types.SurakshaScheme {
	var schemes []types.SurakshaScheme
	k.IterateSchemes(ctx, func(scheme types.SurakshaScheme) bool {
		if scheme.Status == types.StatusActive {
			schemes = append(schemes, scheme)
		}
		return false
	})
	return schemes
}

// UpdateSchemeStatistics updates scheme statistics
func (k Keeper) UpdateSchemeStatistics(ctx sdk.Context, schemeID string) error {
	scheme, found := k.GetScheme(ctx, schemeID)
	if !found {
		return types.ErrSchemeNotFound
	}

	// Get current statistics
	stats, found := k.GetSchemeStatistics(ctx, schemeID)
	if !found {
		stats = types.PensionStatistics{
			SchemeID: schemeID,
		}
	}

	// Count participants by status
	activeCount := uint64(0)
	maturedCount := uint64(0)
	withdrawnCount := uint64(0)
	defaultedCount := uint64(0)
	totalContributed := sdk.NewCoin(scheme.MonthlyContribution.Denom, sdk.ZeroInt())

	k.IterateParticipantsByScheme(ctx, schemeID, func(participant types.SurakshaParticipant) bool {
		switch participant.Status {
		case types.StatusActive:
			activeCount++
		case types.StatusMatured:
			maturedCount++
		case types.StatusWithdrawn:
			withdrawnCount++
		case types.StatusDefaulted:
			defaultedCount++
		}
		totalContributed = totalContributed.Add(participant.TotalContributed)
		return false
	})

	// Update statistics
	stats.TotalParticipants = scheme.CurrentParticipants
	stats.ActiveParticipants = activeCount
	stats.MaturedParticipants = maturedCount
	stats.WithdrawnParticipants = withdrawnCount
	stats.DefaultedParticipants = defaultedCount
	stats.TotalContributed = totalContributed
	stats.LastUpdated = ctx.BlockTime()

	// Calculate rates
	if stats.TotalParticipants > 0 {
		stats.CompletionRate = sdk.NewDec(int64(stats.MaturedParticipants)).Quo(sdk.NewDec(int64(stats.TotalParticipants)))
		stats.DefaultRate = sdk.NewDec(int64(stats.DefaultedParticipants)).Quo(sdk.NewDec(int64(stats.TotalParticipants)))
	}

	// Store updated statistics
	k.SetSchemeStatistics(ctx, stats)

	return nil
}

// SetSchemeStatistics stores scheme statistics
func (k Keeper) SetSchemeStatistics(ctx sdk.Context, stats types.PensionStatistics) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set(types.SchemeStatisticsPrefix.Bytes(append([]byte(stats.SchemeID)), bz)
}

// GetSchemeStatistics retrieves scheme statistics
func (k Keeper) GetSchemeStatistics(ctx sdk.Context, schemeID string) (types.PensionStatistics, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SchemeStatisticsPrefix.Bytes(append([]byte(schemeID)))
	if bz == nil {
		return types.PensionStatistics{}, false
	}

	var stats types.PensionStatistics
	k.cdc.MustUnmarshal(bz, &stats)
	return stats, true
}

// GenerateSchemeID generates a unique scheme ID
func (k Keeper) GenerateSchemeID(ctx sdk.Context, schemeName string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("SCHEME-%s-%d", schemeName[:min(5, len(schemeName))], timestamp)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}