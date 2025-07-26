package tax

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

// InitGenesis initializes the tax module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)
	
	// Set recipient addresses
	k.SetRecipientAddresses(
		genState.RecipientAddresses.NgoWallet,
		genState.RecipientAddresses.ValidatorPool,
		genState.RecipientAddresses.CommunityPool,
		genState.RecipientAddresses.TechInnovation,
		genState.RecipientAddresses.Operations,
		genState.RecipientAddresses.TalentAcquisition,
		genState.RecipientAddresses.StrategicReserve,
		genState.RecipientAddresses.FounderWallet,
		genState.RecipientAddresses.CoFoundersWallet,
		genState.RecipientAddresses.AngelWallet,
	)
}

// ExportGenesis returns the tax module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	
	// Note: We would need to add getter methods for recipient addresses
	// For now, using default genesis addresses
	
	return genesis
}