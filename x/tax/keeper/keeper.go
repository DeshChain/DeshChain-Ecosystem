package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/deshchain/deshchain/x/tax/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	paramSpace    paramtypes.Subspace
	
	// Dependencies
	bankKeeper    types.BankKeeper
	revenueKeeper types.RevenueKeeper
	
	// Tax recipients
	ngoWallet         string
	validatorPool     string
	communityPool     string
	techInnovation    string
	operations        string
	talentAcquisition string
	strategicReserve  string
	founderWallet     string
	coFoundersWallet  string
	angelWallet       string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	revenueKeeper types.RevenueKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    ps,
		bankKeeper:    bankKeeper,
		revenueKeeper: revenueKeeper,
		// Initialize wallet addresses from genesis
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CollectTax deducts tax from a transaction and distributes it
func (k Keeper) CollectTax(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) error {
	params := k.GetParams(ctx)
	if !params.Enabled {
		return nil
	}
	
	// Calculate tax amount (2.5% base rate)
	taxAmount := sdk.Coins{}
	for _, coin := range amount {
		taxCoin := sdk.NewCoin(
			coin.Denom,
			coin.Amount.ToDec().Mul(params.TaxRate).TruncateInt(),
		)
		if taxCoin.Amount.IsPositive() {
			taxAmount = taxAmount.Add(taxCoin)
		}
	}
	
	if taxAmount.IsZero() {
		return nil
	}
	
	// Deduct tax from sender
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, from, types.ModuleName, taxAmount,
	); err != nil {
		return err
	}
	
	// Distribute tax according to distribution percentages
	return k.DistributeTax(ctx, taxAmount)
}

// DistributeTax distributes collected tax to various recipients
func (k Keeper) DistributeTax(ctx sdk.Context, taxAmount sdk.Coins) error {
	moduleAddr := k.GetTaxModuleAccount(ctx)
	
	// Distribution percentages from types/distribution.go
	distribution := types.GetTaxDistribution()
	
	for recipient, percentage := range distribution {
		recipientAmount := sdk.Coins{}
		for _, coin := range taxAmount {
			amt := coin.Amount.ToDec().Mul(percentage).TruncateInt()
			if amt.IsPositive() {
				recipientAmount = recipientAmount.Add(sdk.NewCoin(coin.Denom, amt))
			}
		}
		
		if !recipientAmount.IsZero() {
			var recipientAddr sdk.AccAddress
			var err error
			
			switch recipient {
			case types.RecipientNGO:
				recipientAddr, err = sdk.AccAddressFromBech32(k.ngoWallet)
			case types.RecipientValidators:
				// Send to validator rewards pool
				err = k.SendToValidatorRewards(ctx, recipientAmount)
				continue
			case types.RecipientCommunity:
				// Send to community pool
				err = k.SendToCommunityPool(ctx, recipientAmount)
				continue
			case types.RecipientTechInnovation:
				recipientAddr, err = sdk.AccAddressFromBech32(k.techInnovation)
			case types.RecipientOperations:
				recipientAddr, err = sdk.AccAddressFromBech32(k.operations)
			case types.RecipientTalentAcquisition:
				recipientAddr, err = sdk.AccAddressFromBech32(k.talentAcquisition)
			case types.RecipientStrategicReserve:
				recipientAddr, err = sdk.AccAddressFromBech32(k.strategicReserve)
			case types.RecipientFounder:
				recipientAddr, err = sdk.AccAddressFromBech32(k.founderWallet)
			case types.RecipientCoFounders:
				recipientAddr, err = sdk.AccAddressFromBech32(k.coFoundersWallet)
			case types.RecipientAngelInvestors:
				recipientAddr, err = sdk.AccAddressFromBech32(k.angelWallet)
			}
			
			if err != nil {
				return err
			}
			
			if recipientAddr != nil {
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx, types.ModuleName, recipientAddr, recipientAmount,
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