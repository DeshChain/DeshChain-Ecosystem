package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Validator represents a validator in our system
// This is a simplified version for backward compatibility
type Validator struct {
	OperatorAddress string                     `json:"operator_address"`
	Tokens          sdk.Int                    `json:"tokens"`
	Status          stakingtypes.BondStatus    `json:"status"`
	JoinOrder       uint32                     `json:"join_order,omitempty"`
	StakeInfo       *ValidatorStake            `json:"stake_info,omitempty"`
	NFTInfo         *GenesisValidatorNFT       `json:"nft_info,omitempty"`
}