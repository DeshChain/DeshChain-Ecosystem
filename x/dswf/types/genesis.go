package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		FundGovernance: FundGovernance{
			MinProposalDeposit: sdk.NewCoin("unamo", sdk.NewInt(1000000000)), // 1000 NAMO
			VotingPeriod:       604800, // 7 days in seconds
			Quorum:             sdk.MustNewDecFromStr("0.667"), // 66.7%
			Threshold:          sdk.MustNewDecFromStr("0.6"), // 60%
			MaxAllocation:      sdk.NewCoin("unamo", sdk.NewInt(50000000000000)), // 50M NAMO
			FundManagers:       []string{}, // To be set on chain initialization
			EmergencySignaturesRequired: 3,
			Categories: []AllocationCategory{
				{
					Name:              "infrastructure",
					Description:       "Blockchain infrastructure and development",
					MaxPercentage:     sdk.MustNewDecFromStr("0.30"),
					CurrentAllocation: sdk.NewCoin("unamo", sdk.ZeroInt()),
					ActiveProjects:    0,
				},
				{
					Name:              "ecosystem",
					Description:       "Ecosystem growth and partnerships",
					MaxPercentage:     sdk.MustNewDecFromStr("0.25"),
					CurrentAllocation: sdk.NewCoin("unamo", sdk.ZeroInt()),
					ActiveProjects:    0,
				},
				{
					Name:              "innovation",
					Description:       "Research and new technology development",
					MaxPercentage:     sdk.MustNewDecFromStr("0.20"),
					CurrentAllocation: sdk.NewCoin("unamo", sdk.ZeroInt()),
					ActiveProjects:    0,
				},
				{
					Name:              "marketing",
					Description:       "Marketing and community building",
					MaxPercentage:     sdk.MustNewDecFromStr("0.15"),
					CurrentAllocation: sdk.NewCoin("unamo", sdk.ZeroInt()),
					ActiveProjects:    0,
				},
				{
					Name:              "emergency",
					Description:       "Emergency and crisis management",
					MaxPercentage:     sdk.MustNewDecFromStr("0.10"),
					CurrentAllocation: sdk.NewCoin("unamo", sdk.ZeroInt()),
					ActiveProjects:    0,
				},
			},
		},
		InvestmentPortfolio: InvestmentPortfolio{
			TotalValue:     sdk.NewCoin("unamo", sdk.ZeroInt()),
			LiquidAssets:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			InvestedAssets: sdk.NewCoin("unamo", sdk.ZeroInt()),
			ReservedAssets: sdk.NewCoin("unamo", sdk.ZeroInt()),
			Components:     []PortfolioComponent{},
			TotalReturns:   sdk.NewCoin("unamo", sdk.ZeroInt()),
			AnnualReturnRate: sdk.ZeroDec(),
			RiskScore:       5, // Medium risk default
		},
		Allocations:    []FundAllocation{},
		MonthlyReports: []MonthlyReport{},
	}
}

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{
		Enabled:                 true,
		MinimumFundBalance:      sdk.NewCoin("unamo", sdk.NewInt(1000000000000)), // 1M NAMO
		MaxAllocationPercentage: sdk.MustNewDecFromStr("0.05"), // 5% max per allocation
		InvestmentStrategy: InvestmentStrategy{
			ConservativePercentage: sdk.MustNewDecFromStr("0.30"),
			ModeratePercentage:     sdk.MustNewDecFromStr("0.40"),
			AggressivePercentage:   sdk.MustNewDecFromStr("0.20"),
			RebalancingFrequency:   90, // 90 days
			MinReturnTarget:        sdk.MustNewDecFromStr("0.08"), // 8% minimum
			MaxRiskScore:           7,
		},
		ReportingFrequency:   86400, // Daily blocks (assuming 6s blocks)
		AuthorizedAuditors:   []string{}, // To be set by governance
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	
	// Validate fund governance
	if err := gs.FundGovernance.Validate(); err != nil {
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
	
	// Validate monthly reports
	reportPeriods := make(map[string]bool)
	for _, report := range gs.MonthlyReports {
		if reportPeriods[report.Period] {
			return fmt.Errorf("duplicate report period: %s", report.Period)
		}
		reportPeriods[report.Period] = true
	}
	
	return nil
}

// Validate validates the module parameters
func (p Params) Validate() error {
	if p.MinimumFundBalance.IsNegative() {
		return fmt.Errorf("minimum fund balance cannot be negative")
	}
	
	if p.MaxAllocationPercentage.IsNegative() || p.MaxAllocationPercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("max allocation percentage must be between 0 and 1")
	}
	
	if err := p.InvestmentStrategy.Validate(); err != nil {
		return err
	}
	
	if p.ReportingFrequency <= 0 {
		return fmt.Errorf("reporting frequency must be positive")
	}
	
	return nil
}

// Validate validates the investment strategy
func (is InvestmentStrategy) Validate() error {
	// Check percentages sum to less than or equal to 1
	total := is.ConservativePercentage.Add(is.ModeratePercentage).Add(is.AggressivePercentage)
	if total.GT(sdk.OneDec()) {
		return fmt.Errorf("investment percentages cannot exceed 100%%")
	}
	
	if is.RebalancingFrequency <= 0 {
		return fmt.Errorf("rebalancing frequency must be positive")
	}
	
	if is.MinReturnTarget.IsNegative() {
		return fmt.Errorf("minimum return target cannot be negative")
	}
	
	if is.MaxRiskScore < 1 || is.MaxRiskScore > 10 {
		return fmt.Errorf("max risk score must be between 1 and 10")
	}
	
	return nil
}

// Validate validates the fund governance
func (fg FundGovernance) Validate() error {
	if fg.MinProposalDeposit.IsNegative() {
		return fmt.Errorf("minimum proposal deposit cannot be negative")
	}
	
	if fg.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive")
	}
	
	if fg.Quorum.IsNegative() || fg.Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("quorum must be between 0 and 1")
	}
	
	if fg.Threshold.IsNegative() || fg.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("threshold must be between 0 and 1")
	}
	
	if fg.MaxAllocation.IsNegative() {
		return fmt.Errorf("max allocation cannot be negative")
	}
	
	if fg.EmergencySignaturesRequired <= 0 || fg.EmergencySignaturesRequired > int32(len(fg.FundManagers)) {
		return fmt.Errorf("emergency signatures required must be between 1 and number of fund managers")
	}
	
	// Validate categories
	totalMaxPercentage := sdk.ZeroDec()
	categoryNames := make(map[string]bool)
	for _, cat := range fg.Categories {
		if categoryNames[cat.Name] {
			return fmt.Errorf("duplicate category name: %s", cat.Name)
		}
		categoryNames[cat.Name] = true
		
		if cat.MaxPercentage.IsNegative() || cat.MaxPercentage.GT(sdk.OneDec()) {
			return fmt.Errorf("category max percentage must be between 0 and 1")
		}
		
		totalMaxPercentage = totalMaxPercentage.Add(cat.MaxPercentage)
	}
	
	if totalMaxPercentage.LT(sdk.OneDec()) {
		return fmt.Errorf("total category max percentages must be at least 100%%")
	}
	
	return nil
}

// Validate validates a fund allocation
func (fa FundAllocation) Validate() error {
	if fa.Id == 0 {
		return fmt.Errorf("allocation ID cannot be zero")
	}
	
	if fa.Purpose == "" {
		return fmt.Errorf("allocation purpose cannot be empty")
	}
	
	if fa.Category == "" {
		return fmt.Errorf("allocation category cannot be empty")
	}
	
	if fa.Amount.IsNegative() {
		return fmt.Errorf("allocation amount cannot be negative")
	}
	
	if _, err := sdk.AccAddressFromBech32(fa.Recipient); err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}
	
	for _, approver := range fa.ApprovedBy {
		if _, err := sdk.AccAddressFromBech32(approver); err != nil {
			return fmt.Errorf("invalid approver address: %w", err)
		}
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
		paramtypes.NewParamSetPair(KeyMinimumFundBalance, &p.MinimumFundBalance, validateMinimumFundBalance),
		paramtypes.NewParamSetPair(KeyMaxAllocationPercentage, &p.MaxAllocationPercentage, validateMaxAllocationPercentage),
		paramtypes.NewParamSetPair(KeyInvestmentStrategy, &p.InvestmentStrategy, validateInvestmentStrategy),
		paramtypes.NewParamSetPair(KeyReportingFrequency, &p.ReportingFrequency, validateReportingFrequency),
		paramtypes.NewParamSetPair(KeyAuthorizedAuditors, &p.AuthorizedAuditors, validateAuthorizedAuditors),
	}
}

// Parameter store keys
var (
	KeyEnabled                 = []byte("Enabled")
	KeyMinimumFundBalance      = []byte("MinimumFundBalance")
	KeyMaxAllocationPercentage = []byte("MaxAllocationPercentage")
	KeyInvestmentStrategy      = []byte("InvestmentStrategy")
	KeyReportingFrequency      = []byte("ReportingFrequency")
	KeyAuthorizedAuditors      = []byte("AuthorizedAuditors")
)

// Validation functions
func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinimumFundBalance(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("minimum fund balance cannot be negative")
	}
	return nil
}

func validateMaxAllocationPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() || v.GT(sdk.OneDec()) {
		return fmt.Errorf("max allocation percentage must be between 0 and 1")
	}
	return nil
}

func validateInvestmentStrategy(i interface{}) error {
	v, ok := i.(InvestmentStrategy)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.Validate()
}

func validateReportingFrequency(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("reporting frequency must be positive")
	}
	return nil
}

func validateAuthorizedAuditors(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, addr := range v {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid auditor address: %w", err)
		}
	}
	return nil
}