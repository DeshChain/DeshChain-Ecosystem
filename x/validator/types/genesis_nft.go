package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisValidatorNFT represents the special NFT for genesis validators
type GenesisValidatorNFT struct {
    TokenID         uint64   `json:"token_id"`
    Rank            uint32   `json:"rank"`
    ValidatorAddr   string   `json:"validator_address"`
    SanskritName    string   `json:"sanskrit_name"`
    EnglishName     string   `json:"english_name"`
    Title           string   `json:"title"`
    MintHeight      int64    `json:"mint_height"`
    SpecialPowers   []string `json:"special_powers"`
    Tradeable       bool     `json:"tradeable"`
    CurrentOwner    string   `json:"current_owner"`
    
    // Stake binding fields
    BoundStakeAmount   sdk.Int   `json:"bound_stake_amount"`
    StakeBindingActive bool      `json:"stake_binding_active"`
    
    // Trading history
    TradeCount         uint32    `json:"trade_count"`
    LastTradePrice     sdk.Coins `json:"last_trade_price"`
    LastTradeHeight    int64     `json:"last_trade_height"`
}

// GenesisNFTMetadata defines the metadata for each genesis validator NFT
var GenesisNFTMetadata = map[uint32]GenesisValidatorNFT{
    1: {
        Rank:         1,
        SanskritName: "परम रक्षक",
        EnglishName:  "Param Rakshak",
        Title:        "The Supreme Guardian of DeshChain",
        SpecialPowers: []string{
            "2x Governance Weight",
            "Genesis Crown Badge",
            "Golden UI Theme",
            "Priority Block Proposals",
        },
        Tradeable: true,
    },
    2: {
        Rank:         2,
        SanskritName: "महा सेनानी",
        EnglishName:  "Maha Senani",
        Title:        "The Great General",
        SpecialPowers: []string{
            "1.5x Governance Weight",
            "Battle Armor Theme",
            "War Room Access",
        },
        Tradeable: true,
    },
    3: {
        Rank:         3,
        SanskritName: "धर्म पालक",
        EnglishName:  "Dharma Palak",
        Title:        "Keeper of Righteousness",
        SpecialPowers: []string{
            "1.3x Governance Weight",
            "Justice Scale Badge",
            "Dispute Resolution Priority",
        },
        Tradeable: true,
    },
    // ... continue for all 21
    21: {
        Rank:         21,
        SanskritName: "भारत गौरव",
        EnglishName:  "Bharat Gaurav",
        Title:        "Pride of India",
        SpecialPowers: []string{
            "1.1x Governance Weight",
            "Tricolor Theme",
            "National Holiday Bonuses",
        },
        Tradeable: true,
    },
}

// ValidatorRevenueDistribution calculates revenue share for each validator
type ValidatorRevenueDistribution struct {
    TotalValidators   uint32
    GenesisValidators []string // Top 21 addresses
    TotalRevenue      sdk.Coins
}

// CalculateValidatorShares computes individual validator shares based on count
func (vrd *ValidatorRevenueDistribution) CalculateValidatorShares() map[string]sdk.Coins {
    shares := make(map[string]sdk.Coins)
    validatorCount := vrd.TotalValidators
    
    switch {
    case validatorCount <= 21:
        // Equal distribution among all validators
        sharePerValidator := vrd.TotalRevenue.QuoInt(sdk.NewInt(int64(validatorCount)))
        // Distribute equally
        
    case validatorCount > 21 && validatorCount <= 100:
        // Each validator gets exactly 1%
        onePercent := vrd.TotalRevenue.QuoInt(sdk.NewInt(100))
        // Logic for distribution
        
    case validatorCount > 100:
        // Genesis validators (top 21) get 21% (1% each)
        // Remaining 79% distributed equally among ALL validators
        genesisShare := vrd.TotalRevenue.MulInt(sdk.NewInt(21)).QuoInt(sdk.NewInt(100))
        remainingShare := vrd.TotalRevenue.MulInt(sdk.NewInt(79)).QuoInt(sdk.NewInt(100))
        
        // Each genesis validator gets 1% + their portion of 79%
        genesisIndividual := genesisShare.QuoInt(sdk.NewInt(21))
        equalShare := remainingShare.QuoInt(sdk.NewInt(int64(validatorCount)))
        
        // Set shares
        for i, addr := range vrd.GenesisValidators {
            if i < 21 {
                shares[addr] = genesisIndividual.Add(equalShare...)
            }
        }
        // Other validators get only equal share from 79%
    }
    
    return shares
}

// NFTTradeRequest represents a trade request for genesis validator NFT
type NFTTradeRequest struct {
    TokenID      uint64    `json:"token_id"`
    FromAddress  string    `json:"from_address"`
    ToAddress    string    `json:"to_address"`
    Price        sdk.Coins `json:"price"`
    MinPrice     sdk.Coins `json:"min_price"` // 10,000 NAMO minimum
}

// ValidateNFTTrade ensures trade meets minimum requirements
func ValidateNFTTrade(req NFTTradeRequest) error {
    minPrice := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(10000000000))) // 10,000 NAMO
    if req.Price.IsAllLT(minPrice) {
        return ErrPriceBelowMinimum
    }
    return nil
}