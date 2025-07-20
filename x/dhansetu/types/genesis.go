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

import "fmt"

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() GenesisState {
	return GenesisState{
		Params:              DefaultParams(),
		DhanpataAddresses:   []DhanPataAddress{},
		MitraProfiles:       []EnhancedMitraProfile{},
		KshetraCoins:        []KshetraCoin{},
		CrossModuleBridges:  []CrossModuleBridge{},
		TradeHistory:        []TradeHistoryEntry{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate DhanPata addresses
	dhanpataNames := make(map[string]bool)
	for _, address := range gs.DhanpataAddresses {
		if err := ValidateDhanPataName(address.Name); err != nil {
			return fmt.Errorf("invalid DhanPata address name %s: %w", address.Name, err)
		}
		
		if dhanpataNames[address.Name] {
			return fmt.Errorf("duplicate DhanPata address name: %s", address.Name)
		}
		dhanpataNames[address.Name] = true
	}

	// Validate Kshetra coins
	pincodes := make(map[string]bool)
	for _, coin := range gs.KshetraCoins {
		if err := ValidatePincode(coin.Pincode); err != nil {
			return fmt.Errorf("invalid pincode %s: %w", coin.Pincode, err)
		}
		
		if pincodes[coin.Pincode] {
			return fmt.Errorf("duplicate Kshetra coin for pincode: %s", coin.Pincode)
		}
		pincodes[coin.Pincode] = true
	}

	// Validate Mitra profiles
	mitraIds := make(map[string]bool)
	for _, profile := range gs.MitraProfiles {
		if mitraIds[profile.MitraId] {
			return fmt.Errorf("duplicate mitra ID: %s", profile.MitraId)
		}
		mitraIds[profile.MitraId] = true
		
		if profile.MitraType != MitraTypeIndividual &&
			profile.MitraType != MitraTypeBusiness &&
			profile.MitraType != MitraTypeGlobal {
			return fmt.Errorf("invalid mitra type: %s", profile.MitraType)
		}
		
		if profile.TrustScore < 0 || profile.TrustScore > 100 {
			return fmt.Errorf("invalid trust score for mitra %s: %d", profile.MitraId, profile.TrustScore)
		}
	}

	return nil
}

// GenesisState defines the dhansetu module's genesis state.
type GenesisState struct {
	Params              Params                    `json:"params" yaml:"params"`
	DhanpataAddresses   []DhanPataAddress         `json:"dhanpata_addresses" yaml:"dhanpata_addresses"`
	MitraProfiles       []EnhancedMitraProfile    `json:"mitra_profiles" yaml:"mitra_profiles"`
	KshetraCoins        []KshetraCoin             `json:"kshetra_coins" yaml:"kshetra_coins"`
	CrossModuleBridges  []CrossModuleBridge       `json:"cross_module_bridges" yaml:"cross_module_bridges"`
	TradeHistory        []TradeHistoryEntry       `json:"trade_history" yaml:"trade_history"`
}