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
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/DeshChain/DeshChain-Ecosystem/x/governance/types"
)

// Keeper of the governance store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	govKeeper     types.GovKeeper
	stakingKeeper types.StakingKeeper

	// Authority address for governance actions
	authority string
}

// NewKeeper creates a new governance Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	govKeeper types.GovKeeper,
	stakingKeeper types.StakingKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		govKeeper:     govKeeper,
		stakingKeeper: stakingKeeper,
		authority:     authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the governance module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// InitGenesis initializes the governance module's state from a genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	// Set founder address
	k.SetFounderAddress(ctx, genState.FounderAddress)
	
	// Set current phase
	k.SetGovernancePhase(ctx, genState.CurrentPhase)
	
	// Set genesis time
	k.SetGenesisTime(ctx, genState.GenesisTime)
	
	// Set protected parameters
	for _, param := range genState.ProtectedParams {
		k.SetProtectedParameter(ctx, param)
	}
	
	// Set vetoed proposals
	for _, proposalID := range genState.VetoedProposals {
		k.SetProposalVetoed(ctx, proposalID)
	}
}

// ExportGenesis returns the governance module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		FounderAddress:  k.GetFounderAddress(ctx),
		CurrentPhase:    k.GetGovernancePhase(ctx),
		GenesisTime:     k.GetGenesisTime(ctx),
		ProtectedParams: k.GetAllProtectedParameters(ctx),
		VetoedProposals: k.GetAllVetoedProposals(ctx),
	}
}

// ProcessPhaseTransitions checks and processes governance phase transitions
func (k Keeper) ProcessPhaseTransitions(ctx sdk.Context) {
	currentPhase := k.GetGovernancePhase(ctx)
	genesisTime := k.GetGenesisTime(ctx)
	currentTime := ctx.BlockTime()
	
	timeSinceGenesis := currentTime.Sub(genesisTime)
	yearsSinceGenesis := int(timeSinceGenesis.Hours() / (24 * 365))
	
	newPhase := currentPhase
	
	switch currentPhase {
	case types.GovernancePhase_FOUNDER_CONTROL:
		if yearsSinceGenesis >= 2 {
			newPhase = types.GovernancePhase_SHARED_GOVERNANCE
		}
	case types.GovernancePhase_SHARED_GOVERNANCE:
		if yearsSinceGenesis >= 3 {
			newPhase = types.GovernancePhase_COMMUNITY_GOVERNANCE
		}
	}
	
	if newPhase != currentPhase {
		k.SetGovernancePhase(ctx, newPhase)
		k.Logger(ctx).Info("Governance phase transitioned", 
			"from", currentPhase.String(), 
			"to", newPhase.String(),
			"years_since_genesis", yearsSinceGenesis)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"governance_phase_transition",
				sdk.NewAttribute("from_phase", currentPhase.String()),
				sdk.NewAttribute("to_phase", newPhase.String()),
				sdk.NewAttribute("years_since_genesis", fmt.Sprintf("%d", yearsSinceGenesis)),
			),
		)
	}
}

// ProcessProposals processes governance proposals with founder protections
func (k Keeper) ProcessProposals(ctx sdk.Context) {
	// This is called in EndBlock to process proposals
	// The actual proposal processing is handled by the standard gov module
	// We just add our custom validations via hooks
}

// Store functions

// SetFounderAddress stores the founder address
func (k Keeper) SetFounderAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyFounderAddress, []byte(address))
}

// GetFounderAddress retrieves the founder address
func (k Keeper) GetFounderAddress(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyFounderAddress)
	return string(bz)
}

// SetGovernancePhase stores the current governance phase
func (k Keeper) SetGovernancePhase(ctx sdk.Context, phase types.GovernancePhase) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyGovernancePhase, []byte{byte(phase)})
}

// GetGovernancePhase retrieves the current governance phase
func (k Keeper) GetGovernancePhase(ctx sdk.Context) types.GovernancePhase {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyGovernancePhase)
	if len(bz) == 0 {
		return types.GovernancePhase_FOUNDER_CONTROL
	}
	return types.GovernancePhase(bz[0])
}

// SetGenesisTime stores the genesis time
func (k Keeper) SetGenesisTime(ctx sdk.Context, time time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&time)
	store.Set(types.KeyGenesisTime, bz)
}

// GetGenesisTime retrieves the genesis time
func (k Keeper) GetGenesisTime(ctx sdk.Context) time.Time {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyGenesisTime)
	if bz == nil {
		return time.Time{}
	}
	
	var genesisTime time.Time
	k.cdc.MustUnmarshal(bz, &genesisTime)
	return genesisTime
}

// SetProtectedParameter stores a protected parameter
func (k Keeper) SetProtectedParameter(ctx sdk.Context, param types.ProtectedParameter) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyProtectedParams, []byte(param.Name)...)
	bz := k.cdc.MustMarshal(&param)
	store.Set(key, bz)
}

// GetProtectedParameter retrieves a protected parameter
func (k Keeper) GetProtectedParameter(ctx sdk.Context, name string) (types.ProtectedParameter, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyProtectedParams, []byte(name)...)
	bz := store.Get(key)
	if bz == nil {
		return types.ProtectedParameter{}, false
	}
	
	var param types.ProtectedParameter
	k.cdc.MustUnmarshal(bz, &param)
	return param, true
}

// GetAllProtectedParameters retrieves all protected parameters
func (k Keeper) GetAllProtectedParameters(ctx sdk.Context) []types.ProtectedParameter {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyProtectedParams)
	defer iterator.Close()
	
	var params []types.ProtectedParameter
	for ; iterator.Valid(); iterator.Next() {
		var param types.ProtectedParameter
		k.cdc.MustUnmarshal(iterator.Value(), &param)
		params = append(params, param)
	}
	
	return params
}

// SetProposalVetoed marks a proposal as vetoed
func (k Keeper) SetProposalVetoed(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetProposalVetoKey(proposalID), []byte{1})
}

// IsProposalVetoed checks if a proposal is vetoed
func (k Keeper) IsProposalVetoed(ctx sdk.Context, proposalID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetProposalVetoKey(proposalID))
}

// GetAllVetoedProposals retrieves all vetoed proposal IDs
func (k Keeper) GetAllVetoedProposals(ctx sdk.Context) []uint64 {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyProposalVeto)
	defer iterator.Close()
	
	var proposalIDs []uint64
	for ; iterator.Valid(); iterator.Next() {
		// Extract proposal ID from key
		key := iterator.Key()
		proposalID := sdk.BigEndianToUint64(key[len(types.KeyProposalVeto):])
		proposalIDs = append(proposalIDs, proposalID)
	}
	
	return proposalIDs
}