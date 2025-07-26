package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/app"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dswf/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dswf/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *app.App
	ctx    sdk.Context
	keeper keeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height: 1,
		Time:   time.Now(),
	})
	suite.keeper = suite.app.DSWFKeeper
}

// Test Fund Balance Management
func (suite *KeeperTestSuite) TestFundBalance() {
	// Initial balance should be zero
	balance := suite.keeper.GetFundBalance(suite.ctx)
	suite.Require().Equal(0, len(balance))

	// Add funds to the DSWF
	newBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(1000000000))) // 1000 NAMO
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, newBalance)
	suite.Require().NoError(err)

	// Check balance
	balance = suite.keeper.GetFundBalance(suite.ctx)
	suite.Require().Equal(newBalance, balance)
}

// Test Fund Allocation Creation
func (suite *KeeperTestSuite) TestCreateFundAllocation() {
	// Setup fund balance
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000))) // 10000 NAMO
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Create allocation
	allocation := types.FundAllocation{
		Id:              1,
		Purpose:         "Education Infrastructure",
		Amount:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)), // 1000 NAMO
		Category:        "infrastructure",
		Recipient:       "desh1education123...",
		ProposedBy:      []string{"desh1manager1...", "desh1manager2..."},
		ApprovedBy:      []string{"desh1manager1...", "desh1manager2...", "desh1manager3..."},
		ProposedAt:      suite.ctx.BlockTime(),
		ApprovedAt:      suite.ctx.BlockTime(),
		Status:          "approved",
		ExpectedReturns: sdk.MustNewDecFromStr("1.08"), // 8% expected return
		RiskCategory:    "low",
	}

	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	// Retrieve and verify
	retrieved, found := suite.keeper.GetFundAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(allocation, retrieved)
}

// Test Investment Portfolio Management
func (suite *KeeperTestSuite) TestInvestmentPortfolio() {
	// Initial portfolio should not exist
	_, found := suite.keeper.GetInvestmentPortfolio(suite.ctx)
	suite.Require().False(found)

	// Create portfolio
	portfolio := types.InvestmentPortfolio{
		TotalValue:      sdk.NewCoin("unamo", sdk.NewInt(10000000000)), // 10000 NAMO
		LiquidAssets:    sdk.NewCoin("unamo", sdk.NewInt(3000000000)),  // 3000 NAMO
		InvestedAssets:  sdk.NewCoin("unamo", sdk.NewInt(6000000000)),  // 6000 NAMO
		ReservedAssets:  sdk.NewCoin("unamo", sdk.NewInt(1000000000)),  // 1000 NAMO
		AllocatedAmount: sdk.NewCoin("unamo", sdk.NewInt(500000000)),   // 500 NAMO
		AvailableAmount: sdk.NewCoin("unamo", sdk.NewInt(2500000000)),  // 2500 NAMO
		InvestedAmount:  sdk.NewCoin("unamo", sdk.NewInt(6000000000)),  // 6000 NAMO
		Components: []types.PortfolioComponent{
			{
				AssetType:    "treasury_bonds",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(3240000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.08"),
				RiskRating:   "AAA",
				Maturity:     suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			},
			{
				AssetType:    "equity_index",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(2000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(2300000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.15"),
				RiskRating:   "A",
			},
			{
				AssetType:    "startup_fund",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(1100000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.10"),
				RiskRating:   "B",
			},
		},
		TotalReturns:     sdk.NewCoin("unamo", sdk.NewInt(640000000)),
		AnnualReturnRate: sdk.MustNewDecFromStr("0.107"), // 10.7% weighted average
		RiskScore:        3,                               // Moderate risk
		LastRebalanced:   suite.ctx.BlockTime(),
	}

	suite.keeper.SetInvestmentPortfolio(suite.ctx, portfolio)

	// Retrieve and verify
	retrieved, found := suite.keeper.GetInvestmentPortfolio(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(portfolio, retrieved)
}

// Test Fund Governance
func (suite *KeeperTestSuite) TestFundGovernance() {
	// Initial governance should not exist
	_, found := suite.keeper.GetFundGovernance(suite.ctx)
	suite.Require().False(found)

	// Create governance
	governance := types.FundGovernance{
		FundManagers: []types.FundManager{
			{
				Address:     "desh1manager1...",
				Name:        "Investment Manager 1",
				Expertise:   "Fixed Income",
				AddedAt:     suite.ctx.BlockTime(),
				Performance: sdk.MustNewDecFromStr("1.08"),
			},
			{
				Address:     "desh1manager2...",
				Name:        "Investment Manager 2",
				Expertise:   "Equities",
				AddedAt:     suite.ctx.BlockTime(),
				Performance: sdk.MustNewDecFromStr("1.15"),
			},
		},
		RequiredSignatures: 2,
		ApprovalThreshold:  sdk.MustNewDecFromStr("0.67"),
		InvestmentCommittee: []string{
			"desh1committee1...",
			"desh1committee2...",
			"desh1committee3...",
		},
		RiskOfficer:       "desh1risk...",
		ComplianceOfficer: "desh1compliance...",
		AuditSchedule:     180, // Every 6 months
		LastAudit:         suite.ctx.BlockTime(),
		NextReview:        suite.ctx.BlockTime().Add(90 * 24 * time.Hour),
	}

	suite.keeper.SetFundGovernance(suite.ctx, governance)

	// Retrieve and verify
	retrieved, found := suite.keeper.GetFundGovernance(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(governance, retrieved)
}

// Test Allocation Proposal Validation
func (suite *KeeperTestSuite) TestValidateAllocationProposal() {
	// Setup fund balance
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000))) // 10000 NAMO
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Test valid allocation
	validAmount := sdk.NewCoin("unamo", sdk.NewInt(500000000)) // 500 NAMO (5% of fund)
	err = suite.keeper.ValidateAllocationProposal(suite.ctx, validAmount, "infrastructure")
	suite.Require().NoError(err)

	// Test allocation exceeding limit
	excessAmount := sdk.NewCoin("unamo", sdk.NewInt(1100000000)) // 1100 NAMO (11% of fund)
	err = suite.keeper.ValidateAllocationProposal(suite.ctx, excessAmount, "infrastructure")
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "exceeds maximum allocation limit")

	// Test invalid category
	err = suite.keeper.ValidateAllocationProposal(suite.ctx, validAmount, "invalid_category")
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid allocation category")

	// Test with low fund balance
	// Burn most funds
	burnAmount := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(9500000000)))
	err = suite.app.BankKeeper.BurnCoins(suite.ctx, types.ModuleName, burnAmount)
	suite.Require().NoError(err)

	err = suite.keeper.ValidateAllocationProposal(suite.ctx, validAmount, "infrastructure")
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "fund balance below minimum threshold")
}

// Test Monthly Report Generation
func (suite *KeeperTestSuite) TestGenerateMonthlyReport() {
	// Setup initial data
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Create some allocations
	allocation1 := types.FundAllocation{
		Id:       1,
		Purpose:  "Education Infrastructure",
		Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Category: "infrastructure",
		Status:   "completed",
		Returns: &types.InvestmentReturns{
			ExpectedAmount: sdk.NewCoin("unamo", sdk.NewInt(1080000000)),
			ActualAmount:   sdk.NewCoin("unamo", sdk.NewInt(1100000000)),
			ReturnRate:     sdk.MustNewDecFromStr("0.10"),
			Period:         30,
		},
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation1)

	allocation2 := types.FundAllocation{
		Id:       2,
		Purpose:  "Healthcare Equipment",
		Amount:   sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		Category: "healthcare",
		Status:   "active",
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation2)

	// Generate report
	report := suite.keeper.GenerateMonthlyReport(suite.ctx, "2025-01")
	
	suite.Require().Equal("2025-01", report.Period)
	suite.Require().Equal(fundBalance[0], report.OpeningBalance)
	suite.Require().Equal(int32(1), report.AllocationsApproved)
	suite.Require().Equal(int32(1), report.AllocationsCompleted)
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(1500000000)), report.TotalAllocated)
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(100000000)), report.TotalReturns)
	suite.Require().True(report.AverageReturnRate.GT(sdk.ZeroDec()))
}

// Test Multi-signature Validation
func (suite *KeeperTestSuite) TestValidateMultiSignature() {
	// Setup governance
	governance := types.FundGovernance{
		FundManagers: []types.FundManager{
			{Address: "desh1manager1..."},
			{Address: "desh1manager2..."},
			{Address: "desh1manager3..."},
		},
		RequiredSignatures: 2,
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	// Test valid signatures
	signers := []string{"desh1manager1...", "desh1manager2..."}
	valid := suite.keeper.ValidateMultiSignature(suite.ctx, signers)
	suite.Require().True(valid)

	// Test insufficient signatures
	signers = []string{"desh1manager1..."}
	valid = suite.keeper.ValidateMultiSignature(suite.ctx, signers)
	suite.Require().False(valid)

	// Test invalid signer
	signers = []string{"desh1manager1...", "desh1invalid..."}
	valid = suite.keeper.ValidateMultiSignature(suite.ctx, signers)
	suite.Require().False(valid)

	// Test duplicate signers
	signers = []string{"desh1manager1...", "desh1manager1..."}
	valid = suite.keeper.ValidateMultiSignature(suite.ctx, signers)
	suite.Require().False(valid)
}

// Test Disbursement Execution
func (suite *KeeperTestSuite) TestExecuteDisbursement() {
	// Setup fund balance
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Create allocation with disbursements
	allocation := types.FundAllocation{
		Id:         1,
		Purpose:    "Education Infrastructure",
		Amount:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Category:   "infrastructure",
		Recipient:  "desh1education123...",
		Status:     "active",
		Disbursements: []types.Disbursement{
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				ScheduledDate: suite.ctx.BlockTime(),
				Status:        "pending",
				Milestone:     "Phase 1 - Foundation",
			},
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				ScheduledDate: suite.ctx.BlockTime().Add(30 * 24 * time.Hour),
				Status:        "pending",
				Milestone:     "Phase 2 - Construction",
			},
		},
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	// Execute first disbursement
	err = suite.keeper.ExecuteDisbursement(suite.ctx, 1, 0)
	suite.Require().NoError(err)

	// Verify disbursement status
	updatedAllocation, found := suite.keeper.GetFundAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("completed", updatedAllocation.Disbursements[0].Status)
	suite.Require().Equal("pending", updatedAllocation.Disbursements[1].Status)

	// Try to execute already completed disbursement
	err = suite.keeper.ExecuteDisbursement(suite.ctx, 1, 0)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already executed")
}

// Test Investment Strategy Validation
func (suite *KeeperTestSuite) TestValidateInvestmentStrategy() {
	// Valid conservative strategy
	conservativeComponents := []types.PortfolioComponent{
		{
			AssetType:  "treasury_bonds",
			Amount:     sdk.NewCoin("unamo", sdk.NewInt(7000000000)), // 70%
			RiskRating: "AAA",
		},
		{
			AssetType:  "corporate_bonds",
			Amount:     sdk.NewCoin("unamo", sdk.NewInt(2000000000)), // 20%
			RiskRating: "AA",
		},
		{
			AssetType:  "cash_equivalent",
			Amount:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)), // 10%
			RiskRating: "AAA",
		},
	}

	valid := suite.keeper.ValidateInvestmentStrategy(suite.ctx, conservativeComponents, "conservative")
	suite.Require().True(valid)

	// Invalid growth strategy (too much high-risk)
	invalidGrowthComponents := []types.PortfolioComponent{
		{
			AssetType:  "startup_fund",
			Amount:     sdk.NewCoin("unamo", sdk.NewInt(8000000000)), // 80% - too high
			RiskRating: "C",
		},
		{
			AssetType:  "equity_index",
			Amount:     sdk.NewCoin("unamo", sdk.NewInt(2000000000)), // 20%
			RiskRating: "B",
		},
	}

	valid = suite.keeper.ValidateInvestmentStrategy(suite.ctx, invalidGrowthComponents, "growth")
	suite.Require().False(valid)
}

// Test Fund Allocation Indexing
func (suite *KeeperTestSuite) TestAllocationIndexing() {
	// Create allocations with different statuses and categories
	allocations := []types.FundAllocation{
		{
			Id:       1,
			Category: "infrastructure",
			Status:   "active",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		},
		{
			Id:       2,
			Category: "healthcare",
			Status:   "active",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		},
		{
			Id:       3,
			Category: "infrastructure",
			Status:   "completed",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(800000000)),
		},
	}

	for _, alloc := range allocations {
		suite.keeper.SetFundAllocation(suite.ctx, alloc)
	}

	// Test GetAllocationsByStatus
	activeAllocs := suite.keeper.GetAllocationsByStatus(suite.ctx, "active")
	suite.Require().Len(activeAllocs, 2)

	completedAllocs := suite.keeper.GetAllocationsByStatus(suite.ctx, "completed")
	suite.Require().Len(completedAllocs, 1)

	// Test GetAllocationsByCategory
	infraAllocs := suite.keeper.GetAllocationsByCategory(suite.ctx, "infrastructure")
	suite.Require().Len(infraAllocs, 2)

	healthAllocs := suite.keeper.GetAllocationsByCategory(suite.ctx, "healthcare")
	suite.Require().Len(healthAllocs, 1)
}