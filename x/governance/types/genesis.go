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
	"time"
)

// NewGenesisState creates a new governance genesis state
func NewGenesisState(
	founderAddress string,
	phase GovernancePhase,
	genesisTime time.Time,
) *GenesisState {
	return &GenesisState{
		FounderAddress:     founderAddress,
		CurrentPhase:       phase,
		GenesisTime:        genesisTime,
		ProtectedParams:    DefaultProtectedParams(),
	}
}

// DefaultGenesis returns the default governance genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		FounderAddress:     "", // Must be set at genesis
		CurrentPhase:       GovernancePhase_FOUNDER_CONTROL,
		GenesisTime:        time.Time{},
		ProtectedParams:    DefaultProtectedParams(),
	}
}

// ValidateGenesis validates the governance genesis state
func ValidateGenesis(data *GenesisState) error {
	if data.FounderAddress == "" {
		return fmt.Errorf("founder address cannot be empty")
	}

	if data.GenesisTime.IsZero() {
		return fmt.Errorf("genesis time must be set")
	}

	// Validate protected parameters
	for _, param := range data.ProtectedParams {
		if param.Name == "" {
			return fmt.Errorf("protected parameter name cannot be empty")
		}
		if param.Protection < ProtectionType_NONE || param.Protection > ProtectionType_SUPERMAJORITY {
			return fmt.Errorf("invalid protection type for parameter %s", param.Name)
		}
	}

	return nil
}

// DefaultProtectedParams returns the default protected parameters
func DefaultProtectedParams() []ProtectedParameter {
	return []ProtectedParameter{
		// Immutable parameters
		{Name: "founder_token_allocation", Protection: ProtectionType_IMMUTABLE},
		{Name: "founder_tax_royalty", Protection: ProtectionType_IMMUTABLE},
		{Name: "founder_platform_royalty", Protection: ProtectionType_IMMUTABLE},
		{Name: "founder_inheritance_mechanism", Protection: ProtectionType_IMMUTABLE},
		{Name: "founder_minimum_voting_power", Protection: ProtectionType_IMMUTABLE},
		
		// Founder consent required
		{Name: "chain_upgrade_handler", Protection: ProtectionType_FOUNDER_CONSENT},
		{Name: "crisis_module_permissions", Protection: ProtectionType_FOUNDER_CONSENT},
		{Name: "slashing_parameters", Protection: ProtectionType_FOUNDER_CONSENT},
		{Name: "consensus_parameters", Protection: ProtectionType_FOUNDER_CONSENT},
		
		// Supermajority required
		{Name: "governance_voting_period", Protection: ProtectionType_SUPERMAJORITY},
		{Name: "governance_deposit_amount", Protection: ProtectionType_SUPERMAJORITY},
		{Name: "distribution_community_tax", Protection: ProtectionType_SUPERMAJORITY},
		{Name: "mint_inflation_rate", Protection: ProtectionType_SUPERMAJORITY},
	}
}

// GenesisState defines the governance module's genesis state
type GenesisState struct {
	// Founder address (immutable)
	FounderAddress string `json:"founder_address" yaml:"founder_address"`
	
	// Current governance phase
	CurrentPhase GovernancePhase `json:"current_phase" yaml:"current_phase"`
	
	// Genesis time for phase calculations
	GenesisTime time.Time `json:"genesis_time" yaml:"genesis_time"`
	
	// Protected parameters
	ProtectedParams []ProtectedParameter `json:"protected_params" yaml:"protected_params"`
	
	// Vetoed proposals
	VetoedProposals []uint64 `json:"vetoed_proposals" yaml:"vetoed_proposals"`
}

// ProtectedParameter defines a parameter with protection level
type ProtectedParameter struct {
	Name       string         `json:"name" yaml:"name"`
	Protection ProtectionType `json:"protection" yaml:"protection"`
}