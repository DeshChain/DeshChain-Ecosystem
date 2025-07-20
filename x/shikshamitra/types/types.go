package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LoanStatus represents the status of an education loan
type LoanStatus int32

const (
	LoanStatus_PENDING LoanStatus = iota
	LoanStatus_APPROVED
	LoanStatus_DISBURSED
	LoanStatus_IN_MORATORIUM
	LoanStatus_REPAYING
	LoanStatus_COMPLETED
	LoanStatus_DEFAULTED
	LoanStatus_REJECTED
)

// CourseType represents different education levels
type CourseType int32

const (
	CourseType_DIPLOMA CourseType = iota
	CourseType_UNDERGRADUATE
	CourseType_POSTGRADUATE
	CourseType_DOCTORATE
	CourseType_PROFESSIONAL
	CourseType_VOCATIONAL
	CourseType_CERTIFICATE
	CourseType_ONLINE
)

// InstitutionType represents types of educational institutions
type InstitutionType int32

const (
	InstitutionType_IIT InstitutionType = iota
	InstitutionType_IIM
	InstitutionType_NIT
	InstitutionType_CENTRAL_UNIVERSITY
	InstitutionType_STATE_UNIVERSITY
	InstitutionType_DEEMED_UNIVERSITY
	InstitutionType_PRIVATE_UNIVERSITY
	InstitutionType_FOREIGN_UNIVERSITY
	InstitutionType_VOCATIONAL_INSTITUTE
	InstitutionType_ONLINE_PLATFORM
)

// EducationLoan represents an education loan
type EducationLoan struct {
	ID                  string        `json:"id"`
	StudentID           string        `json:"student_id"`
	Borrower            string        `json:"borrower"`
	DhanPataAddress     string        `json:"dhanpata_address"`
	CoApplicantID       string        `json:"co_applicant_id,omitempty"`
	Amount              sdk.Coin      `json:"amount"`
	InterestRate        sdk.Dec       `json:"interest_rate"`
	Term                int64         `json:"term"` // in months
	MoratoriumPeriod    int64         `json:"moratorium_period"` // in months
	Status              LoanStatus    `json:"status"`
	InstitutionID       string        `json:"institution_id"`
	CourseID            string        `json:"course_id"`
	CourseStartDate     time.Time     `json:"course_start_date"`
	CourseEndDate       time.Time     `json:"course_end_date"`
	CreatedAt           time.Time     `json:"created_at"`
	DisbursedAt         *time.Time    `json:"disbursed_at,omitempty"`
	RepaymentStartDate  *time.Time    `json:"repayment_start_date,omitempty"`
	MaturityDate        *time.Time    `json:"maturity_date,omitempty"`
	RepaidAmount        sdk.Coin      `json:"repaid_amount"`
	LastRepaymentDate   *time.Time    `json:"last_repayment_date,omitempty"`
	PINCode             string        `json:"pin_code"`
	FestivalBonus       bool          `json:"festival_bonus"`
	MeritScholarship    sdk.Dec       `json:"merit_scholarship"`
	InsuranceRequired   bool          `json:"insurance_required"`
	InsurancePremium    sdk.Coin      `json:"insurance_premium"`
}

// StudentProfile represents a student's profile
type StudentProfile struct {
	ID                   string                `json:"id"`
	Name                 string                `json:"name"`
	DhanPataAddress      string                `json:"dhanpata_address"`
	DateOfBirth          time.Time             `json:"date_of_birth"`
	AadhaarHash          string                `json:"aadhaar_hash"`
	Email                string                `json:"email"`
	Phone                string                `json:"phone"`
	PINCode              string                `json:"pin_code"`
	CurrentEducation     string                `json:"current_education"`
	AcademicRecords      []AcademicRecord      `json:"academic_records"`
	EntranceExamScores   []EntranceExamScore   `json:"entrance_exam_scores"`
	ParentIncome         sdk.Coin              `json:"parent_income"`
	FamilyMembers        int32                 `json:"family_members"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	VerificationStatus   string                `json:"verification_status"`
	TotalLoansAvailed    int32                 `json:"total_loans_availed"`
	ActiveLoans          int32                 `json:"active_loans"`
}

// LoanApplication represents an education loan application
type LoanApplication struct {
	ID                   string              `json:"id"`
	StudentID            string              `json:"student_id"`
	Applicant            string              `json:"applicant"`
	DhanPataAddress      string              `json:"dhanpata_address"`
	RequestedAmount      sdk.Coin            `json:"requested_amount"`
	InstitutionID        string              `json:"institution_id"`
	CourseID             string              `json:"course_id"`
	AdmissionStatus      string              `json:"admission_status"`
	AdmissionLetter      string              `json:"admission_letter_hash"`
	FeeStructure         []FeeComponent      `json:"fee_structure"`
	LivingExpenses       sdk.Coin            `json:"living_expenses"`
	BooksMaterials       sdk.Coin            `json:"books_materials"`
	ComputerEquipment    sdk.Coin            `json:"computer_equipment,omitempty"`
	TravelExpenses       sdk.Coin            `json:"travel_expenses,omitempty"`
	CoApplicantDetails   CoApplicant         `json:"co_applicant_details"`
	GuarantorDetails     []Guarantor         `json:"guarantor_details,omitempty"`
	CollateralOffered    string              `json:"collateral_offered,omitempty"`
	AppliedAt            time.Time           `json:"applied_at"`
	ReviewedBy           string              `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time          `json:"reviewed_at,omitempty"`
	Status               LoanStatus          `json:"status"`
	RejectionReason      string              `json:"rejection_reason,omitempty"`
	ProposedInterestRate sdk.Dec             `json:"proposed_interest_rate"`
	ScholarshipDetails   []ScholarshipInfo   `json:"scholarship_details,omitempty"`
}

// Institution represents an educational institution
type Institution struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Type              InstitutionType `json:"type"`
	AccreditationRank string          `json:"accreditation_rank"` // A+, A, B+, etc
	Location          string          `json:"location"`
	Country           string          `json:"country"`
	EstablishedYear   int32           `json:"established_year"`
	IsRecognized      bool            `json:"is_recognized"`
	MaxLoanAmount     sdk.Coin        `json:"max_loan_amount"`
	Courses           []string        `json:"course_ids"`
	PlacementRate     sdk.Dec         `json:"placement_rate"`
	AvgStartingSalary sdk.Coin        `json:"avg_starting_salary"`
	Website           string          `json:"website"`
	ContactEmail      string          `json:"contact_email"`
}

// Course represents an educational course
type Course struct {
	ID                string      `json:"id"`
	InstitutionID     string      `json:"institution_id"`
	Name              string      `json:"name"`
	Type              CourseType  `json:"type"`
	Duration          int32       `json:"duration"` // in months
	TotalFees         sdk.Coin    `json:"total_fees"`
	AnnualFees        sdk.Coin    `json:"annual_fees"`
	EligibilityCriteria string    `json:"eligibility_criteria"`
	MinPercentage     sdk.Dec     `json:"min_percentage"`
	EntranceExamRequired bool     `json:"entrance_exam_required"`
	AcceptedExams     []string    `json:"accepted_exams"`
	MaxAge            int32       `json:"max_age,omitempty"`
	Language          string      `json:"language"`
	Mode              string      `json:"mode"` // regular, distance, online
}

// CoApplicant represents a co-applicant (usually parent/guardian)
type CoApplicant struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Relationship    string   `json:"relationship"`
	DhanPataAddress string   `json:"dhanpata_address"`
	AadhaarHash     string   `json:"aadhaar_hash"`
	Occupation      string   `json:"occupation"`
	AnnualIncome    sdk.Coin `json:"annual_income"`
	IncomeProof     string   `json:"income_proof_hash"`
	PanNumber       string   `json:"pan_number"`
	CreditScore     int32    `json:"credit_score"`
	ExistingLoans   sdk.Coin `json:"existing_loans"`
}

// AcademicRecord represents academic history
type AcademicRecord struct {
	Level           string    `json:"level"` // 10th, 12th, graduation, etc
	Board           string    `json:"board"`
	Institution     string    `json:"institution"`
	YearOfPassing   int32     `json:"year_of_passing"`
	Percentage      sdk.Dec   `json:"percentage"`
	Grade           string    `json:"grade,omitempty"`
	MarksheetHash   string    `json:"marksheet_hash"`
	VerifiedAt      time.Time `json:"verified_at"`
}

// EntranceExamScore represents entrance exam results
type EntranceExamScore struct {
	ExamName        string    `json:"exam_name"` // JEE, NEET, CAT, etc
	Score           sdk.Dec   `json:"score"`
	Rank            int64     `json:"rank,omitempty"`
	Percentile      sdk.Dec   `json:"percentile,omitempty"`
	Year            int32     `json:"year"`
	ScoreCardHash   string    `json:"score_card_hash"`
	VerifiedAt      time.Time `json:"verified_at"`
}

// FeeComponent represents breakdown of fees
type FeeComponent struct {
	Type        string   `json:"type"` // tuition, hostel, lab, library, etc
	Amount      sdk.Coin `json:"amount"`
	Frequency   string   `json:"frequency"` // annual, semester, one-time
	IsMandatory bool     `json:"is_mandatory"`
}

// Guarantor represents a loan guarantor
type Guarantor struct {
	Name            string   `json:"name"`
	DhanPataAddress string   `json:"dhanpata_address"`
	Relationship    string   `json:"relationship"`
	NetWorth        sdk.Coin `json:"net_worth"`
	ConsentHash     string   `json:"consent_hash"`
}

// ScholarshipInfo represents scholarship details
type ScholarshipInfo struct {
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`
	Amount          sdk.Coin `json:"amount"`
	Duration        int32    `json:"duration"` // in months
	Status          string   `json:"status"`
	DocumentHash    string   `json:"document_hash"`
}

// Repayment represents a loan repayment
type Repayment struct {
	ID              string    `json:"id"`
	LoanID          string    `json:"loan_id"`
	Amount          sdk.Coin  `json:"amount"`
	Principal       sdk.Coin  `json:"principal"`
	Interest        sdk.Coin  `json:"interest"`
	PenaltyAmount   sdk.Coin  `json:"penalty_amount,omitempty"`
	PaidBy          string    `json:"paid_by"`
	PaidAt          time.Time `json:"paid_at"`
	TransactionID   string    `json:"transaction_id"`
	PaymentMethod   string    `json:"payment_method"`
	ReceiptNumber   string    `json:"receipt_number"`
}

// EmploymentRecord for loan repayment tracking
type EmploymentRecord struct {
	StudentID       string    `json:"student_id"`
	EmployerName    string    `json:"employer_name"`
	Position        string    `json:"position"`
	JoiningDate     time.Time `json:"joining_date"`
	CurrentSalary   sdk.Coin  `json:"current_salary"`
	EmploymentType  string    `json:"employment_type"` // full-time, part-time, contract
	VerificationDoc string    `json:"verification_doc_hash"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PINCodeEligibility for education loans
type PINCodeEligibility struct {
	PINCode             string   `json:"pin_code"`
	DistrictName        string   `json:"district_name"`
	StateName           string   `json:"state_name"`
	IsEligible          bool     `json:"is_eligible"`
	MaxLoanAmount       sdk.Coin `json:"max_loan_amount"`
	BaseInterestRate    sdk.Dec  `json:"base_interest_rate"`
	RuralArea           bool     `json:"rural_area"`
	BackwardDistrict    bool     `json:"backward_district"`
	EducationHub        bool     `json:"education_hub"`
	LocalInstitutions   []string `json:"local_institutions"`
}