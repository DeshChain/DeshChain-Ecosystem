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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PoolAsset represents a token in an AMM pool
type PoolAsset struct {
	Denom  string  `json:"denom"`
	Amount sdk.Int `json:"amount"`
	Weight sdk.Int `json:"weight"` // For weighted pools (future use)
}

// AMMPoolInfo represents an automated market maker pool
type AMMPoolInfo struct {
	PoolId          uint64      `json:"pool_id"`
	PoolAssets      []PoolAsset `json:"pool_assets"`
	TotalShares     sdk.Int     `json:"total_shares"`
	SwapFee         sdk.Dec     `json:"swap_fee"`
	Creator         string      `json:"creator"`
	Active          bool        `json:"active"`
	CulturalPair    bool        `json:"cultural_pair"`
	FestivalBonus   bool        `json:"festival_bonus"`
	VillagePriority bool        `json:"village_priority"`
}

// GetAMMPoolKey returns the store key for an AMM pool
func GetAMMPoolKey(poolId uint64) []byte {
	return append(KeyPrefixAMMPool, sdk.Uint64ToBigEndian(poolId)...)
}

// Validate validates a PoolAsset
func (p PoolAsset) Validate() error {
	if err := sdk.ValidateDenom(p.Denom); err != nil {
		return err
	}
	if p.Amount.IsNegative() {
		return ErrInvalidPoolAssets
	}
	return nil
}

// Validate validates an AMMPoolInfo
func (p AMMPoolInfo) Validate() error {
	if p.PoolId == 0 {
		return ErrInvalidPoolId
	}
	if len(p.PoolAssets) < 2 {
		return ErrInvalidPoolAssets
	}
	for _, asset := range p.PoolAssets {
		if err := asset.Validate(); err != nil {
			return err
		}
	}
	if p.TotalShares.IsNegative() {
		return ErrInvalidShares
	}
	if p.SwapFee.IsNegative() || p.SwapFee.GTE(sdk.OneDec()) {
		return ErrInvalidSwapFee
	}
	if _, err := sdk.AccAddressFromBech32(p.Creator); err != nil {
		return err
	}
	return nil
}