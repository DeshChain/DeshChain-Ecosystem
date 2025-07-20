package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types
const (
	TypeMsgApplyEducationLoan       = "apply_education_loan"
	TypeMsgApproveLoan              = "approve_loan"
	TypeMsgRejectLoan               = "reject_loan" 
	TypeMsgDisburseLoan             = "disburse_loan"
	TypeMsgRepayLoan                = "repay_loan"
	TypeMsgUpdateAcademicProgress   = "update_academic_progress"
	TypeMsgApplyScholarship         = "apply_scholarship"
	TypeMsgUpdateEmploymentStatus   = "update_employment_status"
)

// MsgApplyEducationLoan - Apply for education loan
type MsgApplyEducationLoan struct {
	Applicant              string            `json:"applicant"`
	DhanPataID             string            `json:"dhanpata_id"`
	StudentName            string            `json:"student_name"`
	CoApplicant            string            `json:"co_applicant"` // Parent/Guardian
	CoApplicantDhanPataID  string            `json:"co_applicant_dhanpata_id"`
	InstitutionName        string            `json:"institution_name"`
	InstitutionType        InstitutionType   `json:"institution_type"`
	CourseType             CourseType        `json:"course_type"`
	CourseName             string            `json:"course_name"`
	CourseDuration         int32             `json:"course_duration"` // in months
	AcademicYear           string            `json:"academic_year"`
	LoanAmount             sdk.Coin          `json:"loan_amount"`
	LoanComponents         []LoanComponent   `json:"loan_components"`
	EntranceExamScore      string            `json:"entrance_exam_score"`
	PreviousAcademicRecord string            `json:"previous_academic_record"` // IPFS hash
	AdmissionLetter        string            `json:"admission_letter"` // IPFS hash
	Pincode                string            `json:"pincode"`
	FamilyIncome           sdk.Coin          `json:"family_income"` // Annual
	CollateralOffered      string            `json:"collateral_offered"`
	CulturalQuote          string            `json:"cultural_quote"`
}

// LoanComponent represents different components of education loan
type LoanComponent struct {
	ComponentType string   `json:"component_type"` // TUITION, HOSTEL, BOOKS, EQUIPMENT, LIVING
	Amount        sdk.Coin `json:"amount"`
	Description   string   `json:"description"`
}

func NewMsgApplyEducationLoan(
	applicant, dhanPataID, studentName, coApplicant, coApplicantDhanPataID string,
	institutionName string, institutionType InstitutionType, courseType CourseType,
	courseName string, courseDuration int32, academicYear string,
	loanAmount sdk.Coin, loanComponents []LoanComponent,
	entranceExamScore, previousAcademicRecord, admissionLetter, pincode string,
	familyIncome sdk.Coin, collateralOffered, culturalQuote string,
) *MsgApplyEducationLoan {
	return &MsgApplyEducationLoan{
		Applicant:              applicant,
		DhanPataID:             dhanPataID,
		StudentName:            studentName,
		CoApplicant:            coApplicant,
		CoApplicantDhanPataID:  coApplicantDhanPataID,
		InstitutionName:        institutionName,
		InstitutionType:        institutionType,
		CourseType:             courseType,
		CourseName:             courseName,
		CourseDuration:         courseDuration,
		AcademicYear:           academicYear,
		LoanAmount:             loanAmount,
		LoanComponents:         loanComponents,
		EntranceExamScore:      entranceExamScore,
		PreviousAcademicRecord: previousAcademicRecord,
		AdmissionLetter:        admissionLetter,
		Pincode:                pincode,
		FamilyIncome:           familyIncome,
		CollateralOffered:      collateralOffered,
		CulturalQuote:          culturalQuote,
	}
}

func (msg *MsgApplyEducationLoan) Route() string { return RouterKey }
func (msg *MsgApplyEducationLoan) Type() string  { return TypeMsgApplyEducationLoan }

func (msg *MsgApplyEducationLoan) GetSigners() []sdk.AccAddress {
	applicant, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{applicant}
}

func (msg *MsgApplyEducationLoan) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyEducationLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.CoApplicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid co-applicant address (%s)", err)
	}

	if msg.StudentName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "student name cannot be empty")
	}

	if msg.InstitutionName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "institution name cannot be empty")
	}

	if !msg.LoanAmount.IsValid() || msg.LoanAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "loan amount must be positive")
	}

	if msg.CourseDuration < 6 || msg.CourseDuration > 96 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "course duration must be between 6-96 months")
	}

	if len(msg.Pincode) != 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid PIN code")
	}

	return nil
}

// MsgUpdateAcademicProgress - Update student's academic progress
type MsgUpdateAcademicProgress struct {
	Authority        string  `json:"authority"`
	LoanID           string  `json:"loan_id"`
	Semester         int32   `json:"semester"`
	GradePointAvg    sdk.Dec `json:"grade_point_avg"`
	Attendance       sdk.Dec `json:"attendance"` // Percentage
	Marksheet        string  `json:"marksheet"` // IPFS hash
	Remarks          string  `json:"remarks"`
	ContinuationStatus string `json:"continuation_status"` // CONTINUING, DROPPED, COMPLETED
}

func NewMsgUpdateAcademicProgress(
	authority, loanID string,
	semester int32,
	gradePointAvg, attendance sdk.Dec,
	marksheet, remarks, continuationStatus string,
) *MsgUpdateAcademicProgress {
	return &MsgUpdateAcademicProgress{
		Authority:          authority,
		LoanID:             loanID,
		Semester:           semester,
		GradePointAvg:      gradePointAvg,
		Attendance:         attendance,
		Marksheet:          marksheet,
		Remarks:            remarks,
		ContinuationStatus: continuationStatus,
	}
}

func (msg *MsgUpdateAcademicProgress) Route() string { return RouterKey }
func (msg *MsgUpdateAcademicProgress) Type() string  { return TypeMsgUpdateAcademicProgress }

func (msg *MsgUpdateAcademicProgress) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateAcademicProgress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAcademicProgress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if msg.GradePointAvg.IsNegative() || msg.GradePointAvg.GT(sdk.NewDec(10)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "GPA must be between 0-10")
	}

	if msg.Attendance.IsNegative() || msg.Attendance.GT(sdk.NewDec(100)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "attendance must be between 0-100%")
	}

	return nil
}

// MsgUpdateEmploymentStatus - Update employment status after course completion
type MsgUpdateEmploymentStatus struct {
	Student          string   `json:"student"`
	LoanID           string   `json:"loan_id"`
	EmploymentStatus string   `json:"employment_status"` // EMPLOYED, SELF_EMPLOYED, HIGHER_STUDIES, SEEKING
	EmployerName     string   `json:"employer_name"`
	JobTitle         string   `json:"job_title"`
	MonthlySalary    sdk.Coin `json:"monthly_salary"`
	JoiningDate      string   `json:"joining_date"`
	OfferLetter      string   `json:"offer_letter"` // IPFS hash
}

func NewMsgUpdateEmploymentStatus(
	student, loanID, employmentStatus, employerName, jobTitle string,
	monthlySalary sdk.Coin,
	joiningDate, offerLetter string,
) *MsgUpdateEmploymentStatus {
	return &MsgUpdateEmploymentStatus{
		Student:          student,
		LoanID:           loanID,
		EmploymentStatus: employmentStatus,
		EmployerName:     employerName,
		JobTitle:         jobTitle,
		MonthlySalary:    monthlySalary,
		JoiningDate:      joiningDate,
		OfferLetter:      offerLetter,
	}
}

func (msg *MsgUpdateEmploymentStatus) Route() string { return RouterKey }
func (msg *MsgUpdateEmploymentStatus) Type() string  { return TypeMsgUpdateEmploymentStatus }

func (msg *MsgUpdateEmploymentStatus) GetSigners() []sdk.AccAddress {
	student, err := sdk.AccAddressFromBech32(msg.Student)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{student}
}

func (msg *MsgUpdateEmploymentStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateEmploymentStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Student)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid student address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if msg.EmploymentStatus == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "employment status cannot be empty")
	}

	if msg.EmploymentStatus == "EMPLOYED" && !msg.MonthlySalary.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid monthly salary")
	}

	return nil
}

// MsgApplyScholarship - Apply for merit-based scholarship
type MsgApplyScholarship struct {
	Applicant         string  `json:"applicant"`
	LoanID            string  `json:"loan_id"`
	ScholarshipType   string  `json:"scholarship_type"` // MERIT, NEED, SPORTS, SPECIAL
	AcademicScore     sdk.Dec `json:"academic_score"`
	Achievements      string  `json:"achievements"` // Description
	Recommendations   string  `json:"recommendations"` // IPFS hash
	FinancialNeed     string  `json:"financial_need"` // Justification
	RequestedAmount   sdk.Coin `json:"requested_amount"`
}

func NewMsgApplyScholarship(
	applicant, loanID, scholarshipType string,
	academicScore sdk.Dec,
	achievements, recommendations, financialNeed string,
	requestedAmount sdk.Coin,
) *MsgApplyScholarship {
	return &MsgApplyScholarship{
		Applicant:       applicant,
		LoanID:          loanID,
		ScholarshipType: scholarshipType,
		AcademicScore:   academicScore,
		Achievements:    achievements,
		Recommendations: recommendations,
		FinancialNeed:   financialNeed,
		RequestedAmount: requestedAmount,
	}
}

func (msg *MsgApplyScholarship) Route() string { return RouterKey }
func (msg *MsgApplyScholarship) Type() string  { return TypeMsgApplyScholarship }

func (msg *MsgApplyScholarship) GetSigners() []sdk.AccAddress {
	applicant, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{applicant}
}

func (msg *MsgApplyScholarship) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyScholarship) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Applicant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid applicant address (%s)", err)
	}

	if msg.LoanID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "loan ID cannot be empty")
	}

	if msg.ScholarshipType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "scholarship type cannot be empty")
	}

	if !msg.RequestedAmount.IsValid() || msg.RequestedAmount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "requested amount must be positive")
	}

	return nil
}