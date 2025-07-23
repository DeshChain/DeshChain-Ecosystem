package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/deshchain/deshchain/x/revenue/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace
	
	// Dependencies
	bankKeeper     types.BankKeeper
	donationKeeper types.DonationKeeper
	
	// Revenue recipients (set in genesis)
	developmentFund   string
	communityTreasury string
	liquidityPool     string
	emergencyReserve  string
	founderRoyalty    string
	validatorPool     string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	donationKeeper types.DonationKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     ps,
		bankKeeper:     bankKeeper,
		donationKeeper: donationKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CollectRevenue collects platform revenue from various sources
func (k Keeper) CollectRevenue(ctx sdk.Context, source string, amount sdk.Coins) error {
	// Send revenue to module account
	moduleAddr := k.GetRevenueModuleAccount(ctx)
	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx, source, types.ModuleName, amount,
	); err != nil {
		return err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRevenueCollected,
			sdk.NewAttribute(types.AttributeKeySource, source),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)
	
	// Distribute revenue immediately
	return k.DistributeRevenue(ctx, amount)
}

// DistributeRevenue distributes collected revenue according to platform distribution
func (k Keeper) DistributeRevenue(ctx sdk.Context, revenue sdk.Coins) error {
	// Platform revenue distribution percentages
	distribution := map[string]sdk.Dec{
		"development":    sdk.MustNewDecFromStr("0.25"), // 25%
		"community":      sdk.MustNewDecFromStr("0.20"), // 20%
		"validators":     sdk.MustNewDecFromStr("0.15"), // 15%
		"liquidity":      sdk.MustNewDecFromStr("0.15"), // 15%
		"operations":     sdk.MustNewDecFromStr("0.10"), // 10%
		"emergency":      sdk.MustNewDecFromStr("0.10"), // 10%
		"founder":        sdk.MustNewDecFromStr("0.05"), // 5%
	}
	
	for recipient, percentage := range distribution {
		recipientAmount := sdk.Coins{}
		for _, coin := range revenue {
			amt := coin.Amount.ToDec().Mul(percentage).TruncateInt()
			if amt.IsPositive() {
				recipientAmount = recipientAmount.Add(sdk.NewCoin(coin.Denom, amt))
			}
		}
		
		if !recipientAmount.IsZero() {
			var recipientAddr sdk.AccAddress
			var err error
			
			switch recipient {
			case "development":
				if k.developmentFund != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.developmentFund)
				}
			case "community":
				if k.communityTreasury != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.communityTreasury)
				}
			case "validators":
				// Send to validator rewards pool
				if k.validatorPool != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.validatorPool)
				}
			case "liquidity":
				if k.liquidityPool != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.liquidityPool)
				}
			case "operations":
				// Operations fund for platform expenses
				// This would typically go to a multisig wallet
				continue
			case "emergency":
				if k.emergencyReserve != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.emergencyReserve)
				}
			case "founder":
				if k.founderRoyalty != "" {
					recipientAddr, err = sdk.AccAddressFromBech32(k.founderRoyalty)
				}
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
	
	// Emit distribution event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRevenueDistributed,
			sdk.NewAttribute(types.AttributeKeyAmount, revenue.String()),
		),
	)
	
	return nil
}

// ProcessGrossProfit handles gross profit and sends 30% to NGOs
func (k Keeper) ProcessGrossProfit(ctx sdk.Context, grossProfit sdk.Coins) error {
	// Calculate 30% for NGO donations
	ngoAmount := sdk.Coins{}
	for _, coin := range grossProfit {
		amt := coin.Amount.ToDec().Mul(sdk.MustNewDecFromStr("0.30")).TruncateInt()
		if amt.IsPositive() {
			ngoAmount = ngoAmount.Add(sdk.NewCoin(coin.Denom, amt))
		}
	}
	
	// Send to donation module for NGO distribution
	if !ngoAmount.IsZero() && k.donationKeeper != nil {
		if err := k.donationKeeper.DistributeToNGOs(ctx, ngoAmount); err != nil {
			return err
		}
	}
	
	// Remaining 70% is net profit for operations
	return nil
}

// GetRevenueModuleAccount returns the revenue module account
func (k Keeper) GetRevenueModuleAccount(ctx sdk.Context) sdk.AccAddress {
	return k.bankKeeper.GetModuleAddress(types.ModuleName)
}

// SetRecipientAddresses sets the recipient addresses for revenue distribution
func (k Keeper) SetRecipientAddresses(
	development, community, liquidity, emergency, founder, validators string,
) {
	k.developmentFund = development
	k.communityTreasury = community
	k.liquidityPool = liquidity
	k.emergencyReserve = emergency
	k.founderRoyalty = founder
	k.validatorPool = validators
}