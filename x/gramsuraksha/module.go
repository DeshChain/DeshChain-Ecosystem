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

package grampension

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the gram pension module.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

// Name returns the gram pension module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the gram pension module's types on the LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes for the gram pension module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the gram pension module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return ValidateGenesis(&genState)
}

// RegisterRESTRoutes registers the REST routes for the gram pension module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	// REST routes will be implemented when needed
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the gram pension module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// gRPC gateway routes will be implemented when needed
}

// GetTxCmd returns the transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	// CLI commands will be implemented when needed
	return nil
}

// GetQueryCmd returns the cli query commands for this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// CLI commands will be implemented when needed
	return nil
}

// AppModule implements the AppModule interface for the gram pension module.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.BinaryCodec, keeper keeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// Name returns the gram pension module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// Service registration will be implemented when needed
}

// RegisterInvariants registers the gram pension module's invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// Invariants will be implemented when needed
}

// Route returns the message routing key for the gram pension module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the gram pension module's query routing key.
func (AppModule) QuerierRoute() string { return types.QuerierRoute }

// LegacyQuerierHandler returns the gram pension module's Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil // Will implement when needed
}

// InitGenesis performs the gram pension module's genesis initialization
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, &genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the gram pension module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock executes all ABCI BeginBlock logic respective to the gram pension module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	// Check for defaulted participants
	am.keeper.CheckParticipantDefaults(ctx)
	
	// Process any matured pensions
	am.keeper.ProcessAllMaturities(ctx)
	
	// Send notifications for upcoming maturities
	am.keeper.NotifyUpcomingMaturities(ctx)
}

// EndBlock executes all ABCI EndBlock logic respective to the gram pension module.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Monthly revenue distribution (check if it's the first day of the month)
	if ctx.BlockTime().Day() == 1 {
		am.keeper.CallHookMonthlyDistribution(ctx)
	}
	
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState creates a randomized GenState of the gram pension module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	// Simulation will be implemented when needed
}

// ProposalContents returns all the gram pension content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  gram pension param changes for the simulator.
func (AppModule) RandomizedParams(_ *simtypes.Rand) []simtypes.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for gram pension module's types
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gram pension module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return nil
}