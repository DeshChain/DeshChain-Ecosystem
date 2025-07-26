package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dswf/types"
)

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()
	
	require.NotNil(t, genesis)
	require.NotNil(t, genesis.Params)
	require.NotNil(t, genesis.FundGovernance)
	require.Empty(t, genesis.Allocations)
	require.Empty(t, genesis.MonthlyReports)
	require.Equal(t, uint64(0), genesis.AllocationCount)
}

func TestGenesisStateValidate(t *testing.T) {
	tests := []struct {
		name        string
		genesis     *types.GenesisState
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid genesis",
			genesis:     types.DefaultGenesis(),
			shouldError: false,
		},
		{
			name: "invalid params - negative min fund balance",
			genesis: &types.GenesisState{
				Params: types.Params{
					MinFundBalance:          sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
					MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
					MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
					RebalancingFrequency:   90,
					AllocationCategories:   []string{"infrastructure"},
				},
				FundGovernance: types.FundGovernance{
					FundManagers:       []types.FundManager{{Address: "desh1test..."}},
					RequiredSignatures: 1,
				},
			},
			shouldError: true,
			errorMsg:    "minimum fund balance cannot be negative",
		},
		{
			name: "invalid params - max allocation percentage > 1",
			genesis: &types.GenesisState{
				Params: types.Params{
					MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
					MaxAllocationPercentage: sdk.MustNewDecFromStr("1.5"),
					MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
					RebalancingFrequency:   90,
					AllocationCategories:   []string{"infrastructure"},
				},
				FundGovernance: types.FundGovernance{
					FundManagers:       []types.FundManager{{Address: "desh1test..."}},
					RequiredSignatures: 1,
				},
			},
			shouldError: true,
			errorMsg:    "max allocation percentage must be between 0 and 1",
		},
		{
			name: "invalid governance - no fund managers",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				FundGovernance: types.FundGovernance{
					FundManagers:       []types.FundManager{},
					RequiredSignatures: 1,
				},
			},
			shouldError: true,
			errorMsg:    "must have at least one fund manager",
		},
		{
			name: "invalid governance - required signatures > managers",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				FundGovernance: types.FundGovernance{
					FundManagers: []types.FundManager{
						{Address: "desh1test1..."},
						{Address: "desh1test2..."},
					},
					RequiredSignatures: 3,
				},
			},
			shouldError: true,
			errorMsg:    "required signatures cannot exceed number of fund managers",
		},
		{
			name: "duplicate allocation IDs",
			genesis: &types.GenesisState{
				Params:         types.DefaultParams(),
				FundGovernance: types.DefaultFundGovernance(),
				Allocations: []types.FundAllocation{
					{
						Id:       1,
						Purpose:  "Test 1",
						Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000)),
						Category: "infrastructure",
						Status:   "active",
					},
					{
						Id:       1, // Duplicate ID
						Purpose:  "Test 2",
						Amount:   sdk.NewCoin("unamo", sdk.NewInt(2000000)),
						Category: "education",
						Status:   "active",
					},
				},
			},
			shouldError: true,
			errorMsg:    "duplicate allocation ID: 1",
		},
		{
			name: "invalid allocation - negative amount",
			genesis: &types.GenesisState{
				Params:         types.DefaultParams(),
				FundGovernance: types.DefaultFundGovernance(),
				Allocations: []types.FundAllocation{
					{
						Id:       1,
						Purpose:  "Test",
						Amount:   sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
						Category: "infrastructure",
						Status:   "active",
					},
				},
			},
			shouldError: true,
			errorMsg:    "allocation amount cannot be negative",
		},
		{
			name: "duplicate monthly report periods",
			genesis: &types.GenesisState{
				Params:         types.DefaultParams(),
				FundGovernance: types.DefaultFundGovernance(),
				MonthlyReports: []types.MonthlyReport{
					{
						Period:         "2025-01",
						OpeningBalance: sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
						ClosingBalance: sdk.NewCoin("unamo", sdk.NewInt(1100000000)),
					},
					{
						Period:         "2025-01", // Duplicate period
						OpeningBalance: sdk.NewCoin("unamo", sdk.NewInt(1100000000)),
						ClosingBalance: sdk.NewCoin("unamo", sdk.NewInt(1200000000)),
					},
				},
			},
			shouldError: true,
			errorMsg:    "duplicate report period: 2025-01",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParamsValidate(t *testing.T) {
	tests := []struct {
		name        string
		params      types.Params
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid params",
			params:      types.DefaultParams(),
			shouldError: false,
		},
		{
			name: "negative min fund balance",
			params: types.Params{
				MinFundBalance:          sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{"infrastructure"},
			},
			shouldError: true,
			errorMsg:    "minimum fund balance cannot be negative",
		},
		{
			name: "max allocation percentage > 1",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("1.1"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{"infrastructure"},
			},
			shouldError: true,
			errorMsg:    "max allocation percentage must be between 0 and 1",
		},
		{
			name: "min liquidity ratio negative",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("-0.1"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{"infrastructure"},
			},
			shouldError: true,
			errorMsg:    "min liquidity ratio must be between 0 and 1",
		},
		{
			name: "zero rebalancing frequency",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   0,
				AllocationCategories:   []string{"infrastructure"},
			},
			shouldError: true,
			errorMsg:    "rebalancing frequency must be positive",
		},
		{
			name: "no allocation categories",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{},
			},
			shouldError: true,
			errorMsg:    "must have at least one allocation category",
		},
		{
			name: "duplicate allocation categories",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{"infrastructure", "education", "infrastructure"},
			},
			shouldError: true,
			errorMsg:    "duplicate allocation category",
		},
		{
			name: "invalid risk score",
			params: types.Params{
				MinFundBalance:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
				MinLiquidityRatio:      sdk.MustNewDecFromStr("0.20"),
				RebalancingFrequency:   90,
				AllocationCategories:   []string{"infrastructure"},
				MaxRiskScore:           11,
			},
			shouldError: true,
			errorMsg:    "max risk score must be between 1 and 10",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFundGovernanceValidate(t *testing.T) {
	tests := []struct {
		name        string
		governance  types.FundGovernance
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid governance",
			governance: types.FundGovernance{
				FundManagers: []types.FundManager{
					{Address: "desh1manager1...", Name: "Manager 1"},
					{Address: "desh1manager2...", Name: "Manager 2"},
				},
				RequiredSignatures: 2,
				ApprovalThreshold:  sdk.MustNewDecFromStr("0.67"),
			},
			shouldError: false,
		},
		{
			name: "no fund managers",
			governance: types.FundGovernance{
				FundManagers:       []types.FundManager{},
				RequiredSignatures: 1,
			},
			shouldError: true,
			errorMsg:    "must have at least one fund manager",
		},
		{
			name: "duplicate manager addresses",
			governance: types.FundGovernance{
				FundManagers: []types.FundManager{
					{Address: "desh1manager1...", Name: "Manager 1"},
					{Address: "desh1manager1...", Name: "Manager 2"}, // Duplicate address
				},
				RequiredSignatures: 1,
			},
			shouldError: true,
			errorMsg:    "duplicate fund manager address",
		},
		{
			name: "empty manager name",
			governance: types.FundGovernance{
				FundManagers: []types.FundManager{
					{Address: "desh1manager1...", Name: ""},
				},
				RequiredSignatures: 1,
			},
			shouldError: true,
			errorMsg:    "fund manager name cannot be empty",
		},
		{
			name: "required signatures > managers",
			governance: types.FundGovernance{
				FundManagers: []types.FundManager{
					{Address: "desh1manager1...", Name: "Manager 1"},
				},
				RequiredSignatures: 2,
			},
			shouldError: true,
			errorMsg:    "required signatures cannot exceed number of fund managers",
		},
		{
			name: "invalid approval threshold",
			governance: types.FundGovernance{
				FundManagers: []types.FundManager{
					{Address: "desh1manager1...", Name: "Manager 1"},
				},
				RequiredSignatures: 1,
				ApprovalThreshold:  sdk.MustNewDecFromStr("1.5"),
			},
			shouldError: true,
			errorMsg:    "approval threshold must be between 0 and 1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.governance.Validate()
			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFundAllocationValidate(t *testing.T) {
	tests := []struct {
		name        string
		allocation  types.FundAllocation
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid allocation",
			allocation: types.FundAllocation{
				Id:       1,
				Purpose:  "Infrastructure Development",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				Category: "infrastructure",
				Status:   "active",
			},
			shouldError: false,
		},
		{
			name: "zero ID",
			allocation: types.FundAllocation{
				Id:       0,
				Purpose:  "Test",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category: "infrastructure",
				Status:   "active",
			},
			shouldError: true,
			errorMsg:    "allocation ID cannot be zero",
		},
		{
			name: "empty purpose",
			allocation: types.FundAllocation{
				Id:       1,
				Purpose:  "",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category: "infrastructure",
				Status:   "active",
			},
			shouldError: true,
			errorMsg:    "allocation purpose cannot be empty",
		},
		{
			name: "negative amount",
			allocation: types.FundAllocation{
				Id:       1,
				Purpose:  "Test",
				Amount:   sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
				Category: "infrastructure",
				Status:   "active",
			},
			shouldError: true,
			errorMsg:    "allocation amount cannot be negative",
		},
		{
			name: "empty category",
			allocation: types.FundAllocation{
				Id:       1,
				Purpose:  "Test",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category: "",
				Status:   "active",
			},
			shouldError: true,
			errorMsg:    "allocation category cannot be empty",
		},
		{
			name: "invalid status",
			allocation: types.FundAllocation{
				Id:       1,
				Purpose:  "Test",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category: "infrastructure",
				Status:   "invalid_status",
			},
			shouldError: true,
			errorMsg:    "invalid allocation status",
		},
		{
			name: "invalid expected returns",
			allocation: types.FundAllocation{
				Id:              1,
				Purpose:         "Test",
				Amount:          sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category:        "infrastructure",
				Status:          "active",
				ExpectedReturns: sdk.MustNewDecFromStr("-0.5"),
			},
			shouldError: true,
			errorMsg:    "expected returns cannot be negative",
		},
		{
			name: "invalid risk category",
			allocation: types.FundAllocation{
				Id:           1,
				Purpose:      "Test",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(1000000)),
				Category:     "infrastructure",
				Status:       "active",
				RiskCategory: "invalid_risk",
			},
			shouldError: true,
			errorMsg:    "invalid risk category",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.allocation.Validate()
			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestInvestmentPortfolioValidate(t *testing.T) {
	tests := []struct {
		name        string
		portfolio   types.InvestmentPortfolio
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid portfolio",
			portfolio: types.InvestmentPortfolio{
				TotalValue:      sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				LiquidAssets:    sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				InvestedAssets:  sdk.NewCoin("unamo", sdk.NewInt(6000000000)),
				ReservedAssets:  sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				Components:      []types.PortfolioComponent{},
				TotalReturns:    sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				AnnualReturnRate: sdk.MustNewDecFromStr("0.08"),
				RiskScore:       3,
				LastRebalanced:  time.Now(),
			},
			shouldError: false,
		},
		{
			name: "component sum mismatch",
			portfolio: types.InvestmentPortfolio{
				TotalValue:     sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				LiquidAssets:   sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				InvestedAssets: sdk.NewCoin("unamo", sdk.NewInt(6000000000)),
				ReservedAssets: sdk.NewCoin("unamo", sdk.NewInt(2000000000)), // Sum > TotalValue
			},
			shouldError: true,
			errorMsg:    "asset components sum exceeds total value",
		},
		{
			name: "negative returns",
			portfolio: types.InvestmentPortfolio{
				TotalValue:     sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				LiquidAssets:   sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				TotalReturns:   sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-100000)},
			},
			shouldError: true,
			errorMsg:    "total returns cannot be negative",
		},
		{
			name: "invalid risk score",
			portfolio: types.InvestmentPortfolio{
				TotalValue:   sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				LiquidAssets: sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
				RiskScore:    11,
			},
			shouldError: true,
			errorMsg:    "risk score must be between 1 and 10",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.portfolio.Validate()
			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}