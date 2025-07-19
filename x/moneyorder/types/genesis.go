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

// GenesisState defines the money order module's genesis state
type GenesisState struct {
	Params             MoneyOrderParams        `json:"params" yaml:"params"`
	NextPoolId         uint64                  `json:"next_pool_id" yaml:"next_pool_id"`
	FixedRatePools     []*FixedRatePool        `json:"fixed_rate_pools" yaml:"fixed_rate_pools"`
	VillagePools       []*VillagePool          `json:"village_pools" yaml:"village_pools"`
	MoneyOrderReceipts []*MoneyOrderReceipt    `json:"money_order_receipts" yaml:"money_order_receipts"`
	VillagePoolMembers []VillagePoolMembership `json:"village_pool_members" yaml:"village_pool_members"`
	UPIAddressMappings []UPIAddressMapping     `json:"upi_address_mappings" yaml:"upi_address_mappings"`
}

// VillagePoolMembership represents a member's association with a village pool
type VillagePoolMembership struct {
	PoolId uint64               `json:"pool_id" yaml:"pool_id"`
	Member *VillagePoolMember   `json:"member" yaml:"member"`
}

// UPIAddressMapping represents a UPI address to account address mapping
type UPIAddressMapping struct {
	UPIAddress string `json:"upi_address" yaml:"upi_address"`
	Address    string `json:"address" yaml:"address"`
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:             DefaultParams(),
		NextPoolId:         1,
		FixedRatePools:     []*FixedRatePool{},
		VillagePools:       []*VillagePool{},
		MoneyOrderReceipts: []*MoneyOrderReceipt{},
		VillagePoolMembers: []VillagePoolMembership{},
		UPIAddressMappings: []UPIAddressMapping{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if err := data.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	
	// Validate next pool ID
	if data.NextPoolId == 0 {
		return fmt.Errorf("next pool ID must be greater than 0")
	}
	
	// Validate fixed rate pools
	poolIds := make(map[uint64]bool)
	for _, pool := range data.FixedRatePools {
		if err := pool.ValidatePool(); err != nil {
			return fmt.Errorf("invalid fixed rate pool %d: %w", pool.PoolId, err)
		}
		if poolIds[pool.PoolId] {
			return fmt.Errorf("duplicate pool ID: %d", pool.PoolId)
		}
		poolIds[pool.PoolId] = true
	}
	
	// Validate village pools
	postalCodes := make(map[string]bool)
	for _, pool := range data.VillagePools {
		if err := pool.ValidateVillagePool(); err != nil {
			return fmt.Errorf("invalid village pool %d: %w", pool.PoolId, err)
		}
		if poolIds[pool.PoolId] {
			return fmt.Errorf("duplicate pool ID: %d", pool.PoolId)
		}
		if postalCodes[pool.PostalCode] {
			return fmt.Errorf("duplicate postal code: %s", pool.PostalCode)
		}
		poolIds[pool.PoolId] = true
		postalCodes[pool.PostalCode] = true
	}
	
	// Validate money order receipts
	orderIds := make(map[string]bool)
	for _, receipt := range data.MoneyOrderReceipts {
		if err := receipt.ValidateReceipt(); err != nil {
			return fmt.Errorf("invalid money order receipt %s: %w", receipt.OrderId, err)
		}
		if orderIds[receipt.OrderId] {
			return fmt.Errorf("duplicate order ID: %s", receipt.OrderId)
		}
		orderIds[receipt.OrderId] = true
	}
	
	// Validate village pool members
	for _, membership := range data.VillagePoolMembers {
		if !poolIds[membership.PoolId] {
			return fmt.Errorf("member references non-existent pool: %d", membership.PoolId)
		}
		if membership.Member == nil {
			return fmt.Errorf("nil member for pool %d", membership.PoolId)
		}
	}
	
	// Validate UPI mappings
	upiAddresses := make(map[string]bool)
	for _, mapping := range data.UPIAddressMappings {
		if mapping.UPIAddress == "" {
			return fmt.Errorf("empty UPI address")
		}
		if mapping.Address == "" {
			return fmt.Errorf("empty address for UPI %s", mapping.UPIAddress)
		}
		if upiAddresses[mapping.UPIAddress] {
			return fmt.Errorf("duplicate UPI address: %s", mapping.UPIAddress)
		}
		upiAddresses[mapping.UPIAddress] = true
	}
	
	return nil
}