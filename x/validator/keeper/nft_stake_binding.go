package keeper

import (
    "fmt"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
)

// BindNFTToStake creates an unbreakable bond between NFT and stake
func (k Keeper) BindNFTToStake(
    ctx sdk.Context,
    nftID uint64,
    validatorAddr string,
    stakeAmount sdk.Int,
) error {
    // Get the NFT
    nft, found := k.GetGenesisNFT(ctx, nftID)
    if !found {
        return fmt.Errorf("NFT %d not found", nftID)
    }
    
    // Verify ownership
    if nft.ValidatorAddr != validatorAddr {
        return fmt.Errorf("NFT does not belong to validator %s", validatorAddr)
    }
    
    // Get the stake
    stake, found := k.GetValidatorStake(ctx, validatorAddr)
    if !found {
        return fmt.Errorf("no stake found for validator %s", validatorAddr)
    }
    
    // Create binding in NFT metadata
    nft.BoundStakeAmount = stakeAmount
    nft.StakeBindingActive = true
    k.SetGenesisNFT(ctx, nft)
    
    // Create reverse binding in stake
    stake.BoundNFTID = nftID
    stake.NFTBindingActive = true
    k.SetValidatorStake(ctx, stake)
    
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "nft_stake_bound",
            sdk.NewAttribute("nft_id", fmt.Sprintf("%d", nftID)),
            sdk.NewAttribute("validator", validatorAddr),
            sdk.NewAttribute("stake_amount", stakeAmount.String()),
        ),
    )
    
    return nil
}

// TransferNFTWithStake handles the atomic transfer of NFT and its bound stake
func (k Keeper) TransferNFTWithStake(
    ctx sdk.Context,
    nftID uint64,
    fromAddr string,
    toAddr string,
    price sdk.Coins,
) error {
    // Get the NFT
    nft, found := k.GetGenesisNFT(ctx, nftID)
    if !found {
        return fmt.Errorf("NFT %d not found", nftID)
    }
    
    // Verify current ownership
    if nft.CurrentOwner != fromAddr {
        return fmt.Errorf("NFT not owned by %s", fromAddr)
    }
    
    // Check if NFT has stake binding
    if !nft.StakeBindingActive {
        return fmt.Errorf("NFT has no active stake binding")
    }
    
    // Get the bound stake
    stake, found := k.GetValidatorStake(ctx, nft.ValidatorAddr)
    if !found {
        return fmt.Errorf("no stake found for NFT's validator")
    }
    
    // Verify stake binding
    if stake.BoundNFTID != nftID || !stake.NFTBindingActive {
        return fmt.Errorf("stake binding mismatch")
    }
    
    // Check minimum holding period (6 months for all NFTs)
    sixMonths := nft.MintHeight + (6 * 30 * 24 * 60 * 60 / 5) // Assuming 5s blocks
    if ctx.BlockHeight() < sixMonths {
        blocksRemaining := sixMonths - ctx.BlockHeight()
        return fmt.Errorf("NFT locked for %d more blocks", blocksRemaining)
    }
    
    // Process payment
    from, _ := sdk.AccAddressFromBech32(fromAddr)
    to, _ := sdk.AccAddressFromBech32(toAddr)
    
    // Calculate 5% royalty to original validator
    royalty := sdk.NewCoins()
    netPrice := sdk.NewCoins()
    
    for _, coin := range price {
        royaltyAmount := coin.Amount.ToDec().Mul(sdk.NewDecWithPrec(5, 2)).TruncateInt()
        royaltyCoin := sdk.NewCoin(coin.Denom, royaltyAmount)
        netCoin := sdk.NewCoin(coin.Denom, coin.Amount.Sub(royaltyAmount))
        
        royalty = royalty.Add(royaltyCoin)
        netPrice = netPrice.Add(netCoin)
    }
    
    // Transfer payment
    if err := k.bankKeeper.SendCoins(ctx, to, from, netPrice); err != nil {
        return fmt.Errorf("failed to send payment: %w", err)
    }
    
    // Pay royalty to original validator
    origValidator, _ := sdk.AccAddressFromBech32(nft.ValidatorAddr)
    if err := k.bankKeeper.SendCoins(ctx, to, origValidator, royalty); err != nil {
        return fmt.Errorf("failed to send royalty: %w", err)
    }
    
    // Transfer 5% to treasury as transfer fee
    treasuryFee := sdk.NewCoins()
    for _, coin := range price {
        feeAmount := coin.Amount.ToDec().Mul(sdk.NewDecWithPrec(5, 2)).TruncateInt()
        treasuryFee = treasuryFee.Add(sdk.NewCoin(coin.Denom, feeAmount))
    }
    
    treasuryAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
    if err := k.bankKeeper.SendCoins(ctx, to, treasuryAddr, treasuryFee); err != nil {
        return fmt.Errorf("failed to send treasury fee: %w", err)
    }
    
    // Now transfer the stake obligations
    // The new owner inherits ALL stake conditions
    stake.ValidatorAddr = toAddr
    // Vesting schedule continues unchanged
    // Lock periods remain the same
    // Performance bond stays locked
    
    // Update NFT ownership
    nft.CurrentOwner = toAddr
    nft.TradeCount++
    nft.LastTradePrice = price
    nft.LastTradeHeight = ctx.BlockHeight()
    
    // Save updates
    k.SetGenesisNFT(ctx, nft)
    k.SetValidatorStake(ctx, stake)
    
    // Transfer validator rights
    if err := k.TransferValidatorRights(ctx, fromAddr, toAddr); err != nil {
        return fmt.Errorf("failed to transfer validator rights: %w", err)
    }
    
    // Emit comprehensive event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeNFTTransferred,
            sdk.NewAttribute(types.AttributeKeyTokenID, fmt.Sprintf("%d", nftID)),
            sdk.NewAttribute(types.AttributeKeyFrom, fromAddr),
            sdk.NewAttribute(types.AttributeKeyTo, toAddr),
            sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
            sdk.NewAttribute(types.AttributeKeyRoyalty, royalty.String()),
            sdk.NewAttribute("treasury_fee", treasuryFee.String()),
            sdk.NewAttribute("stake_transferred", stake.NAMOTokensStaked.String()),
            sdk.NewAttribute("vesting_continues", "true"),
        ),
    )
    
    k.Logger(ctx).Info("NFT transferred with stake",
        "nft_id", nftID,
        "from", fromAddr,
        "to", toAddr,
        "price", price,
        "stake", stake.NAMOTokensStaked)
    
    return nil
}

// ValidateNFTStakeIntegrity ensures NFT and stake remain bound
func (k Keeper) ValidateNFTStakeIntegrity(ctx sdk.Context) error {
    // Iterate through all NFTs
    nfts := k.GetAllGenesisNFTs(ctx)
    
    for _, nft := range nfts {
        if !nft.StakeBindingActive {
            continue
        }
        
        // Check corresponding stake exists
        stake, found := k.GetValidatorStake(ctx, nft.CurrentOwner)
        if !found {
            k.Logger(ctx).Error("Orphaned NFT found",
                "nft_id", nft.TokenID,
                "owner", nft.CurrentOwner)
            continue
        }
        
        // Verify binding integrity
        if stake.BoundNFTID != nft.TokenID {
            k.Logger(ctx).Error("NFT-stake binding mismatch",
                "nft_id", nft.TokenID,
                "stake_bound_nft", stake.BoundNFTID)
        }
        
        // Verify stake amount matches
        if !stake.NAMOTokensStaked.Equal(nft.BoundStakeAmount) {
            k.Logger(ctx).Error("Stake amount mismatch",
                "nft_id", nft.TokenID,
                "nft_stake", nft.BoundStakeAmount,
                "actual_stake", stake.NAMOTokensStaked)
        }
    }
    
    return nil
}

// GetAllGenesisNFTs returns all genesis NFTs
func (k Keeper) GetAllGenesisNFTs(ctx sdk.Context) []types.GenesisValidatorNFT {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, types.GenesisNFTKeyPrefix)
    defer iterator.Close()
    
    var nfts []types.GenesisValidatorNFT
    for ; iterator.Valid(); iterator.Next() {
        var nft types.GenesisValidatorNFT
        k.cdc.MustUnmarshal(iterator.Value(), &nft)
        nfts = append(nfts, nft)
    }
    
    return nfts
}