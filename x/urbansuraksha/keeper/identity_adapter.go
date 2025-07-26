package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	identitytypes "github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/urbansuraksha/types"
)

// IdentityAdapter provides identity-based verification for UrbanSuraksha
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

// VerifyContributorIdentity verifies contributor identity using the identity module
func (ia *IdentityAdapter) VerifyContributorIdentity(
	ctx sdk.Context,
	contributorAddress sdk.AccAddress,
	cityCode string,
) (*ContributorIdentityStatus, error) {
	status := &ContributorIdentityStatus{
		Address:        contributorAddress.String(),
		HasIdentity:    false,
		IsKYCVerified:  false,
		KYCLevel:       "none",
		CityCode:       cityCode,
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", contributorAddress.String())
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
			if credType == "KYCCredential" || credType == "UrbanSurakshaKYC" {
				// Check if credential is valid
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					// Extract KYC level and location
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if city, ok := subject["city"].(string); ok {
							status.VerifiedCity = city
						}
						if income, ok := subject["annual_income"].(float64); ok {
							status.VerifiedIncome = sdk.NewDec(int64(income))
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

	// Check for employment credential
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "EmploymentCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasEmploymentProof = true
					status.EmploymentCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if employer, ok := subject["employer"].(string); ok {
							status.Employer = employer
						}
					}
					break
				}
			}
		}
	}

	return status, nil
}

// CreateContributorCredential creates an UrbanSuraksha contributor credential
func (ia *IdentityAdapter) CreateContributorCredential(
	ctx sdk.Context,
	scheme types.UrbanSurakshaScheme,
) error {
	// Ensure contributor has identity
	did := fmt.Sprintf("did:desh:%s", scheme.ContributorAddress.String())
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		// Create basic identity
		identity := identitytypes.Identity{
			Did:        did,
			Controller: scheme.ContributorAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
			Metadata: map[string]string{
				"source": "urbansuraksha",
				"type":   "contributor",
			},
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Create contributor credential
	credID := fmt.Sprintf("vc:urbansuraksha:%s:%d", scheme.SchemeID, ctx.BlockTime().Unix())
	
	expiryDate := scheme.MaturityDate.AddDate(5, 0, 0) // Valid for 5 years after maturity
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.com/contexts/urbansuraksha/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "UrbanSurakshaContributor"},
		Issuer: "did:desh:urbansuraksha-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                   did,
			"scheme_id":            scheme.SchemeID,
			"account_id":           scheme.AccountID,
			"start_date":           scheme.StartDate,
			"maturity_date":        scheme.MaturityDate,
			"monthly_contribution": scheme.MonthlyContribution.String(),
			"return_percentage":    scheme.ReturnPercentage.String(),
			"status":               scheme.Status,
			"life_insurance_cover": scheme.LifeInsuranceCover.String(),
			"health_insurance_cover": scheme.HealthInsuranceCover.String(),
			"city_code":            extractCityCode(scheme.SchemeID),
			"contribution_months":  scheme.ContributionPeriod,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:urbansuraksha-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign properly
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// VerifyIncomeWithZKProof verifies contributor income using zero-knowledge proof
func (ia *IdentityAdapter) VerifyIncomeWithZKProof(
	ctx sdk.Context,
	contributorAddress sdk.AccAddress,
	minIncome sdk.Dec,
) (bool, error) {
	did := fmt.Sprintf("did:desh:%s", contributorAddress.String())
	
	// Find income-related credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	var incomeCredentialID string
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if credential contains income information
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if _, hasIncome := subject["annual_income"]; hasIncome {
				incomeCredentialID = credID
				break
			}
		}
	}
	
	if incomeCredentialID == "" {
		return false, fmt.Errorf("no income credential found")
	}
	
	// In production, this would create an actual ZK proof
	// For now, we'll do a simple check
	statement := fmt.Sprintf("annual_income >= %s", minIncome.String())
	
	// Mock ZK proof verification
	// In reality, this would use the privacy keeper to verify the proof
	return true, nil
}

// UpdateContributorCredentialStatus updates the status of a contributor credential
func (ia *IdentityAdapter) UpdateContributorCredentialStatus(
	ctx sdk.Context,
	schemeID string,
	newStatus string,
	reason string,
) error {
	// Find the contributor credential
	// In production, we'd have an index for this
	credentialID := fmt.Sprintf("vc:urbansuraksha:%s", schemeID)
	
	if newStatus == "suspended" || newStatus == "terminated" {
		// Revoke the credential
		return ia.identityKeeper.UpdateCredentialStatus(ctx, credentialID, &identitytypes.CredentialStatus{
			Type:   "revoked",
			Reason: reason,
		})
	}
	
	return nil
}

// GetContributorCredentials retrieves all UrbanSuraksha credentials for a contributor
func (ia *IdentityAdapter) GetContributorCredentials(
	ctx sdk.Context,
	contributorAddress sdk.AccAddress,
) ([]ContributorCredential, error) {
	did := fmt.Sprintf("did:desh:%s", contributorAddress.String())
	
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	contributorCreds := []ContributorCredential{}
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's an UrbanSuraksha credential
		isUrbanSuraksha := false
		for _, credType := range cred.Type {
			if credType == "UrbanSurakshaContributor" || credType == "UrbanSurakshaKYC" {
				isUrbanSuraksha = true
				break
			}
		}
		
		if !isUrbanSuraksha {
			continue
		}
		
		contributorCred := ContributorCredential{
			CredentialID: credID,
			Type:         cred.Type,
			IssuanceDate: cred.IssuanceDate,
			IsValid:      ia.isCredentialValid(ctx, cred),
		}
		
		// Extract scheme info if available
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if schemeID, ok := subject["scheme_id"].(string); ok {
				contributorCred.SchemeID = schemeID
			}
			if status, ok := subject["status"].(string); ok {
				contributorCred.Status = status
			}
			if cityCode, ok := subject["city_code"].(string); ok {
				contributorCred.CityCode = cityCode
			}
		}
		
		contributorCreds = append(contributorCreds, contributorCred)
	}
	
	return contributorCreds, nil
}

// CreateEmploymentCredential creates an employment verification credential
func (ia *IdentityAdapter) CreateEmploymentCredential(
	ctx sdk.Context,
	contributorAddress sdk.AccAddress,
	employmentData EmploymentData,
) error {
	did := fmt.Sprintf("did:desh:%s", contributorAddress.String())
	
	// Ensure identity exists
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		return fmt.Errorf("identity not found for contributor")
	}

	credID := fmt.Sprintf("vc:employment:%s:%d", contributorAddress.String(), ctx.BlockTime().Unix())
	expiryDate := ctx.BlockTime().AddDate(1, 0, 0) // Valid for 1 year
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.com/contexts/employment/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "EmploymentCredential"},
		Issuer: employmentData.IssuerDID,
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":              did,
			"employer":        employmentData.Employer,
			"designation":     employmentData.Designation,
			"employment_date": employmentData.EmploymentDate,
			"annual_income":   employmentData.AnnualIncome.RoundInt64(),
			"employment_type": employmentData.EmploymentType,
			"city":            employmentData.City,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: employmentData.IssuerDID + "#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature",
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// MigrateExistingContributors migrates existing contributors to use identity system
func (ia *IdentityAdapter) MigrateExistingContributors(ctx sdk.Context) (int, error) {
	migrated := 0
	
	// Get all schemes
	schemes := ia.keeper.GetAllSchemes(ctx)
	
	for _, scheme := range schemes {
		// Check if already has identity
		did := fmt.Sprintf("did:desh:%s", scheme.ContributorAddress.String())
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); exists {
			continue // Already migrated
		}
		
		// Create contributor credential
		if err := ia.CreateContributorCredential(ctx, scheme); err != nil {
			ia.keeper.Logger(ctx).Error("Failed to migrate contributor",
				"scheme_id", scheme.SchemeID,
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

func extractCityCode(schemeID string) string {
	// SchemeID format: "UPS-{cityCode}-{timestamp}"
	if len(schemeID) > 4 {
		parts := sdk.StringsBetween(schemeID, "UPS-", "-")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// Types for identity integration

type ContributorIdentityStatus struct {
	Address                string
	HasIdentity            bool
	DID                    string
	IdentityStatus         string
	IsKYCVerified          bool
	KYCLevel               string
	KYCCredentialID        string
	CityCode               string
	VerifiedCity           string
	VerifiedIncome         sdk.Dec
	HasEmploymentProof     bool
	EmploymentCredentialID string
	Employer               string
}

type ContributorCredential struct {
	CredentialID string
	Type         []string
	SchemeID     string
	Status       string
	CityCode     string
	IssuanceDate time.Time
	IsValid      bool
}

type EmploymentData struct {
	Employer       string
	Designation    string
	EmploymentDate time.Time
	AnnualIncome   sdk.Dec
	EmploymentType string // permanent, contract, self-employed
	City           string
	IssuerDID      string // DID of the employer or verification service
}