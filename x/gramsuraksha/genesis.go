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

package grampension

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// GenesisState defines the gram pension module's genesis state.
type GenesisState struct {
	Params              types.Params              `json:"params"`
	Schemes             []types.SurakshaScheme     `json:"schemes"`
	Participants        []types.SurakshaParticipant `json:"participants"`
	Contributions       []types.SurakshaContribution `json:"contributions"`
	Maturities          []types.SurakshaMaturity   `json:"maturities"`
	Withdrawals         []types.PensionWithdrawal `json:"withdrawals"`
	Statistics          []types.PensionStatistics `json:"statistics"`
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Schemes:       []types.SurakshaScheme{},
		Participants:  []types.SurakshaParticipant{},
		Contributions: []types.SurakshaContribution{},
		Maturities:    []types.SurakshaMaturity{},
		Withdrawals:   []types.PensionWithdrawal{},
		Statistics:    []types.PensionStatistics{},
	}
}

// DefaultParams returns default parameters for the gram pension module
func DefaultParams() types.Params {
	return types.Params{
		// Default parameters will be defined when the Params type is created
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data *GenesisState) error {
	// Validate schemes
	for _, scheme := range data.Schemes {
		if err := scheme.Validate(); err != nil {
			return err
		}
	}

	// Validate participants
	for _, participant := range data.Participants {
		if err := participant.Validate(); err != nil {
			return err
		}
	}

	// Validate that participants reference valid schemes
	schemeIDs := make(map[string]bool)
	for _, scheme := range data.Schemes {
		schemeIDs[scheme.SchemeID] = true
	}

	for _, participant := range data.Participants {
		if !schemeIDs[participant.SchemeID] {
			return types.ErrSchemeNotFound
		}
	}

	// Validate that contributions reference valid participants
	participantIDs := make(map[string]bool)
	for _, participant := range data.Participants {
		participantIDs[participant.ParticipantID] = true
	}

	for _, contribution := range data.Contributions {
		if !participantIDs[contribution.ParticipantID] {
			return types.ErrParticipantNotFound
		}
	}

	return nil
}

// InitGenesis initializes the gram pension module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *GenesisState) {
	// Initialize schemes
	for _, scheme := range genState.Schemes {
		k.SetScheme(ctx, scheme)
	}

	// Initialize participants
	for _, participant := range genState.Participants {
		k.SetParticipant(ctx, participant)
		
		// Create indexes
		k.SetParticipantByAddress(ctx, participant.Address, participant.SchemeID, participant.ParticipantID)
		k.SetParticipantByScheme(ctx, participant.SchemeID, participant.ParticipantID)
	}

	// Initialize contributions
	for _, contribution := range genState.Contributions {
		k.SetContribution(ctx, contribution)
	}

	// Initialize maturities
	for _, maturity := range genState.Maturities {
		k.SetMaturity(ctx, maturity)
	}

	// Initialize statistics
	for _, stats := range genState.Statistics {
		k.SetSchemeStatistics(ctx, stats)
	}
}

// ExportGenesis returns the gram pension module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *GenesisState {
	genesis := DefaultGenesis()

	// Export schemes
	genesis.Schemes = k.GetAllSchemes(ctx)

	// Export participants
	var participants []types.SurakshaParticipant
	for _, scheme := range genesis.Schemes {
		k.IterateParticipantsByScheme(ctx, scheme.SchemeID, func(participant types.SurakshaParticipant) bool {
			participants = append(participants, participant)
			return false
		})
	}
	genesis.Participants = participants

	// Export contributions
	var contributions []types.SurakshaContribution
	for _, participant := range participants {
		participantContributions := k.GetParticipantContributions(ctx, participant.ParticipantID)
		contributions = append(contributions, participantContributions...)
	}
	genesis.Contributions = contributions

	// Export statistics
	var statistics []types.PensionStatistics
	for _, scheme := range genesis.Schemes {
		if stats, found := k.GetSchemeStatistics(ctx, scheme.SchemeID); found {
			statistics = append(statistics, stats)
		}
	}
	genesis.Statistics = statistics

	return genesis
}