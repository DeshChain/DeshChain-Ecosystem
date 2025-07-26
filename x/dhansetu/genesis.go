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

package dhansetu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dhansetu/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dhansetu/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	k.SetParams(ctx, genState.Params)

	// Initialize DhanPata addresses
	for _, address := range genState.DhanpataAddresses {
		k.SetDhanPataAddress(ctx, address)
		k.SetAddressToDhanPata(ctx, address.Owner, address.Name)
	}

	// Initialize Enhanced Mitra profiles
	for _, profile := range genState.MitraProfiles {
		k.SetEnhancedMitraProfile(ctx, profile)
	}

	// Initialize Kshetra coins
	for _, coin := range genState.KshetraCoins {
		k.SetKshetraCoin(ctx, coin)
	}

	// Initialize cross-module bridges
	for _, bridge := range genState.CrossModuleBridges {
		k.SetCrossModuleBridge(ctx, bridge)
	}

	// Initialize trade history
	for _, trade := range genState.TradeHistory {
		k.RecordTradeHistory(ctx, trade)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export DhanPata addresses
	genesis.DhanpataAddresses = k.GetAllDhanPataAddresses(ctx)

	// Export Kshetra coins
	genesis.KshetraCoins = k.GetAllKshetraCoins(ctx)

	// Note: Other exports (mitra profiles, bridges, trade history) would be implemented
	// with appropriate iteration methods in the keeper

	return genesis
}