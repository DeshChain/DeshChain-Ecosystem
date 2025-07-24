package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "validator"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// InsurancePoolName defines the insurance pool module account name
	InsurancePoolName = "validator_insurance_pool"
)

// Store key prefixes
var (
	ValidatorStakeKeyPrefix     = []byte{0x01}
	GenesisNFTKeyPrefix        = []byte{0x02}
	ValidatorByJoinOrderKey    = []byte{0x03}
	SlashingEventKeyPrefix     = []byte{0x04}
	InsuranceClaimKeyPrefix    = []byte{0x05}
	NextNFTIDKey               = []byte{0x06}
	ValidatorRevenueKeyPrefix  = []byte{0x07}
	NFTTradeHistoryKeyPrefix   = []byte{0x08}
	
	// Referral system keys
	ReferralKeyPrefix          = []byte{0x09}
	ReferralStatsKeyPrefix     = []byte{0x0A}
	ValidatorTokenKeyPrefix    = []byte{0x0B}
	CommissionPayoutKeyPrefix  = []byte{0x0C}
	NextReferralIDKey          = []byte{0x0D}
	NextTokenIDKey             = []byte{0x0E}
	NextPayoutIDKey            = []byte{0x0F}
)

// GetValidatorStakeKey returns the store key for a validator stake
func GetValidatorStakeKey(validatorAddr string) []byte {
	return append(ValidatorStakeKeyPrefix, []byte(validatorAddr)...)
}

// GetGenesisNFTKey returns the store key for a genesis NFT
func GetGenesisNFTKey(tokenID uint64) []byte {
	return append(GenesisNFTKeyPrefix, sdk.Uint64ToBigEndian(tokenID)...)
}

// GetSlashingEventKey returns the store key for slashing events
func GetSlashingEventKey(validatorAddr string, timestamp int64) []byte {
	return append(
		append(SlashingEventKeyPrefix, []byte(validatorAddr)...),
		sdk.Uint64ToBigEndian(uint64(timestamp))...,
	)
}

// GetInsuranceClaimKey returns the store key for insurance claims
func GetInsuranceClaimKey(claimID uint64) []byte {
	return append(InsuranceClaimKeyPrefix, sdk.Uint64ToBigEndian(claimID)...)
}

// Referral system key functions

// GetReferralKey returns the store key for a referral
func GetReferralKey(referralID uint64) []byte {
	return append(ReferralKeyPrefix, sdk.Uint64ToBigEndian(referralID)...)
}

// GetReferralStatsKey returns the store key for referral stats
func GetReferralStatsKey(validatorAddr string) []byte {
	return append(ReferralStatsKeyPrefix, []byte(validatorAddr)...)
}

// GetValidatorTokenKey returns the store key for a validator token
func GetValidatorTokenKey(tokenID uint64) []byte {
	return append(ValidatorTokenKeyPrefix, sdk.Uint64ToBigEndian(tokenID)...)
}

// GetCommissionPayoutKey returns the store key for a commission payout
func GetCommissionPayoutKey(payoutID uint64) []byte {
	return append(CommissionPayoutKeyPrefix, sdk.Uint64ToBigEndian(payoutID)...)
}