package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/kisaanmitra/types"
	identitykeeper "github.com/deshchain/deshchain/x/identity/keeper"
	identitytypes "github.com/deshchain/deshchain/x/identity/types"
)

// IdentityAdapter provides identity integration for KisaanMitra module
type IdentityAdapter struct {
	keeper         *Keeper
	identityKeeper identitykeeper.Keeper
}

// NewIdentityAdapter creates a new identity adapter
func NewIdentityAdapter(keeper *Keeper, identityKeeper identitykeeper.Keeper) *IdentityAdapter {
	return &IdentityAdapter{
		keeper:         keeper,
		identityKeeper: identityKeeper,
	}
}

// BorrowerIdentityStatus represents the identity status of a borrower
type BorrowerIdentityStatus struct {
	HasIdentity      bool
	DID              string
	IsKYCVerified    bool
	KYCLevel         string
	HasLandRecords   bool
	TotalLandArea    sdk.Dec
	RegisteredCrops  []string
	HasCropInsurance bool
	PINCode          string
	VillageCode      string
	ActiveLoans      []string
}

// LandOwnershipStatus represents land ownership verification status
type LandOwnershipStatus struct {
	TotalLandArea     sdk.Dec
	HasSufficientLand bool
	LandCredentials   []string
}

// VerifyBorrowerIdentity verifies a borrower's identity and credentials
func (ia *IdentityAdapter) VerifyBorrowerIdentity(
	ctx sdk.Context,
	borrowerAddress sdk.AccAddress,
	villageCode string,
) (*BorrowerIdentityStatus, error) {
	// Check if identity exists
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddress)
	if !found {
		return &BorrowerIdentityStatus{
			HasIdentity: false,
		}, nil
	}

	status := &BorrowerIdentityStatus{
		HasIdentity: true,
		DID:         identity.Did,
	}

	// Check KYC credentials
	kycCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "KYCCredential")
	for _, cred := range kycCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.IsKYCVerified = true
			if level, ok := cred.CredentialSubject["kyc_level"].(string); ok {
				status.KYCLevel = level
			}
			if pincode, ok := cred.CredentialSubject["pincode"].(string); ok {
				status.PINCode = pincode
			}
			if village, ok := cred.CredentialSubject["village_code"].(string); ok {
				status.VillageCode = village
			}
			break
		}
	}

	// Check farmer credentials
	farmerCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "FarmerCredential")
	for _, cred := range farmerCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if crops, ok := cred.CredentialSubject["registered_crops"].([]string); ok {
				status.RegisteredCrops = crops
			}
			break
		}
	}

	// Check land record credentials
	landCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "LandRecordCredential")
	totalArea := sdk.ZeroDec()
	for _, cred := range landCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.HasLandRecords = true
			if areaStr, ok := cred.CredentialSubject["total_area"].(string); ok {
				area, err := sdk.NewDecFromStr(areaStr)
				if err == nil {
					totalArea = totalArea.Add(area)
				}
			}
		}
	}
	status.TotalLandArea = totalArea

	// Check crop insurance credentials
	insuranceCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "CropInsuranceCredential")
	for _, cred := range insuranceCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			// Check if insurance is still valid
			if validTill, ok := cred.CredentialSubject["valid_till"].(string); ok {
				validTime, err := time.Parse(time.RFC3339, validTill)
				if err == nil && validTime.After(ctx.BlockTime()) {
					status.HasCropInsurance = true
					break
				}
			}
		}
	}

	// Check active loans
	loanCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "AgriculturalLoanCredential")
	for _, cred := range loanCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if loanID, ok := cred.CredentialSubject["loan_id"].(string); ok {
				// Verify loan is still active
				loan, found := ia.keeper.GetLoan(ctx, loanID)
				if found && loan.Status == types.StatusActive {
					status.ActiveLoans = append(status.ActiveLoans, loanID)
				}
			}
		}
	}

	// Verify village code if specified
	if villageCode != "" && status.VillageCode != "" && status.VillageCode != villageCode {
		return nil, fmt.Errorf("borrower village code mismatch: expected %s, got %s", villageCode, status.VillageCode)
	}

	return status, nil
}

// VerifyLandOwnership verifies if a borrower has sufficient land
func (ia *IdentityAdapter) VerifyLandOwnership(
	ctx sdk.Context,
	borrowerAddress sdk.AccAddress,
	requiredArea sdk.Dec,
) (*LandOwnershipStatus, error) {
	// Get borrower's identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddress)
	if !found {
		return &LandOwnershipStatus{
			TotalLandArea:     sdk.ZeroDec(),
			HasSufficientLand: false,
		}, nil
	}

	// Check land record credentials
	landCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "LandRecordCredential")
	totalArea := sdk.ZeroDec()
	var credentialIDs []string

	for _, cred := range landCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if areaStr, ok := cred.CredentialSubject["total_area"].(string); ok {
				area, err := sdk.NewDecFromStr(areaStr)
				if err == nil {
					totalArea = totalArea.Add(area)
					credentialIDs = append(credentialIDs, cred.Id)
				}
			}
		}
	}

	return &LandOwnershipStatus{
		TotalLandArea:     totalArea,
		HasSufficientLand: totalArea.GTE(requiredArea),
		LandCredentials:   credentialIDs,
	}, nil
}

// CreateFarmerCredential creates a farmer credential for a borrower
func (ia *IdentityAdapter) CreateFarmerCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	borrowerAddress sdk.AccAddress,
	pincode string,
) (string, error) {
	// Get or create identity for borrower
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddress)
	if !found {
		// Create new identity
		did := fmt.Sprintf("did:desh:%s", borrowerAddress.String())
		identity = identitytypes.Identity{
			Did:        did,
			Controller: borrowerAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Get borrower from KisaanMitra
	borrowerID := fmt.Sprintf("borrower_%s", borrowerAddress.String())
	borrower, found := ia.keeper.GetBorrower(ctx, borrowerID)
	if !found {
		return "", fmt.Errorf("borrower not found: %s", borrowerID)
	}

	// Create farmer credential
	credentialSubject := map[string]interface{}{
		"id":               identity.Did,
		"address":          borrowerAddress.String(),
		"borrower_id":      borrower.BorrowerID,
		"name":             borrower.Name,
		"pincode":          pincode,
		"village_code":     borrower.VillageCode,
		"village_name":     borrower.VillageName,
		"registered_crops": []string{borrower.PrimaryMobilization},
		"land_size":        borrower.LandSize.String(),
		"land_ownership":   borrower.LandOwnership,
		"phone":            borrower.Phone,
		"bank_account":     borrower.BankAccount,
		"ifsc_code":        borrower.IFSCCode,
	}

	// Add secondary crops if available
	if borrower.SecondaryMobilization != "" {
		crops := credentialSubject["registered_crops"].([]string)
		crops = append(crops, borrower.SecondaryMobilization)
		credentialSubject["registered_crops"] = crops
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "FarmerCredential"},
		credentialSubject,
	)
}

// CreateLandRecordCredential creates a land record credential
func (ia *IdentityAdapter) CreateLandRecordCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	borrowerAddress sdk.AccAddress,
	khataNumber string,
	totalArea sdk.Dec,
) (string, error) {
	// Get identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddress)
	if !found {
		return "", fmt.Errorf("identity not found for address: %s", borrowerAddress.String())
	}

	// Create land record credential
	credentialSubject := map[string]interface{}{
		"id":               identity.Did,
		"address":          borrowerAddress.String(),
		"khata_number":     khataNumber,
		"total_area":       totalArea.String(),
		"area_unit":        "acres",
		"verification_date": ctx.BlockTime().Format(time.RFC3339),
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "LandRecordCredential"},
		credentialSubject,
	)
}

// CreateLoanCredential creates a loan credential
func (ia *IdentityAdapter) CreateLoanCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	loan types.Loan,
) (string, error) {
	// Get borrower identity
	borrowerAddr, err := sdk.AccAddressFromBech32(loan.BorrowerAddress)
	if err != nil {
		return "", err
	}

	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddr)
	if !found {
		return "", fmt.Errorf("identity not found for borrower: %s", loan.BorrowerAddress)
	}

	// Create loan credential
	credentialSubject := map[string]interface{}{
		"id":                identity.Did,
		"loan_id":           loan.LoanID,
		"borrower_id":       loan.BorrowerID,
		"loan_scheme_id":    loan.LoanSchemeID,
		"principal":         loan.Principal.String(),
		"interest_rate":     loan.InterestRate.String(),
		"term_months":       loan.Term,
		"monthly_payment":   loan.MonthlyPayment.String(),
		"disbursed_at":      loan.DisbursedAt.Format(time.RFC3339),
		"maturity_date":     loan.MaturityDate.Format(time.RFC3339),
		"status":            loan.Status,
		"liquidity_source":  loan.LiquiditySource,
	}

	// Note: Crop information fields (CropType, CropSeason) may not exist in current Loan struct
	// These will be added when the loan module is fully implemented with agricultural features

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "AgriculturalLoanCredential"},
		credentialSubject,
	)
}

// MigrateExistingBorrowers migrates existing borrowers to identity system
func (ia *IdentityAdapter) MigrateExistingBorrowers(ctx sdk.Context) error {
	borrowers := ia.keeper.GetAllBorrowers(ctx)
	
	for _, borrower := range borrowers {
		borrowerAddr := borrower.Address
		if borrowerAddr.Empty() {
			continue // Skip invalid addresses
		}

		// Check if already migrated
		_, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddr)
		if found {
			continue
		}

		// Create identity
		did := fmt.Sprintf("did:desh:%s", borrowerAddr.String())
		identity := identitytypes.Identity{
			Did:        did,
			Controller: borrowerAddr.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  borrower.RegisteredAt,
			UpdatedAt:  ctx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(ctx, identity)

		// Create basic KYC credential
		credentialSubject := map[string]interface{}{
			"id":            did,
			"name":          borrower.Name,
			"phone":         borrower.Phone,
			"village_code":  borrower.VillageCode,
			"village_name":  borrower.VillageName,
			"kyc_level":     borrower.KYCStatus,
			"verified":      borrower.KYCStatus == "verified",
		}

		// Self-issue for migration
		ia.identityKeeper.IssueCredential(
			ctx,
			borrowerAddr, // Self-issued during migration
			did,
			[]string{"VerifiableCredential", "KYCCredential"},
			credentialSubject,
		)
	}

	return nil
}

// MigrateExistingLoans migrates existing loans to create loan credentials
func (ia *IdentityAdapter) MigrateExistingLoans(ctx sdk.Context) error {
	loans := ia.keeper.GetAllLoans(ctx)
	
	for _, loan := range loans {
		// Only migrate active loans
		if loan.Status != types.StatusActive {
			continue
		}

		borrowerAddr := loan.BorrowerAddress
		if borrowerAddr.Empty() {
			continue
		}

		// Check if identity exists
		identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, borrowerAddr)
		if !found {
			continue
		}

		// Check if loan credential already exists
		loanCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "AgriculturalLoanCredential")
		credExists := false
		for _, cred := range loanCreds {
			if loanID, ok := cred.CredentialSubject["loan_id"].(string); ok && loanID == loan.LoanID {
				credExists = true
				break
			}
		}

		if !credExists {
			// Create loan credential (self-issued during migration)
			ia.CreateLoanCredential(ctx, borrowerAddr, loan)
		}
	}

	return nil
}