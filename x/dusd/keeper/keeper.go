package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/dusd/types"
	oracletypes "github.com/deshchain/deshchain/x/oracle/types"
	treasurytypes "github.com/deshchain/deshchain/x/treasury/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	memKey       storetypes.StoreKey
	authority    string
	
	// Keepers for integration
	bankKeeper     types.BankKeeper
	oracleKeeper   types.OracleKeeper
	treasuryKeeper types.TreasuryKeeper
	accountKeeper  types.AccountKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	authority string,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	treasuryKeeper types.TreasuryKeeper,
	accountKeeper types.AccountKeeper,
) *Keeper {
	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		authority:      authority,
		bankKeeper:     bankKeeper,
		oracleKeeper:   oracleKeeper,
		treasuryKeeper: treasuryKeeper,
		accountKeeper:  accountKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.DUSDParams, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams(), nil
	}
	
	err = k.cdc.Unmarshal(bz, &params)
	return params, err
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.DUSDParams) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetPosition returns a DUSD position by ID
func (k Keeper) GetPosition(ctx sdk.Context, positionID string) (types.DUSDPosition, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PositionKey)
	bz := store.Get([]byte(positionID))
	if bz == nil {
		return types.DUSDPosition{}, false
	}
	
	var position types.DUSDPosition
	k.cdc.MustUnmarshal(bz, &position)
	return position, true
}

// SetPosition stores a DUSD position
func (k Keeper) SetPosition(ctx sdk.Context, position types.DUSDPosition) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PositionKey)
	bz := k.cdc.MustMarshal(&position)
	store.Set([]byte(position.Id), bz)
}

// DeletePosition removes a DUSD position
func (k Keeper) DeletePosition(ctx sdk.Context, positionID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PositionKey)
	store.Delete([]byte(positionID))
}

// GetAllPositions returns all DUSD positions
func (k Keeper) GetAllPositions(ctx sdk.Context) []types.DUSDPosition {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PositionKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	var positions []types.DUSDPosition
	for ; iterator.Valid(); iterator.Next() {
		var position types.DUSDPosition
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		positions = append(positions, position)
	}
	
	return positions
}

// GetPositionsByOwner returns all positions for a specific owner
func (k Keeper) GetPositionsByOwner(ctx sdk.Context, owner string) []types.DUSDPosition {
	allPositions := k.GetAllPositions(ctx)
	var userPositions []types.DUSDPosition
	
	for _, position := range allPositions {
		if position.Owner == owner {
			userPositions = append(userPositions, position)
		}
	}
	
	return userPositions
}

// GetUSDPrice returns current USD price from oracle (same logic as DINR)
func (k Keeper) GetUSDPrice(ctx sdk.Context) (sdk.Dec, error) {
	// Get USD price from oracle network
	priceData, err := k.oracleKeeper.GetPrice(ctx, "USD", "DINR")
	if err != nil {
		return sdk.ZeroDec(), fmt.Errorf("failed to get USD price: %w", err)
	}
	
	return priceData.Price, nil
}

// CalculateHealthFactor calculates position health factor using same logic as DINR
func (k Keeper) CalculateHealthFactor(ctx sdk.Context, position types.DUSDPosition) (sdk.Dec, error) {
	// Get current USD price
	usdPrice, err := k.GetUSDPrice(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Parse liquidation ratio
	liquidationRatio, err := sdk.NewDecFromStr(params.LiquidationRatio)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Calculate collateral value in USD
	collateralValue := position.CollateralAmount.Amount.ToLegacyDec().Mul(usdPrice)
	
	// Calculate debt value in USD (DUSD is already USD-denominated)
	debtValue := position.MintedAmount.Amount.ToLegacyDec()
	
	// Calculate health factor
	healthFactor := types.CalculateHealthFactor(collateralValue, debtValue, liquidationRatio)
	
	return healthFactor, nil
}

// ValidateCollateral validates collateral requirements using same logic as DINR
func (k Keeper) ValidateCollateral(ctx sdk.Context, collateralAmount sdk.Coin, dusdAmount sdk.Coin) error {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	
	// Check if collateral type is accepted
	collateralAccepted := false
	for _, accepted := range params.AcceptedCollateral {
		if collateralAmount.Denom == accepted {
			collateralAccepted = true
			break
		}
	}
	if !collateralAccepted {
		return fmt.Errorf("collateral type %s not accepted", collateralAmount.Denom)
	}
	
	// Get collateral price
	collateralPrice, err := k.oracleKeeper.GetPrice(ctx, collateralAmount.Denom, "USD")
	if err != nil {
		return fmt.Errorf("failed to get collateral price: %w", err)
	}
	
	// Parse min collateral ratio
	minCollateralRatio, err := sdk.NewDecFromStr(params.MinCollateralRatio)
	if err != nil {
		return err
	}
	
	// Calculate required collateral
	requiredCollateral := types.CalculateCollateralRequired(
		dusdAmount.Amount.ToLegacyDec(),
		collateralPrice.Price,
		minCollateralRatio,
	)
	
	// Check if provided collateral is sufficient
	providedCollateral := collateralAmount.Amount.ToLegacyDec()
	if providedCollateral.LT(requiredCollateral) {
		return fmt.Errorf("insufficient collateral: provided %s, required %s", 
			providedCollateral.String(), requiredCollateral.String())
	}
	
	return nil
}

// MintDUSD creates new DUSD tokens using same algorithm as DINR
func (k Keeper) MintDUSD(ctx sdk.Context, creator string, collateralAmount sdk.Coin, dusdAmount sdk.Coin) (*types.DUSDPosition, error) {
	// Validate inputs
	if err := types.ValidateAddress(creator); err != nil {
		return nil, err
	}
	if err := types.ValidatePositiveAmount(collateralAmount); err != nil {
		return nil, err
	}
	if err := types.ValidatePositiveAmount(dusdAmount); err != nil {
		return nil, err
	}
	
	// Validate collateral requirements
	if err := k.ValidateCollateral(ctx, collateralAmount, dusdAmount); err != nil {
		return nil, err
	}
	
	// Calculate fee using same logic as DINR
	fee := types.CalculateFee(dusdAmount)
	
	// Transfer collateral from user to module
	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	
	if err := k.bankKeeper.SendCoins(ctx, creatorAddr, moduleAddr, sdk.NewCoins(collateralAmount)); err != nil {
		return nil, fmt.Errorf("failed to transfer collateral: %w", err)
	}
	
	// Mint DUSD tokens
	dusdCoin := sdk.NewCoins(dusdAmount)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, dusdCoin); err != nil {
		return nil, fmt.Errorf("failed to mint DUSD: %w", err)
	}
	
	// Transfer DUSD to user (minus fee)
	dusdToUser := dusdAmount.Sub(fee)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(dusdToUser)); err != nil {
		return nil, fmt.Errorf("failed to transfer DUSD to user: %w", err)
	}
	
	// Transfer fee to charity pool (same 40% allocation as other fees)
	if err := k.treasuryKeeper.AddRevenue(ctx, treasurytypes.RevenueSourceDUSDFees, sdk.NewCoins(fee)); err != nil {
		k.Logger(ctx).Error("failed to add DUSD fee to treasury", "error", err)
		// Don't fail the transaction for fee accounting issues
	}
	
	// Create position
	positionID := types.GeneratePositionID(creator, ctx.BlockTime())
	position := types.DUSDPosition{
		Id:               positionID,
		Owner:            creator,
		MintedAmount:     dusdAmount,
		CollateralAmount: collateralAmount,
		CollateralType:   collateralAmount.Denom,
		LastUpdate:       ctx.BlockTime(),
		IsLiquidatable:   false,
	}
	
	// Calculate and set health factor
	healthFactor, err := k.CalculateHealthFactor(ctx, position)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate health factor: %w", err)
	}
	position.HealthFactor = healthFactor.String()
	position.IsLiquidatable = types.IsPositionLiquidatable(healthFactor)
	
	// Store position
	k.SetPosition(ctx, position)
	
	// Update supply statistics
	k.UpdateSupplyStats(ctx)
	
	return &position, nil
}

// UpdateSupplyStats updates total supply statistics
func (k Keeper) UpdateSupplyStats(ctx sdk.Context) {
	totalSupply := k.bankKeeper.GetSupply(ctx, types.DUSDDenom)
	
	store := ctx.KVStore(k.storeKey)
	supplyStats := types.ReserveStats{
		TotalSupply: totalSupply,
		LastUpdated: ctx.BlockTime(),
	}
	
	bz := k.cdc.MustMarshal(&supplyStats)
	store.Set(types.SupplyStatsKey, bz)
}