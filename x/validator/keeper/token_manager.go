package keeper

import (
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
)

// TokenManager handles validator token management operations
type TokenManager struct {
    keeper              Keeper
    sikkebaazIntegration *SikkebaazIntegration
}

// NewTokenManager creates a new token manager
func NewTokenManager(k Keeper) *TokenManager {
    return &TokenManager{
        keeper:              k,
        sikkebaazIntegration: NewSikkebaazIntegration(k),
    }
}

// LaunchValidatorToken manually launches a validator token if eligible
func (tm *TokenManager) LaunchValidatorToken(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
) (uint64, error) {
    // Check eligibility
    stats := tm.keeper.GetReferralStats(ctx, validatorAddr.String())
    if stats.TokenLaunched {
        return 0, fmt.Errorf("validator token already launched")
    }
    
    conditions := types.GetTokenLaunchConditions()
    
    // Check referral count or commission threshold
    eligible := stats.TotalReferrals >= conditions.MinReferrals ||
                stats.TotalCommission.GTE(conditions.MinCommission)
    
    if !eligible {
        return 0, fmt.Errorf("validator not eligible for token launch: need %d referrals OR %s commission, have %d referrals and %s commission",
            conditions.MinReferrals,
            conditions.MinCommission.String(),
            stats.TotalReferrals,
            stats.TotalCommission.String())
    }
    
    // Check NFT ownership
    nfts := tm.keeper.GetAllGenesisNFTs(ctx)
    var validatorNFT *types.GenesisValidatorNFT
    for _, nft := range nfts {
        if nft.CurrentOwner == validatorAddr.String() && nft.Rank <= 21 {
            validatorNFT = &nft
            break
        }
    }
    
    if validatorNFT == nil {
        return 0, fmt.Errorf("validator does not own a genesis NFT")
    }
    
    // Launch token
    return tm.createAndLaunchToken(ctx, validatorAddr, *validatorNFT, "manual")
}

// AirdropTokens distributes validator tokens to specified addresses
func (tm *TokenManager) AirdropTokens(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    recipients []types.AirdropRecipient,
) error {
    // Verify token ownership
    token, found := tm.keeper.GetValidatorToken(ctx, tokenID)
    if !found {
        return fmt.Errorf("token not found: %d", tokenID)
    }
    
    if token.ValidatorAddr != validatorAddr.String() {
        return fmt.Errorf("validator does not own this token")
    }
    
    // Validate recipients
    if len(recipients) == 0 {
        return fmt.Errorf("no recipients specified")
    }
    
    if len(recipients) > 1000 {
        return fmt.Errorf("too many recipients (max 1000)")
    }
    
    // Calculate total airdrop amount
    totalAmount := sdk.ZeroInt()
    for _, recipient := range recipients {
        if recipient.Amount.IsZero() || recipient.Amount.IsNegative() {
            return fmt.Errorf("invalid amount for recipient %s: %s", recipient.Address, recipient.Amount)
        }
        totalAmount = totalAmount.Add(recipient.Amount)
    }
    
    // Check if validator has enough allocation
    remainingAllocation := tm.calculateRemainingAirdropAllocation(ctx, token)
    if totalAmount.GT(remainingAllocation) {
        return fmt.Errorf("insufficient airdrop allocation: requested %s, available %s",
            totalAmount.String(), remainingAllocation.String())
    }
    
    // Execute airdrop
    for _, recipient := range recipients {
        if err := tm.transferTokens(ctx, token, recipient.Address, recipient.Amount); err != nil {
            return fmt.Errorf("failed to airdrop to %s: %w", recipient.Address, err)
        }
    }
    
    // Update token allocation tracking
    token.AirdropAllocation = token.AirdropAllocation.Sub(totalAmount)
    tm.keeper.SetValidatorToken(ctx, token)
    
    // Create airdrop record
    airdropRecord := types.AirdropRecord{
        TokenID:        tokenID,
        ValidatorAddr:  validatorAddr.String(),
        Recipients:     recipients,
        TotalAmount:    totalAmount,
        AirdropTime:    ctx.BlockTime(),
        BlockHeight:    ctx.BlockHeight(),
    }
    tm.keeper.SetAirdropRecord(ctx, airdropRecord)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeTokensAirdropped,
            sdk.NewAttribute(types.AttributeKeyTokenID, fmt.Sprintf("%d", tokenID)),
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeyAmount, totalAmount.String()),
            sdk.NewAttribute("recipients_count", fmt.Sprintf("%d", len(recipients))),
        ),
    )
    
    return nil
}

// UpdateTokenParameters allows validators to modify certain token parameters
func (tm *TokenManager) UpdateTokenParameters(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    params types.TokenUpdateParams,
) error {
    // Verify token ownership
    token, found := tm.keeper.GetValidatorToken(ctx, tokenID)
    if !found {
        return fmt.Errorf("token not found: %d", tokenID)
    }
    
    if token.ValidatorAddr != validatorAddr.String() {
        return fmt.Errorf("validator does not own this token")
    }
    
    // Check launch protection period (24 hours)
    if ctx.BlockTime().Before(token.LaunchedAt.Add(24 * time.Hour)) {
        return fmt.Errorf("token parameters cannot be updated within 24 hours of launch")
    }
    
    // Validate parameter ranges
    if err := tm.validateTokenParameters(params); err != nil {
        return fmt.Errorf("invalid parameters: %w", err)
    }
    
    // Apply updates
    if params.SellTaxPercent != nil {
        token.SellTaxPercent = *params.SellTaxPercent
    }
    if params.BuyTaxPercent != nil {
        token.BuyTaxPercent = *params.BuyTaxPercent
    }
    if params.CooldownSeconds != nil {
        token.CooldownSeconds = *params.CooldownSeconds
    }
    
    // Update token in Sikkebaaz
    if err := tm.sikkebaazIntegration.UpdateTokenParameters(ctx, token, params); err != nil {
        return fmt.Errorf("failed to update token parameters on Sikkebaaz: %w", err)
    }
    
    tm.keeper.SetValidatorToken(ctx, token)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "validator_token_updated",
            sdk.NewAttribute(types.AttributeKeyTokenID, fmt.Sprintf("%d", tokenID)),
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
        ),
    )
    
    return nil
}

// GetTokenPerformance returns comprehensive token performance metrics
func (tm *TokenManager) GetTokenPerformance(
    ctx sdk.Context,
    tokenID uint64,
) (types.TokenPerformance, error) {
    token, found := tm.keeper.GetValidatorToken(ctx, tokenID)
    if !found {
        return types.TokenPerformance{}, fmt.Errorf("token not found: %d", tokenID)
    }
    
    performance := types.TokenPerformance{
        TokenID:       tokenID,
        TokenName:     token.TokenName,
        TokenSymbol:   token.TokenSymbol,
        LaunchedAt:    token.LaunchedAt,
        ValidatorAddr: token.ValidatorAddr,
    }
    
    // Get current price from Sikkebaaz
    currentPrice, err := tm.sikkebaazIntegration.GetTokenPrice(ctx, tokenID)
    if err != nil {
        currentPrice = token.CurrentPrice // Fallback to stored price
    }
    
    performance.CurrentPrice = currentPrice
    performance.MarketCap = token.TotalSupply.ToDec().Mul(currentPrice).TruncateInt()
    
    // Calculate price changes
    initialPrice := sdk.OneDec() // Assuming 1:1 launch price
    performance.PriceChange24h = tm.calculatePriceChange(currentPrice, initialPrice)
    performance.PriceChange7d = tm.calculatePriceChange(currentPrice, initialPrice) // Simplified
    
    // Volume metrics (would integrate with Sikkebaaz for real data)
    performance.Volume24h = sdk.NewInt(1000000) // Placeholder
    performance.Volume7d = sdk.NewInt(7000000)  // Placeholder
    
    // Liquidity metrics
    performance.LiquidityLocked = token.LiquidityAllocation
    performance.LiquidityProvider = token.ValidatorAddr
    
    // Holder metrics (would query Sikkebaaz for real data)
    performance.HolderCount = 150 // Placeholder
    performance.TopHolders = tm.getTopHolders(ctx, tokenID)
    
    // Tax collection metrics
    performance.TaxCollected = tm.calculateTaxCollected(ctx, tokenID)
    performance.TaxDistribution = tm.getTaxDistribution(ctx, tokenID)
    
    return performance, nil
}

// GetTokenManagementInfo returns comprehensive token management information
func (tm *TokenManager) GetTokenManagementInfo(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
) (types.TokenManagementInfo, error) {
    stats := tm.keeper.GetReferralStats(ctx, validatorAddr.String())
    
    info := types.TokenManagementInfo{
        ValidatorAddr: validatorAddr.String(),
        TokenLaunched: stats.TokenLaunched,
    }
    
    if stats.TokenLaunched {
        token, found := tm.keeper.GetValidatorToken(ctx, stats.TokenID)
        if found {
            info.TokenID = stats.TokenID
            info.Token = &token
            
            // Get performance metrics
            performance, err := tm.GetTokenPerformance(ctx, stats.TokenID)
            if err == nil {
                info.Performance = &performance
            }
            
            // Get airdrop history
            info.AirdropHistory = tm.keeper.GetAirdropRecordsByToken(ctx, stats.TokenID)
            
            // Calculate available allocations
            info.AvailableAllocations = types.TokenAllocations{
                ValidatorAllocation: tm.calculateRemainingValidatorAllocation(ctx, token),
                AirdropAllocation:   tm.calculateRemainingAirdropAllocation(ctx, token),
                LiquidityAllocation: token.LiquidityAllocation,
            }
        }
    } else {
        // Check launch eligibility
        conditions := types.GetTokenLaunchConditions()
        info.LaunchEligible = stats.TotalReferrals >= conditions.MinReferrals ||
                             stats.TotalCommission.GTE(conditions.MinCommission)
        info.LaunchConditions = conditions
    }
    
    return info, nil
}

// Helper functions

func (tm *TokenManager) createAndLaunchToken(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    nft types.GenesisValidatorNFT,
    trigger string,
) (uint64, error) {
    tokenID := tm.keeper.GetNextTokenID(ctx)
    totalSupply := sdk.NewInt(1000000000000000) // 1 billion with 6 decimals
    
    token := types.ValidatorToken{
        TokenID:       tokenID,
        ValidatorAddr: validatorAddr.String(),
        TokenName:     fmt.Sprintf("%s Coin", nft.EnglishName),
        TokenSymbol:   tm.generateTokenSymbol(nft.EnglishName),
        TotalSupply:   totalSupply,
        Decimals:      6,
        LogoURI:       fmt.Sprintf("/nfts/%d.png", nft.TokenID),
        
        // Distribution
        ValidatorAllocation:   totalSupply.MulRaw(40).QuoRaw(100),
        LiquidityAllocation:   totalSupply.MulRaw(30).QuoRaw(100),
        AirdropAllocation:     totalSupply.MulRaw(15).QuoRaw(100),
        DevelopmentAllocation: totalSupply.MulRaw(10).QuoRaw(100),
        InitialLiquidity:      totalSupply.MulRaw(5).QuoRaw(100),
        
        // Launch info
        LaunchedAt:    ctx.BlockTime(),
        LaunchTrigger: trigger,
        
        // Anti-dump parameters
        MaxWalletPercent: sdk.NewDecWithPrec(2, 2),   // 2%
        MaxTxPercent:     sdk.NewDecWithPrec(5, 3),   // 0.5%
        SellTaxPercent:   sdk.NewDecWithPrec(5, 2),   // 5%
        BuyTaxPercent:    sdk.NewDecWithPrec(2, 2),   // 2%
        CooldownSeconds:  3600,                        // 1 hour
    }
    
    // Launch on Sikkebaaz
    if err := tm.sikkebaazIntegration.LaunchValidatorToken(ctx, token); err != nil {
        return 0, fmt.Errorf("failed to launch token on Sikkebaaz: %w", err)
    }
    
    tm.keeper.SetValidatorToken(ctx, token)
    tm.keeper.SetNextTokenID(ctx, tokenID+1)
    
    // Update validator stats
    stats := tm.keeper.GetReferralStats(ctx, validatorAddr.String())
    stats.TokenLaunched = true
    stats.TokenID = tokenID
    tm.keeper.SetReferralStats(ctx, stats)
    
    return tokenID, nil
}

func (tm *TokenManager) generateTokenSymbol(nftName string) string {
    // Simple symbol generation from NFT name
    words := []rune(nftName)
    symbol := ""
    
    for i, char := range words {
        if i == 0 || (i > 0 && words[i-1] == ' ') {
            symbol += string(char)
        }
        if len(symbol) >= 4 {
            break
        }
    }
    
    if len(symbol) < 3 {
        symbol = nftName[:min(4, len(nftName))]
    }
    
    return symbol
}

func (tm *TokenManager) validateTokenParameters(params types.TokenUpdateParams) error {
    if params.SellTaxPercent != nil {
        if params.SellTaxPercent.GT(sdk.NewDecWithPrec(10, 2)) {
            return fmt.Errorf("sell tax cannot exceed 10%%")
        }
        if params.SellTaxPercent.LT(sdk.ZeroDec()) {
            return fmt.Errorf("sell tax cannot be negative")
        }
    }
    
    if params.BuyTaxPercent != nil {
        if params.BuyTaxPercent.GT(sdk.NewDecWithPrec(5, 2)) {
            return fmt.Errorf("buy tax cannot exceed 5%%")
        }
        if params.BuyTaxPercent.LT(sdk.ZeroDec()) {
            return fmt.Errorf("buy tax cannot be negative")
        }
    }
    
    if params.CooldownSeconds != nil {
        if *params.CooldownSeconds > 7200 { // 2 hours max
            return fmt.Errorf("cooldown cannot exceed 2 hours")
        }
        if *params.CooldownSeconds < 0 {
            return fmt.Errorf("cooldown cannot be negative")
        }
    }
    
    return nil
}

func (tm *TokenManager) calculateRemainingAirdropAllocation(
    ctx sdk.Context,
    token types.ValidatorToken,
) sdk.Int {
    // Get all airdrop records for this token
    records := tm.keeper.GetAirdropRecordsByToken(ctx, token.TokenID)
    
    totalAirdropped := sdk.ZeroInt()
    for _, record := range records {
        totalAirdropped = totalAirdropped.Add(record.TotalAmount)
    }
    
    originalAllocation := token.TotalSupply.MulRaw(15).QuoRaw(100) // 15%
    remaining := originalAllocation.Sub(totalAirdropped)
    
    if remaining.IsNegative() {
        return sdk.ZeroInt()
    }
    
    return remaining
}

func (tm *TokenManager) calculateRemainingValidatorAllocation(
    ctx sdk.Context,
    token types.ValidatorToken,
) sdk.Int {
    // This would track validator withdrawals/vesting
    // For now, return full allocation minus any claimed amount
    return token.ValidatorAllocation
}

func (tm *TokenManager) transferTokens(
    ctx sdk.Context,
    token types.ValidatorToken,
    recipientAddr string,
    amount sdk.Int,
) error {
    // This would integrate with Sikkebaaz to transfer actual tokens
    // For now, emit event to track the transfer
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "token_transfer",
            sdk.NewAttribute("token_id", fmt.Sprintf("%d", token.TokenID)),
            sdk.NewAttribute("recipient", recipientAddr),
            sdk.NewAttribute("amount", amount.String()),
        ),
    )
    
    return nil
}

func (tm *TokenManager) calculatePriceChange(current, previous sdk.Dec) sdk.Dec {
    if previous.IsZero() {
        return sdk.ZeroDec()
    }
    
    return current.Sub(previous).Quo(previous).Mul(sdk.NewDec(100))
}

func (tm *TokenManager) getTopHolders(ctx sdk.Context, tokenID uint64) []types.TokenHolder {
    // This would query Sikkebaaz for actual holder data
    // For now, return placeholder data
    return []types.TokenHolder{
        {Address: "deshchain1...", Balance: sdk.NewInt(50000000), Percentage: sdk.NewDecWithPrec(5, 2)},
        {Address: "deshchain2...", Balance: sdk.NewInt(30000000), Percentage: sdk.NewDecWithPrec(3, 2)},
        {Address: "deshchain3...", Balance: sdk.NewInt(20000000), Percentage: sdk.NewDecWithPrec(2, 2)},
    }
}

func (tm *TokenManager) calculateTaxCollected(ctx sdk.Context, tokenID uint64) sdk.Int {
    // This would integrate with Sikkebaaz to get actual tax collection data
    return sdk.NewInt(5000000) // Placeholder
}

func (tm *TokenManager) getTaxDistribution(ctx sdk.Context, tokenID uint64) types.TaxDistribution {
    // This would get actual tax distribution from Sikkebaaz
    return types.TaxDistribution{
        LiquidityPercent: sdk.NewDecWithPrec(60, 2), // 60%
        ValidatorPercent: sdk.NewDecWithPrec(30, 2), // 30%
        PlatformPercent:  sdk.NewDecWithPrec(10, 2), // 10%
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}