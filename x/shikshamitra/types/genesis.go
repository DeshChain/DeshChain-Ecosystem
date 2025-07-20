package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Loans: []EducationLoan{},
		StudentProfiles: []StudentProfile{},
		FestivalOffers: []FestivalOffer{},
		LoanCounter: 0,
		ScholarshipCounter: 0,
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Check for duplicate loan IDs
	loanIDs := make(map[string]bool)
	for _, loan := range gs.Loans {
		if loanIDs[loan.LoanID] {
			return fmt.Errorf("duplicate loan ID: %s", loan.LoanID)
		}
		loanIDs[loan.LoanID] = true

		// Validate loan
		if err := validateLoan(loan); err != nil {
			return fmt.Errorf("invalid loan %s: %w", loan.LoanID, err)
		}
	}

	// Check for duplicate student addresses
	studentAddrs := make(map[string]bool)
	for _, profile := range gs.StudentProfiles {
		if studentAddrs[profile.Address] {
			return fmt.Errorf("duplicate student address: %s", profile.Address)
		}
		studentAddrs[profile.Address] = true

		// Validate profile
		if err := validateStudentProfile(profile); err != nil {
			return fmt.Errorf("invalid student profile %s: %w", profile.Address, err)
		}
	}

	// Validate festival offers
	for _, offer := range gs.FestivalOffers {
		if err := validateFestivalOffer(offer); err != nil {
			return fmt.Errorf("invalid festival offer %s: %w", offer.FestivalID, err)
		}
	}

	return nil
}

func validateLoan(loan EducationLoan) error {
	if loan.LoanID == "" {
		return fmt.Errorf("loan ID cannot be empty")
	}
	if loan.Borrower == "" {
		return fmt.Errorf("borrower cannot be empty")
	}
	if loan.StudentName == "" {
		return fmt.Errorf("student name cannot be empty")
	}
	if loan.CoApplicant == "" {
		return fmt.Errorf("co-applicant cannot be empty")
	}
	if !loan.LoanAmount.IsValid() || loan.LoanAmount.IsZero() {
		return fmt.Errorf("invalid loan amount")
	}
	if loan.InterestRate == "" {
		return fmt.Errorf("interest rate cannot be empty")
	}
	return nil
}

func validateStudentProfile(profile StudentProfile) error {
	if profile.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if profile.DhanPataID == "" {
		return fmt.Errorf("DhanPata ID cannot be empty")
	}
	if profile.StudentName == "" {
		return fmt.Errorf("student name cannot be empty")
	}
	if len(profile.Pincode) != 6 {
		return fmt.Errorf("invalid PIN code")
	}
	return nil
}

func validateFestivalOffer(offer FestivalOffer) error {
	if offer.FestivalID == "" {
		return fmt.Errorf("festival ID cannot be empty")
	}
	if offer.FestivalName == "" {
		return fmt.Errorf("festival name cannot be empty")
	}
	if offer.InterestReduction == "" {
		return fmt.Errorf("interest reduction cannot be empty")
	}
	if offer.EndDate.Before(offer.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}
	return nil
}