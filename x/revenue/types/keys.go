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

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "revenue"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_revenue"
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = collections.NewPrefix(0)

	// RevenueStreamKey is the key for revenue streams
	RevenueStreamKey = collections.NewPrefix(1)

	// DistributionRecordKey is the key for distribution records
	DistributionRecordKey = collections.NewPrefix(2)

	// FounderRoyaltyKey is the key for founder royalty configuration
	FounderRoyaltyKey = collections.NewPrefix(3)

	// PlatformRevenueKey is the key for platform revenue tracking
	PlatformRevenueKey = collections.NewPrefix(4)

	// MonthlyRevenueKey is the key for monthly revenue tracking
	MonthlyRevenueKey = collections.NewPrefix(5)

	// RoyaltyClaimKey is the key for royalty claims
	RoyaltyClaimKey = collections.NewPrefix(6)
)

// Revenue stream types
const (
	RevenueStreamDEX         = "dex_trading"
	RevenueStreamNFT         = "nft_marketplace"
	RevenueStreamLaunchpad   = "sikkebaaz_launchpad"
	RevenueStreamPension     = "gram_pension"
	RevenueStreamLending     = "kisaan_mitra"
	RevenueStreamPrivacy     = "privacy_fees"
	RevenueStreamGovernance  = "governance_fees"
	RevenueStreamOther       = "other_revenue"
)