package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// TradeFinanceIntegration provides backward-compatible KYC integration for tradefinance module
type TradeFinanceIntegration struct {
	keeper *Keeper
}

// NewTradeFinanceIntegration creates a new integration adapter
func NewTradeFinanceIntegration(k *Keeper) *TradeFinanceIntegration {
	return &TradeFinanceIntegration{keeper: k}
}

// CreateKYCIdentity creates a DID-based identity from tradefinance KYC data
func (tfi *TradeFinanceIntegration) CreateKYCIdentity(
	ctx sdk.Context,
	customerID string,
	kycData map[string]interface{},
) (string, error) {
	// Generate DID for the customer
	did := fmt.Sprintf("did:desh:%s", customerID)
	
	// Check if identity already exists
	if _, exists := tfi.keeper.GetIdentity(ctx, did); exists {
		return did, nil
	}

	// Create identity from KYC data
	identity := types.Identity{
		Did:        did,
		Controller: customerID,
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
		Metadata: map[string]string{
			"source":      "tradefinance",
			"customer_id": customerID,
		},
	}

	// Store identity
	tfi.keeper.SetIdentity(ctx, identity)

	// Create KYC credential from tradefinance data
	if err := tfi.createKYCCredential(ctx, did, kycData); err != nil {
		return "", fmt.Errorf("failed to create KYC credential: %w", err)
	}

	return did, nil
}

// VerifyKYCStatus checks if a customer has valid KYC using identity module
func (tfi *TradeFinanceIntegration) VerifyKYCStatus(
	ctx sdk.Context,
	customerID string,
) (bool, *types.VerifiableCredential, error) {
	did := fmt.Sprintf("did:desh:%s", customerID)
	
	// Get identity
	identity, exists := tfi.keeper.GetIdentity(ctx, did)
	if !exists {
		return false, nil, nil
	}

	// Check if identity is active
	if identity.Status != types.IdentityStatus_ACTIVE {
		return false, nil, nil
	}

	// Find KYC credential
	credentials := tfi.keeper.GetCredentialsBySubject(ctx, did)
	for _, credID := range credentials {
		cred, found := tfi.keeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a KYC credential
		for _, credType := range cred.Type {
			if credType == "KYCCredential" || credType == "TradeFinanceKYC" {
				// Verify credential status
				if cred.Status == nil || cred.Status.Type != "revoked" {
					// Check expiry
					if cred.ExpirationDate == nil || cred.ExpirationDate.After(ctx.BlockTime()) {
						return true, &cred, nil
					}
				}
			}
		}
	}

	return false, nil, nil
}

// UpdateKYCCredential updates existing KYC credential with new data
func (tfi *TradeFinanceIntegration) UpdateKYCCredential(
	ctx sdk.Context,
	customerID string,
	kycData map[string]interface{},
) error {
	did := fmt.Sprintf("did:desh:%s", customerID)
	
	// Get existing credential
	isValid, existingCred, err := tfi.VerifyKYCStatus(ctx, customerID)
	if err != nil {
		return err
	}

	if isValid && existingCred != nil {
		// Revoke existing credential
		tfi.keeper.UpdateCredentialStatus(ctx, existingCred.ID, &types.CredentialStatus{
			Type:   "revoked",
			Reason: "Updated with new KYC data",
		})
	}

	// Create new credential with updated data
	return tfi.createKYCCredential(ctx, did, kycData)
}

// createKYCCredential creates a verifiable credential from KYC data
func (tfi *TradeFinanceIntegration) createKYCCredential(
	ctx sdk.Context,
	did string,
	kycData map[string]interface{},
) error {
	// Generate credential ID
	credID := fmt.Sprintf("vc:kyc:%s:%d", did, ctx.BlockTime().Unix())
	
	// Set expiry based on KYC type
	expiryDays := uint32(180) // Default 6 months
	if kycLevel, ok := kycData["kyc_level"].(string); ok {
		switch kycLevel {
		case "enhanced":
			expiryDays = 365 // 1 year
		case "basic":
			expiryDays = 90 // 3 months
		}
	}
	
	expiryDate := ctx.BlockTime().AddDate(0, 0, int(expiryDays))

	// Create credential
	credential := types.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.com/contexts/kyc/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "TradeFinanceKYC"},
		Issuer: "did:desh:tradefinance-kyc-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":          did,
			"kyc_data":    kycData,
			"verified_at": ctx.BlockTime().Format(time.RFC3339),
		},
		Proof: &types.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:tradefinance-kyc-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign with issuer's key
		},
	}

	// Store credential
	tfi.keeper.SetCredential(ctx, credential)

	// Link credential to subject
	tfi.keeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// MigrateExistingKYC migrates existing tradefinance KYC profiles to identity module
func (tfi *TradeFinanceIntegration) MigrateExistingKYC(
	ctx sdk.Context,
	profiles []TradeFinanceKYCProfile,
) (int, error) {
	migrated := 0
	
	for _, profile := range profiles {
		// Convert tradefinance KYC data to generic format
		kycData := map[string]interface{}{
			"customer_id":    profile.CustomerID,
			"kyc_level":      profile.KYCLevel,
			"risk_rating":    profile.RiskRating,
			"customer_type":  profile.CustomerType,
			"personal_info":  profile.PersonalInfo,
			"business_info":  profile.BusinessInfo,
			"pep_status":     profile.PEPStatus,
			"sanctions_check": profile.SanctionsScreening,
			"verified_docs":  len(profile.Documents),
			"last_review":    profile.LastReviewDate,
		}

		// Create identity and credential
		if _, err := tfi.CreateKYCIdentity(ctx, profile.CustomerID, kycData); err != nil {
			tfi.keeper.Logger(ctx).Error("Failed to migrate KYC profile", 
				"customer_id", profile.CustomerID, 
				"error", err)
			continue
		}

		migrated++
	}

	return migrated, nil
}

// GetKYCLevel extracts KYC level from identity credentials
func (tfi *TradeFinanceIntegration) GetKYCLevel(
	ctx sdk.Context,
	customerID string,
) (string, error) {
	isValid, cred, err := tfi.VerifyKYCStatus(ctx, customerID)
	if err != nil {
		return "", err
	}

	if !isValid || cred == nil {
		return "none", nil
	}

	// Extract KYC level from credential
	if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
		if kycData, ok := subject["kyc_data"].(map[string]interface{}); ok {
			if level, ok := kycData["kyc_level"].(string); ok {
				return level, nil
			}
		}
	}

	return "basic", nil
}

// Backward compatibility types
type TradeFinanceKYCProfile struct {
	CustomerID          string
	KYCLevel            string
	RiskRating          string
	CustomerType        string
	LastReviewDate      time.Time
	PersonalInfo        map[string]interface{}
	BusinessInfo        map[string]interface{}
	Documents           []interface{}
	PEPStatus           map[string]interface{}
	SanctionsScreening  map[string]interface{}
}