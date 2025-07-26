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

package moneyorder

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// InitGenesis initializes the money order module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	k.SetParams(ctx, genState.Params)
	
	// Set next pool ID
	if genState.NextPoolId > 0 {
		store := ctx.KVStore(k.storeKey)
		store.Set(types.KeyPrefixSequence, sdk.Uint64ToBigEndian(genState.NextPoolId))
	}
	
	// Initialize fixed rate pools
	for _, pool := range genState.FixedRatePools {
		k.SetFixedRatePool(ctx, pool)
	}
	
	// Initialize village pools
	for _, pool := range genState.VillagePools {
		k.SetVillagePool(ctx, pool)
	}
	
	// Initialize money order receipts
	for _, receipt := range genState.MoneyOrderReceipts {
		k.SetMoneyOrderReceipt(ctx, receipt)
	}
	
	// Initialize village pool members
	for _, membership := range genState.VillagePoolMembers {
		k.SetVillagePoolMember(ctx, membership.PoolId, membership.Member)
	}
	
	// Initialize UPI addresses
	for _, upiMapping := range genState.UPIAddressMappings {
		addr, _ := sdk.AccAddressFromBech32(upiMapping.Address)
		k.RegisterUPIAddress(ctx, upiMapping.UPIAddress, addr)
	}
}

// ExportGenesis returns the money order module's exported genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	
	// Export next pool ID
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixSequence)
	if bz != nil {
		genesis.NextPoolId = sdk.BigEndianToUint64(bz)
	}
	
	// Export fixed rate pools
	genesis.FixedRatePools = k.GetAllFixedRatePools(ctx)
	
	// Export village pools
	genesis.VillagePools = k.GetAllVillagePools(ctx)
	
	// Export money order receipts
	// Note: In production, would only export recent/active receipts
	genesis.MoneyOrderReceipts = k.GetAllMoneyOrderReceipts(ctx)
	
	// Export village pool members
	for _, pool := range genesis.VillagePools {
		members := k.GetAllVillagePoolMembers(ctx, pool.PoolId)
		for _, member := range members {
			genesis.VillagePoolMembers = append(genesis.VillagePoolMembers, types.VillagePoolMembership{
				PoolId: pool.PoolId,
				Member: member,
			})
		}
	}
	
	return genesis
}