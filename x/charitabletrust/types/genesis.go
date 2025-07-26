package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		TrustGovernance: TrustGovernance{
			Trustees: []Trustee{}, // To be set on chain initialization
			Quorum:   4,          // 4 out of 7 trustees
			ApprovalThreshold: sdk.MustNewDecFromStr("0.571"), // 4/7 = 57.1%
			AdvisoryCommittee: []AdvisoryMember{},
			TransparencyOfficer: "", // To be appointed
			NextElection: time.Time{}, // Will be set during genesis initialization
		},
		TrustFundBalance: TrustFundBalance{
			TotalBalance:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			AllocatedAmount:  sdk.NewCoin("unamo", sdk.ZeroInt()),
			AvailableAmount:  sdk.NewCoin("unamo", sdk.ZeroInt()),
			TotalDistributed: sdk.NewCoin("unamo", sdk.ZeroInt()),
			LastUpdated:      time.Time{}, // Will be set during genesis initialization
		},
		Allocations:     []CharitableAllocation{},
		Proposals:       []AllocationProposal{},
		ImpactReports:   []ImpactReport{},
		FraudAlerts:     []FraudAlert{},
		AllocationCount: 0,
		ProposalCount:   0,
		ReportCount:     0,
		AlertCount:      0,
	}
}

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{
		Enabled:                    true,
		MinAllocationAmount:        sdk.NewCoin("unamo", sdk.NewInt(100000000)), // 100 NAMO
		MaxMonthlyAllocationPerOrg: sdk.NewCoin("unamo", sdk.NewInt(100000000000)), // 100K NAMO
		ProposalVotingPeriod:       604800, // 7 days in seconds
		FraudInvestigationPeriod:   30,     // 30 days
		ImpactReportFrequency:      30,     // 30 days
		AuthorizedInvestigators:    []string{}, // To be set by governance
		DistributionCategories: []string{
			"education",
			"healthcare",
			"rural_development",
			"women_empowerment",
			"emergency_relief",
		},
		EmergencyPauseAuthorities: []string{}, // To be set by governance
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	
	// Validate trust governance
	if err := gs.TrustGovernance.Validate(); err != nil {
		return err
	}
	
	// Validate allocations
	allocationIDs := make(map[uint64]bool)
	for _, allocation := range gs.Allocations {
		if allocationIDs[allocation.Id] {
			return fmt.Errorf("duplicate allocation ID: %d", allocation.Id)
		}
		allocationIDs[allocation.Id] = true
		
		if err := allocation.Validate(); err != nil {
			return err
		}
	}
	
	// Validate proposals
	proposalIDs := make(map[uint64]bool)
	for _, proposal := range gs.Proposals {
		if proposalIDs[proposal.Id] {
			return fmt.Errorf("duplicate proposal ID: %d", proposal.Id)
		}
		proposalIDs[proposal.Id] = true
		
		if err := proposal.Validate(); err != nil {
			return err
		}
	}
	
	// Validate impact reports
	reportIDs := make(map[uint64]bool)
	for _, report := range gs.ImpactReports {
		if reportIDs[report.Id] {
			return fmt.Errorf("duplicate report ID: %d", report.Id)
		}
		reportIDs[report.Id] = true
	}
	
	// Validate fraud alerts
	alertIDs := make(map[uint64]bool)
	for _, alert := range gs.FraudAlerts {
		if alertIDs[alert.Id] {
			return fmt.Errorf("duplicate alert ID: %d", alert.Id)
		}
		alertIDs[alert.Id] = true
	}
	
	return nil
}

// Validate validates the module parameters
func (p Params) Validate() error {
	if p.MinAllocationAmount.IsNegative() {
		return fmt.Errorf("minimum allocation amount cannot be negative")
	}
	
	if p.MaxMonthlyAllocationPerOrg.IsNegative() {
		return fmt.Errorf("max monthly allocation per org cannot be negative")
	}
	
	if p.ProposalVotingPeriod <= 0 {
		return fmt.Errorf("proposal voting period must be positive")
	}
	
	if p.FraudInvestigationPeriod <= 0 {
		return fmt.Errorf("fraud investigation period must be positive")
	}
	
	if p.ImpactReportFrequency <= 0 {
		return fmt.Errorf("impact report frequency must be positive")
	}
	
	// Validate distribution categories
	if len(p.DistributionCategories) == 0 {
		return fmt.Errorf("must have at least one distribution category")
	}
	
	categoryMap := make(map[string]bool)
	for _, cat := range p.DistributionCategories {
		if categoryMap[cat] {
			return fmt.Errorf("duplicate distribution category: %s", cat)
		}
		categoryMap[cat] = true
	}
	
	// Validate authorized addresses
	for _, addr := range p.AuthorizedInvestigators {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid investigator address: %w", err)
		}
	}
	
	for _, addr := range p.EmergencyPauseAuthorities {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid emergency pause authority address: %w", err)
		}
	}
	
	return nil
}

// Validate validates the trust governance
func (tg TrustGovernance) Validate() error {
	if len(tg.Trustees) == 0 {
		return fmt.Errorf("must have at least one trustee")
	}
	
	if tg.Quorum <= 0 || tg.Quorum > int32(len(tg.Trustees)) {
		return fmt.Errorf("quorum must be between 1 and number of trustees")
	}
	
	if tg.ApprovalThreshold.IsNegative() || tg.ApprovalThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("approval threshold must be between 0 and 1")
	}
	
	// Validate trustees
	trusteeAddresses := make(map[string]bool)
	for _, trustee := range tg.Trustees {
		if trusteeAddresses[trustee.Address] {
			return fmt.Errorf("duplicate trustee address: %s", trustee.Address)
		}
		trusteeAddresses[trustee.Address] = true
		
		if _, err := sdk.AccAddressFromBech32(trustee.Address); err != nil {
			return fmt.Errorf("invalid trustee address: %w", err)
		}
		
		if trustee.Name == "" {
			return fmt.Errorf("trustee name cannot be empty")
		}
		
		if trustee.Role == "" {
			return fmt.Errorf("trustee role cannot be empty")
		}
	}
	
	if tg.TransparencyOfficer != "" {
		if _, err := sdk.AccAddressFromBech32(tg.TransparencyOfficer); err != nil {
			return fmt.Errorf("invalid transparency officer address: %w", err)
		}
	}
	
	return nil
}

// Validate validates a charitable allocation
func (ca CharitableAllocation) Validate() error {
	if ca.Id == 0 {
		return fmt.Errorf("allocation ID cannot be zero")
	}
	
	if ca.CharitableOrgWalletId == 0 {
		return fmt.Errorf("charitable org wallet ID cannot be zero")
	}
	
	if ca.OrganizationName == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	
	if ca.Amount.IsNegative() {
		return fmt.Errorf("allocation amount cannot be negative")
	}
	
	if ca.Purpose == "" {
		return fmt.Errorf("allocation purpose cannot be empty")
	}
	
	if ca.Category == "" {
		return fmt.Errorf("allocation category cannot be empty")
	}
	
	for _, approver := range ca.ApprovedBy {
		if _, err := sdk.AccAddressFromBech32(approver); err != nil {
			return fmt.Errorf("invalid approver address: %w", err)
		}
	}
	
	return nil
}

// Validate validates an allocation proposal
func (ap AllocationProposal) Validate() error {
	if ap.Id == 0 {
		return fmt.Errorf("proposal ID cannot be zero")
	}
	
	if _, err := sdk.AccAddressFromBech32(ap.Proposer); err != nil {
		return fmt.Errorf("invalid proposer address: %w", err)
	}
	
	if ap.Title == "" {
		return fmt.Errorf("proposal title cannot be empty")
	}
	
	if ap.Description == "" {
		return fmt.Errorf("proposal description cannot be empty")
	}
	
	if ap.TotalAmount.IsNegative() {
		return fmt.Errorf("total amount cannot be negative")
	}
	
	if len(ap.Allocations) == 0 {
		return fmt.Errorf("proposal must have at least one allocation")
	}
	
	totalAllocated := sdk.NewCoin(ap.TotalAmount.Denom, sdk.ZeroInt())
	for _, alloc := range ap.Allocations {
		if alloc.CharitableOrgWalletId == 0 {
			return fmt.Errorf("charitable org wallet ID cannot be zero")
		}
		
		if alloc.Amount.IsNegative() {
			return fmt.Errorf("allocation amount cannot be negative")
		}
		
		totalAllocated = totalAllocated.Add(alloc.Amount)
	}
	
	if !totalAllocated.IsEqual(ap.TotalAmount) {
		return fmt.Errorf("sum of allocations doesn't match total amount")
	}
	
	return nil
}

// ParamKeyTable returns the parameter key table
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateEnabled),
		paramtypes.NewParamSetPair(KeyMinAllocationAmount, &p.MinAllocationAmount, validateMinAllocationAmount),
		paramtypes.NewParamSetPair(KeyMaxMonthlyAllocationPerOrg, &p.MaxMonthlyAllocationPerOrg, validateMaxMonthlyAllocationPerOrg),
		paramtypes.NewParamSetPair(KeyProposalVotingPeriod, &p.ProposalVotingPeriod, validateProposalVotingPeriod),
		paramtypes.NewParamSetPair(KeyFraudInvestigationPeriod, &p.FraudInvestigationPeriod, validateFraudInvestigationPeriod),
		paramtypes.NewParamSetPair(KeyImpactReportFrequency, &p.ImpactReportFrequency, validateImpactReportFrequency),
		paramtypes.NewParamSetPair(KeyAuthorizedInvestigators, &p.AuthorizedInvestigators, validateAuthorizedInvestigators),
		paramtypes.NewParamSetPair(KeyDistributionCategories, &p.DistributionCategories, validateDistributionCategories),
		paramtypes.NewParamSetPair(KeyEmergencyPauseAuthorities, &p.EmergencyPauseAuthorities, validateEmergencyPauseAuthorities),
	}
}

// Parameter store keys
var (
	KeyEnabled                    = []byte("Enabled")
	KeyMinAllocationAmount        = []byte("MinAllocationAmount")
	KeyMaxMonthlyAllocationPerOrg = []byte("MaxMonthlyAllocationPerOrg")
	KeyProposalVotingPeriod       = []byte("ProposalVotingPeriod")
	KeyFraudInvestigationPeriod   = []byte("FraudInvestigationPeriod")
	KeyImpactReportFrequency      = []byte("ImpactReportFrequency")
	KeyAuthorizedInvestigators    = []byte("AuthorizedInvestigators")
	KeyDistributionCategories     = []byte("DistributionCategories")
	KeyEmergencyPauseAuthorities  = []byte("EmergencyPauseAuthorities")
)

// Validation functions
func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinAllocationAmount(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("minimum allocation amount cannot be negative")
	}
	return nil
}

func validateMaxMonthlyAllocationPerOrg(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("max monthly allocation per org cannot be negative")
	}
	return nil
}

func validateProposalVotingPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("proposal voting period must be positive")
	}
	return nil
}

func validateFraudInvestigationPeriod(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("fraud investigation period must be positive")
	}
	return nil
}

func validateImpactReportFrequency(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("impact report frequency must be positive")
	}
	return nil
}

func validateAuthorizedInvestigators(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, addr := range v {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid investigator address: %w", err)
		}
	}
	return nil
}

func validateDistributionCategories(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("must have at least one distribution category")
	}
	return nil
}

func validateEmergencyPauseAuthorities(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, addr := range v {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid emergency pause authority address: %w", err)
		}
	}
	return nil
}