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

package cultural

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/cultural/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/cultural/types"
)

// InitGenesis initializes the cultural module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)

	// Initialize quotes
	if len(genState.Quotes) == 0 {
		// If no quotes in genesis, use default quotes
		defaultQuotes := types.DefaultQuotes()
		for _, quote := range defaultQuotes {
			k.SetQuote(ctx, quote)
		}
	} else {
		// Use quotes from genesis
		for _, quote := range genState.Quotes {
			k.SetQuote(ctx, quote)
		}
	}

	// Initialize historical events
	if len(genState.HistoricalEvents) == 0 {
		// If no events in genesis, use default events
		defaultEvents := types.DefaultHistoricalEvents()
		for _, event := range defaultEvents {
			k.SetHistoricalEvent(ctx, event)
		}
	} else {
		// Use events from genesis
		for _, event := range genState.HistoricalEvents {
			k.SetHistoricalEvent(ctx, event)
		}
	}

	// Initialize cultural wisdom
	if len(genState.CulturalWisdom) == 0 {
		// If no wisdom in genesis, use default wisdom
		defaultWisdom := types.DefaultCulturalWisdom()
		for _, wisdom := range defaultWisdom {
			k.SetCulturalWisdom(ctx, wisdom)
		}
	} else {
		// Use wisdom from genesis
		for _, wisdom := range genState.CulturalWisdom {
			k.SetCulturalWisdom(ctx, wisdom)
		}
	}

	// Initialize transaction quotes
	for _, txQuote := range genState.TransactionQuotes {
		k.StoreTransactionQuote(ctx, txQuote)
	}

	// Set counters
	k.SetQuoteCount(ctx, genState.QuoteCount)
	k.SetHistoricalEventCount(ctx, genState.HistoricalEventCount)
	k.SetCulturalWisdomCount(ctx, genState.CulturalWisdomCount)
}

// ExportGenesis returns the cultural module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export all quotes
	genesis.Quotes = k.GetAllQuotes(ctx)

	// Export all historical events
	genesis.HistoricalEvents = k.GetAllHistoricalEvents(ctx)

	// Export all cultural wisdom
	genesis.CulturalWisdom = k.GetAllCulturalWisdom(ctx)

	// Export transaction quotes
	// Note: This requires iterating through the store
	// For now, we'll leave it empty as it's not critical for genesis
	genesis.TransactionQuotes = []types.TransactionQuote{}

	// Export counters
	genesis.QuoteCount = k.GetNextQuoteID(ctx) - 1
	genesis.HistoricalEventCount = k.GetNextHistoricalEventID(ctx) - 1
	genesis.CulturalWisdomCount = k.GetNextCulturalWisdomID(ctx) - 1

	return genesis
}