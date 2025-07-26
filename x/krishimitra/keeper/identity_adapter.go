package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
	"github.com/deshchain/deshchain/x/krishimitra/types"
)

// IdentityAdapter provides identity-based verification for KrishiMitra
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

// VerifyFarmerIdentity verifies farmer identity using the identity module
func (ia *IdentityAdapter) VerifyFarmerIdentity(
	ctx sdk.Context,
	farmerAddress sdk.AccAddress,
	pinCode string,
) (*FarmerIdentityStatus, error) {
	status := &FarmerIdentityStatus{
		Address:            farmerAddress.String(),
		HasIdentity:       false,
		IsKYCVerified:     false,
		KYCLevel:          "none",
		HasLandRecords:    false,
		PINCode:           pinCode,
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", farmerAddress.String())
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
			if credType == "KYCCredential" || credType == "FarmerKYC" {
				// Check if credential is valid
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					// Extract KYC level and farmer details
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if aadhaar, ok := subject["aadhaar_hash"].(string); ok {
							status.AadhaarHash = aadhaar
						}
						if kcc, ok := subject["kisan_credit_card"].(string); ok {
							status.KisanCreditCard = kcc
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

	// Check for land record credentials
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "LandRecordCredential" || credType == "LandOwnershipCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasLandRecords = true
					status.LandCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if area, ok := subject["total_land_area"].(float64); ok {
							status.TotalLandArea = sdk.NewDecFromBigIntWithPrec(sdk.NewInt(int64(area*100)), 2)
						}
						if village, ok := subject["village_code"].(string); ok {
							status.VillageCode = village
						}
						if crops, ok := subject["registered_crops"].([]interface{}); ok {
							status.RegisteredCrops = make([]string, len(crops))
							for i, crop := range crops {
								if cropStr, ok := crop.(string); ok {
									status.RegisteredCrops[i] = cropStr
								}
							}
						}
					}
					break
				}
			}
		}
	}

	// Check for crop insurance credentials
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "CropInsuranceCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasCropInsurance = true
					status.InsuranceCredentialID = credID
					break
				}
			}
		}
	}

	return status, nil
}

// CreateFarmerCredential creates a farmer credential
func (ia *IdentityAdapter) CreateFarmerCredential(
	ctx sdk.Context,
	farmerAddress sdk.AccAddress,
	farmerData FarmerData,
) error {
	// Ensure farmer has identity
	did := fmt.Sprintf("did:desh:%s", farmerAddress.String())
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		// Create basic identity
		identity := identitytypes.Identity{
			Did:        did,
			Controller: farmerAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
			Metadata: map[string]string{
				"source": "krishimitra",
				"type":   "farmer",
			},
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Create farmer profile credential
	credID := fmt.Sprintf("vc:farmer:%s:%d", farmerAddress.String(), ctx.BlockTime().Unix())
	
	expiryDate := ctx.BlockTime().AddDate(3, 0, 0) // Valid for 3 years
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/farmer/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "FarmerProfileCredential"},
		Issuer: "did:desh:krishimitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                     did,
			"name":                   farmerData.Name,
			"pincode":                farmerData.PINCode,
			"village_code":           farmerData.VillageCode,
			"district":               farmerData.District,
			"state":                  farmerData.State,
			"aadhaar_hash":           farmerData.AadhaarHash,
			"kisan_credit_card":      farmerData.KisanCreditCard,
			"total_land_area":        farmerData.TotalLandArea.String(),
			"cultivable_land":        farmerData.CultivableLand.String(),
			"crop_types":             farmerData.CropTypes,
			"is_women_farmer":        farmerData.IsWomenFarmer,
			"is_small_farmer":        farmerData.IsSmallFarmer,
			"credit_score":           farmerData.CreditScore,
			"active_loans":           farmerData.ActiveLoans,
			"repaid_loans":           farmerData.RepaidLoans,
			"defaulted_loans":        farmerData.DefaultedLoans,
			"verification_status":    "verified",
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:krishimitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign properly
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// CreateLandRecordCredential creates a land ownership credential
func (ia *IdentityAdapter) CreateLandRecordCredential(
	ctx sdk.Context,
	farmerAddress sdk.AccAddress,
	landData LandRecordData,
) error {
	did := fmt.Sprintf("did:desh:%s", farmerAddress.String())
	
	// Ensure identity exists
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		return fmt.Errorf("identity not found for farmer")
	}

	credID := fmt.Sprintf("vc:landrecord:%s:%d", farmerAddress.String(), ctx.BlockTime().Unix())
	expiryDate := ctx.BlockTime().AddDate(5, 0, 0) // Valid for 5 years
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/landrecord/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "LandRecordCredential"},
		Issuer: "did:desh:revenue-department",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                  did,
			"khata_number":        landData.KhataNumber,
			"survey_numbers":      landData.SurveyNumbers,
			"total_land_area":     landData.TotalArea.String(),
			"cultivable_area":     landData.CultivableArea.String(),
			"irrigated_area":      landData.IrrigatedArea.String(),
			"soil_type":           landData.SoilType,
			"water_source":        landData.WaterSource,
			"registered_crops":    landData.RegisteredCrops,
			"ownership_type":      landData.OwnershipType,
			"village_code":        landData.VillageCode,
			"tehsil":              landData.Tehsil,
			"district":            landData.District,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:revenue-department#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature",
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// CreateAgriculturalLoanCredential creates a credential for an approved agricultural loan
func (ia *IdentityAdapter) CreateAgriculturalLoanCredential(
	ctx sdk.Context,
	loan types.Loan,
) error {
	did := fmt.Sprintf("did:desh:%s", loan.Borrower)
	
	// Ensure identity exists
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		return fmt.Errorf("identity not found for borrower")
	}

	credID := fmt.Sprintf("vc:agriloan:%s:%d", loan.ID, ctx.BlockTime().Unix())
	expiryDate := loan.MaturityDate.AddDate(5, 0, 0) // Valid for 5 years after maturity
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/agriculturalloan/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "AgriculturalLoanCredential"},
		Issuer: "did:desh:krishimitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                    did,
			"loan_id":               loan.ID,
			"amount":                loan.Amount.String(),
			"interest_rate":         loan.InterestRate.String(),
			"purpose":               loan.Purpose,
			"crop_type":             loan.CropType,
			"land_area":             loan.LandArea.String(),
			"expected_yield":        loan.ExpectedYield.String(),
			"term_months":           loan.Term,
			"disbursed_at":          loan.DisbursedAt,
			"maturity_date":         loan.MaturityDate,
			"status":                loan.Status,
			"village_code":          loan.VillageCode,
			"pin_code":              loan.PINCode,
			"festival_bonus":        loan.FestivalBonus,
			"insurance_required":    loan.InsuranceRequired,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:krishimitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature",
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// VerifyLandOwnership verifies if a farmer owns sufficient land for the loan
func (ia *IdentityAdapter) VerifyLandOwnership(
	ctx sdk.Context,
	farmerAddress sdk.AccAddress,
	requiredArea sdk.Dec,
) (bool, error) {
	did := fmt.Sprintf("did:desh:%s", farmerAddress.String())
	
	// Get land credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	totalArea := sdk.ZeroDec()
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's a land record credential
		isLandRecord := false
		for _, credType := range cred.Type {
			if credType == "LandRecordCredential" {
				isLandRecord = true
				break
			}
		}
		
		if !isLandRecord || !ia.isCredentialValid(ctx, cred) {
			continue
		}
		
		// Extract land area
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if areaStr, ok := subject["total_land_area"].(string); ok {
				area, err := sdk.NewDecFromStr(areaStr)
				if err == nil {
					totalArea = totalArea.Add(area)
				}
			}
		}
	}
	
	return totalArea.GTE(requiredArea), nil
}

// GetFarmerCredentials retrieves all farmer-related credentials
func (ia *IdentityAdapter) GetFarmerCredentials(
	ctx sdk.Context,
	farmerAddress sdk.AccAddress,
) ([]FarmerCredential, error) {
	did := fmt.Sprintf("did:desh:%s", farmerAddress.String())
	
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	farmerCreds := []FarmerCredential{}
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's a farmer-related credential
		isFarmer := false
		for _, credType := range cred.Type {
			if credType == "FarmerProfileCredential" || 
			   credType == "LandRecordCredential" ||
			   credType == "AgriculturalLoanCredential" ||
			   credType == "CropInsuranceCredential" ||
			   credType == "FarmerKYC" {
				isFarmer = true
				break
			}
		}
		
		if !isFarmer {
			continue
		}
		
		farmerCred := FarmerCredential{
			CredentialID: credID,
			Type:         cred.Type,
			IssuanceDate: cred.IssuanceDate,
			IsValid:      ia.isCredentialValid(ctx, cred),
		}
		
		// Extract relevant info
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if loanID, ok := subject["loan_id"].(string); ok {
				farmerCred.LoanID = loanID
			}
			if status, ok := subject["status"].(string); ok {
				farmerCred.Status = status
			}
			if cropType, ok := subject["crop_type"].(string); ok {
				farmerCred.CropType = cropType
			}
			if landArea, ok := subject["land_area"].(string); ok {
				farmerCred.LandArea = landArea
			}
		}
		
		farmerCreds = append(farmerCreds, farmerCred)
	}
	
	return farmerCreds, nil
}

// MigrateExistingFarmers migrates existing farmers to use identity system
func (ia *IdentityAdapter) MigrateExistingFarmers(ctx sdk.Context) (int, error) {
	migrated := 0
	
	// Note: Since FarmerProfile type is not defined in the current codebase,
	// this is a placeholder for the migration logic.
	// In production, this would iterate through all farmer profiles and create credentials.
	
	// Migrate existing loans
	loans := ia.keeper.GetAllLoans(ctx)
	for _, loan := range loans {
		// Check if borrower has identity
		did := fmt.Sprintf("did:desh:%s", loan.Borrower)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
			// Create basic identity for loan borrower
			identity := identitytypes.Identity{
				Did:        did,
				Controller: loan.Borrower,
				Status:     identitytypes.IdentityStatus_ACTIVE,
				CreatedAt:  ctx.BlockTime(),
				UpdatedAt:  ctx.BlockTime(),
				Metadata: map[string]string{
					"source": "krishimitra",
					"type":   "farmer",
				},
			}
			ia.identityKeeper.SetIdentity(ctx, identity)
			migrated++
		}
		
		// Create loan credential
		if err := ia.CreateAgriculturalLoanCredential(ctx, loan); err != nil {
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

type FarmerIdentityStatus struct {
	Address                string
	HasIdentity           bool
	DID                   string
	IdentityStatus        string
	IsKYCVerified         bool
	KYCLevel              string
	KYCCredentialID       string
	HasLandRecords        bool
	LandCredentialID      string
	HasCropInsurance      bool
	InsuranceCredentialID string
	PINCode               string
	VillageCode           string
	AadhaarHash           string
	KisanCreditCard       string
	TotalLandArea         sdk.Dec
	RegisteredCrops       []string
}

type FarmerData struct {
	Name             string
	PINCode          string
	VillageCode      string
	District         string
	State            string
	AadhaarHash      string
	KisanCreditCard  string
	TotalLandArea    sdk.Dec
	CultivableLand   sdk.Dec
	CropTypes        []string
	IsWomenFarmer    bool
	IsSmallFarmer    bool
	CreditScore      int32
	ActiveLoans      int32
	RepaidLoans      int32
	DefaultedLoans   int32
}

type LandRecordData struct {
	KhataNumber     string
	SurveyNumbers   []string
	TotalArea       sdk.Dec
	CultivableArea  sdk.Dec
	IrrigatedArea   sdk.Dec
	SoilType        string
	WaterSource     string
	RegisteredCrops []string
	OwnershipType   string
	VillageCode     string
	Tehsil          string
	District        string
}

type FarmerCredential struct {
	CredentialID string
	Type         []string
	LoanID       string
	Status       string
	CropType     string
	LandArea     string
	IssuanceDate time.Time
	IsValid      bool
}

// FarmerProfile represents a farmer's profile (placeholder for missing type)
type FarmerProfile struct {
	Address            string
	Name               string
	PINCode            string
	LandSize           sdk.Dec
	TotalLandArea      sdk.Dec
	ActiveLoans        int32
	TotalLoans         int32
	DefaultedLoans     int32
	RepaidLoans        int32
	CreditScore        int32
	IsWomenFarmer      bool
	VerificationStatus string
}

// Placeholder methods for FarmerProfile operations
func (k Keeper) GetFarmerProfile(ctx sdk.Context, farmerAddress string) (FarmerProfile, bool) {
	// This is a placeholder - in production, this would retrieve from store
	return FarmerProfile{}, false
}

func (k Keeper) SetFarmerProfile(ctx sdk.Context, profile FarmerProfile) {
	// This is a placeholder - in production, this would store the profile
}

// GetAllLoans retrieves all loans
func (k Keeper) GetAllLoans(ctx sdk.Context) []types.Loan {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.LoanKeyPrefix, nil)
	defer iterator.Close()

	var loans []types.Loan
	for ; iterator.Valid(); iterator.Next() {
		var loan types.Loan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		loans = append(loans, loan)
	}
	return loans
}