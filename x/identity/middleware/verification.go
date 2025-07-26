package middleware

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
)

// IdentityVerificationDecorator provides identity verification middleware
type IdentityVerificationDecorator struct {
	identityKeeper identitykeeper.Keeper
	config         VerificationConfig
}

// VerificationConfig defines configuration for identity verification
type VerificationConfig struct {
	// RequireIdentity specifies if identity is required for all transactions
	RequireIdentity bool
	
	// RequireKYC specifies if KYC verification is required
	RequireKYC bool
	
	// MinKYCLevel specifies minimum KYC level required
	MinKYCLevel string
	
	// HighValueThreshold defines threshold for high-value transactions
	HighValueThreshold sdk.Int
	
	// RequireBiometricForHighValue requires biometric auth for high-value tx
	RequireBiometricForHighValue bool
	
	// ExemptMessages lists message types exempt from verification
	ExemptMessages []string
	
	// ModuleSpecificRules defines per-module verification rules
	ModuleSpecificRules map[string]ModuleVerificationRule
}

// ModuleVerificationRule defines verification rules for a specific module
type ModuleVerificationRule struct {
	RequireIdentity bool
	RequireKYC      bool
	MinKYCLevel     string
	RequireBiometric bool
	CustomRules     []CustomVerificationRule
}

// CustomVerificationRule allows custom verification logic
type CustomVerificationRule struct {
	Name      string
	CheckFunc func(ctx sdk.Context, tx sdk.Tx, signer sdk.AccAddress) error
}

// NewIdentityVerificationDecorator creates a new identity verification decorator
func NewIdentityVerificationDecorator(
	identityKeeper identitykeeper.Keeper,
	config VerificationConfig,
) IdentityVerificationDecorator {
	return IdentityVerificationDecorator{
		identityKeeper: identityKeeper,
		config:         config,
	}
}

// AnteHandle implements the AnteDecorator interface
func (ivd IdentityVerificationDecorator) AnteHandle(
	ctx sdk.Context, 
	tx sdk.Tx, 
	simulate bool, 
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// Skip verification in simulation mode
	if simulate {
		return next(ctx, tx, simulate)
	}

	// Get transaction messages
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return ctx, fmt.Errorf("transaction must contain at least one message")
	}

	// Get signers
	signers := tx.GetSigners()
	if len(signers) == 0 {
		return ctx, fmt.Errorf("transaction must have at least one signer")
	}

	// Check each message and signer
	for _, msg := range msgs {
		// Check if message is exempt
		if ivd.isMessageExempt(msg) {
			continue
		}

		// Get module-specific rules
		moduleRules := ivd.getModuleRules(msg)
		
		// Verify each signer
		for _, signer := range msg.GetSigners() {
			if err := ivd.verifyIdentity(ctx, signer, msg, moduleRules); err != nil {
				return ctx, fmt.Errorf("identity verification failed for %s: %w", signer, err)
			}
		}
	}

	// Check transaction value for high-value verification
	if ivd.config.RequireBiometricForHighValue {
		if err := ivd.checkHighValueTransaction(ctx, tx); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}

// verifyIdentity verifies the identity of a signer
func (ivd IdentityVerificationDecorator) verifyIdentity(
	ctx sdk.Context,
	signer sdk.AccAddress,
	msg sdk.Msg,
	rules ModuleVerificationRule,
) error {
	// Check if identity is required
	requireIdentity := ivd.config.RequireIdentity || rules.RequireIdentity
	if !requireIdentity {
		return nil
	}

	// Get identity
	did := fmt.Sprintf("did:desh:%s", signer.String())
	identity, exists := ivd.identityKeeper.GetIdentity(ctx, did)
	if !exists {
		return fmt.Errorf("identity not found for address %s", signer)
	}

	// Check identity status
	if identity.Status != identitytypes.IdentityStatus_ACTIVE {
		return fmt.Errorf("identity is not active: %s", identity.Status)
	}

	// Check KYC if required
	requireKYC := ivd.config.RequireKYC || rules.RequireKYC
	if requireKYC {
		if err := ivd.verifyKYC(ctx, did, rules.MinKYCLevel); err != nil {
			return err
		}
	}

	// Check biometric if required
	if rules.RequireBiometric {
		if err := ivd.verifyBiometric(ctx, did); err != nil {
			return err
		}
	}

	// Run custom rules
	for _, customRule := range rules.CustomRules {
		if err := customRule.CheckFunc(ctx, nil, signer); err != nil {
			return fmt.Errorf("custom rule '%s' failed: %w", customRule.Name, err)
		}
	}

	// Store verification in context for other modules
	ctx = ctx.WithValue(identityVerifiedKey(signer.String()), true)

	return nil
}

// verifyKYC verifies KYC credentials
func (ivd IdentityVerificationDecorator) verifyKYC(
	ctx sdk.Context,
	did string,
	minLevel string,
) error {
	// Get credentials
	credentials := ivd.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	// Find valid KYC credential
	var validKYC *identitytypes.VerifiableCredential
	for _, credID := range credentials {
		cred, found := ivd.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a KYC credential
		isKYC := false
		for _, credType := range cred.Type {
			if credType == "KYCCredential" {
				isKYC = true
				break
			}
		}

		if !isKYC {
			continue
		}

		// Check if credential is valid
		if cred.Status != nil && cred.Status.Type == "revoked" {
			continue
		}

		// Check expiry
		if cred.ExpirationDate != nil && cred.ExpirationDate.Before(ctx.BlockTime()) {
			continue
		}

		validKYC = &cred
		break
	}

	if validKYC == nil {
		return fmt.Errorf("no valid KYC credential found")
	}

	// Check KYC level if specified
	if minLevel != "" {
		if err := ivd.checkKYCLevel(validKYC, minLevel); err != nil {
			return err
		}
	}

	return nil
}

// verifyBiometric verifies biometric authentication
func (ivd IdentityVerificationDecorator) verifyBiometric(ctx sdk.Context, did string) error {
	// Get recent biometric authentication
	credentials := ivd.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	for _, credID := range credentials {
		cred, found := ivd.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a biometric credential
		isBiometric := false
		for _, credType := range cred.Type {
			if credType == "BiometricCredential" {
				isBiometric = true
				break
			}
		}

		if isBiometric && ivd.isCredentialRecent(ctx, cred, 300) { // 5 minutes
			return nil
		}
	}

	return fmt.Errorf("recent biometric authentication required")
}

// Helper methods

func (ivd IdentityVerificationDecorator) isMessageExempt(msg sdk.Msg) bool {
	msgType := sdk.MsgTypeURL(msg)
	for _, exempt := range ivd.config.ExemptMessages {
		if msgType == exempt {
			return true
		}
	}
	return false
}

func (ivd IdentityVerificationDecorator) getModuleRules(msg sdk.Msg) ModuleVerificationRule {
	// Extract module name from message route
	route := msg.Route()
	if rules, ok := ivd.config.ModuleSpecificRules[route]; ok {
		return rules
	}
	
	// Return default rules
	return ModuleVerificationRule{
		RequireIdentity: ivd.config.RequireIdentity,
		RequireKYC:      ivd.config.RequireKYC,
		MinKYCLevel:     ivd.config.MinKYCLevel,
	}
}

func (ivd IdentityVerificationDecorator) checkKYCLevel(
	cred *identitytypes.VerifiableCredential,
	minLevel string,
) error {
	// Extract KYC level from credential
	if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
		if level, ok := subject["kyc_level"].(string); ok {
			if !ivd.isKYCLevelSufficient(level, minLevel) {
				return fmt.Errorf("KYC level %s is insufficient, requires %s", level, minLevel)
			}
		}
	}
	return nil
}

func (ivd IdentityVerificationDecorator) isKYCLevelSufficient(actual, required string) bool {
	levels := map[string]int{
		"basic":    1,
		"standard": 2,
		"enhanced": 3,
		"premium":  4,
	}
	
	actualLevel, ok1 := levels[actual]
	requiredLevel, ok2 := levels[required]
	
	if !ok1 || !ok2 {
		return false
	}
	
	return actualLevel >= requiredLevel
}

func (ivd IdentityVerificationDecorator) isCredentialRecent(
	ctx sdk.Context,
	cred identitytypes.VerifiableCredential,
	seconds int64,
) bool {
	// Check if credential was issued recently
	timeDiff := ctx.BlockTime().Sub(cred.IssuanceDate)
	return timeDiff.Seconds() <= float64(seconds)
}

func (ivd IdentityVerificationDecorator) checkHighValueTransaction(
	ctx sdk.Context,
	tx sdk.Tx,
) error {
	// Check if transaction contains high-value transfers
	// This is a simplified check - in production, implement comprehensive value checking
	
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil
	}
	
	fee := feeTx.GetFee()
	totalValue := sdk.NewInt(0)
	
	for _, coin := range fee {
		totalValue = totalValue.Add(coin.Amount)
	}
	
	if totalValue.GT(ivd.config.HighValueThreshold) {
		// Require biometric authentication
		signers := tx.GetSigners()
		for _, signer := range signers {
			did := fmt.Sprintf("did:desh:%s", signer.String())
			if err := ivd.verifyBiometric(ctx, did); err != nil {
				return fmt.Errorf("biometric required for high-value transaction: %w", err)
			}
		}
	}
	
	return nil
}

// Context key for storing verification status
type contextKey string

func identityVerifiedKey(address string) contextKey {
	return contextKey("identity_verified_" + address)
}

// IsIdentityVerified checks if an address has been verified in the current context
func IsIdentityVerified(ctx context.Context, address string) bool {
	verified, ok := ctx.Value(identityVerifiedKey(address)).(bool)
	return ok && verified
}

// DefaultVerificationConfig returns a default verification configuration
func DefaultVerificationConfig() VerificationConfig {
	return VerificationConfig{
		RequireIdentity:              false, // Opt-in by default
		RequireKYC:                   false,
		MinKYCLevel:                  "basic",
		HighValueThreshold:           sdk.NewInt(1000000), // 1M NAMO
		RequireBiometricForHighValue: true,
		ExemptMessages: []string{
			// Exempt basic query messages
			"/cosmos.bank.v1beta1.MsgSend",
			"/cosmos.staking.v1beta1.MsgDelegate",
		},
		ModuleSpecificRules: map[string]ModuleVerificationRule{
			"gramsuraksha": {
				RequireIdentity:  true,
				RequireKYC:       true,
				MinKYCLevel:      "standard",
				RequireBiometric: false,
			},
			"tradefinance": {
				RequireIdentity:  true,
				RequireKYC:       true,
				MinKYCLevel:      "enhanced",
				RequireBiometric: true,
			},
			"moneyorder": {
				RequireIdentity:  true,
				RequireKYC:       false,
				MinKYCLevel:      "",
				RequireBiometric: true,
			},
		},
	}
}