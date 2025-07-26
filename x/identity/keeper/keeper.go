package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	
	// Module-specific keepers
	tradeFinanceKeeper types.TradeFinanceKeeper
	moneyOrderKeeper   types.MoneyOrderKeeper
	validatorKeeper    types.ValidatorKeeper
	
	// External services (optional)
	aadhaarService    types.AadhaarService
	digiLockerService types.DigiLockerService
	upiService        types.UPIService
	biometricService  types.BiometricService
}

// NewKeeper creates a new identity Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetModuleKeepers sets the module-specific keepers
func (k *Keeper) SetModuleKeepers(
	tradeFinanceKeeper types.TradeFinanceKeeper,
	moneyOrderKeeper types.MoneyOrderKeeper,
	validatorKeeper types.ValidatorKeeper,
) {
	k.tradeFinanceKeeper = tradeFinanceKeeper
	k.moneyOrderKeeper = moneyOrderKeeper
	k.validatorKeeper = validatorKeeper
}

// SetExternalServices sets the external service interfaces
func (k *Keeper) SetExternalServices(
	aadhaarService types.AadhaarService,
	digiLockerService types.DigiLockerService,
	upiService types.UPIService,
	biometricService types.BiometricService,
) {
	k.aadhaarService = aadhaarService
	k.digiLockerService = digiLockerService
	k.upiService = upiService
	k.biometricService = biometricService
}

// GetParams gets all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MaxDIDDocumentSize(ctx),
		k.MaxCredentialSize(ctx),
		k.MaxProofSize(ctx),
		k.CredentialExpiryDays(ctx),
		k.KYCExpiryDays(ctx),
		k.BiometricExpiryDays(ctx),
		k.MaxRecoveryMethods(ctx),
		k.MaxCredentialsPerIdentity(ctx),
		k.MinAnonymitySetSize(ctx),
		k.ProofExpiryMinutes(ctx),
		k.MaxConsentDurationDays(ctx),
		k.EnableAnonymousCredentials(ctx),
		k.EnableZKProofs(ctx),
		k.EnableIndiaStack(ctx),
		k.RequireKYCForHighValue(ctx),
		k.HighValueThreshold(ctx),
		k.MaxFailedAuthAttempts(ctx),
		k.AuthLockoutDurationHours(ctx),
		k.SupportedDIDMethods(ctx),
		k.SupportedProofSystems(ctx),
		k.TrustedIssuers(ctx),
	)
}

// SetParams sets the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// Individual parameter getters

// MaxDIDDocumentSize returns the maximum DID document size
func (k Keeper) MaxDIDDocumentSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxDIDDocumentSize, &res)
	return
}

// MaxCredentialSize returns the maximum credential size
func (k Keeper) MaxCredentialSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxCredentialSize, &res)
	return
}

// MaxProofSize returns the maximum proof size
func (k Keeper) MaxProofSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxProofSize, &res)
	return
}

// CredentialExpiryDays returns the credential expiry days
func (k Keeper) CredentialExpiryDays(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyCredentialExpiryDays, &res)
	return
}

// KYCExpiryDays returns the KYC expiry days
func (k Keeper) KYCExpiryDays(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyKYCExpiryDays, &res)
	return
}

// BiometricExpiryDays returns the biometric expiry days
func (k Keeper) BiometricExpiryDays(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyBiometricExpiryDays, &res)
	return
}

// MaxRecoveryMethods returns the maximum recovery methods
func (k Keeper) MaxRecoveryMethods(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxRecoveryMethods, &res)
	return
}

// MaxCredentialsPerIdentity returns the maximum credentials per identity
func (k Keeper) MaxCredentialsPerIdentity(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxCredentialsPerIdentity, &res)
	return
}

// MinAnonymitySetSize returns the minimum anonymity set size
func (k Keeper) MinAnonymitySetSize(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMinAnonymitySetSize, &res)
	return
}

// ProofExpiryMinutes returns the proof expiry minutes
func (k Keeper) ProofExpiryMinutes(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyProofExpiryMinutes, &res)
	return
}

// MaxConsentDurationDays returns the maximum consent duration days
func (k Keeper) MaxConsentDurationDays(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxConsentDurationDays, &res)
	return
}

// EnableAnonymousCredentials returns whether anonymous credentials are enabled
func (k Keeper) EnableAnonymousCredentials(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyEnableAnonymousCredentials, &res)
	return
}

// EnableZKProofs returns whether ZK proofs are enabled
func (k Keeper) EnableZKProofs(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyEnableZKProofs, &res)
	return
}

// EnableIndiaStack returns whether India Stack is enabled
func (k Keeper) EnableIndiaStack(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyEnableIndiaStack, &res)
	return
}

// RequireKYCForHighValue returns whether KYC is required for high value
func (k Keeper) RequireKYCForHighValue(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyRequireKYCForHighValue, &res)
	return
}

// HighValueThreshold returns the high value threshold
func (k Keeper) HighValueThreshold(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyHighValueThreshold, &res)
	return
}

// MaxFailedAuthAttempts returns the maximum failed auth attempts
func (k Keeper) MaxFailedAuthAttempts(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxFailedAuthAttempts, &res)
	return
}

// AuthLockoutDurationHours returns the auth lockout duration hours
func (k Keeper) AuthLockoutDurationHours(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyAuthLockoutDurationHours, &res)
	return
}

// SupportedDIDMethods returns the supported DID methods
func (k Keeper) SupportedDIDMethods(ctx sdk.Context) (res []string) {
	k.paramstore.Get(ctx, types.KeySupportedDIDMethods, &res)
	return
}

// SupportedProofSystems returns the supported proof systems
func (k Keeper) SupportedProofSystems(ctx sdk.Context) (res []string) {
	k.paramstore.Get(ctx, types.KeySupportedProofSystems, &res)
	return
}

// TrustedIssuers returns the trusted issuers
func (k Keeper) TrustedIssuers(ctx sdk.Context) (res []string) {
	k.paramstore.Get(ctx, types.KeyTrustedIssuers, &res)
	return
}