package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/governance/types"
)

// CharityAllocationParams defines the graduated charity allocation structure
type CharityAllocationParams struct {
	Year1Percentage  sdk.Dec // 10%
	Year2Percentage  sdk.Dec // 20%
	Year3Percentage  sdk.Dec // 30%
	Year4Percentage  sdk.Dec // 35%
	Year5Percentage  sdk.Dec // 40%
	GenesisTime      time.Time
}

// DefaultCharityAllocationParams returns the default graduated charity allocation
func DefaultCharityAllocationParams(genesisTime time.Time) CharityAllocationParams {
	return CharityAllocationParams{
		Year1Percentage: sdk.NewDecWithPrec(10, 2), // 10%
		Year2Percentage: sdk.NewDecWithPrec(20, 2), // 20%
		Year3Percentage: sdk.NewDecWithPrec(30, 2), // 30%
		Year4Percentage: sdk.NewDecWithPrec(35, 2), // 35%
		Year5Percentage: sdk.NewDecWithPrec(40, 2), // 40%
		GenesisTime:     genesisTime,
	}
}

// GetCurrentCharityPercentage returns the charity percentage based on years since genesis
func (k Keeper) GetCurrentCharityPercentage(ctx sdk.Context) sdk.Dec {
	genesisTime := k.GetGenesisTime(ctx)
	currentTime := ctx.BlockTime()
	
	// Calculate years since genesis
	yearsSinceGenesis := currentTime.Sub(genesisTime).Hours() / (24 * 365)
	
	// Get charity allocation params
	params := k.GetCharityAllocationParams(ctx)
	
	// Return percentage based on year
	switch {
	case yearsSinceGenesis < 1:
		return params.Year1Percentage
	case yearsSinceGenesis < 2:
		return params.Year2Percentage
	case yearsSinceGenesis < 3:
		return params.Year3Percentage
	case yearsSinceGenesis < 4:
		return params.Year4Percentage
	default:
		return params.Year5Percentage
	}
}

// SetCharityAllocationParams stores charity allocation parameters
func (k Keeper) SetCharityAllocationParams(ctx sdk.Context, params CharityAllocationParams) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.KeyCharityAllocationParams, bz)
}

// GetCharityAllocationParams retrieves charity allocation parameters
func (k Keeper) GetCharityAllocationParams(ctx sdk.Context) CharityAllocationParams {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCharityAllocationParams)
	
	if bz == nil {
		// Return defaults if not set
		return DefaultCharityAllocationParams(k.GetGenesisTime(ctx))
	}
	
	var params CharityAllocationParams
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// CalculateCharityAmount calculates charity amount based on current percentage
func (k Keeper) CalculateCharityAmount(ctx sdk.Context, totalRevenue sdk.Coins) sdk.Coins {
	percentage := k.GetCurrentCharityPercentage(ctx)
	
	charityCoins := sdk.Coins{}
	for _, coin := range totalRevenue {
		charityAmount := sdk.NewDecFromInt(coin.Amount).Mul(percentage).TruncateInt()
		if charityAmount.IsPositive() {
			charityCoin := sdk.NewCoin(coin.Denom, charityAmount)
			charityCoins = charityCoins.Add(charityCoin)
		}
	}
	
	return charityCoins
}

// GetCharityAllocationInfo returns current charity allocation information
func (k Keeper) GetCharityAllocationInfo(ctx sdk.Context) types.CharityAllocationInfo {
	currentPercentage := k.GetCurrentCharityPercentage(ctx)
	genesisTime := k.GetGenesisTime(ctx)
	currentTime := ctx.BlockTime()
	yearsSinceGenesis := currentTime.Sub(genesisTime).Hours() / (24 * 365)
	
	return types.CharityAllocationInfo{
		CurrentPercentage:  currentPercentage,
		YearsSinceGenesis:  int32(yearsSinceGenesis),
		NextMilestone:      k.getNextCharityMilestone(yearsSinceGenesis),
		GenesisTime:        genesisTime,
		CurrentTime:        currentTime,
	}
}

// getNextCharityMilestone returns the next charity percentage milestone
func (k Keeper) getNextCharityMilestone(yearsSinceGenesis float64) string {
	switch {
	case yearsSinceGenesis < 1:
		return "Year 2: 20% charity allocation"
	case yearsSinceGenesis < 2:
		return "Year 3: 30% charity allocation"
	case yearsSinceGenesis < 3:
		return "Year 4: 35% charity allocation"
	case yearsSinceGenesis < 4:
		return "Year 5: 40% charity allocation"
	default:
		return "Maximum allocation reached: 40%"
	}
}

// EmitCharityAllocationEvent emits an event for charity allocation changes
func (k Keeper) EmitCharityAllocationEvent(ctx sdk.Context, percentage sdk.Dec, year int) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCharityAllocationUpdate,
			sdk.NewAttribute(types.AttributeKeyCharityPercentage, percentage.String()),
			sdk.NewAttribute(types.AttributeKeyYear, fmt.Sprintf("%d", year)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	)
}

// UpdateCharityAllocationYearly is called during EndBlock to check and update charity percentage
func (k Keeper) UpdateCharityAllocationYearly(ctx sdk.Context) {
	info := k.GetCharityAllocationInfo(ctx)
	
	// Check if we've crossed a year boundary
	lastCheckedYear := k.GetLastCharityUpdateYear(ctx)
	currentYear := int(info.YearsSinceGenesis) + 1
	
	if currentYear > lastCheckedYear {
		// Update the year
		k.SetLastCharityUpdateYear(ctx, currentYear)
		
		// Emit event for the milestone
		k.EmitCharityAllocationEvent(ctx, info.CurrentPercentage, currentYear)
		
		// Log the update
		k.Logger(ctx).Info("Charity allocation updated",
			"year", currentYear,
			"percentage", info.CurrentPercentage.String(),
			"next_milestone", info.NextMilestone,
		)
	}
}

// SetLastCharityUpdateYear stores the last year charity was updated
func (k Keeper) SetLastCharityUpdateYear(ctx sdk.Context, year int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastCharityUpdateYear, sdk.Uint64ToBigEndian(uint64(year)))
}

// GetLastCharityUpdateYear retrieves the last year charity was updated
func (k Keeper) GetLastCharityUpdateYear(ctx sdk.Context) int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastCharityUpdateYear)
	if bz == nil {
		return 0
	}
	return int(sdk.BigEndianToUint64(bz))
}