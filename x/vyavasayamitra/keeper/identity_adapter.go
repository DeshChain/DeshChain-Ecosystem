package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
	"github.com/deshchain/x/vyavasayamitra/types"
)

// IdentityAdapter provides identity-based verification for VyavasayaMitra
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

// VerifyBusinessIdentity verifies business identity using the identity module
func (ia *IdentityAdapter) VerifyBusinessIdentity(
	ctx sdk.Context,
	businessAddress sdk.AccAddress,
	businessType string,
) (*BusinessIdentityStatus, error) {
	status := &BusinessIdentityStatus{
		Address:              businessAddress.String(),
		HasIdentity:         false,
		IsKYCVerified:       false,
		KYCLevel:           "none",
		HasComplianceDocs:  false,
		BusinessType:       businessType,
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", businessAddress.String())
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

		// Check if it's a business KYC credential
		for _, credType := range cred.Type {
			if credType == "KYCCredential" || credType == "BusinessKYC" {
				// Check if credential is valid
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					// Extract KYC level and business details
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if gst, ok := subject["gst_number"].(string); ok {
							status.GSTNumber = gst
						}
						if pan, ok := subject["pan_number"].(string); ok {
							status.PANNumber = pan
						}
						if udyam, ok := subject["udyam_number"].(string); ok {
							status.UdyamNumber = udyam
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

	// Check for compliance documents
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "ComplianceCredential" || credType == "BusinessComplianceCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasComplianceDocs = true
					status.ComplianceCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if licenses, ok := subject["licenses"].([]interface{}); ok {
							status.BusinessLicenses = make([]string, len(licenses))
							for i, license := range licenses {
								if licStr, ok := license.(string); ok {
									status.BusinessLicenses[i] = licStr
								}
							}
						}
					}
					break
				}
			}
		}
	}

	// Check for financial credentials
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "FinancialCredential" || credType == "BusinessFinancialCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasFinancialDocs = true
					status.FinancialCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if revenue, ok := subject["annual_revenue"].(float64); ok {
							status.VerifiedRevenue = sdk.NewDec(int64(revenue))
						}
						if score, ok := subject["credit_score"].(float64); ok {
							status.CreditScore = int32(score)
						}
					}
					break
				}
			}
		}
	}

	return status, nil
}

// CreateBusinessCredential creates a business credential
func (ia *IdentityAdapter) CreateBusinessCredential(
	ctx sdk.Context,
	profile types.BusinessProfile,
) error {
	// Ensure business has identity
	did := fmt.Sprintf("did:desh:%s", profile.DhanPataAddress)
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		// Create basic identity
		identity := identitytypes.Identity{
			Did:        did,
			Controller: profile.DhanPataAddress,
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
			Metadata: map[string]string{
				"source": "vyavasayamitra",
				"type":   "business",
			},
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Create business profile credential
	credID := fmt.Sprintf("vc:business:%s:%d", profile.ID, ctx.BlockTime().Unix())
	
	expiryDate := ctx.BlockTime().AddDate(2, 0, 0) // Valid for 2 years
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/business/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "BusinessProfileCredential"},
		Issuer: "did:desh:vyavasayamitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                      did,
			"business_id":             profile.ID,
			"business_name":           profile.BusinessName,
			"business_type":           profile.BusinessType,
			"registration_date":       profile.RegistrationDate,
			"gst_number":              profile.GSTNumber,
			"pan_number":              profile.PANNumber,
			"address":                 profile.Address,
			"state":                   profile.State,
			"district":                profile.District,
			"pincode":                 profile.Pincode,
			"annual_revenue":          profile.AnnualRevenue.String(),
			"credit_score":            profile.CreditScore,
			"verification_status":     profile.VerificationStatus,
			"merchant_rating":         profile.MerchantRating.String(),
			"total_loans_availed":     profile.TotalLoansAvailed,
			"active_loans":            profile.ActiveLoans,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:vyavasayamitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign properly
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// CreateBusinessLoanCredential creates a credential for an approved business loan
func (ia *IdentityAdapter) CreateBusinessLoanCredential(
	ctx sdk.Context,
	loan types.BusinessLoan,
) error {
	did := fmt.Sprintf("did:desh:%s", loan.Borrower)
	
	// Ensure identity exists
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		return fmt.Errorf("identity not found for borrower")
	}

	credID := fmt.Sprintf("vc:bizloan:%s:%d", loan.ID, ctx.BlockTime().Unix())
	expiryDate := loan.EndDate.AddDate(5, 0, 0) // Valid for 5 years after loan end
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/businessloan/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "BusinessLoanCredential"},
		Issuer: "did:desh:vyavasayamitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                    did,
			"loan_id":               loan.ID,
			"business_id":           loan.BusinessID,
			"amount":                loan.Amount.String(),
			"interest_rate":         loan.InterestRate.String(),
			"purpose":               loan.Purpose,
			"tenure_months":         loan.TenureMonths,
			"start_date":            loan.StartDate,
			"end_date":              loan.EndDate,
			"status":                loan.Status,
			"repayment_frequency":   loan.RepaymentFrequency,
			"disbursed_amount":      loan.DisbursedAmount.String(),
			"collateral_type":       loan.CollateralType,
			"festival_discount":     loan.FestivalDiscount.String(),
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:vyavasayamitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature",
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// VerifyBusinessCompliance verifies if a business has required compliance documents
func (ia *IdentityAdapter) VerifyBusinessCompliance(
	ctx sdk.Context,
	businessAddress sdk.AccAddress,
	requiredLicenses []string,
) (bool, error) {
	did := fmt.Sprintf("did:desh:%s", businessAddress.String())
	
	// Get compliance credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's a compliance credential
		isCompliance := false
		for _, credType := range cred.Type {
			if credType == "ComplianceCredential" || credType == "BusinessComplianceCredential" {
				isCompliance = true
				break
			}
		}
		
		if !isCompliance || !ia.isCredentialValid(ctx, cred) {
			continue
		}
		
		// Extract licenses
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if licenses, ok := subject["licenses"].([]interface{}); ok {
				// Check if all required licenses are present
				hasAllLicenses := true
				for _, required := range requiredLicenses {
					found := false
					for _, license := range licenses {
						if licStr, ok := license.(string); ok && licStr == required {
							found = true
							break
						}
					}
					if !found {
						hasAllLicenses = false
						break
					}
				}
				if hasAllLicenses {
					return true, nil
				}
			}
		}
	}
	
	return false, nil
}

// GetBusinessCredentials retrieves all business-related credentials
func (ia *IdentityAdapter) GetBusinessCredentials(
	ctx sdk.Context,
	businessAddress sdk.AccAddress,
) ([]BusinessCredential, error) {
	did := fmt.Sprintf("did:desh:%s", businessAddress.String())
	
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	businessCreds := []BusinessCredential{}
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's a business-related credential
		isBusiness := false
		for _, credType := range cred.Type {
			if credType == "BusinessProfileCredential" || 
			   credType == "BusinessLoanCredential" ||
			   credType == "BusinessKYC" ||
			   credType == "ComplianceCredential" ||
			   credType == "FinancialCredential" {
				isBusiness = true
				break
			}
		}
		
		if !isBusiness {
			continue
		}
		
		businessCred := BusinessCredential{
			CredentialID: credID,
			Type:         cred.Type,
			IssuanceDate: cred.IssuanceDate,
			IsValid:      ia.isCredentialValid(ctx, cred),
		}
		
		// Extract relevant info
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if loanID, ok := subject["loan_id"].(string); ok {
				businessCred.LoanID = loanID
			}
			if status, ok := subject["status"].(string); ok {
				businessCred.Status = status
			}
			if businessName, ok := subject["business_name"].(string); ok {
				businessCred.BusinessName = businessName
			}
		}
		
		businessCreds = append(businessCreds, businessCred)
	}
	
	return businessCreds, nil
}

// MigrateExistingBusinesses migrates existing businesses to use identity system
func (ia *IdentityAdapter) MigrateExistingBusinesses(ctx sdk.Context) (int, error) {
	migrated := 0
	
	// Get all business profiles
	profiles := ia.keeper.GetAllBusinessProfiles(ctx)
	
	for _, profile := range profiles {
		// Check if already has identity
		did := fmt.Sprintf("did:desh:%s", profile.DhanPataAddress)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); exists {
			continue // Already migrated
		}
		
		// Create business credential
		if err := ia.CreateBusinessCredential(ctx, profile); err != nil {
			ia.keeper.Logger(ctx).Error("Failed to migrate business",
				"business_id", profile.ID,
				"error", err)
			continue
		}
		
		migrated++
	}
	
	// Migrate existing loans
	loans := ia.keeper.GetAllBusinessLoans(ctx)
	for _, loan := range loans {
		// Check if borrower has identity
		did := fmt.Sprintf("did:desh:%s", loan.Borrower)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
			continue // Skip if no identity
		}
		
		// Create loan credential
		if err := ia.CreateBusinessLoanCredential(ctx, loan); err != nil {
			ia.keeper.Logger(ctx).Error("Failed to migrate loan credential",
				"loan_id", loan.ID,
				"error", err)
			continue
		}
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

type BusinessIdentityStatus struct {
	Address                 string
	HasIdentity            bool
	DID                    string
	IdentityStatus         string
	IsKYCVerified          bool
	KYCLevel               string
	KYCCredentialID        string
	HasComplianceDocs      bool
	ComplianceCredentialID string
	HasFinancialDocs       bool
	FinancialCredentialID  string
	BusinessType           string
	GSTNumber              string
	PANNumber              string
	UdyamNumber            string
	BusinessLicenses       []string
	VerifiedRevenue        sdk.Dec
	CreditScore            int32
}

type BusinessCredential struct {
	CredentialID string
	Type         []string
	LoanID       string
	Status       string
	BusinessName string
	IssuanceDate time.Time
	IsValid      bool
}