package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateStudentsToIdentity = "migrate_students_to_identity"
	TypeMsgCreateStudentCredential   = "create_student_credential"
)

var (
	_ sdk.Msg = &MsgMigrateStudentsToIdentity{}
	_ sdk.Msg = &MsgCreateStudentCredential{}
)

// MsgMigrateStudentsToIdentity migrates students to identity system
type MsgMigrateStudentsToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateStudentsToIdentity(authority string) *MsgMigrateStudentsToIdentity {
	return &MsgMigrateStudentsToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateStudentsToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateStudentsToIdentity) Type() string {
	return TypeMsgMigrateStudentsToIdentity
}

func (msg *MsgMigrateStudentsToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateStudentsToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateStudentsToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateStudentCredential creates a credential for a student
type MsgCreateStudentCredential struct {
	Authority string `json:"authority"`
	StudentId string `json:"student_id"`
}

func NewMsgCreateStudentCredential(authority, studentId string) *MsgCreateStudentCredential {
	return &MsgCreateStudentCredential{
		Authority: authority,
		StudentId: studentId,
	}
}

func (msg *MsgCreateStudentCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateStudentCredential) Type() string {
	return TypeMsgCreateStudentCredential
}

func (msg *MsgCreateStudentCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateStudentCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateStudentCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	if msg.StudentId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "student ID cannot be empty")
	}
	
	return nil
}

// Query types for identity integration

// QueryStudentIdentityRequest is the request for student identity
type QueryStudentIdentityRequest struct {
	StudentAddress string `json:"student_address"`
}

// QueryStudentIdentityResponse is the response for student identity
type QueryStudentIdentityResponse struct {
	Address               string  `json:"address"`
	HasIdentity          bool    `json:"has_identity"`
	DID                  string  `json:"did,omitempty"`
	IsKYCVerified        bool    `json:"is_kyc_verified"`
	KYCLevel             string  `json:"kyc_level"`
	HasEducationProfile  bool    `json:"has_education_profile"`
	VerifiedInstitution  string  `json:"verified_institution,omitempty"`
	AcademicScore        sdk.Dec `json:"academic_score,omitempty"`
	ActiveLoans          []string `json:"active_loans"`
}

// QueryCoApplicantIdentityRequest is the request for co-applicant identity
type QueryCoApplicantIdentityRequest struct {
	CoApplicantAddress string `json:"co_applicant_address"`
	StudentAddress     string `json:"student_address"`
}

// QueryCoApplicantIdentityResponse is the response for co-applicant identity
type QueryCoApplicantIdentityResponse struct {
	Address            string  `json:"address"`
	HasIdentity        bool    `json:"has_identity"`
	DID                string  `json:"did,omitempty"`
	IsKYCVerified      bool    `json:"is_kyc_verified"`
	KYCLevel           string  `json:"kyc_level"`
	HasIncomeProof     bool    `json:"has_income_proof"`
	VerifiedIncome     sdk.Dec `json:"verified_income,omitempty"`
	RelationToStudent  string  `json:"relation_to_student,omitempty"`
}

// QueryEducationCredentialsRequest is the request for education credentials
type QueryEducationCredentialsRequest struct {
	StudentAddress string `json:"student_address"`
}

// QueryEducationCredentialsResponse is the response for education credentials
type QueryEducationCredentialsResponse struct {
	Address         string                   `json:"address"`
	DID             string                   `json:"did"`
	Credentials     []EducationCredentialInfo `json:"credentials"`
}

// EducationCredentialInfo contains education credential information
type EducationCredentialInfo struct {
	CredentialID string   `json:"credential_id"`
	Type         []string `json:"type"`
	LoanID       string   `json:"loan_id,omitempty"`
	Status       string   `json:"status"`
	Institution  string   `json:"institution,omitempty"`
	CourseName   string   `json:"course_name,omitempty"`
	IsValid      bool     `json:"is_valid"`
}