package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// StakingKeeper defines the expected staking keeper for validator verification
// type StakingKeeper interface {
// 	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
// 	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
// }

// DistributionKeeper defines the expected distribution keeper
// type DistributionKeeper interface {
// 	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
// }

// GovKeeper defines the expected governance keeper for proposal verification
// type GovKeeper interface {
// 	GetProposal(ctx sdk.Context, proposalID uint64) (proposal govtypes.Proposal, found bool)
// }

// IbcKeeper defines the expected IBC keeper for cross-chain identity
// type IbcKeeper interface {
// 	SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
// }

// Additional expected keepers for DeshChain modules

// TradeFinanceKeeper defines the expected trade finance keeper
type TradeFinanceKeeper interface {
	GetKYCProfile(ctx sdk.Context, address string) (interface{}, bool)
	SetKYCProfile(ctx sdk.Context, profile interface{})
	GetBusinessProfile(ctx sdk.Context, address string) (interface{}, bool)
}

// MoneyOrderKeeper defines the expected money order keeper
type MoneyOrderKeeper interface {
	GetUser(ctx sdk.Context, address string) (interface{}, bool)
	SetUser(ctx sdk.Context, user interface{})
	GetBiometricRegistration(ctx sdk.Context, biometricID string) (interface{}, bool)
}

// ValidatorKeeper defines the expected validator keeper
type ValidatorKeeper interface {
	GetValidatorProfile(ctx sdk.Context, address string) (interface{}, bool)
	GetValidatorKYCStatus(ctx sdk.Context, address string) (bool, string)
}

// CulturalKeeper defines the expected cultural keeper for heritage verification
type CulturalKeeper interface {
	GetCulturalProfile(ctx sdk.Context, address string) (interface{}, bool)
	VerifyCulturalCredentials(ctx sdk.Context, address string, credentialType string) (bool, error)
}

// DonationKeeper defines the expected donation keeper for charity verification
type DonationKeeper interface {
	GetNGOProfile(ctx sdk.Context, address string) (interface{}, bool)
	IsVerifiedNGO(ctx sdk.Context, address string) bool
}

// External service interfaces

// AadhaarService defines the interface for Aadhaar verification
type AadhaarService interface {
	VerifyAadhaar(ctx sdk.Context, aadhaarHash string, otp string) (bool, error)
	GetDemographicData(ctx sdk.Context, aadhaarHash string, consentToken string) (map[string]interface{}, error)
}

// DigiLockerService defines the interface for DigiLocker integration
type DigiLockerService interface {
	AuthenticateUser(ctx sdk.Context, authToken string) (string, error)
	FetchDocument(ctx sdk.Context, userID string, documentURI string) ([]byte, error)
	GetIssuedDocuments(ctx sdk.Context, userID string) ([]string, error)
}

// UPIService defines the interface for UPI verification
type UPIService interface {
	VerifyVPA(ctx sdk.Context, vpa string) (bool, error)
	GetLinkedAccounts(ctx sdk.Context, vpa string, consentToken string) ([]string, error)
}

// BiometricService defines the interface for biometric operations
type BiometricService interface {
	EnrollBiometric(ctx sdk.Context, userID string, biometricType string, template []byte) (string, error)
	VerifyBiometric(ctx sdk.Context, userID string, biometricType string, template []byte) (float64, error)
	GetQualityScore(ctx sdk.Context, template []byte) (float64, error)
}