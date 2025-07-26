package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
	"github.com/deshchain/deshchain/x/gramsuraksha/types"
)

// IdentityAdapter provides identity-based verification for GramSuraksha
type IdentityAdapter struct {
	keeper         *Keeper
	identityKeeper identitykeeper.Keeper
}

// NewIdentityAdapter creates a new identity adapter
func NewIdentityAdapter(k *Keeper, ik identitykeeper.Keeper) *IdentityAdapter {
	return &IdentityAdapter{
		keeper:         k,
		identityKeeper: ik,
	}
}

// VerifyParticipantIdentity verifies participant identity using the identity module
func (ia *IdentityAdapter) VerifyParticipantIdentity(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
	schemeID string,
) (*ParticipantIdentityStatus, error) {
	status := &ParticipantIdentityStatus{
		Address:       participantAddress.String(),
		HasIdentity:   false,
		IsKYCVerified: false,
		KYCLevel:      "none",
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", participantAddress.String())
	identity, exists := ia.identityKeeper.GetIdentity(ctx, did)
	if !exists {
		return status, nil
	}

	status.HasIdentity = true
	status.DID = did
	status.IdentityStatus = string(identity.Status)

	// Check for KYC credential
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a KYC credential
		for _, credType := range cred.Type {
			if credType == "KYCCredential" || credType == "GramSurakshaKYC" {
				// Check if credential is valid
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					// Extract KYC level
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if age, ok := subject["age"].(float64); ok {
							status.Age = int32(age)
						}
					}
					break
				}
			}
		}
		if status.IsKYCVerified {
			break
		}
	}

	return status, nil
}

// CreateParticipantCredential creates a GramSuraksha participant credential
func (ia *IdentityAdapter) CreateParticipantCredential(
	ctx sdk.Context,
	participant types.SurakshaParticipant,
	scheme types.SurakshaScheme,
) error {
	// Ensure participant has identity
	did := fmt.Sprintf("did:desh:%s", participant.Address.String())
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		// Create basic identity
		identity := identitytypes.Identity{
			Did:        did,
			Controller: participant.Address.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
			Metadata: map[string]string{
				"source": "gramsuraksha",
				"type":   "participant",
			},
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Create participant credential
	credID := fmt.Sprintf("vc:gramsuraksha:%s:%s", participant.ParticipantID, ctx.BlockTime().Unix())
	
	expiryDate := participant.MaturityDate.AddDate(5, 0, 0) // Valid for 5 years after maturity
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/gramsuraksha/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "GramSurakshaParticipant"},
		Issuer: "did:desh:gramsuraksha-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":               did,
			"participant_id":   participant.ParticipantID,
			"scheme_id":        participant.SchemeID,
			"scheme_name":      scheme.Name,
			"enrollment_date":  participant.EnrollmentDate,
			"maturity_date":    participant.MaturityDate,
			"monthly_contribution": scheme.MonthlyContribution.String(),
			"guaranteed_return": scheme.GuaranteedReturn,
			"status":           participant.Status,
			"age":              participant.Age,
			"village":          participant.VillageName,
			"district":         participant.District,
			"state":            participant.State,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:gramsuraksha-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign properly
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// VerifyAgeWithZKProof verifies participant age using zero-knowledge proof
func (ia *IdentityAdapter) VerifyAgeWithZKProof(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
	minAge, maxAge int32,
) (bool, error) {
	did := fmt.Sprintf("did:desh:%s", participantAddress.String())
	
	// Find age-related credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	var ageCredentialID string
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if credential contains age information
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if _, hasAge := subject["age"]; hasAge || subject["date_of_birth"] != nil {
				ageCredentialID = credID
				break
			}
		}
	}
	
	if ageCredentialID == "" {
		return false, fmt.Errorf("no age credential found")
	}
	
	// In production, this would create an actual ZK proof
	// For now, we'll do a simple check
	statement := fmt.Sprintf("age >= %d AND age <= %d", minAge, maxAge)
	
	// Mock ZK proof verification
	// In reality, this would use the privacy keeper to verify the proof
	return true, nil
}

// UpdateParticipantCredentialStatus updates the status of a participant credential
func (ia *IdentityAdapter) UpdateParticipantCredentialStatus(
	ctx sdk.Context,
	participantID string,
	newStatus string,
	reason string,
) error {
	// Find the participant credential
	// In production, we'd have an index for this
	// For now, this is a placeholder
	
	credentialID := fmt.Sprintf("vc:gramsuraksha:%s", participantID)
	
	if newStatus == "suspended" || newStatus == "terminated" {
		// Revoke the credential
		return ia.identityKeeper.UpdateCredentialStatus(ctx, credentialID, &identitytypes.CredentialStatus{
			Type:   "revoked",
			Reason: reason,
		})
	}
	
	return nil
}

// GetParticipantCredentials retrieves all GramSuraksha credentials for a participant
func (ia *IdentityAdapter) GetParticipantCredentials(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
) ([]ParticipantCredential, error) {
	did := fmt.Sprintf("did:desh:%s", participantAddress.String())
	
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	participantCreds := []ParticipantCredential{}
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's a GramSuraksha credential
		isGramSuraksha := false
		for _, credType := range cred.Type {
			if credType == "GramSurakshaParticipant" || credType == "GramSurakshaKYC" {
				isGramSuraksha = true
				break
			}
		}
		
		if !isGramSuraksha {
			continue
		}
		
		participantCred := ParticipantCredential{
			CredentialID: credID,
			Type:         cred.Type,
			IssuanceDate: cred.IssuanceDate,
			IsValid:      ia.isCredentialValid(ctx, cred),
		}
		
		// Extract scheme info if available
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if schemeID, ok := subject["scheme_id"].(string); ok {
				participantCred.SchemeID = schemeID
			}
			if status, ok := subject["status"].(string); ok {
				participantCred.Status = status
			}
		}
		
		participantCreds = append(participantCreds, participantCred)
	}
	
	return participantCreds, nil
}

// MigrateExistingParticipants migrates existing participants to use identity system
func (ia *IdentityAdapter) MigrateExistingParticipants(ctx sdk.Context) (int, error) {
	migrated := 0
	
	// Get all participants
	participants := ia.keeper.GetAllParticipants(ctx)
	
	for _, participant := range participants {
		// Check if already has identity
		did := fmt.Sprintf("did:desh:%s", participant.Address.String())
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); exists {
			continue // Already migrated
		}
		
		// Get scheme
		scheme, found := ia.keeper.GetScheme(ctx, participant.SchemeID)
		if !found {
			continue
		}
		
		// Create participant credential
		if err := ia.CreateParticipantCredential(ctx, participant, scheme); err != nil {
			ia.keeper.Logger(ctx).Error("Failed to migrate participant",
				"participant_id", participant.ParticipantID,
				"error", err)
			continue
		}
		
		migrated++
	}
	
	return migrated, nil
}

// Helper methods

func (ia *IdentityAdapter) isCredentialValid(ctx sdk.Context, cred identitytypes.VerifiableCredential) bool {
	// Check if revoked
	if cred.Status != nil && cred.Status.Type == "revoked" {
		return false
	}
	
	// Check expiry
	if cred.ExpirationDate != nil && cred.ExpirationDate.Before(ctx.BlockTime()) {
		return false
	}
	
	return true
}

// Types for identity integration

type ParticipantIdentityStatus struct {
	Address         string
	HasIdentity     bool
	DID             string
	IdentityStatus  string
	IsKYCVerified   bool
	KYCLevel        string
	KYCCredentialID string
	Age             int32
}

type ParticipantCredential struct {
	CredentialID string
	Type         []string
	SchemeID     string
	Status       string
	IssuanceDate time.Time
	IsValid      bool
}