package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
)

// IdentityAdapter provides backward-compatible access to the new identity module
type IdentityAdapter struct {
	keeper         *Keeper
	identityKeeper identitykeeper.Keeper
	integration    *identitykeeper.TradeFinanceIntegration
}

// NewIdentityAdapter creates a new adapter for identity integration
func NewIdentityAdapter(k *Keeper, ik identitykeeper.Keeper) *IdentityAdapter {
	return &IdentityAdapter{
		keeper:         k,
		identityKeeper: ik,
		integration:    identitykeeper.NewTradeFinanceIntegration(&ik),
	}
}

// PerformKYCWithIdentity performs KYC using the new identity module while maintaining compatibility
func (ia *IdentityAdapter) PerformKYCWithIdentity(
	ctx sdk.Context,
	customerID string,
	submittedData CustomerSubmission,
) (*KYCAssessmentResult, error) {
	// First, perform traditional KYC assessment
	kycEngine := NewKYCAMLEngine(ia.keeper)
	result, err := kycEngine.PerformKYCAssessment(ctx, customerID, AssessmentTypeStandard, submittedData)
	if err != nil {
		return nil, err
	}

	// Convert KYC data for identity module
	kycData := ia.convertKYCDataForIdentity(result, submittedData)

	// Create or update identity
	did, err := ia.integration.CreateKYCIdentity(ctx, customerID, kycData)
	if err != nil {
		ia.keeper.Logger(ctx).Error("Failed to create identity", "error", err)
		// Don't fail the KYC process if identity creation fails
		// This maintains backward compatibility
	} else {
		// Store DID reference in result
		result.IdentityDID = did
	}

	return result, nil
}

// VerifyCustomerIdentity checks customer identity using both old and new systems
func (ia *IdentityAdapter) VerifyCustomerIdentity(
	ctx sdk.Context,
	customerID string,
) (*CustomerIdentityStatus, error) {
	status := &CustomerIdentityStatus{
		CustomerID:     customerID,
		HasTraditionalKYC: false,
		HasDIDIdentity:    false,
	}

	// Check traditional KYC
	kycEngine := NewKYCAMLEngine(ia.keeper)
	profile, err := kycEngine.getKYCProfile(ctx, customerID)
	if err == nil && profile != nil {
		status.HasTraditionalKYC = true
		status.TraditionalKYCLevel = profile.KYCLevel.String()
		status.TraditionalRiskRating = profile.RiskRating.String()
	}

	// Check DID-based identity
	isValid, credential, err := ia.integration.VerifyKYCStatus(ctx, customerID)
	if err == nil && isValid && credential != nil {
		status.HasDIDIdentity = true
		status.DID = fmt.Sprintf("did:desh:%s", customerID)
		status.CredentialID = credential.ID
		
		// Extract KYC level from credential
		if level, err := ia.integration.GetKYCLevel(ctx, customerID); err == nil {
			status.DIDKYCLevel = level
		}
	}

	// Determine effective status
	status.IsVerified = status.HasTraditionalKYC || status.HasDIDIdentity
	
	return status, nil
}

// MigrateKYCToIdentity migrates existing KYC profiles to identity module
func (ia *IdentityAdapter) MigrateKYCToIdentity(ctx sdk.Context) (int, error) {
	// Get all KYC profiles
	profiles := ia.getAllKYCProfiles(ctx)
	
	// Convert to integration format
	integrationProfiles := make([]identitykeeper.TradeFinanceKYCProfile, len(profiles))
	for i, profile := range profiles {
		integrationProfiles[i] = ia.convertToIntegrationProfile(profile)
	}

	// Perform migration
	return ia.integration.MigrateExistingKYC(ctx, integrationProfiles)
}

// Enhanced KYC operations using identity module

// IssueKYCCredential issues a verifiable KYC credential
func (ia *IdentityAdapter) IssueKYCCredential(
	ctx sdk.Context,
	customerID string,
	kycLevel string,
	riskRating string,
) error {
	kycData := map[string]interface{}{
		"customer_id":  customerID,
		"kyc_level":    kycLevel,
		"risk_rating":  riskRating,
		"issuer":       "tradefinance",
		"verified_at":  ctx.BlockTime(),
	}

	_, err := ia.integration.CreateKYCIdentity(ctx, customerID, kycData)
	return err
}

// RequestSelectiveDisclosure requests specific KYC attributes from customer
func (ia *IdentityAdapter) RequestSelectiveDisclosure(
	ctx sdk.Context,
	customerID string,
	requiredAttributes []string,
) (*SelectiveDisclosureResponse, error) {
	did := fmt.Sprintf("did:desh:%s", customerID)
	
	// Get customer's credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	response := &SelectiveDisclosureResponse{
		CustomerID: customerID,
		DID:        did,
		Attributes: make(map[string]interface{}),
	}

	// Extract requested attributes from credentials
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Extract attributes from credential subject
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			for _, attr := range requiredAttributes {
				if value, exists := subject[attr]; exists {
					response.Attributes[attr] = value
				}
			}
		}
	}

	return response, nil
}

// Helper methods

func (ia *IdentityAdapter) convertKYCDataForIdentity(
	result *KYCAssessmentResult,
	submission CustomerSubmission,
) map[string]interface{} {
	return map[string]interface{}{
		"customer_id":     result.CustomerID,
		"kyc_level":       result.RecommendedKYCLevel.String(),
		"risk_rating":     result.RiskAssessment.OverallRiskRating.String(),
		"customer_type":   submission.PersonalInfo.FullName,
		"doc_verification": map[string]interface{}{
			"status":          result.DocumentVerification.Status,
			"verified_count":  len(result.DocumentVerification.VerifiedDocs),
			"completion_rate": result.DocumentVerification.CompletionRate,
		},
		"identity_verification": map[string]interface{}{
			"status":      result.IdentityVerification.Status,
			"match_score": result.IdentityVerification.MatchScore,
		},
		"pep_status": result.PEPScreening.IsPEP,
		"sanctions_matches": len(result.SanctionsScreening.Matches),
		"assessment_id":     result.AssessmentID,
		"assessment_date":   result.CompletionTime,
	}
}

func (ia *IdentityAdapter) getAllKYCProfiles(ctx sdk.Context) []*KYCProfile {
	// Implementation to retrieve all KYC profiles from store
	// This is a placeholder - actual implementation would iterate through store
	return []*KYCProfile{}
}

func (ia *IdentityAdapter) convertToIntegrationProfile(profile *KYCProfile) identitykeeper.TradeFinanceKYCProfile {
	return identitykeeper.TradeFinanceKYCProfile{
		CustomerID:     profile.CustomerID,
		KYCLevel:       profile.KYCLevel.String(),
		RiskRating:     profile.RiskRating.String(),
		CustomerType:   fmt.Sprintf("%d", profile.CustomerType),
		LastReviewDate: profile.LastReviewDate,
		PersonalInfo: map[string]interface{}{
			"full_name":    profile.PersonalInfo.FullName,
			"date_of_birth": profile.PersonalInfo.DateOfBirth,
			"nationality":  profile.PersonalInfo.Nationality,
		},
		BusinessInfo: map[string]interface{}{
			"legal_name": profile.BusinessInfo.LegalName,
			"reg_number": profile.BusinessInfo.RegistrationNumber,
		},
		Documents: make([]interface{}, len(profile.Documents)),
		PEPStatus: map[string]interface{}{
			"is_pep": profile.PEPStatus.IsPEP,
			"type":   profile.PEPStatus.PEPType,
		},
		SanctionsScreening: map[string]interface{}{
			"screened": true,
			"matches":  0,
		},
	}
}

// New types for enhanced identity features
type CustomerIdentityStatus struct {
	CustomerID            string
	HasTraditionalKYC     bool
	HasDIDIdentity        bool
	IsVerified            bool
	DID                   string
	CredentialID          string
	TraditionalKYCLevel   string
	TraditionalRiskRating string
	DIDKYCLevel           string
}

type SelectiveDisclosureResponse struct {
	CustomerID string
	DID        string
	Attributes map[string]interface{}
}

// Extend KYCAssessmentResult to include DID
type KYCAssessmentResultWithDID struct {
	*KYCAssessmentResult
	IdentityDID string `json:"identity_did,omitempty"`
}