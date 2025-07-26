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

package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/cultural/types"
)

func TestDefaultQuotes(t *testing.T) {
	quotes := types.DefaultQuotes()
	require.NotEmpty(t, quotes)
	require.GreaterOrEqual(t, len(quotes), 15)

	// Test first quote
	firstQuote := quotes[0]
	require.Equal(t, uint64(1), firstQuote.Id)
	require.Equal(t, types.AuthorGandhi, firstQuote.Author)
	require.NotEmpty(t, firstQuote.Text)
	require.NotEmpty(t, firstQuote.Translation)
}

func TestDefaultHistoricalEvents(t *testing.T) {
	events := types.DefaultHistoricalEvents()
	require.NotEmpty(t, events)
	require.GreaterOrEqual(t, len(events), 8)

	// Test Independence Day event
	independenceEvent := events[0]
	require.Equal(t, uint64(1), independenceEvent.Id)
	require.Equal(t, "Independence Day", independenceEvent.Title)
	require.Equal(t, int32(1947), independenceEvent.Year)
}

func TestDefaultCulturalWisdom(t *testing.T) {
	wisdom := types.DefaultCulturalWisdom()
	require.NotEmpty(t, wisdom)
	require.GreaterOrEqual(t, len(wisdom), 8)

	// Test first wisdom
	firstWisdom := wisdom[0]
	require.Equal(t, uint64(1), firstWisdom.Id)
	require.NotEmpty(t, firstWisdom.Text)
	require.NotEmpty(t, firstWisdom.Translation)
	require.Equal(t, types.TraditionSanskrit, firstWisdom.Tradition)
}

func TestGetQuoteByAmountRange(t *testing.T) {
	testCases := []struct {
		name     string
		amount   sdk.Int
		minCount int
	}{
		{
			name:     "small amount",
			amount:   sdk.NewInt(500),
			minCount: 1,
		},
		{
			name:     "medium amount",
			amount:   sdk.NewInt(5000),
			minCount: 1,
		},
		{
			name:     "large amount",
			amount:   sdk.NewInt(50000),
			minCount: 1,
		},
		{
			name:     "very large amount",
			amount:   sdk.NewInt(500000),
			minCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			quotes := types.GetQuoteByAmountRange(tc.amount)
			require.NotEmpty(t, quotes)
			require.GreaterOrEqual(t, len(quotes), tc.minCount)
		})
	}
}

func TestValidateQuote(t *testing.T) {
	testCases := []struct {
		name      string
		quote     types.Quote
		shouldErr bool
	}{
		{
			name: "valid quote",
			quote: types.Quote{
				Id:                   1,
				Text:                 "Be the change you want to see",
				Author:               "Gandhi",
				DifficultyLevel:      5,
				CulturalSignificance: 8,
			},
			shouldErr: false,
		},
		{
			name: "empty text",
			quote: types.Quote{
				Id:                   1,
				Text:                 "",
				Author:               "Gandhi",
				DifficultyLevel:      5,
				CulturalSignificance: 8,
			},
			shouldErr: true,
		},
		{
			name: "empty author",
			quote: types.Quote{
				Id:                   1,
				Text:                 "Be the change",
				Author:               "",
				DifficultyLevel:      5,
				CulturalSignificance: 8,
			},
			shouldErr: true,
		},
		{
			name: "invalid difficulty",
			quote: types.Quote{
				Id:                   1,
				Text:                 "Be the change",
				Author:               "Gandhi",
				DifficultyLevel:      11,
				CulturalSignificance: 8,
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateQuote(tc.quote)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFormatQuoteForDisplay(t *testing.T) {
	quote := types.Quote{
		Text:   "Be the change you want to see in the world.",
		Author: "Mahatma Gandhi",
		Translation: map[string]string{
			"hindi": "वह परिवर्तन बनो जो तुम दुनिया में देखना चाहते हो।",
		},
	}

	// Test English format
	englishDisplay := types.FormatQuoteForDisplay(quote, "english")
	require.Contains(t, englishDisplay, quote.Text)
	require.Contains(t, englishDisplay, quote.Author)

	// Test Hindi format
	hindiDisplay := types.FormatQuoteForDisplay(quote, "hindi")
	require.Contains(t, hindiDisplay, quote.Translation["hindi"])
	require.Contains(t, hindiDisplay, quote.Author)
}