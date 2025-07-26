package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the identity module
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the identity module's name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the identity module's types on the given LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the identity module
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the identity module
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return data.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the identity module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// TODO: Register query handlers after protobuf generation
	// types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the identity module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	// TODO: Implement CLI commands
	return nil
}

// GetQueryCmd returns no root query command for the identity module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// TODO: Implement CLI queries
	return nil
}

// AppModule implements an application module for the identity module
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// Name returns the identity module's name
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the identity module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO: Register invariants
}

// Route returns the message routing key for the identity module
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the identity module's querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the identity module sdk.Querier
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers module services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// TODO: Register Msg and Query services after protobuf generation
	// types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	// types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the identity module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)
	InitGenesis(ctx, am.keeper, genState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the identity module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock returns the begin blocker for the identity module
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the identity module
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, req, am.keeper)
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the identity module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	// TODO: Implement simulation
}

// ProposalContents doesn't return any content functions for governance proposals
// func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
// 	return nil
// }

// RandomizedParams creates randomized identity param changes for the simulator
// func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
// 	// TODO: Implement simulation param changes
// 	return nil
// }

// RegisterStoreDecoder registers a decoder for identity module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	// TODO: Implement store decoder
}

// WeightedOperations returns the all the identity module operations with their respective weights
// func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
// 	// TODO: Implement weighted operations
// 	return nil
// }