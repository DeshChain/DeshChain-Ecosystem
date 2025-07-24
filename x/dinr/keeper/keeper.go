package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/deshchain/namo/x/dinr/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	oracleKeeper  types.OracleKeeper // For price feeds
	revenueKeeper types.RevenueKeeper // For fee distribution
	taxKeeper     types.TaxKeeper     // For NAMO fee collection
}

// NewKeeper creates new instances of the DINR Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	revenueKeeper types.RevenueKeeper,
	taxKeeper types.TaxKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		oracleKeeper:  oracleKeeper,
		revenueKeeper: revenueKeeper,
		taxKeeper:     taxKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// MintDINR mints new DINR tokens against deposited collateral
func (k Keeper) MintDINR(ctx sdk.Context, minter sdk.AccAddress, collateral sdk.Coin, dinrToMint sdk.Coin) error {
	params := k.GetParams(ctx)
	
	// Check if minting is enabled
	if !params.MintingEnabled {
		return types.ErrMintingDisabled
	}
	
	// Use CollateralManager for validation
	collateralManager := k.GetCollateralManager()
	
	// Validate collateral asset
	_, err := collateralManager.ValidateCollateral(ctx, collateral.Denom)
	if err != nil {
		return err
	}
	
	// Get existing position
	userPosition, found := k.GetUserPosition(ctx, minter.String())
	totalCollateral := sdk.NewCoins(collateral)
	totalDINRToMint := dinrToMint.Amount
	
	if found {
		// Add existing collateral and debt
		totalCollateral = userPosition.Collateral.Add(collateral)
		totalDINRToMint = totalDINRToMint.Add(userPosition.DinrMinted.Amount)
	}
	
	// Validate collateral sufficiency
	err = collateralManager.ValidateCollateralSufficiency(ctx, totalCollateral, totalDINRToMint)
	if err != nil {
		return err
	}
	
	// Calculate collateral value for ratio calculation
	collateralValue, err := collateralManager.CalculateCollateralValue(ctx, totalCollateral)
	if err != nil {
		return err
	}
	
	collateralRatio := uint64(collateralValue.Mul(sdk.NewInt(10000)).Quo(totalDINRToMint).Int64())
	
	// Calculate minting fee
	fee := k.calculateMintingFee(dinrToMint, params.Fees)
	netDINR := dinrToMint.Sub(fee)
	
	// Transfer collateral from user to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, minter, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}
	
	// Mint DINR tokens
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(netDINR))
	if err != nil {
		return err
	}
	
	// Send minted DINR to user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, minter, sdk.NewCoins(netDINR))
	if err != nil {
		return err
	}
	
	// Distribute fees
	if fee.Amount.GT(sdk.ZeroInt()) {
		err = k.distributeFees(ctx, fee)
		if err != nil {
			return err
		}
	}
	
	// Update user position
	k.updateUserPosition(ctx, minter.String(), collateral, dinrToMint, collateralRatio)
	
	// Update stability metrics
	k.updateStabilityMetrics(ctx)
	
	// Emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintDINR,
			sdk.NewAttribute(types.AttributeKeyMinter, minter.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyDINRMinted, netDINR.String()),
			sdk.NewAttribute(types.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralRatio, fmt.Sprintf("%d", collateralRatio)),
		),
	})
	
	return nil
}

// BurnDINR burns DINR tokens and returns collateral
func (k Keeper) BurnDINR(ctx sdk.Context, burner sdk.AccAddress, dinrToBurn sdk.Coin, collateralDenom string) error {
	params := k.GetParams(ctx)
	
	// Check if burning is enabled
	if !params.BurningEnabled {
		return types.ErrBurningDisabled
	}
	
	// Get user position
	userPosition, found := k.GetUserPosition(ctx, burner.String())
	if !found {
		return types.ErrPositionNotFound
	}
	
	// Validate user has enough DINR minted
	if userPosition.DinrMinted.Amount.LT(dinrToBurn.Amount) {
		return types.ErrInsufficientDINR
	}
	
	// Calculate burning fee
	fee := k.calculateBurningFee(dinrToBurn, params.Fees)
	totalDINRNeeded := dinrToBurn.Add(fee)
	
	// Check user has enough DINR balance
	userDINRBalance := k.bankKeeper.GetBalance(ctx, burner, types.DINRDenom)
	if userDINRBalance.Amount.LT(totalDINRNeeded.Amount) {
		return types.ErrInsufficientBalance
	}
	
	// Calculate collateral to return
	collateralToReturn := k.calculateCollateralToReturn(ctx, userPosition, dinrToBurn, collateralDenom)
	
	// Burn DINR from user account
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, burner, types.ModuleName, sdk.NewCoins(totalDINRNeeded))
	if err != nil {
		return err
	}
	
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(dinrToBurn))
	if err != nil {
		return err
	}
	
	// Return collateral to user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, burner, sdk.NewCoins(collateralToReturn))
	if err != nil {
		return err
	}
	
	// Distribute fees
	if fee.Amount.GT(sdk.ZeroInt()) {
		err = k.distributeFees(ctx, fee)
		if err != nil {
			return err
		}
	}
	
	// Update user position
	k.updateUserPositionAfterBurn(ctx, burner.String(), dinrToBurn, collateralToReturn)
	
	// Update stability metrics
	k.updateStabilityMetrics(ctx)
	
	// Emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnDINR,
			sdk.NewAttribute(types.AttributeKeyBurner, burner.String()),
			sdk.NewAttribute(types.AttributeKeyDINRBurned, dinrToBurn.String()),
			sdk.NewAttribute(types.AttributeKeyCollateralReturned, collateralToReturn.String()),
			sdk.NewAttribute(types.AttributeKeyFee, fee.String()),
		),
	})
	
	return nil
}

// DepositCollateral allows users to deposit additional collateral
func (k Keeper) DepositCollateral(ctx sdk.Context, depositor sdk.AccAddress, collateral sdk.Coin) error {
	// Get user position
	userPosition, found := k.GetUserPosition(ctx, depositor.String())
	if !found {
		return types.ErrPositionNotFound
	}

	// Validate collateral is acceptable
	collateralAsset, found := k.GetCollateralAsset(ctx, collateral.Denom)
	if !found || !collateralAsset.IsEnabled {
		return types.ErrInvalidCollateral
	}

	// Transfer collateral from user to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	// Update user position
	userPosition.Collateral = userPosition.Collateral.Add(collateral)
	
	// Recalculate health factor
	totalCollateralValue := k.calculateTotalCollateralValue(ctx, userPosition.Collateral)
	collateralRatio := k.calculateCollateralRatio(totalCollateralValue, userPosition.DinrMinted.Amount)
	userPosition.HealthFactor = k.calculateHealthFactor(collateralRatio)

	k.SetUserPosition(ctx, userPosition)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositCollateral,
			sdk.NewAttribute(types.AttributeKeyDepositor, depositor.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyHealthFactor, userPosition.HealthFactor),
		),
	)

	return nil
}

// WithdrawCollateral allows users to withdraw excess collateral
func (k Keeper) WithdrawCollateral(ctx sdk.Context, withdrawer sdk.AccAddress, collateral sdk.Coin) error {
	params := k.GetParams(ctx)
	
	// Get user position
	userPosition, found := k.GetUserPosition(ctx, withdrawer.String())
	if !found {
		return types.ErrPositionNotFound
	}

	// Check if user has enough collateral of this type
	userCollateralAmount := userPosition.Collateral.AmountOf(collateral.Denom)
	if userCollateralAmount.LT(collateral.Amount) {
		return types.ErrInsufficientCollateral
	}

	// Calculate new collateral value after withdrawal
	newCollateral := userPosition.Collateral.Sub(sdk.NewCoins(collateral))
	newCollateralValue := k.calculateTotalCollateralValue(ctx, newCollateral)
	
	// Check if withdrawal maintains minimum collateral ratio
	newCollateralRatio := k.calculateCollateralRatio(newCollateralValue, userPosition.DinrMinted.Amount)
	if newCollateralRatio < uint64(params.MinCollateralRatio) {
		return types.ErrInsufficientCollateral
	}

	// Transfer collateral from module to user
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawer, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	// Update user position
	userPosition.Collateral = newCollateral
	userPosition.HealthFactor = k.calculateHealthFactor(newCollateralRatio)
	k.SetUserPosition(ctx, userPosition)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCollateral,
			sdk.NewAttribute(types.AttributeKeyWithdrawer, withdrawer.String()),
			sdk.NewAttribute(types.AttributeKeyCollateral, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyHealthFactor, userPosition.HealthFactor),
		),
	)

	return nil
}

// Liquidate allows liquidators to liquidate undercollateralized positions
func (k Keeper) Liquidate(ctx sdk.Context, liquidator sdk.AccAddress, user sdk.AccAddress, dinrToCover sdk.Coin) (sdk.Coins, error) {
	// Use the new CollateralManager for enhanced liquidation logic
	collateralManager := k.GetCollateralManager()
	
	err := collateralManager.ProcessLiquidation(ctx, liquidator, user.String(), dinrToCover.Amount)
	if err != nil {
		return sdk.Coins{}, err
	}
	
	// Update stability metrics
	k.updateStabilityMetrics(ctx)
	
	// Return the collateral that was seized (for compatibility)
	// In practice, the collateral transfer is handled in ProcessLiquidation
	return sdk.Coins{dinrToCover}, nil
}

// Helper functions
func (k Keeper) calculateCollateralValue(collateral sdk.Coin, price sdk.Dec) sdk.Int {
	// Implementation would convert collateral amount to INR value using oracle price
	return collateral.Amount.Mul(price.TruncateInt())
}

func (k Keeper) calculateCollateralRatio(collateralValue, dinrAmount sdk.Int) uint64 {
	// Returns ratio in basis points (e.g., 15000 = 150%)
	if dinrAmount.IsZero() {
		return 0
	}
	ratio := collateralValue.Mul(sdk.NewInt(10000)).Quo(dinrAmount)
	return ratio.Uint64()
}

func (k Keeper) calculateMintingFee(dinr sdk.Coin, fees types.FeeStructure) sdk.Coin {
	// Calculate fee with cap
	feeAmount := dinr.Amount.Mul(sdk.NewInt(int64(fees.MintFee))).Quo(sdk.NewInt(10000))
	
	// Apply cap
	capAmount, _ := sdk.NewIntFromString(fees.MintFeeCap)
	if feeAmount.GT(capAmount) {
		feeAmount = capAmount
	}
	
	return sdk.NewCoin(types.DINRDenom, feeAmount)
}

func (k Keeper) calculateBurningFee(dinr sdk.Coin, fees types.FeeStructure) sdk.Coin {
	// Calculate fee with cap
	feeAmount := dinr.Amount.Mul(sdk.NewInt(int64(fees.BurnFee))).Quo(sdk.NewInt(10000))
	
	// Apply cap
	capAmount, _ := sdk.NewIntFromString(fees.BurnFeeCap)
	if feeAmount.GT(capAmount) {
		feeAmount = capAmount
	}
	
	return sdk.NewCoin(types.DINRDenom, feeAmount)
}

func (k Keeper) distributeFees(ctx sdk.Context, fee sdk.Coin) error {
	// Distribute fees according to platform revenue distribution model
	// This would integrate with the revenue module
	return k.revenueKeeper.DistributePlatformRevenue(ctx, sdk.NewCoins(fee))
}

// Additional helper functions
func (k Keeper) calculateTotalCollateralValue(ctx sdk.Context, collateral sdk.Coins) sdk.Int {
	totalValue := sdk.ZeroInt()
	
	for _, coin := range collateral {
		// Get price from oracle
		price, err := k.oracleKeeper.GetPrice(ctx, coin.Denom)
		if err != nil {
			continue // Skip if no price available
		}
		
		// Convert to INR value
		collateralValue := k.calculateCollateralValue(coin, price)
		totalValue = totalValue.Add(collateralValue)
	}
	
	return totalValue
}

func (k Keeper) calculateHealthFactor(collateralRatio uint64) string {
	// Health factor = collateral ratio / liquidation threshold
	// > 1.0 = healthy, < 1.0 = liquidatable
	params := k.GetParams(sdk.UnwrapSDKContext(nil))
	healthFactor := sdk.NewDec(int64(collateralRatio)).Quo(sdk.NewDec(int64(params.LiquidationThreshold)))
	return healthFactor.String()
}

func (k Keeper) calculateCollateralToReturn(ctx sdk.Context, userPosition types.UserPosition, dinrToBurn sdk.Coin, collateralDenom string) sdk.Coin {
	// Calculate proportional collateral to return
	burnRatio := dinrToBurn.Amount.Mul(sdk.NewInt(10000)).Quo(userPosition.DinrMinted.Amount)
	
	// Find the requested collateral
	for _, coin := range userPosition.Collateral {
		if coin.Denom == collateralDenom {
			returnAmount := coin.Amount.Mul(burnRatio).Quo(sdk.NewInt(10000))
			return sdk.NewCoin(collateralDenom, returnAmount)
		}
	}
	
	return sdk.NewCoin(collateralDenom, sdk.ZeroInt())
}

func (k Keeper) selectCollateralForLiquidation(ctx sdk.Context, collateral sdk.Coins, targetValue sdk.Int) sdk.Coins {
	// Select collateral proportionally to meet target value
	totalValue := k.calculateTotalCollateralValue(ctx, collateral)
	
	if totalValue.IsZero() {
		return sdk.NewCoins()
	}
	
	selectedCollateral := sdk.NewCoins()
	ratio := targetValue.Mul(sdk.NewInt(10000)).Quo(totalValue)
	
	for _, coin := range collateral {
		selectedAmount := coin.Amount.Mul(ratio).Quo(sdk.NewInt(10000))
		if selectedAmount.GT(sdk.ZeroInt()) {
			selectedCollateral = selectedCollateral.Add(sdk.NewCoin(coin.Denom, selectedAmount))
		}
	}
	
	return selectedCollateral
}

func (k Keeper) updateUserPosition(ctx sdk.Context, address string, collateral sdk.Coin, dinrMinted sdk.Coin, collateralRatio uint64) {
	position, found := k.GetUserPosition(ctx, address)
	
	if !found {
		position = types.UserPosition{
			Address:      address,
			Collateral:   sdk.NewCoins(collateral),
			DinrMinted:   dinrMinted,
			HealthFactor: k.calculateHealthFactor(collateralRatio),
			LastUpdate:   ctx.BlockTime(),
		}
	} else {
		position.Collateral = position.Collateral.Add(collateral)
		position.DinrMinted = position.DinrMinted.Add(dinrMinted)
		position.HealthFactor = k.calculateHealthFactor(collateralRatio)
		position.LastUpdate = ctx.BlockTime()
	}
	
	k.SetUserPosition(ctx, position)
}

func (k Keeper) updateUserPositionAfterBurn(ctx sdk.Context, address string, dinrBurned sdk.Coin, collateralReturned sdk.Coin) {
	position, found := k.GetUserPosition(ctx, address)
	if !found {
		return
	}
	
	// Update position
	position.DinrMinted = position.DinrMinted.Sub(dinrBurned)
	position.Collateral = position.Collateral.Sub(sdk.NewCoins(collateralReturned))
	
	// If no more DINR minted, remove position
	if position.DinrMinted.IsZero() {
		k.RemoveUserPosition(ctx, address)
		return
	}
	
	// Recalculate health factor
	totalCollateralValue := k.calculateTotalCollateralValue(ctx, position.Collateral)
	collateralRatio := k.calculateCollateralRatio(totalCollateralValue, position.DinrMinted.Amount)
	position.HealthFactor = k.calculateHealthFactor(collateralRatio)
	position.LastUpdate = ctx.BlockTime()
	
	k.SetUserPosition(ctx, position)
}

func (k Keeper) updateStabilityMetrics(ctx sdk.Context) {
	// This function would update global stability metrics
	// Implementation would track total supply, collateral value, etc.
	// For now, it's a placeholder that would be implemented with oracle integration
}

// Oracle integration methods
func (k Keeper) GetCurrentPrice(ctx sdk.Context, symbol string) (sdk.Dec, error) {
	// Get price from oracle module
	price, err := k.oracleKeeper.GetPrice(ctx, symbol)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	return price, nil
}

// Daily limits tracking
func (k Keeper) GetDailyMintedAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DailyMintedAmountKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	
	amount, ok := sdk.NewIntFromString(string(bz))
	if !ok {
		return sdk.ZeroInt()
	}
	
	return amount
}

func (k Keeper) SetDailyMintedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DailyMintedAmountKey, []byte(amount.String()))
}

func (k Keeper) GetDailyBurnedAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DailyBurnedAmountKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	
	amount, ok := sdk.NewIntFromString(string(bz))
	if !ok {
		return sdk.ZeroInt()
	}
	
	return amount
}

func (k Keeper) SetDailyBurnedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DailyBurnedAmountKey, []byte(amount.String()))
}

// Reset daily amounts at the beginning of each day
func (k Keeper) ResetDailyAmounts(ctx sdk.Context) {
	k.SetDailyMintedAmount(ctx, sdk.ZeroInt())
	k.SetDailyBurnedAmount(ctx, sdk.ZeroInt())
}

// GetStabilityController returns a new stability controller instance
func (k Keeper) GetStabilityController() *StabilityController {
	return NewStabilityController(k)
}

// GetCollateralManager returns a new collateral manager instance
func (k Keeper) GetCollateralManager() *CollateralManager {
	return NewCollateralManager(k)
}