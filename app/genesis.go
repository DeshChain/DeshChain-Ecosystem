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

package app

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	
	namotypes "github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
	taxtypes "github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
	dextypes "github.com/DeshChain/DeshChain-Ecosystem/x/dex/types"
	launchpadtypes "github.com/DeshChain/DeshChain-Ecosystem/x/launchpad/types"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	genesis := ModuleBasics.DefaultGenesis(cdc)
	
	// Configure NAMO token genesis with v2.0 economic model
	genesis = configureNAMOGenesis(cdc, genesis)
	
	// Configure tax distribution genesis
	genesis = configureTaxGenesis(cdc, genesis)
	
	// Configure DEX genesis
	genesis = configureDEXGenesis(cdc, genesis)
	
	// Configure Sikkebaaz launchpad genesis
	genesis = configureLaunchpadGenesis(cdc, genesis)
	
	// Configure bank genesis with initial allocations
	genesis = configureBankGenesis(cdc, genesis)
	
	return genesis
}

// configureNAMOGenesis configures the NAMO token genesis state with v2.0 allocations
func configureNAMOGenesis(cdc codec.JSONCodec, genesis GenesisState) GenesisState {
	// Create initial token distribution events for transparency
	distributionEvents := []namotypes.DistributionEvent{
		{
			Id:          1,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.PublicSalePoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.PublicSaleAllocation)),
			Percentage:  sdk.NewDecWithPrec(20, 2), // 20%
			Description: "Public sale allocation - reduced for scarcity",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          2,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.LiquidityPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.LiquidityAllocation)),
			Percentage:  sdk.NewDecWithPrec(18, 2), // 18%
			Description: "Liquidity provision allocation",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          3,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.CommunityPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CommunityAllocation)),
			Percentage:  sdk.NewDecWithPrec(15, 2), // 15%
			Description: "Community rewards allocation",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          4,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.DevelopmentPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.DevelopmentAllocation)),
			Percentage:  sdk.NewDecWithPrec(15, 2), // 15%
			Description: "Development fund allocation",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          5,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.TeamPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TeamAllocation)),
			Percentage:  sdk.NewDecWithPrec(12, 2), // 12%
			Description: "Team allocation with 24-month vesting and 12-month cliff",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          6,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.FounderPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.FounderAllocation)),
			Percentage:  sdk.NewDecWithPrec(8, 2), // 8%
			Description: "Founder allocation with 48-month vesting and 12-month cliff",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          7,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.DAOTreasuryPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.DAOTreasuryAllocation)),
			Percentage:  sdk.NewDecWithPrec(5, 2), // 5%
			Description: "DAO treasury allocation",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          8,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.CoFounderPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CoFounderAllocation)),
			Percentage:  sdk.NewDecWithPrec(35, 3), // 3.5%
			Description: "Co-founder allocation with 24-month vesting and 12-month cliff",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          9,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.OperationsPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.OperationsAllocation)),
			Percentage:  sdk.NewDecWithPrec(2, 2), // 2%
			Description: "Operations allocation",
			Timestamp:   time.Now().Unix(),
		},
		{
			Id:          10,
			EventType:   namotypes.EventTypeInitialDistribution,
			Recipient:   namotypes.AngelPoolName,
			Amount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.AngelAllocation)),
			Percentage:  sdk.NewDecWithPrec(15, 3), // 1.5%
			Description: "Angel investor allocation with 24-month vesting and 12-month cliff",
			Timestamp:   time.Now().Unix(),
		},
	}
	
	// Create vesting schedules for vested allocations
	vestingSchedules := []namotypes.VestingSchedule{
		{
			Id:               1,
			PoolName:         namotypes.FounderPoolName,
			TotalAmount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.FounderAllocation)),
			VestingMonths:    namotypes.FounderVestingMonths,
			CliffMonths:      namotypes.UniversalCliffMonths,
			StartTime:        time.Now().Unix(),
			CliffEnd:         time.Now().AddDate(0, namotypes.UniversalCliffMonths, 0).Unix(),
			VestingEnd:       time.Now().AddDate(0, namotypes.FounderVestingMonths, 0).Unix(),
			MonthlyRelease:   sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.FounderAllocation/namotypes.FounderVestingMonths)),
			ClaimedAmount:    sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.ZeroInt()),
			LastClaimTime:    0,
		},
		{
			Id:               2,
			PoolName:         namotypes.TeamPoolName,
			TotalAmount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TeamAllocation)),
			VestingMonths:    namotypes.TeamVestingMonths,
			CliffMonths:      namotypes.UniversalCliffMonths,
			StartTime:        time.Now().Unix(),
			CliffEnd:         time.Now().AddDate(0, namotypes.UniversalCliffMonths, 0).Unix(),
			VestingEnd:       time.Now().AddDate(0, namotypes.TeamVestingMonths, 0).Unix(),
			MonthlyRelease:   sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TeamAllocation/namotypes.TeamVestingMonths)),
			ClaimedAmount:    sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.ZeroInt()),
			LastClaimTime:    0,
		},
		{
			Id:               3,
			PoolName:         namotypes.CoFounderPoolName,
			TotalAmount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CoFounderAllocation)),
			VestingMonths:    namotypes.CoFounderVestingMonths,
			CliffMonths:      namotypes.UniversalCliffMonths,
			StartTime:        time.Now().Unix(),
			CliffEnd:         time.Now().AddDate(0, namotypes.UniversalCliffMonths, 0).Unix(),
			VestingEnd:       time.Now().AddDate(0, namotypes.CoFounderVestingMonths, 0).Unix(),
			MonthlyRelease:   sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CoFounderAllocation/namotypes.CoFounderVestingMonths)),
			ClaimedAmount:    sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.ZeroInt()),
			LastClaimTime:    0,
		},
		{
			Id:               4,
			PoolName:         namotypes.AngelPoolName,
			TotalAmount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.AngelAllocation)),
			VestingMonths:    namotypes.AngelVestingMonths,
			CliffMonths:      namotypes.UniversalCliffMonths,
			StartTime:        time.Now().Unix(),
			CliffEnd:         time.Now().AddDate(0, namotypes.UniversalCliffMonths, 0).Unix(),
			VestingEnd:       time.Now().AddDate(0, namotypes.AngelVestingMonths, 0).Unix(),
			MonthlyRelease:   sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.AngelAllocation/namotypes.AngelVestingMonths)),
			ClaimedAmount:    sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.ZeroInt()),
			LastClaimTime:    0,
		},
		{
			Id:               5,
			PoolName:         namotypes.CommunityPoolName,
			TotalAmount:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CommunityAllocation)),
			VestingMonths:    namotypes.CommunityDistribMonths,
			CliffMonths:      0, // No cliff for community distribution
			StartTime:        time.Now().Unix(),
			CliffEnd:         time.Now().Unix(), // Immediate start
			VestingEnd:       time.Now().AddDate(0, namotypes.CommunityDistribMonths, 0).Unix(),
			MonthlyRelease:   sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CommunityAllocation/namotypes.CommunityDistribMonths)),
			ClaimedAmount:    sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.ZeroInt()),
			LastClaimTime:    0,
		},
	}
	
	// Create NAMO genesis state
	namoGenesis := namotypes.GenesisState{
		TokenSupply: namotypes.TokenSupply{
			TotalSupply:       sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TotalSupply)),
			CirculatingSupply: sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.PublicSaleAllocation + namotypes.LiquidityAllocation)),
			LockedSupply:      sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TotalSupply - namotypes.PublicSaleAllocation - namotypes.LiquidityAllocation)),
			LastUpdated:       time.Now().Unix(),
		},
		DistributionEvents:    distributionEvents,
		VestingSchedules:      vestingSchedules,
		VestingScheduleCount:  uint64(len(vestingSchedules)),
		DistributionEventCount: uint64(len(distributionEvents)),
	}
	
	genesis[namotypes.ModuleName] = cdc.MustMarshalJSON(&namoGenesis)
	return genesis
}

// configureTaxGenesis configures the tax distribution genesis state
func configureTaxGenesis(cdc codec.JSONCodec, genesis GenesisState) GenesisState {
	taxGenesis := taxtypes.GenesisState{
		Params: taxtypes.Params{
			TaxDistribution: taxtypes.NewDefaultTaxDistribution(),
			PlatformDistribution: taxtypes.NewDefaultPlatformDistribution(),
		},
	}
	
	genesis[taxtypes.ModuleName] = cdc.MustMarshalJSON(&taxGenesis)
	return genesis
}

// configureDEXGenesis configures the Money Order DEX genesis state
func configureDEXGenesis(cdc codec.JSONCodec, genesis GenesisState) GenesisState {
	dexGenesis := dextypes.GenesisState{
		Params: dextypes.NewDefaultDEXParams(),
	}
	
	genesis[dextypes.ModuleName] = cdc.MustMarshalJSON(&dexGenesis)
	return genesis
}

// configureLaunchpadGenesis configures the Sikkebaaz launchpad genesis state
func configureLaunchpadGenesis(cdc codec.JSONCodec, genesis GenesisState) GenesisState {
	launchpadGenesis := launchpadtypes.GenesisState{
		Params: launchpadtypes.NewDefaultSikkebaazParams(),
	}
	
	genesis[launchpadtypes.ModuleName] = cdc.MustMarshalJSON(&launchpadGenesis)
	return genesis
}

// configureBankGenesis configures the bank genesis state with initial module account balances
func configureBankGenesis(cdc codec.JSONCodec, genesis GenesisState) GenesisState {
	var bankGenesis banktypes.GenesisState
	cdc.MustUnmarshalJSON(genesis[banktypes.ModuleName], &bankGenesis)
	
	// Add initial balances for module accounts based on v2.0 allocation
	moduleBalances := []banktypes.Balance{
		{
			Address: namotypes.PublicSalePoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.PublicSaleAllocation))),
		},
		{
			Address: namotypes.LiquidityPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.LiquidityAllocation))),
		},
		{
			Address: namotypes.CommunityPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CommunityAllocation))),
		},
		{
			Address: namotypes.DevelopmentPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.DevelopmentAllocation))),
		},
		{
			Address: namotypes.TeamPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TeamAllocation))),
		},
		{
			Address: namotypes.FounderPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.FounderAllocation))),
		},
		{
			Address: namotypes.DAOTreasuryPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.DAOTreasuryAllocation))),
		},
		{
			Address: namotypes.CoFounderPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.CoFounderAllocation))),
		},
		{
			Address: namotypes.OperationsPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.OperationsAllocation))),
		},
		{
			Address: namotypes.AngelPoolName,
			Coins:   sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.AngelAllocation))),
		},
	}
	
	// Add balances to existing balances
	bankGenesis.Balances = append(bankGenesis.Balances, moduleBalances...)
	
	// Update total supply
	totalSupply := sdk.NewCoins(sdk.NewCoin(namotypes.DefaultTokenDenom, sdk.NewInt(namotypes.TotalSupply)))
	bankGenesis.Supply = totalSupply
	
	genesis[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenesis)
	return genesis
}

// GetGenesisStateFromAppState unmarshals the GenesisState from the provided app state
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	genesisState := make(GenesisState)
	
	for key, value := range appState {
		genesisState[key] = value
	}
	
	return genesisState
}