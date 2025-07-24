package keeper

import (
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// DistributeValidatorRevenue distributes revenue according to the tiered model
func (k Keeper) DistributeValidatorRevenue(ctx sdk.Context, totalRevenue sdk.Coins) error {
    // Get all active validators
    validators := k.GetAllActiveValidators(ctx)
    validatorCount := len(validators)
    
    if validatorCount == 0 {
        return fmt.Errorf("no active validators found")
    }
    
    // Get genesis validators (first 21 by join order)
    genesisValidators := k.GetGenesisValidators(ctx)
    
    distribution := make(map[string]sdk.Coins)
    
    switch {
    case validatorCount <= 21:
        // Equal distribution among all validators
        sharePerValidator := totalRevenue.QuoInt(sdk.NewInt(int64(validatorCount)))
        
        for _, val := range validators {
            distribution[val.OperatorAddress] = sharePerValidator
        }
        
        k.Logger(ctx).Info("Validator revenue distributed equally",
            "validator_count", validatorCount,
            "share_per_validator", sharePerValidator.String())
    
    case validatorCount > 21 && validatorCount <= 100:
        // Each validator gets exactly 1% of total revenue
        onePercent := totalRevenue.QuoInt(sdk.NewInt(100))
        
        for _, val := range validators {
            distribution[val.OperatorAddress] = onePercent
        }
        
        k.Logger(ctx).Info("Validator revenue distributed at 1% each",
            "validator_count", validatorCount,
            "share_per_validator", onePercent.String())
    
    case validatorCount > 100:
        // Genesis validators (top 21) get 21% guaranteed (1% each)
        // Remaining 79% distributed equally among ALL validators
        
        // Calculate 21% for genesis validators
        genesisShare := totalRevenue.MulInt(sdk.NewInt(21)).QuoInt(sdk.NewInt(100))
        genesisIndividualShare := genesisShare.QuoInt(sdk.NewInt(21))
        
        // Calculate 79% for all validators
        remainingShare := totalRevenue.MulInt(sdk.NewInt(79)).QuoInt(sdk.NewInt(100))
        equalShare := remainingShare.QuoInt(sdk.NewInt(int64(validatorCount)))
        
        // Distribute to all validators
        for _, val := range validators {
            valAddr := val.OperatorAddress
            
            // Check if this is a genesis validator
            isGenesis := false
            for i, genVal := range genesisValidators {
                if i < 21 && genVal.OperatorAddress == valAddr {
                    isGenesis = true
                    break
                }
            }
            
            if isGenesis {
                // Genesis validator gets 1% + equal share from 79%
                distribution[valAddr] = genesisIndividualShare.Add(equalShare...)
            } else {
                // Regular validator gets only equal share from 79%
                distribution[valAddr] = equalShare
            }
        }
        
        k.Logger(ctx).Info("Validator revenue distributed with genesis bonus",
            "total_validators", validatorCount,
            "genesis_bonus", genesisIndividualShare.String(),
            "equal_share", equalShare.String())
    }
    
    // Execute the distribution
    for valAddr, amount := range distribution {
        if err := k.SendCoinsToValidator(ctx, valAddr, amount); err != nil {
            return fmt.Errorf("failed to send coins to validator %s: %w", valAddr, err)
        }
        
        // Emit distribution event
        ctx.EventManager().EmitEvent(
            sdk.NewEvent(
                types.EventTypeValidatorReward,
                sdk.NewAttribute(types.AttributeKeyValidator, valAddr),
                sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
                sdk.NewAttribute(types.AttributeKeyDistributionType, k.getDistributionType(validatorCount)),
            ),
        )
    }
    
    return nil
}

// GetGenesisValidators returns the first 21 validators by join order
func (k Keeper) GetGenesisValidators(ctx sdk.Context) []types.Validator {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, types.ValidatorByJoinOrderKey)
    defer iterator.Close()
    
    var genesisValidators []types.Validator
    count := 0
    
    for ; iterator.Valid() && count < 21; iterator.Next() {
        var val types.Validator
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        genesisValidators = append(genesisValidators, val)
        count++
    }
    
    return genesisValidators
}

// MintGenesisNFT mints the special NFT for genesis validators
func (k Keeper) MintGenesisNFT(ctx sdk.Context, validatorAddr string, rank uint32) (uint64, error) {
    if rank > 21 {
        return 0, fmt.Errorf("only top 21 validators eligible for genesis NFT")
    }
    
    metadata, exists := types.GenesisNFTMetadata[rank]
    if !exists {
        return 0, fmt.Errorf("no metadata found for rank %d", rank)
    }
    
    // Generate unique token ID
    tokenID := k.GetNextNFTID(ctx)
    
    // Create NFT
    nft := types.GenesisValidatorNFT{
        TokenID:       tokenID,
        Rank:          rank,
        ValidatorAddr: validatorAddr,
        SanskritName:  metadata.SanskritName,
        EnglishName:   metadata.EnglishName,
        Title:         metadata.Title,
        MintHeight:    ctx.BlockHeight(),
        SpecialPowers: metadata.SpecialPowers,
        Tradeable:     true,
        CurrentOwner:  validatorAddr,
    }
    
    // Store NFT
    k.SetGenesisNFT(ctx, nft)
    
    // Emit minting event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeNFTMinted,
            sdk.NewAttribute(types.AttributeKeyTokenID, fmt.Sprintf("%d", tokenID)),
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr),
            sdk.NewAttribute(types.AttributeKeyRank, fmt.Sprintf("%d", rank)),
            sdk.NewAttribute(types.AttributeKeyNFTName, metadata.EnglishName),
        ),
    )
    
    k.Logger(ctx).Info("Genesis validator NFT minted",
        "token_id", tokenID,
        "rank", rank,
        "validator", validatorAddr,
        "name", metadata.EnglishName)
    
    return tokenID, nil
}

// TransferGenesisNFT handles the transfer of genesis NFT with validation
func (k Keeper) TransferGenesisNFT(ctx sdk.Context, req types.NFTTradeRequest) error {
    // Validate minimum price
    if err := types.ValidateNFTTrade(req); err != nil {
        return err
    }
    
    // Get NFT
    nft, found := k.GetGenesisNFT(ctx, req.TokenID)
    if !found {
        return fmt.Errorf("NFT with ID %d not found", req.TokenID)
    }
    
    // Verify ownership
    if nft.CurrentOwner != req.FromAddress {
        return fmt.Errorf("sender does not own this NFT")
    }
    
    // Calculate royalty (5% to original validator)
    royalty := req.Price.MulInt(sdk.NewInt(5)).QuoInt(sdk.NewInt(100))
    netPrice := req.Price.Sub(royalty)
    
    // Transfer payment
    if err := k.bankKeeper.SendCoins(ctx, 
        sdk.MustAccAddressFromBech32(req.ToAddress),
        sdk.MustAccAddressFromBech32(req.FromAddress),
        netPrice); err != nil {
        return err
    }
    
    // Pay royalty to original validator
    if err := k.bankKeeper.SendCoins(ctx,
        sdk.MustAccAddressFromBech32(req.ToAddress),
        sdk.MustAccAddressFromBech32(nft.ValidatorAddr),
        royalty); err != nil {
        return err
    }
    
    // Update NFT ownership
    nft.CurrentOwner = req.ToAddress
    k.SetGenesisNFT(ctx, nft)
    
    // Transfer validator rights if applicable
    if nft.Rank <= 21 {
        k.TransferValidatorRights(ctx, req.FromAddress, req.ToAddress)
    }
    
    // Emit transfer event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeNFTTransferred,
            sdk.NewAttribute(types.AttributeKeyTokenID, fmt.Sprintf("%d", req.TokenID)),
            sdk.NewAttribute(types.AttributeKeyFrom, req.FromAddress),
            sdk.NewAttribute(types.AttributeKeyTo, req.ToAddress),
            sdk.NewAttribute(types.AttributeKeyPrice, req.Price.String()),
            sdk.NewAttribute(types.AttributeKeyRoyalty, royalty.String()),
        ),
    )
    
    return nil
}

func (k Keeper) getDistributionType(validatorCount int) string {
    switch {
    case validatorCount <= 21:
        return "equal_all"
    case validatorCount <= 100:
        return "one_percent_each"
    default:
        return "genesis_bonus_plus_equal"
    }
}