package types

const (
	// ModuleName defines the module name
	ModuleName = "shikshamitra"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_shikshamitra"

	// Interest rate constants
	MinInterestRate = "0.04" // 4%
	MaxInterestRate = "0.07" // 7%

	// Loan limits
	MinLoanAmount = "25000"    // ₹25,000
	MaxLoanAmount = "20000000" // ₹20 Lakhs

	// Grace period after course completion
	GracePeriodMonths = 6

	// Festival bonus reduction
	FestivalInterestReduction = "0.0025" // 0.25% reduction during festivals

	// Merit-based reductions
	MeritReduction90Plus  = "0.01"  // 1% for 90%+ marks
	MeritReduction80Plus  = "0.005" // 0.5% for 80-89% marks
	MeritReduction70Plus  = "0.0025" // 0.25% for 70-79% marks
)

// Key prefixes
var (
	LoanKeyPrefix              = []byte{0x01}
	ApplicationKeyPrefix       = []byte{0x02}
	StudentProfilePrefix       = []byte{0x03}
	InstitutionPrefix          = []byte{0x04}
	CoursePrefix               = []byte{0x05}
	RepaymentKeyPrefix         = []byte{0x06}
	ScholarshipPrefix          = []byte{0x07}
	CoApplicantPrefix          = []byte{0x08}
	AcademicRecordPrefix       = []byte{0x09}
	EmploymentRecordPrefix     = []byte{0x0A}
	PINCodeEligiblePrefix      = []byte{0x0B}
)

// GetLoanKey returns the store key for a loan
func GetLoanKey(loanID string) []byte {
	return append(LoanKeyPrefix, []byte(loanID)...)
}

// GetApplicationKey returns the store key for an application
func GetApplicationKey(applicationID string) []byte {
	return append(ApplicationKeyPrefix, []byte(applicationID)...)
}

// GetStudentProfileKey returns the store key for a student profile
func GetStudentProfileKey(studentID string) []byte {
	return append(StudentProfilePrefix, []byte(studentID)...)
}

// GetInstitutionKey returns the store key for an institution
func GetInstitutionKey(institutionID string) []byte {
	return append(InstitutionPrefix, []byte(institutionID)...)
}

// GetCourseKey returns the store key for a course
func GetCourseKey(courseID string) []byte {
	return append(CoursePrefix, []byte(courseID)...)
}

// GetScholarshipKey returns the store key for a scholarship
func GetScholarshipKey(scholarshipID string) []byte {
	return append(ScholarshipPrefix, []byte(scholarshipID)...)
}

// GetCoApplicantKey returns the store key for a co-applicant
func GetCoApplicantKey(coApplicantID string) []byte {
	return append(CoApplicantPrefix, []byte(coApplicantID)...)
}

// GetAcademicRecordKey returns the store key for academic records
func GetAcademicRecordKey(studentID string) []byte {
	return append(AcademicRecordPrefix, []byte(studentID)...)
}

// GetPINCodeKey returns the store key for PIN code eligibility
func GetPINCodeKey(pinCode string) []byte {
	return append(PINCodeEligiblePrefix, []byte(pinCode)...)
}