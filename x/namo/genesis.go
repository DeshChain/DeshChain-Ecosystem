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

package namo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/namo/keeper"
	"github.com/deshchain/deshchain/x/namo/types"
)

// InitGenesis initializes the NAMO module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)

	// Set token supply
	if genState.TokenSupply != nil {
		k.SetTokenSupply(ctx, *genState.TokenSupply)
	} else {
		// Initialize with default supply
		defaultSupply := types.TokenSupply{
			TotalSupply:           types.DefaultTotalSupply,
			FounderAllocation:     types.CalculateAllocation(types.DefaultTotalSupply, types.FounderAllocationPercentage),
			TeamAllocation:        types.CalculateAllocation(types.DefaultTotalSupply, types.TeamAllocationPercentage),
			CommunityAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.CommunityAllocationPercentage),
			DevelopmentAllocation: types.CalculateAllocation(types.DefaultTotalSupply, types.DevelopmentAllocationPercentage),
			LiquidityAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.LiquidityAllocationPercentage),
			PublicSaleAllocation:  types.CalculateAllocation(types.DefaultTotalSupply, types.PublicSaleAllocationPercentage),
			DAOTreasuryAllocation: types.CalculateAllocation(types.DefaultTotalSupply, types.DAOTreasuryAllocationPercentage),
			CoFounderAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.CoFounderAllocationPercentage),
			OperationsAllocation:  types.CalculateAllocation(types.DefaultTotalSupply, types.OperationsAllocationPercentage),
			AngelAllocation:       types.CalculateAllocation(types.DefaultTotalSupply, types.AngelAllocationPercentage),
		}
		k.SetTokenSupply(ctx, defaultSupply)
	}

	// Set vesting schedules
	for _, schedule := range genState.VestingSchedules {
		k.SetVestingSchedule(ctx, schedule)
	}

	// Set token distribution events
	for _, event := range genState.TokenDistributionEvents {
		k.CreateTokenDistributionEvent(ctx, event)
	}

	// Set next distribution event ID
	k.SetNextDistributionEventID(ctx, genState.NextDistributionEventId)
}

// ExportGenesis returns the NAMO module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export token supply
	if tokenSupply, found := k.GetTokenSupply(ctx); found {
		genesis.TokenSupply = &tokenSupply
	}

	// Export vesting schedules
	genesis.VestingSchedules = k.GetAllVestingSchedules(ctx)

	// Export token distribution events
	genesis.TokenDistributionEvents = k.GetAllTokenDistributionEvents(ctx)

	// Export next distribution event ID
	genesis.NextDistributionEventId = k.GetNextDistributionEventID(ctx)

	return genesis
}

// BeginBlocker processes module state at the beginning of each block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Process any begin block logic for NAMO token
	// This could include automated vesting releases, token burns, etc.
	
	// Check for automated token operations
	params := k.GetParams(ctx)
	if params.EnableTokenOperations {
		// Perform any scheduled token operations
		processScheduledOperations(ctx, k)
	}
}

// EndBlocker processes module state at the end of each block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []sdk.ValidatorUpdate {
	// Process any end block logic for NAMO token
	// This could include updating statistics, processing delayed operations, etc.
	
	// No validator updates needed for NAMO module
	return []sdk.ValidatorUpdate{}
}

// processScheduledOperations processes any scheduled token operations
func processScheduledOperations(ctx sdk.Context, k keeper.Keeper) {
	// Check for any vesting schedules that need processing
	// This could be automated claiming for certain scenarios
	
	// For now, this is a placeholder for future automated operations
	// In production, this might include:
	// - Automated founder/team token releases
	// - Scheduled burns
	// - Liquidity operations
}