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
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/keeper"
)

// BeginBlocker processes scheduled orders and other periodic tasks
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Process any scheduled money orders that are due
	k.ProcessScheduledOrders(ctx)
	
	// Reset daily limits at midnight (simplified - in production would use proper time zones)
	if ctx.BlockTime().Hour() == 0 && ctx.BlockTime().Minute() == 0 {
		k.ResetDailyLimits(ctx)
	}
	
	// Check and update festival periods
	k.UpdateFestivalStatus(ctx)
}

// EndBlocker returns the validator set updates
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// No validator updates from this module
	return []abci.ValidatorUpdate{}
}