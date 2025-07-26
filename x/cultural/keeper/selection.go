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
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/cultural/types"
)

// SelectQuoteForTransaction selects an appropriate quote for a transaction
func (k Keeper) SelectQuoteForTransaction(
	ctx sdk.Context,
	amount sdk.Int,
	sender, receiver string,
	txHash string,
) (types.Quote, error) {
	params := k.GetParams(ctx)
	
	// Get all available quotes
	allQuotes := k.GetAllQuotes(ctx)
	if len(allQuotes) == 0 {
		return types.Quote{}, types.ErrNoQuotesAvailable
	}
	
	// Apply selection algorithm based on parameters
	var selectedQuotes []types.Quote
	
	switch params.SelectionAlgorithm {
	case types.SelectionAlgorithmRandom:
		selectedQuotes = allQuotes
	case types.SelectionAlgorithmAmount:
		selectedQuotes = k.filterQuotesByAmount(allQuotes, amount, params.AmountRanges)
	case types.SelectionAlgorithmSeasonal:
		selectedQuotes = k.filterQuotesBySeason(allQuotes, ctx.BlockTime())
	case types.SelectionAlgorithmPersonalized:
		selectedQuotes = k.filterQuotesPersonalized(allQuotes, sender, amount)
	default:
		selectedQuotes = allQuotes
	}
	
	// If no quotes match criteria, fall back to all quotes
	if len(selectedQuotes) == 0 {
		selectedQuotes = allQuotes
	}
	
	// Select a random quote from filtered list
	selectedQuote := k.selectRandomQuote(selectedQuotes, txHash)
	
	// Record the transaction quote
	txQuote := types.TransactionQuote{
		TxHash:    txHash,
		QuoteId:   selectedQuote.Id,
		Timestamp: ctx.BlockTime().Unix(),
		Amount:    amount.String(),
		Sender:    sender,
		Receiver:  receiver,
	}
	k.StoreTransactionQuote(ctx, txQuote)
	
	return selectedQuote, nil
}

// filterQuotesByAmount filters quotes based on transaction amount
func (k Keeper) filterQuotesByAmount(quotes []types.Quote, amount sdk.Int, ranges []types.AmountRange) []types.Quote {
	// If no ranges defined, use default logic from cultural_data.go
	if len(ranges) == 0 {
		return types.GetQuoteByAmountRange(amount)
	}
	
	// Find matching range
	for _, r := range ranges {
		minAmount, _ := sdk.NewIntFromString(r.MinAmount)
		maxAmount, _ := sdk.NewIntFromString(r.MaxAmount)
		
		if amount.GTE(minAmount) && amount.LTE(maxAmount) {
			return k.filterQuotesByCategories(quotes, r.Categories)
		}
	}
	
	return quotes
}

// filterQuotesByCategories filters quotes by specific categories
func (k Keeper) filterQuotesByCategories(quotes []types.Quote, categories []string) []types.Quote {
	if len(categories) == 0 {
		return quotes
	}
	
	var filtered []types.Quote
	for _, quote := range quotes {
		for _, category := range categories {
			if quote.Category == category {
				filtered = append(filtered, quote)
				break
			}
		}
	}
	
	return filtered
}

// filterQuotesBySeason filters quotes based on current season/festival
func (k Keeper) filterQuotesBySeason(quotes []types.Quote, currentTime time.Time) []types.Quote {
	// This is a simplified implementation
	// In production, this would integrate with a festival calendar
	
	month := currentTime.Month()
	var seasonalCategories []string
	
	switch month {
	case time.January:
		// Republic Day month
		seasonalCategories = []string{types.CategoryPatriotism, types.CategoryLeadership}
	case time.March:
		// Holi month
		seasonalCategories = []string{types.CategoryHappiness, types.CategoryUnity}
	case time.August:
		// Independence Day month
		seasonalCategories = []string{types.CategoryPatriotism, types.CategoryFreedom}
	case time.October, time.November:
		// Diwali season
		seasonalCategories = []string{types.CategoryProsperity, types.CategoryWisdom}
	default:
		// General motivational quotes
		seasonalCategories = []string{types.CategoryMotivation, types.CategoryWisdom}
	}
	
	return k.filterQuotesByCategories(quotes, seasonalCategories)
}

// filterQuotesPersonalized filters quotes based on user preferences
func (k Keeper) filterQuotesPersonalized(quotes []types.Quote, userAddr string, amount sdk.Int) []types.Quote {
	// Get user's quote history
	userQuoteHistory := k.GetUserQuoteHistory(userAddr)
	
	// Calculate preferred categories based on history
	categoryCount := make(map[string]int)
	for _, history := range userQuoteHistory {
		if quote, found := k.GetQuote(sdk.Context{}, history.QuoteId); found {
			categoryCount[quote.Category]++
		}
	}
	
	// Find top categories
	var preferredCategories []string
	for category := range categoryCount {
		preferredCategories = append(preferredCategories, category)
	}
	
	// If user has preferences, use them
	if len(preferredCategories) > 0 {
		filtered := k.filterQuotesByCategories(quotes, preferredCategories)
		if len(filtered) > 0 {
			return filtered
		}
	}
	
	// Otherwise, use amount-based filtering
	return types.GetQuoteByAmountRange(amount)
}

// selectRandomQuote selects a random quote using deterministic randomness
func (k Keeper) selectRandomQuote(quotes []types.Quote, seed string) types.Quote {
	if len(quotes) == 0 {
		return types.Quote{}
	}
	
	if len(quotes) == 1 {
		return quotes[0]
	}
	
	// Use transaction hash as seed for deterministic randomness
	h := sha256.Sum256([]byte(seed))
	randSeed := binary.BigEndian.Uint64(h[:8])
	r := rand.New(rand.NewSource(int64(randSeed)))
	
	return quotes[r.Intn(len(quotes))]
}

// GetUserQuoteHistory retrieves quote history for a user
func (k Keeper) GetUserQuoteHistory(userAddr string) []types.TransactionQuote {
	// This is a simplified implementation
	// In production, this would query an index of user transactions
	return []types.TransactionQuote{}
}

// GetDailyQuote returns the quote of the day
func (k Keeper) GetDailyQuote(ctx sdk.Context) (types.Quote, error) {
	// Use block time to ensure consistency across nodes
	wisdom := types.GetDailyWisdom(ctx.BlockTime())
	
	// Convert wisdom to quote format for display
	quote := types.Quote{
		Id:                   uint64(wisdom.Id),
		Text:                 wisdom.Text,
		Author:               string(wisdom.Tradition),
		Language:             wisdom.Language,
		Category:             types.CategorySpirituality,
		Translation:          map[string]string{"english": wisdom.Translation},
		CulturalSignificance: wisdom.SpiritualSignificance,
	}
	
	return quote, nil
}

// GetQuotesByRegion returns quotes from a specific region
func (k Keeper) GetQuotesByRegion(ctx sdk.Context, region string) []types.Quote {
	allQuotes := k.GetAllQuotes(ctx)
	return types.GetRegionalQuote(allQuotes, region)
}

// GetQuotesByAuthor returns all quotes by a specific author
func (k Keeper) GetQuotesByAuthor(ctx sdk.Context, author string) []types.Quote {
	allQuotes := k.GetAllQuotes(ctx)
	
	var authorQuotes []types.Quote
	for _, quote := range allQuotes {
		if quote.Author == author {
			authorQuotes = append(authorQuotes, quote)
		}
	}
	
	return authorQuotes
}

// GetQuotesByCategory returns all quotes in a specific category
func (k Keeper) GetQuotesByCategory(ctx sdk.Context, category string) []types.Quote {
	allQuotes := k.GetAllQuotes(ctx)
	
	var categoryQuotes []types.Quote
	for _, quote := range allQuotes {
		if quote.Category == category {
			categoryQuotes = append(categoryQuotes, quote)
		}
	}
	
	return categoryQuotes
}

// GetPopularQuotes returns the most used quotes
func (k Keeper) GetPopularQuotes(ctx sdk.Context, limit int) []types.Quote {
	allQuotes := k.GetAllQuotes(ctx)
	
	// Sort by usage count
	// In production, this would use a more efficient data structure
	for i := 0; i < len(allQuotes); i++ {
		for j := i + 1; j < len(allQuotes); j++ {
			if allQuotes[j].UsageCount > allQuotes[i].UsageCount {
				allQuotes[i], allQuotes[j] = allQuotes[j], allQuotes[i]
			}
		}
	}
	
	if limit > len(allQuotes) {
		limit = len(allQuotes)
	}
	
	return allQuotes[:limit]
}

// GetQuoteStatistics returns statistics about quote usage
func (k Keeper) GetQuoteStatistics(ctx sdk.Context) map[string]interface{} {
	allQuotes := k.GetAllQuotes(ctx)
	return types.GetQuoteStatistics(allQuotes)
}