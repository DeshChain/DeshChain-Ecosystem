package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		RecipientAddresses: RecipientAddresses{
			NgoWallet:         "", // To be set in genesis
			ValidatorPool:     "",
			CommunityPool:     "",
			TechInnovation:    "",
			Operations:        "",
			TalentAcquisition: "",
			StrategicReserve:  "",
			FounderWallet:     "",
			CoFoundersWallet:  "",
			AngelWallet:       "",
		},
	}
}

// ValidateGenesis validates the provided genesis state
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	
	// Validate recipient addresses
	if err := validateAddress(data.RecipientAddresses.NgoWallet, "NGO wallet"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.ValidatorPool, "Validator pool"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.CommunityPool, "Community pool"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.TechInnovation, "Tech innovation"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.Operations, "Operations"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.TalentAcquisition, "Talent acquisition"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.StrategicReserve, "Strategic reserve"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.FounderWallet, "Founder wallet"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.CoFoundersWallet, "Co-founders wallet"); err != nil {
		return err
	}
	if err := validateAddress(data.RecipientAddresses.AngelWallet, "Angel wallet"); err != nil {
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
		return WrapInvalidAddress(name, err)
	}
	return nil
}

// GenesisState defines the tax module's genesis state.
type GenesisState struct {
	Params             Params             `json:"params" yaml:"params"`
	RecipientAddresses RecipientAddresses `json:"recipient_addresses" yaml:"recipient_addresses"`
}

// RecipientAddresses holds all tax recipient addresses
type RecipientAddresses struct {
	NgoWallet         string `json:"ngo_wallet" yaml:"ngo_wallet"`
	ValidatorPool     string `json:"validator_pool" yaml:"validator_pool"`
	CommunityPool     string `json:"community_pool" yaml:"community_pool"`
	TechInnovation    string `json:"tech_innovation" yaml:"tech_innovation"`
	Operations        string `json:"operations" yaml:"operations"`
	TalentAcquisition string `json:"talent_acquisition" yaml:"talent_acquisition"`
	StrategicReserve  string `json:"strategic_reserve" yaml:"strategic_reserve"`
	FounderWallet     string `json:"founder_wallet" yaml:"founder_wallet"`
	CoFoundersWallet  string `json:"co_founders_wallet" yaml:"co_founders_wallet"`
	AngelWallet       string `json:"angel_wallet" yaml:"angel_wallet"`
}