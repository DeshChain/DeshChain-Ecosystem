package revenue

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/revenue/keeper"
	"github.com/deshchain/deshchain/x/revenue/types"
)

// InitGenesis initializes the revenue module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)
	
	// Set recipient addresses
	k.SetRecipientAddresses(
		genState.RecipientAddresses.DevelopmentFund,
		genState.RecipientAddresses.CommunityTreasury,
		genState.RecipientAddresses.LiquidityPool,
		genState.RecipientAddresses.EmergencyReserve,
		genState.RecipientAddresses.FounderRoyalty,
		genState.RecipientAddresses.ValidatorPool,
	)
}

// ExportGenesis returns the revenue module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	
	// Note: We would need to add getter methods for recipient addresses
	// For now, using default genesis addresses
	
	return genesis
}