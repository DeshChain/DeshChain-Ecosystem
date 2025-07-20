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

package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		Quotes:               []Quote{},
		HistoricalEvents:     []HistoricalEvent{},
		CulturalWisdom:       []CulturalWisdom{},
		TransactionQuotes:    []TransactionQuote{},
		QuoteCount:           0,
		HistoricalEventCount: 0,
		CulturalWisdomCount:  0,
	}
}

// ValidateGenesis validates the cultural module's genesis state
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if err := data.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate quotes
	quoteIDs := make(map[uint64]bool)
	for _, quote := range data.Quotes {
		if err := ValidateQuote(quote); err != nil {
			return fmt.Errorf("invalid quote %d: %w", quote.Id, err)
		}
		if quoteIDs[quote.Id] {
			return fmt.Errorf("duplicate quote ID: %d", quote.Id)
		}
		quoteIDs[quote.Id] = true
	}

	// Validate historical events
	eventIDs := make(map[uint64]bool)
	for _, event := range data.HistoricalEvents {
		if err := validateHistoricalEvent(event); err != nil {
			return fmt.Errorf("invalid historical event %d: %w", event.Id, err)
		}
		if eventIDs[event.Id] {
			return fmt.Errorf("duplicate historical event ID: %d", event.Id)
		}
		eventIDs[event.Id] = true
	}

	// Validate cultural wisdom
	wisdomIDs := make(map[uint64]bool)
	for _, wisdom := range data.CulturalWisdom {
		if err := validateCulturalWisdom(wisdom); err != nil {
			return fmt.Errorf("invalid cultural wisdom %d: %w", wisdom.Id, err)
		}
		if wisdomIDs[wisdom.Id] {
			return fmt.Errorf("duplicate cultural wisdom ID: %d", wisdom.Id)
		}
		wisdomIDs[wisdom.Id] = true
	}

	// Validate transaction quotes
	txHashes := make(map[string]bool)
	for _, txQuote := range data.TransactionQuotes {
		if err := validateTransactionQuote(txQuote); err != nil {
			return fmt.Errorf("invalid transaction quote: %w", err)
		}
		if txHashes[txQuote.TxHash] {
			return fmt.Errorf("duplicate transaction hash: %s", txQuote.TxHash)
		}
		txHashes[txQuote.TxHash] = true
		
		// Verify referenced quote exists
		if !quoteIDs[txQuote.QuoteId] && txQuote.QuoteId > uint64(len(DefaultQuotes())) {
			return fmt.Errorf("transaction quote references non-existent quote ID: %d", txQuote.QuoteId)
		}
	}

	return nil
}

// validateHistoricalEvent validates a historical event
func validateHistoricalEvent(event HistoricalEvent) error {
	if event.Id == 0 {
		return fmt.Errorf("event ID cannot be 0")
	}
	if len(event.Title) == 0 {
		return fmt.Errorf("event title cannot be empty")
	}
	if len(event.Description) == 0 {
		return fmt.Errorf("event description cannot be empty")
	}
	if event.Year < -5000 || event.Year > 2100 {
		return fmt.Errorf("invalid event year: %d", event.Year)
	}
	if len(event.Category) == 0 {
		return fmt.Errorf("event category cannot be empty")
	}
	return nil
}

// validateCulturalWisdom validates cultural wisdom
func validateCulturalWisdom(wisdom CulturalWisdom) error {
	if wisdom.Id == 0 {
		return fmt.Errorf("wisdom ID cannot be 0")
	}
	if len(wisdom.Text) == 0 {
		return fmt.Errorf("wisdom text cannot be empty")
	}
	if len(wisdom.Tradition) == 0 {
		return fmt.Errorf("wisdom tradition cannot be empty")
	}
	if len(wisdom.Language) == 0 {
		return fmt.Errorf("wisdom language cannot be empty")
	}
	if wisdom.SpiritualSignificance < 0 || wisdom.SpiritualSignificance > 10 {
		return fmt.Errorf("spiritual significance must be between 0 and 10")
	}
	return nil
}

// validateTransactionQuote validates a transaction quote
func validateTransactionQuote(txQuote TransactionQuote) error {
	if len(txQuote.TxHash) == 0 {
		return fmt.Errorf("transaction hash cannot be empty")
	}
	if txQuote.QuoteId == 0 {
		return fmt.Errorf("quote ID cannot be 0")
	}
	if txQuote.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be positive")
	}
	if len(txQuote.Sender) == 0 {
		return fmt.Errorf("sender address cannot be empty")
	}
	if len(txQuote.Receiver) == 0 {
		return fmt.Errorf("receiver address cannot be empty")
	}
	return nil
}