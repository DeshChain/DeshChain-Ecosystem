package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		RecipientAddresses: RecipientAddresses{
			DevelopmentFund:   "",
			CommunityTreasury: "",
			LiquidityPool:     "",
			EmergencyReserve:  "",
			FounderRoyalty:    "",
			ValidatorPool:     "",
		},
	}
}

// ValidateGenesis validates the provided genesis state
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	
	// Validate recipient addresses (empty is allowed)
	if err := validateAddress(data.RecipientAddresses.DevelopmentFund, "Development fund"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.CommunityTreasury, "Community treasury"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.LiquidityPool, "Liquidity pool"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.EmergencyReserve, "Emergency reserve"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.FounderRoyalty, "Founder royalty"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.ValidatorPool, "Validator pool"); err != nil {
		return err
	}
	
	return nil
}

func validateAddress(address, name string) error {
	if address == "" {
		// Empty address is allowed in default genesis
		return nil
	}
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return ErrInvalidRecipient
	}
	return nil
}

// GenesisState defines the revenue module's genesis state.
type GenesisState struct {
	Params             Params             `json:"params" yaml:"params"`
	RecipientAddresses RecipientAddresses `json:"recipient_addresses" yaml:"recipient_addresses"`
}

// RecipientAddresses holds all revenue recipient addresses
type RecipientAddresses struct {
	DevelopmentFund   string `json:"development_fund" yaml:"development_fund"`
	CommunityTreasury string `json:"community_treasury" yaml:"community_treasury"`
	LiquidityPool     string `json:"liquidity_pool" yaml:"liquidity_pool"`
	EmergencyReserve  string `json:"emergency_reserve" yaml:"emergency_reserve"`
	FounderRoyalty    string `json:"founder_royalty" yaml:"founder_royalty"`
	ValidatorPool     string `json:"validator_pool" yaml:"validator_pool"`
}