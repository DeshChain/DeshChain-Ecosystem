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

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CollectModuleFee is a helper function for modules to collect and report fees
// This function handles the fee collection, revenue reporting, and distribution
func (k Keeper) CollectModuleFee(ctx sdk.Context, moduleName string, feePayer sdk.AccAddress, fee sdk.Coins, description string) error {
	if fee.IsZero() {
		return nil
	}

	// Transfer fee from payer to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, feePayer, moduleName, fee); err != nil {
		return err
	}

	// Report the revenue
	if err := k.RecordRevenue(ctx, moduleName, "fee", fee, description); err != nil {
		return err
	}

	// Transfer from module to revenue module for distribution
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleName, k.ModuleName, fee); err != nil {
		return err
	}

	// Distribute the platform revenue
	if err := k.DistributePlatformRevenue(ctx, fee); err != nil {
		return err
	}

	return nil
}

// CollectTradingFee is a specialized helper for trading/swap fees
func (k Keeper) CollectTradingFee(ctx sdk.Context, moduleName string, trader sdk.AccAddress, fee sdk.Coins, pair string) error {
	description := "Trading fee for " + pair
	return k.CollectModuleFee(ctx, moduleName, trader, fee, description)
}

// CollectServiceFee is a specialized helper for service fees
func (k Keeper) CollectServiceFee(ctx sdk.Context, moduleName string, user sdk.AccAddress, fee sdk.Coins, service string) error {
	description := "Service fee for " + service
	return k.CollectModuleFee(ctx, moduleName, user, fee, description)
}

// CollectLaunchFee is a specialized helper for token launch fees
func (k Keeper) CollectLaunchFee(ctx sdk.Context, moduleName string, creator sdk.AccAddress, fee sdk.Coins, tokenSymbol string) error {
	description := "Token launch fee for " + tokenSymbol
	return k.CollectModuleFee(ctx, moduleName, creator, fee, description)
}

// CollectLendingFee is a specialized helper for lending protocol fees
func (k Keeper) CollectLendingFee(ctx sdk.Context, moduleName string, borrower sdk.AccAddress, fee sdk.Coins, loanType string) error {
	description := "Lending fee for " + loanType
	return k.CollectModuleFee(ctx, moduleName, borrower, fee, description)
}