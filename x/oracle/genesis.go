package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/oracle/keeper"
	"github.com/deshchain/deshchain/x/oracle/types"
)

// InitGenesis initializes the oracle module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set parameters
	k.SetParams(ctx, genState.Params)

	// Initialize oracle validators
	for _, validator := range genState.OracleValidators {
		k.SetOracleValidator(ctx, validator)
	}

	// Initialize price data
	for _, priceData := range genState.PriceData {
		k.SetPriceData(ctx, priceData)
	}

	// Initialize exchange rates
	for _, exchangeRate := range genState.ExchangeRates {
		store := ctx.KVStore(k.StoreKey())
		bz := k.Codec().MustMarshal(&exchangeRate)
		store.Set(types.ExchangeRateKey(exchangeRate.Base, exchangeRate.Target), bz)
	}
}

// ExportGenesis returns the oracle module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genesis := types.DefaultGenesis()
	
	// Export parameters
	genesis.Params = k.GetParams(ctx)

	// Export oracle validators
	genesis.OracleValidators = k.GetAllOracleValidators(ctx)

	// Export price data
	genesis.PriceData = k.GetAllPriceData(ctx)

	// Export exchange rates
	genesis.ExchangeRates = k.GetAllExchangeRates(ctx)

	return *genesis
}