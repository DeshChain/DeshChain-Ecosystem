package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/deshchain/namo/x/tax/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	paramSpace    paramtypes.Subspace
	
	// Dependencies
	bankKeeper    types.BankKeeper
	revenueKeeper types.RevenueKeeper
	dexKeeper     types.DEXKeeper
	oracleKeeper  types.OracleKeeper
	
	// Managers
	namoBurnManager *NAMOBurnManager
	namoSwapRouter  *NAMOSwapRouter
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	revenueKeeper types.RevenueKeeper,
	dexKeeper types.DEXKeeper,
	oracleKeeper types.OracleKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    ps,
		bankKeeper:    bankKeeper,
		revenueKeeper: revenueKeeper,
		dexKeeper:     dexKeeper,
		oracleKeeper:  oracleKeeper,
	}
	
	// Initialize managers
	k.namoBurnManager = NewNAMOBurnManager(k)
	k.namoSwapRouter = NewNAMOSwapRouter(k)
	
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CollectTax deducts tax from a transaction and distributes it
func (k Keeper) CollectTax(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coin, msgType string) error {
	// Create tax calculator with current config
	taxCalculator := types.NewTaxCalculator(k.GetTaxConfig(ctx))
	
	// Calculate tax using progressive structure
	taxResult, err := taxCalculator.CalculateTax(amount, msgType)
	if err != nil {
		return err
	}
	
	if taxResult.TaxAmount.IsZero() {
		return nil
	}
	
	// Handle NAMO payment for tax
	finalTaxAmount := taxResult.TaxAmount
	if amount.Denom != "namo" {
		// User paying in different token - swap for NAMO
		namoTax, err := k.namoSwapRouter.SwapForNAMOFee(ctx, from, taxResult.TaxAmount, amount, false)
		if err != nil {
			return fmt.Errorf("failed to swap for NAMO tax: %w", err)
		}
		finalTaxAmount = namoTax
	}
	
	// Deduct tax from sender
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, from, types.ModuleName, sdk.NewCoins(finalTaxAmount),
	); err != nil {
		return err
	}
	
	// Emit tax collected event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTaxCollected,
			sdk.NewAttribute(types.AttributeKeyAmount, finalTaxAmount.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(types.AttributeKeyTaxRate, taxResult.EffectiveRate),
		),
	)
	
	// Distribute tax according to distribution percentages
	return k.DistributeTax(ctx, finalTaxAmount)
}

// DistributeTax distributes collected tax to various recipients
func (k Keeper) DistributeTax(ctx sdk.Context, taxAmount sdk.Coin) error {
	// Get tax distribution configuration
	taxDist := types.NewDefaultTaxDistribution()
	
	// Calculate amounts for each recipient
	amounts := taxDist.CalculateTaxAmounts(taxAmount)
	
	// Distribute to each recipient
	for poolName, amount := range amounts {
		if amount.IsZero() {
			continue
		}
		
		// Handle NAMO burn separately
		if poolName == types.NAMOBurnPoolName {
			err := k.namoBurnManager.BurnFromDistribution(ctx, amounts)
			if err != nil {
				return fmt.Errorf("failed to burn NAMO: %w", err)
			}
			continue
		}
		
		// Send to module accounts
		err := k.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			types.ModuleName,
			poolName,
			sdk.NewCoins(amount),
		)
		if err != nil {
			return fmt.Errorf("failed to distribute to %s: %w", poolName, err)
		}
	}
	
	// Emit distribution event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTaxDistributed,
			sdk.NewAttribute(types.AttributeKeyAmount, taxAmount.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				); err != nil {
					return err
				}
			}
		}
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTaxCollected,
			sdk.NewAttribute(types.AttributeKeyAmount, taxAmount.String()),
		),
	)
	
	return nil
}

// SendToValidatorRewards sends coins to validator rewards pool
func (k Keeper) SendToValidatorRewards(ctx sdk.Context, amount sdk.Coins) error {
	// This would integrate with the distribution module
	// For now, send to a designated validator rewards address
	if k.validatorPool != "" {
		valAddr, err := sdk.AccAddressFromBech32(k.validatorPool)
		if err != nil {
			return err
		}
		return k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, valAddr, amount,
		)
	}
	return nil
}

// SendToCommunityPool sends coins to community pool
func (k Keeper) SendToCommunityPool(ctx sdk.Context, amount sdk.Coins) error {
	// This would integrate with the distribution module
	// For now, send to a designated community pool address
	if k.communityPool != "" {
		commAddr, err := sdk.AccAddressFromBech32(k.communityPool)
		if err != nil {
			return err
		}
		return k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, commAddr, amount,
		)
	}
	return nil
}

// GetTaxModuleAccount returns the tax module account
func (k Keeper) GetTaxModuleAccount(ctx sdk.Context) sdk.AccAddress {
	return k.bankKeeper.GetModuleAddress(types.ModuleName)
}

// SetRecipientAddresses sets the recipient addresses for tax distribution
func (k Keeper) SetRecipientAddresses(
	ngo, validators, community, tech, ops, talent, 
	strategic, founder, coFounders, angels string,
) {
	k.ngoWallet = ngo
	k.validatorPool = validators
	k.communityPool = community
	k.techInnovation = tech
	k.operations = ops
	k.talentAcquisition = talent
	k.strategicReserve = strategic
	k.founderWallet = founder
	k.coFoundersWallet = coFounders
	k.angelWallet = angels
}

// GetTaxConfig returns the current tax configuration
func (k Keeper) GetTaxConfig(ctx sdk.Context) *types.TaxConfig {
	params := k.GetParams(ctx)
	return &types.TaxConfig{
		Enabled:                  params.Enabled,
		ExemptMessages:          params.ExemptMessages,
		ExemptAddresses:         params.ExemptAddresses,
		OptimizationEnabled:     true,
		DefaultOptimizationMode: "automatic",
	}
}

// GetModuleAddress returns the module address for a given module name
func (k Keeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return k.bankKeeper.GetModuleAddress(moduleName)
}

// DistributePlatformRevenue distributes platform revenue according to the distribution model
func (k Keeper) DistributePlatformRevenue(ctx sdk.Context, revenueSource string, revenue sdk.Coin) error {
	// Get platform distribution configuration
	platformDist := types.NewDefaultPlatformDistribution()
	
	// Calculate amounts for each recipient
	amounts := platformDist.CalculatePlatformAmounts(revenue)
	
	// Distribute to each recipient
	for poolName, amount := range amounts {
		if amount.IsZero() {
			continue
		}
		
		// Handle NAMO burn separately
		if poolName == types.NAMOBurnPoolName {
			err := k.namoBurnManager.BurnFromPlatformRevenue(ctx, revenueSource, revenue)
			if err != nil {
				return fmt.Errorf("failed to burn NAMO: %w", err)
			}
			continue
		}
		
		// Send to module accounts
		err := k.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			types.ModuleName,
			poolName,
			sdk.NewCoins(amount),
		)
		if err != nil {
			return fmt.Errorf("failed to distribute to %s: %w", poolName, err)
		}
	}
	
	return nil
}

// GetNAMOBurnManager returns the NAMO burn manager
func (k Keeper) GetNAMOBurnManager() *NAMOBurnManager {
	return k.namoBurnManager
}

// GetNAMOSwapRouter returns the NAMO swap router
func (k Keeper) GetNAMOSwapRouter() *NAMOSwapRouter {
	return k.namoSwapRouter
}
}