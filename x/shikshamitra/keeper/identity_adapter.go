package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
	"github.com/deshchain/deshchain/x/shikshamitra/types"
)

// IdentityAdapter provides identity-based verification for ShikshaMitra
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

// VerifyStudentIdentity verifies student identity using the identity module
func (ia *IdentityAdapter) VerifyStudentIdentity(
	ctx sdk.Context,
	studentAddress sdk.AccAddress,
	institutionCode string,
) (*StudentIdentityStatus, error) {
	status := &StudentIdentityStatus{
		Address:               studentAddress.String(),
		HasIdentity:          false,
		IsKYCVerified:        false,
		KYCLevel:             "none",
		HasEducationProfile:  false,
		InstitutionCode:      institutionCode,
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", studentAddress.String())
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
			if credType == "KYCCredential" || credType == "StudentKYC" {
				// Check if credential is valid
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					// Extract KYC level and details
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if age, ok := subject["age"].(float64); ok {
							status.Age = int32(age)
						}
						if aadhaarHash, ok := subject["aadhaar_hash"].(string); ok {
							status.AadhaarHash = aadhaarHash
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

	// Check for education profile credential
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "EducationProfileCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasEducationProfile = true
					status.EducationCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if grade, ok := subject["current_grade"].(string); ok {
							status.CurrentGrade = grade
						}
						if inst, ok := subject["institution"].(string); ok {
							status.VerifiedInstitution = inst
						}
						if score, ok := subject["academic_score"].(float64); ok {
							status.AcademicScore = sdk.NewDecFromBigIntWithPrec(sdk.NewInt(int64(score*100)), 2)
						}
					}
					break
				}
			}
		}
	}

	return status, nil
}

// VerifyCoApplicantIdentity verifies co-applicant (parent/guardian) identity
func (ia *IdentityAdapter) VerifyCoApplicantIdentity(
	ctx sdk.Context,
	coApplicantAddress sdk.AccAddress,
	studentAddress sdk.AccAddress,
) (*CoApplicantIdentityStatus, error) {
	status := &CoApplicantIdentityStatus{
		Address:              coApplicantAddress.String(),
		StudentAddress:       studentAddress.String(),
		HasIdentity:         false,
		IsKYCVerified:       false,
		KYCLevel:            "none",
		HasIncomeProof:      false,
	}

	// Check for DID
	did := fmt.Sprintf("did:desh:%s", coApplicantAddress.String())
	identity, exists := ia.identityKeeper.GetIdentity(ctx, did)
	if !exists {
		return status, nil
	}

	status.HasIdentity = true
	status.DID = did
	status.IdentityStatus = string(identity.Status)

	// Check for KYC and income credentials
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	
	// Check KYC
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "KYCCredential" || credType == "GuardianKYC" {
				if ia.isCredentialValid(ctx, cred) {
					status.IsKYCVerified = true
					status.KYCCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if level, ok := subject["kyc_level"].(string); ok {
							status.KYCLevel = level
						}
						if relation, ok := subject["student_relation"].(string); ok {
							status.RelationToStudent = relation
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

	// Check income proof
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "IncomeCredential" || credType == "EmploymentCredential" {
				if ia.isCredentialValid(ctx, cred) {
					status.HasIncomeProof = true
					status.IncomeCredentialID = credID
					
					if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
						if income, ok := subject["annual_income"].(float64); ok {
							status.VerifiedIncome = sdk.NewDec(int64(income))
						}
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

// CreateStudentCredential creates a student education credential
func (ia *IdentityAdapter) CreateStudentCredential(
	ctx sdk.Context,
	profile types.StudentProfile,
) error {
	// Ensure student has identity
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
				"source": "shikshamitra",
				"type":   "student",
			},
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Create student education profile credential
	credID := fmt.Sprintf("vc:education:%s:%d", profile.ID, ctx.BlockTime().Unix())
	
	expiryDate := ctx.BlockTime().AddDate(4, 0, 0) // Valid for 4 years
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/education/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "EducationProfileCredential"},
		Issuer: "did:desh:shikshamitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                     did,
			"student_id":             profile.ID,
			"name":                   profile.Name,
			"date_of_birth":          profile.DateOfBirth,
			"aadhaar_hash":           profile.AadhaarHash,
			"current_education":      profile.CurrentEducation,
			"institution":            profile.Institution,
			"enrollment_number":      profile.EnrollmentNumber,
			"current_grade":          profile.CurrentGrade,
			"academic_year":          profile.AcademicYear,
			"total_loans_availed":    profile.TotalLoansAvailed,
			"active_loans":           profile.ActiveLoans,
			"verification_status":    profile.VerificationStatus,
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:shikshamitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature", // In production, sign properly
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	return nil
}

// CreateEducationLoanCredential creates a credential for an approved education loan
func (ia *IdentityAdapter) CreateEducationLoanCredential(
	ctx sdk.Context,
	loan types.EducationLoan,
) error {
	did := fmt.Sprintf("did:desh:%s", loan.Borrower)
	
	// Ensure identity exists
	if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
		return fmt.Errorf("identity not found for borrower")
	}

	credID := fmt.Sprintf("vc:edloan:%s:%d", loan.ID, ctx.BlockTime().Unix())
	expiryDate := loan.ExpectedCompletionDate.AddDate(10, 0, 0) // Valid for 10 years after completion
	
	credential := identitytypes.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.bharat/contexts/educationloan/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "EducationLoanCredential"},
		Issuer: "did:desh:shikshamitra-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":                       did,
			"loan_id":                  loan.ID,
			"student_id":               loan.StudentID,
			"institution_id":           loan.InstitutionID,
			"course_name":              loan.CourseName,
			"course_type":              loan.CourseType,
			"course_duration":          loan.CourseDuration,
			"total_semesters":          loan.TotalSemesters,
			"total_loan_amount":        loan.TotalLoanAmount.String(),
			"interest_rate":            loan.InterestRate.String(),
			"repayment_start_date":     loan.RepaymentStartDate,
			"expected_completion_date": loan.ExpectedCompletionDate,
			"status":                   loan.Status,
			"semester_funding_enabled": loan.SemesterFunding,
			"disbursed_amount":         loan.DisbursedAmount.String(),
		},
		Proof: &identitytypes.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:shikshamitra-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-signature",
		},
	}

	// Store credential
	ia.identityKeeper.SetCredential(ctx, credential)
	ia.identityKeeper.AddCredentialToSubject(ctx, did, credID)

	// Also add to co-applicant if exists
	if loan.CoApplicantID != "" {
		coApplicantDID := fmt.Sprintf("did:desh:%s", loan.CoApplicantID)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, coApplicantDID); exists {
			ia.identityKeeper.AddCredentialToSubject(ctx, coApplicantDID, credID)
		}
	}

	return nil
}

// VerifyInstitutionCredential verifies if an institution is accredited
func (ia *IdentityAdapter) VerifyInstitutionCredential(
	ctx sdk.Context,
	institutionID string,
) (bool, error) {
	// In production, this would check for institution accreditation credentials
	// For now, we'll do a simple check
	did := fmt.Sprintf("did:desh:institution:%s", institutionID)
	
	identity, exists := ia.identityKeeper.GetIdentity(ctx, did)
	if !exists {
		return false, nil
	}

	if identity.Status != identitytypes.IdentityStatus_ACTIVE {
		return false, nil
	}

	// Check for accreditation credential
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		for _, credType := range cred.Type {
			if credType == "AccreditationCredential" {
				if ia.isCredentialValid(ctx, cred) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// GetStudentEducationCredentials retrieves all education-related credentials
func (ia *IdentityAdapter) GetStudentEducationCredentials(
	ctx sdk.Context,
	studentAddress sdk.AccAddress,
) ([]EducationCredential, error) {
	did := fmt.Sprintf("did:desh:%s", studentAddress.String())
	
	credentials := ia.identityKeeper.GetCredentialsBySubject(ctx, did)
	educationCreds := []EducationCredential{}
	
	for _, credID := range credentials {
		cred, found := ia.identityKeeper.GetCredential(ctx, credID)
		if !found {
			continue
		}
		
		// Check if it's an education-related credential
		isEducation := false
		for _, credType := range cred.Type {
			if credType == "EducationProfileCredential" || 
			   credType == "EducationLoanCredential" ||
			   credType == "StudentKYC" {
				isEducation = true
				break
			}
		}
		
		if !isEducation {
			continue
		}
		
		educationCred := EducationCredential{
			CredentialID: credID,
			Type:         cred.Type,
			IssuanceDate: cred.IssuanceDate,
			IsValid:      ia.isCredentialValid(ctx, cred),
		}
		
		// Extract relevant info
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if loanID, ok := subject["loan_id"].(string); ok {
				educationCred.LoanID = loanID
			}
			if status, ok := subject["status"].(string); ok {
				educationCred.Status = status
			}
			if institution, ok := subject["institution"].(string); ok {
				educationCred.Institution = institution
			}
			if course, ok := subject["course_name"].(string); ok {
				educationCred.CourseName = course
			}
		}
		
		educationCreds = append(educationCreds, educationCred)
	}
	
	return educationCreds, nil
}

// MigrateExistingStudents migrates existing students to use identity system
func (ia *IdentityAdapter) MigrateExistingStudents(ctx sdk.Context) (int, error) {
	migrated := 0
	
	// Get all student profiles
	profiles := ia.keeper.GetAllStudentProfiles(ctx)
	
	for _, profile := range profiles {
		// Check if already has identity
		did := fmt.Sprintf("did:desh:%s", profile.DhanPataAddress)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); exists {
			continue // Already migrated
		}
		
		// Create student credential
		if err := ia.CreateStudentCredential(ctx, profile); err != nil {
			ia.keeper.Logger(ctx).Error("Failed to migrate student",
				"student_id", profile.ID,
				"error", err)
			continue
		}
		
		migrated++
	}
	
	// Migrate existing loans
	loans := ia.keeper.GetAllEducationLoans(ctx)
	for _, loan := range loans {
		// Check if borrower has identity
		did := fmt.Sprintf("did:desh:%s", loan.Borrower)
		if _, exists := ia.identityKeeper.GetIdentity(ctx, did); !exists {
			continue // Skip if no identity
		}
		
		// Create loan credential
		if err := ia.CreateEducationLoanCredential(ctx, loan); err != nil {
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

type StudentIdentityStatus struct {
	Address                string
	HasIdentity            bool
	DID                    string
	IdentityStatus         string
	IsKYCVerified          bool
	KYCLevel               string
	KYCCredentialID        string
	HasEducationProfile    bool
	EducationCredentialID  string
	Age                    int32
	AadhaarHash           string
	InstitutionCode       string
	VerifiedInstitution   string
	CurrentGrade          string
	AcademicScore         sdk.Dec
}

type CoApplicantIdentityStatus struct {
	Address               string
	StudentAddress        string
	HasIdentity          bool
	DID                  string
	IdentityStatus       string
	IsKYCVerified        bool
	KYCLevel             string
	KYCCredentialID      string
	HasIncomeProof       bool
	IncomeCredentialID   string
	VerifiedIncome       sdk.Dec
	Employer             string
	RelationToStudent    string
}

type EducationCredential struct {
	CredentialID string
	Type         []string
	LoanID       string
	Status       string
	Institution  string
	CourseName   string
	IssuanceDate time.Time
	IsValid      bool
}