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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/cultural/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetEnableQuotes(ctx),
		k.GetQuoteSelectionAlgorithm(ctx),
		k.GetAmountRanges(ctx),
		k.GetSeasonalQuotesEnabled(ctx),
		k.GetUserPreferencesEnabled(ctx),
		k.GetMinQuoteLength(ctx),
		k.GetMaxQuoteLength(ctx),
		k.GetDefaultLanguage(ctx),
		k.GetAvailableLanguages(ctx),
		k.GetQuoteRefreshInterval(ctx),
		k.GetMaxQuotesPerUser(ctx),
		k.GetEnableNFTCreation(ctx),
		k.GetNFTMintingFee(ctx),
		k.GetNFTRoyaltyPercentage(ctx),
		k.GetEnableStatistics(ctx),
		k.GetStatisticsRetentionDays(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetEnableQuotes returns whether quotes are enabled
func (k Keeper) GetEnableQuotes(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyEnableQuotes, &res)
	return res
}

// GetQuoteSelectionAlgorithm returns the quote selection algorithm
func (k Keeper) GetQuoteSelectionAlgorithm(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyQuoteSelectionAlgorithm, &res)
	return res
}

// GetAmountRanges returns the amount ranges for quote selection
func (k Keeper) GetAmountRanges(ctx sdk.Context) []types.AmountRange {
	var res []types.AmountRange
	k.paramstore.Get(ctx, types.KeyAmountRanges, &res)
	return res
}

// GetSeasonalQuotesEnabled returns whether seasonal quotes are enabled
func (k Keeper) GetSeasonalQuotesEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeySeasonalQuotesEnabled, &res)
	return res
}

// GetUserPreferencesEnabled returns whether user preferences are enabled
func (k Keeper) GetUserPreferencesEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyUserPreferencesEnabled, &res)
	return res
}

// GetMinQuoteLength returns the minimum quote length
func (k Keeper) GetMinQuoteLength(ctx sdk.Context) uint32 {
	var res uint32
	k.paramstore.Get(ctx, types.KeyMinQuoteLength, &res)
	return res
}

// GetMaxQuoteLength returns the maximum quote length
func (k Keeper) GetMaxQuoteLength(ctx sdk.Context) uint32 {
	var res uint32
	k.paramstore.Get(ctx, types.KeyMaxQuoteLength, &res)
	return res
}

// GetDefaultLanguage returns the default language
func (k Keeper) GetDefaultLanguage(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyDefaultLanguage, &res)
	return res
}

// GetAvailableLanguages returns available languages
func (k Keeper) GetAvailableLanguages(ctx sdk.Context) []string {
	var res []string
	k.paramstore.Get(ctx, types.KeyAvailableLanguages, &res)
	return res
}

// GetQuoteRefreshInterval returns the quote refresh interval
func (k Keeper) GetQuoteRefreshInterval(ctx sdk.Context) int64 {
	var res int64
	k.paramstore.Get(ctx, types.KeyQuoteRefreshInterval, &res)
	return res
}

// GetMaxQuotesPerUser returns the max quotes per user
func (k Keeper) GetMaxQuotesPerUser(ctx sdk.Context) uint32 {
	var res uint32
	k.paramstore.Get(ctx, types.KeyMaxQuotesPerUser, &res)
	return res
}

// GetEnableNFTCreation returns whether NFT creation is enabled
func (k Keeper) GetEnableNFTCreation(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyEnableNFTCreation, &res)
	return res
}

// GetNFTMintingFee returns the NFT minting fee
func (k Keeper) GetNFTMintingFee(ctx sdk.Context) sdk.Coin {
	var res sdk.Coin
	k.paramstore.Get(ctx, types.KeyNFTMintingFee, &res)
	return res
}

// GetNFTRoyaltyPercentage returns the NFT royalty percentage
func (k Keeper) GetNFTRoyaltyPercentage(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyNFTRoyaltyPercentage, &res)
	return res
}

// GetEnableStatistics returns whether statistics are enabled
func (k Keeper) GetEnableStatistics(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyEnableStatistics, &res)
	return res
}

// GetStatisticsRetentionDays returns the statistics retention days
func (k Keeper) GetStatisticsRetentionDays(ctx sdk.Context) uint32 {
	var res uint32
	k.paramstore.Get(ctx, types.KeyStatisticsRetentionDays, &res)
	return res
}