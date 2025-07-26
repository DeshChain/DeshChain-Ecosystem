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
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// Keeper of the gram pension store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	memKey   sdk.StoreKey

	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	moneyOrderKeeper types.MoneyOrderKeeper
	culturalKeeper   types.CulturalKeeper
	taxKeeper        types.TaxKeeper
	donationKeeper   types.DonationKeeper
	kycKeeper        types.KYCKeeper
	revenueKeeper    types.RevenueKeeper

	// Identity integration
	identityAdapter *IdentityAdapter

	// Hooks
	hooks types.GramPensionHooks

	// Authority address for governance
	authority string
}

// NewKeeper creates a new gram pension Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	moneyOrderKeeper types.MoneyOrderKeeper,
	culturalKeeper types.CulturalKeeper,
	taxKeeper types.TaxKeeper,
	donationKeeper types.DonationKeeper,
	kycKeeper types.KYCKeeper,
	revenueKeeper types.RevenueKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		moneyOrderKeeper: moneyOrderKeeper,
		culturalKeeper:   culturalKeeper,
		taxKeeper:        taxKeeper,
		donationKeeper:   donationKeeper,
		kycKeeper:        kycKeeper,
		revenueKeeper:    revenueKeeper,
		authority:        authority,
	}
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the gram pension hooks
func (k *Keeper) SetHooks(gh types.GramPensionHooks) {
	if k.hooks != nil {
		panic("cannot set gram pension hooks twice")
	}
	k.hooks = gh
}

// SetIdentityKeeper sets the identity keeper for integration
func (k *Keeper) SetIdentityKeeper(identityKeeper interface{}) {
	// Type assertion to avoid circular imports
	if ik, ok := identityKeeper.(interface {
		GetIdentity(sdk.Context, string) (interface{}, bool)
		SetIdentity(sdk.Context, interface{})
		GetCredential(sdk.Context, string) (interface{}, bool)
		SetCredential(sdk.Context, interface{})
		GetCredentialsBySubject(sdk.Context, string) []string
		AddCredentialToSubject(sdk.Context, string, string)
		UpdateCredentialStatus(sdk.Context, string, interface{}) error
	}); ok {
		// Create wrapper that implements our expected interface
		wrapper := &identityKeeperWrapper{keeper: ik}
		k.identityAdapter = NewIdentityAdapter(k, wrapper)
	}
}
}

// GetHooks returns the gram pension hooks
func (k Keeper) GetHooks() types.GramPensionHooks {
	return k.hooks
}

// CallHookAfterContribution calls the after contribution hook if set
func (k Keeper) CallHookAfterContribution(ctx sdk.Context, pensionAccountId string, contributor sdk.AccAddress, contribution sdk.Coin, villagePostalCode string) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.AfterSurakshaContribution(ctx, pensionAccountId, contributor, contribution, villagePostalCode)
}

// CallHookAfterMaturity calls the after maturity hook if set
func (k Keeper) CallHookAfterMaturity(ctx sdk.Context, pensionAccountId string, beneficiary sdk.AccAddress, maturityAmount sdk.Coin) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.AfterSurakshaMaturity(ctx, pensionAccountId, beneficiary, maturityAmount)
}

// CallHookMonthlyDistribution calls the monthly distribution hook if set
func (k Keeper) CallHookMonthlyDistribution(ctx sdk.Context) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.MonthlyRevenueDistribution(ctx)
}